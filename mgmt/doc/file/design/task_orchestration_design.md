# 任务编排管理设计

> 关联需求文档：[task_orchestration_requirement.md](./task_orchestration_requirement.md)

## 1. 背景与目标

当前平台已有"定时任务"（`schedule`）能力：一个任务可以关联多个数据文件或场景，按 cron 定时执行，并生成任务级报告。但一个完整的回归/版本验证往往需要**多个任务按顺序组合执行**，例如：

- 登录冒烟任务 → 订单主流程任务 → 报表回归任务。

现状只能逐个任务手动运行，缺少"把任务组装成编排、一次触发、统一查看结果"的能力。

本方案目标：

- 新增**任务编排管理页面**，把已有定时任务按顺序组装成编排；
- 编排支持**页面「立即执行」**，按任务顺序串行执行；
- 普通失败**继续执行**后续步骤；`scene` / `data` 步骤返回 **401 鉴权失败时停止**后续步骤；
- 每次执行生成**编排级汇总报告**，并入「结果报告列表」，可下钻到各任务原有报告；
- 支持任务编排及其完整依赖的**导入/导出**；
- **只新增 `task_orchestration` 一张表**，执行记录复用现有 `dashboard` 与历史文件，不引入新的执行记录表；
- 不改变任务本身的执行逻辑与既有报告体系。

## 2. 已确认口径

| 项目 | 结论 |
|---|---|
| 编排组成 | 基于已有定时任务（`schedule`）组装，任务可被多个编排复用 |
| 执行方式 | **仅串行**；并行执行不纳入本期，模型不保留执行方式字段 |
| 触发方式 | **仅页面「立即执行」**；不做运行按钮、不做 cron 调度、不做 CLI/CI |
| 失败策略 | **普通失败继续**；`scene` / `data` 步骤返回 401 鉴权失败时停止后续步骤。不保留失败策略字段 |
| 编排存储 | `task_orchestration.task_list`（longtext），**每行一个任务名称，行序即步骤顺序** |
| 数据模型 | **只新增 `task_orchestration` 一张表**；不建步骤表、不建执行记录表 |
| 执行记录 | 复用定时任务的既有方式：任务级 `dashboard` + 编排级 `dashboard`（`report_type=orchestration`）+ `HistoryBasePath` 历史文件 |
| 编排与报告关联 | 不改 `dashboard` 表结构，用 `related_task_ids = orch:<orchestrationId>:<runTag>` 约定过滤 |
| 产品配置 | **编排级覆盖任务级**：编排配置了产品则所有步骤按编排产品执行；未配置时沿用各任务自身产品。`data` 类型步骤只使用第一个产品 |
| 鉴权刷新 | **本期不需要**，编排执行前不做统一鉴权刷新 |
| 编排报告 | 并入「结果报告列表」（`report_type=orchestration`），可下钻任务报告 |
| 页面形态 | 新增独立管理页，列表能力对齐定时任务（查询、复制、报告、导入导出等） |
| 任务改动 | 不改任务定义；仅抽取任务执行函数供编排复用 |

### 2.1 不在本期范围

- **并行执行**（并发上限、排队、并发失败策略全部不做）；
- **cron / 每天 / 每周定时触发**，以及服务重启后的调度恢复；
- **「运行」与「停止」按钮**（含批量停止）：本期只有「立即执行」，执行开始后不提供中断；
- **CLI / CI 接口**（`/cli/orchestrationRun` 等）；
- **编排执行前的统一产品鉴权刷新**；
- 编排中直接引用数据文件/场景（仍通过任务间接引用）；
- 步骤级产品覆盖、步骤级重试、步骤级超时；
- 邮件/IM 通知、审批流。

## 3. 概念模型

```text
任务编排 task_orchestration
   │  task_list（换行分隔任务名称，行序即步骤顺序）
   ├── 1:N ── 定时任务 schedule（按名称解析为 task_id）
   │
   └── 每次「立即执行」
          ├── 1 条编排报告 dashboard(report_type=orchestration, related_task_ids=orch:<id>:<runTag>)
          └── N 条任务报告 dashboard(report_type=task)，由每个步骤的任务执行产生
```

- **编排**：名称 + 任务编排 + 编排级产品 + 备注；不复制任务配置，任务修改后编排下次执行即生效。
- **步骤**：`task_list` 的一行，运行时按任务名称解析为 `schedule.id`；不单独建表。
- **执行记录**：不建执行记录表，通过编排报告的 `dashboard` 记录 + 报告内的步骤明细 JSON 表达；步骤明细指向任务报告 `dashboard.id` 用于下钻。

## 4. 产品设计

### 4.1 导航与页面

- 菜单位置：**任务 → 任务编排**（与"定时任务"同级，父菜单 id=203）。
- 路由：`/admin/info/task_orchestration`（GoAdmin 自动生成列表/新建/编辑/详情）。
- 编排报告：并入「结果报告列表」，统一使用 `/admin/report_dashboard?id=<dashboard_id>` 查看；`/admin/orchestration_dashboard?id=<orchestration_id>` 作为兼容入口，跳转该编排最近一次报告。

### 4.2 列表页

字段：

| 列 | 说明 |
|---|---|
| ID | 隐藏，可筛选 |
| 编排名称 | 模糊筛选；支持多语种展示（沿用任务名称的 i18n 机制） |
| 步骤数 | `task_list` 非空行数 |
| 关联产品 | 编排级产品 `task_orchestration.product_list`；未配置时展示步骤任务的产品并集（只展示） |
| 状态 | 未执行 / 执行中 / 已结束 |
| 上次执行时间 | 可排序 |
| 创建人 | 模糊筛选 |
| 创建时间 | 时间范围筛选 |
| 备注 | 隐藏，详情展示 |

筛选：编排名称、状态、创建人、创建时间。

行操作：

- **立即执行**：立即按串行方式执行一轮，不涉及任何调度状态。
- **复制**：复制编排字段（含 `task_list`），名称追加 `_副本`。
- **报告**：跳转 `/admin/orchestration_dashboard?id={{.Id}}`。
- **导出**：导出当前编排及完整依赖。
- **编辑 / 删除**：删除为软删除；执行中的编排禁止删除。

