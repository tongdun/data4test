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

// SyncEntityI18n 按选中任务导出中文骨架：合并选中任务的任务名 + 关联场景 + 关联数据（保留历史条目，不覆盖）
// 返回（任务数, 数据数, 场景数, error）；仅合并 zh-CN.yaml，保留 en-US.yaml 等其他语种文件
func SyncEntityI18n(idStr string) (taskCount, dataCount, playbookCount int, err error) {
	keys, err := collectEntityI18nKeys(idStr)
	if err != nil {
		return 0, 0, 0, err
	}

	if err := writeEntityI18n(i18nTaskDir(), keys.taskSet); err != nil {
		return 0, 0, 0, err
	}
	if err := writeEntityI18n(i18nPlaybookDir(), keys.playbookSet); err != nil {
		return 0, 0, 0, err
	}
	if err := writeEntityI18n(i18nDataDir(), keys.dataSet); err != nil {
		return 0, 0, 0, err
	}

	ReloadEntityI18n()
	return keys.taskCount, keys.dataCount, keys.playbookCount, nil
}

// entityI18nKeys 任务链路（task → playbook → data）上用到的语料 key 集合（中文原文）及命中记录数。
type entityI18nKeys struct {
	taskSet       map[string]string
	playbookSet   map[string]string
	dataSet       map[string]string
	taskCount     int
	dataCount     int
	playbookCount int
}

// collectEntityI18nKeys 收集选中任务链路（task → playbook → data）上的语料 key 集合。
// playbook.last_file 必然落在 data_file_list 中，故无需单独计算。
// 供 SyncEntityI18n（写中文骨架）与离线导出（按需过滤语料）共用。
func collectEntityI18nKeys(idStr string) (keys entityI18nKeys, err error) {
	ids := parseIds(idStr)
	if len(ids) == 0 {
		return keys, E("common.btn_select_first")
	}

	var schedules []Schedule
	models.Orm.Table("schedule").Where("id IN (?)", ids).Find(&schedules)
	if len(schedules) == 0 {
		return keys, E("error.task_info_not_found", idStr)
	}

	keys.taskSet = make(map[string]string)
	keys.playbookSet = make(map[string]string)
	keys.dataSet = make(map[string]string)

	add := func(set map[string]string, s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		set[s] = s
	}

	var sceneNames, fileNames []string
	var playbooks []DbScene
	var datas []DbSceneData

	for _, s := range schedules {
		add(keys.taskSet, s.TaskName)
		add(keys.taskSet, s.Remark)
		for _, name := range splitNames(s.SceneList) {
			add(keys.playbookSet, name)
			sceneNames = append(sceneNames, name)
		}
		for _, name := range splitNames(s.DataList) {
			add(keys.dataSet, name)
			fileNames = append(fileNames, name)
		}
	}

	// 关联场景的描述 + 场景关联的数据文件（链路：task → playbook → data）
	if len(sceneNames) > 0 {
		models.Orm.Table("playbook").Where("name IN (?)", sceneNames).Find(&playbooks)
		for _, p := range playbooks {
			add(keys.playbookSet, p.Name)
			add(keys.playbookSet, p.Remark)
			for _, name := range splitNames(p.DataFileList) {
				add(keys.dataSet, name)
				fileNames = append(fileNames, name)
			}
		}
	}
	// 关联数据的名称与描述
	if len(fileNames) > 0 {
		models.Orm.Table("scene_data").Where("file_name IN (?)", fileNames).Find(&datas)
		for _, d := range datas {
			add(keys.dataSet, d.FileName)
			add(keys.dataSet, d.Name)
			add(keys.dataSet, d.Remark)
		}
	}

	keys.taskCount = len(schedules)
	keys.dataCount = len(datas)
	keys.playbookCount = len(playbooks)
	return keys, nil
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

// writeEntityI18n 合并写某个目录下的 zh-CN.yaml（身份映射骨架），保留其他语种文件
// 采用合并而非覆盖，避免跨多次同步（每次只选部分任务）时丢失历史任务/场景/数据名的译文
func writeEntityI18n(dir string, data map[string]string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		Logger.Error("entity_i18n sync mkdir %s failed: %v", dir, err)
		return err
	}
	merged, err := readFlatYaml(filepath.Join(dir, "zh-CN.yaml"))
	if err != nil {
		Logger.Error("entity_i18n sync read zh-CN.yaml failed: %v", err)
		merged = map[string]string{}
	}
	for k, v := range data {
		merged[k] = v
	}
	bytes, err := yaml.Marshal(merged)
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
