package biz

import (
	"data4test/models"
	"data4test/sorting"
	"fmt"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetAPITypeCount(mode, appName string) (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {
	appList := strings.Split(appName, ",")
	if mode == "all" {
		models.Orm.Table("api_definition").Group("http_method").Pluck("http_method", &infos)
	} else {
		models.Orm.Table("api_definition").Group("http_method").Where("app in (?)", appList).Pluck("http_method", &infos)
	}

	if len(infos) == 0 {
		infos = []string{T("report.no_data")}
		counts = []float64{0}
		colors = []chartjs.Color{"rgb(255, 205, 86)"}
		return
	}

	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}
	colorNames := []string{"yellow", "blue", "red", "green", "black"}

	for index, item := range infos {
		var apiCount float64
		labelInfo := make(map[string]string)
		if mode == "all" {
			models.Orm.Table("api_definition").Where("http_method = ?", item).Count(&apiCount)
		} else if mode == "product" || mode == "app" {
			models.Orm.Table("api_definition").Where("http_method = ? and app in (?)", item, appList).Count(&apiCount)
		}

		counts = append(counts, apiCount)
		labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(apiCount))
		if index < len(defineColors) {
			colors = append(colors, defineColors[index])
			labelInfo["color"] = colorNames[index]
		} else {
			colors = append(colors, defineColors[0])
			labelInfo["color"] = colorNames[0]
		}

		labels = append(labels, labelInfo)

	}

	return
}

func GetAPISpecCount(mode, appName string) (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {
	infos = []string{"pass", "fail", "unknown"}
	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}
	colorNames := []string{"yellow", "blue", "red", "green", "black"}
	appList := strings.Split(appName, ",")

	for index, item := range infos {
		var itemCount float64
		labelInfo := make(map[string]string)

		if item == "unknown" {
			var allCount, unknownCount float64
			if mode == "all" {
				models.Orm.Table("api_definition").Count(&allCount)
			} else if mode == "product" || mode == "app" {
				models.Orm.Table("api_definition").Where("app in (?)", appList).Count(&allCount)
			}

			unknownCount = allCount - counts[0] - counts[1]
			counts = append(counts, unknownCount)
			labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(unknownCount))
			labelInfo["color"] = colorNames[index]
			colors = append(colors, defineColors[index])
		} else {
			if mode == "all" {
				models.Orm.Table("api_definition").Where("`check` = ?", item).Count(&itemCount)
			} else if mode == "product" || mode == "app" {
				models.Orm.Table("api_definition").Where("`check` = ? and app in (?)", item, appList).Count(&itemCount)
			}
			counts = append(counts, itemCount)
			labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(itemCount))
			if index < len(defineColors) {
				colors = append(colors, defineColors[index])
				labelInfo["color"] = colorNames[index]
			} else {
				colors = append(colors, defineColors[0])
				labelInfo["color"] = colorNames[0]
			}
		}
		labels = append(labels, labelInfo)
	}

	return
}

func GetAutoAPICount(mode, appName string) (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {
	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}
	colorNames := []string{"yellow", "blue", "red", "green", "black"}
	type ApiSingle struct {
		ApiId string `gorm:"column:api_id" json:"api_id"`
	}
	infos = []string{"yes", "no"}
	appList := strings.Split(appName, ",")

	for index, item := range infos {
		var apiCount, allCount, noNeedAutoCount float64
		labelInfo := make(map[string]string)
		if item == "yes" {
			if mode == "all" {
				models.Orm.Table("api_definition").Group("api_id").Count(&apiCount)
				if apiCount == 0 {
					models.Orm.Table("scene_data").Group("api_id").Count(&apiCount)
				}
			} else if mode == "product" || mode == "app" {
				models.Orm.Table("api_definition").Group("api_id").Where("app in (?) and is_auto = '1'", appList).Count(&apiCount)
				if apiCount == 0 {
					models.Orm.Table("scene_data").Group("api_id").Where("app in (?)", appList).Count(&apiCount)
				}
			}
		} else {
			if mode == "all" {
				models.Orm.Table("api_definition").Count(&allCount)
				models.Orm.Table("api_definition").Where("is_need_auto = 0").Count(&noNeedAutoCount)
			} else if mode == "product" || mode == "app" {
				models.Orm.Table("api_definition").Where("app in (?)", appList).Count(&allCount)
				models.Orm.Table("api_definition").Where("app in (?) and is_need_auto = 0", appList).Count(&noNeedAutoCount)
			}

			if allCount == 0 {
				apiCount = 0
			} else {
				apiCount = allCount - noNeedAutoCount - counts[0]
			}
		}

		counts = append(counts, apiCount)
		labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(apiCount))
		if index < len(defineColors) {
			colors = append(colors, defineColors[index])
			labelInfo["color"] = colorNames[index]
		} else {
			colors = append(colors, defineColors[0])
			labelInfo["color"] = colorNames[0]
		}

		labels = append(labels, labelInfo)

	}

	return
}

