package pages

import (
	"data4test/biz"
	"data4test/models"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/GoAdminGroup/go-admin/template/chartjs"

	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/chart_legend"
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
	components := tmpl.Default()
	colComp := components.Col()

	apiTypeDistribution := template.HTML(biz.T("dashboard.api_type_distribution"))
	viewAllApis := template.HTML(biz.T("dashboard.view_all_apis"))
	apiSpecCheck := template.HTML(biz.T("dashboard.api_spec_check"))
	apiIsAutomation := template.HTML(biz.T("dashboard.api_is_automation"))
	viewAllDataFiles := template.HTML(biz.T("dashboard.view_all_data_files"))
	dataExecTrend30d := template.HTML(biz.T("dashboard.data_exec_trend_30d"))
	dataExecDist30d := template.HTML(biz.T("dashboard.data_exec_dist_30d"))
	viewAllHistoryData := template.HTML(biz.T("dashboard.view_all_history_data"))
	moduleApiOverview := template.HTML(biz.T("dashboard.module_api_overview"))
	appDashboardTitle := template.HTML(biz.T("app_dashboard.title"))
	appDashboardDesc := template.HTML(biz.T("app_dashboard.description"))

	// 报表生成时间头部
	renderGenTime := report.CreatedAt
	//if len(renderGenTime) == 0 {
	//	renderGenTime = "未知"
	//}
	//renderHeaderHtml := fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-default" style="border-top:none;box-shadow:none;background:#fff;margin-bottom:10px"><div class="box-body" style="padding:10px 20px"><strong>应用:</strong> %s | <strong>报表生成时间:</strong> %s</div></div></div></div>`, appName, renderGenTime)

	apiMethods := appReportData.APITypeCount.Infos
	apiCounts := appReportData.APITypeCount.Counts
	colors := appReportData.APITypeCount.Colors
	labels := appReportData.APITypeCount.Labels
	pie1 := chartjs.Pie().
		SetHeight(120).
		SetLabels(apiMethods).
		SetID("pieChart1").
		AddDataSet(apiMethods[0]).
		DSData(apiCounts).
		DSBackgroundColor(colors).
		GetContent()
	legend1 := chart_legend.New().SetData(labels).GetContent()
	boxDanger1 := components.Box().SetTheme("danger1").WithHeadBorder().SetHeader(apiTypeDistribution).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie1).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend1).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/api_definition" class="uppercase1">` + viewAllApis + `</a></p>`).
		GetContent()
	col1 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger1).GetContent()
	specInfos := appReportData.APISpecCount.Infos
	specCounts := appReportData.APITypeCount.Counts
	colors2 := appReportData.APITypeCount.Colors
	labels2 := appReportData.APITypeCount.Labels
	pie2 := chartjs.Pie().
		SetHeight(120).
		SetLabels(specInfos).
		SetID("pieChart2").
		AddDataSet(specInfos[0]).
		DSData(specCounts).
		DSBackgroundColor(colors2).
		GetContent()
	legend2 := chart_legend.New().SetData(labels2).GetContent()
	boxDanger2 := components.Box().SetTheme("danger2").WithHeadBorder().SetHeader(apiSpecCheck).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie2).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend2).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/api_definition" class="uppercase2">` + viewAllApis + `</a></p>`).
		GetContent()
	col2 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger2).GetContent()

	autoInfos := appReportData.AutoAPICount.Infos
	autoCounts := appReportData.AutoAPICount.Counts
	colors3 := appReportData.AutoAPICount.Colors
	labels3 := appReportData.AutoAPICount.Labels
	pie3 := chartjs.Pie().
		SetHeight(120).
		SetLabels(autoInfos).
		SetID("pieChart3").
		AddDataSet(autoInfos[0]).
		DSData(autoCounts).
		DSBackgroundColor(colors3).
		GetContent()
	legend3 := chart_legend.New().SetData(labels3).GetContent()
	boxDanger3 := components.Box().SetTheme("danger3").WithHeadBorder().SetHeader(apiIsAutomation).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie3).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend3).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">` + viewAllDataFiles + `</a></p>`).
		GetContent()
	col3 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger3).GetContent()
	row1 := components.Row().SetContent(col1 + col2 + col3).GetContent()

	trendTitle := appReportData.APIRunResultCount.Title
	trendInfos := appReportData.APIRunResultCount.Infos
	trendCounts := appReportData.APIRunResultCount.Counts
	dayLabels := appReportData.APIRunResultCount.DayList
	line1 := chartjs.Line()
	lineChart1 := line1.
		SetID("dataChart").
		SetHeight(320).
		SetTitle(trendTitle).
		SetLabels(dayLabels)
	if len(trendInfos) > 0 && len(trendCounts) > 0 {
		lineChart1 = lineChart1.AddDataSet(trendInfos[0]).
			DSData(trendCounts[0]).
			DSFill(false).
			DSBorderColor(chartjs.Color("rgb(255, 205, 86)")).
			DSLineTension(0.1)
	}
	if len(trendInfos) > 1 && len(trendCounts) > 1 {
		lineChart1 = lineChart1.AddDataSet(trendInfos[1]).
			DSData(trendCounts[1]).
			DSFill(false).
			DSBorderColor(chartjs.Color("rgb(54, 162, 235)")).
			DSLineTension(0.1)
	}
	if len(trendInfos) > 2 && len(trendCounts) > 2 {
		lineChart1 = lineChart1.AddDataSet(trendInfos[2]).
			DSData(trendCounts[2]).
			DSFill(true).
			DSBorderColor(chartjs.Color("rgba(60,141,188,1)")).
			DSLineTension(0.1)
	}
	lineChartHTML := lineChart1.GetContent()
	boxInternalCol1 := colComp.SetContent(lineChartHTML).SetSize(types.SizeMD(12)).GetContent()
	boxInternalRow1 := components.Row().SetContent(boxInternalCol1).GetContent()
	box1 := components.Box().WithHeadBorder().SetHeader(dataExecTrend30d).
		SetBody(boxInternalRow1).
		GetContent()
	boxcol1 := colComp.SetContent(box1).SetSize(types.SizeMD(8)).GetContent()

	distInfos := appReportData.DaysAPIResultCount.Infos
	distCounts := appReportData.DaysAPIResultCount.Counts
	colors4 := appReportData.DaysAPIResultCount.Colors
	labels4 := appReportData.DaysAPIResultCount.Labels
	pie4 := chartjs.Pie().
		SetHeight(120).
		SetLabels(distInfos).
		SetID("pieChart4").
		AddDataSet(distInfos[0]).
		DSData(distCounts).
		DSBackgroundColor(colors4).
		GetContent()
	legend4 := chart_legend.New().SetData(labels4).GetContent()
	boxDanger4 := components.Box().SetTheme("danger4").WithHeadBorder().SetHeader(dataExecDist30d).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie4).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend4).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data_test_history" class="uppercase3">` + viewAllHistoryData + `</a></p>`).
		GetContent()
	col4 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger4).GetContent()
	row2 := components.Row().SetContent(boxcol1 + col4).GetContent()

	contents := appReportData.AppModuleTableCount.Contents
	headers := appReportData.AppModuleTableCount.Headers
	table := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
	row4 := components.Box().
		WithHeadBorder().
		SetHeader(moduleApiOverview).
		SetHeadColor("#f7f7f7").
		SetBody(table).
		GetContent()
	desc := template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s - %s</span><span style="color:#888">%s: %s</span></div>`, string(appDashboardDesc), appName, biz.T("schedule_report.generated_at"), renderGenTime))

	return types.Panel{
		Content:     row1 + row2 + row4,
		Title:       appDashboardTitle,
		Description: desc,
	}, nil

}
