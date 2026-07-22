package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	// editType "github.com/GoAdminGroup/go-admin/template/types/table"
)

func GetFuzzingDefinitionTable(ctx *context.Context) table.Table {

	fuzzingDefinition := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := fuzzingDefinition.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField(biz.T("common.value"), "value", db.Longtext)
	info.AddField(biz.T("common.value_type"), "type", db.Enum).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("fuzzing_definition").SetTitle(biz.T("fuzzing_definition.title")).SetDescription(biz.T("fuzzing_definition.description"))

	formList := fuzzingDefinition.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text)
	formList.AddField(biz.T("common.value"), "value", db.Longtext, form.Text)
	formList.AddField(biz.T("common.value_type"), "type", db.Enum, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("fuzzing_definition").SetTitle(biz.T("fuzzing_definition.title")).SetDescription(biz.T("fuzzing_definition.description"))

	return fuzzingDefinition
}
