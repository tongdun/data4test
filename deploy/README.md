##### 环境要求
- linux系统/mac系统/windows系统
- mysql5.7数据库

##### 生产环境
- 下载release包
- 解压tar包：tar -xvf data4test_20XX0X0X.tgz
- 切换到文件目录：cd deploy
- 默认包提供的是linux x86版本，如有其他版本需要，请下载其他包进行替换即可
- 下载需要环境的执行文件到./deploy目录
- 变更包名：mv data4test_xxx data4test

##### 全新环境部署
- 创建数据库：create database data4test;
- 导入SQL文件：
  - 初始化SQL: mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/init_all_XXX.sql
- 更改配置文件：config.json, 所有占位符按实际情况填写
- vim config.json  # 变更数据库信息，路径信息按实际填写
- vi config.json   # 部署系统未安装vim,使用vi
- 监听端口默认9088，若有冲突，可根据实际情况变更


##### 已有环境更新
- 变更日志：./mgmt/doc/file/update/change_log.md
- 若有定时任务，运行中的任务重新部署后会自动拉起
- 导入SQL文件:
  - 更新SQL: mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/update.sql (根据日期执行对应的SQL)
- kill已有进程：ps aux | grep -i data4test | head -n 1 | awk '{print $2}' | xargs kill -9

#### 部署系统
- nohup ./data4test  >/dev/null 2>&1 &   # 不写nohup.out日志
- nohup ./data4test &   # 自动写nohup.out日志

##### 访问登录
- 访问：http://10.0.X.X:9088
- 默认用户：admin/ admin

#### 参数变更
- 如需使用智能化测试，请前往[环境-系统参数]，对"aiRunEngine"参数，根据大模型的实际情况进行配置
- 如需执行非标准化的数据， 请前往[环境-系统参数]，对"scriptRunEngine"参数，根据执行引擎的实际情况进行配置

#### deploy包注解
```
.
├── README.md      // 部署文档总体说明
├── config.json    // 配置文件
├── data4test      // 二进制执行程序
├── html         
│   └── index.html
├── mgmt
│   ├── README.md    // 文档目录说明
│   ├── api
│   ├── case
│   │   └── project_V1.0.0_testcase_demo.xmind   // 用例Xmind模板
│   ├── common
│   │   ├── image
│   │   │   ├── arch.jpg
│   │   │   ├── 数据流程图.png
│   │   │   ├── 全局使用流程图.jpg
│   │   │   └── 系统数据关系图.jpg
│   │   ├── 模板使用说明.json
│   │   └── 模板使用说明.yml
│   ├── data
│   │   └── 示例-用户管理-新建用户.yml
│   ├── doc
│   │   └── file
│   │       ├── README.md      // 在线使用文档总体说明
│   │       ├── arch
│   │       │   ├── arch.md
│   │       │   └── arch_en.md
│   │       ├── design
│   │       │   ├── action_design.md
│   │       │   ├── action_design_en.md
│   │       │   ├── api_mgmt_design.md
│   │       │   ├── api_mgmt_design_en.md
│   │       │   ├── assert_design.md
│   │       │   ├── assert_design_en.md
│   │       │   ├── console_design.md
│   │       │   ├── console_design_en.md
│   │       │   ├── data_file_design.md
│   │       │   ├── data_file_design_en.md
│   │       │   ├── db_design.md
│   │       │   ├── fuzzing_design.md
│   │       │   ├── mock_design.md
│   │       │   ├── mock_design_en.md
│   │       │   ├── parameter_design.md
│   │       │   ├── parameter_design_en.md
│   │       │   ├── perf_design.md
│   │       │   ├── playbook_design.md
│   │       │   ├── playbook_design_en.md
│   │       │   ├── script_design.md
│   │       │   ├── script_design_en.md
│   │       │   ├── task_design.md
│   │       │   ├── task_design_en.md
│   │       │   ├── testcase_design.md
│   │       │   └── testcase_design_en.md
│   │       ├── development
│   │       │   ├── dev_env.md
│   │       │   ├── dev_env_en.md
│   │       │   ├── must_know.md
│   │       │   └── must_know_en.md
│   │       ├── function
│   │       │   ├── feature_introduction.md
│   │       │   ├── feature_introduction_en.md
│   │       │   ├── module_function.md
│   │       │   └── module_function_en.md
│   │       ├── nav.yml
│   │       ├── plan
│   │       │   ├── blue_print.md
│   │       │   ├── blue_print_en.md
│   │       │   ├── todo.md
│   │       │   └── todo_en.md
│   │       ├── question
│   │       │   ├── FAQ.md
│   │       │   └── FAQ_EN.md
│   │       └── update
│   │           ├── change_log.md
│   │           ├── change_log_en.md
│   │           └── release.md
│   ├── download
│   ├── history
│   ├── log
│   ├── old
│   ├── sql                      // SQL文件
│   │   ├── init.sql       // 初始化SQL
│   │   └── update.sql     // 更新SQL
│   └── upload                   // 三方数据示例模板
│       ├── create_template.json
│       ├── create_template.txt
│       ├── create_template.xml
│       ├── create_template.yml
│       ├── record_template.json
│       ├── record_template.txt
│       ├── record_template.xml
│       └── record_template.yml
└── web                               // 控制台前端静态资源文件
└── static
└── js
├── app.min.js
└── vendor.min.js

25 directories, 71 files
```

#### 配置文件config.json注解
```
{
	"custom_foot_html": "./html/tongdun.html",
	"footer_info": "./html/tongdun.html",
	"open_admin_api": true,
	"database": {      // 数据库信息，根据实际情况填写
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
	"server_port": 9088,     // 系统端口，如有冲突，可自行变更
	"log_level": "release",
	"cicd_host": "X.X.X.X:8088",   // CICD自动触发，无定制化，无需关注
	"swagger_path": "http://{host:port}/api/metadata/rest/docs?group=group1",
	"redirect_path": "{\"Administrator\":\"/admin/info/schedule\", \"Operator\":\"/admin/info/schedule\", \"ApiManage\":\"/admin/likePostman\",\"Download\":\"/admin/fm/common/list\"}"    // 初次进入页面定义，可根据高频使用，自行定义，默认为任务列表
}
```