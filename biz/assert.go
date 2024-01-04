package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func (assert SceneAssert) GetAssertValue(lang string) (out string) {
	rawStr, _ := Interface2Str(assert.Value)
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
				err := fmt.Errorf("未关联到断言值定义 :%s", name)
				Logger.Warning("%s", err)
				return
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
				return
			}

		}
	} else {
		out = rawStr
	}
	return
}

func (sceneAssert SceneAssert) GetOutput(data map[string]interface{}) (outputDict map[string][]interface{}, err error) {
	// 如果校验类型不是 ouput ，则直接返回
	if sceneAssert.Type != "output" {
		return
	}

	// 如果返回的数据为空，则直接返回
	if len(data) == 0 {
		err = fmt.Errorf("无返回信息，无法解析输出参数")
		Logger.Error("%s", err)
		return
	}
	var tmpInterface interface{}
	outputDict = make(map[string][]interface{})
	targetValueStr, err := Interface2Str(sceneAssert.Value)
	if err != nil {
		return
	}
	// 解析定义的校验参数
	is4Split := Is2Split(sceneAssert.Source)
	if strings.Contains(sceneAssert.Source, "-") && is4Split { // 字典 key 分隔
		items := strings.SplitN(sceneAssert.Source, "-", 2)
		sceneAssert.Source = items[1]
		if data[items[0]] == nil {
			err = fmt.Errorf("未解析到[%v]的值，请核对", items[0])
			Logger.Error("%s", err)
			return
		} else {
			varType := fmt.Sprintf("%T", data[items[0]])
			if varType == "string" {
				tmpStr := data[items[0]].(string)
				if strings.Contains(tmpStr, "{") {
					var tmpMap map[string]interface{}
					json.Unmarshal([]byte(tmpStr), &tmpMap)
					tmpInterface = tmpMap
				}
			} else if varType == "[]interface {}" {
				tmpInterface = data[items[0]].([]interface{})[0]
			} else if varType != "map[string]interface {}" {
				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
				Logger.Error("%s", err)
				return
			}
			tmpInterface = data[items[0]].(interface{})
		}

		varType := fmt.Sprintf("%T", tmpInterface)
		if varType == "string" {
			tmpStr := tmpInterface.(string)
			if strings.Contains(tmpStr, "{") {
				var tmpMap map[string]interface{}
				json.Unmarshal([]byte(tmpStr), &tmpMap)
				return sceneAssert.GetOutput(tmpMap)
			}
		} else if varType == "map[string]interface {}" {
			// 进入下一层数据解析，递归调用
			return sceneAssert.GetOutput(tmpInterface.(map[string]interface{}))
		} else if varType == "[]interface {}" {
			// 进入下一层数据解析，递归调用
			newData := tmpInterface.([]interface{})[0]
			return sceneAssert.GetOutput(newData.(map[string]interface{}))
		} else {
			err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
			Logger.Debug("varType: %v", varType)
			Logger.Error("%s", err)
			return
		}

	} else if strings.Contains(sceneAssert.Source, "**") {
		items := strings.SplitN(sceneAssert.Source, "**", 2)
		sceneAssert.Source = items[1]
		if data[items[0]] != nil {
			varType := fmt.Sprintf("%T", data[items[0]])
			if varType != "[]interface {}" {
				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
				Logger.Error("%s", err)
				return
			}
			tmpInterface = data[items[0]].([]interface{})[0]
			varSubType := fmt.Sprintf("%T", tmpInterface)
			if varSubType != "[]interface {}" {
				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
				Logger.Error("%s", err)
				return
			}
			newInterface := tmpInterface.([]interface{})[0]
			return sceneAssert.GetOutput(newInterface.(map[string]interface{}))
		}
	} else if strings.Contains(sceneAssert.Source, "*") { // 对数组内的字典进行 key 分隔
		items := strings.SplitN(sceneAssert.Source, "*", 2)
		sceneAssert.Source = items[1]
		if data[items[0]] != nil {
			if strings.Contains(items[1], "*") {
				isHit, index, keyName := GetSlicesIndex(items[0])
				if !isHit {
					keyName = items[0]
				}
				varType := fmt.Sprintf("%T", data[keyName])
				if varType != "[]interface {}" {
					err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
					Logger.Error("%s", err)
					return
				}

				if len(data[keyName].([]interface{})) > index {
					if index < 0 {
						tmpInterface = data[keyName].([]interface{})[len(data[keyName].([]interface{}))+index]
					} else {
						tmpInterface = data[keyName].([]interface{})[index]
					}

				} else {
					Logger.Warning("索引:%d超出数据范围，自动取第0个数据", index)
					tmpInterface = data[keyName].([]interface{})[0]
				}
				return sceneAssert.GetOutput(tmpInterface.(map[string]interface{}))
			}

			isHit, index, keyName := GetSlicesIndex(items[1])
			if !isHit {
				if strings.Contains(keyName, "-") {
					return sceneAssert.GetOutput(data[items[0]].(map[string]interface{}))
				} else {
					keyName = items[1]
				}
			}
			if isHit {
				dataLen := len(data[items[0]].([]interface{}))
				if dataLen == 0 {
					Logger.Warning("返回的数据列表为空，请核对")
				} else {
					var tmpDict map[string]interface{}
					if dataLen > index {
						if index < 0 {
							newIndex := dataLen + index
							if newIndex < 0 {
								err = fmt.Errorf("提供的索引超过数据范围了，数据长度: %v, 索引: %v", dataLen, index)
								Logger.Error("%s", err)
								return
							}
							tmpDict = data[items[0]].([]interface{})[dataLen+index].(map[string]interface{})
							outputDict[targetValueStr] = append(outputDict[targetValueStr], tmpDict[keyName])
						} else {
							tmpDict = data[items[0]].([]interface{})[index].(map[string]interface{})
							outputDict[targetValueStr] = append(outputDict[targetValueStr], tmpDict[keyName])
						}
					} else {
						Logger.Warning("索引:%d超出数据范围，自动取第0个数据", index)
						tmpDict = data[items[0]].([]interface{})[0].(map[string]interface{})
						outputDict[targetValueStr] = append(outputDict[targetValueStr], tmpDict[keyName])
					}
				}
			} else {
				var tmpDict map[string]interface{}
				varType := fmt.Sprintf("%T", data[items[0]])
				if varType != "[]interface {}" {
					err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
					Logger.Error("%s", err)
					return
				}

				// 先判断数组是否为空，若为空，不再往后走
				valueList := data[items[0]].([]interface{})
				if len(valueList) == 0 {
					err = fmt.Errorf("未获取到[%v]即[%v]的值，请核对~", sceneAssert.Source, sceneAssert.Value)
					Logger.Error("%s", err)
					return
				}

				for _, tmpInfo := range valueList {
					varType := fmt.Sprintf("%T", tmpInfo)

					if varType != "map[string]interface {}" {
						err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
						Logger.Error("%s", err)
						return
					}
					if strings.Contains(keyName, "-") {

						sceneAssert.Source = keyName

						return sceneAssert.GetOutput(tmpInfo.(map[string]interface{}))
					}

					tmpDict = tmpInfo.(map[string]interface{})
					outputDict[targetValueStr] = append(outputDict[targetValueStr], tmpDict[items[1]])
				}
			}
		}
	} else { // 获取最终的值
		if value, ok := data[sceneAssert.Source]; ok {
			subType := fmt.Sprintf("%T", data[sceneAssert.Source])
			if subType == "[]interface {}" {
				outputDict[targetValueStr] = data[sceneAssert.Source].([]interface{})
			} else {
				strValue, err1 := Interface2Str(value)
				if err1 != nil {
					err = err1
					Logger.Error("%s", err)
					return
				}
				outputDict[targetValueStr] = append(outputDict[targetValueStr], strValue)
			}
		} else {
			err1 := fmt.Errorf("未获取到字段[%s]即[%v]的值, 请核对~", sceneAssert.Source, sceneAssert.Value)
			err = err1
			Logger.Error("%s", err)
			return
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

	keyName, err = Interface2Str(sceneAssert.Value)
	if err != nil {
		err = fmt.Errorf("正则取值定义未定义输出变量名，请核对~")
		Logger.Error("%s", err)
	}

	// 解析定义的校验参数
	sourceStr, err := Interface2Str(sceneAssert.Source)
	comReg := regexp.MustCompile(sourceStr)
	comMatch := comReg.FindAllSubmatch(raw, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
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

func (sceneAssert SceneAssert) SpecialAssertValue(data map[string]interface{}) (outputDict map[string][]interface{}, err error) {
	//支持断言非output的下标取值运算

	var tmpInterface interface{}
	outputDict = make(map[string][]interface{})
	targetValueStr := "_"
	if err != nil {
		return
	}
	// 解析定义的校验参数
	is4Split := Is2Split(sceneAssert.Source)
	if strings.Contains(sceneAssert.Source, "-") && is4Split { // 字典 key 分隔
		items := strings.SplitN(sceneAssert.Source, "-", 2)
		sceneAssert.Source = items[1]
		if data[items[0]] == nil {
			err = fmt.Errorf("未解析到[%v]的值，请核对", items[0])
			Logger.Error("%s", err)
			return
		} else {
			varType := fmt.Sprintf("%T", data[items[0]])
			if varType == "string" {
				tmpStr := data[items[0]].(string)
				if strings.Contains(tmpStr, "{") {
					var tmpMap map[string]interface{}
					json.Unmarshal([]byte(tmpStr), &tmpMap)
					tmpInterface = tmpMap
				}
			} else if varType == "[]interface {}" {
				tmpInterface = data[items[0]].([]interface{})[0]
			} else if varType != "map[string]interface {}" {
				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
				Logger.Error("%s", err)
				return
			}
			tmpInterface = data[items[0]].(interface{})
		}

		varType := fmt.Sprintf("%T", tmpInterface)
		if varType == "string" {
			tmpStr := tmpInterface.(string)
			if strings.Contains(tmpStr, "{") {
				var tmpMap map[string]interface{}
				json.Unmarshal([]byte(tmpStr), &tmpMap)
				return sceneAssert.SpecialAssertValue(tmpMap)
			}
		} else if varType == "map[string]interface {}" {
			// 进入下一层数据解析，递归调用
			return sceneAssert.SpecialAssertValue(tmpInterface.(map[string]interface{}))
		} else if varType == "[]interface {}" {
			// 进入下一层数据解析，递归调用
			newData := tmpInterface.([]interface{})[0]
			return sceneAssert.SpecialAssertValue(newData.(map[string]interface{}))
		} else {
			err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
			Logger.Debug("varType: %v", varType)
			Logger.Error("%s", err)
			return
		}

	} else if strings.Contains(sceneAssert.Source, "**") {
		items := strings.SplitN(sceneAssert.Source, "**", 2)
		sceneAssert.Source = items[1]
		if data[items[0]] != nil {
			varType := fmt.Sprintf("%T", data[items[0]])
			if varType != "[]interface {}" {
				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
				Logger.Error("%s", err)
				return
			}
			tmpInterface = data[items[0]].([]interface{})[0]
			varSubType := fmt.Sprintf("%T", tmpInterface)
			if varSubType != "[]interface {}" {
				err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
				Logger.Error("%s", err)
				return
			}
			newInterface := tmpInterface.([]interface{})[0]
			return sceneAssert.SpecialAssertValue(newInterface.(map[string]interface{}))
		}
	} else if strings.Contains(sceneAssert.Source, "*") { // 对数组内的字典进行 key 分隔
		items := strings.SplitN(sceneAssert.Source, "*", 2)
		sceneAssert.Source = items[1]
		if data[items[0]] != nil {
			if strings.Contains(items[1], "*") {
				isHit, index, keyName := GetSlicesIndex(items[0])
				if !isHit {
					keyName = items[0]
				}
				varType := fmt.Sprintf("%T", data[keyName])
				if varType != "[]interface {}" {
					err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
					Logger.Error("%s", err)
					return
				}

				if len(data[keyName].([]interface{})) > index {
					if index < 0 {
						tmpInterface = data[keyName].([]interface{})[len(data[keyName].([]interface{}))+index]
					} else {
						tmpInterface = data[keyName].([]interface{})[index]
					}

				} else {
					Logger.Warning("索引:%d超出数据范围，自动取第0个数据", index)
					tmpInterface = data[keyName].([]interface{})[0]
				}
				return sceneAssert.SpecialAssertValue(tmpInterface.(map[string]interface{}))
			}

			isHit, index, keyName := GetSlicesIndex(items[1])
			if !isHit {
				if strings.Contains(keyName, "-") {
					return sceneAssert.SpecialAssertValue(data[items[0]].(map[string]interface{}))
				} else {
					keyName = items[1]
				}

			}
			if isHit {
				dataLen := len(data[items[0]].([]interface{}))
				if dataLen == 0 {
					Logger.Warning("返回的数据列表为空，请核对")
				} else {
					var tmpDict map[string]interface{}
					if dataLen > index {
						if index < 0 {
							newIndex := dataLen + index
							if newIndex < 0 {
								err = fmt.Errorf("提供的索引超过数据范围了，数据长度: %v, 索引: %v", dataLen, index)
								Logger.Error("%s", err)
								return
							}
							tmpDict = data[items[0]].([]interface{})[dataLen+index].(map[string]interface{})
							outputDict[targetValueStr] = append(outputDict[targetValueStr], tmpDict[keyName])
						} else {
							tmpDict = data[items[0]].([]interface{})[index].(map[string]interface{})
							outputDict[targetValueStr] = append(outputDict[targetValueStr], tmpDict[keyName])
						}
					} else {
						Logger.Warning("索引:%d超出数据范围，自动取第0个数据", index)
						tmpDict = data[items[0]].([]interface{})[0].(map[string]interface{})
						outputDict[targetValueStr] = append(outputDict[targetValueStr], tmpDict[keyName])
					}
				}
			} else {
				var tmpDict map[string]interface{}
				varType := fmt.Sprintf("%T", data[items[0]])

				if varType != "[]interface {}" {
					err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
					Logger.Error("%s", err)
					return
				}
				for _, tmpInfo := range data[items[0]].([]interface{}) {
					varType := fmt.Sprintf("%T", tmpInfo)

					if varType != "map[string]interface {}" {
						err = fmt.Errorf("断言定义与实际返回结构不一致，请核对~")
						Logger.Error("%s", err)
						return
					}
					if strings.Contains(keyName, "-") {

						sceneAssert.Source = keyName

						return sceneAssert.SpecialAssertValue(tmpInfo.(map[string]interface{}))
					}

					tmpDict = tmpInfo.(map[string]interface{})
					outputDict[targetValueStr] = append(outputDict[targetValueStr], tmpDict[items[1]])
				}
			}
		}
	} else { // 获取最终的值
		if value, ok := data[sceneAssert.Source]; ok {
			subType := fmt.Sprintf("%T", data[sceneAssert.Source])
			if subType == "[]interface {}" {
				outputDict[targetValueStr] = data[sceneAssert.Source].([]interface{})
			} else {
				strValue, err1 := Interface2Str(value)
				if err1 != nil {
					err = err1
					Logger.Error("%s", err)
					return
				}
				outputDict[targetValueStr] = append(outputDict[targetValueStr], strValue)
			}
		} else {
			err1 := fmt.Errorf("未找到对应字段[%s]的值, 请核对~", sceneAssert.Source)
			err = err1
			Logger.Error("%s", err)
			return
		}
	}
	return
}

func (sceneAssert SceneAssert) AssertResult(data map[string]interface{}, inOutPutDict map[string][]interface{}) (b bool, err error) {
	// 如果返回的数据为空，则直接返回
	if len(data) == 0 {
		err = fmt.Errorf("无返回信息，无法解析输出参数")
		Logger.Error("%s", err)
		return
	}

	var tmpInterface interface{}

	targetValueStr, err := Interface2Str(sceneAssert.Value)

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if strings.Contains(targetValueStr, "{") && strings.Contains(targetValueStr, "}") {
		if strings.HasPrefix(targetValueStr, "{") && strings.HasSuffix(targetValueStr, "}") {
			//保留原逻辑
			keyNameTmp := strings.TrimLeft(targetValueStr, "{")
			keyName := strings.TrimRight(keyNameTmp, "}")
			if _, ok := inOutPutDict[keyName]; ok {
				targetValueStr, _ = Interface2Str(inOutPutDict[keyName][0])
			} else {
				err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", keyName)
				Logger.Error("%s", err)
				return
			}
		} else {
			// 替换类似 "{V1} + 1 + {V2}" 格式中的变量
			strReg := regexp.MustCompile(`\{.*?\}`)
			strMatch := strReg.FindAllString(targetValueStr, -1)
			for _, v := range strMatch {
				KeyStr := strings.Replace(v, "{", "", -1)
				KeyStr = strings.Replace(KeyStr, "}", "", -1)
				if _, ok := inOutPutDict[KeyStr]; ok {
					valueStr, _ := Interface2Str(inOutPutDict[KeyStr][0])
					targetValueStr = strings.Replace(targetValueStr, v, valueStr, -1)
				} else {
					err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", KeyStr)
					Logger.Error("%s", err)
					return
				}
			}

		}

	}

	strReg := regexp.MustCompile(`\[(\d+)\]|\[(-\d+)\]`)
	strMatch := strReg.FindAllString(sceneAssert.Source, -1)
	if strings.Contains(sceneAssert.Source, "-") && len(strMatch) >= 1 {
		outputTmp, err1 := sceneAssert.SpecialAssertValue(data)
		if err1 != nil {
			err = fmt.Errorf("返回数据与定义提取的数据结构不一致，请核对")
			return
		} else {
			for _, v := range outputTmp {
				expectValue, err2 := Interface2Str(v[0])
				if err2 != nil {
					err = err2
					Logger.Error("%s", err2)
					return
				}
				b, err = StrComparion(sceneAssert.Type, expectValue, targetValueStr)
				if err != nil {
					Logger.Error("%s", err)
					return
				}
			}
			return
		}
	}

	commonErr := fmt.Errorf("返回数据与断言定义提取的数据结构不一致，请核对")

	// 解析定义的校验参数
	if strings.Contains(sceneAssert.Source, "-") { // 字典 key 分隔
		items := strings.SplitN(sceneAssert.Source, "-", 2)
		sceneAssert.Source = items[1]
		if data[items[0]] == nil {
			err = commonErr
			//Logger.Debug("commonErr: %v", commonErr)
			return
		}
		tmpInterface = data[items[0]].(interface{})
		// 进入下一层数据解析，递归调用

		tmpType := fmt.Sprintf("%T", tmpInterface)
		if tmpType != "map[string]interface{}" && tmpType != "map[string]interface {}" {
			err = fmt.Errorf("%v, %v", commonErr, tmpType)
			return
		}

		return sceneAssert.AssertResult(tmpInterface.(map[string]interface{}), inOutPutDict)
	} else if strings.Contains(sceneAssert.Source, "*") { // 对数组内的字典进行 key 分隔
		items := strings.SplitN(sceneAssert.Source, "*", 2)
		sceneAssert.Source = items[1]
		if data[items[0]] != nil {
			tmpInterface = data[items[0]].(interface{})
			var tmpDict map[string]interface{}
			for _, tmpInfo := range tmpInterface.([]interface{}) {
				tmpType := fmt.Sprintf("%T", tmpInfo)
				if tmpType != "map[string]interface{}" && tmpType != "map[string]interface {}" {
					err = fmt.Errorf("%v, %v", commonErr, tmpType)
					Logger.Error("err: %v", err)
					return
				}
				tmpDict = tmpInfo.(map[string]interface{})
				curValueStr, err1 := Interface2Str(tmpDict[items[1]])
				if err1 != nil {
					err = err1
					Logger.Error("%s", err)
					return
				}

				// 比较当前值和目标值，返回 true 或 false
				b, err = StrComparion(sceneAssert.Type, curValueStr, targetValueStr)
				if err != nil {
					Logger.Error("%s", err)
					return
				}
			}
		}
	} else { // 获取最终的值
		if value, ok := data[sceneAssert.Source]; ok {
			curValueStr, err1 := Interface2Str(value)
			if err1 != nil {
				err = err1
				Logger.Error("%s", err)
				return
			}
			// 比较当前值和目标值，返回 true 或 false
			b, err = StrComparion(sceneAssert.Type, curValueStr, targetValueStr)
			if err != nil {
				return
			}
		} else {
			err1 := fmt.Errorf("未找到对应字段[%s]的值, 请核对~", sceneAssert.Source)
			err = err1
			Logger.Error("%s", err)
			return
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
