# Data4Test (盾测)
- [中文文档](./README.md)
- [EnglishDoc](./README_EN.md)

### Preface
Data4Test (盾测) enables quick implementation of automated testing and management of interfaces, supporting rich data generation. It is suitable for various testing tasks such as functional, concurrent, exceptional, fuzzy, scenario-based, long-duration, internationalization, big data, and performance testing.

### Background
- Existing testing tools cannot support the invocation and execution of multiple application interfaces within a single playbook.
- Single-machine testing tools like Postman, Jmeter cannot facilitate quick testing data sharing among roles such as development, testing, and implementation.
- Interface changes are unnoticed; there is awareness of changes but an inability to quickly pinpoint the modified interfaces, relying on manual integration is unreliable.
- Decision engine system scenarios are complex with link dependencies reaching 20+ or more pre-existing data. Automated case maintenance is difficult, script writing costs are high, and environment change failure rates are also high.
- Risk control system interface request data fields are numerous, ranging from 20+ to 100+ or more. Manually inputting data matching features and manually constructing time costs are excessively high.
- Statistical functions require long-term data accumulation, testing data for various time dimensions, and different frequency scheduled task executions.
- Existing testing tools face difficulties in environment playback after testing data changes, requiring data case idempotent execution, and swift landing of data for reproduction after environment change.
- Real-time, offline, batch conversion flow, external data, and other multi-dimensional data characteristics need to remain consistent, and data values need to be correlated.
- Low-concurrency testing needs to be normalized; it is impossible manually, and implementation and maintenance costs are too high through scripts.
- The tested system supports internationalization, multiple languages need testing data, and existing cases should be directly reusable to reduce construction costs.
- And so on, multiple reasons led to the birth of this system.

### System
#### Quick Trial
- Download docker-compose.yml
- Switch to the download directory
- Start the service: docker-compose up -d
- Default access: http://127.0.0.1:9088
- Default user: admin / admin

#### Arch Design
- [Arch Diagram](./mgmt/doc/file/arch/arch.md)

####  System Pages
- [Console Operation Domain](./image/控制台运行域页面.jpg)
- [Console Playbook Domain](./image/控制台场景域页面.png)
- [Console Data History Domain](./image/控制台数据历史域页面.png)
- [Management Domain Product List](./image/管理域产品页面功能介绍.png)
- [Management Domain Data List](./image/管理域数据列表功能介绍.png)
- [Management Domain Scenario List](./image/管理域场景列表功能介绍.png)
- [Management Domain Scheduled Tasks](./image/管理域定时任务列表功能介绍.png)

#### System Introduction
- [Feature Introduction](./mgmt/doc/file/function/feature_introduction_en.md)
- [Module Functions](./mgmt/doc/file/function/module_function_en.md)
- [Playbook Combination](./mgmt/doc/file/design/scene_design_en.md)
- [Data Cases](./mgmt/common/模板使用说明.yml)
- [Parameter Application](./mgmt/doc/file/design/parameter_design_en.md)
- [Action Application](./mgmt/doc/file/design/action_design_en.md)
- [Assertion Application](./mgmt/doc/file/design/assert_design_en.md)
- [Mock Application](./mgmt/doc/file/design/mock_design_en.md)
- [Fuzzy Testing](./mgmt/doc/file/design/relation_design_en.md)

#### Near-Term and Long-Term Plans
- [Long-Term Plan](./mgmt/doc/file/plan/blue_print_en.md)
- [Near-Term Plan](./mgmt/doc/file/plan/todo_en.md)
- [Performance Testing](./mgmt/doc/file/design/perf_design_en.md)  To be enhanced or integrated with k6, Jmeter, or others

#### Development Notes
- [Change Notes](./mgmt/doc/file/development/must_know_en.md)
- [Update Records](./mgmt/doc/file/update/change_log_en.md)

#### Production Environment
- [Production Environment Deployment](./deploy/README_EN.md)

#### Development Environment
- [Development Environment Deployment](./mgmt/doc/file/development/dev_env_en.md)

### Features
#### Applicable Test Types
- Functional Testing: Positive path functional testing, custom or automatically generated test data that fits specific features.
- Concurrent Testing: Parallel execution of single-interface multiple test data, parallel execution of multiple interfaces, and parallel execution at the scenario level.
- Exception Testing: Quickly construct extremely long boundary values, special characters, etc., using placeholders.
- Fuzzy Testing: Automatically generate fuzzy data, enabling robustness testing (functionality to be fully verified).
- Scenario Testing: Cross-application, multi-interface, multi-authentication, multi-environment testing, supporting real-time, offline, external data, and aligning with data from multiple sources.
- Long-Duration Testing: Timed tasks for continuous test data construction.
- Mock Testing: Construct data with specific features for use by the tested system when external data is given.
- Internationalization Testing: Automatically generate test data corresponding to the requested language, supporting multiple language data definition and assertion judgment without the need to write multiple data cases.
- Big Data Testing: Automatically generate massive test data through actions, with control over the number of executions for data and scenarios, achieving real-time and offline large data volumes.
- Performance Testing: Supports controlling concurrency, conducting performance stress testing (functionality to be enhanced).


### Other
#### Peripheral Technologies
- YAML file syntax: https://www.runoob.com/w3cnote/yaml-intro.html
- Regular expression writing: https://www.runoob.com/regexp/regexp-syntax.html

#### WeChat Official Account
- Official Account: Data4Test  
Welcome to scan the code to follow Shield Test updates and tips on system usage, which will be released one after another.  
<img src="./image/Data4Test公众号.jpg" width=30% />

#### Community WeChat Group
Feel free to scan and join our open-source community WeChat group for communication and discussion.  
<img src="./image/微信社区交流群.jpg" width=30% />  
If you need to contact the development author, please add WeChat: liuhuocjx or scan the QR Code.   
<img src="./image/沟通联系方式.jpg" width=30% />  
(When applying, please mention the keyword "data4test")  