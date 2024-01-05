#### 开发环境
#### 控制台前端
- 命令1：yarn install
- 命令2：yarn run  选择dev模式

#### 后端
##### 导入MySQL
- [InitSQL](../../../sql/init.sql)
- [UpdateSQL](../../sql/update.sql) 

##### 更新配置文件
- [配置文件](../../../../config.json)  各项信息根据实际情况填写

##### 代码开启服务
- 命令：go run main.go / sudo go run main.go
- 访问：http://127.0.0.1:9088

##### Docker开启服务
- 命令1：docker run -p 3306:3306 --name data4test -e MYSQL_ROOT_PASSWORD=data4test -d --restart always mysql:5.7
- 命令2: mysql -h 127.0.0.1 -P 3306 -u root -p data4test < ./mgmt/sql/init.sql
- 命令3：mysql -h 127.0.0.1 -P 3306 -u root -p data4test < ./mgmt/sql/update.sql
- 命令4：

##### 登录
- 默认用户：admin/ admin

##### 其他
- 产品列表若使用Xmind导入用例功能，环境需先安装xmind2case
- xmind2case安装：https://pypi.org/project/xmind2case/
- xmind文件使用比较老版本的Xmind软件保存，最新的xmind2case解析不出来