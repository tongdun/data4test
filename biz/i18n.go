package biz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/GoAdminGroup/go-admin/modules/config"
)

// Translator 多语言翻译引擎
type Translator struct {
	locale   string
	messages map[string]string
	aliases  map[string]string // key → canonical key，值相同但已从 JSON 中清理的 key
	mu       sync.RWMutex
}

var globalTranslator *Translator

// InitI18n 初始化翻译引擎, locale: "zh-CN" / "en-US"
func InitI18n(locale string) {
	globalTranslator = &Translator{
		locale:   locale,
		messages: make(map[string]string),
		aliases:  make(map[string]string),
	}

	globalTranslator.loadLocale(locale)
}

// loadLocale 从 JSON 文件加载语言包, 深度展开嵌套 key
func (t *Translator) loadLocale(locale string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 清空已有翻译
	t.messages = make(map[string]string)
	t.aliases = make(map[string]string)
	var filePath string
	i18nBasePath := fmt.Sprintf("%s/i18n", BASEPATH)
	switch locale {
	case "en-US", "en":
		filePath = fmt.Sprintf("%s/en-US.json", i18nBasePath)
	default:
		filePath = fmt.Sprintf("%s/zh-CN.json", i18nBasePath)

	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Logger.Error("err: %s", err)
		return
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	rawMessages := make(map[string]interface{})
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return
	}

	// 展开 JSON → 平面 map
	flattenJSON(rawMessages, "", t.messages)

	// 自动构建别名表：遍历已加载的 common.* key，记录其值
	// 允许 biz.T("模块.key") 在 JSON 中已删除该 key 时，仍能找到对应的 common.* 值
	commonValueMap := make(map[string]string) // value → canonical common key
	for k, v := range t.messages {
		if len(k) >= 7 && k[:7] == "common." {
			// 如果同一个值有多个 common.* key，取最短 key 作为 canonical
			if existing, ok := commonValueMap[v]; !ok || len(k) < len(existing) {
				commonValueMap[v] = k
			}
		}
	}
	// 对非 common key，如果值匹配 common.*，则记录别名
	for k, v := range t.messages {
		if len(k) >= 7 && k[:7] == "common." {
			continue
		}
		if canonical, ok := commonValueMap[v]; ok {
			t.aliases[k] = canonical
		}
	}
	t.locale = locale
}

// flattenJSON 将嵌套 JSON 展开为 "a.b.c": "value" 平面 map
func flattenJSON(data map[string]interface{}, prefix string, result map[string]string) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		switch v := value.(type) {
		case string:
			result[fullKey] = v
		case map[string]interface{}:
			flattenJSON(v, fullKey, result)
		default:
			result[fullKey] = fmt.Sprintf("%v", v)
		}
	}
}

// SyncLocale 同步项目翻译器与 go-admin 的当前语言设置
// 当用户在管理后台切换语种时自动更新项目内部的 globalTranslator
func SyncLocale() {
	if globalTranslator == nil {
		return
	}
	gaLang := config.GetLanguage()
	if gaLang == "" {
		return
	}

	// 统一语言 key 格式
	normalizedLang := gaLang
	switch gaLang {
	case "en", "en-US":
		normalizedLang = "en-US"
	case "cn", "zh-CN", "zh-cn":
		normalizedLang = "zh-CN"
	}
	// 快速检查：语言未变化则跳过（不影响并发安全）
	globalTranslator.mu.RLock()
	same := globalTranslator.locale == normalizedLang
	globalTranslator.mu.RUnlock()
	if same {
		return
	}
	// 语言已变化，重新加载翻译文件
	globalTranslator.loadLocale(normalizedLang)
	// 同时刷新菜单翻译，确保新增菜单也被注册到 language.Lang
	RefreshMenuI18n()
}

// T 翻译 key, 支持 %s %d 等格式化占位符
// 查找顺序：messages → aliases → 返回 key 原字符串
func T(key string, args ...interface{}) string {
	// 确保翻译器语言与 go-admin 当前语言一致
	SyncLocale()
	if globalTranslator == nil {
		return key
	}
	globalTranslator.mu.RLock()
	msg, ok := globalTranslator.messages[key]
	if !ok {
		// 尝试通过别名查找（该 key 已从 JSON 中清理，值同 common.*）
		if canonical, hasAlias := globalTranslator.aliases[key]; hasAlias {
			msg, ok = globalTranslator.messages[canonical]
		}
	}
	globalTranslator.mu.RUnlock()
	if !ok {
		return key
	}
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// E 创建翻译后的错误
func E(key string, args ...interface{}) error {
	return fmt.Errorf(T(key, args...))
}

// GetLocale 获取当前语言
func GetLocale() string {
	if globalTranslator == nil {
		return "zh-CN"
	}
	globalTranslator.mu.RLock()
	defer globalTranslator.mu.RUnlock()
	return globalTranslator.locale
}

// SetLocale 切换语言，如果语种未变化则跳过加载
// locale 支持 GoAdmin 的语种值: "en" (→"en-US"), "zh-CN" (→"zh-CN")
func SetLocale(locale string) {
	if globalTranslator == nil {
		return
	}
	normalized := locale
	switch locale {
	case "en", "en-US":
		normalized = "en-US"
	case "cn", "zh-CN", "zh-cn":
		normalized = "zh-CN"
	}
	if normalized == "" {
		return
	}
	globalTranslator.mu.RLock()
	same := globalTranslator.locale == normalized
	globalTranslator.mu.RUnlock()
	if same {
		return
	}
	globalTranslator.loadLocale(normalized)
}
