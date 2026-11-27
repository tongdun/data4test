package biz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func remainingP1Bundle(t *testing.T, files map[string]string, manifest string) []byte {
	t.Helper()
	var tarBytes bytes.Buffer
	tw := tar.NewWriter(&tarBytes)
	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0644, Size: int64(len(content))}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	var compressed bytes.Buffer
	gw := gzip.NewWriter(&compressed)
	if _, err := gw.Write(tarBytes.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	return p1Zip(t, map[string]string{"orchestration_info.yaml": manifest, "bundle/tasks.tgz": string(compressed.Bytes())})
}
func remainingP1Check(t *testing.T, files map[string]string) map[string]interface{} {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bundle.zip")
	if err := os.WriteFile(path, remainingP1Bundle(t, files, "- name: orch\n  tasks: [task]\n"), 0644); err != nil {
		t.Fatal(err)
	}
	checked, err := ImportTaskOrchestrationCheck(path, "tester")
	if err != nil {
		t.Fatal(err)
	}
	return checked
}
func TestP1BlockedImportCannotBeConfirmed(t *testing.T) {
	for _, kind := range []string{"task", "scene", "nested-data"} {
		t.Run(kind, func(t *testing.T) {
			db := p1Database(t)
			files := map[string]string{"task_info.yaml": "task:\n  任务类型: playbook\n  关联场景: [scene]\n"}
			missing := "scene"
			if kind == "task" {
				files["task_info.yaml"] = "other:\n  任务类型: data\n"
				missing = "task"
			}
			if kind == "nested-data" {
				files["playbook_info.yaml"] = "scene:\n  执行类型: continue\n  关联数据: [absent.yaml]\n"
				missing = "absent.yaml"
			}
			checked := remainingP1Check(t, files)
			if !checked["blocked"].(bool) {
				t.Fatal("missing dependency not blocked")
			}
			id := checked["import_id"].(string)
			state, err := loadOrchestrationImportState(id)
			if err != nil {
				t.Fatal(err)
			}
			if !state.Blocked || !strings.Contains(state.DependencyError, missing) {
				t.Fatalf("blocked state not persisted: %+v", state)
			}
			for _, mode := range []string{"skip", "overwrite"} {
				err = ImportTaskOrchestrationConfirm(id, mode, "tester")
				if err == nil || !strings.Contains(err.Error(), missing) {
					t.Fatalf("%s confirmation did not reject missing %s: %v", mode, missing, err)
				}
				for _, table := range []string{"schedule", "playbook", "scene_data", "task_orchestration"} {
					var count int
					db.Table(table).Count(&count)
					if count != 0 {
						t.Fatalf("blocked import wrote %s", table)
					}
				}
			}
			// Old or stale state files without the cached blocked flag must also be checked.
			state.Blocked = false
			state.DependencyError = ""
			if err = saveOrchestrationImportState(state); err != nil {
				t.Fatal(err)
			}
			if err = ImportTaskOrchestrationConfirm(id, "skip", "tester"); err == nil {
				t.Fatal("stale unblocked state bypassed dependency check")
			}
			if err = ImportTaskOrchestrationConfirm(id, "cancel", "tester"); err != nil {
				t.Fatal(err)
			}
			if _, err = os.Stat(state.TempDir); !os.IsNotExist(err) {
				t.Fatal("cancel did not clean blocked outer staging")
			}
			if _, err = os.Stat(filepath.Join(DownloadBasePath, "import_"+state.ScheduleImportId+".json")); !os.IsNotExist(err) {
				t.Fatal("cancel did not clean inner state")
			}
		})
	}
}
func TestP1ImportKeepsCompleteOrchestrationWhenSiblingMissing(t *testing.T) {
	db := p1Database(t)
	bundle := remainingP1Bundle(t,
		map[string]string{"task_info.yaml": "good-task:\n  任务类型: data\n  关联数据: [one.yaml]\n", "one.yaml": "name: one\n"},
		"- name: good-orch\n  tasks: [good-task]\n- name: broken-orch\n  tasks: [missing-task]\n",
	)
	path := filepath.Join(t.TempDir(), "multi.zip")
	if err := os.WriteFile(path, bundle, 0644); err != nil {
		t.Fatal(err)
	}
	checked, err := ImportTaskOrchestrationCheck(path, "tester")
	if err != nil {
		t.Fatal(err)
	}
	// 编排之间互相独立：一条缺失不应阻断同包其他完整编排。
	if checked["blocked"].(bool) {
		t.Fatal("partial missing dependencies blocked the whole package")
	}
	if checked["missing_count"].(int) != 1 || checked["importable_count"].(int) != 1 {
		t.Fatalf("unexpected counts: %v/%v", checked["missing_count"], checked["importable_count"])
	}
	id := checked["import_id"].(string)
	if err = ImportTaskOrchestrationConfirm(id, "skip", "tester"); err != nil {
		t.Fatal(err)
	}
	var names []string
	if err = db.Table("task_orchestration").Pluck("orchestration_name", &names).Error; err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 || names[0] != "good-orch" {
		t.Fatalf("unexpected orchestrations: %v", names)
	}
	var taskCount int
	db.Table("schedule").Where("task_name = 'good-task'").Count(&taskCount)
	if taskCount != 1 {
		t.Fatal("complete orchestration dependencies were not imported")
	}
	state, err := loadOrchestrationImportState(id)
	if err != nil {
		t.Fatal(err)
	}
	if state.Status != "done" || !strings.Contains(state.Message, "broken-orch") {
		t.Fatalf("skipped orchestration not reported: %+v", state)
	}
}

