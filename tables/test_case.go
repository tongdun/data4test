package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
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
	info.AddField("用例编号", "case_number", db.Varchar).
		FieldWidth(150).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField("用例名称", "case_name", db.Varchar).
		FieldWidth(150).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("用例模块", "module", db.Varchar).
		FieldWidth(120).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("用例类型", "case_type", db.Varchar).
		FieldWidth(120).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField("优先级", "priority", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField("预置条件", "pre_condition", db.Longtext).
		FieldWidth(120).
		FieldHide()
	info.AddField("测试范围", "test_range", db.Longtext).
		FieldWidth(120).
		FieldHide()
	info.AddField("测试步骤", "test_steps", db.Longtext).
		FieldWidth(200).
		FieldHide()
	info.AddField("预期结果", "expect_result", db.Longtext).
		FieldWidth(200).
		FieldHide()
	info.AddField("自动化", "auto", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return "否"
			} else if model.Value == "1" {
				return "是"
			} else if model.Value == "2" {
				return "部分是"
			}
			return "是"
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "0", Text: "否"},
		{Value: "1", Text: "是"},
		{Value: "2", Text: "部分是"},
	})
	info.AddField("功能开发者", "fun_developer", db.Varchar).
		FieldHide().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("用例设计者", "case_designer", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldHide()
	info.AddField("用例执行者", "case_executor", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldHide()
	info.AddField("测试结果", "test_result", db.Varchar).
		FieldWidth(120).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "pass", Text: "通过"},
			{Value: "part", Text: "部分通过"},
			{Value: "fail", Text: "失败"},
			{Value: "untest", Text: "未测试"},
			{Value: "deprecated", Text: "已废弃"},
		}).
		FieldEditAble(editType.Select).
		FieldEditOptions(types.FieldOptions{
			{Value: "pass", Text: "通过"},
			{Value: "part", Text: "部分通过"},
			{Value: "fail", Text: "失败"},
			{Value: "untest", Text: "未测试"},
			{Value: "deprecated", Text: "已废弃"},
		}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "pass" {
				return "通过"
			} else if model.Value == "part" {
				return "部分通过"
			} else if model.Value == "fail" {
				return "失败"
			} else if model.Value == "untest" {
				return "未测试"
			} else if model.Value == "deprecated" {
				return "已废弃"
			}
			return "未测试"
		})
	info.AddField("测试时间", "test_time", db.Varchar).
		FieldWidth(120)
	info.AddField("引入版本", "intro_version", db.Varchar)
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

	info.AddButton("导出Markdown", icon.File, action.Ajax("test_case_export_markdown",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再导出"
				return false, status, ""
			}

			fileName, err := biz.ExportTestCase2Markdown(idStr)
			if err != nil {
				status = fmt.Sprintf("导出失败: %s", err)
				return false, status, ""
			}

			status = fmt.Sprintf("导出成功\n请至[文件-用例文件]下载\n文件名为: %s", fileName)
			return true, status, ""
		}))

	info.SetTable("test_case").SetTitle("测试用例").SetDescription("测试用例")

	formList := testCase.GetForm()

	formList.AddField("自动主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("用例编号", "case_number", db.Varchar, form.Text)
	formList.AddField("用例名称", "case_name", db.Varchar, form.Text)
	formList.AddField("用例模块", "module", db.Varchar, form.Text)
	formList.AddField("用例类型", "case_type", db.Varchar, form.Select).
		FieldOptions(caseTypes)
	formList.AddField("优先级", "priority", db.Varchar, form.Text)
	formList.AddField("预置条件", "pre_condition", db.Longtext, form.TextArea)
	formList.AddField("测试范围", "test_range", db.Longtext, form.TextArea)
	formList.AddField("测试步骤", "test_steps", db.Longtext, form.TextArea)
	formList.AddField("预期结果", "expect_result", db.Longtext, form.TextArea)
	formList.AddField("是否自动化", "auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "0", Text: "否"},
			{Value: "1", Text: "是"},
			{Value: "2", Text: "部分是"},
		}).FieldDefault("0")
	formList.AddField("功能开发者", "fun_developer", db.Varchar, form.Text)
	formList.AddField("用例设计者", "case_designer", db.Varchar, form.Text)
	formList.AddField("用例执行者", "case_executor", db.Varchar, form.Text)
	formList.AddField("测试结果", "test_result", db.Varchar, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "pass", Text: "通过"},
			{Value: "part", Text: "部分通过"},
			{Value: "fail", Text: "失败"},
			{Value: "untest", Text: "未测试"},
			{Value: "deprecated", Text: "已废弃"},
		}).FieldDefault("untest")
	formList.AddField("测试时间", "test_time", db.Varchar, form.Datetime)

	formList.AddField("引入版本", "intro_version", db.Varchar, form.Text)
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
	detail.AddField("用例模块", "module", db.Varchar)
	detail.AddField("用例类型", "case_type", db.Varchar)
	detail.AddField("优先级", "priority", db.Varchar)
	detail.AddField("预置条件", "pre_condition", db.Longtext)
	detail.AddField("测试范围", "test_range", db.Longtext)
	detail.AddField("测试步骤", "test_steps", db.Longtext)
	detail.AddField("预期结果", "expect_result", db.Longtext)
	detail.AddField("是否自动化", "auto", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return "否"
			} else if model.Value == "1" {
				return "是"
			} else if model.Value == "2" {
				return "部分是"
			}
			return "是"
		})
	detail.AddField("功能开发者", "fun_developer", db.Varchar)
	detail.AddField("用例设计者", "case_designer", db.Varchar)
	detail.AddField("用例执行者", "case_executor", db.Varchar)
	detail.AddField("测试时间", "test_time", db.Varchar)
	detail.AddField("测试结果", "test_result", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "pass" {
				return "通过"
			} else if model.Value == "part" {
				return "部分通过"
			} else if model.Value == "fail" {
				return "失败"
			} else if model.Value == "untest" {
				return "未测试"
			} else if model.Value == "deprecated" {
				return "已废弃"
			}
			return "是"
		})
	detail.AddField("引入版本", "intro_version", db.Varchar)
	detail.AddField("关联场景", "scene", db.Varchar)
	detail.AddField("关联产品", "product", db.Varchar)
	detail.AddField("备注", "remark", db.Varchar)

	detail.SetTable("test_case").SetTitle("测试用例").SetDescription("测试用例")

	return testCase
}
