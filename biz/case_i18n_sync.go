package biz

import (
	"data4test/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// caseI18nRow 同步导出用的用例行（7 个可翻译字段 + 用例编号）
type caseI18nRow struct {
	CaseNumber   string `gorm:"column:case_number"`
	CaseName     string `gorm:"column:case_name"`
	Module       string `gorm:"column:module"`
	CaseType     string `gorm:"column:case_type"`
	PreCondition string `gorm:"column:pre_condition"`
	TestRange    string `gorm:"column:test_range"`
	TestSteps    string `gorm:"column:test_steps"`
	ExpectResult string `gorm:"column:expect_result"`
}

// collectLeafModules 收集统计定义树中所有叶子模块名（去重、去空）
func collectLeafModules(def CaseStatisticsDefinition) []string {
	seen := make(map[string]struct{})
	modules := make([]string, 0)
	add := func(m string) {
		m = strings.TrimSpace(m)
		if m == "" {
			return
		}
		if _, ok := seen[m]; ok {
			return
		}
		seen[m] = struct{}{}
		modules = append(modules, m)
	}
	for _, l1 := range def.Tree {
		for _, m := range l1.Modules {
			add(m)
		}
		for _, l2 := range l1.Children {
			for _, m := range l2.Modules {
				add(m)
			}
		}
	}
	return modules
}

// collectCaseCountNames 收集统计定义树中的一级/二级菜单名 + 三级模块名（去重、去空）
func collectCaseCountNames(def CaseStatisticsDefinition) map[string]string {
	set := make(map[string]string)
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		set[s] = s
	}
	for _, l1 := range def.Tree {
		add(l1.Name)
		for _, m := range l1.Modules {
			add(m)
		}
		for _, l2 := range l1.Children {
			add(l2.Name)
			for _, m := range l2.Modules {
				add(m)
			}
		}
	}
	return set
}

// writeCaseCountI18n 写 i18n_caseCount 目录下的 zh-CN.yaml（一级/二级菜单名 + 三级模块名身份映射骨架），保留其他语种文件
func writeCaseCountI18n(data map[string]string) error {
	dir := caseCountI18nDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		Logger.Error("case_count i18n sync mkdir %s failed: %v", dir, err)
		return err
	}
	bytes, err := yaml.Marshal(data)
	if err != nil {
		Logger.Error("case_count i18n sync marshal failed: %v", err)
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "zh-CN.yaml"), bytes, 0644); err != nil {
		Logger.Error("case_count i18n sync write zh-CN.yaml failed: %v", err)
		return err
	}
	return nil
}

// SyncCaseI18nFromStatistics 将选中统计定义的叶子模块内全部用例原始数据导出为 zh-CN.yaml
// 支持一次选中多个定义（逗号分隔 id），跨定义去重后按模块写文件
// 返回（模块数, 用例数, error）；仅重写 zh-CN.yaml，保留 en-US.yaml 等其他语种文件
func SyncCaseI18nFromStatistics(idStr string) (int, int, error) {
	ids := parseIds(idStr)
	if len(ids) == 0 {
		return 0, 0, E("common.btn_select_first")
	}

	var css []CaseStatistics
	models.Orm.Table("case_statistics").Where("id IN (?)", ids).Find(&css)
	if len(css) == 0 {
		return 0, 0, E("case_statistics.not_found")
	}

	// 汇总所有选中定义涉及的叶子模块（跨定义去重）+ 菜单树名称
	moduleSet := make(map[string]struct{})
	caseCountSet := make(map[string]string)
	for _, cs := range css {
		def, err := ParseCaseStatisticsDefinition(cs.Definition)
		if err != nil {
			Logger.Error("i18n同步解析统计定义 %d 失败: %v", cs.Id, err)
			continue
		}
		for _, m := range collectLeafModules(def) {
			moduleSet[m] = struct{}{}
		}
		for k, v := range collectCaseCountNames(def) {
			caseCountSet[k] = v
		}
	}
	if len(moduleSet) == 0 {
		return 0, 0, E("case_statistics.no_module")
	}
	modules := make([]string, 0, len(moduleSet))
	for m := range moduleSet {
		modules = append(modules, m)
	}

	dir := caseI18nDir()
	totalCases := 0
	for _, m := range modules {
		var rows []caseI18nRow
		if err := models.Orm.Table("test_case").
			Where("module = ? AND deleted_at IS NULL", m).
			Select("case_number, case_name, module, case_type, pre_condition, test_range, test_steps, expect_result").
			Find(&rows).Error; err != nil {
			Logger.Error("i18n同步查询模块 %s 用例失败: %v", m, err)
			continue
		}

		data := make(map[string]map[string]string, len(rows))
		for _, r := range rows {
			if strings.TrimSpace(r.CaseNumber) == "" {
				continue
			}
			data[r.CaseNumber] = map[string]string{
				FieldCaseName:     r.CaseName,
				FieldCaseModule:   r.Module,
				FieldCaseType:     r.CaseType,
				FieldPreCondition: r.PreCondition,
				FieldTestRange:    r.TestRange,
				FieldTestSteps:    r.TestSteps,
				FieldExpectResult: r.ExpectResult,
			}
		}

		bytes, err := yaml.Marshal(data)
		if err != nil {
			Logger.Error("i18n同步序列化模块 %s 失败: %v", m, err)
			continue
		}

		moduleDir := filepath.Join(dir, m)
		if err := os.MkdirAll(moduleDir, 0755); err != nil {
			Logger.Error("i18n同步创建目录 %s 失败: %v", moduleDir, err)
			continue
		}
		if err := ioutil.WriteFile(filepath.Join(moduleDir, "zh-CN.yaml"), bytes, 0644); err != nil {
			Logger.Error("i18n同步写模块 %s 失败: %v", m, err)
			continue
		}
		totalCases += len(data)
		Logger.Info("i18n同步导出模块 %s：%d 条用例", m, len(data))
	}

	// 一级/二级菜单名 + 三级模块名 → i18n_caseCount/zh-CN.yaml
	if err := writeCaseCountI18n(caseCountSet); err != nil {
		return 0, 0, err
	}
	ReloadCaseCountI18n()

	ReloadCaseI18n()
	return len(modules), totalCases, nil
}
