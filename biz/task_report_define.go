package biz

// DashboardReport 执行报告汇总
type DashboardReport struct {
	Id              string `gorm:"column:id" json:"id"`
	ReportName      string `gorm:"column:report_name" json:"report_name"`
	ReportType      string `gorm:"column:report_type" json:"report_type"`
	RelatedTaskIds  string `gorm:"column:related_task_ids" json:"related_task_ids"`
	RelatedProducts string `gorm:"column:related_products" json:"related_products"`
	RelatedApps     string `gorm:"column:related_apps" json:"related_apps"`
	TimeRangeStart  string `gorm:"column:time_range_start" json:"time_range_start"`
	TimeRangeEnd    string `gorm:"column:time_range_end" json:"time_range_end"`
	Status          string `gorm:"column:status" json:"status"`
	Creator         string `gorm:"column:creator" json:"creator"`
	Remark          string `gorm:"column:remark" json:"remark"`
	ReportData      string `gorm:"column:report_data" json:"report_data"`
	CreatedAt       string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       string `gorm:"column:updated_at" json:"updated_at"`
	//DeletedAt       string `gorm:"column:deleted_at" json:"deleted_at"`
}

// TaskReportData 任务报告数据JSON结构
type TaskReportData struct {
	Overview struct {
		TaskName        string  `json:"task_name"`
		TaskType        string  `json:"task_type"`
		Environment     string  `json:"environment"`
		ExecutionTime   string  `json:"execution_time"`
		StartTime       string  `json:"start_time"`
		EndTime         string  `json:"end_time"`
		DurationSeconds int     `json:"duration_seconds"`
		TotalExpected   int     `json:"total_expected"`
		TotalExecuted   int     `json:"total_executed"`
		NotExecuted     int     `json:"not_executed"`
		SuccessCount    int     `json:"success_count"`
		FailCount       int     `json:"fail_count"`
		PassRate        float64 `json:"pass_rate"`
		ExecuteRate     float64 `json:"execute_rate"`
		Executor        string  `json:"executor"`
	} `json:"overview"`
	SceneStats struct {
		Total    int     `json:"total"`
		Pass     int     `json:"pass"`
		Fail     int     `json:"fail"`
		PassRate float64 `json:"pass_rate"`
	} `json:"scene_stats"`
	DataStats struct {
		Total    int     `json:"total"`
		Pass     int     `json:"pass"`
		Fail     int     `json:"fail"`
		PassRate float64 `json:"pass_rate"`
	} `json:"data_stats"`
	APITypeDistribution []CountItem         `json:"api_type_distribution"`
	SceneDetails        []SceneDetail       `json:"scene_details"`
	DataDetails         []DataDetail        `json:"data_details"`
	ByProduct           []ProductReportItem `json:"by_product"`
	FailItems           []FailItem          `json:"fail_items"`
	Trend               []TrendItem         `json:"trend"`
}

type CountItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type SceneDetail struct {
	Name       string `json:"name"`
	Result     string `json:"result"`
	FailReason string `json:"fail_reason,omitempty"`
}

type DataDetail struct {
	SceneName  string `json:"scene_name,omitempty"`
	Name       string `json:"name"`
	ApiId      string `json:"api_id,omitempty"`
	Result     string `json:"result"`
	FailReason string `json:"fail_reason,omitempty"`
}

type ProductReportItem struct {
	Product  string  `json:"product"`
	Total    int     `json:"total"`
	Pass     int     `json:"pass"`
	Fail     int     `json:"fail"`
	PassRate float64 `json:"pass_rate"`
}

type FailItem struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	APIId  string `json:"api_id,omitempty"`
	Reason string `json:"reason"`
}

type TrendItem struct {
	ExecutionTime string  `json:"execution_time"`
	Total         int     `json:"total"`
	Pass          int     `json:"pass"`
	Fail          int     `json:"fail"`
	PassRate      float64 `json:"pass_rate"`
}

