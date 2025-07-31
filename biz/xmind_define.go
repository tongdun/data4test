package biz

type TestCase struct {
	CaseNumber   string `gorm:"column:case_number" json:"case_number"`
	CaseName     string `gorm:"column:case_name" json:"case_name"`
	CaseType     string `gorm:"column:case_type" json:"case_type"`
	Priority     string `gorm:"column:priority" json:"priority"`
	PreCondition string `gorm:"column:pre_condition" json:"pre_condition"`
	TestRange    string `gorm:"column:test_range" json:"test_range"`
	TestSteps    string `gorm:"column:test_steps" json:"test_steps"`
	ExpectResult string `gorm:"column:expect_result" json:"expect_result"`
	Auto         string `gorm:"column:auto" json:"auto"`
	Scene        string `gorm:"column:scene" json:"scene"`
	FunDeveloper string `gorm:"column:fun_developer" json:"fun_developer"`
	CaseDesigner string `gorm:"column:case_designer" json:"case_designer"`
	CaseExecutor string `gorm:"column:case_executor" json:"case_executor"`
	// TestTime     string `gorm:"column:test_time" json:"test_time"`
	TestResult   string `gorm:"column:test_result" json:"test_result"`
	Module       string `gorm:"column:module" json:"module"`
	IntroVersion string `gorm:"column:intro_version" json:"intro_version"`
	//UpdatedAt    string `gorm:"column:updated_at;autoUpdateTime:nano;<-:update"`
	Product string `gorm:"column:product" json:"product"`
	Remark  string `gorm:"column:remark" json:"remark"`
}

type XmindTestCase struct {
	Name                  string      `json:"name"`
	Version               int         `json:"version"`
	Summary               string      `json:"summary"`
	Preconditions         string      `json:"preconditions"`
	ExecutionType         int         `json:"execution_type"`
	Importance            int         `json:"importance"`
	EstimatedExecDuration int         `json:"estimated_exec_duration"`
	Status                int         `json:"status"`
	Result                int         `json:"result"`
	Steps                 []TestSteps `json:"steps"`
	Product               string      `json:"product"`
	Suite                 string      `json:"suite"`
}

type TestSteps struct {
	StepNumber      int    `json:"step_number"`
	Actions         string `json:"actions"`
	ExpectedResults string `json:"expectedresults"`
	ExecutionType   int    `json:"execution_type"`
	Result          int    `json:"result"`
}

// XMindContent represents the structure of content.json
type XMindContent []struct {
	ID        string `json:"id"`
	Class     string `json:"class"`
	Title     string `json:"title"`
	RootTopic struct {
		ID             string `json:"id"`
		Class          string `json:"class"`
		Title          string `json:"title"`
		Href           string `json:"href"`
		StructureClass string `json:"structureClass"`
		Children       struct {
			Attached []Topic `json:"attached"`
		} `json:"children"`
	} `json:"rootTopic"`
}

type Topic struct {
	Title    string `json:"title"`
	ID       string `json:"id"`
	Href     string `json:"href"`
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position"`
	Children struct {
		Attached []*Topic `json:"attached"`
	} `json:"children"`
	Branch    string     `json:"branch"`
	Markers   []*Marker  `json:"markers"`
	Summaries []*Summary `json:"summaries"`
}

type Marker struct {
	MarkerID string `json:"markerId"`
}

type Summary struct {
	Range   string `json:"range"`
	TopicID string `json:"topicId"`
}
