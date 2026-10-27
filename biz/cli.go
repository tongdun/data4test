package biz

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"data4test/models"
)

// cliDDLStatements 本地 SQLite 建表语句，由 mgmt/sql/init.sql 裁剪而来。
// 仅保留离线执行与报告所需的 10 张表，首次运行时自动建表。
var cliDDLStatements = []string{
	`CREATE TABLE IF NOT EXISTS schedule (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_name TEXT, task_mode TEXT, crontab TEXT, threading TEXT, task_type TEXT,
		data_number TEXT, data_list TEXT, scene_number TEXT, scene_list TEXT,
		product_list TEXT, task_status TEXT, time4week TEXT, time4day TEXT, week TEXT,
		remark TEXT, last_at TEXT, next_at TEXT, user_name TEXT,
		created_at TEXT, updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS playbook (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT, data_number TEXT, data_file_list TEXT, last_file TEXT,
		scene_type INTEGER, priority INTEGER, run_time INTEGER, result TEXT,
		fail_reason TEXT, remark TEXT, user_name TEXT, product TEXT,
		created_at TEXT, updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS scene_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT, api_id TEXT, app TEXT, file_name TEXT, file_type INTEGER,
		content TEXT, run_time INTEGER, result TEXT, fail_reason TEXT,
		remark TEXT, user_name TEXT, created_at TEXT, updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS product (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product TEXT, ip TEXT, protocol TEXT, threading TEXT, thread_number INTEGER,
		auth TEXT, testmode TEXT, apps TEXT, env_type INTEGER,
		private_app_prefix TEXT, created_at TEXT, private_parameter TEXT,
		remark TEXT, updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS env_config (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product TEXT, app TEXT, ip TEXT, protocol TEXT, prepath TEXT,
		threading TEXT, thread_number INTEGER, auth TEXT, testmode TEXT,
		swagger_path TEXT, remark TEXT, created_at TEXT, updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS sys_parameter (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT, value_list TEXT, remark TEXT, created_at TEXT, updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS assert_template (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT, value TEXT, remark TEXT, user_name TEXT,
		created_at TEXT, updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS scene_test_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT, name TEXT, data_file_list TEXT, last_file TEXT,
		scene_type INTEGER, result TEXT, fail_reason TEXT, product TEXT,
		user_name TEXT, env_type INTEGER, remark TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')), updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS scene_data_test_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT, name TEXT, api_id TEXT, app TEXT, content TEXT,
		result TEXT, fail_reason TEXT, env_type INTEGER, file_type INTEGER,
		product TEXT, user_name TEXT, remark TEXT,
		created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')), updated_at TEXT, deleted_at TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS dashboard (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		report_name TEXT, report_type TEXT, related_task_ids TEXT,
		related_products TEXT, related_apps TEXT, time_range_start TEXT,
		time_range_end TEXT, status TEXT, creator TEXT, remark TEXT,
		report_data TEXT, created_at TEXT, updated_at TEXT, deleted_at TEXT
	)`,
}

var cliParamRe = regexp.MustCompile(`\{([A-Za-z_][A-Za-z0-9_.\-]*)\}`)

// ExportReportHTMLFunc 由 main.go 注入（pages.ExportReportHTML），用于 report <ID> 导出离线 HTML。
// 通过函数注入而非直接 import pages，避免 biz → pages 循环依赖。
var ExportReportHTMLFunc func(reportId, outDir string) (string, error)

// RunCLI 离线 CLI 入口，返回进程退出码。
// args[0] 为子命令（run/report/help），其余为子命令参数。
func RunCLI(args []string) int {
	if len(args) == 0 {
		printCLIHelp()
		return 1
	}
	switch args[0] {
	case "help", "-h", "--help":
		printCLIHelp()
		return 0
	case "run", "report":
	default:
		printCLIHelp()
		return 1
	}

	if err := cliInitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %s: %v\n", T("cli.init_db_failed"), err)
		return 1
	}

	switch args[0] {
	case "run":
		return runCLI(args[1:])
	case "report":
		return reportCLI(args[1:])
	}
	return 1
}

// cliInitDB 打开（或创建）本地 SQLite 库并自动建表。
func cliInitDB() error {
	if err := os.MkdirAll(BASEPATH, 0755); err != nil {
		return err
	}
	dbPath := filepath.Join(BASEPATH, "cli.db")
	if err := models.InitSQLite(dbPath); err != nil {
		return err
	}
	for _, ddl := range cliDDLStatements {
		if err := models.Orm.Exec(ddl).Error; err != nil {
			return err
		}
	}
	return nil
}