func GetAppAPIRunCount() (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {

	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}

	colorNames := []string{"yellow", "blue", "red", "green", "black"}

	models.Orm.Table("scene_data_test_history").Group("app").Pluck("app", &infos)

	if len(infos) == 0 {
		infos = append(infos, T("report.no_data"))
		counts = append(counts, 0)
		colors = append(colors, defineColors[0])
		labelInfo := make(map[string]string)
		labelInfo["label"] = T("report.no_data")
		labelInfo["color"] = colorNames[0]
		labels = append(labels, labelInfo)
		return
	}

	before6MonthTimestamp := time.Now().Unix() - int64(86400*180)
	timeStr := fmt.Sprintf(time.Unix(before6MonthTimestamp, 0).Format("2006-01-02"))

	var allCounts []int
	for _, item := range infos {
		var allCount int
		models.Orm.Table("scene_data_test_history").Where("app = ? and created_at > ?", item, timeStr).Count(&allCount)
		allCounts = append(allCounts, allCount)
	}

	sort.Sort(sorting.Lex(sort.IntSlice(allCounts), sort.StringSlice(infos)))
	Reverse(&infos)

	if len(infos) > 5 {
		infos = infos[:5]
	}

	for index, item := range infos {
		var apiCount float64
		labelInfo := make(map[string]string)
		models.Orm.Table("scene_data_test_history").Where("app = ? and created_at > ?", item, timeStr).Count(&apiCount)
		counts = append(counts, apiCount)
		labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(apiCount))
		if index < len(defineColors) {
			colors = append(colors, defineColors[index])
			labelInfo["color"] = colorNames[index]
		} else {
			colors = append(colors, defineColors[0])
			labelInfo["color"] = colorNames[0]
		}
		labels = append(labels, labelInfo)
	}

	return
}

func GetSceneRunCount() (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {

	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}

	colorNames := []string{"yellow", "blue", "red", "green", "black"}
	models.Orm.Table("scene_test_history").Group("product").Pluck("product", &infos)

	if len(infos) == 0 {
		infos = append(infos, T("report.no_data"))
		counts = append(counts, 0)
		colors = append(colors, defineColors[0])
		labelInfo := make(map[string]string)
		labelInfo["label"] = T("report.no_data")
		labelInfo["color"] = colorNames[0]
		labels = append(labels, labelInfo)
		return
	}

	before6MonthTimestamp := time.Now().Unix() - int64(86400*180)
	timeStr := fmt.Sprintf(time.Unix(before6MonthTimestamp, 0).Format("2006-01-02"))
	var allCounts []int
	for _, item := range infos {
		var allCount int
		models.Orm.Table("scene_test_history").Where("product = ? and created_at > ?", item, timeStr).Count(&allCount)
		allCounts = append(allCounts, allCount)
	}

	sort.Sort(sorting.Lex(sort.IntSlice(allCounts), sort.StringSlice(infos)))
	Reverse(&infos)
	if len(infos) > 5 {
		infos = infos[:5]
	}

	for index, item := range infos {
		var apiCount float64
		labelInfo := make(map[string]string)

		models.Orm.Table("scene_test_history").Where("product = ? and created_at > ?", item, timeStr).Count(&apiCount)

		counts = append(counts, apiCount)
		labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(apiCount))
		if index < len(defineColors) {
			colors = append(colors, defineColors[index])
			labelInfo["color"] = colorNames[index]
		} else {
			colors = append(colors, defineColors[0])
			labelInfo["color"] = colorNames[0]
		}

		labels = append(labels, labelInfo)

	}

	return
}

func GetSceneResultCount() (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {

	infos = []string{"pass", "fail", "unknown"}

	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}

	colorNames := []string{"yellow", "blue", "red", "green", "black"}

	for index, item := range infos {
		var itemCount float64
		labelInfo := make(map[string]string)

		if item == "unknown" {
			var allCount, unknownCount float64
			models.Orm.Table("scene_test_history").Count(&allCount)
			unknownCount = allCount - counts[0] - counts[1]
			if unknownCount == 0 {
				continue
			}
			infos = append(infos, item)
			counts = append(counts, unknownCount)
			labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(unknownCount))
			labelInfo["color"] = colorNames[index]
			colors = append(colors, defineColors[index])
		} else {
			models.Orm.Table("scene_test_history").Where("result = ?", item).Count(&itemCount)
			counts = append(counts, itemCount)
			labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(itemCount))
			if index < len(defineColors) {
				colors = append(colors, defineColors[index])
				labelInfo["color"] = colorNames[index]
			} else {
				colors = append(colors, defineColors[0])
				labelInfo["color"] = colorNames[0]
			}
		}

		labels = append(labels, labelInfo)
	}

	return
}

func GetSceneDataResultCount() (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {

	infos = []string{"pass", "fail", "unknown"}

	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}

	colorNames := []string{"yellow", "blue", "red", "green", "black"}

	for index, item := range infos {
		var itemCount float64
		labelInfo := make(map[string]string)

		if item == "unknown" {
			var allCount, unknownCount float64
			models.Orm.Table("scene_data_test_history").Count(&allCount)
			unknownCount = allCount - counts[0] - counts[1]
			counts = append(counts, unknownCount)
			labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(unknownCount))
			labelInfo["color"] = colorNames[index]
			colors = append(colors, defineColors[index])
		} else {
			models.Orm.Table("scene_data_test_history").Where("result = ?", item).Count(&itemCount)
			counts = append(counts, itemCount)
			labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(itemCount))
			if index < len(defineColors) {
				colors = append(colors, defineColors[index])
				labelInfo["color"] = colorNames[index]
			} else {
				colors = append(colors, defineColors[0])
				labelInfo["color"] = colorNames[0]
			}
		}

		labels = append(labels, labelInfo)
	}

	return
}

func GetScheduleTypeCount() (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {

	infos = []string{T("schedule.task_mode_cron"), T("schedule.task_mode_once"), T("schedule.task_mode_day"), T("schedule.task_mode_week")}
	infosEn := []string{"cron", "once", "day", "week"}

	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}

	colorNames := []string{"yellow", "blue", "red", "green", "black"}

	for index, item := range infosEn {
		var itemCount float64
		labelInfo := make(map[string]string)

		models.Orm.Table("schedule").Where("task_mode = ?", item).Count(&itemCount)
		counts = append(counts, itemCount)
		labelInfo["label"] = fmt.Sprintf(" %s - %d", infos[index], int(itemCount))
		if index < len(defineColors) {
			colors = append(colors, defineColors[index])
			labelInfo["color"] = colorNames[index]
		} else {
			colors = append(colors, defineColors[0])
			labelInfo["color"] = colorNames[0]
		}

		labels = append(labels, labelInfo)
	}

	return
}

