package biz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/GoAdminGroup/go-admin/template/types"
	"gopkg.in/yaml.v2"
)

// translateTemplateType AI 模板类型「翻译」（ai_template.template_type）
const translateTemplateType = "8"

// buildTranslateQuery 用翻译模板内容拼接 LLM query，替换 {目标语种} / {待翻译内容} 占位符
func buildTranslateQuery(content, lang, jsonStr string) string {
	q := strings.Replace(content, "{目标语种}", lang, -1)
	return strings.Replace(q, "{待翻译内容}", jsonStr, -1)
}

// GetTranslateTargetLanguages 返回可用于翻译的目标语种选项（排除源语种 zh-CN）
func GetTranslateTargetLanguages() []types.FieldOption {
	all := GetSupportLanguages()
	out := make([]types.FieldOption, 0, len(all))
	for _, l := range all {
		if l.Value == "zh-CN" {
			continue
		}
		out = append(out, l)
	}
	return out
}

// filterNonZh 过滤掉 zh-CN（翻译目标不能是源语种）
func filterNonZh(langs []string) []string {
	out := make([]string, 0, len(langs))
	for _, l := range langs {
		if l == "zh-CN" {
			continue
		}
		out = append(out, l)
	}
	return out
}

// extractJSON 从 LLM 回复中剥离 ```json 围栏/多余文字，截取第一个 { 到最后一个 }
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, "```"); i >= 0 {
		s = s[i+3:]
		s = strings.TrimPrefix(strings.TrimSpace(s), "json")
		s = strings.TrimSpace(s)
		if j := strings.Index(s, "```"); j >= 0 {
			s = s[:j]
		}
	}
	s = strings.TrimSpace(s)
	if l := strings.Index(s, "{"); l >= 0 {
		if r := strings.LastIndex(s, "}"); r > l {
			return s[l : r+1]
		}
	}
	return s
}

// translateToJSON 把结构 JSON 化后交给 AI 翻译 value，返回同构的 JSON 字节
func translateToJSON(content, lang, platform, userCode string, v interface{}) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	replies, err := ConnetAIModel(buildTranslateQuery(content, lang, string(raw)), "", "", platform, userCode)
	if err != nil {
		return nil, err
	}
	if len(replies) == 0 || strings.TrimSpace(replies[0]) == "" {
		return nil, E("error.no_reply_info")
	}
	return []byte(extractJSON(replies[0])), nil
}

// writeTranslateYaml 序列化并写盘
func writeTranslateYaml(path string, v interface{}) error {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

// caseSnapDir 返回用例翻译源内容快照目录（位于 i18n/ 之外，避免被 loader 扫描误读）
func caseSnapDir() string {
	return fmt.Sprintf("%s/i18n_snap/i18n_case", BASEPATH)
}

// readCaseModuleYaml 读用例模块 yaml（map[用例编号]map[字段]值），文件不存在返回空 map
func readCaseModuleYaml(path string) (map[string]map[string]string, error) {
	m := map[string]map[string]string{}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		return m, err
	}
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return m, err
	}
	return m, nil
}

// readFlatYaml 读扁平 yaml（map[string]string），文件不存在返回空 map
func readFlatYaml(path string) (map[string]string, error) {
	m := map[string]string{}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		return m, err
	}
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return m, err
	}
	return m, nil
}

// writeSnapYaml 写源内容快照（自动建目录）
func writeSnapYaml(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return writeTranslateYaml(path, v)
}

// StartTranslateIfEnabled 校验「是否翻译/目标语种/AI引擎/翻译模板」，满足则取模板内容后台 go fn
// 返回 translating 表示是否已启动后台翻译；err 表示校验失败（模板缺失/语种或引擎未选）
func StartTranslateIfEnabled(translateEnable, langsStr, platform, templateName, userCode string,
	fn func(userCode, platform, content string, langs []string)) (translating bool, err error) {
	if translateEnable != "1" {
		return false, nil
	}
	langs := filterNonZh(splitNames(langsStr))
	if len(langs) == 0 {
		return false, E("common.select_lang_first")
	}
	if strings.TrimSpace(platform) == "" {
		return false, E("common.select_platform_first")
	}
	content, _, err := GetAiTemplateByName(strings.TrimSpace(templateName), translateTemplateType)
	if err != nil {
		return false, err
	}
	go fn(userCode, platform, content, langs)
	return true, nil
}

