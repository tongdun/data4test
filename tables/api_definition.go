package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
)

func GetApiDefinitionTable(ctx *context.Context) table.Table {

	apiDefinition := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	user := auth.Auth(ctx)
	userName := user.Name

	// 全局域的在有数据更新时，子弹框的不会实时更新
	//aiCaseTemplates := biz.GetAiTemplateOptions("1")
	//aiDataTemplates := biz.GetAiTemplateOptions("2")
	//products := biz.GetProducts()
	//aiPlatforms := biz.GetAiCreatePlatform()

	info := apiDefinition.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField("自增主键", "id", db.Int).
		FieldHide()
	info.AddField("接口ID", "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldWidth(150).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiAutoDataList(model.Value, model.ID)
		})
	info.AddField("所属模块", "api_module", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)
	info.AddField("接口描述", "api_desc", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(150)
	info.AddField("请求方法", "http_method", db.Enum).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldEditAble().FieldWidth(120)
	info.AddField("请求路径", "path", db.Varchar).FieldWidth(120).FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("Header参数", "header", db.JSON).FieldWidth(300).
		FieldHide()
	info.AddField("Path参数", "path_variable", db.JSON).FieldWidth(150).
		FieldHide()
	info.AddField("Query参数", "query_parameter", db.JSON).FieldWidth(300).
		FieldHide()
	info.AddField("Body参数", "body", db.JSON).FieldWidth(300).
		FieldHide()
	info.AddField("Resp参数", "response", db.JSON).FieldWidth(300).
		FieldHide()
	info.AddField("接口版本", "version", db.Int).FieldWidth(80)
	info.AddField("接口状态", "api_status", db.Enum).FieldWidth(100).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "新增"},
		{Value: "2", Text: "被删除"},
		{Value: "3", Text: "被修改"},
		{Value: "4", Text: "保持原样"},
		{Text: "Header被修改", Value: "30"},
		{Text: "Path被修改", Value: "31"},
		{Text: "Query被修改", Value: "32"},
		{Text: "Body被修改", Value: "33"},
		{Text: "Resp被修改", Value: "34"},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return "新增"
		} else if model.Value == "2" {
			return "被删除"
		} else if model.Value == "3" {
			return "被修改"
		} else if model.Value == "4" {
			return "保持原样"
		} else if model.Value == "30" {
			return "Header被修改"
		} else if model.Value == "31" {
			return "Path被修改"
		} else if model.Value == "32" {
			return "Query被修改"
		} else if model.Value == "33" {
			return "Body被修改"
		} else if model.Value == "34" {
			return "Resp被修改"
		}
		return "新增"
	})
	info.AddField("变更内容", "change_content", db.JSON).
		FieldHide()
	info.AddField("规范检查", "check", db.Varchar).FieldWidth(120)
	info.AddField("规范检查失败原因", "api_check_fail_reason", db.JSON).
		FieldHide()
	info.AddField("是否需自动化", "is_need_auto", db.Enum).
		FieldEditAble(editType.Select).FieldEditOptions(types.FieldOptions{
		{Value: "-1", Text: "否"},
		{Value: "1", Text: "是"},
	}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "-1", Text: "否"},
		{Value: "1", Text: "是"},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return "是"
		} else if model.Value == "-1" {
			return "否"
		}
		return "是"
	})
	info.AddField("是否已自动化", "is_auto", db.Enum).
		FieldEditAble(editType.Select).FieldEditOptions(types.FieldOptions{
		{Value: "-1", Text: "否"},
		{Value: "1", Text: "是"},
	}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "-1", Text: "否"},
		{Value: "1", Text: "是"},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return "是"
		} else if model.Value == "-1" {
			return "否"
		}
		return "否"
	})

	info.AddField("所属应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).FieldWidth(120)

	info.AddField("备注", "remark", db.JSON).
		FieldHide()
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddActionButton("生成数据用例", action.Ajax("create_scene_raw_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneDataFromRaw(id, ""); err == nil {
				status = "生成完成，请前往[数据文件]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton("生成JSON数据用例", action.Ajax("create_scene_raw_data_json",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateSceneDataFromRaw(id, "json"); err == nil {
				status = "生成完成，请前往[数据文件]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton("生成测试数据", action.Ajax("create_test_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateTestData(id); err == nil {
				status = "生成完成，请前往[测试数据]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton("生成模糊数据", action.Ajax("create_fuzzing_data",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreateFuzzingData(id); err == nil {
				status = "生成完成，请前往[模糊数据]列表查看"
			} else {
				status = fmt.Sprintf("生成失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton("测试数据测试", action.Ajax("test_data_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTestData(id); err == nil {
				status = "测试完成，请前往[结果详情]列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	//info.AddButton("模糊数据测试", icon.Android, action.Ajax("fuzzing_data_batch_test",
	//	func(ctx *context.Context) (success bool, msg string, data interface{}) {
	//		idStr := ctx.FormValue("ids")
	//		var status string
	//		if idStr == "," {
	//			status = "请先选择数据再测试"
	//			return false, status, ""
	//		}
	//		ids := strings.Split(idStr, ",")
	//		for _, id := range ids {
	//			if len(id) == 0 {
	//				status = "测试完成，请前往[结果详情]列表查看"
	//				continue
	//			}
	//			if err := biz.RunTestData(id); err == nil {
	//				status = "测试完成，请前往[结果详情]列表查看"
	//			} else {
	//				status = fmt.Sprintf("测试失败：%s: %s", id, err)
	//				return false, status, ""
	//			}
	//		}
	//		return true, status, ""
	//	}))

	info.AddActionButton("模糊数据测试", action.Ajax("fuzzing_data_test",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.RunTestData(id); err == nil {
				status = "测试完成，请前往[结果详情]列表查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton("生成场景用例", action.Ajax("create_playbook_case",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.CreatePlaybookByAPIId(id, userName); err == nil {
				status = "生成完成，请前往[场景列表]查看"
			} else {
				status = fmt.Sprintf("测试失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.AddButton("AI生成用例", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_create_case_by_api_define",
		Title:  "AI生成用例",
		Width:  "900px",
		Height: "540px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		products := biz.GetProducts()
		aiCaseTemplates := biz.GetAiTemplateOptions("1")
		aiPlatforms := biz.GetAiCreatePlatform()
		panel.AddField("已选择编号", "ids", db.Varchar, form.Text).
			FieldDefault(ids).
			FieldHide()
		panel.AddField("智能模板", "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiCaseTemplates).FieldDefault(aiCaseTemplates[0].Value)
		panel.AddField("引入版本", "intro_version", db.Varchar, form.Text)
		panel.AddField("所属产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField("生成平台", "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.EnableAjax("生成任务已在后台运行,请稍后前往[助手-智能用例]列表查看生成用例")
		return panel
	}, "/ai_create_case_by_api_define"))

	info.AddButton("AI生成数据", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_create_data_by_api_define",
		Title:  "AI生成数据",
		Width:  "900px",
		Height: "540px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		products := biz.GetProducts()
		aiDataTemplates := biz.GetAiTemplateOptions("2")
		aiPlatforms := biz.GetAiCreatePlatform()
		panel.AddField("已选择编号", "ids", db.Varchar, form.Text).
			FieldDefault(ids).
			FieldHide()
		panel.AddField("智能模板", "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiDataTemplates).
			FieldDefault(aiDataTemplates[0].Value)
		panel.AddField("生成平台", "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).
			FieldDefault(aiPlatforms[0].Value)
		panel.AddField("引入版本", "intro_version", db.Varchar, form.Text).
			FieldHelpMsg("若提供，则相关描述带版本后缀信息")
		panel.AddField("关联产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).
			FieldDefault(aiPlatforms[0].Value).
			FieldHelpMsg("用于生成智能场景")
		panel.EnableAjax("生成任务已在后台运行,请稍后前往[助手-智能数据/智能场景]列表查看生成数据/场景")
		return panel
	}, "/ai_create_data_by_api_define"))

	apps := biz.GetApps()
	info.AddSelectBox("关联应用", apps, action.FieldFilter("app"))

	info.AddSelectBox("请求方法", types.FieldOptions{
		{Value: "get", Text: "get"},
		{Value: "post", Text: "post"},
		{Value: "delete", Text: "delete"},
		{Value: "put", Text: "put"},
	}, action.FieldFilter("http_method"))

	info.AddSelectBox("规范检查", types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	}, action.FieldFilter("check"))

	info.SetTable("api_definition").SetTitle("接口定义").SetDescription("接口定义")

	formList := apiDefinition.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("接口ID", "api_id", db.Varchar, form.Text)
	formList.AddField("所属模块", "api_module", db.Varchar, form.Text)
	formList.AddField("接口描述", "api_desc", db.Varchar, form.Text)
	formList.AddField("请求方法", "http_method", db.Enum, form.Text).
		FieldDefault("get")
	formList.AddField("请求路径", "path", db.Varchar, form.Text)
	formList.AddField("Header参数", "header", db.Longtext, form.TextArea)
	formList.AddField("Path参数", "path_variable", db.Longtext, form.TextArea)
	formList.AddField("Query参数", "query_parameter", db.Longtext, form.TextArea)
	formList.AddField("Body参数", "body", db.Longtext, form.TextArea)
	formList.AddField("Resp参数", "response", db.Longtext, form.TextArea)
	formList.AddField("接口版本", "version", db.Int, form.Number).
		FieldDefault("1").
		FieldDisplayButCanNotEditWhenUpdate().
		FieldDisplayButCanNotEditWhenCreate()
	formList.AddField("接口状态", "api_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "新增", Value: "1"},
			{Text: "被删除", Value: "2"},
			{Text: "被修改", Value: "3"},
			{Text: "保持原样", Value: "4"},
			{Text: "Header被修改", Value: "30"},
			{Text: "Path被修改", Value: "31"},
			{Text: "Query被修改", Value: "32"},
			{Text: "Body被修改", Value: "33"},
			{Text: "Resp被修改", Value: "34"},
		}).FieldDefault("1")
	formList.AddField("变更内容", "change_content", db.Longtext, form.TextArea)
	formList.AddField("规范检查", "check", db.Varchar, form.Text)
	formList.AddField("规范检查失败原因", "api_check_fail_reason", db.Longtext, form.TextArea)
	formList.AddField("是否需自动化", "is_need_auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "否", Value: "-1"},
			{Text: "是", Value: "1"},
		}).FieldDefault("1")
	formList.AddField("是否已自动化", "is_auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "否", Value: "-1"},
			{Text: "是", Value: "1"},
		}).FieldDefault("-1")
	formList.AddField("所属应用", "app", db.Varchar, form.Text)
	formList.AddField("备注", "remark", db.Longtext, form.TextArea)
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("api_definition").SetTitle("接口定义").SetDescription("接口定义")
	formList.SetPostHook(func(values form2.Values) (err error) {
		id := values["id"][0]
		err = biz.UpdateApiDefVer(id)
		return
	})

	detail := apiDefinition.GetDetail()
	detail.AddField("自增主键", "id", db.Int)
	detail.AddField("接口ID", "api_id", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetApiAutoDataList(model.Value, model.ID)
		})
	detail.AddField("所属模块", "api_module", db.Varchar)
	detail.AddField("接口描述", "api_desc", db.Varchar)
	detail.AddField("请求方法", "http_method", db.Enum)
	detail.AddField("请求路径", "path", db.Varchar)
	detail.AddField("Header参数", "header", db.JSON)
	detail.AddField("Path参数", "path_variable", db.JSON)
	detail.AddField("Query参数", "query_parameter", db.JSON)
	detail.AddField("Body参数", "body", db.JSON)
	detail.AddField("Resp参数", "response", db.JSON)
	detail.AddField("接口版本", "version", db.Int)
	detail.AddField("接口状态", "api_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "新增"
			} else if model.Value == "2" {
				return "被删除"
			} else if model.Value == "3" {
				return "被修改"
			} else if model.Value == "4" {
				return "保持原样"
			} else if model.Value == "30" {
				return "Header被修改"
			} else if model.Value == "31" {
				return "Path被修改"
			} else if model.Value == "32" {
				return "Query被修改"
			} else if model.Value == "33" {
				return "Body被修改"
			} else if model.Value == "34" {
				return "Resp被修改"
			}
			return "新增"
		})
	detail.AddField("变更内容", "change_content", db.Longtext)
	detail.AddField("规范检查", "check", db.Varchar)
	detail.AddField("规范检查失败原因", "api_check_fail_reason", db.JSON)
	detail.AddField("是否需自动化", "is_need_auto", db.Enum).
		FieldEditAble(editType.Select).FieldEditOptions(types.FieldOptions{
		{Value: "-1", Text: "否"},
		{Value: "1", Text: "是"},
	}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "-1", Text: "否"},
		{Value: "1", Text: "是"},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return "是"
		} else if model.Value == "-1" {
			return "否"
		}
		return "是"
	})
	detail.AddField("是否已自动化", "is_auto", db.Enum).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "-1", Text: "否"},
		{Value: "1", Text: "是"},
	}).FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "1" {
			return "是"
		} else if model.Value == "-1" {
			return "否"
		}
		return "否"
	})

	detail.AddField("所属应用", "app", db.Varchar)
	detail.AddField("备注", "remark", db.JSON)
	detail.AddField("创建时间", "created_at", db.Timestamp)
	detail.AddField("更新时间", "updated_at", db.Timestamp)
	detail.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	return apiDefinition
}
