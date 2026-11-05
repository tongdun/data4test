# 测试计划管理设计

> 关联需求文档：[test_plan_requirement.md](./test_plan_requirement.md)

## 1. 背景与目标

当前平台已有"定时任务"（`schedule`）能力：一个任务可以关联多个数据文件或场景，按 cron 定时执行，并生成任务级报告。但一个完整的回归/版本验证往往需要**多个任务按一定顺序或并发组合执行**，例如：

- 登录冒烟任务 → 订单主流程任务 → 报表回归任务；
- 三条互不依赖的冒烟任务并发执行，缩短总耗时。

现状只能逐个任务手动运行，或在 cron 上分别配置，缺少"把任务组装成计划、统一触发、统一查看结果"的能力。

本方案目标：

- 新增**测试计划管理页面**，列表能力对齐定时任务（查询、运行、停止、复制、报告等）；
- 计划由**已有定时任务**按顺序组装成步骤，可配置串行/并发执行、失败即停/失败继续；
- 计划支持**手动执行、cron 定时触发、CLI/CI 触发**；
- 每次计划执行生成**计划级汇总报告**，可下钻到各任务原有报告；
- 支持计划及其完整依赖的**导入/导出**，用于跨环境迁移与配置复用；
- 不改变任务本身的执行逻辑与既有报告体系。

## 2. 已确认的需求口径

| 项目 | 结论 |
|---|---|
| 计划组成 | 基于已有定时任务（`schedule`）组装，任务可被多个计划复用 |
| 执行方式 | **计划级可配置**：串行 / 并发 |
| 触发方式 | **手动 + 定时（cron） + CLI/CI** |
| 失败策略 | **计划级可配置**：失败即停 / 失败继续，默认失败即停 |
| 产品配置 | **计划级覆盖任务级**：计划配置了产品则所有步骤按计划产品执行；未配置时沿用各任务自身产品。`data` 类型步骤只使用第一个产品（见第 9 节） |
| 鉴权刷新 | 本期只提供 `PrepareProductAuth` 空 stub，真实刷新逻辑待 `auth_refresh_design.md` 落地后接入（见 6.8） |
| 计划报告 | **计划级汇总报告**，并入「结果报告列表」（`report_type=plan`），可下钻任务报告 |
| 页面形态 | 新增独立管理页，列表功能同定时任务 |
| 导入/导出 | **纳入本期**；导出计划 + 完整依赖，导入采用检查后选择跳过/覆盖，步骤任务缺失时整项计划失败 |
| 任务改动 | 不改任务定义；仅抽取任务执行函数供计划复用 |

### 2.1 不在本期范围

- 计划中直接引用数据文件/场景（仍通过任务间接引用）；
- 步骤级产品覆盖、步骤级重试、步骤级超时；
- 产品鉴权刷新的真实实现（本期只提供空 stub，见 6.8）；
- 邮件/IM 通知、审批流。

## 3. 概念模型

```text
测试计划 test_plan
   │  1:N（有序）
   ├── 计划步骤 test_plan_step ──N:1── 定时任务 schedule
   │
   └── 计划执行记录 test_plan_run
          │  1:N
          └── 步骤执行记录 test_plan_step_run ──0:1── 任务报告 dashboard
```

- **计划**：调度与执行策略的载体（cron 配置、串行/并发、失败策略）。
- **步骤**：对任务的引用 + 顺序号；不复制任务配置，任务修改后计划下次执行即生效。
- **执行记录**：每次触发产生一条，保存策略快照、状态、耗时与汇总 JSON。
- **步骤执行记录**：记录每个步骤本次执行结果，并关联任务级 `dashboard` 报告用于下钻。

## 4. 产品设计

### 4.1 导航与页面

- 菜单位置：**任务 → 测试计划**（与"定时任务"同级，父菜单 id=203）。
- 路由：`/admin/info/test_plan`（GoAdmin 自动生成列表/新建/编辑/详情）。
- 计划报告：并入「结果报告列表」，统一使用 `/admin/report_dashboard?id=<dashboard_id>` 查看；`/admin/plan_dashboard?id=<plan_id>` 作为兼容入口，跳转该计划最近一次报告。

### 4.2 列表页

字段：

| 列 | 说明 |
|---|---|
| ID | 隐藏，可筛选 |
| 计划名称 | 模糊筛选；支持多语种展示（沿用任务名称的 i18n 机制） |
| 计划模式 | cron / once / day / week（同定时任务） |
| 调度配置 | cron 表达式或每周/每天时间，表头隐藏，详情展示 |
| 执行方式 | 串行 / 并发 |
| 失败策略 | 失败即停 / 失败继续 |
| 步骤数 | 计划内启用步骤数（展示 `N`） |
| 关联产品 | 计划级产品 `test_plan.product_list`；未配置时展示步骤任务的产品并集（只展示） |
| 状态 | 运行中 / 已停止 / 已完成 / 未开始（与任务状态语义一致） |
| 上次执行时间 | 可排序 |
| 下次执行时间 | 可排序 |
| 创建人 | 模糊筛选 |
| 创建时间 | 时间范围筛选 |
| 备注 | 隐藏，详情展示 |

筛选：计划名称、计划模式、状态、执行方式、失败策略、创建人、创建时间。

行操作：

- **运行**：与定时任务一致——`once` 立即执行一次；`cron/day/week` 注册调度并置为运行中。
- **立即执行**：不改变调度状态，立刻执行一次（`trigger_type=manual`）。
- **停止**：取消 cron 调度并置为已停止；若计划正在执行，请求停止本轮（见 6.6）。
- **复制**：复制计划字段 + 全部步骤，名称追加 `_副本`。
- **报告**：跳转 `/admin/plan_dashboard?id={{.Id}}`。
- **编辑 / 删除**：删除为软删除；执行中的计划禁止删除。

批量操作：批量运行、批量停止、批量复制、批量导出。

导入/导出入口：

