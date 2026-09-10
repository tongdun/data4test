package tables

import (
	"data4test/biz"
	"html/template"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// GetCaseStatisticsReportTable 用例统计报告列表（独立表 case_statistics_report，不共用 dashboard）
func GetCaseStatisticsReportTable(ctx *context.Context) table.Table {
	report := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := report.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldFilterable()
	info.AddField(biz.T("schedule_report.report_name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(220)
	info.AddField(biz.T("case_statistics.related_stat"), "related_stat_id", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			id, err := strconv.Atoi(model.Value)
			if err != nil {
				return model.Value
			}
			if name := biz.GetCaseStatisticsName(id); name != "" {
				return name
			}
			return model.Value
		}).FieldWidth(110)
	info.AddField(biz.T("case_statistics.intro_versions"), "intro_versions", db.Varchar).
		FieldWidth(160)
	info.AddField(biz.T("case_statistics.product"), "product", db.Varchar).
		FieldWidth(160)
	info.AddField(biz.T("case_statistics.stat_start_time"), "stat_start_time", db.Varchar).
		FieldSortable().FieldWidth(130)
	info.AddField(biz.T("case_statistics.stat_end_time"), "stat_end_time", db.Varchar).
		FieldSortable().FieldWidth(130)
	info.AddField(biz.T("common.status"), "status", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			switch model.Value {
			case "generating":
				return "<span style=\"color:orange;font-weight:bold\">" + biz.T("common.status_generating") + "</span>"
			case "finished":
				return "<span style=\"color:green;font-weight:bold\">" + biz.T("common.status_finished") + "</span>"
			case "failed":
				return "<span style=\"color:red;font-weight:bold\">" + biz.T("common.status_failed") + "</span>"
			}
			return model.Value
		}).FieldWidth(90)
	info.AddField(biz.T("common.user_name"), "creator", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(90)
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldHide()
	info.AddField(biz.T("schedule_report.report_data"), "report_data", db.Longtext).
		FieldHide()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(130)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	// 查看报告
	info.AddActionButton(template.HTML(biz.T("common.btn_report")), action.Jump("/admin/case_statistics_report_detail?id={{.Id}}"))

	info.SetTable("case_statistics_report").SetTitle(biz.T("case_statistics.report_list_title")).SetDescription(biz.T("case_statistics.report_list_description"))

	formList := report.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("schedule_report.report_name"), "name", db.Varchar, form.Text)
	formList.AddField(biz.T("case_statistics.related_stat"), "related_stat_id", db.Int, form.Text)
	formList.AddField(biz.T("common.status"), "status", db.Varchar, form.Text)
	formList.AddField(biz.T("common.user_name"), "creator", db.Varchar, form.Text)
	formList.AddField(biz.T("schedule_report.report_data"), "report_data", db.Longtext, form.TextArea).
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("case_statistics_report").SetTitle(biz.T("case_statistics.report_list_title")).SetDescription(biz.T("case_statistics.report_list_description"))

	detail := report.GetDetail()
	detail.AddField(biz.T("schedule_report.report_name"), "name", db.Varchar)
	detail.AddField(biz.T("case_statistics.related_stat"), "related_stat_id", db.Int)
	detail.AddField(biz.T("common.status"), "status", db.Varchar)
	detail.AddField(biz.T("common.user_name"), "creator", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.remark"), "remark", db.Longtext)

	detail.SetTable("case_statistics_report").SetTitle(biz.T("case_statistics.report_list_title")).SetDescription(biz.T("case_statistics.report_list_description"))

	return report
}
