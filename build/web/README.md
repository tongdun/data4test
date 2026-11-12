## 环境要求
- Linux / macOS / Windows 系统
- MySQL 5.7 数据库

## 包说明

`data4test-web_YYYYMMDD.tgz` 为「单目录多二进制」发布包，解压后内含各平台二进制与跨平台启动脚本 `start.sh`，无需再单独下载、替换其它平台二进制。

```bash
tar -xzf data4test-web_YYYYMMDD.tgz
cd data4test-web
```

## 全新环境部署
1. 创建数据库：`create database data4test;`
2. 导入 SQL 文件（全新部署直接导入完整初始化脚本，无需再执行 update.sql）：
   - 初始化 SQL：`mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/init_all_XXXX.sql`（包内仅保留最新一份初始化 SQL）
3. 修改配置文件 `config.json`：数据库信息、路径信息等所有占位符按实际情况填写（`vim config.json`；未装 vim 用 `vi`）
4. 启动服务（默认监听端口 9088，如有冲突在 config.json 变更）：
   - `nohup ./data4test >/dev/null 2>&1 &`   # 直接运行 linux/amd64 二进制
   - `nohup ./start.sh >/dev/null 2>&1 &`    # 跨平台启动脚本，自动选对应二进制

## 已有环境更新
1. 变更日志：`./mgmt/doc/file/update/change_log.md`
2. 若有定时任务，运行中的任务重新部署后会自动拉起
3. 增量导入 SQL：打开 `./mgmt/sql/update.sql`（文件内以 `# 日期` 分块记录各次变更），从「上次更新的日期」之后开始，依次执行其后的所有 SQL 块：
   - `mysql -h x.x.x.x -u user -p data4test`（进入交互后粘贴对应 SQL 块），或
   - 将对应 SQL 块另存为临时文件后执行 `mysql -h x.x.x.x -u user -p data4test < tmp.sql`
4. kill 已有进程：`ps aux | grep -i data4test | head -n 1 | awk '{print $2}' | xargs kill -9`

## 访问登录
- 访问：`http://10.0.X.X:9088`
- 默认用户：`admin / admin`

## 参数变更
- 如需使用智能化测试，请前往【环境-系统参数】，对 `aiRunEngine` 参数根据大模型实际情况配置
- 如需执行非标准化的数据，请前往【环境-系统参数】，对 `scriptRunEngine` 参数根据执行引擎实际情况配置

## 多平台二进制

默认平台 linux/amd64 直接产出裸 `data4test`，其余平台产出 `data4test_<os>_<arch>[.exe]`；`start.sh` 检测当前平台运行对应二进制。

| 平台 | 二进制名 |
|---|---|
| linux/amd64 | `data4test` |
| linux/386 | `data4test_linux_i386` |
| linux/arm64 | `data4test_linux_aarch64` |
| darwin/amd64（Intel Mac） | `data4test_darwin_x86_64` |
| darwin/arm64（Apple Silicon） | `data4test_darwin_aarch64` |
| windows/amd64 | `data4test_windows_x86_64.exe` |
| windows/386 | `data4test_windows_i386.exe` |

> 仅打包指定平台时，其余平台二进制不会出现在包内，`start.sh` 会提示「binary not found」。单平台包直接运行 `data4test` 即可。

## 目录结构

```
data4test-web/
├── README.md                // 部署文档总体说明
├── config.json              // 配置文件（数据库、端口、语种等）
├── data4test                // Web 服务二进制（linux/amd64）
├── data4test_<os>_<arch>    // 其它平台二进制（仅多平台包）
├── start.sh                 // 跨平台启动脚本，自动选对应二进制
├── html                     // 控制台前端页面
│   ├── index.html
│   ├── case_import.html
│   └── schedule_import.html
├── mgmt
│   ├── sql                  // SQL 文件
│   │   ├── init_all_XXXX.sql   // 初始化 SQL（仅保留最新一份）
│   │   └── update.sql          // 更新 SQL
│   ├── case
│   │   └── project_V1.0.0_testcase_demo.xmind   // 用例 Xmind 模板
│   ├── common               // 模板使用说明 + 文档图片
│   │   ├── 模板使用说明.json
│   │   ├── 模板使用说明.yml
│   │   └── image
│   ├── data
│   │   └── 示例-用户管理-新建用户.yml
│   ├── doc                  // 在线使用文档（含 design / update / change_log）
│   ├── upload               // 三方数据示例模板（*_template.*）
│   └── ...                  // 运行数据目录（log / api / history / download / i18n 等）
└── web
    └── static
        └── js
            ├── app.min.js
            └── vendor.min.js
```

## 配置文件 config.json 注解

```
{
	"custom_foot_html": "./html/tongdun.html",
	"footer_info": "./html/tongdun.html",
	"open_admin_api": true,
	"database": {                    // 数据库信息，根据实际情况填写
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
	"server_port": 9088,             // 系统端口，如有冲突可自行变更
	"log_level": "release",
	"cicd_host": "X.X.X.X:8088",     // CICD 自动触发，无定制化无需关注
	"swagger_path": "http://{host:port}/api/metadata/rest/docs?group=group1",
	"redirect_path": "{\"Administrator\":\"/admin/info/schedule\", \"Operator\":\"/admin/info/schedule\", \"ApiManage\":\"/admin/likePostman\"}"   // 初次进入页面定义，可按高频使用自行定义，默认为任务列表
}
```
