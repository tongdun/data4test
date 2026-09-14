package pages

import (
	"data4test/biz"
	"data4test/models"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"
)

// GetCaseStatisticsReportContent 根据用例统计定义 id 渲染最新的统计报告
func GetCaseStatisticsReportContent(ctx *gin.Context) (types.Panel, error) {
	id := ctx.Query("id")
	var report biz.CaseStatisticsReport
	models.Orm.Table("case_statistics_report").
		Where("status = ? AND related_stat_id = ?", "finished", id).
		Order("created_at desc").
		Limit(1).
		Find(&report)

	if report.Id == 0 || len(report.ReportData) == 0 {
		return noCaseStatisticsReport(), nil
	}

	var csData biz.CaseStatisticsReportData
	if err := json.Unmarshal([]byte(report.ReportData), &csData); err != nil {
		biz.Logger.Error("解析用例统计报告失败: %v", err)
		return types.Panel{}, err
	}
	return renderCaseStatisticsReport(ctx, csData, report)
}

// GetCaseStatisticsReportById 根据报告 id 渲染指定用例统计报告（供报告列表查看）
func GetCaseStatisticsReportById(ctx *gin.Context) (types.Panel, error) {
	id := ctx.Query("id")
	var report biz.CaseStatisticsReport
	models.Orm.Table("case_statistics_report").
		Where("id = ?", id).
		Find(&report)

	if report.Id == 0 || len(report.ReportData) == 0 {
		return noCaseStatisticsReport(), nil
	}

	var csData biz.CaseStatisticsReportData
	if err := json.Unmarshal([]byte(report.ReportData), &csData); err != nil {
		biz.Logger.Error("解析用例统计报告失败: %v", err)
		return types.Panel{}, err
	}
	return renderCaseStatisticsReport(ctx, csData, report)
}

// noCaseStatisticsReport 无报告时的占位面板
func noCaseStatisticsReport() types.Panel {
	return types.Panel{
		Content:     template.HTML(`<div style='padding:40px;text-align:center'><h3>` + biz.T("case_statistics.no_report") + `</h3></div>`),
		Title:       template.HTML(biz.T("case_statistics.report_type")),
		Description: template.HTML(biz.T("case_statistics.description")),
	}
}

