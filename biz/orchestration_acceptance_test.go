package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

// 本文件按需求文档第 10 章验收标准（AC-01～AC-25）组织自动化验证。
// 运行方式（仓库根目录）：go test -vet=off -exec "$PWD/build/test-exec.sh" ./biz -run TestAC

func acSeedOrchestration(t *testing.T, name, tasks string) {
	t.Helper()
	p1Exec(t, models.Orm, "INSERT INTO task_orchestration (orchestration_name,task_list,orchestration_status) VALUES (?,?,'not_started')", name, tasks)
}

// AC-01 新建计划（3 个步骤）：列表步骤数与详情顺序
func TestAC01CreateOrchestrationWithThreeSteps(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type,task_status) VALUES ('t1','data','not_started'),('t2','data','not_started'),('t3','data','not_started')")

	if err := ValidateOrchestrationTaskList("ac01", "t1\nt2\nt3"); err != nil {
		t.Fatalf("valid task list rejected: %v", err)
	}
	orch := TaskOrchestration{OrchestrationName: "ac01", TaskList: NormalizeTaskList(" t1 \r\nt2\n\nt3 "), OrchestrationStatus: OrchestrationStatusNotStarted}
	if err := models.Orm.Table("task_orchestration").Omit("deleted_at", "last_at").Create(&orch).Error; err != nil {
		t.Fatal(err)
	}
	var saved DbTaskOrchestration
	models.Orm.Table("task_orchestration").Where("orchestration_name = ?", "ac01").Find(&saved)
	names := ParseTaskList(saved.TaskList)
	if len(names) != 3 || names[0] != "t1" || names[1] != "t2" || names[2] != "t3" {
		t.Fatalf("step count/order incorrect: %v", names)
	}
}

// AC-02 名称重复 / 任务不存在 / 任务重复 / 无步骤
func TestAC02ValidationRejectsInvalidInput(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type) VALUES ('dup','data')")
	acSeedOrchestration(t, "existing", "dup")

	if err := ValidateOrchestrationName("existing", ""); err == nil {
		t.Fatal("duplicate name accepted")
	}
	if err := ValidateOrchestrationTaskList("ac02", ""); err == nil {
		t.Fatal("empty task list accepted")
	}
	if err := ValidateOrchestrationTaskList("ac02", "dup\ndup"); err == nil {
		t.Fatal("duplicate task accepted")
	}
	if err := ValidateOrchestrationTaskList("ac02", "dup\nmissing"); err == nil {
		t.Fatal("missing task accepted")
	}
}

// AC-03 / AC-04 串行顺序 + 普通失败继续
func TestAC03AC04SerialOrderAndContinueOnFailure(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type) VALUES ('a','data'),('b','data'),('c','data')")
	acSeedOrchestration(t, "serial", "a\nb\nc")

	var order []string
	reportId, err := runOrchestration("1", "test", func(id, _ string, _ TaskRunOptions) (TaskRunResult, error) {
		order = append(order, id)
		if id == "2" {
			return TaskRunResult{Status: "fail", Err: fmt.Errorf("case failed")}, fmt.Errorf("case failed")
		}
		return TaskRunResult{Status: "success"}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(order, ",") != "1,2,3" {
		t.Fatalf("not executed serially in order: %v", order)
	}

	var report DashboardReport
	models.Orm.Table("dashboard").Where("id = ?", reportId).Find(&report)
	var data OrchestrationReportData
	if err = json.Unmarshal([]byte(report.ReportData), &data); err != nil {
		t.Fatal(err)
	}
	if data.Summary.Success != 2 || data.Summary.Fail != 1 || data.Summary.Skipped != 0 {
		t.Fatalf("failure did not continue: %+v", data.Summary)
	}
}

// AC-05 场景步骤 401 停止后续
func TestAC05AuthFailureStopsRemainingSteps(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type) VALUES ('a','scene'),('b','scene')")
	acSeedOrchestration(t, "auth", "a\nb")

	calls := 0
	reportId, err := runOrchestration("1", "test", func(_ string, _ string, _ TaskRunOptions) (TaskRunResult, error) {
		calls++
		return TaskRunResult{Status: "fail", AuthFailed: true, Err: fmt.Errorf("401")}, fmt.Errorf("401")
	})
	if err != nil {
		t.Fatal(err)
	}
	var report DashboardReport
	models.Orm.Table("dashboard").Where("id = ?", reportId).Find(&report)
	var data OrchestrationReportData
	_ = json.Unmarshal([]byte(report.ReportData), &data)
	if calls != 1 || data.Steps[1].Status != OrchestrationStepSkipped || data.Summary.Skipped != 1 {
		t.Fatalf("auth failure did not stop: calls=%d data=%+v", calls, data)
	}
}

// AC-06 data 步骤 401 同样可识别
func TestAC06DataAuthFailureRecognized(t *testing.T) {
	zh := fmt.Errorf("请求失败，返回码: 401, 返回信息: unauthorized")
	en := fmt.Errorf("request failed, status code: 401, response message: unauthorized")
	other := fmt.Errorf("test fail")
	if !isAuthFailureError(zh) || !isAuthFailureError(en) {
		t.Fatal("401 not recognized across locales")
	}
	if isAuthFailureError(other) {
		t.Fatal("normal failure misjudged as auth failure")
	}
}

