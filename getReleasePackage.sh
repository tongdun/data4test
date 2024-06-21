#!/bin/bash

# 当你执行该脚本时，请先执行getDeployContent.sh脚本，获取deploy目录的内容
# 再执行该脚本进行打包和数据清理

# Release打包
curDate=`date +"%Y%m%d"`
tar -cvf ./release/data4test_${curDate}.tgz deploy

# 后端各个架构的打包
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/data4test_linux_x86_64_${curDate} main.go
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./release/data4test_linux_i386_${curDate} main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/data4test_darwin_x86_64_${curDate} main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./release/data4test_linux_aarch64_${curDate} main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/data4test_windows_x86_64_${curDate}.exe main.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./release/data4test_windows_i386_${curDate}.exe main.go


cd deploy
rm -rf `ls | grep -v "README"`
cd -
