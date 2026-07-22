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

func GetPlaybookTable(ctx *context.Context) table.Table {

	playbook := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	products := biz.GetProducts()
	partProducts := biz.GetProductsByUpdateTime(1)
	//info := playbook.GetInfo().HideFilterArea()
	info := playbook.GetInfo()
	user := auth.Auth(ctx)
	userName := user.Name

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldWidth(150).
		FieldFilterable()
	info.AddField(biz.T("dashboard.task_id"), "task_id", db.Varchar).
		FieldWidth(150)
	info.AddField(biz.T("common.name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().
		FieldWidth(220).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetPlaybookUsedInTaskList(model.Value, model.ID)
		})
	info.AddField(biz.T("common.data_file_list"), "data_file_list", db.Longtext).
		FieldWidth(600).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetDataFileLinkByDataStr(model.Value)
		})

	info.AddField(biz.T("common.last_file"), "last_file", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return template.Default().
				Link().
				SetURL("/admin/fm/data/preview?path=/" + value.Value).
				SetContent(template2.HTML(value.Value)).
				OpenInNewTab().
				SetTabTitle(template2.HTML(biz.T("common.data_file"))).
				GetContent()
		}).FieldWidth(160)

	pTypes := types.FieldOptions{
		{Value: "1", Text: biz.T("scene_test_history.scene_type_1")},
		{Value: "2", Text: biz.T("scene_test_history.scene_type_2")},
		{Value: "3", Text: biz.T("scene_test_history.scene_type_3")},
		{Value: "4", Text: biz.T("scene_test_history.scene_type_4")},
		{Value: "5", Text: biz.T("scene_test_history.scene_type_5")},
	}

	info.AddField(biz.T("common.scene_type"), "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("scene_test_history.scene_type_1")
			} else if model.Value == "2" {
				return biz.T("scene_test_history.scene_type_2")
			} else if model.Value == "3" {
				return biz.T("scene_test_history.scene_type_3")
			} else if model.Value == "4" {
				return biz.T("scene_test_history.scene_type_4")
			} else if model.Value == "5" {
				return biz.T("scene_test_history.scene_type_5")
			}
			return biz.T("scene_test_history.scene_type_1")
		}).FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(pTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(pTypes).
		FieldWidth(80)

	info.AddField(biz.T("common.priority"), "priority", db.Int).
		FieldFilterable(types.FilterType{FormType: form.Number}).
		FieldSortable().FieldWidth(80).
		FieldEditAble(editType.Text)
	info.AddField(biz.T("common.run_time"), "run_time", db.Int).
		FieldFilterable(types.FilterType{FormType: form.Number}).
		FieldSortable().
		FieldEditAble(editType.Text)

	info.AddField(biz.T("common.test_result"), "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}).FieldWidth(80)
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext).FieldWidth(160).FieldHide()
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().
		FieldWidth(80).
		FieldHide()
	info.AddField(biz.T("common.product"), "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(products).
		FieldEditAble(editType.Select).
		FieldEditOptions(partProducts).
		FieldWidth(220)

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

	info.AddButton(template2.HTML(biz.T("common.btn_copy")), icon.Android, action.Ajax("playbook_batch_copy",
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
				if err := biz.CopyPlaybook(id, userNameSub); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_copy")), action.Ajax("playbook_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopyPlaybook(id, userNameSub); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("common.btn_create_task")), icon.Android, action.Ajax("playbook_batch_task_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			user := auth.Auth(ctx)
			userNameSub := user.Name
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			err := biz.AutoCreateSchedule(idStr, userNameSub, "scene")
			if err != nil {
				status = biz.T("error.create_task_fail")
				return false, status, fmt.Sprintf("%s", err)
			}
			status = biz.T("schedule.task_added")
			return true, status, ""

		}))

	info.AddButton(template2.HTML(biz.T("common.btn_run")), icon.Android, action.Ajax("playbook_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			user := auth.Auth(ctx)
			userNameSub := user.Name
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids, err := biz.GetPriority(idStr)
			if err != nil {
				status = fmt.Sprintf("%s", err)
				return false, status, ""
			}

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}

				if err := biz.RunPlaybookFromMgmt(id, "start", "", "playbook", userNameSub, ""); err == nil {
					status = biz.T("scene_test_history.title") + biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_run")), action.Ajax("playbook_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			user := auth.Auth(ctx)
			userNameSub := user.Name
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunPlaybookFromMgmt(id, "start", "", "playbook", userNameSub, ""); err == nil {
				status = biz.T("scene_test_history.title") + biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("common.btn_continue")), icon.Android, action.Ajax("playbook_batch_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			user := auth.Auth(ctx)
			userNameSub := user.Name
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids, err := biz.GetPriority(idStr)
			if err != nil {
				status = fmt.Sprintf("%s", err)
				return false, status, ""
			}

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.RunPlaybookFromMgmt(id, "continue", "", "playbook", userNameSub, ""); err == nil {
					status = biz.T("scene_test_history.title") + biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_continue")), action.Ajax("playbook_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.RunPlaybookFromMgmt(id, "continue", "", "playbook", userNameSub, ""); err == nil {
				status = biz.T("scene_test_history.title") + biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddSelectBox(biz.T("common.product"), products, action.FieldFilter("product"))

	info.AddSelectBox(biz.T("common.test_result"), types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))
	info.SetTable("playbook").SetTitle(biz.T("playbook.title")).SetDescription(biz.T("playbook.description"))

	files := biz.GetFilesFromMySQL()

	formList := playbook.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text)

	sceneTypeMsg := template2.HTML(biz.T("scene_test_history.scene_type_help"))
	dataHelp := template2.HTML(biz.T("common.help_data_required"))

	formList.AddField(biz.T("common.edit_type"), "edit_type", db.Enum, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "input", Text: biz.T("playbook.edit_type_input")},
		}).
		FieldOnChooseHide("input", "select_table").
		FieldOnChooseShow("input", "input_list").
		FieldDefault("input").
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetPlaybookEditTypeById(model.ID)
		})

	formList.AddField(biz.T("playbook.data_input"), "input_list", db.Varchar, form.TextArea).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetPlaybookApiStr(model.ID)
		}).FieldHelpMsg(dataHelp)

	formList.AddField(biz.T("common.last_file"), "last_file", db.Varchar, form.SelectSingle).
		FieldOptions(files)
	formList.AddField(biz.T("common.scene_type"), "scene_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("scene_test_history.scene_type_1")},
			{Value: "2", Text: biz.T("scene_test_history.scene_type_2")},
			{Value: "3", Text: biz.T("scene_test_history.scene_type_3")},
			{Value: "4", Text: biz.T("scene_test_history.scene_type_4")},
			{Value: "5", Text: biz.T("scene_test_history.scene_type_5")},
		}).FieldDefault("1").FieldHelpMsg(sceneTypeMsg)
	formList.AddField(biz.T("common.priority"), "priority", db.Int, form.Number).FieldDefault("1")
	formList.AddField(biz.T("common.run_time"), "run_time", db.Int, form.Number).FieldDefault("1")
	formList.AddField(biz.T("common.test_result"), "result", db.Varchar, form.Text)
	formList.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.product"), "product", db.Varchar, form.Select).
		FieldOptions(products)
	formList.AddField(biz.T("common.user_name"), "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("playbook").SetTitle(biz.T("playbook.title")).SetDescription(biz.T("playbook.description"))

	formList.SetPostHook(func(values form2.Values) (err error) {
		if _, ok := values["edit_type"]; !ok {
			return
		}

		pEditType := values["edit_type"][0]
		inputInfo := values["input_list"][0]
		var apiList []string
		if pEditType == "select" {
			apiList = values["select_list"]
		} else {
			apiList = strings.Split(inputInfo, "\r\n")
		}
		id := values["id"][0]
		numList := values["data_number"]
		err = biz.UpdatePlaybookApiList(id, apiList, numList)
		return
	})

	detail := playbook.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.name"), "name", db.Varchar)
	detail.AddField(biz.T("common.data_file_list"), "data_file_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetDataDetailLinkByDataStr(model.Value)
		})
	detail.AddField(biz.T("common.last_file"), "last_file", db.Varchar)
	detail.AddField(biz.T("common.scene_type"), "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("scene_test_history.scene_type_1")
			}
			if model.Value == "2" {
				return biz.T("scene_test_history.scene_type_2")
			}
			if model.Value == "3" {
				return biz.T("scene_test_history.scene_type_3")
			}
			if model.Value == "4" {
				return biz.T("scene_test_history.scene_type_4")
			}
			if model.Value == "5" {
				return biz.T("scene_test_history.scene_type_5")
			}
			return biz.T("scene_test_history.scene_type_1")
		})
	detail.AddField(biz.T("common.priority"), "priority", db.Int)
	detail.AddField(biz.T("common.run_time"), "run_time", db.Int)
	detail.AddField(biz.T("common.test_result"), "result", db.Varchar)
	detail.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext)
	detail.AddField(biz.T("common.remark"), "remark", db.Longtext)
	detail.AddField(biz.T("common.product"), "product", db.Varchar)
	detail.AddField(biz.T("common.user_name"), "user_name", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).FieldHide()

	detail.SetTable("playbook").SetTitle(biz.T("common.detail_title")).SetDescription(biz.T("common.detail_description"))

	return playbook
}
