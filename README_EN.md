# Data4Test (盾测)

- [中文文档](./README.md)
- [EnglishDoc](./README_EN.md)

### Preface
- Data4Test (ShieldTest) is an automated testing platform designed to simplify complex system testing by enabling users to write structured test data or test instructions through declarative statements. It makes automated testing for complex systems easier and can be used for functional, concurrency, exception, fuzz, scenario, long-running, internationalization, big data, and performance testing.

### Background

#### Application Background
- Data4Test is an automated testing platform specifically designed to address the challenges of testing complex business systems. It has been deeply applied in decision engine systems and risk control business systems.
- Over 150+ test tasks managed, 1500+ automated scenario test cases, 5000+ automated data file test cases, 1 million+ automated data case executions, and 100+ users.
- Iterated for over 3 years, spanning dozens of applications and multiple products, supporting the testing efforts of the company's ToB product line, and implemented and applied on-site at multiple B-end customers to support their testing acceptance and daily iteration work.
- Roles involved include testers, developers, implementers, customers, and product managers. The system enables easy automated testing and access to rich test data.

#### Origin Background
- 1. Existing testing tools cannot quickly support the invocation and execution of multiple application interfaces in a single scenario.
- 2. Local testing tools like Postman and JMeter cannot quickly share test data among developers, testers, and implementers.
- 3. Interface changes are not immediately apparent, and it's difficult to quickly locate changed interfaces, relying on manual coordination is unreliable.
- 4. Decision engine systems have complex scenarios with over 20+ or more dependent pre-data, making automated test case maintenance difficult, scripting costly, and environmental change failure rates high.
- 5. Risk control system interface request data fields are numerous, with a minimum of 20+ and up to 100+ or more, making manual input of characteristic data and construction time-consuming and costly.
- 6. Statistical functions require long-term data accumulation, testing data across various time dimensions, and scheduled tasks of different frequencies.
- 7. Existing testing tools have difficulties replaying test data changes across environments, requiring idempotent execution of data cases and quick environment switching for data reproduction.
- 8. Real-time, offline, batch-to-stream, and external data features need to be consistent, and data values need to be correlated.
- 9. Low concurrency testing needs to be normalized, which is impossible with manual efforts and too costly with scripts.
- 10. The tested system supports internationalization and multiple languages, requiring multilingual test data, and existing test cases can be directly reused to reduce development costs.
- 11. Some interfaces have encryption or logic written in the frontend, requiring UI automation or other scripts for management and execution.
- And many more reasons contributed to the birth and continuous iteration of this system.

### System

#### Online Experience
- Access URL: http://101.35.3.10:9088/
- Default Username and Password: admin / admin (Recommended to create a personal user after the first login for better experience)
- (Note: Not for production use)

#### Quick Local Trial
- Download docker-compose.yml
- Switch to the download directory
- Start the service: docker-compose up -d
- Default access: http://127.0.0.1:9088
- Default Username: admin / admin

#### System Features
- 1. User-friendly: Visual interface console similar to Postman, easy to learn and use with no barriers.
- 2. Easy to Write: Standardized data files in YAML format, declarative statements, quick batch writing without scripting.
- 3. Flexible Expansion: Customizable execution engines for various scripts, scalable on demand.
- 4. Operation-friendly: Logs can be viewed online, and system or execution issues can be debugged online.
- 5. Documentation-friendly: User manuals and guides are online for easy access and use.
- 6. Cross-platform Compatible: Golang program, compilable into various cross-platform executables.

### Functions

#### Functional Features
- 1. Supports interface management, change tracking, and Swagger interface one-click import and specification checks.
- 2. Supports Postman-style visual writing and YAML file batch writing for test data.
- 3. Supports rich assertion types for standardized data files for result verification.
- 4. Supports rich built-in characteristic data auto-generation and assembly for standardized data files.
- 5. Supports quick verification of exported standardized data files (e.g., CSV/EXCEL/YAML/JSON).
- 6. Supports N-level nested variable substitution for JSON format input parameters in standardized data files.
- 7. Supports correlated input parameters in standardized data files, using output List variables as references.
- 8. Supports management and execution of non-standardized script data files (e.g., Python, Shell, JMeter, DOS, etc.), with extensible script execution engines.
- 9. Supports scenario orchestration of standardized data files and non-standardized script files for management and execution.
- 10. Supports multiple types of scenario control: serial interruption, serial comparison, ordinary concurrency, and concurrency comparison.
- 11. Supports various task management: custom, one-time, daily, and weekly.
- 12. Supports scenario and task execution and scheduling across multiple environments simultaneously.
- 13. Supports historical data replay, re-execution of historical data/scenarios, and scenario continuation testing.
- 14. Supports the generation and online invocation of various custom format third-party Mock data.
- 15. Supports the generation of controllable polymorphic test data for upstream and downstream interfaces or function calls.
- 16. Supports one-click viewing of non-automated interfaces.
- 17. And more...

