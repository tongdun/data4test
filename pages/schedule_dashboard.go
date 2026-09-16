package pages

import (
	"data4test/biz"
	"encoding/json"
	"fmt"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
	"time"

	"data4test/models"
	"github.com/GoAdminGroup/go-admin/template/types"
)

func GetDashboardByReportId(ctx *gin.Context) (types.Panel, error) {
	reportId := ctx.Query("id")
	cookie, _ := ctx.Cookie("data_locale")
	dataLocale := biz.GetDataLocale(cookie, biz.GetLocale())
	var report biz.DashboardReport
	models.Orm.Table("dashboard").
		Where("id = ?", reportId).
		Find(&report)
	switch report.ReportType {
	case "task":
		if strings.Contains(report.RelatedTaskIds, ",") {
			return renderMultiTaskReport(report, dataLocale)
		}

		return renderTaskReport(report, "", dataLocale)
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
}

func GetScheduleReportContent(ctx *gin.Context) (types.Panel, error) {
	taskId := ctx.Query("id")
	cookie, _ := ctx.Cookie("data_locale")
	dataLocale := biz.GetDataLocale(cookie, biz.GetLocale())
	var report biz.DashboardReport
	if len(taskId) > 0 {
		models.Orm.Table("dashboard").
			Where("report_type = ? and status = ? and related_task_ids like CONCAT(?, '\\_%')", "task", "finished", taskId).
			Order("created_at desc").
			Limit(1).
			Find(&report)
		if len(report.Id) == 0 || len(report.ReportData) == 0 {
			// 无执行报告 → 自动生成覆盖范围报告
			coverageData, err := biz.BuildTaskCoverageReportData(taskId)
			if err != nil {
				biz.Logger.Warning("生成覆盖范围报告失败: %s", err)
				return types.Panel{
					Content: template.HTML(`<div style='padding:40px;text-align:center'>
							<h3>` + biz.T("schedule_report.no_report") + `</h3>
							<p style='color:#888;margin-top:20px'>` + biz.T("schedule_report.view_report_empty") + `</p>
						</div>`),
					Title:       template.HTML(biz.T("schedule_report.page_title")),
					Description: template.HTML(biz.T("schedule_report.description")),
				}, nil
			}
			jsonBytes, _ := json.Marshal(coverageData)
			fakeReport := biz.DashboardReport{
				ReportName: coverageData.Overview.TaskName,
				ReportType: "task",
				ReportData: string(jsonBytes),
				Status:     "coverage",
				Creator:    coverageData.Overview.Executor,
				CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
			}
			return renderTaskReport(fakeReport, taskId, dataLocale)
		}

	}
	return renderTaskReport(report, taskId, dataLocale)
}

// renderTaskReport 解析 TaskReportData JSON 并渲染完整报告页面
func renderTaskReport(report biz.DashboardReport, scheduleId string, dataLocale string) (types.Panel, error) {
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
	headerInfo := buildTaskHeader(reportData, dataLocale)

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

	// ====== 第3行: 执行明细表(场景+数据合并, 折叠展开, 独占一行) ======
	sceneTable := buildTaskSceneGroupTable(reportData, dataLocale)
	row3 := ""
	if len(sceneTable) > 0 {
		row3 = fmt.Sprintf(`<div class="row"><div class="col-md-12">%s</div></div>`, sceneTable)
	}

	// ====== 第4行: 失败明细表 ======
	failHTML := buildTaskFailTable(reportData, dataLocale)

	content := string(headerInfo) + string(kpiCards) + row1 + row2 + row3 + string(failHTML)
	content += string(reportStyle()) + reportCollapseScript()

	return types.Panel{
		Content:     template.HTML(content),
		Title:       template.HTML(biz.T("schedule_report.page_title")),
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s</span><span style="color:#888">%s: %s</span></div>`, report.ReportName, biz.T("schedule_report.generated_at"), report.CreatedAt)),
	}, nil
}

// ==================== KPI卡片 ====================

func buildTaskKpiCards(data biz.TaskReportData) template.HTML {
	s := data.SceneStats
	d := data.DataStats
	scenePassRate := fmt.Sprintf("%.1f%%", s.PassRate)
	dataPassRate := fmt.Sprintf("%.1f%%", d.PassRate)

	// 第一行：场景执行统计（执行数/通过数/失败数/通过率）
	row1 := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		reportInfoBox(fmt.Sprintf("%d", s.Total), biz.T("schedule_report.scene_exec_count"), "fa-cubes", "blue"),
		reportInfoBox(fmt.Sprintf("%d", s.Pass), biz.T("schedule_report.pass_count"), "fa-check-circle", "green"),
		reportInfoBox(fmt.Sprintf("%d", s.Fail), biz.T("schedule_report.fail_count"), "fa-times-circle", "red"),
		reportInfoBox(scenePassRate, biz.T("schedule_report.pass_rate_label"), "fa-percent", "green"))

	// 第二行：数据执行统计（执行数/通过数/失败数/通过率）
	row2 := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		reportInfoBox(fmt.Sprintf("%d", d.Total), biz.T("schedule_report.data_exec_count"), "fa-file-text", "blue"),
		reportInfoBox(fmt.Sprintf("%d", d.Pass), biz.T("schedule_report.pass_count"), "fa-check-circle", "green"),
		reportInfoBox(fmt.Sprintf("%d", d.Fail), biz.T("schedule_report.fail_count"), "fa-times-circle", "red"),
		reportInfoBox(dataPassRate, biz.T("schedule_report.pass_rate_label"), "fa-percent", "green"))

	return template.HTML(row1 + row2)
}

// buildTaskHeader 顶部执行信息（白色背景，两行排版）
func buildTaskHeader(data biz.TaskReportData, dataLocale string) template.HTML {
	o := data.Overview
	envInfo := o.Environment
	if len(envInfo) == 0 {
		envInfo = biz.T("schedule_report.unknown")
	}
	durationStr := formatDuration(o.DurationSeconds)
	statusColor := resultColor(o.Executor)

	html := fmt.Sprintf(`<div class="row" style="margin-bottom:8px">
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
		</div>`,
		biz.T("schedule_report.task_name_label"), biz.GetTaskLocalized(o.TaskName, dataLocale),
		biz.T("schedule_report.task_type_label"), taskTypeLabel(o.TaskType),
		biz.T("schedule_report.environment"), envInfo,
		biz.T("schedule_report.executor"), statusColor, o.Executor,
		biz.T("schedule_report.exec_time"), o.StartTime, o.EndTime,
		biz.T("schedule_report.duration_label"), durationStr,
		biz.T("schedule_report.expected_total"), o.TotalExpected,
		biz.T("schedule_report.executed"), o.TotalExecuted,
		biz.T("schedule_report.not_executed"), o.NotExecuted,
		biz.T("schedule_report.exec_rate"), o.ExecuteRate)

	return template.HTML(reportCard("", "", html))
}

// ==================== 饼图：API类型分布 ====================

func buildTaskAPIPie(data biz.TaskReportData) template.HTML {
	title := biz.T("schedule_report.api_type_dist")
	if len(data.APITypeDistribution) == 0 {
		return template.HTML(`<div class="col-md-4">` + reportCard(title, "", fmt.Sprintf(`<div style="text-align:center;padding:20px;color:#aaa">%s</div>`, biz.T("schedule_report.no_data"))) + `</div>`)
	}

	labels := make([]string, 0, len(data.APITypeDistribution))
	counts := make([]float64, 0, len(data.APITypeDistribution))
	colors := make([]chartjs.Color, 0, len(data.APITypeDistribution))
	legendItems := make([]map[string]string, 0, len(data.APITypeDistribution))
	for i, item := range data.APITypeDistribution {
		labels = append(labels, item.Name)
		counts = append(counts, float64(item.Count))
		colors = append(colors, identityColor(i))
		legendItems = append(legendItems, map[string]string{
			"label": fmt.Sprintf(" %s - %d", item.Name, item.Count),
			"color": identityLegendColor(i),
		})
	}
	return template.HTML(reportPie("apiPie", title, "", labels, counts, colors, legendItems))
}

// ==================== 饼图：场景执行结果分布 ====================

func buildTaskSceneResultPie(data biz.TaskReportData, taskId string) template.HTML {
	title := biz.T("schedule_report.scene_result_dist")
	s := data.SceneStats
	if s.Total == 0 {
		return template.HTML(`<div class="col-md-4">` + reportCard(title, "", fmt.Sprintf(`<div style="height:180px;display:flex;align-items:center;justify-content:center;color:#aaa">%s</div>`, biz.T("schedule_report.no_scene_data"))) + `</div>`)
	}

	labels := []string{biz.T("common.pass"), biz.T("common.fail")}
	counts := []float64{float64(s.Pass), float64(s.Fail)}
	colors := []chartjs.Color{colorPass, colorFail}
	legendItems := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.pass"), s.Pass), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.fail"), s.Fail), "color": "red"},
	}

	extra := ""
	if len(taskId) > 0 {
		extra = fmt.Sprintf(`<a href="/admin/info/scene_test_history?task_id=%s" target="_blank" style="font-size:12px">%s &raquo;</a>`, taskId, biz.T("schedule_report.view_detail"))
	}
	return template.HTML(reportPie("scenePie", title, extra, labels, counts, colors, legendItems))
}

// ==================== 饼图：数据执行结果分布 ====================

func buildTaskDataResultPie(data biz.TaskReportData, taskId string) template.HTML {
	title := biz.T("schedule_report.data_result_dist")
	d := data.DataStats
	if d.Total == 0 {
		return template.HTML(`<div class="col-md-4">` + reportCard(title, "", fmt.Sprintf(`<div style="height:180px;display:flex;align-items:center;justify-content:center;color:#aaa">%s</div>`, biz.T("schedule_report.no_data_record"))) + `</div>`)
	}

	labels := []string{biz.T("common.pass"), biz.T("common.fail")}
	counts := []float64{float64(d.Pass), float64(d.Fail)}
	colors := []chartjs.Color{colorPass, colorFail}
	legendItems := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.pass"), d.Pass), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.fail"), d.Fail), "color": "red"},
	}

	extra := ""
	if len(taskId) > 0 {
		extra = fmt.Sprintf(`<a href="/admin/info/scene_data_test_history?task_id=%s" target="_blank" style="font-size:12px">%s &raquo;</a>`, taskId, biz.T("schedule_report.view_detail"))
	}
	return template.HTML(reportPie("dataPie", title, extra, labels, counts, colors, legendItems))
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

	series := []reportLineSeries{
		{Label: biz.T("common.pass"), Data: passCounts, Color: colorPass},
		{Label: biz.T("common.fail"), Data: failCounts, Color: colorFail},
		{Label: biz.T("schedule_report.status_not_executed"), Data: notExecutedCounts, Color: colorUntest},
		{Label: biz.T("schedule_report.total_label"), Data: totalCounts, Color: colorInfo},
	}
	return template.HTML(reportLine("trendChart", biz.T("schedule_report.trend_chart"), dayLabels, series))
}

// ==================== 执行明细表（场景+数据合并, 可折叠展开） ====================

// buildTaskSceneGroupTable 将场景明细与其关联数据合并为一张可折叠展开的「执行明细」表：
// 每个场景一行头行，点击展开其关联数据子行；默认全部收起。
func buildTaskSceneGroupTable(data biz.TaskReportData, dataLocale string) template.HTML {
	if len(data.SceneDetails) == 0 && len(data.DataDetails) == 0 {
		return template.HTML("")
	}

	// 数据明细按场景名分组
	grouped := make(map[string][]biz.DataDetail)
	for _, d := range data.DataDetails {
		grouped[d.SceneName] = append(grouped[d.SceneName], d)
	}

	headers := []string{
		biz.T("schedule_report.scene_data_col"),       // 场景 / 数据
		biz.T("schedule_report.api_id_col"),           // API_ID
		biz.T("schedule_report.test_result"),          // 执行结果
		biz.T("schedule_report.related_data_count"),   // 关联数据
		biz.T("common.fail_reason"),                   // 失败原因
	}

	expandAll := fmt.Sprintf(`<span class="sc-btn" onclick="setAllScenes(true)">%s</span><span class="sc-btn" onclick="setAllScenes(false)">%s</span>`,
		biz.T("schedule_report.expand_all"), biz.T("schedule_report.collapse_all"))

	var rows []string
	idx := 0

	// 场景头行 + 其关联数据子行
	for _, s := range data.SceneDetails {
		children := grouped[s.Name]
		delete(grouped, s.Name)
		pass, fail := countDataResult(children)
		rows = append(rows, buildSceneGroupHeaderRow(idx, biz.GetPlaybookLocalized(s.Name, dataLocale), s.Result, len(children), pass, fail, s.FailReason))
		for _, c := range children {
			rows = append(rows, buildSceneGroupChildRow(idx, c, dataLocale))
		}
		idx++
	}

	// 孤立数据：SceneName 未匹配到任何场景（含 scene_name 为空），兜底归入一个分组避免丢数据
	for sceneName, children := range grouped {
		if len(children) == 0 {
			continue
		}
		pass, fail := countDataResult(children)
		orphanResult := "未执行"
		if fail > 0 {
			orphanResult = "fail"
		} else if pass > 0 {
			orphanResult = "pass"
		}
		displayName := sceneName
		if len(displayName) == 0 {
			displayName = biz.T("schedule_report.orphan_data")
		} else {
			displayName = biz.GetPlaybookLocalized(displayName, dataLocale)
		}
		rows = append(rows, buildSceneGroupHeaderRow(idx, displayName, orphanResult, len(children), pass, fail, ""))
		for _, c := range children {
			rows = append(rows, buildSceneGroupChildRow(idx, c, dataLocale))
		}
		idx++
	}

	return template.HTML(reportCard(biz.T("schedule_report.exec_detail"), expandAll, reportTable(headers, rows)))
}

// buildSceneGroupHeaderRow 渲染一个场景头行（可点击展开其关联数据子行）。
func buildSceneGroupHeaderRow(idx int, name, result string, childCount, pass, fail int, reason string) string {
	return fmt.Sprintf(`<tr class="scene-header" data-group="%d" onclick="toggleSceneGroup(%d, this)"><td><span class="scene-arrow">▸</span> <span class="scene-name">%s</span></td><td>—</td><td>%s</td><td class="sc-count">%s</td>%s</tr>`,
		idx, idx, name, sceneResultBadge(result), biz.T("schedule_report.related_data_summary", childCount, pass, fail), sceneReasonCell(reason))
}

// buildSceneGroupChildRow 渲染场景下的一个关联数据子行。
func buildSceneGroupChildRow(idx int, d biz.DataDetail, dataLocale string) string {
	apiID := d.ApiId
	if len(apiID) == 0 {
		apiID = "—"
	}
	return fmt.Sprintf(`<tr class="scene-child" data-group="%d"><td class="sc-data-name">%s</td><td class="sc-api">%s</td><td>%s</td><td>—</td>%s</tr>`,
		idx, biz.GetDataLocalized(d.Name, dataLocale), apiID, sceneResultBadge(d.Result), sceneReasonCell(d.FailReason))
}

// sceneResultBadge 结果徽标（通过绿 / 失败红 / 未执行灰）。
func sceneResultBadge(result string) string {
	label := biz.T("common.pass")
	cls := "sc-badge sc-badge-pass"
	if result == "fail" {
		label = biz.T("common.fail")
		cls = "sc-badge sc-badge-fail"
	} else if result == "未执行" {
		label = biz.T("schedule_report.status_not_executed")
		cls = "sc-badge sc-badge-na"
	}
	return fmt.Sprintf(`<span class="%s">%s</span>`, cls, label)
}

// sceneReasonCell 失败原因单元格（截断 + hover 全文）。
func sceneReasonCell(reason string) string {
	if len(reason) == 0 || reason == " " {
		return `<td>—</td>`
	}
	return fmt.Sprintf(`<td><div class="sc-td" style="max-width:400px"><span class="sc-truncate">%s</span><div class="sc-full">%s</div></div></td>`, reason, reason)
}

// countDataResult 统计一组数据记录的通过/失败数。
func countDataResult(children []biz.DataDetail) (pass, fail int) {
	for _, c := range children {
		if c.Result == "pass" {
			pass++
		} else if c.Result == "fail" {
			fail++
		}
	}
	return
}

// ==================== 失败明细表 ====================

func buildTaskFailTable(data biz.TaskReportData, dataLocale string) template.HTML {
	if len(data.FailItems) == 0 {
		return template.HTML("")
	}

	headers := []string{
		biz.T("common.name"),
		biz.T("schedule_report.type_col"),
		biz.T("schedule_report.api_id_col"),
		biz.T("common.fail_reason"),
	}
	var rows []string
	for _, f := range data.FailItems {
		typeStr := biz.T("type.scene")
		nameStr := biz.GetPlaybookLocalized(f.Name, dataLocale)
		if f.Type == "data" {
			typeStr = biz.T("type.data")
			nameStr = biz.GetDataLocalized(f.Name, dataLocale)
		}
		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td><div class="sc-td" style="max-width:400px;color:red"><span class="sc-truncate">%s</span><div class="sc-full" style="color:#333">%s</div></div></td></tr>`,
			nameStr, typeStr, f.APIId, f.Reason, f.Reason))
	}
	return template.HTML(reportCard(biz.T("schedule_report.fail_detail"), "", reportTable(headers, rows)))
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
func renderMultiTaskReport(report biz.DashboardReport, dataLocale string) (types.Panel, error) {
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
	kpi_data := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		reportInfoBox(fmt.Sprintf("%d", dataTotal), biz.T("schedule_report.data_exec_count"), "fa-file-text", "blue"),
		reportInfoBox(fmt.Sprintf("%d", dataPass), biz.T("schedule_report.pass_count"), "fa-check-circle", "green"),
		reportInfoBox(fmt.Sprintf("%d", dataFail), biz.T("schedule_report.fail_count"), "fa-times-circle", "red"),
		reportInfoBox(fmt.Sprintf("%.1f%%", dataPassRate), biz.T("schedule_report.pass_rate_label"), "fa-percent", "green"))

	// ====== KPI 卡片行2：场景执行数/通过数/失败数/通过率 ======
	kpi_scene := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		reportInfoBox(fmt.Sprintf("%d", sceneTotal), biz.T("schedule_report.scene_exec_count"), "fa-cubes", "blue"),
		reportInfoBox(fmt.Sprintf("%d", scenePass), biz.T("schedule_report.pass_count"), "fa-check-circle", "green"),
		reportInfoBox(fmt.Sprintf("%d", sceneFail), biz.T("schedule_report.fail_count"), "fa-times-circle", "red"),
		reportInfoBox(fmt.Sprintf("%.1f%%", scenePassRate), biz.T("schedule_report.pass_rate_label"), "fa-percent", "green"))

	// ====== KPI 卡片行3：任务数/场景数/数据文件数/API数 ======
	kpi2 := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		reportInfoBox(fmt.Sprintf("%d", o.TaskCount), biz.T("schedule_report.task_list_col"), "fa-tasks", "purple"),
		reportInfoBox(fmt.Sprintf("%d", o.SceneCount), biz.T("schedule_report.scene_count"), "fa-play-circle", "aqua"),
		reportInfoBox(fmt.Sprintf("%d", o.DataCount), biz.T("schedule_report.data_count"), "fa-file-text", "lightblue"),
		reportInfoBox(fmt.Sprintf("%d", o.APICount), biz.T("schedule_report.api_count"), "fa-plug", "maroon"))

	// ====== 饼图行：API类型分布 + 场景执行结果 + 数据执行结果 (3并列) ======
	apiPie := buildMultiTaskAPIPie(reportData)
	scenePie := buildMultiTaskSceneResultPie(reportData)
	dataPie := buildMultiTaskDataResultPie(reportData)
	row1 := fmt.Sprintf(`<div class="row">%s%s%s</div>`, scenePie, dataPie, apiPie)

	// ====== 各任务统计表（移到饼图下方） ======
	taskTable := buildMultiTaskStatsTable(reportData, dataLocale)

	// ====== 场景执行明细 ======
	sceneTable := buildMultiTaskSceneTable(reportData, dataLocale)

	// ====== 数据执行明细 ======
	dataTable := buildMultiTaskDataTable(reportData, dataLocale)

	content := headerInfo + template.HTML(kpi2+kpi_scene+kpi_data) + template.HTML(row1) + taskTable + sceneTable + dataTable
	content += reportStyle()

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

	html := fmt.Sprintf(`<div class="row" style="margin-bottom:8px">
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
	</div>`,
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

	return template.HTML(reportCard("", "", html))
}

