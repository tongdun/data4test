package biz

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"data4test/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v3"
)

// Run from the repository root with build/test-exec.sh: package init reads config.json.
func p1Database(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	previous := models.Orm
	models.Orm = db
	db.LogMode(false)
	oldBase, oldDownload, oldUpload, oldData := BASEPATH, DownloadBasePath, UploadBasePath, DataBasePath
	BASEPATH = t.TempDir()
	DownloadBasePath, UploadBasePath, DataBasePath = filepath.Join(BASEPATH, "download"), filepath.Join(BASEPATH, "upload"), filepath.Join(BASEPATH, "data")
	for _, dir := range []string{DownloadBasePath, UploadBasePath, DataBasePath} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		db.Close()
		models.Orm = previous
		BASEPATH, DownloadBasePath, UploadBasePath, DataBasePath = oldBase, oldDownload, oldUpload, oldData
	})
	tables := map[string]interface{}{"schedule": &DbSchedule{}, "playbook": &DbScene{}, "scene_data": &DbSceneData{}, "product": &DbProduct{}, "env_config": &EnvConfig{}, "sys_parameter": &SysParameter{}, "assert_template": &AssertValueDefine{}, "task_orchestration": &DbTaskOrchestration{}, "dashboard": &DashboardReport{}}
	for table, model := range tables {
		columns := map[string]bool{"id": true, "deleted_at": true, "created_at": true, "updated_at": true}
		definition := []string{"id INTEGER PRIMARY KEY AUTOINCREMENT", "deleted_at TEXT", "created_at TEXT DEFAULT CURRENT_TIMESTAMP", "updated_at TEXT"}
		for _, field := range db.NewScope(model).Fields() {
			if !field.IsNormal || columns[field.DBName] {
				continue
			}
			columns[field.DBName] = true
			definition = append(definition, "\""+field.DBName+"\" TEXT")
		}
		if err := db.Exec("CREATE TABLE " + table + " (" + strings.Join(definition, ",") + ")").Error; err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func p1Exec(t *testing.T, db *gorm.DB, query string, args ...interface{}) {
	t.Helper()
	if err := db.Exec(query, args...).Error; err != nil {
		t.Fatal(err)
	}
}
func p1Import(t *testing.T, result ImportResult, mode string) {
	t.Helper()
	raw, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(DownloadBasePath, "import_"+result.ImportId+".json"), raw, 0644); err != nil {
		t.Fatal(err)
	}
	if err = ImportScheduleConfirm(result.ImportId, mode, "cli"); err != nil {
		t.Fatal(err)
	}
}

func TestP1ImportPreservesTaskDependencies(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO playbook (name) VALUES (?)", "shared")
	products := "A,B"
	result := ImportResult{ImportId: "new", ProductName: "default", Conflicts: map[string]string{"shared": "scene"}, Tasks: []KTask{
		{Name: "scene-task", TaskType: "scene", PlaybookList: []string{"shared"}, ProductList: &products},
		{Name: "data-task", TaskType: "data", PlaybookList: []string{"one.yaml", "two.yaml"}, ProductList: &products},
	}}
	p1Import(t, result, "skip")
	var scene, data DbSchedule
	db.Table("schedule").Where("task_name = ?", "scene-task").Find(&scene)
	db.Table("schedule").Where("task_name = ?", "data-task").Find(&data)
	if scene.SceneList != "shared" || scene.ProductList != products {
		t.Fatalf("lost reused scene/products: %+v", scene.Schedule)
	}
	if data.DataList != "one.yaml,two.yaml" || data.SceneList != "" || data.ProductList != products {
		t.Fatalf("lost data/products: %+v", data.Schedule)
	}
	result.ImportId = "overwrite"
	result.Tasks = []KTask{{Name: "scene-task", TaskType: "data", PlaybookList: []string{"replacement.yaml"}, ProductList: &products}}
	result.Conflicts = map[string]string{"scene-task": "task"}
	p1Import(t, result, "overwrite")
	db.Table("schedule").Where("task_name = ?", "scene-task").Find(&scene)
	if scene.DataList != "replacement.yaml" || scene.SceneList != "" || scene.TaskType != "data" {
		t.Fatalf("incorrect type overwrite: %+v", scene.Schedule)
	}
}

