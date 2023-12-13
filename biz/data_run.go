package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GetFilePath(id string) (app, filePath string, err error) {
	var dbSceneData DbSceneData
	s, _ := strconv.Atoi(id)
	models.Orm.Table("scene_data").Where("id = ?", s).Find(&dbSceneData)
	if len(dbSceneData.ApiId) == 0 {
		err = fmt.Errorf("未找到对应[%v]的场景数据，请核对", s)
		Logger.Error("%s", err)
		return
	}
	filePath = fmt.Sprintf("%s/%s", DataBasePath, dbSceneData.FileName)
	app = dbSceneData.App
	_, err = os.Stat(filePath)

	content := GetStrFromHtml(dbSceneData.Content)

	if os.IsNotExist(err) {
		var dataFile DataFile
		var dataInfo []byte
		var err1 error
		if strings.HasSuffix(filePath, ".json") {
			err = json.Unmarshal([]byte(content), &dataFile)
			dataFile.Version = 1
			dataInfo, err1 = json.MarshalIndent(dataFile, "", "    ")
		} else {
			err = yaml.Unmarshal([]byte(content), &dataFile)
			dataFile.Version = 1
			dataInfo, err1 = yaml.Marshal(dataFile)
		}

		if err1 != nil {
			Logger.Error("%s", err1)
			err = err1
			return
		}
		err2 := ioutil.WriteFile(filePath, dataInfo, 0644)
		if err2 != nil {
			Logger.Error("%s", err2)
			err = err2
			return
		}
	}
	return
}

func GetResultFilePath(src string) (dst string, err error) {
	var dirName, targetDirName string
	fileName := path.Base(src)

	b, num := IsStrEndWithTimeFormat(fileName)
	suffix := GetStrSuffix(fileName)

	if b {
		dirName = fileName[:len(fileName)-num-len(suffix)]
	} else {
		dirName = fileName[:len(fileName)-len(suffix)]
	}

	curTime := time.Now().Format("20060102_150405.999999999")

	dst = fmt.Sprintf("%s/%s/%s_%s%s", HistoryBasePath, dirName, dirName, curTime, suffix)

	targetDirName = fmt.Sprintf("%s/%s", HistoryBasePath, dirName)

	// 判断目录是否存在，若不存在，则自动新建
	_, err = os.Stat(targetDirName)
	if os.IsNotExist(err) {
		//mask := syscall.Umask(0)  // 后续根据系统不同做适配
		//defer syscall.Umask(mask)
		err = os.MkdirAll(targetDirName, 0644)
		if err != nil {
			Logger.Error("%s", err)
			return
		}
	}

	return
}

func GetRunTimeData(id string) (dataInfo DbSceneData, appInfo EnvConfig, err error) {
	models.Orm.Table("scene_data").Where("id = ?", id).Find(&dataInfo)
	if len(dataInfo.ApiId) == 0 {
		err = fmt.Errorf("未找到对应的场景数据，请核对:%s", id)
		Logger.Error("%s", err)
		return
	}

	models.Orm.Table("env_config").Where("app = ?", dataInfo.App).Find(&appInfo)

	return
}

func (dataFile DataFile) GetIsParallel() (b bool) {
	if value := dataFile.IsParallel; value == "yes" {
		b = true
	} else {
		b = false
	}
	return
}

func (dataFile DataFile) GetIsRunPreApis() (b bool) {
	if value := dataFile.IsRunPreApis; value == "yes" {
		b = true
	} else {
		b = false
	}
	return
}

func (dataFile DataFile) GetIsRunPostApis() (b bool) {
	if value := dataFile.IsRunPostApis; value == "yes" {
		b = true
	} else {
		b = false
	}
	return
}

func (dataFile DataFile) GetHeader(envConfig EnvConfig) (header map[string]interface{}, err error) {
	header = make(map[string]interface{})
	header = CopyMapInterface(dataFile.Single.Header)

	if len(envConfig.Auth) != 0 {
		strTmp := GetStrFromHtml(envConfig.Auth)
		err = json.Unmarshal([]byte(strTmp), &header)
		if err != nil {
			Logger.Error("%s", err)
			return
		}
	}

	if dataFile.IsUseEnvConfig == "yes" {
		for k, v := range dataFile.Single.Header {
			if _, ok := header[k]; !ok {
				header[k] = v
			}
			if k == "Content-Type" {
				header[k] = v
			}
		}
	} else {
		for k, v := range dataFile.Single.Header {
			if _, ok := header[k]; ok {
				header[k] = v
			}
		}
	}

	lang := GetRequestLangage(header)

	for k, v := range dataFile.Single.Header {
		valueDef, err := Interface2Str(v)
		if err != nil {
			Logger.Error("%s", err)
		}

		if strings.Contains(valueDef, "{") && strings.Contains(valueDef, "}") {
			strByte := []byte(valueDef)
			strReg := regexp.MustCompile(`\{(.+)\}`)
			strMatch := strReg.FindAllSubmatch(strByte, -1)
			var key, rawDef string
			for _, item := range strMatch {
				key = string(item[1])
				rawDef = string(item[0])
				break
			}
			value := GetValueFromSysParameter(lang, key)
			if len(value) > 0 {
				header[k] = strings.ReplaceAll(valueDef, rawDef, value)
			}
		}
	}

	return
}

func GetSceneContent(fileName string) (sceneContent DataFile, err error) {
	var mode string
	if strings.HasSuffix(fileName, ".yml") {
		mode = "yaml"
	} else if strings.HasSuffix(fileName, ".json") {
		mode = "json"
	} else {
		err = fmt.Errorf("当前只支持.yml和.json尾缀文件，请核对文件名称: %s", fileName)
		Logger.Error("%s", err)
		return
	}
	filePath := fmt.Sprintf("%s/%s", DataBasePath, fileName)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("%s, 请核对文件: %s", err, fileName)
		Logger.Error("%s", err)
		return
	}
	if mode == "json" {
		err = json.Unmarshal([]byte(content), &sceneContent)
	} else if mode == "yaml" {
		err = yaml.Unmarshal([]byte(content), &sceneContent)
	}

	if err != nil {
		err = fmt.Errorf("%s, 请核对文件: %s", err, fileName)
		Logger.Error("%s", err)
		return
	}
	return
}