批量操作：**批量立即执行、批量复制、批量导出**（不含批量停止）。

导入/导出入口：

- 行操作「导出」：导出当前编排及完整依赖；
- 批量操作「批量导出」：把选中编排及去重后的共享依赖打入一个归档包；
- 顶部按钮「导入」：进入 `/admin/task_orchestration_import`，执行上传检查、冲突选择和确认导入。

**导出包结构**（实际实现，`orchestration_<yyyymmddHHMMSS>.zip`）：

```text
orchestration_<ts>.zip
├── orchestration_info.yaml        # 编排清单：name / tasks[] / products / remark
├── products_info.yaml             # 编排及任务使用的全部产品与应用环境配置（新增，可选，兼容旧包）
└── bundle/
    ├── <数据文件包>.tgz            # 由现有 ExportSchedule 生成，含 task_info.yaml、playbook_info.yaml 与数据文件
    ├── <附件包>.tgz                # 有上传附件时存在
    └── <SQL 导出>.sql              # 仅作参考，导入流程不执行
```

- 导出复用任务导出器生成任务/场景/数据/断言/系统参数/附件依赖，再与清单合成一个 zip；`products_info.yaml` 单独保存全部编排覆盖产品及任务产品的配置，避免原单产品 CLI 配置丢失多产品；`task_info.yaml` 新增可选 `product_list` 字段，旧包仍回退默认产品。
- 任务导入按任务类型恢复 `data_list` 或 `scene_list`，同名依赖被跳过时仍保留关联；`skip` 保留同名产品、应用、系统参数、断言模板、附件与语料键，`overwrite` 才更新。产品和应用的冲突纳入检查结果。
- `env_config.yaml` 始终导出，即使编排和任务均未选择产品，也迁移应用、系统参数和断言；应用可携带可选 `product` 与 `thread_number`，保留独立配置，旧包仍回退原默认值。元信息数据库写入以单次配置为事务，写入失败返回错误、回滚该次配置，编排导入状态记录为 `error`，不创建依赖该配置的编排；不增加跨编排事务。
- 导出前做依赖校验：**任务、场景、数据文件**缺失直接终止导出并返回具体缺失项；**产品、环境配置**（按产品关联的 app 判断）缺失只写告警日志，不阻断导出；
- 导入**不重新实现依赖解析**，而是从 `bundle/` 中挑出真正包含 `task_info.yaml` 的归档，交给现有 `ImportScheduleCheck` / `ImportScheduleConfirm` 处理。

**导入两阶段与状态**：

| 阶段 | 接口 | 行为 |
|---|---|---|
| 检查 | `POST /task_orchestration_import_check` | 解压 zip、解析清单、调用任务导入检查；返回编排同名冲突、任务映射（已存在/将导入/缺失）、缺失场景/数据文件（阻断）与缺失产品（仅提示）；**存在阻断项时返回 `blocked=true`，页面隐藏确认区** |
| 确认 | `POST /task_orchestration_import_confirm` | 先执行 `ImportScheduleConfirm` 落依赖，再按任务名映射逐条创建/跳过/覆盖编排；同名编排在 `skip` 模式下内容不一致时返回错误、要求显式覆盖；**单个编排失败只回滚本条**，不影响其他编排 |
| 状态 | `GET /task_orchestration_import_status?import_id=` | 状态文件 `DownloadBasePath/orchestration_import_<id>.json`，取值 `pending / processing / done / error` |
| 文案 | `GET /task_orchestration_import_labels` | 导入页按当前 `data_locale` 取文案；静态页保留中文兜底，仅在 labels 加载失败时显示 |

实现约束（沿用定时任务现有写法）：

- 导入页用 `r.GET` + `requireOrchestrationImport` 注册（登录 + `GET /task_orchestration_import` 权限 309），校验通过后 `c.File` 返回静态页；
- 行操作/顶部按钮用 `info.AddActionButton` / `info.AddButton`，批量导出从 `ctx.FormValue("ids")` 取参、单条导出从 `ctx.FormValue("id")` 取参（参照 `tables/schedule.go:196-254`）；
- 导出是 `action.Ajax` 返回可点击下载链接，不是文件流响应；导出文件必须落在 `mgmt/download` 下，filemanager 只放行路径含 `mgmt` 的下载（`main.go:232-237`）；
- action id 使用 `task_orchestration_export` / `task_orchestration_batch_export`，不要与现有 `schedule_export` / `schedule_batch_export` 重名。

在**定时任务列表**增加：

- 行操作/批量操作 **「加入编排」**：弹窗选择目标任务编排，把选中任务名称按顺序追加到 `task_list` 末尾（已存在则跳过）。
- 该入口是"基于任务组装编排"的快捷方式，编排页内的步骤编排为主入口。
- 两个加入编排接口均要求有效登录，并校验 `POST /edit/task_orchestration` 权限；只读角色或仅有任务编辑权限的角色不可通过接口修改编排。

### 4.3 新建 / 编辑表单

| 字段 | 控件 | 规则 |
|---|---|---|
| 编排名称 | Text | 必填、唯一 |
| 关联任务 | TextArea | 每行一个任务名称，行顺序即步骤顺序；保存时校验任务存在、去重、至少 1 个 |
| 产品配置 | Select（多选） | 选填；选项来自产品列表，存 `product_list`（逗号分隔）。**配置后覆盖所有步骤任务的产品**；为空时沿用各任务自身产品 |
| 备注 | TextArea | 选填 |
| 状态/上次执行时间/创建人/时间 | 只读 | 由系统维护 |

本期表单**不展示**编排模式、cron、执行方式、失败策略等字段（数据库中保留预留列，见 5.1）。

保存时：

- 校验编排名称不重复（`deleted_at IS NULL` 范围内）；
- 校验任务名称全部存在于 `schedule` 且未删除；
- 校验 `task_list` 内无重复任务名；
- 单条写入 `task_orchestration`（`task_list` 按行归一化：去空行、去首尾空格），单语句本身原子；校验在写入前执行，失败不落库；
- 编排执行中（存在 `orchestration_status=running`）时禁止修改步骤，避免执行中配置漂移。

