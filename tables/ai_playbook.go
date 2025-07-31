package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	template2 "html/template"
	"strings"
)

func GetAiPlaybookTable(ctx *context.Context) table.Table {

	aiPlaybook := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiPlaybook.GetInfo().HideFilterArea()
	info.SetFilterFormLayout(form.LayoutThreeCol)

	user := auth.Auth(ctx)
	userName := user.Name
	products := biz.GetProducts()
	pTypes := types.FieldOptions{
		{Value: "1", Text: "串行中断"},
		{Value: "2", Text: "串行比较"},
		{Value: "3", Text: "串行继续"},
		{Value: "4", Text: "普通并发"},
		{Value: "5", Text: "并发比较"},
	}

	aiAnalysisTemplates := biz.GetAiTemplateOptions("7")
	aiPlatforms := biz.GetAiCreatePlatform()

	info.AddField("自增主键", "id", db.Int).
		FieldFilterable().
		FieldWidth(80)
	info.AddField("场景描述", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldWidth(220)
	info.AddField("数据文件列表", "data_file_list", db.Text).
		FieldWidth(600).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetAiDataFileLinkByDataStr(model.Value)
		})
	info.AddField("最近数据文件", "last_file", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return template.Default().
				Link().
				SetURL("/admin/fm/ai_data/preview?path=/" + value.Value).
				SetContent(template2.HTML(value.Value)).
				OpenInNewTab().
				SetTabTitle(template.HTML("数据文件")).
				GetContent()
		}).FieldWidth(160)
	info.AddField("场景类型", "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "串行中断"
			} else if model.Value == "2" {
				return "串行比较"
			} else if model.Value == "3" {
				return "串行继续"
			} else if model.Value == "4" {
				return "普通并发"
			} else if model.Value == "5" {
				return "并发比较"
			}
			return "串行中断"
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(pTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(pTypes).
		FieldWidth(80)
	info.AddField("优先级", "priority", db.Int).
		FieldFilterable(types.FilterType{FormType: form.Number}).
		FieldSortable().FieldWidth(80).
		FieldEditAble(editType.Text)
	info.AddField("生成来源", "source", db.Varchar).
		FieldFilterable().
		FieldWidth(80)
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
	info.AddField("测试结果", "result", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "pass", Text: "pass"},
		{Value: "fail", Text: "fail"},
	})
	info.AddField("失败原因", "fail_reason", db.Text)
	info.AddField("所属产品", "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).
		FieldFilterOptions(products).
		FieldEditAble(editType.Select).
		FieldEditOptions(products).
		FieldWidth(220)
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
		Id:     "/ai_playbook_import",
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
	}, "/ai_playbook_import"))

	//info.AddButton("AI生成", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
	//	Id:     "/ai_playbook_create_by_create_desc",
	//	Title:  "AI生成数据",
	//	Width:  "900px",
	//	Height: "680px", // TextArea
	//}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
	//	panel.AddField("智能模板", "ai_template", db.Varchar, form.SelectSingle).
	//		FieldOptions(aiDataTemplates).FieldDefault(aiDataTemplates[0].Value)
	//	panel.AddField("生成平台", "create_platform", db.Varchar, form.SelectSingle).
	//		FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
	//	panel.AddField("引入版本", "intro_version", db.Varchar, form.Text).
	//		FieldHelpMsg("若提供，则相关描述带版本后缀信息")
	//	panel.AddField("所属产品", "product", db.Varchar, form.SelectSingle).
	//		FieldOptions(products).FieldDefault(products[0].Value)
	//	panel.AddField("生成指令", "create_desc", db.Varchar, form.TextArea)
	//	panel.AddField("上传文件", "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
	//		"maxFileCount": 10,
	//	}).FieldHelpMsg("需生成的平台支持文档图片等的解析")
	//	panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
	//	return panel
	//}, "/ai_playbook_create_by_create_desc"))

	//info.AddButton("AI优化", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
	//	Id:     "/ai_playbook_optimize",
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
	//}, "/ai_playbook_optimize"))

	info.AddButton("AI分析", icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/ai_playbook_test_and_analysis",
		Title:  "测试执行后进行AI结果分析",
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		panel.AddField("已选择编号", "ids", db.Varchar, form.Text).
			FieldDefault(ids).
			FieldHide()
		panel.AddField("分析模板", "ai_template", db.Varchar, form.SelectSingle).
			FieldOptions(aiAnalysisTemplates).
			FieldDefault(aiAnalysisTemplates[0].Value)
		panel.AddField("分析平台", "analysis_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).
			FieldDefault(aiPlatforms[0].Value)
		panel.AddField("关联产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).
			FieldDefault(products[0].Value).
			FieldHelpMsg("用于执行智能数据")
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_playbook_test_and_analysis"))

	info.AddButton("测试", icon.Android, action.Ajax("ai_playbook_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择场景再测试"
				return false, status, ""
			}
			idList := strings.Split(idStr, ",")
			for _, id := range idList {
				if len(id) == 0 {
					continue
				}

				if err := biz.RunPlaybookFromMgmt(id, "start", "", "ai_playbook"); err == nil {
					status = "测试完成，请前往[结果详情]列表查看"
				} else {
					status = fmt.Sprintf("测试失败：%s: %s", id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddButton("取用", icon.Android, action.Ajax("ai_playbook_batch_use",
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
			if err := biz.UseAiPlaybook(ids, userNameSub); err == nil {
				status = "取用成功，请前往[场景-场景列表]查看"
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
				status = "取用成功，请前往[场景-场景列表]查看"
			} else {
				status = fmt.Sprintf("取用失败：%s: %s", id, err)
			}
			return true, status, ""
		}))

	info.SetTable("ai_playbook").SetTitle("智能场景").SetDescription("智能场景")
	playbookTypeMsg := template.HTML("默认值为: 串行中断<br>串行中断: 场景内的数据用例若存在执行失败，失败后的数据用例不再执行<br>串行比较: 场景内的数据用例串行执行完成后，对各数据用例中输出的同名变量进行值比较，相等则通过<br>串行继续: 场景内的数据用例串行执行存在失败数据用例，失败后的数据用例继续执行<br>普通并发: 场景内的数据用例并发执行<br>并发比较: 场景内的数据用例并发执行完成后，对各数据用例中输出的同名变量进行值比较，相等则通过")
	dataHelp := template.HTML("关联数据必填")

	formList := aiPlaybook.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("场景描述", "name", db.Varchar, form.Text)
	formList.AddField("数据文件列表", "data_file_list", db.Text, form.TextArea).
		FieldHelpMsg(dataHelp)
	formList.AddField("场景类型", "scene_type", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "串行中断"},
			{Value: "2", Text: "串行比较"},
			{Value: "3", Text: "串行继续"},
			{Value: "4", Text: "普通并发"},
			{Value: "5", Text: "并发比较"},
		}).FieldDefault("1").FieldHelpMsg(playbookTypeMsg)
	formList.AddField("优先级", "priority", db.Int, form.Number)
	formList.AddField("生成来源", "source", db.Varchar, form.Text).
		FieldDisplayButCanNotEditWhenUpdate()
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
	formList.AddField("测试结果", "result", db.Varchar, form.Text)
	formList.AddField("失败原因", "fail_reason", db.Text, form.RichText)
	formList.AddField("所属产品", "product", db.Varchar, form.SelectSingle).
		FieldOptions(products)
	formList.AddField("创建人", "create_user", db.Varchar, form.Text).
		FieldDefault(userName).FieldHide()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("ai_playbook").SetTitle("智能场景").SetDescription("智能场景")

	detail := aiPlaybook.GetDetail()
	detail.AddField("唯一标识", "id", db.Int)
	detail.AddField("场景描述", "name", db.Varchar)
	detail.AddField("数据文件列表", "data_file_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetAiDataDetailLinkByDataStr(model.Value)
		})
	detail.AddField("最近数据文件", "last_file", db.Varchar)
	detail.AddField("场景类型", "scene_type", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "串行中断"
			}
			if model.Value == "2" {
				return "串行比较"
			}
			if model.Value == "3" {
				return "串行继续"
			}
			if model.Value == "4" {
				return "普通并发"
			}
			if model.Value == "5" {
				return "并发比较"
			}
			return "串行中断"
		})
	detail.AddField("优先级", "priority", db.Int)
	detail.AddField("执行次数", "run_time", db.Int)
	detail.AddField("测试结果", "result", db.Varchar)
	detail.AddField("失败原因", "fail_reason", db.Text)
	detail.AddField("备注", "remark", db.Longtext)
	detail.AddField("所属产品", "product", db.Varchar)
	detail.AddField("创建人", "create_user", db.Varchar)
	detail.AddField("创建时间", "created_at", db.Timestamp)
	detail.AddField("更新时间", "updated_at", db.Timestamp)
	detail.AddField("删除时间", "deleted_at", db.Timestamp).FieldHide()

	detail.SetTable("ai_playbook").SetTitle("智能场景详情").SetDescription("智能场景详情")

	return aiPlaybook
}