func GetAppSceneDataRunCount() (title template.HTML, getMonthLabels, infos []string, counts [][]float64) {
	now := time.Now()
	year := now.Year()
	before6MonthTimestamp := time.Now().Unix() - int64(86400*180)
	before6MonthStr := fmt.Sprintf(time.Unix(before6MonthTimestamp, 0).Format("2006") + T("report.year") + "01" + T("report.month"))
	curMonth := time.Now().Unix()
	curYearAndMonthStr := fmt.Sprintf(time.Unix(curMonth, 0).Format("2006") + T("report.year") + "01" + T("report.month") + "02" + T("report.day"))
	titleTemp := fmt.Sprintf("%s01"+T("report.day")+" - %s", before6MonthStr, curYearAndMonthStr)
	title = template.HTML(titleTemp)

	allMonthLabel := []string{
		T("report.month_1"), T("report.month_2"), T("report.month_3"), T("report.month_4"),
		T("report.month_5"), T("report.month_6"), T("report.month_7"), T("report.month_8"),
		T("report.month_9"), T("report.month_10"), T("report.month_11"), T("report.month_12"),
		T("report.month_1"), T("report.month_2"), T("report.month_3"), T("report.month_4"), T("report.month_5"),
	}
	allMonthInt := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "1", "2", "3", "4", "5"}
	var allYearAndMonths []string
	for _, item := range allMonthInt {
		var yearAndMonth string
		yearAndMonth = fmt.Sprintf("%d-%s-01", year, item)
		allYearAndMonths = append(allYearAndMonths, yearAndMonth)
	}

	before6MonthAndYear := fmt.Sprintf(time.Unix(before6MonthTimestamp, 0).Format("2006-01"))
	beginMonthStr := strings.Split(before6MonthAndYear, "-")[1]
	beginMonth, _ := strconv.Atoi(beginMonthStr)

	monthLabels := allMonthInt[beginMonth-1 : beginMonth+5]

	var curYearAndMonths []string
	for _, item := range monthLabels {
		itemInt, _ := strconv.Atoi(item)
		curYearAndMonths = append(curYearAndMonths, allYearAndMonths[itemInt-1])
		getMonthLabels = append(getMonthLabels, allMonthLabel[itemInt-1])
	}

	models.Orm.Table("scene_data_test_history").Group("app").Pluck("app", &infos)
	if len(infos) == 0 {
		return
	}

	var allCounts []int
	for _, item := range infos {
		var allCount int
		models.Orm.Table("scene_data_test_history").Where("app = ? and created_at > ?", item, curYearAndMonths[0]).Count(&allCount)
		allCounts = append(allCounts, allCount)
	}
	sort.Sort(sorting.Lex(sort.IntSlice(allCounts), sort.StringSlice(infos)))
	Reverse(&infos)

	if len(infos) > 5 {
		infos = infos[:5]
	}

	for _, item := range infos {
		var monthsCount []float64
		var monthCount float64
		for subIndex, subItem := range curYearAndMonths {
			if subIndex == len(curYearAndMonths)-1 {
				models.Orm.Table("scene_data_test_history").Where("app = ? and created_at > ?", item, subItem).Count(&monthCount)
			} else {
				models.Orm.Table("scene_data_test_history").Where("app = ? and created_at > ? and created_at < ?", item, subItem, curYearAndMonths[subIndex+1]).Count(&monthCount)

			}
			monthsCount = append(monthsCount, monthCount)
		}

		counts = append(counts, monthsCount)

	}
	return
}

func GetProductSceneRunCount() (title template.HTML, getMonthLabels, infos []string, counts [][]float64) {
	now := time.Now()
	year := now.Year()
	before6MonthTimestamp := time.Now().Unix() - int64(86400*180)
	before6MonthStr := fmt.Sprintf(time.Unix(before6MonthTimestamp, 0).Format("2006") + T("report.year") + "01" + T("report.month"))
	curMonth := time.Now().Unix()
	curYearAndMonthStr := fmt.Sprintf(time.Unix(curMonth, 0).Format("2006") + T("report.year") + "01" + T("report.month") + "02" + T("report.day"))
	titleTemp := fmt.Sprintf("%s01"+T("report.day")+" - %s", before6MonthStr, curYearAndMonthStr)
	title = template.HTML(titleTemp)

	allMonthLabel := []string{
		T("report.month_1"), T("report.month_2"), T("report.month_3"), T("report.month_4"),
		T("report.month_5"), T("report.month_6"), T("report.month_7"), T("report.month_8"),
		T("report.month_9"), T("report.month_10"), T("report.month_11"), T("report.month_12"),
		T("report.month_1"), T("report.month_2"), T("report.month_3"), T("report.month_4"), T("report.month_5"),
	}
	allMonthInt := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "1", "2", "3", "4", "5"}
	var allYearAndMonths []string
	for _, item := range allMonthInt {
		var yearAndMonth string
		yearAndMonth = fmt.Sprintf("%d-%s-01", year, item)
		allYearAndMonths = append(allYearAndMonths, yearAndMonth)
	}

	before6MonthAndYear := fmt.Sprintf(time.Unix(before6MonthTimestamp, 0).Format("2006-01"))
	beginMonthStr := strings.Split(before6MonthAndYear, "-")[1]
	beginMonth, _ := strconv.Atoi(beginMonthStr)

	monthLabels := allMonthInt[beginMonth-1 : beginMonth+5]

	var curYearAndMonths []string
	for _, item := range monthLabels {
		itemInt, _ := strconv.Atoi(item)
		curYearAndMonths = append(curYearAndMonths, allYearAndMonths[itemInt-1])
		getMonthLabels = append(getMonthLabels, allMonthLabel[itemInt-1])
	}
	models.Orm.Table("scene_test_history").Group("product").Where("product <> ?", "").Pluck("product", &infos)
	if len(infos) == 0 {
		return
	}
	var allCounts []int
	for _, item := range infos {
		var allCount int
		models.Orm.Table("scene_test_history").Where("product = ? and created_at > ?", item, curYearAndMonths[0]).Count(&allCount)
		allCounts = append(allCounts, allCount)
	}

	sort.Sort(sorting.Lex(sort.IntSlice(allCounts), sort.StringSlice(infos)))
	Reverse(&infos)
	if len(infos) > 5 {
		infos = infos[:5]
	}

	for _, item := range infos {
		var monthsCount []float64
		var monthCount float64
		for subIndex, subItem := range curYearAndMonths {
			if subIndex == len(curYearAndMonths)-1 {
				models.Orm.Table("scene_test_history").Where("product = ? and created_at > ?", item, subItem).Count(&monthCount)
			} else {
				models.Orm.Table("scene_test_history").Where("product = ? and created_at > ? and created_at < ?", item, subItem, curYearAndMonths[subIndex+1]).Count(&monthCount)
			}
			monthsCount = append(monthsCount, monthCount)
		}
		counts = append(counts, monthsCount)

	}
	return
}

