package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/template/types"
)

// GenerateTaskReport 生成任务类型的执行报告
func GenerateTaskReport(taskDB DbSchedule, historyId, taskTag, userName string,
	totalExpected, successCount, failCount int, durationSeconds int,
	startTime, endTime string) {

	executedCount := successCount + failCount
	notExecuted := totalExpected - executedCount
	if notExecuted < 0 {
		notExecuted = 0
	}

	// 计算执行率
	var executeRate float64
	if executedCount > 0 {
		executeRate = float64(executedCount) / float64(totalExpected) * 100
	}

	// 从历史表中查询详细统计数据
	sceneTotal, scenePass, sceneFail := querySceneStatsByTaskId(taskTag)
	dataTotal, dataPass, dataFail := queryDataStatsByTaskId(taskTag)

	var productStats []ProductReportItem
	products := parseProductList(taskDB.ProductList)
	for _, product := range products {
		pTotal, pPass, pFail := queryStatsByProductAndTaskId(product, taskTag)
		total := pTotal
		pRate := 0.0

		if total > 0 {
			pRate = float64(pPass) / float64(total) * 100
		}
		productStats = append(productStats, ProductReportItem{
			Product: product,

			Total:    total,
			Pass:     pPass,
			Fail:     pFail,
			PassRate: pRate,
		})
	}

	// 查询失败项明细
	//failItems := queryFailItemsByTaskId(taskTag)

	// 查询历史趋势（同schedule最近7次执行）
	trendItems := queryExecutionTrend(taskDB.Id, historyId)

	// 组装report_data

	t1, _ := time.Parse(baseFormat, startTime)
	t2, _ := time.Parse(baseFormat, endTime)
	durationSeconds = int(t2.Sub(t1).Seconds())

	reportData := TaskReportData{}
	reportData.Overview.TaskName = taskDB.TaskName
	reportData.Overview.TaskType = taskDB.TaskType
	reportData.Overview.ExecutionTime = startTime
	reportData.Overview.Environment = strings.Split(taskDB.ProductList, ",")[0]
	reportData.Overview.StartTime = startTime
	reportData.Overview.EndTime = endTime
	reportData.Overview.DurationSeconds = durationSeconds
	reportData.Overview.TotalExpected = totalExpected
	reportData.Overview.TotalExecuted = executedCount
	reportData.Overview.NotExecuted = notExecuted
	reportData.Overview.SuccessCount = successCount
	reportData.Overview.FailCount = failCount
	reportData.Overview.ExecuteRate = executeRate
	reportData.Overview.Executor = userName

	reportData.SceneStats.Total = sceneTotal
	reportData.SceneStats.Pass = scenePass
	reportData.SceneStats.Fail = sceneFail
	if sceneTotal > 0 {
		reportData.SceneStats.PassRate = float64(scenePass) / float64(sceneTotal) * 100
	}

	reportData.DataStats.Total = dataTotal
	reportData.DataStats.Pass = dataPass
	reportData.DataStats.Fail = dataFail
	if dataTotal > 0 {
		reportData.DataStats.PassRate = float64(dataPass) / float64(dataTotal) * 100
	}

	if productStats == nil {
		productStats = []ProductReportItem{}
	}
	reportData.APITypeDistribution = queryAPITypeDistribution(taskTag)
	// 查询已执行的场景明细，并追加未执行的场景
	executedScenes := querySceneDetails(taskTag)
	if len(taskDB.SceneList) > 0 {
		expectedSceneNames := parseLineList(taskDB.SceneList)
		executedMap := make(map[string]bool, len(executedScenes))
		for _, s := range executedScenes {
			executedMap[s.Name] = true
		}
		for _, name := range expectedSceneNames {
			name = strings.TrimSpace(name)
			if len(name) > 0 && !executedMap[name] {
				executedScenes = append(executedScenes, SceneDetail{
					Name:   name,
					Result: "未执行",
				})
			}
		}
	}
	reportData.SceneDetails = executedScenes

	reportData.DataDetails = queryDataDetails(taskTag)

	reportData.ByProduct = productStats

	if trendItems == nil {
		trendItems = []TrendItem{}
	}
	reportData.Trend = trendItems

	// 序列化JSON
	jsonBytes, err := json.MarshalIndent(reportData, "", "    ")
	if err != nil {
		Logger.Error("序列化报告数据失败: %s", err)
		return
	}

	// 写入schedule_report表
	now := time.Now().Format(baseFormat)
	report := DashboardReport{
		ReportName:      fmt.Sprintf("%s_%s", taskDB.TaskName, startTime),
		ReportType:      "task",
		RelatedTaskIds:  taskDB.Id,
		RelatedProducts: taskDB.ProductList,
		TimeRangeStart:  startTime,
		TimeRangeEnd:    endTime,
		Status:          "finished",
		Creator:         userName,
		ReportData:      string(jsonBytes),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// 更新已存在的 dashboard 记录（由 CreateDashboardRecord 创建）
	if len(historyId) > 0 {
		err = models.Orm.Table("dashboard").Where("id = ?", historyId).Update(map[string]interface{}{
			"report_type": "task",
			"report_name": report.ReportName,
			"report_data": string(jsonBytes),
			"status":      "finished",
			"updated_at":  now,
		}).Error
		if err != nil {
			Logger.Error("更新任务报告失败: %s", err)
		}
	}
}

// querySceneStatsByTaskId 按taskId查询场景执行统计
func querySceneStatsByTaskId(taskId string) (total, pass, fail int) {
	if len(taskId) == 0 {
		return
	}
	var allCount int64
	models.Orm.Table("scene_test_history").Where("task_id = ?", taskId).Count(&allCount)
	total = int(allCount)

	if total > 0 {
		var passCount int64
		models.Orm.Table("scene_test_history").Where("task_id = ? and result = ?", taskId, "pass").Count(&passCount)
		pass = int(passCount)

		var failCount int64
		models.Orm.Table("scene_test_history").Where("task_id = ? and result = ?", taskId, "fail").Count(&failCount)
		fail = int(failCount)
	}
	return
}

// queryDataStatsByTaskId 按taskId查询数据执行统计
func queryDataStatsByTaskId(taskId string) (total, pass, fail int) {
	if len(taskId) == 0 {
		return
	}
	var allCount int64
	models.Orm.Table("scene_data_test_history").Where("task_id = ?", taskId).Count(&allCount)
	total = int(allCount)

	if total > 0 {
		var passCount int64
		models.Orm.Table("scene_data_test_history").Where("task_id = ? and result = ?", taskId, "pass").Count(&passCount)
		pass = int(passCount)

		var failCount int64
		models.Orm.Table("scene_data_test_history").Where("task_id = ? and result = ?", taskId, "fail").Count(&failCount)
		fail = int(failCount)
	}
	return
}

// queryStatsByProductAndTaskId 按产品和taskId查询统计
func queryStatsByProductAndTaskId(product, taskId string) (total, pass, fail int) {
	if len(taskId) == 0 || len(product) == 0 {
		return
	}
	var sceneCount, dataCount int64

	models.Orm.Table("scene_test_history").
		Where("task_id = ? and product = ?", taskId, product).Count(&sceneCount)
	models.Orm.Table("scene_data_test_history").
		Where("task_id = ? and product = ?", taskId, product).Count(&dataCount)

	total = int(sceneCount + dataCount)

	if total > 0 {
		var scenePass, dataPass, sceneFail, dataFail int64
		models.Orm.Table("scene_test_history").
			Where("task_id = ? and product = ? and result = ?", taskId, product, "pass").Count(&scenePass)
		models.Orm.Table("scene_data_test_history").
			Where("task_id = ? and product = ? and result = ?", taskId, product, "pass").Count(&dataPass)
		models.Orm.Table("scene_test_history").
			Where("task_id = ? and product = ? and result = ?", taskId, product, "fail").Count(&sceneFail)
		models.Orm.Table("scene_data_test_history").
			Where("task_id = ? and product = ? and result = ?", taskId, product, "fail").Count(&dataFail)
		pass = int(scenePass + dataPass)
		fail = int(sceneFail + dataFail)
	}
	return
}

// queryFailItemsByTaskId 按taskId查询失败项
func queryFailItemsByTaskId(taskId string) (items []FailItem) {
	if len(taskId) == 0 {
		return
	}
	type sceneFailResult struct {
		Name       string `gorm:"column:name"`
		Product    string `gorm:"column:product"`
		ApiList    string `gorm:"column:data_file_list"`
		FailReason string `gorm:"column:fail_reason"`
	}
	var sceneFails []sceneFailResult
	models.Orm.Table("scene_test_history").
		Select("name, product, data_file_list, fail_reason").
		Where("task_id = ? and result = ? and fail_reason is not null and fail_reason <> '' and fail_reason <> ' '", taskId, "fail").
		Find(&sceneFails)
	for _, f := range sceneFails {
		items = append(items, FailItem{
			Name:   f.Name,
			Type:   "scene",
			Reason: f.FailReason,
		})
	}

	type dataFailResult struct {
		Name       string `gorm:"column:name"`
		Product    string `gorm:"column:product"`
		ApiId      string `gorm:"column:api_id"`
		FailReason string `gorm:"column:fail_reason"`
	}
	var dataFails []dataFailResult
	models.Orm.Table("scene_data_test_history").
		Select("name, product, api_id, fail_reason").
		Where("task_id = ? and result = ? and fail_reason is not null and fail_reason <> '' and fail_reason <> ' '", taskId, "fail").
		Find(&dataFails)
	for _, f := range dataFails {
		items = append(items, FailItem{
			Name:   f.Name,
			Type:   "data",
			APIId:  f.ApiId,
			Reason: f.FailReason,
		})
	}
	return
}

// queryExecutionTrend 查询同schedule最近N次执行趋势
func queryExecutionTrend(scheduleId, excludeHistoryId string) (items []TrendItem) {
	if len(scheduleId) == 0 {
		return
	}
	type historyResult struct {
		Id         string `gorm:"column:id"`
		ReportData string `gorm:"column:report_data"`
		CreatedAt  string `gorm:"column:created_at"`
	}
	var histories []historyResult
	models.Orm.Table("dashboard").
		Select("id, report_data, created_at").
		Where("related_task_ids like CONCAT(?, '\\_%') and id <> ? and status = ?", scheduleId, excludeHistoryId, "finished").
		Order("created_at desc").
		Limit(7).
		Find(&histories)

	// 反转顺序（从旧到新）
	for i := len(histories) - 1; i >= 0; i-- {
		h := histories[i]
		var trendData struct {
			Overview struct {
				TotalExpected int `json:"total_expected"`
				SuccessCount  int `json:"success_count"`
				FailCount     int `json:"fail_count"`
			} `json:"overview"`
		}
		if err := json.Unmarshal([]byte(h.ReportData), &trendData); err != nil {
			continue
		}
		total := trendData.Overview.TotalExpected
		passRate := 0.0
		if total > 0 {
			passRate = float64(trendData.Overview.SuccessCount) / float64(total) * 100
		}
		items = append(items, TrendItem{
			ExecutionTime: h.CreatedAt,
			Total:         total,
			Pass:          trendData.Overview.SuccessCount,
			Fail:          trendData.Overview.FailCount,
			PassRate:      passRate,
		})
	}
	return
}

// parseProductList 解析产品列表
func parseProductList(productList string) (products []string) {
	if len(productList) == 0 {
		return
	}
	for _, p := range strings.Split(productList, ",") {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			products = append(products, p)
		}
	}
	return
}

