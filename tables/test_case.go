package tables

import (
	"data4perf/biz"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"

	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	//"strings"
	"github.com/GoAdminGroup/go-admin/context"
)

func GetTestCaseTable(ctx *context.Context) table.Table {
	products := biz.GetProducts()
	sceneNames := biz.GetScenes()
	caseTypes := biz.GetTestcaseType()
	testCase := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := testCase.GetInfo().HideFilterArea().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	info.AddField("自动主键", "id", db.Int).
		FieldHide()
	info.AddField("用例编号", "case_number", db.Varchar).FieldWidth(150).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldSortable()
	info.AddField("用例名称", "case_name", db.Varchar).FieldWidth(150).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("用例类型", "case_type", db.Varchar).FieldWidth(120).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldSortable()
	info.AddField("优先级", "priority", db.Varchar).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldSortable()
	info.AddField("预置条件", "pre_condition", db.Longtext).FieldWidth(120).
		FieldHide()
	info.AddField("测试范围", "test_range", db.Longtext).FieldWidth(120).
		FieldHide()
	info.AddField("测试步骤", "test_steps", db.Longtext).FieldWidth(200).
		FieldHide()
	info.AddField("预期结果", "expect_result", db.Longtext).FieldWidth(200).
		FieldHide()
	info.AddField("自动化", "auto", db.Varchar).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("功能开发者", "fun_developer", db.Varchar).
		FieldHide().FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("用例设计者", "case_designer", db.Varchar).
		FieldHide().FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("用例执行者", "case_executor", db.Varchar).
		FieldHide().FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("测试时间", "test_time", db.Varchar).FieldWidth(120).FieldEditAble(editType.Text)
	info.AddField("测试结果", "test_result", db.Varchar).
		FieldWidth(120).
		FieldEditAble(editType.Select).FieldEditOptions(types.FieldOptions{
		{Value: "pass", Text: "通过"},
		{Value: "part", Text: "部分通过"},
		{Value: "fail", Text: "失败"},
		{Value: "untest", Text: "未测试"},
		{Value: "deprecated", Text: "已废弃"},
	}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "通过"},
		{Value: "part", Text: "部分通过"},
		{Value: "fail", Text: "失败"},
		{Value: "untest", Text: "未测试"},
		{Value: "deprecated", Text: "已废弃"},
	})

	info.AddField("用例模块", "module", db.Varchar).FieldWidth(120).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("引入版本", "version", db.Varchar)
	info.AddField("关联场景", "scene", db.Varchar).
		FieldWidth(120).
		FieldEditAble(editType.Select).
		FieldEditOptions(sceneNames).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(sceneNames).
		FieldHide()

	info.AddField("关联产品", "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(products).FieldWidth(120)
	info.AddField("备注", "remark", db.Varchar).
		FieldHide().FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("创建时间", "created_at", db.Timestamp).FieldWidth(120)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddSelectBox("关联产品", products, action.FieldFilter("product"))

	info.AddSelectBox("测试结果", types.FieldOptions{
		{Value: "pass", Text: "通过"},
		{Value: "fail", Text: "失败"},
		{Value: "untest", Text: "未测试"},
		{Value: "deprecated", Text: "已废弃"},
	}, action.FieldFilter("test_result"))

	info.SetTable("test_case").SetTitle("测试用例").SetDescription("测试用例")

	formList := testCase.GetForm()

	formList.AddField("自动主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("用例编号", "case_number", db.Varchar, form.Text)
	formList.AddField("用例名称", "case_name", db.Varchar, form.Text)
	formList.AddField("用例类型", "case_type", db.Varchar, form.Select).FieldOptions(caseTypes).FieldDefault("功能")
	formList.AddField("优先级", "priority", db.Varchar, form.Text)
	formList.AddField("预置条件", "pre_condition", db.Longtext, form.RichText)
	formList.AddField("测试范围", "test_range", db.Longtext, form.RichText)
	formList.AddField("测试步骤", "test_steps", db.Longtext, form.RichText)
	formList.AddField("预期结果", "expect_result", db.Longtext, form.RichText)
	formList.AddField("是否自动化", "auto", db.Varchar, form.Text)
	formList.AddField("功能开发者", "fun_developer", db.Varchar, form.Text)
	formList.AddField("用例设计者", "case_designer", db.Varchar, form.Text)
	formList.AddField("用例执行者", "case_executor", db.Varchar, form.Text)
	formList.AddField("测试时间", "test_time", db.Varchar, form.Date)
	formList.AddField("测试结果", "test_result", db.Varchar, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "pass", Text: "通过"},
			{Value: "part", Text: "部分通过"},
			{Value: "fail", Text: "失败"},
			{Value: "untest", Text: "未测试"},
			{Value: "deprecated", Text: "已废弃"},
		}).FieldDefault("untest")

	formList.AddField("用例模块", "module", db.Varchar, form.Text)
	formList.AddField("引入版本", "version", db.Varchar, form.Text)
	formList.AddField("关联场景", "scene", db.Varchar, form.Select).
		FieldOptions(sceneNames)
	formList.AddField("关联产品", "product", db.Varchar, form.SelectSingle).
		FieldOptions(products)
	formList.AddField("备注", "remark", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("test_case").SetTitle("测试用例").SetDescription("测试用例")

	detail := testCase.GetDetail()
	detail.AddField("自动主键", "id", db.Int)
	detail.AddField("用例编号", "case_number", db.Varchar)
	detail.AddField("用例名称", "case_name", db.Varchar)
	detail.AddField("用例类型", "case_type", db.Varchar)
	detail.AddField("优先级", "priority", db.Varchar)
	detail.AddField("预置条件", "pre_condition", db.Longtext)
	detail.AddField("测试范围", "test_range", db.Longtext)
	detail.AddField("测试步骤", "test_steps", db.Longtext)
	detail.AddField("预期结果", "expect_result", db.Longtext)
	detail.AddField("是否自动化", "auto", db.Varchar)
	detail.AddField("功能开发者", "fun_developer", db.Varchar)
	detail.AddField("用例设计者", "case_designer", db.Varchar)
	detail.AddField("用例执行者", "case_executor", db.Varchar)
	detail.AddField("测试时间", "test_time", db.Varchar)
	detail.AddField("测试结果", "test_result", db.Varchar)

	detail.AddField("用例模块", "module", db.Varchar)
	detail.AddField("引入版本", "version", db.Varchar)
	detail.AddField("关联场景", "scene", db.Varchar)
	detail.AddField("关联产品", "product", db.Varchar)
	detail.AddField("备注", "remark", db.Varchar)

	detail.SetTable("test_case").SetTitle("测试用例").SetDescription("测试用例")

	return testCase
}
