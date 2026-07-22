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

func GetApiTestDetailTable(ctx *context.Context) table.Table {

	apiTestDetail := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiTestDetail.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.api_id"), "api_id", db.Varchar).FieldWidth(150).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar).FieldWidth(250).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.data_desc"), "data_desc", db.Varchar).FieldWidth(250).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("api_test_detail.header"), "header", db.Longtext).FieldWidth(500).
		FieldHide()
	info.AddField("URL", "url", db.Varchar).FieldWidth(300)
	info.AddField(biz.T("api_test_detail.body"), "body", db.Longtext).FieldWidth(300)
	info.AddField(biz.T("common.response_data"), "response", db.Longtext).FieldWidth(300)
	info.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext)
	info.AddField(biz.T("api_test_detail.test_result"), "test_result", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(200)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton(template.HTML(biz.T("common.btn_again")), icon.Android, action.Ajax("apitest_batch_again",
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
				if err := biz.RunAgain(id); err == nil {
					status = biz.T("common.operate_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("common.btn_again")), action.Ajax("apitest_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunAgain(id); err == nil {
				status = biz.T("common.operate_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("error.exec_fail"), id, err)
			}
			return true, status, ""
		}))

	apps := biz.GetApps()
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))

	info.AddSelectBox(biz.T("api_test_detail.test_result"), types.FieldOptions{
		{Value: "success", Text: "success"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("test_result"))

	info.SetTable("api_test_detail").SetTitle(biz.T("api_test_detail.title")).SetDescription(biz.T("api_test_detail.description"))

	formList := apiTestDetail.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.api_id"), "api_id", db.Varchar, form.Text)
	formList.AddField(biz.T("common.api_desc"), "api_desc", db.Varchar, form.Text)
	formList.AddField(biz.T("common.data_desc"), "data_desc", db.Varchar, form.Text)
	formList.AddField(biz.T("api_test_detail.header"), "header", db.Longtext, form.Text)
	formList.AddField("URL", "url", db.Varchar, form.Text)
	formList.AddField(biz.T("api_test_detail.body"), "body", db.Longtext, form.Text)
	formList.AddField(biz.T("common.response_data"), "response", db.Longtext, form.Text)
	formList.AddField(biz.T("common.fail_reason"), "fail_reason", db.Longtext, form.Text)
	formList.AddField(biz.T("api_test_detail.test_result"), "test_result", db.Varchar, form.Text)
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.Text)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_test_detail").SetTitle(biz.T("api_test_detail.title")).SetDescription(biz.T("api_test_detail.description"))

	return apiTestDetail
}