- 行操作「导出」：导出当前计划及完整依赖；
- 批量操作「批量导出」：把选中计划及去重后的共享依赖打入一个归档包；
- 顶部按钮「导入」：进入 `/admin/test_plan_import`，执行上传检查、冲突选择和确认导入。

实现约束（沿用定时任务现有写法）：

- 导入页用 `eng.HTMLFile("GET", "/admin/test_plan_import", "./html/test_plan_import.html", nil)` 注册（参照 `main.go:153-155`），不用 `r.GET`；
- 行操作/顶部按钮用 `info.AddActionButton` / `info.AddButton`，批量导出从 `ctx.FormValue("ids")` 取参、单条导出从 `ctx.FormValue("id")` 取参（参照 `tables/schedule.go:196-254`）；
- 导出是 `action.Ajax` 返回可点击下载链接，不是文件流响应；导出文件必须落在 `mgmt/download` 下，filemanager 只放行路径含 `mgmt` 的下载（`main.go:232-237`）；
- action id 不要与现有 `schedule_export` / `schedule_batch_export` 重名。

在**定时任务列表**增加：

- 行操作/批量操作 **「加入计划」**：弹窗选择目标计划，把选中任务按顺序追加到计划末尾（已存在则跳过）。
- 该入口是"基于任务组装计划"的快捷方式，计划页内的步骤编排为主入口。

### 4.3 新建 / 编辑表单

| 字段 | 控件 | 规则 |
|---|---|---|
| 计划名称 | Text | 必填、唯一 |
| 计划模式 | SelectSingle | cron / once / day / week，联动显示 crontab / week / time4week / time4day（同定时任务） |
| cron 表达式 | Text | 模式为 cron 时必填，使用 `robfig/cron` 校验 |
| 每周 / 每天时间 | Select | 复用 `Get7No()` / `Get24No()` |
| 执行方式 | Radio | 串行（默认）/ 并发 |
| 失败策略 | Radio | 失败即停（默认）/ 失败继续 |
| 关联任务 | TextArea | 每行一个任务名称，行顺序即步骤顺序；保存时校验任务存在、去重、至少 1 个 |
| 产品配置 | Select（多选） | 选填；选项来自产品列表，存 `product_list`（逗号分隔）。**配置后覆盖所有步骤任务的产品**；为空时沿用各任务自身产品 |
| 备注 | TextArea | 选填 |
| 状态/上次/下次/创建人/时间 | 只读 | 由系统维护 |

保存时：

- 校验计划名称不重复；
- 校验任务名称全部存在于 `schedule` 且未删除；
- 事务内写 `test_plan` 与 `test_plan_step`（先删后插，保持行顺序为 `step_order`）；
- 计划执行中（存在 status=running 的 run）时禁止修改步骤，避免执行中配置漂移。

### 4.4 步骤编排

- P0：表单内 TextArea，每行 `任务名称`（可带行号展示），保存后写入 `test_plan_step`，`step_order` 从 1 递增。
- 详情页展示有序步骤表：序号、任务名称（链接到任务详情）、任务类型、关联产品、当前是否可用（任务是否被删除）。
- 运行前对步骤做**快照**写入 `test_plan_step_run`（记录 task_id、task_name、step_order），执行中修改计划不影响本轮。
- P1：独立"步骤编排"页面，支持拖拽排序、启用/禁用单个步骤。

### 4.5 详情页

- 基本信息：同表单字段（只读）。
- 步骤列表：有序展示。
- 最近执行记录：最近 10 次 `test_plan_run`（时间、触发方式、状态、耗时、成功/失败步骤数），点击进入对应报告。

### 4.6 执行与停止（产品行为）

- **串行 + 失败即停**：某步骤失败后，剩余步骤标记"跳过"，计划状态 `failed`。
- **串行 + 失败继续**：失败步骤记录后继续执行，最终状态仍为 `failed`（有失败），但全部步骤都执行过。
- **并发 + 失败即停**：已启动的步骤继续执行到结束；**尚未启动**的排队步骤取消并标记"跳过"（受并发上限约束，见 6.4）。
- **并发 + 失败继续**：所有步骤都执行。
- **停止**：仅阻止后续未启动步骤；已开始的任务执行到结束（与现有任务不支持中断执行中的任务一致），计划状态记为 `stopped`。
- 同一计划同一时间只允许一轮执行；重复触发时：
  - 页面手动：提示"计划正在执行"；
  - cron：跳过本次并记录日志；
  - CLI：返回错误码 400。

### 4.7 计划报告

入口：

- **「结果报告列表」**（`/admin/info/dashboard`）：报告类型新增「计划」，计划执行完成后自动出现在列表中；
- 计划列表行操作「报告」：跳转该计划最近一次执行的报告；
- 计划详情「最近执行记录」：选择具体某一次执行。

页面路由：复用 `/admin/report_dashboard?id=<dashboard_id>`，`GetDashboardByReportId` 增加 `case "plan"` 渲染计划报告；`/admin/plan_dashboard?id=<plan_id>` 作为兼容入口，直接跳转该计划最近一次执行的 `dashboard_id`。

页面内容：

1. **概览**：计划名称、触发方式、执行方式、失败策略、开始/结束时间、总耗时、步骤总数、成功/失败/跳过数、执行率。
2. **步骤明细**：序号、任务名称、状态（成功/失败/跳过/停止）、开始/结束/耗时、通过/失败用例数、错误信息、「查看任务报告」链接（`/admin/report_dashboard?id=<dashboard_id>`）。
3. **产品维度汇总**：产品、总数、通过、失败、通过率（按本轮实际使用的产品聚合）。
4. **历史趋势**（P1）：最近 10 次计划执行的通过率。

报告数据落库在 `dashboard.report_data`（`report_type=plan`），`test_plan_run.dashboard_id` 关联本轮的 `dashboard` 记录；计划报告与任务报告统一在「结果报告列表」中查询，下钻关系保持不变。

### 4.8 权限与 i18n

