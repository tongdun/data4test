package biz

import (
	"data4perf/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
	//xlsx3 "github.com/tealeg/xlsx/v3"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func (assert SceneAssert) GetAssertValue(lang string) (out string) {
	rawStr := Interface2Str(assert.Value)
	strByte := []byte(rawStr)

	comReg := regexp.MustCompile(`\{(.+)\}`) // `\{([A-Z][A-Za-z]+)\}`匹配不上，先用当前方案
	comMatch := comReg.FindAllSubmatch(strByte, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
			name := string(comMatch[i][1])
			rawStrDef := string(comMatch[i][0])
			var assertValueDefine AssertValueDefine
			models.Orm.Table("assert_template").Where("name = ?", name).Find(&assertValueDefine)
			valueRaw := assertValueDefine.Value
			if len(valueRaw) == 0 {
				err := fmt.Errorf("未关联到断言值定义:%s", name)
				Logger.Warning("%s", err)
				//return
			}

			if len(valueRaw) > 0 {
				if !(strings.Contains(valueRaw, "{") && strings.Contains(valueRaw, "}")) {
					out = strings.Replace(rawStr, rawStrDef, valueRaw, -1)
					continue
				}

				valueDefine := make(map[string]string)
				json.Unmarshal([]byte(valueRaw), &valueDefine)
				if value, ok := valueDefine[lang]; ok {
					out = strings.Replace(rawStr, rawStrDef, value, -1)
				} else if value1, ok1 := valueDefine["default"]; ok1 {
					out = strings.Replace(rawStr, rawStrDef, value1, -1)
				} else if value2, ok2 := valueDefine["ch"]; ok2 {
					out = strings.Replace(rawStr, rawStrDef, value2, -1)
				} else if value3, ok3 := valueDefine["en"]; ok3 {
					out = strings.Replace(rawStr, rawStrDef, value3, -1)
				}

				if len(out) == 0 {
					for _, v := range valueDefine {
						if len(v) > 0 {
							out = strings.Replace(rawStr, rawStrDef, v, -1)
							break
						}
					}
				}

			} else {
				err := fmt.Errorf("未关联到断言值定义 :%s", name)
				Logger.Warning("%s", err)
			}

		}
	}

	if len(out) == 0 {
		out = rawStr
	}

	return
}