// CreateScheduleExecutionHistory 创建执行报告
func CreateDashboardRecord(task DbSchedule, taskTag, userName string) (id string, err error) {

	curTime := time.Now().Format(baseFormat)
	var dr DashboardReport
	dr.ReportName = fmt.Sprintf("任务报告_%s_%s", taskTag[len(taskTag)-8], GetRandomStr(4, ""))
	dr.RelatedProducts = task.ProductList
	dr.ReportData = ""
	dr.Status = "running"
	dr.RelatedTaskIds = taskTag
	dr.TimeRangeStart = curTime
	dr.TimeRangeEnd = curTime
	dr.Creator = userName
	dr.CreatedAt = curTime
	dr.UpdatedAt = curTime

	err = models.Orm.Table("dashboard").Create(&dr).Error
	if err != nil {
		Logger.Error("保存执行报告失败: %s", err)
		return "", err
	}

	// 读取回填的ID
	var savedHistory DashboardReport
	models.Orm.Table("dashboard").Where("report_name = ?", dr.ReportName).Find(&savedHistory)
	return savedHistory.Id, nil
}

// FinishScheduleExecutionHistory 完成执行报告（更新结果）
func FinishScheduleExecutionHistory(id string) (err error) {
	if len(id) == 0 {
		return fmt.Errorf("执行报告ID为空")
	}

	curTime := time.Now().Format(baseFormat)
	updates := map[string]interface{}{
		"status":     "finished",
		"updated_at": curTime,
	}
	//if len(failReason) > 0 {
	//	updates["fail_reason"] = failReason
	//}

	err = models.Orm.Table("dashboard").Where("id = ?", id).Update(updates).Error
	if err != nil {
		Logger.Error("更新执行报告失败: %s", err)
	}
	return
}