- 新增菜单后，在 GoAdmin「权限管理」中补充：
  - `GET /info/test_plan`（查询）、`/info/test_plan/edit`、`/info/test_plan/new`；
  - `POST /new/test_plan`、`/edit/test_plan`、`/delete/test_plan`、`/export/test_plan`；
  - 计划运行/停止/复制等 Ajax 操作跟随页面权限；
  - `/admin/plan_dashboard` 页面权限；
  - 在「角色管理」中分配给对应角色。
- 菜单与权限种子以 `INSERT INTO goadmin_menu / goadmin_permissions / goadmin_role_permissions` 形式**追加到 `mgmt/sql/update.sql`**，不修改 `init_all_*.sql`。现有种子位于 `init_all_20260916.sql:1143/1153/1183`；权限模板应参照 `test_case` 的 7 条（同文件 1153 区块 80-86 行），`case_statistics` 只有菜单行、没有权限行。新增 id 需避开现有最大值（menu ≥2765、permission ≥301），并显式分配给目标角色。
- 菜单英文翻译：`mgmt/i18n/i18n_menu/en-US.yaml` 增加 `测试计划: Test Plan`。
- 业务文案：`mgmt/i18n/i18n_system/zh-CN.json`、`en-US.json` 新增 `test_plan.*`、`error.plan_*`（见第 8 节）。

## 5. 数据模型设计

### 5.1 `test_plan`

