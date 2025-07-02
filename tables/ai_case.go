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
	"html/template"
	"strings"
)

func GetAiCaseTable(ctx *context.Context) table.Table {

	aiCase := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	info := aiCase.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name

	info.SetFilterFormLayout(form.LayoutThreeCol)

	aiPlatforms := biz.GetAiCreatePlatform()
	aiTemplates := biz.GetAiTemplateOptions("1")
	products := biz.GetProducts()
	caseTypes := biz.GetTestcaseType()

	info.AddField("自增主键", "id", db.Int).
		FieldFilterable()
	info.AddField("用例编号", "case_number", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("用例名称", "case_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField("所属模块", "module", db.Varchar).
		FieldFilterable()
	info.AddField("用例类型", "case_type", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(caseTypes)
	info.AddField("优先级", "priority", db.Varchar).
		FieldFilterable()
	info.AddField("前置条件", "pre_condition", db.Varchar)
	info.AddField("测试范围", "test_range", db.Varchar)
	info.AddField("测试步骤", "test_steps", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTMLEscapeString(model.Value)
		})
	info.AddField("预期结果", "expect_result", db.Varchar)
	info.AddField("是否自动化", "auto", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return "否"
			} else if model.Value == "1" {
				return "是"
			} else if model.Value == "2" {
				return "部分是"
			}
			return "是"
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "0", Text: "否"},
		{Value: "1", Text: "是"},
		{Value: "2", Text: "部分是"},
	})
	info.AddField("引入版本", "intro_version", db.Varchar).
		FieldFilterable()
	//info.AddField("用例版本", "case_version", db.Int)
	info.AddField("关联产品", "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(products)
	info.AddField("生成来源", "source", db.Varchar).
		FieldFilterable()
	info.AddField("取用状态", "use_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "初始"
			} else if model.Value == "2" {
				return "取用"
			} else if model.Value == "3" {
				return "废弃"
			}
			return "初始"
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "初始"},
		{Value: "2", Text: "取用"},
		{Value: "3", Text: "废弃"},
	})
	info.AddField("改造状态", "modify_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "初始"
			} else if model.Value == "2" {
				return "人工改造"
			} else if model.Value == "3" {
				return "自动改造"
			}
			return "初始"
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "初始"},
		{Value: "2", Text: "人工改造"},
		{Value: "3", Text: "自动改造"},
	})
	info.AddField("创建人", "create_user", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(80)
	info.AddField("创建时间", "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField("更新时间", "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange}).
		FieldHide()
	info.AddField("删除时间", "deleted_at", db.Timestamp).
		FieldHide()

	info.AddButton("AI导入", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_case_import",
		Title:  "AI导入用例",
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		panel.AddField("引入版本", "intro_version", db.Varchar, form.Text)
		panel.AddField("所属产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField("生成平台", "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField("会话ID", "conversation_id", db.Varchar, form.Text).
			FieldHelpMsg("生成平台的会话ID，会话ID和原生回复二选一")
		panel.AddField("原生回复", "raw_reply", db.Varchar, form.TextArea).
			FieldHelpMsg("生成平台的回复原文")

		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_case_import"))

	info.AddButton("AI生成", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_case_create_by_create_desc",
		Title:  "AI生成用例",
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		panel.AddField("智能模板", "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiTemplates).FieldDefault(aiTemplates[0].Value)
		panel.AddField("引入版本", "intro_version", db.Varchar, form.Text)
		panel.AddField("所属产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField("生成平台", "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField("生成指令", "create_desc", db.Varchar, form.TextArea)
		//panel.AddField("生成指令", "create_desc", db.Varchar, form.RichText)   // 大模型不支持富文本的解析和调用
		panel.AddField("上传文件", "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
			"maxFileCount": 10,
		}).FieldHelpMsg("需生成的平台支持文档图片等的解析")
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/ai_case_create_by_create_desc"))

	info.AddButton("AI优化", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_case_optimize",
		Title:  "AI优化用例",
		Width:  "900px",
		Height: "480px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		panel.AddField("已选择编号", "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
		panel.AddField("优化平台", "optimize_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField("优化指令", "optimize_desc", db.Varchar, form.TextArea)
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/ai_case_optimize"))

	info.AddButton("取用", icon.Android, action.Ajax("ai_case_batch_use",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if idStr == "," {
				status = "请先选择数据再取用"
				return false, status, ""
			}
			ids := strings.Trim(idStr, ",")
			if err := biz.UseAiCase(ids, userNameSub); err == nil {
				status = "取用成功，请前往[用例-测试用例]列表查看"
			} else {
				status = fmt.Sprintf("取用失败：%s: %s", ids, err)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.AddActionButton("取用", action.Ajax("ai_case_use",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.UseAiCase(id, userNameSub); err == nil {
				status = "取用成功，请前往[用例-测试用例]列表查看"
			} else {
				status = fmt.Sprintf("取用失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.SetTable("ai_case").SetTitle("智能用例").SetDescription("智能用例")

	formList := aiCase.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("用例编号", "case_number", db.Varchar, form.Text)
	formList.AddField("用例名称", "case_name", db.Varchar, form.Text)
	formList.AddField("所属模块", "module", db.Varchar, form.Text)
	formList.AddField("用例类型", "case_type", db.Varchar, form.Select).
		FieldOptions(caseTypes)
	formList.AddField("优先级", "priority", db.Varchar, form.Text)
	formList.AddField("前置条件", "pre_condition", db.Varchar, form.TextArea)
	formList.AddField("测试范围", "test_range", db.Varchar, form.TextArea)
	formList.AddField("测试步骤", "test_steps", db.Varchar, form.TextArea)
	formList.AddField("预期结果", "expect_result", db.Varchar, form.TextArea)
	formList.AddField("是否自动化", "auto", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "0", Text: "否"},
			{Value: "1", Text: "是"},
			{Value: "2", Text: "部分是"},
		}).FieldDefault("0")
	formList.AddField("引入版本", "intro_version", db.Varchar, form.Text)
	//formList.AddField("用例版本", "case_version", db.Int, form.Number)
	formList.AddField("生成来源", "source", db.Varchar, form.Text).
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("取用状态", "use_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "初始"},
			{Value: "2", Text: "取用"},
			{Value: "3", Text: "废弃"},
		}).
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenUpdate().
		FieldDefault("1")
	formList.AddField("改造状态", "modify_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "初始"},
			{Value: "2", Text: "人工改造"},
			{Value: "3", Text: "自动改造"},
		}).
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenUpdate().
		FieldDefault("1")
	formList.AddField("关联产品", "product", db.Varchar, form.SelectSingle).
		FieldOptions(products)
	formList.AddField("创建人", "create_user", db.Varchar, form.Text).
		FieldDefault(userName).FieldHide()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("ai_case").SetTitle("智能用例").SetDescription("智能用例")

	formList.SetPostHook(func(values form2.Values) (err error) {
		id := values["id"][0]
		err = biz.UpdateAiCaseStatus(id)
		return
	})

	detail := aiCase.GetDetail()
	detail.AddField("自增主键", "id", db.Int)
	detail.AddField("用例编号", "case_number", db.Varchar)
	detail.AddField("用例名称", "case_name", db.Varchar)
	detail.AddField("所属模块", "module", db.Varchar)
	detail.AddField("用例类型", "case_type", db.Varchar)
	detail.AddField("优先级", "priority", db.Varchar)
	detail.AddField("前置条件", "pre_condition", db.Text)
	detail.AddField("测试范围", "test_range", db.Text)
	detail.AddField("测试步骤", "test_steps", db.Text)
	detail.AddField("预期结果", "expect_result", db.Text)
	detail.AddField("是否自动化", "auto", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return "否"
			}
			if model.Value == "1" {
				return "是"
			}
			if model.Value == "2" {
				return "部分是"
			}
			return "是"
		})
	detail.AddField("引入版本", "intro_version", db.Varchar)
	//detail.AddField("用例版本", "case_version", db.Int)
	detail.AddField("生成来源", "source", db.Varchar)
	detail.AddField("取用状态", "use_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "初始"
			}
			if model.Value == "2" {
				return "取用"
			}
			if model.Value == "3" {
				return "废弃"
			}
			return "初始"
		})
	detail.AddField("改造状态", "modify_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "初始"
			}
			if model.Value == "2" {
				return "人工改造"
			}
			if model.Value == "3" {
				return "自动改造"
			}
			return "初始"
		})
	detail.AddField("所属产品", "product", db.Varchar)
	detail.AddField("创建人", "create_user", db.Varchar)
	detail.AddField("创建时间", "created_at", db.Timestamp)
	detail.AddField("更新时间", "updated_at", db.Timestamp)
	detail.AddField("删除时间", "deleted_at", db.Timestamp).FieldHide()

	detail.SetTable("ai_case").SetTitle("智能用例").SetDescription("智能用例")
	return aiCase
}