### 4.4 步骤编排

- 表单内 TextArea，每行 `任务名称`（可带行号展示），保存后整体写入 `task_orchestration.task_list`，行序即步骤顺序。
- 详情页解析 `task_list` 展示有序步骤表：序号、任务名称（链接到任务详情）、任务类型、关联产品、当前是否可用（任务是否被删除或改名）。
- 执行时按行顺序逐个解析为 `schedule.id`；任一任务无法解析则该轮不执行并报错（见 6.3）。
- 每个步骤执行前完整解析该任务配置的场景或数据文件；任一依赖缺失或软删除时，该步骤失败且不执行部分依赖，保留错误与任务报告，并沿用普通失败继续策略。空执行项或未知任务类型同样返回失败。
- 后续增强（P1）：独立编排页、拖拽排序、步骤启用/禁用。

### 4.5 详情页

- 基本信息：编排名称、关联任务、产品配置、备注、状态、上次执行时间、创建人、创建时间。
- 步骤列表：解析 `task_list` 有序展示。
- 最近执行记录：查询 `dashboard`（`report_type='orchestration' AND related_task_ids LIKE 'orch:<orchestrationId>:%'`）最近 10 次，展示时间、状态、耗时、成功/失败/跳过步骤数，点击进入对应报告。

### 4.6 执行与失败处理（产品行为）

- **立即执行**：点击后立即开始一轮串行执行，不改变任何调度状态（本期无调度）。
- **串行**：按 `task_list` 行序逐个执行，上一步结束后执行下一步。
- **普通失败继续**：某步骤用例失败、接口非 401 错误、脚本失败时，记录该步骤失败并**继续**执行后续步骤。
- **401 停止**：`scene` 或 `data` 类型步骤返回 401 鉴权失败时，记录失败，**后续步骤标记"跳过"并停止执行**。
- **不提供停止/中断**：执行开始后只能等待完成，与现有任务不支持中断执行保持一致。
- 同一编排同一时间只允许一轮执行；重复点击提示"编排正在执行"。

### 4.7 编排报告

入口：

- **「结果报告列表」**（`/admin/info/dashboard`）：报告类型新增「任务编排」，每轮执行完成后自动出现在列表中；
- 编排列表行操作「报告」：跳转该编排最近一次执行的报告；
- 编排详情「最近执行记录」：选择具体某一次执行。

页面路由：复用 `/admin/report_dashboard?id=<dashboard_id>`，`GetDashboardByReportId` 增加 `case "orchestration"` 渲染编排报告；`/admin/orchestration_dashboard?id=<orchestration_id>` 作为兼容入口，直接跳转该编排最近一次执行的 `dashboard_id`。

页面内容：

1. **概览**：编排名称、触发方式（manual）、开始/结束时间、总耗时、步骤总数、成功/失败/跳过数、执行率。
2. **步骤明细**：序号、任务名称、状态（成功/失败/跳过）、开始/结束/耗时、通过/失败用例数、错误信息、「查看任务报告」链接（`/admin/report_dashboard?id=<dashboard_id>`）。
3. **产品维度汇总**：产品、总数、通过、失败、通过率（按本轮实际使用的产品聚合）。
4. **历史执行选择**：最近 20 次执行记录切换。
5. **趋势图**（P1）。

编排报告与任务报告统一在「结果报告列表」中查询，下钻关系保持不变。

### 4.8 权限与 i18n

- 新增菜单后，在 GoAdmin「权限管理」中补充：
  - `GET /info/task_orchestration`（查询）、`/info/task_orchestration/edit`、`/info/task_orchestration/new`；
  - `POST /new/task_orchestration`、`/edit/task_orchestration`、`/delete/task_orchestration`、`/export/task_orchestration`；
  - 编排立即执行/复制/导入导出等 Ajax 操作跟随页面权限；
  - `/admin/orchestration_dashboard`、`/admin/task_orchestration_import` 页面权限；
  - 在「角色管理」中分配给对应角色。
- 菜单与权限种子以 `INSERT INTO goadmin_menu / goadmin_permissions / goadmin_role_permissions` 形式**追加到 `mgmt/sql/update.sql`**，不修改 `init_all_*.sql`。现有种子位于 `init_all_20260916.sql:1143/1153/1183`；权限模板应参照 `test_case` 的 7 条（同文件 1153 区块 80-86 行），`case_statistics` 只有菜单行、没有权限行。新增 id 需避开现有最大值（menu ≥2765、permission ≥301），并显式分配给目标角色。
- 菜单英文翻译：`mgmt/i18n/i18n_menu/en-US.yaml` 增加 `任务编排: Task Orchestration`。
- 业务文案：`mgmt/i18n/i18n_system/zh-CN.json`、`en-US.json` 新增 `task_orchestration.*`、`error.orchestration_*`（见第 8 节）。

## 5. 数据模型设计

### 5.1 `task_orchestration`

