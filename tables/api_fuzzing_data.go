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

func GetApiFuzzingDataTable(ctx *context.Context) table.Table {

	apiFuzzingData := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := apiFuzzingData.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("数据描述", "data_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(200)
	info.AddField("接口描述", "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(200)
	info.AddField("关联接口", "api_id", db.Varchar).FieldWidth(200).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(300)
	info.AddField("所属模块", "api_module", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)
	info.AddField("Header", "header", db.Longtext)
	info.AddField("UrlQuery", "url_query", db.Longtext).FieldWidth(250)
	info.AddField("Body", "body", db.Longtext).FieldWidth(300)
	info.AddField("执行次数", "run_num", db.Int).FieldEditAble().
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("预期结果", "expected_result", db.Varchar)
	info.AddField("实际结果", "actual_result", db.Varchar)
	info.AddField("测试结果", "result", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("失败原因", "fail_reason", db.Longtext)
	info.AddField("Response", "response", db.Longtext)
	info.AddField("关联应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160)
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()
	info.AddButton("生成场景数据", icon.Android, action.Ajax("create_batch_scene_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再生成"
				return false, status, ""
			}
			if err := biz.CreateSceneData(idStr, "fuzzing", ""); err == nil {
				status = "生成完成，请前往[场景数据]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", idStr, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddActionButton("生成场景数据", action.Ajax("create_scene_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneData(id, "fuzzing", ""); err == nil {
				status = "生成完成，请前往[场景数据]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("生成JSON场景数据", icon.Android, action.Ajax("create_batch_scene_fuzzing_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再生成"
				return false, status, ""
			}
			if err := biz.CreateSceneData(idStr, "fuzzing", "json"); err == nil {
				status = "生成完成，请前往[场景数据]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", idStr, err)
				return false, status, ""
			}
			return true, status, ""
		}))

	info.AddActionButton("生成JSON场景数据", action.Ajax("create_scene_fuzzing_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneData(id, "fuzzing", "json"); err == nil {
				status = "生成完成，请前往[场景数据]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("测试", icon.Android, action.Ajax("fuzzingdata_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再测试"
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					status = "测试完成，请前往[结果详情]列表查看"
					continue
				}
				if err := biz.RunFuzzData(id); err == nil {
					status = "测试完成，请前往[结果详情]列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton("测试", action.Ajax("fuzzingdata_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunFuzzData(id); err == nil {
				status = "测试完成，请前往[结果详情]列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))


	apps := biz.GetApps()
	info.AddSelectBox("关联应用", apps, action.FieldFilter("app"))

	info.AddSelectBox("测试结果", types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("result"))

	info.SetTable("api_fuzzing_data").SetTitle("模糊数据").SetDescription("模糊数据")

	formList := apiFuzzingData.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("数据描述", "data_desc", db.Varchar, form.Text)
	formList.AddField("接口描述", "api_desc", db.Varchar, form.Text)
	formList.AddField("关联接口", "api_id", db.Varchar, form.Text)
	formList.AddField("所属模块", "api_module", db.Varchar, form.Text)
	formList.AddField("Header", "header", db.Longtext, form.Text)
	formList.AddField("UrlQuery", "url_query", db.Longtext, form.Text)
	formList.AddField("Body", "body", db.Longtext, form.Text)
	formList.AddField("执行次数", "run_num", db.Int, form.Number)
	formList.AddField("预期结果", "expected_result", db.Varchar, form.Text)
	formList.AddField("实际结果", "actual_result", db.Varchar, form.Text)
	formList.AddField("测试结果", "result", db.Varchar, form.Text)
	formList.AddField("失败原因", "fail_reason", db.Longtext, form.Text)
	formList.AddField("Response", "response", db.Longtext, form.Text)
	formList.AddField("关联应用", "app", db.Varchar, form.Text)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_fuzzing_data").SetTitle("模糊数据").SetDescription("模糊数据")

	return apiFuzzingData
}