//func (sceneAssert SceneAssert) GetOutput(data map[string]interface{}) (keyName string, values []interface{}, err error) {
//	// 如果返回的数据为空，则直接返回
//	if len(data) == 0 {
//		err = fmt.Errorf("无返回信息，无法解析输出参数")
//		Logger.Error("%s", err)
//		return
//	}
//	var tmpInterface interface{}
//
//	targetValueFlowStr := Interface2Str(sceneAssert.Value)
//	if sceneAssert.Type == "output" {
//		keyName = targetValueFlowStr
//	} else {
//		keyName = fmt.Sprintf("flowVar_%s", targetValueFlowStr)
//	}
//
//	// 解析定义的校验参数
//	is4Split := Is2Split(sceneAssert.Source)
//	if strings.Contains(sceneAssert.Source, "-") && is4Split { // 字典 key 分隔
//		items := strings.SplitN(sceneAssert.Source, "-", 2)
//		sceneAssert.Source = items[1]
//		if data[items[0]] == nil {
//			err = fmt.Errorf("未解析到[%v]的值，请核对", items[0])
//			Logger.Error("%s", err)
//			return
//		} else {
//			varType := fmt.Sprintf("%T", data[items[0]])
//			if varType == "string" {
//				tmpStr := data[items[0]].(string)
//				if strings.Contains(tmpStr, "{") {
//					var tmpMap map[string]interface{}
//					json.Unmarshal([]byte(tmpStr), &tmpMap)
//					tmpInterface = tmpMap
//				}
//			} else if varType == "[]interface {}" {
//				getInterface := data[items[0]].([]interface{})
//				if len(getInterface) > 0 {
//					tmpInterface = data[items[0]].([]interface{})[0]
//				} else {
//					err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//					Logger.Error("%s", err)
//					return
//				}
//			} else if varType != "map[string]interface {}" {
//				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//				Logger.Error("%s", err)
//				return
//			}
//			tmpInterface = data[items[0]].(interface{})
//		}
//
//		varType := fmt.Sprintf("%T", tmpInterface)
//		if varType == "string" {
//			tmpStr := tmpInterface.(string)
//			if strings.Contains(tmpStr, "{") {
//				var tmpMap map[string]interface{}
//				json.Unmarshal([]byte(tmpStr), &tmpMap)
//				return sceneAssert.GetOutput(tmpMap)
//			}
//		} else if varType == "map[string]interface {}" {
//			// 进入下一层数据解析，递归调用
//			return sceneAssert.GetOutput(tmpInterface.(map[string]interface{}))
//		} else if varType == "[]interface {}" {
//			// 进入下一层数据解析，递归调用
//			getInterface := tmpInterface.([]interface{})
//			if len(getInterface) > 0 {
//				newData := tmpInterface.([]interface{})[0]
//				return sceneAssert.GetOutput(newData.(map[string]interface{}))
//			} else {
//				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//				Logger.Error("%s", err)
//				return
//			}
//
//		} else {
//			err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//			Logger.Debug("varType: %v", varType)
//			Logger.Error("%s", err)
//			return
//		}
//
//	} else if strings.Contains(sceneAssert.Source, "**") {
//		items := strings.SplitN(sceneAssert.Source, "**", 2)
//		sceneAssert.Source = items[1]
//		if data[items[0]] != nil {
//			varType := fmt.Sprintf("%T", data[items[0]])
//			if varType != "[]interface {}" {
//				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//				Logger.Error("%s", err)
//				return
//			}
//			tmpInterface = data[items[0]].([]interface{})[0]
//			varSubType := fmt.Sprintf("%T", tmpInterface)
//			if varSubType != "[]interface {}" {
//				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//				Logger.Error("%s", err)
//				return
//			}
//			newInterface := tmpInterface.([]interface{})[0]
//			return sceneAssert.GetOutput(newInterface.(map[string]interface{}))
//		}
//	} else if strings.Contains(sceneAssert.Source, "*") { // 对数组内的字典进行 key 分隔
//		items := strings.SplitN(sceneAssert.Source, "*", 2)
//		sceneAssert.Source = items[1]
//		if data[items[0]] != nil {
//			if strings.Contains(items[1], "*") {
//				isHit, index, keyName := GetSlicesIndex(items[0])
//				Logger.Debug("items[1]: %v", items[1])
//				Logger.Debug("items[0]: %v", items[0])
//				Logger.Debug("isHit: %v", isHit)
//				Logger.Debug("index: %v", index)
//				Logger.Debug("keyName: %v", keyName)
//				if !isHit {
//					keyName = items[0]
//				}
//				Logger.Debug("keyName: %v", keyName)
//				varType := fmt.Sprintf("%T", data[keyName])
//				if varType != "[]interface {}" {
//					err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//					Logger.Error("%s", err)
//					return keyName, values, err
//				}
//
//				if len(data[keyName].([]interface{})) > index {
//					if index < 0 {
//						tmpInterface = data[keyName].([]interface{})[len(data[keyName].([]interface{}))+index]
//					} else {
//						tmpInterface = data[keyName].([]interface{})[index]
//					}
//				} else {
//					Logger.Warning("索引:%d超出数据范围，自动取第0个数据", index)
//					tmpInterface = data[keyName].([]interface{})[0]
//				}
//				Logger.Debug("tmpInterface: %v", tmpInterface)
//				Logger.Debug("tmpInterface: %T", tmpInterface)
//				return sceneAssert.GetOutput(tmpInterface.(map[string]interface{}))
//			}
//
//			isHit, index, keyName := GetSlicesIndex(items[1])
//			if !isHit {
//				if strings.Contains(keyName, "-") {
//					return sceneAssert.GetOutput(data[items[0]].(map[string]interface{}))
//				} else {
//					keyName = items[1]
//				}
//			}
//
//			if isHit {
//				dataLen := len(data[items[0]].([]interface{}))
//				if dataLen == 0 {
//					Logger.Warning("返回的数据列表为空，请核对")
//				} else {
//					var tmpDict map[string]interface{}
//					if dataLen > index {
//						if index < 0 {
//							newIndex := dataLen + index
//							if newIndex < 0 {
//								err = fmt.Errorf("提供的索引超过数据范围了，数据长度: %v, 索引: %v", dataLen, index)
//								Logger.Error("%s", err)
//								return keyName, values, err
//							}
//							tmpDict = data[items[0]].([]interface{})[dataLen+index].(map[string]interface{})
//						} else {
//							tmpDict = data[items[0]].([]interface{})[index].(map[string]interface{})
//						}
//
//						if strings.Contains(items[1], "-") {
//							subItems := strings.Split(items[1], "-")
//							sceneAssert.Source = subItems[1]
//							varType := fmt.Sprintf("%T", tmpDict[keyName])
//							if varType != "map[string]interface {}" {
//								err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//								return keyName, values, err
//							}
//							subDict := tmpDict[keyName].(map[string]interface{})
//							return sceneAssert.GetOutput(subDict)
//						}
//						values = append(values, tmpDict[keyName])
//					} else {
//						Logger.Warning("索引:%d超出数据范围，自动取第0个数据", index)
//						tmpDict = data[items[0]].([]interface{})[0].(map[string]interface{})
//						values = append(values, tmpDict[keyName])
//					}
//				}
//			} else {
//				var tmpDict map[string]interface{}
//				varType := fmt.Sprintf("%T", data[items[0]])
//				if varType != "[]interface {}" {
//					err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//					return keyName, values, err
//				}
//
//				// 先判断数组是否为空，若为空，不再往后走
//				valueList := data[items[0]].([]interface{})
//				if len(valueList) == 0 {
//					err = fmt.Errorf("未获取到[%v]即[%v]的值，请核对~", sceneAssert.Source, sceneAssert.Value)
//					Logger.Error("%s", err)
//					return keyName, values, err
//				}
//
//				for _, tmpInfo := range valueList {
//					varType := fmt.Sprintf("%T", tmpInfo)
//
//					if varType != "map[string]interface {}" {
//						err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
//						Logger.Error("%s", err)
//						return keyName, values, err
//					}
//					if strings.Contains(keyName, "-") {
//						sceneAssert.Source = keyName
//						return sceneAssert.GetOutput(tmpInfo.(map[string]interface{}))
//					}
//
//					tmpDict = tmpInfo.(map[string]interface{})
//					values = append(values, tmpDict[items[1]])
//				}
//			}
//		}
//	} else { // 获取最终的值
//		if value, ok := data[sceneAssert.Source]; ok {
//			subType := fmt.Sprintf("%T", data[sceneAssert.Source])
//			if subType == "[]interface {}" {
//				values = data[sceneAssert.Source].([]interface{})
//			} else {
//				strValue := Interface2Str(value)
//				values = append(values, strValue)
//			}
//		} else {
//			Logger.Debug("sceneAssert.Source: %v", sceneAssert.Source)
//			Logger.Debug("data: %+v", data)
//			err1 := fmt.Errorf("未获取到字段[%s]即[%v]的值, 请核对~", sceneAssert.Source, sceneAssert.Value)
//			err = err1
//			Logger.Error("%s", err)
//			return
//		}
//	}
//
//	return
//}

