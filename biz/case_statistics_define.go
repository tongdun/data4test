package biz

// CaseStatistics 用例统计定义表行
type CaseStatistics struct {
	Id            int    `gorm:"column:id" json:"id"`
	Name          string `gorm:"column:name" json:"name"`
	Definition    string `gorm:"column:definition" json:"definition"`
	IntroVersions string `gorm:"column:intro_versions" json:"intro_versions"`
	Product       string `gorm:"column:product" json:"product"`
	Status        string `gorm:"column:status" json:"status"`
	Creator       string `gorm:"column:creator" json:"creator"`
	Remark        string `gorm:"column:remark" json:"remark"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     string `gorm:"column:deleted_at" json:"deleted_at"`
}

// CaseStatisticsReport 用例统计报告表行
type CaseStatisticsReport struct {
	Id            int    `gorm:"column:id" json:"id"`
	Name          string `gorm:"column:name" json:"name"`
	RelatedStatId int    `gorm:"column:related_stat_id" json:"related_stat_id"`
	IntroVersions string `gorm:"column:intro_versions" json:"intro_versions"`
	Product       string `gorm:"column:product" json:"product"`
	StatStartTime string `gorm:"column:stat_start_time" json:"stat_start_time"`
	StatEndTime   string `gorm:"column:stat_end_time" json:"stat_end_time"`
	Status        string `gorm:"column:status" json:"status"`
	Creator       string `gorm:"column:creator" json:"creator"`
	ReportData    string `gorm:"column:report_data" json:"report_data"`
	Remark        string `gorm:"column:remark" json:"remark"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at" json:"updated_at"`
}

// CaseStatisticsDefinition 统计定义 YAML 解析结果
// YAML 为嵌套 map：顶层键 = 一级菜单名（ext_info 为保留键，可选），二级键 = 二级菜单名，值为模块名列表
type CaseStatisticsDefinition struct {
	ExtInfo []string           // 扩展信息中要统计的字段 key（来自 ext_info 保留键）
	Tree    []CaseStatTreeNode // 一级 → 二级 → 三级(模块名)
}

// CaseStatTreeNode 一级/二级菜单节点
type CaseStatTreeNode struct {
	Name     string             // 一级/二级菜单名
	Children []CaseStatTreeNode // 二级菜单
	Modules  []string           // 三级 = test_case.module 值
}

// CaseStatisticsReportData 用例统计报告数据（序列化到 dashboard.report_data）
type CaseStatisticsReportData struct {
	Overview           CaseStatOverview       `json:"overview"`
	ModuleStats        []CaseStatModuleItem   `json:"module_stats"`
	ResultDistribution []CountItem            `json:"result_distribution"`
	ByFunDeveloper     []CountItem            `json:"by_fun_developer"`
	ByCaseDesigner     []CountItem            `json:"by_case_designer"`
	ByPriority         []CountItem            `json:"by_priority"`
	ByCaseExecutor     []CountItem            `json:"by_case_executor"`
	ExecutorStats      []CaseStatExecutorItem `json:"executor_stats"`
	ExtInfoStats       []CaseStatExtInfoItem  `json:"ext_info_stats"`
	FailCases          []CaseStatCaseDetail   `json:"fail_cases"`
	UntestCases        []CaseStatCaseDetail   `json:"untest_cases"`
}

// CaseStatOverview 报告概览
type CaseStatOverview struct {
	Name          string  `json:"name"`
	IntroVersions string  `json:"intro_versions"`
	Product       string  `json:"product"`
	ModuleCount   int     `json:"module_count"`
	TotalCases    int     `json:"total_cases"`
	Pass          int     `json:"pass"`
	Fail          int     `json:"fail"`
	Untest        int     `json:"untest"`
	Part          int     `json:"part"`
	Deprecated    int     `json:"deprecated"`
	Unmerged      int     `json:"unmerged"`
	PassRate      float64 `json:"pass_rate"`
	ExecRate      float64 `json:"exec_rate"`
}

// CaseStatModuleItem 单个三级模块的统计项
type CaseStatModuleItem struct {
	Level1     string  `json:"level1"`
	Level2     string  `json:"level2"`
	Module     string  `json:"module"`
	Total      int     `json:"total"`
	Pass       int     `json:"pass"`
	Fail       int     `json:"fail"`
	Untest     int     `json:"untest"`
	Part       int     `json:"part"`
	Deprecated int     `json:"deprecated"`
	Unmerged   int     `json:"unmerged"`
	PassRate   float64 `json:"pass_rate"`
	ExecRate   float64 `json:"exec_rate"`
}

// CaseStatExtInfoItem 单个扩展字段的统计
type CaseStatExtInfoItem struct {
	Key   string      `json:"key"`
	Items []CountItem `json:"items"`
}

// CaseStatCaseDetail 失败/未执行用例明细
type CaseStatCaseDetail struct {
	CaseNumber string `json:"case_number"`
	CaseName   string `json:"case_name"`
	Module     string `json:"module"`
	Remark     string `json:"remark"`
}

// CaseStatExecutorItem 用例执行者统计项（含已执行/通过/失败/废弃/暂未合入/未执行/执行率）
type CaseStatExecutorItem struct {
	Name       string  `json:"name"`
	Total      int     `json:"total"`
	Executed   int     `json:"executed"`
	Pass       int     `json:"pass"`
	Fail       int     `json:"fail"`
	Deprecated int     `json:"deprecated"`
	Unmerged   int     `json:"unmerged"`
	Untest     int     `json:"untest"`
	ExecRate   float64 `json:"exec_rate"`
}
