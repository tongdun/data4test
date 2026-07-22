package pages

import (
	"data4test/biz"
	"data4test/models"
	"encoding/json"
	"fmt"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"html/template"
	"time"

	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/chart_legend"
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
	globalDR.PlaybookResultCount.Infos, globalDR.PlaybookResultCount.Counts, globalDR.PlaybookResultCount.Colors, globalDR.PlaybookResultCount.Labels = biz.GetSceneRunCount()
	globalDR.ProductsTableCount.Contents, globalDR.ProductsTableCount.Headers = biz.GetProductsTableCount()
	globalDR.PlaybookResultCount.Infos, globalDR.PlaybookResultCount.Counts, globalDR.PlaybookResultCount.Colors, globalDR.PlaybookResultCount.Labels = biz.GetSceneResultCount()
	globalDR.TestDataResultCount.Infos, globalDR.TestDataResultCount.Counts, globalDR.TestDataResultCount.Colors, globalDR.TestDataResultCount.Labels = biz.GetSceneDataResultCount()
	globalDR.ScheduleResultCount.Infos, globalDR.ScheduleResultCount.Counts, globalDR.ScheduleResultCount.Colors, globalDR.ScheduleResultCount.Labels = biz.GetScheduleTypeCount()

	globalDRStr, _ := json.MarshalIndent(globalDR, "", "    ")

	nowStr := time.Now().Format("20060102150405")
	now := time.Now().Format("2006-01-02 15:04:05")
	reportName := fmt.Sprintf("%s报告_%s", "全局", nowStr)
	report := biz.DashboardReport{
		ReportName:      reportName,
		ReportType:      "global",
		RelatedProducts: "",
		RelatedApps:     "",
		TimeRangeStart:  now,
		TimeRangeEnd:    now,
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

//func GetDashBoardContent(ctx *gin.Context, userName string) (types.Panel, error) {
//	// 优先从 dashboard 表读取预计算数据
//	var globalReportData GlobalDashboardReport
//	var dr biz.DashboardReport
//	models.Orm.Table("dashboard").
//		Where("report_type = ?", "global").
//		Order("created_at desc").
//		Limit(1).
//		Find(&dr)
//
//	if len(dr.Id) > 0 {
//		err := json.Unmarshal([]byte(dr.ReportData), &globalReportData)
//		if err != nil {
//			biz.Logger.Error("err: %v", err)
//			biz.Logger.Debug("%v", dr.ReportData)
//		}
//	} else {
//		err := GetGlobalReportData(userName)
//		if err != nil {
//			biz.Logger.Error("err: %v", err)
//			return types.Panel{}, err
//		}
//		models.Orm.Table("dashboard").
//			Where("report_type = ?", "global").
//			Order("created_at desc").
//			Limit(1).
//			Find(&dr)
//
//		if len(dr.Id) > 0 {
//			err = json.Unmarshal([]byte(dr.ReportData), &globalReportData)
//			if err != nil {
//				biz.Logger.Error("err: %v", err)
//				biz.Logger.Debug("%+v", globalReportData)
//			}
//		}
//	}
//	components := tmpl.Default()
//	colComp := components.Col()
//	apiMethods := globalReportData.APITypeCount.Infos
//	apiCounts := globalReportData.APITypeCount.Counts
//	colors := globalReportData.APITypeCount.Colors
//	labels := globalReportData.APITypeCount.Labels
//	pie1 := chartjs.Pie().
//		SetHeight(120).
//		SetLabels(apiMethods).
//		SetID("pieChart1").
//		AddDataSet(apiMethods[0]).
//		DSData(apiCounts).
//		DSBackgroundColor(colors).
//		GetContent()
//
//	legend1 := chart_legend.New().SetData(labels).GetContent()
//
//	apiTypeDistribution := template.HTML(biz.T("dashboard.api_type_distribution"))
//	viewAllApis := template.HTML(biz.T("dashboard.view_all_apis"))
//	apiSpecCheck := template.HTML(biz.T("dashboard.api_spec_check"))
//	apiIsAutomation := template.HTML(biz.T("dashboard.api_is_automation"))
//	viewAllApiData := template.HTML(biz.T("dashboard.view_all_api_data"))
//	dataExecTrend6m := template.HTML(biz.T("dashboard.data_exec_trend_6m"))
//	dataExecDist6m := template.HTML(biz.T("dashboard.data_exec_dist_6m"))
//	viewAllHistoryData := template.HTML(biz.T("dashboard.view_all_history_data"))
//	sceneExecTrend6m := template.HTML(biz.T("dashboard.scene_exec_trend_6m"))
//	sceneExecDist6m := template.HTML(biz.T("dashboard.scene_exec_dist_6m"))
//	viewAllHistoryScene := template.HTML(biz.T("dashboard.view_all_history_scene"))
//	productList := template.HTML(biz.T("common.product_list"))
//	gotoProductDetail := template.HTML(biz.T("dashboard.goto_product_detail"))
//	appList := template.HTML(biz.T("dashboard.app_list"))
//	gotoAppDetail := template.HTML(biz.T("dashboard.goto_app_detail"))
//	sceneExecStatus := template.HTML(biz.T("dashboard.scene_exec_status"))
//	viewAllSceneRecords := template.HTML(biz.T("dashboard.view_all_scene_records"))
//	dataExecStatus := template.HTML(biz.T("dashboard.data_exec_status"))
//	viewAllDataRecords := template.HTML(biz.T("dashboard.view_all_data_records"))
//	taskTypeDist := template.HTML(biz.T("dashboard.task_type_dist"))
//	viewAllTasks := template.HTML(biz.T("dashboard.view_all_tasks"))
//	dashboardTitle := template.HTML(biz.T("dashboard.title"))
//	dashboardDesc := template.HTML(biz.T("dashboard.description"))
//
//	boxDanger1 := components.Box().SetTheme("danger1").WithHeadBorder().SetHeader(apiTypeDistribution).
//		SetBody(components.Row().
//			SetContent(colComp.SetSize(types.SizeMD(8)).
//				SetContent(pie1).
//				GetContent() + colComp.SetSize(types.SizeMD(4)).
//				SetContent(legend1).
//				GetContent()).GetContent()).
//		SetFooter(`<p class="text-center"><a href="/admin/info/api_definition" class="uppercase1">` + viewAllApis + `</a></p>`).
//		GetContent()
//	col1 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger1).GetContent()
//	infos := globalReportData.APISpecCount.Infos
//	counts := globalReportData.APITypeCount.Counts
//	colors = globalReportData.APITypeCount.Colors
//	labels = globalReportData.APITypeCount.Labels
//	pie2 := chartjs.Pie().
//		SetHeight(120).
//		SetLabels(infos).
//		SetID("pieChart2").
//		AddDataSet(infos[0]).
//		DSData(counts).
//		DSBackgroundColor(colors).
//		GetContent()
//
//	legend2 := chart_legend.New().SetData(labels).GetContent()
//
//	boxDanger2 := components.Box().SetTheme("danger2").WithHeadBorder().SetHeader(apiSpecCheck).
//		SetBody(components.Row().
//			SetContent(colComp.SetSize(types.SizeMD(8)).
//				SetContent(pie2).
//				GetContent() + colComp.SetSize(types.SizeMD(4)).
//				SetContent(legend2).
//				GetContent()).GetContent()).
//		SetFooter(`<p class="text-center"><a href="/admin/info/api_definition" class="uppercase2">` + viewAllApis + `</a></p>`).
//		GetContent()
//
//	col2 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger2).GetContent()
//	infos = globalReportData.AutoAPICount.Infos
//	counts = globalReportData.AutoAPICount.Counts
//	colors = globalReportData.AutoAPICount.Colors
//	labels = globalReportData.AutoAPICount.Labels
//	pie3 := chartjs.Pie().
//		SetHeight(120).
//		SetLabels(infos).
//		SetID("pieChart3").
//		AddDataSet(infos[0]).
//		DSData(counts).
//		DSBackgroundColor(colors).
//		GetContent()
//
//	legend3 := chart_legend.New().SetData(labels).GetContent()
//
//	boxDanger3 := components.Box().SetTheme("danger3").WithHeadBorder().SetHeader(apiIsAutomation).
//		SetBody(components.Row().
//			SetContent(colComp.SetSize(types.SizeMD(8)).
//				SetContent(pie3).
//				GetContent() + colComp.SetSize(types.SizeMD(4)).
//				SetContent(legend3).
//				GetContent()).GetContent()).
//		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">` + viewAllApiData + `</a></p>`).
//		GetContent()
//
//	col3 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger3).GetContent()
//
//	row1 := components.Row().SetContent(col1 + col2 + col3).GetContent()
//
//	line1 := chartjs.Line()
//	title := globalReportData.AppTestDataRunCount.Title
//	infos = globalReportData.AppTestDataRunCount.Infos
//	monthCounts := globalReportData.AppTestDataRunCount.Counts
//	monthLable := globalReportData.AppTestDataRunCount.DayList
//	var lineChart1 template.HTML
//	if len(infos) == 1 {
//		lineChart1 = line1.
//			SetID("dataChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	} else if len(infos) == 2 {
//		lineChart1 = line1.
//			SetID("dataChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgb(255, 205, 86)").
//			DSLineTension(0.1).
//			AddDataSet(infos[1]).
//			DSData(monthCounts[1]).
//			DSFill(false).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	} else if len(infos) == 3 {
//		lineChart1 = line1.
//			SetID("dataChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgb(255, 205, 86)").
//			DSLineTension(0.1).
//			AddDataSet(infos[1]).
//			DSData(monthCounts[1]).
//			DSFill(false).
//			DSBorderColor("rgb(54, 162, 235)").
//			DSLineTension(0.1).
//			AddDataSet(infos[2]).
//			DSData(monthCounts[2]).
//			DSFill(true).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	} else if len(infos) == 4 {
//		lineChart1 = line1.
//			SetID("dataChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgb(255, 205, 86)").
//			DSLineTension(0.1).
//			AddDataSet(infos[1]).
//			DSData(monthCounts[1]).
//			DSFill(false).
//			DSBorderColor("rgb(54, 162, 235)").
//			DSLineTension(0.1).
//			AddDataSet(infos[2]).
//			DSData(monthCounts[2]).
//			DSFill(false).
//			DSBorderColor("rgb(238,232,170)").
//			DSLineTension(0.1).
//			AddDataSet(infos[3]).
//			DSData(monthCounts[3]).
//			DSFill(false).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	} else if len(infos) >= 5 {
//		lineChart1 = line1.
//			SetID("dataChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgb(255, 205, 86)").
//			DSLineTension(0.1).
//			AddDataSet(infos[1]).
//			DSData(monthCounts[1]).
//			DSFill(false).
//			DSBorderColor("rgb(54, 162, 235)").
//			DSLineTension(0.1).
//			AddDataSet(infos[2]).
//			DSData(monthCounts[2]).
//			DSFill(false).
//			DSBorderColor("rgb(238,232,170)").
//			DSLineTension(0.1).
//			AddDataSet(infos[3]).
//			DSData(monthCounts[3]).
//			DSFill(false).
//			DSBorderColor("rgb(189,183,107)").
//			DSLineTension(0.1).
//			AddDataSet(infos[4]).
//			DSData(monthCounts[4]).
//			DSFill(false).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	}
//	boxInternalCol1 := colComp.SetContent(lineChart1).SetSize(types.SizeMD(12)).GetContent()
//	boxInternalRow1 := components.Row().SetContent(boxInternalCol1).GetContent()
//	box1 := components.Box().WithHeadBorder().SetHeader(dataExecTrend6m).
//		SetBody(boxInternalRow1).
//		GetContent()
//
//	boxcol1 := colComp.SetContent(box1).SetSize(types.SizeMD(8)).GetContent()
//	infos = globalReportData.AppAPIRunCount.Infos
//	counts = globalReportData.AppAPIRunCount.Counts
//	colors = globalReportData.AppAPIRunCount.Colors
//	labels = globalReportData.AppAPIRunCount.Labels
//	pie4 := chartjs.Pie().
//		SetHeight(120).
//		SetLabels(infos).
//		SetID("pieChart4").
//		AddDataSet(infos[0]).
//		DSData(counts).
//		DSBackgroundColor(colors).
//		GetContent()
//
//	legend4 := chart_legend.New().SetData(labels).GetContent()
//
//	boxDanger4 := components.Box().SetTheme("danger4").WithHeadBorder().SetHeader(dataExecDist6m).
//		SetBody(components.Row().
//			SetContent(colComp.SetSize(types.SizeMD(8)).
//				SetContent(pie4).
//				GetContent() + colComp.SetSize(types.SizeMD(4)).
//				SetContent(legend4).
//				GetContent()).GetContent()).
//		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data_test_history" class="uppercase3">` + viewAllHistoryData + `</a></p>`).
//		GetContent()
//
//	col4 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger4).GetContent()
//
//	row2 := components.Row().SetContent(boxcol1 + col4).GetContent()
//
//	line2 := chartjs.Line()
//	title = globalReportData.ProductPlaybookResultCount.Title
//	infos = globalReportData.ProductPlaybookResultCount.Infos
//	monthCounts = globalReportData.ProductPlaybookResultCount.Counts
//	monthLable = globalReportData.ProductPlaybookResultCount.DayList
//	var lineChart2 template.HTML
//	if len(infos) == 1 {
//		lineChart2 = line2.
//			SetID("sceneChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	} else if len(infos) == 2 {
//		lineChart2 = line2.
//			SetID("sceneChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgb(255, 205, 86)").
//			DSLineTension(0.1).
//			AddDataSet(infos[1]).
//			DSData(monthCounts[1]).
//			DSFill(false).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	} else if len(infos) == 3 {
//		lineChart2 = line2.
//			SetID("sceneChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgb(255, 205, 86)").
//			DSLineTension(0.1).
//			AddDataSet(infos[1]).
//			DSData(monthCounts[1]).
//			DSFill(false).
//			DSBorderColor("rgb(54, 162, 235)").
//			DSLineTension(0.1).
//			AddDataSet(infos[2]).
//			DSData(monthCounts[2]).
//			DSFill(true).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	} else if len(infos) == 4 {
//		lineChart2 = line2.
//			SetID("sceneChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgb(255, 205, 86)").
//			DSLineTension(0.1).
//			AddDataSet(infos[1]).
//			DSData(monthCounts[1]).
//			DSFill(false).
//			DSBorderColor("rgb(54, 162, 235)").
//			DSLineTension(0.1).
//			AddDataSet(infos[2]).
//			DSData(monthCounts[2]).
//			DSFill(false).
//			DSBorderColor("rgb(238,232,170)").
//			DSLineTension(0.1).
//			AddDataSet(infos[3]).
//			DSData(monthCounts[3]).
//			DSFill(false).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	} else if len(infos) >= 5 {
//		lineChart2 = line2.
//			SetID("sceneChart").
//			SetHeight(320).
//			SetTitle(title).
//			SetLabels(monthLable).
//			AddDataSet(infos[0]).
//			DSData(monthCounts[0]).
//			DSFill(false).
//			DSBorderColor("rgb(255, 205, 86)").
//			DSLineTension(0.1).
//			AddDataSet(infos[1]).
//			DSData(monthCounts[1]).
//			DSFill(false).
//			DSBorderColor("rgb(54, 162, 235)").
//			DSLineTension(0.1).
//			AddDataSet(infos[2]).
//			DSData(monthCounts[2]).
//			DSFill(false).
//			DSBorderColor("rgb(238,232,170)").
//			DSLineTension(0.1).
//			AddDataSet(infos[3]).
//			DSData(monthCounts[3]).
//			DSFill(false).
//			DSBorderColor("rgb(189,183,107)").
//			DSLineTension(0.1).
//			AddDataSet(infos[4]).
//			DSData(monthCounts[4]).
//			DSFill(false).
//			DSBorderColor("rgba(60,141,188,1)").
//			DSLineTension(0.1).
//			GetContent()
//	}
//	boxInternalCol2 := colComp.SetContent(lineChart2).SetSize(types.SizeMD(12)).GetContent()
//	boxInternalRow2 := components.Row().SetContent(boxInternalCol2).GetContent()
//	box2 := components.Box().WithHeadBorder().SetHeader(sceneExecTrend6m).
//		SetBody(boxInternalRow2).
//		GetContent()
//
//	boxcol2 := colComp.SetContent(box2).SetSize(types.SizeMD(8)).GetContent()
//	infos = globalReportData.PlaybookResultCount.Infos
//	counts = globalReportData.PlaybookResultCount.Counts
//	colors = globalReportData.PlaybookResultCount.Colors
//	labels = globalReportData.PlaybookResultCount.Labels
//	pie5 := chartjs.Pie().
//		SetHeight(120).
//		SetLabels(infos).
//		SetID("pieChart5").
//		AddDataSet(infos[0]).
//		DSData(counts).
//		DSBackgroundColor(colors).
//		GetContent()
//	legend5 := chart_legend.New().SetData(labels).GetContent()
//	boxDanger5 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(sceneExecDist6m).
//		SetBody(components.Row().
//			SetContent(colComp.SetSize(types.SizeMD(8)).
//				SetContent(pie5).
//				GetContent() + colComp.SetSize(types.SizeMD(4)).
//				SetContent(legend5).
//				GetContent()).GetContent()).
//		SetFooter(`<p class="text-center"><a href="/admin/info/sence_test_history" class="uppercase3">` + viewAllHistoryScene + `</a></p>`).
//		GetContent()
//	col5 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger5).GetContent()
//	row3 := components.Row().SetContent(boxcol2 + col5).GetContent()
//	contents := globalReportData.ProductsTableCount.Contents
//	headers := globalReportData.ProductsTableCount.Headers
//	table := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
//	boxInfo := components.Box().
//		WithHeadBorder().
//		SetHeader(productList).
//		SetHeadColor("#f7f7f7").
//		SetBody(table).
//		SetFooter(`<div class="clearfix"><a href="/admin/info/product" class="btn btn-sm btn-default btn-flat pull-right">` + gotoProductDetail + `</a> </div>`).
//		GetContent()
//	tableCol := colComp.SetSize(types.SizeMD(12)).SetContent(boxInfo).GetContent()
//	row4 := components.Row().SetContent(tableCol).GetContent()
//	var tableApp template.HTML
//	aContents, aHeaders := biz.GetAppsTableCount()
//	if len(aContents) > 0 {
//		tableApp = components.Table().SetInfoList(aContents).SetThead(aHeaders).GetContent()
//	}
//
//	boxAppInfo := components.Box().
//		WithHeadBorder().
//		SetHeader(appList).
//		SetHeadColor("#f7f7f7").
//		SetBody(tableApp).
//		SetFooter(`<div class="clearfix"><a href="/admin/info/env_config" class="btn btn-sm btn-default btn-flat pull-right">` + gotoAppDetail + `</a> </div>`).
//		GetContent()
//	tableAppCol := colComp.SetSize(types.SizeMD(12)).SetContent(boxAppInfo).GetContent()
//	row5 := components.Row().SetContent(tableAppCol).GetContent()
//	infos = globalReportData.PlaybookResultCount.Infos
//	counts = globalReportData.PlaybookResultCount.Counts
//	colors = globalReportData.PlaybookResultCount.Colors
//	labels = globalReportData.PlaybookResultCount.Labels
//	pie6 := chartjs.Pie().
//		SetHeight(120).
//		SetLabels(infos).
//		SetID("pieChart6").
//		AddDataSet(infos[0]).
//		DSData(counts).
//		DSBackgroundColor(colors).
//		GetContent()
//
//	legend6 := chart_legend.New().SetData(labels).GetContent()
//
//	boxDanger6 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(sceneExecStatus).
//		SetBody(components.Row().
//			SetContent(colComp.SetSize(types.SizeMD(8)).
//				SetContent(pie6).
//				GetContent() + colComp.SetSize(types.SizeMD(4)).
//				SetContent(legend6).
//				GetContent()).GetContent()).
//		SetFooter(`<p class="text-center"><a href="/admin/info/scene_test_history" class="uppercase3">` + viewAllSceneRecords + `</a></p>`).
//		GetContent()
//
//	col6 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger6).GetContent()
//	infos = globalReportData.TestDataResultCount.Infos
//	counts = globalReportData.TestDataResultCount.Counts
//	colors = globalReportData.TestDataResultCount.Colors
//	labels = globalReportData.TestDataResultCount.Labels
//	pie7 := chartjs.Pie().
//		SetHeight(120).
//		SetLabels(infos).
//		SetID("pieChart7").
//		AddDataSet(infos[0]).
//		DSData(counts).
//		DSBackgroundColor(colors).
//		GetContent()
//
//	legend7 := chart_legend.New().SetData(labels).GetContent()
//
//	boxDanger7 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(dataExecStatus).
//		SetBody(components.Row().
//			SetContent(colComp.SetSize(types.SizeMD(8)).
//				SetContent(pie7).
//				GetContent() + colComp.SetSize(types.SizeMD(4)).
//				SetContent(legend7).
//				GetContent()).GetContent()).
//		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">` + viewAllDataRecords + `</a></p>`).
//		GetContent()
//
//	col7 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger7).GetContent()
//	infos = globalReportData.ScheduleResultCount.Infos
//	counts = globalReportData.ScheduleResultCount.Counts
//	colors = globalReportData.ScheduleResultCount.Colors
//	labels = globalReportData.ScheduleResultCount.Labels
//	pie8 := chartjs.Pie().
//		SetHeight(120).
//		SetLabels(infos).
//		SetID("pieChart8").
//		AddDataSet(infos[0]).
//		DSData(counts).
//		DSBackgroundColor(colors).
//		GetContent()
//
//	legend8 := chart_legend.New().SetData(labels).GetContent()
//
//	boxDanger8 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(taskTypeDist).
//		SetBody(components.Row().
//			SetContent(colComp.SetSize(types.SizeMD(8)).
//				SetContent(pie8).
//				GetContent() + colComp.SetSize(types.SizeMD(4)).
//				SetContent(legend8).
//				GetContent()).GetContent()).
//		SetFooter(`<p class="text-center"><a href="/admin/info/schedule" class="uppercase3">` + viewAllTasks + `</a></p>`).
//		GetContent()
//
//	col8 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger8).GetContent()
//	row6 := components.Row().SetContent(col6 + col7 + col8).GetContent()
//	return types.Panel{
//		Content:     row1 + row2 + row3 + row6 + row5 + row4,
//		Title:       dashboardTitle,
//		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s</span><span style="color:#888">生成时间: %s</span></div>`, dashboardDesc, report.CreatedAt)),
//	}, nil
//}

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
			biz.Logger.Debug("%v", dr.ReportData)
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
	return renderGlobalReport(dr)
	//components := tmpl.Default()
	//colComp := components.Col()
	//apiMethods := globalReportData.APITypeCount.Infos
	//apiCounts := globalReportData.APITypeCount.Counts
	//colors := globalReportData.APITypeCount.Colors
	//labels := globalReportData.APITypeCount.Labels
	//pie1 := chartjs.Pie().
	//	SetHeight(120).
	//	SetLabels(apiMethods).
	//	SetID("pieChart1").
	//	AddDataSet(apiMethods[0]).
	//	DSData(apiCounts).
	//	DSBackgroundColor(colors).
	//	GetContent()
	//
	//legend1 := chart_legend.New().SetData(labels).GetContent()
	//
	//apiTypeDistribution := template.HTML(biz.T("dashboard.api_type_distribution"))
	//viewAllApis := template.HTML(biz.T("dashboard.view_all_apis"))
	//apiSpecCheck := template.HTML(biz.T("dashboard.api_spec_check"))
	//apiIsAutomation := template.HTML(biz.T("dashboard.api_is_automation"))
	//viewAllApiData := template.HTML(biz.T("dashboard.view_all_api_data"))
	//dataExecTrend6m := template.HTML(biz.T("dashboard.data_exec_trend_6m"))
	//dataExecDist6m := template.HTML(biz.T("dashboard.data_exec_dist_6m"))
	//viewAllHistoryData := template.HTML(biz.T("dashboard.view_all_history_data"))
	//sceneExecTrend6m := template.HTML(biz.T("dashboard.scene_exec_trend_6m"))
	//sceneExecDist6m := template.HTML(biz.T("dashboard.scene_exec_dist_6m"))
	//viewAllHistoryScene := template.HTML(biz.T("dashboard.view_all_history_scene"))
	//productList := template.HTML(biz.T("common.product_list"))
	//gotoProductDetail := template.HTML(biz.T("dashboard.goto_product_detail"))
	//appList := template.HTML(biz.T("dashboard.app_list"))
	//gotoAppDetail := template.HTML(biz.T("dashboard.goto_app_detail"))
	//sceneExecStatus := template.HTML(biz.T("dashboard.scene_exec_status"))
	//viewAllSceneRecords := template.HTML(biz.T("dashboard.view_all_scene_records"))
	//dataExecStatus := template.HTML(biz.T("dashboard.data_exec_status"))
	//viewAllDataRecords := template.HTML(biz.T("dashboard.view_all_data_records"))
	//taskTypeDist := template.HTML(biz.T("dashboard.task_type_dist"))
	//viewAllTasks := template.HTML(biz.T("dashboard.view_all_tasks"))
	//dashboardTitle := template.HTML(biz.T("dashboard.title"))
	//dashboardDesc := template.HTML(biz.T("dashboard.description"))
	//
	//boxDanger1 := components.Box().SetTheme("danger1").WithHeadBorder().SetHeader(apiTypeDistribution).
	//	SetBody(components.Row().
	//		SetContent(colComp.SetSize(types.SizeMD(8)).
	//			SetContent(pie1).
	//			GetContent() + colComp.SetSize(types.SizeMD(4)).
	//			SetContent(legend1).
	//			GetContent()).GetContent()).
	//	SetFooter(`<p class="text-center"><a href="/admin/info/api_definition" class="uppercase1">` + viewAllApis + `</a></p>`).
	//	GetContent()
	//col1 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger1).GetContent()
	//infos := globalReportData.APISpecCount.Infos
	//counts := globalReportData.APITypeCount.Counts
	//colors = globalReportData.APITypeCount.Colors
	//labels = globalReportData.APITypeCount.Labels
	//pie2 := chartjs.Pie().
	//	SetHeight(120).
	//	SetLabels(infos).
	//	SetID("pieChart2").
	//	AddDataSet(infos[0]).
	//	DSData(counts).
	//	DSBackgroundColor(colors).
	//	GetContent()
	//
	//legend2 := chart_legend.New().SetData(labels).GetContent()
	//
	//boxDanger2 := components.Box().SetTheme("danger2").WithHeadBorder().SetHeader(apiSpecCheck).
	//	SetBody(components.Row().
	//		SetContent(colComp.SetSize(types.SizeMD(8)).
	//			SetContent(pie2).
	//			GetContent() + colComp.SetSize(types.SizeMD(4)).
	//			SetContent(legend2).
	//			GetContent()).GetContent()).
	//	SetFooter(`<p class="text-center"><a href="/admin/info/api_definition" class="uppercase2">` + viewAllApis + `</a></p>`).
	//	GetContent()
	//
	//col2 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger2).GetContent()
	//infos = globalReportData.AutoAPICount.Infos
	//counts = globalReportData.AutoAPICount.Counts
	//colors = globalReportData.AutoAPICount.Colors
	//labels = globalReportData.AutoAPICount.Labels
	//pie3 := chartjs.Pie().
	//	SetHeight(120).
	//	SetLabels(infos).
	//	SetID("pieChart3").
	//	AddDataSet(infos[0]).
	//	DSData(counts).
	//	DSBackgroundColor(colors).
	//	GetContent()
	//
	//legend3 := chart_legend.New().SetData(labels).GetContent()
	//
	//boxDanger3 := components.Box().SetTheme("danger3").WithHeadBorder().SetHeader(apiIsAutomation).
	//	SetBody(components.Row().
	//		SetContent(colComp.SetSize(types.SizeMD(8)).
	//			SetContent(pie3).
	//			GetContent() + colComp.SetSize(types.SizeMD(4)).
	//			SetContent(legend3).
	//			GetContent()).GetContent()).
	//	SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">` + viewAllApiData + `</a></p>`).
	//	GetContent()
	//
	//col3 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger3).GetContent()
	//
	//row1 := components.Row().SetContent(col1 + col2 + col3).GetContent()
	//
	//line1 := chartjs.Line()
	//title := globalReportData.AppTestDataRunCount.Title
	//infos = globalReportData.AppTestDataRunCount.Infos
	//monthCounts := globalReportData.AppTestDataRunCount.Counts
	//monthLable := globalReportData.AppTestDataRunCount.DayList
	//var lineChart1 template.HTML
	//if len(infos) == 1 {
	//	lineChart1 = line1.
	//		SetID("dataChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//} else if len(infos) == 2 {
	//	lineChart1 = line1.
	//		SetID("dataChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgb(255, 205, 86)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[1]).
	//		DSData(monthCounts[1]).
	//		DSFill(false).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//} else if len(infos) == 3 {
	//	lineChart1 = line1.
	//		SetID("dataChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgb(255, 205, 86)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[1]).
	//		DSData(monthCounts[1]).
	//		DSFill(false).
	//		DSBorderColor("rgb(54, 162, 235)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[2]).
	//		DSData(monthCounts[2]).
	//		DSFill(true).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//} else if len(infos) == 4 {
	//	lineChart1 = line1.
	//		SetID("dataChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgb(255, 205, 86)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[1]).
	//		DSData(monthCounts[1]).
	//		DSFill(false).
	//		DSBorderColor("rgb(54, 162, 235)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[2]).
	//		DSData(monthCounts[2]).
	//		DSFill(false).
	//		DSBorderColor("rgb(238,232,170)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[3]).
	//		DSData(monthCounts[3]).
	//		DSFill(false).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//} else if len(infos) >= 5 {
	//	lineChart1 = line1.
	//		SetID("dataChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgb(255, 205, 86)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[1]).
	//		DSData(monthCounts[1]).
	//		DSFill(false).
	//		DSBorderColor("rgb(54, 162, 235)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[2]).
	//		DSData(monthCounts[2]).
	//		DSFill(false).
	//		DSBorderColor("rgb(238,232,170)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[3]).
	//		DSData(monthCounts[3]).
	//		DSFill(false).
	//		DSBorderColor("rgb(189,183,107)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[4]).
	//		DSData(monthCounts[4]).
	//		DSFill(false).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//}
	//boxInternalCol1 := colComp.SetContent(lineChart1).SetSize(types.SizeMD(12)).GetContent()
	//boxInternalRow1 := components.Row().SetContent(boxInternalCol1).GetContent()
	//box1 := components.Box().WithHeadBorder().SetHeader(dataExecTrend6m).
	//	SetBody(boxInternalRow1).
	//	GetContent()
	//
	//boxcol1 := colComp.SetContent(box1).SetSize(types.SizeMD(8)).GetContent()
	//infos = globalReportData.AppAPIRunCount.Infos
	//counts = globalReportData.AppAPIRunCount.Counts
	//colors = globalReportData.AppAPIRunCount.Colors
	//labels = globalReportData.AppAPIRunCount.Labels
	//pie4 := chartjs.Pie().
	//	SetHeight(120).
	//	SetLabels(infos).
	//	SetID("pieChart4").
	//	AddDataSet(infos[0]).
	//	DSData(counts).
	//	DSBackgroundColor(colors).
	//	GetContent()
	//
	//legend4 := chart_legend.New().SetData(labels).GetContent()
	//
	//boxDanger4 := components.Box().SetTheme("danger4").WithHeadBorder().SetHeader(dataExecDist6m).
	//	SetBody(components.Row().
	//		SetContent(colComp.SetSize(types.SizeMD(8)).
	//			SetContent(pie4).
	//			GetContent() + colComp.SetSize(types.SizeMD(4)).
	//			SetContent(legend4).
	//			GetContent()).GetContent()).
	//	SetFooter(`<p class="text-center"><a href="/admin/info/scene_data_test_history" class="uppercase3">` + viewAllHistoryData + `</a></p>`).
	//	GetContent()
	//
	//col4 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger4).GetContent()
	//
	//row2 := components.Row().SetContent(boxcol1 + col4).GetContent()
	//
	//line2 := chartjs.Line()
	//title = globalReportData.ProductPlaybookResultCount.Title
	//infos = globalReportData.ProductPlaybookResultCount.Infos
	//monthCounts = globalReportData.ProductPlaybookResultCount.Counts
	//monthLable = globalReportData.ProductPlaybookResultCount.DayList
	//var lineChart2 template.HTML
	//if len(infos) == 1 {
	//	lineChart2 = line2.
	//		SetID("sceneChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//} else if len(infos) == 2 {
	//	lineChart2 = line2.
	//		SetID("sceneChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgb(255, 205, 86)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[1]).
	//		DSData(monthCounts[1]).
	//		DSFill(false).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//} else if len(infos) == 3 {
	//	lineChart2 = line2.
	//		SetID("sceneChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgb(255, 205, 86)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[1]).
	//		DSData(monthCounts[1]).
	//		DSFill(false).
	//		DSBorderColor("rgb(54, 162, 235)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[2]).
	//		DSData(monthCounts[2]).
	//		DSFill(true).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//} else if len(infos) == 4 {
	//	lineChart2 = line2.
	//		SetID("sceneChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgb(255, 205, 86)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[1]).
	//		DSData(monthCounts[1]).
	//		DSFill(false).
	//		DSBorderColor("rgb(54, 162, 235)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[2]).
	//		DSData(monthCounts[2]).
	//		DSFill(false).
	//		DSBorderColor("rgb(238,232,170)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[3]).
	//		DSData(monthCounts[3]).
	//		DSFill(false).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//} else if len(infos) >= 5 {
	//	lineChart2 = line2.
	//		SetID("sceneChart").
	//		SetHeight(320).
	//		SetTitle(title).
	//		SetLabels(monthLable).
	//		AddDataSet(infos[0]).
	//		DSData(monthCounts[0]).
	//		DSFill(false).
	//		DSBorderColor("rgb(255, 205, 86)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[1]).
	//		DSData(monthCounts[1]).
	//		DSFill(false).
	//		DSBorderColor("rgb(54, 162, 235)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[2]).
	//		DSData(monthCounts[2]).
	//		DSFill(false).
	//		DSBorderColor("rgb(238,232,170)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[3]).
	//		DSData(monthCounts[3]).
	//		DSFill(false).
	//		DSBorderColor("rgb(189,183,107)").
	//		DSLineTension(0.1).
	//		AddDataSet(infos[4]).
	//		DSData(monthCounts[4]).
	//		DSFill(false).
	//		DSBorderColor("rgba(60,141,188,1)").
	//		DSLineTension(0.1).
	//		GetContent()
	//}
	//boxInternalCol2 := colComp.SetContent(lineChart2).SetSize(types.SizeMD(12)).GetContent()
	//boxInternalRow2 := components.Row().SetContent(boxInternalCol2).GetContent()
	//box2 := components.Box().WithHeadBorder().SetHeader(sceneExecTrend6m).
	//	SetBody(boxInternalRow2).
	//	GetContent()
	//
	//boxcol2 := colComp.SetContent(box2).SetSize(types.SizeMD(8)).GetContent()
	//infos = globalReportData.PlaybookResultCount.Infos
	//counts = globalReportData.PlaybookResultCount.Counts
	//colors = globalReportData.PlaybookResultCount.Colors
	//labels = globalReportData.PlaybookResultCount.Labels
	//pie5 := chartjs.Pie().
	//	SetHeight(120).
	//	SetLabels(infos).
	//	SetID("pieChart5").
	//	AddDataSet(infos[0]).
	//	DSData(counts).
	//	DSBackgroundColor(colors).
	//	GetContent()
	//legend5 := chart_legend.New().SetData(labels).GetContent()
	//boxDanger5 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(sceneExecDist6m).
	//	SetBody(components.Row().
	//		SetContent(colComp.SetSize(types.SizeMD(8)).
	//			SetContent(pie5).
	//			GetContent() + colComp.SetSize(types.SizeMD(4)).
	//			SetContent(legend5).
	//			GetContent()).GetContent()).
	//	SetFooter(`<p class="text-center"><a href="/admin/info/sence_test_history" class="uppercase3">` + viewAllHistoryScene + `</a></p>`).
	//	GetContent()
	//col5 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger5).GetContent()
	//row3 := components.Row().SetContent(boxcol2 + col5).GetContent()
	//contents := globalReportData.ProductsTableCount.Contents
	//headers := globalReportData.ProductsTableCount.Headers
	//table := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
	//boxInfo := components.Box().
	//	WithHeadBorder().
	//	SetHeader(productList).
	//	SetHeadColor("#f7f7f7").
	//	SetBody(table).
	//	SetFooter(`<div class="clearfix"><a href="/admin/info/product" class="btn btn-sm btn-default btn-flat pull-right">` + gotoProductDetail + `</a> </div>`).
	//	GetContent()
	//tableCol := colComp.SetSize(types.SizeMD(12)).SetContent(boxInfo).GetContent()
	//row4 := components.Row().SetContent(tableCol).GetContent()
	//var tableApp template.HTML
	//aContents, aHeaders := biz.GetAppsTableCount()
	//if len(aContents) > 0 {
	//	tableApp = components.Table().SetInfoList(aContents).SetThead(aHeaders).GetContent()
	//}
	//
	//boxAppInfo := components.Box().
	//	WithHeadBorder().
	//	SetHeader(appList).
	//	SetHeadColor("#f7f7f7").
	//	SetBody(tableApp).
	//	SetFooter(`<div class="clearfix"><a href="/admin/info/env_config" class="btn btn-sm btn-default btn-flat pull-right">` + gotoAppDetail + `</a> </div>`).
	//	GetContent()
	//tableAppCol := colComp.SetSize(types.SizeMD(12)).SetContent(boxAppInfo).GetContent()
	//row5 := components.Row().SetContent(tableAppCol).GetContent()
	//infos = globalReportData.PlaybookResultCount.Infos
	//counts = globalReportData.PlaybookResultCount.Counts
	//colors = globalReportData.PlaybookResultCount.Colors
	//labels = globalReportData.PlaybookResultCount.Labels
	//pie6 := chartjs.Pie().
	//	SetHeight(120).
	//	SetLabels(infos).
	//	SetID("pieChart6").
	//	AddDataSet(infos[0]).
	//	DSData(counts).
	//	DSBackgroundColor(colors).
	//	GetContent()
	//
	//legend6 := chart_legend.New().SetData(labels).GetContent()
	//
	//boxDanger6 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(sceneExecStatus).
	//	SetBody(components.Row().
	//		SetContent(colComp.SetSize(types.SizeMD(8)).
	//			SetContent(pie6).
	//			GetContent() + colComp.SetSize(types.SizeMD(4)).
	//			SetContent(legend6).
	//			GetContent()).GetContent()).
	//	SetFooter(`<p class="text-center"><a href="/admin/info/scene_test_history" class="uppercase3">` + viewAllSceneRecords + `</a></p>`).
	//	GetContent()
	//
	//col6 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger6).GetContent()
	//infos = globalReportData.TestDataResultCount.Infos
	//counts = globalReportData.TestDataResultCount.Counts
	//colors = globalReportData.TestDataResultCount.Colors
	//labels = globalReportData.TestDataResultCount.Labels
	//pie7 := chartjs.Pie().
	//	SetHeight(120).
	//	SetLabels(infos).
	//	SetID("pieChart7").
	//	AddDataSet(infos[0]).
	//	DSData(counts).
	//	DSBackgroundColor(colors).
	//	GetContent()
	//
	//legend7 := chart_legend.New().SetData(labels).GetContent()
	//
	//boxDanger7 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(dataExecStatus).
	//	SetBody(components.Row().
	//		SetContent(colComp.SetSize(types.SizeMD(8)).
	//			SetContent(pie7).
	//			GetContent() + colComp.SetSize(types.SizeMD(4)).
	//			SetContent(legend7).
	//			GetContent()).GetContent()).
	//	SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">` + viewAllDataRecords + `</a></p>`).
	//	GetContent()
	//
	//col7 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger7).GetContent()
	//infos = globalReportData.ScheduleResultCount.Infos
	//counts = globalReportData.ScheduleResultCount.Counts
	//colors = globalReportData.ScheduleResultCount.Colors
	//labels = globalReportData.ScheduleResultCount.Labels
	//pie8 := chartjs.Pie().
	//	SetHeight(120).
	//	SetLabels(infos).
	//	SetID("pieChart8").
	//	AddDataSet(infos[0]).
	//	DSData(counts).
	//	DSBackgroundColor(colors).
	//	GetContent()
	//
	//legend8 := chart_legend.New().SetData(labels).GetContent()
	//
	//boxDanger8 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(taskTypeDist).
	//	SetBody(components.Row().
	//		SetContent(colComp.SetSize(types.SizeMD(8)).
	//			SetContent(pie8).
	//			GetContent() + colComp.SetSize(types.SizeMD(4)).
	//			SetContent(legend8).
	//			GetContent()).GetContent()).
	//	SetFooter(`<p class="text-center"><a href="/admin/info/schedule" class="uppercase3">` + viewAllTasks + `</a></p>`).
	//	GetContent()
	//
	//col8 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger8).GetContent()
	//row6 := components.Row().SetContent(col6 + col7 + col8).GetContent()
	//return types.Panel{
	//	Content:     row1 + row2 + row3 + row6 + row5 + row4,
	//	Title:       dashboardTitle,
	//	Description: dashboardDesc,
	//}, nil
}

func renderGlobalReport(report biz.DashboardReport) (types.Panel, error) {
	var globalReportData GlobalDashboardReport
	err := json.Unmarshal([]byte(report.ReportData), &globalReportData)
	if err != nil {
		return types.Panel{
			Content:     template.HTML(fmt.Sprintf("<div style='padding:20px;color:red'>%s</div>", biz.T("schedule_report.parse_error", err))),
			Title:       template.HTML(biz.T("schedule_report.page_title")),
			Description: template.HTML(biz.T("schedule_report.description")),
		}, nil
	}
	components := tmpl.Default()
	colComp := components.Col()
	apiMethods := globalReportData.APITypeCount.Infos
	apiCounts := globalReportData.APITypeCount.Counts
	colors := globalReportData.APITypeCount.Colors
	labels := globalReportData.APITypeCount.Labels
	pie1 := chartjs.Pie().
		SetHeight(120).
		SetLabels(apiMethods).
		SetID("pieChart1").
		AddDataSet(apiMethods[0]).
		DSData(apiCounts).
		DSBackgroundColor(colors).
		GetContent()

	legend1 := chart_legend.New().SetData(labels).GetContent()

	apiTypeDistribution := template.HTML(biz.T("dashboard.api_type_distribution"))
	viewAllApis := template.HTML(biz.T("dashboard.view_all_apis"))
	apiSpecCheck := template.HTML(biz.T("dashboard.api_spec_check"))
	apiIsAutomation := template.HTML(biz.T("dashboard.api_is_automation"))
	viewAllApiData := template.HTML(biz.T("dashboard.view_all_api_data"))
	dataExecTrend6m := template.HTML(biz.T("dashboard.data_exec_trend_6m"))
	dataExecDist6m := template.HTML(biz.T("dashboard.data_exec_dist_6m"))
	viewAllHistoryData := template.HTML(biz.T("dashboard.view_all_history_data"))
	sceneExecTrend6m := template.HTML(biz.T("dashboard.scene_exec_trend_6m"))
	sceneExecDist6m := template.HTML(biz.T("dashboard.scene_exec_dist_6m"))
	viewAllHistoryScene := template.HTML(biz.T("dashboard.view_all_history_scene"))
	productList := template.HTML(biz.T("common.product_list"))
	gotoProductDetail := template.HTML(biz.T("dashboard.goto_product_detail"))
	appList := template.HTML(biz.T("dashboard.app_list"))
	gotoAppDetail := template.HTML(biz.T("dashboard.goto_app_detail"))
	sceneExecStatus := template.HTML(biz.T("dashboard.scene_exec_status"))
	viewAllSceneRecords := template.HTML(biz.T("dashboard.view_all_scene_records"))
	dataExecStatus := template.HTML(biz.T("dashboard.data_exec_status"))
	viewAllDataRecords := template.HTML(biz.T("dashboard.view_all_data_records"))
	taskTypeDist := template.HTML(biz.T("dashboard.task_type_dist"))
	viewAllTasks := template.HTML(biz.T("dashboard.view_all_tasks"))
	dashboardTitle := template.HTML(biz.T("dashboard.title"))
	dashboardDesc := template.HTML(biz.T("dashboard.description"))

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
	infos := globalReportData.APISpecCount.Infos
	counts := globalReportData.APITypeCount.Counts
	colors = globalReportData.APITypeCount.Colors
	labels = globalReportData.APITypeCount.Labels
	pie2 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart2").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend2 := chart_legend.New().SetData(labels).GetContent()

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
	infos = globalReportData.AutoAPICount.Infos
	counts = globalReportData.AutoAPICount.Counts
	colors = globalReportData.AutoAPICount.Colors
	labels = globalReportData.AutoAPICount.Labels
	pie3 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart3").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend3 := chart_legend.New().SetData(labels).GetContent()

	boxDanger3 := components.Box().SetTheme("danger3").WithHeadBorder().SetHeader(apiIsAutomation).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie3).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend3).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">` + viewAllApiData + `</a></p>`).
		GetContent()

	col3 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger3).GetContent()

	row1 := components.Row().SetContent(col1 + col2 + col3).GetContent()

	line1 := chartjs.Line()
	title := globalReportData.AppTestDataRunCount.Title
	infos = globalReportData.AppTestDataRunCount.Infos
	monthCounts := globalReportData.AppTestDataRunCount.Counts
	monthLable := globalReportData.AppTestDataRunCount.DayList
	var lineChart1 template.HTML
	if len(infos) == 1 {
		lineChart1 = line1.
			SetID("dataChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	} else if len(infos) == 2 {
		lineChart1 = line1.
			SetID("dataChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgb(255, 205, 86)").
			DSLineTension(0.1).
			AddDataSet(infos[1]).
			DSData(monthCounts[1]).
			DSFill(false).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	} else if len(infos) == 3 {
		lineChart1 = line1.
			SetID("dataChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgb(255, 205, 86)").
			DSLineTension(0.1).
			AddDataSet(infos[1]).
			DSData(monthCounts[1]).
			DSFill(false).
			DSBorderColor("rgb(54, 162, 235)").
			DSLineTension(0.1).
			AddDataSet(infos[2]).
			DSData(monthCounts[2]).
			DSFill(true).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	} else if len(infos) == 4 {
		lineChart1 = line1.
			SetID("dataChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgb(255, 205, 86)").
			DSLineTension(0.1).
			AddDataSet(infos[1]).
			DSData(monthCounts[1]).
			DSFill(false).
			DSBorderColor("rgb(54, 162, 235)").
			DSLineTension(0.1).
			AddDataSet(infos[2]).
			DSData(monthCounts[2]).
			DSFill(false).
			DSBorderColor("rgb(238,232,170)").
			DSLineTension(0.1).
			AddDataSet(infos[3]).
			DSData(monthCounts[3]).
			DSFill(false).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	} else if len(infos) >= 5 {
		lineChart1 = line1.
			SetID("dataChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgb(255, 205, 86)").
			DSLineTension(0.1).
			AddDataSet(infos[1]).
			DSData(monthCounts[1]).
			DSFill(false).
			DSBorderColor("rgb(54, 162, 235)").
			DSLineTension(0.1).
			AddDataSet(infos[2]).
			DSData(monthCounts[2]).
			DSFill(false).
			DSBorderColor("rgb(238,232,170)").
			DSLineTension(0.1).
			AddDataSet(infos[3]).
			DSData(monthCounts[3]).
			DSFill(false).
			DSBorderColor("rgb(189,183,107)").
			DSLineTension(0.1).
			AddDataSet(infos[4]).
			DSData(monthCounts[4]).
			DSFill(false).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	}
	boxInternalCol1 := colComp.SetContent(lineChart1).SetSize(types.SizeMD(12)).GetContent()
	boxInternalRow1 := components.Row().SetContent(boxInternalCol1).GetContent()
	box1 := components.Box().WithHeadBorder().SetHeader(dataExecTrend6m).
		SetBody(boxInternalRow1).
		GetContent()

	boxcol1 := colComp.SetContent(box1).SetSize(types.SizeMD(8)).GetContent()
	infos = globalReportData.AppAPIRunCount.Infos
	counts = globalReportData.AppAPIRunCount.Counts
	colors = globalReportData.AppAPIRunCount.Colors
	labels = globalReportData.AppAPIRunCount.Labels
	pie4 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart4").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend4 := chart_legend.New().SetData(labels).GetContent()

	boxDanger4 := components.Box().SetTheme("danger4").WithHeadBorder().SetHeader(dataExecDist6m).
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

	line2 := chartjs.Line()
	title = globalReportData.ProductPlaybookResultCount.Title
	infos = globalReportData.ProductPlaybookResultCount.Infos
	monthCounts = globalReportData.ProductPlaybookResultCount.Counts
	monthLable = globalReportData.ProductPlaybookResultCount.DayList
	var lineChart2 template.HTML
	if len(infos) == 1 {
		lineChart2 = line2.
			SetID("sceneChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	} else if len(infos) == 2 {
		lineChart2 = line2.
			SetID("sceneChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgb(255, 205, 86)").
			DSLineTension(0.1).
			AddDataSet(infos[1]).
			DSData(monthCounts[1]).
			DSFill(false).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	} else if len(infos) == 3 {
		lineChart2 = line2.
			SetID("sceneChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgb(255, 205, 86)").
			DSLineTension(0.1).
			AddDataSet(infos[1]).
			DSData(monthCounts[1]).
			DSFill(false).
			DSBorderColor("rgb(54, 162, 235)").
			DSLineTension(0.1).
			AddDataSet(infos[2]).
			DSData(monthCounts[2]).
			DSFill(true).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	} else if len(infos) == 4 {
		lineChart2 = line2.
			SetID("sceneChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgb(255, 205, 86)").
			DSLineTension(0.1).
			AddDataSet(infos[1]).
			DSData(monthCounts[1]).
			DSFill(false).
			DSBorderColor("rgb(54, 162, 235)").
			DSLineTension(0.1).
			AddDataSet(infos[2]).
			DSData(monthCounts[2]).
			DSFill(false).
			DSBorderColor("rgb(238,232,170)").
			DSLineTension(0.1).
			AddDataSet(infos[3]).
			DSData(monthCounts[3]).
			DSFill(false).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	} else if len(infos) >= 5 {
		lineChart2 = line2.
			SetID("sceneChart").
			SetHeight(320).
			SetTitle(title).
			SetLabels(monthLable).
			AddDataSet(infos[0]).
			DSData(monthCounts[0]).
			DSFill(false).
			DSBorderColor("rgb(255, 205, 86)").
			DSLineTension(0.1).
			AddDataSet(infos[1]).
			DSData(monthCounts[1]).
			DSFill(false).
			DSBorderColor("rgb(54, 162, 235)").
			DSLineTension(0.1).
			AddDataSet(infos[2]).
			DSData(monthCounts[2]).
			DSFill(false).
			DSBorderColor("rgb(238,232,170)").
			DSLineTension(0.1).
			AddDataSet(infos[3]).
			DSData(monthCounts[3]).
			DSFill(false).
			DSBorderColor("rgb(189,183,107)").
			DSLineTension(0.1).
			AddDataSet(infos[4]).
			DSData(monthCounts[4]).
			DSFill(false).
			DSBorderColor("rgba(60,141,188,1)").
			DSLineTension(0.1).
			GetContent()
	}
	boxInternalCol2 := colComp.SetContent(lineChart2).SetSize(types.SizeMD(12)).GetContent()
	boxInternalRow2 := components.Row().SetContent(boxInternalCol2).GetContent()
	box2 := components.Box().WithHeadBorder().SetHeader(sceneExecTrend6m).
		SetBody(boxInternalRow2).
		GetContent()

	boxcol2 := colComp.SetContent(box2).SetSize(types.SizeMD(8)).GetContent()
	infos = globalReportData.PlaybookResultCount.Infos
	counts = globalReportData.PlaybookResultCount.Counts
	colors = globalReportData.PlaybookResultCount.Colors
	labels = globalReportData.PlaybookResultCount.Labels
	pie5 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart5").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()
	legend5 := chart_legend.New().SetData(labels).GetContent()
	boxDanger5 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(sceneExecDist6m).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie5).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend5).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/sence_test_history" class="uppercase3">` + viewAllHistoryScene + `</a></p>`).
		GetContent()
	col5 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger5).GetContent()
	row3 := components.Row().SetContent(boxcol2 + col5).GetContent()
	contents := globalReportData.ProductsTableCount.Contents
	headers := globalReportData.ProductsTableCount.Headers
	table := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
	boxInfo := components.Box().
		WithHeadBorder().
		SetHeader(productList).
		SetHeadColor("#f7f7f7").
		SetBody(table).
		SetFooter(`<div class="clearfix"><a href="/admin/info/product" class="btn btn-sm btn-default btn-flat pull-right">` + gotoProductDetail + `</a> </div>`).
		GetContent()
	tableCol := colComp.SetSize(types.SizeMD(12)).SetContent(boxInfo).GetContent()
	row4 := components.Row().SetContent(tableCol).GetContent()
	var tableApp template.HTML
	aContents, aHeaders := biz.GetAppsTableCount()
	if len(aContents) > 0 {
		tableApp = components.Table().SetInfoList(aContents).SetThead(aHeaders).GetContent()
	}

	boxAppInfo := components.Box().
		WithHeadBorder().
		SetHeader(appList).
		SetHeadColor("#f7f7f7").
		SetBody(tableApp).
		SetFooter(`<div class="clearfix"><a href="/admin/info/env_config" class="btn btn-sm btn-default btn-flat pull-right">` + gotoAppDetail + `</a> </div>`).
		GetContent()
	tableAppCol := colComp.SetSize(types.SizeMD(12)).SetContent(boxAppInfo).GetContent()
	row5 := components.Row().SetContent(tableAppCol).GetContent()
	infos = globalReportData.PlaybookResultCount.Infos
	counts = globalReportData.PlaybookResultCount.Counts
	colors = globalReportData.PlaybookResultCount.Colors
	labels = globalReportData.PlaybookResultCount.Labels
	pie6 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart6").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend6 := chart_legend.New().SetData(labels).GetContent()

	boxDanger6 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(sceneExecStatus).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie6).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend6).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_test_history" class="uppercase3">` + viewAllSceneRecords + `</a></p>`).
		GetContent()

	col6 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger6).GetContent()
	infos = globalReportData.TestDataResultCount.Infos
	counts = globalReportData.TestDataResultCount.Counts
	colors = globalReportData.TestDataResultCount.Colors
	labels = globalReportData.TestDataResultCount.Labels
	pie7 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart7").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend7 := chart_legend.New().SetData(labels).GetContent()

	boxDanger7 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(dataExecStatus).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie7).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend7).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">` + viewAllDataRecords + `</a></p>`).
		GetContent()

	col7 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger7).GetContent()
	infos = globalReportData.ScheduleResultCount.Infos
	counts = globalReportData.ScheduleResultCount.Counts
	colors = globalReportData.ScheduleResultCount.Colors
	labels = globalReportData.ScheduleResultCount.Labels
	pie8 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart8").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend8 := chart_legend.New().SetData(labels).GetContent()

	boxDanger8 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(taskTypeDist).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie8).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend8).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/schedule" class="uppercase3">` + viewAllTasks + `</a></p>`).
		GetContent()

	col8 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger8).GetContent()
	row6 := components.Row().SetContent(col6 + col7 + col8).GetContent()
	return types.Panel{
		Content:     row1 + row2 + row3 + row6 + row5 + row4,
		Title:       dashboardTitle,
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s</span><span style="color:#888">%s: %s</span></div>`, dashboardDesc, biz.T("schedule_report.generated_at"), report.CreatedAt)),
	}, nil
}
