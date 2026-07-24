package pages

import (
	"data4test/biz"
	"encoding/json"
	"fmt"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"

	"data4test/models"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/chart_legend"
)

func GetDashboardByReportId(ctx *gin.Context) (types.Panel, error) {
	reportId := ctx.Query("id")
	var report biz.DashboardReport
	models.Orm.Table("dashboard").
		Where("id = ?", reportId).
		Find(&report)
	switch report.ReportType {
	case "task":
		if strings.Contains(report.RelatedTaskIds, ",") {
			return renderMultiTaskReport(report)
		}

		return renderTaskReport(report, "")
	case "product":
		return renderProductReport(report.RelatedProducts, report)
	case "app":
		return renderApplicationReport(report.RelatedApps, report)
	case "global":
		var globalReportData GlobalDashboardReport
		err := json.Unmarshal([]byte(report.ReportData), &globalReportData)
		if err != nil {
			biz.Logger.Error("err: %v", err)
			return types.Panel{}, err
		}
		return renderGlobalReport(globalReportData, report)
	default:
		desc := template.HTML(biz.T("schedule_report.description"))
		return types.Panel{
			Content: template.HTML(`<div style='padding:40px;text-align:center'>
							<h3>` + biz.T("schedule_report.not_supported") + `</h3>
							</div>`),
			Title:       template.HTML(biz.T("schedule_report.page_title")),
			Description: template.HTML(desc),
		}, nil

	}

	return renderTaskReport(report, "")
}

func GetScheduleReportContent(ctx *gin.Context) (types.Panel, error) {
	taskId := ctx.Query("id")
	var report biz.DashboardReport
	if len(taskId) > 0 {
		models.Orm.Table("dashboard").
			Where("report_type = ? and status = ? and related_task_ids like CONCAT(?, '\\_%')", "task", "finished", taskId).
			Order("created_at desc").
			Limit(1).
			Find(&report)
		if len(report.Id) == 0 || len(report.ReportData) == 0 {
			return types.Panel{
				Content: template.HTML(`<div style='padding:40px;text-align:center'>
							<h3>` + biz.T("schedule_report.no_report") + `</h3>
							<p style='color:#888;margin-top:20px'>` + biz.T("schedule_report.view_report_empty") + `</p>
						</div>`),
				Title:       template.HTML(biz.T("schedule_report.page_title")),
				Description: template.HTML(biz.T("schedule_report.description")),
			}, nil
		}

	}
	return renderTaskReport(report, taskId)
}

// renderTaskReport 解析 TaskReportData JSON 并渲染完整报告页面
func renderTaskReport(report biz.DashboardReport, scheduleId string) (types.Panel, error) {
	var reportData biz.TaskReportData
	err := json.Unmarshal([]byte(report.ReportData), &reportData)
	if err != nil {
		return types.Panel{
			Content:     template.HTML(fmt.Sprintf("<div style='padding:20px;color:red'>%s</div>", biz.T("schedule_report.parse_error", err))),
			Title:       template.HTML(biz.T("schedule_report.page_title")),
			Description: template.HTML(biz.T("schedule_report.description")),
		}, nil
	}

	// 获取执行任务 ID 用于场景历史跳转链接
	// 优先使用 report.RelatedTaskIds（执行 taskTag），为空时用 scheduleId 反查
	taskId := report.RelatedTaskIds
	if len(taskId) == 0 && len(scheduleId) > 0 {
		var latestReport biz.DashboardReport
		models.Orm.Table("dashboard").
			Select("related_task_ids").
			Where("report_type = ? and status = ? and related_task_ids like CONCAT(?, '\\_%')", "task", "finished", scheduleId).
			Order("created_at desc").
			Limit(1).
			Find(&latestReport)
		taskId = latestReport.RelatedTaskIds
	}

	// ====== KPI卡片 + 执行信息 ======
	headerInfo := buildTaskHeader(reportData)

	kpiCards := buildTaskKpiCards(reportData)

	// ====== 第1行: API类型分布 + 场景执行结果 + 数据执行结果 (3饼图) ======
	playbookPie := buildTaskSceneResultPie(reportData, taskId)
	dataPie := buildTaskDataResultPie(reportData, taskId)
	apiPie := buildTaskAPIPie(reportData)
	row1 := fmt.Sprintf(`<div class="row">%s%s%s</div>`, playbookPie, dataPie, apiPie)

	// ====== 第2行: 执行趋势(折线图) ======
	trendHTML := buildTaskTrendChart(reportData)
	row2 := ""
	if len(trendHTML) > 0 {
		row2 = fmt.Sprintf(`<div class="row"><div class="col-md-12">%s</div></div>`, trendHTML)
	}

	// ====== 第3行: 场景明细表(独占一行) ======
	sceneTable := buildTaskSceneTable(reportData)
	row3 := ""
	if len(sceneTable) > 0 {
		row3 = fmt.Sprintf(`<div class="row"><div class="col-md-12">%s</div></div>`, sceneTable)
	}

	// ====== 第4行: 数据文件明细表(独占一行) ======
	dataTable := buildTaskDataTable(reportData)
	row4 := ""
	if len(dataTable) > 0 {
		row4 = fmt.Sprintf(`<div class="row"><div class="col-md-12">%s</div></div>`, dataTable)
	}

	// ====== 第5行: 失败明细表 ======
	failHTML := buildTaskFailTable(reportData)

	content := string(headerInfo) + string(kpiCards) + row1 + row2 + row3 + row4 + string(failHTML)
	styleBlock := `<style>.sc-td{position:relative;cursor:pointer}.sc-td .sc-truncate{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.sc-td .sc-full{display:none;position:absolute;right:5%;top:50%;background:#fff;border:2px solid #666;padding:15px;z-index:9999;max-width:600px;max-height:80vh;overflow-y:auto;white-space:pre-wrap;word-break:break-all;box-shadow:0 4px 20px rgba(0,0,0,0.3);border-radius:4px;font-size:13px;line-height:1.4}.sc-td:hover .sc-full{display:block!important}</style>`
	content += styleBlock

	return types.Panel{
		Content:     template.HTML(content),
		Title:       template.HTML(biz.T("schedule_report.page_title")),
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s</span><span style="color:#888">%s: %s</span></div>`, report.ReportName, biz.T("schedule_report.generated_at"), report.CreatedAt)),
	}, nil
}

// ==================== KPI卡片 ====================

func buildTaskKpiCards(data biz.TaskReportData) template.HTML {
	o := data.Overview
	passRate := fmt.Sprintf("%.1f%%", o.PassRate)

	// 第一行：场景执行统计（执行数/通过/失败/通过率）
	card1 := fmt.Sprintf(`<div class="col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-blue"><i class="fa fa-cubes"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>`, biz.T("schedule_report.scene_exec_count"), o.TotalExecuted)
	card2 := fmt.Sprintf(`<div class="col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-green"><i class="fa fa-check-circle"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>`, biz.T("schedule_report.pass_count"), o.SuccessCount)
	card3 := fmt.Sprintf(`<div class="col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-red"><i class="fa fa-times-circle"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>`, biz.T("common.fail_count"), o.FailCount)
	card4 := fmt.Sprintf(`<div class="col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-yellow"><i class="fa fa-percent"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%s</span></div></div></div>`, biz.T("schedule_report.pass_rate_label"), passRate)

	row1 := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`, card1, card2, card3, card4)

	return template.HTML(row1)
}