#### Applicable Test Types
- 1. Functional Testing: Positive path functional testing, custom or automatically generated characteristic test data.
- 2. Concurrency Testing: Parallel execution of multiple test data for a single interface, parallel execution of multiple interfaces, and scenario-level parallel execution.
- 3. Exception Testing: Quickly construct ultra-long boundary values, special characters, etc., through placeholders.
- 4. Fuzz Testing: Automatically generate fuzz data to initiate robustness testing.
- 5. Scenario Testing: Complex scenario visual orchestration, supporting cross-application, multi-interface, multi-authentication, and multi-environment testing, while aligning real-time, offline, and external data.
- 6. Long-running Testing: Scheduled tasks for continuous test data generation.
- 7. Mock Testing: Construct data with specified characteristics for use by the tested system as external data.
- 8. Internationalization Testing: Automatically generate test data in corresponding languages based on request languages, supporting multilingual data definition and assertion judgments without writing multiple data cases.
- 9. Big Data Testing: Automatically generate massive test data through actions, with execution count control for data and scenarios, achieving real-time and offline big data volumes.
- 10. Performance Testing: Supports controlling concurrency and invoking JMeter scripts for performance stress testing.
- 11. UI Testing: Completes UI automation through scripts.
- 12. Other Automated Tasks: Test report generation, efficiency-enhancing scripts, etc.

#### Architecture Design
- [Architecture Diagram](./mgmt/doc/file/arch/arch.md)

#### System Pages
- [Console Operation Domain](./image/控制台运行域页面.jpg)
- [Console Scenario Domain](./image/控制台场景域页面.png)
- [Console Data History Domain](./image/控制台数据历史域页面.png)
- [Management Domain Product List](./image/管理域产品页面功能介绍.png)
- [Management Domain Data List](./image/管理域数据列表功能介绍.png)
- [Management Domain Scenario List](./image/管理域场景列表功能介绍.png)
- [Management Domain Scheduled Tasks](./image/管理域定时任务列表功能介绍.png)

#### System Introduction
- [Feature Introduction](./mgmt/doc/file/function/feature_introduction_en.md)
- [Module Functions](./mgmt/doc/file/function/module_function_en.md)
- [API Management](./mgmt/doc/file/design/api_mgmt_design_en.md)
- [Test Case Management](./mgmt/doc/file/design/testcase_design_en.md)
- [Scenario Management](./mgmt/doc/file/design/playbook_design_en.md)
- [Data File Management](./mgmt/doc/file/design/data_file_design_en.md)
- [Task Management](./mgmt/doc/file/design/task_design_en.md)
- [Mock Application](./mgmt/doc/file/design/mock_design_en.md)

#### Short-term and Long-term Plans
- [Long-term Plan](./mgmt/doc/file/plan/blue_print_en.md)
- [Short-term Plan](./mgmt/doc/file/plan/todo_en.md)
- [Performance Testing](./mgmt/doc/file/design/perf_design_en.md)

#### Development Guidelines
- [Change Notice](./mgmt/doc/file/development/must_know_en.md)
- [Update Log](./mgmt/doc/file/update/change_log_en.md)
- [Release Notes](./mgmt/doc/file/update/release_en.md)

#### Production Environment
[Production Environment Deployment](./deploy/README_EN.md)

#### Development Environment
[Development Environment Deployment](./mgmt/doc/file/development/dev_env_en.md)

### Other
#### Surrounding Technologies
- YAML File Syntax: https://www.runoob.com/w3cnote/yaml-intro.html
- Regular Expression Writing: https://www.runoob.com/regexp/regexp-syntax.html

#### WeChat Official Account
WeChat Official Account: Data4Test    
Welcome to scan the QR code and follow for the latest updates and system usage tips. More tips will be coming soon!  
<img src="./image/Data4Test公众号.jpg" width="30%" />

#### Community WeChat Group
Welcome to scan the QR code and join our open-source community WeChat group for communication and discussion.  
<img src="./image/微信社区交流群.jpg" width="30%" />  
If you need to contact the developer, please add WeChat: liuhuocjx or scan the QR code below.  
<img src="./image/沟通联系方式.jpg" width="30%" />  
(Please mention "data4test" when adding)