func TestP1SkipPreservesMetadataAndUploads(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO product (product,ip) VALUES ('A','target')")
	p1Exec(t, db, "INSERT INTO env_config (app,ip) VALUES ('api','target')")
	p1Exec(t, db, "INSERT INTO sys_parameter (name,value_list) VALUES ('token','target')")
	p1Exec(t, db, "INSERT INTO assert_template (name,value) VALUES ('check','target')")
	cfg := CliConfig{Product: CliProduct{Name: "A", Ip: "source"}, App: map[string]CliApp{"api": {Ip: "source"}, "new-api": {Ip: "source"}}, SysParameter: map[string]CliJSONField{"token": "source"}, AssertTemplates: map[string]CliJSONField{"check": "source"}}
	for _, mode := range []string{"skip", "overwrite"} {
		if err := upsertCliConfigMode(cfg, cfg.SysParameter, mode); err != nil {
			t.Fatal(err)
		}
		expected := "target"
		if mode == "overwrite" {
			expected = "source"
		}
		for _, query := range []string{"SELECT ip FROM product WHERE product='A'", "SELECT ip FROM env_config WHERE app='api'", "SELECT value_list FROM sys_parameter WHERE name='token'", "SELECT value FROM assert_template WHERE name='check'"} {
			var got string
			if err := db.Raw(query).Row().Scan(&got); err != nil {
				t.Fatal(err)
			}
			if got != expected {
				t.Fatalf("%s: got %s want %s", query, got, expected)
			}
		}
	}
	stage := t.TempDir()
	os.Mkdir(filepath.Join(stage, "upload"), 0755)
	os.WriteFile(filepath.Join(stage, "upload", "attachment.txt"), []byte("source"), 0644)
	os.WriteFile(filepath.Join(UploadBasePath, "attachment.txt"), []byte("target"), 0644)
	if err := copyCliUploadsMode(stage, "skip"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(UploadBasePath, "attachment.txt"))
	if string(got) != "target" {
		t.Fatal("skip overwrote attachment")
	}
	if err := copyCliUploadsMode(stage, "overwrite"); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(filepath.Join(UploadBasePath, "attachment.txt"))
	if string(got) != "source" {
		t.Fatal("overwrite did not replace attachment")
	}
}

func TestP1ExportChecksTransitiveDependencies(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO task_orchestration (orchestration_name,task_list) VALUES ('orch','task')")
	p1Exec(t, db, "INSERT INTO schedule (task_name,task_type,scene_list) VALUES ('task','scene','missing-scene')")
	if _, err := ExportTaskOrchestration("1", "test"); err == nil {
		t.Fatal("export accepted missing scene")
	}
	p1Exec(t, db, "UPDATE schedule SET scene_list='scene'")
	p1Exec(t, db, "INSERT INTO playbook (name,data_file_list) VALUES ('scene','missing-data.yaml')")
	if _, err := ExportTaskOrchestration("1", "test"); err == nil {
		t.Fatal("export accepted scene with missing data")
	}
	p1Exec(t, db, "UPDATE schedule SET task_type='data', data_list='missing-data.yaml', scene_list=''")
	if _, err := ExportTaskOrchestration("1", "test"); err == nil {
		t.Fatal("export accepted missing direct data")
	}
}

func TestP1ProductConfigurationsRoundTrip(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO product (product,ip,apps) VALUES ('A','host-A','api-A'),('B','host-B','api-B'),('override','host-O','api-O')")
	p1Exec(t, db, "INSERT INTO env_config (app,ip) VALUES ('api-A','app-A'),('api-B','app-B'),('api-O','app-O')")
	orch := []DbTaskOrchestration{{TaskOrchestration: TaskOrchestration{ProductList: "override"}}}
	tasks := []DbSchedule{{Schedule: Schedule{ProductList: "B,A,missing"}}}
	configs, err := collectOrchestrationProductConfigs(orch, tasks)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := yaml.Marshal(configs)
	if err != nil {
		t.Fatal(err)
	}
	var restored []CliConfig
	if err = yaml.Unmarshal(raw, &restored); err != nil {
		t.Fatal(err)
	}
	p1Exec(t, db, "DELETE FROM product")
	p1Exec(t, db, "DELETE FROM env_config")
	for _, cfg := range restored {
		if err = upsertCliConfigMode(cfg, cfg.SysParameter, "skip"); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"A", "B", "override"} {
		var count int
		db.Table("product").Where("product = ?", name).Count(&count)
		if count != 1 {
			t.Fatalf("missing product %s", name)
		}
	}
	if len(restored) != 3 {
		t.Fatalf("unexpected configs: %s", raw)
	}
	var count int
	db.Table("env_config").Count(&count)
	if count != 3 {
		t.Fatal("lost app configurations")
	}
}

