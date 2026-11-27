package biz

import (
	"encoding/json"
	"html"

	"data4test/models"
	"fmt"
	"github.com/GoAdminGroup/go-admin/template/types"
	"strings"
	"time"
)

// ParseTaskList 解析 task_list：按行拆分、去空行、去首尾空格，保留行序
func ParseTaskList(taskList string) []string {
	var names []string
	for _, line := range strings.Split(strings.ReplaceAll(taskList, "\r\n", "\n"), "\n") {
		name := strings.TrimSpace(line)
		if len(name) == 0 {
			continue
		}
		names = append(names, name)
	}
	return names
}

// NormalizeTaskList 归一化 task_list：统一换行、去空行、去首尾空格
func NormalizeTaskList(taskList string) string {
	return strings.Join(ParseTaskList(taskList), "\n")
}

// GetOrchestrationTaskCount 步骤数
func GetOrchestrationTaskCount(taskList string) int {
	return len(ParseTaskList(taskList))
}

// GetTaskOrchestrationById 按 id 查询编排（仅未删除）
func GetTaskOrchestrationById(id string) (orch DbTaskOrchestration, err error) {
	if len(strings.TrimSpace(id)) == 0 {
		err = fmt.Errorf(T("error.orchestration_not_found"), id)
		return
	}
	models.Orm.Table("task_orchestration").
		Where("id = ? and deleted_at is null", id).
		Find(&orch)
	if len(orch.Id) == 0 {
		err = fmt.Errorf(T("error.orchestration_not_found"), id)
	}
	return
}

// GetOrchestrationTaskListById 查询编排的 task_list
func GetOrchestrationTaskListById(id string) string {
	var orch DbTaskOrchestration
	models.Orm.Table("task_orchestration").Select("task_list").Where("id = ? and deleted_at is null", id).Find(&orch)
	return orch.TaskList
}

// ValidateOrchestrationName 校验编排名称唯一（excludeId 非空时排除自身）
func ValidateOrchestrationName(name, excludeId string) error {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return fmt.Errorf(T("error.orchestration_name_required"))
	}
	query := models.Orm.Table("task_orchestration").Where("orchestration_name = ? and deleted_at is null", name)
	if len(excludeId) > 0 {
		query = query.Where("id != ?", excludeId)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf(T("error.orchestration_name_exists"), name)
	}
	return nil
}

// ValidateOrchestrationTaskList 校验关联任务：非空、无重复、全部存在且未删除
func ValidateOrchestrationTaskList(orchestrationName, taskList string) error {
	names := ParseTaskList(taskList)
	if len(names) == 0 {
		return fmt.Errorf(T("error.orchestration_no_tasks"), orchestrationName)
	}

	seen := make(map[string]bool, len(names))
	for _, name := range names {
		if seen[name] {
			return fmt.Errorf(T("error.orchestration_task_duplicate"), orchestrationName, name)
		}
		seen[name] = true
	}

	var rows []DbSchedule
	if err := models.Orm.Table("schedule").
		Select("task_name").
		Where("task_name in (?) and deleted_at is null", names).
		Find(&rows).Error; err != nil {
		return err
	}
	existing := make(map[string]bool, len(rows))
	for _, row := range rows {
		existing[row.TaskName] = true
	}
	for _, name := range names {
		if !existing[name] {
			return fmt.Errorf(T("error.orchestration_task_missing"), orchestrationName, name)
		}
	}
	return nil
}

// GetOrchestrationTasks 按 task_list 行序解析任务，任一任务缺失即返回错误
func GetOrchestrationTasks(orch *DbTaskOrchestration) (tasks []DbSchedule, err error) {
	names := ParseTaskList(orch.TaskList)
	if len(names) == 0 {
		err = fmt.Errorf(T("error.orchestration_no_tasks"), orch.OrchestrationName)
		return
	}

	var rows []DbSchedule
	if err = models.Orm.Table("schedule").
		Where("task_name in (?) and deleted_at is null", names).
		Find(&rows).Error; err != nil {
		return
	}
	byName := make(map[string]DbSchedule, len(rows))
	for _, row := range rows {
		byName[row.TaskName] = row
	}

	for _, name := range names {
		task, ok := byName[name]
		if !ok {
			err = fmt.Errorf(T("error.orchestration_task_missing"), orch.OrchestrationName, name)
			return
		}
		tasks = append(tasks, task)
	}
	return
}

