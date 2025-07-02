package tables

import (
	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"strings"
)

func GetAiTemplateTable(ctx *context.Context) table.Table {

	aiTemplate := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiTemplate.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name
	aiPlatform := biz.GetAiPlatform()

	info.AddField("自增主键", "id", db.Int)
	info.AddField("模板名称", "template_name", db.Varchar)
	pTypes := types.FieldOptions{
		{Value: "1", Text: "用例"},
		{Value: "2", Text: "数据"},
		{Value: "3", Text: "场景"},
		{Value: "4", Text: "任务"},
		{Value: "5", Text: "Issue"},
		{Value: "6", Text: "报告"},
		{Value: "7", Text: "分析"},
	}
	info.AddField("模板类型", "template_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "用例"
			} else if model.Value == "2" {
				return "数据"
			} else if model.Value == "3" {
				return "场景"
			} else if model.Value == "4" {
				return "任务"
			} else if model.Value == "5" {
				return "Issue"
			} else if model.Value == "6" {
				return "报告"
			} else if model.Value == "7" {
				return "分析"
			}
			return "用例"
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(pTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(pTypes).
		FieldWidth(80)

	info.AddField("模板内容", "template_content", db.Text).FieldHide()
	info.AddField("生效状态", "use_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "apply" {
				return "启用"
			} else if model.Value == "edit" {
				return "编辑"
			}
			return "编辑"
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "apply", Text: "启用"},
		{Value: "edit", Text: "编辑"},
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "apply", Text: "启用"},
		{Value: "edit", Text: "编辑"},
	})

	info.AddField("适用平台", "applicable_platform", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(aiPlatform).
		FieldEditAble(editType.Select).
		FieldEditOptions(aiPlatform)
	info.AddField("创建人", "create_user", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	//info.AddField("修改人", "modify_user", db.Varchar).
	//	FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
	//	FieldTrimSpace().FieldWidth(80)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.SetTable("ai_template").SetTitle("智能模板").SetDescription("AiTemplate")
	appendConversionHelp := template.HTML("选填，按需追加对话轮次")
	formList := aiTemplate.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("模板名称", "template_name", db.Varchar, form.Text).
		FieldMust()
	formList.AddField("模板类型", "template_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "用例"},
			{Value: "2", Text: "数据"},
			{Value: "3", Text: "场景"},
			{Value: "4", Text: "任务"},
			{Value: "5", Text: "Issue"},
			{Value: "6", Text: "报告"},
			{Value: "7", Text: "分析"},
		})
	formList.AddField("模板内容", "template_content", db.Text, form.TextArea).
		FieldMust()
	formList.AddField("追加对话", "append_conversion", db.Text, form.Array).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Split(model.Value, "|**|")
		}).FieldHelpMsg(appendConversionHelp)
	formList.AddField("适用平台", "applicable_platform", db.Varchar, form.Radio).
		FieldOptions(aiPlatform)
	formList.AddField("生效状态", "use_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "apply", Text: "启用"},
			{Value: "edit", Text: "编辑"},
		}).FieldDefault("edit")

	formList.AddField("创建人", "create_user", db.Varchar, form.Text).
		FieldDefault(userName).
		FieldHideWhenCreate().
		FieldDisplayButCanNotEditWhenUpdate()
	//formList.AddField("修改人", "modify_user", db.Varchar, form.Text).
	//	FieldDefault(userName).
	//	FieldHideWhenCreate().
	//	FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("ai_template").SetTitle("智能模板").SetDescription("智能模板")

	formList.SetPostHook(func(values form2.Values) (err error) {
		appendContentList := values["append_conversion[values][]"]
		id := values["id"][0]
		err = biz.UpdateAiTemplateAppendContend(appendContentList, id)
		return
	})
	return aiTemplate
}