// runCLI 执行 ./data4test run 的全部/指定任务。
func runCLI(args []string) int {
	opts := parseCliRunFlags(args)

	// 语种：--language 覆盖 config.json 默认（biz init 已按默认语种初始化）
	if opts.language != "" {
		SetLocale(opts.language)
	}

	// 选择任务包
	pkg, err := selectTaskPkg(opts.pkgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}

	// 解包 CLI 专属文件：env_config.yaml / upload/
	tmpDir, err := extractCliMeta(pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	// 安装任务包附带的实体语料到 mgmt/i18n 下并重载索引（须在导入/执行之前，供报告与显示翻译使用）
	installCliI18n(tmpDir)

	// 首次运行 / --reconfig：把包内 env_config.yaml 复制到当前目录，作为可编辑本地配置
	if err := seedCliConfigFromPackage(tmpDir, opts.reconfig); err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}

	// 读取 env_config.yaml（本地覆盖优先，否则用任务包内占位符版本作底）
	cfg, err := loadCliConfig(tmpDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}

	// 系统参数与断言模板均取自 env_config.yaml
	sysParam := cfg.SysParameter

	// 检测并交互填写 {Name} 占位符（替代原 --ip 覆盖与 ${ENV} 展开）
	if err := fillCliPlaceholders(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}

	// 必填校验
	if cfg.Product.Name == "" {
		fmt.Fprintf(os.Stderr, "[cli] %s\n", T("cli.product_required"))
		return 1
	}
	if cfg.Product.Ip == "" {
		fmt.Fprintf(os.Stderr, "[cli] %s\n", T("cli.ip_required"))
		return 1
	}
	if cfg.Product.Protocol == "" {
		cfg.Product.Protocol = "http"
	}

	// 复用 Web 导入逻辑：写入 schedule / playbook / scene_data 与数据文件
	importResult, err := ImportScheduleCheck(pkg, "cli")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}
	if err := ImportScheduleConfirm(importResult.ImportId, "overwrite", "cli"); err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}

	// 同步 env_config.yaml 到本地库
	if err := upsertCliConfig(cfg, sysParam); err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}

	// 复制 upload/*
	if err := copyCliUploads(tmpDir); err != nil {
		Logger.Warning("cli copy upload failed: %s", err)
	}

	// 本地 Mock 服务：--mock / mock.enable 显式开启，或检测到数据引用 /mock/ 时自动启动
	mockPort := cfg.Mock.Port
	if mockPort == 0 {
		mockPort = 9066
	}
	mockUsed := detectCliMockUsage()
	if opts.mock || cfg.Mock.Enable || mockUsed {
		if srv := startCliMock(mockPort); srv != nil {
			defer shutdownCliMock(srv)
		}
		if cfg.Mock.BaseURL != "" {
			replaceCliMockURLs(cfg.Mock.BaseURL, mockPort)
			fmt.Printf("[cli] "+T("cli.mock_started")+"\n", mockPort)
		} else if mockUsed {
			fmt.Printf("[cli] %s\n", T("cli.mock_no_base_url"))
		}
	}

	// 选择要执行的任务
	taskNames, err := selectCliTasks(importResult.Tasks, opts.taskFilter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %v\n", err)
		return 1
	}
	if len(taskNames) == 0 {
		fmt.Fprintf(os.Stderr, "[cli] %s\n", T("cli.no_task"))
		return 1
	}

	fmt.Printf("[cli] "+T("cli.pkg_summary")+"\n", pkg, len(importResult.Tasks), len(importResult.Playbooks), len(importResult.DataFiles))

	// 逐个执行
	hasFail := false
	var taskIds []string
	for i, name := range taskNames {
		id, sp, st, dp, dt, runErr := runCliTask(name, cfg.Product.Name)
		if runErr != nil {
			Logger.Warning("task %s run error: %s", name, runErr)
		}
		taskIds = append(taskIds, id)
		dispName := GetTaskLocalized(name, GetLocale()) // 显示名：英文模式翻译，中文模式返回原值
		failCnt := (st - sp) + (dt - dp)
		if failCnt > 0 {
			hasFail = true
			fmt.Printf(T("cli.task_failed")+"\n", i+1, len(taskNames), dispName, cfg.Product.Name, sp, st, dp, dt)
		} else {
			fmt.Printf(T("cli.task_done")+"\n", i+1, len(taskNames), dispName, cfg.Product.Name, sp, st, dp, dt, cliPassRate(sp+dp, st+dt))
		}
	}

	// 多任务报告（复用 Web 逻辑，尽力而为）
	if len(taskIds) > 1 {
		func() {
			defer func() {
				if r := recover(); r != nil {
					Logger.Warning("cli multi-task report skipped: %v", r)
				}
			}()
			if e := GenerateMultiTaskReports(strings.Join(taskIds, ","), cfg.Product.Name, "cli"); e != nil {
				Logger.Warning("cli multi-task report failed: %s", e)
			}
		}()
	}

	// 等待异步单任务报告落地（供 report 命令使用，超时不阻塞退出码）
	waitDashboards(taskIds, 30*time.Second)

	if hasFail {
		fmt.Printf("[cli] %s\n", T("cli.has_fail"))
		return 3
	}
	fmt.Printf("[cli] %s\n", T("cli.all_pass"))
	return 0
}

// parseCliRunFlags 解析 run 子命令参数。
func parseCliRunFlags(args []string) cliRunOptions {
	var opts cliRunOptions
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--task":
			if i+1 < len(args) {
				opts.taskFilter = args[i+1]
				i++
			}
		case "--pkg":
			if i+1 < len(args) {
				opts.pkgPath = args[i+1]
				i++
			}
		case "--language":
			if i+1 < len(args) {
				opts.language = args[i+1]
				i++
			}
		case "--mock":
			opts.mock = true
		case "--reconfig":
			opts.reconfig = true
		}
	}
	return opts
}