func (dataFile DataFile) GetUrl(envConfig EnvConfig) (rawUrls []string, err error) {
	var rawUrl string
	envInfo := []string{envConfig.Protocol, "://", envConfig.Ip, envConfig.Prepath, dataFile.Api.Path}
	sceneInfo := []string{dataFile.Env.Protocol, "://", dataFile.Env.Host, envConfig.Prepath, dataFile.Api.Path}
	tag := 0
	for i := 0; i < len(envInfo); i++ {
		if dataFile.IsUseEnvConfig == "no" {
			if len(sceneInfo[i]) == 0 && i != 3 {
				tag++
			}
			rawUrl = rawUrl + sceneInfo[i]
		} else {
			if len(envInfo[i]) == 0 && i != 3 {
				tag++
			}
			rawUrl = rawUrl + envInfo[i]
		}
	}

	if tag != 0 {
		err1 := fmt.Errorf("环境信息不完善,请检查, URL: %s，", rawUrl)
		err = err1
		Logger.Error("%s", err)
		return
	}

	if strings.Contains(rawUrl, "{") {
		if dataFile.Single.Path == nil && dataFile.Multi.Path == nil {
			err = fmt.Errorf("未进行Path变量定义，请先定义")
			Logger.Error("%s", err)
			return
		}
		pathVarsReg := regexp.MustCompile(`{[[:alpha:]]+}`)
		var pathVars []string
		pathVars = pathVarsReg.FindAllString(rawUrl, -1)
		for _, v := range pathVars {
			str1 := v[1 : len(v)-1]
			var tag int
			tag = 0
			if value, ok := dataFile.Single.Path[str1]; ok {
				valueStr, _ := Interface2Str(value)
				url := strings.Replace(rawUrl, v, valueStr, -1)
				rawUrls = append(rawUrls, url)
				tag = tag + 1
			}

			if values, ok := dataFile.Multi.Path[str1]; ok {
				for _, value := range values {
					valueStr, _ := Interface2Str(value)
					url := strings.Replace(rawUrl, v, valueStr, -1)
					rawUrls = append(rawUrls, url)
					Logger.Debug("rawUrl: %v", rawUrl)
				}
				tag = tag + 1
			}

			if tag == 0 {
				err = fmt.Errorf("未找到Path:%s变量的定义值，请先进行定义", v)
				Logger.Error("%s", err)
				return
			}
		}
	} else {
		rawUrls = append(rawUrls, rawUrl)
	}
	return
}

func (dataFile DataFile) GetQuery(lang string, depOutVars map[string][]interface{}) (querys []map[string]interface{}, err error) {
	var query map[string]interface{}
	query = make(map[string]interface{})

	if len(dataFile.Single.Query) > 0 {
		for k, v := range dataFile.Single.Query {
			strK, _ := Interface2Str(v)
			t, subV, allDef := GetStrType(strK, lang)

			if t == 1 {
				if value, ok := depOutVars[k]; ok {
					query[k] = value[0]
				} else {
					err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", k)
					Logger.Error("%s", err)
					if len(depOutVars) > 0 {
						Logger.Info("outputDict: %v", depOutVars)
					}
					return
				}
			} else if t == 2 {
				query[k] = subV
			} else if t == 3 {
				newStr := strK
				for defKey, defValue := range allDef {
					var tmpStrV string
					if value, ok := depOutVars[defKey]; ok {
						tmpStrV, _ = Interface2Str(value[0])
					} else {
						err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", defKey)
						Logger.Error("%s", err)
						return
					}
					newStr = strings.Replace(newStr, defValue, tmpStrV, -1)
				}
				query[k] = newStr
			} else {
				query[k] = v
			}
		}
	}

	if len(dataFile.Multi.Query) > 0 {
		minLen := GetSliceMinLen(dataFile.Multi.Query)
		for i := 0; i < minLen; i++ {
			for k, v := range dataFile.Multi.Query {
				strK, err1 := Interface2Str(v[i])
				if err1 != nil {
					err = err1
					Logger.Error("%s", err)
					return
				}

				t, subV, allDef := GetStrType(strK, lang)

				if t == 1 {
					if value, ok := depOutVars[k]; ok {
						if len(value) > i {
							query[k] = value[i]
						} else {
							query[k] = value[0]
						}
					} else {
						err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", k)
						Logger.Error("%s", err)
						return
					}
				} else if t == 2 {
					query[k] = subV
				} else if t == 3 {
					for defKey, defValue := range allDef {
						var tmpStrV string
						if value, ok := depOutVars[defKey]; ok {
							if len(value) > i {
								tmpStrV, _ = Interface2Str(value[i])
							} else {
								tmpStrV, _ = Interface2Str(value[0])
							}
						} else {
							err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", defKey)
							Logger.Error("%s", err)
							return
						}
						newStr := strings.Replace(strK, defValue, tmpStrV, -1)
						query[k] = newStr
					}
				} else {
					query[k] = v[i]
				}
			}
			queryTmp := CopyMap(query)
			querys = append(querys, queryTmp)
		}
	} else {
		if len(query) > 0 {
			querys = append(querys, query)
		}
	}

	return
}

