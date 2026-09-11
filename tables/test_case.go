package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"html/template"

	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	template2 "html/template"
)

func GetTestCaseTable(ctx *context.Context) table.Table {
	products := biz.GetCaseProducts()
	sceneNames := biz.GetScenes()
	caseTypes := biz.GetTestcaseType()
	dataLocale := biz.GetDataLocale(ctx.Cookie("data_locale"), biz.GetLocale())
	userName := auth.Auth(ctx).Name
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
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldDisplay(caseFieldDisplay(dataLocale, "case_name", false))
	info.AddField(biz.T("common.case_module"), "module", db.Varchar).
		FieldWidth(120).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldDisplay(caseFieldDisplay(dataLocale, "module", false))
	info.AddField(biz.T("common.case_type"), "case_type", db.Varchar).
		FieldWidth(120).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable().
		FieldHide()
	info.AddField(biz.T("common.case_level"), "priority", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable().
		FieldHide()
	info.AddField(biz.T("test_case.precondition"), "pre_condition", db.Longtext).
		FieldWidth(120).
		FieldDisplay(caseFieldDisplay(dataLocale, "pre_condition", true))
	info.AddField(biz.T("common.test_range"), "test_range", db.Longtext).
		FieldWidth(120).
		FieldDisplay(caseFieldDisplay(dataLocale, "test_range", false))
	info.AddField(biz.T("common.test_step"), "test_steps", db.Longtext).
		FieldWidth(250).
		FieldDisplay(caseFieldDisplay(dataLocale, "test_steps", true))
	info.AddField(biz.T("common.expected_result"), "expect_result", db.Longtext).
		FieldWidth(250).
		FieldDisplay(caseFieldDisplay(dataLocale, "expect_result", true))
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
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "0", Text: biz.T("common.no")},
			{Value: "1", Text: biz.T("common.yes")},
			{Value: "2", Text: biz.T("common.part")},
		}).
		FieldHide()
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
			{Value: "fail", Text: biz.T("common.fail")},
			{Value: "untest", Text: biz.T("test_case.result_untest")},
			{Value: "deprecated", Text: biz.T("common.status_discarded")},
			{Value: "unmerged", Text: biz.T("test_case.result_unmerged")},
		}).
		FieldEditAble(editType.Select).
		FieldEditOptions(types.FieldOptions{
			{Value: "pass", Text: biz.T("common.pass")},
			{Value: "fail", Text: biz.T("common.fail")},
			{Value: "untest", Text: biz.T("test_case.result_untest")},
			{Value: "deprecated", Text: biz.T("common.status_discarded")},
			{Value: "unmerged", Text: biz.T("test_case.result_unmerged")},
		}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "pass" {
				return biz.T("common.pass")
			} else if model.Value == "fail" {
				return biz.T("common.fail")
			} else if model.Value == "untest" {
				return biz.T("test_case.result_untest")
			} else if model.Value == "deprecated" {
				return biz.T("common.status_discarded")
			} else if model.Value == "unmerged" {
				return biz.T("test_case.result_unmerged")
			}
			return biz.T("test_case.result_untest")
		})
	info.AddField(biz.T("test_case.test_time"), "test_time", db.Varchar).
		FieldWidth(120).
		FieldHide()
	info.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar).
		FieldHide()
	info.AddField(biz.T("common.related_scene"), "scene", db.Varchar).
		FieldWidth(120).
		FieldEditAble(editType.Select).
		FieldEditOptions(sceneNames).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(sceneNames).
		FieldHide()

	info.AddField(biz.T("common.product_line"), "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(products).
		FieldWidth(120).
		FieldHide()
	info.AddField(biz.T("common.remark"), "remark", db.Varchar).
		FieldHide().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("test_case.test_process"), "test_process", db.Longtext).
		FieldWidth(150).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTMLEscapeString(model.Value)
		}).
		FieldHide()
	info.AddField(biz.T("test_case.ext_info"), "ext_info", db.Longtext).
		FieldWidth(150).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldHide()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).FieldWidth(120)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddSelectBox(biz.T("common.product_line"), products, action.FieldFilter("product"))

	info.AddSelectBox(biz.T("common.test_result"), types.FieldOptions{
		{Value: "pass", Text: biz.T("common.pass")},
		{Value: "fail", Text: biz.T("common.fail")},
		{Value: "untest", Text: biz.T("test_case.result_untest")},
		{Value: "deprecated", Text: biz.T("common.status_discarded")},
		{Value: "unmerged", Text: biz.T("test_case.result_unmerged")},
	}, action.FieldFilter("test_result"))

	info.AddButton(template2.HTML(biz.T("test_case.btn_switch_data_lang")), icon.Language, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/setDataLocale",
		Title:  biz.T("test_case.title_switch_data_lang"),
		Width:  "600px",
		Height: "240px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		langOptions := append(types.FieldOptions{
			{Value: "auto", Text: biz.T("test_case.data_lang_follow_ui")},
		}, biz.GetSupportLanguages()...)
		panel.AddField(biz.T("test_case.title_switch_data_lang"), "lang", db.Varchar, form.SelectSingle).
			FieldOptions(langOptions).FieldDefault("auto")
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/setDataLocale"))

	addDropdownButton(info, template2.HTML(biz.T("test_case.btn_export_group")), icon.Download, []DropdownItem{
		{Label: template2.HTML(biz.T("common.btn_export_md")), Icon: icon.File, Action: action.Ajax("test_case_export_markdown",
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
			})},

		{Label: template2.HTML(biz.T("common.btn_export_xmind")), Icon: icon.FolderO, Action: action.PopUpWithCtxForm(action.PopUpData{
			Id:     "/test_case_export_xmind",
			Title:  biz.T("common.tilte_export_xmind"),
			Width:  "900px",
			Height: "720px", // TextArea
		}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
			ids := ctx.FormValue("ids")
			products := biz.GetCaseProducts() // 子域
			panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldDisplayButCanNotEditWhenCreate().
				FieldHelpMsg(template.HTML(biz.T("common.help_select_or_filter")))
			panel.AddField(biz.T("common.product_line"), "product", db.Varchar, form.SelectSingle).
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
		}, "/test_case_export_xmind")},

		{Label: template2.HTML(biz.T("test_case.btn_xmind2excel")), Icon: icon.FolderO, Action: action.PopUpWithCtxForm(action.PopUpData{
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
		}, "/testcase_xmind2excel")},

		{Label: template2.HTML(biz.T("test_case.btn_export_excel")), Icon: icon.FolderO, Action: action.PopUpWithCtxForm(action.PopUpData{
			Id:     "/test_case_export_excel",
			Title:  biz.T("test_case.title_export_excel"),
			Width:  "900px",
			Height: "960px",
		}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
			ids := ctx.FormValue("ids")
			products := biz.GetCaseProducts()
			templates := biz.GetTestCaseExportTemplates()
			languages := biz.GetSupportLanguages()
			panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldDisplayButCanNotEditWhenCreate().
				FieldHelpMsg(template.HTML(biz.T("common.help_select_or_filter")))
			panel.AddField(biz.T("test_case.export_template"), "template", db.Varchar, form.SelectSingle).
				FieldOptions(templates).FieldDefault(templates[0].Value)
			panel.AddField(biz.T("test_case.export_lang"), "lang", db.Varchar, form.SelectSingle).
				FieldOptions(languages).FieldDefault(biz.GetLocale())
			panel.AddField(biz.T("test_case.screenshot_mode"), "screenshot_mode", db.Varchar, form.SelectSingle).
				FieldOptions(types.FieldOptions{
					{Value: "", Text: biz.T("test_case.screenshot_mode_none")},
					{Value: "path", Text: biz.T("test_case.screenshot_mode_path")},
					{Value: "embed", Text: biz.T("test_case.screenshot_mode_embed")},
				}).FieldDefault("").
				// 图片模式联动：嵌入图片→显示 EXCEL类型；打包图片→显示 打包格式
				// 注：导出弹窗是 PopUpWithCtxForm，不渲染 panel.FooterHtml（AddJS / FieldOnChooseShow 都会丢弃），
				// 故改用 FieldFoot 把脚本内联到该字段，才能随弹窗一起渲染。
				FieldFoot(template.HTML(`<script>
$(function () {
    function toggleExportPicFields() {
        var v = $('select.screenshot_mode').val() || '';
        if (v === 'embed') {
            $("label[for='excel_type']").parent().show();
        } else {
            $("label[for='excel_type']").parent().hide();
        }
        if (v === 'path') {
            $("label[for='pack_format']").parent().show();
        } else {
            $("label[for='pack_format']").parent().hide();
        }
    }
    toggleExportPicFields();
    $('select.screenshot_mode').on('select2:select change', toggleExportPicFields);
});
</script>`))
			panel.AddField(biz.T("test_case.excel_type"), "excel_type", db.Varchar, form.SelectSingle).
				FieldOptions(types.FieldOptions{
					{Value: "wps", Text: biz.T("test_case.excel_type_wps")},
					{Value: "office", Text: biz.T("test_case.excel_type_office")},
				}).FieldDefault("wps")
			panel.AddField(biz.T("test_case.pack_format"), "pack_format", db.Varchar, form.SelectSingle).
				FieldOptions(types.FieldOptions{
					{Value: "tgz", Text: biz.T("test_case.pack_tgz")},
					{Value: "zip", Text: biz.T("test_case.pack_zip")},
				}).FieldDefault("tgz")
			panel.AddField(biz.T("common.product_line"), "product", db.Varchar, form.SelectSingle).FieldOptions(products)
			panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text).
				FieldHelpMsg(template.HTML(biz.T("common.help_version_multi")))
			panel.AddField(biz.T("common.case_designer"), "case_designer", db.Varchar, form.Text)
			panel.AddField(biz.T("common.case_module"), "module", db.Varchar, form.Text).
				FieldHelpMsg(template.HTML(biz.T("common.help_module_multi")))
			panel.AddField(biz.T("common.created_at"), "created_at", db.Varchar, form.DatetimeRange)
			panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
			return panel
		}, "/test_case_export_excel")}})

	info.AddButton(template2.HTML(biz.T("test_case.btn_batch_update")), icon.Edit, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/test_case_batch_update_form",
		Title:  biz.T("test_case.title_batch_update"),
		Width:  "700px",
		Height: "800px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).
			FieldDefault(ids).FieldDisplayButCanNotEditWhenCreate().
			FieldHelpMsg(template.HTML(biz.T("common.help_select_or_filter")))
		panel.AddField(biz.T("common.case_module"), "module", db.Varchar, form.Text)
		panel.AddField(biz.T("common.test_result"), "test_result", db.Varchar, form.SelectSingle).
			FieldOptions(types.FieldOptions{
				{Value: "", Text: biz.T("test_case.no_change")},
				{Value: "pass", Text: biz.T("common.pass")},
				{Value: "fail", Text: biz.T("common.fail")},
				{Value: "untest", Text: biz.T("test_case.result_untest")},
				{Value: "deprecated", Text: biz.T("common.status_discarded")},
				{Value: "unmerged", Text: biz.T("test_case.result_unmerged")},
			}).
			FieldDefault("")
		panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text)
		panel.AddField(biz.T("case.func_devleoper"), "fun_developer", db.Varchar, form.Text)
		panel.AddField(biz.T("common.case_designer"), "case_designer", db.Varchar, form.Text)
		panel.AddField(biz.T("test_case.case_executor"), "case_executor", db.Varchar, form.Text)
		panel.AddField(biz.T("test_case.test_time"), "test_time", db.Varchar, form.Datetime)
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/test_case_batch_update"))

	addDropdownButton(info, template2.HTML(biz.T("test_case.btn_import_group")), icon.Upload, []DropdownItem{
		{Label: template2.HTML(biz.T("test_case.btn_import_xmind")), Icon: icon.FolderO, Action: action.PopUpWithCtxForm(action.PopUpData{
			Id:     "/testcase_xmind2import",
			Title:  biz.T("test_case.btn_import_xmind"),
			Width:  "900px",
			Height: "680px", // TextArea
		}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
			info.AddField(biz.T("common.product_line"), "product", db.Varchar).
				FieldFilterable(types.FilterType{FormType: form.Select}).
				FieldFilterOptions(products)
			info.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar).
				FieldFilterable()
			panel.AddField(biz.T("common.upload_file"), "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
				"maxFileCount": 1,
			}).FieldHelpMsg(template2.HTML(biz.T("common.help_xmind")))
			panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
			return panel
		}, "/testcase_xmind2import")},

		{Label: template2.HTML(biz.T("test_case.btn_import_excel")), Icon: icon.FolderO, Action: action.Jump("/admin/case_import")}})

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
	// 关联场景：仅当"是否自动化"选择为 是(1) / 部分是(2) 时才显示
	formList.AddField(biz.T("common.related_scene"), "scene", db.Varchar, form.Select).
		FieldOptions(sceneNames)
	formList.AddJS(template.JS(`
$(function () {
    function toggleRelatedScene() {
        var v = $('input[name="auto"]:checked').val();
        if (v === '1' || v === '2') {
            $("label[for='scene']").parent().show();
        } else {
            $("label[for='scene']").parent().hide();
        }
    }
    toggleRelatedScene();
    $('input[name="auto"]').on('ifChanged change', toggleRelatedScene);
});
`))
	formList.AddField(biz.T("case.func_devleoper"), "fun_developer", db.Varchar, form.Text)
	formList.AddField(biz.T("common.case_designer"), "case_designer", db.Varchar, form.Text)
	formList.AddField(biz.T("test_case.case_executor"), "case_executor", db.Varchar, form.Text)
	formList.AddField(biz.T("common.test_result"), "test_result", db.Varchar, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "pass", Text: biz.T("common.pass")},
			{Value: "fail", Text: biz.T("common.fail")},
			{Value: "untest", Text: biz.T("test_case.result_untest")},
			{Value: "deprecated", Text: biz.T("common.status_discarded")},
			{Value: "unmerged", Text: biz.T("test_case.result_unmerged")},
		}).FieldDefault("untest")
	formList.AddField(biz.T("test_case.test_time"), "test_time", db.Varchar, form.Datetime)
	formList.AddField(biz.T("test_case.result_history"), "result_history", db.Longtext, form.Custom).
		FieldCustomContent(template.HTML("{{ .Value }}")).
		FieldDisplay(resultHistoryDisplay).
		FieldDisableWhenCreate()

	formList.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text)
	formList.AddField(biz.T("common.product_line"), "product", db.Varchar, form.Select).
		FieldOptions(products)
	formList.AddField(biz.T("test_case.ext_info"), "ext_info", db.Longtext, form.TextArea).
		FieldHelpMsg(template2.HTML(biz.T("test_case.help_ext_info")))
	formList.AddField(biz.T("test_case.test_process"), "test_process", db.Longtext, form.RichText).
		FieldEnableFileUpload(config.Url("/operation/test_case_process_image_upload"), biz.TestCaseImageUploadHandler).
		FieldOptionExtJS(template.JS(fmt.Sprintf(`
	test_processeditor.customConfig.uploadImgServer = '%s';
	test_processeditor.customConfig.uploadImgMaxSize = 3 * 1024 * 1024;
	test_processeditor.customConfig.uploadImgMaxLength = 10;
	test_processeditor.customConfig.uploadFileName = 'file';
	test_processeditor.customConfig.uploadImgParams = {
	    case_number: $('input[name="case_number"]').val() || '',
	    module: $('input[name="module"]').val() || ''
	};
	`, config.Url("/operation/test_case_process_image_upload")))).
		FieldHelpMsg(template2.HTML(biz.T("test_case.help_test_process")))
	formList.AddField(biz.T("common.remark"), "remark", db.Varchar, form.TextArea)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("test_case").SetTitle(biz.T("test_case.title")).SetDescription(biz.T("test_case.description"))

	// 保存前自动填充：操作测试结果 → 更新测试时间 + 执行者；操作测试时间 → 更新执行者
	formList.SetPreProcessFn(func(values form2.Values) form2.Values {
		// 同步测试过程图片：删除已移除图片 + 剩余图片按顺序重命名
		if html := values.Get("test_process"); html != "" {
			if newHTML := biz.SyncTestCaseProcessImages(html, values.Get("case_number"), values.Get("module")); newHTML != html {
				values.Add("test_process", newHTML)
			}
		}

		id := values.Get("id")
		if id == "" || id == "0" {
			return values
		}
		updates := biz.AutoFillCaseExecutorTime(id, values.Get("test_result"), values.Get("test_time"), userName, values.IsSingleUpdatePost())
		for k, v := range updates {
			values.Add(k, v)
		}
		return values
	})

	detail := testCase.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.case_number"), "case_number", db.Varchar)
	detail.AddField(biz.T("common.case_title"), "case_name", db.Varchar).
		FieldDisplay(caseFieldDisplay(dataLocale, "case_name", false))
	detail.AddField(biz.T("common.case_module"), "module", db.Varchar).
		FieldDisplay(caseFieldDisplay(dataLocale, "module", false))
	detail.AddField(biz.T("common.case_type"), "case_type", db.Varchar).
		FieldDisplay(caseFieldDisplay(dataLocale, "case_type", false))
	detail.AddField(biz.T("common.case_level"), "priority", db.Varchar)
	detail.AddField(biz.T("test_case.precondition"), "pre_condition", db.Longtext).
		FieldDisplay(caseFieldDisplay(dataLocale, "pre_condition", true))
	detail.AddField(biz.T("common.test_range"), "test_range", db.Longtext).
		FieldDisplay(caseFieldDisplay(dataLocale, "test_range", false))
	detail.AddField(biz.T("common.test_step"), "test_steps", db.Longtext).
		FieldDisplay(caseFieldDisplay(dataLocale, "test_steps", true))
	detail.AddField(biz.T("common.expected_result"), "expect_result", db.Longtext).
		FieldDisplay(caseFieldDisplay(dataLocale, "expect_result", true))
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
			} else if model.Value == "fail" {
				return biz.T("common.fail")
			} else if model.Value == "untest" {
				return biz.T("test_case.result_untest")
			} else if model.Value == "deprecated" {
				return biz.T("common.status_discarded")
			} else if model.Value == "unmerged" {
				return biz.T("test_case.result_unmerged")
			}
			return biz.T("common.yes")
		})
	detail.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar)
	detail.AddField(biz.T("common.related_scene"), "scene", db.Varchar)
	detail.AddField(biz.T("common.product_line"), "product", db.Varchar)
	detail.AddField(biz.T("common.remark"), "remark", db.Varchar)
	detail.AddField(biz.T("test_case.test_process"), "test_process", db.Longtext)
	detail.AddField(biz.T("test_case.ext_info"), "ext_info", db.Longtext)
	detail.AddField(biz.T("test_case.result_history"), "result_history", db.Longtext).
		FieldDisplay(resultHistoryDisplay)

	detail.SetTable("test_case").SetTitle(biz.T("test_case.title")).SetDescription(biz.T("test_case.description"))

	return testCase
}

