package biz

import "testing"

// 产品维度任务成功/失败按产品配置隔离：任务因其他产品失败而整体失败时，
// 本产品执行项全部通过仍应计为成功任务。
func TestOrchestrationProductTaskStatsIsolatedByProduct(t *testing.T) {
	db := p1Database(t)
	p1Exec(t, db, "CREATE TABLE scene_test_history (id INTEGER PRIMARY KEY AUTOINCREMENT, task_id TEXT, product TEXT, result TEXT)")
	p1Exec(t, db, "CREATE TABLE scene_data_test_history (id INTEGER PRIMARY KEY AUTOINCREMENT, task_id TEXT, product TEXT, result TEXT)")

	for _, query := range []string{
		"INSERT INTO scene_test_history (task_id,product,result) VALUES ('t1','A','pass'),('t1','A','pass'),('t2','A','fail')",
		"INSERT INTO scene_data_test_history (task_id,product,result) VALUES ('t1','A','pass'),('t1','B','fail')",
	} {
		p1Exec(t, db, query)
	}

	steps := []OrchestrationStep{
		{StepOrder: 1, TaskId: "1", TaskTag: "t1", Status: OrchestrationStepFail},
		{StepOrder: 2, TaskId: "2", TaskTag: "t2", Status: OrchestrationStepFail},
	}
	reports := map[int]TaskReportData{
		0: {ByProduct: []ProductReportItem{
			{Product: "A", Total: 3, Pass: 3},
			{Product: "B", Total: 1, Fail: 1},
		}},
		1: {ByProduct: []ProductReportItem{
			{Product: "A", Total: 1, Fail: 1},
		}},
	}

	products := collectOrchestrationProducts(steps, reports, "", "orch1")
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
	a, b := products[0], products[1]
	if a.Product != "A" || a.TaskTotal != 2 || a.TaskSuccess != 1 || a.TaskFail != 1 {
		t.Fatalf("product A task stats = %+v", a)
	}
	if a.SceneTotal != 3 || a.SceneSuccess != 2 || a.SceneFail != 1 {
		t.Fatalf("product A scene stats = %+v", a)
	}
	if a.DataTotal != 1 || a.DataSuccess != 1 || a.DataFail != 0 || a.DataPassRate != 100 {
		t.Fatalf("product A data stats = %+v", a)
	}
	if b.Product != "B" || b.TaskTotal != 1 || b.TaskSuccess != 0 || b.TaskFail != 1 {
		t.Fatalf("product B task stats = %+v", b)
	}
	if b.DataTotal != 1 || b.DataSuccess != 0 || b.DataFail != 1 || b.DataPassRate != 0 {
		t.Fatalf("product B data stats = %+v", b)
	}
}
