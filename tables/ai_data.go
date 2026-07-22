package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	template2 "html/template"
	"strings"
)

func GetAiDataTable(ctx *context.Context) table.Table {

	aiData := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	user := auth.Auth(ctx)
	userName := user.Name
	info := aiData.GetInfo().HideFilterArea()
	info.SetFilterFormLayout(form.LayoutThreeCol)

	apps := biz.GetApps()

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldFilterable()
	info.AddField(biz.T("common.name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetAiDataUsedInPlaybookList(model.Value, model.ID)
		})
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(apps)
	info.AddField(biz.T("common.source"), "source", db.Varchar).
		FieldFilterable()
	info.AddField(biz.T("common.file_name"), "file_name", db.Text).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return template.Default().
				Link().
				SetURL("/admin/fm/ai_data/preview?path=/" + value.Value).
				SetContent(template2.HTML(value.Value)).
				OpenInNewTab().
				SetTabTitle(template.HTML(biz.T("common.data_file"))).
				GetContent()
		}).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.file_type"), "file_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_data.type_standard")
			} else if model.Value == "2" {
				return biz.T("ai_data.type_python")
			} else if model.Value == "3" {
				return biz.T("ai_data.type_shell")
			} else if model.Value == "4" {
				return biz.T("ai_data.type_bat")
			} else if model.Value == "5" {
				return biz.T("ai_data.type_jmeter")
			} else if model.Value == "99" {
				return biz.T("ai_data.type_other")
			}
			return biz.T("ai_data.type_standard")
		}).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("ai_data.type_standard")},
			{Value: "2", Text: biz.T("ai_data.type_python")},
			{Value: "3", Text: biz.T("ai_data.type_shell")},
			{Value: "4", Text: biz.T("ai_data.type_bat")},
			{Value: "5", Text: biz.T("ai_data.type_jmeter")},
			{Value: "99", Text: biz.T("ai_data.type_other")},
		})
	info.AddField(biz.T("common.data_content"), "content", db.Text).
		FieldHide()
	info.AddField(biz.T("common.test_result"), "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "pass", Text: biz.T("common.pass")},
			{Value: "fail", Text: biz.T("common.fail")},
		})
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Text)
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
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("common.status_initial")},
			{Value: "2", Text: biz.T("common.status_in_use")},
			{Value: "3", Text: biz.T("common.status_discarded")},
		})
	info.AddField(biz.T("common.modify_status"), "modify_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_data.modify_initial")
			} else if model.Value == "2" {
				return biz.T("common.modify_manual")
			} else if model.Value == "3" {
				return biz.T("common.modify_auto")
			}
			return biz.T("ai_data.modify_initial")
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("ai_data.modify_initial")},
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

	info.AddButton(template.HTML(biz.T("common.btn_import")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_data_import",
		Title:  biz.T("common.btn_import"),
		Width:  "900px",
		Height: "680px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		aiPlatforms := biz.GetAiCreatePlatform()
		products := biz.GetProducts()
		panel.AddField(biz.T("common.create_platform"), "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("common.help_version_suffix")))
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField(biz.T("common.conversation_id"), "conversation_id", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("common.help_conversation_or_raw")))
		panel.AddField(biz.T("common.raw_reply"), "raw_reply", db.Varchar, form.TextArea).
			FieldHelpMsg(template.HTML(biz.T("common.help_raw_reply")))

		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_data_import"))

	info.AddButton(template.HTML(biz.T("ai_data.btn_analysis")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_data_test_and_analysis",
		Title:  biz.T("ai_data.btn_analysis"),
		Width:  "900px",
		Height: "680px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		aiPlatforms := biz.GetAiCreatePlatform()
		products := biz.GetProducts()
		aiAnalysisTemplates := biz.GetAiTemplateOptions("7")
		panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
		panel.AddField(biz.T("common.ai_template"), "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiAnalysisTemplates).FieldDefault(aiAnalysisTemplates[0].Value)
		panel.AddField(biz.T("common.analysis_platform"), "analysis_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).
			FieldDefault(aiPlatforms[0].Value)
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).
			FieldDefault(products[0].Value).
			FieldHelpMsg(template.HTML(biz.T("common.help_execute_data")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_data_test_and_analysis"))

	info.AddButton(template.HTML(biz.T("common.btn_test")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_data_test",
		Title:  biz.T("common.btn_test"),
		Width:  "900px",
		Height: "300px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		products := biz.GetProducts()
		panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).
			FieldDefault(products[0].Value).
			FieldHelpMsg(template.HTML(biz.T("ai_data.help_select_env")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_data_test"))

	info.AddButton(template.HTML(biz.T("common.btn_use")), icon.Android, action.Ajax("ai_data_batch_use",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Trim(idStr, ",")
			if err := biz.UseAiData(ids, userNameSub); err == nil {
				status = biz.T("common.msg_use_success")
			} else {
				status = fmt.Sprintf(biz.T("common.msg_use_failed"), ids, err)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_use")), action.Ajax("ai_data_use",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.UseAiData(id, userNameSub); err == nil {
				status = biz.T("common.msg_use_success")
			} else {
				status = fmt.Sprintf(biz.T("common.msg_use_failed"), id, err)
			}
			return true, status, ""
		}))

	info.SetTable("ai_data").SetTitle(biz.T("ai_data.title")).SetDescription(biz.T("ai_data.description"))
	fileNameHelp := template.HTML(biz.T("scene_data.help_file_name"))
	fileTypeMsg := template.HTML(biz.T("scene_data.help_file_type"))

	formList := aiData.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_id"), "api_id", db.Varchar, form.Text)
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.SelectSingle).
		FieldOptions(apps)
	formList.AddField(biz.T("common.source"), "source", db.Varchar, form.Text).
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.file_name"), "file_name", db.Varchar, form.Url).
		FieldHelpMsg(fileNameHelp)
	formList.AddField(biz.T("common.file_type"), "file_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("ai_data.type_standard")},
			{Value: "2", Text: biz.T("ai_data.type_python")},
			{Value: "3", Text: biz.T("ai_data.type_shell")},
			{Value: "4", Text: biz.T("ai_data.type_bat")},
			{Value: "5", Text: biz.T("ai_data.type_jmeter")},
			{Value: "99", Text: biz.T("ai_data.type_other")},
		}).FieldDefault("1").FieldHelpMsg(fileTypeMsg)
	formList.AddField(biz.T("common.data_content"), "content", db.Text, form.TextArea)
	formList.AddField(biz.T("common.test_result"), "result", db.Varchar, form.Text)
	formList.AddField(biz.T("common.fail_reason"), "fail_reason", db.Text, form.TextArea)
	formList.AddField(biz.T("common.use_status"), "use_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("common.status_initial")},
			{Value: "2", Text: biz.T("common.status_in_use")},
			{Value: "3", Text: biz.T("common.status_discarded")},
		}).
		FieldDefault("1").
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.modify_status"), "modify_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("ai_data.modify_initial")},
			{Value: "2", Text: biz.T("common.modify_manual")},
			{Value: "3", Text: biz.T("common.modify_auto")},
		}).
		FieldDefault("1").
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.user_name"), "create_user", db.Varchar, form.Text).
		FieldDefault(userName).FieldHide()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("ai_data").SetTitle(biz.T("ai_data.title")).SetDescription(biz.T("ai_data.description"))

	formList.SetPostHook(func(values form2.Values) (err error) {
		content := values["content"][0]
		fileName := values["file_name"][0]
		id := values["id"][0]
		err = biz.BakOldAiDataVer(id, content, fileName)
		err = biz.UpdateAiDataStatus(id)
		return
	})

	detail := aiData.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.name"), "name", db.Varchar)
	detail.AddField(biz.T("common.api_id"), "api_id", db.Varchar)
	detail.AddField(biz.T("common.app"), "app", db.Varchar)
	detail.AddField(biz.T("common.file_name"), "file_name", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			linkStr := fmt.Sprintf("<a href=\"/admin/fm/ai_data/preview?path=/%s\">%s</a>", model.Value, model.Value)
			return linkStr
		})
	detail.AddField(biz.T("common.source"), "source", db.Varchar)
	detail.AddField(biz.T("common.file_type"), "file_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_data.type_standard")
			}
			if model.Value == "2" {
				return biz.T("ai_data.type_python")
			}
			if model.Value == "3" {
				return biz.T("ai_data.type_shell")
			}
			if model.Value == "4" {
				return biz.T("ai_data.type_bat")
			}
			if model.Value == "5" {
				return biz.T("ai_data.type_jmeter")
			}
			if model.Value == "99" {
				return biz.T("ai_data.type_other")
			}
			return biz.T("ai_data.type_standard")
		})
	detail.AddField(biz.T("common.test_result"), "result", db.Varchar)
	detail.AddField(biz.T("common.fail_reason"), "fail_reason", db.Text)
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
				return biz.T("ai_data.modify_initial")
			}
			if model.Value == "2" {
				return biz.T("common.modify_manual")
			}
			if model.Value == "3" {
				return biz.T("common.modify_auto")
			}
			return biz.T("ai_data.modify_initial")
		})
	detail.AddField(biz.T("common.user_name"), "create_user", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).FieldHide()

	detail.SetTable("ai_data").SetTitle(biz.T("ai_data.title")).SetDescription(biz.T("ai_data.description"))

	return aiData
}