func (dataFile DataFile) GetBody(lang string, depOutVars map[string][]interface{}) (bodys []map[string]interface{}, bodyAfterList []interface{}, err error) {
	if dataFile.Api.Method == "get" {
		return
	}

	if dataFile.Single.BodyList != nil {
		dst := CopyList(dataFile.Single.BodyList)
		bodyAfterList = make([]interface{}, 0, len(dataFile.Single.BodyList))
		bodyAfterList, err = GetAfterListBody(lang, dst, depOutVars)
	} else {
		var bodyTmp, body map[string]interface{}
		bodyTmp = make(map[string]interface{})
		if dataFile.Single.Body != nil {
			bodyTmp = CopyMap(dataFile.Single.Body)
		}
		body, err = GetAfterBody(lang, bodyTmp, depOutVars)
		if err != nil {
			return
		}

		if len(dataFile.Multi.Body) > 0 {
			minLen := GetSliceMinLen(dataFile.Multi.Body)
			for i := 0; i < minLen; i++ {
				for k, v := range dataFile.Multi.Body {
					strK, err1 := Interface2Str(v[i])
					if err1 != nil {
						err = err1
						Logger.Error("%s", err)
						return
					}

					t, subV, allDef := GetStrType(strK, lang)

					if t == 1 {
						if value, ok := depOutVars[k]; ok {
							if len(value) > i {
								body[k] = value[i]
							} else {
								body[k] = value[0]
							}
						} else {
							err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", k)
							Logger.Error("%s", err)
							return
						}
					} else if t == 2 {
						body[k] = subV
					} else if t == 3 {
						var newStr string
						if len(subV) > 0 {
							newStr = subV
						} else {
							newStr = strK
						}
						for defKey, defValue := range allDef {
							var tmpStrV string
							if value, ok := depOutVars[defKey]; ok {
								if len(value) > i {
									tmpStrV, _ = Interface2Str(value[i])
								} else {
									tmpStrV, _ = Interface2Str(value[0])
								}
							} else {
								err = fmt.Errorf("未找到变量[%s]定义，请先关联或定义", defKey)
								Logger.Error("%s", err)
								return
							}
							newStr = strings.Replace(newStr, defValue, tmpStrV, -1)
						}
						body[k] = newStr
					} else {
						body[k] = v[i]
					}
				}
				bodyTmp := CopyMap(body)
				bodys = append(bodys, bodyTmp)
			}
		} else {
			if len(body) > 0 {
				bodys = append(bodys, body)
			}
		}
	}
	return
}

func (dataFile DataFile) CreatActionData(lang string, depOutVars map[string][]interface{}) (err error) {
	if len(dataFile.Action) > 0 {
		fileName := ""
		dataCount := 0
		//Logger.Debug("开始生成文件")
		for _, item := range dataFile.Action {
			if item.Type == "create_csv" {
				tmpValue, _ := Interface2Str(item.Value)
				strList := strings.Split(tmpValue, ":")
				if len(strList) == 0 {
					err = fmt.Errorf("creat_csv的值未定义，请先定义")
					return
				}
				if !strings.Contains(strList[0], ".csv") {
					fileName = fmt.Sprintf("%s.csv", strList[0])
				} else {
					fileName = strList[0]
				}

				if len(strList) >= 2 {
					dataCount, _ = strconv.Atoi(strList[1])
				} else {
					dataCount = 100
				}
				break
			}
		}

		if dataCount > 0 {
			filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
			var keyList []string
			Logger.Debug("dataCount: %v", dataCount)
			for i := 0; i < dataCount; i++ {
				bodys, _, errTmp := dataFile.GetBody(lang, depOutVars)
				if errTmp != nil {
					err = errTmp
					return
				}
				if i == 0 {
					for index, item := range bodys {
						keyStr := ""
						valueStr := ""
						for k, v := range item {
							keyList = append(keyList, k)
							if len(keyStr) == 0 {
								keyStr = fmt.Sprintf("%v", k)
								valueStr = fmt.Sprintf("%v", v)
							} else {
								keyStr = fmt.Sprintf("%s,%v", keyStr, k)
								valueStr = fmt.Sprintf("%s,%v", valueStr, v)
							}
						}
						if index == 0 {
							tStr := fmt.Sprintf("%s\n", keyStr)
							_ = WriteDataInCSV(filePath, tStr)
							tStr = fmt.Sprintf("%s\n", valueStr)
							_ = WriteDataInCSV(filePath, tStr)
						} else {
							tStr := fmt.Sprintf("%s\n", valueStr)
							_ = WriteDataInCSV(filePath, tStr)
						}

					}
				} else {
					for _, item := range bodys {
						valueStr := ""
						for _, key := range keyList {
							if len(valueStr) == 0 {
								valueStr = fmt.Sprintf("%v", item[key])
							} else {
								valueStr = fmt.Sprintf("%s,%v", valueStr, item[key])
							}
						}
						tStr := fmt.Sprintf("%s\n", valueStr)
						_ = WriteDataInCSV(filePath, tStr)

					}
				}
			}
		}
	}

	if len(dataFile.Action) > 0 {
		fileName := ""
		dataCount := 0
		for _, item := range dataFile.Action {
			if item.Type == "create_excel" || item.Type == "create_xls" || item.Type == "create_xlsx" {
				tmpValue, _ := Interface2Str(item.Value)
				strList := strings.Split(tmpValue, ":")
				if len(strList) == 0 {
					err = fmt.Errorf("creat_excel的值未定义，请先定义")
					return
				}
				if !strings.Contains(strList[0], ".xlsx") {
					fileName = fmt.Sprintf("%s.xlsx", strList[0])
				} else {
					fileName = strList[0]
				}

				if len(strList) >= 2 {
					dataCount, _ = strconv.Atoi(strList[1])
				} else {
					dataCount = 100
				}
				break
			}
		}

		if dataCount > 0 {
			filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
			var keyList []string
			Logger.Debug("dataCount: %v", dataCount)
			for i := 0; i < dataCount; i++ {
				bodys, _, errTmp := dataFile.GetBody(lang, depOutVars)
				if errTmp != nil {
					err = errTmp
					return
				}
				if i == 0 {
					for _, item := range bodys {
						var valueList []string
						for k, v := range item {
							keyList = append(keyList, k)
							vStr, _ := Interface2Str(v)
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, keyList)
						_ = WriteDataInXls(filePath, valueList)
					}
				} else {
					for _, item := range bodys {
						var valueList []string
						for _, key := range keyList {
							vStr, _ := Interface2Str(item[key])
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, valueList)
					}
				}
			}
		}
		//Logger.Debug("结束生成文件")
	}

	if len(dataFile.Action) > 0 {
		fileName := ""
		dataCount := 0
		for _, item := range dataFile.Action {
			if item.Type == "create_excel" || item.Type == "create_xls" || item.Type == "create_xlsx" {
				tmpValue, _ := Interface2Str(item.Value)
				strList := strings.Split(tmpValue, ":")
				if len(strList) == 0 {
					err = fmt.Errorf("creat_excel的值未定义，请先定义")
					return
				}
				if !strings.Contains(strList[0], ".xlsx") {
					fileName = fmt.Sprintf("%s.xlsx", strList[0])
				} else {
					fileName = strList[0]
				}

				if len(strList) >= 2 {
					dataCount, _ = strconv.Atoi(strList[1])
				} else {
					dataCount = 100
				}
				break
			}
		}

		if dataCount > 0 {
			filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
			var keyList []string
			Logger.Debug("dataCount: %v", dataCount)
			for i := 0; i < dataCount; i++ {
				bodys, _, errTmp := dataFile.GetBody(lang, depOutVars)
				if errTmp != nil {
					err = errTmp
					return
				}
				if i == 0 {
					for _, item := range bodys {
						var valueList []string
						for k, v := range item {
							keyList = append(keyList, k)
							vStr, _ := Interface2Str(v)
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, keyList)
						_ = WriteDataInXls(filePath, valueList)
					}
				} else {
					for _, item := range bodys {
						var valueList []string
						for _, key := range keyList {
							vStr, _ := Interface2Str(item[key])
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, valueList)
					}
				}
			}
		}
		//Logger.Debug("结束生成文件")
	}

	return
}

