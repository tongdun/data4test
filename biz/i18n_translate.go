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
func translateToJSON(content, lang, platform, userCode string, aiConnect AIConnect, v interface{}) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	replies, err := ConnetAIModelWithConnect(buildTranslateQuery(content, lang, string(raw)), "", "", userCode, platform, aiConnect)
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
// 每个模块一个 zh-CN.yaml，记录「上次翻译时的源内容」，用于差异化比较（当前源 vs 快照）
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

// StartTranslateIfEnabled 校验「是否翻译/目标语种/AI引擎/翻译模板」，满足且存在待翻译内容时后台 go fn
// hasPending 可选：返回是否存在待翻译增量；为 nil 表示不预检（始终翻译）
// 返回 translating 表示已启动后台翻译；err 表示校验失败；开启翻译但 hasPending 判定无增量时返回 (false, nil)
func StartTranslateIfEnabled(translateEnable, langsStr, platform, templateName, userCode string,
	override *AIConnectOverride,
	fn func(userCode, platform string, aiConnect AIConnect, content string, langs []string),
	hasPending func(langs []string) bool) (translating bool, err error) {
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
	if hasPending != nil && !hasPending(langs) {
		return false, nil
	}
	aiConnect, err := ResolveAIConnect(strings.TrimSpace(platform), override)
	if err != nil {
		return false, err
	}
	go fn(userCode, platform, aiConnect, content, langs)
	return true, nil
}

