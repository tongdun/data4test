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

func GetAiDataTable(ctx *context.Context) table.Table {

	aiData := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	user := auth.Auth(ctx)
	userName := user.Name
	info := aiData.GetInfo().HideFilterArea()
	info.SetFilterFormLayout(form.LayoutThreeCol)

	apps := biz.GetApps()
	aiAnalysisTemplates := biz.GetAiTemplateOptions("7")
	//aiDataTemplates := biz.GetAiTemplateOptions("2")
	aiPlatforms := biz.GetAiCreatePlatform()
	products := biz.GetProducts()

	info.AddField("自增主键", "id", db.Int).
		FieldFilterable()
	info.AddField("数据描述", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("接口ID", "api_id", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("所属应用", "app", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(apps)
	info.AddField("生成来源", "source", db.Varchar).
		FieldFilterable()
	info.AddField("文件名称", "file_name", db.Text).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return template.Default().
				Link().
				SetURL("/admin/fm/ai_data/preview?path=/" + value.Value).
				SetContent(template2.HTML(value.Value)).
				OpenInNewTab().
				SetTabTitle(template.HTML("数据文件")).
				GetContent()
		}).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField("文件类型", "file_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "标准"
			} else if model.Value == "2" {
				return "Python"
			} else if model.Value == "3" {
				return "Shell"
			} else if model.Value == "4" {
				return "Bat"
			} else if model.Value == "5" {
				return "Jmeter"
			} else if model.Value == "99" {
				return "其他"
			}
			return "标准"
		}).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "1", Text: "标准"},
			{Value: "2", Text: "Python"},
			{Value: "3", Text: "Shell"},
			{Value: "4", Text: "Bat"},
			{Value: "5", Text: "Jmeter"},
			{Value: "99", Text: "其他"},
		})
	info.AddField("文件内容", "content", db.Text).
		FieldHide()
	info.AddField("测试结果", "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "pass", Text: "pass"},
			{Value: "fail", Text: "fail"},
		})
	info.AddField("失败原因", "fail_reason", db.Text)
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
		}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
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
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
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
		Id:     "/ai_data_import",
		Title:  "AI导入数据",
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		panel.AddField("生成平台", "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField("引入版本", "intro_version", db.Varchar, form.Text).
			FieldHelpMsg("若提供，则相关描述带版本后缀信息")
		panel.AddField("所属产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField("会话ID", "conversation_id", db.Varchar, form.Text).
			FieldHelpMsg("生成平台的会话ID，会话ID和原生回复二选一")
		panel.AddField("原生回复", "raw_reply", db.Varchar, form.TextArea).
			FieldHelpMsg("生成平台的回复原文")

		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_data_import"))

	//info.AddButton("AI生成", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
	//	Id:     "/ai_data_create_by_create_desc",
	//	Title:  "AI生成数据",
	//	Width:  "900px",
	//	Height: "720px", // TextArea
	//}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
	//	panel.AddField("智能模板", "ai_template", db.Varchar, form.SelectSingle).
	//		FieldOptions(aiDataTemplates).FieldDefault(aiDataTemplates[0].Value)
	//	panel.AddField("生成平台", "create_platform", db.Varchar, form.SelectSingle).
	//		FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
	//	panel.AddField("引入版本", "intro_version", db.Varchar, form.Text).
	//		FieldHelpMsg("若提供，则相关描述带版本后缀信息")
	//	panel.AddField("关联产品", "product", db.Varchar, form.SelectSingle).
	//		FieldOptions(products).
	//		FieldDefault(aiPlatforms[0].Value).
	//		FieldHelpMsg("用于生成智能场景")
	//	panel.AddField("生成指令", "create_desc", db.Varchar, form.TextArea)
	//	panel.AddField("上传文件", "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
	//		"maxFileCount": 10,
	//	}).FieldHelpMsg("需生成的平台支持文档图片等的解析")
	//	panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
	//	return panel
	//}, "/ai_data_create_by_create_desc"))

	//info.AddButton("AI优化", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
	//	Id:     "/ai_data_optimize",
	//	Title:  "AI优化数据",
	//	Width:  "900px",
	//	Height: "480px", // TextArea
	//}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
	//	ids := ctx.FormValue("ids")
	//	panel.AddField("已选择编号", "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
	//	panel.AddField("优化平台", "optimize_platform", db.Varchar, form.SelectSingle).
	//		FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
	//	panel.AddField("优化指令", "optimize_desc", db.Varchar, form.TextArea)
	//	panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
	//	return panel
	//}, "/ai_data_optimize"))

	info.AddButton("AI分析", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_data_test_and_analysis",
		Title:  "测试执行后进行AI结果分析",
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		panel.AddField("已选择编号", "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
		panel.AddField("分析模板", "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiAnalysisTemplates).FieldDefault(aiAnalysisTemplates[0].Value)
		panel.AddField("分析平台", "analysis_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).
			FieldDefault(aiPlatforms[0].Value)
		panel.AddField("关联产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).
			FieldDefault(products[0].Value).
			FieldHelpMsg("用于执行智能数据")
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_data_test_and_analysis"))

	info.AddButton("测试", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_data_test",
		Title:  "数据测试",
		Width:  "900px",
		Height: "300px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		panel.AddField("已选择编号", "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
		panel.AddField("关联产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).
			FieldDefault(products[0].Value).
			FieldHelpMsg("选择执行环境")
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_data_test"))

	info.AddButton("取用", icon.Android, action.Ajax("ai_data_batch_use",
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
			if err := biz.UseAiData(ids, userNameSub); err == nil {
				status = "取用成功，请前往[用例-测试用例]列表查看"
			} else {
				status = fmt.Sprintf("取用失败：%s: %s", ids, err)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.AddActionButton("取用", action.Ajax("ai_data_use",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.UseAiData(id, userNameSub); err == nil {
				status = "取用成功，请前往[用例-测试用例]列表查看"
			} else {
				status = fmt.Sprintf("取用失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.SetTable("ai_data").SetTitle("智能数据").SetDescription("智能数据")
	fileNameHelp := template.HTML("e.g.: 类型-模块-功能描述.yml / 类型-模块-功能描述.json / 类型-模块-功能描述.py / 类型-模块-功能描述.sh / 类型-模块-功能描述.jmx")
	fileTypeMsg := template.HTML("默认值: 标准数据<br>标准数据: 推荐优先使用，结构化编写，简单高效<br>Python脚本: 标准数据无法支持的场景<br>Shell脚本: 标准数据无法支持的场景<br>Bat脚本: 适用于windows系统,.bat文件<br>Jmeter脚本: 适用于性能测试，.jmt文件，系统参数JmeterRunConfig可控制执行参数(待合入)<br>其他脚本: 根据文件后缀从系统参数中获取执行引擎或脚本中自定义有执行引擎<br>脚本执行引擎优先级: 系统参数scriptRunEngine定义 > 脚本中首行定义<br>执行引擎以文件名后缀为key可任意扩展，执行引擎需自行配置环境")

	formList := aiData.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("数据描述", "name", db.Varchar, form.Text)
	formList.AddField("接口ID", "api_id", db.Varchar, form.Text)
	formList.AddField("所属应用", "app", db.Varchar, form.SelectSingle).
		FieldOptions(apps)
	formList.AddField("生成来源", "source", db.Varchar, form.Text).
		FieldDisplayButCanNotEditWhenUpdate()
	formList.AddField("文件名称", "file_name", db.Varchar, form.Url).
		FieldHelpMsg(fileNameHelp)
	formList.AddField("文件类型", "file_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "标准"},
			{Value: "2", Text: "Python"},
			{Value: "3", Text: "Shell"},
			{Value: "4", Text: "Bat"},
			{Value: "5", Text: "Jmeter"},
			{Value: "99", Text: "其他"},
		}).FieldDefault("1").FieldHelpMsg(fileTypeMsg)
	formList.AddField("文件内容", "content", db.Text, form.TextArea)
	formList.AddField("测试结果", "result", db.Varchar, form.Text)
	formList.AddField("失败原因", "fail_reason", db.Text, form.TextArea)
	formList.AddField("取用状态", "use_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "初始"},
			{Value: "2", Text: "取用"},
			{Value: "3", Text: "废弃"},
		}).
		FieldDefault("1").
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenUpdate()
	formList.AddField("改造状态", "modify_status", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "初始"},
			{Value: "2", Text: "人工改造"},
			{Value: "3", Text: "自动改造"},
		}).
		FieldDefault("1").
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenUpdate()
	formList.AddField("创建人", "create_user", db.Varchar, form.Text).
		FieldDefault(userName).FieldHide()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("ai_data").SetTitle("智能数据").SetDescription("智能数据")

	formList.SetPostHook(func(values form2.Values) (err error) {
		content := values["content"][0]
		fileName := values["file_name"][0]
		id := values["id"][0]
		err = biz.BakOldAiDataVer(id, content, fileName)
		err = biz.UpdateAiDataStatus(id)
		return
	})

	detail := aiData.GetDetail()
	detail.AddField("自增主键", "id", db.Int)
	detail.AddField("数据描述", "name", db.Varchar)
	detail.AddField("接口ID", "api_id", db.Varchar)
	detail.AddField("所属应用", "app", db.Varchar)
	detail.AddField("文件名称", "file_name", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			linkStr := fmt.Sprintf("<a href=\"/admin/fm/ai_data/preview?path=/%s\">%s</a>", model.Value, model.Value)
			return linkStr
		})
	detail.AddField("生成来源", "source", db.Varchar)
	detail.AddField("文件类型", "file_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "标准"
			}
			if model.Value == "2" {
				return "Python"
			}
			if model.Value == "3" {
				return "Shell"
			}
			if model.Value == "4" {
				return "Bat"
			}
			if model.Value == "5" {
				return "Jmeter"
			}
			if model.Value == "99" {
				return "其他"
			}
			return "标准"
		})
	//detail.AddField("文件内容", "content", db.Text)
	detail.AddField("测试结果", "result", db.Varchar)
	detail.AddField("失败原因", "fail_reason", db.Text)
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
	detail.AddField("创建人", "create_user", db.Varchar)
	detail.AddField("创建时间", "created_at", db.Timestamp)
	detail.AddField("更新时间", "updated_at", db.Timestamp)
	detail.AddField("删除时间", "deleted_at", db.Timestamp).FieldHide()

	detail.SetTable("ai_data").SetTitle("智能数据").SetDescription("智能数据")

	return aiData
}
