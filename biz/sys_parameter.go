package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"strings"
)

func IsValueInSysParameter(sourceName, targetName string) (b bool, err error) {
	var sysParameter SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", sourceName).Find(&sysParameter)

	if len(sysParameter.ValueList) == 0 {
		err1 := fmt.Errorf("未找到系统参数[%s]，请先定义，再使用", sourceName)
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

func GetRUID(keyName string) (isExist bool) {
	var sysParameter SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", "RUID").Find(&sysParameter)

	if len(sysParameter.ValueList) == 0 {
		err := fmt.Errorf("未找到系统参数[%s]，请先定义，再使用", "RUID")
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
		err := fmt.Errorf("未找到系统参数[%s]，请先定义，再使用", "Router4Add")
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
			err = fmt.Errorf("[%s]参数中未定义[%s]的值，请核对~", parameterName, keyName)
			Logger.Warning("%s", err)
		}
	} else {
		err = fmt.Errorf("系统参数中未定义[%s]参数的值，请核对~", parameterName)
		Logger.Warning("%s", err)
	}

	return
}
