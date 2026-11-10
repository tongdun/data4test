package biz

import (
	"data4test/models"
	"fmt"
	"strings"
	"time"
)

// RunOrchestration 同步执行一轮任务编排（串行）
func RunOrchestration(orchestrationId, operator string) (reportId string, err error) {
	orch, err := GetTaskOrchestrationById(orchestrationId)
	if err != nil {
		return "", err
	}
	if _, loaded := orchestrationRunning.LoadOrStore(orchestrationId, struct{}{}); loaded {
		return "", fmt.Errorf(T("error.orchestration_running"), orch.OrchestrationName)
	}
	defer orchestrationRunning.Delete(orchestrationId)

	return runOrchestration(orchestrationId, operator, RunTaskOnce)
}

// StartOrchestrationAsync 同步抢锁后异步执行；抢锁失败立即返回，保证页面能收到"正在执行"
func StartOrchestrationAsync(orchestrationId, operator string) error {
	if err := CheckOrchestrationStartable(orchestrationId); err != nil {
		return err
	}
	orch, err := GetTaskOrchestrationById(orchestrationId)
	if err != nil {
		return err
	}
	if _, loaded := orchestrationRunning.LoadOrStore(orchestrationId, struct{}{}); loaded {
		return fmt.Errorf(T("error.orchestration_running"), orch.OrchestrationName)
	}

	go func() {
		defer orchestrationRunning.Delete(orchestrationId)
		defer func() {
			if e := recover(); e != nil {
				Logger.Error("run orchestration panic: %v", e)
			}
		}()
		if _, e := runOrchestration(orchestrationId, operator, RunTaskOnce); e != nil {
			Logger.Error("%s", e)
		}
	}()
	return nil
}

func runOrchestration(orchestrationId, operator string, runner func(string, string, TaskRunOptions) (TaskRunResult, error)) (reportId string, err error) {
	orch, err := GetTaskOrchestrationById(orchestrationId)
	if err != nil {
		return "", err
	}

	tasks, err := GetOrchestrationTasks(&orch)
	if err != nil {
		return "", err
	}

	// 报告快照用「本轮涉及的产品」；但只有当编排显式配置产品时才覆盖各步骤任务
	products := GetOrchestrationEffectiveProducts(&orch, tasks)
	productOverride := ""
	if len(strings.TrimSpace(orch.ProductList)) > 0 {
		productOverride = orch.ProductList
	}
	startTime := time.Now()
	runTag := fmt.Sprintf("orch%s_%s", orchestrationId, startTime.Format("20060102150405"))
	reportTag := fmt.Sprintf("orch:%s:%s", orchestrationId, startTime.Format("20060102150405"))

	reportId, err = CreateOrchestrationReportRecord(reportTag, operator, orch, products)
	if err != nil {
		return "", err
	}

	steps := make([]OrchestrationStep, len(tasks))
	for index, task := range tasks {
		steps[index] = OrchestrationStep{StepOrder: index + 1, TaskId: task.Id, TaskName: task.TaskName, TaskType: task.TaskType, Status: OrchestrationStepSkipped}
	}
	active := -1
	// 状态收尾单独 defer，即使报告聚合异常也不能遗留 running。
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("orchestration finalization panic: %v", recovered)
		}
		end := time.Now().Format(baseFormat)
		if updateErr := models.Orm.Table("task_orchestration").Where("id = ?", orchestrationId).
			Updates(map[string]interface{}{"orchestration_status": OrchestrationStatusFinished, "last_at": end, "updated_at": end}).Error; updateErr != nil {
			err = combineRunErrors(err, updateErr)
		}
	}()
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("orchestration execution panic: %v", recovered)
			if active >= 0 {
				steps[active].Status = OrchestrationStepFail
				steps[active].ErrorMsg = err.Error()
				steps[active].EndAt = time.Now().Format(baseFormat)
			}
		}
		endTime := time.Now()
		report := BuildOrchestrationReportData(&orch, runTag, products, steps, startTime, endTime)
		if saveErr := SaveOrchestrationReport(reportId, report, startTime, endTime); saveErr != nil {
			err = combineRunErrors(err, saveErr)
		}
	}()
	if err = models.Orm.Table("task_orchestration").Where("id = ?", orchestrationId).
		Update("orchestration_status", OrchestrationStatusRunning).Error; err != nil {
		return reportId, err
	}

	for index, task := range tasks {
		active = index
		step := &steps[index]
		step.StartAt = time.Now().Format(baseFormat)
		result, runErr := runOrchestrationTask(task.Id, operator, TaskRunOptions{
			RunTag: fmt.Sprintf("%s_%s_%d", task.Id, runTag, index+1), ProductOverride: productOverride, UpdateTaskStatus: false, SyncReport: true, StrictExecution: true,
		}, runner)
		step.StartAt = result.StartAt
		step.EndAt = result.EndAt
		step.DurationSeconds = result.DurationSeconds
		step.DashboardId = result.DashboardId
		step.Total = result.Total
		step.Success = result.SuccessCount
		step.Fail = result.FailCount
		if runErr != nil {
			step.Status = OrchestrationStepFail
			step.ErrorMsg = runErr.Error()
			if result.AuthFailed {
				break
			}
		} else {
			step.Status = OrchestrationStepSuccess
		}
	}
	active = -1

	return reportId, err
}

// RunOrchestrations 批量立即执行：逐个同步抢锁，失败项立即返回
func RunOrchestrations(ids, operator string) (started []string, err error) {
	for _, id := range strings.Split(ids, ",") {
		id = strings.TrimSpace(id)
		if len(id) == 0 {
			continue
		}
		if e := StartOrchestrationAsync(id, operator); e != nil {
			return started, e
		}
		started = append(started, id)
	}

	if len(started) == 0 {
		return started, fmt.Errorf(T("error.orchestration_select_first"))
	}
	return started, nil
}

// runOrchestrationTask treats a task panic as an ordinary failed step.
func runOrchestrationTask(id, operator string, opts TaskRunOptions, run func(string, string, TaskRunOptions) (TaskRunResult, error)) (result TaskRunResult, err error) {
	start := time.Now()
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("task execution panic: %v", recovered)
			result.Status = OrchestrationStepFail
			result.Err = err
			result.StartAt = start.Format(baseFormat)
			result.EndAt = time.Now().Format(baseFormat)
			result.DurationSeconds = int(time.Since(start).Seconds())
		}
	}()
	return run(id, operator, opts)
}
