package pages

import (
	"data4test/biz"
	"data4test/models"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"
)

// GetOrchestrationReportContent 按编排 id 展示最近一次编排报告
func GetOrchestrationReportContent(ctx *gin.Context) (types.Panel, error) {
	orchestrationId := ctx.Query("id")
	reportId, err := biz.GetOrchestrationLatestReport(orchestrationId)
	if err != nil {
		return types.Panel{
			Content: template.HTML(`<div style='padding:40px;text-align:center'>
					<h3>` + biz.T("schedule_report.no_report") + `</h3>
				</div>`),
			Title:       template.HTML(biz.T("task_orchestration.report_title")),
			Description: template.HTML(biz.T("task_orchestration.description")),
		}, nil
	}

	var report biz.DashboardReport
	models.Orm.Table("dashboard").Where("id = ?", reportId).Find(&report)
	return renderOrchestrationReport(report)
}

// renderOrchestrationReport 渲染编排报告，视觉风格与定时任务报告保持一致：
// 顶部信息卡 + KPI 卡片 + 饼图卡片 + 统一表格（状态徽标 / hover 展开失败原因）。
func renderOrchestrationReport(report biz.DashboardReport) (types.Panel, error) {
	if len(report.ReportData) == 0 {
		message := biz.T("task_orchestration.report_generating")
		if report.Status == "failed" {
			message = biz.T("task_orchestration.report_interrupted")
		}
		return types.Panel{
			Content: template.HTML(`<div style='padding:40px;text-align:center'>
					<h3>` + message + `</h3>
				</div>`),
			Title:       template.HTML(biz.T("task_orchestration.report_title")),
			Description: template.HTML(biz.T("task_orchestration.description")),
		}, nil
	}

	var data biz.OrchestrationReportData
	if err := json.Unmarshal([]byte(report.ReportData), &data); err != nil {
		biz.Logger.Error("%s", err)
		return types.Panel{}, err
	}
	// 历史报告未保存场景/数据文件及产品维度统计，展示时按步骤任务报告回溯
	biz.EnsureOrchestrationReportStats(&data)

	var b strings.Builder

	// 历史执行切换（最近 20 次），沿用编排报告原有能力，样式收敛到报告统一风格
	if reports := biz.GetOrchestrationRecentReports(data.OrchestrationId, 20); len(reports) > 1 {
		b.WriteString(`<div style="margin:0 0 12px;font-size:13px;color:#555">` + biz.T("task_orchestration.history") + `: <select class="form-control" style="display:inline-block;width:auto;height:30px;padding:0 8px" onchange="location.href='/admin/report_dashboard?id='+this.value">`)
		for _, r := range reports {
			selected := ""
			if r.Id == report.Id {
				selected = " selected"
			}
			b.WriteString(fmt.Sprintf(`<option value="%s"%s>%s (%s)</option>`, r.Id, selected, template.HTMLEscapeString(r.ReportName), template.HTMLEscapeString(r.TimeRangeStart)))
		}
		b.WriteString(`</select></div>`)
	}

	b.WriteString(string(buildOrchestrationHeader(data, report)))
	b.WriteString(string(buildOrchestrationKpiCards(data)))
	b.WriteString(string(buildOrchestrationResourceKpiCards(data)))
	b.WriteString(string(buildOrchestrationStepsTable(data)))
	if len(data.Products) > 0 {
		b.WriteString(string(buildOrchestrationProductsTable(data)))
	}
	b.WriteString(string(reportStyle()))

	return types.Panel{
		Content: template.HTML(b.String()),
		Title:   template.HTML(biz.T("task_orchestration.report_title")),
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s</span><span style="color:#888">%s: %s</span></div>`,
			template.HTMLEscapeString(data.OrchestrationName), biz.T("schedule_report.generated_at"), template.HTMLEscapeString(report.CreatedAt))),
	}, nil
}

// buildOrchestrationHeader 顶部执行信息卡（白色背景，两行排版，对齐任务报告表头）。
func buildOrchestrationHeader(data biz.OrchestrationReportData, report biz.DashboardReport) template.HTML {
	products := strings.TrimSpace(data.ProductList)
	if len(products) == 0 {
		products = "—"
	}
	executor := report.Creator
	if len(executor) == 0 {
		executor = "—"
	}
	s := data.Summary
	html := fmt.Sprintf(`<div class="row" style="margin-bottom:8px">
			<div class="col-md-4"><strong>%s</strong> %s</div>
			<div class="col-md-2"><strong>%s</strong> %s</div>
			<div class="col-md-4"><strong>%s</strong> %s</div>
			<div class="col-md-2"><strong>%s</strong> %s</div>
		</div>
		<div class="row">
			<div class="col-md-4"><strong>%s</strong> %s ~ %s</div>
			<div class="col-md-2"><strong>%s</strong> %s</div>
			<div class="col-md-4"><strong>%s</strong> %d | <strong>%s</strong> %d | <strong>%s</strong> %d | <strong>%s</strong> %d</div>
			<div class="col-md-2"><strong>%s</strong> <span style="color:%s;font-weight:bold">%.1f%%</span></div>
		</div>`,
		biz.T("task_orchestration.orchestration_name"), template.HTMLEscapeString(data.OrchestrationName),
		biz.T("task_orchestration.trigger_type"), template.HTMLEscapeString(orchestrationTriggerText(data.TriggerType)),
		biz.T("task_orchestration.product"), template.HTMLEscapeString(products),
		biz.T("schedule_report.executor"), template.HTMLEscapeString(executor),
		biz.T("task_orchestration.start_time"), template.HTMLEscapeString(report.TimeRangeStart), template.HTMLEscapeString(report.TimeRangeEnd),
		biz.T("task_orchestration.summary_duration"), formatDuration(s.DurationSeconds),
		biz.T("task_orchestration.report_task_total"), s.Total,
		biz.T("task_orchestration.summary_success"), s.Success,
		biz.T("task_orchestration.summary_fail"), s.Fail,
		biz.T("task_orchestration.summary_skipped"), s.Skipped,
		biz.T("task_orchestration.summary_pass_rate"), orchestrationPassRateColor(s), s.PassRate)

	return template.HTML(reportCard("", "", html))
}

// buildOrchestrationKpiCards 任务级 KPI 卡片：任务总数 / 成功 / 失败 / 通过率（对齐任务报告的 info-box 风格）。
func buildOrchestrationKpiCards(data biz.OrchestrationReportData) template.HTML {
	s := data.Summary
	return template.HTML(fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		reportInfoBox(fmt.Sprintf("%d", s.Total), biz.T("task_orchestration.report_task_total"), "fa-cubes", "blue"),
		reportInfoBox(fmt.Sprintf("%d", s.Success), biz.T("task_orchestration.summary_success"), "fa-check-circle", "green"),
		reportInfoBox(fmt.Sprintf("%d", s.Fail), biz.T("task_orchestration.summary_fail"), "fa-times-circle", "red"),
		reportInfoBox(fmt.Sprintf("%.1f%%", s.PassRate), biz.T("schedule_report.pass_rate_label"), "fa-percent", "green")))
}

// buildOrchestrationResourceKpiCards 场景 / 数据文件的执行统计（各一行 KPI 卡片）。
func buildOrchestrationResourceKpiCards(data biz.OrchestrationReportData) template.HTML {
	s := data.SceneStats
	d := data.DataStats
	sceneRow := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		reportInfoBox(fmt.Sprintf("%d", s.Total), biz.T("task_orchestration.scene_total"), "fa-cubes", "blue"),
		reportInfoBox(fmt.Sprintf("%d", s.Success), biz.T("task_orchestration.success_count"), "fa-check-circle", "green"),
		reportInfoBox(fmt.Sprintf("%d", s.Fail), biz.T("task_orchestration.fail_count"), "fa-times-circle", "red"),
		reportInfoBox(fmt.Sprintf("%.1f%%", s.PassRate), biz.T("schedule_report.pass_rate_label"), "fa-percent", "green"))
	dataRow := fmt.Sprintf(`<div class="row">%s%s%s%s</div>`,
		reportInfoBox(fmt.Sprintf("%d", d.Total), biz.T("task_orchestration.data_total"), "fa-file-text", "blue"),
		reportInfoBox(fmt.Sprintf("%d", d.Success), biz.T("task_orchestration.success_count"), "fa-check-circle", "green"),
		reportInfoBox(fmt.Sprintf("%d", d.Fail), biz.T("task_orchestration.fail_count"), "fa-times-circle", "red"),
		reportInfoBox(fmt.Sprintf("%.1f%%", d.PassRate), biz.T("schedule_report.pass_rate_label"), "fa-percent", "green"))
	return template.HTML(sceneRow + dataRow)
}

// buildOrchestrationStepsTable 步骤明细表（统一表格 + 状态徽标 + 失败步骤下方整行展示错误原因）。
func buildOrchestrationStepsTable(data biz.OrchestrationReportData) template.HTML {
	headers := []string{
		biz.T("task_orchestration.step_order"),
		biz.T("task_orchestration.step_task"),
		biz.T("task_orchestration.step_status"),
		biz.T("task_orchestration.step_time"),
		biz.T("task_orchestration.step_duration"),
		biz.T("task_orchestration.step_cases"),
		biz.T("task_orchestration.step_fail"),
		biz.T("task_orchestration.step_report"),
	}
	var rows []string
	for _, step := range data.Steps {
		reportCell := `<td>—</td>`
		if len(step.DashboardId) > 0 {
			reportCell = fmt.Sprintf(`<td><a href="/admin/report_dashboard?id=%s" target="_blank">%s</a></td>`,
				template.HTMLEscapeString(step.DashboardId), template.HTMLEscapeString(biz.T("task_orchestration.step_report")))
		}
		rows = append(rows, fmt.Sprintf(`<tr><td>%d</td><td>%s</td><td>%s</td><td><div style="white-space:nowrap">%s</div><div style="white-space:nowrap;color:#888">%s</div></td><td>%s</td><td>%d/%d</td><td>%d</td>%s</tr>`,
			step.StepOrder,
			template.HTMLEscapeString(step.TaskName),
			orchestrationStatusBadge(step.Status),
			template.HTMLEscapeString(step.StartAt),
			template.HTMLEscapeString(step.EndAt),
			formatDuration(step.DurationSeconds),
			step.Success, step.Total,
			step.Fail,
			reportCell))
		// 失败步骤在下一行整行展示错误原因，避免长文本挤压表格列宽
		if msg := strings.TrimSpace(step.ErrorMsg); len(msg) > 0 {
			rows = append(rows, fmt.Sprintf(`<tr><td></td><td colspan="7"><div style="color:#c0392b;font-size:12px;white-space:pre-wrap;word-break:break-all;max-height:160px;overflow:auto;background:#fdf3f2;border-radius:4px;padding:8px">%s</div></td></tr>`,
				template.HTMLEscapeString(msg)))
		}
	}
	return template.HTML(reportCard(biz.T("task_orchestration.report_steps"), "", reportTable(headers, rows)))
}

// buildOrchestrationProductsTable 产品维度汇总表（统一表格）。
func buildOrchestrationProductsTable(data biz.OrchestrationReportData) template.HTML {
	headers := []string{
		biz.T("task_orchestration.product"),
		biz.T("task_orchestration.report_task_total"),
		biz.T("task_orchestration.task_success_count"),
		biz.T("task_orchestration.task_fail_count"),
		biz.T("task_orchestration.scene_total"),
		biz.T("task_orchestration.scene_success_count"),
		biz.T("task_orchestration.scene_fail_count"),
		biz.T("task_orchestration.data_total"),
		biz.T("task_orchestration.data_success_count"),
		biz.T("task_orchestration.data_fail_count"),
		biz.T("task_orchestration.data_pass_rate"),
	}
	var rows []string
	for _, p := range data.Products {
		rows = append(rows, fmt.Sprintf(`<tr><td>%s</td><td>%d</td><td><span style="color:#00a65a">%d</span></td><td><span style="color:#dd4b39">%d</span></td>`+
			`<td>%d</td><td><span style="color:#00a65a">%d</span></td><td><span style="color:#dd4b39">%d</span></td>`+
			`<td>%d</td><td><span style="color:#00a65a">%d</span></td><td><span style="color:#dd4b39">%d</span></td><td>%.1f%%</td></tr>`,
			template.HTMLEscapeString(p.Product),
			p.TaskTotal, p.TaskSuccess, p.TaskFail,
			p.SceneTotal, p.SceneSuccess, p.SceneFail,
			p.DataTotal, p.DataSuccess, p.DataFail, p.DataPassRate))
	}
	return template.HTML(reportCard(biz.T("task_orchestration.report_products"), "", reportTable(headers, rows)))
}

// orchestrationStatusBadge 步骤状态徽标（成功绿 / 失败红 / 跳过灰）。
func orchestrationStatusBadge(status string) string {
	label := biz.T("task_orchestration.step_success")
	cls := "sc-badge sc-badge-pass"
	switch status {
	case biz.OrchestrationStepFail:
		label = biz.T("common.fail")
		cls = "sc-badge sc-badge-fail"
	case biz.OrchestrationStepSkipped:
		label = biz.T("task_orchestration.summary_skipped")
		cls = "sc-badge sc-badge-na"
	}
	return fmt.Sprintf(`<span class="%s">%s</span>`, cls, template.HTMLEscapeString(label))
}

// orchestrationPassRateColor 通过率颜色：有失败=红，有跳过=橙，全通过=绿。
func orchestrationPassRateColor(s biz.OrchestrationSummary) string {
	if s.Fail > 0 {
		return "red"
	}
	if s.Skipped > 0 {
		return "orange"
	}
	return "green"
}

func orchestrationTriggerText(trigger string) string {
	switch trigger {
	case "manual":
		return biz.T("task_orchestration.trigger_manual")
	case "cron":
		return biz.T("task_orchestration.trigger_cron")
	default:
		return trigger
	}
}
