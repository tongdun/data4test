#!/bin/bash

# 前端代码打包
yarn install
yarn build

# 后端代码打包
# Linux
GOOS=linux GOARCH=amd64 go build -o ./deploy/data4test main.go

# 部署打包
cp -rf hmtl deploy/
cp -rf mgmt deploy/
cp -f config.json deploy
mkdir -p deploy/web/static
cp -rf web/static deploy/web/static
rm -f deploy/log/*
rm -f deploy/api/*
rm -rf deploy/old/*
rm -rf deploy/history/*
rm -rf node_modules
rm `ls deploy/mgmt/data | grep -v "示例-用户管理-新建用户.yml"`
rm `ls deploy/mgmt/case | grep -v "project_V1.0.0_testcase_demo.xmind"`
rm `ls deploy/mgmt/common | grep -v "模板使用说明."`
rm `ls deploy/mgmt/upload | grep -v "_template."`

# Release打包
curDate=`date +"%Y%m%d"`
tar -cvf ./release/data4test_${curDate}.tgz deploy


# 后端各个架构的打包
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/data4test_linux_x86_64 main.go
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./release/data4test_linux_i386 main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/data4test_darwin_x86_64 main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./release/data4test_linux_aarch64 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/data4test_windows_x86_64.exe main.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./release/data4test_windows_i386.exe main.go


# 线上环境打包
tar -cvf data4test_20XX0X0X.tgz deploy --exclude=./deploy/mgmt/history/ --exclude=./deploy/mgmt/old/ --exclude=./deploy/mgmt/log/  --exclude=./deploy/download

# 变更包名称
mv data4test_XX_xXX_XX data4test


# 系统管理数据库手动批量新建
mkdir -p ./mgmt/{api,case,common,log,upload,history,old,download}

# 变更配置文件
vim config.json  # 变更数据库信息，路径信息按实际填写
vi config.json  # 部署系统未安装vim,使用vi

# 启动系统
nohup ./data4test  >/dev/null 2>&1 &   # 不写nohup.out日志
nohup ./data4test &   # 自动写nohup.out日志