func (dataFile DataFile) CreatActionDataOrderByKey(lang string, depOutVars map[string][]interface{}) (err error) {
	var isCreateCSV, isCreateXLS, isCreateHiveSQL bool
	var csvValue, xlsValue, hiveSQLValue interface{}

	if len(dataFile.Action) > 0 {
		for _, item := range dataFile.Action {
			if item.Type == "create_csv" {
				isCreateCSV = true
				csvValue = item.Value
			} else if item.Type == "create_excel" || item.Type == "create_xls" || item.Type == "create_xlsx" {
				isCreateXLS = true
				xlsValue = item.Value
			} else if item.Type == "create_hive_table_sql" {
				isCreateHiveSQL = true
				hiveSQLValue = item.Value
			}
		}
	}

	if isCreateCSV {
		fileName := ""
		dataCount := 0
		tmpValue, _ := Interface2Str(csvValue)
		strList := strings.Split(tmpValue, ":")
		if len(strList) == 0 {
			err = fmt.Errorf("creat_csv的值未定义，请先定义")
			return
		}
		if !strings.Contains(strList[0], ".csv") {
			fileName = fmt.Sprintf("%s.csv", strList[0])
		} else {
			fileName = strList[0]
		}

		if len(strList) >= 2 {
			dataCount, _ = strconv.Atoi(strList[1])
		} else {
			dataCount = 100
		}

		if dataCount > 0 {
			filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
			var keyList []string
			tmpBody := make(map[string]interface{})
			splitTag := ","
			for i := 0; i < dataCount; i++ {
				bodys, _, errTmp := dataFile.GetBody(lang, depOutVars)
				if errTmp != nil {
					err = errTmp
					return
				}

				if i == 0 {
					tmpBody = bodys[0]
					for k, v := range tmpBody {
						keyList = append(keyList, k)
						vStr, _ := Interface2Str(v)
						if strings.Contains(vStr, ",") {
							splitTag = "|"
						}
					}
					sort.Strings(keyList)
				}

				if i == 0 {
					for index, item := range bodys {
						keyStr := ""
						valueStr := ""
						for _, k := range keyList {
							if len(keyStr) == 0 {
								keyStr = fmt.Sprintf("%v", k)
								valueStr = fmt.Sprintf("%v", item[k])
							} else {
								keyStr = fmt.Sprintf("%s%s%v", keyStr, splitTag, k)
								valueStr = fmt.Sprintf("%s%s%v", valueStr, splitTag, item[k])
							}
						}
						if index == 0 {
							tStr := fmt.Sprintf("%s\n", keyStr)
							_ = WriteDataInCSV(filePath, tStr)
							tStr = fmt.Sprintf("%s\n", valueStr)
							_ = WriteDataInCSV(filePath, tStr)
						} else {
							tStr := fmt.Sprintf("%s\n", valueStr)
							_ = WriteDataInCSV(filePath, tStr)
						}

					}
				} else {
					for _, item := range bodys {
						valueStr := ""
						for _, key := range keyList {
							if len(valueStr) == 0 {
								valueStr = fmt.Sprintf("%v", item[key])
							} else {
								valueStr = fmt.Sprintf("%s%s%v", valueStr, splitTag, item[key])
							}
						}
						tStr := fmt.Sprintf("%s\n", valueStr)
						_ = WriteDataInCSV(filePath, tStr)

					}
				}
			}
		}
	}

	if isCreateXLS {
		fileName := ""
		dataCount := 0
		tmpValue, _ := Interface2Str(xlsValue)
		strList := strings.Split(tmpValue, ":")
		if len(strList) == 0 {
			err = fmt.Errorf("creat_excel的值未定义，请先定义")
			return
		}
		if !strings.Contains(strList[0], ".xlsx") {
			fileName = fmt.Sprintf("%s.xlsx", strList[0])
		} else {
			fileName = strList[0]
		}

		if len(strList) >= 2 {
			dataCount, _ = strconv.Atoi(strList[1])
		} else {
			dataCount = 100
		}

		if dataCount > 0 {
			filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
			var keyList []string
			tmpBody := make(map[string]interface{})
			for i := 0; i < dataCount; i++ {
				bodys, _, errTmp := dataFile.GetBody(lang, depOutVars)
				if errTmp != nil {
					err = errTmp
					return
				}
				if i == 0 {
					tmpBody = bodys[0]
					for k, _ := range tmpBody {
						keyList = append(keyList, k)
					}
					sort.Strings(keyList)
				}

				if i == 0 {
					for _, item := range bodys {
						var valueList []string
						for _, k := range keyList {
							keyList = append(keyList, k)
							vStr, _ := Interface2Str(item[k])
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, keyList)
						_ = WriteDataInXls(filePath, valueList)
					}
				} else {
					for _, item := range bodys {
						var valueList []string
						for _, key := range keyList {
							vStr, _ := Interface2Str(item[key])
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, valueList)
					}
				}
			}
		}
	}

	if isCreateHiveSQL {
		fileName := ""
		tmpValue, _ := Interface2Str(hiveSQLValue)
		if len(tmpValue) == 0 {
			err = fmt.Errorf("creat_excel的值未定义，请先定义")
			return
		}
		if !strings.Contains(tmpValue, ".sql") {
			fileName = fmt.Sprintf("%s.sql", tmpValue)
		} else {
			fileName = tmpValue
		}

		filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
		var keyList []string
		tmpBody := make(map[string]interface{})

		bodys, _, errTmp := dataFile.GetBody(lang, depOutVars)
		if errTmp != nil {
			err = errTmp
			return
		}

		tmpBody = bodys[0]
		for k, _ := range tmpBody {
			keyList = append(keyList, k)
		}
		sort.Strings(keyList)

		var parameterType, sqlStr string
		headStr := "CREATE TABLE IF NOT EXISTS default.all_type_data\n(\n"
		tailStr := "\n) PARTITIONED BY (ds STRING) STORED AS TEXTFILE;"
		midStr := ""
		for _, k := range keyList {
			varType := fmt.Sprintf("%T", tmpBody[k])
			switch varType {
			case "float64":
				parameterType = "DOUBLE"
			case "string":
				parameterType = "STRING"
			case "bool":
				parameterType = "BOOLEAN"
			case "int":
				parameterType = "INT"
			default:
				parameterType = "STRING"
			}

			if len(midStr) == 0 {
				midStr = fmt.Sprintf("  %s %s", k, parameterType)
			} else {
				midStr = fmt.Sprintf("%s,\n  %s %s", midStr, k, parameterType)
			}
		}
		sqlStr = fmt.Sprintf("%s%s%s", headStr, midStr, tailStr)

		_ = WriteDataInCSV(filePath, sqlStr)
	}

	return
}