// buildMultiTaskStatsTable 构建各任务统计表
func buildMultiTaskStatsTable(data biz.MultiTaskReportData, dataLocale string) template.HTML {
	if len(data.ByTask) == 0 {
		return template.HTML("")
	}
	headers := []string{
		biz.T("common.task_name"), biz.T("common.task_type"),
		biz.T("schedule_report.exec_time_col"), biz.T("schedule_report.duration_col"),
		biz.T("schedule_report.scene_pass_fail_col"), biz.T("schedule_report.data_pass_fail_col"),
		biz.T("schedule_report.scene_pass_rate"), biz.T("schedule_report.data_pass_rate"),
	}
	var rows []string
	for _, t := range data.ByTask {
		sceneRate := ""
		if t.SceneTotal > 0 {
			sceneRate = fmt.Sprintf("%.1f%%", float64(t.ScenePass)/float64(t.SceneTotal)*100)
		}
		dataRate := ""
		if t.DataTotal > 0 {
			dataRate = fmt.Sprintf("%.1f%%", float64(t.DataPass)/float64(t.DataTotal)*100)
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
		rows = append(rows, fmt.Sprintf(`<tr>
			<td>%s</td><td>%s</td>
			<td>%s</td><td>%s</td>
			<td>%d / %d / %d</td>
			<td>%d / %d / %d</td>
			<td>%s</td><td>%s</td>
		</tr>`, biz.GetTaskLocalized(t.TaskName, dataLocale), taskTypeLabel, timeStr, durStr,
			t.ScenePass, t.SceneFail, t.SceneTotal,
			t.DataPass, t.DataFail, t.DataTotal,
			sceneRate, dataRate))
	}
	return template.HTML(reportCard(biz.T("schedule_report.task_detail"), "", reportTable(headers, rows)))
}

// buildMultiTaskAPIPie 构建多任务聚合的API类型分布饼图
func buildMultiTaskAPIPie(data biz.MultiTaskReportData) template.HTML {
	if len(data.APITypeDistribution) == 0 {
		return template.HTML("")
	}
	labels := make([]string, 0, len(data.APITypeDistribution))
	counts := make([]float64, 0, len(data.APITypeDistribution))
	colors := make([]chartjs.Color, 0, len(data.APITypeDistribution))
	legendItems := make([]map[string]string, 0, len(data.APITypeDistribution))
	for i, item := range data.APITypeDistribution {
		labels = append(labels, item.Name)
		counts = append(counts, float64(item.Count))
		colors = append(colors, identityColor(i))
		legendItems = append(legendItems, map[string]string{
			"label": fmt.Sprintf(" %s - %d", item.Name, item.Count),
			"color": identityLegendColor(i),
		})
	}
	return template.HTML(reportPie("multiApiPie", biz.T("schedule_report.api_type_dist"), "", labels, counts, colors, legendItems))
}

// buildMultiTaskSceneResultPie 多任务场景执行结果分布饼图
func buildMultiTaskSceneResultPie(data biz.MultiTaskReportData) template.HTML {
	title := biz.T("schedule_report.scene_result_dist")
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
		return template.HTML(`<div class="col-md-4">` + reportCard(title, "", fmt.Sprintf(`<div style="text-align:center;padding:20px;color:#aaa">%s</div>`, biz.T("schedule_report.no_scene_data"))) + `</div>`)
	}

	labels := []string{biz.T("common.pass"), biz.T("common.fail")}
	counts := []float64{float64(pass), float64(fail)}
	colors := []chartjs.Color{colorPass, colorFail}
	legendItems := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.pass"), pass), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.fail"), fail), "color": "red"},
	}
	return template.HTML(reportPie("multiScenePie", title, "", labels, counts, colors, legendItems))
}