func GetAppsTableCount() (contents []map[string]types.InfoItem, headers types.Thead) {
	headers = types.Thead{
		{Head: T("report.app_name")},
		{Head: T("report.app_orig_api_count"), Sortable: true},
		{Head: T("report.app_covered_api_count"), Sortable: true},
		{Head: T("report.app_data_file_count"), Sortable: true},
		{Head: T("report.app_history_total"), Sortable: true},
		{Head: T("report.app_history_pass"), Sortable: true},
		{Head: T("report.app_history_fail"), Sortable: true},
		{Head: T("report.app_history_unknown"), Sortable: true},
	}

	var infos []string
	models.Orm.Table("env_config").Order("created_at DESC").Pluck("app", &infos)

	if len(infos) == 0 {
		return
	}

	for _, item := range infos {
		content := make(map[string]types.InfoItem)

		itemHtml := template.HTML(item)
		content[T("report.app_name")] = types.InfoItem{Content: itemHtml}

		var itemCount, allCount, allRunCount, passRunCount, failRunCount, unknownRunCount int
		var apiIds []string

		models.Orm.Table("api_definition").Where("app = ?", item).Count(&itemCount)
		itemCountHtml := template.HTML(fmt.Sprintf("%d", itemCount))
		content[T("report.app_orig_api_count")] = types.InfoItem{Content: itemCountHtml}

		models.Orm.Table("scene_data").Group("api_id").Where("app = ?", item).Pluck("api_id", &apiIds).Count(&itemCount)
		itemCountHtml = template.HTML(fmt.Sprintf("%d", itemCount))
		content[T("report.app_covered_api_count")] = types.InfoItem{Content: itemCountHtml}

		models.Orm.Table("scene_data").Where("app = ?", item).Count(&allCount)
		itemCountHtml = template.HTML(fmt.Sprintf("%d", allCount))
		content[T("report.app_data_file_count")] = types.InfoItem{Content: itemCountHtml}

		models.Orm.Table("scene_data_test_history").Where("app = ?", item).Count(&allRunCount)
		itemCountHtml = template.HTML(fmt.Sprintf("%d", allRunCount))
		content[T("report.app_history_total")] = types.InfoItem{Content: itemCountHtml}

		models.Orm.Table("scene_data_test_history").Where("app = ? and result = ?", item, "pass").Count(&passRunCount)
		itemCountHtml = template.HTML(fmt.Sprintf("%d", passRunCount))
		content[T("report.app_history_pass")] = types.InfoItem{Content: itemCountHtml}

		models.Orm.Table("scene_data_test_history").Where("app = ? and result = ?", item, "fail").Count(&failRunCount)
		itemCountHtml = template.HTML(fmt.Sprintf("%d", failRunCount))
		content[T("report.app_history_fail")] = types.InfoItem{Content: itemCountHtml}

		unknownRunCount = allRunCount - passRunCount - failRunCount
		itemCountHtml = template.HTML(fmt.Sprintf("%d", unknownRunCount))
		content[T("report.app_history_unknown")] = types.InfoItem{Content: itemCountHtml}
		contents = append(contents, content)

	}

	return
}

func GetProductsTableCount() (contents []map[string]types.InfoItem, headers types.Thead) {
	headers = types.Thead{
		{Head: T("report.product_name")},
		{Head: T("report.product_app_count"), Sortable: true},
		{Head: T("report.product_scene_total"), Sortable: true},
		{Head: T("report.product_scene_pass"), Sortable: true},
		{Head: T("report.product_scene_fail"), Sortable: true},
		{Head: T("report.product_scene_history_total"), Sortable: true},
		{Head: T("report.product_scene_history_pass"), Sortable: true},
		{Head: T("report.product_scene_history_fail"), Sortable: true},
		{Head: T("report.product_data_total"), Sortable: true},
		{Head: T("report.product_data_pass"), Sortable: true},
		{Head: T("report.product_data_fail"), Sortable: true},
		{Head: T("report.product_data_history_total"), Sortable: true},
		{Head: T("report.product_data_history_pass"), Sortable: true},
		{Head: T("report.product_data_history_fail"), Sortable: true},
	}

	var infos []string
	curTimestamp := time.Now().Unix() - int64(86400*365*1)
	yearBefore := fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
	models.Orm.Table("product").Where("updated_at >= ? OR (updated_at is null AND created_at >= ?)", yearBefore, yearBefore).Order("created_at desc").Pluck("product", &infos)
	if len(infos) == 0 {
		return
	}

	for _, item := range infos {
		content, err := GetSumOfProduct(item)
		if err != nil {
			continue
		}
		contents = append(contents, content)

	}

	return
}

