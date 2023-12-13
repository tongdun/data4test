package tables

import (
	"data4perf/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"strings"
)

func GetSceneTestHistoryTable(ctx *context.Context) table.Table {

	senceTestHistory := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := senceTestHistory.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField("唯一标识", "id", db.Int).
		FieldFilterable().
		FieldTrimSpace()
	info.AddField("场景描述", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("数据文件列表", "api_list", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("最近数据文件", "last_file", db.Longtext).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("场景类型", "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "默认"
			}
			if model.Value == "2" {
				return "比较"
			}
			return "默认"
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "默认"},
		{Value: "2", Text: "比较"},
	}).
		FieldHide()
	info.AddField("测试结果", "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	})
	info.AddField("失败原因", "fail_reason", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
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
	info.AddField("备注", "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("所属产品", "product", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton("再来一次", icon.Android, action.Ajax("historyPlaybook_batch_again",
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
				if err := biz.RunHistoryPlaybook(id, "again"); err == nil {
					status = "测试完成，请刷新列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("再来一次", action.Ajax("historyPlaybook_again",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunHistoryPlaybook(id, "again"); err == nil {
				status = "测试完成，请刷新列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("继续", icon.Android, action.Ajax("historyPlaybook_batch_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再继续"
				return false, status, ""
			}

			ids := strings.Split(idStr, ",")

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.RunHistoryPlaybook(id, "continue"); err == nil {
					status = "测试完成，请前往[结果详情]列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("继续", action.Ajax("historyPlaybook_continue",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string

			if err := biz.RunHistoryPlaybook(id, "continue"); err == nil {
				status = "测试完成，请前往[结果详情]列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	products := biz.GetProducts()
	info.AddSelectBox("关联产品", products, action.FieldFilter("product"))

	info.AddSelectBox("测试结果", types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))

	info.SetTable("scene_test_history").SetTitle("场景测试历史").SetDescription("场景测试历史")

	formList := senceTestHistory.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("场景描述", "name", db.Varchar, form.Text).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("数据文件列表", "api_list", db.Longtext, form.RichText).FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("最近数据文件", "last_file", db.Varchar, form.Text)
	formList.AddField("场景类型", "scene_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "默认", Value: "1"},
			{Text: "比较", Value: "2"},
		}).FieldDefault("1")
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
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("所属产品", "product", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldNowWhenInsert().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldNowWhenUpdate().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("scene_test_history").SetTitle("场景测试历史").SetDescription("场景测试历史")

	detail := senceTestHistory.GetDetail()
	detail.AddField("唯一标识", "id", db.Int)
	detail.AddField("场景描述", "name", db.Varchar)
	detail.AddField("数据文件列表", "api_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Replace(model.Value, ",", ",<br>", -1)
		})
	detail.AddField("最近数据文件", "last_file", db.Varchar)
	detail.AddField("场景类型", "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "默认"
			}
			if model.Value == "2" {
				return "比较"
			}
			return "默认"
		})
	detail.AddField("测试结果", "result", db.Varchar)
	detail.AddField("失败原因", "fail_reason", db.Longtext)
	detail.AddField("环境类型", "env_type", db.Int).
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
		})
	detail.AddField("备注", "remark", db.Longtext)
	detail.AddField("所属产品", "product", db.Varchar)
	detail.AddField("创建时间", "created_at", db.Timestamp)
	detail.AddField("更新时间", "updated_at", db.Timestamp)
	detail.AddField("删除时间", "deleted_at", db.Timestamp)

	detail.SetTable("scene_test_history").SetTitle("场景测试历史").SetDescription("场景测试历史")

	return senceTestHistory
}