// renderCaseStatisticsReport 渲染用例统计报告（供定义页最新报告 + 报告列表历史报告共用）
func renderCaseStatisticsReport(ctx *gin.Context, csData biz.CaseStatisticsReportData, report biz.CaseStatisticsReport) (types.Panel, error) {
	cookie, _ := ctx.Cookie("data_locale")
	dataLocale := biz.GetDataLocale(cookie, biz.GetLocale())

	// ===== 概览 KPI =====
	kpi := caseStatKpiCards(csData.Overview)

	// ===== 测试结果分布饼图 =====
	resultPie := caseStatResultPie(csData.ResultDistribution)

	// ===== 测试是否执行饼图 =====
	execPie := caseStatExecPie(csData.Overview)

	// ===== 功能开发者 / 用例设计者 / 用例执行者 饼图 =====
	personPies := `<div class="row">` +
		caseStatCountPie("csFunDevPie", biz.T("case_statistics.by_fun_developer"), csData.ByFunDeveloper) +
		caseStatCountPie("csDesignerPie", biz.T("case_statistics.by_case_designer"), csData.ByCaseDesigner) +
		caseStatCountPie("csExecutorPie", biz.T("case_statistics.by_case_executor"), csData.ByCaseExecutor) +
		`</div>`

	// ===== 模块维度总表 =====
	moduleTable := caseStatModuleTable(csData.ModuleStats, dataLocale)

	// ===== 用例执行者统计列表 =====
	executorTable := caseStatExecutorTable(csData.ExecutorStats)

	// ===== 扩展信息维度 =====
	var extHTML string
	if len(csData.ExtInfoStats) > 0 {
		extHTML += `<div class="row">`
		for _, ext := range csData.ExtInfoStats {
			extHTML += caseStatCountTable(biz.T("case_statistics.ext_info_dimension")+" · "+ext.Key, ext.Items)
		}
		extHTML += `</div>`
	}

	// ===== 失败 / 未执行用例下钻 =====
	drilldown := caseStatDrilldown(csData, dataLocale)

	pieRow := `<div class="row">` + execPie + resultPie +
		caseStatCountPie("csPriorityPie", biz.T("case_statistics.by_priority"), csData.ByPriority) +
		`</div>`

	content := kpi + pieRow + personPies + moduleTable + executorTable + drilldown + extHTML + string(reportStyle())

	reportName := csData.Overview.Name
	if reportName == "" {
		reportName = biz.T("case_statistics.description")
	}

	return types.Panel{
		Content:     template.HTML(content),
		Title:       template.HTML(biz.T("case_statistics.report_type")),
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s</span><span style="color:#888">%s: %s</span></div>`, template.HTMLEscapeString(reportName), biz.T("schedule_report.generated_at"), report.CreatedAt)),
	}, nil
}

// caseStatKpiCards 概览 KPI 卡片（两行布局，info-box 风格，与定时任务报告一致）
func caseStatKpiCards(ov biz.CaseStatOverview) string {
	box := func(icon, color, num, label string) string {
		return reportInfoBox(num, label, icon, color)
	}
	row1 := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		box("fa-folder", "blue", fmt.Sprintf("%d", ov.ModuleCount), biz.T("case_statistics.module_count")),
		box("fa-list", "aqua", fmt.Sprintf("%d", ov.TotalCases), biz.T("case_statistics.total_cases")),
		box("fa-minus-circle", "yellow", fmt.Sprintf("%d", ov.Untest), biz.T("test_case.result_untest")),
		box("fa-percent", "teal", fmt.Sprintf("%.1f%%", ov.ExecRate), biz.T("case_statistics.exec_rate")))
	row2 := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		box("fa-check-circle", "green", fmt.Sprintf("%d", ov.Pass), biz.T("common.pass")),
		box("fa-times-circle", "red", fmt.Sprintf("%d", ov.Fail), biz.T("common.fail")),
		box("fa-code-fork", "purple", fmt.Sprintf("%d", ov.Unmerged), biz.T("test_case.result_unmerged")),
		box("fa-percent", "green", fmt.Sprintf("%.1f%%", ov.PassRate), biz.T("case_statistics.pass_rate")))
	return row1 + row2
}

// caseStatResultPie 测试结果分布饼图
func caseStatResultPie(distribution []biz.CountItem) string {
	labels := make([]string, 0, len(distribution))
	counts := make([]float64, 0, len(distribution))
	colors := make([]chartjs.Color, 0, len(distribution))
	legendLabels := make([]map[string]string, 0, len(distribution))
	for _, item := range distribution {
		name := caseStatResultLabel(item.Name)
		labels = append(labels, name)
		counts = append(counts, float64(item.Count))
		colors = append(colors, caseStatResultColor(item.Name))
		legendLabels = append(legendLabels, map[string]string{
			"label": fmt.Sprintf(" %s - %d", name, item.Count),
			"color": caseStatResultLegendColor(item.Name),
		})
	}
	if len(labels) == 0 {
		return `<div class="col-md-4">` + reportCard(biz.T("case_statistics.result_distribution"), "", "") + `</div>`
	}
	return reportPie("csResultPie", biz.T("case_statistics.result_distribution"), "", labels, counts, colors, legendLabels)
}

// caseStatExecPie 测试是否执行饼图（已执行 vs 未执行）
func caseStatExecPie(ov biz.CaseStatOverview) string {
	executed := ov.TotalCases - ov.Untest
	labels := []string{biz.T("case_statistics.executed"), biz.T("case_statistics.not_executed")}
	counts := []float64{float64(executed), float64(ov.Untest)}
	colors := []chartjs.Color{colorPass, colorUntest}
	legendLabels := []map[string]string{
		{"label": fmt.Sprintf(" %s - %d", biz.T("case_statistics.executed"), executed), "color": "green"},
		{"label": fmt.Sprintf(" %s - %d", biz.T("case_statistics.not_executed"), ov.Untest), "color": "black"},
	}
	return reportPie("csExecPie", biz.T("case_statistics.exec_or_not"), "", labels, counts, colors, legendLabels)
}

// caseStatCountPie 通用计数饼图（用于人员/扩展等维度）
func caseStatCountPie(id, title string, items []biz.CountItem) string {
	labels := make([]string, 0, len(items))
	counts := make([]float64, 0, len(items))
	colors := make([]chartjs.Color, 0, len(items))
	legendLabels := make([]map[string]string, 0, len(items))
	for i, item := range items {
		name := caseStatCountLabel(item.Name)
		labels = append(labels, name)
		counts = append(counts, float64(item.Count))
		pieColor := identityColor(i)
		legendColor := identityLegendColor(i)
		if item.Name == biz.ExtEmptyPlaceholder {
			// 空值统一灰色配色
			pieColor = colorUntest
			legendColor = "black"
		}
		colors = append(colors, pieColor)
		legendLabels = append(legendLabels, map[string]string{
			"label": fmt.Sprintf(" %s - %d", name, item.Count),
			"color": legendColor,
		})
	}
	if len(labels) == 0 {
		return `<div class="col-md-4">` + reportCard(title, "", "") + `</div>`
	}
	return reportPie(id, title, "", labels, counts, colors, legendLabels)
}

// caseStatCountLabel 计数饼图切片名：空值占位符翻译为「空值」
func caseStatCountLabel(name string) string {
	if name == biz.ExtEmptyPlaceholder {
		return biz.T("case_statistics.empty_value")
	}
	return name
}

// caseStatModuleTable 按模块统计用例执行表
func caseStatModuleTable(items []biz.CaseStatModuleItem, dataLocale string) string {
	headers := []string{
		biz.T("case_statistics.level1"),
		biz.T("case_statistics.level2"),
		biz.T("common.case_module"),
		biz.T("case_statistics.total"),
		biz.T("case_statistics.executed"),
		biz.T("common.pass"),
		biz.T("common.fail"),
		biz.T("common.status_discarded"),
		biz.T("test_case.result_unmerged"),
		biz.T("case_statistics.not_executed"),
		biz.T("case_statistics.pass_rate"),
		biz.T("case_statistics.exec_rate"),
	}

	var rows []string
	for _, item := range items {
		level1 := item.Level1
		if level1 == "" {
			level1 = biz.T("case_statistics.ungrouped")
		} else {
			level1 = biz.GetCaseCountLocalized(level1, dataLocale)
		}
		level2 := biz.GetCaseCountLocalized(item.Level2, dataLocale)
		module := biz.GetCaseCountLocalized(item.Module, dataLocale)

		failCell := fmt.Sprintf("%d", item.Fail)
		if item.Fail > 0 {
			failCell = fmt.Sprintf(`<a href="/admin/info/test_case?test_result=fail&module=%s" target="_blank" style="color:red">%d</a>`, template.URLQueryEscaper(item.Module), item.Fail)
		}
		untestCell := fmt.Sprintf("%d", item.Untest)
		if item.Untest > 0 {
			untestCell = fmt.Sprintf(`<a href="/admin/info/test_case?untest_broad=1&module=%s" target="_blank">%d</a>`, template.URLQueryEscaper(item.Module), item.Untest)
		}

		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%s</td><td>%d</td><td>%d</td><td>%s</td><td>%.1f%%</td><td>%.1f%%</td></tr>`,
			template.HTMLEscapeString(level1),
			template.HTMLEscapeString(level2),
			template.HTMLEscapeString(module),
			item.Total, item.Total-item.Untest, item.Pass, failCell, item.Deprecated, item.Unmerged, untestCell, item.PassRate, item.ExecRate))
	}
	return reportCard(biz.T("case_statistics.module_dimension"), "", reportTable(headers, rows))
}

