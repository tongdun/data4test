#### Environment Requirements
- Linux system / macOS / Windows system
- MySQL 5.7 database

#### Production Environment Deployment
- Download the release package
- Unpack the tar file: `tar -xvf data4test_20XX0X0X.tgz`
- Navigate to the file directory: `cd deploy`
- By default, the package provides a Linux x86 version. If you need another version, download and replace the corresponding package.
- Download the necessary environment executables to the `./deploy` directory
- Rename the package: `mv data4test_xxx data4test`

#### Deploying in a New Environment
- Create the database: `create database data4test;`
- Import SQL files:
    - Initialization SQL: `mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/init.sql`
    - Update SQL: `mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/update.sql`
- Modify the configuration file: `config.json`, filling in all placeholders with actual values
- Use `vim config.json` to edit the database and path information (if `vim` is not installed, use `vi`)
- The default listening port is 9088. If there's a conflict, change it as needed.

#### Updating an Existing Environment
- Check the change log: `./mgmt/doc/file/update/change_log.md`
- If there are scheduled tasks, running tasks will be automatically restarted after redeployment
- Import SQL files:
    - Update SQL: `mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/update.sql` (execute the corresponding SQL based on the date)
- Kill existing processes: `ps aux | grep -i data4test | head -n 1 | awk '{print $2}' | xargs kill -9`

#### Deploying the System
- Run without logging to nohup.out: `nohup ./data4test >/dev/null 2>&1 &`
- Run with automatic logging to nohup.out: `nohup ./data4test &`

#### Access and Login
- Access: `http://10.0.X.X:9088`
- Default user: `admin/admin`

#### Deploy Package Annotations
```
.
├── README.md      // Overall deployment documentation
├── config.json    // Configuration file
├── data4test      // Binary executable
├── html         
│   └── index.html
├── mgmt
│   ├── README.md    // Documentation directory description
│   ├── api
│   ├── case
│   │   └── project_V1.0.0_testcase_demo.xmind   // Test case Xmind template
│   ├── common
│   │   ├── image
│   │   │   ├── arch.jpg
│   │   │   ├── 数据流程图.png
│   │   │   ├── 全局使用流程图.jpg
│   │   │   └── 系统数据关系图.jpg
│   │   ├── 模板使用说明.json
│   │   └── 模板使用说明.yml
│   ├── data
│   │   └── 示例-用户管理-新建用户.yml
│   ├── doc
│   │   └── file
│   │       ├── README.md      // Overall online documentation description
│   │       ├── ... (Various documentation files)
│   ├── download
│   ├── history
│   ├── log
│   ├── old
│   ├── sql                      // SQL files
│   │   ├── init.sql       // Initialization SQL
│   │   └── update.sql     // Update SQL
│   └── upload                   // Third-party data example templates
│       ├── ... (Various template files)
└── web                               // Frontend static resources for the console
└── static
└── js
    ├── app.min.js
    └── vendor.min.js

25 directories, 71 files
```

#### Configuration File `config.json` Annotations
```json
{
    "custom_foot_html": "./html/tongdun.html",
    "footer_info": "./html/tongdun.html",
    "open_admin_api": true,
    "database": {      // Database information, fill in according to actual values
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
    "server_port": 9088,     // System port, can be changed if there's a conflict
    "log_level": "debug",
    "cicd_host": "X.X.X.X:8088",   // CICD auto-trigger, no customization required
    "swagger_path": "http://{host:port}/api/metadata/rest/docs?group=group1",
    "redirect_path": "{\"Administrator\":\"/admin/info/schedule\", \"Operator\":\"/admin/info/schedule\", \"ApiManage\":\"/admin/likePostman\",\"Download\":\"/admin/fm/common/list\"}"    // Initial page definition, can be customized based on frequent usage
}
```