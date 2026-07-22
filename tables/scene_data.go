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

func GetSceneDataTable(ctx *context.Context) table.Table {

	sceneData := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	apps := biz.GetApps()
	info := sceneData.GetInfo()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	user := auth.Auth(ctx)
	userName := user.Name

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldFilterable()
	info.AddField(biz.T("common.name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetDataUsedInPlaybookList(model.Value, model.ID)
		})
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(apps)
	info.AddField(biz.T("common.file_name"), "file_name", db.Longtext).FieldDisplay(func(value types.FieldModel) interface{} {
		return template.Default().
			Link().
			SetURL("/admin/fm/data/preview?path=/" + value.Value).
			SetContent(template2.HTML(value.Value)).
			OpenInNewTab().
			SetTabTitle(template2.HTML(biz.T("common.data_file"))).
			GetContent()
	}).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.file_type"), "file_type", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("scene_data.file_type_1")
			} else if model.Value == "2" {
				return biz.T("scene_data.file_type_2")
			} else if model.Value == "3" {
				return biz.T("scene_data.file_type_3")
			} else if model.Value == "4" {
				return biz.T("scene_data.file_type_4")
			} else if model.Value == "5" {
				return biz.T("scene_data.file_type_5")
			} else if model.Value == "99" {
				return biz.T("scene_data.file_type_99")
			}
			return biz.T("scene_data.file_type_1")
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("scene_data.file_type_1")},
		{Value: "2", Text: biz.T("scene_data.file_type_2")},
		{Value: "3", Text: biz.T("scene_data.file_type_3")},
		{Value: "4", Text: biz.T("scene_data.file_type_4")},
		{Value: "5", Text: biz.T("scene_data.file_type_5")},
		{Value: "99", Text: biz.T("scene_data.file_type_99")},
	})
	info.AddField(biz.T("common.data_content"), "content", db.Longtext).FieldHide()
	info.AddField(biz.T("common.run_time"), "run_time", db.Int).
		FieldFilterable(types.FilterType{FormType: form.Number}).FieldSortable().
		FieldEditAble(editType.Text)
	info.AddField(biz.T("common.test_result"), "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	})
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext).
		FieldWidth(120)
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldHide()
	info.AddField(biz.T("common.user_name"), "user_name", db.Varchar).
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

	info.AddButton(template2.HTML(biz.T("common.btn_copy")), icon.Android, action.Ajax("scenedata_batch_copy",
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
				if err := biz.CopySceneData(id, userNameSub); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))
	info.AddActionButton(template2.HTML(biz.T("common.btn_copy")), action.Ajax("data_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopySceneData(id, userNameSub); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err)
			}
			return true, status, ""
		}))
	info.AddButton(template2.HTML(biz.T("common.btn_create_task")), icon.Android, action.Ajax("data_batch_task_create_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			err := biz.AutoCreateSchedule(idStr, userNameSub, "data")
			if err != nil {
				status = biz.T("error.create_playbook_fail")
				return false, status, fmt.Sprintf("%s", err)
			}

			status = biz.T("schedule.task_added")
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("common.btn_create_scene")), icon.Android, action.Ajax("data_batch_playbook_create_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			err := biz.AutoCreatePlaybook(idStr, userNameSub)
			if err != nil {
				status = biz.T("error.create_playbook_fail")
				return false, status, fmt.Sprintf("%s", err)
			}

			status = biz.T("common.operate_success")
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("common.btn_run")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/data_batch_run",
		Title:  biz.T("scene_data.test_title"),
		Width:  "900px",
		Height: "300px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		products := biz.GetProducts()
		ids := ctx.FormValue("ids")
		panel.AddField(biz.T("common.id"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
		panel.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).
			FieldDefault(products[0].Value).
			FieldHelpMsg(template.HTML(biz.T("scene_data.env_help")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/data_batch_run"))

	info.AddActionButton(template2.HTML(biz.T("common.btn_run")), action.Ajax("scenedata_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			var status string
			id := ctx.FormValue("id")
			if err := biz.RepeatRunDataFile(userName, id, "", "data", ""); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))
	info.AddSelectBox(biz.T("common.test_result"), types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))
	info.SetTable("scene_data").SetTitle(biz.T("common.data_file")).SetDescription(biz.T("common.data_file"))

	fileNameHelp := template2.HTML(biz.T("scene_data.help_file_name"))
	fileTypeMsg := template2.HTML(biz.T("scene_data.help_file_type"))

	formList := sceneData.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_id"), "api_id", db.Varchar, form.Text)
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.SelectSingle).
		FieldOptions(apps)
	formList.AddField(biz.T("common.file_name"), "file_name", db.Longtext, form.Url).
		FieldHelpMsg(fileNameHelp)
	formList.AddField(biz.T("common.file_type"), "file_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("scene_data.file_type_1")},
			{Value: "2", Text: biz.T("scene_data.file_type_2")},
			{Value: "3", Text: biz.T("scene_data.file_type_3")},
			{Value: "4", Text: biz.T("scene_data.file_type_4")},
			{Value: "5", Text: biz.T("scene_data.file_type_5")},
			{Value: "99", Text: biz.T("scene_data.file_type_99")},
		}).FieldDefault("1").FieldHelpMsg(fileTypeMsg)
	formList.AddField(biz.T("common.data_content"), "content", db.Longtext, form.TextArea).FieldDefault(biz.T("data_content.sample"))
	formList.AddField(biz.T("common.run_time"), "run_time", db.Int, form.Number).
		FieldDefault("1")
	formList.AddField(biz.T("common.test_result"), "result", db.Varchar, form.Text)
	formList.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.user_name"), "user_name", db.Varchar, form.Text).
		FieldDefault(userName).
		FieldDisplayButCanNotEditWhenCreate().
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()
	formList.SetTable("scene_data").SetTitle(biz.T("common.data_file")).SetDescription(biz.T("common.data_file"))

	formList.SetPostHook(func(values form2.Values) (err error) {
		content := values["content"][0]
		fileName := values["file_name"][0]
		id := values["id"][0]
		err = biz.BakOldVer(id, content, fileName)
		return
	})

	detail := sceneData.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.name"), "name", db.Varchar)
	detail.AddField(biz.T("common.api_id"), "api_id", db.Varchar)
	detail.AddField(biz.T("common.app"), "app", db.Varchar)
	detail.AddField(biz.T("common.file_name"), "file_name", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			linkStr := fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", model.Value, model.Value)
			return linkStr
		})
	detail.AddField(biz.T("common.file_type"), "file_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("scene_data.file_type_1")
			}
			if model.Value == "2" {
				return biz.T("scene_data.file_type_2")
			}
			if model.Value == "3" {
				return biz.T("scene_data.file_type_3")
			}
			if model.Value == "4" {
				return biz.T("scene_data.file_type_4")
			}
			if model.Value == "5" {
				return biz.T("scene_data.file_type_5")
			}
			if model.Value == "99" {
				return biz.T("scene_data.file_type_99")
			}
			return biz.T("scene_data.file_type_1")
		})
	detail.AddField(biz.T("common.run_time"), "run_time", db.Int)
	detail.AddField(biz.T("common.test_result"), "result", db.Varchar)
	detail.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext)
	detail.AddField(biz.T("common.remark"), "remark", db.Longtext)
	detail.AddField(biz.T("common.user_name"), "user_name", db.Varchar)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp)

	detail.SetTable("scene_data").SetTitle(biz.T("common.detail_title")).SetDescription(biz.T("common.detail_description"))

	return sceneData
}