func TestP1TaskPanicFinishesOrchestration(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO schedule (task_name,task_type) VALUES ('first','data'),('second','scene')")
	p1Exec(t, db, "INSERT INTO task_orchestration (orchestration_name,task_list,orchestration_status) VALUES ('panic-orch',?,'not_started')", "first\nsecond")
	calls := 0
	id, err := runOrchestration("1", "test", func(_ string, _ string, _ TaskRunOptions) (TaskRunResult, error) {
		calls++
		if calls == 1 {
			panic("injected failure")
		}
		return TaskRunResult{Status: "success"}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	var orch DbTaskOrchestration
	db.Table("task_orchestration").Where("id=1").Find(&orch)
	var report DashboardReport
	db.Table("dashboard").Where("id = ?", id).Find(&report)
	var data OrchestrationReportData
	if err = json.Unmarshal([]byte(report.ReportData), &data); err != nil {
		t.Fatal(err)
	}
	if calls != 2 || orch.OrchestrationStatus != "finished" || report.Status != "finished" || data.Steps[0].Status != "fail" || data.Steps[1].Status != "success" {
		t.Fatalf("panic did not finish correctly: calls=%d orch=%+v data=%+v", calls, orch, data)
	}
	if _, loaded := orchestrationRunning.Load("1"); loaded {
		t.Fatal("mutex leaked")
	}
}

func TestP1ConcurrentErrorsRetain401(t *testing.T) {
	var errors concurrentRunErrors
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i == 0 {
				errors.add(fmt.Errorf("auth 401"))
			} else if i%3 == 0 {
				errors.add(fmt.Errorf("ordinary failure"))
			} else {
				errors.add(nil)
			}
		}(i)
	}
	wg.Wait()
	if err := errors.get(); err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("lost 401: %v", err)
	}
}

