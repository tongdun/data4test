package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetApiIdCountTable(ctx *context.Context) table.Table {

	apiIdCount := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiIdCount.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.run_time"), "run_times", db.Int)
	info.AddField(biz.T("api_id_count.test_times"), "test_times", db.Int)
	info.AddField(biz.T("api_id_count.pass_times"), "pass_times", db.Int)
	info.AddField(biz.T("api_id_count.fail_times"), "fail_times", db.Int)
	info.AddField(biz.T("api_id_count.untest_times"), "untest_times", db.Int)
	info.AddField(biz.T("api_id_count.test_result"), "test_result", db.Char).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext)
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	apps := biz.GetApps()
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))

	info.SetTable("api_id_count").SetTitle(biz.T("api_id_count.title")).SetDescription(biz.T("api_id_count.description"))

	formList := apiIdCount.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.api_id"), "api_id", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar, form.Text)
	formList.AddField(biz.T("common.run_time"), "run_times", db.Int, form.Number)
	formList.AddField(biz.T("api_id_count.test_times"), "test_times", db.Int, form.Number)
	formList.AddField(biz.T("api_id_count.pass_times"), "pass_times", db.Int, form.Number)
	formList.AddField(biz.T("api_id_count.fail_times"), "fail_times", db.Int, form.Number)
	formList.AddField(biz.T("api_id_count.untest_times"), "untest_times", db.Int, form.Number)
	formList.AddField(biz.T("api_id_count.test_result"), "test_result", db.Char, form.Text)
	formList.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext, form.Text)
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_id_count").SetTitle(biz.T("api_id_count.title")).SetDescription(biz.T("api_id_count.description"))

	return apiIdCount
}
