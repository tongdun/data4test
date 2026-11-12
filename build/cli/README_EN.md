# data4test CLI Execution Tool

## Quick Start

1. Place the task package (`*.tgz`) exported from the Web console into the `task/` directory.
2. On first run, `env_config.yaml` is automatically copied from the task package to the current directory; edit it directly to override any site-specific fields (required fields listed below).
3. Run: `./data4test run` (for multi-platform packages, use `./start.sh run` to auto-select the platform; if the config contains `{Name}` placeholders, you'll be prompted to fill them in one by one).
4. View reports:
   - `./data4test report` lists recent reports
   - `./data4test report <ID>` exports HTML and opens it in the browser (output to `report/`)

## Directory Structure

```
data4test-cli/
├── data4test                # CLI executable (linux/amd64)
├── data4test_<os>_<arch>    # binaries for other platforms (multi-platform package only)
├── start.sh                 # cross-platform launcher, auto-selects the matching binary
├── config.json              # runtime config (language etc., do not delete)
├── env_config.example.yaml  # env_config example (provided for reference only)
├── env_config.yaml          # active config (auto-copied from the task package on first run, editable)
├── task/                    # task packages (*.tgz)
├── report/                  # exported HTML reports
└── mgmt/                    # runtime data (cli.db, history, uploads, etc.)
```

## Supported Platforms

The default platform linux/amd64 produces a bare `data4test`; other platforms produce `data4test_<os>_<arch>[.exe]`. `start.sh` detects the current platform and runs the matching binary.

| Platform | Binary |
|---|---|
| linux/amd64 | `data4test` |
| linux/386 | `data4test_linux_i386` |
| linux/arm64 | `data4test_linux_aarch64` |
| darwin/amd64 (Intel Mac) | `data4test_darwin_x86_64` |
| darwin/arm64 (Apple Silicon) | `data4test_darwin_aarch64` |
| windows/amd64 | `data4test_windows_x86_64.exe` |
| windows/386 | `data4test_windows_i386.exe` |

> A single-platform package contains only one `data4test` — run it directly; for multi-platform packages, use `./start.sh` to skip platform detection.

## Commands

| Command | Description |
|---|---|
| `./data4test run [options]` | Execute all or specified tasks in the task package |
| `./data4test report` | List the 20 most recent reports |
| `./data4test report <ID>` | Export that report as HTML and open it in the browser |
| `./data4test report --multi --task <id1>,<id2>` | Regenerate a multi-task report from report IDs (see the "ID" column in the report list; must be single-task reports of different tasks) |
| `./data4test help` | Show help |

### run options

| Option | Description |
|---|---|
| `--task TaskA,TaskB` | Execute only the specified tasks (default: all) |
| `--pkg task/xxx.tgz` | Explicitly specify the task package (default: the latest one under `task/`) |
| `--reconfig` | Re-copy the original env_config.yaml from the task package to overwrite the local one and re-enter variables |
| `--param Name=value` | Inject placeholder values from the command line; multiple comma-separated and repeatable (e.g. `--param A=1,B=2`), auto-written back so the next run is prompt-free |
| `--language en` | Language for this run |
| `--mock` | Start a local Mock service in the background |

## Variable Placeholders {Name}

Config values may use `{Name}` placeholders (e.g. `{HostIP}`, `{Auth}`). At runtime, detected placeholders are prompted one by one (input is echoed in full); non-empty input is auto-written back to `env_config.yaml` for future runs; pressing Enter skips and keeps the placeholder, continuing execution as-is.

Three ways to change variables:

- Edit `env_config.yaml` before execution, replacing `{Name}` with literal values;
- Inject via `--param Name=value` from the command line (multiple comma-separated, repeatable); injected values are auto-written back to `env_config.yaml` for future runs. Injecting all of them runs directly; partial injection continues to prompt for the remaining placeholders;
- Enter values interactively at runtime (each detected placeholder is prompted until none remain).

`--reconfig` re-copies the original `env_config.yaml` from the task package to overwrite the local one and re-trigger prompting.

## Configuration Requirements

- `env_config.yaml` sits alongside the `data4test` binary, in YAML format; it is auto-copied from the task package on first run, or re-copied with `--reconfig`.
- Field names are consistent with the `product` / `env_config` / `sys_parameter` tables; a single product maintains one `product`.
- The default language comes from `language` in `config.json`; change it with `--language en`.

## Local Mock Service (third-party interfaces)
- If a data file references a `/mock/*` interface, the CLI automatically starts this Mock service
- It is exposed at `http://<host-ip>:<mock.port>` for the system under test to call

## Exit Codes

| Exit code | Meaning |
|---|---|
| 0 | All passed |
| 3 | Some failures |
| 1 | Execution error |