// buildTaskHeader 顶部执行信息（白色背景，两行排版）
func buildTaskHeader(data biz.TaskReportData) template.HTML {
	o := data.Overview
	envInfo := o.Environment
	if len(envInfo) == 0 {
		envInfo = biz.T("schedule_report.unknown")
	}
	durationStr := formatDuration(o.DurationSeconds)
	statusColor := resultColor(o.Executor)

	html := fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-default" style="border-top:none;box-shadow:none;background:#fff;margin-bottom:10px"><div class="box-body" style="padding:15px 20px">
			<div class="row" style="margin-bottom:8px">
				<div class="col-md-4"><strong>%s</strong> %s</div>
				<div class="col-md-2"><strong>%s</strong> %s</div>
				<div class="col-md-4"><strong>%s</strong> %s</div>
				<div class="col-md-2"><strong>%s</strong> <span style="color:%s;font-weight:bold">%s</span></div>
			</div>
			<div class="row">
				<div class="col-md-4"><strong>%s</strong> %s ~ %s</div>
				<div class="col-md-2"><strong>%s</strong> %s</div>
				<div class="col-md-4"><strong>%s</strong> %d | <strong>%s</strong> %d | <strong>%s</strong> %d</div>
				<div class="col-md-2"><strong>%s</strong> <span style="color:green">%.1f%%</span></div>
			</div>
		</div></div></div></div>`,
		biz.T("schedule_report.task_name_label"), o.TaskName,
		biz.T("schedule_report.task_type_label"), taskTypeLabel(o.TaskType),
		biz.T("schedule_report.environment"), envInfo,
		biz.T("schedule_report.executor"), statusColor, o.Executor,
		biz.T("schedule_report.exec_time"), o.StartTime, o.EndTime,
		biz.T("schedule_report.duration_label"), durationStr,
		biz.T("schedule_report.expected_total"), o.TotalExpected,
		biz.T("schedule_report.executed"), o.TotalExecuted,
		biz.T("schedule_report.not_executed"), o.NotExecuted,
		biz.T("schedule_report.exec_rate"), o.ExecuteRate)

	return template.HTML(html)
}

// ==================== 饼图：API类型分布 ====================

func buildTaskAPIPie(data biz.TaskReportData) template.HTML {
	if len(data.APITypeDistribution) == 0 {
		return template.HTML(`<div class="col-md-4"><div class="box box-info"><div class="box-header with-border"><h3 class="box-title">` + biz.T("schedule_report.api_type_dist") + `</h3></div><div class="box-body"><div style="text-align:center;padding:20px;color:#aaa">` + biz.T("schedule_report.no_data") + `</div></div></div></div>`)
	}

	var infos, colorNames []string
	var counts []float64
	colors := []chartjs.Color{"rgb(255,205,86)", "rgb(54,162,235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}
	colorNameMap := []string{"yellow", "blue", "red", "green", "black"}

	for i, item := range data.APITypeDistribution {
		infos = append(infos, item.Name)
		counts = append(counts, float64(item.Count))
		c := colors[i%len(colors)]
		cn := colorNameMap[i%len(colorNameMap)]
		colorNames = append(colorNames, string(c))
		_ = cn
	}

	pie := chartjs.Pie().
		SetHeight(180).
		SetLabels(infos).
		SetID("apiPie").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	var labels []map[string]string
	for i, item := range data.APITypeDistribution {
		labels = append(labels, map[string]string{
			"label": fmt.Sprintf(" %s - %d", item.Name, item.Count),
			"color": colorNameMap[i%len(colorNameMap)],
		})
	}
	legend := chart_legend.New().SetData(labels).GetContent()
	boxContent := fmt.Sprintf(`<div class="col-md-8">%s</div><div class="col-md-4">%s</div>`, pie, legend)
	return template.HTML(fmt.Sprintf(`<div class="col-md-4"><div class="box box-info"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body"><div class="row">%s</div></div></div></div>`, biz.T("schedule_report.api_type_dist"), boxContent))
}

