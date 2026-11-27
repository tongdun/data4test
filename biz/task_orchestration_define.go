package biz

import "sync"

// TaskOrchestration 任务编排
type TaskOrchestration struct {
	OrchestrationName   string `gorm:"column:orchestration_name" json:"orchestration_name" yaml:"orchestration_name"`
	TaskList            string `gorm:"column:task_list" json:"task_list" yaml:"task_list"`
	ProductList         string `gorm:"column:product_list" json:"product_list" yaml:"product_list"`
	OrchestrationStatus string `gorm:"column:orchestration_status" json:"orchestration_status" yaml:"orchestration_status"`
	Remark              string `gorm:"column:remark" json:"remark" yaml:"remark"`
	LastAt              string `gorm:"column:last_at" json:"last_at" yaml:"last_at"`
	UserName            string `gorm:"column:user_name" json:"user_name" yaml:"user_name"`
	CreatedAt           string `gorm:"column:created_at" json:"created_at" yaml:"created_at"`
	UpdatedAt           string `gorm:"column:updated_at" json:"updated_at" yaml:"updated_at"`
	DeletedAt           string `gorm:"column:deleted_at" json:"deleted_at" yaml:"deleted_at"`
}

// DbTaskOrchestration 带主键的任务编排
type DbTaskOrchestration struct {
	Id string `gorm:"column:id" json:"id"`
	TaskOrchestration
}

const (
	OrchestrationStatusNotStarted = "not_started"
	OrchestrationStatusRunning    = "running"
	OrchestrationStatusFinished   = "finished"
	OrchestrationStatusStopped    = "stopped"

	OrchestrationStepSuccess = "success"
	OrchestrationStepFail    = "fail"
	OrchestrationStepSkipped = "skipped"

	OrchestrationReportType    = "orchestration"
	OrchestrationTriggerManual = "manual"

	// OrchestrationReportStatsVersion 报告统计口径版本；低于该版本的历史报告在展示时回溯重算
	OrchestrationReportStatsVersion = 1
)

// OrchestrationSummary 编排报告概览
type OrchestrationSummary struct {
	Total           int     `json:"total"`
	Success         int     `json:"success"`
	Fail            int     `json:"fail"`
	Skipped         int     `json:"skipped"`
	DurationSeconds int     `json:"duration_seconds"`
	PassRate        float64 `json:"pass_rate"`
}

// OrchestrationStep 编排报告中的单个步骤
type OrchestrationStep struct {
	StepOrder       int    `json:"step_order"`
	TaskId          string `json:"task_id"`
	TaskName        string `json:"task_name"`
	TaskType        string `json:"task_type"`
	Status          string `json:"status"` // success / fail / skipped
	StartAt         string `json:"start_at"`
	EndAt           string `json:"end_at"`
	DurationSeconds int    `json:"duration_seconds"`
	DashboardId     string `json:"dashboard_id"`
	TaskTag         string `json:"task_tag"` // 本轮任务执行标记（历史报告可能为空，可按下标回溯）
	Total           int    `json:"total"`
	Success         int    `json:"success"`
	Fail            int    `json:"fail"`
	ErrorMsg        string `json:"error_msg"`
}

// OrchestrationResourceStats 编排报告中场景 / 数据文件的执行统计
type OrchestrationResourceStats struct {
	Total    int     `json:"total"`
	Success  int     `json:"success"`
	Fail     int     `json:"fail"`
	PassRate float64 `json:"pass_rate"`
}

// OrchestrationProduct 编排报告中的产品维度汇总
// 任务维度按步骤状态统计；场景 / 数据文件维度按该产品在本轮执行历史统计
type OrchestrationProduct struct {
	Product      string  `json:"product"`
	TaskTotal    int     `json:"task_total"`
	TaskSuccess  int     `json:"task_success"`
	TaskFail     int     `json:"task_fail"`
	SceneTotal   int     `json:"scene_total"`
	SceneSuccess int     `json:"scene_success"`
	SceneFail    int     `json:"scene_fail"`
	DataTotal    int     `json:"data_total"`
	DataSuccess  int     `json:"data_success"`
	DataFail     int     `json:"data_fail"`
	DataPassRate float64 `json:"data_pass_rate"`
}

// OrchestrationReportData 编排报告数据（存 dashboard.report_data）
type OrchestrationReportData struct {
	OrchestrationId   string                     `json:"orchestration_id"`
	OrchestrationName string                     `json:"orchestration_name"`
	RunTag            string                     `json:"run_tag"`
	TriggerType       string                     `json:"trigger_type"`
	ProductList       string                     `json:"product_list"`
	StatsVersion      int                        `json:"stats_version"` // 统计口径版本，见 OrchestrationReportStatsVersion
	Summary           OrchestrationSummary       `json:"summary"`
	SceneStats        OrchestrationResourceStats `json:"scene_stats"`
	DataStats         OrchestrationResourceStats `json:"data_stats"`
	Steps             []OrchestrationStep        `json:"steps"`
	Products          []OrchestrationProduct     `json:"products"`
}

// orchestrationRunning 同一编排同一时间只允许一轮执行
var orchestrationRunning sync.Map
