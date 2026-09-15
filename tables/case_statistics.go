package tables

import (
	"data4test/biz"
	"fmt"
	"html/template"
	"strconv"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetCaseStatisticsTable(ctx *context.Context) table.Table {
	caseStat := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	userName := auth.Auth(ctx).Name

	versions := make([]types.FieldOption, 0, 4)
	for _, v := range biz.GetIntroVersions() {
		versions = append(versions, types.FieldOption{Value: v, Text: v})
	}

	products := make([]types.FieldOption, 0, 4)
	for _, v := range biz.GetCaseProductList() {
		products = append(products, types.FieldOption{Value: v, Text: v})
	}

	info := caseStat.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).FieldWidth(100)
	info.AddField(biz.T("case_statistics.name"), "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("case_statistics.intro_versions"), "intro_versions", db.Varchar).
		FieldWidth(160)
	info.AddField(biz.T("case_statistics.product"), "product", db.Varchar).
		FieldWidth(160)
	info.AddField(biz.T("case_statistics.status"), "status", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			switch model.Value {
			case "generating":
				return "<span style=\"color:orange;font-weight:bold\">" + biz.T("case_statistics.status_generating") + "</span>"
			case "finished":
				return "<span style=\"color:green;font-weight:bold\">" + biz.T("case_statistics.status_finished") + "</span>"
			case "failed":
				return "<span style=\"color:red;font-weight:bold\">" + biz.T("case_statistics.status_failed") + "</span>"
			}
			return biz.T("case_statistics.status_none")
		}).FieldWidth(100)
	info.AddField(biz.T("common.user_name"), "creator", db.Varchar).
		FieldWidth(100)
	info.AddField(biz.T("common.remark"), "remark", db.Varchar).
		FieldHide()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable()
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	// 开始统计（异步触发）
	info.AddActionButton(template.HTML(biz.T("case_statistics.btn_start")), action.Ajax("case_statistics_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return false, biz.T("case_statistics.invalid_id"), ""
			}
			go biz.GenerateCaseStatisticsReport(id, auth.Auth(ctx).Name)
			return true, biz.T("case_statistics.stat_running_background"), ""
		}))

	// 查看报告（最新）
	info.AddActionButton(template.HTML(biz.T("case_statistics.btn_report")), action.Jump("/admin/case_statistics_report?id={{.Id}}"))

	// i18n同步（批量）：导出选中定义的叶子模块用例原始数据 → mgmt/i18n/i18n_case/<模块>/zh-CN.yaml
	info.AddButton(template.HTML(biz.T("case_statistics.btn_i18n_sync")), icon.Android, action.Ajax("case_statistics_i18n_sync",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			if idStr == "," {
				return false, biz.T("common.btn_select_first"), ""
			}
			mc, cc, err := biz.SyncCaseI18nFromStatistics(idStr)
			if err != nil {
				return false, fmt.Sprintf("%s: %v", biz.T("common.operate_fail"), err), ""
			}
			return true, fmt.Sprintf(biz.T("case_statistics.i18n_sync_done"), mc, cc), ""
		}))

	// 导出Excel（批量）：选中定义覆盖的全量用例按模块导出，一个模块一个 xlsx，打包为单个 zip
	info.AddButton(template.HTML(biz.T("case_statistics.btn_export_excel")), icon.Download, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/case_statistics_export_excel",
		Title:  biz.T("case_statistics.title_export_excel"),
		Width:  "900px",
		Height: "480px",
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		templates := biz.GetTestCaseExportTemplates()
		languages := biz.GetSupportLanguages()
		panel.AddField(biz.T("common.selected_ids"), "ids", db.Varchar, form.Text).
			FieldDefault(ids).FieldDisplayButCanNotEditWhenCreate().
			FieldHelpMsg(template.HTML(biz.T("common.help_select_or_filter")))
		panel.AddField(biz.T("test_case.export_template"), "template", db.Varchar, form.SelectSingle).
			FieldOptions(templates).FieldDefault(templates[0].Value)
		panel.AddField(biz.T("test_case.export_lang"), "lang", db.Varchar, form.SelectSingle).
			FieldOptions(languages).FieldDefault(biz.GetLocale())
		panel.AddField(biz.T("test_case.screenshot_mode"), "screenshot_mode", db.Varchar, form.SelectSingle).
			FieldOptions(types.FieldOptions{
				{Value: "", Text: biz.T("test_case.screenshot_mode_none")},
				{Value: "path", Text: biz.T("test_case.screenshot_mode_path")},
				{Value: "embed", Text: biz.T("test_case.screenshot_mode_embed")},
			}).FieldDefault("")
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/case_statistics_export_excel"))

	info.SetTable("case_statistics").SetTitle(biz.T("case_statistics.title")).SetDescription(biz.T("case_statistics.description"))

	formList := caseStat.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("case_statistics.name"), "name", db.Varchar, form.Text)
	formList.AddField(biz.T("case_statistics.definition"), "definition", db.Longtext, form.TextArea).
		FieldDefault(biz.T("case_statistics.sample_definition")).
		FieldHelpMsg(template.HTML(biz.T("case_statistics.help_definition")))
	formList.AddField(biz.T("case_statistics.intro_versions"), "intro_versions", db.Varchar, form.Select).
		FieldOptions(versions).
		FieldHelpMsg(template.HTML(biz.T("case_statistics.help_intro_versions")))
	formList.AddField(biz.T("case_statistics.product"), "product", db.Varchar, form.Select).
		FieldOptions(products).
		FieldHelpMsg(template.HTML(biz.T("case_statistics.help_product")))
	formList.AddField(biz.T("common.remark"), "remark", db.Varchar, form.TextArea)
	formList.AddField(biz.T("common.user_name"), "creator", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("case_statistics").SetTitle(biz.T("case_statistics.title")).SetDescription(biz.T("case_statistics.description"))

	detail := caseStat.GetDetail()
	detail.AddField(biz.T("case_statistics.name"), "name", db.Varchar)
	detail.AddField(biz.T("case_statistics.intro_versions"), "intro_versions", db.Varchar)
	detail.AddField(biz.T("case_statistics.product"), "product", db.Varchar)
	detail.AddField(biz.T("case_statistics.definition"), "definition", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return fmt.Sprintf("<pre style='white-space:pre-wrap'>%s</pre>", template.HTMLEscapeString(model.Value))
		})
	detail.AddField(biz.T("case_statistics.status"), "status", db.Varchar)
	detail.AddField(biz.T("common.user_name"), "creator", db.Varchar)
	detail.AddField(biz.T("common.remark"), "remark", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)

	detail.SetTable("case_statistics").SetTitle(biz.T("case_statistics.title")).SetDescription(biz.T("case_statistics.description"))

	return caseStat
}
