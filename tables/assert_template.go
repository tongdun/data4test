package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	template2 "html/template"
)

func GetAssertTemplateTable(ctx *context.Context) table.Table {

	assertTemplate := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := assertTemplate.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldSortable()
	info.AddField(biz.T("common.name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.define_info"), "value", db.Longtext)
	info.AddField(biz.T("common.remark"), "remark", db.Varchar)
	info.AddField(biz.T("common.user_name"), "user_name", db.Varchar).
		FieldFilterable()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("assert_template").SetTitle(biz.T("assert_template.title")).SetDescription(biz.T("assert_template.description"))

	formList := assertTemplate.GetForm()

	helpMsg := template2.HTML(biz.T("assert_template.help_value"))

	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text)
	formList.AddField(biz.T("common.define_info"), "value", db.Longtext, form.TextArea).FieldHelpMsg(helpMsg)
	formList.AddField(biz.T("common.remark"), "remark", db.Varchar, form.Text)
	formList.AddField(biz.T("common.user_name"), "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.SetTable("assert_template").SetTitle(biz.T("assert_template.title")).SetDescription(biz.T("assert_template.description"))

	return assertTemplate
}