```sql
CREATE TABLE `test_plan` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `plan_name` varchar(255) NOT NULL COMMENT '计划名称',
  `plan_mode` enum('cron','once','day','week') NOT NULL DEFAULT 'once' COMMENT '计划模式',
  `crontab` varchar(64) DEFAULT NULL COMMENT 'cron表达式',
  `week` varchar(64) DEFAULT NULL COMMENT '每周几',
  `time4week` varchar(64) DEFAULT NULL COMMENT '每周执行小时',
  `time4day` varchar(64) DEFAULT NULL COMMENT '每天执行小时',
  `execute_mode` enum('serial','concurrent') NOT NULL DEFAULT 'serial' COMMENT '执行方式',
  `fail_strategy` enum('stop','continue') NOT NULL DEFAULT 'stop' COMMENT '失败策略',
  `product_list` varchar(64) DEFAULT NULL COMMENT '计划级产品(逗号分隔)，非空时覆盖所有步骤任务的产品',
  `plan_status` enum('running','stopped','finished','not_started') NOT NULL DEFAULT 'not_started' COMMENT '计划状态',
  `remark` longtext COMMENT '备注',
  `last_at` datetime DEFAULT NULL COMMENT '上次执行时间',
  `next_at` datetime DEFAULT NULL COMMENT '下次执行时间',
  `user_name` varchar(100) DEFAULT NULL COMMENT '创建人',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plan_name` (`plan_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试计划';
```

`plan_name` 唯一约束叠加软删除会导致"删除后同名计划无法重建"（`schedule` 表目前同样存在该问题）。实现约定：删除时先把 `plan_name` 改写为 `<原名>_deleted_<yyyymmddHHMMSS>` 再置 `deleted_at`，既保留唯一键语义，又允许重建同名计划。

### 5.2 `test_plan_step`

```sql
CREATE TABLE `test_plan_step` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `plan_id` int unsigned NOT NULL COMMENT '计划ID',
  `task_id` int unsigned NOT NULL COMMENT '任务ID(schedule.id)',
  `task_name` varchar(255) NOT NULL COMMENT '任务名称快照',
  `step_order` int NOT NULL DEFAULT 0 COMMENT '步骤顺序(从1开始)',
  `enabled` enum('yes','no') NOT NULL DEFAULT 'yes' COMMENT '是否启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_plan_task` (`plan_id`,`task_id`),
  KEY `idx_plan_order` (`plan_id`,`step_order`),
  KEY `idx_task` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试计划步骤';
```

`task_name` 作为快照便于任务改名/删除后仍能展示；执行时以 `task_id` 为准。`enabled` 字段 P0 固定为 `yes`，启用/禁用能力随 P1 编排页一并实现。

### 5.3 `test_plan_run`

```sql
CREATE TABLE `test_plan_run` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `plan_id` int unsigned NOT NULL,
  `plan_name` varchar(255) NOT NULL COMMENT '计划名称快照',
  `run_tag` varchar(100) NOT NULL COMMENT 'pl_<planId>_<yyyymmddHHMMSS>',
  `trigger_type` enum('manual','cron','cli') NOT NULL DEFAULT 'manual' COMMENT '触发方式',
  `execute_mode` enum('serial','concurrent') NOT NULL COMMENT '执行方式快照',
  `fail_strategy` enum('stop','continue') NOT NULL COMMENT '失败策略快照',
  `product_list` varchar(64) DEFAULT NULL COMMENT '本轮实际使用的产品快照(逗号分隔)',
  `status` enum('running','finished','failed','stopped') NOT NULL DEFAULT 'running' COMMENT '本轮状态',
  `start_at` datetime DEFAULT NULL,
  `end_at` datetime DEFAULT NULL,
  `duration_seconds` int DEFAULT 0,
  `total_steps` int DEFAULT 0,
  `success_steps` int DEFAULT 0,
  `fail_steps` int DEFAULT 0,
  `skipped_steps` int DEFAULT 0,
  `dashboard_id` int DEFAULT NULL COMMENT '计划报告dashboard.id(report_type=plan)',
  `creator` varchar(100) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_plan_created` (`plan_id`,`created_at`),
  KEY `idx_run_tag` (`run_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试计划执行记录';
```

### 5.4 `test_plan_step_run`

```sql
CREATE TABLE `test_plan_step_run` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `plan_run_id` int unsigned NOT NULL,
  `plan_id` int unsigned NOT NULL,
  `step_id` int unsigned NOT NULL,
  `task_id` int unsigned NOT NULL,
  `task_name` varchar(255) NOT NULL,
  `step_order` int NOT NULL DEFAULT 0,
  `status` enum('not_started','running','success','fail','skipped','stopped') NOT NULL DEFAULT 'not_started',
  `dashboard_id` int DEFAULT NULL COMMENT '任务报告dashboard.id',
  `start_at` datetime DEFAULT NULL,
  `end_at` datetime DEFAULT NULL,
  `duration_seconds` int DEFAULT 0,
  `error_msg` text,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_run_order` (`plan_run_id`,`step_order`),
  KEY `idx_task_run` (`task_id`,`plan_run_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试计划步骤执行记录';
```

### 5.5 与现有表的关系

| 表 | 关系 |
|---|---|
| `schedule` | 计划步骤通过 `task_id` 引用；任务删除后计划运行前置校验失败 |
| `dashboard` | ① 每个步骤调用一次任务执行，产生一条任务报告，`test_plan_step_run.dashboard_id` 关联；② 每轮计划执行生成一条 `report_type=plan` 的计划报告，`test_plan_run.dashboard_id` 关联，并统一展示在「结果报告列表」 |
| `goadmin_menu` | 新增"测试计划"菜单 |
| `goadmin_permissions` / `goadmin_role_permissions` | 新增页面与 CRUD 权限并分配角色 |

`dashboard` 表变更（承载计划报告）：

```sql
ALTER TABLE dashboard
  ADD COLUMN `related_plan_id` int unsigned DEFAULT 0 COMMENT '关联计划ID' AFTER `related_task_ids`,
  ADD COLUMN `related_plan_run_id` int unsigned DEFAULT 0 COMMENT '关联计划执行记录ID' AFTER `related_plan_id`;

-- report_type 为 varchar(20)，新增枚举值 'plan' 无需改列类型，只需补充代码枚举与 i18n
```

迁移：DDL 与菜单/权限种子**只追加到 `mgmt/sql/update.sql` 末尾**；`mgmt/sql/init_all_*.sql` 已定型发布，不修改。

## 6. 执行引擎设计

### 6.1 代码组织

按项目约定成对新增：

- `biz/test_plan_define.go`：结构体、常量、`test_plan`/`step`/`run` 类型；
- `biz/test_plan.go`：计划 CRUD、步骤保存、复制、停止；
- `biz/test_plan_run.go`：执行引擎（串行/并发、失败策略、停止、互斥）；
- `biz/test_plan_report.go`：计划报告数据聚合；
- `tables/test_plan.go`：GoAdmin 列表/表单/详情；
- `pages/test_plan_dashboard.go`：计划报告渲染（供 `GetDashboardByReportId` 调用）与最近一次报告跳转。

### 6.2 核心类型与函数

```go
// test_plan_define.go
type DbTestPlan struct {
	Id           string `gorm:"column:id" json:"id"`
	PlanName     string `gorm:"column:plan_name" json:"plan_name"`
	PlanMode     string `gorm:"column:plan_mode" json:"plan_mode"`
	Crontab      string `gorm:"column:crontab" json:"crontab"`
	Week         string `gorm:"column:week" json:"week"`
	Time4week    string `gorm:"column:time4week" json:"time4week"`
	Time4day     string `gorm:"column:time4day" json:"time4day"`
	ExecuteMode  string `gorm:"column:execute_mode" json:"execute_mode"`
	FailStrategy string `gorm:"column:fail_strategy" json:"fail_strategy"`
	ProductList  string `gorm:"column:product_list" json:"product_list"`
	PlanStatus   string `gorm:"column:plan_status" json:"plan_status"`
	Remark       string `gorm:"column:remark" json:"remark"`
	LastAt       string `gorm:"column:last_at" json:"last_at"`
	NextAt       string `gorm:"column:next_at" json:"next_at"`
	UserName     string `gorm:"column:user_name" json:"user_name"`
	CreatedAt    string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    string `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    string `gorm:"column:deleted_at" json:"deleted_at"`
}

type DbTestPlanRun struct {
	Id              string `gorm:"column:id" json:"id"`
	PlanId          string `gorm:"column:plan_id" json:"plan_id"`
	PlanName        string `gorm:"column:plan_name" json:"plan_name"`
	RunTag          string `gorm:"column:run_tag" json:"run_tag"`
	TriggerType     string `gorm:"column:trigger_type" json:"trigger_type"`
	ExecuteMode     string `gorm:"column:execute_mode" json:"execute_mode"`
	FailStrategy    string `gorm:"column:fail_strategy" json:"fail_strategy"`
	ProductList     string `gorm:"column:product_list" json:"product_list"`
	Status          string `gorm:"column:status" json:"status"`
	StartAt         string `gorm:"column:start_at" json:"start_at"`
	EndAt           string `gorm:"column:end_at" json:"end_at"`
	DurationSeconds int    `gorm:"column:duration_seconds" json:"duration_seconds"`
	TotalSteps      int    `gorm:"column:total_steps" json:"total_steps"`
	SuccessSteps    int    `gorm:"column:success_steps" json:"success_steps"`
	FailSteps       int    `gorm:"column:fail_steps" json:"fail_steps"`
	SkippedSteps    int    `gorm:"column:skipped_steps" json:"skipped_steps"`
	DashboardId     string `gorm:"column:dashboard_id" json:"dashboard_id"`
	Creator         string `gorm:"column:creator" json:"creator"`
}

type TestPlanRunOptions struct {
	TriggerType string // manual/cron/cli
	Operator    string
}

// 核心函数
func RunPlan(planId string, opts TestPlanRunOptions) (runId string, err error)
func RunPlans(planIds, triggerType, operator string) ([]string, error) // 批量，异步启动
func StopPlan(ids string) error                                      // 取消调度 + 请求停止本轮
func CopyTestPlan(id, userName string) error
func SaveTestPlanSteps(planId string, taskNames []string) error
func GetPlanRunReport(planId, runId string) (*PlanReportData, error)
```

任务侧抽取（关键改造）：

```go
// biz/biz_task.go
type TaskRunOptions struct {
	RunTag           string // 为空时沿用现有规则：<taskId>_<lastAt 格式化后的 YYYYMMDDHHMMSS>
	ProductOverride  string // 计划级产品（逗号分隔），非空时覆盖任务自身 product_list
	SkipAuthRefresh  bool   // 计划统一刷新鉴权后，步骤内跳过
	UpdateTaskStatus bool   // 计划步骤不修改 schedule.task_status
}

func RunTaskOnce(taskId, userName string, opts TaskRunOptions) (TaskRunResult, error)

// 原 RunOnceTask 保留为兼容包装
func RunOnceTask(id, userName string) error {
	_, err := RunTaskOnce(id, userName, TaskRunOptions{UpdateTaskStatus: true})
	return err
}
```

`TaskRunResult` 返回：`TaskTag`、`DashboardId`、`Status`、`Total/SuccessCount/FailCount`、`StartAt/EndAt/DurationSeconds`、`Err`。

### 6.3 执行流程

```text
触发（页面/ cron / CLI）
   │
   ├─ 计划互斥：同一 plan 已有 run → 拒绝/跳过
   ├─ 读取计划 + 启用步骤（按 step_order）
   ├─ 前置校验：步骤非空；任务均存在；cron 合法（启动调度时）
   ├─ 解析本轮产品：plan.product_list 非空 → 计划产品；为空 → 步骤任务产品并集
   ├─ 创建 test_plan_run(status=running) 并快照本轮产品；步骤快照写入 test_plan_step_run
   ├─ PrepareProductAuth(本轮产品, "task")   // 本期为空 stub，见 6.8
   ├─ 按 execute_mode 执行
   │    ├─ serial：逐个 RunTaskOnce；失败即停时剩余标记 skipped
   │    └─ concurrent：固定并发上限，失败即停时取消排队步骤
   ├─ 汇总：更新 run 状态/计数/耗时，写入本轮报告（dashboard，report_type=plan）
   └─ 更新计划 last_at；一次计划置 finished，cron 计划保持 running
```

### 6.4 并发与失败策略

- 并发实现：`sync.WaitGroup` + 信号量（带缓冲 channel），并发上限取系统参数 `planMaxConcurrency`，默认 5。
- 失败即停（concurrent）：收到失败信号后不再向信号量投递新步骤；已在执行的步骤等待完成，排队步骤标记 `skipped`。
- 失败即停（serial）：当前步骤结束后立即跳出循环。
- 失败继续：所有步骤执行完后统一汇总，计划状态为 `failed`（只要存在失败步骤）。
- 步骤之间不共享变量；结果只写 `test_plan_step_run`，避免并发写竞争。

### 6.5 定时调度

- 新增独立调度器 `PlanCronHandle = NewCrontab()`（复用 `Crontab` 结构），与任务调度器隔离，避免 id 冲突。
- 新 job 类型：

```go
type RunPlanSchedule struct {
	Id       string
	UserName string
}

func (r *RunPlanSchedule) Run() {
	RunPlan(r.Id, TestPlanRunOptions{TriggerType: "cron", Operator: r.UserName})
	PlanCronHandle.UpdateTime(r.Id)
}
```

- `GetPlanCron(planId, mode)` 镜像 `GetTaskCron`（`biz/biz_task.go:247`）：返回 cron 表达式并置 `plan_status=running`。注意 `GetTaskCron` 内部用 `strconv.Atoi` 解析 id，`test_plan.id` 若不是数字字符串不能照抄。
- 启动恢复：`biz.Init()`（`biz/biz_task.go:32`，其中 `ServiceRestart()` 在 `:34`）后增加 `PlanServiceRestart()`，把 `plan_status=running` 且非 once 的计划重新注册。
- 每次调度执行完成后，`PlanCronHandle.UpdateTime` 回写 `last_at/next_at`（参照 `Crontab.UpdateTime`，`biz/biz_task.go:56`）。
- 注册完成后必须显式调用 `PlanCronHandle.Start()`（参照 `biz/biz_task.go:483` 的 `CronHandle.Start()`），否则调度不会生效。

### 6.6 停止与互斥

- 计划互斥：包级 `sync.Map` 用 `LoadOrStore` 占位（`planId -> *struct{}`），`loaded=true` 即视为已有执行，执行结束 `Delete`。**不使用 `sync.Mutex.TryLock`**——该 API 要求 Go 1.18+，而当前 `go.mod` 为 `go 1.17`，新工具链下会报 `requires go1.18 or later`。系统重启后内存状态自然释放，`plan_run.status=running` 的历史记录在启动时标记为 `stopped`（防止永远 running）。
- 停止运行中的计划：包级 `sync.Map` 保存 `planId -> *atomic.Bool` 停止标记；
  - serial：每步执行前检查，置位后剩余步骤 `skipped`；
  - concurrent：检查发生在投递新步骤前，排队步骤 `skipped`；
  - 已开始的任务不中断（现有任务执行无法安全中断），计划最终状态 `stopped`。
- `StopPlan(ids)`：
  - 非 once 计划：`PlanCronHandle.DelByID`；
  - `plan_status=stopped`；
  - 若存在执行中的 run，置停止标记。

### 6.7 复用任务执行与改造

- 计划步骤统一按**任务立即执行语义**运行（抽取后的 `RunTaskOnce`），不使用 `OneTask`（cron 专用路径），保证行为可预期。注意：管理页"运行"按钮入口实际是 `RunTask`（`tables/schedule.go:305`），只有 `once` 模式才落到 `RunOnceTask`，`cron/day/week` 只是注册调度。
- `RunTaskOnce` 内部：
  - `effectiveProductList := opts.ProductOverride`；为空时取 `taskDb.ProductList`。执行与任务报告（`CreateDashboardRecord`）均使用该有效产品，保证"计划级产品优先"；
  - `taskTag` 支持外部传入，格式 `<taskId>_<planRunTag>_<stepOrder>`，仍以 `<taskId>_` 开头，保证现有"任务报告按 related_task_ids like `<taskId>_%`"的查询不受影响；
  - 返回本次任务报告 `dashboard_id`，供计划步骤下钻；
  - `UpdateTaskStatus=false` 时步骤执行不改写 `schedule.task_status`，避免计划执行影响任务自身的调度状态。
- 现有 `RunOnceTask` 改为包装函数，行为保持不变，避免影响其他调用方。
- **产品覆盖对 `data` 类型步骤的约束**：`data` 链路 `RepeatRunDataFile → RunDataFile → GetEnvConfig` 按单产品精确等值查询（`biz/app.go:55-57`），传入逗号分隔的多产品会查不到环境、丢失 host/auth。因此 `data` 步骤只取有效产品的**第一个**；`scene` 步骤由 `RunPlaybookFromMgmt` 遍历全部产品（`biz/biz_playbook.go:103-118`）。
- `CreateDashboardRecord` 与 `GenerateTaskReport` 都从 `DbSchedule.ProductList` 读取产品（`biz/task_report.go:322`、`:34`、`:175`），因此实现时复制一份 `taskDb` 并改写其 `ProductList` 即可贯穿执行与报告，**无需修改这两个函数的签名**。

### 6.8 与鉴权刷新方案的协同

按 `auth_refresh_design.md` 的口径，**本期先落 stub，不等该方案**：

- 新增 `PrepareProductAuth(products []string, source string) error` 空实现（直接 `return nil`），作为计划与鉴权刷新方案之间的稳定接入点；`auth_refresh_design.md` 落地后只替换该函数体，计划侧不再改动；
- 计划执行开始时确定本轮产品：`plan.product_list` 非空时直接使用；为空时按步骤任务依赖产品求**并集去重**（可复用 `GetScheduleProductListByIds`，但它只读 `schedule.product_list`、不含场景自身产品，需评估是否补齐）；
- 步骤执行传入 `SkipAuthRefresh=true`，待真实刷新落地后启用；
- 真实刷新失败时整个计划本轮不执行任何步骤，计划记为 `failed` 并写入错误信息（stub 阶段不会触发）。

### 6.9 CLI / CI

复用现有 CLI 认证（`/cli/login` 获取 JWT，`biz.CliAuthMiddleware`，定义在 `biz/web_cli_middleware.go:13`）。HTTP 路由仍注册在 `main.go` 的 `startServer()` 内，与 `/cli/dataRun`（`main.go:2347`）、`/cli/playbookRun`（`main.go:2378`）并列；`biz/cli.go` 是离线命令行实现，不注册这些路由。

| 接口 | 方法 | 参数 | 返回 |
|---|---|---|---|
| `/cli/planRun` | POST | `plan_name` 或 `plan_id`；可选 `wait`（默认 false） | `{code:200,data:{run_id,plan_id,status}}` |
| `/cli/planRunStatus` | POST | `run_id` | `{code:200,data:{status,steps:[...],report_url}}` |
| `/cli/planStop` | POST | `plan_id` 或 `run_id` | `{code:200,msg}` |

- `planRun` 默认异步启动（`go RunPlan(...)`），立即返回 `run_id`；
- `wait=true` 时同步等待执行结束再返回汇总；
- 若计划已在执行，返回 `{code:400,msg:"计划正在执行"}`；
- 响应格式沿用 `{code,msg,data}`。

## 7. 报告设计

报告统一存储在 `dashboard.report_data`（`report_type=plan`），`test_plan_run.dashboard_id` 关联本轮记录，`test_plan_run` 不再重复存储报告 JSON。

JSON 结构：

```json
{
  "plan_id": "3",
  "plan_name": "版本回归计划",
  "run_tag": "pl_3_20261023103000",
  "trigger_type": "cron",
  "execute_mode": "serial",
  "fail_strategy": "stop",
  "product_list": "产品A,产品B",
  "summary": {
    "total": 4,
    "success": 3,
    "fail": 1,
    "skipped": 0,
    "duration_seconds": 1234,
    "pass_rate": 75.0
  },
  "steps": [
    {
      "step_order": 1,
      "task_id": "12",
      "task_name": "登录冒烟",
      "status": "success",
      "start_at": "2026-10-23 10:30:00",
      "end_at": "2026-10-23 10:30:30",
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

- 生成时机：所有步骤结束后同步聚合，写入一条 `dashboard` 记录：
  - `report_type = 'plan'`
  - `report_name = biz.T("dashboard.plan_report_name")` 填充，中文为 `计划[%s]执行报告_%s`，与现有 `dashboard.task_report_name`（`biz/task_report.go:321`）同机制，不硬编码中文
  - `related_products = 本轮实际使用的产品`
  - `related_task_ids = 步骤任务ID（去重，逗号分隔的纯 ID）`。单任务报告才是 `<taskId>_<yyyymmddHHMMSS>` 形式并参与 `like '<taskId>_%'` 查询；计划报告与多任务报告都是纯 ID 列表，**不参与**该 like 查询
  - `related_plan_id = plan.id`、`related_plan_run_id = test_plan_run.id`
  - `report_data = 上述 JSON`，`status = finished`
  - 回写 `test_plan_run.dashboard_id`。
- 结果报告列表：`tables/dashboard.go` 需改**三处**——列表类型展示 switch（`:31-42`）、筛选 `FieldFilterOptions`（`:43-48`）、详情页类型展示 switch（`:131-143`），新增 `plan → 计划`，并补 i18n key `schedule_report.report_type_plan`。
- 展示：`pages.GetDashboardByReportId`（`pages/schedule_dashboard.go:17`）增加 `case "plan"`，渲染计划报告；步骤行链接到任务报告 `/admin/report_dashboard?id=<dashboard_id>`。
- 离线导出：`pages/report_export.go:60-69` 的 `renderReportPanel` 目前只有 `case "task"`，其余类型返回 `report_type_unsupported`；计划报告会出现在 CLI `report` 列表（`biz/cli.go:1168-1184`），必须补 `case "plan"`。渲染函数还需识别 `offlineReportExport` 标记（`pages/report_export.go:26`），保证离线链接行为一致。
- GoAdmin 详情页：`tables/dashboard.go:150-153` 把 `report_data` 渲染成任务专用的 `/admin/schedule_report?report_id=`，计划报告需按报告类型分流到 `/admin/report_dashboard?id=`。
- 计划列表「报告」：查询该计划最近一条 `test_plan_run.dashboard_id` 并跳转；无执行记录时提示暂无报告。
- 历史执行选择：在计划报告页内提供最近 20 次执行的下拉切换（通过 `test_plan_run` 查询对应 `dashboard_id`）。
- 复用现有图表能力（Chart.js，`pages/schedule_dashboard.go` 已使用）展示步骤成功/失败分布与通过率趋势（P1）。

## 8. 接入点与文件清单

新增：

- `biz/test_plan_define.go`、`biz/test_plan.go`、`biz/test_plan_run.go`、`biz/test_plan_report.go`
- `biz/test_plan_transfer_define.go`、`biz/test_plan_transfer.go`（导入/导出 DTO、依赖闭包、检查与确认）
- `tables/test_plan.go`
- `pages/test_plan_dashboard.go`（计划报告渲染 + 最近一次报告跳转）
- `html/test_plan_import.html`（上传检查、冲突选择、导入结果）

修改：

- `tables/tables.go`：`"test_plan": GetTestPlanTable`
- `tables/dashboard.go`：报告类型展示/筛选新增 `plan → 计划`（列表 switch、筛选 options、详情 switch 三处）；详情页链接按报告类型分流
- `pages/schedule_dashboard.go`：`GetDashboardByReportId` 增加 `case "plan"` 分支
- `pages/report_export.go`：`renderReportPanel` 增加 `case "plan"`，并处理 `offlineReportExport` 链接
- `biz/task_report_define.go`：`DashboardReport` 增加 `RelatedPlanId` / `RelatedPlanRunId` 字段（否则写不进新列）
- `biz/cli.go`：离线 SQLite `dashboard` 建表语句（`:82-88`）增加两列，保证 CLI 本地库可读计划报告
- `main.go`：
  - `GET /admin/plan_dashboard`：按 plan_id 跳转最近一次计划报告
  - `GET /admin/test_plan_import`：用 `eng.HTMLFile` 注册计划导入页面
  - `POST /test_plan_import_check`、`POST /test_plan_import_confirm`、`POST /test_plan_import_status`
  - `POST /cli/planRun`、`POST /cli/planRunStatus`、`POST /cli/planStop`
- `biz/biz_task.go`：抽取 `RunTaskOnce`（支持 `ProductOverride`），`Init()` 增加 `PlanServiceRestart`
- `biz/task_define.go`：`TaskRunOptions` / `TaskRunResult`（或放入 `test_plan_define.go`）
- `tables/schedule.go`：新增「加入计划」行操作/批量操作
- `mgmt/sql/update.sql`：4 张表 DDL + `dashboard` 新增字段 + 菜单 + 权限种子（仅追加到该文件，不改动 `init_all_*.sql`）
- `mgmt/i18n/i18n_system/zh-CN.json`、`en-US.json`
- `mgmt/i18n/i18n_menu/en-US.yaml`
- 实施完成后更新 `mgmt/doc/file/update/change_log.md`

i18n key 建议：

```text
test_plan.title                        测试计划
test_plan.description                  将多个定时任务组装为可串行/并发执行的测试计划
test_plan.detail_title                 测试计划详情
test_plan.plan_name                    计划名称
test_plan.plan_mode                    计划模式
test_plan.execute_mode                 执行方式
test_plan.execute_mode_serial          串行
test_plan.execute_mode_concurrent      并发
test_plan.fail_strategy                失败策略
test_plan.fail_strategy_stop           失败即停
test_plan.fail_strategy_continue       失败继续
test_plan.product_list                 产品配置
test_plan.steps                        关联任务
test_plan.step_count                   步骤数
schedule_report.report_type_plan       计划
dashboard.plan_report_name             计划[%s]执行报告_%s
test_plan.btn_run                      运行
test_plan.btn_run_once                 立即执行
test_plan.btn_stop                     停止
test_plan.btn_copy                     复制
test_plan.btn_report                   报告
test_plan.btn_add_to_plan              加入计划
test_plan.msg_run_started              计划已开始执行
test_plan.msg_stopped                  计划已停止
error.plan_not_found                   计划[%s]不存在
error.plan_no_steps                    计划[%s]未配置任务
error.plan_task_missing                计划[%s]中的任务[%s]不存在
error.plan_running                     计划[%s]正在执行，请稍后再试
error.plan_save_steps                  保存计划步骤失败: %s
```

## 9. 边界与限制

- **执行中停止**：只能阻止未开始的步骤，已开始的任务执行到结束（不改造任务执行链路的前提下）。
- **并发共享资源**：并发步骤可能同时访问同一产品环境，token、测试数据可能相互影响；页面提示"仅建议对相互独立的任务使用并发"。
- **计划级产品优先**：计划配置了产品时，步骤实际执行与任务报告均使用计划产品；未配置时才使用任务自身产品。计划产品与任务产品不一致时，以计划为准。
- **`data` 步骤只取首个产品**：`data` 执行链路按单产品精确匹配取环境配置，计划配置多个产品时 `data` 步骤只用第一个；多产品覆盖实际只在 `scene` 步骤完整生效。
- **同名计划重建**：`plan_name` 唯一索引叠加软删除，删除时需把名称改写为 `<原名>_deleted_<时间戳>`，否则同名计划无法重建（`schedule` 表目前就有此限制）。
- **鉴权刷新未实装**：本期 `PrepareProductAuth` 为空 stub，计划级统一刷新不生效；多步骤共用同一产品时仍按各步骤原有逻辑处理。
- **任务删除**：计划引用的任务被删除后，下次执行前置校验失败，计划直接失败且不执行任何步骤；建议后续给任务删除增加"被计划引用"校验。
- **状态语义**：`plan_status=running` 表示"调度已启用"，不代表"正在执行"；实际执行状态看 `test_plan_run.status`，与任务当前语义保持一致。
- **计划修改**：执行中禁止修改步骤；基本信息修改不影响已启动的本轮（本轮使用快照）。
- **报告保留**：`test_plan_run` 随执行次数增长，P1 增加按计划保留最近 N 次（系统参数 `planRunHistoryKeep`）。
- **权限**：CLI 接口使用 JWT，不受页面 RBAC 限制；需在网关/CI 侧保管账号密码。

## 10. 验证方案

1. **编译**：`make build` 或 `CGO_ENABLED=0 go build .`。
2. **页面 CRUD**：新建计划（3 个步骤）、编辑、复制、删除；校验任务不存在/重复/空计划被拒绝。
3. **串行 + 失败即停**：第 2 步任务构造失败，确认第 3 步为 skipped，计划 failed，报告可下钻。
4. **串行 + 失败继续**：全部步骤执行，最终 failed，步骤状态齐全。
5. **并发 + 失败即停/继续**：设置 `planMaxConcurrency=2`、5 个步骤，确认并发上限生效、排队步骤在失败即停时被 skipped。
6. **定时触发**：cron 计划启动后 `plan_status=running`、`next_at` 正确；到达时间触发；重启服务后调度自动恢复。
7. **手动立即执行**：不改变 cron 调度状态，产生一条 `trigger_type=manual` 的执行记录。
8. **停止**：serial 计划在第 1 步执行时停止，确认后续步骤 skipped、状态 stopped；调度被取消。
9. **CLI**：`/cli/login` → `/cli/planRun` → 轮询 `/cli/planRunStatus`，状态与页面一致；重复触发返回"计划正在执行"。
10. **互斥**：同一计划并发触发两次，仅一轮执行。
11. **鉴权刷新 stub**：本期 `PrepareProductAuth` 为空实现，确认计划执行不因缺少该能力报错；真实刷新落地后再补"多步骤共用产品只登录一次"的验证。
12. **计划级产品覆盖**：计划配置产品 A、步骤任务配置产品 B，执行后请求发往 A，任务报告 `related_products` 为 A；清空计划产品后回退到 B。另验证 data 类型步骤只使用计划产品中的第一个且能取到环境配置。
13. **报告统一**：计划执行完成后，「结果报告列表」出现类型为"计划"的报告，可筛选、打开并下钻任务报告。
14. **导出完整性**：单条/批量导出后检查归档清单，确认计划、步骤任务、场景、数据、产品、配置及附件齐全，且关联不依赖源环境数据库 ID。
15. **导入检查**：上传归档后仅返回新增、冲突、缺失及映射结果，不写业务表；分别验证 skip、overwrite、cancel。
16. **导入原子性**：构造步骤任务缺失，确认该计划及其步骤均不落库；同包其他完整计划可成功导入。
17. **归档安全**：验证路径穿越、非法类型、超限文件数/大小及损坏归档均被拒绝，临时目录被清理。
18. **回归**：未使用计划的既有任务执行、报告、cron 行为不变。
19. **离线报告导出**：`data4test report` 命令能列出并导出 `report_type=plan` 的报告为离线 HTML，步骤链接可正常点击。
20. **同名计划重建**：删除计划后可用同一名称重新创建成功。
21. **互斥实现**：在 `go 1.17` 指令下编译通过，重复触发同一计划时第二个请求被拒绝。

## 11. 实施拆分

**阶段一：数据与页面（可独立验收）**

- 4 张表 DDL + 菜单/权限/i18n；
- `tables/test_plan.go` 列表/表单/详情 + 步骤保存；
- 复制/删除。

**阶段二：执行引擎**

- 抽取 `RunTaskOnce`；
- 串行执行、失败策略、计划/步骤执行记录；
- 手动运行/立即执行/停止/互斥。

**阶段三：调度与 CLI**

- `PlanCronHandle`、启动恢复、时间回写；
- `/cli/planRun`、`/cli/planRunStatus`、`/cli/planStop`；
- 并发执行与 `planMaxConcurrency`。

**阶段四：报告与增强**

- 计划报告页、下钻、产品汇总；
- `tables/dashboard.go` 三处类型展示 + `pages/report_export.go` 离线导出分支；
- 任务列表「加入计划」、趋势图、历史清理（P1）。

**阶段五：导入/导出**

- 计划及任务依赖闭包收集、归档生成；
- 上传检查、冲突分类与名称映射；
- 异步确认导入、**导入任务状态持久化与页面轮询**（现有 `/schedule_import_confirm` 为 fire-and-forget，无状态接口可复用）、失败回滚及归档安全校验。

## 12. 已确认口径与后续

2026-10-23 首次确认，2026-11-05 补充确认，以下口径**保持当前设计**，本期不再变更：

1. 同一任务在一个计划中只允许出现一次（`uk_plan_task` 唯一约束）。
2. **支持计划级产品覆盖**：`plan.product_list` 非空时覆盖所有步骤任务的产品；为空时回退到各任务自身产品。
3. 并发 + 失败即停时，仅取消**排队中**的步骤；已启动步骤执行到结束（与"停止"语义一致）。
4. **计划报告并入「结果报告列表」**：新增 `report_type=plan`，报告数据存 `dashboard.report_data`，通过 `related_plan_id` / `related_plan_run_id` 关联计划与执行记录。
5. **计划导入/导出纳入本期**：导出计划、步骤任务及完整依赖；导入采用两阶段检查并支持跳过/覆盖冲突；步骤任务无法导入或映射时，该计划整体失败。
6. 执行记录长期保留，本期不做自动清理；后续如需限制保留条数，用系统参数 `planRunHistoryKeep` 增量支持。
7. **`data` 类型步骤只取首个产品**：`ProductOverride` 仍存完整逗号分隔串，但 data 执行链路按单产品精确匹配，只用第一个；多产品完整生效仅限 scene 步骤。
8. **鉴权刷新本期只落 stub**：新增 `PrepareProductAuth` 空实现并保留调用点，真实逻辑待 `auth_refresh_design.md` 落地后替换函数体。
9. **同名计划可重建**：删除时改写 `plan_name` 为 `<原名>_deleted_<时间戳>`，避免唯一索引与软删除冲突。
10. **互斥不用 `TryLock`**：改用 `sync.Map.LoadOrStore` 占位，兼容当前 `go 1.17` 指令。
