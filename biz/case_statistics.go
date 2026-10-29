package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// caseStatRow 统计所需的用例字段
type caseStatRow struct {
	Module       string `gorm:"column:module"`
	CaseNumber   string `gorm:"column:case_number"`
	CaseName     string `gorm:"column:case_name"`
	TestResult   string `gorm:"column:test_result"`
	TestTime     string `gorm:"column:test_time"`
	FunDeveloper string `gorm:"column:fun_developer"`
	CaseDesigner string `gorm:"column:case_designer"`
	Priority     string `gorm:"column:priority"`
	CaseExecutor string `gorm:"column:case_executor"`
	ExtInfo      string `gorm:"column:ext_info"`
	Remark       string `gorm:"column:remark"`
}

// ExtEmptyPlaceholder 统计维度中空值的占位显示（扩展字段/优先级/执行者等）
const ExtEmptyPlaceholder = "-"

// GetIntroVersions 获取 test_case 中所有非空引入版（用于统计定义表单多选）
func GetIntroVersions() (versions []string) {
	models.Orm.Table("test_case").
		Where("intro_version IS NOT NULL AND intro_version <> ''").
		Group("intro_version").
		Order("intro_version").
		Pluck("intro_version", &versions)
	return
}

// GetCaseStatisticsName 根据统计定义 id 获取名称（报告列表展示关联统计定义用）
func GetCaseStatisticsName(id int) string {
	var names []string
	models.Orm.Table("case_statistics").Where("id = ?", id).Pluck("name", &names)
	if len(names) > 0 {
		return names[0]
	}
	return ""
}

// ParseCaseStatisticsDefinition 解析统计定义 YAML
// 支持格式（ext_info 为保留键，可选；其余顶层键 = 一级菜单名）：
//
//	ext_info:              # 可选：扩展信息中要统计的字段 key
//	  - priority
//	一级菜单A:
//	  二级菜单A1:
//	    - 模块名1
//	    - 模块名2
//	一级菜单B:             # 允许跳过二级，一级直接挂模块
//	  - 模块名3
func ParseCaseStatisticsDefinition(definition string) (CaseStatisticsDefinition, error) {
	var raw yaml.MapSlice
	if err := yaml.Unmarshal([]byte(definition), &raw); err != nil {
		return CaseStatisticsDefinition{}, err
	}
	var def CaseStatisticsDefinition
	for _, item := range raw {
		key := strings.TrimSpace(fmt.Sprintf("%v", item.Key))
		if key == "" {
			continue
		}
		if key == "ext_info" {
			def.ExtInfo = yamlStringList(item.Value)
			continue
		}
		node := CaseStatTreeNode{Name: key}
		switch v := item.Value.(type) {
		case yaml.MapSlice:
			for _, l2 := range v {
				l2Name := strings.TrimSpace(fmt.Sprintf("%v", l2.Key))
				if l2Name == "" {
					continue
				}
				node.Children = append(node.Children, CaseStatTreeNode{
					Name:    l2Name,
					Modules: yamlStringList(l2.Value),
				})
			}
		case []interface{}:
			node.Modules = yamlStringList(v)
		}
		def.Tree = append(def.Tree, node)
	}
	return def, nil
}