// MultiTaskReportData 多任务报告的 report_data JSON 结构
type MultiTaskReportData struct {
	Overview            MultiTaskOverview   `json:"overview"`
	ByTask              []TaskReportItem    `json:"by_task"`
	SceneList           []ResourceItem      `json:"scene_list"`            // 聚合: 所有关联场景
	DataList            []ResourceItem      `json:"data_list"`             // 聚合: 所有关联数据文件
	APIList             []ResourceItem      `json:"api_list"`              // 聚合: 所有关联API接口
	APITypeDistribution []CountItem         `json:"api_type_distribution"` // 聚合所有任务的API分布
	SceneDetails        []SceneDetailWithTask `json:"scene_details"`         // 所有任务的场景明细(含任务名)
	DataDetails         []DataDetailWithTask  `json:"data_details"`          // 所有任务的数据文件明细(含任务名)
	ByProduct           []ProductReportItem `json:"by_product"`
	FailItems           []FailItem          `json:"fail_items"` // 所有任务的失败明细
	Trend               []TrendItem         `json:"trend"`
}

// MultiTaskOverview 多任务报告概览
type MultiTaskOverview struct {
	ReportName      string  `json:"report_name"`
	Product         string  `json:"product"`
	TaskCount       int     `json:"task_count"`
	TotalCases      int     `json:"total_cases"`
	PassCases       int     `json:"pass_cases"`
	FailCases       int     `json:"fail_cases"`
	PassRate        float64 `json:"pass_rate"`
	SceneCount      int     `json:"scene_count"` // 关联场景总数(去重)
	DataCount       int     `json:"data_count"`  // 关联数据文件总数(去重)
	APICount        int     `json:"api_count"`   // 关联API接口总数(去重)
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	DurationSeconds int     `json:"duration_seconds"`
	ExecuteRate     float64 `json:"execute_rate"`
}

// TaskReportItem 多任务报告中单个任务的统计信息
type TaskReportItem struct {
	TaskId          string         `json:"task_id"`
	TaskName        string         `json:"task_name"`
	TaskType        string         `json:"task_type"` // scene / data
	SceneTotal      int            `json:"scene_total"`
	ScenePass       int            `json:"scene_pass"`
	SceneFail       int            `json:"scene_fail"`
	DataTotal       int            `json:"data_total"`
	DataPass        int            `json:"data_pass"`
	DataFail        int            `json:"data_fail"`
	Total           int            `json:"total"`
	Pass            int            `json:"pass"`
	Fail            int            `json:"fail"`
	PassRate        float64        `json:"pass_rate"`
	StartTime       string         `json:"start_time"`
	EndTime         string         `json:"end_time"`
	DurationSeconds int            `json:"duration_seconds"`
	Scenes          []ResourceItem `json:"scenes"` // 该任务关联的场景名
	Datas           []ResourceItem `json:"datas"`  // 该任务关联的数据文件名
	APIs            []ResourceItem `json:"apis"`   // 该任务关联的API ID
}

// ResourceItem 资源项（场景/数据文件/API）
type ResourceItem struct {
	Name  string `json:"name"`
	Count int    `json:"count,omitempty"`
}

// SceneDetailWithTask 多任务报告中带任务信息的场景明细
type SceneDetailWithTask struct {
	TaskName   string `json:"task_name"`
	Name       string `json:"name"`
	Result     string `json:"result"`
	FailReason string `json:"fail_reason,omitempty"`
}

// DataDetailWithTask 多任务报告中带任务信息的数据文件明细
type DataDetailWithTask struct {
	TaskName   string `json:"task_name"`
	Name       string `json:"name"`
	ApiId      string `json:"api_id,omitempty"`
	Result     string `json:"result"`
	FailReason string `json:"fail_reason,omitempty"`
}

// FailItemWithTask 多任务报告中带任务信息的失败项
type FailItemWithTask struct {
	TaskName string `json:"task_name"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	APIId    string `json:"api_id,omitempty"`
	Reason   string `json:"reason"`
}
