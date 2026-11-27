package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// CreateOrchestrationReportRecord 创建编排报告记录（状态 generating）
func CreateOrchestrationReportRecord(reportTag, operator string, orch DbTaskOrchestration, products string) (id string, err error) {
	now := time.Now().Format(baseFormat)
	var dr DashboardReport
	dr.ReportName = fmt.Sprintf(T("dashboard.orchestration_report_name"), orch.OrchestrationName, time.Now().Format("20060102150405"))
	dr.ReportType = OrchestrationReportType
	dr.RelatedTaskIds = reportTag
	dr.RelatedProducts = products
	dr.ReportData = ""
	dr.Status = "generating"
	dr.TimeRangeStart = now
	dr.TimeRangeEnd = now
	dr.Creator = operator
	dr.CreatedAt = now
	dr.UpdatedAt = now

	if err = models.Orm.Table("dashboard").Create(&dr).Error; err != nil {
		Logger.Error("保存编排报告失败: %s", err)
		return "", err
	}

	var saved DashboardReport
	if err = models.Orm.Table("dashboard").
		Where("report_name = ? and report_type = ?", dr.ReportName, OrchestrationReportType).
		Find(&saved).Error; err != nil {
		return "", err
	}
	return saved.Id, nil
}

// BuildOrchestrationReportData 汇总编排报告数据
func BuildOrchestrationReportData(orch *DbTaskOrchestration, runTag, products string, steps []OrchestrationStep, startTime, endTime time.Time) OrchestrationReportData {
	summary := OrchestrationSummary{Total: len(steps)}
	for _, step := range steps {
		switch step.Status {
		case OrchestrationStepSuccess:
			summary.Success++
		case OrchestrationStepFail:
			summary.Fail++
		case OrchestrationStepSkipped:
			summary.Skipped++
		}
	}
	summary.DurationSeconds = int(endTime.Sub(startTime).Seconds())
	if summary.Total > 0 {
		summary.PassRate = roundRate(float64(summary.Success) / float64(summary.Total) * 100)
	}

	stepReports := loadOrchestrationStepReports(steps)
	sceneStats, dataStats := collectOrchestrationResourceStats(stepReports)

	return OrchestrationReportData{
		OrchestrationId:   orch.Id,
		OrchestrationName: orch.OrchestrationName,
		RunTag:            runTag,
		TriggerType:       OrchestrationTriggerManual,
		ProductList:       products,
		StatsVersion:      OrchestrationReportStatsVersion,
		Summary:           summary,
		SceneStats:        sceneStats,
		DataStats:         dataStats,
		Steps:             steps,
		Products:          collectOrchestrationProducts(steps, stepReports, products, runTag),
	}
}

// loadOrchestrationStepReports 读取各步骤的任务报告数据（按步骤下标存放，跳过无报告或解析失败的步骤）
func loadOrchestrationStepReports(steps []OrchestrationStep) map[int]TaskReportData {
	reports := make(map[int]TaskReportData, len(steps))
	for index, step := range steps {
		if len(step.DashboardId) == 0 {
			continue
		}
		var row DashboardReport
		models.Orm.Table("dashboard").Select("report_data").Where("id = ?", step.DashboardId).Find(&row)
		if len(row.ReportData) == 0 {
			continue
		}
		var data TaskReportData
		if err := json.Unmarshal([]byte(row.ReportData), &data); err != nil {
			continue
		}
		reports[index] = data
	}
	return reports
}

// EnsureOrchestrationReportStats 兼容历史报告：统计口径版本较低时，按各步骤任务报告 / 执行历史回溯重算
func EnsureOrchestrationReportStats(data *OrchestrationReportData) {
	if data.StatsVersion >= OrchestrationReportStatsVersion {
		return
	}
	reports := loadOrchestrationStepReports(data.Steps)
	data.SceneStats, data.DataStats = collectOrchestrationResourceStats(reports)
	data.Products = collectOrchestrationProducts(data.Steps, reports, data.ProductList, data.RunTag)
	data.StatsVersion = OrchestrationReportStatsVersion
}

// collectOrchestrationResourceStats 汇总各任务的场景 / 数据文件执行统计
func collectOrchestrationResourceStats(reports map[int]TaskReportData) (sceneStats, dataStats OrchestrationResourceStats) {
	for _, data := range reports {
		sceneStats.Total += data.SceneStats.Total
		sceneStats.Success += data.SceneStats.Pass
		sceneStats.Fail += data.SceneStats.Fail

		dataStats.Total += data.DataStats.Total
		dataStats.Success += data.DataStats.Pass
		dataStats.Fail += data.DataStats.Fail
	}
	sceneStats.PassRate = resourcePassRate(sceneStats)
	dataStats.PassRate = resourcePassRate(dataStats)
	return sceneStats, dataStats
}

func resourcePassRate(stats OrchestrationResourceStats) float64 {
	if stats.Total == 0 {
		return 0
	}
	return roundRate(float64(stats.Success) / float64(stats.Total) * 100)
}

