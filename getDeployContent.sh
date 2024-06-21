#!/bin/bash

# 前端代码打包
yarn install
yarn build

# 后端代码打包
# Linux
GOOS=linux GOARCH=amd64 go build -o ./deploy/data4test main.go

# 部署打包
cp -rf html deploy/
cp -rf mgmt deploy/
cp -f config.json deploy
mkdir -p deploy/web/static
cp -rf web/static deploy/web/
sudo rm -f deploy/mgmt/log/*
sudo rm -f deploy/mgmt/api/*
sudo rm -rf deploy/mgmt/old/*
sudo rm -rf deploy/mgmt/history/*
cd deploy/mgmt/data
sudo rm -rf `ls | grep -v "示例-用户管理-新建用户.yml"`
cd -
cd deploy/mgmt/case
sudo rm -rf `ls | grep -v "project_V1.0.0_testcase_demo.xmind"`
cd -
cd deploy/mgmt/common
sudo rm -rf `ls| grep -v "模板使用说明."`
cd -
cd deploy/mgmt/upload
sudo rm -rf `ls | grep -v "_template."`
cd -
cd deploy/mgmt/download
sudo rm -rf `ls`
cd -

mv deploy/mgmt/doc/file/arch/image deploy/mgmt/common

#linux版本
#sed -i 's/..\/..\/..\/upload/\/admin\/fm\/upload\/preview?path="/g' deploy/mgmt/doc/file/design/action_design.md
#sed -i 's/..\/..\/..\/upload/\/admin\/fm\/upload\/preview?path="/g' deploy/mgmt/doc/file/design/mock_design.md
#sed -i 's/<img src=".\/image\/arch.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/arch.jpg)/g' deploy/mgmt/doc/file/arch/arch.md
#sed -i 's/<img src=".\/image\/全局使用流程图.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/全局使用流程图.jpg)/g' deploy/mgmt/doc/file/arch/arch.md
#sed -i 's/<img src=".\/image\/数据流程图.png">/[链接](\/admin\/fm\/common\/preview?path=\/image\/数据流程图.png)/g' deploy/mgmt/doc/file/arch/arch.md
#sed -i 's/<img src=".\/image\/系统数据关系图.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/系统数据关系图.jpg)/g' deploy/mgmt/doc/file/arch/arch.md


#Mac版本
sed -i '' 's/..\/..\/..\/upload/\/admin\/fm\/upload\/preview?path=/g' deploy/mgmt/doc/file/design/action_design.md
sed -i '' 's/..\/..\/..\/upload/\/admin\/fm\/upload\/preview?path=/g' deploy/mgmt/doc/file/design/mock_design.md
sed -i '' 's/<img src=".\/image\/arch.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/arch.jpg)/g' deploy/mgmt/doc/file/arch/arch.md
sed -i '' 's/<img src=".\/image\/全局使用流程图.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/全局使用流程图.jpg)/g' deploy/mgmt/doc/file/arch/arch.md
sed -i '' 's/<img src=".\/image\/数据流程图.png">/[链接](\/admin\/fm\/common\/preview?path=\/image\/数据流程图.png)/g' deploy/mgmt/doc/file/arch/arch.md
sed -i '' 's/<img src=".\/image\/系统数据关系图.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/系统数据关系图.jpg)/g' deploy/mgmt/doc/file/arch/arch.md

sed -i '' 's/parameter_design.md/parameter_design/g' deploy/mgmt/doc/file/design/data_file_design.md
sed -i '' 's/action_design.md/action_design/g' deploy/mgmt/doc/file/design/data_file_design.md
sed -i '' 's/assert_design.md/assert_design/g' deploy/mgmt/doc/file/design/data_file_design.md
sed -i '' 's/script_design.md/script_design/g' deploy/mgmt/doc/file/design/data_file_design.md


echo '''
{
	"custom_foot_html": "./html/tongdun.html",
	"footer_info": "./html/tongdun.html",
	"open_admin_api": true,
	"database": {
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
	"server_port": 9088,
	"log_level": "debug",
	"cicd_host": "X.X.X.X:8088",
	"swagger_path": "http://{host:port}/api/metadata/rest/docs?group=group1"
}
''' > deploy/config.json