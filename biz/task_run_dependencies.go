package biz

import (
	"data4test/models"
	"fmt"
	"strings"
)

// Resolve the full list before executing anything; a missing dependency must not
// turn a partial or empty execution into a successful step.
func resolveTaskRunItems(task DbSchedule) (ids []string, total int, err error) {
	var configured string
	switch task.TaskType {
	case "data":
		configured = task.DataList
	case "scene":
		configured = task.SceneList
	default:
		return nil, 0, fmt.Errorf(T("error.task_type_invalid"), task.TaskType)
	}
	var names []string
	for _, name := range strings.Split(configured, ",") {
		if name = strings.TrimSpace(name); name != "" {
			names = append(names, name)
		}
	}
	total = len(names)
	if total == 0 {
		return nil, total, fmt.Errorf(T("error.task_items_empty"), task.TaskName)
	}
	for _, name := range names {
		var rows []struct {
			Id string `gorm:"column:id"`
		}
		table, column, key := "scene_data", "file_name", "error.data_name_not_found"
		if task.TaskType == "scene" {
			table, column, key = "playbook", "name", "error.scene_not_found"
		}
		if queryErr := models.Orm.Table(table).Select("id").Where(column+" = ? and deleted_at is null", name).Find(&rows).Error; queryErr != nil {
			return nil, total, queryErr
		}
		if len(rows) == 0 {
			return nil, total, fmt.Errorf(T(key), name)
		}
		ids = append(ids, rows[0].Id)
	}
	return ids, total, nil
}