func p1Zip(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var b bytes.Buffer
	w := zip.NewWriter(&b)
	for name, body := range entries {
		f, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = f.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
func p1Tar(t *testing.T, names []string, content string) []byte {
	t.Helper()
	var b bytes.Buffer
	w := tar.NewWriter(&b)
	for _, name := range names {
		if err := w.WriteHeader(&tar.Header{Name: name, Mode: 0644, Size: int64(len(content))}); err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestP1ArchiveLimitsAndTraversal(t *testing.T) {
	for _, format := range []string{"zip", "tar"} {
		t.Run(format, func(t *testing.T) {
			extract := func(names []string, body string, files int, size int64) error {
				if format == "tar" {
					return extractBoundedTar(bytes.NewReader(p1Tar(t, names, body)), t.TempDir(), files, size)
				}
				entries := map[string]string{}
				for _, name := range names {
					entries[name] = body
				}
				return extractBoundedZip(p1Zip(t, entries), t.TempDir(), files, size)
			}
			if err := extract([]string{"nested/valid.yaml"}, "valid", 2, 10); err != nil {
				t.Fatal(err)
			}
			if err := extract([]string{"a", "b"}, "x", 1, 10); err == nil {
				t.Fatal("file limit bypassed")
			}
			if err := extract([]string{"a", "b"}, "123456", 2, 10); err == nil {
				t.Fatal("size limit bypassed")
			}
			if err := extract([]string{"../outside"}, "x", 2, 10); err == nil {
				t.Fatal("traversal accepted")
			}
		})
	}
}

func TestP1NestedArchiveLimitCleansStaging(t *testing.T) {
	p1Database(t)
	names := make([]string, archiveMaxFiles+1)
	for i := range names {
		names[i] = fmt.Sprintf("file-%d.yaml", i)
	}
	var gz bytes.Buffer
	w := gzip.NewWriter(&gz)
	if _, err := w.Write(p1Tar(t, names, "x")); err != nil {
		t.Fatal(err)
	}
	w.Close()
	bundle := filepath.Join(t.TempDir(), "bundle.tgz")
	os.WriteFile(bundle, gz.Bytes(), 0644)
	if _, err := ImportScheduleCheck(bundle, "test"); err == nil {
		t.Fatal("nested archive count limit bypassed")
	}
	outer := p1Zip(t, map[string]string{"orchestration_info.yaml": "- name: invalid\n  tasks: [task]\n", "bundle/tasks.tgz": string(gz.Bytes())})
	zipPath := filepath.Join(t.TempDir(), "outer.zip")
	if err := os.WriteFile(zipPath, outer, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ImportTaskOrchestrationCheck(zipPath, "test"); err == nil {
		t.Fatal("outer package accepted oversized inner archive")
	}
	outerDirs, _ := filepath.Glob(filepath.Join(DownloadBasePath, "orch_import_*"))
	if len(outerDirs) != 0 {
		t.Fatalf("outer staging leaked: %v", outerDirs)
	}
	matches, _ := filepath.Glob(filepath.Join(DownloadBasePath, "import_stage_*"))
	if len(matches) != 0 {
		t.Fatalf("staging leaked: %v", matches)
	}
}

func TestP1OrchestrationPackageRoundTrip(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO product (product,ip,apps) VALUES ('A','source-A','api'),('B','source-B','api'),('override','source-O','api')")
	p1Exec(t, db, "INSERT INTO env_config (app,ip) VALUES ('api','source-app')")
	p1Exec(t, db, "INSERT INTO scene_data (file_name,name,content,app,file_type,run_time) VALUES ('one.yaml','one','name: one','api',1,1)")
	p1Exec(t, db, "INSERT INTO schedule (task_name,task_type,task_mode,data_list,product_list) VALUES ('data-task','data','once','one.yaml','A,B')")
	p1Exec(t, db, "INSERT INTO task_orchestration (orchestration_name,task_list,product_list) VALUES ('round-trip','data-task','override')")
	name, err := ExportTaskOrchestration("1", "tester")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(DownloadBasePath, name)
	for _, table := range []string{"task_orchestration", "schedule", "scene_data", "product", "env_config"} {
		p1Exec(t, db, "DELETE FROM "+table)
	}
	p1Exec(t, db, "INSERT INTO product (product,ip) VALUES ('A','target-A')")
	checked, err := ImportTaskOrchestrationCheck(path, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if len(checked["orchestration_conflicts"].(map[string]string)) == 0 {
		t.Fatal("product conflict not shown")
	}
	if err = ImportTaskOrchestrationConfirm(checked["import_id"].(string), "skip", "tester"); err != nil {
		t.Fatal(err)
	}
	var task DbSchedule
	db.Table("schedule").Where("task_name='data-task'").Find(&task)
	if task.DataList != "one.yaml" || task.ProductList != "A,B" {
		t.Fatalf("task changed during migration: %+v", task.Schedule)
	}
	var orch DbTaskOrchestration
	db.Table("task_orchestration").Where("orchestration_name='round-trip'").Find(&orch)
	if orch.TaskList != "data-task" || orch.ProductList != "override" {
		t.Fatalf("orchestration changed: %+v", orch)
	}
	var a Product
	db.Table("product").Where("product='A'").Find(&a)
	if a.Ip != "target-A" {
		t.Fatal("skip overwrote target product")
	}
	for _, product := range []string{"B", "override"} {
		var found Product
		db.Table("product").Where("product=?", product).Find(&found)
		if found.Name == "" {
			t.Fatalf("lost product %s", product)
		}
	}
}

func TestP1AuthFailureSkipsRemainingSteps(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO schedule (task_name,task_type) VALUES ('auth','data'),('after','scene')")
	p1Exec(t, db, "INSERT INTO task_orchestration (orchestration_name,task_list,orchestration_status) VALUES ('auth-orch',?,'not_started')", "auth\nafter")
	calls := 0
	id, err := runOrchestration("1", "test", func(_ string, _ string, _ TaskRunOptions) (TaskRunResult, error) {
		calls++
		return TaskRunResult{AuthFailed: true}, fmt.Errorf("401")
	})
	if err != nil {
		t.Fatal(err)
	}
	var report DashboardReport
	db.Table("dashboard").Where("id=?", id).Find(&report)
	var data OrchestrationReportData
	if err = json.Unmarshal([]byte(report.ReportData), &data); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || data.Summary.Fail != 1 || data.Summary.Skipped != 1 || data.Steps[1].Status != "skipped" {
		t.Fatalf("401 did not skip remaining steps: %+v", data)
	}
}

func TestP1DeleteOnlyBlocksRunningOrchestrations(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO task_orchestration (orchestration_name,task_list,orchestration_status) VALUES ('new','task','not_started'),('done','task','finished'),('active','task','running')")
	if err := PreDeleteTaskOrchestrations([]string{"1", "2"}); err != nil {
		t.Fatal(err)
	}
	if err := PreDeleteTaskOrchestrations([]string{"2", "3"}); err == nil {
		t.Fatal("running orchestration accepted for deletion")
	}
	if err := SoftDeleteTaskOrchestrations([]string{"1", "2"}); err != nil {
		t.Fatal(err)
	}
	var count int
	db.Table("task_orchestration").Where("id in (1,2) and deleted_at is not null").Count(&count)
	if count != 2 {
		t.Fatalf("deleted %d orchestrations", count)
	}
	p1Exec(t, db, "INSERT INTO task_orchestration (orchestration_name,task_list,orchestration_status) VALUES ('done','task','not_started')")
}

func TestP1MissingTaskDependenciesFailWithoutExecuting(t *testing.T) {
	for _, test := range []struct {
		name, kind, items string
		deleted           bool
	}{
		{"missing-data", "data", "missing.yaml", false}, {"partial-data", "data", "valid.yaml,missing.yaml", false},
		{"missing-scene", "scene", "missing-scene", false}, {"partial-scene", "scene", "valid-scene,missing-scene", false},
		{"deleted-data", "data", "valid.yaml", true}, {"deleted-scene", "scene", "valid-scene", true},
		{"empty", "data", "", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			db := p1Database(t)
			p1Exec(t, db, "INSERT INTO scene_data (file_name,name) VALUES ('valid.yaml','valid')")
			p1Exec(t, db, "INSERT INTO playbook (name) VALUES ('valid-scene')")
			if test.deleted {
				p1Exec(t, db, "UPDATE scene_data SET deleted_at='2026-01-01'")
				p1Exec(t, db, "UPDATE playbook SET deleted_at='2026-01-01'")
			}
			column := "data_list"
			if test.kind == "scene" {
				column = "scene_list"
			}
			p1Exec(t, db, "INSERT INTO schedule (task_name,task_type,task_status,"+column+") VALUES (?,?,'not_started',?)", test.name, test.kind, test.items)
			result, err := RunTaskOnce("1", "test", TaskRunOptions{RunTag: "missing-check", SyncReport: true, StrictExecution: true})
			if err == nil || result.Status != "fail" || result.SuccessCount != 0 || result.FailCount != 0 {
				t.Fatalf("err=%v result=%+v", err, result)
			}
			if result.DashboardId == "" {
				t.Fatal("failed step has no task report")
			}
			var task DbSchedule
			db.Table("schedule").Where("id=1").Find(&task)
			if task.TaskStatus != "not_started" {
				t.Fatal("orchestration changed task status")
			}
		})
	}
}

func TestP1MetadataFailuresPropagateAndRollback(t *testing.T) {
	for _, table := range []string{"product", "env_config", "sys_parameter", "assert_template"} {
		t.Run(table, func(t *testing.T) {
			db := p1Database(t)
			cfg := CliConfig{Product: CliProduct{Name: "A", Ip: "new"}, App: map[string]CliApp{"api": {Ip: "new"}}, SysParameter: map[string]CliJSONField{"token": "new"}, AssertTemplates: map[string]CliJSONField{"check": "new"}}
			p1Exec(t, db, "CREATE TRIGGER reject_insert BEFORE INSERT ON "+table+" BEGIN SELECT RAISE(ABORT, 'injected write failure'); END")
			err := upsertCliConfigMode(cfg, cfg.SysParameter, "skip")
			if err == nil || !strings.Contains(err.Error(), table) {
				t.Fatalf("write error not propagated: %v", err)
			}
			for _, entity := range []string{"product", "env_config", "sys_parameter", "assert_template"} {
				var count int
				db.Table(entity).Count(&count)
				if count != 0 {
					t.Fatalf("partial metadata remained in %s", entity)
				}
			}
			p1Exec(t, db, "DROP TRIGGER reject_insert")
			if err := upsertCliConfigMode(cfg, cfg.SysParameter, "skip"); err != nil {
				t.Fatal(err)
			}
			p1Exec(t, db, "CREATE TRIGGER reject_update BEFORE UPDATE ON "+table+" BEGIN SELECT RAISE(ABORT, 'injected update failure'); END")
			cfg.Product.Ip = "changed"
			cfg.App["api"] = CliApp{Ip: "changed"}
			cfg.SysParameter["token"] = "changed"
			cfg.AssertTemplates["check"] = "changed"
			if err := upsertCliConfigMode(cfg, cfg.SysParameter, "overwrite"); err == nil {
				t.Fatal("overwrite error not propagated")
			}
			for _, query := range []string{"SELECT ip FROM product WHERE product='A'", "SELECT ip FROM env_config WHERE app='api'", "SELECT value_list FROM sys_parameter WHERE name='token'", "SELECT value FROM assert_template WHERE name='check'"} {
				var got string
				if err := db.Raw(query).Row().Scan(&got); err != nil {
					t.Fatal(err)
				}
				if got != "new" {
					t.Fatalf("failed overwrite changed metadata: %s=%s", query, got)
				}
			}
		})
	}
}

func TestP1ImportStatusReportsMetadataFailure(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO schedule (task_name) VALUES ('task')")
	p1Exec(t, db, "CREATE TRIGGER reject_product BEFORE INSERT ON product BEGIN SELECT RAISE(ABORT, 'injected write failure'); END")
	stage := t.TempDir()
	raw, err := yaml.Marshal(CliConfig{Product: CliProduct{Name: "A"}})
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(stage, "env_config.yaml"), raw, 0644); err != nil {
		t.Fatal(err)
	}
	inner := ImportResult{ImportId: "11111111-1111-1111-1111-111111111111", StagingDir: stage, Conflicts: map[string]string{}}
	raw, err = json.Marshal(inner)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(DownloadBasePath, "import_"+inner.ImportId+".json"), raw, 0644); err != nil {
		t.Fatal(err)
	}
	outerStage, err := os.MkdirTemp(DownloadBasePath, "orch_test_")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(outerStage) })
	state := OrchestrationImportState{ImportId: "22222222-2222-2222-2222-222222222222", ScheduleImportId: inner.ImportId, TempDir: outerStage, Status: "pending", Orchestrations: []KOrchestration{{Name: "orch", Tasks: []string{"task"}, Products: "A"}}}
	if err = saveOrchestrationImportState(state); err != nil {
		t.Fatal(err)
	}
	if err = ImportTaskOrchestrationConfirm(state.ImportId, "skip", "tester"); err == nil {
		t.Fatal("import succeeded after metadata write failed")
	}
	got, err := loadOrchestrationImportState(state.ImportId)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "error" || !strings.Contains(got.Message, "injected write failure") {
		t.Fatalf("incorrect failure status %+v", got)
	}
	var count int
	db.Table("task_orchestration").Count(&count)
	if count != 0 {
		t.Fatal("orchestration saved despite failed metadata")
	}
}

func TestP1ProductlessPackagePreservesEnvironment(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "INSERT INTO env_config (app,ip,thread_number) VALUES ('api','source-app',7)")
	p1Exec(t, db, "INSERT INTO scene_data (file_name,name,content,app,file_type,run_time) VALUES ('one.yaml','one','name: one','api',1,1)")
	p1Exec(t, db, "INSERT INTO sys_parameter (name,value_list) VALUES ('scriptRunEngine','source-engine')")
	p1Exec(t, db, "INSERT INTO schedule (task_name,task_type,task_mode,data_list,product_list) VALUES ('data-task','data','once','one.yaml','')")
	p1Exec(t, db, "INSERT INTO task_orchestration (orchestration_name,task_list,product_list) VALUES ('no-product','data-task','')")
	name, err := ExportTaskOrchestration("1", "test")
	if err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"task_orchestration", "schedule", "scene_data", "env_config", "sys_parameter"} {
		p1Exec(t, db, "DELETE FROM "+table)
	}
	checked, err := ImportTaskOrchestrationCheck(filepath.Join(DownloadBasePath, name), "test")
	if err != nil {
		t.Fatal(err)
	}
	if err = ImportTaskOrchestrationConfirm(checked["import_id"].(string), "skip", "test"); err != nil {
		t.Fatal(err)
	}
	var app EnvConfig
	db.Table("env_config").Where("app='api'").Find(&app)
	if app.Ip != "source-app" || app.MaxThreadNum != 7 {
		t.Fatalf("lost app configuration: %+v", app)
	}
	var param SysParameter
	db.Table("sys_parameter").Where("name='scriptRunEngine'").Find(&param)
	if param.ValueList != "source-engine" {
		t.Fatalf("lost system parameter: %+v", param)
	}
	var products int
	db.Table("product").Count(&products)
	if products != 0 {
		t.Fatal("created an unnamed product")
	}
}