// translateFlatI18n 增量翻译扁平身份映射（map[string]string，key 即原文）到目标语种
// 只翻新增/空译文，移除源里已删的 key；未变的复用已有译文
func translateFlatI18n(zhPath, langPath, content, lang, platform, userCode string) error {
	source, err := readFlatYaml(zhPath)
	if err != nil {
		Logger.Error("i18n translate read %s failed: %v", zhPath, err)
		return nil
	}
	existing, err := readFlatYaml(langPath)
	if err != nil {
		Logger.Error("i18n translate read %s failed: %v", langPath, err)
		return nil
	}

	toTranslate := map[string]string{}
	for k, v := range source {
		if tv, ok := existing[k]; !ok || strings.TrimSpace(tv) == "" {
			toTranslate[k] = v
		}
	}
	toRemove := make([]string, 0)
	for k := range existing {
		if _, ok := source[k]; !ok {
			toRemove = append(toRemove, k)
		}
	}

	if len(toTranslate) == 0 && len(toRemove) == 0 {
		return nil
	}

	result := map[string]string{}
	for k, v := range existing {
		result[k] = v
	}
	for _, k := range toRemove {
		delete(result, k)
	}

	translated := map[string]string{}
	if len(toTranslate) > 0 {
		out, err := translateToJSON(content, lang, platform, userCode, toTranslate)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(out, &translated); err != nil {
			return err
		}
	}
	// 仅合并预期 key（忽略 AI 幻觉的多余 key）；翻译失败的 key 保留旧译文（下次重试）
	for k := range toTranslate {
		if tv, ok := translated[k]; ok && strings.TrimSpace(tv) != "" {
			result[k] = tv
		}
	}

	if err := writeTranslateYaml(langPath, result); err != nil {
		Logger.Error("i18n translate write %s failed: %v", langPath, err)
	}
	return nil
}