// collectOrchestrationProducts 按产品汇总任务 / 场景 / 数据文件的执行情况
func collectOrchestrationProducts(steps []OrchestrationStep, reports map[int]TaskReportData, defaultProducts, runTag string) []OrchestrationProduct {
	merged := make(map[string]*OrchestrationProduct)
	var order []string
	ensure := func(name string) *OrchestrationProduct {
		if len(name) == 0 {
			return nil
		}
		item, ok := merged[name]
		if !ok {
			item = &OrchestrationProduct{Product: name}
			merged[name] = item
			order = append(order, name)
		}
		return item
	}

	type stepProduct struct {
		name     string
		total    int
		fail     int
		hasStats bool
	}

	for index, step := range steps {
		stepProducts := make([]stepProduct, 0)
		if data, ok := reports[index]; ok {
			for _, p := range data.ByProduct {
				stepProducts = append(stepProducts, stepProduct{name: p.Product, total: p.Total, fail: p.Fail, hasStats: true})
			}
		}
		if len(stepProducts) == 0 {
			for _, name := range parseProductList(defaultProducts) {
				stepProducts = append(stepProducts, stepProduct{name: name})
			}
		}
		taskTag := step.TaskTag
		if len(taskTag) == 0 {
			// 历史报告未保存任务执行标记，按编排标记与步骤下标回溯
			taskTag = fmt.Sprintf("%s_%s_%d", step.TaskId, runTag, step.StepOrder)
		}
		for _, sp := range stepProducts {
			item := ensure(sp.name)
			if item == nil {
				continue
			}
			item.TaskTotal++
			// 任务成功/失败按产品维度隔离：该产品下没有失败执行项即视为成功，
			// 避免任务因其他产品失败而被整体计为失败
			if sp.hasStats && sp.total > 0 {
				if sp.fail > 0 {
					item.TaskFail++
				} else {
					item.TaskSuccess++
				}
			} else {
				switch step.Status {
				case OrchestrationStepSuccess:
					item.TaskSuccess++
				case OrchestrationStepFail:
					item.TaskFail++
				}
			}

			sTotal, sSuccess, sFail := queryProductResourceStats("scene_test_history", taskTag, sp.name)
			item.SceneTotal += sTotal
			item.SceneSuccess += sSuccess
			item.SceneFail += sFail

			dTotal, dSuccess, dFail := queryProductResourceStats("scene_data_test_history", taskTag, sp.name)
			item.DataTotal += dTotal
			item.DataSuccess += dSuccess
			item.DataFail += dFail
		}
	}

	products := make([]OrchestrationProduct, 0, len(order))
	for _, name := range order {
		item := merged[name]
		if item.DataTotal > 0 {
			item.DataPassRate = roundRate(float64(item.DataSuccess) / float64(item.DataTotal) * 100)
		}
		products = append(products, *item)
	}
	return products
}

// queryProductResourceStats 按产品与任务执行标记统计场景 / 数据文件的执行结果
func queryProductResourceStats(table, taskTag, product string) (total, success, fail int) {
	if len(taskTag) == 0 || len(product) == 0 {
		return
	}
	var totalCount int64
	models.Orm.Table(table).Where("task_id = ? and product = ?", taskTag, product).Count(&totalCount)
	total = int(totalCount)
	if total == 0 {
		return
	}
	var successCount, failCount int64
	models.Orm.Table(table).Where("task_id = ? and product = ? and result = ?", taskTag, product, "pass").Count(&successCount)
	models.Orm.Table(table).Where("task_id = ? and product = ? and result = ?", taskTag, product, "fail").Count(&failCount)
	success, fail = int(successCount), int(failCount)
	return
}

// SaveOrchestrationReport 回填编排报告内容并置为已完成
func SaveOrchestrationReport(reportId string, report OrchestrationReportData, startTime, endTime time.Time) error {
	content, err := json.Marshal(report)
	if err != nil {
		Logger.Error("%s", err)
		return err
	}

	return models.Orm.Table("dashboard").Where("id = ?", reportId).
		Updates(map[string]interface{}{
			"report_data":      string(content),
			"status":           "finished",
			"time_range_start": startTime.Format(baseFormat),
			"time_range_end":   endTime.Format(baseFormat),
			"updated_at":       time.Now().Format(baseFormat),
		}).Error
}

// GetOrchestrationLatestReport 查询编排最近一次报告 id
func GetOrchestrationLatestReport(orchestrationId string) (dashboardId string, err error) {
	var row DashboardReport
	err = models.Orm.Table("dashboard").
		Select("id").
		Where("report_type = ? and related_task_ids like ?", OrchestrationReportType, fmt.Sprintf("orch:%s:%%", orchestrationId)).
		Order("created_at desc").
		Limit(1).
		Find(&row).Error
	if err != nil {
		return "", err
	}
	if len(row.Id) == 0 {
		return "", fmt.Errorf(T("error.orchestration_report_not_found"))
	}
	return row.Id, nil
}

func roundRate(value float64) float64 {
	return math.Round(value*100) / 100
}
