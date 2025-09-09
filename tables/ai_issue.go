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
)

func GetAiIssueTable(ctx *context.Context) table.Table {

	aiIssue := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := aiIssue.GetInfo().HideFilterArea()
	user := auth.Auth(ctx)
	userName := user.Name
	products := biz.GetProducts()
	aiPlatforms := biz.GetAiCreatePlatform()
	issueSourceTypes := types.FieldOptions{
		{Value: "1", Text: "数据"},
		{Value: "2", Text: "场景"},
		{Value: "3", Text: "手动输入"},
	}
	confirmStatusTypes := types.FieldOptions{
		{Value: "1", Text: "BUG"},
		{Value: "2", Text: "优化"},
		{Value: "3", Text: "误判"},
	}
	solveStatusTypes := types.FieldOptions{
		{Value: "1", Text: "创建"},
		{Value: "2", Text: "解决中"},
		{Value: "3", Text: "修复完成"},
		{Value: "4", Text: "验证完成"},
		{Value: "5", Text: "不处理"},
	}
	againTestTypes := types.FieldOptions{
		{Value: "0", Text: "失败"},
		{Value: "1", Text: "成功"},
		{Value: "2", Text: "未知"},
	}

	info.AddField("自增主键", "id", db.Int)
	info.AddField("问题名称", "issue_name", db.Varchar)
	info.AddField("问题级别", "issue_level", db.Varchar)
	info.AddField("问题来源", "issue_source", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "数据"
			} else if model.Value == "2" {
				return "场景"
			} else if model.Value == "3" {
				return "手动输入"
			}
			return "数据"
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(issueSourceTypes)

	info.AddField("来源名称", "source_name", db.Varchar).
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
	info.AddField("分析平台", "source", db.Varchar)
	info.AddField("请求数据", "request_data", db.Text).
		FieldHide()
	info.AddField("返回数据", "response_data", db.Text).
		FieldHide()
	info.AddField("问题详情", "issue_detail", db.Text).
		FieldHide()
	info.AddField("确认状态", "confirm_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "BUG"
			} else if model.Value == "2" {
				return "优化"
			} else if model.Value == "3" {
				return "误判"
			}
			return "BUG"
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(confirmStatusTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(confirmStatusTypes)

	info.AddField("问题原因推测", "root_cause", db.Text).
		FieldHide()
	info.AddField("影响范围分析", "impact_scope_analysis", db.Text).
		FieldHide()
	info.AddField("受影响的场景推测", "impact_playbook", db.Text).
		FieldHide()
	info.AddField("受影响的数据推测", "impact_data", db.Text).
		FieldHide()
	info.AddField("解决状态", "resolution_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return "创建"
			} else if model.Value == "2" {
				return "解决中"
			} else if model.Value == "3" {
				return "修复完成"
			} else if model.Value == "4" {
				return "验证完成"
			} else if model.Value == "5" {
				return "不处理"
			}
			return "创建"
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(solveStatusTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(solveStatusTypes)

	info.AddField("回归测试结果", "again_test_result", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return "失败"
			} else if model.Value == "1" {
				return "成功"
			} else if model.Value == "2" {
				return "未知"
			}
			return "未知"
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(againTestTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(againTestTypes)

	info.AddField("受影响模块回归测试结果", "impact_test_result", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "0" {
				return "失败"
			} else if model.Value == "1" {
				return "成功"
			} else if model.Value == "2" {
				return "未知"
			}
			return "未知"
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(againTestTypes).
		FieldEditAble(editType.Select).
		FieldEditOptions(againTestTypes)

	info.AddField("关联产品", "product_list", db.Text).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(products).
		FieldEditAble(editType.Select).
		FieldEditOptions(products)
	info.AddField("创建人", "create_user", db.Varchar)
	info.AddField("修改人", "modify_user", db.Varchar)
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
		Id:     "/ai_issue_import",
		Title:  "AI导入Issue",
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		products := biz.GetProducts()
		panel.AddField("分析平台", "create_platform", db.Varchar, form.SelectSingle).
			FieldOptions(aiPlatforms).FieldDefault(aiPlatforms[0].Value)
		panel.AddField("所属产品", "product", db.Varchar, form.SelectSingle).
			FieldOptions(products).FieldDefault(products[0].Value)
		panel.AddField("会话ID", "conversation_id", db.Varchar, form.Text).
			FieldHelpMsg("生成平台的会话ID，会话ID和原生回复二选一")
		panel.AddField("原生回复", "raw_reply", db.Varchar, form.TextArea).
			FieldHelpMsg("生成平台的回复原文")

		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)

		return panel
	}, "/ai_issue_import"))

	info.AddButton("回归测试", icon.Android, action.Ajax("source_again_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = "请先选择数据再测试"
				return false, status, ""
			}
			if err := biz.SourceAgainTest(idStr); err == nil {
				status = "测试完成，请刷新列表查看测试结果"
			} else {
				status = fmt.Sprintf("测试失败: %s", err)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.SetTable("ai_issue").SetTitle("智能分析").SetDescription("智能分析")

	formList := aiIssue.GetForm()
	formList.AddField("自增主键", "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField("问题名称", "issue_name", db.Varchar, form.Text)
	formList.AddField("问题级别", "issue_level", db.Varchar, form.Text)
	formList.AddField("问题来源", "issue_source", db.Enum, form.Radio).
		FieldOptions(solveStatusTypes).
		FieldDefault("3")
	formList.AddField("来源名称", "source_name", db.Varchar, form.Text)
	formList.AddField("分析平台", "source", db.Varchar, form.SelectSingle).
		FieldOptions(aiPlatforms)
	formList.AddField("请求数据", "request_data", db.Text, form.TextArea)
	formList.AddField("返回数据", "response_data", db.Text, form.TextArea)
	formList.AddField("问题详情", "issue_detail", db.Text, form.TextArea)
	formList.AddField("确认状态", "confirm_status", db.Enum, form.Radio).
		FieldOptions(confirmStatusTypes).
		FieldDefault("1")
	formList.AddField("问题原因推测", "root_cause", db.Text, form.TextArea)
	formList.AddField("影响范围分析", "impact_scope_analysis", db.Text, form.TextArea)
	formList.AddField("受影响的场景推测", "impact_playbook", db.Text, form.TextArea)
	formList.AddField("受影响的数据推测", "impact_data", db.Text, form.TextArea)
	formList.AddField("解决状态", "resolution_status", db.Enum, form.Radio).
		FieldOptions(solveStatusTypes).
		FieldDefault("1")
	formList.AddField("回归测试结果", "again_test_result", db.Enum, form.Radio).
		FieldOptions(againTestTypes).
		FieldDefault("2")
	formList.AddField("受影响模块回归测试结果", "impact_test_result", db.Enum, form.Radio).
		FieldOptions(againTestTypes).
		FieldDefault("2")
	formList.AddField("关联产品", "product_list", db.Text, form.SelectSingle).
		FieldOptions(products)
	formList.AddField("创建人", "create_user", db.Varchar, form.Text).
		FieldDefault(userName).
		FieldHideWhenCreate().
		FieldDisableWhenCreate()
	formList.AddField("修改人", "modify_user", db.Varchar, form.Text).
		FieldDefault(userName).
		FieldHideWhenCreate().
		FieldHideWhenUpdate().
		FieldDisableWhenCreate()
	formList.AddField("创建时间", "created_at", db.Timestamp, form.Datetime).
		FieldHide().
		FieldNowWhenInsert().
		FieldDisableWhenCreate()
	formList.AddField("更新时间", "updated_at", db.Timestamp, form.Datetime).
		FieldHide().
		FieldNowWhenUpdate().
		FieldDisableWhenCreate()
	formList.AddField("删除时间", "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()

	formList.SetTable("ai_issue").SetTitle("智能分析").SetDescription("智能分析")

	return aiIssue
}