// rowString 将 model.Row 中的值安全转 string（兼容 string / []byte / nil）
func rowString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", t)
	}
}

// caseFieldDisplay 生成按数据语种取词的字段显示函数；html=true 时对 \n 换行
func caseFieldDisplay(dataLocale, field string, html bool) func(types.FieldModel) interface{} {
	return func(model types.FieldModel) interface{} {
		v := biz.GetCaseLocalized(rowString(model.Row["case_number"]), rowString(model.Row["module"]), dataLocale, field)
		if v == "" {
			v = model.Value
		}
		if html {
			return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(v), "\n", "<br/>"))
		}
		return v
	}
}

// resultHistoryDisplay 渲染结果历史：HTML 表格，结果枚举按当前语言翻译
func resultHistoryDisplay(model types.FieldModel) interface{} {
	raw := rowString(model.Value)
	if raw == "" {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(`<table class="table table-bordered" style="margin-bottom:0;width:100%;">`)
	sb.WriteString(`<thead><tr>`)
	sb.WriteString(`<th>` + template.HTMLEscapeString(biz.T("test_case.history_seq")) + `</th>`)
	sb.WriteString(`<th>` + template.HTMLEscapeString(biz.T("test_case.history_result")) + `</th>`)
	sb.WriteString(`<th>` + template.HTMLEscapeString(biz.T("test_case.history_executor")) + `</th>`)
	sb.WriteString(`<th>` + template.HTMLEscapeString(biz.T("test_case.history_time")) + `</th>`)
	sb.WriteString(`</tr></thead><tbody>`)
	for _, line := range strings.Split(raw, "\n") {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		sb.WriteString(`<tr>`)
		sb.WriteString(`<td>` + template.HTMLEscapeString(parts[0]) + `</td>`)
		sb.WriteString(`<td>` + template.HTMLEscapeString(testResultLabel(parts[1])) + `</td>`)
		sb.WriteString(`<td>` + template.HTMLEscapeString(parts[2]) + `</td>`)
		sb.WriteString(`<td>` + template.HTMLEscapeString(parts[3]) + `</td>`)
		sb.WriteString(`</tr>`)
	}
	sb.WriteString(`</tbody></table>`)
	return template.HTML(sb.String())
}

// testResultLabel 结果枚举 → 当前语言文案
func testResultLabel(v string) string {
	switch v {
	case "pass":
		return biz.T("common.pass")
	case "fail":
		return biz.T("common.fail")
	case "untest":
		return biz.T("test_case.result_untest")
	case "deprecated":
		return biz.T("common.status_discarded")
	case "unmerged":
		return biz.T("test_case.result_unmerged")
	}
	return v
}