func GetAPIRunResultCount(mode, name string) (title template.HTML, dayList, infos []string, counts [][]float64) {
	now := time.Now()
	before30DaysTimestamp := time.Now().Unix() - int64(86400*31)
	before6MonthStr := fmt.Sprintf(time.Unix(before30DaysTimestamp, 0).Format("2006") + T("report.year") + "01" + T("report.month"))
	curMonth := time.Now().Unix()
	curYearAndMonthStr := fmt.Sprintf(time.Unix(curMonth, 0).Format("2006") + T("report.year") + "01" + T("report.month") + "02" + T("report.day"))
	titleTemp := fmt.Sprintf("%s01"+T("report.day")+" - %s", before6MonthStr, curYearAndMonthStr)
	title = template.HTML(titleTemp)

	endDay := now.Format("2006-01-02")
	startDay := fmt.Sprintf(time.Unix(before30DaysTimestamp, 0).Format("2006-01-02"))
	dayList = GetBetweenDates(startDay, endDay)
	infos = []string{"pass", "fail", "unknown"}

	var appList []string
	if mode != "product" {
		appList = strings.Split(name, ",")
	}

	for _, item := range infos {
		var daysCount []float64
		var dayCount float64
		for subIndex, subItem := range dayList[1 : len(dayList)-1] {

			if mode == "product" {
				models.Orm.Table("scene_data_test_history").Where("result = ? and created_at > ? and created_at < ? and product = ?", item, subItem, dayList[subIndex+2], name).Count(&dayCount)
			} else {
				models.Orm.Table("scene_data_test_history").Where("result = ? and created_at > ? and created_at < ? and app in (?)", item, subItem, dayList[subIndex+2], appList).Count(&dayCount)
			}
			daysCount = append(daysCount, dayCount)
		}
		counts = append(counts, daysCount)
	}

	return
}

func GetProductPlaybookRunResultCount(product string) (title template.HTML, dayList, infos []string, counts [][]float64) {
	now := time.Now()
	before30DaysTimestamp := time.Now().Unix() - int64(86400*31)
	before6MonthStr := fmt.Sprintf(time.Unix(before30DaysTimestamp, 0).Format("2006") + T("report.year") + "01" + T("report.month"))
	curMonth := time.Now().Unix()
	curYearAndMonthStr := fmt.Sprintf(time.Unix(curMonth, 0).Format("2006") + T("report.year") + "01" + T("report.month") + "02" + T("report.day"))
	titleTemp := fmt.Sprintf("%s01"+T("report.day")+" - %s", before6MonthStr, curYearAndMonthStr)
	title = template.HTML(titleTemp)

	endDay := now.Format("2006-01-02")
	startDay := fmt.Sprintf(time.Unix(before30DaysTimestamp, 0).Format("2006-01-02"))
	dayList = GetBetweenDates(startDay, endDay)
	infos = []string{"pass", "fail", "unknown"}

	for _, item := range infos {
		var daysCount []float64
		var dayCount float64
		for subIndex, subItem := range dayList[1 : len(dayList)-1] {
			models.Orm.Table("scene_test_history").Where("product = ? and result = ? and created_at > ? and created_at < ?", product, item, subItem, dayList[subIndex+2]).Count(&dayCount)
			daysCount = append(daysCount, dayCount)
		}
		counts = append(counts, daysCount)
	}
	return
}

func GetDaysAPIResultCount(mode, name string, day int) (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {

	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}

	colorNames := []string{"yellow", "blue", "red", "green", "black"}
	infos = []string{"pass", "fail", "unknown"}

	beforeDayTimestamp := time.Now().Unix() - int64(86400*day)
	timeStr := fmt.Sprintf(time.Unix(beforeDayTimestamp, 0).Format("2006-01-02 15:04:05"))
	var appList []string
	if mode != "product" {
		appList = strings.Split(name, ",")
	}

	for index, item := range infos {
		var apiCount float64
		labelInfo := make(map[string]string)

		if mode == "product" {
			models.Orm.Table("scene_data_test_history").Where("result = ? and created_at > ? and product = ? ", item, timeStr, name).Count(&apiCount)
		} else {
			models.Orm.Table("scene_data_test_history").Where("result = ? and created_at > ? and app in (?)", item, timeStr, appList).Count(&apiCount)
		}

		counts = append(counts, apiCount)
		labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(apiCount))
		if index < len(defineColors) {
			colors = append(colors, defineColors[index])
			labelInfo["color"] = colorNames[index]
		} else {
			colors = append(colors, defineColors[0])
			labelInfo["color"] = colorNames[0]
		}

		labels = append(labels, labelInfo)

	}

	return
}