// AC-07 / AC-13 / AC-14 立即执行：1 条编排报告 + N 条任务报告，步骤可下钻
func TestAC07ReportCountsAndDrilldown(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type) VALUES ('a','data'),('b','data')")
	acSeedOrchestration(t, "reports", "a\nb")

	var taskReportIds []string
	for _, taskId := range []string{"1", "2"} {
		var task DbSchedule
		models.Orm.Table("schedule").Where("id = ?", taskId).Find(&task)
		id, err := CreateDashboardRecord(taskId+"_ac", "test", task)
		if err != nil {
			t.Fatal(err)
		}
		taskReportIds = append(taskReportIds, id)
	}

	index := 0
	reportId, err := runOrchestration("1", "test", func(_ string, _ string, _ TaskRunOptions) (TaskRunResult, error) {
		id := taskReportIds[index]
		index++
		return TaskRunResult{Status: "success", DashboardId: id}, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	var taskCount, orchestrationCount int
	models.Orm.Table("dashboard").Where("report_type = ?", "task").Count(&taskCount)
	models.Orm.Table("dashboard").Where("report_type = ?", OrchestrationReportType).Count(&orchestrationCount)
	if taskCount != 2 || orchestrationCount != 1 {
		t.Fatalf("report counts incorrect: task=%d orchestration=%d", taskCount, orchestrationCount)
	}

	var report DashboardReport
	models.Orm.Table("dashboard").Where("id = ?", reportId).Find(&report)
	if report.ReportType != OrchestrationReportType || !strings.HasPrefix(report.RelatedTaskIds, "orch:1:") {
		t.Fatalf("orchestration report linkage incorrect: %+v", report)
	}
	var data OrchestrationReportData
	_ = json.Unmarshal([]byte(report.ReportData), &data)
	if data.Steps[0].DashboardId != taskReportIds[0] || data.Steps[1].DashboardId != taskReportIds[1] {
		t.Fatalf("step drilldown missing: %+v", data.Steps)
	}
}

// AC-08 同一计划并发触发仅一轮
func TestAC08MutexRejectsConcurrentRun(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type) VALUES ('a','data')")
	acSeedOrchestration(t, "lock", "a")

	orchestrationRunning.Store("1", struct{}{})
	defer orchestrationRunning.Delete("1")
	if _, err := RunOrchestration("1", "test"); err == nil || !strings.Contains(err.Error(), "正在执行") {
		t.Fatalf("concurrent run not rejected: %v", err)
	}
}

// AC-09 / AC-10 编排级产品覆盖与回退
func TestAC09AC10ProductOverrideAndFallback(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type,product_list) VALUES ('a','scene','TA'),('b','scene','TB')")

	acSeedOrchestration(t, "override", "a\nb")
	p1Exec(t, models.Orm, "UPDATE task_orchestration SET product_list='PA,PB' WHERE id=1")
	var captured []string
	if _, err := runOrchestration("1", "test", func(_ string, _ string, opts TaskRunOptions) (TaskRunResult, error) {
		captured = append(captured, opts.ProductOverride)
		return TaskRunResult{Status: "success"}, nil
	}); err != nil {
		t.Fatal(err)
	}
	if captured[0] != "PA,PB" || captured[1] != "PA,PB" {
		t.Fatalf("orchestration products not applied: %v", captured)
	}

	p1Exec(t, models.Orm, "DELETE FROM task_orchestration")
	acSeedOrchestration(t, "fallback", "a\nb")
	var fallback DbTaskOrchestration
	models.Orm.Table("task_orchestration").Where("orchestration_name = ?", "fallback").Find(&fallback)
	captured = nil
	if _, err := runOrchestration(fallback.Id, "test", func(_ string, _ string, opts TaskRunOptions) (TaskRunResult, error) {
		captured = append(captured, opts.ProductOverride)
		return TaskRunResult{Status: "success"}, nil
	}); err != nil {
		t.Fatal(err)
	}
	if captured[0] != "" || captured[1] != "" {
		t.Fatalf("fallback should not override task products: %v", captured)
	}
}

// AC-11 data 步骤只取首个产品
func TestAC11FirstProductForData(t *testing.T) {
	if got := firstProduct("A,B"); got != "A" {
		t.Fatalf("firstProduct(A,B)=%s", got)
	}
	if got := firstProduct(" B , A "); got != "B" {
		t.Fatalf("firstProduct with spaces=%s", got)
	}
	if got := firstProduct(""); got != "" {
		t.Fatalf("firstProduct empty=%s", got)
	}
}

// AC-12 编排执行不改写任务状态
func TestAC12TaskStatusUnchanged(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type,task_status,data_list) VALUES ('a','data','not_started','missing.yaml')")
	if _, err := RunTaskOnce("1", "test", TaskRunOptions{RunTag: "ac12", SyncReport: true, StrictExecution: true}); err == nil {
		t.Fatal("missing dependency should fail in strict mode")
	}
	var task DbSchedule
	models.Orm.Table("schedule").Where("id = 1").Find(&task)
	if task.TaskStatus != "not_started" {
		t.Fatalf("orchestration changed task status: %s", task.TaskStatus)
	}
}