// ==================== 饼图：场景执行结果分布 ====================

func buildTaskSceneResultPie(data biz.TaskReportData, taskId string) template.HTML {
	s := data.SceneStats
	if s.Total == 0 {
		return template.HTML(`<div class="col-md-4"><div class="box box-primary"><div class="box-header with-border"><h3 class="box-title">` + biz.T("schedule_report.scene_result_dist") + `</h3></div><div class="box-body"><div style="text-align:center;padding:20px;color:#aaa">` + biz.T("schedule_report.no_scene_data") + `</div></div></div></div>`)
	}

	infos := []string{biz.T("common.pass"), biz.T("common.fail")}
	counts := []float64{float64(s.Pass), float64(s.Fail)}
	colors := []chartjs.Color{"rgb(0, 166, 90)", "rgb(221, 75, 57)"}
	labels := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.pass"), s.Pass), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.fail"), s.Fail), "color": "red"},
	}

	pie := chartjs.Pie().
		SetHeight(180).
		SetLabels(infos).
		SetID("scenePie").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend := chart_legend.New().SetData(labels).GetContent()
	boxContent := fmt.Sprintf(`<div class="col-md-8">%s</div><div class="col-md-4">%s</div>`, pie, legend)
	historyLink := ""
	if len(taskId) > 0 {
		historyLink = fmt.Sprintf(`<a href="/admin/info/scene_test_history?task_id=%s" class="pull-right" target="_blank" style="font-size:12px;margin-top:3px">%s &raquo;</a>`, taskId, biz.T("schedule_report.view_detail"))
	}
	return template.HTML(fmt.Sprintf(`<div class="col-md-4"><div class="box box-primary"><div class="box-header with-border"><h3 class="box-title">%s</h3>%s</div><div class="box-body"><div class="row">%s</div></div></div></div>`, biz.T("schedule_report.scene_result_dist"), historyLink, boxContent))
}

// ==================== 饼图：数据执行结果分布 ====================

func buildTaskDataResultPie(data biz.TaskReportData, taskId string) template.HTML {
	d := data.DataStats
	if d.Total == 0 {
		return template.HTML(`<div class="col-md-4"><div class="box box-success"><div class="box-header with-border"><h3 class="box-title">` + biz.T("schedule_report.data_result_dist") + `</h3></div><div class="box-body"><div style="text-align:center;padding:20px;color:#aaa">` + biz.T("schedule_report.no_data_record") + `</div></div></div></div>`)
	}

	infos := []string{biz.T("common.pass"), biz.T("common.fail")}
	counts := []float64{float64(d.Pass), float64(d.Fail)}
	colors := []chartjs.Color{"rgb(0, 166, 90)", "rgb(221, 75, 57)"}
	labels := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.pass"), d.Pass), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.fail"), d.Fail), "color": "red"},
	}

	pie := chartjs.Pie().
		SetHeight(180).
		SetLabels(infos).
		SetID("dataPie").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend := chart_legend.New().SetData(labels).GetContent()
	boxContent := fmt.Sprintf(`<div class="col-md-8">%s</div><div class="col-md-4">%s</div>`, pie, legend)
	historyLink := ""
	if len(taskId) > 0 {
		historyLink = fmt.Sprintf(`<a href="/admin/info/scene_data_test_history?task_id=%s" class="pull-right" target="_blank" style="font-size:12px;margin-top:3px">%s &raquo;</a>`, taskId, biz.T("schedule_report.view_detail"))
	}
	return template.HTML(fmt.Sprintf(`<div class="col-md-4"><div class="box box-success"><div class="box-header with-border"><h3 class="box-title">%s</h3>%s</div><div class="box-body"><div class="row">%s</div></div></div></div>`, biz.T("schedule_report.data_result_dist"), historyLink, boxContent))
}

// ==================== 执行趋势(折线图) ====================

func buildTaskTrendChart(data biz.TaskReportData) template.HTML {
	if len(data.Trend) == 0 {
		return template.HTML("")
	}

	var dayLabels []string
	var passCounts, failCounts, totalCounts, notExecutedCounts []float64

	for _, t := range data.Trend {
		dayLabels = append(dayLabels, t.ExecutionTime)
		totalCounts = append(totalCounts, float64(t.Total))
		passCounts = append(passCounts, float64(t.Pass))
		failCounts = append(failCounts, float64(t.Fail))
		notExecuted := t.Total - t.Pass - t.Fail
		if notExecuted < 0 {
			notExecuted = 0
		}
		notExecutedCounts = append(notExecutedCounts, float64(notExecuted))
	}

	line := chartjs.Line().
		SetID("trendChart").
		SetHeight(320).
		SetTitle(template.HTML(biz.T("schedule_report.trend_7"))).
		SetLabels(dayLabels).
		AddDataSet(biz.T("common.pass")).
		DSData(passCounts).
		DSFill(false).
		DSBorderColor("rgb(0, 166, 90)").
		DSLineTension(0.1).
		AddDataSet(biz.T("common.fail")).
		DSData(failCounts).
		DSFill(false).
		DSBorderColor("rgb(221, 75, 57)").
		DSLineTension(0.1).
		AddDataSet(biz.T("schedule_report.status_not_executed")).
		DSData(notExecutedCounts).
		DSFill(false).
		DSBorderColor("rgb(169, 169, 169)").
		DSLineTension(0.1).
		AddDataSet(biz.T("schedule_report.total_label")).
		DSData(totalCounts).
		DSFill(false).
		DSBorderColor("rgb(54, 162, 235)").
		DSLineTension(0.1).
		GetContent()

	return template.HTML(fmt.Sprintf(`<div class="box box-info"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body">%s</div></div>`, biz.T("schedule_report.trend_chart"), line))
}

