package biz

import (
	"bytes"
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
)

// CliJSONField 表示 env_config.yaml 中以 JSON 字符串承载的字段（product.auth / private_parameter /
// private_app_prefix，以及 sys_parameter / assert_template 中值为 JSON 的项）。
// 导出时多行 JSON 用 json.Indent 格式化并以 YAML 块标量（|-）展示，普通字符串原样输出；
// 读取时 yaml.Unmarshal 仍按字符串解析（命名 string 类型），语义不变。
type CliJSONField string

// MarshalYAML 使该字段在 YAML 中按内容分两种方式输出（yaml.v3）：
//   - 多行 JSON（对象/数组）：缩进后以块标量（|-）展示，便于阅读；
//   - 单行 JSON 标量或非 JSON 字符串：按普通字符串原样输出，不做块标量转义。
func (j CliJSONField) MarshalYAML() (interface{}, error) {
	s := string(j)
	if s == "" {
		return "", nil
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(s), "", "  "); err != nil {
		// 非合法 JSON：按普通字符串原样输出
		return s, nil
	}
	pretty := buf.String()
	if !strings.Contains(pretty, "\n") {
		// 单行 JSON 标量（数字/布尔/null/字符串）：无需块标量，原样输出
		return pretty, nil
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: pretty, Style: yaml.LiteralStyle}, nil
}

// CliProduct 对应 env_config.yaml 的 product 块（字段名与 product 表保持一致）
type CliProduct struct {
	Name             string       `yaml:"name" json:"name"`
	Protocol         string       `yaml:"protocol" json:"protocol"`
	Ip               string       `yaml:"ip" json:"ip"`
	Auth             CliJSONField `yaml:"auth" json:"auth"`
	EnvType          int          `yaml:"env_type" json:"env_type"`
	Threading        string       `yaml:"threading" json:"threading"`
	ThreadNumber     int          `yaml:"thread_number" json:"thread_number"`
	Apps             string       `yaml:"apps" json:"apps"`
	Testmode         string       `yaml:"testmode" json:"testmode"`
	PrivateAppPrefix CliJSONField `yaml:"private_app_prefix" json:"private_app_prefix"`
	PrivateParameter CliJSONField `yaml:"private_parameter" json:"private_parameter"`
	Remark           string       `yaml:"remark" json:"remark"`
}

// CliApp 对应 env_config.yaml 的 app.* 块（字段名与 env_config 表保持一致）
type CliApp struct {
	Protocol  string       `yaml:"protocol" json:"protocol"`
	Ip        string       `yaml:"ip" json:"ip"`
	Prepath   string       `yaml:"prepath" json:"prepath"`
	Auth      CliJSONField `yaml:"auth" json:"auth"`
	Threading string       `yaml:"threading" json:"threading"`
	Testmode  string       `yaml:"testmode" json:"testmode"`
	Swagger   string       `yaml:"swagger_path" json:"swagger_path"`
	Remark    string       `yaml:"remark" json:"remark"`
}

// CliMock 对应 env_config.yaml 的 mock 块
type CliMock struct {
	Enable  bool   `yaml:"enable" json:"enable"`
	Port    int    `yaml:"port" json:"port"`
	BaseURL string `yaml:"base_url" json:"base_url"`
}

// CliConfig env_config.yaml 根结构（用户唯一编辑的文件，环境类信息全部集中在此）
type CliConfig struct {
	Product         CliProduct              `yaml:"product" json:"product"`
	App             map[string]CliApp       `yaml:"app" json:"app"`
	SysParameter    map[string]CliJSONField `yaml:"sys_parameter" json:"sys_parameter"`
	AssertTemplates map[string]CliJSONField `yaml:"assert_template" json:"assert_template"`
	Mock            CliMock                 `yaml:"mock" json:"mock"`
}

// cliRunOptions 解析 ./data4test run 子命令后的运行选项
type cliRunOptions struct {
	taskFilter string // 逗号分隔的任务名，空表示全部
	pkgPath    string // 显式指定的任务包路径，空表示自动选择 task/ 下最新包
	language   string // 覆盖 config.json 默认语种
	mock       bool   // 显式开启本地 mock 服务
	reconfig   bool   // 从任务包重新复制原始 env_config.yaml 覆盖本地
}
