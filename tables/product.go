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

func GetProductTable(ctx *context.Context) table.Table {

	product := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	apps := biz.GetApps()
	info := product.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)

	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).
		FieldFilterable()
	info.AddField(biz.T("common.product_name"), "product", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField(biz.T("common.ip"), "ip", db.Char).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldSortable()
	info.AddField(biz.T("common.protocol"), "protocol", db.Enum).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})
	info.AddField(biz.T("common.threading"), "threading", db.Enum).
		FieldFilterable().FieldDisplay(func(model types.FieldModel) interface{} {
		if model.Value == "yes" {
			return "yes"
		}
		if model.Value == "no" {
			return "no"
		}
		return "unknown"
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
	info.AddField(biz.T("common.testmode"), "testmode", db.Enum)
	info.AddField(biz.T("product.apps"), "apps", db.Longtext).FieldWidth(250)
	info.AddField(biz.T("common.env_type_label"), "env_type", db.Int).
		FieldDisplay(func(model types.FieldModel) interface{} {
			if model.Value == "1" {
				return biz.T("common.env_type._1")
			} else if model.Value == "2" {
				return biz.T("common.env_type._2")
			} else if model.Value == "3" {
				return biz.T("common.env_type._3")
			} else if model.Value == "4" {
				return biz.T("common.env_type._4")
			} else if model.Value == "5" {
				return biz.T("common.env_type._5")
			}
			return ""
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: biz.T("common.env_type._1")},
		{Value: "2", Text: biz.T("common.env_type._2")},
		{Value: "3", Text: biz.T("common.env_type._3")},
		{Value: "4", Text: biz.T("common.env_type._4")},
		{Value: "5", Text: biz.T("common.env_type._5")},
	})
	info.AddField(biz.T("product.private_app_prefix"), "private_app_prefix", db.Longtext).
		//FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().
		FieldHide()
	info.AddField(biz.T("product.private_parameter"), "private_parameter", db.Longtext).
		//FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().
		FieldHide()
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldHide()
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).
		FieldHide()
	info.AddButton(template2.HTML(biz.T("common.btn_copy")), icon.Android, action.Ajax("product_batch_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if idStr == "," {
				status = biz.T("common.btn_select_first")
				return false, status, ""
			}

			ids := strings.Split(idStr, ",")

			for _, id := range ids {
				if len(id) == 0 {
					continue
				}
				if err := biz.CopyProduct(id, userNameSub); err == nil {
					status = biz.T("product.copy_success")
				} else {
					status = fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_copy")), action.Ajax("product_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string
			user := auth.Auth(ctx)
			userNameSub := user.Name
			if err := biz.CopyProduct(id, userNameSub); err == nil {
				status = biz.T("product.copy_success")
			} else {
				status = fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err)
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("common.btn_report")), action.Jump("/admin/product_dashboard?id={{.Id}}"))

	info.AddButton(template2.HTML(biz.T("product.btn_refresh_report")), icon.Refresh, action.Ajax("product_batch_refresh_report",
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
					productName, err := biz.GetProductName(id)
					if err != nil {
						biz.Logger.Error("刷新产品报告失败[%s]: %s", id, err)
						continue
					}
					appName, err := biz.GetProductApps(id)
					if err != nil {
						biz.Logger.Error("获取产品应用失败[%s]: %s", id, err)
						continue
					}
					if err := pages.GetProductReportData(productName, appName, userNameSub); err != nil {
						biz.Logger.Error("刷新产品报告失败[%s]: %s", productName, err)
					}
				}
			}()
			status = biz.T("product.report_refreshing")
			return true, status, ""
		}))

	info.AddButton(template2.HTML(biz.T("product.btn_import_scene")), icon.Android, action.Ajax("import_batch_scene",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			idStr := ctx.FormValue("ids")
			var status string
			if idStr == "," {
				status = biz.T("product.import_select_first")
				return false, status, ""
			}
			ids := strings.Split(idStr, ",")
			for _, id := range ids {
				if len(id) == 0 {
					continue
				}

				if newCount, oldCount, err := biz.ImportPlaybookFromExcel(id); err == nil {
					status = fmt.Sprintf(biz.T("product.import_done")+": %d "+biz.T("product.added")+", %d "+biz.T("product.existing"), newCount, oldCount)
				} else {
					status = fmt.Sprintf(biz.T("product.sync_failed")+": %s, "+biz.T("product.added")+": %d, "+biz.T("product.existing")+": %d", err, newCount, oldCount)
					return false, status, ""
				}
			}
			return true, status, ""
		}))

	info.AddActionButton(template2.HTML(biz.T("product.btn_import_scene")), action.Ajax("import_scene",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			id := ctx.FormValue("id")
			var status string

			if newCount, oldCount, err := biz.ImportPlaybookFromExcel(id); err == nil {
				status = fmt.Sprintf(biz.T("product.import_done")+": %d "+biz.T("product.added")+", %d "+biz.T("product.existing"), newCount, oldCount)
			} else {
				status = fmt.Sprintf(biz.T("product.sync_failed")+": %s, "+biz.T("product.added")+": %d, "+biz.T("product.existing")+": %d", err, newCount, oldCount)
				return false, status, ""
			}

			return true, status, ""
		}))

	info.SetTable("product").SetTitle(biz.T("product.title")).SetDescription(biz.T("product.description"))
	helpMsg := template.HTML(biz.T("product.help_private_parameter"))
	helpPrefixMsg := template.HTML(biz.T("product.help_private_prefix"))

	formList := product.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).
		FieldDisableWhenCreate()
	formList.AddField(biz.T("common.product_name"), "product", db.Varchar, form.Text)
	formList.AddField(biz.T("common.ip"), "ip", db.Char, form.Ip)
	formList.AddField(biz.T("common.protocol"), "protocol", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: "http", Value: "http"},
			{Text: "https", Value: "https"},
		}).FieldDefault("http")
	formList.AddField(biz.T("common.threading"), "threading", db.Enum, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.yes"), Value: "yes"},
			{Text: biz.T("common.no"), Value: "no"},
		}).FieldDefault("no")
	formList.AddField(biz.T("common.thread_number"), "thread_number", db.Int, form.Number).FieldDefault("1")
	formList.AddField(biz.T("common.auth"), "auth", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.testmode"), "testmode", db.Enum, form.Radio).FieldOptions(types.FieldOptions{
		{Text: biz.T("product.testmode_custom"), Value: "custom"},
		{Text: biz.T("product.testmode_fuzzing"), Value: "fuzzing"},
		{Text: biz.T("product.testmode_all"), Value: "all"},
	}).FieldDefault("custom")
	formList.AddField(biz.T("product.apps"), "apps", db.Longtext, form.Select).
		FieldOptions(apps)
	formList.AddField(biz.T("common.env_type_label"), "env_type", db.Int, form.Radio).
		FieldOptions(types.FieldOptions{
			{Text: biz.T("common.env_type._1"), Value: "1"},
			{Text: biz.T("common.env_type._2"), Value: "2"},
			{Text: biz.T("common.env_type._3"), Value: "3"},
			{Text: biz.T("common.env_type._4"), Value: "4"},
			{Text: biz.T("common.env_type._5"), Value: "5"},
		}).FieldDefault("2")
	formList.AddField(biz.T("product.private_app_prefix"), "private_app_prefix", db.Longtext, form.TextArea).FieldHelpMsg(helpPrefixMsg)

	formList.AddField(biz.T("product.private_parameter"), "private_parameter", db.Longtext, form.TextArea).FieldHelpMsg(helpMsg)
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("product").SetTitle(biz.T("product.title")).SetDescription(biz.T("product.description"))

	return product
}
