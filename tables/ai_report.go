package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetAiReportTable(ctx *context.Context) table.Table {

	aiReport := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiReport.GetInfo()

	info.AddField(biz.T("common.id"), "id", db.Int)
	info.AddField(biz.T("ai_report.report_name"), "report_name", db.Varchar)
	info.AddField(biz.T("ai_report.demand"), "demand", db.Text)
	info.AddField(biz.T("common.source"), "source", db.Varchar)
	info.AddField(biz.T("common.use_status"), "use_status", db.Enum)
	info.AddField(biz.T("common.modify_status"), "modify_status", db.Enum)
	info.AddField(biz.T("ai_report.report_link"), "report_link", db.Varchar)
	info.AddField(biz.T("common.user_name"), "create_user", db.Varchar)
	info.AddField(biz.T("common.modify_user"), "modify_user", db.Varchar)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp)

	info.SetTable("ai_report").SetTitle(biz.T("ai_report.title")).SetDescription(biz.T("ai_report.description"))

	formList := aiReport.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("ai_report.report_name"), "report_name", db.Varchar, form.Text)
	formList.AddField(biz.T("ai_report.demand"), "demand", db.Text, form.RichText)
	formList.AddField(biz.T("common.source"), "source", db.Varchar, form.Text)
	formList.AddField(biz.T("common.use_status"), "use_status", db.Enum, form.Text)
	formList.AddField(biz.T("common.modify_status"), "modify_status", db.Enum, form.Text)
	formList.AddField(biz.T("ai_report.report_link"), "report_link", db.Varchar, form.Text)
	formList.AddField(biz.T("common.user_name"), "create_user", db.Varchar, form.Text)
	formList.AddField(biz.T("common.modify_user"), "modify_user", db.Varchar, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime)

	formList.SetTable("ai_report").SetTitle(biz.T("ai_report.title")).SetDescription(biz.T("ai_report.description"))

	return aiReport
}
