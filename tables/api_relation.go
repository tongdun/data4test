package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	// "github.com/GoAdminGroup/go-admin/template/icon"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	// "strings"
	// "fmt"
)

func GetApiRelationTable(ctx *context.Context) table.Table {

	apiRelation := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiRelation.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("接口ID", "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("接口描述", "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("所属模块", "api_module", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("所属应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("是否自动化", "auto", db.Enum).
		FieldFilterable().FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "yes" {
			return "是"
		}
		if model.Value == "no" {
			return "否"
		}
		return "否"
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "yes", Text: "是"},
		{Value: "no", Text: "否"},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "yes", Text: "是"},
		{Value: "no", Text: "否"},
	})
	info.AddField("前置接口", "pre_apis", db.Varchar)
	info.AddField("提供变量关系", "out_vars", db.Varchar)
	info.AddField("校验变量转换关系", "check_vars", db.Varchar)
	info.AddField("依赖参数关联接口", "param_apis", db.Varchar)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	apps := biz.GetApps()
	info.AddSelectBox("所属应用", apps, action.FieldFilter("app"))

	info.SetTable("api_relation").SetTitle("接口关系").SetDescription("接口关系")

	formList := apiRelation.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("接口ID", "api_id", db.Varchar, form.Text)
	formList.AddField("接口描述", "api_desc", db.Varchar, form.Text)
	formList.AddField("所属模块", "api_module", db.Varchar, form.Text)
	formList.AddField("所属应用", "app", db.Varchar, form.Text)
	formList.AddField("是否自动化", "auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "是", Value: "yes"},
			{Text: "否", Value: "no"},
		}).FieldDefault("yes")

	formList.AddField("前置接口", "pre_apis", db.Varchar, form.Text)
	formList.AddField("提供变量关系", "out_vars", db.Varchar, form.Text)
	formList.AddField("校验变量转换关系", "check_vars", db.Varchar, form.Text)
	formList.AddField("依赖参数关联接口", "param_apis", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_relation").SetTitle("接口关系").SetDescription("接口关系")

	return apiRelation
}
