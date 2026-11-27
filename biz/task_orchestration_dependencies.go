package biz

import (
	"data4test/models"
	"fmt"
	"strings"
)

type importTaskDependency struct {
	Id, TaskName, TaskType, SceneList, DataList, ProductList string
	DeletedAt                                                *string
}
type importSceneDependency struct {
	Name, DataFileList string
	DeletedAt          *string
}

// Resolve every task -> scene -> data edge. Confirm uses the same resolver before
// writing dependencies, including dependencies changed since the upload check.
func checkOrchestrationImportDependencies(orchestrations []KOrchestration, imported ImportResult, mode string) ([]OrchestrationCheckDetail, error) {
	tasks := make(map[string]KTask)
	scenes := make(map[string]KPlaybook)
	data := make(map[string]bool)
	for _, task := range imported.Tasks {
		tasks[task.Name] = task
	}
	for _, scene := range imported.Playbooks {
		scenes[scene.Name] = scene
	}
	for _, file := range imported.DataFiles {
		data[file.FileName] = true
	}
	var details []OrchestrationCheckDetail
	for _, orch := range orchestrations {
		detail := OrchestrationCheckDetail{Name: orch.Name}
		var count int
		if err := models.Orm.Table("task_orchestration").Where("orchestration_name = ? and deleted_at is null", orch.Name).Count(&count).Error; err != nil {
			return nil, err
		}
		detail.Exists = count > 0
		products := splitDependencyNames(orch.Products)
		checkData := func(name string) error {
			var rows []struct {
				FileName  string
				DeletedAt *string
			}
			if err := models.Orm.Unscoped().Table("scene_data").Where("file_name = ?", name).Find(&rows).Error; err != nil {
				return err
			}
			usable := false
			for _, row := range rows {
				if row.DeletedAt == nil || *row.DeletedAt == "" {
					usable = true
				}
			}
			// Existing soft-deleted names cannot be recreated by the dependency importer.
			if !usable && (len(rows) > 0 || !data[name]) {
				detail.MissingData = appendUniqueDependency(detail.MissingData, name)
			}
			return nil
		}
		for _, name := range orch.Tasks {
			name = strings.TrimSpace(name)
			var rows []importTaskDependency
			if err := models.Orm.Unscoped().Table("schedule").Where("task_name = ?", name).Find(&rows).Error; err != nil {
				return nil, err
			}
			var target *importTaskDependency
			for i := range rows {
				if rows[i].DeletedAt == nil || *rows[i].DeletedAt == "" {
					target = &rows[i]
					break
				}
			}
			task, inPackage := tasks[name]
			mapping := OrchestrationTaskMapping{TaskName: name, Source: "missing"}
			switch {
			case target != nil:
				mapping.Source = "existing"
				mapping.TaskId = target.Id
			case len(rows) == 0 && inPackage:
				mapping.Source = "import"
			default:
				detail.Missing = appendUniqueDependency(detail.Missing, name)
			}
			detail.Tasks = append(detail.Tasks, mapping)
			if mapping.Source == "missing" {
				continue
			}
			if target != nil && (mode == "skip" || !inPackage) {
				refs := target.DataList
				if target.TaskType == "scene" {
					refs = target.SceneList
				}
				productsValue := target.ProductList
				task = KTask{Name: name, TaskType: target.TaskType, PlaybookList: splitDependencyNames(refs), ProductList: &productsValue}
			}
			if task.ProductList != nil {
				products = append(products, splitDependencyNames(*task.ProductList)...)
			}
			for _, ref := range task.PlaybookList {
				ref = strings.TrimSpace(ref)
				if ref == "" {
					continue
				}
				if task.TaskType != "scene" && task.TaskType != "playbook" {
					if err := checkData(ref); err != nil {
						return nil, err
					}
					continue
				}
				var rows []importSceneDependency
				if err := models.Orm.Unscoped().Table("playbook").Where("name = ?", ref).Find(&rows).Error; err != nil {
					return nil, err
				}
				var targetScene *importSceneDependency
				for i := range rows {
					if rows[i].DeletedAt == nil || *rows[i].DeletedAt == "" {
						targetScene = &rows[i]
						break
					}
				}
				scene, inPackage := scenes[ref]
				if targetScene == nil && (len(rows) > 0 || !inPackage) {
					detail.MissingScenes = appendUniqueDependency(detail.MissingScenes, ref)
					continue
				}
				files := scene.DataList
				if targetScene != nil && (mode == "skip" || !inPackage) {
					files = splitDependencyNames(targetScene.DataFileList)
				}
				for _, file := range files {
					if file = strings.TrimSpace(file); file != "" {
						if err := checkData(file); err != nil {
							return nil, err
						}
					}
				}
			}
		}
		for _, product := range splitDependencyNames(strings.Join(products, ",")) {
			var count int
			if err := models.Orm.Table("product").Where("product = ?", product).Count(&count).Error; err != nil {
				return nil, err
			}
			if count == 0 {
				detail.MissingProducts = append(detail.MissingProducts, product)
			}
		}
		details = append(details, detail)
	}
	return details, nil
}

func appendUniqueDependency(names []string, name string) []string {
	for _, existing := range names {
		if existing == name {
			return names
		}
	}
	return append(names, name)
}

// orchestrationDetailMissing 单条编排是否缺少必需依赖（任务/场景/数据文件）。
// 产品缺失仅提示，不参与阻断判断。
func orchestrationDetailMissing(detail OrchestrationCheckDetail) bool {
	return len(detail.Missing) > 0 || len(detail.MissingScenes) > 0 || len(detail.MissingData) > 0
}

func orchestrationDependencyError(details []OrchestrationCheckDetail) error {
	var missing []string
	for _, detail := range details {
		for _, item := range []struct {
			label string
			names []string
		}{
			{"task_orchestration.check_missing", detail.Missing},
			{"task_orchestration.check_missing_scene", detail.MissingScenes},
			{"task_orchestration.check_missing_data", detail.MissingData},
		} {
			if len(item.names) > 0 {
				missing = append(missing, fmt.Sprintf("%s: %s [%s]", detail.Name, T(item.label), strings.Join(item.names, ", ")))
			}
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf(T("error.orchestration_import_dependencies_missing"), strings.Join(missing, "; "))
}
