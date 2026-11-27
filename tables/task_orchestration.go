package tables

import (
	"data4test/biz"
	"data4test/models"
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
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

	// 筛选区域默认展开展示（不使用 HideFilterArea 折叠）
	info := orchestration.GetInfo()
	info.SetFilterFormHeadWidth(4)
	info.SetFilterFormInputWidth(8)
	info.SetFilterFormLayout(form.LayoutThreeCol)

	info.AddField(biz.T("common.id"), "id", db.Int).FieldHide()
	info.AddField(biz.T("task_orchestration.orchestration_name"), "orchestration_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(200)
	info.AddField(biz.T("task_orchestration.step_count"), "task_list", db.Longtext).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return fmt.Sprintf("%d", biz.GetOrchestrationTaskCount(model.Value))
		}).FieldWidth(100)
	info.AddField(biz.T("task_orchestration.product_list"), "product_list", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldFilterOptions(products).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return biz.GetOrchestrationDisplayProducts(model.ID, model.Value, biz.GetOrchestrationTaskListById(model.ID))
		}).FieldWidth(200)
	info.AddField(biz.T("common.task_status"), "orchestration_status", db.Enum).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return orchestrationStatusText(model.Value)
		}).FieldFilterable(types.FilterType{FormType: form.Select}).FieldFilterOptions(types.FieldOptions{
		{Value: biz.OrchestrationStatusRunning, Text: biz.T("common.running")},
		{Value: biz.OrchestrationStatusFinished, Text: biz.T("common.finished")},
		{Value: biz.OrchestrationStatusNotStarted, Text: biz.T("common.not_started")},
	}).FieldWidth(110)
	info.AddField(biz.T("task_orchestration.last_at"), "last_at", db.Timestamp).FieldWidth(140).FieldSortable()
	info.AddField(biz.T("common.remark"), "remark", db.Longtext).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldHide()
	info.AddField(biz.T("common.user_name"), "user_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike}).
		FieldTrimSpace().FieldWidth(100)
	info.AddField(biz.T("common.created_at"), "created_at", db.Timestamp).
		FieldSortable().FieldWidth(120).
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

	// 列表自管查询：产品配置筛选需按「产品配置」列的展示口径匹配（未配置编排级产品时取任务产品并集），
	// GoAdmin 的字段筛选只能比对 product_list 原始列，会漏掉从任务继承产品的编排，故用 GetDataFn 统一实现。
	info.SetGetDataFn(orchestrationListData)

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

// orchestrationSortFields 列表排序字段白名单，避免 __sort 参数直接拼进 SQL
var orchestrationSortFields = map[string]bool{
	"id": true, "orchestration_name": true, "product_list": true, "orchestration_status": true,
	"last_at": true, "user_name": true, "created_at": true, "updated_at": true,
}

// orchestrationFilteredQuery 按列表筛选项构建查询，口径与 GoAdmin 原生筛选一致（模糊/等值/时间范围/软删除）
func orchestrationFilteredQuery(params parameter.Parameters) *gorm.DB {
	q := models.Orm.Table("task_orchestration").Where("deleted_at is null")
	q = applyOrchestrationEqualFilter(q, params, "id")
	q = applyOrchestrationLikeFilter(q, params, "orchestration_name")
	// 任务状态为多选筛选：多值走 in (...)，与 GoAdmin 原生筛选口径一致
	q = applyOrchestrationEqualFilter(q, params, "orchestration_status")
	q = applyOrchestrationLikeFilter(q, params, "user_name")
	q = applyOrchestrationLikeFilter(q, params, "remark")
	if v := strings.TrimSpace(params.GetFilterFieldValueStart("created_at")); v != "" {
		q = q.Where("created_at >= ?", v)
	}
	if v := strings.TrimSpace(params.GetFilterFieldValueEnd("created_at")); v != "" {
		q = q.Where("created_at <= ?", v)
	}
	// 产品配置：按列表展示口径匹配 —— 编排级产品优先，未配置时取步骤任务产品并集；多值时任一命中即可
	if values := orchestrationFilterValues(params, "product_list"); len(values) > 0 {
		seen := make(map[string]bool)
		var ids []string
		for _, keyword := range values {
			matched, err := biz.SearchOrchestrationIdsByProduct(keyword)
			if err != nil {
				biz.Logger.Error("任务编排产品筛选失败: %v", err)
				return q.Where("1 = 0")
			}
			for _, id := range matched {
				if !seen[id] {
					seen[id] = true
					ids = append(ids, id)
				}
			}
		}
		if len(ids) == 0 {
			return q.Where("1 = 0")
		}
		q = q.Where("id in (?)", ids)
	}
	return q
}

// orchestrationFilterValues 读取筛选值并过滤空白；多选控件（如任务状态）会给出多个值
func orchestrationFilterValues(params parameter.Parameters, field string) []string {
	raw := params.GetFieldValues(field)
	values := make([]string, 0, len(raw))
	for _, value := range raw {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}

// applyOrchestrationEqualFilter 等值筛选：单值 =，多值 in（对应 GoAdmin 原生筛选语义）
func applyOrchestrationEqualFilter(q *gorm.DB, params parameter.Parameters, field string) *gorm.DB {
	values := orchestrationFilterValues(params, field)
	switch len(values) {
	case 0:
		return q
	case 1:
		return q.Where(field+" = ?", values[0])
	default:
		return q.Where(field+" in (?)", values)
	}
}

// applyOrchestrationLikeFilter 模糊筛选：单值 like；多值时任一命中即可
func applyOrchestrationLikeFilter(q *gorm.DB, params parameter.Parameters, field string) *gorm.DB {
	values := orchestrationFilterValues(params, field)
	switch len(values) {
	case 0:
		return q
	case 1:
		return q.Where(field+" like ?", "%"+values[0]+"%")
	default:
		conditions := make([]string, 0, len(values))
		args := make([]interface{}, 0, len(values))
		for _, value := range values {
			conditions = append(conditions, field+" like ?")
			args = append(args, "%"+value+"%")
		}
		return q.Where("("+strings.Join(conditions, " or ")+")", args...)
	}
}

// orchestrationListData 任务编排列表数据：筛选、排序、分页与 GoAdmin 原生列表行为保持一致
func orchestrationListData(params parameter.Parameters) ([]map[string]interface{}, int) {
	var total int64
	if err := orchestrationFilteredQuery(params).Count(&total).Error; err != nil {
		biz.Logger.Error("查询任务编排列表失败: %v", err)
		return nil, 0
	}

	sortField := params.SortField
	if !orchestrationSortFields[sortField] {
		sortField = "id"
	}
	sortType := strings.ToLower(params.SortType)
	if sortType != "asc" && sortType != "desc" {
		sortType = "desc"
	}

	q := orchestrationFilteredQuery(params).Order(sortField + " " + sortType)
	if !params.IsAll() {
		page, size := params.PageInt, params.PageSizeInt
		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 10
		}
		q = q.Offset((page - 1) * size).Limit(size)
	}

	var rows []biz.DbTaskOrchestration
	if err := q.Find(&rows).Error; err != nil {
		biz.Logger.Error("查询任务编排列表失败: %v", err)
		return nil, 0
	}

	data := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		id, _ := strconv.ParseInt(row.Id, 10, 64)
		data = append(data, map[string]interface{}{
			"id":                   id,
			"orchestration_name":   row.OrchestrationName,
			"task_list":            row.TaskList,
			"product_list":         row.ProductList,
			"orchestration_status": row.OrchestrationStatus,
			"last_at":              row.LastAt,
			"remark":               row.Remark,
			"user_name":            row.UserName,
			"created_at":           row.CreatedAt,
			"updated_at":           row.UpdatedAt,
			"deleted_at":           row.DeletedAt,
		})
	}
	return data, int(total)
}
