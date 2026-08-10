# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## How We Use Claude Code

### 需求阶段
- 对不明确定义或不确定的问题，用 AskUserQuestion 显式列出，不替用户做假设
- 需要理解多处代码才能回答的问题，先派 Agent 去探索再汇报结论，不要一次性读大量文件

### 设计阶段
- 复杂变更先进入 Plan mode，输出实现方案后 ExitPlanMode 征求用户确认
- 多方案并存时，用 AskUserQuestion 的 preview 字段让用户并排比较；先定方向，再写正式代码
- 新增页面/功能需要同时配好：后端路由(main.go) + 业务逻辑(biz/xxx.go) + 管理面板(tables/xxx.go) + i18n翻译(mgmt/i18n/*.json)

### 实现阶段
- ❗ 人工改过的代码不要随便改回去，人工删除的代码不要随意加回。修改前先确认文件状态，区分"用户有意为之"和"遗漏"
- `define.go`与`xxx.go`成对出现：先确保`define.go`中类型完整，再写实现逻辑
- 涉及数据库查询的统一用`models.Orm.Table("table_name")`，不引入额外 ORM 层
- 用户可见的字符串都用`biz.T("key")`包裹，翻译 key 按层级用点号命名（如 `"error.api_id_not_found"`）
- API 响应统一格式：`{ code: 200|400, msg: biz.T("..."), data: {...} }`
- 异步任务（AI生成、报告生成等）用 `go func() { ... }()` ，耗时操作不要阻塞 HTTP handler
- 新增 biz/ 下的功能模块，文件名和包名保持一致

### 验证阶段
- 运行是验证：改完代码后执行 `make build` 确保编译通过，或用 `go run .` 启动验证
- 端到端验证优先于单测 — 这是一个运行时平台，行为正确性需要实际运行确认

### 变更日志
- 写变更日志时，先通过 `git log --since/--until` 查看指定日期范围的提交记录，再按日期分组整理
- 日志格式：`#### 2026年X月X日` 日期标题 + 编号列表，每项以 `[Feature]`/`[Bug]`/`[Optimize]` 分类
- 写入文件 `mgmt/doc/file/update/change_log.md`，按日期倒序插入新条目

## Project Overview

**Data4Test (盾测)** — An automated testing platform for complex systems written in Go. Built on GoAdmin (admin panel framework) + Gin (HTTP router) + GORM (ORM). Supports functional, concurrency, fuzz, scenario, performance, internationalization, and AI-driven testing.

## Architecture

```
main.go                   — Entry point: flag parsing, Gin routes + GoAdmin panel setup + server start
├── biz/                  — Business logic (~70 files), layered by domain
│   ├── xxx_define.go     — Data structures/types for each domain
│   └── xxx.go            — Implementation
├── tables/               — GoAdmin table/CRUD generators (~30 files)
│   └── tables.go         — Generators map (table prefix → generator function)
├── models/base.go        — GORM Orm singleton (*gorm.DB)
├── pages/                — GoAdmin dashboard page content generators
├── web/                  — Frontend assets (static/ts/vue)
├── html/                 — HTML template files
└── mgmt/                 — Documentation, SQL migrations, uploads, i18n translations
```

### Key biz/ Modules

| Area | Files | Responsibility |
|------|-------|---------------|
| API Mgmt | `api_manage.go`, `api_check.go`, `swagger.go` | Interface definitions, Swagger import, change detection |
| Test Data | `data.go`, `data_run.go`, `data_test_history.go` | Structured test data files (YAML), execution, history |
| Playbook | `biz_playbook.go`, `playbook_history.go` | Scenario orchestration: sequential/concurrent, compare modes |
| Task | `biz_task.go`, `task_report.go`, `cicd.go` | Scheduled tasks, task report generation, CI/CD hook |
| Assertion | `assert.go` | Status code / body structure / dynamic data assert engine |
| AI | `ai_*.go` (case/data/playbook/issue/template) | LLM-driven test generation & analysis |
| Fuzzing | `fuzzing.go` | Fuzz data generation for robustness testing |
| Mock | `mock.go` | Dynamic mock service for third-party data |
| Report | `report.go`, `dashboard.go` | Chart.js dashboards and statistics |

### Key Patterns

- **Routing**: All routes defined in `main.go` as `r.GET()`/`r.POST()` (JSON API) or `eng.HTMLFile()` (admin pages). GoAdmin auto-generates CRUD routes for each table generator in `tables.Generators`.
- **Business Logic**: Every domain has `<name>.go` + `<name>_define.go` — types in define, logic in the other.
- **Database**: `models.Orm.Table("table_name").Where(...).Find(&dest)` — raw GORM, no repository layer.
- **i18n**: `biz.T("key.with.dots")` reads from `mgmt/i18n/{locale}.json`. Key uses flattened dot-notation. Always wrap user-facing strings.
- **Async Tasks**: `go func() { ... }()` in route handlers. AI generation, report generation run in background.
- **API Response**: Uniform `{ code: 200|400, msg: "...", data: {...} }`.
- **Logging**: `biz.Logger.Info/Warning/Error/Debug(format, args...)`.
- **Panic Recovery**: `main()` defers `recover()` — auto-calls `startServer()` on panic.
- **History Files**: Test execution results saved as timestamped files under `HistoryBasePath`.

## Build & Development Commands

```bash
make serve                   # go run . (uses config.json, default port 9066)
make build                   # cross-compile → deploy/ (linux/darwin/windows)
CGO_ENABLED=0 go build .     # single binary (macOS)
make test                    # testpre → go test -v ./... → testafter
go test -v ./biz/...         # single package tests
make install                 # go mod tidy
make generate                # generate GoAdmin table code (from adm.ini)
docker-compose up -d         # quick local start (MySQL + app)
docker build -f appDockerfile -t data4test:5.0 .  # production image
```