// taskInfo 任务信息（用于报告生成）
type taskInfo struct {
	Id          string
	TaskName    string
	ProductList string
	TaskType    string
	DataList    string
	SceneList   string
}

// GenerateMultiTaskReports 生成多任务合并报告
// idStr: 逗号分隔的schedule_id列表
// reportProducts: 报告关联的产品（逗号分隔），为空则取所有任务的产品并集
// reportUser: 创建人
func GenerateMultiTaskReports(idStr, reportProducts, reportUser string) (err error) {
	if len(idStr) == 0 || idStr == "," {
		return fmt.Errorf("未选择任务")
	}

	ids := strings.Split(idStr, ",")
	if len(ids) == 0 {
		return fmt.Errorf("未选择任务")
	}

	// 查询选择的schedule完整信息
	var tasks []taskInfo
	for _, sid := range ids {
		sid = strings.TrimSpace(sid)
		if len(sid) == 0 {
			continue
		}
		var t taskInfo
		models.Orm.Table("schedule").Select("id, task_name, product_list, task_type, data_list, scene_list").
			Where("id = ?", sid).Find(&t)
		if len(t.Id) > 0 {
			tasks = append(tasks, t)
		}
	}

	if len(tasks) == 0 {
		return fmt.Errorf("未找到有效的任务")
	}

	// 确定产品列表
	var targetProducts []string
	if len(reportProducts) > 0 {
		targetProducts = parseProductList(reportProducts)
	} else {
		// 取所有任务的产品并集
		productSet := make(map[string]bool)
		for _, t := range tasks {
			for _, p := range parseProductList(t.ProductList) {
				productSet[p] = true
			}
		}
		for p := range productSet {
			targetProducts = append(targetProducts, p)
		}
	}

	if len(targetProducts) == 0 {
		return fmt.Errorf("未选择产品，无法生成报告")
	}

	now := time.Now().Format(baseFormat)
	reportTime := time.Now().Format("20060102150405")

	// 收集 taskIds 和 taskNames
	var taskIds []string
	var taskNames []string
	for _, t := range tasks {
		taskIds = append(taskIds, t.Id)
		taskNames = append(taskNames, t.TaskName)
	}

	// 生成报告名称：多任务用"第一个任务名等N个任务的报告"
	var reportNamePrefix string
	if len(tasks) == 1 {
		reportNamePrefix = tasks[0].TaskName
	} else {
		reportNamePrefix = fmt.Sprintf(T("schedule_report.multi_task_suffix"), len(tasks))
		if len(tasks) > 0 && len(tasks[0].TaskName) > 0 {
			reportNamePrefix = tasks[0].TaskName + " " + reportNamePrefix
		}
	}

	// 为每个产品生成一个报告
	for _, product := range targetProducts {
		err = generateSingleProductReport(tasks, product, reportUser, now, reportTime, reportNamePrefix)
		if err != nil {
			Logger.Error("生成产品[%s]报告失败: %s", product, err)
		}
	}
	return nil
}

