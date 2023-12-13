package pages

import (
	"data4perf/biz"
	"fmt"
	"html/template"

	//"data4perf/biz"
	"github.com/GoAdminGroup/go-admin/template/chartjs"

	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/chart_legend"
	//"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	//"github.com/GoAdminGroup/themes/adminlte/components/smallbox"
	"github.com/gin-gonic/gin"
)

func GetDashBoard2Content(ctx *gin.Context) (types.Panel, error) {
	id := ctx.Query("id")
	productName, err := biz.GetProductName(id)
	if err != nil {
		return types.Panel{}, err
	}
	appName, err := biz.GetProductApps(id)
	if err != nil {
		return types.Panel{}, err
	}

	components := tmpl.Default()
	colComp := components.Col()

	apiMethods, apiCounts, colors, labels := biz.GetAPITypeCount("product", appName)
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

	infos, counts, colors, labels := biz.GetAPISpecCount("product", appName)
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

	infos, counts, colors, labels = biz.GetAutoAPICount("product", appName)
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
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data" class="uppercase3">查看全部数据文件</a></p>`).
		GetContent()

	col3 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger3).GetContent()

	row1 := components.Row().SetContent(col1 + col2 + col3).GetContent()

	line1 := chartjs.Line()
	title, dayLable, infos, dayCounts := biz.GetAPIRunResultCount("product", productName)
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
	box1 := components.Box().WithHeadBorder().SetHeader("最近30天数据文件执行趋势图").
		SetBody(boxInternalRow1).
		GetContent()

	boxcol1 := colComp.SetContent(box1).SetSize(types.SizeMD(8)).GetContent()
	infos, counts, colors, labels = biz.GetDaysAPIResultCount("product", productName, 30)
	pie4 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart4").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend4 := chart_legend.New().SetData(labels).GetContent()

	boxDanger4 := components.Box().SetTheme("danger4").WithHeadBorder().SetHeader("最近30天数据文件执行分布").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie4).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend4).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_data_test_history" class="uppercase3">查看全部历史数据记录</a></p>`).
		GetContent()

	col4 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger4).GetContent()

	row2 := components.Row().SetContent(boxcol1 + col4).GetContent()

	line2 := chartjs.Line()
	title, dayLable, infos, dayCounts = biz.GetProductPlaybookRunResultCount(productName)
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
	box2 := components.Box().WithHeadBorder().SetHeader("最近30天产品场景执行趋势图").
		SetBody(boxInternalRow2).
		GetContent()

	boxcol2 := colComp.SetContent(box2).SetSize(types.SizeMD(8)).GetContent()
	infos, counts, colors, labels = biz.GetDaysSceneResultCount(productName, 30)
	pie5 := chartjs.Pie().
		SetHeight(120).
		SetLabels(infos).
		SetID("pieChart5").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend5 := chart_legend.New().SetData(labels).GetContent()

	boxDanger5 := components.Box().SetTheme("danger5").WithHeadBorder().SetHeader("最近30天产品场景执行分布").
		SetBody(components.Row().
			SetContent(colComp.SetSize(types.SizeMD(8)).
				SetContent(pie5).
				GetContent() + colComp.SetSize(types.SizeMD(4)).
				SetContent(legend5).
				GetContent()).GetContent()).
		SetFooter(`<p class="text-center"><a href="/admin/info/scene_test_history" class="uppercase3">查看全部历史场景记录</a></p>`).
		GetContent()

	col5 := colComp.SetSize(types.SizeMD(4)).SetContent(boxDanger5).GetContent()

	row3 := components.Row().SetContent(boxcol2 + col5).GetContent()

	contents, headers := biz.GetProductAppTableCount(appName)
	table := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
	row4 := components.Box().
		WithHeadBorder().
		SetHeader("涉及应用接口总览").
		SetHeadColor("#f7f7f7").
		SetBody(table).
		GetContent()

	desc := fmt.Sprintf("产品维度 - %s", productName)
	descHtml := template.HTML(desc)
	return types.Panel{
		Content:     row1 + row2 + row3 + row4,
		Title:       "统计报告",
		Description: descHtml,
	}, nil
}