// selectTaskPkg 返回要执行的任务包路径：显式指定或 task/ 下最新修改的包。
func selectTaskPkg(explicit string) (string, error) {
	if explicit != "" {
		if _, err := os.Stat(explicit); err != nil {
			return "", fmt.Errorf(T("cli.pkg_not_found"))
		}
		return explicit, nil
	}
	dir := "task"
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf(T("cli.pkg_not_found"))
	}
	var newest string
	var newestTime time.Time
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".tgz" && ext != ".tar" && ext != ".gz" && ext != ".zip" {
			continue
		}
		if e.ModTime().After(newestTime) {
			newestTime = e.ModTime()
			newest = filepath.Join(dir, e.Name())
		}
	}
	if newest == "" {
		return "", fmt.Errorf(T("cli.pkg_not_found"))
	}
	return newest, nil
}

// extractCliMeta 解包任务包到临时目录，供后续读取 env_config.yaml / upload/ 与语料。
func extractCliMeta(tgz string) (tmpDir string, err error) {
	tmpDir, err = ioutil.TempDir("", "cli_meta_")
	if err != nil {
		return tmpDir, err
	}
	data, err := ioutil.ReadFile(tgz)
	if err != nil {
		return tmpDir, err
	}
	if len(data) >= 2 && data[0] == 0x50 && data[1] == 0x4b {
		err = extractZip(data, tmpDir)
	} else if len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		err = extractGzipTar(data, tmpDir)
	} else {
		err = extractRawTar(data, tmpDir)
	}
	if err != nil {
		return tmpDir, err
	}
	normalizeCliExtractDir(tmpDir)
	return tmpDir, nil
}

// normalizeCliExtractDir 归一化解包目录：部分导出包会在包根下再包一层目录
// （如「数据文件包_xxx/」），此时把该层目录内容上提到 tmpDir，保证后续按平铺路径读取
// env_config.yaml / i18n / upload。平铺结构或有多个顶层项时保持原样。
func normalizeCliExtractDir(tmpDir string) {
	entries, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return
	}
	var dirs []os.FileInfo
	hasFile := false
	for _, e := range entries {
		// 跳过隐藏/元数据文件（.DS_Store、AppleDouble ._* 等），它们不是包内容。
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if e.IsDir() {
			dirs = append(dirs, e)
		} else {
			hasFile = true
		}
	}
	if hasFile || len(dirs) != 1 {
		return
	}
	top := filepath.Join(tmpDir, dirs[0].Name())
	subs, err := ioutil.ReadDir(top)
	if err != nil {
		return
	}
	for _, s := range subs {
		_ = os.Rename(filepath.Join(top, s.Name()), filepath.Join(tmpDir, s.Name()))
	}
	_ = os.Remove(top)
}

// installCliI18n 把任务包 i18n/ 下附带的实体语料（i18n_task / i18n_playbook / i18n_data）
// 挪到 BASEPATH/i18n 下并重载实体翻译索引，供英文模式翻译业务名称。语料目录不存在则跳过。
func installCliI18n(extractDir string) {
	for _, name := range []string{"i18n_task", "i18n_playbook", "i18n_data"} {
		src := filepath.Join(extractDir, "i18n", name)
		if _, err := os.Stat(src); err != nil {
			continue // 包内没有该语料，跳过
		}
		dst := filepath.Join(BASEPATH, "i18n", name)
		_ = os.RemoveAll(dst) // 覆盖旧语料，保证幂等
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			continue
		}
		if err := os.Rename(src, dst); err != nil {
			// 跨设备 rename 失败（如 /tmp 与工作目录不同挂载点）则逐文件复制
			copyCliI18nDir(src, dst)
		}
	}
	ReloadEntityI18n()
}

// copyCliI18nDir 跨设备 rename 失败时的兜底：逐文件复制语料目录。
func copyCliI18nDir(src, dst string) {
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}
	_ = os.MkdirAll(dst, 0755)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		b, err := ioutil.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			continue
		}
		_ = ioutil.WriteFile(filepath.Join(dst, e.Name()), b, 0644)
	}
}

// loadCliConfig 读取 env_config.yaml：任务包内 env_config.yaml（可选，含完整 app 块与占位符）作底，
// 本地 env_config.yaml（可选）作为覆盖层叠加，本地仅覆盖显式填写的字段，不再整体替换包内配置。
// 两者都缺失时返回空配置，交由 runCLI 的必填校验报错。
func loadCliConfig(tmpDir string) (cfg CliConfig, err error) {
	pkgPath := filepath.Join(tmpDir, "env_config.yaml")
	if _, e := os.Stat(pkgPath); e == nil {
		if err = readCliConfigFile(pkgPath, &cfg); err != nil {
			return cfg, err
		}
	}

	localPath := "env_config.yaml"
	if _, e := os.Stat(localPath); e != nil {
		return cfg, nil
	}
	var local CliConfig
	if err = readCliConfigFile(localPath, &local); err != nil {
		return cfg, err
	}
	cfg = mergeCliConfig(cfg, local)
	return cfg, nil
}

// readCliConfigFile 读取并解析一个 env_config.yaml。
func readCliConfigFile(path string, cfg *CliConfig) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf(T("cli.config_load_failed"), err)
	}
	return nil
}