```sql
CREATE TABLE `task_orchestration` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `orchestration_name` varchar(255) NOT NULL COMMENT '编排名称',
  `task_list` longtext COMMENT '关联任务名称，换行分隔，行序即步骤顺序',
  `product_list` varchar(64) DEFAULT NULL COMMENT '编排级产品(逗号分隔)，非空时覆盖所有步骤任务的产品',
  `orchestration_status` enum('running','stopped','finished','not_started') NOT NULL DEFAULT 'not_started' COMMENT '编排状态',
  `orchestration_mode` enum('cron','once','day','week') NOT NULL DEFAULT 'once' COMMENT '预留：编排模式(本期仅手动立即执行)',
  `crontab` varchar(64) DEFAULT NULL COMMENT '预留：cron表达式',
  `week` varchar(64) DEFAULT NULL COMMENT '预留：每周几',
  `time4week` varchar(64) DEFAULT NULL COMMENT '预留：每周执行小时',
  `time4day` varchar(64) DEFAULT NULL COMMENT '预留：每天执行小时',
  `remark` longtext COMMENT '备注',
  `last_at` datetime DEFAULT NULL COMMENT '上次执行时间',
  `next_at` datetime DEFAULT NULL COMMENT '预留：下次执行时间',
  `user_name` varchar(100) DEFAULT NULL COMMENT '创建人',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orchestration_name` (`orchestration_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='任务编排';
```

- 调度相关列（`orchestration_mode` / `crontab` / `week` / `time4week` / `time4day` / `next_at`）为**预留字段**，本期不使用、不展示，后续做定时触发时直接启用。
- 不建步骤表：编排信息全部在 `task_list`；重复任务由保存校验拦截。
- 不建执行记录表：执行记录见 5.2。
- `orchestration_status` 本期取值：`not_started`（未执行）、`running`（执行中）、`finished`（本轮结束，失败也置 finished，结果以报告为准）；`stopped` 仅预留。
- `orchestration_name` 唯一约束叠加软删除会导致"删除后同名编排无法重建"（`schedule` 表目前同样存在该问题）。实现约定：删除时先把 `orchestration_name` 改写为 `<原名>_deleted_<yyyymmddHHMMSS>` 再置 `deleted_at`。

### 5.2 与现有表/存储的关系

| 表/存储 | 关系 |
|---|---|
| `schedule` | 步骤通过 `task_list` 中的任务名称解析为 `task_id`；任务被删除或改名后该步骤不可用 |
| `dashboard` | ① 每个步骤执行一个任务，产生一条任务报告（`report_type=task`，由现有 `CreateDashboardRecord` 生成），报告中的 `steps[].dashboard_id` 记录其 id 用于下钻；② 每轮编排执行生成**一条**编排报告（`report_type=orchestration`），`related_task_ids = orch:<orchestrationId>:<runTag>` 用于按编排检索 |
| `HistoryBasePath` 历史文件 | 每个数据文件/场景执行时照旧写入，粒度与任务执行一致，不额外改动 |
| `goadmin_menu` / `goadmin_permissions` | 新增菜单与页面/CRUD 权限并分配角色 |

**`dashboard` 不新增任何列**。编排报告通过 `related_task_ids` 的 `orch:` 前缀与编排关联；查询方式：

```sql
SELECT * FROM dashboard
WHERE report_type = 'orchestration' AND related_task_ids LIKE 'orch:<orchestrationId>:%'
ORDER BY created_at DESC;
```

现有任务报告查询使用 `related_task_ids LIKE '<taskId>!_%' ESCAPE '!'`（`pages/schedule_dashboard.go:64`、`biz/task_report.go:268`），任务 id 为纯数字，`orch:` 前缀不会与之冲突。

### 5.3 迁移

DDL 与菜单/权限种子**只追加到 `mgmt/sql/update.sql` 末尾**；`mgmt/sql/init_all_*.sql` 已定型发布，不修改。

## 6. 执行引擎设计

### 6.1 代码组织

按项目约定成对新增：

- `biz/task_orchestration_define.go`：结构体、常量、`task_orchestration` 类型与报告 JSON 类型；
- `biz/task_orchestration.go`：编排 CRUD、`task_list` 解析与保存、复制、删除；
- `biz/task_orchestration_run.go`：执行引擎（串行、失败继续、401 停止、互斥）；
- `biz/task_orchestration_report.go`：编排报告数据聚合与 `dashboard` 写入；
- `biz/task_orchestration_transfer_define.go`、`biz/task_orchestration_transfer.go`：导入/导出 DTO、依赖闭包、检查与确认；
- `biz/archive_import.go`：归档解压与安全限额（同时被任务导入复用）；
- `biz/task_run_dependencies.go`：任务执行前的依赖解析与校验；
- `biz/run_errors.go`：执行错误合并工具；
- `orchestration_join.go`：任务列表「加入编排」接口（含权限守卫）；
- `tables/task_orchestration.go`：GoAdmin 列表/表单/详情；
- `pages/task_orchestration_dashboard.go`：编排报告渲染与最近一次报告跳转；
- `html/task_orchestration_import.html`：导入页。

### 6.2 核心类型与函数

```go
// task_orchestration_define.go
type DbTaskOrchestration struct {
	Id                  string `gorm:"column:id" json:"id"`
	OrchestrationName   string `gorm:"column:orchestration_name" json:"orchestration_name"`
	TaskList            string `gorm:"column:task_list" json:"task_list"`
	ProductList         string `gorm:"column:product_list" json:"product_list"`
	OrchestrationStatus string `gorm:"column:orchestration_status" json:"orchestration_status"`
	Remark              string `gorm:"column:remark" json:"remark"`
	LastAt              string `gorm:"column:last_at" json:"last_at"`
	UserName            string `gorm:"column:user_name" json:"user_name"`
	CreatedAt           string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           string `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt           string `gorm:"column:deleted_at" json:"deleted_at"`
}

// 编排报告 JSON（存 dashboard.report_data）
type OrchestrationReportData struct {
	OrchestrationId   string                `json:"orchestration_id"`
	OrchestrationName string                `json:"orchestration_name"`
	RunTag            string                `json:"run_tag"`
	TriggerType       string                `json:"trigger_type"` // 本期固定 manual
	ProductList       string                `json:"product_list"`
	Summary           OrchestrationSummary  `json:"summary"`
	Steps             []OrchestrationStep   `json:"steps"`
	Products          []OrchestrationProd   `json:"products"`
}

// 核心函数
func RunOrchestration(orchestrationId, operator string) (dashboardId string, err error)
func RunOrchestrations(orchestrationIds, operator string) ([]string, error) // 批量立即执行
func CopyTaskOrchestration(id, userName string) error
func DeleteTaskOrchestration(id, userName string) error
func ParseTaskList(taskList string) []string
func GetOrchestrationLatestReport(orchestrationId string) (string, error)
```

任务侧抽取（关键改造）：

```go
// biz/biz_task.go
type TaskRunOptions struct {
	RunTag           string // 编排传入 <taskId>_<orchRunTag>_<stepOrder>；为空时沿用现有规则
	ProductOverride  string // 编排级产品（逗号分隔），非空时覆盖任务自身 product_list
	UpdateTaskStatus bool   // 编排步骤传 false，不改写 schedule.task_status
	SyncReport       bool   // 编排步骤传 true，同步生成任务报告
	StrictExecution  bool   // 编排步骤传 true：依赖强校验、data 取首产品、报告创建失败即中止
}