func (dataFile DataFile) GetDepParams() (depOutDict map[string][]interface{}, err error) {
	depOutDict = make(map[string][]interface{})
	var tmpDict map[string][]interface{}
	tmpDict = make(map[string][]interface{})
	var fileNames []string

	fileNames = append(dataFile.Api.ParamApis, dataFile.Api.PreApi...)

	if len(fileNames) > 0 {
		for _, fileName := range fileNames {
			sceneContent, err1 := GetSceneContent(fileName)
			if err1 != nil {
				err = err1
				Logger.Error("%s", err)
				return
			}
			tmpDict = CopyMapList(sceneContent.Output)
			for k, v := range tmpDict {
				depOutDict[k] = v
			}
		}
	} // else {   // 提示信息暂且不展示
	//	Logger.Warning("如需其他API 提供关联参数，请先关联前置API(pre_apis)或关联API参数(para_apis)")
	//}

	return
}

func (dataFile DataFile) GetResult(lang, source, filePath string, header map[string]interface{}, isThread string, res [][]byte, inOutPutDict map[string][]interface{}, errs []error) (result, dst string, err error) {
	var outputDict map[string][]interface{}
	outputDict = make(map[string][]interface{})
	isPass := 0

	dst, err = GetResultFilePath(filePath)

	if err != nil {
		return
	}

	for i := 0; i < len(res); i++ {
		if i == 0 {
			dataFile.Response = []string{string(res[i])}
			// 请求时有返回 Error 信息，结果设置为失败，不再走后续流程
			if len(errs) > i {
				if errs[i] != nil {
					dataFile.TestResult = []string{"fail"}
					isPass++
					if err != nil {
						err = fmt.Errorf("%v,%v", err, errs[i])
					} else {
						err = errs[i]
					}
					continue
				}
			}
		} else {
			dataFile.Response = append(dataFile.Response, string(res[i]))

			if len(errs) > i {
				if errs[i] != nil {
					dataFile.TestResult = append(dataFile.TestResult, "fail")
					isPass++
					if err != nil {
						err = fmt.Errorf("%v,%v", err, errs[i])
					} else {
						err = errs[i]
					}
					continue
				}
			}
		}

		// 未设置断言时，结果设置为成功，不再走后续流程
		if len(dataFile.Assert) == 0 {
			if i == 0 {
				dataFile.TestResult = []string{"pass"}
			} else {
				dataFile.TestResult = append(dataFile.TestResult, "pass")
			}
			continue
		}

		// 加载返回信息Response，若不是标准的 json 格式，则结果设置为失败，不再走后续流程
		var resDict map[string]interface{}
		resDict = make(map[string]interface{})
		var errTmp error
		if len(res[i]) == 0 {
			errTmp = fmt.Errorf("Response为空，无法做数据校验，请核对")
		} else {
			errTmp = json.Unmarshal(res[i], &resDict)
		}

		if errTmp != nil {
			if i == 0 {
				dataFile.TestResult = []string{"fail"}
			} else {
				dataFile.TestResult = append(dataFile.TestResult, "fail")
			}
			isPass++
			if err != nil {
				err = fmt.Errorf("%v,%v", err, errTmp)
			} else {
				err = errTmp
			}

			continue
		}
		for _, assert := range dataFile.Assert {
			aType := assert.Type
			if isPass != 0 && aType == "output" {
				continue
			}

			assert.Value = assert.GetAssertValue(lang)

			if assert.Source == "raw" {
				b, err1 := RawStrComparion(aType, string(res[i]), assert.Value)
				if !b {
					if i == 0 {
						dataFile.TestResult = []string{"fail"}
					} else {
						dataFile.TestResult = append(dataFile.TestResult, "fail")
					}
					isPass++
					if err != nil {
						err = fmt.Errorf("%s; %s", err, err1)
					} else {
						err = fmt.Errorf("Response: %s; %s", string(res[i]), err1)
					}
				} else {
					if i == 0 {
						dataFile.TestResult = []string{"pass"}
					} else {
						dataFile.TestResult = append(dataFile.TestResult, "pass")
					}
				}

				if !b {
					isPass++
				}

			} else {
				switch aType {
				case "output":
					outputTmp, err1 := assert.GetOutput(resDict)
					if err1 != nil {
						if i == 0 {
							dataFile.TestResult = []string{"fail"}
						} else {
							dataFile.TestResult = append(dataFile.TestResult, "fail")
						}
						isPass++
						if err != nil {
							err = fmt.Errorf("%s, %s", err, err1)
						} else {
							err = err1
						}
						continue
					} else {
						if i == 0 {
							dataFile.TestResult = []string{"pass"}
						} else {
							dataFile.TestResult = append(dataFile.TestResult, "pass")
						}
					}
					for k, v := range outputTmp {
						outputDict[k] = append(outputDict[k], v...)
					}
				default:
					b, err1 := assert.AssertResult(resDict, inOutPutDict)
					if err1 != nil {
						if i == 0 {
							dataFile.TestResult = []string{"fail"}
						} else {
							dataFile.TestResult = append(dataFile.TestResult, "fail")
						}
						if err != nil {
							err = fmt.Errorf("%s, %s", err, err1)
						} else {
							err = err1
						}
						isPass++
					} else {
						if i == 0 {
							dataFile.TestResult = []string{"pass"}
						} else {
							dataFile.TestResult = append(dataFile.TestResult, "pass")
						}
					}
					if !b {
						isPass++
					}
				}
			}

		}
	}

	dataFile.Output = outputDict

	if isThread == "yes" {
		Logger.Debug("dst: %s", dst)
	}

	var dataInfo, dataWithHeaher []byte
	var errTmp error
	if strings.HasSuffix(filePath, ".json") {
		dataInfo, errTmp = json.MarshalIndent(dataFile, "", "    ")
	} else {
		dataInfo, errTmp = yaml.Marshal(dataFile)
	}

	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%v,%v", err, errTmp)
		} else {
			err = errTmp
		}
	}

	dataFile.Single.Header = header
	if strings.HasSuffix(filePath, ".json") {
		dataWithHeaher, errTmp = json.MarshalIndent(dataFile, "", "    ")
	} else {
		dataWithHeaher, errTmp = yaml.Marshal(dataFile)
	}

	errTmp = ioutil.WriteFile(dst, dataWithHeaher, 0644)
	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%v,%v", err, errTmp)
		} else {
			err = errTmp
		}
	}

	if source == "data" || source == "playbook" {
		errTmp = ioutil.WriteFile(filePath, dataInfo, 0644)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%v,%v", err, errTmp)
			} else {
				err = errTmp
			}
		}
	}

	if isPass != 0 {
		result = "fail"
	} else {
		result = "pass"
	}

	return
}

