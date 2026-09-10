package biz

// CaseI18nIndex 多语种索引：[lang][module][caseNumber][field] → 译文
type CaseI18nIndex map[string]map[string]map[string]map[string]string

// 参与数据国际化的字段（与 DB 列名一致）
const (
	FieldCaseName     = "case_name"
	FieldCaseModule   = "module"
	FieldCaseType     = "case_type"
	FieldPreCondition = "pre_condition"
	FieldTestRange    = "test_range"
	FieldTestSteps    = "test_steps"
	FieldExpectResult = "expect_result"
)
