##### 环境要求
- linux系统/mac系统/windows系统
- mysql数据库

##### 生产环境
- 下载release包
- 解压tar包：tar -xvf data4test_20XX0X0X.tgz
- 切换到文件目录：cd deploy
- 默认包提供的是linux x86版本，如有其他版本需要，请下载其他包进行替换即可
- 下载需要环境的执行文件到./deploy目录
- 变更包名：mv data4test_xxx data4test

##### 全新环境部署
- 创建数据库：create database data4test;
- 导入初始化数据库文件：mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/init.sql
- 更改配置文件：config.json, 所有占位符按实际情况填写
- vim config.json  # 变更数据库信息，路径信息按实际填写
- vi config.json  # 部署系统未安装vim,使用vi
- 监听端口默认9088，若有冲突，可根据实际情况变更


##### 已有环境更新
- 变更日志：./mgmt/doc/file/update/change_log.md
- 若有定时任务，在新的环境中无需执行，可以把运行中的任务暂停，运行中的任务重新部署后会自动拉起
- 导入最新变动的SQL: mysql -h x.x.x.x -u user -p data4test < ./mgmt/sql/update.sql (拉取对应日志的SQL进行更新)
- kill已有进程：ps aux | grep -i data4test | head -n 1 | awk '{print $2}' | xargs kill -9

#### 部署系统
- nohup ./data4test  >/dev/null 2>&1 &   # 不写nohup.out日志
- nohup ./data4test &   # 自动写nohup.out日志

##### 访问登录
- 访问：http://10.0.X.X:9088
- 默认用户：admin/ admin