// ==================== 场景明细表 ====================

func buildTaskSceneTable(data biz.TaskReportData) template.HTML {
	if len(data.SceneDetails) == 0 {
		return template.HTML("")
	}

	rows := ""
	for _, s := range data.SceneDetails {
		label := biz.T("common.pass")
		color := "green"
		if s.Result == "fail" {
			color = "red"
			label = biz.T("common.fail")
		} else if s.Result == "未执行" {
			color = "gray"
			label = biz.T("schedule_report.status_not_executed")
		}
		reason := ""
		if len(s.FailReason) > 0 && s.FailReason != " " {
			reason = s.FailReason
		}
		rows += fmt.Sprintf(`<tr><td>%s</td><td style="color:%s">%s</td><td><div class="sc-td" style="max-width:400px"><span class="sc-truncate">%s</span><div class="sc-full">%s</div></div></td></tr>`,
			s.Name, color, label, reason, reason)
	}

	table := fmt.Sprintf(`<table class="table table-bordered table-hover"><thead><tr>
			<th>%s</th><th>%s</th><th>%s</th>
		</tr></thead><tbody>%s</tbody></table>`,
		biz.T("schedule_report.scene_name"), biz.T("schedule_report.test_result"), biz.T("common.fail_reason"), rows)

	return template.HTML(fmt.Sprintf(`<div class="box box-primary"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body">%s</div></div>`, biz.T("schedule_report.scene_detail"), table))
}

// ==================== 数据文件明细表 ====================

func buildTaskDataTable(data biz.TaskReportData) template.HTML {
	if len(data.DataDetails) == 0 {
		return template.HTML("")
	}

	rows := ""
	colorA := "#ffffff"
	colorB := "#f7f9fc"
	currentColor := colorA
	prevScene := ""

	for _, d := range data.DataDetails {
		// 相邻场景切换时交替底色
		if len(prevScene) > 0 && d.SceneName != prevScene {
			if currentColor == colorA {
				currentColor = colorB
			} else {
				currentColor = colorA
			}
		}
		if len(d.SceneName) > 0 {
			prevScene = d.SceneName
		}

		label := biz.T("common.pass")
		color := "green"
		if d.Result == "fail" {
			color = "red"
			label = biz.T("common.fail")
		} else if d.Result == "未执行" {
			color = "gray"
			label = biz.T("schedule_report.status_not_executed")
		}
		reason := ""
		if len(d.FailReason) > 0 && d.FailReason != " " {
			reason = d.FailReason
		}
		rows += fmt.Sprintf(`<tr style="background-color:%s"><td>%s</td><td>%s</td><td>%s</td><td style="color:%s">%s</td><td><div class="sc-td" style="max-width:400px"><span class="sc-truncate">%s</span><div class="sc-full">%s</div></div></td></tr>`,
			currentColor, d.SceneName, d.Name, d.ApiId, color, label, reason, reason)
	}

	table := fmt.Sprintf(`<table class="table table-bordered table-hover"><thead><tr>
			<th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th>
		</tr></thead><tbody>%s</tbody></table>`,
		biz.T("schedule_report.scene_name"), biz.T("schedule_report.data_name"), biz.T("schedule_report.api_id_col"), biz.T("schedule_report.test_result"), biz.T("common.fail_reason"), rows)

	return template.HTML(fmt.Sprintf(`<div class="box box-success"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body">%s</div></div>`, biz.T("schedule_report.data_detail"), table))
}

// ==================== 失败明细表 ====================

func buildTaskFailTable(data biz.TaskReportData) template.HTML {
	if len(data.FailItems) == 0 {
		return template.HTML("")
	}

	rows := ""
	for _, f := range data.FailItems {
		typeStr := biz.T("type.scene")
		if f.Type == "data" {
			typeStr = biz.T("type.data")
		}
		rows += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td><div class="sc-td" style="max-width:400px;color:red"><span class="sc-truncate">%s</span><div class="sc-full" style="color:#333">%s</div></div></td></tr>`,
			f.Name, typeStr, f.APIId, f.Reason, f.Reason)
	}

	table := fmt.Sprintf(`<table class="table table-bordered table-striped"><thead><tr>
			<th>%s</th><th>%s</th><th>%s</th><th>%s</th>
		</tr></thead><tbody>%s</tbody></table>`,
		biz.T("common.name"), biz.T("schedule_report.type_col"), biz.T("schedule_report.api_id_col"), biz.T("common.fail_reason"), rows)

	return template.HTML(fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-danger"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body">%s</div></div></div></div>`, biz.T("schedule_report.fail_detail"), table))
}

// ==================== 通用辅助函数 ====================

func formatDuration(seconds int) string {
	if seconds >= 3600 {
		return biz.T("schedule_report.hour_min_sec", seconds/3600, (seconds%3600)/60, seconds%60)
	} else if seconds >= 60 {
		return biz.T("schedule_report.min_sec", seconds/60, seconds%60)
	}
	return biz.T("schedule_report.sec", seconds)
}

func resultColor(result string) string {
	switch result {
	case "pass":
		return "green"
	case "fail":
		return "red"
	case "partial":
		return "orange"
	}
	return "gray"
}

func taskTypeLabel(t string) string {
	switch t {
	case "scene":
		return biz.T("type.scene")
	case "data":
		return biz.T("type.data")
	}
	return t
}

