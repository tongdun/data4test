package pages

import (
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"html/template"
)

type BaseCount struct {
	Infos  []string            `json:"infos"`
	Counts []float64           `json:"counts"`
	Colors []chartjs.Color     `json:"colors"`
	Labels []map[string]string `json:"labels"`
}

type AppModuleCount struct {
	Contents []map[string]types.InfoItem `json:"contents"`
	Headers  types.Thead                 `json:"headers"`
}

type AppDashboardReport struct {
	APITypeCount        BaseCount      `json:"api_type_count"`
	APISpecCount        BaseCount      `json:"api_spec_count"`
	AutoAPICount        BaseCount      `json:"auto_api_count"`
	APIRunResultCount   DayRunResult   `json:"api_run_result_count"`
	DaysAPIResultCount  BaseCount      `json:"days_api_result_count"`
	AppModuleTableCount AppModuleCount `json:"app_module_table_count"`
}

type ProductDashboardReport struct {
	APITypeCount               BaseCount      `json:"api_type_count"`
	APISpecCount               BaseCount      `json:"api_spec_count"`
	AutoAPICount               BaseCount      `json:"auto_api_count"`
	APIRunResultCount          DayRunResult   `json:"api_run_result_count"`
	DaysAPIResultCount         BaseCount      `json:"days_api_result_count"`
	ProductPlaybookResultCount DayRunResult   `json:"product_playbook_result_count"`
	DaysSceneResultCount       BaseCount      `json:"days_scene_result_count"`
	ProductAppModuleTableCount AppModuleCount `json:"product_app_module_table_count"`
}

type GlobalDashboardReport struct {
	APITypeCount               BaseCount      `json:"api_type_count"`
	APISpecCount               BaseCount      `json:"api_spec_count"`
	AutoAPICount               BaseCount      `json:"auto_api_count"`
	AppTestDataRunCount        DayRunResult   `json:"app_test_data_run_count"`
	AppAPIRunCount             BaseCount      `json:"app_api_run_count"`
	ProductPlaybookResultCount DayRunResult   `json:"product_playbook_result_count"`
	ProductSceneRunCount       BaseCount      `json:"product_scene_run_count"`
	ProductsTableCount         AppModuleCount `json:"products_table_count"`
	PlaybookResultCount        BaseCount      `json:"playbook_result_count"`
	TestDataResultCount        BaseCount      `json:"test_data_result_count"`
	ScheduleResultCount        BaseCount      `json:"schedule_result_count"`
}

type TaskDashboardReport struct {
	APITypeCount               BaseCount      `json:"api_type_count"`
	APISpecCount               BaseCount      `json:"api_spec_count"`
	AutoAPICount               BaseCount      `json:"auto_api_count"`
	APIRunResultCount          DayRunResult   `json:"api_run_result_count"`
	DaysAPIResultCount         BaseCount      `json:"days_api_result_count"`
	ProductPlaybookResultCount DayRunResult   `json:"product_playbook_result_count"`
	DaysSceneResultCount       BaseCount      `json:"days_scene_result_count"`
	ProductAppModuleTableCount AppModuleCount `json:"product_app_module_table_count"`
}

type DayRunResult struct {
	Title   template.HTML `json:"title"`
	DayList []string      `json:"day_list"`
	Infos   []string      `json:"infos"`
	Counts  [][]float64   `json:"counts"`
}