// translateCaseModuleI18n 增量翻译某个用例模块到目标语种（含源内容快照，检测新增/内容变更/删除）
func translateCaseModuleI18n(moduleDir, content, lang, platform, userCode string) error {
	zhPath := filepath.Join(moduleDir, "zh-CN.yaml")
	langPath := filepath.Join(moduleDir, lang+".yaml")
	snapPath := filepath.Join(caseSnapDir(), filepath.Base(moduleDir), lang+".yaml")

	source, err := readCaseModuleYaml(zhPath)
	if err != nil {
		Logger.Error("case i18n translate read %s failed: %v", zhPath, err)
		return nil
	}
	existing, err := readCaseModuleYaml(langPath)
	if err != nil {
		Logger.Error("case i18n translate read %s failed: %v", langPath, err)
		return nil
	}
	snap, err := readCaseModuleYaml(snapPath)
	if err != nil {
		Logger.Error("case i18n translate read snapshot %s failed: %v", snapPath, err)
		return nil
	}

	// toTranslate：新增（无译文）或内容变更（有旧快照且源值变了）
	// 冷启动（有译文但无快照）不重翻，由下方回填快照兜底
	toTranslate := map[string]map[string]string{}
	for c, sv := range source {
		if len(existing[c]) == 0 {
			toTranslate[c] = sv
		} else if oldSnap, ok := snap[c]; ok && !reflect.DeepEqual(oldSnap, sv) {
			toTranslate[c] = sv
		}
	}
	// toRemove：源里已删的 key
	toRemove := make([]string, 0)
	for c := range existing {
		if _, ok := source[c]; !ok {
			toRemove = append(toRemove, c)
		}
	}

	result := map[string]map[string]string{}
	for c, tv := range existing {
		result[c] = tv
	}
	for _, c := range toRemove {
		delete(result, c)
	}

	// 新快照：以旧快照为底去掉源里已删的 key；冷启动（有译文但无快照）按当前源值回填；
	// 之后成功翻译的新增/变更 key 再覆盖为当前源值
	newSnap := map[string]map[string]string{}
	for c, sv := range snap {
		newSnap[c] = sv
	}
	for c := range newSnap {
		if _, ok := source[c]; !ok {
			delete(newSnap, c)
		}
	}
	for c, sv := range source {
		if len(existing[c]) > 0 {
			if _, ok := newSnap[c]; !ok {
				newSnap[c] = sv
			}
		}
	}

	translated := map[string]map[string]string{}
	if len(toTranslate) > 0 {
		out, err := translateToJSON(content, lang, platform, userCode, toTranslate)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(out, &translated); err != nil {
			return err
		}
	}
	// 仅合并预期 key（忽略 AI 幻觉）；翻译失败的 key 保留旧译文、旧快照（下次重试）
	for c := range toTranslate {
		if tv, ok := translated[c]; ok && len(tv) > 0 {
			result[c] = tv
			newSnap[c] = source[c]
		}
	}

	if err := writeTranslateYaml(langPath, result); err != nil {
		Logger.Error("case i18n translate write %s of %s failed: %v", lang, filepath.Base(moduleDir), err)
		return nil
	}
	if err := writeSnapYaml(snapPath, newSnap); err != nil {
		Logger.Error("case i18n translate write snapshot %s failed: %v", snapPath, err)
	}
	return nil
}

// TranslateCaseI18n 后台把用例统计语料（i18n_case/<模块>/ + i18n_caseCount/）增量翻译到目标语种
func TranslateCaseI18n(userCode, platform, content string, langs []string) {
	defer func() {
		if e := recover(); e != nil {
			Logger.Error("case i18n translate panic: %v\n%s", e, debug.Stack())
		}
	}()
	for _, lang := range langs {
		modules, err := ioutil.ReadDir(caseI18nDir())
		if err != nil {
			Logger.Error("case i18n translate read dir %s failed: %v", caseI18nDir(), err)
			continue
		}
		for _, mod := range modules {
			if !mod.IsDir() {
				continue
			}
			if err := translateCaseModuleI18n(filepath.Join(caseI18nDir(), mod.Name()), content, lang, platform, userCode); err != nil {
				Logger.Error("case i18n translate aborted at %s: %v", mod.Name(), err)
				return
			}
		}

		// 菜单/模块名扁平语料
		if err := translateFlatI18n(filepath.Join(caseCountI18nDir(), "zh-CN.yaml"),
			filepath.Join(caseCountI18nDir(), lang+".yaml"), content, lang, platform, userCode); err != nil {
			Logger.Error("case i18n translate aborted at caseCount: %v", err)
			return
		}
	}
	ReloadCaseI18n()
	ReloadCaseCountI18n()
}

// TranslateEntityI18n 后台把任务实体语料（i18n_task/ i18n_playbook/ i18n_data/）增量翻译到目标语种
func TranslateEntityI18n(userCode, platform, content string, langs []string) {
	defer func() {
		if e := recover(); e != nil {
			Logger.Error("entity i18n translate panic: %v\n%s", e, debug.Stack())
		}
	}()
	for _, lang := range langs {
		for _, dir := range []string{i18nTaskDir(), i18nPlaybookDir(), i18nDataDir()} {
			if err := translateFlatI18n(filepath.Join(dir, "zh-CN.yaml"),
				filepath.Join(dir, lang+".yaml"), content, lang, platform, userCode); err != nil {
				Logger.Error("entity i18n translate aborted at %s: %v", filepath.Base(dir), err)
				return
			}
		}
	}
	ReloadEntityI18n()
}