// renderMultiTaskReport 渲染多任务报告页面
func renderMultiTaskReport(report biz.DashboardReport) (types.Panel, error) {
	var reportData biz.MultiTaskReportData
	err := json.Unmarshal([]byte(report.ReportData), &reportData)
	if err != nil {
		return types.Panel{
			Content:     template.HTML(fmt.Sprintf("<div style='padding:20px;color:red'>%s</div>", biz.T("schedule_report.parse_error", err))),
			Title:       template.HTML(biz.T("schedule_report.page_title")),
			Description: template.HTML(biz.T("schedule_report.description")),
		}, nil
	}

	o := reportData.Overview

	// ====== 顶部任务信息 ======
	headerInfo := buildMultiTaskHeader(reportData, report)

	// ====== 从 ByTask 聚合数据维度和场景维度统计 ======
	var dataTotal, dataPass, dataFail, sceneTotal, scenePass, sceneFail int
	for _, t := range reportData.ByTask {
		dataTotal += t.DataTotal
		dataPass += t.DataPass
		dataFail += t.DataFail
		sceneTotal += t.SceneTotal
		scenePass += t.ScenePass
		sceneFail += t.SceneFail
	}
	dataPassRate := 0.0
	if dataTotal > 0 {
		dataPassRate = float64(dataPass) / float64(dataTotal) * 100
	}
	scenePassRate := 0.0
	if sceneTotal > 0 {
		scenePassRate = float64(scenePass) / float64(sceneTotal) * 100
	}

	// ====== KPI 卡片行1：数据执行数/通过数/失败数/通过率 ======
	kpi_data := fmt.Sprintf(`<div class="row">
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-blue"><i class="fa fa-file-text"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-green"><i class="fa fa-check-circle"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-red"><i class="fa fa-times-circle"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-yellow"><i class="fa fa-percent"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%.1f%%</span></div></div></div>
	</div>`,
		biz.T("schedule_report.data_exec_count"), dataTotal,
		biz.T("schedule_report.pass_count"), dataPass,
		biz.T("schedule_report.fail_count"), dataFail,
		biz.T("schedule_report.pass_rate_label"), dataPassRate)

	// ====== KPI 卡片行2：场景执行数/通过数/失败数/通过率 ======
	kpi_scene := fmt.Sprintf(`<div class="row">
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-blue"><i class="fa fa-cubes"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-green"><i class="fa fa-check-circle"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-red"><i class="fa fa-times-circle"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-yellow"><i class="fa fa-percent"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%.1f%%</span></div></div></div>
	</div>`,
		biz.T("schedule_report.scene_exec_count"), sceneTotal,
		biz.T("schedule_report.pass_count"), scenePass,
		biz.T("schedule_report.fail_count"), sceneFail,
		biz.T("schedule_report.pass_rate_label"), scenePassRate)

	// ====== KPI 卡片行2：任务数/场景数/数据文件数/API数 ======
	kpi2 := fmt.Sprintf(`<div class="row">
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-purple"><i class="fa fa-tasks"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-aqua"><i class="fa fa-play-circle"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-light-blue"><i class="fa fa-file-text"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
		<div class="col-lg-3 col-md-3 col-sm-6 col-xs-12"><div class="info-box"><span class="info-box-icon bg-maroon"><i class="fa fa-plug"></i></span><div class="info-box-content"><span class="info-box-text">%s</span><span class="info-box-number">%d</span></div></div></div>
	</div>`,
		biz.T("schedule_report.task_list_col"), o.TaskCount,
		biz.T("schedule_report.scene_count"), o.SceneCount,
		biz.T("schedule_report.data_count"), o.DataCount,
		biz.T("schedule_report.api_count"), o.APICount)

	// ====== 饼图行：API类型分布 + 场景执行结果 + 数据执行结果 (3并列) ======
	apiPie := buildMultiTaskAPIPie(reportData)
	scenePie := buildMultiTaskSceneResultPie(reportData)
	dataPie := buildMultiTaskDataResultPie(reportData)
	row1 := fmt.Sprintf(`<div class="row">%s%s%s</div>`, scenePie, dataPie, apiPie)

	// ====== 各任务统计表（移到饼图下方） ======
	taskTable := buildMultiTaskStatsTable(reportData)

	// ====== 场景执行明细 ======
	sceneTable := buildMultiTaskSceneTable(reportData)

	// ====== 数据执行明细 ======
	dataTable := buildMultiTaskDataTable(reportData)

	content := headerInfo + template.HTML(kpi2+kpi_scene+kpi_data) + template.HTML(row1) + taskTable + sceneTable + dataTable
	styleBlock := `<style>.sc-td{position:relative;cursor:pointer}.sc-td .sc-truncate{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.sc-td .sc-full{display:none;position:absolute;right:5%;top:50%;background:#fff;border:2px solid #666;padding:15px;z-index:9999;max-width:600px;max-height:80vh;overflow-y:auto;white-space:pre-wrap;word-break:break-all;box-shadow:0 4px 20px rgba(0,0,0,0.3);border-radius:4px;font-size:13px;line-height:1.4}.sc-td:hover .sc-full{display:block!important}</style>`
	content += template.HTML(styleBlock)

	return types.Panel{
		Content:     template.HTML(content),
		Title:       template.HTML(biz.T("schedule_report.page_title")),
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s</span><span style="color:#888">%s: %s</span></div>`, report.ReportName, biz.T("schedule_report.generated_at"), report.CreatedAt)),
	}, nil
}

// buildMultiTaskHeader 多任务报告顶部信息
func buildMultiTaskHeader(data biz.MultiTaskReportData, report biz.DashboardReport) template.HTML {
	o := data.Overview

	// 计算场景维度的预计总数 / 已执行 / 未执行
	var sceneTotalExpected, sceneTotalExecuted int
	for _, t := range data.ByTask {
		sceneTotalExpected += t.SceneTotal
		sceneTotalExecuted += t.ScenePass + t.SceneFail
	}
	notExecuted := sceneTotalExpected - sceneTotalExecuted
	if notExecuted < 0 {
		notExecuted = 0
	}
	executeRate := 0.0
	if sceneTotalExpected > 0 {
		executeRate = float64(sceneTotalExecuted) / float64(sceneTotalExpected) * 100
	}

	// 执行时间
	startTime := o.StartTime
	endTime := o.EndTime
	if len(startTime) == 0 {
		startTime = biz.T("schedule_report.unknown")
	}
	if len(endTime) == 0 {
		endTime = biz.T("schedule_report.unknown")
	}
	timeRange := fmt.Sprintf("%s ~ %s", startTime, endTime)

	// 执行耗时
	durationStr := ""
	if o.DurationSeconds > 0 {
		durationStr = formatDuration(o.DurationSeconds)
	} else {
		durationStr = biz.T("schedule_report.unknown")
	}

	// 执行率
	execRateDisplay := ""
	if o.ExecuteRate > 0 {
		execRateDisplay = fmt.Sprintf("%.1f%%", o.ExecuteRate)
	} else if executeRate > 0 {
		execRateDisplay = fmt.Sprintf("%.1f%%", executeRate)
	} else {
		execRateDisplay = biz.T("schedule_report.unknown")
	}

	html := fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-default" style="border-top:none;box-shadow:none;background:#fff;margin-bottom:10px"><div class="box-body" style="padding:15px 20px">
		<div class="row" style="margin-bottom:8px">
			<div class="col-md-4"><strong>%s</strong> %s</div>
			<div class="col-md-2"><strong>%s</strong> %d</div>
			<div class="col-md-4"><strong>%s</strong> %s</div>
			<div class="col-md-2"><strong>%s</strong> %s</div>
		</div>
		<div class="row">
			<div class="col-md-4"><strong>%s</strong> %s</div>
			<div class="col-md-2"><strong>%s</strong> %s</div>
			<div class="col-md-4"><strong>%s</strong> %d | <strong>%s</strong> %d | <strong>%s</strong> %d</div>
			<div class="col-md-2"><strong>%s</strong> <span style="color:green">%s</span></div>
		</div>
	</div></div></div></div>`,
		biz.T("schedule_report.task_name_label"), report.ReportName,
		biz.T("schedule_report.task_count_label"), o.TaskCount,
		biz.T("schedule_report.environment"), o.Product,
		biz.T("schedule_report.executor"), report.Creator,
		biz.T("schedule_report.exec_time"), timeRange,
		biz.T("schedule_report.duration_label"), durationStr,
		biz.T("schedule_report.expected_total"), sceneTotalExpected,
		biz.T("schedule_report.executed"), sceneTotalExecuted,
		biz.T("schedule_report.not_executed"), notExecuted,
		biz.T("schedule_report.exec_rate"), execRateDisplay)

	return template.HTML(html)
}

