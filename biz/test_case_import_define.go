package biz

// TestCaseImportItem 单条待导入用例（字段与 test_case 表列一致，另含截图信息）
type TestCaseImportItem struct {
	CaseNumber   string   `gorm:"column:case_number" json:"case_number"`
	CaseName     string   `gorm:"column:case_name" json:"case_name"`
	CaseType     string   `gorm:"column:case_type" json:"case_type"`
	Priority     string   `gorm:"column:priority" json:"priority"`
	PreCondition string   `gorm:"column:pre_condition" json:"pre_condition"`
	TestRange    string   `gorm:"column:test_range" json:"test_range"`
	TestSteps    string   `gorm:"column:test_steps" json:"test_steps"`
	ExpectResult string   `gorm:"column:expect_result" json:"expect_result"`
	Auto         string   `gorm:"column:auto" json:"auto"`
	Scene        string   `gorm:"column:scene" json:"scene"`
	FunDeveloper string   `gorm:"column:fun_developer" json:"fun_developer"`
	CaseDesigner string   `gorm:"column:case_designer" json:"case_designer"`
	CaseExecutor string   `gorm:"column:case_executor" json:"case_executor"`
	TestTime     string   `gorm:"column:test_time" json:"test_time"`
	TestResult   string   `gorm:"column:test_result" json:"test_result"`
	Module       string   `gorm:"column:module" json:"module"`
	IntroVersion string   `gorm:"column:intro_version" json:"intro_version"`
	Product      string   `gorm:"column:product" json:"product"`
	Remark       string   `gorm:"column:remark" json:"remark"`
	ExtInfo      string   `gorm:"column:ext_info" json:"ext_info"`
	TestProcess  string   `gorm:"column:test_process" json:"test_process"`
	ImageFiles   []string `gorm:"-" json:"image_files"` // 关联截图文件名，如 [TC_1_1.png, TC_1_2.png]
	Row          int      `gorm:"-" json:"row"`         // Excel 行号（1 起），用于重复提示定位
}

// ConflictInfo 冲突明细（用于前端展示「编号 + 名称」）
type ConflictInfo struct {
	CaseNumber string `json:"case_number"`
	Module     string `json:"module"`
	ExistName  string `json:"exist_name"`  // 库中已有用例名称
	ImportName string `json:"import_name"` // 导入文件中的用例名称
}

// DuplicateInfo Excel 内部用例编号重复明细（case_number + module 相同）
type DuplicateInfo struct {
	CaseNumber string `json:"case_number"` // 重复的用例编号
	Module     string `json:"module"`      // 所属模块
	Count      int    `json:"count"`       // 出现次数
	Rows       []int  `json:"rows"`        // 出现行号（1 起）
}

// TestCaseImportResult 导入检查结果
type TestCaseImportResult struct {
	ImportId   string               `json:"import_id"`
	UserName   string               `json:"user_name"`
	Template   string               `json:"template"`
	Cases      []TestCaseImportItem `json:"cases"`
	Conflicts  []ConflictInfo       `json:"conflicts"`
	Duplicates []DuplicateInfo      `json:"duplicates"`
	NewCount   int                  `json:"new_count"`
}

// LastCaseImportId 最近一次导入检查的 ID，供确认导入阶段使用
var LastCaseImportId string