func (sceneAssert SceneAssert) GetOutput(data interface{}) (keyName string, values []interface{}, err error) {
	// 如果返回的数据为空，则直接返回
	if data == nil {
		err = fmt.Errorf("无返回信息，无法解析输出参数")
		Logger.Error("%s", err)
		return
	}
	var tmpInterface interface{}

	targetValueFlowStr := Interface2Str(sceneAssert.Value)
	if sceneAssert.Type == "output" {
		keyName = targetValueFlowStr
	} else {
		keyName = fmt.Sprintf("flowVar_%s", targetValueFlowStr)
	}

	// 解析定义的校验参数
	splitIndexType, splitTag := GetSplitIndexType(sceneAssert.Source)
	var keyRawName string
	if len(splitTag) > 0 {
		items := strings.SplitN(sceneAssert.Source, splitTag, 2)
		keyRawName = items[0]
		sceneAssert.Source = items[1]
	} else {
		keyRawName = sceneAssert.Source
		sceneAssert.Source = ""
	}

	if splitIndexType == "list" {
		keyTmpName, index := GetSlicesIndex(keyRawName)
		tmpInterface = data.(map[string]interface{})[keyTmpName]
		listInterface := tmpInterface.([]interface{})
		var targetInterface interface{}
		varType := fmt.Sprintf("%T", tmpInterface)
		if varType != "[]interface {}" {
			err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
			Logger.Error("%s", err)
			return keyName, values, err
		}

		if len(listInterface) > index {
			if index < 0 {
				targetInterface = listInterface[len(listInterface)+index]
			} else {
				targetInterface = listInterface[index]
			}
		} else {
			Logger.Warning("索引:%d超出数据范围，自动取第0个数据", index)
			targetInterface = listInterface[0]
		}
		return sceneAssert.GetOutput(targetInterface)
	} else if splitIndexType == "listAll" {
		items := strings.SplitN(keyRawName, "[", 2)
		keyTmpName := items[0]
		sceneAssert.Source = ""
		tmpInterface = data.(map[string]interface{})[keyTmpName]
		listInterface := tmpInterface.([]interface{})
		varType := fmt.Sprintf("%T", tmpInterface)
		if varType != "[]interface {}" {
			err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
			Logger.Error("%s", err)
			return keyName, values, err
		}
		return keyName, listInterface, err
	} else {
		varType := fmt.Sprintf("%T", data)
		if len(sceneAssert.Source) == 0 {
			if varType == "[]interface {}" {
				return keyName, data.([]interface{}), err
			}

			if value, ok := data.(map[string]interface{})[keyRawName]; ok {
				subType := fmt.Sprintf("%T", data.(map[string]interface{})[keyRawName])
				if subType == "[]interface {}" {
					values = data.(map[string]interface{})[keyRawName].([]interface{})
				} else {
					strValue := Interface2Str(value)
					values = append(values, strValue)
				}
			} else {
				err1 := fmt.Errorf("未获取到字段[%s]即[%v]的值, 请核对~", keyRawName, sceneAssert.Value)
				Logger.Debug("response: %v", data)
				err = err1
				Logger.Error("%s", err)
				return
			}
		} else {
			if varType == "[]interface {}" {
				data = data.([]interface{})[0] //如果是数据，未定义索引，自动取第0个值
			}

			tmpType := fmt.Sprintf("%T", data)
			if varType == "string" { // 若取值遇到字符串，是标准的JSON格式，进行再次序列号
				tmpStr := data.(string)
				listIndex := strings.Index(tmpStr, "[")
				boundIndex := strings.Index(tmpStr, "{")
				if boundIndex < listIndex && boundIndex != -1 {
					var tmpMap map[string]interface{}
					json.Unmarshal([]byte(tmpStr), &tmpMap)
					return sceneAssert.GetOutput(tmpMap[keyRawName])
				} else if listIndex < boundIndex && listIndex != -1 {
					keyTmpName, index := GetSlicesIndex(keyRawName)
					var tmpMap []interface{}
					json.Unmarshal([]byte(tmpStr), &tmpMap)
					return sceneAssert.GetOutput(tmpMap[index].(map[string]interface{})[keyTmpName])
				} else if boundIndex != -1 {
					var tmpMap map[string]interface{}
					json.Unmarshal([]byte(tmpStr), &tmpMap)
					return sceneAssert.GetOutput(tmpMap[keyRawName])
				} else if listIndex != -1 {
					keyTmpName, index := GetSlicesIndex(keyRawName)
					var tmpMap []interface{}
					json.Unmarshal([]byte(tmpStr), &tmpMap)
					return sceneAssert.GetOutput(tmpMap[index].(map[string]interface{})[keyTmpName])
				}
			} else if tmpType != "map[string]interface {}" {
				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
				Logger.Error("%s", err)
				return
			}

			if value, ok := data.(map[string]interface{})[keyRawName]; ok {
				return sceneAssert.GetOutput(value)
			} else {
				err1 := fmt.Errorf("未获取到字段[%s]即[%v]的值, 请核对~", keyRawName, sceneAssert.Value)
				err = err1
				Logger.Error("%s", err)
				return
			}
		}

	}

	return
}