// buildMultiTaskStatsTable 构建各任务统计表
func buildMultiTaskStatsTable(data biz.MultiTaskReportData) template.HTML {
	if len(data.ByTask) == 0 {
		return template.HTML("")
	}
	rows := ""
	for _, t := range data.ByTask {
		passRate := ""
		if t.Total > 0 {
			passRate = fmt.Sprintf("%.1f%%", t.PassRate)
		}
		taskTypeLabel := taskTypeLabel(t.TaskType)
		timeStr := t.StartTime + " ~ " + t.EndTime
		if len(t.StartTime) == 0 || len(t.EndTime) == 0 {
			timeStr = "-"
		}
		durStr := ""
		if t.DurationSeconds > 0 {
			durStr = formatDuration(t.DurationSeconds)
		} else {
			durStr = "-"
		}
		rows += fmt.Sprintf(`<tr>
			<td>%s</td><td>%s</td>
			<td>%s</td><td>%s</td>
			<td>%d / %d / %d</td>
			<td>%d / %d / %d</td>
			<td>%d</td><td>%s</td>
		</tr>`, t.TaskName, taskTypeLabel, timeStr, durStr,
			t.ScenePass, t.SceneFail, t.SceneTotal,
			t.DataPass, t.DataFail, t.DataTotal,
			t.Total, passRate)
	}
	tableHTML := fmt.Sprintf(`<table class="table table-bordered table-striped"><thead><tr>
		<th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th>
	</tr></thead><tbody>%s</tbody></table>`,
		biz.T("common.task_name"), biz.T("common.task_type"),
		biz.T("schedule_report.exec_time_col"), biz.T("schedule_report.duration_col"),
		biz.T("schedule_report.scene_pass_fail_col"), biz.T("schedule_report.data_pass_fail_col"),
		biz.T("schedule_report.total_exec"), biz.T("schedule_report.pass_rate_label"), rows)
	return template.HTML(fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-primary"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body">%s</div></div></div></div>`,
		biz.T("schedule_report.task_detail"), tableHTML))
}

// buildMultiTaskAPIPie 构建多任务聚合的API类型分布饼图
func buildMultiTaskAPIPie(data biz.MultiTaskReportData) template.HTML {
	if len(data.APITypeDistribution) == 0 {
		return template.HTML("")
	}
	var infos []string
	var counts []float64
	colorNames := []string{"yellow", "blue", "red", "green", "black"}
	colors := []chartjs.Color{"rgb(255,205,86)", "rgb(54,162,235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}
	for i, item := range data.APITypeDistribution {
		infos = append(infos, item.Name)
		counts = append(counts, float64(item.Count))
		_ = colors[i%len(colors)]
		_ = colorNames[i%len(colorNames)]
	}
	pie := chartjs.Pie().SetHeight(180).SetLabels(infos).SetID("multiApiPie").
		AddDataSet(infos[0]).DSData(counts).DSBackgroundColor(colors).GetContent()
	var labels []map[string]string
	for i, item := range data.APITypeDistribution {
		labels = append(labels, map[string]string{
			"label": fmt.Sprintf(" %s - %d", item.Name, item.Count),
			"color": colorNames[i%len(colorNames)],
		})
	}
	legend := chart_legend.New().SetData(labels).GetContent()
	boxContent := fmt.Sprintf(`<div class="col-md-8">%s</div><div class="col-md-4">%s</div>`, pie, legend)
	return template.HTML(fmt.Sprintf(`<div class="row"><div class="col-md-4"><div class="box box-info"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body"><div class="row">%s</div></div></div></div></div>`, biz.T("schedule_report.api_type_dist"), boxContent))
}

// buildMultiTaskSceneResultPie 多任务场景执行结果分布饼图
func buildMultiTaskSceneResultPie(data biz.MultiTaskReportData) template.HTML {
	pass, fail := 0, 0
	for _, s := range data.SceneDetails {
		if s.Result == "pass" {
			pass++
		} else if s.Result == "fail" {
			fail++
		}
	}
	total := pass + fail
	if total == 0 {
		return template.HTML(`<div class="col-md-4"><div class="box box-primary"><div class="box-header with-border"><h3 class="box-title">` + biz.T("schedule_report.scene_result_dist") + `</h3></div><div class="box-body"><div style="text-align:center;padding:20px;color:#aaa">` + biz.T("schedule_report.no_scene_data") + `</div></div></div></div>`)
	}

	infos := []string{biz.T("common.pass"), biz.T("common.fail")}
	counts := []float64{float64(pass), float64(fail)}
	colors := []chartjs.Color{"rgb(0, 166, 90)", "rgb(221, 75, 57)"}
	labels := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.pass"), pass), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.fail"), fail), "color": "red"},
	}

	pie := chartjs.Pie().
		SetHeight(180).
		SetLabels(infos).
		SetID("multiScenePie").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend := chart_legend.New().SetData(labels).GetContent()
	boxContent := fmt.Sprintf(`<div class="col-md-8">%s</div><div class="col-md-4">%s</div>`, pie, legend)
	return template.HTML(fmt.Sprintf(`<div class="col-md-4"><div class="box box-primary"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body"><div class="row">%s</div></div></div></div>`, biz.T("schedule_report.scene_result_dist"), boxContent))
}

