package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetAiTaskTable(ctx *context.Context) table.Table {

	aiTask := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiTask.GetInfo()

	info.AddField("自增主键", "id", db.Int)
	info.AddField("任务名称", "task_name", db.Varchar)
	info.AddField("任务模式", "task_mode", db.Enum)
	info.AddField("任务类型", "task_type", db.Enum)
	info.AddField("生成来源", "source", db.Varchar)
	info.AddField("任务状态", "task_status", db.Enum)
	info.AddField("关联数据", "data_list", db.Text)
	info.AddField("关联场景", "playbook_list", db.Text)
	info.AddField("取用状态", "use_status", db.Enum)
	info.AddField("改造状态", "modify_status", db.Enum)
	info.AddField("所属产品", "product", db.Varchar)
	info.AddField("创建人", "create_user", db.Varchar)
	info.AddField("创建时间", "created_at", db.Timestamp)
	info.AddField("更新时间", "updated_at", db.Timestamp)
	info.AddField("删除时间", "deleted_at", db.Timestamp)

	info.SetTable("ai_task").SetTitle("智能任务").SetDescription("AiTask")

	formList := aiTask.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("任务名称", "task_name", db.Varchar, form.Text)
	formList.AddField("任务模式", "task_mode", db.Enum, form.Text)
	formList.AddField("任务类型", "task_type", db.Enum, form.Text)
	formList.AddField("生成来源", "source", db.Varchar, form.Text)
	formList.AddField("任务状态", "task_status", db.Enum, form.Text)
	formList.AddField("关联数据", "data_list", db.Text, form.RichText)
	formList.AddField("关联场景", "playbook_list", db.Text, form.RichText)
	formList.AddField("取用状态", "use_status", db.Enum, form.Text)
	formList.AddField("改造状态", "modify_status", db.Enum, form.Text)
	formList.AddField("所属产品", "product", db.Varchar, form.Text)
	formList.AddField("创建人", "create_user", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime)

	formList.SetTable("ai_task").SetTitle("智能任务").SetDescription("AiTask")

	return aiTask
}