// mergeCliConfig 将本地 env_config.yaml 作为覆盖层叠加到任务包 env_config.yaml（作为底）。
// 规则：本地非空字符串 / 非零整数覆盖，map（app、sys_parameter、assert_template）按 key 合并。
func mergeCliConfig(base, local CliConfig) CliConfig {
	out := base

	if local.Product.Name != "" {
		out.Product.Name = local.Product.Name
	}
	if local.Product.Protocol != "" {
		out.Product.Protocol = local.Product.Protocol
	}
	if local.Product.Ip != "" {
		out.Product.Ip = local.Product.Ip
	}
	if local.Product.Auth != "" {
		out.Product.Auth = local.Product.Auth
	}
	if local.Product.EnvType != 0 {
		out.Product.EnvType = local.Product.EnvType
	}
	if local.Product.Threading != "" {
		out.Product.Threading = local.Product.Threading
	}
	if local.Product.ThreadNumber != 0 {
		out.Product.ThreadNumber = local.Product.ThreadNumber
	}
	if local.Product.Apps != "" {
		out.Product.Apps = local.Product.Apps
	}
	if local.Product.Testmode != "" {
		out.Product.Testmode = local.Product.Testmode
	}
	if local.Product.PrivateAppPrefix != "" {
		out.Product.PrivateAppPrefix = local.Product.PrivateAppPrefix
	}
	if local.Product.PrivateParameter != "" {
		out.Product.PrivateParameter = local.Product.PrivateParameter
	}
	if local.Product.Remark != "" {
		out.Product.Remark = local.Product.Remark
	}

	if out.App == nil {
		out.App = make(map[string]CliApp)
	}
	for k, la := range local.App {
		ba, ok := out.App[k]
		if !ok {
			out.App[k] = la
			continue
		}
		if la.Protocol != "" {
			ba.Protocol = la.Protocol
		}
		if la.Ip != "" {
			ba.Ip = la.Ip
		}
		if la.Prepath != "" {
			ba.Prepath = la.Prepath
		}
		if la.Auth != "" {
			ba.Auth = la.Auth
		}
		if la.Threading != "" {
			ba.Threading = la.Threading
		}
		if la.Testmode != "" {
			ba.Testmode = la.Testmode
		}
		if la.Swagger != "" {
			ba.Swagger = la.Swagger
		}
		if la.Remark != "" {
			ba.Remark = la.Remark
		}
		out.App[k] = ba
	}

	if out.SysParameter == nil {
		out.SysParameter = make(map[string]CliJSONField)
	}
	for k, v := range local.SysParameter {
		out.SysParameter[k] = v
	}

	if out.AssertTemplates == nil {
		out.AssertTemplates = make(map[string]CliJSONField)
	}
	for k, v := range local.AssertTemplates {
		out.AssertTemplates[k] = v
	}

	if local.Mock.Enable {
		out.Mock = local.Mock
	} else {
		if local.Mock.Port != 0 {
			out.Mock.Port = local.Mock.Port
		}
		if local.Mock.BaseURL != "" {
			out.Mock.BaseURL = local.Mock.BaseURL
		}
	}

	return out
}

// seedCliConfigFromPackage 在首次运行（本地无 env_config.yaml）或 --reconfig 时，
// 把任务包内的 env_config.yaml 复制到当前目录（与 data4test 二进制同级）。
// force 为 true（--reconfig）时直接覆盖，不备份。
func seedCliConfigFromPackage(tmpDir string, force bool) error {
	dst := "env_config.yaml"
	if !force {
		if _, err := os.Stat(dst); err == nil {
			return nil // 已存在，不覆盖
		}
	}
	src := filepath.Join(tmpDir, "env_config.yaml")
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return nil // 包内无 env_config.yaml，交由 loadCliConfig 报错
	}
	if err := ioutil.WriteFile(dst, data, 0644); err != nil {
		return err
	}
	if force {
		fmt.Printf("[cli] %s\n", T("cli.reconfig_done"))
	} else {
		fmt.Printf("[cli] %s\n", T("cli.seed_config"))
	}
	return nil
}

