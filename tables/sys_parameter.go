package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"html/template"
)

func GetSysParameterTable(ctx *context.Context) table.Table {

	sysParameter := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := sysParameter.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("参数名称", "name", db.Varchar)
	info.AddField("值定义", "value_list", db.Longtext)
	info.AddField("备注", "remark", db.Longtext)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("sys_parameter").SetTitle("系统参数").SetDescription("系统参数")

	formList := sysParameter.GetForm()

	helpMsg := template.HTML("JSON格式，e.g.: {\"default\": \"v1,v2,v3,v4,……\", \"ch\": \"v1,v2,v3,v4,……\", \"en\": \"v1,v2,v3,v4,……\"}<br>普通格式, e.g.: v1,v2,v3,v4,……<br>根据请求头里的语言取对应语言的值，如未定义语言模式，优先级: default > ch > en > other 即：默认 > 中文 > 英文 > 其他语言<br>可以根据需要进行不同语言系统参数的定义<br>普通格式的定义会直接把整体值当时default取用")

	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("参数名称", "name", db.Varchar, form.Text)
	formList.AddField("值定义", "value_list", db.Longtext, form.TextArea).FieldHelpMsg(helpMsg)
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("sys_parameter").SetTitle("系统参数").SetDescription("系统参数")

	return sysParameter
}
