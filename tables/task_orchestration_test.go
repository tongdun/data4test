package tables

import (
	"data4test/models"
	"net/http"
	"testing"

	"github.com/GoAdminGroup/go-admin/context"
	adminModels "github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	tmpProductEnv   = "kbank_169项目环境"
	tmpProductEn    = "KBank英文环境"
	tmpProductOther = "重构二期测试环境_中文_变更"
)

func setupOrchestrationTableDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	previous := models.Orm
	models.Orm = db
	t.Cleanup(func() {
		db.Close()
		models.Orm = previous
	})

	stmts := []string{
		"CREATE TABLE product (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, updated_at TEXT)",
		"CREATE TABLE schedule (id INTEGER PRIMARY KEY AUTOINCREMENT, task_name TEXT, product_list TEXT, deleted_at TEXT)",
		"CREATE TABLE task_orchestration (id INTEGER PRIMARY KEY AUTOINCREMENT, orchestration_name TEXT, task_list TEXT, " +
			"product_list TEXT, orchestration_status TEXT, last_at TEXT, remark TEXT, user_name TEXT, " +
			"created_at TEXT, updated_at TEXT, deleted_at TEXT)",
		"INSERT INTO schedule (task_name, product_list) VALUES ('T1','" + tmpProductEnv + "')",
		"INSERT INTO schedule (task_name, product_list) VALUES ('T2','" + tmpProductEn + "')",
		"INSERT INTO schedule (task_name, product_list) VALUES ('T4','" + tmpProductOther + "')",
		// 1 未配置产品，任务产品命中；2 编排级产品命中；3 任务产品不命中；4 任务缺失；5 多任务并集命中
		// 状态分别为 running/finished/stopped/not_started/running，用于覆盖任务状态多选筛选
		"INSERT INTO task_orchestration (orchestration_name, task_list, product_list, orchestration_status, user_name, created_at) VALUES ('步骤产品命中','T1','','running','u1','2026-11-01 10:00:00')",
		"INSERT INTO task_orchestration (orchestration_name, task_list, product_list, orchestration_status, user_name, created_at) VALUES ('编排产品命中','T2','" + tmpProductEnv + "," + tmpProductEn + "','finished','u1','2026-11-02 10:00:00')",
		"INSERT INTO task_orchestration (orchestration_name, task_list, product_list, orchestration_status, user_name, created_at) VALUES ('不命中','T2','','stopped','u2','2026-11-03 10:00:00')",
		"INSERT INTO task_orchestration (orchestration_name, task_list, product_list, orchestration_status, user_name, created_at) VALUES ('任务缺失','T3','','not_started','u2','2026-11-04 10:00:00')",
		"INSERT INTO task_orchestration (orchestration_name, task_list, product_list, orchestration_status, user_name, created_at) VALUES ('多任务并集','T1\nT2','','running','u3','2026-11-05 10:00:00')",
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("执行 %q 失败: %v", stmt, err)
		}
	}
}

func orchestrationTableInfo(t *testing.T) *types.InfoPanel {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/info/task_orchestration", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.NewContext(req)
	ctx.SetUserValue("user", adminModels.UserModel{})
	return GetTaskOrchestrationTable(ctx).GetInfo()
}

func listParams(fields map[string][]string, page, size int) parameter.Parameters {
	params := parameter.Parameters{
		Page: "1", PageSize: "10", PageInt: page, PageSizeInt: size,
		SortField: "id", SortType: "asc", Fields: fields,
	}
	return params
}

func listedIDs(data []map[string]interface{}) []int64 {
	ids := make([]int64, 0, len(data))
	for _, row := range data {
		id, _ := row["id"].(int64)
		ids = append(ids, id)
	}
	return ids
}

// 任务编排列表必须走自管查询：筛选、排序、分页由参数生成，产品配置按展示口径匹配。
// 不能退回 SetQueryFilterFn —— stopQuery 会只按 pk in (...) 查询并丢弃筛选与分页。
func TestTaskOrchestrationListKeepsNativeFilterQuery(t *testing.T) {
	setupOrchestrationTableDB(t)
	info := orchestrationTableInfo(t)

	if info.IsHideFilterArea {
		t.Fatal("任务编排列表的筛选区域应默认展开，不再折叠")
	}
	if info.QueryFilterFn != nil {
		t.Fatal("任务编排列表不应使用 QueryFilterFn，否则 stopQuery 会丢弃筛选条件与分页")
	}
	if info.GetDataFn == nil {
		t.Fatal("任务编排列表应使用 GetDataFn 以支持产品配置按展示口径筛选")
	}

	filterable := make(map[string]bool)
	for _, field := range info.FieldList {
		if field.Filterable {
			filterable[field.Field] = true
		}
	}
	for _, field := range []string{"orchestration_name", "product_list", "orchestration_status", "user_name", "created_at"} {
		if !filterable[field] {
			t.Errorf("字段 %s 应支持筛选", field)
		}
	}
	if filterable["id"] {
		t.Error("自增主键 id 不应作为筛选项")
	}

	// 任务状态筛选项不含「已暂停」
	for _, field := range info.FieldList {
		if field.Field != "orchestration_status" {
			continue
		}
		if len(field.FilterFormFields) == 0 {
			t.Fatal("任务状态字段缺少筛选表单配置")
		}
		for _, option := range field.FilterFormFields[0].Options {
			if option.Value == "stopped" {
				t.Error("任务状态筛选项不应包含「已暂停」")
			}
		}
		if len(field.FilterFormFields[0].Options) != 3 {
			t.Errorf("任务状态筛选项应为 3 个，实际 %d 个", len(field.FilterFormFields[0].Options))
		}
	}
}