// collectParamNames 从合并后配置的所有字符串字段中收集 {Name} 占位符，按名去重排序。
func collectParamNames(cfg CliConfig) []string {
	seen := make(map[string]bool)
	add := func(s string) {
		for _, m := range cliParamRe.FindAllStringSubmatch(s, -1) {
			seen[m[1]] = true
		}
	}
	add(cfg.Product.Name)
	add(cfg.Product.Protocol)
	add(cfg.Product.Ip)
	add(string(cfg.Product.Auth))
	add(cfg.Product.Threading)
	add(cfg.Product.Apps)
	add(cfg.Product.Testmode)
	add(string(cfg.Product.PrivateAppPrefix))
	add(string(cfg.Product.PrivateParameter))
	add(cfg.Product.Remark)
	for _, a := range cfg.App {
		add(a.Protocol)
		add(a.Ip)
		add(a.Prepath)
		add(string(a.Auth))
		add(a.Threading)
		add(a.Testmode)
		add(a.Swagger)
		add(a.Remark)
	}
	for _, v := range cfg.SysParameter {
		add(string(v))
	}
	for _, v := range cfg.AssertTemplates {
		add(string(v))
	}
	add(cfg.Mock.BaseURL)

	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// stdinIsTTY 判断标准输入是否为终端（决定是否交互式询问变量）。
func stdinIsTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// fillCliPlaceholders 检测配置中的 {Name} 占位符并交互式引导填写。
// 无占位符或非交互环境直接返回；非空输入写回本地 env_config.yaml，回车跳过保留占位符。
// 随后保留原 expandCliConfigEnv 的 app 级 ip/auth 继承 product 语义。
func fillCliPlaceholders(cfg *CliConfig) error {
	names := collectParamNames(*cfg)
	if len(names) > 0 {
		if !stdinIsTTY() {
			fmt.Printf("[cli] %s\n", T("cli.param_non_tty"))
		} else {
			fmt.Println(T("cli.param_title"))
			for _, n := range names {
				fmt.Printf("{%s}\n", n)
			}
			fmt.Println(T("cli.param_hint"))

			reader := bufio.NewReader(os.Stdin)
			filled := make(map[string]string)
			var written []string
			for _, n := range names {
				fmt.Printf("  "+T("cli.param_input"), n)
				line, err := reader.ReadString('\n')
				if err != nil {
					break // 读取失败/EOF，停止询问
				}
				// 归一化非法 UTF-8 字节（粘贴/终端编码残留），避免脏字节进入配置破坏 JSON。
				line = strings.ToValidUTF8(strings.TrimSpace(line), "")
				if line == "" {
					continue // 回车跳过，保留占位符
				}
				filled[n] = line
				written = append(written, n)
			}
			if len(filled) > 0 {
				replaceCliPlaceholders(cfg, filled)
				if err := writeBackPlaceholders(filled); err != nil {
					Logger.Warning("write back env_config.yaml failed: %s", err)
				} else {
					fmt.Printf("[cli] %s\n", fmt.Sprintf(T("cli.param_written"), strings.Join(written, ", ")))
				}
			}
		}
	}
	inheritAppFromProduct(cfg)
	return nil
}

// inheritAppFromProduct 保留原 expandCliConfigEnv 语义：app 级 ip/auth 为空或仍为 {Name} 占位符时，
// 继承 product 的同名字段（与 GetEnvConfig 运行时 product 覆盖 app 的语义一致）。
func inheritAppFromProduct(cfg *CliConfig) {
	for k, a := range cfg.App {
		if a.Ip == "" || cliParamRe.MatchString(a.Ip) {
			a.Ip = cfg.Product.Ip
		}
		if a.Auth == "" || cliParamRe.MatchString(string(a.Auth)) {
			a.Auth = cfg.Product.Auth
		}
		cfg.App[k] = a
	}
}

// replaceCliPlaceholders 用 m 中的 name→value 替换 cfg 中所有 {name} 占位符（内存替换，供本次运行）。
func replaceCliPlaceholders(cfg *CliConfig, m map[string]string) {
	sub := func(s string) string {
		for name, val := range m {
			s = strings.ReplaceAll(s, "{"+name+"}", val)
		}
		return s
	}
	cfg.Product.Name = sub(cfg.Product.Name)
	cfg.Product.Protocol = sub(cfg.Product.Protocol)
	cfg.Product.Ip = sub(cfg.Product.Ip)
	cfg.Product.Auth = CliJSONField(sub(string(cfg.Product.Auth)))
	cfg.Product.Threading = sub(cfg.Product.Threading)
	cfg.Product.Apps = sub(cfg.Product.Apps)
	cfg.Product.Testmode = sub(cfg.Product.Testmode)
	cfg.Product.PrivateAppPrefix = CliJSONField(sub(string(cfg.Product.PrivateAppPrefix)))
	cfg.Product.PrivateParameter = CliJSONField(sub(string(cfg.Product.PrivateParameter)))
	cfg.Product.Remark = sub(cfg.Product.Remark)
	for k, a := range cfg.App {
		a.Protocol = sub(a.Protocol)
		a.Ip = sub(a.Ip)
		a.Prepath = sub(a.Prepath)
		a.Auth = CliJSONField(sub(string(a.Auth)))
		a.Threading = sub(a.Threading)
		a.Testmode = sub(a.Testmode)
		a.Swagger = sub(a.Swagger)
		a.Remark = sub(a.Remark)
		cfg.App[k] = a
	}
	for k, v := range cfg.SysParameter {
		cfg.SysParameter[k] = CliJSONField(sub(string(v)))
	}
	for k, v := range cfg.AssertTemplates {
		cfg.AssertTemplates[k] = CliJSONField(sub(string(v)))
	}
	cfg.Mock.BaseURL = sub(cfg.Mock.BaseURL)
}

// writeBackPlaceholders 把 m 中的 name→value 写回本地 env_config.yaml（字符串替换，保留注释与排版）。
// 带引号的占位符（"{name}"）用 JSON 编码值替换：JSON 转义（\" \\ \uXXXX 等）是 YAML 双引号转义的子集，
// 二者都合法；且 json.Marshal 会把非法 UTF-8 字节归一为 �，不会像 strconv.Quote 那样产生
// JSON 非法的 \xNN 转义（会破坏 private_parameter / auth 等 JSON 字段）。裸占位符兜底按原值替换。
func writeBackPlaceholders(m map[string]string) error {
	const path = "env_config.yaml"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	text := string(data)
	for name, val := range m {
		quoted, _ := json.Marshal(val)
		text = strings.ReplaceAll(text, `"`+"{"+name+`}"`, string(quoted))
		text = strings.ReplaceAll(text, "{"+name+"}", val)
	}
	return ioutil.WriteFile(path, []byte(text), 0644)
}

// upsertCliConfig 同步 env_config.yaml 到本地库（同名覆盖）。
func upsertCliConfig(cfg CliConfig, sysParam map[string]CliJSONField) error {
	cur := time.Now().Format(baseFormat)

	// product：同名覆盖
	if cfg.Product.Testmode == "" {
		cfg.Product.Testmode = "fuzzing"
	}
	p := Product{
		Name:             cfg.Product.Name,
		Ip:               cfg.Product.Ip,
		Protocol:         cfg.Product.Protocol,
		Threading:        cfg.Product.Threading,
		ThreadNumber:     cfg.Product.ThreadNumber,
		Auth:             string(cfg.Product.Auth),
		Apps:             cfg.Product.Apps,
		Testmode:         cfg.Product.Testmode,
		EnvType:          cfg.Product.EnvType,
		PrivateAppPrefix: string(cfg.Product.PrivateAppPrefix),
		PrivateParameter: string(cfg.Product.PrivateParameter),
		Remark:           cfg.Product.Remark,
	}
	var exist Product
	models.Orm.Table("product").Where("product = ?", p.Name).Find(&exist)
	if exist.Name == "" {
		models.Orm.Table("product").Create(&p)
	} else {
		models.Orm.Table("product").Where("product = ?", p.Name).Updates(map[string]interface{}{
			"ip": p.Ip, "protocol": p.Protocol, "threading": p.Threading,
			"thread_number": p.ThreadNumber, "auth": p.Auth, "apps": p.Apps,
			"testmode": p.Testmode, "env_type": p.EnvType,
			"private_app_prefix": p.PrivateAppPrefix, "private_parameter": p.PrivateParameter,
			"remark": p.Remark, "updated_at": cur,
		})
	}

	// env_config：按 app 同名覆盖
	for appName, a := range cfg.App {
		if a.Testmode == "" {
			a.Testmode = "fuzzing"
		}
		var e EnvConfig
		models.Orm.Table("env_config").Where("app = ?", appName).Find(&e)
		if e.App == "" {
			models.Orm.Table("env_config").Create(&EnvConfig{
				App: appName, Product: cfg.Product.Name, Ip: a.Ip, Protocol: a.Protocol,
				Prepath: a.Prepath, Auth: string(a.Auth), Threading: a.Threading,
				MaxThreadNum: cfg.Product.ThreadNumber, Testmode: a.Testmode,
				SwaggerPath: a.Swagger, Remark: a.Remark,
			})
		} else {
			models.Orm.Table("env_config").Where("app = ?", appName).Updates(map[string]interface{}{
				"product": cfg.Product.Name, "ip": a.Ip, "protocol": a.Protocol,
				"prepath": a.Prepath, "auth": string(a.Auth), "threading": a.Threading,
				"testmode": a.Testmode, "swagger_path": a.Swagger, "remark": a.Remark,
				"updated_at": cur,
			})
		}
	}

	// sys_parameter：按 name 同名覆盖
	for name, val := range sysParam {
		var sp SysParameter
		models.Orm.Table("sys_parameter").Where("name = ?", name).Find(&sp)
		if sp.Name == "" {
			models.Orm.Table("sys_parameter").Create(&SysParameter{Name: name, ValueList: string(val)})
		} else {
			models.Orm.Table("sys_parameter").Where("name = ?", name).Updates(map[string]interface{}{
				"value_list": string(val), "updated_at": cur,
			})
		}
	}

	// assert_template：取自 env_config.yaml，按 name 同名覆盖
	for name, val := range cfg.AssertTemplates {
		var cnt int
		models.Orm.Table("assert_template").Where("name = ?", name).Count(&cnt)
		if cnt == 0 {
			models.Orm.Table("assert_template").Exec("INSERT INTO assert_template (name, value, created_at, updated_at) VALUES (?, ?, ?, ?)",
				name, string(val), cur, cur)
		} else {
			models.Orm.Table("assert_template").Where("name = ?", name).Updates(map[string]interface{}{
				"value": string(val), "updated_at": cur,
			})
		}
	}

	return nil
}

// copyCliUploads 复制任务包内 upload/ 目录到 UploadBasePath。
func copyCliUploads(tmpDir string) error {
	src := filepath.Join(tmpDir, "upload")
	if _, err := os.Stat(src); err != nil {
		return err
	}
	if err := os.MkdirAll(UploadBasePath, 0755); err != nil {
		return err
	}
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		b, e := ioutil.ReadFile(filepath.Join(src, f.Name()))
		if e != nil {
			continue
		}
		_ = ioutil.WriteFile(filepath.Join(UploadBasePath, f.Name()), b, 0644)
	}
	return nil
}

