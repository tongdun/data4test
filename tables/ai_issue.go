package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	template2 "html/template"
)

func GetAiIssueTable(ctx *context.Context) table.Table {

	aiIssue := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiIssue.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name
	products := biz.GetProducts()
	partProducts := biz.GetProductsByUpdateTime(1)
	aiPlatforms := biz.GetAiCreatePlatform()
	issueSourceTypes := types.FieldOptions{
		{Value: "1", Text: biz.T("ai_issue.source_data")},
		{Value: "2", Text: biz.T("ai_issue.source_scene")},
		{Value: "3", Text: biz.T("ai_issue.source_manual")},
	}
	confirmStatusTypes := types.FieldOptions{
		{Value: "1", Text: biz.T("ai_issue.confirm_bug")},
		{Value: "2", Text: biz.T("ai_issue.confirm_optimize")},
		{Value: "3", Text: biz.T("ai_issue.confirm_false")},
	}
	solveStatusTypes := types.FieldOptions{
		{Value: "1", Text: biz.T("ai_issue.status_created")},
		{Value: "2", Text: biz.T("ai_issue.status_solving")},
		{Value: "3", Text: biz.T("ai_issue.status_fixed")},
		{Value: "4", Text: biz.T("ai_issue.status_verified")},
		{Value: "5", Text: biz.T("ai_issue.status_ignored")},
	}
	againTestTypes := types.FieldOptions{
		{Value: "0", Text: biz.T("common.msg_test_failed")},
		{Value: "1", Text: biz.T("ai_issue.test_success")},
		{Value: "2", Text: biz.T("ai_issue.test_unknown")},
	}

	info.AddField(biz.T("common.id"), "id", db.Int)
	info.AddField(biz.T("ai_issue.issue_name"), "issue_name", db.Varchar)
	info.AddField(biz.T("ai_issue.issue_level"), "issue_level", db.Varchar)
	info.AddField(biz.T("ai_issue.issue_type"), "issue_source", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_issue.source_data")
			} else if model.Value == "2" {
				return biz.T("ai_issue.source_scene")
			} else if model.Value == "3" {
				return biz.T("ai_issue.source_manual")
			}
			return biz.T("ai_issue.source_data")
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(issueSourceTypes)

	info.AddField(biz.T("ai_issue.source_name"), "source_name", db.Varchar).
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
	info.AddField(biz.T("common.analysis_platform"), "source", db.Varchar)
	info.AddField(biz.T("ai_issue.request_data"), "request_data", db.Text).
		FieldHide()
	info.AddField(biz.T("common.response_data"), "response_data", db.Text).
		FieldHide()
	info.AddField(biz.T("ai_issue.issue_detail"), "issue_detail", db.Text).
		FieldHide()
	info.AddField(biz.T("ai_issue.confirm_status"), "confirm_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_issue.confirm_bug")
			} else if model.Value == "2" {
				return biz.T("ai_issue.confirm_optimize")
			} else if model.Value == "3" {
				return biz.T("ai_issue.confirm_false")
			}
			return biz.T("ai_issue.confirm_bug")
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(confirmStatusTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(confirmStatusTypes)

	info.AddField(biz.T("ai_issue.root_cause"), "root_cause", db.Text).
		FieldHide()
	info.AddField(biz.T("ai_issue.impact_scope"), "impact_scope_analysis", db.Text).
		FieldHide()
	info.AddField(biz.T("ai_issue.impact_playbook"), "impact_playbook", db.Text).
		FieldHide()
	info.AddField(biz.T("ai_issue.impact_data"), "impact_data", db.Text).
		FieldHide()
	info.AddField(biz.T("ai_issue.solve_status"), "resolution_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_issue.status_created")
			} else if model.Value == "2" {
				return biz.T("ai_issue.status_solving")
			} else if model.Value == "3" {
				return biz.T("ai_issue.status_fixed")
			} else if model.Value == "4" {
				return biz.T("ai_issue.status_verified")
			} else if model.Value == "5" {
				return biz.T("ai_issue.status_ignored")
			}
			return biz.T("ai_issue.status_created")
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(solveStatusTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(solveStatusTypes)

	info.AddField(biz.T("ai_issue.again_test_result"), "again_test_result", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return biz.T("common.msg_test_failed")
			} else if model.Value == "1" {
				return biz.T("ai_issue.test_success")
			} else if model.Value == "2" {
				return biz.T("ai_issue.test_unknown")
			}
			return biz.T("ai_issue.test_unknown")
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(againTestTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(againTestTypes)

	info.AddField(biz.T("ai_issue.impact_test_result"), "impact_test_result", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return biz.T("common.msg_test_failed")
			} else if model.Value == "1" {
				return biz.T("ai_issue.test_success")
			} else if model.Value == "2" {
				return biz.T("ai_issue.test_unknown")
			}
			return biz.T("ai_issue.test_unknown")
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(againTestTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(againTestTypes)

	info.AddField(biz.T("common.product_list"), "product_list", db.Text).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(products).
		FieldEditAble(editType.Select).
		FieldEditOptions(partProducts)
	info.AddField(biz.T("common.user_name"), "create_user", db.Varchar)
	info.AddField(biz.T("common.modify_user"), "modify_user", db.Varchar)
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
		Id:     "/ai_issue_import",
		Title:  biz.T("common.btn_import"),
		Width:  "900px",
		Height: "680px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		products := biz.GetProducts()
		panel.AddField(biz.T("common.create_platform"), "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField(biz.T("common.conversation_id"), "conversation_id", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("common.help_conversation_or_raw")))
		panel.AddField(biz.T("common.raw_reply"), "raw_reply", db.Varchar, form.TextArea).
			FieldHelpMsg(template.HTML(biz.T("common.help_raw_reply")))

		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_issue_import"))

	info.AddButton(template.HTML(biz.T("ai_issue.btn_regression_test")), icon.Android, action.Ajax("source_again_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			user := auth.Auth(ctx)
			userName := user.Name
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			if err := biz.SourceAgainTest(userName, idStr); err == nil {
				status = biz.T("common.msg_test_completed")
			} else {
				status = fmt.Sprintf(biz.T("common.msg_test_failed"), err)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.SetTable("ai_issue").SetTitle(biz.T("ai_issue.title")).SetDescription(biz.T("ai_issue.description"))

	formList := aiIssue.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("ai_issue.issue_name"), "issue_name", db.Varchar, form.Text)
	formList.AddField(biz.T("ai_issue.issue_level"), "issue_level", db.Varchar, form.Text)
	formList.AddField(biz.T("ai_issue.issue_type"), "issue_source", db.Enum, form.Radio).
		FieldOptions(solveStatusTypes).
		FieldDefault("3")
	formList.AddField(biz.T("ai_issue.source_name"), "source_name", db.Varchar, form.Text)
	formList.AddField(biz.T("common.analysis_platform"), "source", db.Varchar, form.SelectSingle).
		FieldOptions(aiPlatforms)
	formList.AddField(biz.T("ai_issue.request_data"), "request_data", db.Text, form.TextArea)
	formList.AddField(biz.T("common.response_data"), "response_data", db.Text, form.TextArea)
	formList.AddField(biz.T("ai_issue.issue_detail"), "issue_detail", db.Text, form.TextArea)
	formList.AddField(biz.T("ai_issue.confirm_status"), "confirm_status", db.Enum, form.Radio).
		FieldOptions(confirmStatusTypes).
		FieldDefault("1")
	formList.AddField(biz.T("ai_issue.root_cause"), "root_cause", db.Text, form.TextArea)
	formList.AddField(biz.T("ai_issue.impact_scope"), "impact_scope_analysis", db.Text, form.TextArea)
	formList.AddField(biz.T("ai_issue.impact_playbook"), "impact_playbook", db.Text, form.TextArea)
	formList.AddField(biz.T("ai_issue.impact_data"), "impact_data", db.Text, form.TextArea)
	formList.AddField(biz.T("ai_issue.solve_status"), "resolution_status", db.Enum, form.Radio).
		FieldOptions(solveStatusTypes).
		FieldDefault("1")
	formList.AddField(biz.T("ai_issue.again_test_result"), "again_test_result", db.Enum, form.Radio).
		FieldOptions(againTestTypes).
		FieldDefault("2")
	formList.AddField(biz.T("ai_issue.impact_test_result"), "impact_test_result", db.Enum, form.Radio).
		FieldOptions(againTestTypes).
		FieldDefault("2")
	formList.AddField(biz.T("common.product_list"), "product_list", db.Text, form.SelectSingle).
		FieldOptions(products)
	formList.AddField(biz.T("common.user_name"), "create_user", db.Varchar, form.Text).
		FieldDefault(userName).
		FieldHideWhenCreate().
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.modify_user"), "modify_user", db.Varchar, form.Text).
		FieldDefault(userName).
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().
		FieldNowWhenInsert().
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().
		FieldNowWhenUpdate().
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()

	formList.SetTable("ai_issue").SetTitle(biz.T("ai_issue.title")).SetDescription(biz.T("ai_issue.description"))

	return aiIssue
}
