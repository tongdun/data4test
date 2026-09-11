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

// SyncCaseI18nFromStatistics 将统计定义叶子模块内的全部用例原始数据导出为 zh-CN.yaml
// 返回（模块数, 用例数, error）；仅重写 zh-CN.yaml，保留 en-US.yaml 等其他语种文件
func SyncCaseI18nFromStatistics(caseStatId int) (int, int, error) {
	var cs CaseStatistics
	models.Orm.Table("case_statistics").Where("id = ?", caseStatId).Find(&cs)
	if cs.Id == 0 {
		return 0, 0, E("case_statistics.not_found")
	}

	def, err := ParseCaseStatisticsDefinition(cs.Definition)
	if err != nil {
		return 0, 0, E("case_statistics.definition_invalid")
	}

	modules := collectLeafModules(def)
	if len(modules) == 0 {
		return 0, 0, E("case_statistics.no_module")
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

	ReloadCaseI18n()
	return len(modules), totalCases, nil
}