// buildMultiTaskDataResultPie 多任务数据执行结果分布饼图
func buildMultiTaskDataResultPie(data biz.MultiTaskReportData) template.HTML {
	pass, fail := 0, 0
	for _, d := range data.DataDetails {
		if d.Result == "pass" {
			pass++
		} else if d.Result == "fail" {
			fail++
		}
	}
	total := pass + fail
	if total == 0 {
		return template.HTML(`<div class="col-md-4"><div class="box box-success"><div class="box-header with-border"><h3 class="box-title">` + biz.T("schedule_report.data_result_dist") + `</h3></div><div class="box-body"><div style="text-align:center;padding:20px;color:#aaa">` + biz.T("schedule_report.no_data_record") + `</div></div></div></div>`)
	}

	infos := []string{biz.T("common.pass"), biz.T("common.fail")}
	counts := []float64{float64(pass), float64(fail)}
	colors := []chartjs.Color{"rgb(0, 166, 90)", "rgb(221, 75, 57)"}
	labels := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.pass"), pass), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.fail"), fail), "color": "red"},
	}

	pie := chartjs.Pie().
		SetHeight(180).
		SetLabels(infos).
		SetID("multiDataPie").
		AddDataSet(infos[0]).
		DSData(counts).
		DSBackgroundColor(colors).
		GetContent()

	legend := chart_legend.New().SetData(labels).GetContent()
	boxContent := fmt.Sprintf(`<div class="col-md-8">%s</div><div class="col-md-4">%s</div>`, pie, legend)
	return template.HTML(fmt.Sprintf(`<div class="col-md-4"><div class="box box-success"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body"><div class="row">%s</div></div></div></div>`, biz.T("schedule_report.data_result_dist"), boxContent))
}

// buildMultiTaskSceneTable 构建多任务聚合的场景/数据明细表
func buildMultiTaskSceneTable(data biz.MultiTaskReportData) template.HTML {
	if len(data.SceneDetails) == 0 {
		return template.HTML("")
	}
	rows := ""
	for _, s := range data.SceneDetails {
		label := biz.T("common.pass")
		color := "green"
		if s.Result == "fail" {
			color = "red"
			label = biz.T("common.fail")
		} else if s.Result == "未执行" {
			color = "gray"
			label = biz.T("schedule_report.status_not_executed")
		}
		reason := ""
		if len(s.FailReason) > 0 && s.FailReason != " " {
			reason = s.FailReason
		}
		rows += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td style="color:%s">%s</td><td><div class="sc-td" style="max-width:400px"><span class="sc-truncate">%s</span><div class="sc-full">%s</div></div></td></tr>`, s.TaskName, s.Name, color, label, reason, reason)
	}
	tableHTML := fmt.Sprintf(`<table class="table table-bordered table-hover"><thead><tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th></tr></thead><tbody>%s</tbody></table>`,
		biz.T("common.task_name"), biz.T("schedule_report.scene_name"), biz.T("schedule_report.test_result"), biz.T("common.fail_reason"), rows)
	return template.HTML(fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-primary"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body">%s</div></div></div></div>`, biz.T("schedule_report.scene_detail"), tableHTML))
}

