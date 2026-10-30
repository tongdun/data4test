#!/usr/bin/env bash
# 统一发布入口：一键产出 data4test-cli 与 data4test-web 两个「单目录多二进制」发布包。
#
# 用法：
#   ./build/release.sh                                        # 默认 CLI + Web，平台 linux/amd64
#   ./build/release.sh --cli                                  # 只打 CLI
#   ./build/release.sh --web                                  # 只打 Web
#   ./build/release.sh --all                                  # 全 7 平台（CLI + Web）
#   ./build/release.sh --platforms linux/amd64,darwin/amd64   # 指定平台（GOOS/GOARCH 逗号分隔）
#   ./build/release.sh --cli --all                            # 只打 CLI，全平台
#
# 可选平台（GOOS/GOARCH，逗号分隔；默认 linux/amd64）：
#   linux/amd64  darwin/amd64  darwin/arm64  linux/386  linux/arm64  windows/amd64  windows/386
#   darwin/amd64 → data4test_darwin_x86_64（Intel Mac）；darwin/arm64 → data4test_darwin_aarch64（Apple Silicon）
#
# 说明：
#   - sqlite 使用 mattn/go-sqlite3（CGO），CLI/Web 均 CGO_ENABLED=1 编译；交叉编译需各平台可用的 C 编译器（可用 CC=... 指定）。
#   - Web 收集逻辑内联自 getDeployContent.sh（yarn build + 清理 + sed + 生产 config.json + 示例 yml），
#     但 yarn install/build 降级为 best-effort：失败不中断，回退复用仓库中已构建好的 web/static。
#     该旧脚本与 getReleasePackage.sh 不动，待新脚本跑通后人工处理；build/cli_release.sh 已废弃删除（功能已并入本脚本）。
#   - Web 清理沿用 getDeployContent.sh 的清单（无需 sudo）。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

RELEASE_DIR="$ROOT/release"
CLI_DIR="$RELEASE_DIR/data4test-cli"
WEB_DIR="$RELEASE_DIR/data4test-web"
curDate="$(date +"%Y%m%d")"

ALL_PLATFORMS=(linux/amd64 darwin/amd64 darwin/arm64 linux/386 linux/arm64 windows/amd64 windows/386)

usage() {
  sed -n '2,14p' "$0"
}

# ---------------------------------------------------------------------------
# 参数解析
# ---------------------------------------------------------------------------
DO_CLI=0
DO_WEB=0
DO_ALL=0
PLATFORMS=""

while [ $# -gt 0 ]; do
  case "$1" in
    --cli)       DO_CLI=1 ;;
    --web)       DO_WEB=1 ;;
    --all)       DO_ALL=1 ;;
    --platforms) PLATFORMS="${2:-}"; shift ;;
    -h|--help)   usage; exit 0 ;;
    *) echo "未知参数: $1" >&2; usage >&2; exit 1 ;;
  esac
  shift
done

# 未指定 --cli/--web 时默认双包
if [ "$DO_CLI" -eq 0 ] && [ "$DO_WEB" -eq 0 ]; then
  DO_CLI=1
  DO_WEB=1
fi

# ---------------------------------------------------------------------------
# 平台选择
# ---------------------------------------------------------------------------
if [ "$DO_ALL" -eq 1 ]; then
  SELECTED=("${ALL_PLATFORMS[@]}")
elif [ -n "$PLATFORMS" ]; then
  IFS=',' read -r -a SELECTED <<< "$PLATFORMS"
else
  SELECTED=(linux/amd64)
fi

# ---------------------------------------------------------------------------
# 工具函数
# ---------------------------------------------------------------------------
# 二进制命名：data4test_<os>_<arch>[.exe]，arch 映射 386→i386 / amd64→x86_64 / arm64→aarch64
bin_name() {
  local goos="$1" goarch="$2" arch
  case "$goarch" in
    386)   arch=i386 ;;
    amd64) arch=x86_64 ;;
    arm64) arch=aarch64 ;;
    *)     arch="$goarch" ;;
  esac
  local name="data4test_${goos}_${arch}"
  [ "$goos" = "windows" ] && name="${name}.exe"
  echo "$name"
}

