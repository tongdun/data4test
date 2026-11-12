## Environment Requirements
- Linux / macOS / Windows
- MySQL 5.7 database

## Package Description

`data4test-web_YYYYMMDD.tgz` is a single-directory multi-binary release package. After extraction it contains binaries for all platforms plus a cross-platform launcher `start.sh`, so there is no need to download or replace binaries for other platforms separately.

```bash
tar -xzf data4test-web_YYYYMMDD.tgz
cd data4test-web
```

## Fresh Environment Deployment
1. Create the database: `create database data4test;`
2. Import the SQL file (for a fresh deployment, import the complete initialization script directly — no need to run `update.sql`):
   - Initialization SQL: `mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/init_all_XXXX.sql` (only the latest initialization SQL is kept in the package)
3. Modify `config.json`: fill in the database information, paths, and all other placeholders with actual values (`vim config.json`; use `vi` if `vim` is not installed)
4. Start the service (default listening port 9088; change it in `config.json` if there is a conflict):
   - `nohup ./data4test >/dev/null 2>&1 &`   # run the linux/amd64 binary directly
   - `nohup ./start.sh >/dev/null 2>&1 &`    # cross-platform launcher, auto-selects the matching binary

## Updating an Existing Environment
1. Change log: `./mgmt/doc/file/update/change_log.md`
2. If there are scheduled tasks, running tasks will be automatically restarted after redeployment
3. Incremental SQL import: open `./mgmt/sql/update.sql` (the file records each change in blocks marked by `# date`); starting from the date after your last update, execute all subsequent SQL blocks in order:
   - `mysql -h x.x.x.x -u user -p data4test` (paste the corresponding SQL blocks in the interactive shell), or
   - save the corresponding SQL blocks to a temp file and run `mysql -h x.x.x.x -u user -p data4test < tmp.sql`
4. Kill existing processes: `ps aux | grep -i data4test | head -n 1 | awk '{print $2}' | xargs kill -9`

## Access and Login
- Access: `http://10.0.X.X:9088`
- Default user: `admin / admin`

## Parameter Changes
- To use AI-powered testing, go to 【Environment - System Parameters】 and configure the `aiRunEngine` parameter according to your LLM setup
- To execute non-standardized data, go to 【Environment - System Parameters】 and configure the `scriptRunEngine` parameter according to your execution engine setup

## Multi-platform Binaries

The default platform linux/amd64 produces a bare `data4test`; other platforms produce `data4test_<os>_<arch>[.exe]`. `start.sh` detects the current platform and runs the matching binary.

| Platform | Binary |
|---|---|
| linux/amd64 | `data4test` |
| linux/386 | `data4test_linux_i386` |
| linux/arm64 | `data4test_linux_aarch64` |
| darwin/amd64 (Intel Mac) | `data4test_darwin_x86_64` |
| darwin/arm64 (Apple Silicon) | `data4test_darwin_aarch64` |
| windows/amd64 | `data4test_windows_x86_64.exe` |
| windows/386 | `data4test_windows_i386.exe` |

> When only specific platforms are built, binaries for the other platforms are absent, and `start.sh` reports "binary not found". For a single-platform package, run `data4test` directly.

## Directory Structure

```
data4test-web/
├── README.md                // Overall deployment documentation
├── config.json              // Configuration file (database, port, language, etc.)
├── data4test                // Web service binary (linux/amd64)
├── data4test_<os>_<arch>    // Binaries for other platforms (multi-platform package only)
├── start.sh                 // Cross-platform launcher, auto-selects the matching binary
├── html                     // Console frontend pages
│   ├── index.html
│   ├── case_import.html
│   └── schedule_import.html
├── mgmt
│   ├── sql                  // SQL files
│   │   ├── init_all_XXXX.sql   // Initialization SQL (only the latest is kept)
│   │   └── update.sql          // Update SQL (date-partitioned)
│   ├── case
│   │   └── project_V1.0.0_testcase_demo.xmind   // Test case Xmind template
│   ├── common               // Template usage instructions + doc images
│   │   ├── 模板使用说明.json
│   │   ├── 模板使用说明.yml
│   │   └── image
│   ├── data
│   │   └── 示例-用户管理-新建用户.yml
│   ├── doc                  // Online documentation (design / update / change_log)
│   ├── upload               // Third-party data example templates (*_template.*)
│   └── ...                  // Runtime data directories (log / api / history / download / i18n, etc.)
└── web
    └── static
        └── js
            ├── app.min.js
            └── vendor.min.js
```

## Configuration File `config.json` Annotations

```
{
	"custom_foot_html": "./html/tongdun.html",
	"footer_info": "./html/tongdun.html",
	"open_admin_api": true,
	"database": {                    // Database information, fill in according to actual values
		"default": {
			"host": "db",
			"port": "3306",
			"user": "root",
			"pwd": "password",
			"name": "data4test",
			"max_idle_con": 5,
			"max_open_con": 10,
			"driver": "mysql",
			"parse_time": true
		}
	},
	"app_id": "uPkhI73C0y3p",
	"language": "cn",
	"prefix": "admin",
	"theme": "sword",
	"store": {
		"path": "./mgmt/upload",
		"prefix": "uploads"
	},
	"title": "盾测-自动化",
	"logo": "盾测-自动化",
	"mini_logo": "盾测",
	"index": "/",
	"login_url": "/login",
	"debug": true,
	"sql_log": false,
	"env": "local",
	"info_log": "./mgmt/log/info.log",
	"error_log": "./mgmt/log/error.log",
	"access_log": "./mgmt/log/access.log",
	"session_life_time": 7200,
	"file_upload_engine": {
		"name": "local"
	},
	"login_title": "盾测-自动化",
	"login_logo": "盾测-自动化",
	"auth_user_table": "goadmin_users",
	"bootstrap_file_path": "./bootstrap.go",
	"go_mod_file_path": "./go.mod",
	"asset_root_path": "./public/",
	"file_base_path": "./mgmt",
	"server_port": 9088,             // System port, can be changed if there's a conflict
	"log_level": "release",
	"cicd_host": "X.X.X.X:8088",     // CICD auto-trigger, no customization required
	"swagger_path": "http://{host:port}/api/metadata/rest/docs?group=group1",
	"redirect_path": "{\"Administrator\":\"/admin/info/schedule\", \"Operator\":\"/admin/info/schedule\", \"ApiManage\":\"/admin/likePostman\"}"   // Initial page definition, can be customized based on frequent usage, defaults to the task list
}
```
