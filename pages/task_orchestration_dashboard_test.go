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
		Summary:           biz.OrchestrationSummary{Total: 1, Success: 1, PassRate: 100},
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
	for _, expected := range []string{"ac15-orchestration", "task-a", "/admin/report_dashboard?id=88"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("rendered report missing %q", expected)
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
