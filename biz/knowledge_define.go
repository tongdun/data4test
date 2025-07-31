package biz

type KPlaybook struct {
	Name     string   `json:"场景名称"`
	DataList []string `json:"关联数据"`
}

type KData struct {
	Name     string `json:"数据描述"`
	FileName string `json:"数据名称"`
	Content  string `json:"数据详情"`
}

type KTask struct {
	Name         string   `json:"任务名称"`
	PlaybookList []string `json:"关联场景"`
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