build() {
  local goos="$1" goarch="$2" out="$3"
  echo "    build $goos/$goarch -> ${out#"$ROOT"/}"
  CGO_ENABLED=1 GOOS="$goos" GOARCH="$goarch" go build -o "$out" .
}

# 按所选平台循环编译到指定目录。
# 默认平台 linux/amd64 直接产出裸 data4test（无后缀），其余平台产出 data4test_<os>_<arch>[.exe]，
# 避免 linux/amd64 同时出现 data4test 与 data4test_linux_x86_64 两份重复二进制。
build_for_dir() {
  local dir="$1" p goos goarch name
  for p in "${SELECTED[@]}"; do
    goos="${p%/*}"
    goarch="${p#*/}"
    if [ "$p" = "linux/amd64" ]; then
      name="data4test"
    else
      name="$(bin_name "$goos" "$goarch")"
    fi
    build "$goos" "$goarch" "$dir/$name"
  done
}

# 生成跨平台启动脚本 start.sh：检测当前平台，运行对应二进制（免认平台）。
# linux/amd64 对应裸 data4test，其余平台对应 data4test_<os>_<arch>[.exe]。
write_launcher() {
  local dir="$1"
  cat > "$dir/start.sh" <<'EOF'
#!/usr/bin/env bash
# 启动脚本：检测当前平台，运行对应二进制。
set -euo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"

case "$(uname -s)" in
  Linux)  os=linux ;;
  Darwin) os=darwin ;;
  *) echo "start.sh: unsupported os $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64)        arch=x86_64 ;;
  aarch64|arm64)       arch=aarch64 ;;
  i386|i486|i586|i686) arch=i386 ;;
  *) echo "start.sh: unsupported arch $(uname -m)" >&2; exit 1 ;;
esac

if [ "$os" = "linux" ] && [ "$arch" = "x86_64" ]; then
  bin="$DIR/data4test"
else
  bin="$DIR/data4test_${os}_${arch}"
fi

if [ ! -x "$bin" ]; then
  echo "start.sh: binary not found: $bin" >&2
  exit 1
fi

exec "$bin" "$@"
EOF
  chmod +x "$dir/start.sh"
}

# 清理并重建输出目录；遇 root 所有的残留目录时显式报错提示手动清理
clean_dir() {
  local d="$1"
  if [ -e "$d" ] && ! rm -rf "$d" 2>/dev/null; then
    echo "无法删除 ${d}（可能为 root 所有），请先手动执行: sudo rm -rf ${d}" >&2
    exit 1
  fi
  mkdir -p "$d"
}

# 打成压缩包：tar -czf release/<prefix>_<curDate>.tgz，归档内为 <dir>/ 顶层目录。
# Web 用 data4test-web_${curDate}.tgz，CLI 用 data4test-cli_${curDate}.tgz，便于区分。
pack_tgz() {
  local dir="$1" prefix="$2"
  local out="$RELEASE_DIR/${prefix}_${curDate}.tgz"
  tar -czf "$out" -C "$RELEASE_DIR" "$dir"
  echo "  ✓ $out"
}

