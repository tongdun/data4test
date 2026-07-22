package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"html/template"
)

func GetDashboardTable(ctx *context.Context) table.Table {

	report := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	apps := biz.GetApps()
	info := report.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int)

	info.AddField(biz.T("schedule_report.report_name"), "report_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(200)

	info.AddField(biz.T("schedule_report.report_type"), "report_type", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			switch model.Value {
			case "global":
				return biz.T("schedule_report.report_type_global")
			case "product":
				return biz.T("schedule_report.report_type_product")
			case "app":
				return biz.T("schedule_report.report_type_app")
			case "task":
				return biz.T("schedule_report.report_type_task")
			}
			return model.Value
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "global", Text: biz.T("schedule_report.report_type_global")},
		{Value: "product", Text: biz.T("schedule_report.report_type_product")},
		{Value: "app", Text: biz.T("schedule_report.report_type_app")},
		{Value: "task", Text: biz.T("schedule_report.report_type_task")},
	}).FieldWidth(100)

	info.AddField(biz.T("common.product_list"), "related_products", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(150)
	info.AddField(biz.T("product.apps"), "related_apps", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(150)
	info.AddField(biz.T("common.related_task"), "related_task_ids", db.Varchar)

	info.AddField(biz.T("common.user_name"), "creator", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)

	info.AddField(biz.T("schedule_report.time_range_start"), "time_range_start", db.Datetime).
		FieldWidth(110)

	info.AddField(biz.T("schedule_report.time_range_end"), "time_range_end", db.Datetime).
		FieldWidth(110)

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
		}).FieldWidth(80)

	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldHide()

	info.AddField(biz.T("schedule_report.report_data"), "report_data", db.Longtext).
		FieldHide()

	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110)

	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()

	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	// 查看报告按钮
	info.AddActionButton(template.HTML(biz.T("common.btn_report")), action.Jump("/admin/report_dashboard?id={{.Id}}"))

	info.SetTable("dashboard").SetTitle(biz.T("schedule_report.title")).SetDescription(biz.T("schedule_report.description"))

	formList := report.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("schedule_report.report_name"), "report_name", db.Varchar, form.Text)
	formList.AddField(biz.T("schedule_report.report_type"), "report_type", db.Varchar, form.Text).FieldHide()
	formList.AddField(biz.T("common.product_list"), "related_products", db.Text, form.Text)
	formList.AddField(biz.T("product.apps"), "related_apps", db.Text, form.Select).
		FieldOptions(apps)
	formList.AddField(biz.T("common.related_task"), "related_task_ids", db.Text, form.Text)
	formList.AddField(biz.T("common.user_name"), "creator", db.Varchar, form.Text)
	formList.AddField(biz.T("schedule_report.time_range_start"), "time_range_start", db.Datetime, form.Datetime)
	formList.AddField(biz.T("schedule_report.time_range_end"), "time_range_end", db.Datetime, form.Datetime)
	formList.AddField(biz.T("common.status"), "status", db.Varchar, form.Text)
	formList.AddField(biz.T("schedule_report.report_data"), "report_data", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("dashboard").SetTitle(biz.T("schedule_report.title")).SetDescription(biz.T("schedule_report.description"))

	detail := report.GetDetail()
	detail.AddField(biz.T("schedule_report.report_name"), "report_name", db.Varchar)
	detail.AddField(biz.T("schedule_report.report_type"), "report_type", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			switch model.Value {
			case "global":
				return biz.T("schedule_report.report_type_global")
			case "product":
				return biz.T("schedule_report.report_type_product")
			case "app":
				return biz.T("schedule_report.report_type_app")
			case "task":
				return biz.T("schedule_report.report_type_task")
			}
			return model.Value
		})
	detail.AddField(biz.T("common.product_list"), "related_products", db.Varchar)
	detail.AddField(biz.T("product.apps"), "related_apps", db.Varchar)
	detail.AddField(biz.T("common.related_task"), "related_task_ids", db.Varchar)
	detail.AddField(biz.T("schedule_report.time_range_start"), "time_range_start", db.Datetime)
	detail.AddField(biz.T("schedule_report.time_range_end"), "time_range_end", db.Datetime)
	detail.AddField(biz.T("common.status"), "status", db.Varchar)
	detail.AddField(biz.T("schedule_report.report_data"), "report_data", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return fmt.Sprintf("<a href=\"/admin/schedule_report?report_id=%s\">%s</a>", model.Row["id"], biz.T("schedule_report.btn_view_report"))
		})
	detail.AddField(biz.T("common.remark"), "remark", db.Longtext)
	detail.AddField(biz.T("common.user_name"), "creator", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)

	detail.SetTable("dashboard").SetTitle(biz.T("schedule_report.title")).SetDescription(biz.T("schedule_report.detail_title"))

	return report
}
