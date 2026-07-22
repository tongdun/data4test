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

func GetSceneDataTestHistoryTable(ctx *context.Context) table.Table {

	dataTestHistory := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := dataTestHistory.GetInfo()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	user := auth.Auth(ctx)
	userName := user.Name
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldFilterable().
		FieldTrimSpace()
	info.AddField(biz.T("dashboard.task_id"), "task_id", db.Varchar).
		FieldWidth(150)
	info.AddField(biz.T("common.name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.data_content"), "content", db.Longtext).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return template.Default().
				Link().
				SetURL("/admin/fm/history/preview?path=" + value.Value).
				SetContent(template2.HTML(value.Value)).
				OpenInNewTab().
				SetTabTitle(template2.HTML(biz.T("data_test_history.history_record"))).
				GetContent()
		})
	info.AddField(biz.T("common.test_result"), "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "pass", Text: "pass"},
			{Value: "fail", Text: "fail"},
		})
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext).
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
	})
	info.AddField(biz.T("common.product"), "product", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("common.user_name"), "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton(template2.HTML(biz.T("common.btn_again")), icon.Android, action.Ajax("historyData_batch_again",
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
					status = biz.T("common.operate_success")
					continue
				}
				if err := biz.HistoryDataRunAgain(id, userName); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_again")), action.Ajax("historyData_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			user := auth.Auth(ctx)
			userName := user.Name
			var status string
			if err := biz.HistoryDataRunAgain(id, userName); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
			}
			return true, status, ""
		}))

	apps := biz.GetApps()
	products := biz.GetProducts()
	info.AddSelectBox(biz.T("common.product"), products, action.FieldFilter("product"))
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))
	info.AddSelectBox(biz.T("common.test_result"), types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))

	info.SetTable("scene_data_test_history").SetTitle(biz.T("data_test_history.title")).SetDescription(biz.T("scene_data_test_history.description"))

	formList := dataTestHistory.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.name"), "name", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.api_id"), "api_id", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField(biz.T("common.data_content"), "content", db.Longtext, form.Text)
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
	formList.AddField(biz.T("common.product"), "product", db.Varchar, form.Text)
	formList.AddField(biz.T("common.user_name"), "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("scene_data_test_history").SetTitle(biz.T("data_test_history.title")).SetDescription(biz.T("scene_data_test_history.description"))

	formList.SetPostHook(func(values form2.Values) (err error) {
		fileName := values["content"][0]
		id := values["id"][0]
		err = biz.ModifyEditedData(id, fileName)
		return
	})

	return dataTestHistory
}