// caseStatExecutorTable 按人员统计用例执行列表（总用例数/已执行/通过/失败/废弃/暂未合入/未执行/执行率）
func caseStatExecutorTable(items []biz.CaseStatExecutorItem) string {
	headers := []string{
		biz.T("case_statistics.by_case_executor"),
		biz.T("case_statistics.executor_total"),
		biz.T("case_statistics.executor_executed"),
		biz.T("case_statistics.executor_pass"),
		biz.T("case_statistics.executor_fail"),
		biz.T("case_statistics.executor_deprecated"),
		biz.T("case_statistics.executor_unmerged"),
		biz.T("case_statistics.executor_untest"),
		biz.T("case_statistics.exec_rate"),
	}
	var rows []string
	for _, it := range items {
		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%.1f%%</td></tr>`,
			template.HTMLEscapeString(caseStatCountLabel(it.Name)), it.Total, it.Executed, it.Pass, it.Fail, it.Deprecated, it.Unmerged, it.Untest, it.ExecRate))
	}
	return reportCard(biz.T("case_statistics.executor_stat_title"), "", reportTable(headers, rows))
}

// caseStatCountTable 通用计数表（功能开发者/用例设计者/用例执行者/扩展字段）
func caseStatCountTable(title string, items []biz.CountItem) string {
	headers := []string{biz.T("case_statistics.dim_name"), biz.T("case_statistics.dim_count")}
	var rows []string
	for _, it := range items {
		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%d</td></tr>`, template.HTMLEscapeString(it.Name), it.Count))
	}
	return `<div class="col-md-4">` + reportCard(title, "", reportTable(headers, rows)) + `</div>`
}