// selectCliTasks 根据导入结果与 --task 过滤，返回待执行任务名（排序后）。
func selectCliTasks(tasks []KTask, filter string) ([]string, error) {
	all := make([]string, 0, len(tasks))
	for _, t := range tasks {
		if t.Name != "" {
			all = append(all, t.Name)
		}
	}
	sort.Strings(all)

	if filter == "" {
		return all, nil
	}

	want := make(map[string]bool)
	for _, n := range strings.Split(filter, ",") {
		n = strings.TrimSpace(n)
		if n != "" {
			want[n] = true
		}
	}
	var selected []string
	for _, n := range all {
		if want[n] {
			selected = append(selected, n)
		}
	}
	for n := range want {
		found := false
		for _, a := range all {
			if a == n {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf(T("cli.task_not_found"), n, strings.Join(all, ", "))
		}
	}
	return selected, nil
}

// runCliTask 执行单个任务，返回 schedule id 与场景/数据统计。
func runCliTask(taskName, product string) (id string, scenePass, sceneTotal, dataPass, dataTotal int, err error) {
	var db DbSchedule
	models.Orm.Table("schedule").Where("task_name = ?", taskName).Find(&db)
	if db.Id == "" {
		return "", 0, 0, 0, 0, fmt.Errorf(T("cli.task_not_found"), taskName, "")
	}

	now := time.Now()
	lastAt := now.Format(baseFormat)
	timeTag := now.Format("20060102150405")
	taskTag := fmt.Sprintf("%s_%s", db.Id, timeTag)

	models.Orm.Table("schedule").Where("id = ?", db.Id).Updates(map[string]interface{}{
		"product_list": product,
		"last_at":      lastAt,
	})

	err = RunOnceTask(db.Id, "cli")

	sceneTotal, scenePass, _ = querySceneStatsByTaskId(taskTag)
	dataTotal, dataPass, _ = queryDataStatsByTaskId(taskTag)
	return db.Id, scenePass, sceneTotal, dataPass, dataTotal, err
}

// waitDashboards 等待给定任务的异步报告状态离开 generating（超时即返回）。
func waitDashboards(taskIds []string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		allDone := true
		for _, id := range taskIds {
			var cnt int
			models.Orm.Table("dashboard").
				Where("related_task_ids like ? and status = ?", id+"_%", "generating").
				Count(&cnt)
			if cnt > 0 {
				allDone = false
				break
			}
		}
		if allDone {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// cliPassRate 计算通过率字符串。
func cliPassRate(pass, total int) string {
	if total == 0 {
		return "0.0%"
	}
	return fmt.Sprintf("%.1f%%", float64(pass)/float64(total)*100)
}

// padDisplayLeft 按显示宽度左对齐补齐到 width 列（中文等宽字符按 2 列、ASCII 按 1 列）。
func padDisplayLeft(s string, width int) string {
	w := 0
	for _, r := range s {
		w += runeDisplayWidth(r)
	}
	if pad := width - w; pad > 0 {
		return s + strings.Repeat(" ", pad)
	}
	return s
}

// reportCLI 无参时列出最近报告；`report <ID>` 导出指定报告为独立 HTML 并在浏览器中打开；
// `report --multi --task 任务A,任务B` 基于本地库历史结果补生成多任务报告。
// 三种形式均支持 --language：查看/生成时按指定语种渲染，系统信息与报告内容同步翻译（与 Web 一致）。
func reportCLI(args []string) int {
	// 先统一解析参数（--language / --multi / --task / <ID> 位置任意）
	var language, taskFilter, reportId string
	var multi bool
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--language":
			if i+1 < len(args) {
				language = args[i+1]
				i++
			}
		case "--multi":
			multi = true
		case "--task":
			if i+1 < len(args) {
				taskFilter = args[i+1]
				i++
			}
		default:
			if !strings.HasPrefix(args[i], "-") && reportId == "" {
				reportId = args[i]
			}
		}
	}

	// 语种：--language 覆盖 config.json 默认，报告系统信息(biz.T)与内容(实体语料)同步翻译
	if language != "" {
		SetLocale(language)
	}

	if reportId != "" {
		return exportReportHTML(reportId)
	}
	if multi {
		return genMultiTaskReportCLI(taskFilter)
	}

	var list []DashboardReport
	models.Orm.Table("dashboard").Order("id desc").Limit(20).Find(&list)
	if len(list) == 0 {
		fmt.Printf("[cli] %s\n", T("cli.report_none"))
		return 0
	}
	fmt.Printf("[cli] %s\n", T("cli.report_list_hint"))
	fmt.Printf("  %s  %s  %s  %s\n",
		padDisplayLeft(T("cli.report_col_id"), 4),
		padDisplayLeft(T("cli.report_col_name"), 64),
		padDisplayLeft(T("cli.report_col_product"), 16),
		T("cli.report_col_time"))
	for _, d := range list {
		fmt.Printf("  %s  %s  %s  %s\n",
			padDisplayLeft(d.Id, 4),
			padDisplayLeft(d.ReportName, 64),
			padDisplayLeft(d.RelatedProducts, 16),
			d.CreatedAt)
	}
	return 0
}

// genMultiTaskReportCLI 解析报告编号列表为 schedule id，去重后复用 Web 逻辑补生成多任务报告。
// 报告编号来自 `report` 列表第一列（dashboard.id）；同一任务的多个报告不可重复选择。
func genMultiTaskReportCLI(reportIds string) int {
	if reportIds == "" {
		fmt.Fprintf(os.Stderr, "[cli] %s\n", T("cli.multi_task_required"))
		return 1
	}

	seen := make(map[string]bool)
	var ids []string
	for _, n := range strings.Split(reportIds, ",") {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		var d DashboardReport
		models.Orm.Table("dashboard").Where("id = ?", n).Find(&d)
		if d.Id == "" {
			fmt.Fprintf(os.Stderr, "[cli] %s\n", fmt.Sprintf(T("cli.report_id_not_found"), n))
			return 1
		}
		// 多任务报告（related_task_ids 含逗号，如 "2,1"）已合并了多个任务，
		// 不能再作为"单个任务"参与合并，否则会和其中已包含的任务重复。
		if strings.Contains(d.RelatedTaskIds, ",") {
			fmt.Fprintf(os.Stderr, "[cli] %s\n", fmt.Sprintf(T("cli.multi_task_is_multi"), n, d.RelatedTaskIds))
			return 1
		}
		base := baseTaskIds(d.RelatedTaskIds)
		if len(base) == 0 {
			fmt.Fprintf(os.Stderr, "[cli] %s\n", fmt.Sprintf(T("cli.report_id_not_found"), n))
			return 1
		}
		taskID := base[0]
		if seen[taskID] {
			fmt.Fprintf(os.Stderr, "[cli] %s\n", fmt.Sprintf(T("cli.multi_task_dup"), taskID))
			return 1
		}
		seen[taskID] = true
		ids = append(ids, taskID)
	}
	if len(ids) < 2 {
		fmt.Fprintf(os.Stderr, "[cli] %s\n", T("cli.multi_task_min"))
		return 1
	}

	if err := GenerateMultiTaskReports(strings.Join(ids, ","), "", "cli"); err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %s: %v\n", T("cli.multi_task_failed"), err)
		return 1
	}
	fmt.Printf("[cli] %s\n", T("cli.multi_task_done"))
	return 0
}

