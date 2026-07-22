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

func GetProductReportData(productName, appName, userName string) (err error) {
	var productDR ProductDashboardReport

	productDR.APITypeCount.Infos, productDR.APITypeCount.Counts, productDR.APITypeCount.Colors, productDR.APITypeCount.Labels = biz.GetAPITypeCount("product", appName)
	productDR.APISpecCount.Infos, productDR.APISpecCount.Counts, productDR.APISpecCount.Colors, productDR.APISpecCount.Labels = biz.GetAPISpecCount("product", appName)
	productDR.AutoAPICount.Infos, productDR.AutoAPICount.Counts, productDR.AutoAPICount.Colors, productDR.AutoAPICount.Labels = biz.GetAutoAPICount("product", appName)
	productDR.APIRunResultCount.Title, productDR.APIRunResultCount.DayList, productDR.APIRunResultCount.Infos, productDR.APIRunResultCount.Counts = biz.GetAPIRunResultCount("product", productName)
	productDR.DaysAPIResultCount.Infos, productDR.DaysAPIResultCount.Counts, productDR.DaysAPIResultCount.Colors, productDR.DaysAPIResultCount.Labels = biz.GetDaysAPIResultCount("product", productName, 30)
	productDR.ProductPlaybookResultCount.Title, productDR.ProductPlaybookResultCount.DayList, productDR.ProductPlaybookResultCount.Infos, productDR.ProductPlaybookResultCount.Counts = biz.GetProductPlaybookRunResultCount(productName)
	productDR.DaysSceneResultCount.Infos, productDR.DaysSceneResultCount.Counts, productDR.DaysSceneResultCount.Colors, productDR.DaysSceneResultCount.Labels = biz.GetDaysSceneResultCount(productName, 30)
	productDR.ProductAppModuleTableCount.Contents, productDR.ProductAppModuleTableCount.Headers = biz.GetProductAppTableCount(appName)

	productDRStr, _ := json.MarshalIndent(productDR, "", "    ")

	nowStr := time.Now().Format("20060102150405")
	now := time.Now().Format("2006-01-02 15:04:05")
	reportName := biz.T("product.report_name", nowStr)
	report := biz.DashboardReport{
		ReportName:      reportName,
		ReportType:      "product",
		RelatedProducts: productName,
		RelatedApps:     appName,
		TimeRangeStart:  now,
		TimeRangeEnd:    now,
		Status:          "finished",
		Creator:         userName,
		ReportData:      string(productDRStr),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = models.Orm.Table("dashboard").Create(&report).Error
	if err != nil {
		return fmt.Errorf(biz.T("product.save_report_failed"), err)
	}
	return err
}

func GetDashBoard2Content(ctx *gin.Context, userName string) (types.Panel, error) {
	id := ctx.Query("id")
	productName, err := biz.GetProductName(id)
	if err != nil {
		return types.Panel{}, err
	}
	appName, err := biz.GetProductApps(id)
	if err != nil {
		return types.Panel{}, err
	}

	// 优先从 dashboard 表读取预计算数据
	var productReportData ProductDashboardReport
	var dr biz.DashboardReport
	err = models.Orm.Table("dashboard").
		Where("related_products =  ?", productName).
		Order("created_at desc").
		Limit(1).
		Find(&dr).Error

	if len(dr.Id) > 0 {
		err := json.Unmarshal([]byte(dr.ReportData), &productReportData)
		if err != nil {
			biz.Logger.Error("err: %v", err)
			biz.Logger.Debug("%v", dr.ReportData)
		}
	} else {
		err = GetProductReportData(productName, appName, userName)
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
			err = json.Unmarshal([]byte(dr.ReportData), &productReportData)
			if err != nil {
				biz.Logger.Error("err: %v", err)
				biz.Logger.Debug("%+v", productReportData)
			}
		}
	}
	return renderProductReport(productName, dr)

}