// OrchestrationIsRunning 编排是否正在执行（以库状态为准，跨进程有效）
func OrchestrationIsRunning(orchestrationId string) bool {
	var row DbTaskOrchestration
	models.Orm.Table("task_orchestration").
		Select("orchestration_status").
		Where("id = ?", orchestrationId).
		Find(&row)
	return row.OrchestrationStatus == OrchestrationStatusRunning
}

// CheckOrchestrationStartable 执行前校验：存在、未在执行、任务全部可解析
func CheckOrchestrationStartable(orchestrationId string) error {
	orch, err := GetTaskOrchestrationById(orchestrationId)
	if err != nil {
		return err
	}
	if orch.OrchestrationStatus == OrchestrationStatusRunning {
		return fmt.Errorf(T("error.orchestration_running"), orch.OrchestrationName)
	}
	_, err = GetOrchestrationTasks(&orch)
	return err
}

// CheckOrchestrationRunnable 执行前校验：编排存在、任务非空且全部可解析
func CheckOrchestrationRunnable(orchestrationId string) error {
	orch, err := GetTaskOrchestrationById(orchestrationId)
	if err != nil {
		return err
	}
	_, err = GetOrchestrationTasks(&orch)
	return err
}

// GetOrchestrationEffectiveProducts 本轮实际使用的产品：编排级优先，否则取步骤任务产品并集（去重）
func GetOrchestrationEffectiveProducts(orch *DbTaskOrchestration, tasks []DbSchedule) string {
	if len(strings.TrimSpace(orch.ProductList)) > 0 {
		return strings.TrimSpace(orch.ProductList)
	}

	productSet := make(map[string]bool)
	var products []string
	for _, task := range tasks {
		for _, p := range strings.Split(task.ProductList, ",") {
			p = strings.TrimSpace(p)
			if len(p) == 0 || productSet[p] {
				continue
			}
			productSet[p] = true
			products = append(products, p)
		}
	}
	return strings.Join(products, ",")
}

// GetOrchestrationDisplayProducts 列表展示用：编排级产品；未配置时展示步骤任务产品并集
func GetOrchestrationDisplayProducts(id, productList, taskList string) string {
	if len(strings.TrimSpace(productList)) > 0 {
		return productList
	}
	orch := DbTaskOrchestration{Id: id, TaskOrchestration: TaskOrchestration{TaskList: taskList}}
	tasks, err := GetOrchestrationTasks(&orch)
	if err != nil {
		return ""
	}
	return GetOrchestrationEffectiveProducts(&orch, tasks)
}

// SearchOrchestrationIdsByProduct 按产品关键字匹配编排，口径与列表「产品配置」列展示一致：
// 配置了编排级产品时按编排产品匹配；未配置时按步骤任务的产品并集匹配；
// 任务缺失时列表展示为空，也不参与匹配。
func SearchOrchestrationIdsByProduct(keyword string) (ids []string, err error) {
	keyword = strings.TrimSpace(keyword)
	if len(keyword) == 0 {
		return
	}
	var orchs []DbTaskOrchestration
	if err = models.Orm.Table("task_orchestration").Where("deleted_at is null").Find(&orchs).Error; err != nil {
		return
	}
	nameSet := make(map[string]bool)
	for _, orch := range orchs {
		for _, name := range ParseTaskList(orch.TaskList) {
			nameSet[name] = true
		}
	}
	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	byName := make(map[string]DbSchedule, len(names))
	if len(names) > 0 {
		var rows []DbSchedule
		if err = models.Orm.Table("schedule").Where("task_name in (?) and deleted_at is null", names).Find(&rows).Error; err != nil {
			return
		}
		for _, row := range rows {
			byName[row.TaskName] = row
		}
	}
	lowered := strings.ToLower(keyword)
	for i := range orchs {
		if strings.Contains(strings.ToLower(orchestrationDisplayProducts(&orchs[i], byName)), lowered) {
			ids = append(ids, orchs[i].Id)
		}
	}
	return
}