// buildMultiTaskDataResultPie 多任务数据执行结果分布饼图
func buildMultiTaskDataResultPie(data biz.MultiTaskReportData) template.HTML {
	title := biz.T("schedule_report.data_result_dist")
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
		return template.HTML(`<div class="col-md-4">` + reportCard(title, "", fmt.Sprintf(`<div style="text-align:center;padding:20px;color:#aaa">%s</div>`, biz.T("schedule_report.no_data_record"))) + `</div>`)
	}

	labels := []string{biz.T("common.pass"), biz.T("common.fail")}
	counts := []float64{float64(pass), float64(fail)}
	colors := []chartjs.Color{colorPass, colorFail}
	legendItems := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.pass"), pass), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("common.fail"), fail), "color": "red"},
	}
	return template.HTML(reportPie("multiDataPie", title, "", labels, counts, colors, legendItems))
}

// buildMultiTaskSceneTable 构建多任务聚合的场景/数据明细表
func buildMultiTaskSceneTable(data biz.MultiTaskReportData, dataLocale string) template.HTML {
	if len(data.SceneDetails) == 0 {
		return template.HTML("")
	}
	headers := []string{
		biz.T("common.task_name"), biz.T("schedule_report.scene_name"), biz.T("schedule_report.test_result"), biz.T("common.fail_reason"),
	}
	var rows []string
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
		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td style="color:%s">%s</td><td><div class="sc-td" style="max-width:400px"><span class="sc-truncate">%s</span><div class="sc-full">%s</div></div></td></tr>`, biz.GetTaskLocalized(s.TaskName, dataLocale), biz.GetPlaybookLocalized(s.Name, dataLocale), color, label, reason, reason))
	}
	return template.HTML(reportCard(biz.T("schedule_report.scene_detail"), "", reportTable(headers, rows)))
}

// buildMultiTaskDataTable 构建多任务聚合的数据文件明细表
func buildMultiTaskDataTable(data biz.MultiTaskReportData, dataLocale string) template.HTML {
	if len(data.DataDetails) == 0 {
		return template.HTML("")
	}
	headers := []string{
		biz.T("common.task_name"), biz.T("schedule_report.data_name"), biz.T("schedule_report.api_id_col"), biz.T("schedule_report.test_result"), biz.T("common.fail_reason"),
	}
	var rows []string
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
		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td style="color:%s">%s</td><td><div class="sc-td" style="max-width:400px"><span class="sc-truncate">%s</span><div class="sc-full">%s</div></div></td></tr>`, biz.GetTaskLocalized(d.TaskName, dataLocale), biz.GetDataLocalized(d.Name, dataLocale), d.ApiId, color, label, reason, reason))
	}
	return template.HTML(reportCard(biz.T("schedule_report.data_detail"), "", reportTable(headers, rows)))
}