func renderProductReport(productName string, report biz.DashboardReport) (types.Panel, error) {
	var productReportData ProductDashboardReport
	err := json.Unmarshal([]byte(report.ReportData), &productReportData)
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
	sceneExecTrend30d := template.HTML(biz.T("dashboard.scene_exec_trend_30d"))
	sceneExecDist30d := template.HTML(biz.T("dashboard.scene_exec_dist_30d"))
	viewAllHistoryScene := template.HTML(biz.T("dashboard.view_all_history_scene"))
	appApiOverview := template.HTML(biz.T("dashboard.app_api_overview"))
	productDashboardTitle := template.HTML(biz.T("product_dashboard.title"))
	productDashboardDesc := template.HTML(biz.T("product_dashboard.description"))

	apiMethods := productReportData.APITypeCount.Infos
	apiCounts := productReportData.APITypeCount.Counts
	colors := productReportData.APITypeCount.Colors
	labels := productReportData.APITypeCount.Labels
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

	infos := productReportData.APISpecCount.Infos
	counts := productReportData.APITypeCount.Counts
	colors = productReportData.APITypeCount.Colors
	labels = productReportData.APITypeCount.Labels
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

	infos = productReportData.AutoAPICount.Infos
	counts = productReportData.AutoAPICount.Counts
	colors = productReportData.AutoAPICount.Colors
	labels = productReportData.AutoAPICount.Labels
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
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">` + viewAllDataFiles + `</a></p>`).
		GetContent()

	col3 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger3).GetContent()

	row1 := components.Row().SetContent(col1 + col2 + col3).GetContent()

	line1 := chartjs.Line()
	title := productReportData.APIRunResultCount.Title
	infos = productReportData.APIRunResultCount.Infos
	dayCounts := productReportData.APIRunResultCount.Counts
	dayLable := productReportData.APIRunResultCount.DayList
	lineChart1 := line1.
		SetID("dataChart").
		SetHeight(320).
		SetTitle(title).
		SetLabels(dayLable).
		AddDataSet(infos[0]).
		DSData(dayCounts[0]).
		DSFill(false).
		DSBorderColor("rgb(255, 205, 86)").
		DSLineTension(0.1).
		AddDataSet(infos[1]).
		DSData(dayCounts[1]).
		DSFill(false).
		DSBorderColor("rgb(54, 162, 235)").
		DSLineTension(0.1).
		AddDataSet(infos[2]).
		DSData(dayCounts[2]).
		DSFill(true).
		DSBorderColor("rgba(60,141,188,1)").
		DSLineTension(0.1).
		GetContent()
	boxInternalCol1 := colComp.SetContent(lineChart1).SetSize(types.SizeMD(12)).GetContent()
	boxInternalRow1 := components.Row().SetContent(boxInternalCol1).GetContent()
	box1 := components.Box().WithHeadBorder().SetHeader(dataExecTrend30d).
		SetBody(boxInternalRow1).
		GetContent()

	boxcol1 := colComp.SetContent(box1).SetSize(types.SizeMD(8)).GetContent()
	infos = productReportData.DaysAPIResultCount.Infos
	counts = productReportData.DaysAPIResultCount.Counts
	colors = productReportData.DaysAPIResultCount.Colors
	labels = productReportData.DaysAPIResultCount.Labels
	pie4 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart4").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend4 := chart_legend.New().SetData(labels).GetContent()

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

	line2 := chartjs.Line()
	title = productReportData.ProductPlaybookResultCount.Title
	infos = productReportData.ProductPlaybookResultCount.Infos
	dayCounts = productReportData.ProductPlaybookResultCount.Counts
	dayLable = productReportData.ProductPlaybookResultCount.DayList
	lineChart2 := line2.
		SetID("sceneChart").
		SetHeight(320).
		SetTitle(title).
		SetLabels(dayLable).
		AddDataSet(infos[0]).
		DSData(dayCounts[0]).
		DSFill(false).
		DSBorderColor("rgb(255, 205, 86)").
		DSLineTension(0.1).
		AddDataSet(infos[1]).
		DSData(dayCounts[1]).
		DSFill(false).
		DSBorderColor("rgb(54, 162, 235)").
		DSLineTension(0.1).
		AddDataSet(infos[2]).
		DSData(dayCounts[2]).
		DSFill(true).
		DSBorderColor("rgba(60,141,188,1)").
		DSLineTension(0.1).
		GetContent()
	boxInternalCol2 := colComp.SetContent(lineChart2).SetSize(types.SizeMD(12)).GetContent()
	boxInternalRow2 := components.Row().SetContent(boxInternalCol2).GetContent()
	box2 := components.Box().WithHeadBorder().SetHeader(sceneExecTrend30d).
		SetBody(boxInternalRow2).
		GetContent()

	boxcol2 := colComp.SetContent(box2).SetSize(types.SizeMD(8)).GetContent()
	infos = productReportData.DaysSceneResultCount.Infos
	counts = productReportData.DaysSceneResultCount.Counts
	colors = productReportData.DaysSceneResultCount.Colors
	labels = productReportData.DaysSceneResultCount.Labels
	pie5 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart5").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend5 := chart_legend.New().SetData(labels).GetContent()

	boxDanger5 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader(sceneExecDist30d).
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie5).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend5).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_test_history" class="uppercase3">` + viewAllHistoryScene + `</a></p>`).
		GetContent()

	col5 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger5).GetContent()

	row3 := components.Row().SetContent(boxcol2 + col5).GetContent()

	contents := productReportData.ProductAppModuleTableCount.Contents
	headers := productReportData.ProductAppModuleTableCount.Headers
	table := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
	row4 := components.Box().
		WithHeadBorder().
		SetHeader(appApiOverview).
		SetHeadColor("#f7f7f7").
		SetBody(table).
		GetContent()

	return types.Panel{
		Content:     row1 + row2 + row3 + row4,
		Title:       productDashboardTitle,
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s - %s</span><span style="color:#888">%s: %s</span></div>`, productDashboardDesc, productName, biz.T("schedule_report.generated_at"), report.CreatedAt)),
	}, nil
}
