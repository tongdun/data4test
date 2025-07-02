package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetAiReportTable(ctx *context.Context) table.Table {

	aiReport := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiReport.GetInfo()

	info.AddField("自增主键", "id", db.Int)
	info.AddField("报告名称", "report_name", db.Varchar)
	info.AddField("报告需求", "demand", db.Text)
	info.AddField("生成来源", "source", db.Varchar)
	info.AddField("取用状态", "use_status", db.Enum)
	info.AddField("改造状态", "modify_status", db.Enum)
	info.AddField("报告详情", "report_link", db.Varchar)
	info.AddField("创建人", "create_user", db.Varchar)
	info.AddField("修改人", "modify_user", db.Varchar)
	info.AddField("创建时间", "created_at", db.Timestamp)
	info.AddField("更新时间", "updated_at", db.Timestamp)
	info.AddField("删除时间", "deleted_at", db.Timestamp)

	info.SetTable("ai_report").SetTitle("智能报告").SetDescription("AiReport")

	formList := aiReport.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("报告名称", "report_name", db.Varchar, form.Text)
	formList.AddField("报告需求", "demand", db.Text, form.RichText)
	formList.AddField("生成来源", "source", db.Varchar, form.Text)
	formList.AddField("取用状态", "use_status", db.Enum, form.Text)
	formList.AddField("改造状态", "modify_status", db.Enum, form.Text)
	formList.AddField("报告详情", "report_link", db.Varchar, form.Text)
	formList.AddField("创建人", "create_user", db.Varchar, form.Text)
	formList.AddField("修改人", "modify_user", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime)

	formList.SetTable("ai_report").SetTitle("智能报告").SetDescription("AiReport")

	return aiReport
}