func GetDaysSceneResultCount(product string, day int) (infos []string, counts []float64, colors []chartjs.Color, labels []map[string]string) {

	defineColors := []chartjs.Color{"rgb(255, 205, 86)", "rgb(54, 162, 235)", "rgb(238,232,170)", "rgb(189,183,107)", "rgb(255,228,181)"}

	colorNames := []string{"yellow", "blue", "red", "green", "black"}

	infos = []string{"pass", "fail", "unknown"}

	beforeDayTimestamp := time.Now().Unix() - int64(86400*day)
	timeStr := fmt.Sprintf(time.Unix(beforeDayTimestamp, 0).Format("2006-01-02 15:04:05"))

	for index, item := range infos {
		var apiCount float64
		labelInfo := make(map[string]string)
		models.Orm.Table("scene_test_history").Where("result = ? and product = ? and created_at > ?", item, product, timeStr).Count(&apiCount).Limit(5)
		counts = append(counts, apiCount)
		labelInfo["label"] = fmt.Sprintf(" %s - %d", item, int(apiCount))
		if index < len(defineColors) {
			colors = append(colors, defineColors[index])
			labelInfo["color"] = colorNames[index]
		} else {
			colors = append(colors, defineColors[0])
			labelInfo["color"] = colorNames[0]
		}

		labels = append(labels, labelInfo)

	}

	return
}

func GetAppModuleTableCount(appName string) (contents []map[string]types.InfoItem, headers types.Thead) {
	headers = types.Thead{
		{Head: T("report.module_name")},
		{Head: T("report.module_api_total"), Sortable: true},
		{Head: T("report.module_auto_count"), Sortable: true},
		{Head: T("report.module_not_auto_count"), Sortable: true},
		{Head: T("report.module_data_file_count"), Sortable: true},
	}

	appList := strings.Split(appName, ",")
	var infos, httpMethods []string
	models.Orm.Table("api_definition").Group("api_module").Where("app in (?)", appList).Pluck("api_module", &infos)

	models.Orm.Table("api_definition").Group("http_method").Where("app in (?)", appList).Pluck("http_method", &httpMethods)

	if len(infos) == 0 {
		return
	}

	for _, item := range httpMethods {
		httpHtml := types.TheadItem{Head: item, Sortable: true}
		headers = append(headers, httpHtml)
	}
	type APIIdStruct struct {
		ApiId string `gorm:"column:api_id" json:"api_id"`
	}
	var apiIds []APIIdStruct

	for _, item := range infos {
		content := make(map[string]types.InfoItem)

		itemHtml := template.HTML(item)
		content[T("report.module_name")] = types.InfoItem{Content: itemHtml}

		var allCount, autoCount, notAutoCount, dataFileCount int
		var apiList []string
		models.Orm.Table("api_definition").Where("api_module = ? and app in (?)", item, appList).Count(&allCount).Select("api_id").Find(&apiIds)
		allCountHtml := template.HTML(fmt.Sprintf("%d", allCount))
		content[T("report.module_api_total")] = types.InfoItem{Content: allCountHtml}

		for _, subItem := range apiIds {
			apiList = append(apiList, subItem.ApiId)
		}
		models.Orm.Table("scene_data").Group("api_id").Where("api_id in (?) and app in (?)", apiList, appList).Count(&autoCount)
		autoCountHtml := template.HTML(fmt.Sprintf("%d", autoCount))
		content[T("report.module_auto_count")] = types.InfoItem{Content: autoCountHtml}

		if allCount > 0 && allCount > autoCount {
			notAutoCount = allCount - autoCount
		} else {
			notAutoCount = 0
		}
		notAutoCountHtml := template.HTML(fmt.Sprintf("%d", notAutoCount))
		content[T("report.module_not_auto_count")] = types.InfoItem{Content: notAutoCountHtml}

		if len(appList) > 0 {
			models.Orm.Table("scene_data").Where("api_id in (?) and app in (?)", apiList, appList).Select("api_id").Count(&dataFileCount)
		} else {
			dataFileCount = 0
		}
		dataFileCountHtml := template.HTML(fmt.Sprintf("%d", dataFileCount))
		content[T("report.module_data_file_count")] = types.InfoItem{Content: dataFileCountHtml}

		for _, subItem := range httpMethods {
			var apiIds []string
			var methodCount int
			models.Orm.Table("api_definition").Group("api_id").Where("api_module = ? and http_method = ? and app in (?)", item, subItem, appList).Pluck("api_id", &apiIds).Count(&methodCount)
			methodCountHtml := template.HTML(fmt.Sprintf("%d", methodCount))
			content[subItem] = types.InfoItem{Content: methodCountHtml}
		}

		contents = append(contents, content)

	}

	return
}

// ModuleRow 模块统计行数据
type ModuleRow struct {
	ModuleName      string         `json:"module_name"`
	ApiTotal        int            `json:"api_total"`
	AutoCount       int            `json:"auto_count"`
	NotAutoCount    int            `json:"not_auto_count"`
	DataFileCount   int            `json:"data_file_count"`
	MethodBreakdown map[string]int `json:"method_breakdown"`
}

