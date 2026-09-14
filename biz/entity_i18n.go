package biz

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	dataI18nIndex       map[string]map[string]string // [lang][中文] → 译文
	playbookI18nIndex   map[string]map[string]string
	taskI18nIndex       map[string]map[string]string
	entityI18nMu        sync.RWMutex
	entityI18nLoaded    bool
	entityI18nLastCheck time.Time
)

const entityI18nReloadTTL = 5 * time.Second

// 三个实体各自独立的多语种数据目录
func i18nDataDir() string     { return fmt.Sprintf("%s/i18n/i18n_data", BASEPATH) }
func i18nPlaybookDir() string { return fmt.Sprintf("%s/i18n/i18n_playbook", BASEPATH) }
func i18nTaskDir() string     { return fmt.Sprintf("%s/i18n/i18n_task", BASEPATH) }

// loadEntityI18nDir 扫描 <dir>/<语种>.yaml，构建 [lang][中文]→译文 索引
// 每语种一个文件，内容为扁平「中文 → 译文」映射
func loadEntityI18nDir(dir string) map[string]map[string]string {
	index := make(map[string]map[string]string)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return index
	}
	for _, f := range files {
		name := f.Name()
		if f.IsDir() || (!strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml")) {
			continue
		}
		base := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")
		lang := normalizeLang(base)

		data, err := ioutil.ReadFile(filepath.Join(dir, name))
		if err != nil {
			Logger.Error("entity_i18n read %s failed: %v", name, err)
			continue
		}
		var m map[string]string
		if err := yaml.Unmarshal(data, &m); err != nil {
			Logger.Error("entity_i18n parse %s failed: %v", name, err)
			continue
		}
		index[lang] = m
	}
	return index
}

// ReloadEntityI18n 重新加载场景/数据/任务列表多语种索引
func ReloadEntityI18n() {
	entityI18nMu.Lock()
	dataI18nIndex = loadEntityI18nDir(i18nDataDir())
	playbookI18nIndex = loadEntityI18nDir(i18nPlaybookDir())
	taskI18nIndex = loadEntityI18nDir(i18nTaskDir())
	entityI18nLoaded = true
	entityI18nLastCheck = time.Now()
	entityI18nMu.Unlock()
}

func ensureEntityI18nLoaded() {
	entityI18nMu.RLock()
	loaded := entityI18nLoaded
	lastCheck := entityI18nLastCheck
	entityI18nMu.RUnlock()
	if loaded && time.Since(lastCheck) < entityI18nReloadTTL {
		return
	}
	ReloadEntityI18n()
}

// getLocalized 从指定索引取译文，未命中回退中文原值
func getLocalized(index map[string]map[string]string, zhValue, lang string) string {
	if zhValue == "" {
		return ""
	}
	if lang == "" || lang == "zh-CN" {
		return zhValue
	}
	ensureEntityI18nLoaded()
	entityI18nMu.RLock()
	defer entityI18nMu.RUnlock()
	if byLang, ok := index[lang]; ok {
		if v, ok := byLang[zhValue]; ok && v != "" {
			return v
		}
	}
	return zhValue
}

// GetDataLocalized 数据名/数据文件名/数据描述
func GetDataLocalized(zhValue, lang string) string {
	return getLocalized(dataI18nIndex, zhValue, lang)
}

// GetPlaybookLocalized 场景名/场景描述
func GetPlaybookLocalized(zhValue, lang string) string {
	return getLocalized(playbookI18nIndex, zhValue, lang)
}

// GetTaskLocalized 任务名/任务描述
func GetTaskLocalized(zhValue, lang string) string {
	return getLocalized(taskI18nIndex, zhValue, lang)
}
