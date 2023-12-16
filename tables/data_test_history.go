package tables

import (
	"data4perf/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
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

	senceDataTestHistory := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := senceDataTestHistory.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField("唯一标识", "id", db.Int).
		FieldFilterable().
		FieldTrimSpace()
	info.AddField("数据描述", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("接口ID", "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("所属应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("数据文件", "content", db.Varchar).FieldDisplay(func(value types.FieldModel) interface{} {
		return template.Default().
			Link().
			SetURL("/admin/fm/history/preview?path=" + value.Value).
			SetContent(template2.HTML(value.Value)).
			OpenInNewTab().
			SetTabTitle(template.HTML("历史记录")).
			GetContent()
	}).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("测试结果", "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	})
	info.AddField("失败原因", "fail_reason", db.Longtext)
	info.AddField("环境类型", "env_type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "开发"
			} else if model.Value == "2" {
				return "测试"
			} else if model.Value == "3" {
				return "预发"
			} else if model.Value == "4" {
				return "演示"
			} else if model.Value == "5" {
				return "生产"
			}
			return ""
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "开发"},
		{Value: "2", Text: "测试"},
		{Value: "3", Text: "预发"},
		{Value: "4", Text: "演示"},
		{Value: "5", Text: "生产"},
	})
	info.AddField("所属产品", "product", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("备注", "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton("再来一次", icon.Android, action.Ajax("historyData_batch_again",
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
				if err := biz.HistoryDataRunAgain(id); err == nil {
					status = "测试完成，请刷新列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("再来一次", action.Ajax("historyData_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.HistoryDataRunAgain(id); err == nil {
				status = "测试完成，请刷新列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	apps := biz.GetApps()
	products := biz.GetProducts()
	info.AddSelectBox("所属产品", products, action.FieldFilter("product"))
	info.AddSelectBox("所属应用", apps, action.FieldFilter("app"))
	info.AddSelectBox("测试结果", types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))

	info.SetTable("scene_data_test_history").SetTitle("数据测试历史").SetDescription("场景数据测试历史")

	formList := senceDataTestHistory.GetForm()
	formList.AddField("唯一标识", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("数据描述", "name", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("接口ID", "api_id", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("所属应用", "app", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("数据文件", "content", db.Longtext, form.Text)
	formList.AddField("测试结果", "result", db.Varchar, form.Text)
	formList.AddField("失败原因", "fail_reason", db.Longtext, form.TextArea)
	formList.AddField("环境类型", "env_type", db.Int, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "开发", Value: "1"},
			{Text: "测试", Value: "2"},
			{Text: "预发", Value: "3"},
			{Text: "演示", Value: "4"},
			{Text: "生产", Value: "5"},
		}).FieldDefault("2")
	formList.AddField("所属产品", "product", db.Varchar, form.Text)
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("scene_data_test_history").SetTitle("数据测试历史").SetDescription("场景数据测试历史")

	formList.SetPostHook(func(values form2.Values) (err error) {
		fileName := values["content"][0]
		id := values["id"][0]
		err = biz.ModifyEditedData(id, fileName)
		return
	})

	return senceDataTestHistory
}