func TestP1ImportRechecksDependenciesBeforeWriting(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO scene_data (file_name,name) VALUES ('external.yaml','external')")
	p1Exec(t, db, "INSERT INTO playbook (name,data_file_list) VALUES ('scene','external.yaml')")
	checked := remainingP1Check(t, map[string]string{"task_info.yaml": "task:\n  任务类型: playbook\n  关联场景: [scene]\n"})
	if checked["blocked"].(bool) {
		t.Fatal("complete existing dependencies blocked")
	}
	p1Exec(t, db, "UPDATE scene_data SET deleted_at='2026-11-10'")
	if err := ImportTaskOrchestrationConfirm(checked["import_id"].(string), "skip", "tester"); err == nil || !strings.Contains(err.Error(), "external.yaml") {
		t.Fatalf("removed dependency accepted: %v", err)
	}
	var count int
	db.Table("schedule").Count(&count)
	if count != 0 {
		t.Fatal("dependencies were written before revalidation")
	}
}
func TestP1ImportValidatesReusedTaskDependencies(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO schedule (task_name,task_type,scene_list) VALUES ('task','scene','target-scene')")
	p1Exec(t, db, "INSERT INTO playbook (name,data_file_list) VALUES ('target-scene','missing-target.yaml')")
	checked := remainingP1Check(t, map[string]string{"task_info.yaml": "task:\n  任务类型: data\n  关联场景: [source.yaml]\n", "source.yaml": "name: source\n"})
	if checked["blocked"].(bool) {
		t.Fatal("complete source package was blocked")
	}
	if err := ImportTaskOrchestrationConfirm(checked["import_id"].(string), "skip", "tester"); err == nil || !strings.Contains(err.Error(), "missing-target.yaml") {
		t.Fatalf("skip accepted broken reused task: %v", err)
	}
	if err := ImportTaskOrchestrationConfirm(checked["import_id"].(string), "overwrite", "tester"); err != nil {
		t.Fatal(err)
	}
	var task DbSchedule
	db.Table("schedule").Where("task_name='task'").Find(&task)
	if task.TaskType != "data" || task.DataList != "source.yaml" {
		t.Fatalf("valid overwrite failed: %+v", task)
	}
}
func TestP1ImportValidSceneDependencyClosure(t *testing.T) {
	db := p1Database(t)
	checked := remainingP1Check(t, map[string]string{"task_info.yaml": "task:\n  任务类型: playbook\n  关联场景: [scene]\n", "playbook_info.yaml": "scene:\n  执行类型: continue\n  关联数据: [one.yaml]\n", "one.yaml": "name: one\n"})
	if checked["blocked"].(bool) {
		t.Fatal("complete scene dependency closure blocked")
	}
	if err := ImportTaskOrchestrationConfirm(checked["import_id"].(string), "skip", "tester"); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"schedule", "playbook", "scene_data", "task_orchestration"} {
		var count int
		db.Table(table).Count(&count)
		if count != 1 {
			t.Fatalf("%s count=%d", table, count)
		}
	}
}

type remainingP1ReadError struct{}