// AC-15 离线渲染由 pages 包测试覆盖（TestAC15...），此处不重复。

// AC-16 加入编排：追加并跳过重复
func TestAC16AddTasksToOrchestration(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type) VALUES ('a','data'),('b','data')")
	acSeedOrchestration(t, "join", "a")

	added, skipped, err := AddTasksToOrchestration("1", "1,2")
	if err != nil {
		t.Fatal(err)
	}
	if added != 1 || skipped != 1 {
		t.Fatalf("join counts incorrect: added=%d skipped=%d", added, skipped)
	}
	var orch DbTaskOrchestration
	models.Orm.Table("task_orchestration").Where("id = 1").Find(&orch)
	if NormalizeTaskList(orch.TaskList) != "a\nb" {
		t.Fatalf("task_list not appended: %q", orch.TaskList)
	}

	p1Exec(t, models.Orm, "UPDATE task_orchestration SET orchestration_status='running' WHERE id=1")
	if _, _, err = AddTasksToOrchestration("1", "2"); err == nil {
		t.Fatal("running orchestration accepted join")
	}
}

// AC-17 权限与 i18n：中英 key 对齐、权限种子存在
func TestAC17I18nAndPermissionSeeds(t *testing.T) {
	zh := loadJSONKeys(t, "mgmt/i18n/i18n_system/zh-CN.json")
	en := loadJSONKeys(t, "mgmt/i18n/i18n_system/en-US.json")
	for key := range zh {
		if strings.HasPrefix(key, "task_orchestration.") || strings.HasPrefix(key, "error.orchestration_") {
			if !en[key] {
				t.Fatalf("missing en-US key: %s", key)
			}
		}
	}
	for key := range en {
		if strings.HasPrefix(key, "task_orchestration.") || strings.HasPrefix(key, "error.orchestration_") {
			if !zh[key] {
				t.Fatalf("missing zh-CN key: %s", key)
			}
		}
	}

	raw, err := os.ReadFile("mgmt/sql/update.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := string(raw)
	for _, token := range []string{"(309,", "goadmin_role_menu", "task_orchestration_import"} {
		if !strings.Contains(sql, token) {
			t.Fatalf("permission seed missing: %s", token)
		}
	}
}

// AC-18 删除后可重建同名编排
func TestAC18RecreateAfterSoftDelete(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type) VALUES ('a','data')")
	acSeedOrchestration(t, "same-name", "a")

	if err := SoftDeleteTaskOrchestrations([]string{"1"}); err != nil {
		t.Fatal(err)
	}
	var deleted DbTaskOrchestration
	models.Orm.Table("task_orchestration").Unscoped().Where("id = 1").Find(&deleted)
	if !strings.Contains(deleted.OrchestrationName, "_deleted_") || deleted.DeletedAt == "" {
		t.Fatalf("soft delete did not rewrite name: %+v", deleted)
	}

	orch := TaskOrchestration{OrchestrationName: "same-name", TaskList: "a", OrchestrationStatus: OrchestrationStatusNotStarted}
	if err := models.Orm.Table("task_orchestration").Omit("deleted_at", "last_at").Create(&orch).Error; err != nil {
		t.Fatalf("recreate failed: %v", err)
	}
}

// AC-19 存量任务保持旧行为：缺失依赖空跑成功
func TestAC19LegacyTaskBehaviorUnchanged(t *testing.T) {
	p1Database(t)
	p1Exec(t, models.Orm, "INSERT INTO schedule (task_name,task_type,task_status,data_list) VALUES ('legacy','data','not_started','missing.yaml')")

	result, err := RunTaskOnce("1", "test", TaskRunOptions{UpdateTaskStatus: true, SyncReport: true})
	if err != nil {
		t.Fatalf("legacy task should not fail on missing dependency: %v", err)
	}
	if result.Status != "success" || result.Total != 0 {
		t.Fatalf("legacy behavior changed: %+v", result)
	}
}

// AC-24 导入时任务缺失：整个编排不落库
func TestAC24MissingTaskRejectsWholeOrchestration(t *testing.T) {
	p1Database(t)
	state := OrchestrationImportState{
		ImportId: "33333333-3333-3333-3333-333333333333",
		UserName: "test",
	}
	if err := saveOrchestrationImportState(state); err != nil {
		t.Fatal(err)
	}
	err := importSingleOrchestration(KOrchestration{Name: "broken", Tasks: []string{"missing-task"}}, "skip", "test")
	if err == nil {
		t.Fatal("missing task accepted")
	}
	var count int
	models.Orm.Table("task_orchestration").Count(&count)
	if count != 0 {
		t.Fatalf("orchestration persisted despite missing task: %d", count)
	}
}

func loadJSONKeys(t *testing.T, path string) map[string]bool {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]interface{}
	if err = json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}
	keys := make(map[string]bool, len(data))
	for key := range data {
		keys[key] = true
	}
	return keys
}
