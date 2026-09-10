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
	caseI18nIndex     CaseI18nIndex
	caseI18nMu        sync.RWMutex
	caseI18nLoaded    bool
	caseI18nLastCheck time.Time
)

const caseI18nReloadTTL = 5 * time.Second

// caseI18nDir 多语种数据目录
func caseI18nDir() string {
	return fmt.Sprintf("%s/i18n_case", BASEPATH)
}

// normalizeLang 将语种后缀规范化为 BCP-47 标准大小写（en-Us/en-us → en-US，zh-cn → zh-CN）
// 系统标准语种为 "en-US"/"zh-CN"，文件名大小写不一致时应归一化后再做索引 key
func normalizeLang(lang string) string {
	parts := strings.Split(lang, "-")
	if len(parts) == 0 || parts[0] == "" {
		return lang
	}
	parts[0] = strings.ToLower(parts[0])
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.ToUpper(parts[i])
	}
	return strings.Join(parts, "-")
}

// loadCaseI18nDir 扫描 <dir>/<模块>/<语种>.yaml，构建 [lang][module][caseNumber][field] 索引（纯函数）
// 每个模块一个子目录，目录名 = DB module 值；目录内每语种一个文件，避免语种增多后单目录文件过多
func loadCaseI18nDir(dir string) CaseI18nIndex {
	index := make(CaseI18nIndex)
	moduleEntries, err := ioutil.ReadDir(dir)
	if err != nil {
		return index
	}
	for _, mod := range moduleEntries {
		if !mod.IsDir() {
			continue
		}
		module := mod.Name()
		moduleDir := filepath.Join(dir, module)

		files, err := ioutil.ReadDir(moduleDir)
		if err != nil {
			Logger.Error("case_i18n read dir %s failed: %v", moduleDir, err)
			continue
		}
		for _, f := range files {
			name := f.Name()
			if f.IsDir() || (!strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml")) {
				continue
			}
			base := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")
			// 文件名形如 <语种>.yaml，或冗余带模块前缀的 <模块>.<语种>.yaml，取最后一段作语种
			lang := base
			if idx := strings.LastIndex(base, "."); idx >= 0 && idx < len(base)-1 {
				lang = base[idx+1:]
			}
			lang = normalizeLang(lang)

			data, err := ioutil.ReadFile(filepath.Join(moduleDir, name))
			if err != nil {
				Logger.Error("case_i18n read %s failed: %v", name, err)
				continue
			}
			var m map[string]map[string]string
			if err := yaml.Unmarshal(data, &m); err != nil {
				Logger.Error("case_i18n parse %s failed: %v", name, err)
				continue
			}
			if index[lang] == nil {
				index[lang] = make(map[string]map[string]map[string]string)
			}
			index[lang][module] = m
		}
	}
	return index
}

// ReloadCaseI18n 重新加载多语种索引（启动 / 手动刷新调用）
func ReloadCaseI18n() {
	caseI18nMu.Lock()
	caseI18nIndex = loadCaseI18nDir(caseI18nDir())
	caseI18nLoaded = true
	caseI18nLastCheck = time.Now()
	caseI18nMu.Unlock()
}

// ensureCaseI18nLoaded 惰性加载：TTL 内不重复扫描，过期后重扫（捕获文件内容变更）
func ensureCaseI18nLoaded() {
	caseI18nMu.RLock()
	loaded := caseI18nLoaded
	lastCheck := caseI18nLastCheck
	caseI18nMu.RUnlock()
	if loaded && time.Since(lastCheck) < caseI18nReloadTTL {
		return
	}
	ReloadCaseI18n()
}

// GetCaseLocalized 取用例字段译文，未命中返回 ""（调用方回退 DB 中文值）
func GetCaseLocalized(caseNumber, module, lang, field string) string {
	ensureCaseI18nLoaded()
	caseI18nMu.RLock()
	defer caseI18nMu.RUnlock()
	if byModule, ok := caseI18nIndex[lang]; ok {
		if byCase, ok := byModule[module]; ok {
			if byField, ok := byCase[caseNumber]; ok {
				return byField[field]
			}
		}
	}
	return ""
}

// GetDataLocale 解析有效数据语种：cookie 覆盖 → UI 语种 → zh-CN
func GetDataLocale(cookieValue, uiLocale string) string {
	if cookieValue != "" && cookieValue != "auto" {
		return cookieValue
	}
	if uiLocale != "" {
		return uiLocale
	}
	return "zh-CN"
}
