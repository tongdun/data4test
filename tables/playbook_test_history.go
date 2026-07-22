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
	template2 "html/template"
	"strings"
)

func GetSceneTestHistoryTable(ctx *context.Context) table.Table {

	playbookTestHistory := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := playbookTestHistory.GetInfo()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	user := auth.Auth(ctx)
	userName := user.Name
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldFilterable().
		FieldTrimSpace().FieldWidth(60)
	info.AddField(biz.T("dashboard.task_id"), "task_id", db.Varchar)
	info.AddField(biz.T("common.name"), "name", db.Varchar).FieldWidth(160).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.data_file_list"), "data_file_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetHistoryDataLinkByDataStr(model.Value)
		})
	info.AddField(biz.T("common.last_file"), "last_file", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldDisplay(func(value types.FieldModel) interface{} {
			b, num := biz.IsStrEndWithTimeFormat(value.Value)
			suffix := biz.GetStrSuffix(value.Value)
			if b {
				dirName := value.Value[:len(value.Value)-num-len(suffix)]
				return template.Default().
					Link().
					SetURL("/admin/fm/history/preview?path=/" + dirName + "/" + value.Value).
					SetContent(template2.HTML(value.Value)).
					OpenInNewTab().
					SetTabTitle(template2.HTML(biz.T("test_execute_history_file.title"))).
					GetContent()
			} else {
				return template.Default().
					Link().
					SetURL("/admin/fm/data/preview?path=/" + value.Value).
					SetContent(template2.HTML(value.Value)).
					OpenInNewTab().
					SetTabTitle(template2.HTML(biz.T("common.data_file"))).
					GetContent()
			}
		}).FieldWidth(160)
	info.AddField(biz.T("common.scene_type"), "scene_type", db.Enum).
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
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("scene_test_history.scene_type_1")},
		{Value: "2", Text: biz.T("scene_test_history.scene_type_2")},
		{Value: "3", Text: biz.T("scene_test_history.scene_type_3")},
		{Value: "4", Text: biz.T("scene_test_history.scene_type_4")},
		{Value: "5", Text: biz.T("scene_test_history.scene_type_5")},
	}).FieldWidth(60)
	info.AddField(biz.T("common.test_result"), "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}).FieldWidth(70)
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext).
		FieldWidth(200).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.env_type_label"), "env_type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("common.env_type._1")
			} else if model.Value == "2" {
				return biz.T("common.env_type._2")
			} else if model.Value == "3" {
				return biz.T("common.env_type._3")
			} else if model.Value == "4" {
				return biz.T("common.env_type._4")
			} else if model.Value == "5" {
				return biz.T("common.env_type._5")
			}
			return ""
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("common.env_type._1")},
		{Value: "2", Text: biz.T("common.env_type._2")},
		{Value: "3", Text: biz.T("common.env_type._3")},
		{Value: "4", Text: biz.T("common.env_type._4")},
		{Value: "5", Text: biz.T("common.env_type._5")},
	}).FieldWidth(60)
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(120)
	info.AddField(biz.T("common.product"), "product", db.Varchar).
		FieldWidth(60).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.user_name"), "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(100).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton(template2.HTML(biz.T("common.btn_again")), icon.Android, action.Ajax("historyPlaybook_batch_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			user := auth.Auth(ctx)
			userName := user.Name
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = biz.T("common.operate_success")
					continue
				}
				if err := biz.RunHistoryPlaybook(id, "again", userName); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_again")), action.Ajax("historyPlaybook_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			user := auth.Auth(ctx)
			userName := user.Name
			var status string
			if err := biz.RunHistoryPlaybook(id, "again", userName); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("common.btn_continue")), icon.Android, action.Ajax("historyPlaybook_batch_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			user := auth.Auth(ctx)
			userName := user.Name
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			ids := strings.Split(idStr, ",")

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.RunHistoryPlaybook(id, "continue", userName); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_continue")), action.Ajax("historyPlaybook_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userName := user.Name
			if err := biz.RunHistoryPlaybook(id, "continue", userName); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	products := biz.GetProducts()
	info.AddSelectBox(biz.T("common.product"), products, action.FieldFilter("product"))

	info.AddSelectBox(biz.T("common.test_result"), types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))

	info.SetTable("scene_test_history").SetTitle(biz.T("playbook_test_history.title")).SetDescription(biz.T("scene_test_history.description"))

	sceneTypeMsg := template2.HTML(biz.T("scene_test_history.scene_type_help"))
	formList := playbookTestHistory.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.data_file_list"), "data_file_list", db.Longtext, form.RichText).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.last_file"), "last_file", db.Varchar, form.Text)
	formList.AddField(biz.T("common.scene_type"), "scene_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: biz.T("scene_test_history.scene_type_1")},
			{Value: "2", Text: biz.T("scene_test_history.scene_type_2")},
			{Value: "3", Text: biz.T("scene_test_history.scene_type_3")},
			{Value: "4", Text: biz.T("scene_test_history.scene_type_4")},
			{Value: "5", Text: biz.T("scene_test_history.scene_type_5")},
		}).FieldDefault("1").FieldHelpMsg(sceneTypeMsg)
	formList.AddField(biz.T("common.test_result"), "result", db.Varchar, form.Text)
	formList.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.env_type_label"), "env_type", db.Int, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.env_type._1"), Value: "1"},
			{Text: biz.T("common.env_type._2"), Value: "2"},
			{Text: biz.T("common.env_type._3"), Value: "3"},
			{Text: biz.T("common.env_type._4"), Value: "4"},
			{Text: biz.T("common.env_type._5"), Value: "5"},
		}).FieldDefault("2")
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.product"), "product", db.Varchar, form.Text)
	formList.AddField(biz.T("common.user_name"), "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldNowWhenUpdate().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("scene_test_history").SetTitle(biz.T("playbook_test_history.title")).SetDescription(biz.T("scene_test_history.description"))

	detail := playbookTestHistory.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("common.name"), "name", db.Varchar)
	detail.AddField(biz.T("common.data_file_list"), "data_file_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Replace(model.Value, ",", ",<br>", -1)
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
	detail.AddField(biz.T("common.test_result"), "result", db.Varchar)
	detail.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext)
	detail.AddField(biz.T("common.env_type_label"), "env_type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("common.env_type._1")
			} else if model.Value == "2" {
				return biz.T("common.env_type._2")
			} else if model.Value == "3" {
				return biz.T("common.env_type._3")
			} else if model.Value == "4" {
				return biz.T("common.env_type._4")
			} else if model.Value == "5" {
				return biz.T("common.env_type._5")
			}
			return ""
		})
	detail.AddField(biz.T("common.remark"), "remark", db.Longtext)
	detail.AddField(biz.T("common.product"), "product", db.Varchar)
	detail.AddField(biz.T("common.user_name"), "user_name", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp)

	detail.SetTable("scene_test_history").SetTitle(biz.T("playbook_test_history.title")).SetDescription(biz.T("scene_test_history.description"))

	return playbookTestHistory
}