// orchestrationDisplayProducts 与 GetOrchestrationDisplayProducts 同口径，复用已加载的任务，避免逐行查库
func orchestrationDisplayProducts(orch *DbTaskOrchestration, byName map[string]DbSchedule) string {
	if len(strings.TrimSpace(orch.ProductList)) > 0 {
		return strings.TrimSpace(orch.ProductList)
	}
	var tasks []DbSchedule
	for _, name := range ParseTaskList(orch.TaskList) {
		task, ok := byName[name]
		if !ok {
			return ""
		}
		tasks = append(tasks, task)
	}
	return GetOrchestrationEffectiveProducts(orch, tasks)
}

// CopyTaskOrchestration 复制编排，名称追加"_副本"，状态重置为未执行
func CopyTaskOrchestration(id, userName string) (err error) {
	orch, err := GetTaskOrchestrationById(id)
	if err != nil {
		return
	}

	copied := TaskOrchestration{
		OrchestrationName:   nextOrchestrationCopyName(fmt.Sprintf("%s_%s", orch.OrchestrationName, T("task_orchestration.copy_suffix"))),
		TaskList:            NormalizeTaskList(orch.TaskList),
		ProductList:         orch.ProductList,
		OrchestrationStatus: OrchestrationStatusNotStarted,
		Remark:              orch.Remark,
		UserName:            userName,
		CreatedAt:           time.Now().Format(baseFormat),
		UpdatedAt:           time.Now().Format(baseFormat),
	}

	if err = models.Orm.Table("task_orchestration").Omit("deleted_at", "last_at").Create(&copied).Error; err != nil {
		Logger.Error("%s", err)
	}
	return
}

// nextOrchestrationCopyName 名称冲突时追加时间戳，保证复制可用
func nextOrchestrationCopyName(name string) string {
	var count int64
	models.Orm.Table("task_orchestration").Where("orchestration_name = ?", name).Count(&count)
	if count == 0 {
		return name
	}
	return fmt.Sprintf("%s_%s", name, time.Now().Format("20060102150405"))
}

// PreDeleteTaskOrchestrations 删除前校验：执行中的编排不允许删除
func PreDeleteTaskOrchestrations(ids []string) error {
	var running []DbTaskOrchestration
	err := models.Orm.Table("task_orchestration").
		Where("id in (?) and orchestration_status = ?", ids, OrchestrationStatusRunning).
		Find(&running).Error
	if err != nil {
		return err
	}
	if len(running) > 0 {
		return fmt.Errorf(T("error.orchestration_delete_running"), running[0].OrchestrationName)
	}
	return nil
}

// SoftDeleteTaskOrchestrations 软删除：改写名称（保证同名可重建）并置 deleted_at
func SoftDeleteTaskOrchestrations(ids []string) error {
	deletedAt := time.Now().Format(baseFormat)
	suffix := "_deleted_" + time.Now().Format("20060102150405")
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if len(id) == 0 {
			continue
		}
		orch, err := GetTaskOrchestrationById(id)
		if err != nil {
			return err
		}
		name := fmt.Sprintf("%s%s", orch.OrchestrationName, suffix)
		if err = models.Orm.Table("task_orchestration").Where("id = ?", id).
			Updates(map[string]interface{}{
				"orchestration_name": name,
				"deleted_at":         deletedAt,
				"updated_at":         deletedAt,
			}).Error; err != nil {
			Logger.Error("%s", err)
			return err
		}
	}
	return nil
}

