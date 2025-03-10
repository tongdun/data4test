package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"html/template"
)

func GetAssertTemplateTable(ctx *context.Context) table.Table {

	assertTemplate := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := assertTemplate.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name

	info.AddField("自增主键", "id", db.Int).
		FieldSortable()
	info.AddField("模板名称", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("定义信息", "value", db.Longtext)
	info.AddField("备注", "remark", db.Varchar)
	info.AddField("创建人", "user_name", db.Varchar).
		FieldFilterable()
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("assert_template").SetTitle("断言值模板").SetDescription("断言值模板")

	formList := assertTemplate.GetForm()

	helpMsg := template.HTML("JSON格式，e.g.: {\"default\": \"value1\", \"ch\": \"value2\", \"en\": \"value3\"}<br>普通格式, e.g.: value<br>根据请求头里的语言取对应语言的值，如未定义语言模式，优先级: default > ch > en > other 即：默认 > 中文 > 英文 > 其他语言<br>可以根据需要进行不同语言断言值的定义<br>普通格式的定义会直接把整体值当时default取用")

	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("模板名称", "name", db.Varchar, form.Text)
	formList.AddField("定义信息", "value", db.Longtext, form.TextArea).FieldHelpMsg(helpMsg)
	formList.AddField("备注", "remark", db.Varchar, form.Text)
	formList.AddField("创建人", "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.SetTable("assert_template").SetTitle("断言值模板").SetDescription("断言值模板")

	return assertTemplate
}