type TaskRunResult struct {
	TaskTag         string
	DashboardId     string
	Status          string // success / fail
	Total           int
	SuccessCount    int
	FailCount       int
	StartAt         string
	EndAt           string
	DurationSeconds int
	Err             error // 执行错误，必须保留原始错误（见 6.4）
}

func RunTaskOnce(taskId, userName string, opts TaskRunOptions) (TaskRunResult, error)

// 原 RunOnceTask 保留为兼容包装，行为不变
func RunOnceTask(id, userName string) error {
	_, err := RunTaskOnce(id, userName, TaskRunOptions{UpdateTaskStatus: true})
	return err
}
```

抽取后 `RunTaskOnce` 需返回步骤结果与任务报告 id，供编排报告聚合与下钻。

### 6.3 执行流程

```text
点击「立即执行」
   │
   ├─ 编排互斥：同一编排已有执行 → 提示"编排正在执行"
   ├─ 读取编排；解析 task_list（去空行）
   ├─ 前置校验：步骤非空；任务名全部能解析到未删除的 schedule；否则整轮不执行并报错
   ├─ 解析本轮产品：product_list 非空 → 编排产品；为空 → 步骤任务 product_list 并集去重
   ├─ 生成 orchRunTag = orch<orchestrationId>_<yyyymmddHHMMSS>
   ├─ 创建编排报告 dashboard(report_type=orchestration, status=generating,
   │     related_task_ids=orch:<orchestrationId>:<yyyymmddHHMMSS>, related_products=本轮产品)
   ├─ 置 orchestration_status=running
   ├─ 逐个步骤串行执行 RunTaskOnce(taskId, operator, opts{
   │       RunTag: "<taskId>_<orchRunTag>_<stepOrder>",
   │       ProductOverride: 本轮产品,
   │       UpdateTaskStatus: false })
   │    ├─ 成功 → 记录步骤成功，继续
   │    ├─ 普通失败 → 记录步骤失败，继续
   │    └─ 步骤命中 401（scene / data）→ 记录失败，后续步骤标记 skipped，停止
   ├─ 汇总写回编排报告 dashboard：report_data=OrchestrationReportData、status=finished
   ├─ 更新 task_orchestration：orchestration_status=finished、last_at=结束时间
   └─ 释放编排互斥