func generateSingleProductReport(tasks []taskInfo, product, reportUser, now, reportTime string, reportNamePrefix string) error {
	totalScenePass, totalSceneFail, totalScene := 0, 0, 0
	totalDataPass, totalDataFail, totalData := 0, 0, 0
	totalExecuted := 0
	totalFail := 0
	totalPass := 0
	var taskIds []string
	var taskNames []string

	// 收集所有任务的资源信息
	var allTaskItems []TaskReportItem
	sceneSet := make(map[string]bool) // 去重场景
	dataSet := make(map[string]bool)  // 去重数据文件
	apiSet := make(map[string]bool)   // 去重API

	// 聚合所有任务的明细
	var allSceneDetails []SceneDetailWithTask
	var allDataDetails []DataDetailWithTask
	var allFailItems []FailItem
	var allAPIDist []CountItem
	var allTrendItems []TrendItem

	for _, t := range tasks {
		taskIds = append(taskIds, t.Id)
		taskNames = append(taskNames, t.TaskName)

		// 获取该任务的资源关联
		scenes, datas, apis := getTaskResources(t)
		item := TaskReportItem{
			TaskId:   t.Id,
			TaskName: t.TaskName,
			TaskType: t.TaskType,
			Scenes:   scenes,
			Datas:    datas,
			APIs:     apis,
		}

		// 去重收集全局资源
		for _, s := range scenes {
			sceneSet[s.Name] = true
		}
		for _, d := range datas {
			dataSet[d.Name] = true
		}
		for _, a := range apis {
			apiSet[a.Name] = true
		}

		// 查找该任务在该产品下的最新执行历史
		var history DashboardReport
		models.Orm.Table("dashboard").
			Where("related_task_ids in (?) and related_products like ? and status = ?", t.Id, "%"+product+"%", "finished").
			Order("created_at desc").
			Limit(1).
			Find(&history)

		if len(history.Id) == 0 {
			allTaskItems = append(allTaskItems, item)
			continue
		}

		// 按taskId统计
		curTaskId := fmt.Sprintf("%s_%s", t.Id, strings.Replace(history.TimeRangeStart, " ", "", -1))
		curTaskId = fmt.Sprintf("%s_%s", t.Id, strings.Replace(strings.Replace(history.TimeRangeStart, "-", "", -1), ":", "", -1))
		curTaskId = strings.Replace(curTaskId, " ", "", -1)

		sPass, sFail, sTotal := 0, 0, 0
		dPass, dFail, dTotal := 0, 0, 0

		// 查询场景统计
		var sc, sp, sf int64
		models.Orm.Table("scene_test_history").
			Where("task_id = ? and product = ?", curTaskId, product).Count(&sc)
		models.Orm.Table("scene_test_history").
			Where("task_id = ? and product = ? and result = ?", curTaskId, product, "pass").Count(&sp)
		models.Orm.Table("scene_test_history").
			Where("task_id = ? and product = ? and result = ?", curTaskId, product, "fail").Count(&sf)

		sTotal = int(sc)
		sPass = int(sp)
		sFail = int(sf)

		models.Orm.Table("scene_data_test_history").
			Where("task_id = ? and product = ?", curTaskId, product).Count(&sc)
		models.Orm.Table("scene_data_test_history").
			Where("task_id = ? and product = ? and result = ?", curTaskId, product, "pass").Count(&sp)
		models.Orm.Table("scene_data_test_history").
			Where("task_id = ? and product = ? and result = ?", curTaskId, product, "fail").Count(&sf)

		dTotal = int(sc)
		dPass = int(sp)
		dFail = int(sf)

		totalScene += sTotal
		totalScenePass += sPass
		totalSceneFail += sFail
		totalData += dTotal
		totalDataPass += dPass
		totalDataFail += dFail
		totalExecuted += sTotal + dTotal
		totalPass += sPass + dPass
		totalFail += sFail + dFail

		item.SceneTotal = sTotal
		item.ScenePass = sPass
		item.SceneFail = sFail
		item.DataTotal = dTotal
		item.DataPass = dPass
		item.DataFail = dFail
		item.Total = sTotal + dTotal
		item.Pass = sPass + dPass
		item.Fail = sFail + dFail
		if item.Total > 0 {
			item.PassRate = float64(item.Pass) / float64(item.Total) * 100
		}

		// 查询该任务的场景明细
		taskScenes := querySceneDetailsForTask(curTaskId, product)
		item.Scenes = mergeResourceItems(item.Scenes, countResources(taskScenes))
		for _, sd := range taskScenes {
			allSceneDetails = append(allSceneDetails, SceneDetailWithTask{
				TaskName:   t.TaskName,
				Name:       sd.Name,
				Result:     sd.Result,
				FailReason: sd.FailReason,
			})
		}

		// 查询该任务的数据文件明细
		taskDatas := queryDataDetailsForTask(curTaskId, product)
		for _, dd := range taskDatas {
			allDataDetails = append(allDataDetails, DataDetailWithTask{
				TaskName:   t.TaskName,
				Name:       dd.Name,
				ApiId:      dd.ApiId,
				Result:     dd.Result,
				FailReason: dd.FailReason,
			})
		}

		// 查询该任务的失败项
		taskFails := queryFailItemsForTask(curTaskId, product)
		for _, fi := range taskFails {
			allFailItems = append(allFailItems, FailItem{
				Name:   fmt.Sprintf("[%s] %s", t.TaskName, fi.Name),
				Type:   fi.Type,
				APIId:  fi.APIId,
				Reason: fi.Reason,
			})
		}

		// 聚合API类型分布
		taskAPIs := queryAPITypeDistributionForTask(curTaskId, product)
		allAPIDist = mergeCountItems(allAPIDist, taskAPIs)

		// 查询该任务下场景执行的时间区间
		type execTime struct {
			MinTime string `gorm:"column:min_time"`
			MaxTime string `gorm:"column:max_time"`
		}
		var et execTime
		models.Orm.Table("scene_test_history").
			Select("MIN(created_at) as min_time, MAX(created_at) as max_time").
			Where("task_id = ? and product = ?", curTaskId, product).
			Find(&et)
		if len(et.MinTime) == 0 {
			models.Orm.Table("scene_data_test_history").
				Select("MIN(created_at) as min_time, MAX(created_at) as max_time").
				Where("task_id = ? and product = ?", curTaskId, product).
				Find(&et)
		}
		item.StartTime = et.MinTime
		item.EndTime = et.MaxTime
		if len(et.MinTime) > 0 && len(et.MaxTime) > 0 {
			layouts := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", "2006-01-02T15:04:05"}
			var startT, endT time.Time
			var err1, err2 error
			for _, layout := range layouts {
				startT, err1 = time.Parse(layout, et.MinTime)
				if err1 == nil {
					break
				}
			}
			for _, layout := range layouts {
				endT, err2 = time.Parse(layout, et.MaxTime)
				if err2 == nil {
					break
				}
			}
			if err1 == nil && err2 == nil {
				diff := int(endT.Sub(startT).Seconds())
				if diff < 0 {
					diff = 0
				}
				item.DurationSeconds = diff
			}
		}

		allTaskItems = append(allTaskItems, item)
	}

	// 构建报告数据
	reportData := MultiTaskReportData{}
	if len(reportNamePrefix) > 0 {
		reportData.Overview.ReportName = fmt.Sprintf("%s_%s_%s", reportNamePrefix, product, reportTime)
	} else {
		reportData.Overview.ReportName = fmt.Sprintf("多任务报告_%s_%s", product, reportTime)
	}
	reportData.Overview.Product = product
	reportData.Overview.TaskCount = len(tasks)
	reportData.Overview.TotalCases = totalExecuted
	reportData.Overview.PassCases = totalPass
	reportData.Overview.FailCases = totalFail
	if totalExecuted > 0 {
		reportData.Overview.PassRate = float64(totalPass) / float64(totalExecuted) * 100
	}
	reportData.Overview.SceneCount = len(sceneSet)
	reportData.Overview.DataCount = len(dataSet)
	reportData.Overview.APICount = len(apiSet)

	// 计算执行时间区间：取第一个任务的开始时间 ~ 最后一个任务的结束时间
	if len(allTaskItems) > 0 {
		reportData.Overview.StartTime = allTaskItems[0].StartTime
		reportData.Overview.EndTime = allTaskItems[len(allTaskItems)-1].EndTime
		// 执行耗时 = 各任务执行耗时相加
		totalDuration := 0
		for _, item := range allTaskItems {
			totalDuration += item.DurationSeconds
		}
		reportData.Overview.DurationSeconds = totalDuration
	}
	if totalExecuted > 0 {
		notExec := totalExecuted - (totalPass + totalFail)
		if notExec < 0 {
			notExec = 0
		}
		reportData.Overview.ExecuteRate = float64(totalExecuted) / float64(totalExecuted+notExec) * 100
	}
	reportData.ByTask = allTaskItems
	reportData.SceneList = mapToResourceItems(sceneSet)
	reportData.DataList = mapToResourceItems(dataSet)
	reportData.APIList = mapToResourceItems(apiSet)
	reportData.SceneDetails = allSceneDetails
	reportData.DataDetails = allDataDetails
	reportData.FailItems = allFailItems
	reportData.APITypeDistribution = allAPIDist
	reportData.Trend = allTrendItems

	jsonBytes, err := json.Marshal(reportData)
	if err != nil {
		return err
	}

	reportName := fmt.Sprintf("%s-%s", reportNamePrefix, reportTime)

	// 写入dashboard表
	report := DashboardReport{
		ReportName:      reportName,
		ReportType:      "task",
		RelatedTaskIds:  strings.Join(taskIds, ","),
		RelatedProducts: product,
		TimeRangeStart:  now,
		TimeRangeEnd:    now,
		Status:          "finished",
		Creator:         reportUser,
		ReportData:      string(jsonBytes),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = models.Orm.Table("dashboard").Create(&report).Error
	if err != nil {
		return fmt.Errorf("保存报告失败: %s", err)
	}

	return nil
}

// getTaskResources 查询任务关联的场景、数据文件、API接口
func getTaskResources(t taskInfo) (scenes []ResourceItem, datas []ResourceItem, apis []ResourceItem) {
	scenes = make([]ResourceItem, 0)
	datas = make([]ResourceItem, 0)
	apis = make([]ResourceItem, 0)

	if t.TaskType == "scene" && len(t.SceneList) > 0 {
		// 解析场景名列表
		sceneNames := parseLineList(t.SceneList)
		for _, name := range sceneNames {
			name = strings.TrimSpace(name)
			if len(name) == 0 {
				continue
			}
			scenes = append(scenes, ResourceItem{Name: name})
		}
		// 通过 playbook 关联的数据文件查找 API
		if len(sceneNames) > 0 {
			// 从 playbook 获取关联的数据文件名
			type playbookData struct {
				DataName string `gorm:"column:data_file_list"`
			}
			var pbDataList []playbookData
			models.Orm.Table("playbook").
				Select("data_file_list").
				Where("name in (?)", sceneNames).
				Find(&pbDataList)
			var dataNames []string
			for _, pd := range pbDataList {
				if len(pd.DataName) > 0 {
					parts := strings.Split(pd.DataName, ",")
					dataNames = append(dataNames, parts...)
				}
			}
			// 将数据文件名加入 datas 返回
			for _, dn := range dataNames {
				dn = strings.TrimSpace(dn)
				if len(dn) > 0 {
					datas = append(datas, ResourceItem{Name: dn})
				}
			}
			// 从 scene_data 查 API
			if len(dataNames) > 0 {
				var apiIds []string
				models.Orm.Table("scene_data").
					Where("file_name in (?)", dataNames).
					Group("api_id").
					Pluck("api_id", &apiIds)
				for _, apiId := range apiIds {
					apis = append(apis, ResourceItem{Name: apiId})
				}
			}
		}
	} else if t.TaskType == "data" && len(t.DataList) > 0 {
		dataNames := parseLineList(t.DataList)
		cleanNames := make([]string, 0, len(dataNames))
		for _, name := range dataNames {
			name = strings.TrimSpace(name)
			if len(name) == 0 {
				continue
			}
			datas = append(datas, ResourceItem{Name: name})
			cleanNames = append(cleanNames, name)
		}
		// 从 scene_data 查 API
		if len(cleanNames) > 0 {
			var apiIds []string
			models.Orm.Table("scene_data").
				Where("file_name in (?)", cleanNames).
				Group("api_id").
				Pluck("api_id", &apiIds)
			for _, apiId := range apiIds {
				apis = append(apis, ResourceItem{Name: apiId})
			}
		}
	}
	return
}

// parseLineList 解析换行分隔的列表
func parseLineList(listStr string) []string {
	if len(listStr) == 0 {
		return nil
	}
	// 尝试逗号分隔
	if strings.Contains(listStr, ",") {
		return strings.Split(listStr, ",")
	}
	// 尝试换行分隔
	return strings.Split(listStr, "\r\n")
}

// mapToResourceItems 将 set 转为 ResourceItem 列表
func mapToResourceItems(m map[string]bool) []ResourceItem {
	items := make([]ResourceItem, 0, len(m))
	for name := range m {
		items = append(items, ResourceItem{Name: name})
	}
	return items
}

// countResources 统计场景明细中的资源计数
func countResources(details []SceneDetail) []ResourceItem {
	countMap := make(map[string]int)
	for _, d := range details {
		countMap[d.Name]++
	}
	items := make([]ResourceItem, 0, len(countMap))
	for name, count := range countMap {
		items = append(items, ResourceItem{Name: name, Count: count})
	}
	return items
}

// countResourcesForTaskScene 统计场景明细资源（SceneDetailWithTask 版本）
func countResourcesForTaskScene(details []SceneDetailWithTask) []ResourceItem {
	countMap := make(map[string]int)
	for _, d := range details {
		countMap[d.Name]++
	}
	items := make([]ResourceItem, 0, len(countMap))
	for name, count := range countMap {
		items = append(items, ResourceItem{Name: name, Count: count})
	}
	return items
}

// mergeResourceItems 合并资源列表（去重）
func mergeResourceItems(a, b []ResourceItem) []ResourceItem {
	seen := make(map[string]bool)
	var result []ResourceItem
	for _, item := range a {
		if !seen[item.Name] {
			seen[item.Name] = true
			result = append(result, item)
		}
	}
	for _, item := range b {
		if !seen[item.Name] {
			seen[item.Name] = true
			result = append(result, item)
		}
	}
	return result
}

// mergeCountItems 合并 CountItem 列表（按 Name 聚合）
func mergeCountItems(a, b []CountItem) []CountItem {
	countMap := make(map[string]int)
	for _, item := range a {
		countMap[item.Name] += item.Count
	}
	for _, item := range b {
		countMap[item.Name] += item.Count
	}
	var result []CountItem
	for name, count := range countMap {
		result = append(result, CountItem{Name: name, Count: count})
	}
	return result
}

// querySceneDetailsForTask 按taskId和product查询场景明细
func querySceneDetailsForTask(taskId, product string) (items []SceneDetail) {
	if len(taskId) == 0 {
		return
	}
	type sceneResult struct {
		Name       string `gorm:"column:name"`
		Result     string `gorm:"column:result"`
		FailReason string `gorm:"column:fail_reason"`
	}
	var results []sceneResult
	models.Orm.Table("scene_test_history").
		Select("name, result, fail_reason").
		Where("task_id = ? and product = ?", taskId, product).
		Find(&results)
	for _, r := range results {
		items = append(items, SceneDetail{
			Name:       r.Name,
			Result:     r.Result,
			FailReason: r.FailReason,
		})
	}
	return
}

// queryDataDetailsForTask 按taskId和product查询数据文件明细
func queryDataDetailsForTask(taskId, product string) (items []DataDetail) {
	if len(taskId) == 0 {
		return
	}
	type dataResult struct {
		Name       string `gorm:"column:name"`
		ApiId      string `gorm:"column:api_id"`
		Result     string `gorm:"column:result"`
		FailReason string `gorm:"column:fail_reason"`
	}
	var results []dataResult
	models.Orm.Table("scene_data_test_history").
		Select("name, api_id, result, fail_reason").
		Where("task_id = ? and product = ?", taskId, product).
		Find(&results)
	for _, r := range results {
		items = append(items, DataDetail{
			Name:       r.Name,
			ApiId:      r.ApiId,
			Result:     r.Result,
			FailReason: r.FailReason,
		})
	}
	return
}

// queryFailItemsForTask 按taskId和product查询失败项
func queryFailItemsForTask(taskId, product string) (items []FailItem) {
	if len(taskId) == 0 {
		return
	}
	// 场景失败项
	type sceneFail struct {
		Name       string `gorm:"column:name"`
		FailReason string `gorm:"column:fail_reason"`
	}
	var sfList []sceneFail
	models.Orm.Table("scene_test_history").
		Select("name, fail_reason").
		Where("task_id = ? and product = ? and result = ?", taskId, product, "fail").
		Find(&sfList)
	for _, sf := range sfList {
		items = append(items, FailItem{
			Name:   sf.Name,
			Type:   "scene",
			Reason: sf.FailReason,
		})
	}
	// 数据失败项
	type dataFail struct {
		Name       string `gorm:"column:name"`
		ApiId      string `gorm:"column:api_id"`
		FailReason string `gorm:"column:fail_reason"`
	}
	var dfList []dataFail
	models.Orm.Table("scene_data_test_history").
		Select("name, api_id, fail_reason").
		Where("task_id = ? and product = ? and result = ?", taskId, product, "fail").
		Find(&dfList)
	for _, df := range dfList {
		items = append(items, FailItem{
			Name:   df.Name,
			Type:   "data",
			APIId:  df.ApiId,
			Reason: df.FailReason,
		})
	}
	return
}

// queryAPITypeDistributionForTask 按taskId和product查询API类型分布
// 与单任务 queryAPITypeDistribution 逻辑一致，查 api_definition 表获取真实 http_method
func queryAPITypeDistributionForTask(taskId, product string) (items []CountItem) {
	if len(taskId) == 0 {
		return
	}

	var apiList []string
	models.Orm.Table("scene_data_test_history").
		Group("api_id").
		Where("task_id = ? and product = ? and api_id IS NOT NULL and api_id <> ''", taskId, product).
		Pluck("api_id", &apiList)
	if len(apiList) == 0 {
		return
	}

	var httpMethods []string
	models.Orm.Table("api_definition").
		Group("http_method").
		Where("api_id in (?)", apiList).
		Pluck("http_method", &httpMethods)

	for _, method := range httpMethods {
		var cntResult struct {
			Cnt int64
		}
		models.Orm.Raw(
			"SELECT COUNT(DISTINCT api_id) as cnt FROM api_definition WHERE http_method = ? AND api_id IN (?)",
			method, apiList).Scan(&cntResult)
		methodName := method
		cnt := cntResult.Cnt
		if len(methodName) == 0 {
			methodName = "其他"
		}
		items = append(items, CountItem{Name: methodName, Count: int(cnt)})
	}

	if len(items) == 0 {
		var total int64
		models.Orm.Table("scene_data_test_history").
			Where("task_id = ? and product = ?", taskId, product).Count(&total)
		if total > 0 {
			items = append(items, CountItem{Name: "全部", Count: int(total)})
		}
	}
	return
}

// GenerateManualReport 手动生成全局/产品/应用报告
func GenerateManualReport(reportType, products, genUser string) (err error) {
	now := time.Now().Format(baseFormat)
	nowStr := time.Now().Format("20060102150405")

	// 汇总数据：复用已有的统计查询函数，异步执行
	var reportData map[string]interface{}

	switch reportType {
	case "global":
		// 全局报告：聚合所有产品的数据
		reportData = make(map[string]interface{})
		reportData["type"] = "global"
		reportData["generated_at"] = now
		reportData["time_range"] = "最近6个月"

		// API类型分布
		apiMethods, apiCounts, _, _ := GetAPITypeCount("all", "")
		reportData["api_type_distribution"] = buildCountMap(apiMethods, apiCounts)

		// API规范检查
		specInfos, specCounts, _, _ := GetAPISpecCount("all", "")
		reportData["api_spec_check"] = buildCountMap(specInfos, specCounts)

		// API自动化覆盖
		autoInfos, autoCounts, _, _ := GetAutoAPICount("all", "")
		reportData["api_auto_coverage"] = buildCountMap(autoInfos, autoCounts)

		// 场景结果统计
		sceneInfos, sceneCounts, _, _ := GetSceneResultCount()
		reportData["scene_result_distribution"] = buildCountMap(sceneInfos, sceneCounts)

		// 数据结果统计
		dataInfos, dataCounts, _, _ := GetSceneDataResultCount()
		reportData["data_result_distribution"] = buildCountMap(dataInfos, dataCounts)

		// 任务类型分布
		schedInfos, schedCounts, _, _ := GetScheduleTypeCount()
		reportData["schedule_type_distribution"] = buildCountMap(schedInfos, schedCounts)

		// 数据执行趋势(6个月)
		_, dataMonthLabels, dataInfos2, dataCounts2 := GetAppSceneDataRunCount()
		reportData["data_trend"] = buildTrendData(dataMonthLabels, dataInfos2, dataCounts2)

		// 数据执行分布(6个月)
		appRunInfos, appRunCounts, _, _ := GetAppAPIRunCount()
		reportData["data_exec_dist"] = buildCountMap(appRunInfos, appRunCounts)

		// 场景执行趋势(6个月)
		_, monthLabels, sceneInfos2, sceneCounts2 := GetProductSceneRunCount()
		reportData["scene_trend"] = buildTrendData(monthLabels, sceneInfos2, sceneCounts2)

		// 场景执行分布(6个月)
		sceneRunInfos, sceneRunCounts, _, _ := GetSceneRunCount()
		reportData["scene_exec_dist"] = buildCountMap(sceneRunInfos, sceneRunCounts)
	}

	jsonBytes, err := json.Marshal(reportData)
	if err != nil {
		return err
	}

	reportName := fmt.Sprintf("%s报告_%s", reportType, nowStr)
	report := DashboardReport{
		ReportName:      reportName,
		ReportType:      reportType,
		RelatedProducts: products,
		RelatedApps:     products,
		TimeRangeStart:  now,
		TimeRangeEnd:    now,
		Status:          "finished",
		Creator:         genUser,
		ReportData:      string(jsonBytes),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = models.Orm.Table("dashboard").Create(&report).Error
	if err != nil {
		return fmt.Errorf("保存报告失败: %s", err)
	}

	return nil
}

// buildCountMap 构建计数映射
func buildCountMap(infos []string, counts []float64) []map[string]interface{} {
	result := make([]map[string]interface{}, len(infos))
	for i, info := range infos {
		count := 0.0
		if i < len(counts) {
			count = counts[i]
		}
		result[i] = map[string]interface{}{
			"name":  info,
			"count": count,
		}
	}
	return result
}

// buildTableData 构建表格数据
func buildTableData(contents []map[string]types.InfoItem, headers types.Thead) map[string]interface{} {
	return map[string]interface{}{
		"headers": headers,
		"rows":    contents,
	}
}

// buildTrendData 构建趋势数据
func buildTrendData(labels []string, infos []string, counts [][]float64) map[string]interface{} {
	return map[string]interface{}{
		"labels": labels,
		"series": infos,
		"data":   counts,
	}
}

// getTaskEnvironment 获取任务执行环境描述
func getTaskEnvironment(productList string) string {
	if len(productList) == 0 {
		return ""
	}
	products := parseProductList(productList)
	if len(products) == 0 {
		return ""
	}

	for _, p := range products {
		var dbProducts []DbProduct
		models.Orm.Table("product").Where("product = ?", p).Find(&dbProducts)
		return fmt.Sprintf("%s", p)
	}
	return ""
}

// queryAPITypeDistribution 按taskId查询API类型分布
func queryAPITypeDistribution(taskId string) (items []CountItem) {
	if len(taskId) == 0 {
		return
	}
	type apiMethodResult struct {
		Method string
		Count  int
	}
	var results []apiMethodResult

	var apiList, httpMethods []string
	models.Orm.Table("scene_data_test_history").Group("api_id").Where("task_id = ?", taskId).Pluck("api_id", &apiList)
	models.Orm.Table("api_definition").Group("http_method").Where("api_id in (?)", apiList).Pluck("http_method", &httpMethods)
	for _, method := range httpMethods {
		var methodCount apiMethodResult
		methodCount.Method = method
		models.Orm.Raw(
			"SELECT COUNT(DISTINCT api_id) as count FROM api_definition WHERE http_method = ? AND api_id IN (?)",
			method, apiList).Scan(&methodCount)
		results = append(results, methodCount)
	}

	if len(results) == 0 {
		var total int64
		models.Orm.Table("scene_data_test_history").Where("task_id = ?", taskId).Count(&total)
		if total > 0 {
			items = append(items, CountItem{Name: "全部", Count: int(total)})
		}
		return
	}

	for _, r := range results {
		method := r.Method
		if len(method) == 0 {
			method = "其他"
		}
		items = append(items, CountItem{Name: method, Count: r.Count})
	}
	return
}

// querySceneDetails 按taskId查询场景明细
func querySceneDetails(taskId string) (items []SceneDetail) {
	if len(taskId) == 0 {
		return
	}
	type sceneResult struct {
		Name       string `gorm:"column:name"`
		Result     string `gorm:"column:result"`
		FailReason string `gorm:"column:fail_reason"`
	}
	var results []sceneResult
	models.Orm.Table("scene_test_history").
		Select("name, result, fail_reason").
		Where("task_id = ?", taskId).
		Order("created_at asc").
		Find(&results)
	for _, r := range results {
		items = append(items, SceneDetail{
			Name:       r.Name,
			Result:     r.Result,
			FailReason: r.FailReason,
		})
	}
	return
}

// queryDataDetails 按taskId查询数据文件执行明细
func queryDataDetails(taskId string) (items []DataDetail) {
	if len(taskId) == 0 {
		return
	}
	type dataResult struct {
		Name       string `gorm:"column:name"`
		ApiId      string `gorm:"column:api_id"`
		Result     string `gorm:"column:result"`
		FailReason string `gorm:"column:fail_reason"`
	}
	var results []dataResult
	models.Orm.Table("scene_data_test_history").
		Select("name, api_id, result, fail_reason").
		Where("task_id = ?", taskId).
		Find(&results)
	for _, r := range results {
		items = append(items, DataDetail{
			Name:       r.Name,
			ApiId:      r.ApiId,
			Result:     r.Result,
			FailReason: r.FailReason,
		})
	}
	return
}