// baseTaskIds 从 related_task_ids 提取基础任务 id（schedule.id）。
// 单任务报告形如 "1_20261022154438"（id_时间戳），多任务报告形如 "1,2"。
func baseTaskIds(relatedTaskIds string) []string {
	var ids []string
	for _, part := range strings.Split(relatedTaskIds, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if i := strings.Index(part, "_"); i > 0 {
			part = part[:i]
		}
		ids = append(ids, part)
	}
	return ids
}

// exportReportHTML 导出指定报告为独立 HTML 文件并打开浏览器。
// 同时启动本地「查看详情」服务（含 mock，合并到同一端口），阻塞到 Ctrl+C 才退出，
// 供报告里两个饼图的「查看详情」跳转到本地页面（与 Web 端交互一致）。
func exportReportHTML(reportId string) int {
	if ExportReportHTMLFunc == nil {
		fmt.Fprintf(os.Stderr, "[cli] %s\n", T("cli.report_export_failed"))
		return 1
	}

	srv := startCliView(CliViewPort)
	if srv != nil {
		defer shutdownCliMock(srv)
		fmt.Printf("[cli] "+T("cli.view_server_started")+"\n", CliViewBaseURL())
	}

	outDir := "report"
	file, err := ExportReportHTMLFunc(reportId, outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[cli] %s: %v\n", T("cli.report_export_failed"), err)
		return 1
	}
	fmt.Printf("[cli] %s\n", T("cli.report_html_written", file))
	openBrowser(file)

	waitForSignal() // 阻塞保活，Ctrl+C 退出（defer 关闭服务）
	return 0
}

