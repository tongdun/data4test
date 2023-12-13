#!/bin/bash

# 线上环境打包
tar -cvf data4test_20XX0X0X.tgz mgmt --exclude=mgmt/history/ --exclude=mgmt/old/ --exclude=mgmt/log/  --exclude=scene/common/

# 前端代码打包
yarn install
yarn run


# 前端包更新
rm -rf mgmt/static/  #删除原有的文件
cp -rf ./web/static mgmt/static         #归到当前环境mgmt包下

# 部署打包
mkdir deploy
cp -rf hmtl deploy/
cp -rf mgmt deploy/
rm -rf deploy/log/*

# 后端包更新-MacOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o data4test_mac main.go

# 后端包更新-MacOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o data4test_darwin_x86_64 main.go

# 后端包更新-Linux
GOOS=linux GOARCH=amd64 go build -o data4test_linux_x86_64 main.go

# 变更包名称
mv data4test_XX_xXX_XX data4test

# 变更配置文件
vim config.json  # 变更数据库信息，路径信息按实际填写
vi config.json  # 部署系统未安装vim,使用vi

# 启动系统
nohup ./data4test  >/dev/null 2>&1 &   # 不写nohup.out日志
nohup ./data4test &   # 自动写nohup.out日志

# 系统管理数据库手动批量新建
mkdir -p ./mgmt/{api,case,common,log,upload,history,old,download}