// buildMultiTaskFailTable 构建多任务聚合的失败明细表
func buildMultiTaskFailTable(data biz.MultiTaskReportData) template.HTML {
	if len(data.FailItems) == 0 {
		return template.HTML("")
	}
	headers := []string{
		biz.T("schedule_report.name_with_task"), biz.T("schedule_report.type_col"), biz.T("schedule_report.api_id_col"), biz.T("common.fail_reason"),
	}
	var rows []string
	for _, f := range data.FailItems {
		typeStr := biz.T("common.scene")
		if f.Type == "data" {
			typeStr = biz.T("common.data")
		}
		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td><div class="sc-td" style="max-width:400px;color:red"><span class="sc-truncate">%s</span><div class="sc-full" style="color:#333">%s</div></div></td></tr>`, f.Name, typeStr, f.APIId, f.Reason, f.Reason))
	}
	return template.HTML(reportCard(biz.T("schedule_report.fail_detail"), "", reportTable(headers, rows)))
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
	sceneBox := `<div class="col-md-4">` + reportCard(fmt.Sprintf("%s (%d)", biz.T("schedule_report.related_scenes"), len(data.SceneList)), "", fmt.Sprintf(`<div style="max-height:300px;overflow-y:auto"><ul style="padding-left:20px">%s</ul></div>`, sceneItems)) + `</div>`

	dataItems := ""
	if len(data.DataList) > 0 {
		for _, d := range data.DataList {
			dataItems += fmt.Sprintf("<li>%s</li>", d.Name)
		}
	} else {
		dataItems = "<li style='color:#aaa'>" + biz.T("schedule_report.no_related_data") + "</li>"
	}
	dataBox := `<div class="col-md-4">` + reportCard(fmt.Sprintf("%s (%d)", biz.T("schedule_report.related_data"), len(data.DataList)), "", fmt.Sprintf(`<div style="max-height:300px;overflow-y:auto"><ul style="padding-left:20px">%s</ul></div>`, dataItems)) + `</div>`

	apiItems := ""
	if len(data.APIList) > 0 {
		for _, a := range data.APIList {
			apiItems += fmt.Sprintf("<li>%s</li>", a.Name)
		}
	} else {
		apiItems = "<li style='color:#aaa'>" + biz.T("schedule_report.no_related_api") + "</li>"
	}
	apiBox := `<div class="col-md-4">` + reportCard(fmt.Sprintf("%s (%d)", biz.T("schedule_report.related_apis"), len(data.APIList)), "", fmt.Sprintf(`<div style="max-height:300px;overflow-y:auto"><ul style="padding-left:20px">%s</ul></div>`, apiItems)) + `</div>`

	return template.HTML(reportCard(biz.T("schedule_report.resource_stats"), "", fmt.Sprintf(`<div class="row">%s%s%s</div>`, sceneBox, dataBox, apiBox)))
}
