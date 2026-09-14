package pages

import (
	"data4test/biz"
	"data4test/models"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"
)

func GetProductReportData(productName, appName, userName string) (err error) {
	var productDR ProductDashboardReport

	productDR.APITypeCount.Infos, productDR.APITypeCount.Counts, productDR.APITypeCount.Colors, productDR.APITypeCount.Labels = biz.GetAPITypeCount("product", appName)
	productDR.APISpecCount.Infos, productDR.APISpecCount.Counts, productDR.APISpecCount.Colors, productDR.APISpecCount.Labels = biz.GetAPISpecCount("product", appName)
	productDR.AutoAPICount.Infos, productDR.AutoAPICount.Counts, productDR.AutoAPICount.Colors, productDR.AutoAPICount.Labels = biz.GetAutoAPICount("product", appName)
	productDR.APIRunResultCount.Title, productDR.APIRunResultCount.DayList, productDR.APIRunResultCount.Infos, productDR.APIRunResultCount.Counts = biz.GetAPIRunResultCount("product", productName)
	productDR.DaysAPIResultCount.Infos, productDR.DaysAPIResultCount.Counts, productDR.DaysAPIResultCount.Colors, productDR.DaysAPIResultCount.Labels = biz.GetDaysAPIResultCount("product", productName, 30)
	productDR.ProductPlaybookResultCount.Title, productDR.ProductPlaybookResultCount.DayList, productDR.ProductPlaybookResultCount.Infos, productDR.ProductPlaybookResultCount.Counts = biz.GetProductPlaybookRunResultCount(productName)
	productDR.DaysSceneResultCount.Infos, productDR.DaysSceneResultCount.Counts, productDR.DaysSceneResultCount.Colors, productDR.DaysSceneResultCount.Labels = biz.GetDaysSceneResultCount(productName, 30)
	productDR.ProductAppModuleTableCount.Contents, productDR.ProductAppModuleTableCount.Headers = biz.GetProductAppTableCount(appName)

	productDRStr, _ := json.MarshalIndent(productDR, "", "    ")

	nowStr := time.Now().Format("20060102150405")
	now := time.Now().Format("2006-01-02 15:04:05")

	// 从关联数据执行历史中获取实际时间区间作为统计时间
	var startTime, endTime string
	models.Orm.Table("scene_data_test_history").
		Select("MIN(created_at) as start_time, MAX(created_at) as end_time").
		Where("product = ?", productName).
		Row().Scan(&startTime, &endTime)
	if len(startTime) == 0 || len(endTime) == 0 {
		// 如果 scene_data_test_history 没有数据，尝试从 scene_test_history 查询
		models.Orm.Table("scene_test_history").
			Select("MIN(created_at) as start_time, MAX(created_at) as end_time").
			Where("product = ?", productName).
			Row().Scan(&startTime, &endTime)
	}
	if len(startTime) == 0 || len(endTime) == 0 {
		startTime = now
		endTime = now
	}

	reportName := biz.T("product.report_name", productName, nowStr)
	report := biz.DashboardReport{
		ReportName:      reportName,
		ReportType:      "product",
		RelatedProducts: productName,
		RelatedApps:     appName,
		TimeRangeStart:  startTime,
		TimeRangeEnd:    endTime,
		Status:          "finished",
		Creator:         userName,
		ReportData:      string(productDRStr),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = models.Orm.Table("dashboard").Create(&report).Error
	if err != nil {
		return fmt.Errorf(biz.T("product.save_report_failed"), err)
	}
	return err
}

func GetDashBoard2Content(ctx *gin.Context, userName string) (types.Panel, error) {
	id := ctx.Query("id")
	productName, err := biz.GetProductName(id)
	if err != nil {
		return types.Panel{}, err
	}
	appName, err := biz.GetProductApps(id)
	if err != nil {
		return types.Panel{}, err
	}

	// 优先从 dashboard 表读取预计算数据
	var productReportData ProductDashboardReport
	var dr biz.DashboardReport
	err = models.Orm.Table("dashboard").
		Where("related_products =  ?", productName).
		Order("created_at desc").
		Limit(1).
		Find(&dr).Error

	if len(dr.Id) > 0 {
		err := json.Unmarshal([]byte(dr.ReportData), &productReportData)
		if err != nil {
			biz.Logger.Error("err: %v", err)
			biz.Logger.Debug("%v", dr.ReportData)
		}
	} else {
		err = GetProductReportData(productName, appName, userName)
		if err != nil {
			biz.Logger.Error("err: %v", err)
			return types.Panel{}, err
		}
		err = models.Orm.Table("dashboard").
			Where("related_apps =  ?", appName).
			Order("created_at desc").
			Limit(1).
			Find(&dr).Error

		if len(dr.Id) > 0 {
			err = json.Unmarshal([]byte(dr.ReportData), &productReportData)
			if err != nil {
				biz.Logger.Error("err: %v", err)
				biz.Logger.Debug("%+v", productReportData)
			}
		}
	}
	return renderProductReport(productName, dr)

}

func renderProductReport(productName string, report biz.DashboardReport) (types.Panel, error) {
	var productReportData ProductDashboardReport
	err := json.Unmarshal([]byte(report.ReportData), &productReportData)
	if err != nil {
		return types.Panel{
			Content:     template.HTML(fmt.Sprintf("<div style='padding:20px;color:red'>%s</div>", biz.T("schedule_report.parse_error", err))),
			Title:       template.HTML(biz.T("schedule_report.page_title")),
			Description: template.HTML(biz.T("schedule_report.description")),
		}, nil
	}

	// 第1行：API 类型分布 / API 规范 / API 自动化
	pie2Data := BaseCount{
		Infos:  productReportData.APISpecCount.Infos,
		Counts: productReportData.APITypeCount.Counts,
		Colors: productReportData.APITypeCount.Colors,
		Labels: productReportData.APITypeCount.Labels,
	}
	row1 := `<div class="row">` +
		globalPie("pieChart1", biz.T("dashboard.api_type_distribution"), productReportData.APITypeCount, "/admin/info/api_definition", biz.T("dashboard.view_all_apis")) +
		globalPie("pieChart2", biz.T("dashboard.api_spec_check"), pie2Data, "/admin/info/api_definition", biz.T("dashboard.view_all_apis")) +
		globalPie("pieChart3", biz.T("dashboard.api_is_automation"), productReportData.AutoAPICount, "/admin/info/scene_data", biz.T("dashboard.view_all_data_files")) +
		`</div>`

	// 第2行：数据执行趋势（折线）+ 数据执行分布（饼）
	line1 := globalLine("dataChart", biz.T("dashboard.data_exec_trend_30d"), productReportData.APIRunResultCount)
	pie4 := globalPie("pieChart4", biz.T("dashboard.data_exec_dist_30d"), productReportData.DaysAPIResultCount, "/admin/info/scene_data_test_history", biz.T("dashboard.view_all_history_data"))
	row2 := `<div class="row"><div class="col-md-8">` + line1 + `</div>` + pie4 + `</div>`

	// 第3行：场景执行趋势（折线）+ 场景执行分布（饼）
	line2 := globalLine("sceneChart", biz.T("dashboard.scene_exec_trend_30d"), productReportData.ProductPlaybookResultCount)
	pie5 := globalPie("pieChart5", biz.T("dashboard.scene_exec_dist_30d"), productReportData.DaysSceneResultCount, "/admin/info/scene_test_history", biz.T("dashboard.view_all_history_scene"))
	row3 := `<div class="row"><div class="col-md-8">` + line2 + `</div>` + pie5 + `</div>`

	// 应用 API 概览表
	row4 := reportCard(biz.T("dashboard.app_api_overview"), "", reportTableFromInfo(productReportData.ProductAppModuleTableCount.Contents, productReportData.ProductAppModuleTableCount.Headers))

	content := row1 + row2 + row3 + row4 + string(reportStyle())

	return types.Panel{
		Content:     template.HTML(content),
		Title:       template.HTML(biz.T("product_dashboard.title")),
		Description: template.HTML(fmt.Sprintf(`<div style="display:flex;justify-content:space-between"><span>%s - %s</span><span style="color:#888">%s: %s</span></div>`, biz.T("product_dashboard.description"), productName, biz.T("schedule_report.generated_at"), report.CreatedAt)),
	}, nil
}
