package pages

import (
	"data4test/biz"
	"data4test/models"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

// AC-15 离线渲染能力：构造一条编排报告，确认可输出 HTML
func TestAC15OrchestrationReportRenders(t *testing.T) {
	db, err := gorm.Open("sqlite3", filepath.Join(t.TempDir(), "page.db"))
	if err != nil {
		t.Fatal(err)
	}
	previous := models.Orm
	models.Orm = db
	t.Cleanup(func() {
		db.Close()
		models.Orm = previous
	})
	if err = db.CreateTable(&biz.DashboardReport{}).Error; err != nil {
		t.Fatal(err)
	}

	data := biz.OrchestrationReportData{
		OrchestrationId:   "1",
		OrchestrationName: "ac15-orchestration",
		TriggerType:       "manual",
		StatsVersion:      biz.OrchestrationReportStatsVersion,
		Summary:           biz.OrchestrationSummary{Total: 1, Success: 1, PassRate: 100},
		SceneStats:        biz.OrchestrationResourceStats{Total: 2, Success: 1, Fail: 1, PassRate: 50},
		DataStats:         biz.OrchestrationResourceStats{Total: 5, Success: 4, Fail: 1, PassRate: 80},
		Products: []biz.OrchestrationProduct{{
			Product: "product-a", TaskTotal: 1, TaskSuccess: 1,
			SceneTotal: 2, SceneSuccess: 1, SceneFail: 1,
			DataTotal: 5, DataSuccess: 4, DataFail: 1, DataPassRate: 90,
		}},
		Steps: []biz.OrchestrationStep{{
			StepOrder: 1, TaskName: "task-a", Status: biz.OrchestrationStepSuccess, DashboardId: "88",
		}},
	}
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	panel, err := renderOrchestrationReport(biz.DashboardReport{ReportData: string(raw), Status: "finished"})
	if err != nil {
		t.Fatal(err)
	}
	content := string(panel.Content)
	for _, expected := range []string{
		"ac15-orchestration", "task-a", "/admin/report_dashboard?id=88",
		biz.T("task_orchestration.scene_total"), biz.T("task_orchestration.data_total"),
		biz.T("task_orchestration.success_count"), biz.T("task_orchestration.fail_count"),
		biz.T("task_orchestration.task_success_count"), biz.T("task_orchestration.task_fail_count"),
		biz.T("task_orchestration.scene_success_count"), biz.T("task_orchestration.scene_fail_count"),
		biz.T("task_orchestration.data_success_count"), biz.T("task_orchestration.data_fail_count"),
		biz.T("task_orchestration.data_pass_rate"),
		"product-a", "90.0%", "50.0%", "80.0%",
	} {
		if !strings.Contains(content, expected) {
			t.Fatalf("rendered report missing %q", expected)
		}
	}
	for _, removed := range []string{"orchResultPie", "orchTypePie", "orchDurationPie"} {
		if strings.Contains(content, removed) {
			t.Fatalf("rendered report should not contain pie chart %q", removed)
		}
	}

	interrupted, err := renderOrchestrationReport(biz.DashboardReport{Status: "failed"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(interrupted.Content), biz.T("task_orchestration.report_interrupted")) {
		t.Fatal("failed report should show interrupted message")
	}
}