// GetAppModuleTableData 获取模块统计纯数据（不含 HTML 包装）
func GetAppModuleTableData(appName string) (headers []string, rows []ModuleRow) {
	appList := strings.Split(appName, ",")
	var infos, httpMethods []string
	models.Orm.Table("api_definition").Group("api_module").Where("app in (?)", appList).Pluck("api_module", &infos)
	models.Orm.Table("api_definition").Group("http_method").Where("app in (?)", appList).Pluck("http_method", &httpMethods)

	headers = append(headers, T("report.module_name"), T("report.module_api_total"),
		T("report.module_auto_count"), T("report.module_not_auto_count"),
		T("report.module_data_file_count"))
	for _, m := range httpMethods {
		headers = append(headers, m)
	}

	type APIIdStruct struct {
		ApiId string `gorm:"column:api_id" json:"api_id"`
	}
	var apiIds []APIIdStruct

	for _, item := range infos {
		var row ModuleRow
		row.ModuleName = item
		row.MethodBreakdown = make(map[string]int)

		var allCount, autoCount, dataFileCount int
		var apiList []string
		models.Orm.Table("api_definition").Where("api_module = ? and app in (?)", item, appList).Count(&allCount).Select("api_id").Find(&apiIds)
		row.ApiTotal = allCount

		for _, subItem := range apiIds {
			apiList = append(apiList, subItem.ApiId)
		}
		models.Orm.Table("scene_data").Group("api_id").Where("api_id in (?) and app in (?)", apiList, appList).Count(&autoCount)
		row.AutoCount = autoCount

		if allCount > 0 && allCount > autoCount {
			row.NotAutoCount = allCount - autoCount
		}

		if len(appList) > 0 {
			models.Orm.Table("scene_data").Where("api_id in (?) and app in (?)", apiList, appList).Select("api_id").Count(&dataFileCount)
		}
		row.DataFileCount = dataFileCount

		for _, subItem := range httpMethods {
			var methodCount int
			models.Orm.Table("api_definition").Group("api_id").Where("api_module = ? and http_method = ? and app in (?)", item, subItem, appList).Pluck("api_id", &apiIds).Count(&methodCount)
			row.MethodBreakdown[subItem] = methodCount
		}

		rows = append(rows, row)
	}
	return
}

func GetProductAppTableCount(appName string) (contents []map[string]types.InfoItem, headers types.Thead) {
	headers = types.Thead{
		{Head: T("report.product_app_name")},
		{Head: T("report.product_app_api_total")},
		{Head: T("report.product_app_auto_count")},
		{Head: T("report.product_app_not_auto_count")},
		{Head: T("report.product_app_module_count")},
		{Head: T("report.product_app_data_file_count")},
		{Head: T("report.product_app_history_total"), Sortable: true},
		{Head: T("report.product_app_history_pass"), Sortable: true},
		{Head: T("report.product_app_history_fail"), Sortable: true},
		{Head: T("report.product_app_history_unknown"), Sortable: true},
	}

	appList := strings.Split(appName, ",")

	type APIIdStruct struct {
		ApiId string `gorm:"column:api_id" json:"api_id"`
	}
	var apiIds []APIIdStruct

	for _, item := range appList {
		content := make(map[string]types.InfoItem)
		itemHtml := template.HTML(item)
		content[T("report.product_app_name")] = types.InfoItem{Content: itemHtml}

		var allCount, autoCount, notAutoCount, dataFileCount, moduleCount, allRunCount, passRunCount, failRunCount, unknownRunCount int
		var apiList []string

		models.Orm.Table("api_definition").Where("app = ?", item).Count(&allCount).Select("api_id").Find(&apiIds)
		allCountHtml := template.HTML(fmt.Sprintf("%d", allCount))
		content[T("report.product_app_api_total")] = types.InfoItem{Content: allCountHtml}

		for _, subItem := range apiIds {
			apiList = append(apiList, subItem.ApiId)
		}
		models.Orm.Table("scene_data").Group("api_id").Where("app = ? and api_id in (?)", item, apiList).Count(&autoCount)
		autoCountHtml := template.HTML(fmt.Sprintf("%d", autoCount))
		content[T("report.product_app_auto_count")] = types.InfoItem{Content: autoCountHtml}

		if allCount > 0 && allCount > autoCount {
			notAutoCount = allCount - autoCount
		} else {
			notAutoCount = 0
		}
		notAutoCountHtml := template.HTML(fmt.Sprintf("%d", notAutoCount))
		content[T("report.product_app_not_auto_count")] = types.InfoItem{Content: notAutoCountHtml}

		var apiModules []string
		models.Orm.Table("api_definition").Group("api_module").Where("app = ?", item).Pluck("api_module", &apiModules).Count(&moduleCount)
		moduleCountHtml := template.HTML(fmt.Sprintf("%d", moduleCount))
		content[T("report.product_app_module_count")] = types.InfoItem{Content: moduleCountHtml}

		models.Orm.Table("scene_data").Where("app = ?", item).Select("api_id").Count(&dataFileCount)

		dataFileCountHtml := template.HTML(fmt.Sprintf("%d", dataFileCount))
		content[T("report.product_app_data_file_count")] = types.InfoItem{Content: dataFileCountHtml}

		models.Orm.Table("scene_data_test_history").Where("app = ?", item).Count(&allRunCount)
		itemCountHtml := template.HTML(fmt.Sprintf("%d", allRunCount))
		content[T("report.product_app_history_total")] = types.InfoItem{Content: itemCountHtml}

		models.Orm.Table("scene_data_test_history").Where("app = ? and result = ?", item, "pass").Count(&passRunCount)
		itemCountHtml = template.HTML(fmt.Sprintf("%d", passRunCount))
		content[T("report.product_app_history_pass")] = types.InfoItem{Content: itemCountHtml}

		models.Orm.Table("scene_data_test_history").Where("app = ? and result = ?", item, "fail").Count(&failRunCount)
		itemCountHtml = template.HTML(fmt.Sprintf("%d", failRunCount))
		content[T("report.product_app_history_fail")] = types.InfoItem{Content: itemCountHtml}

		unknownRunCount = allRunCount - passRunCount - failRunCount
		itemCountHtml = template.HTML(fmt.Sprintf("%d", unknownRunCount))
		content[T("report.product_app_history_unknown")] = types.InfoItem{Content: itemCountHtml}

		contents = append(contents, content)
	}
	return
}

