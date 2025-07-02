package biz

type AITestCase struct {
	CaseBase
	CommonExtend
}

type DbAiCase struct {
	Id string `gorm:"column:id" json:"id"`
	AITestCase
}

type InputCase struct {
	AiTemplate string `gorm:"column:ai_template" json:"ai_template"`
	CommonExtend
}

type CaseBase struct {
	CaseNumber   string `gorm:"column:case_number" json:"case_number"`
	CaseName     string `gorm:"column:case_name" json:"case_name"`
	Module       string `gorm:"column:module" json:"module"`
	CaseType     string `gorm:"column:case_type" json:"case_type"`
	Priority     string `gorm:"column:priority" json:"priority"`
	PreCondition string `gorm:"column:pre_condition" json:"pre_condition"`
	TestRange    string `gorm:"column:test_range" json:"test_range"`
	TestSteps    string `gorm:"column:test_steps" json:"test_steps"`
	ExpectResult string `gorm:"column:expect_result" json:"expect_result"`
	Auto         string `gorm:"column:auto" json:"auto"`
}
