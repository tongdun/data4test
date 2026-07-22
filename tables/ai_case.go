package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"html/template"
	"strings"
)

func GetAiCaseTable(ctx *context.Context) table.Table {

	aiCase := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	info := aiCase.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name

	info.SetFilterFormLayout(form.LayoutThreeCol)

	//aiPlatforms := biz.GetAiCreatePlatform()
	//aiTemplates := biz.GetAiTemplateOptions("1")
	products := biz.GetProducts() // 全局域
	caseTypes := biz.GetTestcaseType()

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldFilterable()
	info.AddField(biz.T("common.case_number"), "case_number", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("common.case_title"), "case_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("common.module"), "module", db.Varchar).
		FieldFilterable()
	info.AddField(biz.T("common.case_type"), "case_type", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(caseTypes)
	info.AddField(biz.T("common.case_level"), "priority", db.Varchar).
		FieldFilterable()
	info.AddField(biz.T("common.test_step"), "pre_condition", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTMLEscapeString(model.Value)
		})
	info.AddField(biz.T("common.test_range"), "test_range", db.Varchar)
	info.AddField(biz.T("common.test_step"), "test_steps", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTMLEscapeString(model.Value)
		})
	info.AddField(biz.T("common.expected_result"), "expect_result", db.Varchar)
	info.AddField(biz.T("common.is_auto"), "auto", db.Enum).
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
	info.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar).
		FieldFilterable()
	//info.AddField("用例版本", "case_version", db.Int)
	info.AddField(biz.T("common.product"), "product", db.Varchar).
		//FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldFilterOptions(products)
	info.AddField(biz.T("common.source"), "source", db.Varchar).
		FieldFilterable()
	info.AddField(biz.T("common.use_status"), "use_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("common.status_initial")
			} else if model.Value == "2" {
				return biz.T("common.status_in_use")
			} else if model.Value == "3" {
				return biz.T("common.status_discarded")
			}
			return biz.T("common.status_initial")
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("common.status_initial")},
		{Value: "2", Text: biz.T("common.status_in_use")},
		{Value: "3", Text: biz.T("common.status_discarded")},
	})
	info.AddField(biz.T("common.modify_status"), "modify_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_case.modify_initial")
			} else if model.Value == "2" {
				return biz.T("common.modify_manual")
			} else if model.Value == "3" {
				return biz.T("common.modify_auto")
			}
			return biz.T("ai_case.modify_initial")
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("ai_case.modify_initial")},
		{Value: "2", Text: biz.T("common.modify_manual")},
		{Value: "3", Text: biz.T("common.modify_auto")},
	})
	info.AddField(biz.T("common.user_name"), "create_user", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton(template.HTML(biz.T("common.btn_export_md")), icon.File, action.Ajax("ai_test_case_export_markdown",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			fileName, err := biz.ExportTestCase2Markdown(idStr, "ai_case")
			if err != nil {
				status = fmt.Sprintf(biz.T("error.export_fail"), err)
				return false, status, ""
			}
			hostIp := ctx.Request.Host
			downloadUrl := fmt.Sprintf("http://%s/admin/fm/case/download?path=/%s", hostIp, fileName)
			status = fmt.Sprintf(biz.T("common.operate_success"), downloadUrl)
			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("common.btn_export_xmind")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_case_export_xmind",
		Title:  biz.T("common.btn_export_xmind"),
		Width:  "900px",
		Height: "720px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		aiPlatforms := biz.GetAiCreatePlatform()
		products := biz.GetProducts()
		panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldDisplayButCanNotEditWhenCreate().
			FieldHelpMsg(template.HTML(biz.T("common.help_select_or_filter")))
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products)
		panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("common.help_version_multi")))
		panel.AddField(biz.T("common.create_platform"), "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms)
		panel.AddField(biz.T("common.user_name"), "create_user", db.Varchar, form.Text)
		panel.AddField(biz.T("common.module"), "module", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("common.help_module_multi")))
		panel.AddField(biz.T("common.created_at"), "created_at", db.Varchar, form.DatetimeRange)
		panel.AddField(biz.T("common.source"), "source", db.Varchar, form.Text).FieldDefault("ai_case").FieldHide()

		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_case_export_xmind"))

	info.AddButton(template.HTML(biz.T("common.btn_import")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_case_import",
		Title:  biz.T("common.btn_import"),
		Width:  "900px",
		Height: "680px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		aiPlatforms := biz.GetAiCreatePlatform()
		products := biz.GetProducts()
		panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text)
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField(biz.T("common.create_platform"), "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField(biz.T("common.conversation_id"), "conversation_id", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("common.help_conversation_or_raw")))
		panel.AddField(biz.T("common.raw_reply"), "raw_reply", db.Varchar, form.TextArea).
			FieldHelpMsg(template.HTML(biz.T("common.help_raw_reply")))

		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_case_import"))

	info.AddButton(template.HTML(biz.T("common.btn_ai_generate")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_case_create_by_create_desc",
		Title:  biz.T("common.btn_ai_generate"),
		Width:  "900px",
		Height: "680px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		aiPlatforms := biz.GetAiCreatePlatform()
		aiTemplates := biz.GetAiTemplateOptions("1")
		products := biz.GetProducts()
		panel.AddField(biz.T("common.ai_template"), "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiTemplates).FieldDefault(aiTemplates[0].Value)
		panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text)
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField(biz.T("common.create_platform"), "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField(biz.T("common.create_desc"), "create_desc", db.Varchar, form.TextArea)
		panel.AddField(biz.T("common.upload_file"), "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
			"maxFileCount": 10,
		}).FieldHelpMsg(template.HTML(biz.T("ai_case.help_upload_file")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/ai_case_create_by_create_desc"))

	info.AddButton(template.HTML(biz.T("common.btn_ai_optimize")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_case_optimize",
		Title:  biz.T("common.btn_ai_optimize"),
		Width:  "900px",
		Height: "480px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		aiPlatforms := biz.GetAiCreatePlatform()
		panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
		panel.AddField(biz.T("common.optimize_platform"), "optimize_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField(biz.T("common.optimize_desc"), "optimize_desc", db.Varchar, form.TextArea)
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/ai_case_optimize"))

	info.AddButton(template.HTML(biz.T("common.btn_use")), icon.Android, action.Ajax("ai_case_batch_use",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if idStr == "," {
				status = biz.T("common.msg_select_first_for_use")
				return false, status, ""
			}
			ids := strings.Trim(idStr, ",")
			if err := biz.UseAiCase(ids, userNameSub); err == nil {
				status = biz.T("common.msg_use_success")
			} else {
				status = fmt.Sprintf(biz.T("common.msg_use_failed"), ids, err)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_use")), action.Ajax("ai_case_use",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.UseAiCase(id, userNameSub); err == nil {
				status = biz.T("common.msg_use_success")
			} else {
				status = fmt.Sprintf(biz.T("common.msg_use_failed"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("ai_case.btn_xmind_import_and_use")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/testcase_xmind_import_and_use",
		Title:  biz.T("ai_case.btn_xmind_import_and_use"),
		Width:  "900px",
		Height: "680px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		info.AddField(biz.T("common.product"), "product", db.Varchar).
			FieldFilterable(types.FilterType{FormType: form.Select}).
			FieldFilterOptions(products)
		info.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar).
			FieldFilterable()
		panel.AddField(biz.T("ai_case.upload_xmind"), "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
			"maxFileCount": 1,
		}).FieldHelpMsg(template.HTML(biz.T("common.help_xmind")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/testcase_xmind_import_and_use"))

	info.SetTable("ai_case").SetTitle(biz.T("ai_case.title")).SetDescription(biz.T("ai_case.description"))

	formList := aiCase.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.case_number"), "case_number", db.Varchar, form.Text)
	formList.AddField(biz.T("common.case_title"), "case_name", db.Varchar, form.Text)
	formList.AddField(biz.T("common.module"), "module", db.Varchar, form.Text)
	formList.AddField(biz.T("common.case_type"), "case_type", db.Varchar, form.Select).
		FieldOptions(caseTypes)
	formList.AddField(biz.T("common.case_level"), "priority", db.Varchar, form.Text)
	formList.AddField(biz.T("common.test_step"), "pre_condition", db.Varchar, form.TextArea)
	formList.AddField(biz.T("common.test_range"), "test_range", db.Varchar, form.TextArea)
	formList.AddField(biz.T("common.test_step"), "test_steps", db.Varchar, form.TextArea)
	formList.AddField(biz.T("common.expected_result"), "expect_result", db.Varchar, form.TextArea)
	formList.AddField(biz.T("common.is_auto"), "auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "0", Text: biz.T("common.no")},
			{Value: "1", Text: biz.T("common.yes")},
			{Value: "2", Text: biz.T("common.part")},
		}).FieldDefault("0")
	formList.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text)
	formList.AddField(biz.T("common.source"), "source", db.Varchar, form.Text).
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.use_status"), "use_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("common.status_initial")},
			{Value: "2", Text: biz.T("common.status_in_use")},
			{Value: "3", Text: biz.T("common.status_discarded")},
		}).
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenUpdate().
		FieldDefault("1")
	formList.AddField(biz.T("common.modify_status"), "modify_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("ai_case.modify_initial")},
			{Value: "2", Text: biz.T("common.modify_manual")},
			{Value: "3", Text: biz.T("common.modify_auto")},
		}).
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenUpdate().
		FieldDefault("1")
	formList.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
		FieldOptions(products)
	formList.AddField(biz.T("common.user_name"), "create_user", db.Varchar, form.Text).
		FieldDefault(userName).FieldHide()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("ai_case").SetTitle(biz.T("ai_case.title")).SetDescription(biz.T("ai_case.description"))

	formList.SetPostHook(func(values form2.Values) (err error) {
		id := values["id"][0]
		err = biz.UpdateAiCaseStatus(id)
		return
	})

	detail := aiCase.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.case_number"), "case_number", db.Varchar)
	detail.AddField(biz.T("common.case_title"), "case_name", db.Varchar)
	detail.AddField(biz.T("common.module"), "module", db.Varchar)
	detail.AddField(biz.T("common.case_type"), "case_type", db.Varchar)
	detail.AddField(biz.T("common.case_level"), "priority", db.Varchar)
	detail.AddField(biz.T("common.test_step"), "pre_condition", db.Text)
	detail.AddField(biz.T("common.test_range"), "test_range", db.Text)
	detail.AddField(biz.T("common.test_step"), "test_steps", db.Text)
	detail.AddField(biz.T("common.expected_result"), "expect_result", db.Text)
	detail.AddField(biz.T("common.is_auto"), "auto", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return biz.T("common.no")
			}
			if model.Value == "1" {
				return biz.T("common.yes")
			}
			if model.Value == "2" {
				return biz.T("common.part")
			}
			return biz.T("common.yes")
		})
	detail.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar)
	detail.AddField(biz.T("common.source"), "source", db.Varchar)
	detail.AddField(biz.T("common.use_status"), "use_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("common.status_initial")
			}
			if model.Value == "2" {
				return biz.T("common.status_in_use")
			}
			if model.Value == "3" {
				return biz.T("common.status_discarded")
			}
			return biz.T("common.status_initial")
		})
	detail.AddField(biz.T("common.modify_status"), "modify_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_case.modify_initial")
			}
			if model.Value == "2" {
				return biz.T("common.modify_manual")
			}
			if model.Value == "3" {
				return biz.T("common.modify_auto")
			}
			return biz.T("ai_case.modify_initial")
		})
	detail.AddField(biz.T("common.product"), "product", db.Varchar)
	detail.AddField(biz.T("common.user_name"), "create_user", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).FieldHide()

	detail.SetTable("ai_case").SetTitle(biz.T("ai_case.title")).SetDescription(biz.T("ai_case.description"))
	return aiCase
}
