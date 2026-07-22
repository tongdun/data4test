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
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"html/template"
)

func GetApiDefinitionTable(ctx *context.Context) table.Table {

	apiDefinition := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	user := auth.Auth(ctx)
	userName := user.Name

	info := apiDefinition.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldWidth(150).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiAutoDataList(model.Value, model.ID)
		})
	info.AddField(biz.T("common.api_module"), "api_module", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)
	info.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(150)
	info.AddField(biz.T("common.http_method"), "http_method", db.Enum).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldEditAble().FieldWidth(120)
	info.AddField(biz.T("common.path"), "path", db.Varchar).FieldWidth(120).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.header_params"), "header", db.JSON).FieldWidth(300).
		FieldHide()
	info.AddField(biz.T("api_definition.path_vars"), "path_variable", db.JSON).FieldWidth(150).
		FieldHide()
	info.AddField(biz.T("common.query_params"), "query_parameter", db.JSON).FieldWidth(300).
		FieldHide()
	info.AddField(biz.T("common.body_params"), "body", db.JSON).FieldWidth(300).
		FieldHide()
	info.AddField(biz.T("api_definition.resp_vars"), "response", db.JSON).FieldWidth(300).
		FieldHide()
	info.AddField(biz.T("api_definition.version"), "version", db.Int).FieldWidth(80)
	info.AddField(biz.T("api_definition.api_status"), "api_status", db.Enum).FieldWidth(100).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("api_definition.api_status_1")},
		{Value: "2", Text: biz.T("api_definition.api_status_2")},
		{Value: "3", Text: biz.T("api_definition.api_status_3")},
		{Value: "4", Text: biz.T("api_definition.api_status_4")},
		{Text: biz.T("api_definition.api_status_30"), Value: "30"},
		{Text: biz.T("api_definition.api_status_31"), Value: "31"},
		{Text: biz.T("api_definition.api_status_32"), Value: "32"},
		{Text: biz.T("api_definition.api_status_33"), Value: "33"},
		{Text: biz.T("api_definition.api_status_34"), Value: "34"},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return biz.T("api_definition.api_status_1")
		} else if model.Value == "2" {
			return biz.T("api_definition.api_status_2")
		} else if model.Value == "3" {
			return biz.T("api_definition.api_status_3")
		} else if model.Value == "4" {
			return biz.T("api_definition.api_status_4")
		} else if model.Value == "30" {
			return biz.T("api_definition.api_status_30")
		} else if model.Value == "31" {
			return biz.T("api_definition.api_status_31")
		} else if model.Value == "32" {
			return biz.T("api_definition.api_status_32")
		} else if model.Value == "33" {
			return biz.T("api_definition.api_status_33")
		} else if model.Value == "34" {
			return biz.T("api_definition.api_status_34")
		}
		return biz.T("api_definition.api_status_1")
	})
	info.AddField(biz.T("api_definition.change_content"), "change_content", db.JSON).
		FieldHide()
	info.AddField(biz.T("api_definition.check"), "check", db.Varchar).FieldWidth(120)
	info.AddField(biz.T("api_definition.api_check_fail_reason"), "api_check_fail_reason", db.JSON).
		FieldHide()
	info.AddField(biz.T("api_definition.is_need_auto"), "is_need_auto", db.Enum).
		FieldEditAble(editType.Select).FieldEditOptions(types.FieldOptions{
		{Value: "-1", Text: biz.T("common.no")},
		{Value: "1", Text: biz.T("common.yes")},
	}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "-1", Text: biz.T("common.no")},
		{Value: "1", Text: biz.T("common.yes")},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return biz.T("common.yes")
		} else if model.Value == "-1" {
			return biz.T("common.no")
		}
		return biz.T("common.yes")
	})
	info.AddField(biz.T("common.is_auto"), "is_auto", db.Enum).
		FieldEditAble(editType.Select).FieldEditOptions(types.FieldOptions{
		{Value: "-1", Text: biz.T("common.no")},
		{Value: "1", Text: biz.T("common.yes")},
	}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "-1", Text: biz.T("common.no")},
		{Value: "1", Text: biz.T("common.yes")},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return biz.T("common.yes")
		} else if model.Value == "-1" {
			return biz.T("common.no")
		}
		return biz.T("common.no")
	})

	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)

	info.AddField(biz.T("common.remark"), "remark", db.JSON).
		FieldHide()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddActionButton(template.HTML(biz.T("common.btn_create_data")), action.Ajax("create_scene_raw_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneDataFromRaw(id, ""); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("api_definition.create_case")), action.Ajax("create_scene_raw_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneDataFromRaw(id, "json"); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_create_scene")), action.Ajax("create_test_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateTestData(id); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_create_data")), action.Ajax("create_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateFuzzingData(id); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_create_scene")), action.Ajax("test_data_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTestData(id); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_create_data")), action.Ajax("fuzzing_data_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTestData(id); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.operate_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_create_scene")), action.Ajax("create_playbook_case",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreatePlaybookByAPIId(id, userName); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("api_definition.create_case")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_create_case_by_api_define",
		Title:  biz.T("api_definition.create_case"),
		Width:  "900px",
		Height: "540px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		products := biz.GetProducts()
		aiCaseTemplates := biz.GetAiTemplateOptions("1")
		aiPlatforms := biz.GetAiCreatePlatform()
		panel.AddField(biz.T("common.id"), "ids", db.Varchar, form.Text).
			FieldDefault(ids).
			FieldHide()
		panel.AddField(biz.T("common.name"), "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiCaseTemplates).FieldDefault(aiCaseTemplates[0].Value)
		panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text)
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField(biz.T("common.create_platform"), "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.EnableAjax(biz.T("common.operate_success"))
		return panel
	}, "/ai_create_case_by_api_define"))

	info.AddButton(template.HTML(biz.T("common.btn_create_data")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_create_data_by_api_define",
		Title:  biz.T("common.btn_create_data"),
		Width:  "900px",
		Height: "540px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		products := biz.GetProducts()
		aiDataTemplates := biz.GetAiTemplateOptions("2")
		aiPlatforms := biz.GetAiCreatePlatform()
		panel.AddField(biz.T("common.id"), "ids", db.Varchar, form.Text).
			FieldDefault(ids).
			FieldHide()
		panel.AddField(biz.T("common.name"), "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiDataTemplates).
			FieldDefault(aiDataTemplates[0].Value)
		panel.AddField(biz.T("common.create_platform"), "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).
			FieldDefault(aiPlatforms[0].Value)
		panel.AddField(biz.T("common.intro_version"), "intro_version", db.Varchar, form.Text).
			FieldHelpMsg(template.HTML(biz.T("api_definition.intro_version_help")))
		panel.AddField(biz.T("common.product_list"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).
			FieldDefault(aiPlatforms[0].Value).
			FieldHelpMsg(template.HTML(biz.T("api_definition.product_help")))
		panel.EnableAjax(biz.T("common.operate_success"))
		return panel
	}, "/ai_create_data_by_api_define"))

	apps := biz.GetApps()
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))

	info.AddSelectBox(biz.T("common.http_method"), types.FieldOptions{
		{Value: "get", Text: "get"},
		{Value: "post", Text: "post"},
		{Value: "delete", Text: "delete"},
		{Value: "put", Text: "put"},
	}, action.FieldFilter("http_method"))

	info.AddSelectBox(biz.T("api_definition.check"), types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("check"))

	info.SetTable("api_definition").SetTitle(biz.T("api_definition.title")).SetDescription(biz.T("api_definition.description"))

	formList := apiDefinition.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.api_id"), "api_id", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_module"), "api_module", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar, form.Text)
	formList.AddField(biz.T("common.http_method"), "http_method", db.Enum, form.Text).
		FieldDefault("get")
	formList.AddField(biz.T("common.path"), "path", db.Varchar, form.Text)
	formList.AddField(biz.T("common.header_params"), "header", db.Longtext, form.TextArea)
	formList.AddField(biz.T("api_definition.path_vars"), "path_variable", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.query_params"), "query_parameter", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.body_params"), "body", db.Longtext, form.TextArea)
	formList.AddField(biz.T("api_definition.resp_vars"), "response", db.Longtext, form.TextArea)
	formList.AddField(biz.T("api_definition.version"), "version", db.Int, form.Number).
		FieldDefault("1").
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("api_definition.api_status"), "api_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("api_definition.api_status_1"), Value: "1"},
			{Text: biz.T("api_definition.api_status_2"), Value: "2"},
			{Text: biz.T("api_definition.api_status_3"), Value: "3"},
			{Text: biz.T("api_definition.api_status_4"), Value: "4"},
			{Text: biz.T("api_definition.api_status_30"), Value: "30"},
			{Text: biz.T("api_definition.api_status_31"), Value: "31"},
			{Text: biz.T("api_definition.api_status_32"), Value: "32"},
			{Text: biz.T("api_definition.api_status_33"), Value: "33"},
			{Text: biz.T("api_definition.api_status_34"), Value: "34"},
		}).FieldDefault("1")
	formList.AddField(biz.T("api_definition.change_content"), "change_content", db.Longtext, form.TextArea)
	formList.AddField(biz.T("api_definition.check"), "check", db.Varchar, form.Text)
	formList.AddField(biz.T("api_definition.api_check_fail_reason"), "api_check_fail_reason", db.Longtext, form.TextArea)
	formList.AddField(biz.T("api_definition.is_need_auto"), "is_need_auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.no"), Value: "-1"},
			{Text: biz.T("common.yes"), Value: "1"},
		}).FieldDefault("1")
	formList.AddField(biz.T("common.is_auto"), "is_auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.no"), Value: "-1"},
			{Text: biz.T("common.yes"), Value: "1"},
		}).FieldDefault("-1")
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.Text)
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_definition").SetTitle(biz.T("api_definition.title")).SetDescription(biz.T("api_definition.description"))
	formList.SetPostHook(func(values form2.Values) (err error) {
		id := values["id"][0]
		err = biz.UpdateApiDefVer(id)
		return
	})

	detail := apiDefinition.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.api_id"), "api_id", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiAutoDataList(model.Value, model.ID)
		})
	detail.AddField(biz.T("common.api_module"), "api_module", db.Varchar)
	detail.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar)
	detail.AddField(biz.T("common.http_method"), "http_method", db.Enum)
	detail.AddField(biz.T("common.path"), "path", db.Varchar)
	detail.AddField(biz.T("common.header_params"), "header", db.JSON)
	detail.AddField(biz.T("api_definition.path_vars"), "path_variable", db.JSON)
	detail.AddField(biz.T("common.query_params"), "query_parameter", db.JSON)
	detail.AddField(biz.T("common.body_params"), "body", db.JSON)
	detail.AddField(biz.T("api_definition.resp_vars"), "response", db.JSON)
	detail.AddField(biz.T("api_definition.version"), "version", db.Int)
	detail.AddField(biz.T("api_definition.api_status"), "api_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("api_definition.api_status_1")
			} else if model.Value == "2" {
				return biz.T("api_definition.api_status_2")
			} else if model.Value == "3" {
				return biz.T("api_definition.api_status_3")
			} else if model.Value == "4" {
				return biz.T("api_definition.api_status_4")
			} else if model.Value == "30" {
				return biz.T("api_definition.api_status_30")
			} else if model.Value == "31" {
				return biz.T("api_definition.api_status_31")
			} else if model.Value == "32" {
				return biz.T("api_definition.api_status_32")
			} else if model.Value == "33" {
				return biz.T("api_definition.api_status_33")
			} else if model.Value == "34" {
				return biz.T("api_definition.api_status_34")
			}
			return biz.T("api_definition.api_status_1")
		})
	detail.AddField(biz.T("api_definition.change_content"), "change_content", db.Longtext)
	detail.AddField(biz.T("api_definition.check"), "check", db.Varchar)
	detail.AddField(biz.T("api_definition.api_check_fail_reason"), "api_check_fail_reason", db.JSON)
	detail.AddField(biz.T("api_definition.is_need_auto"), "is_need_auto", db.Enum).
		FieldEditAble(editType.Select).FieldEditOptions(types.FieldOptions{
		{Value: "-1", Text: biz.T("common.no")},
		{Value: "1", Text: biz.T("common.yes")},
	}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "-1", Text: biz.T("common.no")},
		{Value: "1", Text: biz.T("common.yes")},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return biz.T("common.yes")
		} else if model.Value == "-1" {
			return biz.T("common.no")
		}
		return biz.T("common.yes")
	})
	detail.AddField(biz.T("common.is_auto"), "is_auto", db.Enum).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "-1", Text: biz.T("common.no")},
		{Value: "1", Text: biz.T("common.yes")},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return biz.T("common.yes")
		} else if model.Value == "-1" {
			return biz.T("common.no")
		}
		return biz.T("common.no")
	})

	detail.AddField(biz.T("common.app"), "app", db.Varchar)
	detail.AddField(biz.T("common.remark"), "remark", db.JSON)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	return apiDefinition
}
