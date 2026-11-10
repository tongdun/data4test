package tables

import (
	"data4test/biz"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

func orchestrationStatusText(value string) string {
	switch value {
	case biz.OrchestrationStatusRunning:
		return biz.T("common.running")
	case biz.OrchestrationStatusStopped:
		return biz.T("common.stopped")
	case biz.OrchestrationStatusFinished:
		return biz.T("common.finished")
	default:
		return biz.T("common.not_started")
	}
}

func orchestrationActionAllowed(ctx *context.Context, path, method string) bool {
	user := auth.Auth(ctx)
	if user.Id == 0 {
		return false
	}
	return auth.CheckPermissions(user, config.Url(path), method, url.Values{})
}

func GetTaskOrchestrationTable(ctx *context.Context) table.Table {

	orchestration := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))
	products := biz.GetProducts()
	userName := auth.Auth(ctx).Name
	dataLocale := biz.GetDataLocale(ctx.Cookie("data_locale"), biz.GetLocale())

	info := orchestration.GetInfo().HideFilterArea()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).FieldHide().FieldFilterable()
	info.AddField(biz.T("task_orchestration.orchestration_name"), "orchestration_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(180)
	info.AddField(biz.T("task_orchestration.step_count"), "task_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetOrchestrationTaskCount(model.Value)
		}).FieldWidth(80)
	info.AddField(biz.T("task_orchestration.product_list"), "product_list", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldFilterOptions(products).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetOrchestrationDisplayProducts(model.ID, model.Value, biz.GetOrchestrationTaskListById(model.ID))
		}).FieldWidth(220)
	info.AddField(biz.T("common.task_status"), "orchestration_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return orchestrationStatusText(model.Value)
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: biz.OrchestrationStatusRunning, Text: biz.T("common.running")},
		{Value: biz.OrchestrationStatusFinished, Text: biz.T("common.finished")},
		{Value: biz.OrchestrationStatusStopped, Text: biz.T("common.stopped")},
		{Value: biz.OrchestrationStatusNotStarted, Text: biz.T("common.not_started")},
	}).FieldWidth(100)
	info.AddField(biz.T("task_orchestration.last_at"), "last_at", db.Timestamp).FieldWidth(110).FieldSortable()
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldHide()
	info.AddField(biz.T("common.user_name"), "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(90)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(110).
		FieldFilterable(types.FilterType{FormType: form.DatetimeRange})
	info.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp).FieldHide()
	info.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp).FieldHide()

	// 批量立即执行
	info.AddButton(template2.HTML(biz.T("task_orchestration.btn_run")), icon.Android, action.Ajax("orchestration_batch_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			if !orchestrationActionAllowed(ctx, "/edit/task_orchestration", http.MethodPost) {
				return false, biz.T("error.permission_denied"), ""
			}
			idStr := ctx.FormValue("ids")
			if idStr == "," || len(strings.TrimSpace(idStr)) == 0 {
				return false, biz.T("error.orchestration_select_first"), ""
			}
			user := auth.Auth(ctx)
			if _, err := biz.RunOrchestrations(idStr, user.Name); err != nil {
				return false, fmt.Sprintf("%s: %s", biz.T("error.orchestration_run_fail"), err), ""
			}
			return true, biz.T("task_orchestration.msg_run_started"), ""
		}))

	// 批量复制
	info.AddButton(template2.HTML(biz.T("common.btn_copy")), icon.Android, action.Ajax("orchestration_batch_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			if !orchestrationActionAllowed(ctx, "/edit/task_orchestration", http.MethodPost) {
				return false, biz.T("error.permission_denied"), ""
			}
			idStr := ctx.FormValue("ids")
			if idStr == "," || len(strings.TrimSpace(idStr)) == 0 {
				return false, biz.T("error.orchestration_select_first"), ""
			}
			user := auth.Auth(ctx)
			for _, id := range strings.Split(idStr, ",") {
				id = strings.TrimSpace(id)
				if len(id) == 0 {
					continue
				}
				if err := biz.CopyTaskOrchestration(id, user.Name); err != nil {
					return false, fmt.Sprintf("%s: %s: %s", biz.T("common.copy_fail"), id, err), ""
				}
			}
			return true, biz.T("task_orchestration.copy_success"), ""
		}))

	// 行操作：立即执行
	info.AddActionButton(template2.HTML(biz.T("task_orchestration.btn_run")), action.Ajax("orchestration_run",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			if !orchestrationActionAllowed(ctx, "/edit/task_orchestration", http.MethodPost) {
				return false, biz.T("error.permission_denied"), ""
			}
			id := ctx.FormValue("id")
			user := auth.Auth(ctx)
			if err := biz.StartOrchestrationAsync(id, user.Name); err != nil {
				return false, fmt.Sprintf("%s: %s", biz.T("error.orchestration_run_fail"), err), ""
			}
			return true, biz.T("task_orchestration.msg_run_started"), ""
		}))

	// 行操作：复制
	info.AddActionButton(template2.HTML(biz.T("common.btn_copy")), action.Ajax("orchestration_copy",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			if !orchestrationActionAllowed(ctx, "/edit/task_orchestration", http.MethodPost) {
				return false, biz.T("error.permission_denied"), ""
			}
			id := ctx.FormValue("id")
			user := auth.Auth(ctx)
			if err := biz.CopyTaskOrchestration(id, user.Name); err != nil {
				return false, fmt.Sprintf("%s: %s", biz.T("common.copy_fail"), err), ""
			}
			return true, biz.T("task_orchestration.copy_success"), ""
		}))

	// 顶部：导入
	info.AddButton(template2.HTML(biz.T("task_orchestration.import_btn")), icon.FolderO, action.Jump("/admin/task_orchestration_import"))

	// 批量导出
	info.AddButton(template2.HTML(biz.T("task_orchestration.export_btn")), icon.Android, action.Ajax("orchestration_batch_export",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			if !orchestrationActionAllowed(ctx, "/export/task_orchestration", http.MethodPost) {
				return false, biz.T("error.permission_denied"), ""
			}
			idStr := ctx.FormValue("ids")
			if idStr == "," || len(strings.TrimSpace(idStr)) == 0 {
				return false, biz.T("task_orchestration.export_select_first"), ""
			}
			user := auth.Auth(ctx)
			fileName, err := biz.ExportTaskOrchestration(idStr, user.Name)
			if err != nil {
				return false, fmt.Sprintf("%s: %s", biz.T("task_orchestration.export_fail"), err), ""
			}
			downloadUrl := fmt.Sprintf("http://%s/admin/fm/download/download?path=/%s", ctx.Request.Host, fileName)
			return true, fmt.Sprintf("%s\n%s", biz.T("task_orchestration.export_success"), downloadUrl), ""
		}))

	// 行操作：导出
	info.AddActionButton(template2.HTML(biz.T("task_orchestration.export_btn")), action.Ajax("orchestration_export",
		func(ctx *context.Context) (success bool, msg string, data interface{}) {
			if !orchestrationActionAllowed(ctx, "/export/task_orchestration", http.MethodPost) {
				return false, biz.T("error.permission_denied"), ""
			}
			id := ctx.FormValue("id")
			user := auth.Auth(ctx)
			fileName, err := biz.ExportTaskOrchestration(id, user.Name)
			if err != nil {
				return false, fmt.Sprintf("%s: %s", biz.T("task_orchestration.export_fail"), err), ""
			}
			downloadUrl := fmt.Sprintf("http://%s/admin/fm/download/download?path=/%s", ctx.Request.Host, fileName)
			return true, fmt.Sprintf("%s\n%s", biz.T("task_orchestration.export_success"), downloadUrl), ""
		}))

	// 行操作：报告
	info.AddActionButton(template2.HTML(biz.T("common.btn_report")), action.Jump("/admin/orchestration_dashboard?id={{.Id}}"))

	info.SetTable("task_orchestration").SetTitle(biz.T("task_orchestration.title")).SetDescription(biz.T("task_orchestration.description"))

	// 删除：执行中禁止删除；删除为软删除并改写名称，保证同名可重建
	info.SetPreDeleteFn(func(ids []string) error {
		return biz.PreDeleteTaskOrchestrations(ids)
	})
	info.SetDeleteFn(func(ids []string) error {
		return biz.SoftDeleteTaskOrchestrations(ids)
	})

	formList := orchestration.GetForm()
	formList.AddField(biz.T("common.id"), "id", db.Int, form.Default).FieldDisableWhenCreate()
	formList.AddField(biz.T("task_orchestration.orchestration_name"), "orchestration_name", db.Varchar, form.Text)
	formList.AddField(biz.T("task_orchestration.task_list"), "task_list", db.Longtext, form.TextArea).
		FieldHelpMsg(template2.HTML(biz.T("task_orchestration.task_list_help"))).
		FieldMust()
	formList.AddField(biz.T("task_orchestration.product_list"), "product_list", db.Varchar, form.Select).
		FieldOptions(products)
	formList.AddField(biz.T("common.remark"), "remark", db.Longtext, form.TextArea)
	formList.AddField(biz.T("common.task_status"), "orchestration_status", db.Enum, form.Text).
		FieldHide().FieldDefault(biz.OrchestrationStatusNotStarted).FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("task_orchestration.last_at"), "last_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()
	formList.AddField(biz.T("common.user_name"), "user_name", db.Varchar, form.Text).
		FieldDefault(userName).FieldDisplayButCanNotEditWhenUpdate().FieldDisplayButCanNotEditWhenCreate()
	formList.AddField(biz.T("common.created_at"), "created_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenInsert().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp, form.Datetime).
		FieldHide().FieldNowWhenUpdate().FieldDisableWhenCreate()
	formList.AddField(biz.T("common.deleted_at"), "deleted_at", db.Timestamp, form.Datetime).
		FieldHide().FieldDisableWhenCreate().FieldDisableWhenUpdate()

	formList.SetTable("task_orchestration").SetTitle(biz.T("task_orchestration.title")).SetDescription(biz.T("task_orchestration.description"))

	formList.SetPreProcessFn(func(values form2.Values) form2.Values {
		if taskList, ok := values["task_list"]; ok && len(taskList) > 0 {
			values["task_list"] = []string{biz.NormalizeTaskList(taskList[0])}
		}
		if productList, ok := values["product_list"]; ok && len(productList) > 0 {
			values["product_list"] = []string{biz.NormalizeProductList(productList[0])}
		}
		return values
	})

	formList.SetPostValidator(func(values form2.Values) error {
		id := values.Get("id")
		name := values.Get("orchestration_name")
		if err := biz.ValidateOrchestrationName(name, id); err != nil {
			return err
		}
		if len(id) > 0 && biz.OrchestrationIsRunning(id) {
			return fmt.Errorf(biz.T("error.orchestration_edit_running"), name)
		}
		if err := biz.ValidateOrchestrationProducts(name, values.Get("product_list")); err != nil {
			return err
		}
		return biz.ValidateOrchestrationTaskList(name, values.Get("task_list"))
	})

	detail := orchestration.GetDetail()
	detail.AddField(biz.T("common.id"), "id", db.Int)
	detail.AddField(biz.T("task_orchestration.orchestration_name"), "orchestration_name", db.Varchar)
	detail.AddField(biz.T("task_orchestration.task_list"), "task_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return template.HTML(biz.BuildOrchestrationTaskListHTML(model.Value, dataLocale) +
				biz.BuildOrchestrationRecentReportsHTML(model.ID, 10))
		})
	detail.AddField(biz.T("task_orchestration.product_list"), "product_list", db.Varchar)
	detail.AddField(biz.T("common.task_status"), "orchestration_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return orchestrationStatusText(model.Value)
		})
	detail.AddField(biz.T("common.remark"), "remark", db.Longtext)
	detail.AddField(biz.T("task_orchestration.last_at"), "last_at", db.Timestamp)
	detail.AddField(biz.T("common.user_name"), "user_name", db.Varchar)
	detail.AddField(biz.T("common.created_at"), "created_at", db.Timestamp)
	detail.AddField(biz.T("common.updated_at"), "updated_at", db.Timestamp)
	detail.SetTable("task_orchestration").SetTitle(biz.T("task_orchestration.title")).SetDescription(biz.T("task_orchestration.detail_title"))

	return orchestration
}