func (sceneAssert SceneAssert) GetOutputRe(raw []byte) (keyName string, values []interface{}, err error) {
	// 如果校验类型不是 ouput ，则直接返回
	if sceneAssert.Type != "output_re" {
		return
	}

	// 如果返回的数据为空，则直接返回
	if len(raw) == 0 {
		err = fmt.Errorf("无返回信息，无法解析输出参数")
		Logger.Error("%s", err)
		return
	}

	keyName = Interface2Str(sceneAssert.Value)

	// 解析定义的校验参数
	sourceStr := Interface2Str(sceneAssert.Source)
	comReg := regexp.MustCompile(sourceStr)
	comMatch := comReg.FindAllSubmatch(raw, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
			if len(comMatch[i]) <= 1 {
				err = fmt.Errorf("source: %s, value: %s, 正则取值定义有误，请核对~", sourceStr, keyName)
				Logger.Error("%s", err)
				return
			}
			value := string(comMatch[i][1])
			values = append(values, value)
		}
	} else {
		err = fmt.Errorf("source: %s, value: %s, 正则取值定义与实际返回信息未匹配上，请核对~", sourceStr, keyName)
		Logger.Error("%s", err)
		return
	}

	return
}

func (sceneAssert SceneAssert) AssertResult(data map[string]interface{}, inOutPutDict map[string][]interface{}) (b bool, err error) {
	_, flowValues, err := sceneAssert.GetOutput(data)
	if err != nil {
		return
	}

	targetValueStr := Interface2Str(sceneAssert.Value)
	targetValueStr, notDefVars, falseCount, errTmp := GetIndexStr("", targetValueStr, "", "", inOutPutDict)
	if falseCount > 0 {
		if errTmp != nil {
			err = fmt.Errorf("%s; 存在未定义参数: %s，请先定义或关联", errTmp, notDefVars)
		} else {
			err = fmt.Errorf("存在未定义参数: %s，请先定义或关联", notDefVars)
		}
		Logger.Error("%s", err)
		return
	}

	for _, subV := range flowValues {
		expectValue := Interface2Str(subV)
		errTmp := sceneAssert.AsserValueComparion(expectValue)
		if errTmp != nil {
			if err == nil {
				err = errTmp
			} else {
				err = fmt.Errorf("%s\n%s", err, errTmp)
			}
		}
	}

	return
}

