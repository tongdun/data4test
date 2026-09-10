package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/GoAdminGroup/go-admin/template/types"
)

func IsValueInSysParameter(sourceName, targetName string) (b bool, err error) {
	var sysParameter SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", sourceName).Find(&sysParameter)

	if len(sysParameter.ValueList) == 0 {
		err1 := fmt.Errorf(T("error.sys_param_not_found"), sourceName)
		err = err1
		Logger.Error("%s", err)
		return
	}

	tmps := strings.Split(sysParameter.ValueList, ",")
	for _, item := range tmps {
		info := strings.TrimSpace(item)
		if info == targetName {
			b = true
			return
		}
	}

	return
}

func GetValueFromSysParameter(lang, src string) (dst string) {
	var sysParameter SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", src).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		return
	}

	valueRaw := sysParameter.ValueList

	var valueTmp string
	if len(valueRaw) > 0 {
		if strings.Contains(valueRaw, "{") && strings.Contains(valueRaw, "}") {
			valueDefine := make(map[string]string)
			tag := 0
			json.Unmarshal([]byte(valueRaw), &valueDefine)
			if v, ok := valueDefine[lang]; ok {
				valueTmp = v
			} else if v1, ok1 := valueDefine["default"]; ok1 {
				valueTmp = v1
			} else if v2, ok2 := valueDefine["ch"]; ok2 {
				valueTmp = v2
			} else if v3, ok3 := valueDefine["en"]; ok3 {
				valueTmp = v3
			}

			if len(valueTmp) > 0 {
				dst = GetOneValueFromStringList(valueTmp)
			}

			if tag == 0 {
				for _, v := range valueDefine {
					if len(v) > 0 {
						dst = GetOneValueFromStringList(v)
						break
					}
				}
			}
		} else {
			dst = GetOneValueFromStringList(valueRaw)
		}
	}

	return
}

func GetValuesFromSysParameter(lang, src string) (values []string) {
	var sysParameter SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", src).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		return
	}

	valueRaw := sysParameter.ValueList

	var valueTmp string
	if len(valueRaw) > 0 {
		if strings.Contains(valueRaw, "{") && strings.Contains(valueRaw, "}") {
			valueDefine := make(map[string]string)
			json.Unmarshal([]byte(valueRaw), &valueDefine)
			if v, ok := valueDefine[lang]; ok {
				valueTmp = v
			} else if v1, ok1 := valueDefine["default"]; ok1 {
				valueTmp = v1
			} else if v2, ok2 := valueDefine["ch"]; ok2 {
				valueTmp = v2
			} else if v3, ok3 := valueDefine["en"]; ok3 {
				valueTmp = v3
			}

			if len(valueTmp) > 0 {
				values = GetValuesFromStringList(valueTmp)
			}
		} else {
			values = GetValuesFromStringList(valueRaw)
		}
	}

	return
}

func GetRUID(keyName string) (isExist bool) {
	var sysParameter SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", "RUID").Find(&sysParameter)

	if len(sysParameter.ValueList) == 0 {
		err := fmt.Errorf(T("error.sys_param_not_found"), "RUID")
		Logger.Error("%s", err)
		return
	}

	tmps := strings.Split(sysParameter.ValueList, ",")
	for _, item := range tmps {
		info := strings.TrimSpace(item)
		if info == keyName {
			return true
		}
	}

	return false
}

func IsInRouter4Add(path string) (isIn bool) {
	var sysParameter SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", "Router4Add").Find(&sysParameter)

	if len(sysParameter.ValueList) == 0 {
		err := fmt.Errorf(T("error.sys_param_not_found"), "Router4Add")
		Logger.Error("%s", err)
		return
	}

	tmps := strings.Split(sysParameter.ValueList, ",")
	for _, item := range tmps {
		info := strings.TrimSpace(item)
		if strings.Contains(path, info) {
			return true
		}
	}

	return false
}

func GetValueFromMapDef(parameterName, keyName string) (value string, err error) {
	var sysParameter SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", parameterName).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		return
	}

	valueRaw := sysParameter.ValueList

	if len(sysParameter.ValueList) > 0 {
		valueDefine := make(map[string]string)
		errTmp := json.Unmarshal([]byte(valueRaw), &valueDefine)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			return
		}
		if v, ok := valueDefine[keyName]; ok {
			value = v
		} else {
			err = fmt.Errorf(T("error.param_value_not_defined"), parameterName, keyName)
			Logger.Warning("%s", err)
		}
	} else {
		err = fmt.Errorf(T("error.sys_param_value_not_defined"), parameterName)
		Logger.Warning("%s", err)
	}

	return
}

// GetTestCaseExportTemplates 返回导出模板名列表（下拉选项），仿 GetAiCreatePlatform
func GetTestCaseExportTemplates() (templates []types.FieldOption) {
	var sysParameter SysParameter
	parameterName := "testCaseExportTemplates"
	models.Orm.Table("sys_parameter").Where("name = ?", parameterName).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		templates = GetNoSelectOption(T("test_case.define_export_template"))
		return
	}

	templateMap := make(map[string]TestCaseExportTemplate)
	if err := json.Unmarshal([]byte(sysParameter.ValueList), &templateMap); err != nil {
		Logger.Error(T("error.param_definition_error"), parameterName, err)
		templates = GetNoSelectOption(T("test_case.define_export_template"))
		return
	}

	names := make([]string, 0, len(templateMap))
	for name := range templateMap {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		templates = append(templates, types.FieldOption{Value: name, Text: name})
	}

	if len(templates) == 0 {
		templates = GetNoSelectOption(T("test_case.define_export_template"))
	}
	return
}

// GetTestCaseExportTemplateByName 按名取单个导出模板
func GetTestCaseExportTemplateByName(name string) (template TestCaseExportTemplate, err error) {
	var sysParameter SysParameter
	parameterName := "testCaseExportTemplates"
	models.Orm.Table("sys_parameter").Where("name = ?", parameterName).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		err = fmt.Errorf(T("error.undefined_system_parameter"), parameterName)
		return
	}

	templateMap := make(map[string]TestCaseExportTemplate)
	if errTmp := json.Unmarshal([]byte(sysParameter.ValueList), &templateMap); errTmp != nil {
		err = fmt.Errorf(T("error.param_definition_error"), parameterName, errTmp)
		return
	}

	var ok bool
	template, ok = templateMap[name]
	if !ok {
		err = fmt.Errorf(T("error.template_not_found"), name)
		return
	}
	if len(template.Columns) == 0 {
		err = fmt.Errorf(T("test_case.template_no_columns"), name)
	}
	return
}

// GetSupportLanguages 返回导出语种选项（value 为存入字段 JSON 的语种 key）
func GetSupportLanguages() (languages []types.FieldOption) {
	var sysParameter SysParameter
	parameterName := "supportLanguages"
	models.Orm.Table("sys_parameter").Where("name = ?", parameterName).Find(&sysParameter)

	langMap := make(map[string]string)
	if len(sysParameter.ValueList) > 0 {
		if err := json.Unmarshal([]byte(sysParameter.ValueList), &langMap); err != nil {
			Logger.Error(T("error.param_definition_error"), parameterName, err)
		}
	}
	if len(langMap) == 0 {
		langMap = map[string]string{"zh-CN": "中文", "en-US": "English"}
	}

	langs := make([]string, 0, len(langMap))
	for lang := range langMap {
		langs = append(langs, lang)
	}
	sort.Strings(langs)
	for _, lang := range langs {
		languages = append(languages, types.FieldOption{Value: lang, Text: langMap[lang]})
	}
	return
}
