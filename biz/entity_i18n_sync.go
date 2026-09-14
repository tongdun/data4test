package biz

import (
	"data4test/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// SyncEntityI18n 按选中任务导出中文骨架：仅覆盖选中任务的任务名 + 关联场景 + 关联数据
// 返回（任务数, 数据数, 场景数, error）；仅重写 zh-CN.yaml，保留 en-US.yaml 等其他语种文件
func SyncEntityI18n(idStr string) (taskCount, dataCount, playbookCount int, err error) {
	ids := parseIds(idStr)
	if len(ids) == 0 {
		return 0, 0, 0, E("common.btn_select_first")
	}

	var schedules []Schedule
	models.Orm.Table("schedule").Where("id IN (?)", ids).Find(&schedules)
	if len(schedules) == 0 {
		return 0, 0, 0, E("error.task_info_not_found", idStr)
	}

	taskSet := make(map[string]string)
	playbookSet := make(map[string]string)
	dataSet := make(map[string]string)
	var sceneNames, fileNames []string

	add := func(set map[string]string, s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		set[s] = s
	}

	for _, s := range schedules {
		add(taskSet, s.TaskName)
		add(taskSet, s.Remark)
		for _, name := range splitNames(s.SceneList) {
			add(playbookSet, name)
			sceneNames = append(sceneNames, name)
		}
		for _, name := range splitNames(s.DataList) {
			add(dataSet, name)
			fileNames = append(fileNames, name)
		}
	}

	// 关联场景的描述 + 场景关联的数据文件（链路：task → playbook → data）
	if len(sceneNames) > 0 {
		var playbooks []DbScene
		models.Orm.Table("playbook").Where("name IN (?)", sceneNames).Find(&playbooks)
		for _, p := range playbooks {
			add(playbookSet, p.Name)
			add(playbookSet, p.Remark)
			for _, name := range splitNames(p.DataFileList) {
				add(dataSet, name)
				fileNames = append(fileNames, name)
			}
			if strings.TrimSpace(p.LastFile) != "" {
				add(dataSet, p.LastFile)
				fileNames = append(fileNames, p.LastFile)
			}
		}
	}
	// 关联数据的名称与描述
	if len(fileNames) > 0 {
		var datas []DbSceneData
		models.Orm.Table("scene_data").Where("file_name IN (?)", fileNames).Find(&datas)
		for _, d := range datas {
			add(dataSet, d.FileName)
			add(dataSet, d.Name)
			add(dataSet, d.Remark)
		}
	}

	if err := writeEntityI18n(i18nTaskDir(), taskSet); err != nil {
		return 0, 0, 0, err
	}
	if err := writeEntityI18n(i18nPlaybookDir(), playbookSet); err != nil {
		return 0, 0, 0, err
	}
	if err := writeEntityI18n(i18nDataDir(), dataSet); err != nil {
		return 0, 0, 0, err
	}

	ReloadEntityI18n()
	return len(taskSet), len(dataSet), len(playbookSet), nil
}

func parseIds(idStr string) []int {
	var ids []int
	for _, s := range strings.Split(idStr, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if n, err := strconv.Atoi(s); err == nil {
			ids = append(ids, n)
		}
	}
	return ids
}

func splitNames(s string) []string {
	var out []string
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

// writeEntityI18n 写某个目录下的 zh-CN.yaml（身份映射骨架），保留其他语种文件
func writeEntityI18n(dir string, data map[string]string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		Logger.Error("entity_i18n sync mkdir %s failed: %v", dir, err)
		return err
	}
	bytes, err := yaml.Marshal(data)
	if err != nil {
		Logger.Error("entity_i18n sync marshal failed: %v", err)
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "zh-CN.yaml"), bytes, 0644); err != nil {
		Logger.Error("entity_i18n sync write zh-CN.yaml failed: %v", err)
		return err
	}
	return nil
}
