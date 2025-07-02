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
	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("接口ID", "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("接口描述", "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("执行次数", "run_times", db.Int)
	info.AddField("测试次数", "test_times", db.Int)
	info.AddField("通过次数", "pass_times", db.Int)
	info.AddField("失败次数", "fail_times", db.Int)
	info.AddField("未测试次数", "untest_times", db.Int)
	info.AddField("测试结果", "test_result", db.Char).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("失败原因", "fail_reason", db.Longtext)
	info.AddField("关联应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	apps := biz.GetApps()
	info.AddSelectBox("关联应用", apps, action.FieldFilter("app"))

	info.SetTable("api_id_count").SetTitle("接口统计").SetDescription("接口统计")

	formList := apiIdCount.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("接口ID", "api_id", db.Varchar, form.Text)
	formList.AddField("接口描述", "api_desc", db.Varchar, form.Text)
	formList.AddField("执行次数", "run_times", db.Int, form.Number)
	formList.AddField("测试次数", "test_times", db.Int, form.Number)
	formList.AddField("通过次数", "pass_times", db.Int, form.Number)
	formList.AddField("失败次数", "fail_times", db.Int, form.Number)
	formList.AddField("未测试次数", "untest_times", db.Int, form.Number)
	formList.AddField("测试结果", "test_result", db.Char, form.Text)
	formList.AddField("失败原因", "fail_reason", db.Longtext, form.Text)
	formList.AddField("关联应用", "app", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_id_count").SetTitle("接口统计").SetDescription("接口统计")

	return apiIdCount
}
