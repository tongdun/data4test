package biz

type AIIssue struct {
	IssueBase
	CreateBase
}

type DbAIIssue struct {
	Id string `gorm:"column:id" json:"id"`
	IssueBase
	CreateBase
}

type IssueBase struct {
	IssueName           string `gorm:"column:issue_name" json:"issue_name"`
	IssueLevel          string `gorm:"column:issue_level" json:"issue_level"`
	IssueSource         string `gorm:"column:issue_source" json:"issue_source"`
	SourceName          string `gorm:"column:source_name" json:"source_name"`
	RequestData         string `gorm:"column:request_data" json:"request_data"`
	ResponseData        string `gorm:"column:response_data" json:"response_data"`
	IssueDetail         string `gorm:"column:issue_detail" json:"issue_detail"`
	ConfirmStatus       string `gorm:"column:confirm_status" json:"confirm_status"`
	RootCause           string `gorm:"column:root_cause" json:"root_cause"`
	ImpactScopeAnalysis string `gorm:"column:impact_scope_analysis" json:"impact_scope_analysis"`
	ImpactPlaybook      string `gorm:"column:impact_playbook" json:"impact_playbook"`
	ImpactData          string `gorm:"column:impact_data" json:"impact_data"`
	ResolutionStatus    string `gorm:"column:resolution_status" json:"resolution_status"`
	AgainTestResult     string `gorm:"column:again_test_result" json:"again_test_result"`
	ImpactTestResult    string `gorm:"column:impact_test_result" json:"impact_test_result"`
	ProductList         string `gorm:"column:product_list" json:"product_list"`
}
