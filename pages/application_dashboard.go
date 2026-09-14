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

func GetAppReportData(appName, userName string) (err error) {
	var appDR AppDashboardReport
	appDR.APITypeCount.Infos, appDR.APITypeCount.Counts, appDR.APITypeCount.Colors, appDR.APITypeCount.Labels = biz.GetAPITypeCount("app", appName)
	appDR.APISpecCount.Infos, appDR.APISpecCount.Counts, appDR.APISpecCount.Colors, appDR.APISpecCount.Labels = biz.GetAPISpecCount("app", appName)
	appDR.AutoAPICount.Infos, appDR.AutoAPICount.Counts, appDR.AutoAPICount.Colors, appDR.AutoAPICount.Labels = biz.GetAutoAPICount("app", appName)
	appDR.APIRunResultCount.Title, appDR.APIRunResultCount.DayList, appDR.APIRunResultCount.Infos, appDR.APIRunResultCount.Counts = biz.GetAPIRunResultCount("app", appName)
	appDR.DaysAPIResultCount.Infos, appDR.DaysAPIResultCount.Counts, appDR.DaysAPIResultCount.Colors, appDR.DaysAPIResultCount.Labels = biz.GetDaysAPIResultCount("app", appName, 30)
	appDR.AppModuleTableCount.Contents, appDR.AppModuleTableCount.Headers = biz.GetAppModuleTableCount(appName)

	appDRStr, _ := json.MarshalIndent(appDR, "", "    ")

	nowStr := time.Now().Format("20060102150405")
	now := time.Now().Format("2006-01-02 15:04:05")

	// 从关联数据执行历史中获取实际时间区间作为统计时间
	var startTime, endTime string
	models.Orm.Table("scene_data_test_history").
		Select("MIN(created_at) as start_time, MAX(created_at) as end_time").
		Where("app = ?", appName).
		Row().Scan(&startTime, &endTime)
	if len(startTime) == 0 || len(endTime) == 0 {
		startTime = now
		endTime = now
	}

	reportName := biz.T("app.report_name", appName, nowStr)
	report := biz.DashboardReport{
		ReportName:      reportName,
		ReportType:      "app",
		RelatedProducts: "",
		RelatedApps:     appName,
		TimeRangeStart:  startTime,
		TimeRangeEnd:    endTime,
		Status:          "finished",
		Creator:         userName,
		ReportData:      string(appDRStr),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = models.Orm.Table("dashboard").Create(&report).Error
	if err != nil {
		return fmt.Errorf(biz.T("product.save_report_failed"), err)
	}
	return err
}

func GetDashBoard3Content(ctx *gin.Context, userName string) (types.Panel, error) {
	id := ctx.Query("id")
	appName, err := biz.GetAppName(id)
	if err != nil {
		return types.Panel{}, err
	}

	// 优先从 dashboard 表读取预计算数据
	var appReportData AppDashboardReport
	var dr biz.DashboardReport
	err = models.Orm.Table("dashboard").
		Where("related_apps =  ?", appName).
		Order("created_at desc").
		Limit(1).
		Find(&dr).Error

	if len(dr.Id) > 0 {
		err := json.Unmarshal([]byte(dr.ReportData), &appReportData)
		if err != nil {
			biz.Logger.Error("err: %v", err)
			biz.Logger.Debug("%v", dr.ReportData)
		}
	} else {
		err = GetAppReportData(appName, userName)
		if err != nil {
			biz.Logger.Error("err: %v", err)
			return types.Panel{}, err
		}
		err = models.Orm.Table("dashboard").
			Where("related_apps =  ?", appName).
			Order("created_at desc").
			Limit(1).
			Find(&dr).Error

		if len(dr.Id) > 0 {
			err = json.Unmarshal([]byte(dr.ReportData), &appReportData)
			if err != nil {
				biz.Logger.Error("err: %v", err)
				biz.Logger.Debug("%+v", appReportData)
			}
		}
	}
	return renderApplicationReport(appName, dr)
}

func renderApplicationReport(appName string, report biz.DashboardReport) (types.Panel, error) {
	var appReportData AppDashboardReport
	err := json.Unmarshal([]byte(report.ReportData), &appReportData)
	if err != nil {
		return types.Panel{
			Content:     template.HTML(fmt.Sprintf("<div style='padding:20px;color:red'>%s</div>", biz.T("schedule_report.parse_error", err))),
			Title:       template.HTML(biz.T("schedule_report.page_title")),
			Description: template.HTML(biz.T("schedule_report.description")),
		}, nil
	}

	// 第1行：API 类型分布 / API 规范 / API 自动化
	pie2Data := BaseCount{
		Infos:  appReportData.APISpecCount.Infos,
		Counts: appReportData.APITypeCount.Counts,
		Colors: appReportData.APITypeCount.Colors,
		Labels: appReportData.APITypeCount.Labels,
	}
	row1 := `<div class="row">` +
		globalPie("pieChart1", biz.T("dashboard.api_type_distribution"), appReportData.APITypeCount, "/admin/info/api_definition", biz.T("dashboard.view_all_apis")) +
		globalPie("pieChart2", biz.T("dashboard.api_spec_check"), pie2Data, "/admin/info/api_definition", biz.T("dashboard.view_all_apis")) +
		globalPie("pieChart3", biz.T("dashboard.api_is_automation"), appReportData.AutoAPICount, "/admin/info/scene_data", biz.T("dashboard.view_all_data_files")) +
		`</div>`

	// 第2行：数据执行趋势（折线）+ 数据执行分布（饼）
	line1 := globalLine("dataChart", biz.T("dashboard.data_exec_trend_30d"), appReportData.APIRunResultCount)
	pie4 := globalPie("pieChart4", biz.T("dashboard.data_exec_dist_30d"), appReportData.DaysAPIResultCount, "/admin/info/scene_data_test_history", biz.T("dashboard.view_all_history_data"))
	row2 := `<div class="row"><div class="col-md-8">` + line1 + `</div>` + pie4 + `</div>`

	// 模块 API 概览表
	row4 := reportCard(biz.T("dashboard.module_api_overview"), "", reportTableFromInfo(appReportData.AppModuleTableCount.Contents, appReportData.AppModuleTableCount.Headers))

	desc := template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s - %s</span><span style="color:#888">%s: %s</span></div>`, biz.T("app_dashboard.description"), appName, biz.T("schedule_report.generated_at"), report.CreatedAt))

	content := row1 + row2 + row4 + string(reportStyle())

	return types.Panel{
		Content:     template.HTML(content),
		Title:       template.HTML(biz.T("app_dashboard.title")),
		Description: desc,
	}, nil
}
