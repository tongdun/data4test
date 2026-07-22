package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
)

func GetApiRelationTable(ctx *context.Context) table.Table {

	apiRelation := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiRelation.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.api_module"), "api_module", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.auto_label"), "auto", db.Enum).
		FieldFilterable().FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "yes" {
			return biz.T("common.yes")
		}
		if model.Value == "no" {
			return biz.T("common.no")
		}
		return biz.T("common.no")
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "yes", Text: biz.T("common.yes")},
		{Value: "no", Text: biz.T("common.no")},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "yes", Text: biz.T("common.yes")},
		{Value: "no", Text: biz.T("common.no")},
	})
	info.AddField(biz.T("api_relation.pre_apis"), "pre_apis", db.Varchar)
	info.AddField(biz.T("common.output_params"), "out_vars", db.Varchar)
	info.AddField(biz.T("api_relation.check_vars"), "check_vars", db.Varchar)
	info.AddField(biz.T("api_relation.param_apis"), "param_apis", db.Varchar)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	apps := biz.GetApps()
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))

	info.SetTable("api_relation").SetTitle(biz.T("api_relation.title")).SetDescription(biz.T("api_relation.description"))

	formList := apiRelation.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.api_id"), "api_id", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_module"), "api_module", db.Varchar, form.Text)
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.Text)
	formList.AddField(biz.T("common.auto_label"), "auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.yes"), Value: "yes"},
			{Text: biz.T("common.no"), Value: "no"},
		}).FieldDefault("yes")

	formList.AddField(biz.T("api_relation.pre_apis"), "pre_apis", db.Varchar, form.Text)
	formList.AddField(biz.T("common.output_params"), "out_vars", db.Varchar, form.Text)
	formList.AddField(biz.T("api_relation.check_vars"), "check_vars", db.Varchar, form.Text)
	formList.AddField(biz.T("api_relation.param_apis"), "param_apis", db.Varchar, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_relation").SetTitle(biz.T("api_relation.title")).SetDescription(biz.T("api_relation.description"))

	return apiRelation
}
