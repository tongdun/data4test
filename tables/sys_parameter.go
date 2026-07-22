package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	template2 "html/template"
)

func GetSysParameterTable(ctx *context.Context) table.Table {

	sysParameter := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := sysParameter.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	kTypes := biz.GetKnowledgeType()
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("sys_parameter.value_list"), "value_list", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton(template2.HTML(biz.T("sys_parameter.btn_sync_knowledge")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/sync_knowledge",
		Title:  biz.T("sys_parameter.btn_sync_knowledge"),
		Width:  "900px",
		Height: "340px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		panel.AddField(biz.T("sys_parameter.sync_type"), "k_type", db.Varchar, form.SelectSingle).
			FieldOptions(kTypes)
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/sync_knowledge"))

	info.SetTable("sys_parameter").SetTitle(biz.T("sys_parameter.title")).SetDescription(biz.T("sys_parameter.description"))

	formList := sysParameter.GetForm()

	helpMsg := template2.HTML(biz.T("sys_parameter.help_value_list"))

	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text)
	formList.AddField(biz.T("sys_parameter.value_list"), "value_list", db.Longtext, form.TextArea).FieldHelpMsg(helpMsg)
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("sys_parameter").SetTitle(biz.T("sys_parameter.title")).SetDescription(biz.T("sys_parameter.description"))

	return sysParameter
}