func GetAssertTemplateList() (depAssertTmpList []DepAssertModel) {
	var values []string
	models.Orm.Table("assert_template").Order("created_at desc").Pluck("name", &values)
	for _, item := range values {
		var depAssert DepAssertModel
		depAssert.Name = fmt.Sprintf("{%s}", item)
		depAssertTmpList = append(depAssertTmpList, depAssert)
	}
	return
}

func GetAssertTemplateValue(lang, key string) (value string, err error) {
	var assertValueDefine AssertValueDefine
	models.Orm.Table("assert_template").Where("name = ?", key).Find(&assertValueDefine)
	valueRaw := assertValueDefine.Value
	if len(valueRaw) == 0 {
		errTmp := fmt.Errorf("未关联到断言值定义:%s", key)
		Logger.Warning("%s", errTmp)
		err = errTmp
		return
	}

	if len(valueRaw) > 0 {
		valueDefine := make(map[string]string)
		err = json.Unmarshal([]byte(valueRaw), &valueDefine)
		if err != nil {
			value = valueRaw
			return
		}

		if v, ok := valueDefine[lang]; ok {
			value = v
		} else if v1, ok1 := valueDefine["default"]; ok1 {
			value = v1
		} else if v2, ok2 := valueDefine["ch"]; ok2 {
			value = v2
		} else if v3, ok3 := valueDefine["en"]; ok3 {
			value = v3
		}

		if len(value) == 0 {
			for _, v := range valueDefine {
				value = v
				break
			}
		}
	}

	return
}

func GetAssertTemplateAllValue(lang string) (allValue map[string][]interface{}) {
	var assertList []AssertValueDefine
	models.Orm.Table("assert_template").Find(&assertList)

	if len(assertList) == 0 {
		return
	}

	for _, item := range assertList {
		valueDefine := make(map[string]string)
		err := json.Unmarshal([]byte(item.Value), &valueDefine)
		if err != nil {
			allValue[item.Name] = append(allValue[item.Name], item.Value)
			continue
		}

		if v, ok := valueDefine[lang]; ok {
			allValue[item.Name] = append(allValue[item.Name], v)
		} else if v1, ok1 := valueDefine["default"]; ok1 {
			allValue[item.Name] = append(allValue[item.Name], v1)
		} else if v2, ok2 := valueDefine["ch"]; ok2 {
			allValue[item.Name] = append(allValue[item.Name], v2)
		} else if v3, ok3 := valueDefine["en"]; ok3 {
			allValue[item.Name] = append(allValue[item.Name], v3)
		}

		if len(allValue[item.Name]) == 0 {
			for _, v := range valueDefine {
				allValue[item.Name] = append(allValue[item.Name], v)
				break
			}
		}
	}

	return
}

