package tables

import (
	"data4test/biz"
	"data4test/pages"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"html/template"
	"strings"
)

func GetEnvConfigTable(ctx *context.Context) table.Table {

	envConfig := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	products := biz.GetProducts()
	info := envConfig.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)
	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldHide()
	info.AddField(biz.T("common.product"), "product", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(products)
	info.AddField(biz.T("common.app"), "app", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField(biz.T("common.ip"), "ip", db.Char).
		FieldSortable()
	info.AddField(biz.T("common.protocol"), "protocol", db.Enum)
	info.AddField(biz.T("common.prefix"), "prepath", db.Varchar)
	info.AddField(biz.T("common.threading"), "threading", db.Enum).
		FieldFilterable().FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "yes" {
			return biz.T("common.yes")
		}
		if model.Value == "no" {
			return biz.T("common.no")
		}
		return biz.T("common.no")
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "yes", Text: biz.T("common.yes")},
		{Value: "no", Text: biz.T("common.no")},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "yes", Text: biz.T("common.yes")},
		{Value: "no", Text: biz.T("common.no")},
	})

	info.AddField(biz.T("common.thread_number"), "thread_number", db.Int)

	info.AddField(biz.T("common.auth"), "auth", db.Longtext).FieldWidth(300).
		FieldHide()
	info.AddField(biz.T("common.testmode"), "testmode", db.Enum).
		FieldFilterable()
	info.AddField(biz.T("common.swagger_path"), "swagger_path", db.Text).
		FieldHide()
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(160).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(160).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()

	info.AddActionButton(template2.HTML(biz.T("common.btn_report")), action.Jump("/admin/app_dashboard?id={{.Id}}"))

	info.AddButton(template2.HTML(biz.T("product.btn_refresh_report")), icon.Refresh, action.Ajax("env_batch_refresh_report",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}
			user := auth.Auth(ctx)
			userNameSub := user.Name
			ids := strings.Split(idStr, ",")
			go func() {
				for _, id := range ids {
					if len(id) == 0 {
						continue
					}
					appName, err := biz.GetAppName(id)
					if err != nil {
						biz.Logger.Error("刷新应用报告失败[%s]: %s", id, err)
						continue
					}
					if err := pages.GetAppReportData(appName, userNameSub); err != nil {
						biz.Logger.Error("刷新应用报告失败[%s]: %s", appName, err)
					}
				}
			}()
			status = biz.T("product.report_refreshing")
			return true, status, ""
		}))

		info.AddActionButton(template2.HTML(biz.T("product.btn_refresh_report")), action.Ajax("env_refresh_report",
			func(ctx *context.Context) (success bool, msg string, data interface{}) {
				id := ctx.FormValue("id")
				var status string
				user := auth.Auth(ctx)
				userNameSub := user.Name
				appName, err := biz.GetAppName(id)
				if err != nil {
					status = fmt.Sprintf("刷新应用报告失败: %%s", err)
					return false, status, ""
				}
				go func() {
					if err := pages.GetAppReportData(appName, userNameSub); err != nil {
						biz.Logger.Error("刷新应用报告失败[%%s]: %%s", appName, err)
					}
				}()
				status = biz.T("product.report_refreshing")
				return true, status, ""
			}))

	info.AddButton(template.HTML(biz.T("env_config.btn_new_import")), icon.FolderO, action.PopUpWithCtxForm(action.PopUpData{
		Id:     "/api_define_import",
		Title:  biz.T("env_config.popup_title"),
		Width:  "900px",
		Height: "680px", // TextArea
	}, func(ctx *context.Context, panel *types.FormPanel) *types.FormPanel {
		ids := ctx.FormValue("ids")
		swaggerPath := biz.GetSwaggerPath(ids)
		panel.AddField(biz.T("env_config.field_selected_ids"), "ids", db.Varchar, form.Text).FieldDefault(ids).FieldHide()
		panel.AddField(biz.T("common.swagger_path"), "swagger_path", db.Varchar, form.Text).
			FieldDefault(swaggerPath).
			FieldHelpMsg(template.HTML(biz.T("env_config.help_swagger_or_upload")))
		panel.AddField(biz.T("env_config.field_upload_doc"), "upload_file", db.Varchar, form.Multifile).FieldOptionExt(map[string]interface{}{
			"maxFileCount": 1,
		}).FieldHelpMsg(template.HTML(biz.T("env_config.help_swagger_upload")))
		panel.EnableAjax(ctx.Response.Status, ctx.Response.Status)
		return panel
	}, "/api_define_import"))

	info.AddButton(template.HTML(biz.T("env_config.btn_old_import")), icon.Android, action.Ajax("autogeneration_batch",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("env_config.msg_select_first_import")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if failCount, err := biz.GetSwagger(id); err == nil {
					if failCount > 0 {
						status = fmt.Sprintf(biz.T("env_config.msg_import_done_with_fail"), failCount)
					} else {
						status = biz.T("env_config.msg_import_done")
					}
					_ = biz.UpdateApiChangeByAppId(id)
				} else {
					status = fmt.Sprintf(biz.T("env_config.msg_import_fail"), err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("env_config.btn_old_import_single")), action.Ajax("autogeneration",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if failCount, err := biz.GetSwagger(id); err == nil {
				if failCount > 0 {
					status = fmt.Sprintf(biz.T("env_config.msg_import_done_with_fail"), failCount)
				} else {
					status = biz.T("env_config.msg_import_done")
				}
				_ = biz.UpdateApiChangeByAppId(id)
			} else {
				status = fmt.Sprintf(biz.T("env_config.msg_import_fail"), err)
			}

			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("env_config.btn_update_status")), icon.Android, action.Ajax("autoupdatestatus_batch",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("env_config.msg_select_first_update")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.UpdateApiAutoStatus(id); err == nil {
					status = biz.T("env_config.msg_update_done")
				} else {
					status = fmt.Sprintf(biz.T("env_config.msg_update_fail"), err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("env_config.btn_update_status")), action.Ajax("autoupdatestatus",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			if err := biz.UpdateApiAutoStatus(id); err == nil {
				status = biz.T("env_config.msg_update_done")
			} else {
				status = fmt.Sprintf(biz.T("env_config.msg_update_fail"), err)
			}

			return true, status, ""
		}))

	info.AddButton(template.HTML(biz.T("env_config.btn_update_auth")), icon.Android, action.Ajax("update_batch_auth",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			//idStr := ctx.FormValue("ids")
			var status string
			status = biz.T("env_config.msg_not_implemented")
			//if idStr == "," {
			//	//status = biz.T("env_config.msg_select_first_update")
			//	status = biz.T("env_config.msg_not_implemented")
			//	return false, status, ""
			//}
			//ids := strings.Split(idStr, ",")
			//for _, id := range ids {
			//	if len(id) == 0 {
			//		//status = "更新完成，请开始接口测试之旅"
			//		status = biz.T("env_config.msg_not_implemented")
			//		continue
			//	}
			//	if err := biz.GetAuth(id); err == nil {
			//		status = "更新完成，请开始接口测试之旅"
			//	} else {
			//		status = fmt.Sprintf(biz.T("env_config.msg_update_fail"), err)
			//		return false, status, ""
			//	}
			//}
			return false, status, ""
		}))

	info.AddActionButton(template.HTML(biz.T("env_config.btn_update_auth")), action.Ajax("update_auth",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			//id := ctx.FormValue("id")
			var status string
			status = biz.T("env_config.msg_not_implemented")
			//if err := biz.GetAuth(id); err == nil {
			//	status = "更新完成，请开始接口测试之旅"
			//} else {
			//	status = fmt.Sprintf(biz.T("env_config.msg_update_fail"), err)
			//}

			return false, status, ""
		}))

	info.AddSelectBox(biz.T("common.product"), products, action.FieldFilter("product"))
	apps := biz.GetApps()
	info.AddSelectBox(biz.T("common.app"), apps, action.FieldFilter("app"))
	info.SetTable("env_config").SetTitle(biz.T("env_config.title")).SetDescription(biz.T("env_config.description"))

	formList := envConfig.GetForm()

	helpMsg := template.HTML(biz.T("env_config.help_swagger_route"))
	authHelpMsg := template.HTML(biz.T("env_config.help_auth_header"))
	appHelpMsg := template.HTML(biz.T("env_config.help_app_name"))
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate().
		FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.product"), "product", db.Varchar, form.SelectSingle).
		FieldOptions(products)
	formList.AddField(biz.T("common.app"), "app", db.Varchar, form.Text).
		FieldDisableWhenUpdate().
		FieldHelpMsg(appHelpMsg)
	formList.AddField(biz.T("common.ip"), "ip", db.Char, form.Ip)
	formList.AddField(biz.T("common.protocol"), "protocol", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "http", Value: "http"},
			{Text: "https", Value: "https"},
		}).FieldDefault("http")

	formList.AddField(biz.T("common.prefix"), "prepath", db.Varchar, form.Text)
	formList.AddField(biz.T("common.threading"), "threading", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.yes"), Value: "yes"},
			{Text: biz.T("common.no"), Value: "no"},
		}).FieldDefault("no")
	formList.AddField(biz.T("common.thread_number"), "thread_number", db.Int, form.Number).FieldDefault("1")
	formList.AddField(biz.T("common.auth"), "auth", db.Longtext, form.TextArea).
		FieldHelpMsg(authHelpMsg)
	formList.AddField(biz.T("common.testmode"), "testmode", db.Enum, form.Radio).FieldOptions(types.FieldOptions{
		{Text: biz.T("product.testmode_custom"), Value: "custom"},
		{Text: biz.T("product.testmode_fuzzing"), Value: "fuzzing"},
		{Text: biz.T("product.testmode_all"), Value: "all"},
	}).FieldDefault("custom")
	formList.AddField(biz.T("common.swagger_path"), "swagger_path", db.Varchar, form.Text).
		FieldHelpMsg(helpMsg)
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("env_config").SetTitle(biz.T("env_config.title")).SetDescription(biz.T("env_config.description"))

	return envConfig
}
