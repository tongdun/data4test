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

	return OrchestrationReportData{
		OrchestrationId:   orch.Id,
		OrchestrationName: orch.OrchestrationName,
		RunTag:            runTag,
		TriggerType:       OrchestrationTriggerManual,
		ProductList:       products,
		Summary:           summary,
		Steps:             steps,
		Products:          collectOrchestrationProducts(steps),
	}
}

// collectOrchestrationProducts 从各步骤的任务报告中汇总产品维度数据
func collectOrchestrationProducts(steps []OrchestrationStep) []OrchestrationProduct {
	merged := make(map[string]*OrchestrationProduct)
	var order []string
	for _, step := range steps {
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
		for _, p := range data.ByProduct {
			item, ok := merged[p.Product]
			if !ok {
				item = &OrchestrationProduct{Product: p.Product}
				merged[p.Product] = item
				order = append(order, p.Product)
			}
			item.Total += p.Total
			item.Pass += p.Pass
			item.Fail += p.Fail
		}
	}

	products := make([]OrchestrationProduct, 0, len(order))
	for _, name := range order {
		item := merged[name]
		if item.Total > 0 {
			item.PassRate = roundRate(float64(item.Pass) / float64(item.Total) * 100)
		}
		products = append(products, *item)
	}
	return products
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
