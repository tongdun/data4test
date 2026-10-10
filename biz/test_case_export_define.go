package biz

// ExportColumn 导出模板中的一列定义
type ExportColumn struct {
	Title string `json:"title"` // Excel 表头
	Field string `json:"field"` // 字段映射：系统列名 / test_process / ext_info.<key> / ext_info
}

// TestCaseExportTemplate 单个导出模板
type TestCaseExportTemplate struct {
	Columns []ExportColumn `json:"columns"`
}

// TestCaseExportRow 导出用的测试用例记录（含新增字段）
type TestCaseExportRow struct {
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
	TestTime     string `gorm:"column:test_time" json:"test_time"`
	TestResult   string `gorm:"column:test_result" json:"test_result"`
	Module       string `gorm:"column:module" json:"module"`
	IntroVersion string `gorm:"column:intro_version" json:"intro_version"`
	Product      string `gorm:"column:product" json:"product"`
	Remark       string `gorm:"column:remark" json:"remark"`
	TestProcess  string `gorm:"column:test_process" json:"test_process"` // 截图路径列表 JSON 数组
	ExtInfo      string `gorm:"column:ext_info" json:"ext_info"`         // 扩展信息 YAML
}

// exportImage 打包模式（screenshot_mode=path）下的一张待打包图片
type exportImage struct {
	absPath   string // 本地绝对路径
	fileName  string // 重命名后的文件名（含扩展名），由 resolveScreenshotValue 生成
	entryName string // 归档条目名（模块目录/文件名），由 buildTestCaseXlsx 填充
}