func WriteDataResultByFile(src, result, dst, product string, envType int, errIn error) (err error) {
	var dataFile DataFile
	content, err := ioutil.ReadFile(src)
	if err != nil {
		Logger.Error("Error: %s, filePath: %s", err, src)
		return
	}
	if strings.HasSuffix(src, ".yml") || strings.HasSuffix(src, ".yaml") {
		err = yaml.Unmarshal([]byte(content), &dataFile)
	} else if strings.HasSuffix(src, ".json") {
		err = json.Unmarshal([]byte(content), &dataFile)
	} else {
		err = fmt.Errorf("不支持当前文件: %s, 请核对", src)
		Logger.Error("%s", err)
	}

	var sceneDataRecord SceneDataRecord

	dirName := GetHistoryDataDirName(path.Base(dst))

	apiStr := fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, path.Base(dst), path.Base(dst))

	sceneDataRecord.Content = apiStr
	sceneDataRecord.Name = dataFile.Name

	sceneDataRecord.ApiId = dataFile.ApiId
	sceneDataRecord.App = dataFile.Api.App
	sceneDataRecord.Result = result
	sceneDataRecord.EnvType = envType
	sceneDataRecord.Product = product

	if errIn != nil {
		sceneDataRecord.FailReason = fmt.Sprintf("%s", errIn)
	}

	err = models.Orm.Table("scene_data_test_history").Create(sceneDataRecord).Error

	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func WriteSceneDataResult(id string, result, dst, product string, envType int, errIn error) (err error) {
	var dbSceneData DbSceneData
	s, _ := strconv.Atoi(id)
	models.Orm.Table("scene_data").Where("id = ?", s).Find(&dbSceneData)
	if len(dbSceneData.ApiId) == 0 {
		return
	}

	if len(result) == 0 {
		dbSceneData.Result = " "
	} else {
		dbSceneData.Result = result
	}

	if errIn != nil {
		dbSceneData.FailReason = fmt.Sprintf("%s", errIn)
	} else {
		dbSceneData.FailReason = " " // 用空字符串刷新数据结果
	}
	err = models.Orm.Table("scene_data").Where("id = ?", dbSceneData.Id).Update(&dbSceneData).Error
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	var sceneDataRecord SceneDataRecord
	sceneDataRecord.Content = path.Base(dst)
	sceneDataRecord.Name = dbSceneData.Name
	sceneDataRecord.ApiId = dbSceneData.ApiId
	sceneDataRecord.App = dbSceneData.App
	sceneDataRecord.Result = dbSceneData.Result
	sceneDataRecord.FailReason = dbSceneData.FailReason
	sceneDataRecord.EnvType = envType
	sceneDataRecord.Product = product

	err = models.Orm.Table("scene_data_test_history").Create(sceneDataRecord).Error

	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func RunSceneContent(app, filePath, product, isThread string) (result, dst string, err error) {
	var dataFile DataFile
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if strings.HasSuffix(filePath, ".json") {
		err = json.Unmarshal([]byte(content), &dataFile)
	} else {
		err = yaml.Unmarshal([]byte(content), &dataFile)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if len(dataFile.Api.PreApi) > 0 && dataFile.IsRunPreApis == "yes" {
		for _, preApiFile := range dataFile.Api.PreApi {
			preFilePath := fmt.Sprintf("%s/%s", DataBasePath, preApiFile)
			Logger.Debug("开始执行前置用例: %v", preFilePath)
			result, dst, err = RunSceneContent(app, preFilePath, product, isThread)
			if err != nil {
				Logger.Error("%s", err)
				return
			}
			if result == "fail" {
				return
			}
		}
	}

	var envConfig EnvConfig
	if dataFile.IsUseEnvConfig == "no" {
		envConfig.App = dataFile.Api.App
		hJson, _ := json.Marshal(dataFile.Single.Header)
		envConfig.Auth = string(hJson)
		envConfig.Prepath = dataFile.Env.Prepath
		envConfig.Protocol = dataFile.Env.Protocol
		envConfig.MaxThreadNum = 1
	} else {
		if len(app) == 0 {
			envConfig, err = GetEnvConfig(dataFile.Api.App, "data")
		} else {
			envConfig, err = GetEnvConfig(app, "data")
		}

		if err != nil {
			Logger.Error("%s", err)
		}

		if len(product) > 0 {
			envConfig.Product = product
			sceneEnvConfig, errTmp := GetEnvConfig(envConfig.Product, "scene")
			if errTmp != nil {
				Logger.Error("%s", errTmp)
			}
			envConfig.Ip = sceneEnvConfig.Ip
			envConfig.Auth = sceneEnvConfig.Auth

		}

	}

	header, err := dataFile.GetHeader(envConfig)

	lang := GetRequestLangage(header)

	if err != nil {
		Logger.Debug("%s", err)
		return
	}
	depOutVars, err1 := dataFile.GetDepParams()
	if err1 != nil {
		err = err1
		Logger.Debug("%s", err)
		return
	}

	if len(product) > 0 {
		dbProduct, err := GetProductInfo(product)
		if err != nil {
			Logger.Error("%v", err)
		}
		privateParameter := dbProduct.GetPrivateParameter()
		for k, v := range privateParameter {
			depOutVars[k] = append(depOutVars[k], v)
		}

	}

	urls, err := dataFile.GetUrl(envConfig)
	dataFile.Urls = urls
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	querys, err := dataFile.GetQuery(lang, depOutVars)
	if err != nil {
		Logger.Error("err: %s", err)
		return
	}

	bodys, bodyList, err := dataFile.GetBody(lang, depOutVars)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	var resList [][]byte
	var errs []error
	tag := 0

	if dataFile.GetIsParallel() {
		wg := sync.WaitGroup{}
		for _, url := range urls {
			if len(querys) > 0 {
				for _, data := range querys {
					dJson, _ := json.Marshal(data)
					if tag == 0 {
						dataFile.Request = []string{string(dJson)}
					} else {
						dataFile.Request = append(dataFile.Request, string(dJson))
					}
					tag++
					wg.Add(1)
					go func(method, url string, data map[string]interface{}, header map[string]interface{}) {
						defer wg.Add(-1)
						res, err := RunHttp(method, url, data, header)
						resList = append(resList, res)
						errs = append(errs, err)
					}(dataFile.Api.Method, url, data, header)
				}
			} else if len(bodys) > 0 || len(bodyList) > 0 {
				if len(bodyList) > 0 {
					if len(bodyList) > 0 {
						var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
						readerNew, _ := jsonNew.Marshal(&bodyList)
						dataFile.Request = []string{string(readerNew)}
						res, err := RunHttpJsonList(dataFile.Api.Method, url, bodyList, header)
						if err != nil {
							Logger.Debug("%s", err)
						}
						resList = append(resList, res)
						errs = append(errs, err)
					}
				} else {
					for _, data := range bodys {
						dJson, _ := json.Marshal(data)
						if tag == 0 {
							dataFile.Request = []string{string(dJson)}
						} else {
							dataFile.Request = append(dataFile.Request, string(dJson))
						}
						tag++
						wg.Add(1)
						go func(method, url string, data map[string]interface{}, header map[string]interface{}) {
							defer wg.Add(-1)
							res, err := RunHttp(method, url, data, header)
							resList = append(resList, res)
							errs = append(errs, err)
						}(dataFile.Api.Method, url, data, header)
					}
				}
			} else {
				dataFile.Request = []string{}
				wg.Add(1)
				go func(method, url string, header map[string]interface{}) {
					res, err := RunHttp(method, url, nil, header)
					resList = append(resList, res)
					errs = append(errs, err)
				}(dataFile.Api.Method, url, header)
			}
			wg.Wait()
		}
	} else {
		for _, url := range urls {
			if len(querys) > 0 {
				for _, data := range querys {
					dJson, _ := json.Marshal(data)
					if tag == 0 {
						dataFile.Request = []string{string(dJson)}
					} else {
						dataFile.Request = append(dataFile.Request, string(dJson))
					}
					tag++
					if dataFile.Api.Method == "delete" {
						subTag := 0
						for k, v := range data {
							strV, _ := Interface2Str(v)
							if subTag == 0 {
								url = fmt.Sprintf("%s?%s=%s", url, k, strV)
							} else {
								url = fmt.Sprintf("%s&%s=%s", url, k, strV)
							}
							subTag++
						}
						res, err := RunHttp(dataFile.Api.Method, url, nil, header)
						resList = append(resList, res)
						errs = append(errs, err)
					} else {
						res, err := RunHttp(dataFile.Api.Method, url, data, header)
						if err != nil {
							Logger.Debug("%s", err)
						}
						resList = append(resList, res)
						errs = append(errs, err)
					}

					_ = dataFile.SetSleepAction()
				}
			} else if len(bodys) > 0 {
				for _, data := range bodys {
					var dJson []byte
					dJson, errTmp := json.Marshal(data)
					if errTmp != nil {
						var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
						dJsonTmp, err2 := jsonNew.Marshal(&data)
						if err2 != nil {
							Logger.Error("%s", err2)
							err = err2
							return
						}
						dJson = dJsonTmp
					}
					if tag == 0 {
						dataFile.Request = []string{string(dJson)}
					} else {
						dataFile.Request = append(dataFile.Request, string(dJson))
					}
					tag++
					res, err := RunHttp(dataFile.Api.Method, url, data, header)
					if err != nil {
						Logger.Debug("%s", err)
					}
					resList = append(resList, res)
					errs = append(errs, err)
					_ = dataFile.SetSleepAction()
				}
			} else {
				dataFile.Request = []string{}
				res, err := RunHttp(dataFile.Api.Method, url, nil, header)
				if err != nil {
					Logger.Debug("%s", err)
				}
				resList = append(resList, res)
				errs = append(errs, err)
			}
		}
	}

	result, dst, err = dataFile.GetResult(lang, "", filePath, header, isThread, resList, depOutVars, errs)

	if result != "pass" {
		for _, item := range errs {
			if item != nil {
				if err != nil {
					err = fmt.Errorf("%s;%s", err, item)
				} else {
					err = fmt.Errorf("%s", item)
				}
			}
		}
	}

	if err != nil {
		Logger.Error("%s", err)
	}

	if result != "pass" {
		Logger.Debug("header: %v", header)
	}
	return
}

func RunSceneData(id, product string) (err error) {
	dataInfo, appInfo, err := GetRunTimeData(id)
	var envType, maxThreadNum int
	var isThread string
	if len(product) > 0 {
		productTaskInfo, _ := GetProductInfo(product)
		envType = productTaskInfo.EnvType
		isThread = productTaskInfo.Threading
		maxThreadNum = productTaskInfo.ThreadNumber
	} else {
		envType = 2
		isThread = appInfo.Threading
		maxThreadNum = appInfo.MaxThreadNum
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	// 数据执行为多次性时，自动生成任务，方便管理  //功能先屏蔽
	//if runNum > 10 && mode != "task"{
	//	err = AutoCreateSchedule(id, "data")
	//	return
	//}

	if dataInfo.RunTime == 0 {
		err = fmt.Errorf("执行次数为: %d, 如需运行，请把执行次数置为大于0的数", dataInfo.RunTime)
		Logger.Warning("%s", err)
		return
	}

	app, filePath, err := GetFilePath(id)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	if isThread == "yes" && dataInfo.RunTime > 1 && maxThreadNum > 1 {
		if dataInfo.RunTime > maxThreadNum {
			loopNum := dataInfo.RunTime/maxThreadNum + 1
			count := 1
			Logger.Info("共执行次数：%v", dataInfo.RunTime)
			for i := 0; i < loopNum; i++ {
				Logger.Info("并发模式-最大执行数:%d,总循环次数:%d,当前循环第%d次", maxThreadNum, loopNum, i+1)
				wg := sync.WaitGroup{}
				for j := 0; j < maxThreadNum; j++ {
					if count > dataInfo.RunTime {
						break
					}
					wg.Add(1)
					go func() {
						result, dst, err1 := RunSceneContent(app, filePath, product, isThread)
						if err1 != nil {
							err = err1
							err = WriteSceneDataResult(id, result, dst, product, envType, err1)
							return
						}
						err = WriteSceneDataResult(id, result, dst, product, envType, err1)
						wg.Done()
					}()
					count++
				}
				wg.Wait()
			}
		} else {
			wg := sync.WaitGroup{}
			Logger.Info("共执行次数：%v", dataInfo.RunTime)
			for i := 0; i < dataInfo.RunTime; i++ {
				wg.Add(1)
				go func() {
					result, dst, err1 := RunSceneContent(app, filePath, product, isThread)
					if err1 != nil {
						err = WriteSceneDataResult(id, result, dst, product, envType, err1)
						return
					}
					err = WriteSceneDataResult(id, result, dst, product, envType, err1)
					wg.Done()
				}()
			}
			wg.Wait()
		}
	} else {
		for i := 0; i < dataInfo.RunTime; i++ {
			if i > 0 {
				Logger.Info("串行模式-执行次数:%d", i+1)
			}
			result, dst, err1 := RunSceneContent(app, filePath, product, isThread)
			if err1 != nil {
				Logger.Error("%s", err1)
				err = err1
			}
			err = WriteSceneDataResult(id, result, dst, product, envType, err1)
			if err != nil {
				Logger.Error("%s", err)
				return
			}

			if result != "pass" {
				err = fmt.Errorf("test %v", result)
				return
			}
		}

	}
	return
}

func RunSceneDataOnce(id, product string) (err error) {
	dataInfo, _, err := GetRunTimeData(id)
	var envType int
	var isThread string
	if len(product) > 0 {
		productTaskInfo, _ := GetProductInfo(product)
		envType = productTaskInfo.EnvType
	} else {
		envType = 2
	}
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if dataInfo.RunTime == 0 {
		err = fmt.Errorf("执行次数为: %d, 如需运行，请把执行次数置为大于0的数", dataInfo.RunTime)
		Logger.Warning("%s", err)
		return
	}
	app, filePath, err := GetFilePath(id)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	result, dst, err1 := RunSceneContent(app, filePath, product, isThread)
	if err1 != nil {
		Logger.Error("%s", err1)
		err = err1
	}
	err = WriteSceneDataResult(id, result, dst, product, envType, err1)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if result != "pass" {
		err = fmt.Errorf("test %v", result)
		return
	}

	return
}

func (dataFile DataFile) SetSleepAction() (err error) {
	if len(dataFile.Action) > 0 {
		for _, item := range dataFile.Action {
			if item.Type == "sleep" {
				valueType := fmt.Sprintf("%T", item.Value)
				var sleepSecond int
				if valueType == "string" {
					sleepSecondStr := item.Value.(string)
					sleepSecond, _ = strconv.Atoi(sleepSecondStr)
				} else {
					sleepSecond = item.Value.(int)
				}
				Logger.Debug("开始Sleep")
				time.Sleep(time.Duration(sleepSecond) * time.Second)
				Logger.Debug("结束Sleep")
			}
		}
	}

	return
}
