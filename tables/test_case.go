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
	"html/template"

	//"strings"
	"github.com/GoAdminGroup/go-admin/context"
	template2 "html/template"
)

func GetTestCaseTable(ctx *context.Context) table.Table {
	products := biz.GetProducts()
	sceneNames := biz.GetScenes()
	caseTypes := biz.GetTestcaseType()
	testCase := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := testCase.GetInfo().HideFilterArea().SetFilterFormLayout(form.LayoutThreeCol)
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.case_number"), "case_number", db.Varchar).
		FieldWidth(150).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField(biz.T("common.case_title"), "case_name", db.Varchar).
		FieldWidth(150).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.case_module"), "module", db.Varchar).
		FieldWidth(120).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.case_type"), "case_type", db.Varchar).
		FieldWidth(120).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField(biz.T("common.case_level"), "priority", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField(biz.T("test_case.precondition"), "pre_condition", db.Longtext).
		FieldWidth(120).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTMLEscapeString(model.Value)
		})
	info.AddField(biz.T("common.test_range"), "test_range", db.Longtext).
		FieldWidth(120)
	info.AddField(biz.T("common.test_step"), "test_steps", db.Longtext).
		FieldWidth(250).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTMLEscapeString(model.Value)
		})
	info.AddField(biz.T("common.expected_result"), "expect_result", db.Longtext).
		FieldWidth(250)
	info.AddField(biz.T("common.auto_label"), "auto", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return biz.T("common.no")
			} else if model.Value == "1" {
				return biz.T("common.yes")
			} else if model.Value == "2" {
				return biz.T("common.part")
			}
			return biz.T("common.yes")
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "0", Text: biz.T("common.no")},
		{Value: "1", Text: biz.T("common.yes")},
		{Value: "2", Text: biz.T("common.part")},
	})
	info.AddField(biz.T("case.func_devleoper"), "fun_developer", db.Varchar).
		FieldHide().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.case_designer"), "case_designer", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldHide()
	info.AddField(biz.T("test_case.case_executor"), "case_executor", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldHide()
	info.AddField(biz.T("common.test_result"), "test_result", db.Varchar).
		FieldWidth(120).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "pass", Text: biz.T("common.pass")},
			{Value: "part", Text: biz.T("test_case.result_part")},
			{Value: "fail", Text: biz.T("common.fail")},
			{Value: "untest", Text: biz.T("test_case.result_untest")},
			{Value: "deprecated", Text: biz.T("common.status_discarded")},
		}).
		FieldEditAble(editType.Select).
		FieldEditOptions(types.FieldOptions{
			{Value: "pass", Text: biz.T("common.pass")},
			{Value: "part", Text: biz.T("test_case.result_part")},
			{Value: "fail", Text: biz.T("common.fail")},
			{Value: "untest", Text: biz.T("test_case.result_untest")},
			{Value: "deprecated", Text: biz.T("common.status_discarded")},
		}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "pass" {
				return biz.T("common.pass")
			} else if model.Value == "part" {
				return biz.T("test_case.result_part")
			} else if model.Value == "fail" {
				return biz.T("common.fail")
			} else if model.Value == "untest" {
				return biz.T("test_case.result_untest")
			} else if model.Value == "deprecated" {
				return biz.T("common.status_discarded")
			}
			return biz.T("test_case.result_untest")
		})
	info.AddField(biz.T("test_case.test_time"), "test_time", db.Varchar).
		FieldWidth(120)
	info.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar)
	info.AddField(biz.T("common.scene"), "scene", db.Varchar).
		FieldWidth(120).
		FieldEditAble(editType.Select).
		FieldEditOptions(sceneNames).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(sceneNames).
		FieldHide()

	info.AddField(biz.T("common.product"), "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(products).FieldWidth(120)
	info.AddField(biz.T("common.remark"), "remark", db.Varchar).
		FieldHide().FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).FieldWidth(120)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddSelectBox(biz.T("common.product"), products, action.FieldFilter("product"))

	info.AddSelectBox(biz.T("common.test_result"), types.FieldOptions{
		{Value: "pass", Text: biz.T("common.pass")},
		{Value: "fail", Text: biz.T("common.fail")},
		{Value: "untest", Text: biz.T("test_case.result_untest")},
		{Value: "deprecated", Text: biz.T("common.status_discarded")},
	}, action.FieldFilter("test_result"))

	info.AddButton(template2.HTML(biz.T("common.btn_export_md")), icon.File, action.Ajax("test_case_export_markdown",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			fileName, err := biz.ExportTestCase2Markdown(idStr, "test_case")
			if err != nil {
				status = fmt.Sprintf(biz.T("error.export_fail"), err)
				return false, status, ""
			}
			hostIp := ctx.Request.Host
			//status = fmt.Sprintf("导出成功\n请至[文件-用例文件]下载\n文件名为: %s", fileName)
			downloadUrl := fmt.Sprintf("http://%s/admin/fm/case/download?path=/%s", hostIp, fileName)
			status = fmt.Sprintf(biz.T("common.operate_success"), downloadUrl)
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("common.btn_export_xmind")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/test_case_export_xmind",
		Title:  biz.T("common.tilte_export_xmind"),
		Width:  "900px",
		Height: "720px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		products := biz.GetProducts() // 子域
		panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldDisplayButCanNotEditWhenCreate().
			FieldHelpMsg(template.HTML(biz.T("common.help_select_or_filter")))
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products)
		panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("common.help_version_multi")))
		panel.AddField(biz.T("common.case_designer"), "case_designer", db.Varchar, form.Text)
		panel.AddField(biz.T("common.case_module"), "module", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("common.help_module_multi")))
		panel.AddField(biz.T("common.created_at"), "created_at", db.Varchar, form.DatetimeRange)
		panel.AddField(biz.T("common.case_source"), "source", db.Varchar, form.Text).FieldDefault("test_case").FieldHide()

		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/test_case_export_xmind"))

	info.AddButton(template2.HTML(biz.T("test_case.btn_xmind2excel")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/testcase_xmind2excel",
		Title:  biz.T("test_case.title_xmind2excel"),
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		panel.AddField(biz.T("common.upload_file"), "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
			"maxFileCount": 1,
		}).FieldHelpMsg(template2.HTML(biz.T("common.help_xmind")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/testcase_xmind2excel"))

	info.AddButton(template2.HTML(biz.T("test_case.btn_import_xmind")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/testcase_xmind2import",
		Title:  biz.T("test_case.btn_import_xmind"),
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		info.AddField(biz.T("common.product"), "product", db.Varchar).
			FieldFilterable(types.FilterType{FormType: form.Select}).
			FieldFilterOptions(products)
		info.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar).
			FieldFilterable()
		panel.AddField(biz.T("common.upload_file"), "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
			"maxFileCount": 1,
		}).FieldHelpMsg(template2.HTML(biz.T("common.help_xmind")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/testcase_xmind2import"))

	info.SetTable("test_case").SetTitle(biz.T("test_case.title")).SetDescription(biz.T("test_case.description"))

	formList := testCase.GetForm()

	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.case_number"), "case_number", db.Varchar, form.Text)
	formList.AddField(biz.T("common.case_title"), "case_name", db.Varchar, form.Text)
	formList.AddField(biz.T("common.case_module"), "module", db.Varchar, form.Text)
	formList.AddField(biz.T("common.case_type"), "case_type", db.Varchar, form.Select).
		FieldOptions(caseTypes)
	formList.AddField(biz.T("common.case_level"), "priority", db.Varchar, form.Text)
	formList.AddField(biz.T("test_case.precondition"), "pre_condition", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.test_range"), "test_range", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.test_step"), "test_steps", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.expected_result"), "expect_result", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.is_auto"), "auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "0", Text: biz.T("common.no")},
			{Value: "1", Text: biz.T("common.yes")},
			{Value: "2", Text: biz.T("common.part")},
		}).FieldDefault("0")
	formList.AddField(biz.T("case.func_devleoper"), "fun_developer", db.Varchar, form.Text)
	formList.AddField(biz.T("common.case_designer"), "case_designer", db.Varchar, form.Text)
	formList.AddField(biz.T("test_case.case_executor"), "case_executor", db.Varchar, form.Text)
	formList.AddField(biz.T("common.test_result"), "test_result", db.Varchar, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "pass", Text: biz.T("common.pass")},
			{Value: "part", Text: biz.T("test_case.result_part")},
			{Value: "fail", Text: biz.T("common.fail")},
			{Value: "untest", Text: biz.T("test_case.result_untest")},
			{Value: "deprecated", Text: biz.T("common.status_discarded")},
		}).FieldDefault("untest")
	formList.AddField(biz.T("test_case.test_time"), "test_time", db.Varchar, form.Datetime)

	formList.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text)
	formList.AddField(biz.T("common.scene"), "scene", db.Varchar, form.Select).
		FieldOptions(sceneNames)
	formList.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
		FieldOptions(products)
	formList.AddField(biz.T("common.remark"), "remark", db.Varchar, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("test_case").SetTitle(biz.T("test_case.title")).SetDescription(biz.T("test_case.description"))

	detail := testCase.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.case_number"), "case_number", db.Varchar)
	detail.AddField(biz.T("common.case_title"), "case_name", db.Varchar)
	detail.AddField(biz.T("common.case_module"), "module", db.Varchar)
	detail.AddField(biz.T("common.case_type"), "case_type", db.Varchar)
	detail.AddField(biz.T("common.case_level"), "priority", db.Varchar)
	detail.AddField(biz.T("test_case.precondition"), "pre_condition", db.Longtext)
	detail.AddField(biz.T("common.test_range"), "test_range", db.Longtext)
	detail.AddField(biz.T("common.test_step"), "test_steps", db.Longtext)
	detail.AddField(biz.T("common.expected_result"), "expect_result", db.Longtext)
	detail.AddField(biz.T("common.auto_label"), "auto", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return biz.T("common.no")
			} else if model.Value == "1" {
				return biz.T("common.yes")
			} else if model.Value == "2" {
				return biz.T("common.part")
			}
			return biz.T("common.yes")
		})
	detail.AddField(biz.T("case.func_devleoper"), "fun_developer", db.Varchar)
	detail.AddField(biz.T("common.case_designer"), "case_designer", db.Varchar)
	detail.AddField(biz.T("test_case.case_executor"), "case_executor", db.Varchar)
	detail.AddField(biz.T("test_case.test_time"), "test_time", db.Varchar)
	detail.AddField(biz.T("common.test_result"), "test_result", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "pass" {
				return biz.T("common.pass")
			} else if model.Value == "part" {
				return biz.T("test_case.result_part")
			} else if model.Value == "fail" {
				return biz.T("common.fail")
			} else if model.Value == "untest" {
				return biz.T("test_case.result_untest")
			} else if model.Value == "deprecated" {
				return biz.T("common.status_discarded")
			}
			return biz.T("common.yes")
		})
	detail.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar)
	detail.AddField(biz.T("common.scene"), "scene", db.Varchar)
	detail.AddField(biz.T("common.product"), "product", db.Varchar)
	detail.AddField(biz.T("common.remark"), "remark", db.Varchar)

	detail.SetTable("test_case").SetTitle(biz.T("test_case.title")).SetDescription(biz.T("test_case.description"))

	return testCase
}