// ProductAppRow 产品下APP统计行数据
type ProductAppRow struct {
	AppName        string `json:"app_name"`
	ApiTotal       int    `json:"api_total"`
	AutoCount      int    `json:"auto_count"`
	NotAutoCount   int    `json:"not_auto_count"`
	ModuleCount    int    `json:"module_count"`
	DataFileCount  int    `json:"data_file_count"`
	HistoryTotal   int    `json:"history_total"`
	HistoryPass    int    `json:"history_pass"`
	HistoryFail    int    `json:"history_fail"`
	HistoryUnknown int    `json:"history_unknown"`
}

// GetProductAppTableData 获取产品下APP统计纯数据
func GetProductAppTableData(appName string) (headers []string, rows []ProductAppRow) {
	headers = []string{
		T("report.product_app_name"), T("report.product_app_api_total"),
		T("report.product_app_auto_count"), T("report.product_app_not_auto_count"),
		T("report.product_app_module_count"), T("report.product_app_data_file_count"),
		T("report.product_app_history_total"), T("report.product_app_history_pass"),
		T("report.product_app_history_fail"), T("report.product_app_history_unknown"),
	}

	appList := strings.Split(appName, ",")
	type APIIdStruct struct {
		ApiId string `gorm:"column:api_id" json:"api_id"`
	}
	var apiIds []APIIdStruct

	for _, item := range appList {
		var row ProductAppRow
		row.AppName = item

		var allCount, autoCount, notAutoCount, dataFileCount, moduleCount, allRunCount, passRunCount, failRunCount, unknownRunCount int
		var apiList []string

		models.Orm.Table("api_definition").Where("app = ?", item).Count(&allCount).Select("api_id").Find(&apiIds)
		row.ApiTotal = allCount

		for _, subItem := range apiIds {
			apiList = append(apiList, subItem.ApiId)
		}
		models.Orm.Table("scene_data").Group("api_id").Where("app = ? and api_id in (?)", item, apiList).Count(&autoCount)
		row.AutoCount = autoCount

		if allCount > 0 && allCount > autoCount {
			notAutoCount = allCount - autoCount
		}
		row.NotAutoCount = notAutoCount

		var apiModules []string
		models.Orm.Table("api_definition").Group("api_module").Where("app = ?", item).Pluck("api_module", &apiModules).Count(&moduleCount)
		row.ModuleCount = moduleCount

		models.Orm.Table("scene_data").Where("app = ?", item).Select("api_id").Count(&dataFileCount)
		row.DataFileCount = dataFileCount

		models.Orm.Table("scene_data_test_history").Where("app = ?", item).Count(&allRunCount)
		row.HistoryTotal = allRunCount

		models.Orm.Table("scene_data_test_history").Where("app = ? and result = ?", item, "pass").Count(&passRunCount)
		row.HistoryPass = passRunCount

		models.Orm.Table("scene_data_test_history").Where("app = ? and result = ?", item, "fail").Count(&failRunCount)
		row.HistoryFail = failRunCount

		unknownRunCount = allRunCount - passRunCount - failRunCount
		row.HistoryUnknown = unknownRunCount

		rows = append(rows, row)
	}
	return
}

func GetSumOfProduct(name string) (content map[string]types.InfoItem, err error) {
	itemHtml := template.HTML(name)
	content = make(map[string]types.InfoItem)
	content[T("report.product_name")] = types.InfoItem{Content: itemHtml}
	var allCount, passCount, failCount, appNum int

	var appList []string
	models.Orm.Table("product").Where("product = ? and apps IS NOT NULL", name).Pluck("apps", &appList)
	if len(appList) == 0 {
		Logger.Warning(T("report.product_no_app"), name)
	} else {
		appArray := strings.Split(appList[0], ",")
		appNum = len(appArray)
	}

	models.Orm.Table("scene_data_test_history").Where("product = ?", name).Count(&allCount)
	if allCount == 0 {
		err = fmt.Errorf(T("report.product_no_data"), name)
		Logger.Warning("%s", err)
		return
	}

	itemCountHtml := template.HTML(fmt.Sprintf("%d", appNum))
	content[T("report.product_app_count")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_test_history").Where("product = ?", name).Group("name").Count(&allCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", allCount))
	content[T("report.product_scene_total")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_test_history").Where("product = ? and result = ?", name, "pass").Group("name").Count(&passCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", passCount))
	content[T("report.product_scene_pass")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_test_history").Where("product = ? and result = ?", name, "fail").Group("name").Count(&failCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", failCount))
	content[T("report.product_scene_fail")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_test_history").Where("product = ?", name).Count(&allCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", allCount))
	content[T("report.product_scene_history_total")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_test_history").Where("product = ? and result = ?", name, "pass").Count(&passCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", passCount))
	content[T("report.product_scene_history_pass")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_test_history").Where("product = ? and result = ?", name, "fail").Count(&failCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", failCount))
	content[T("report.product_scene_history_fail")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_data_test_history").Where("product = ?", name).Group("name").Count(&passCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", passCount))
	content[T("report.product_data_total")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_data_test_history").Where("product = ? and result = ?", name, "pass").Group("name").Count(&passCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", passCount))
	content[T("report.product_data_pass")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_data_test_history").Where("product = ? and result = ?", name, "fail").Group("name").Count(&failCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", failCount))
	content[T("report.product_data_fail")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_data_test_history").Where("product = ?", name).Count(&failCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", failCount))
	content[T("report.product_data_history_total")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_data_test_history").Where("product = ? and result = ?", name, "pass").Count(&passCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", passCount))
	content[T("report.product_data_history_pass")] = types.InfoItem{Content: itemCountHtml}

	models.Orm.Table("scene_data_test_history").Where("product = ? and result = ?", name, "fail").Count(&failCount)
	itemCountHtml = template.HTML(fmt.Sprintf("%d", failCount))
	content[T("report.product_data_history_fail")] = types.InfoItem{Content: itemCountHtml}

	return
}