func (assert SceneAssert) AsserValueComparion(curStr string) (err error) {
	var b bool
	var rawTargetStr string
	targetStr := Interface2Str(assert.Value)
	switch assert.Type {
	case "=", "equal", "!=", "not_equal", ">", "larger_than", "greater_than", ">=", "larger_equal", "greater_equal", "<", "less_than", "<=", "less_equal":
		expression, errTmp := govaluate.NewEvaluableExpression(targetStr)
		if errTmp == nil {
			parameters := make(map[string]interface{})
			newTarget, errTmp := expression.Evaluate(parameters)
			if errTmp == nil {
				rawTargetStr = targetStr
				targetStr = Interface2Str(newTarget)

			}
		}
	}

	switch assert.Type {
	case "=", "equal":
		if curStr == targetStr {
			b = true
		} else {
			b = false
		}
	case "!=", "not_equal":
		if curStr != targetStr {
			b = true
		} else {
			b = false
		}
	case "in", "contain":
		if strings.Contains(curStr, targetStr) {
			b = true
		} else {
			b = false
		}
	case "!in", "not_in", "not_contain":
		if !strings.Contains(curStr, targetStr) {
			b = true
		} else {
			b = false
		}
	case "re", "regex", "regexp":
		re := regexp.MustCompile(targetStr)
		result := re.FindStringSubmatch(curStr)
		if len(result) > 0 {
			b = true
		} else {
			b = false
		}
	case "null", "empty":
		if len(curStr) == 0 {
			b = true
		} else {
			b = false
		}
	case "!null", "!empty", "not_null", "not_empty":
		if len(curStr) > 0 {
			b = true
		} else {
			b = false
		}
	case ">", "larger_than", "greater_than":
		targetInt, err1 := strconv.Atoi(targetStr)
		curInt, err2 := strconv.Atoi(curStr)
		if err1 != nil || err2 != nil {
			b = false
		} else {
			if curInt > targetInt {
				b = true
			} else {
				b = false
			}
		}
	case ">=", "larger_equal", "greater_equal":
		targetInt, err1 := strconv.Atoi(targetStr)
		curInt, err2 := strconv.Atoi(curStr)
		if err1 != nil || err2 != nil {
			b = false
		} else {
			if curInt >= targetInt {
				b = true
			} else {
				b = false
			}
		}
	case "<", "less_than":
		targetInt, err1 := strconv.Atoi(targetStr)
		curInt, err2 := strconv.Atoi(curStr)
		if err1 != nil || err2 != nil {
			b = false
		} else {
			if curInt < targetInt {
				b = true
			} else {
				b = false
			}
		}
	case "<=", "less_equal":
		targetInt, err1 := strconv.Atoi(targetStr)
		curInt, err2 := strconv.Atoi(curStr)
		if err1 != nil || err2 != nil {
			b = false
		} else {
			if curInt <= targetInt {
				b = true
			} else {
				b = false
			}
		}
	default:
		err = fmt.Errorf("不支持%s类型的比较，如有需要请反馈致相关人员", assert.Type)
	}

	if !b {
		var expectPrompt string
		switch assert.Source {
		case "raw", "ResponseBody", "Response":
			if len(rawTargetStr) > 0 {
				expectPrompt = fmt.Sprintf("预期: ResponseBody %s %s", assert.Type, rawTargetStr)
			} else {
				expectPrompt = fmt.Sprintf("预期: ResponseBody %s %s", assert.Type, targetStr)
			}
			actualPrompt := fmt.Sprintf("实际: ResponseBody = %s", curStr)
			err = fmt.Errorf("%s\n%s\n断言: ResponseBody %s %s 结果:fail", expectPrompt, actualPrompt, assert.Type, targetStr)
		default:
			if len(rawTargetStr) > 0 {
				expectPrompt = fmt.Sprintf("预期: %s %s %s", assert.Source, assert.Type, rawTargetStr)
			} else {
				expectPrompt = fmt.Sprintf("预期: %s %s %s", assert.Source, assert.Type, targetStr)
			}
			actualPrompt := fmt.Sprintf("实际: %s = %s", assert.Source, curStr)
			err = fmt.Errorf("%s\n%s\n断言: %s %s %s 结果:fail", expectPrompt, actualPrompt, curStr, assert.Type, targetStr)
		}
	}

	return
}

