# Data4Test (盾测)

### 前言
Data4Test(盾测) 可以快速实现接口的自动化测试和管理，支持丰富的数据生成，支持复杂场景用例编排，适用于功能，并发，异常，模糊，场景，长时间，国际化，大数据，性能等方面的测试工作。

### 背景
- 已有的测试工具无法在一个场景里支持多应用接口的调用和执行
- Postman, Jmeter等单机版的测试工具无法快速在开发，测试，实施等多个角色间进行测试数据共享
- 接口变更无感，知道有变动，但无法快速定位到变更的接口，靠人工对接不靠谱
- 决策引擎系统场景复杂，链路依赖达20+或更多前置数据，自动化用例维护困难，编写脚本成本过高，变更环境失败比率也较高
- 风控系统接口请求数据字段过多，少则20+，多则100+或更多，人工输入符合特征的数据，人工构造时间成本过高
- 统计类功能需要长时间的数据积累，需各个时间维度的测试数据，需不同频度的定时任务执行
- 已有测试工具测试数据变更环境回放困难，需要数据用例幂等执行，且更换环境能快速落地数据进行复现
- 实时，离线，批转流，外部数据等多方数据特征需保持一致，且数据值需关联上
- 低并发测试需要常态化，靠手工不可能，靠脚本实现和维护成本过高
- 被测系统支持国际化，支持多语种，需要多语种的测试数据，且已建设的用例能直接复用，减少建设成本
- 等等，多个原因促成了本系统的诞生

### 系统
#### 快速试用
- 下载docker-compose.yml
- 切换到下载目录
- 启动服务：docker-compose up -d
- 默认访问：http://127.0.0.1:9088
- 默认用户：admin / admin

#### 架构设计
- [架构图](./mgmt/doc/file/arch/arch.md)

#### 系统页面
- [控制台运行域](./image/控制台运行域页面.jpg)
- [控制台场景域](./image/控制台场景域页面.png)
- [控制台数据历史域](./image/控制台数据历史域页面.png)
- [管理域产品列表](./image/管理域产品页面功能介绍.png)
- [管理域数据列表](./image/管理域数据列表功能介绍.png)
- [管理域场景列表](./image/管理域场景列表功能介绍.png)
- [管理域定时任务](./image/管理域定时任务列表功能介绍.png)

#### 系统介绍
- [特性简介](./mgmt/doc/file/function/feature_introduction.md)
- [模块功能](./mgmt/doc/file/function/module_function.md)
- [场景组合](./mgmt/doc/file/design/scene_design.md)
- [数据用例](./mgmt/common/模板使用说明.yml)
- [参数应用](./mgmt/doc/file/design/parameter_design.md)
- [动作应用](./mgmt/doc/file/design/action_design.md)
- [断言应用](./mgmt/doc/file/design/assert_design.md)
- [Mock应用](./mgmt/doc/file/design/mock_design.md)
- [模糊测试](./mgmt/doc/file/design/relation_design.md)

#### 近长期规划
- [长期规划](./mgmt/doc/file/plan/blue_print.md)
- [近期规划](./mgmt/doc/file/plan/todo.md)
- [性能测试](./mgmt/doc/file/design/perf_design.md)  待增强，或集成k6或Jmeter或其他

#### 开发须知
- [变更须知](./mgmt/doc/file/development/must_know.md)
- [更新记录](./mgmt/doc/file/update/change_log.md)

#### 生产环境
[生产环境部署](./deploy/README.md)

#### 开发环境
[开发环境部署](./mgmt/doc/file/development/dev_env.md)

### 功能
#### 适用测试类型
- 功能测试: 正向路径功能测试，自定义或自动生成符合特征的测试数据
- 并发测试: 单接口多测试数据的并行执行，多接口的并行执行，以及场景维度的并行执行等
- 异常测试: 通过占位符，快速构造超长边界值，特殊字符等
- 模糊测试: 自动生成模糊数据，开启健壮性测试(功能待充分验证)
- 场景测试: 跨应用，多接口，多鉴权，多环境测试，同时支持实时，离线，外部数据等多方数据对齐
- 长时间测试: 定时任务，持续构造测试数据
- Mock测试: 构造指定特征的数据，当外部数据给被测系统使用
- 国际化测试：根据请求语种，自动生成对应语种的测试数据，同时支持多语种的数据定义和断言判断，无需编写多个数据用例
- 大数据测试：通过动作自动生成海量的测试数据，以及数据和场景支持执行次数控制，实现实时和离线的大数据量
- 性能测试：支持控制并发数，开展性能压测(功能待增强)

### 其他
#### 周边技术
- yaml文件语法：https://www.runoob.com/w3cnote/yaml-intro.html
- 正则编写：https://www.runoob.com/regexp/regexp-syntax.html

#### 微信公众号
公众号：Data4Test  
欢迎扫码关注盾测更新动态以及系统使用小技巧会陆续上线    
<img src="./image/Data4Test公众号.jpg" width=30% />  

#### 社区微信群  
欢迎扫码，邀请加入我们的开源社区微信群，进行沟通交流  
<img src="./image/微信社区交流群.jpg" width=30% />  
如需联系开发作者，请添加微信：liuhuocjx 或 扫一扫添加  
<img src="./image/沟通联系方式.jpg" width=30% />    
(申请时备注“data4test”字样)  