```

### 6.4 失败处理与 401 停止

| 情形 | 行为 |
|---|---|
| 步骤成功 | 继续下一步骤 |
| 步骤普通失败（用例失败、非 401 接口错误、脚本失败） | 记录失败，**继续**执行后续步骤 |
| 步骤返回 401 鉴权失败（`scene` / `data`） | 记录失败，**后续步骤标记 skipped 并停止** |

判定实现（**沿用当前系统 scene 与 data 的既有实现**）：

1. **统一判定方式**：优先匹配当前语种文案 `biz.T("error.auth_fail_prefix")`，并回退到**与语种无关**的正则 `(?i)(返回码|status code)\s*[:：]?\s*401`，避免异步执行期间其他请求切换全局语种导致漏判；
2. **scene 路径**：`RunPlaybookFromMgmt` 返回的 error 已保留 401 原文，直接判定；任务内部"仅首个场景判定 401"的既有逻辑保持不变。
3. **data 路径**：`RepeatRunDataFile` 当前会在 `result != "pass"` 时用 `test %s` 覆盖原始错误（串行分支 `biz/data_run.go:955-957`），并行分支同样丢弃 `err1`（`:900-903`）。按现有错误保留方式修正：**仅在 `err == nil` 时才写 `test %s`**，并行分支保留 `err1`，使 data 与 scene 一样能把 401 原文带回上层。
4. **`RunTaskOnce` 返回执行错误**：当前 `biz/biz_task.go:418-419` 用 `schedule` 更新的 error 覆盖了执行错误，必须改为独立变量保存执行错误并返回，否则编排层拿不到 401 信号。
5. **编排层判定范围**：只看该步骤整体返回的 error，任一步骤（scene 或 data）命中前缀即停止后续步骤；不改动任务自身的执行语义。
6. 已知限制：判定依赖 i18n 文案前缀。执行与判定在同一进程、使用同一 locale 的 `biz.T()`，文案一致；若后续修改该文案需同步调整。

### 6.5 复用任务执行与改造

- 编排步骤统一按**任务立即执行语义**运行（抽取后的 `RunTaskOnce`），不使用 `OneTask`（cron 专用路径），保证行为可预期且能拿到任务报告 id。注意：管理页"运行"按钮入口实际是 `RunTask`（`tables/schedule.go:305`），只有 `once` 模式才落到 `RunOnceTask`。
- `RunTaskOnce` 内部：
  - `effectiveProductList := opts.ProductOverride`；为空时取 `taskDb.ProductList`。执行与任务报告（`CreateDashboardRecord`）均使用该有效产品，保证"编排级产品优先"；
  - `taskTag` 支持外部传入，格式 `<taskId>_<orchRunTag>_<stepOrder>`，仍以 `<taskId>_` 开头，保证现有"任务报告按 `related_task_ids like '<taskId>_%'`"的查询不受影响；
  - 返回本次任务报告 `dashboard_id`，供编排步骤下钻；
  - `UpdateTaskStatus=false` 时步骤执行不改写 `schedule.task_status`，避免编排执行影响任务自身状态（现有 `RunOnceTask` 是**无条件**写 `finished`，该分支为新增逻辑）。
- 每个步骤仍会产生一条**任务报告** `dashboard`（沿用 `CreateDashboardRecord` + 异步 `GenerateTaskReport`），这是下钻的数据来源，属于"统一数据模型"的预期行为。
- **产品覆盖对 `data` 类型步骤的约束**：`data` 链路 `RepeatRunDataFile → RunDataFile → GetEnvConfig` 按单产品精确等值查询（`biz/app.go:55-57`），传入逗号分隔的多产品会查不到环境、丢失 host/auth。因此 `data` 步骤只取有效产品的**第一个**；`scene` 步骤由 `RunPlaybookFromMgmt` 遍历全部产品（`biz/biz_playbook.go:103-118`）。
- `CreateDashboardRecord` 与 `GenerateTaskReport` 都从 `DbSchedule.ProductList` 读取产品（`biz/task_report.go:322`、`:34`、`:175`），因此实现时复制一份 `taskDb` 并改写其 `ProductList` 即可贯穿执行与报告，**无需修改这两个函数的签名**。
- 现有 `RunOnceTask` 改为包装函数，**调用方式保持不变**；新增 `StrictExecution` 选项：仅编排步骤传 `true`（依赖强校验、`data` 取首产品、报告创建失败即中止），普通任务与 CLI 保持原有行为，详见 §9 边界与限制；

### 6.6 互斥

- 包级 `sync.Map` 用 `LoadOrStore` 占位（`orchestrationId -> *struct{}`），`loaded=true` 即视为已有执行，执行结束 `Delete`。
- **不使用 `sync.Mutex.TryLock`**：该 API 要求 Go 1.18+，而当前 `go.mod` 为 `go 1.17`，新工具链下会报 `requires go1.18 or later`。
- 服务重启后内存状态自然释放；若有编排在重启前处于 `running`，启动时把 `orchestration_status=running` 重置为 `finished`，避免状态悬挂；同时把遗留的 `dashboard.status=generating` 编排报告置为 `failed`，页面显示"执行已中断"。
- 页面「立即执行」调用 `StartOrchestrationAsync`：**同步抢锁成功后才返回"已开始"**，抢锁失败直接返回"编排正在执行"；批量执行逐个同步抢锁，任一条失败立即返回。

## 7. 报告设计

编排报告统一存储在 `dashboard.report_data`（`report_type=orchestration`），不新增表、不新增列。

JSON 结构：

```json
{
  "orchestration_id": "3",
  "orchestration_name": "版本回归编排",
  "run_tag": "orch3_20261105103000",
  "trigger_type": "manual",
  "product_list": "产品A,产品B",
  "summary": {
    "total": 3,
    "success": 2,
    "fail": 1,
    "skipped": 0,
    "duration_seconds": 1234,
    "pass_rate": 66.7
  },
  "steps": [
    {
      "step_order": 1,
      "task_id": "12",
      "task_name": "登录冒烟",
      "status": "success",
      "start_at": "2026-11-05 10:30:00",
      "end_at": "2026-11-05 10:30:30",
      "duration_seconds": 30,
      "dashboard_id": "88",
      "total": 10,
      "success": 10,
      "fail": 0,
      "error_msg": ""
    }
  ],
  "products": [
    {"product": "产品A", "total": 40, "pass": 36, "fail": 4, "pass_rate": 90.0}
  ]
}
```

- 生成时机：所有步骤结束后同步聚合，先建记录再回填：
  - `report_type = 'orchestration'`
  - `report_name = biz.T("dashboard.orchestration_report_name")` 填充，中文为 `任务编排[%s]执行报告_%s`，与现有 `dashboard.task_report_name`（`biz/task_report.go:321`）同机制，不硬编码中文
  - `related_products = 本轮实际使用的产品`
  - `related_task_ids = orch:<orchestrationId>:<yyyymmddHHMMSS>`（编排标识，**不是**任务 id 列表）
  - `report_data = 上述 JSON`，`status` 从 `generating` 变为 `finished`
  - `creator = 操作人`
- 结果报告列表：`tables/dashboard.go` 需改**三处**——列表类型展示 switch（`:31-42`）、筛选 `FieldFilterOptions`（`:43-48`）、详情页类型展示 switch（`:131-143`），新增 `orchestration → 任务编排`，并补 i18n key `schedule_report.report_type_orchestration`。
- 展示：`pages.GetDashboardByReportId`（`pages/schedule_dashboard.go:17`）增加 `case "orchestration"`，渲染编排报告；步骤行链接到任务报告 `/admin/report_dashboard?id=<dashboard_id>`。
- 离线导出：`pages/report_export.go:60-69` 的 `renderReportPanel` 已增加 `case "orchestration"`，渲染函数识别 `offlineReportExport` 标记（`pages/report_export.go:26`），离线渲染能力就绪。**本期 CLI 无编排执行入口**，编排报告不会进入 CLI 本地库，端到端导出留待后续版本。
- GoAdmin 详情页：`tables/dashboard.go:150-153` 把 `report_data` 渲染成任务专用的 `/admin/schedule_report?report_id=`，编排报告需按报告类型分流到 `/admin/report_dashboard?id=`。
- 编排列表「报告」：查询该编排最近一条编排报告并跳转；无执行记录时提示暂无报告。
- 历史执行选择：在编排报告页内提供最近 20 次执行的下拉切换（按 `related_task_ids LIKE 'orch:<orchestrationId>:%'` 查询）。
- 复用现有图表能力（Chart.js，`pages/schedule_dashboard.go` 已使用）展示步骤成功/失败分布与通过率趋势（P1）。

## 8. 接入点与文件清单

新增：

- `biz/task_orchestration_define.go`、`biz/task_orchestration.go`、`biz/task_orchestration_run.go`、`biz/task_orchestration_report.go`
- `biz/task_orchestration_transfer_define.go`、`biz/task_orchestration_transfer.go`（导入/导出 DTO、依赖闭包、检查与确认）
- `biz/archive_import.go`、`biz/task_run_dependencies.go`、`biz/run_errors.go`、`orchestration_join.go`
- `biz/orchestration_p1_test.go`、`orchestration_join_test.go`（编排单测）
- `tables/task_orchestration.go`
- `pages/task_orchestration_dashboard.go`（编排报告渲染 + 最近一次报告跳转）
- `html/task_orchestration_import.html`（上传检查、冲突选择、导入结果）

修改：

- `tables/tables.go`：`"task_orchestration": GetTaskOrchestrationTable`
- `tables/dashboard.go`：报告类型展示/筛选新增 `orchestration → 任务编排`（列表 switch、筛选 options、详情 switch 三处）；详情页链接按报告类型分流
- `pages/schedule_dashboard.go`：`GetDashboardByReportId` 增加 `case "orchestration"` 分支
- `pages/report_export.go`：`renderReportPanel` 增加 `case "orchestration"`，并处理 `offlineReportExport` 链接
- `biz/cli.go`：`merge` 命令在解析 `related_task_ids` 前按 `report_type='task'` 过滤，避免把 `orch:` 开头的编排报告当作任务报告
- `main.go`：
  - `GET /admin/orchestration_dashboard`：按 orchestration_id 跳转最近一次编排报告
  - `GET /admin/task_orchestration_import`：用 `eng.HTMLFile` 注册编排导入页面
  - `POST /task_orchestration_import_check`、`POST /task_orchestration_import_confirm`
  - `GET /task_orchestration_import_status`：导入状态查询
  - `GET /task_orchestration_import_labels`：导入页文案（按当前语种）
- `biz/biz_task.go`：抽取 `RunTaskOnce`（支持 `RunTag` / `ProductOverride` / `UpdateTaskStatus`），`Init()` 中增加编排 `running` 状态重置
- `biz/data_run.go`：`RepeatRunDataFile` 保留原始错误（串行 `:955-957`、并行 `:900-903`）
- `biz/task_define.go`：`TaskRunOptions` / `TaskRunResult`（或放入 `task_orchestration_define.go`）
- `tables/schedule.go`：新增「加入编排」行操作/批量操作
- `mgmt/sql/update.sql`：`task_orchestration` 建表 DDL + 菜单(2765) + 权限(301-309，含导入页) + `goadmin_role_permissions` / `goadmin_role_menu` 授权（仅追加到该文件，不改动 `init_all_*.sql`）
- `mgmt/i18n/i18n_system/zh-CN.json`、`en-US.json`
- `mgmt/i18n/i18n_menu/en-US.yaml`
- 实施完成后更新 `mgmt/doc/file/update/change_log.md`

i18n key 建议：

```text
task_orchestration.title                         任务编排
task_orchestration.description                   将多个定时任务组装为按顺序执行的编排
task_orchestration.detail_title                  任务编排详情
task_orchestration.orchestration_name            编排名称
task_orchestration.task_list                     关联任务
task_orchestration.step_count                    步骤数
task_orchestration.product_list                  产品配置
task_orchestration.last_at                       上次执行时间
task_orchestration.btn_run                       立即执行
task_orchestration.btn_copy                      复制
task_orchestration.btn_report                    报告
task_orchestration.btn_add_to_orchestration      加入编排
task_orchestration.msg_run_started               编排已开始执行
schedule_report.report_type_orchestration        任务编排
dashboard.orchestration_report_name              任务编排[%s]执行报告_%s
error.orchestration_not_found                    编排[%s]不存在
error.orchestration_no_tasks                     编排[%s]未配置任务
error.orchestration_task_missing                 编排[%s]中的任务[%s]不存在
error.orchestration_task_duplicate               编排[%s]中存在重复任务[%s]
error.orchestration_running                      编排[%s]正在执行，请稍后再试
error.orchestration_save_tasks                   保存编排任务失败: %s
```

## 9. 边界与限制

- **存量任务行为保持不变**：`StrictExecution` 仅在编排步骤为 `true`——编排步骤做依赖强校验（缺失/为空即失败）、`data` 只取首产品、任务报告创建失败即中止；普通任务与 CLI 仍沿用旧行为（缺失依赖截断执行、空列表空跑成功、报告创建失败继续执行、`data` 传完整产品列表），保证 AC-19 存量零回归。
- **编排步骤同步生成任务报告**：编排执行时以 `SyncReport=true` 同步生成任务报告，以便汇总产品维度；普通任务执行仍为异步生成。
- **仅串行、仅手动**：本期没有并行执行、没有 cron 调度、没有 CLI 触发；调度相关列仅预留。
- **执行不可中断**：开始后只能等待本轮结束，不提供停止。
- **失败继续 + 401 停止**：普通失败不阻断后续步骤；`scene` 与 `data` 步骤的 401 鉴权失败都会停止后续步骤。
- **401 判定与现有实现对齐**：`scene` 沿用 `biz/biz_task.go:408` 的文案前缀判断；`data` 需修正 `RepeatRunDataFile` 覆盖原始错误的问题，二者使用同一判定方式。
- **执行错误不能被覆盖**：`RunOnceTask` 末尾用 `schedule` 更新的 error 覆盖执行错误，抽取 `RunTaskOnce` 时必须拆分变量，否则编排层无法判断 401。
- **编排级产品优先**：编排配置了产品时，步骤执行与任务报告均使用编排产品；未配置时才使用任务自身产品。`data` 步骤只取第一个产品，多产品完整生效仅限 `scene` 步骤。
- **任务改名/删除**：`task_list` 存任务名称，任务改名后该步骤无法解析，编排执行前置校验失败；建议后续给任务改名/删除增加"被编排引用"提示。
- **同名编排重建**：删除时改写 `orchestration_name` 为 `<原名>_deleted_<时间戳>`，否则同名编排无法重建。
- **报告关联**：编排报告不写任务 id 到 `related_task_ids`，而写 `orch:<orchestrationId>:<runTag>`；任何按任务 id 解析该字段的代码必须先按 `report_type` 过滤。
- **导入依赖包选择**：导出 zip 中可能同时存在附件包与数据包，导入端按"归档内含 `task_info.yaml`"挑选，不依赖文件前缀（前缀随语种变化）。
- **归档安全**：外层 ZIP 与内层任务依赖包均限制路径穿越、文件数（2000）与解压总大小（500MB），拒绝链接、重复文件路径与损坏条目；依赖包探测过程也受限；**未做扩展名白名单**，因为平台本身需要导入 `.py/.sh/.bat/.jmx` 等脚本文件。
- **导入接口权限**：导入页与 `check / confirm / status / labels` 均校验登录及 `GET /task_orchestration_import` 权限（309）；`import_id` 强制 UUID，避免路径穿越。
- **动作权限**：立即执行/复制按 `/edit/task_orchestration`、导出按 `/export/task_orchestration` 校验权限，无权限返回 403。
- **历史功能不改造**：定时任务导入（`/schedule_import_check`、`/schedule_import_confirm`）属于历史功能，沿用既有实现（无鉴权、无 import_id 校验），**本期不做优化**；鉴权与路径加固仅覆盖任务编排导入。
- **CLI 限制**：CLI 使用本地 `cli.db`，本期没有编排执行入口，编排报告实际不会进入 CLI 本地库；`report_export` 的 `case "orchestration"` 仅作为后续可用能力保留。
- **权限**：页面操作全部受 RBAC 控制；导入页已补权限与角色菜单。

## 10. 验证方案

1. **编译**：`make build` 或 `CGO_ENABLED=0 go build .`。P1 回归测试：`go test -vet=off -exec "$PWD/build/test-exec.sh" ./...`（需 CGO/SQLite；执行器让包初始化从仓库根目录读取配置）。
2. **页面 CRUD**：新建编排（3 个步骤）、编辑、复制、删除；校验任务不存在/重复/空编排被拒绝。
3. **串行执行**：3 个步骤按 `task_list` 行序执行，步骤之间不并发。
4. **失败继续**：第 2 步构造普通失败，确认第 3 步仍执行，编排报告显示 1 失败 2 成功。
5. **401 停止**：第 2 步为 `scene` 任务并构造 401，确认第 3 步为 skipped，报告可见失败原因。
6. **data 步骤 401**：构造 data 任务返回 401，确认同样触发停止，且错误原文未被 `test <result>` 覆盖。
7. **立即执行**：只产生一条 `report_type=orchestration` 报告与 N 条任务报告，不涉及任何调度状态。
8. **互斥**：同一编排并发触发两次，仅一轮执行，另一次提示"编排正在执行"。
9. **编排级产品覆盖**：编排配置产品 A、步骤任务配置产品 B，执行后请求发往 A，任务报告 `related_products` 为 A；清空编排产品后回退到 B。另验证 data 步骤只使用第一个产品且能取到环境配置。
10. **任务状态不受影响**：编排执行后 `schedule.task_status` 保持不变。
11. **报告统一**：编排执行完成后，「结果报告列表」出现类型为"任务编排"的报告，可筛选、打开并下钻任务报告。
12. **离线渲染能力**：构造一条 `report_type=orchestration` 报告，确认 `renderReportPanel` 能正常输出 HTML（CLI 本地库入口留待后续，不在本期验收范围）。
13. **同名编排重建**：删除编排后可用同一名称重新创建成功。
14. **加入编排**：任务列表选中任务追加到编排末尾，重复任务被跳过。
15. **导出依赖校验**：单条/批量导出后检查 zip 结构（`orchestration_info.yaml` + `bundle/` 下含 `task_info.yaml` 的数据包），确认编排、步骤任务、场景、数据、产品、配置及附件齐全，且关联不依赖源环境数据库 ID；分别构造缺失场景、缺失数据文件、缺失产品，确认前两者终止导出、后者仅告警。
16. **导入检查**：上传归档后返回同名冲突、任务映射（已存在 / 将导入 / 缺失）与缺失项，不写业务表；构造缺失任务时检查阶段即阻断；分别验证 skip、overwrite、cancel。
17. **导入独立性**：构造某条编排任务缺失，确认该编排不落库，同包其他完整编排仍成功导入。
18. **归档安全**：验证路径穿越、文件数超 2000、解压后超 500MB 及损坏归档均被拒绝，临时目录被清理。
19. **回归**：未使用编排的既有任务执行、报告、cron 行为不变。

## 11. 实施拆分

**阶段一：数据与页面（可独立验收）**

- `task_orchestration` DDL + 菜单/权限/i18n；
- `tables/task_orchestration.go` 列表/表单/详情 + `task_list` 解析与保存；
- 复制/删除（含删除改名）。

**阶段二：执行引擎**

- 抽取 `RunTaskOnce`（返回步骤结果、任务报告 id 与执行错误）；
- 串行执行、失败继续、`scene` / `data` 的 401 停止；
- 编排报告记录创建与汇总写回、互斥、`last_at` 更新。

**阶段三：报告接入**

- `pages/task_orchestration_dashboard.go` 渲染 + `GetDashboardByReportId` 的 `case "orchestration"`；
- `tables/dashboard.go` 三处类型展示 + `pages/report_export.go` 离线导出分支；
- 编排列表/详情「报告」入口、产品汇总。

**阶段四：导入/导出**

- 编排及任务依赖闭包收集、归档生成；
- 上传检查、冲突分类与名称映射；
- 异步确认导入、导入任务状态持久化与页面轮询、失败回滚及归档安全校验。

## 12. 已确认口径与后续

2026-10-23 首次确认，2026-11-05 按新范围重定（概念统一为"任务编排"），以下口径本期不再变更：

1. **只新增 `task_orchestration` 一张表**：任务编排存 `task_list`（换行分隔任务名），不建步骤表、不建执行记录表。
2. **执行记录复用定时任务模型**：每步一条任务报告 + 每轮一条编排报告（`report_type=orchestration`），`dashboard` 不新增列。
3. **仅串行、仅立即执行**：不做并行、不做 cron、不做运行/停止按钮、不做 CLI；调度字段在表中预留。
4. **失败继续 + 401 停止**：不保留失败策略字段；`scene` 与 `data` 步骤的 401 均触发停止，判定沿用系统现有的文案前缀实现。
5. **不需要统一鉴权刷新**：编排执行前不做产品鉴权刷新。
6. **编排级产品优先**：编排未配置时回退任务自身产品；`data` 步骤只取首个产品。
7. **同名编排可重建**：删除时改写 `orchestration_name` 追加删除时间戳。
8. **编排导入/导出纳入本期**：导出编排、步骤任务及完整依赖；导入两阶段检查，支持跳过/覆盖；步骤任务无法映射时该编排整体失败。