func (assert SceneAssert) GetValueFromFile(fileName string) (targetList []string, err error) {
	filePath := fmt.Sprintf("%s/%s", DownloadBasePath, fileName) // 下载文件统一保存至下载文件路径
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		err = fmt.Errorf("%s不存在，请核对", filePath)
		return
	}

	dataAnchor := strings.Split(assert.Source, ":")
	if len(dataAnchor) <= 1 {
		err = fmt.Errorf("source: %s, 未定义正常的文件取数信息，请核对", assert.Source)
		return
	}

	fileType := strings.ToUpper(dataAnchor[1])
	switch fileType {
	case "CSV":
		targetList, err = GetTargetValueFromStructFile("csv", assert.Source, filePath)
	case "EXCEL":
		targetList, err = GetTargetValueFromStructFile("excel", assert.Source, filePath)
	//case "TXT":  // 待实现
	case "JSON":
		targetList, err = assert.GetTargetValueFromNoStructFile("json", filePath)
	case "YML", "YAML":
		targetList, err = assert.GetTargetValueFromNoStructFile("yml", filePath)
	//case "XML":    // 待实现
	default:
		err = fmt.Errorf("不支持[%s]文件类型的取数校对，如有需要，请联系管理员~", fileType)
		return
	}

	return
}

func GetTargetValueFromCSV(filePath, columnName, splitTag string, lineNo, columnNo int) (target []string, err error) {
	var tagRune rune

	runeList := []rune(splitTag)

	for _, item := range runeList {
		tagRune = tagRune + item
	}

	fh, errTmp := os.Open(filePath)
	if errTmp != nil {
		err = errTmp
		return
	}

	defer fh.Close()

	reader := csv.NewReader(fh)
	//unicode.UTF8BOM.NewDecoder()

	if len(splitTag) > 0 {
		reader.Comma = tagRune
	} else {
		reader.Comma = ','
	}
	reader.LazyQuotes = true
	//reader.FieldsPerRecord = -1

	curLine := 0
	if lineNo > 0 {
		for {
			rawRecord, errTmp := reader.Read()
			if errTmp == io.EOF {
				break
			}
			curLine++
			var record []string
			for _, item := range rawRecord {
				if strings.Contains(item, "\"") {
					newItem := strings.Trim(item, "\"")
					record = append(record, newItem)
				} else {
					record = append(record, item)
				}
			}

			if len(columnName) > 0 && curLine == 1 {
				for index, item := range record {
					if columnName == item {
						columnNo = index
					}
				}
			}

			if curLine == lineNo {
				if len(record) < columnNo && len(columnName) == 0 {
					err = fmt.Errorf("列号: %d超出索引范围，请核对", columnNo)
					return
				}
				if columnNo > 0 {
					target = []string{record[columnNo-1]}
				} else if columnNo == -1 {
					target = record
				}
				break
			}
		}
	} else if lineNo == -1 {
		for {
			record, errTmp := reader.Read()
			if errTmp == io.EOF {
				break
			}

			if len(columnName) > 0 && curLine == 1 {
				for index, item := range record {
					if columnName == item {
						columnNo = index
					}
				}
				continue
			}

			if len(record) < columnNo {
				err = fmt.Errorf("列号: %d超出索引范围，请核对")
				return
			}

			for index, item := range record {
				if index == columnNo {
					target = append(target, item)
				}
			}
		}
	}

	return
}

func GetTargetValueFromXLSX(filePath, columnName string, lineNo, columnNo int) (target []string, err error) {
	fh, err := xlsx.OpenFile(filePath)
	if err != nil {
		Logger.Debug("filePath: %s", filePath)
		Logger.Error("%v", err)
		return
	}

	sheet := fh.Sheets[0]

	if len(columnName) > 0 {
		titles := sheet.Row(1).Cells
		for index, item := range titles {
			if item.Value == columnName {
				columnNo = index
			}
		}
	}

	maxRowNo := sheet.MaxRow
	maxColNo := sheet.MaxCol

	if lineNo > maxRowNo {
		err = fmt.Errorf("列号: %d超出索引范围，请核对", lineNo)
		return
	}

	if columnNo > maxColNo {
		err = fmt.Errorf("行号: %d超出索引范围，请核对", columnNo)
		return
	}

	if lineNo > 0 {
		if columnNo > 0 {
			targetSingle := sheet.Cell(lineNo-1, columnNo-1)
			target = []string{targetSingle.String()}
		} else {
			for i := 0; i < maxColNo; i++ {
				targetSingle := sheet.Cell(lineNo-1, i)
				target = append(target, targetSingle.String())
			}
		}
	} else if lineNo == -1 {
		for i := 1; i < maxRowNo; i++ {
			targetSingle := sheet.Cell(i, columnNo-1)
			target = append(target, targetSingle.String())
		}
	}
	return
}

