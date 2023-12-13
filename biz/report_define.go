package biz

type ProductCount struct {
	AllCount           int     `gorm:"column:all_count" json:"all_count"`
	AutomatableCount   int     `gorm:"column:automatable_count" json:"automatable_count"`
	UnautomatableCount int     `gorm:"column:unautomatable_count" json:"unautomatable_count"`
	AutoTestCount      int     `gorm:"column:auto_test_count" json:"auto_test_count"`
	UntestCount        int     `gorm:"column:untest_count" json:"untest_count"`
	PassCount          int     `gorm:"column:pass_count" json:"pass_count"`
	FailCount          int     `gorm:"column:fail_count" json:"fail_count"`
	AutoPer            float64 `gorm:"column:auto_per" json:"auto_per"`
	PassPer            float64 `gorm:"column:pass_per" json:"pass_per"`
	FailPer            float64 `gorm:"column:fail_per" json:"fail_per"`
	Product            string  `gorm:"column:product" json:"product"`
	App        string `gorm:"column:app" json:"app"`
}

type APICount struct {
	ApiId      string `gorm:"column:api_id" json:"api_id"`
	ApiDesc string `gorm:"column:api_desc" json:"api_desc"`
	RunTimes    int    `gorm:"column:run_times" json:"run_times"`
	TestTimes   int    `gorm:"column:test_times" json:"test_times"`
	PassTimes   int    `gorm:"column:pass_times" json:"pass_times"`
	FailTimes   int    `gorm:"column:fail_times" json:"fail_times"`
	UntestTimes int    `gorm:"column:untest_times" json:"untest_times"`
	TestResult  string `gorm:"column:test_result" json:"test_result"`
	FailReason  string `gorm:"column:fail_reason" json:"fail_reason"`
	App        string `gorm:"column:app" json:"app"`
}