// buildMultiTaskDataTable 构建多任务聚合的数据文件明细表
func buildMultiTaskDataTable(data biz.MultiTaskReportData) template.HTML {
	if len(data.DataDetails) == 0 {
		return template.HTML("")
	}
	rows := ""
	for _, d := range data.DataDetails {
		label := biz.T("common.pass")
		color := "green"
		if d.Result == "fail" {
			color = "red"
			label = biz.T("common.fail")
		}
		reason := ""
		if len(d.FailReason) > 0 && d.FailReason != " " {
			reason = d.FailReason
		}
		rows += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td style="color:%s">%s</td><td><div class="sc-td" style="max-width:400px"><span class="sc-truncate">%s</span><div class="sc-full">%s</div></div></td></tr>`, d.TaskName, d.Name, d.ApiId, color, label, reason, reason)
	}
	tableHTML := fmt.Sprintf(`<table class="table table-bordered table-hover"><thead><tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th></tr></thead><tbody>%s</tbody></table>`,
		biz.T("common.task_name"), biz.T("schedule_report.data_name"), biz.T("schedule_report.api_id_col"), biz.T("schedule_report.test_result"), biz.T("common.fail_reason"), rows)
	return template.HTML(fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-success"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body">%s</div></div></div></div>`, biz.T("schedule_report.data_detail"), tableHTML))
}

// buildMultiTaskFailTable 构建多任务聚合的失败明细表
func buildMultiTaskFailTable(data biz.MultiTaskReportData) template.HTML {
	if len(data.FailItems) == 0 {
		return template.HTML("")
	}
	rows := ""
	for _, f := range data.FailItems {
		typeStr := biz.T("common.scene")
		if f.Type == "data" {
			typeStr = biz.T("common.data")
		}
		rows += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td><div class="sc-td" style="max-width:400px;color:red"><span class="sc-truncate">%s</span><div class="sc-full" style="color:#333">%s</div></div></td></tr>`, f.Name, typeStr, f.APIId, f.Reason, f.Reason)
	}
	tableHTML := fmt.Sprintf(`<table class="table table-bordered table-striped"><thead><tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th></tr></thead><tbody>%s</tbody></table>`,
		biz.T("schedule_report.name_with_task"), biz.T("schedule_report.type_col"), biz.T("schedule_report.api_id_col"), biz.T("common.fail_reason"), rows)
	return template.HTML(fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-danger"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body">%s</div></div></div></div>`, biz.T("schedule_report.fail_detail"), tableHTML))
}

// buildMultiTaskResourceBox 构建资源关联统计区（场景/数据/API 三个列表卡片并排）
func buildMultiTaskResourceBox(data biz.MultiTaskReportData) template.HTML {
	sceneItems := ""
	if len(data.SceneList) > 0 {
		for _, s := range data.SceneList {
			sceneItems += fmt.Sprintf("<li>%s</li>", s.Name)
		}
	} else {
		sceneItems = "<li style='color:#aaa'>" + biz.T("schedule_report.no_related_scene") + "</li>"
	}
	sceneBox := fmt.Sprintf(`<div class="col-md-4"><div class="box box-default"><div class="box-header with-border"><h3 class="box-title">%s (%d)</h3></div><div class="box-body" style="max-height:300px;overflow-y:auto"><ul style="padding-left:20px">%s</ul></div></div></div>`,
		biz.T("schedule_report.related_scenes"), len(data.SceneList), sceneItems)

	dataItems := ""
	if len(data.DataList) > 0 {
		for _, d := range data.DataList {
			dataItems += fmt.Sprintf("<li>%s</li>", d.Name)
		}
	} else {
		dataItems = "<li style='color:#aaa'>" + biz.T("schedule_report.no_related_data") + "</li>"
	}
	dataBox := fmt.Sprintf(`<div class="col-md-4"><div class="box box-default"><div class="box-header with-border"><h3 class="box-title">%s (%d)</h3></div><div class="box-body" style="max-height:300px;overflow-y:auto"><ul style="padding-left:20px">%s</ul></div></div></div>`,
		biz.T("schedule_report.related_data"), len(data.DataList), dataItems)

	apiItems := ""
	if len(data.APIList) > 0 {
		for _, a := range data.APIList {
			apiItems += fmt.Sprintf("<li>%s</li>", a.Name)
		}
	} else {
		apiItems = "<li style='color:#aaa'>" + biz.T("schedule_report.no_related_api") + "</li>"
	}
	apiBox := fmt.Sprintf(`<div class="col-md-4"><div class="box box-default"><div class="box-header with-border"><h3 class="box-title">%s (%d)</h3></div><div class="box-body" style="max-height:300px;overflow-y:auto"><ul style="padding-left:20px">%s</ul></div></div></div>`,
		biz.T("schedule_report.related_apis"), len(data.APIList), apiItems)

	return template.HTML(fmt.Sprintf(`<div class="row"><div class="col-md-12"><div class="box box-info"><div class="box-header with-border"><h3 class="box-title">%s</h3></div><div class="box-body"><div class="row">%s%s%s</div></div></div></div></div>`,
		biz.T("schedule_report.resource_stats"), sceneBox, dataBox, apiBox))
}
