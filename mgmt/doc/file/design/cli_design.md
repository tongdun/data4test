#### 运行形态
- CLI（命令行）运行模式是 Data4Test 的离线运行形态，与 Web 形态（GoAdmin + Gin）互补
- 使用 SQLite 本地存储，无需 MySQL，即可离线执行数据/场景并查看报告
- 任务包由 Web 端「任务一键导出」生成，包含 task_info.yaml、playbook_info.yaml、数据文件、env_config.yaml、upload/ 与 i18n 语料

#### 启动方式
- 二进制执行 `./data4test run|report|help` 时进入 CLI 分支，不连 MySQL、不启动 GoAdmin
- 首次运行自动在本地创建 SQLite 库 cli.db 并建表（10 张核心表）

#### 子命令
- run：执行任务包中的任务
- report：列出最近报告 / 导出报告 HTML / 补生成多任务报告
- help：显示帮助

#### run 命令
- 用法：`./data4test run [选项]`
- 选项：
  - `--task 任务A,任务B`    只执行指定任务（默认执行全部）
  - `--pkg task/xxx.tgz`    显式指定任务包（默认取 task/ 下最新修改的包）
  - `--reconfig`            从任务包重新复制原始 env_config.yaml 覆盖本地，重置环境信息
  - `--param Name=值`       命令行注入占位符值，多个逗号分隔、可重复，如 `--param A=1,B=2`
  - `--language en`         本次运行语言，覆盖 config.json 默认语种
  - `--mock`                后台启动本地 Mock 服务（检测到 /mock/ 引用时自动启动）
- 执行流程：选包 → 解包 env_config.yaml/upload/i18n → 首次运行或 --reconfig 落本地 env_config.yaml → 读配置 → 注入 --param → 补填剩余占位符（交互式）→ 必填校验 → 复用 Web 导入逻辑落库 → 同步 env_config.yaml → 复制 upload → 按需启动 Mock → 按顺序执行任务 → 多任务补生成报告
- 退出码：0 全部通过，3 存在失败，1 执行异常

#### report 命令
- `./data4test report`                     列出最近报告
- `./data4test report <ID>`                导出指定报告为独立 HTML 并打开浏览器
- `./data4test report --multi --task 1,2`  用报告编号补生成多任务报告
- 选项：`--language en`  查看/生成报告的语种

#### 占位符注入
- 配置文件（env_config.yaml）中若有 {NameXXX} 占位符，执行前检测并提示补充
- 可用 `--param Name=值` 从命令行直接注入（多个逗号分隔、可重复），注入后自动回写、下次免填
- 回车跳过则保持原样并继续执行；填写后自动回写配置文件
- 可手动编辑 env_config.yaml，把 {NameXXX} 占位符改成实际值

#### env_config.yaml 结构
- product：产品环境（name / ip / protocol / auth / threading / apps / private_parameter 等）
- app：按应用维度的环境配置（ip / prepath / auth 等），app 级 ip / auth 为空时继承 product
- sys_parameter：系统参数（多语种、特征数据、级联数据等）
- assert_template：断言值模板
- mock：本地 Mock 服务配置（enable / port / base_url）

#### 鉴权
- CLI 执行直接使用 env_config.yaml 中的鉴权信息
- 或在数据中配置output2auth动作获取鉴权信息