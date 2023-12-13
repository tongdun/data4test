package pages

import (
	"data4perf/biz"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"html/template"

	//"data4perf/biz"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/chart_legend"
	//"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	//"github.com/GoAdminGroup/themes/adminlte/components/smallbox"
	"github.com/gin-gonic/gin"
	//"github.com/GoAdminGroup/components/echarts"
)

func GetDashBoardContent(ctx *gin.Context) (types.Panel, error) {

	components := tmpl.Default()
	colComp := components.Col()
	apiMethods, apiCounts, colors, labels := biz.GetAPITypeCount("all", "")
	pie1 := chartjs.Pie().
		SetHeight(120).
		SetLabels(apiMethods).
		SetID("pieChart1").
		AddDataSet(apiMethods[0]).
		DSData(apiCounts).
		DSBackgroundColor(colors).
		GetContent()

	legend1 := chart_legend.New().SetData(labels).GetContent()

	boxDanger1 := components.Box().SetTheme("danger1").WithHeadBorder().SetHeader("接口类型分布").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie1).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend1).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/api_definition" class="uppercase1">查看全部接口</a></p>`).
		GetContent()
	col1 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger1).GetContent()


	infos, counts, colors, labels := biz.GetAPISpecCount("all", "")
	pie2 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart2").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend2 := chart_legend.New().SetData(labels).GetContent()

	boxDanger2 := components.Box().SetTheme("danger2").WithHeadBorder().SetHeader("接口规范检查").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie2).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend2).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/api_definition" class="uppercase2">查看全部接口</a></p>`).
		GetContent()

	col2 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger2).GetContent()

	infos, counts, colors, labels = biz.GetAutoAPICount("all", "")
	pie3 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart3").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend3 := chart_legend.New().SetData(labels).GetContent()

	boxDanger3 := components.Box().SetTheme("danger3").WithHeadBorder().SetHeader("接口是否自动化").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie3).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend3).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">查看全部接口数据文件</a></p>`).
		GetContent()

	col3 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger3).GetContent()

	row1 := components.Row().SetContent(col1 + col2 + col3).GetContent()

	line1 := chartjs.Line()
	title, monthLable, infos, monthCounts := biz.GetAppSceneDataRunCount()
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
	box1 := components.Box().WithHeadBorder().SetHeader("最近6个月数据文件执行趋势图").
		SetBody(boxInternalRow1).
		GetContent()

	boxcol1 := colComp.SetContent(box1).SetSize(types.SizeMD(8)).GetContent()
	infos, counts, colors, labels = biz.GetAppAPIRunCount()
	pie4 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart4").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend4 := chart_legend.New().SetData(labels).GetContent()

	boxDanger4 := components.Box().SetTheme("danger4").WithHeadBorder().SetHeader("最近6个月数据文件执行分布").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie4).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend4).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data_test_history" class="uppercase3">查看全部历史数据记录</a></p>`).
		GetContent()

	col4 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger4).GetContent()


	row2 := components.Row().SetContent(boxcol1+col4).GetContent()

	line2 := chartjs.Line()
	title, monthLable, infos, monthCounts = biz.GetProductSceneRunCount()
	//biz.Logger.Debug("title: %v, monthLable: %v, infos: %v, monthCounts: %v",title, monthLable, infos, monthCounts)
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
	}  else if len(infos) >= 5 {
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
	box2 := components.Box().WithHeadBorder().SetHeader("最近6个月产品场景执行趋势图").
		SetBody(boxInternalRow2).
		GetContent()

	boxcol2 := colComp.SetContent(box2).SetSize(types.SizeMD(8)).GetContent()

	infos, counts, colors, labels = biz.GetSceneRunCount()
	pie5 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart5").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()
	legend5 := chart_legend.New().SetData(labels).GetContent()
	boxDanger5 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader("最近6个月产品场景执行分布").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie5).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend5).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/sence_test_history" class="uppercase3">查看全部历史场景记录</a></p>`).
		GetContent()
	col5 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger5).GetContent()
	row3 := components.Row().SetContent(boxcol2+col5).GetContent()

	contents, headers := biz.GetProductsTableCount()
	tableProduct := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
	boxInfo := components.Box().
		WithHeadBorder().
		SetHeader("产品列表").
		SetHeadColor("#f7f7f7").
		SetBody(tableProduct).
		SetFooter(`<div class="clearfix"><a href="/admin/info/product" class="btn btn-sm btn-default btn-flat pull-right">跳转产品列表详情</a> </div>`).
		GetContent()
	tableCol := colComp.SetSize(types.SizeMD(12)).SetContent(boxInfo).GetContent()
	row4 := components.Row().SetContent(tableCol).GetContent()

	contents, headers = biz.GetAppsTableCount()
	tableApp := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
	boxAppInfo := components.Box().
		WithHeadBorder().
		SetHeader("应用列表").
		SetHeadColor("#f7f7f7").
		SetBody(tableApp).
		SetFooter(`<div class="clearfix"><a href="/admin/info/env_config" class="btn btn-sm btn-default btn-flat pull-right">跳转应用列表详情</a> </div>`).
		GetContent()
	tableAppCol := colComp.SetSize(types.SizeMD(12)).SetContent(boxAppInfo).GetContent()
	row5 := components.Row().SetContent(tableAppCol).GetContent()


	infos, counts, colors, labels = biz.GetSceneResultCount()
	pie6 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart6").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend6 := chart_legend.New().SetData(labels).GetContent()

	boxDanger6 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader("场景执行状态分布").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie6).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend6).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_test_history" class="uppercase3">查看全部场景执行记录</a></p>`).
		GetContent()

	col6 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger6).GetContent()

	infos, counts, colors, labels = biz.GetSceneDataResultCount()
	pie7 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart7").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend7 := chart_legend.New().SetData(labels).GetContent()

	boxDanger7 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader("数据文件执行状态分布").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie7).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend7).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">查看全部数据文件执行记录</a></p>`).
		GetContent()

	col7 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger7).GetContent()

	infos, counts, colors, labels = biz.GetScheduleTypeCount()
	pie8 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart8").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend8 := chart_legend.New().SetData(labels).GetContent()

	boxDanger8 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader("任务类型分布").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie8).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend8).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/schedule" class="uppercase3">查看全部任务</a></p>`).
		GetContent()

	col8 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger8).GetContent()
	row6 := components.Row().SetContent(col6+col7+col8).GetContent()

	return types.Panel{
		Content:     row1+row2+row3+row6+row5+row4,
		Title:       "统计报告",
		Description: "全局",
	}, nil
}