// ResetRunningOrchestrations 服务重启时把遗留的 running 状态重置为 finished
func ResetRunningOrchestrations() {
	err := models.Orm.Table("task_orchestration").
		Where("orchestration_status = ?", OrchestrationStatusRunning).
		Update("orchestration_status", OrchestrationStatusFinished).Error
	if err != nil {
		Logger.Error("重置任务编排运行状态失败: %s", err)
	}

	// 服务重启时，遗留的"生成中"编排报告标记为失败，避免永远显示生成中
	if err = models.Orm.Table("dashboard").
		Where("report_type = ? and status = ?", OrchestrationReportType, "generating").
		Update("status", "failed").Error; err != nil {
		Logger.Error("重置编排报告状态失败: %s", err)
	}
}

// GetOrchestrationOptions 编排下拉选项（列表表单/加入编排弹窗使用）
func GetOrchestrationOptions() (options []types.FieldOption) {
	var rows []DbTaskOrchestration
	models.Orm.Table("task_orchestration").
		Select("id, orchestration_name").
		Where("deleted_at is null").
		Order("updated_at desc").
		Find(&rows)
	for _, row := range rows {
		options = append(options, types.FieldOption{Value: row.Id, Text: row.OrchestrationName})
	}
	return
}

// AddTasksToOrchestration 把任务追加到编排末尾，已存在则跳过
func AddTasksToOrchestration(orchestrationId, taskIds string) (added, skipped int, err error) {
	orch, err := GetTaskOrchestrationById(orchestrationId)
	if err != nil {
		return
	}
	if orch.OrchestrationStatus == OrchestrationStatusRunning {
		return 0, 0, fmt.Errorf(T("error.orchestration_edit_running"), orch.OrchestrationName)
	}

	var ids []string
	for _, id := range strings.Split(taskIds, ",") {
		id = strings.TrimSpace(id)
		if len(id) > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return 0, 0, fmt.Errorf(T("error.orchestration_select_task_first"))
	}

	var tasks []DbSchedule
	if err = models.Orm.Table("schedule").
		Select("id, task_name").
		Where("id in (?) and deleted_at is null", ids).
		Order("id asc").
		Find(&tasks).Error; err != nil {
		return
	}

	names := ParseTaskList(orch.TaskList)
	existing := make(map[string]bool, len(names))
	for _, name := range names {
		existing[name] = true
	}
	for _, task := range tasks {
		if existing[task.TaskName] {
			skipped++
			continue
		}
		existing[task.TaskName] = true
		names = append(names, task.TaskName)
		added++
	}
	if added == 0 {
		return
	}

	err = models.Orm.Table("task_orchestration").Where("id = ?", orchestrationId).
		Updates(map[string]interface{}{
			"task_list":  strings.Join(names, "\n"),
			"updated_at": time.Now().Format(baseFormat),
		}).Error
	return
}