func GetTargetValueFromXLS(filePath, columnName string, lineNo, columnNo int) (target []string, err error) {
	fh, err := xls.Open(filePath, "utf-8")
	if err != nil {
		Logger.Debug("filePath: %s", filePath)
		Logger.Error("%v", err)
		return
	}

	sheet := fh.GetSheet(0)
	titles := fh.ReadAllCells(1)[0]
	if len(columnName) > 0 {
		for index, item := range titles {
			if item == columnName {
				columnNo = index
			}
		}
	}

	maxRowNo := int(sheet.MaxRow)
	maxColNo := len(titles)
	if lineNo > maxRowNo {
		err = fmt.Errorf("列号: %d超出索引范围，请核对", lineNo)
		return
	}

	if columnNo > maxColNo {
		err = fmt.Errorf("行号: %d超出索引范围，请核对", columnNo)
		return
	}

	if lineNo > 0 {
		if columnNo > 0 {
			targetSingle := sheet.Row(lineNo - 1).Col(columnNo - 1)
			target = []string{targetSingle}
		} else {
			for i := 0; i < maxColNo; i++ {
				targetSingle := sheet.Row(lineNo - 1).Col(i)
				target = append(target, targetSingle)
			}
		}
	} else if lineNo == -1 {
		for i := 1; i < maxRowNo; i++ {
			targetSingle := sheet.Row(i).Col(columnNo - 1)
			target = append(target, targetSingle)
		}
	}
	return
}

func GetTargetValueFromStructFile(fileType, source, filePath string) (target []string, err error) {
	var splitTag, columnName string
	var lineNo, columnNo int
	var errTmp error

	dataAnchor := strings.Split(source, ":")

	if len(dataAnchor) < 4 {
		err = fmt.Errorf("source: %s, 取数信息定义不全，请核对", source)
		return
	}

	if len(dataAnchor[2]) == 0 {
		lineNo = -1
	} else {
		lineNo, errTmp = strconv.Atoi(dataAnchor[2])
		if errTmp != nil {
			Logger.Error("%v", errTmp)
			err = fmt.Errorf("行号: %v, 无法转换为整数，请核对", dataAnchor[2])
			return
		}
	}

	if len(dataAnchor[3]) == 0 {
		columnNo = -1
	} else {
		columnNo, errTmp = strconv.Atoi(dataAnchor[3])
		if errTmp != nil {
			Logger.Warning("%s", errTmp)
			columnName = dataAnchor[3]
		}
	}

	if fileType == "excel" {
		//target, err = GetTargetValueFromEXCEL(filePath, columnName, lineNo, columnNo)
		if strings.HasSuffix(filePath, ".xlsx") {
			target, err = GetTargetValueFromXLSX(filePath, columnName, lineNo, columnNo)
		} else if strings.HasSuffix(filePath, ".xls") {
			target, err = GetTargetValueFromXLS(filePath, columnName, lineNo, columnNo)
		} else {
			Logger.Warning("%s文件格式不支持，请核对", filePath)
			return
		}
	} else if fileType == "csv" {
		if len(dataAnchor) >= 5 {
			splitTag = dataAnchor[4]
		} else {
			splitTag = ","
		}
		target, err = GetTargetValueFromCSV(filePath, columnName, splitTag, lineNo, columnNo)
	}

	return
}

func (assert SceneAssert) GetTargetValueFromNoStructFile(fileType, filePath string) (target []string, err error) {
	dataAnchor := strings.Split(assert.Source, ":")
	if len(dataAnchor) < 3 {
		err = fmt.Errorf("source: %s, 取数信息定义不全，请核对", assert.Source)
		return
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	var resDict map[string]interface{}
	resDict = make(map[string]interface{})

	if fileType == "json" {
		err = json.Unmarshal([]byte(content), &resDict)
	} else if fileType == "yml" {
		err = yaml.Unmarshal([]byte(content), &resDict)
	}

	if err != nil {
		Logger.Error("err: %v", err)
		return
	}

	var valueList []interface{}

	assert.Source = dataAnchor[2]
	switch assert.Type {
	case "output_re":
		_, valueList, err = assert.GetOutputRe(content)
	default:
		_, valueList, err = assert.GetOutput(resDict)
	}

	for _, item := range valueList {
		vStr := Interface2Str(item)
		target = append(target, vStr)
	}

	return
}