# ---------------------------------------------------------------------------
# CLI 包
# ---------------------------------------------------------------------------
build_cli() {
  echo "[cli] 生成 $CLI_DIR"
  clean_dir "$CLI_DIR"

  # 运行骨架目录
  mkdir -p "$CLI_DIR/mgmt/log" \
           "$CLI_DIR/mgmt/data" \
           "$CLI_DIR/mgmt/upload" \
           "$CLI_DIR/mgmt/download" \
           "$CLI_DIR/mgmt/history" \
           "$CLI_DIR/task" \
           "$CLI_DIR/report"

  # UI 文案（biz.T 只读 i18n_system/{locale}.json；业务名语料由 run 时 installCliI18n 自动安装）
  mkdir -p "$CLI_DIR/mgmt/i18n/i18n_system"
  cp "$ROOT/mgmt/i18n/i18n_system/"*.json "$CLI_DIR/mgmt/i18n/i18n_system/" 2>/dev/null || true

  # 精简 config.json：CLI 不连 MySQL，去除数据库凭据，仅保留运行骨架所需字段
  cat > "$CLI_DIR/config.json" <<'EOF'
{
	"file_base_path": "./mgmt",
	"server_port": 9088,
	"log_level": "debug",
	"language": "cn",
	"redirect_path": "{}"
}
EOF

  # env_config.example.yaml 示例（占位符语法 {Name}，运行时逐个提示填写）
  cat > "$CLI_DIR/env_config.example.yaml" <<'EOF'
product:
    name: "{Env_Name}"  // 类{NameXX}的变量，执行前需替换
    protocol: http
    ip: "{IP_PORT}"
    auth: |-
        {
          "Content-Type": "application/x-www-form-urlencoded",
          "Cookie": "lang=cn;  _qjt_ac_={qjt_ac}",
          "Referer": "http://{IP_PORT}",
          "X-Cf-Random": "{XCf_Random}"
        }
    env_type: 2
    threading: "no"
    thread_number: 1
    apps: appName1,appName2
    testmode: custom
    private_app_prefix: |-
        {
          "appName1": "/appName1Api",
          "appName2": "/appName2Api"
        }
    private_parameter: |-
        {
          "parameter1": "value1",
          "parameter2": "value2",
          "parameter3": "parameter3"
        }
    remark: ""
app:
    appName1:
        protocol: http
        ip: "192.168.x.x:80xx"
        prepath: /appName1Api
        auth: '{}'
        threading: "no"
        testmode: custom
        swagger_path: http://192.168.x.x/api/v3/api-docs
        remark: ""
    appName2:
        protocol: http
        ip: ""
        prepath: /noahApi
        auth: '{}'
        threading: "no"
        testmode: fuzzing
        swagger_path: http://92.168.x.x:80xx/api/v2/api-docs?group=group1
        remark: ""
sys_parameter:
    Bool: true, false
    MockTemplateFile: mock_file_1.json,mock_file_2.csv
    Sex: 男, 女, 未知, man, woman
    YorN: 是, 否, YES, NO, Yes, No
    fileName: file,pickleFiles, templateFile
    scriptRunEngine: |-
        {
          ".py": "/usr/bin/python3",
          ".sh": "/bin/bash",
          ".jmx": "xxxx",
          ".bat": "xxxx"
        }
assert_template:
    getSuccess: |-
        {
          "ch": "成功",
          "en": "success|Success|succeeded|Operation successful"
        }
    updateMsg: |-
        {
          "ch": "成功|重复",
          "en": "success|Success|ok|Duplicate tag:"
        }
mock:
    enable: false
    port: 9088
    base_url: ""
EOF

  cp "$ROOT/build/cli/README.md" "$CLI_DIR/README.md"

  build_for_dir "$CLI_DIR"
  write_launcher "$CLI_DIR"
  echo "  ✓ $CLI_DIR"
  pack_tgz "data4test-cli" "data4test-cli"
}