func (remainingP1ReadError) Read([]byte) (int, error) {
	return 0, errors.New("injected upload read error")
}
func TestP1UploadIsIsolatedAndCleaned(t *testing.T) {
	p1Database(t)
	attachment := filepath.Join(UploadBasePath, "archive.zip")
	if err := os.WriteFile(attachment, []byte("existing attachment"), 0644); err != nil {
		t.Fatal(err)
	}
	valid := remainingP1Bundle(t, map[string]string{"task_info.yaml": "task:\n  任务类型: data\n  关联场景: [one.yaml]\n", "one.yaml": "name: one\n"}, "- name: orch\n  tasks: [task]\n")
	for _, test := range []struct {
		name      string
		reader    io.Reader
		wantError bool
	}{
		{"invalid", strings.NewReader("not an archive"), true},
		{"read-error", remainingP1ReadError{}, true},
		{"valid", bytes.NewReader(valid), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := CheckTaskOrchestrationUpload(test.reader, "tester")
			if (err != nil) != test.wantError {
				t.Fatalf("err=%v", err)
			}
			content, err := os.ReadFile(attachment)
			if err != nil || string(content) != "existing attachment" {
				t.Fatal("upload overwrote existing attachment")
			}
			entries, err := os.ReadDir(UploadBasePath)
			if err != nil || len(entries) != 1 {
				t.Fatal("upload polluted attachment directory")
			}
			rawDirs, _ := filepath.Glob(filepath.Join(DownloadBasePath, "orch_upload_*"))
			if len(rawDirs) != 0 {
				t.Fatalf("raw upload leaked: %v", rawDirs)
			}
		})
	}
}
func TestP1ConcurrentSceneRepeatsRetainErrors(t *testing.T) {
	for _, count := range []int{2, 5} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			var mu sync.Mutex
			active, maxActive, called := 0, 0, 0
			result, _, err := runConcurrentPlaybookRepeats(Playbook{Apis: []string{"one"}}, count, 2, func(item Playbook) (string, string, error) {
				mu.Lock()
				index := called
				called++
				active++
				if active > maxActive {
					maxActive = active
				}
				mu.Unlock()
				defer func() { mu.Lock(); active--; mu.Unlock() }()
				item.Apis[0] = "private copy"
				time.Sleep(time.Millisecond)
				if index == 0 {
					return "fail", "", fmt.Errorf("%s expired", T("error.auth_fail_prefix"))
				}
				if index == 1 {
					panic("injected scene worker panic")
				}
				return "pass", "last", nil
			})
			if result != "fail" || err == nil || !isAuthFailureError(err) || !strings.Contains(err.Error(), "injected scene worker panic") {
				t.Fatalf("result=%s err=%v", result, err)
			}
			if called != count || active != 0 || maxActive > 2 {
				t.Fatalf("called=%d active=%d maximum=%d", called, active, maxActive)
			}
		})
	}
}
func TestP1ConcurrentRepeatedScene401FinishesOrchestration(t *testing.T) {
	db := p1Database(t)
	for _, ddl := range cliDDLStatements {
		p1Exec(t, db, ddl)
	}
	oldHistory := HistoryBasePath
	HistoryBasePath = filepath.Join(t.TempDir(), "history")
	t.Cleanup(func() { HistoryBasePath = oldHistory })
	if err := os.MkdirAll(HistoryBasePath, 0755); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error":"expired"}`)
	}))
	defer server.Close()
	host := strings.TrimPrefix(server.URL, "http://")
	p1Exec(t, db, "INSERT INTO product (product,ip,protocol,testmode,env_type,threading,thread_number,apps) VALUES ('A',?,'http','single',1,'yes',2,'app')", host)
	p1Exec(t, db, "INSERT INTO env_config (app,ip,protocol,threading,thread_number) VALUES ('app',?,'http','no',1)", host)
	content := "name: probe\napi_id: get_probe\nis_use_env_config: yes\nenv:\n  protocol: http\n  host: 127.0.0.1\napi:\n  app: app\n  method: get\n  path: /probe\nsingle:\n  header:\n    Content-Type: application/json\n  query:\n    probe: test\n  body: {}\nmulti: {}\naction: []\nassert: []\noutput: {}\n"
	p1Exec(t, db, "INSERT INTO scene_data (name,file_name,api_id,app,content,file_type,run_time) VALUES ('probe','probe.yaml','get_probe','app',?,1,1)", content)
	if err := os.WriteFile(filepath.Join(DataBasePath, "probe.yaml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	p1Exec(t, db, "INSERT INTO playbook (name,data_file_list,run_time,scene_type,product) VALUES ('scene','probe.yaml',2,1,'A')")
	p1Exec(t, db, "INSERT INTO schedule (task_name,task_type,scene_list,product_list,task_status) VALUES ('first','scene','scene','A','not_started'),('after','data','','A','not_started')")
	p1Exec(t, db, "INSERT INTO task_orchestration (orchestration_name,task_list,orchestration_status) VALUES ('orch',?,'not_started')", "first\nafter")
	done := make(chan error, 1)
	go func() { _, err := RunOrchestration("1", "tester"); done <- err }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("concurrent 401 left orchestration stuck")
	}
	var orch DbTaskOrchestration
	db.Table("task_orchestration").Where("id=1").Find(&orch)
	if orch.OrchestrationStatus != "finished" {
		t.Fatal("orchestration did not finish")
	}
	var row DashboardReport
	db.Table("dashboard").Where("report_type='orchestration'").Find(&row)
	report, err := parseRemainingP1Report(row.ReportData)
	if err != nil {
		t.Fatal(err)
	}
	if report.Summary.Fail != 1 || report.Summary.Skipped != 1 || !strings.Contains(report.Steps[0].ErrorMsg, "401") || report.Steps[0].DashboardId == "" {
		t.Fatalf("incorrect report %+v", report)
	}
}

func parseRemainingP1Report(raw string) (report OrchestrationReportData, err error) {
	err = json.Unmarshal([]byte(raw), &report)
	return
}