// translateFlatI18n 增量翻译扁平身份映射（map[string]string，key 即原文）到目标语种
// 只翻新增/空译文，移除源里已删的 key；未变的复用已有译文
// exclude 中的 key 不参与 LLM 翻译，也不会被当作「已删除」移除（用于保留由其它语料回填的译文，如模块名）
func translateFlatI18n(zhPath, langPath, content, lang, platform, userCode string, aiConnect AIConnect, exclude map[string]bool) error {
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
		if exclude[k] {
			continue
		}
		if tv, ok := existing[k]; !ok || strings.TrimSpace(tv) == "" {
			toTranslate[k] = v
		}
	}
	toRemove := make([]string, 0)
	for k := range existing {
		if exclude[k] {
			continue
		}
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
		out, err := translateToJSON(content, lang, platform, userCode, aiConnect, toTranslate)
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

// translateCaseModuleAllLangs 增量翻译某个用例模块到多个目标语种（含模块级源内容快照 zh-CN.yaml）
// 源变更集只算一次（当前源 vs 共享快照），复用到所有语种；翻译完统一推进快照
// 快照只推进「所有选中语种都已有译文」的用例，单个语种翻译失败/缺失时保持旧值，下次重试
func translateCaseModuleAllLangs(moduleDir, content string, langs []string, platform, userCode string, aiConnect AIConnect) error {
	zhPath := filepath.Join(moduleDir, "zh-CN.yaml")
	snapPath := filepath.Join(caseSnapDir(), filepath.Base(moduleDir), "zh-CN.yaml")

	source, err := readCaseModuleYaml(zhPath)
	if err != nil {
		Logger.Error("case i18n translate read %s failed: %v", zhPath, err)
		return nil
	}
	snap, err := readCaseModuleYaml(snapPath)
	if err != nil {
		Logger.Error("case i18n translate read snapshot %s failed: %v", snapPath, err)
		return nil
	}

	// 源变更集（相对共享快照）：快照里存在且源值变了 → 各语种都需重翻
	changed := map[string]bool{}
	for c, sv := range source {
		if oldSnap, ok := snap[c]; ok && !reflect.DeepEqual(oldSnap, sv) {
			changed[c] = true
		}
	}

	// 各语种最终译文，用于推进快照时判断用例是否在所有选中语种都有译文
	final := make(map[string]map[string]map[string]string, len(langs))
	for _, lang := range langs {
		langPath := filepath.Join(moduleDir, lang+".yaml")
		existing, err := readCaseModuleYaml(langPath)
		if err != nil {
			Logger.Error("case i18n translate read %s failed: %v", langPath, err)
			return nil
		}

		// toTranslate：新增（无译文）或源变更（shared changed）
		toTranslate := map[string]map[string]string{}
		for c, sv := range source {
			if len(existing[c]) == 0 {
				toTranslate[c] = sv
			} else if changed[c] {
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

		translated := map[string]map[string]string{}
		if len(toTranslate) > 0 {
			out, err := translateToJSON(content, lang, platform, userCode, aiConnect, toTranslate)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(out, &translated); err != nil {
				return err
			}
		}
		// 仅合并预期 key（忽略 AI 幻觉）；翻译失败的 key 保留旧译文（下次重试）
		for c := range toTranslate {
			if tv, ok := translated[c]; ok && len(tv) > 0 {
				result[c] = tv
			}
		}

		if err := writeTranslateYaml(langPath, result); err != nil {
			Logger.Error("case i18n translate write %s of %s failed: %v", lang, filepath.Base(moduleDir), err)
			return nil
		}
		final[lang] = result
	}

	// 推进快照：只推进在所有选中语种都有译文的用例（任一语种缺失则保持旧值，下次重试）
	newSnap := map[string]map[string]string{}
	for c, sv := range source {
		complete := true
		for _, lang := range langs {
			if len(final[lang][c]) == 0 {
				complete = false
				break
			}
		}
		if complete {
			newSnap[c] = sv
		}
	}
	if err := writeSnapYaml(snapPath, newSnap); err != nil {
		Logger.Error("case i18n translate write snapshot %s failed: %v", snapPath, err)
	}
	return nil
}

// moduleNameTranslation 从已翻译的用例模块语料中取模块名译文（FieldCaseModule 字段，取首个非空值）
func moduleNameTranslation(module, lang string) string {
	m, err := readCaseModuleYaml(filepath.Join(caseI18nDir(), module, lang+".yaml"))
	if err != nil {
		return ""
	}
	for _, fields := range m {
		if v := strings.TrimSpace(fields[FieldCaseModule]); v != "" {
			return v
		}
	}
	return ""
}

// propagateModuleNamesToCaseCount 把模块名译文从用例内容语料回填到 i18n_caseCount/<lang>.yaml
// 模块名随内容一起翻译，不再单独调用大模型，保证两处译文一致
func propagateModuleNamesToCaseCount(modules []string, lang string) error {
	langPath := filepath.Join(caseCountI18nDir(), lang+".yaml")
	existing, err := readFlatYaml(langPath)
	if err != nil {
		Logger.Error("case i18n translate read %s failed: %v", langPath, err)
		return nil
	}
	changed := false
	for _, m := range modules {
		tv := moduleNameTranslation(m, lang)
		if tv == "" {
			continue
		}
		if existing[m] == tv {
			continue
		}
		existing[m] = tv
		changed = true
	}
	if !changed {
		return nil
	}
	if err := writeTranslateYaml(langPath, existing); err != nil {
		Logger.Error("case i18n translate write %s failed: %v", langPath, err)
	}
	return nil
}

// hasCaseModulePending 判断某模块某语种是否存在增量差异（新增/变更/删除），与 translateCaseModuleAllLangs 的判定保持一致
func hasCaseModulePending(module, lang string) bool {
	zhPath := filepath.Join(caseI18nDir(), module, "zh-CN.yaml")
	langPath := filepath.Join(caseI18nDir(), module, lang+".yaml")
	snapPath := filepath.Join(caseSnapDir(), module, "zh-CN.yaml")
	source, _ := readCaseModuleYaml(zhPath)
	existing, _ := readCaseModuleYaml(langPath)
	snap, _ := readCaseModuleYaml(snapPath)
	for c, sv := range source {
		if len(existing[c]) == 0 {
			return true
		}
		if oldSnap, ok := snap[c]; ok && !reflect.DeepEqual(oldSnap, sv) {
			return true
		}
	}
	for c := range existing {
		if _, ok := source[c]; !ok {
			return true
		}
	}
	return false
}

// hasFlatPending 判断扁平语料（排除 exclude 中的 key）是否存在增量差异，与 translateFlatI18n 的判定保持一致
func hasFlatPending(dir, lang string, exclude map[string]bool) bool {
	source, _ := readFlatYaml(filepath.Join(dir, "zh-CN.yaml"))
	existing, _ := readFlatYaml(filepath.Join(dir, lang+".yaml"))
	for k := range source {
		if exclude[k] {
			continue
		}
		if tv, ok := existing[k]; !ok || strings.TrimSpace(tv) == "" {
			return true
		}
	}
	for k := range existing {
		if exclude[k] {
			continue
		}
		if _, ok := source[k]; !ok {
			return true
		}
	}
	return false
}

// HasPendingCaseI18nTranslate 判断本次选中模块（含一级/二级菜单名）是否存在待翻译增量，不做任何翻译
func HasPendingCaseI18nTranslate(modules, langs []string) bool {
	exclude := make(map[string]bool, len(modules))
	for _, m := range modules {
		exclude[m] = true
	}
	for _, lang := range langs {
		for _, m := range modules {
			if hasCaseModulePending(m, lang) {
				return true
			}
		}
		if hasFlatPending(caseCountI18nDir(), lang, exclude) {
			return true
		}
	}
	return false
}

// TranslateCaseI18n 后台把「本次同步选中的模块」用例内容增量翻译到目标语种，并翻译一级/二级菜单名
// modules 为本次选中的叶子模块清单（收窄范围，避免翻译 i18n_case/ 下所有历史模块目录）
// 模块名随用例内容 FieldCaseModule 一起翻译，再回填到 i18n_caseCount/，不再单独调用大模型
func TranslateCaseI18n(userCode, platform string, aiConnect AIConnect, content string, langs, modules []string) {
	defer func() {
		if e := recover(); e != nil {
			Logger.Error("case i18n translate panic: %v\n%s", e, debug.Stack())
		}
	}()
	// 用例内容：按模块翻译（源变更集在模块内只算一次，复用到所有语种）
	for _, m := range modules {
		if err := translateCaseModuleAllLangs(filepath.Join(caseI18nDir(), m), content, langs, platform, userCode, aiConnect); err != nil {
			Logger.Error("case i18n translate aborted at %s: %v", m, err)
			return
		}
	}

	// 模块名译文回填 + 一级/二级菜单名扁平语料（按语种，排除模块名）
	exclude := make(map[string]bool, len(modules))
	for _, m := range modules {
		exclude[m] = true
	}
	for _, lang := range langs {
		// 模块名译文随内容回填（不另起 LLM 调用）
		propagateModuleNamesToCaseCount(modules, lang)

		// 一级/二级菜单名扁平语料（排除模块名，一次独立小调用）
		if err := translateFlatI18n(filepath.Join(caseCountI18nDir(), "zh-CN.yaml"),
			filepath.Join(caseCountI18nDir(), lang+".yaml"), content, lang, platform, userCode, aiConnect, exclude); err != nil {
			Logger.Error("case i18n translate aborted at caseCount: %v", err)
			return
		}
	}
	ReloadCaseI18n()
	ReloadCaseCountI18n()
}

// TranslateEntityI18n 后台把任务实体语料（i18n_task/ i18n_playbook/ i18n_data/）增量翻译到目标语种
func TranslateEntityI18n(userCode, platform string, aiConnect AIConnect, content string, langs []string) {
	defer func() {
		if e := recover(); e != nil {
			Logger.Error("entity i18n translate panic: %v\n%s", e, debug.Stack())
		}
	}()
	for _, lang := range langs {
		for _, dir := range []string{i18nTaskDir(), i18nPlaybookDir(), i18nDataDir()} {
			if err := translateFlatI18n(filepath.Join(dir, "zh-CN.yaml"),
				filepath.Join(dir, lang+".yaml"), content, lang, platform, userCode, aiConnect, nil); err != nil {
				Logger.Error("entity i18n translate aborted at %s: %v", filepath.Base(dir), err)
				return
			}
		}
	}
	ReloadEntityI18n()
}