// openBrowser 在系统默认浏览器中打开本地文件。
func openBrowser(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	_ = cmd.Start()
}

// printCLIHelp 输出离线 CLI 帮助文档（必填项与配置要求）。
func printCLIHelp() {
	fmt.Print(`data4test CLI 执行工具

用法:
  ./data4test run [选项]          执行任务包中的任务
  ./data4test report              列出最近生成的报告
  ./data4test report <ID>         导出指定报告为独立 HTML 并打开浏览器
  ./data4test report --multi --task 1,2          用报告编号补生成多任务报告（编号见 report 列表）
  ./data4test help                显示本帮助

run 选项:
  --task 任务A,任务B    只执行指定任务（默认执行全部）
  --pkg task/xxx.tgz    显式指定任务包（默认取 task/ 下最新修改的包）
  --reconfig            从任务包重新复制原始 env_config.yaml 覆盖本地，重置环境信息
  --language en        本次运行语言，覆盖 config.json 默认语种
  --mock                后台启动本地 Mock 服务（默认检测到 /mock/ 引用时自动启动）

report 选项:
  --language en        查看/补生成报告的语种，报告系统信息与内容同步翻译（默认取 config.json 的 language）

占位符说明:
  配置文件中如有{NameXXX} 占位符, 执行前检测并提示补充；
  回车跳过则保持原样并继续执行；
  填写后自动回写配置文件，下次免填；
  可手动编辑env_config.yaml, 把 {NameXXX} 占位符改成实际值。

退出码: 0 全部通过, 3 存在失败, 1 执行异常
`)
}
