# data4test CLI 执行工具

## 快速开始

1. 将 Web 端导出的任务包（`*.tgz`）放入 `task/` 目录。
2. 首次执行时自动从任务包复制 `env_config.yaml` 到当前目录；如需覆盖现场差异字段，直接编辑该文件（必填项见下）。
3. 执行：`./data4test run`（多平台包可用 `./start.sh run` 自动选平台；若配置含 `{Name}` 占位符，会提示逐个填写）。
4. 查看报告：
   - `./data4test report` 列出最近报告
   - `./data4test report <ID>` 导出 HTML 并在浏览器打开（输出到 `report/`）

## 目录结构

```
data4test-cli/
├── data4test                # CLI 可执行文件（linux/amd64）
├── data4test_<os>_<arch>    # 其它平台二进制（仅多平台包）
├── start.sh                 # 跨平台启动脚本，自动选对应二进制
├── config.json              # 运行配置（语种等，勿删）
├── env_config.example.yaml  # env_config 示例（人工提供，仅参考）
├── env_config.yaml          # 生效配置（首次运行自动从任务包复制，可编辑）
├── task/                    # 任务包（*.tgz）
├── report/                  # 导出的 HTML 报告
└── mgmt/                    # 运行数据（cli.db、历史、上传文件等）
```

## 命令

| 命令 | 说明 |
|---|---|
| `./data4test run [选项]` | 执行任务包中的全部/指定任务 |
| `./data4test report` | 列出最近 20 条报告 |
| `./data4test report <ID>` | 导出该报告 HTML 并打开浏览器 |
| `./data4test report --multi --task <编号1>,<编号2>` | 用报告编号补生成多任务报告（编号见 report 列表的「编号」列，须为不同任务的单任务报告） |
| `./data4test help` | 查看帮助 |

### run 选项

| 选项 | 说明 |
|---|---|
| `--task 任务A,任务B` | 只执行指定任务（默认全部） |
| `--pkg task/xxx.tgz` | 显式指定任务包（默认取 task/ 下最新） |
| `--reconfig` | 从任务包重新复制原始 env_config.yaml 覆盖本地，重新填写变量 |
| `--language en` | 本次运行语言 |
| `--mock` | 后台启动本地 Mock 服务 |

## 变量占位符 {Name}

配置值中可用 `{Name}` 占位符（如 `{HostIP}`、`{Auth}`）。运行检测到占位符时会逐个提示填写（输入全回显），非空输入自动写回 `env_config.yaml`、下次免填；回车跳过则保留占位符、照常执行。

改变量有两种方式：

- 执行前手动编辑 `env_config.yaml`，把 `{Name}` 改成字面值；
- 运行时按提示交互输入（每次检测到占位符都会提示，直到不再残留占位符）。

`--reconfig` 会从任务包重新复制原始 `env_config.yaml` 覆盖本地，重新触发填写。

## 配置要求

- `env_config.yaml` 与 `data4test` 二进制同级，YAML 格式；首次运行自动从任务包复制，或用 `--reconfig` 重拷。
- 字段名与 `product` / `env_config` / `sys_parameter` 表保持一致；单产品只维护一个 `product`。
- 语种默认取 `config.json` 的 `language`，变更用 `--language en`。

## 本地 Mock 服务（三方接口）
- 数据文件中若有引用 `/mock/*`的接口，CLI自动启动本Mock服务
- 通过`http://<本机IP>:<mock.port>`，供被测系统调用

## 退出码

| 退出码 | 含义 |
|---|---|
| 0 | 全部通过 |
| 3 | 存在失败 |
| 1 | 执行异常 |