# ---------------------------------------------------------------------------
# Web 包
# ---------------------------------------------------------------------------
build_web() {
  echo "[web] 生成 $WEB_DIR"
  clean_dir "$WEB_DIR"

  # 前端构建：尽力执行 yarn install / yarn build，失败不中断（前端工具链在当前环境可能失效，
  # 如 package.json version "1.0" 非法 semver、statik 未安装、webpack 3 UglifyJs 无法压缩 ES6），
  # 回退使用仓库中已构建好的 web/static。
  echo "[web] yarn install && yarn build ..."
  yarn install || echo "[web] yarn install 失败，改用已有静态资源" >&2
  yarn build   || echo "[web] yarn build 失败，改用已有静态资源" >&2
  if [ ! -f "$ROOT/web/static/js/app.min.js" ] || [ ! -f "$ROOT/web/static/js/vendor.min.js" ]; then
    echo "[web] 警告: 未找到 web/static/js/{app,vendor}.min.js（前端未构建），web 包将缺少前端静态资源。" >&2
  fi

  # 部署打包（拷贝骨架；源 mgmt 含历史运行数据，随后清理）
  cp -rf "$ROOT/html" "$WEB_DIR/" 2>/dev/null || true
  cp -rf "$ROOT/mgmt" "$WEB_DIR/" 2>/dev/null || true
  mkdir -p "$WEB_DIR/web/static"
  cp -rf "$ROOT/web/static" "$WEB_DIR/web/" 2>/dev/null || true

  # 清理运行时临时目录（沿用 getDeployContent.sh 的清单）。
  # 注：文件均由 cp 以当前用户复制而来，无需 sudo；改用普通 rm，避免非交互环境下 sudo 失败导致清理未生效。
  rm -f  "$WEB_DIR"/mgmt/log/*                          || true
  rm -f  "$WEB_DIR"/mgmt/api/*                          || true
  rm -rf "$WEB_DIR"/mgmt/old/*                          || true
  rm -rf "$WEB_DIR"/mgmt/history/*                      || true
  rm -rf "$WEB_DIR"/mgmt/download/*                     || true
  rm -rf "$WEB_DIR"/mgmt/ai_data/*                      || true
  rm -rf "$WEB_DIR"/mgmt/i18n/i18n_case/*               || true
  rm -rf "$WEB_DIR"/mgmt/i18n/i18n_caseCount/*          || true
  rm -rf "$WEB_DIR"/mgmt/i18n/i18n_data/*               || true
  rm -rf "$WEB_DIR"/mgmt/i18n/i18n_playbook/*           || true
  rm -rf "$WEB_DIR"/mgmt/i18n/i18n_task/*               || true
  rm -rf "$WEB_DIR"/mgmt/i18n_snap/*                    || true
  rm -rf "$WEB_DIR"/mgmt/knowledge/*                    || true
  rm -rf "$WEB_DIR"/mgmt/upload/*                       || true
  rm -rf "$WEB_DIR"/mgmt/data/*                         || true

  # case/common/upload 内仅保留 demo / 模板说明 / _template
  ( cd "$WEB_DIR/mgmt/case"   && rm -rf $(ls | grep -v "project_V1.0.0_testcase_demo.xmind") ) || true
  ( cd "$WEB_DIR/mgmt/common" && rm -rf $(ls | grep -v "模板使用说明.") )                    || true
  ( cd "$WEB_DIR/mgmt/upload" && rm -rf $(ls | grep -v "_template.") )                     || true

  # 清理 macOS 元数据 .DS_Store
  find "$WEB_DIR" -name '.DS_Store' -type f -delete 2>/dev/null || true

  # sql 目录仅保留最新的 init_all_*.sql 与 update.sql，其余 SQL 文件清理
  keep_init="$(ls "$WEB_DIR"/mgmt/sql/init_all_*.sql 2>/dev/null | sort | tail -1)"
  if [ -n "$keep_init" ]; then
    find "$WEB_DIR"/mgmt/sql -maxdepth 1 -type f \
      ! -name 'update.sql' \
      ! -name "$(basename "$keep_init")" \
      -delete 2>/dev/null || true
  fi

  mv "$WEB_DIR/mgmt/doc/file/image" "$WEB_DIR/mgmt/common" 2>/dev/null || true

  # 文档图片路径 / 链接替换（macOS 版 sed）
  sed -i '' 's/..\/..\/..\/upload/\/admin\/fm\/upload\/preview?path=/g' "$WEB_DIR/mgmt/doc/file/design/action_design.md" 2>/dev/null || true
  sed -i '' 's/..\/..\/..\/upload/\/admin\/fm\/upload\/preview?path=/g' "$WEB_DIR/mgmt/doc/file/design/mock_design.md" 2>/dev/null || true
  sed -i '' 's/<img src="..\/image\/arch.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/arch.jpg)/g' "$WEB_DIR/mgmt/doc/file/arch/arch.md" 2>/dev/null || true
  sed -i '' 's/<img src="..\/image\/全局使用流程图.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/全局使用流程图.jpg)/g' "$WEB_DIR/mgmt/doc/file/arch/arch.md" 2>/dev/null || true
  sed -i '' 's/<img src="..\/image\/数据流程图.png">/[链接](\/admin\/fm\/common\/preview?path=\/image\/数据流程图.png)/g' "$WEB_DIR/mgmt/doc/file/arch/arch.md" 2>/dev/null || true
  sed -i '' 's/<img src="..\/image\/系统数据关系图.jpg">/[链接](\/admin\/fm\/common\/preview?path=\/image\/系统数据关系图.jpg)/g' "$WEB_DIR/mgmt/doc/file/arch/arch.md" 2>/dev/null || true
  sed -i '' 's/<img src="..\/image\/智能数据流转图.png">/[链接](\/admin\/fm\/common\/preview?path=\/image\/智能数据流转图.png)/g' "$WEB_DIR/mgmt/doc/file/ai/ai_analysis_design.md" 2>/dev/null || true
  sed -i '' 's/<img src="..\/image\/智能数据流转图.png">/[链接](\/admin\/fm\/common\/preview?path=\/image\/智能数据流转图.png)/g' "$WEB_DIR/mgmt/doc/file/ai/ai_data_design.md" 2>/dev/null || true
  sed -i '' 's/<img src="..\/image\/智能数据流转图.png">/[链接](\/admin\/fm\/common\/preview?path=\/image\/智能数据流转图.png)/g' "$WEB_DIR/mgmt/doc/file/ai/ai_playbook_design.md" 2>/dev/null || true
  sed -i '' 's/<img src="..\/image\/智能用例流转图.png">/[链接](\/admin\/fm\/common\/preview?path=\/image\/智能用例流转图.png)/g' "$WEB_DIR/mgmt/doc/file/ai/ai_testcase_design.md" 2>/dev/null || true
  sed -i '' 's/parameter_design.md/parameter_design/g' "$WEB_DIR/mgmt/doc/file/design/data_file_design.md" 2>/dev/null || true
  sed -i '' 's/action_design.md/action_design/g' "$WEB_DIR/mgmt/doc/file/design/data_file_design.md" 2>/dev/null || true
  sed -i '' 's/assert_design.md/assert_design/g' "$WEB_DIR/mgmt/doc/file/design/data_file_design.md" 2>/dev/null || true
  sed -i '' 's/script_design.md/script_design/g' "$WEB_DIR/mgmt/doc/file/design/data_file_design.md" 2>/dev/null || true

  # 生产 config.json（DB 指向容器内 db:3306 等）
  cat > "$WEB_DIR/config.json" <<'EOF'
{
	"custom_foot_html": "./html/tongdun.html",
	"footer_info": "./html/tongdun.html",
	"open_admin_api": true,
	"database": {
		"default": {
			"host": "db",
			"port": "3306",
			"user": "root",
			"pwd": "password",
			"name": "data4test",
			"max_idle_con": 5,
			"max_open_con": 10,
			"driver": "mysql",
			"parse_time": true
		}
	},
	"app_id": "uPkhI73C0y3p",
	"language": "cn",
	"prefix": "admin",
	"theme": "sword",
	"store": {
		"path": "./mgmt/upload",
		"prefix": "uploads"
	},
	"title": "盾测-自动化",
	"logo": "盾测-自动化",
	"mini_logo": "盾测",
	"index": "/",
	"login_url": "/login",
	"debug": true,
	"sql_log": false,
	"env": "local",
	"info_log": "./mgmt/log/info.log",
	"error_log": "./mgmt/log/error.log",
	"access_log": "./mgmt/log/access.log",
	"session_life_time": 7200,
	"file_upload_engine": {
		"name": "local"
	},
	"login_title": "盾测-自动化",
	"login_logo": "盾测-自动化",
	"auth_user_table": "goadmin_users",
	"bootstrap_file_path": "./bootstrap.go",
	"go_mod_file_path": "./go.mod",
	"asset_root_path": "./public/",
	"file_base_path": "./mgmt",
	"server_port": 9088,
	"log_level": "release",
	"cicd_host": "X.X.X.X:8088",
	"swagger_path": "http://{host:port}/api/metadata/rest/docs?group=group1",
	"redirect_path": "{\"Administrator\":\"/admin/info/schedule\", \"Operator\":\"/admin/info/schedule\", \"ApiManage\":\"/admin/likePostman\"}"
}
EOF

  # 示例数据 yml
  cat > "$WEB_DIR/mgmt/data/示例-用户管理-新建用户.yml" <<'EOF'
---
# 用例信息
name: 示例-用户管理-新建用户 # 数据用例名称，e.g.: 类型-模块-用例， 类型：功能/性能/异常/内置/……， 模块：用户管理/规则管理/……
api_id: post_/path        # 用例ID, method_path组合，后续做数据联动使用，数据统计使用
version: 1              # 数据用例版本，后续可以进行数据升级
is_run_pre_apis: "no"     # 是否跑前置用例，选项：yes / no,  默认 no， 功能未开发
is_run_post_apis: "no"    # 是否跑后置用例，选项：yes / no,  默认 no， 功能未开发
is_parallel: "no"         # 是否并行跑数据，选项：yes / no,  默认 no，
is_use_env_config: "yes"  # 是否使用公共环境，选项：yes / no,  默认 yes

# 环境信息
env:
  protocol: http        # http 或 https，请求协议
  host: X.X.X.X:8089    # 环境IP 或 环境域名 或 环境IP:端口
  prepath: /prefix      # 路由前缀，公共部分可以抽出来

# API 基本信息
api:
  description: 新建用户   # API用途
  module: 用户管理        # API所属模块
  app: appName           # API所属应用
  method: post           # （注意：保证正确） API请求方法
  path: /path            # （注意：保证正确）API请求路由，路由前缀抽离到prepath下时或公共环境中已定义prepath时，这里无需再写路由前缀
  pre_apis: []           # 调试时，依赖前置用例时，可以把关联前置文件写上，功能未充分验证
  param_apis: []         # 调试时，依赖其他用例的参数时，可以把关联文件写上，功能未充分验证
  post_apis: []          # 调试时，测试跑完后需要跑的用例，可以把关联文件写上，功能未充分验证

# 定义单值参数，如果is_use_env_config值为no, 需要定义此处的 header
single:
  header:
    Content-Type: multipart/form-data   # 如果api为导入文件功能，需要把Content-Type定义为multipart/form-data进行公用环境值的覆盖，优化级：数据文件>应用配置>产品配置
  query: {}                             # GET请求时，请求参数定义，只定义一个值，共用的参数放在这里，无需反复定义
  path: {}                              # PATH 变量参数定义，只定义一个值
  body:
    condition: "value":"{nameList}"  # {nameList} 代表字符串里有需替换的变更，nameList为 ouput 中输出的参数名字，在前置的用例中有定义同名变量，即会替换

# 定义多值参数
multi:
  query: {}                   # GET请求时，请求参数定义，定义的值为列表
  path: {}                    # PATH 变更参数定义，定义的值为列表
  body:
    description:              # 定义多值时，取各项定义的个数最少的数据，一一对应
      - "{Rune(128)}"    # 获取设置长度的汉字
      - "{Str(64)}"      # 获取设置长度的字符串
      - "{Int(10,100)}"  # 获取设置范围内的整数
    displayName:
      - "{Date(-2)}"      # 获取两天前的日期
      - "{Date(2)}"       # 获取两天后的日期
      - "{Timestamp(-2)}" # 获取两天前的时间戳
    name:
      - "{IDNo}"          # 获取身份证字符串
      - "{Name}"          # 获取名字字符串
      - "{Address}"       # 获取地址字符串
      - "{BankNo}"        # 获取银行卡号字符串

# 断言，数据校验，根据需要写不同类型的断言，不写断言，只要返回为200，即算 PASS
assert:
  - type: equal   # 验证code的值等于200
    source: code    # 返回的json信息，取key为code的值
    value: 200
  - type: "!=" # 验证code的值不等于200
    source: code    # 返回的json信息，取key为code的值
    value: 200
  - type: ">="    # 验证source字段大于等于1
    source: data-total     # 返回的json信息，data字典-取出productDesc的值
    value: 1
  - type:  contain
    source: data-contents*productDesc  # 返回的json信息，data字典-content字典-字典列表，取出productDesc的值, 并校验是否包含 value中的值
    value: 待删除的产品描述
  - type: "!in"   # 验证取到的productName的值包含删除
    source: data-contents*productName  # 返回的json信息，data字典-content字典-字典列表，取出productName的值, 不包含 value 中的值
    value: 删除
  - type: not_contain   # 验证取到的productName的值不包含删除
    source: data-contents*productName  # 返回的json信息，data字典-content字典-字典列表，取出productName的值
    value: 产品
  - type: re
    source: message
    value: 成功|重复|已存在
  - type: output  # 从返回的json 信息取取出 uuid 的值，并命名为uuid
    source: data-contents*uuid
    value: uuid
  - type: output  # 从返回的json信息取出uuid的值，并重命名为ProductUuid
    source: data-contents*uuid
    value: ProductUuid
# 数据执行后的动作
action:
  - type: sleep
    value: 1    // 表示等待1秒种，时间可根据需要自动设置，单位为秒
  - type: create_csv
    value: name:number    // 生成文件名:生成的数据条数，默认生成10条
  - type: create_xls
    value: name:number    // 生成文件名:生成的数据条数, 默认生成10条
# 输出其他接口需要的依赖数据, 由断言中类型为 ouput 定义，自动生成, 定义为'{self}'或 '{uuid}' 从此处取值
output:
  uuid:
    - XXX
    - XXX

# 测试结果：pass, fail, untest, 自动生成，断言全部符合要求设为pass, 请求若返回非200，直接置为 fail, 如果执行次数测试为0，测置为 untest
# 保留最新测试结果
test_result:
  - pass
  - fail
  - untest

# 请求 URL，自动生成， 保留最新测试结果
urls:
  - http://X.X.X.X:8089/creditApi/loanProduct/pageProductList

# 请求数据，body, query, 自动生成, 保留最新测试结果
requests:
  - {"curPage":"1","endEntryTime":"1627095420000","pageSize":"10","searchOption":"{}"startEntryTime":"1626749820000","timeType":"1"}

# 返回信息, 自动生成， 保留最新测试结果
response:
  - "response1"
  - "response2"
EOF

  build_for_dir "$WEB_DIR"
  write_launcher "$WEB_DIR"
  echo "  ✓ $WEB_DIR"
  pack_tgz "data4test-web" "data4test-web"
}

# ---------------------------------------------------------------------------
# 主流程
# ---------------------------------------------------------------------------
echo "打包平台: ${SELECTED[*]}"
[ "$DO_CLI" -eq 1 ] && build_cli
[ "$DO_WEB" -eq 1 ] && build_web

echo ""
echo "发布完成:"
[ "$DO_CLI" -eq 1 ] && echo "  CLI: $CLI_DIR"
[ "$DO_WEB" -eq 1 ] && echo "  Web: $WEB_DIR"
