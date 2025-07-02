package pages

import (
	"data4test/biz"
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/template/chartjs"

	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/chart_legend"
	"github.com/gin-gonic/gin"
)

func GetDashBoard3Content(ctx *gin.Context) (types.Panel, error) {
	id := ctx.Query("id")
	appName, err := biz.GetAppName(id)
	if err != nil {
		return types.Panel{}, err
	}

	components := tmpl.Default()
	colComp := components.Col()

	apiMethods, apiCounts, colors, labels := biz.GetAPITypeCount("app", appName)
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

	infos, counts, colors, labels := biz.GetAPISpecCount("app", appName)
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

	infos, counts, colors, labels = biz.GetAutoAPICount("app", appName)
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
	title, dayLable, infos, dayCounts := biz.GetAPIRunResultCount("app", appName)
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

	infos, counts, colors, labels = biz.GetDaysAPIResultCount("app", appName, 30)
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

	contents, headers := biz.GetAppModuleTableCount(appName)
	table := components.Table().SetInfoList(contents).SetThead(headers).GetContent()
	row4 := components.Box().
		WithHeadBorder().
		SetHeader("模块接口总览").
		SetHeadColor("#f7f7f7").
		SetBody(table).
		GetContent()

	desc := fmt.Sprintf("应用维度 - %s", appName)
	descHtml := template.HTML(desc)

	return types.Panel{
		Content:     row1 + row2 + row4,
		Title:       "统计报告",
		Description: descHtml,
	}, nil
}
