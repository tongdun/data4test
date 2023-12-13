package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"data4perf/biz"
	"strings"
	"fmt"
)

func GetApiTestDetailTable(ctx *context.Context) table.Table {

	apiTestDetail := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiTestDetail.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("接口ID", "api_id", db.Varchar).FieldWidth(150).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("接口描述", "api_desc", db.Varchar).FieldWidth(250).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("数据描述", "data_desc", db.Varchar).FieldWidth(250).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Header", "header", db.Longtext).FieldWidth(500).
		FieldHide()
	info.AddField("URL", "url", db.Varchar).FieldWidth(300)
	info.AddField("Body", "body", db.Longtext).FieldWidth(300)
	info.AddField("Response", "response", db.Longtext).FieldWidth(300)
	info.AddField("失败原因", "fail_reason", db.Longtext)
	info.AddField("测试结果", "test_result", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("关联应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(200)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton("再来一次", icon.Android, action.Ajax("apitest_batch_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再执行"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "测试完成，请刷新列表查看"
					continue
				}
				if err := biz.RunAgain(id); err == nil {
					status = "测试完成，请刷新列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("再来一次", action.Ajax("apitest_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunAgain(id); err == nil {
				status = "测试完成，请刷新列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	apps := biz.GetApps()
	info.AddSelectBox("关联应用", apps, action.FieldFilter("app"))

	info.AddSelectBox("测试结果", types.FieldOptions{
		{Value: "success", Text: "success"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("test_result"))

	info.SetTable("api_test_detail").SetTitle("结果详情").SetDescription("结果详情")

	formList := apiTestDetail.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("接口ID", "api_id", db.Varchar, form.Text)
	formList.AddField("接口描述", "api_desc", db.Varchar, form.Text)
	formList.AddField("数据描述", "data_desc", db.Varchar, form.Text)
	formList.AddField("Header", "header", db.Longtext, form.Text)
	formList.AddField("URL", "url", db.Varchar, form.Text)
	formList.AddField("Body", "body", db.Longtext, form.Text)
	formList.AddField("Response", "response", db.Longtext, form.Text)
	formList.AddField("失败原因", "fail_reason", db.Longtext, form.Text)
	formList.AddField("测试结果", "test_result", db.Varchar, form.Text)
	formList.AddField("关联应用", "app", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_test_detail").SetTitle("结果详情").SetDescription("结果详情")

	return apiTestDetail
}
