package biz

// SceneTypeRunTypeMap SceneType 数值到 YAML 执行类型字符串的映射
var SceneTypeToRunType = map[int]string{
	1: "stop",
	2: "compare",
	3: "continue",
	4: "parallel",
	5: "parallel_compare",
}

// RunTypeToSceneType YAML 执行类型字符串到 SceneType 数值的映射
var RunTypeToSceneType = map[string]int{
	"stop":             1,
	"compare":          2,
	"continue":         3,
	"parallel":         4,
	"parallel_compare": 5,
}

// TaskTypeToYaml DB 中的任务类型值转换为 YAML 中展示的任务类型值
var TaskTypeToYaml = map[string]string{
	"scene": "playbook",
	"data":  "data",
}

// TaskTypeFromYaml YAML 中的任务类型值转换为 DB 中存储的任务类型值
var TaskTypeFromYaml = map[string]string{
	"playbook": "scene",
	"data":     "data",
}

type KPlaybook struct {
	Name     string   `json:"name" yaml:"-"`
	DataList []string `json:"关联数据" yaml:"关联数据"`
	RunType  string   `json:"执行类型" yaml:"执行类型"`
}

type KData struct {
	Name     string `json:"数据描述"`
	FileName string `json:"数据名称"`
	Content  string `json:"数据详情"`
}

type KTask struct {
	Name         string   `json:"name" yaml:"-"`
	TaskMode     string   `json:"任务模式" yaml:"任务模式"`
	TaskType     string   `json:"任务类型" yaml:"任务类型"`
	PlaybookList []string `json:"关联场景" yaml:"关联场景"`
	Remark       string   `json:"备注" yaml:"备注"`
}

type KCase struct {
	CaseNumber   string `json:"测试编号"`
	CaseName     string `json:"用例名称"`
	CaseType     string `json:"用例类型"`
	Priority     string `json:"优先级"`
	Module       string `json:"所属模块"`
	PreCondition string `json:"预置条件"`
	TestRange    string `json:"测试范围"`
	TestSteps    string `json:"测试步骤"`
	ExpectResult string `json:"预期结果"`
	Auto         string `json:"是否支持自动化"`
}
