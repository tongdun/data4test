package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"html/template"
	"strings"
)

func GetApiFuzzingDataTable(ctx *context.Context) table.Table {

	apiFuzzingData := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiFuzzingData.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.data_desc"), "data_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(200)
	info.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(200)
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).FieldWidth(200).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(300)
	info.AddField(biz.T("common.api_module"), "api_module", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)
	info.AddField(biz.T("api_test_detail.header"), "header", db.Longtext)
	info.AddField(biz.T("common.url"), "url_query", db.Longtext).FieldWidth(250)
	info.AddField(biz.T("api_test_detail.body"), "body", db.Longtext).FieldWidth(300)
	info.AddField(biz.T("api_fuzzing_data.run_num"), "run_num", db.Int).FieldEditAble().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.expected_result"), "expected_result", db.Varchar)
	info.AddField(biz.T("api_fuzzing_data.actual_result"), "actual_result", db.Varchar)
	info.AddField(biz.T("common.test_result"), "result", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext)
	info.AddField(biz.T("common.response_data"), "response", db.Longtext)
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()
	info.AddButton(template.HTML(biz.T("custom.btn_create_scene_data")), icon.Android, action.Ajax("create_batch_scene_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			if err := biz.CreateSceneData(idStr, "fuzzing", ""); err == nil {
				status = biz.T("custom.msg_generated")
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", idStr, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("custom.btn_create_scene_data")), action.Ajax("create_scene_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneData(id, "fuzzing", ""); err == nil {
				status = biz.T("custom.msg_generated")
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("custom.btn_create_json_scene_data")), icon.Android, action.Ajax("create_batch_scene_fuzzing_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			if err := biz.CreateSceneData(idStr, "fuzzing", "json"); err == nil {
				status = biz.T("custom.msg_generated")
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", idStr, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("custom.btn_create_json_scene_data")), action.Ajax("create_scene_fuzzing_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneData(id, "fuzzing", "json"); err == nil {
				status = biz.T("custom.msg_generated")
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("common.btn_run")), icon.Android, action.Ajax("fuzzingdata_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
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
				if err := biz.RunFuzzData(id); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_run")), action.Ajax("fuzzingdata_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunFuzzData(id); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
			}
			return true, status, ""
		}))

	apps := biz.GetApps()
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))

	info.AddSelectBox(biz.T("common.test_result"), types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))

	info.SetTable("api_fuzzing_data").SetTitle(biz.T("api_fuzzing_data.title")).SetDescription(biz.T("api_fuzzing_data.description"))

	formList := apiFuzzingData.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.data_desc"), "data_desc", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_id"), "api_id", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_module"), "api_module", db.Varchar, form.Text)
	formList.AddField(biz.T("api_test_detail.header"), "header", db.Longtext, form.Text)
	formList.AddField(biz.T("common.url"), "url_query", db.Longtext, form.Text)
	formList.AddField(biz.T("api_test_detail.body"), "body", db.Longtext, form.Text)
	formList.AddField(biz.T("api_fuzzing_data.run_num"), "run_num", db.Int, form.Number)
	formList.AddField(biz.T("common.expected_result"), "expected_result", db.Varchar, form.Text)
	formList.AddField(biz.T("api_fuzzing_data.actual_result"), "actual_result", db.Varchar, form.Text)
	formList.AddField(biz.T("common.test_result"), "result", db.Varchar, form.Text)
	formList.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext, form.Text)
	formList.AddField(biz.T("common.response_data"), "response", db.Longtext, form.Text)
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_fuzzing_data").SetTitle(biz.T("api_fuzzing_data.title")).SetDescription(biz.T("api_fuzzing_data.description"))

	return apiFuzzingData
}
