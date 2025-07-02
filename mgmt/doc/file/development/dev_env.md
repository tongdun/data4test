#### 开发环境
#### 控制台前端
- 命令1：yarn install
- 命令2：yarn run  选择dev模式

#### 后端
##### 导入MySQL
- 全量模式
  - [InitAllSQL](../../../sql/init_all_XXX.sql)
- 增量模式
  - [InitSQL](../../../sql/init.sql)
  - [UpdateSQL](../../sql/update.sql)

##### 更新配置文件
- [配置文件](../../../../config.json)  各项信息根据实际情况填写

##### 代码开启服务
- 命令：go run main.go / sudo go run main.go

##### 登录
- 默认访问：http://127.0.0.1:9088
- 默认用户：admin/ admin

#### 参数变更
- 如需使用智能化测试，请前往[环境-系统参数]，对"aiRunEngine"参数，根据大模型的实际情况进行配置
- 如需执行非标准化的数据， 请前往[环境-系统参数]，对"scriptRunEngine"参数，根据执行引擎的实际情况进行配置

##### 其他
- 产品列表若使用Xmind导入用例功能，环境需先安装xmind2case
- xmind2case安装：https://pypi.org/project/xmind2case/
- xmind文件使用比较老版本的Xmind软件保存，最新的xmind2case解析不出来