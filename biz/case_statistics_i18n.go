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
	caseCountI18nIndex     map[string]map[string]string // [lang][中文] → 译文
	caseCountI18nMu        sync.RWMutex
	caseCountI18nLoaded    bool
	caseCountI18nLastCheck time.Time
)

const caseCountI18nReloadTTL = 5 * time.Second

// caseCountI18nDir 用例统计多语种数据目录
func caseCountI18nDir() string {
	return fmt.Sprintf("%s/i18n/i18n_caseCount", BASEPATH)
}

// loadCaseCountI18nDir 扫描 <dir>/<语种>.yaml，构建 [lang][中文]→译文 索引
// 每语种一个文件，内容为扁平「中文 → 译文」映射，覆盖一级/二级菜单名 + 三级模块名
func loadCaseCountI18nDir(dir string) map[string]map[string]string {
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
			Logger.Error("case_statistics i18n read %s failed: %v", name, err)
			continue
		}
		var m map[string]string
		if err := yaml.Unmarshal(data, &m); err != nil {
			Logger.Error("case_statistics i18n parse %s failed: %v", name, err)
			continue
		}
		index[lang] = m
	}
	return index
}

// ReloadCaseCountI18n 重新加载用例统计多语种索引
func ReloadCaseCountI18n() {
	caseCountI18nMu.Lock()
	caseCountI18nIndex = loadCaseCountI18nDir(caseCountI18nDir())
	caseCountI18nLoaded = true
	caseCountI18nLastCheck = time.Now()
	caseCountI18nMu.Unlock()
}

func ensureCaseCountI18nLoaded() {
	caseCountI18nMu.RLock()
	loaded := caseCountI18nLoaded
	lastCheck := caseCountI18nLastCheck
	caseCountI18nMu.RUnlock()
	if loaded && time.Since(lastCheck) < caseCountI18nReloadTTL {
		return
	}
	ReloadCaseCountI18n()
}

// GetCaseCountLocalized 取一级/二级菜单名、三级模块名的译文，未命中回退中文原值
func GetCaseCountLocalized(zhValue, lang string) string {
	if zhValue == "" {
		return ""
	}
	if lang == "" || lang == "zh-CN" {
		return zhValue
	}
	ensureCaseCountI18nLoaded()
	caseCountI18nMu.RLock()
	defer caseCountI18nMu.RUnlock()
	if byLang, ok := caseCountI18nIndex[lang]; ok {
		if v, ok := byLang[zhValue]; ok && v != "" {
			return v
		}
	}
	return zhValue
}
