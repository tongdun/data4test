package pages

import (
	"data4test/biz"
	"data4test/models"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"
)

func GetGlobalReportData(userName string) (err error) {
	var globalDR GlobalDashboardReport

	globalDR.APITypeCount.Infos, globalDR.APITypeCount.Counts, globalDR.APITypeCount.Colors, globalDR.APITypeCount.Labels = biz.GetAPITypeCount("all", "")
	globalDR.APISpecCount.Infos, globalDR.APISpecCount.Counts, globalDR.APISpecCount.Colors, globalDR.APISpecCount.Labels = biz.GetAPISpecCount("all", "")
	globalDR.AutoAPICount.Infos, globalDR.AutoAPICount.Counts, globalDR.AutoAPICount.Colors, globalDR.AutoAPICount.Labels = biz.GetAutoAPICount("all", "")
	globalDR.AppTestDataRunCount.Title, globalDR.AppTestDataRunCount.DayList, globalDR.AppTestDataRunCount.Infos, globalDR.AppTestDataRunCount.Counts = biz.GetAppSceneDataRunCount()
	globalDR.AppAPIRunCount.Infos, globalDR.AppAPIRunCount.Counts, globalDR.AppAPIRunCount.Colors, globalDR.AppAPIRunCount.Labels = biz.GetAppAPIRunCount()
	globalDR.ProductPlaybookResultCount.Title, globalDR.ProductPlaybookResultCount.DayList, globalDR.ProductPlaybookResultCount.Infos, globalDR.ProductPlaybookResultCount.Counts = biz.GetProductSceneRunCount()
	globalDR.ProductsTableCount.Contents, globalDR.ProductsTableCount.Headers = biz.GetProductsTableCount()
	globalDR.AppTableCount.Contents, globalDR.AppTableCount.Headers = biz.GetAppsTableCount()
	globalDR.PlaybookResultCount.Infos, globalDR.PlaybookResultCount.Counts, globalDR.PlaybookResultCount.Colors, globalDR.PlaybookResultCount.Labels = biz.GetSceneResultCount()
	globalDR.TestDataResultCount.Infos, globalDR.TestDataResultCount.Counts, globalDR.TestDataResultCount.Colors, globalDR.TestDataResultCount.Labels = biz.GetSceneDataResultCount()
	globalDR.ScheduleResultCount.Infos, globalDR.ScheduleResultCount.Counts, globalDR.ScheduleResultCount.Colors, globalDR.ScheduleResultCount.Labels = biz.GetScheduleTypeCount()

	globalDRStr, _ := json.MarshalIndent(globalDR, "", "    ")

	nowStr := time.Now().Format("20060102150405")
	now := time.Now().Format("2006-01-02 15:04:05")

	// 从关联数据执行历史中获取实际时间区间作为统计时间
	var startTime, endTime string
	models.Orm.Table("scene_data_test_history").
		Select("MIN(created_at) as start_time, MAX(created_at) as end_time").
		Row().Scan(&startTime, &endTime)
	if len(startTime) == 0 || len(endTime) == 0 {
		models.Orm.Table("scene_test_history").
			Select("MIN(created_at) as start_time, MAX(created_at) as end_time").
			Row().Scan(&startTime, &endTime)
	}
	if len(startTime) == 0 || len(endTime) == 0 {
		startTime = now
		endTime = now
	}

	reportName := fmt.Sprintf("%s_%s", biz.T("dashboard.global_report"), nowStr)
	report := biz.DashboardReport{
		ReportName:      reportName,
		ReportType:      "global",
		RelatedProducts: "",
		RelatedApps:     "",
		TimeRangeStart:  startTime,
		TimeRangeEnd:    endTime,
		Status:          "finished",
		Creator:         userName,
		ReportData:      string(globalDRStr),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = models.Orm.Table("dashboard").Create(&report).Error
	if err != nil {
		return fmt.Errorf("保存报告失败: %s", err)
	}
	return err
}

func GetDashBoardContent(ctx *gin.Context, userName string) (types.Panel, error) {
	// 优先从 dashboard 表读取预计算数据
	var globalReportData GlobalDashboardReport
	var dr biz.DashboardReport
	models.Orm.Table("dashboard").
		Where("report_type = ?", "global").
		Order("created_at desc").
		Limit(1).
		Find(&dr)

	if len(dr.Id) > 0 {
		err := json.Unmarshal([]byte(dr.ReportData), &globalReportData)
		if err != nil {
			biz.Logger.Error("err: %v", err)
			return types.Panel{}, err
		}
	} else {
		err := GetGlobalReportData(userName)
		if err != nil {
			biz.Logger.Error("err: %v", err)
			return types.Panel{}, err
		}
		models.Orm.Table("dashboard").
			Where("report_type = ?", "global").
			Order("created_at desc").
			Limit(1).
			Find(&dr)

		if len(dr.Id) > 0 {
			err = json.Unmarshal([]byte(dr.ReportData), &globalReportData)
			if err != nil {
				biz.Logger.Error("err: %v", err)
				biz.Logger.Debug("%+v", globalReportData)
			}
		}
	}
	return renderGlobalReport(globalReportData, dr)
}

// globalPie 大盘饼图（col-md-4 一栏，含可选 footer 链接）。
func globalPie(id, title string, data BaseCount, footerURL, footerText string) string {
	if len(data.Infos) == 0 {
		return `<div class="col-md-4">` + reportCard(title, "", "") + `</div>`
	}
	infos := localizeInfos(data.Infos)
	pie := reportPieChart(id, infos, data.Counts, data.Colors)
	legend := reportLegend(rebuildLegendLabels(data.Labels, infos, data.Counts))
	body := fmt.Sprintf(`<div class="row"><div class="col-md-8">%s</div><div class="col-md-4">%s</div></div>`, pie, legend)
	if footerURL != "" {
		body += fmt.Sprintf(`<div class="report-card-footer"><a href="%s">%s &raquo;</a></div>`, footerURL, footerText)
	}
	return `<div class="col-md-4">` + reportCard(title, "", body) + `</div>`
}

// globalLine 大盘折线图（按 DayRunResult 系列循环取色）。
// X 轴标签 DayList 存的是 i18n key（如 report.month_N）或日期字符串，渲染时本地化。
func globalLine(id, cardTitle string, dr DayRunResult) string {
	if len(dr.Infos) == 0 {
		return reportCard(cardTitle, "", "")
	}
	labels := localizeInfos(dr.DayList)
	series := make([]reportLineSeries, 0, len(dr.Infos))
	for i := range dr.Infos {
		var data []float64
		if i < len(dr.Counts) {
			data = dr.Counts[i]
		}
		series = append(series, reportLineSeries{
			Label: dr.Infos[i],
			Data:  data,
			Color: categoricalColor(i),
		})
	}
	return reportLine(id, cardTitle, labels, series)
}

func renderGlobalReport(globalReportData GlobalDashboardReport, report biz.DashboardReport) (types.Panel, error) {
	// 第1行：API 类型分布 / API 规范 / API 自动化
	row1 := `<div class="row">` +
		globalPie("pieChart1", biz.T("dashboard.api_type_distribution"), globalReportData.APITypeCount, "/admin/info/api_definition", biz.T("dashboard.view_all_apis")) +
		globalPie("pieChart2", biz.T("dashboard.api_spec_check"), globalReportData.APISpecCount, "/admin/info/api_definition", biz.T("dashboard.view_all_apis")) +
		globalPie("pieChart3", biz.T("dashboard.api_is_automation"), globalReportData.AutoAPICount, "/admin/info/scene_data", biz.T("dashboard.view_all_api_data")) +
		`</div>`

	// 第2行：数据执行趋势（折线）+ 数据执行分布（饼）
	line1 := globalLine("dataChart", biz.T("dashboard.data_exec_trend_6m"), globalReportData.AppTestDataRunCount)
	pie4 := globalPie("pieChart4", biz.T("dashboard.data_exec_dist_6m"), globalReportData.AppAPIRunCount, "/admin/info/scene_data_test_history", biz.T("dashboard.view_all_history_data"))
	row2 := `<div class="row"><div class="col-md-8">` + line1 + `</div>` + pie4 + `</div>`

	// 第3行：场景执行趋势（折线）+ 场景执行分布（饼）
	line2 := globalLine("sceneChart", biz.T("dashboard.scene_exec_trend_6m"), globalReportData.ProductPlaybookResultCount)
	pie5 := globalPie("pieChart5", biz.T("dashboard.scene_exec_dist_6m"), globalReportData.PlaybookResultCount, "/admin/info/sence_test_history", biz.T("dashboard.view_all_history_scene"))
	row3 := `<div class="row"><div class="col-md-8">` + line2 + `</div>` + pie5 + `</div>`

	// 第6行：场景执行状态 / 数据执行状态 / 任务类型分布
	row6 := `<div class="row">` +
		globalPie("pieChart6", biz.T("dashboard.scene_exec_status"), globalReportData.PlaybookResultCount, "/admin/info/scene_test_history", biz.T("dashboard.view_all_scene_records")) +
		globalPie("pieChart7", biz.T("dashboard.data_exec_status"), globalReportData.TestDataResultCount, "/admin/info/scene_data", biz.T("dashboard.view_all_data_records")) +
		globalPie("pieChart8", biz.T("dashboard.task_type_dist"), globalReportData.ScheduleResultCount, "/admin/info/schedule", biz.T("dashboard.view_all_tasks")) +
		`</div>`

	// 产品列表 / 应用列表
	productsBody := reportTableFromInfo(globalReportData.ProductsTableCount.Contents, globalReportData.ProductsTableCount.Headers) +
		fmt.Sprintf(`<div class="report-card-footer"><a href="/admin/info/product">%s &raquo;</a></div>`, biz.T("dashboard.goto_product_detail"))
	row4 := reportCard(biz.T("common.product_list"), "", productsBody)

	appsBody := reportTableFromInfo(globalReportData.AppTableCount.Contents, globalReportData.AppTableCount.Headers) +
		fmt.Sprintf(`<div class="report-card-footer"><a href="/admin/info/env_config">%s &raquo;</a></div>`, biz.T("dashboard.goto_app_detail"))
	row5 := reportCard(biz.T("dashboard.app_list"), "", appsBody)

	content := row1 + row2 + row3 + row6 + row5 + row4 + string(reportStyle())

	return types.Panel{
		Content:     template.HTML(content),
		Title:       template.HTML(biz.T("dashboard.title")),
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s</span><span style="color:#888">%s: %s</span></div>`, biz.T("dashboard.description"), biz.T("schedule_report.generated_at"), report.CreatedAt)),
	}, nil
}