// caseStatDrilldown 失败用例下钻（含备注）
func caseStatDrilldown(csData biz.CaseStatisticsReportData, dataLocale string) string {
	return caseStatDetailBox(biz.T("case_statistics.fail_cases"), "/admin/info/test_case?test_result=fail", dataLocale, csData.FailCases)
}

// caseStatDetailBox 单个用例明细列表（失败/未执行），含备注列
func caseStatDetailBox(title, filterURL, dataLocale string, cases []biz.CaseStatCaseDetail) string {
	link := fmt.Sprintf(`<a href="%s" target="_blank">%s (%d) &raquo;</a>`, filterURL, title, len(cases))

	headers := []string{
		biz.T("common.case_number"),
		biz.T("common.case_title"),
		biz.T("common.case_module"),
		biz.T("common.remark"),
	}
	var rows []string
	for _, c := range cases {
		caseName := biz.GetCaseLocalized(c.CaseNumber, c.Module, dataLocale, biz.FieldCaseName)
		if caseName == "" {
			caseName = c.CaseName
		}
		module := biz.GetCaseLocalized(c.CaseNumber, c.Module, dataLocale, biz.FieldCaseModule)
		if module == "" {
			module = c.Module
		}
		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,
			template.HTMLEscapeString(c.CaseNumber),
			template.HTMLEscapeString(caseName),
			template.HTMLEscapeString(module),
			template.HTMLEscapeString(c.Remark)))
	}
	return reportCard(title, link, reportTable(headers, rows))
}

func caseStatResultLabel(name string) string {
	switch name {
	case "pass":
		return biz.T("common.pass")
	case "fail":
		return biz.T("common.fail")
	case "untest":
		return biz.T("test_case.result_untest")
	case "part":
		return biz.T("test_case.result_part")
	case "deprecated":
		return biz.T("common.status_discarded")
	case "unmerged":
		return biz.T("test_case.result_unmerged")
	}
	return name
}

func caseStatResultColor(name string) chartjs.Color {
	switch name {
	case "pass":
		return colorPass
	case "fail":
		return colorFail
	case "untest":
		return colorUntest
	case "part":
		return colorPart
	case "deprecated":
		return colorDeprecated
	case "unmerged":
		return colorUnmerged
	}
	return "rgba(120,120,120,1)"
}

func caseStatResultLegendColor(name string) string {
	switch name {
	case "pass":
		return "green"
	case "fail":
		return "red"
	case "part":
		return "yellow"
	case "unmerged":
		return "blue"
	default:
		return "black"
	}
}