// BuildOrchestrationTaskListHTML 详情页步骤表：序号、任务、类型、产品、可用性
func BuildOrchestrationTaskListHTML(taskList, dataLocale string) string {
	names := ParseTaskList(taskList)
	if len(names) == 0 {
		return ""
	}

	var rows []DbSchedule
	models.Orm.Table("schedule").Where("task_name in (?) and deleted_at is null", names).Find(&rows)
	byName := make(map[string]DbSchedule, len(rows))
	for _, row := range rows {
		byName[row.TaskName] = row
	}

	var b strings.Builder
	b.WriteString(`<table class="table table-bordered"><thead><tr>`)
	for _, h := range []string{"#", T("common.task_name"), T("common.task_type"), T("common.product_list")} {
		b.WriteString(`<th>` + h + `</th>`)
	}
	b.WriteString(`</tr></thead><tbody>`)
	for i, name := range names {
		task, ok := byName[name]
		b.WriteString(`<tr>`)
		b.WriteString(fmt.Sprintf(`<td>%d</td>`, i+1))
		if ok {
			display := GetTaskLocalized(task.TaskName, dataLocale)
			if len(display) == 0 {
				display = task.TaskName
			}
			b.WriteString(fmt.Sprintf(`<td><a href="/admin/info/schedule/detail?__goadmin_detail_pk=%s">%s</a></td>`, html.EscapeString(task.Id), html.EscapeString(display)))
			taskType := T("common.data")
			if task.TaskType == "scene" {
				taskType = T("common.scene")
			}
			b.WriteString(`<td>` + taskType + `</td>`)
			b.WriteString(`<td>` + html.EscapeString(task.ProductList) + `</td>`)
		} else {
			b.WriteString(`<td>` + html.EscapeString(name) + ` <span style="color:#c0392b">(` + T("task_orchestration.task_unavailable") + `)</span></td>`)
			b.WriteString(`<td>-</td><td>-</td>`)
		}
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return b.String()
}

// GetOrchestrationRecentReports 最近 N 次编排报告
func GetOrchestrationRecentReports(orchestrationId string, limit int) (reports []DashboardReport) {
	models.Orm.Table("dashboard").
		Where("report_type = ? and related_task_ids like ?", OrchestrationReportType, fmt.Sprintf("orch:%s:%%", orchestrationId)).
		Order("created_at desc").
		Limit(limit).
		Find(&reports)
	return
}

// BuildOrchestrationRecentReportsHTML 详情页「最近执行记录」
func BuildOrchestrationRecentReportsHTML(orchestrationId string, limit int) string {
	reports := GetOrchestrationRecentReports(orchestrationId, limit)
	if len(reports) == 0 {
		return `<div style="margin-top:16px;color:#888">` + T("task_orchestration.no_report_yet") + `</div>`
	}
	var b strings.Builder
	b.WriteString(`<div style="margin-top:16px;font-weight:600">` + T("task_orchestration.recent_runs") + `</div>`)
	b.WriteString(`<table class="table table-bordered" style="margin-top:8px"><thead><tr>`)
	for _, h := range []string{
		T("schedule_report.report_name"),
		T("common.status"),
		T("schedule_report.time_range_start"),
		T("task_orchestration.summary_duration"),
		T("task_orchestration.summary_success"),
		T("task_orchestration.summary_fail"),
		T("task_orchestration.summary_skipped"),
	} {
		b.WriteString(`<th>` + html.EscapeString(h) + `</th>`)
	}
	b.WriteString(`</tr></thead><tbody>`)
	for _, r := range reports {
		var summary OrchestrationSummary
		if len(r.ReportData) > 0 {
			var data OrchestrationReportData
			if json.Unmarshal([]byte(r.ReportData), &data) == nil {
				summary = data.Summary
			}
		}
		b.WriteString(`<tr>`)
		b.WriteString(fmt.Sprintf(`<td><a href="/admin/report_dashboard?id=%s">%s</a></td>`, html.EscapeString(r.Id), html.EscapeString(r.ReportName)))
		b.WriteString(`<td>` + html.EscapeString(orchestrationReportStatusText(r.Status)) + `</td>`)
		b.WriteString(`<td>` + html.EscapeString(r.TimeRangeStart) + `</td>`)
		b.WriteString(fmt.Sprintf(`<td>%ds</td><td>%d</td><td>%d</td><td>%d</td>`,
			summary.DurationSeconds, summary.Success, summary.Fail, summary.Skipped))
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return b.String()
}

func orchestrationReportStatusText(status string) string {
	switch status {
	case "generating":
		return T("task_orchestration.report_generating")
	case "finished":
		return T("common.finished")
	default:
		return status
	}
}

// NormalizeProductList 产品列表去重、去空格、逗号分隔
func NormalizeProductList(productList string) string {
	return strings.Join(splitDependencyNames(productList), ",")
}

// ValidateOrchestrationProducts 校验编排配置的产品存在
func ValidateOrchestrationProducts(orchestrationName, productList string) error {
	products := splitDependencyNames(productList)
	if len(products) == 0 {
		return nil
	}
	var rows []Product
	if err := models.Orm.Table("product").Where("product in (?)", products).Find(&rows).Error; err != nil {
		return err
	}
	found := make(map[string]bool, len(rows))
	for _, row := range rows {
		found[row.Name] = true
	}
	for _, name := range products {
		if !found[name] {
			return fmt.Errorf(T("error.orchestration_product_missing"), orchestrationName, name)
		}
	}
	return nil
}