// 产品配置筛选必须覆盖「未配置编排产品、产品来自步骤任务」的编排
func TestOrchestrationProductFilterMatchesInheritedProducts(t *testing.T) {
	setupOrchestrationTableDB(t)
	info := orchestrationTableInfo(t)

	data, total := info.GetDataFn(listParams(map[string][]string{"product_list": {tmpProductEnv}}, 1, 10))
	if total != 3 {
		t.Fatalf("产品配置筛选应返回 3 条，实际 %d 条: %v", total, listedIDs(data))
	}
	got := listedIDs(data)
	want := []int64{1, 2, 5}
	if len(got) != len(want) {
		t.Fatalf("id 不匹配，want %v got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("id 不匹配，want %v got %v", want, got)
		}
	}
}

// 任务状态为多选筛选：多值必须生成 in (...)，不能只按第一个值筛选
func TestOrchestrationStatusFilterSupportsMultipleValues(t *testing.T) {
	setupOrchestrationTableDB(t)
	info := orchestrationTableInfo(t)

	data, total := info.GetDataFn(listParams(map[string][]string{
		"orchestration_status": {"running", "finished"},
	}, 1, 10))
	if total != 3 {
		t.Fatalf("多选任务状态应返回 3 条，实际 total=%d ids=%v", total, listedIDs(data))
	}
	want := []int64{1, 2, 5}
	got := listedIDs(data)
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("多选任务状态 id 不匹配，want %v got %v", want, got)
		}
	}

	data, total = info.GetDataFn(listParams(map[string][]string{
		"orchestration_status": {"not_started"},
	}, 1, 10))
	if total != 1 || len(data) != 1 || listedIDs(data)[0] != 4 {
		t.Fatalf("单选任务状态应命中 id=4，实际 total=%d ids=%v", total, listedIDs(data))
	}

	// 多选状态与产品配置筛选叠加：产品 env 命中 1/2/5，running 命中 1/5
	data, total = info.GetDataFn(listParams(map[string][]string{
		"orchestration_status": {"running", "finished"},
		"product_list":         {tmpProductEnv},
	}, 1, 10))
	if total != 3 || len(data) != 3 {
		t.Fatalf("多选状态+产品筛选应返回 3 条，实际 total=%d ids=%v", total, listedIDs(data))
	}
}

// 产品配置筛选与其他筛选、分页、排序可叠加
func TestOrchestrationFiltersCombineWithProductFilter(t *testing.T) {
	setupOrchestrationTableDB(t)
	info := orchestrationTableInfo(t)

	data, total := info.GetDataFn(listParams(map[string][]string{
		"product_list":       {tmpProductEnv},
		"orchestration_name": {"步骤"},
	}, 1, 10))
	if total != 1 || len(data) != 1 || listedIDs(data)[0] != 1 {
		t.Fatalf("产品筛选+名称筛选应命中 id=1，实际 total=%d ids=%v", total, listedIDs(data))
	}

	data, total = info.GetDataFn(listParams(map[string][]string{"product_list": {tmpProductEnv}}, 2, 2))
	if total != 3 || len(data) != 1 || listedIDs(data)[0] != 5 {
		t.Fatalf("产品筛选分页应返回第 2 页 id=5，实际 total=%d ids=%v", total, listedIDs(data))
	}

	data, total = info.GetDataFn(parameter.Parameters{
		Page: "1", PageSize: "10", PageInt: 1, PageSizeInt: 10,
		SortField: "orchestration_name", SortType: "desc",
		Fields: map[string][]string{"product_list": {tmpProductEnv}},
	})
	if total != 3 || len(data) != 3 {
		t.Fatalf("产品筛选排序应返回 3 条，实际 total=%d ids=%v", total, listedIDs(data))
	}
	// 按名称排序生效即可，不依赖具体排序规则（sqlite 为二进制序，MySQL 为排序规则序）
	for i := 1; i < len(data); i++ {
		prev, _ := data[i-1]["orchestration_name"].(string)
		cur, _ := data[i]["orchestration_name"].(string)
		if prev < cur {
			t.Fatalf("名称倒序未生效: %v", []interface{}{data[i-1]["orchestration_name"], data[i]["orchestration_name"]})
		}
	}
}
