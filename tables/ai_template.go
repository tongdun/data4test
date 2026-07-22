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
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	template2 "html/template"
	"strings"
)

func GetAiTemplateTable(ctx *context.Context) table.Table {

	aiTemplate := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiTemplate.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name
	aiPlatform := biz.GetAiPlatform()

	info.AddField(biz.T("common.id"), "id", db.Int)
	info.AddField(biz.T("common.name"), "template_name", db.Varchar)
	pTypes := types.FieldOptions{
		{Value: "1", Text: biz.T("ai_template.type_case")},
		{Value: "2", Text: biz.T("common.data")},
		{Value: "3", Text: biz.T("common.scene")},
		{Value: "4", Text: biz.T("type.task")},
		{Value: "5", Text: biz.T("ai_template.type_issue")},
		{Value: "6", Text: biz.T("ai_template.type_report")},
		{Value: "7", Text: biz.T("ai_template.type_analysis")},
	}
	info.AddField(biz.T("ai_template.category"), "template_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("ai_template.type_case")
			} else if model.Value == "2" {
				return biz.T("common.data")
			} else if model.Value == "3" {
				return biz.T("common.scene")
			} else if model.Value == "4" {
				return biz.T("type.task")
			} else if model.Value == "5" {
				return biz.T("ai_template.type_issue")
			} else if model.Value == "6" {
				return biz.T("ai_template.type_report")
			} else if model.Value == "7" {
				return biz.T("ai_template.type_analysis")
			}
			return biz.T("ai_template.type_case")
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(pTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(pTypes).
		FieldWidth(80)

	info.AddField(biz.T("common.data_content"), "template_content", db.Text).FieldHide()
	info.AddField(biz.T("common.use_status"), "use_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "apply" {
				return biz.T("ai_template.status_apply")
			} else if model.Value == "edit" {
				return biz.T("ai_template.status_edit")
			}
			return biz.T("ai_template.status_edit")
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "apply", Text: biz.T("ai_template.status_apply")},
		{Value: "edit", Text: biz.T("ai_template.status_edit")},
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "apply", Text: biz.T("ai_template.status_apply")},
		{Value: "edit", Text: biz.T("ai_template.status_edit")},
	})

	info.AddField(biz.T("ai_template.platform"), "applicable_platform", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(aiPlatform).
		FieldEditAble(editType.Select).
		FieldEditOptions(aiPlatform)
	info.AddField(biz.T("common.user_name"), "create_user", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton(template2.HTML(biz.T("common.btn_copy")), icon.Android, action.Ajax("ai_template_batch_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			ids := strings.Split(idStr, ",")

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.CopyAiTemplate(id, userNameSub); err == nil {
					status = biz.T("ai_template.msg_copy_success")
				} else {
					status = fmt.Sprintf(biz.T("common.copy_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_copy")), action.Ajax("ai_template_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopyAiTemplate(id, userNameSub); err == nil {
				status = biz.T("ai_template.msg_copy_success")
			} else {
				status = fmt.Sprintf(biz.T("common.copy_fail"), id, err)
			}
			return true, status, ""
		}))

	info.SetTable("ai_template").SetTitle(biz.T("common.ai_template")).SetDescription(biz.T("ai_template.description"))
	appendConversionHelp := template.HTML(biz.T("ai_template.help_append_conversion"))
	formList := aiTemplate.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "template_name", db.Varchar, form.Text).
		FieldMust()
	formList.AddField(biz.T("ai_template.category"), "template_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("ai_template.type_case")},
			{Value: "2", Text: biz.T("common.data")},
			{Value: "3", Text: biz.T("common.scene")},
			{Value: "4", Text: biz.T("type.task")},
			{Value: "5", Text: biz.T("ai_template.type_issue")},
			{Value: "6", Text: biz.T("ai_template.type_report")},
			{Value: "7", Text: biz.T("ai_template.type_analysis")},
		})
	formList.AddField(biz.T("common.data_content"), "template_content", db.Text, form.TextArea).
		FieldMust()
	formList.AddField(biz.T("ai_template.append_conversion"), "append_conversion", db.Text, form.Array).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Split(model.Value, "|**|")
		}).FieldHelpMsg(appendConversionHelp)
	formList.AddField(biz.T("ai_template.platform"), "applicable_platform", db.Varchar, form.Radio).
		FieldOptions(aiPlatform)
	formList.AddField(biz.T("common.use_status"), "use_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "apply", Text: biz.T("ai_template.status_apply")},
			{Value: "edit", Text: biz.T("ai_template.status_edit")},
		}).FieldDefault("edit")

	formList.AddField(biz.T("common.user_name"), "create_user", db.Varchar, form.Text).
		FieldDefault(userName).
		FieldHideWhenCreate().
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("ai_template").SetTitle(biz.T("common.ai_template")).SetDescription(biz.T("common.ai_template"))

	formList.SetPostHook(func(values form2.Values) (err error) {
		appendContentList := values["append_conversion[values][]"]
		id := values["id"][0]
		err = biz.UpdateAiTemplateAppendContend(appendContentList, id)
		return
	})
	return aiTemplate
}
