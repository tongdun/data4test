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

// renderOrchestrationReport 渲染编排报告
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

	var b strings.Builder
	b.WriteString(`<div style="padding:10px 20px 30px">`)

	// 概览
	b.WriteString(`<h4 style="margin:10px 0">` + template.HTMLEscapeString(data.OrchestrationName) + `</h4>`)
	b.WriteString(`<div style="display:flex;flex-wrap:wrap;gap:16px;margin-bottom:20px">`)
	// 历史执行切换（最近 20 次）
	if reports := biz.GetOrchestrationRecentReports(data.OrchestrationId, 20); len(reports) > 1 {
		b.WriteString(`<div style="margin:4px 0 12px">` + biz.T("task_orchestration.history") + `: <select onchange="location.href='/admin/report_dashboard?id='+this.value">`)
		for _, r := range reports {
			selected := ""
			if r.Id == report.Id {
				selected = " selected"
			}
			b.WriteString(fmt.Sprintf(`<option value="%s"%s>%s (%s)</option>`, r.Id, selected, template.HTMLEscapeString(r.ReportName), r.TimeRangeStart))
		}
		b.WriteString(`</select></div>`)
	}

	overviewItems := []struct {
		label string
		value string
	}{
		{biz.T("task_orchestration.trigger_type"), orchestrationTriggerText(data.TriggerType)},
		{biz.T("task_orchestration.start_time"), report.TimeRangeStart},
		{biz.T("task_orchestration.end_time"), report.TimeRangeEnd},
		{biz.T("task_orchestration.step_count"), fmt.Sprintf("%d", data.Summary.Total)},
		{biz.T("task_orchestration.summary_success"), fmt.Sprintf("%d", data.Summary.Success)},
		{biz.T("task_orchestration.summary_fail"), fmt.Sprintf("%d", data.Summary.Fail)},
		{biz.T("task_orchestration.summary_skipped"), fmt.Sprintf("%d", data.Summary.Skipped)},
		{biz.T("task_orchestration.summary_duration"), fmt.Sprintf("%ds", data.Summary.DurationSeconds)},
		{biz.T("task_orchestration.summary_pass_rate"), fmt.Sprintf("%.2f%%", data.Summary.PassRate)},
	}
	for _, item := range overviewItems {
		b.WriteString(`<div style="min-width:120px;padding:10px 14px;border:1px solid #e8e8e8;border-radius:4px">`)
		b.WriteString(`<div style="color:#888;font-size:12px">` + template.HTMLEscapeString(item.label) + `</div>`)
		b.WriteString(`<div style="font-size:18px;margin-top:4px">` + template.HTMLEscapeString(item.value) + `</div>`)
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div>`)

	// 步骤明细
	b.WriteString(`<h4 style="margin:20px 0 10px">` + biz.T("task_orchestration.report_steps") + `</h4>`)
	b.WriteString(`<table class="table table-bordered" style="background:#fff"><thead><tr>`)
	for _, h := range []string{
		biz.T("task_orchestration.step_order"),
		biz.T("task_orchestration.step_task"),
		biz.T("task_orchestration.step_status"),
		biz.T("task_orchestration.start_time"),
		biz.T("task_orchestration.end_time"),
		biz.T("task_orchestration.step_duration"),
		biz.T("task_orchestration.step_cases"),
		biz.T("task_orchestration.step_fail"),
		biz.T("task_orchestration.step_error"),
		biz.T("task_orchestration.step_report"),
	} {
		b.WriteString(`<th>` + template.HTMLEscapeString(h) + `</th>`)
	}
	b.WriteString(`</tr></thead><tbody>`)

	for _, step := range data.Steps {
		b.WriteString(`<tr>`)
		b.WriteString(fmt.Sprintf(`<td>%d</td>`, step.StepOrder))
		b.WriteString(`<td>` + template.HTMLEscapeString(step.TaskName) + `</td>`)
		b.WriteString(`<td>` + template.HTMLEscapeString(orchestrationStepStatusText(step.Status)) + `</td>`)
		b.WriteString(`<td>` + template.HTMLEscapeString(step.StartAt) + `</td>`)
		b.WriteString(`<td>` + template.HTMLEscapeString(step.EndAt) + `</td>`)
		b.WriteString(fmt.Sprintf(`<td>%ds</td>`, step.DurationSeconds))
		b.WriteString(fmt.Sprintf(`<td>%d/%d</td>`, step.Success, step.Total))
		b.WriteString(fmt.Sprintf(`<td>%d</td>`, step.Fail))
		b.WriteString(`<td style="max-width:360px;word-break:break-all">` + template.HTMLEscapeString(step.ErrorMsg) + `</td>`)
		if len(step.DashboardId) > 0 {
			b.WriteString(`<td><a href="/admin/report_dashboard?id=` + template.HTMLEscapeString(step.DashboardId) + `" target="_blank">` +
				template.HTMLEscapeString(biz.T("task_orchestration.step_report")) + `</a></td>`)
		} else {
			b.WriteString(`<td>-</td>`)
		}
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)

	// 产品维度汇总
	if len(data.Products) > 0 {
		b.WriteString(`<h4 style="margin:20px 0 10px">` + biz.T("task_orchestration.report_products") + `</h4>`)
		b.WriteString(`<table class="table table-bordered" style="background:#fff"><thead><tr>`)
		for _, h := range []string{
			biz.T("task_orchestration.product"),
			biz.T("task_orchestration.summary_total"),
			biz.T("task_orchestration.summary_success"),
			biz.T("task_orchestration.summary_fail"),
			biz.T("task_orchestration.summary_pass_rate"),
		} {
			b.WriteString(`<th>` + template.HTMLEscapeString(h) + `</th>`)
		}
		b.WriteString(`</tr></thead><tbody>`)
		for _, p := range data.Products {
			b.WriteString(`<tr>`)
			b.WriteString(`<td>` + template.HTMLEscapeString(p.Product) + `</td>`)
			b.WriteString(fmt.Sprintf(`<td>%d</td><td>%d</td><td>%d</td><td>%.2f%%</td>`, p.Total, p.Pass, p.Fail, p.PassRate))
			b.WriteString(`</tr>`)
		}
		b.WriteString(`</tbody></table>`)
	}

	b.WriteString(`</div>`)

	return types.Panel{
		Content:     template.HTML(b.String()),
		Title:       template.HTML(biz.T("task_orchestration.report_title")),
		Description: template.HTML(biz.T("task_orchestration.description")),
	}, nil
}

func orchestrationStepStatusText(status string) string {
	switch status {
	case biz.OrchestrationStepSuccess:
		return biz.T("task_orchestration.step_success")
	case biz.OrchestrationStepFail:
		return biz.T("task_orchestration.step_fail")
	case biz.OrchestrationStepSkipped:
		return biz.T("task_orchestration.step_skipped")
	default:
		return status
	}
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