// yamlStringList 将 YAML 值转为 []string（非列表/空值返回 nil）
func yamlStringList(v interface{}) []string {
	list, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, e := range list {
		if s := strings.TrimSpace(fmt.Sprintf("%v", e)); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// ParseCaseStatisticsPlan 解析计划排期配置（case_statistics.plan YAML）
func ParseCaseStatisticsPlan(plan string) (CaseStatisticsPlan, error) {
	var p CaseStatisticsPlan
	if err := yaml.Unmarshal([]byte(plan), &p); err != nil {
		return p, err
	}
	return p, nil
}

// GenerateCaseStatisticsReport 生成用例统计报告（异步调用）
func GenerateCaseStatisticsReport(caseStatId int, userName string) {
	var cs CaseStatistics
	models.Orm.Table("case_statistics").Where("id = ?", caseStatId).Find(&cs)
	if cs.Id == 0 {
		Logger.Error("用例统计定义不存在: %d", caseStatId)
		return
	}

	models.Orm.Table("case_statistics").Where("id = ?", caseStatId).
		UpdateColumn("status", "generating")

	def, err := ParseCaseStatisticsDefinition(cs.Definition)
	if err != nil {
		Logger.Error("解析用例统计定义失败: %s", err)
		models.Orm.Table("case_statistics").Where("id = ?", caseStatId).
			UpdateColumn("status", "failed")
		return
	}

	// 查询用例数据（按引入版/产品线过滤，逗号分隔多选）
	query := models.Orm.Table("test_case").Where("deleted_at IS NULL")
	if versions := splitTrim(cs.IntroVersions); len(versions) > 0 {
		query = query.Where("intro_version IN (?)", versions)
	}
	if products := splitTrim(cs.Product); len(products) > 0 {
		query = query.Where("product IN (?)", products)
	}
	var rows []caseStatRow
	if err := query.Select("module, case_number, case_name, test_result, test_time, fun_developer, case_designer, priority, case_executor, ext_info, remark").
		Find(&rows).Error; err != nil {
		Logger.Error("查询测试用例失败: %s", err)
		models.Orm.Table("case_statistics").Where("id = ?", caseStatId).
			UpdateColumn("status", "failed")
		return
	}

	plan, _ := ParseCaseStatisticsPlan(cs.Plan) // 解析失败/为空 → 零值 → 计划块不启用
	reportData := buildCaseStatisticsReport(def, cs.IntroVersions, rows, plan)
	reportData.Overview.Name = cs.Name
	reportData.Overview.Product = cs.Product

	// 统计时间范围：所选范围内测试用例中最早/最晚的测试时间
	minTestTime, maxTestTime := "", ""
	for _, r := range rows {
		t := strings.TrimSpace(r.TestTime)
		if t == "" {
			continue
		}
		if minTestTime == "" || t < minTestTime {
			minTestTime = t
		}
		if maxTestTime == "" || t > maxTestTime {
			maxTestTime = t
		}
	}

	jsonBytes, err := json.Marshal(reportData)
	if err != nil {
		Logger.Error("序列化用例统计报告失败: %s", err)
		models.Orm.Table("case_statistics").Where("id = ?", caseStatId).
			UpdateColumn("status", "failed")
		return
	}

	// 写入用例统计报告表
	now := time.Now().Format(baseFormat)
	nowStr := time.Now().Format("20060102150405")
	report := CaseStatisticsReport{
		Name:          fmt.Sprintf("%s_%s", cs.Name, nowStr),
		RelatedStatId: cs.Id,
		IntroVersions: cs.IntroVersions,
		Product:       cs.Product,
		StatStartTime: minTestTime,
		StatEndTime:   maxTestTime,
		Status:        "finished",
		Creator:       userName,
		ReportData:    string(jsonBytes),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := models.Orm.Table("case_statistics_report").Create(&report).Error; err != nil {
		Logger.Error("保存用例统计报告失败: %s", err)
		models.Orm.Table("case_statistics").Where("id = ?", caseStatId).
			UpdateColumn("status", "failed")
		return
	}

	models.Orm.Table("case_statistics").Where("id = ?", caseStatId).
		UpdateColumn("status", "finished")
	Logger.Info("用例统计报告生成完成: %s", cs.Name)
}

// buildCaseStatisticsReport 聚合生成统计报告数据
func buildCaseStatisticsReport(def CaseStatisticsDefinition, introVersions string, rows []caseStatRow, plan CaseStatisticsPlan) CaseStatisticsReportData {
	var report CaseStatisticsReportData
	report.Overview.IntroVersions = introVersions

	// 建立 模块名 -> (一级, 二级) 路径映射（树内先出现者优先）
	type path struct{ Level1, Level2 string }
	modulePath := make(map[string]path)
	for _, l1 := range def.Tree {
		for _, m := range l1.Modules {
			if m = strings.TrimSpace(m); m == "" {
				continue
			}
			if _, exists := modulePath[m]; !exists {
				modulePath[m] = path{Level1: l1.Name, Level2: ""}
			}
		}
		for _, l2 := range l1.Children {
			for _, m := range l2.Modules {
				if m = strings.TrimSpace(m); m == "" {
					continue
				}
				if _, exists := modulePath[m]; !exists {
					modulePath[m] = path{Level1: l1.Name, Level2: l2.Name}
				}
			}
		}
	}

	moduleMap := make(map[string]*CaseStatModuleItem)
	resultCount := map[string]int{}
	devCount := map[string]int{}
	designerCount := map[string]int{}
	priorityCount := map[string]int{}
	executorTotal := map[string]int{}
	executorPass := map[string]int{}
	executorFail := map[string]int{}
	executorDeprecated := map[string]int{}
	executorUnmerged := map[string]int{}
	executorUntest := map[string]int{}
	extCount := map[string]map[string]int{}
	var failCases, untestCases []CaseStatCaseDetail
	var matchedRows []caseStatRow
	var totalCases int

	seen := make(map[string]bool)
	var order []string

	for _, r := range rows {
		module := strings.TrimSpace(r.Module)
		p, ok := modulePath[module]
		if !ok {
			continue // 未关联到统计定义的模块不纳入统计
		}
		matchedRows = append(matchedRows, r)
		key := p.Level1 + "\x00" + p.Level2 + "\x00" + module
		if !seen[key] {
			seen[key] = true
			order = append(order, key)
		}
		item := moduleMap[key]
		if item == nil {
			item = &CaseStatModuleItem{Level1: p.Level1, Level2: p.Level2, Module: module}
			moduleMap[key] = item
		}

		item.Total++
		totalCases++
		switch r.TestResult {
		case "pass":
			item.Pass++
			resultCount["pass"]++
		case "part":
			item.Part++
			resultCount["part"]++
		case "fail":
			item.Fail++
			resultCount["fail"]++
		case "deprecated":
			item.Deprecated++
			resultCount["deprecated"]++
		case "unmerged":
			item.Unmerged++
			resultCount["unmerged"]++
		default:
			item.Untest++
			resultCount["untest"]++
		}

		dev := strings.TrimSpace(r.FunDeveloper)
		if dev == "" {
			dev = ExtEmptyPlaceholder
		}
		devCount[dev]++

		des := strings.TrimSpace(r.CaseDesigner)
		if des == "" {
			des = ExtEmptyPlaceholder
		}
		designerCount[des]++

		// 优先级计数（空值归入占位符）
		prio := strings.TrimSpace(r.Priority)
		if prio == "" {
			prio = ExtEmptyPlaceholder
		}
		priorityCount[prio]++

		// 用例执行者计数（空值归入占位符，同时用于饼图与人员统计表）
		exec := strings.TrimSpace(r.CaseExecutor)
		if exec == "" {
			exec = ExtEmptyPlaceholder
		}
		executorTotal[exec]++
		switch r.TestResult {
		case "pass":
			executorPass[exec]++
		case "fail":
			executorFail[exec]++
		case "deprecated":
			executorDeprecated[exec]++
		case "unmerged":
			executorUnmerged[exec]++
		}
		if isUntestResult(r.TestResult) {
			executorUntest[exec]++
		}

		for _, k := range def.ExtInfo {
			if k = strings.TrimSpace(k); k == "" {
				continue
			}
			v := GetExtInfoValue(r.ExtInfo, k, "zh-CN")
			if v == "" {
				v = ExtEmptyPlaceholder
			}
			if extCount[k] == nil {
				extCount[k] = map[string]int{}
			}
			extCount[k][v]++
		}

		if r.TestResult == "fail" {
			failCases = append(failCases, CaseStatCaseDetail{
				CaseNumber: r.CaseNumber, CaseName: r.CaseName, Module: module, Remark: r.Remark,
			})
		} else if r.TestResult == "untest" {
			untestCases = append(untestCases, CaseStatCaseDetail{
				CaseNumber: r.CaseNumber, CaseName: r.CaseName, Module: module, Remark: r.Remark,
			})
		}
	}

	for _, key := range order {
		if item := moduleMap[key]; item != nil {
			if item.Total > 0 {
				item.PassRate = float64(item.Pass) / float64(item.Total) * 100
				item.ExecRate = float64(item.Total-item.Untest) / float64(item.Total) * 100
			}
			report.ModuleStats = append(report.ModuleStats, *item)
		}
	}

	report.Overview.ModuleCount = len(report.ModuleStats)
	report.Overview.TotalCases = totalCases
	report.Overview.Pass = resultCount["pass"]
	report.Overview.Fail = resultCount["fail"]
	report.Overview.Untest = resultCount["untest"]
	report.Overview.Part = resultCount["part"]
	report.Overview.Deprecated = resultCount["deprecated"]
	report.Overview.Unmerged = resultCount["unmerged"]
	if report.Overview.TotalCases > 0 {
		report.Overview.PassRate = float64(report.Overview.Pass) / float64(report.Overview.TotalCases) * 100
		report.Overview.ExecRate = float64(report.Overview.TotalCases-report.Overview.Untest) / float64(report.Overview.TotalCases) * 100
	}

	for _, r := range []string{"pass", "fail", "deprecated", "unmerged"} {
		report.ResultDistribution = append(report.ResultDistribution, CountItem{Name: r, Count: resultCount[r]})
	}
	report.ByFunDeveloper = sortCountMap(devCount)
	report.ByCaseDesigner = sortCountMap(designerCount)
	report.ByPriority = sortCountMap(priorityCount)
	report.ByCaseExecutor = sortCountMap(executorTotal)

	// 用例执行者统计（含已执行/未执行/执行率），按总用例数降序
	executorNames := make([]string, 0, len(executorTotal))
	for name := range executorTotal {
		executorNames = append(executorNames, name)
	}
	sort.Slice(executorNames, func(i, j int) bool {
		// 空值占位符始终排在最后
		if executorNames[i] == ExtEmptyPlaceholder {
			return false
		}
		if executorNames[j] == ExtEmptyPlaceholder {
			return true
		}
		if executorTotal[executorNames[i]] != executorTotal[executorNames[j]] {
			return executorTotal[executorNames[i]] > executorTotal[executorNames[j]]
		}
		return executorNames[i] < executorNames[j]
	})
	for _, name := range executorNames {
		total := executorTotal[name]
		untest := executorUntest[name]
		item := CaseStatExecutorItem{
			Name:       name,
			Total:      total,
			Executed:   total - untest,
			Pass:       executorPass[name],
			Fail:       executorFail[name],
			Deprecated: executorDeprecated[name],
			Unmerged:   executorUnmerged[name],
			Untest:     untest,
		}
		if total > 0 {
			item.ExecRate = float64(total-untest) / float64(total) * 100
		}
		report.ExecutorStats = append(report.ExecutorStats, item)
	}

	extKeys := make([]string, 0, len(extCount))
	for k := range extCount {
		extKeys = append(extKeys, k)
	}
	sort.Strings(extKeys)
	for _, k := range extKeys {
		report.ExtInfoStats = append(report.ExtInfoStats, CaseStatExtInfoItem{
			Key:   k,
			Items: sortCountMap(extCount[k]),
		})
	}

	report.FailCases = failCases
	report.UntestCases = untestCases
	report.TrendDaily, report.TrendWeekly = buildCaseStatisticsTrend(matchedRows)
	report.PlanProgress = buildCaseStatPlanProgress(plan, matchedRows, totalCases)
	return report
}

// caseStatTrendClass 趋势归类：pass/fail/other；空或 untest（未执行）返回空串。
func caseStatTrendClass(result string) string {
	switch result {
	case "pass":
		return "pass"
	case "fail":
		return "fail"
	case "part", "deprecated", "unmerged":
		return "other"
	}
	return ""
}

// trendBucket 单个时间桶的新增计数
type trendBucket struct {
	pass, fail, other int
}

// buildCaseStatisticsTrend 构建测试执行趋势（按日/按周两套桶）。
// 纯快照口径：每个用例只取最新 test_result + test_time；未执行（空或 untest）不参与。
// 空 test_time 或无法解析的时间直接跳过。
func buildCaseStatisticsTrend(rows []caseStatRow) (daily, weekly []CaseStatTrendPoint) {
	dayMap := map[string]*trendBucket{}
	weekMap := map[string]*trendBucket{}
	for _, r := range rows {
		cls := caseStatTrendClass(r.TestResult)
		if cls == "" {
			continue
		}
		t := strings.TrimSpace(r.TestTime)
		if len(t) < len("2006-01-02") {
			continue
		}
		tt, err := time.Parse(baseFormat, t)
		if err != nil {
			continue
		}
		day := t[:10]
		week := tt.AddDate(0, 0, -((int(tt.Weekday())+6)%7)).Format("2006-01-02")

		db := dayMap[day]
		if db == nil {
			db = &trendBucket{}
			dayMap[day] = db
		}
		wb := weekMap[week]
		if wb == nil {
			wb = &trendBucket{}
			weekMap[week] = wb
		}
		switch cls {
		case "pass":
			db.pass++
			wb.pass++
		case "fail":
			db.fail++
			wb.fail++
		case "other":
			db.other++
			wb.other++
		}
	}
	daily = trendCumulate(dayMap)
	weekly = trendCumulate(weekMap)
	return
}

// trendCumulate 将分桶计数按桶名（YYYY-MM-DD）升序做前缀累计
func trendCumulate(m map[string]*trendBucket) []CaseStatTrendPoint {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var passCum, failCum, otherCum int
	out := make([]CaseStatTrendPoint, 0, len(keys))
	for _, k := range keys {
		b := m[k]
		passCum += b.pass
		failCum += b.fail
		otherCum += b.other
		out = append(out, CaseStatTrendPoint{
			Label:    k,
			ExecCum:  passCum + failCum + otherCum,
			PassCum:  passCum,
			FailCum:  failCum,
			OtherCum: otherCum,
			PassNew:  b.pass,
			FailNew:  b.fail,
			OtherNew: b.other,
		})
	}
	return out
}

// isUntestResult 判断结果是否属于未执行（非 pass/part/fail/deprecated/unmerged 均视为未执行）
func isUntestResult(result string) bool {
	switch result {
	case "pass", "part", "fail", "deprecated", "unmerged":
		return false
	}
	return true
}

// sortCountMap 将计数 map 转为按数量降序的 CountItem 列表
func sortCountMap(m map[string]int) []CountItem {
	items := make([]CountItem, 0, len(m))
	for k, v := range m {
		items = append(items, CountItem{Name: k, Count: v})
	}
	sort.Slice(items, func(i, j int) bool {
		// 空值占位符始终排在最后
		if items[i].Name == ExtEmptyPlaceholder {
			return false
		}
		if items[j].Name == ExtEmptyPlaceholder {
			return true
		}
		if items[i].Count != items[j].Count {
			return items[i].Count > items[j].Count
		}
		return items[i].Name < items[j].Name
	})
	return items
}

// splitTrim 拆分逗号分隔字符串并去空
func splitTrim(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

// stringSet 将 []string 转为 map 用于 O(1) 命中判断
func stringSet(list []string) map[string]bool {
	set := make(map[string]bool, len(list))
	for _, s := range list {
		if s = strings.TrimSpace(s); s != "" {
			set[s] = true
		}
	}
	return set
}

// caseStatPlanIsWorkDay 判定某日期是否工作日（优先级：workdays > holidays > rest_mode）
func caseStatPlanIsWorkDay(d time.Time, restMode string, holidays, workdays map[string]bool) bool {
	dateStr := d.Format("2006-01-02")
	if workdays[dateStr] {
		return true // 补班/加班日
	}
	if holidays[dateStr] {
		return false // 节假日
	}
	switch restMode {
	case "none":
		return true
	case "sun": // 单休，休周日
		return d.Weekday() != time.Sunday
	case "sat": // 单休，休周六
		return d.Weekday() != time.Saturday
	default: // weekend / 未识别：双休，去周六日
		return d.Weekday() != time.Saturday && d.Weekday() != time.Sunday
	}
}

// buildCaseStatPlanProgress 构建计划 vs 实际进度（plan 为空/非法 → Enabled=false）
func buildCaseStatPlanProgress(plan CaseStatisticsPlan, rows []caseStatRow, totalCases int) CaseStatPlanProgress {
	progress := CaseStatPlanProgress{
		StartDate:  plan.StartDate,
		EndDate:    plan.EndDate,
		DailyCount: plan.DailyCount,
		RestMode:   plan.RestMode,
		Total:      totalCases,
	}

	start, err1 := time.Parse("2006-01-02", plan.StartDate)
	end, err2 := time.Parse("2006-01-02", plan.EndDate)
	if err1 != nil || err2 != nil || start.After(end) || totalCases == 0 {
		return progress // 日期缺失/非法/无用例 → 未启用
	}

	holidays := stringSet(plan.Holidays)
	workdays := stringSet(plan.Workdays)

	// 工作日数与每日步长
	workDays := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if caseStatPlanIsWorkDay(d, plan.RestMode, holidays, workdays) {
			workDays++
		}
	}
	if workDays == 0 {
		return progress // 无工作日 → 未启用
	}
	daily := plan.DailyCount
	if daily <= 0 {
		daily = (totalCases + workDays - 1) / workDays // ceil(总量/工作日)
	}

	// 实际每日执行量 + 已完成(已执行)用例数
	actualDaily := map[string]int{}
	executed := 0
	for _, r := range rows {
		if isUntestResult(r.TestResult) {
			continue
		}
		executed++
		t := strings.TrimSpace(r.TestTime)
		if len(t) >= len("2006-01-02") {
			actualDaily[t[:10]]++
		}
	}

	// 逐日累计
	progress.Enabled = true
	progress.WorkDays = workDays
	var planCum, actCum int
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		work := caseStatPlanIsWorkDay(d, plan.RestMode, holidays, workdays)
		pd, ad := 0, 0
		if work {
			pd = daily
			ad = actualDaily[dateStr]
			planCum += pd
			if planCum > totalCases {
				planCum = totalCases
			}
			actCum += ad
			if actCum > totalCases {
				actCum = totalCases
			}
			if progress.PlannedFinishDate == "" && planCum >= totalCases {
				progress.PlannedFinishDate = dateStr
			}
		}
		progress.Points = append(progress.Points, CaseStatPlanPoint{
			Date:         dateStr,
			IsWorkDay:    work,
			PlannedDaily: pd,
			ActualDaily:  ad,
			PlannedCum:   planCum,
			ActualCum:    actCum,
		})
	}
	// 实际累计曲线仅在工作日推进，无时间戳/休息日/窗口外的已执行用例未体现在曲线中，
	// 补齐到最后一个点，保证曲线终点与执行进度一致（同样受总量上限约束）
	if len(progress.Points) > 0 {
		lastCum := executed
		if lastCum > totalCases {
			lastCum = totalCases
		}
		progress.Points[len(progress.Points)-1].ActualCum = lastCum
	}
	if totalCases > 0 {
		progress.Progress = float64(executed) / float64(totalCases) * 100
	}
	return progress
}
