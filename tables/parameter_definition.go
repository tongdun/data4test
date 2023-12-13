package tables

import (
	"data4perf/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetParameterDefinitionTable(ctx *context.Context) table.Table {

	parameterDefinition := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := parameterDefinition.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("参数名称/接口ID", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("值", "value", db.Longtext)
	info.AddField("备注", "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("关联应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110)
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	apps := biz.GetApps()
	info.AddSelectBox("关联应用", apps, action.FieldFilter("app"))

	info.SetTable("parameter_definition").SetTitle("参数定义").SetDescription("参数定义")

	formList := parameterDefinition.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("参数名称/接口ID", "name", db.Varchar, form.Text)
	formList.AddField("值", "value", db.Longtext, form.Text)
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("关联应用", "app", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("parameter_definition").SetTitle("参数定义").SetDescription("参数定义")

	return parameterDefinition
}
