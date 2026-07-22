package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func CompareParameterDef(new, old []VarDefModel) (isChanged bool, newList, deletedList, changedList, oldList []VarDefModel) {
	for _, item := range old {
		for index, subItem := range new {
			if item.Name == subItem.Name {
				if item.ValueType != subItem.ValueType || item.Desc != subItem.Desc || item.EgValue != subItem.EgValue {
					changedList = append(changedList, subItem) //记录变更定义
					oldList = append(oldList, item)            // 记录原有定义
				}
				break
			} else if index == len(new)-1 {
				deletedList = append(deletedList, item)
			}
		}
	}

	for _, item := range new {
		for index, subItem := range old {
			if item.Name == subItem.Name {
				break
			} else if index == len(old)-1 {
				newList = append(newList, item)
			}
		}
	}

	if len(newList) > 0 || len(deletedList) > 0 || len(changedList) > 0 {
		isChanged = true
	}

	return
}

func ApiDetailCheck(method, path string, bodyList, pathList, queryList, respList []VarDefModel) (checkTag bool, err error) {
	if method == "get" {
		if len(bodyList) == 0 {
			checkTag = true
		} else {
			if err != nil {
				err = E("error.get_body_param_invalid", err)
			} else {
				err = E("error.get_body_param", method, path)
			}
		}

		tmpDict := make(map[string]interface{})
		for _, item := range respList {
			tmpDict[item.Name] = item.ValueType
		}

		if value, ok := tmpDict["data"]; ok {
			if value == nil {
				checkTag = false
				if err != nil {
					err = E("error.get_missing_data_field_invalid", err)
				} else {
					err = E("error.get_missing_data_field", method, path)
				}
			}
		}
	} else {
		if len(queryList) == 0 || len(pathList) == 0 {
			checkTag = true
		}
	}

	if method == "post" && IsInRouter4Add(path) {
		var respTag bool
		for _, item := range respList {
			isExist := GetRUID(item.Name)
			if isExist {
				respTag = true
				break
			}
		}
		if !respTag {
			checkTag = false
			if err != nil {
				err = E("error.missing_unique_id_invalid", err)
			} else {
				err = E("error.missing_unique_id", method, path)
			}
		}
	}

	if strings.Contains(path, "{") {
		if len(pathList) == 0 {
			checkTag = false
			if err != nil {
				err = E("error.missing_path_variable_invalid", err)
			} else {
				err = E("error.missing_path_variable", method, path)
			}
		} else {
			checkTag = true
		}
	}

	if err != nil {
		Logger.Error("%v", err)
	}

	return
}

func GetApiChkResult(apiCheck ApiCheck) (status string, err error) {
	var envConfig EnvConfig
	status = "WARNING"
	baseUrl := fmt.Sprintf("http://%s/aries/cicd/updateNodeStatus", CICD_HOST)
	nodeId := apiCheck.NodeId
	if len(apiCheck.DevEnvHost) == 0 {
		err = E("error.config_api_doc_path")
		go CallBack4Cicd(baseUrl, nodeId, status)
		return status, err
	}
	models.Orm.Table("env_config").Where("app = ?", apiCheck.App).Find(&envConfig)
	if len(envConfig.App) == 0 {
		envConfig.App = apiCheck.App
		envConfig.Threading = "no"
		envConfig.Testmode = "fuzzing"
		envConfig.Protocol = "http"
		envConfig.SwaggerPath = strings.Replace(SWAGGER_PATH, "{host:port}", apiCheck.DevEnvHost, -1)
		err = models.Orm.Table("env_config").Create(&envConfig).Error
		if err != nil {
			Logger.Error("%s", err)
			go CallBack4Cicd(baseUrl, nodeId, status)
			return status, err
		}
	} else {
		if len(envConfig.SwaggerPath) > 0 {
			if strings.Contains(envConfig.SwaggerPath, "{host:port}") {
				envConfig.SwaggerPath = strings.Replace(envConfig.SwaggerPath, "{host:port}", apiCheck.DevEnvHost, -1)
			} else {
				strReg := regexp.MustCompile(`http(s)?://([\w-.:])+/`)
				strMatch := strReg.FindAllSubmatch([]byte(envConfig.SwaggerPath), -1)
				if len(strMatch) > 0 {
					index, _ := strconv.Atoi(string(strMatch[0][2]))
					indexEnd := len(strMatch[0][0]) - 1
					sIsExsit := string(strMatch[0][1])
					var oldStr string
					if len(sIsExsit) > 0 {
						oldStr = string(strMatch[0][0][index:indexEnd])
					} else {
						oldStr = string(strMatch[0][0][index-1 : indexEnd])
					}
					envConfig.SwaggerPath = strings.Replace(envConfig.SwaggerPath, oldStr, apiCheck.DevEnvHost, -1)
				}
			}
		} else {
			envConfig.SwaggerPath = strings.Replace(SWAGGER_PATH, "{host:port}", apiCheck.DevEnvHost, -1)
		}
		err = models.Orm.Table("env_config").Where("app = ?", apiCheck.App).Update(&envConfig).Error
		if err != nil {
			Logger.Error("%s", err)
			go CallBack4Cicd(baseUrl, nodeId, status)
			return status, err
		}
	}

	if failCount, err := GetSwagger(envConfig.Id); err == nil {
		err = models.Orm.Table("api_definition").Where("app = ? and api_status not in (1,3,4)", envConfig.App).UpdateColumn(&ApiStringDefinition{ApiStatus: "2"}).Error
		if err != nil {
			Logger.Error("%v", err)
			go CallBack4Cicd(baseUrl, nodeId, status)
			return status, err
		}

		var appApiChange AppApiChange
		appApiChange.App = apiCheck.App
		appApiChange.Branch = apiCheck.Branch

		UpdateApiChangeLog(appApiChange)

		if failCount > 0 {
			err = E("error.import_fail_count", failCount)
			go CallBack4Cicd(baseUrl, nodeId, status)
			return status, err
		} else {
			if appApiChange.DeletedApiSum > 0 || appApiChange.ChangedApiSum > 0 {
				status = "FAILED"
			} else if (appApiChange.NewApiSum > 0) && (appApiChange.NewApiSum != appApiChange.CurApiSum) {
				status = "WARNING"
			} else {
				status = "SUCCESSFUL"
			}

			go CallBack4Cicd(baseUrl, nodeId, status)
			return status, err
		}
	} else {
		err = E("error.import_failed", err)
		go CallBack4Cicd(baseUrl, nodeId, status)
		return status, err
	}

	go CallBack4Cicd(baseUrl, nodeId, status)
	return status, err
}

func GetChangedContent(infoType string, newList, deletedList, changedList, oldList []VarDefModel) (changedContent string) {
	changedContent = fmt.Sprintf("%s:", infoType)
	if len(newList) > 0 {
		newByte, _ := json.Marshal(newList)
		newStr := string(newByte)
		changedContent = fmt.Sprintf(T("label.new"), changedContent, newStr)
	}

	if len(deletedList) > 0 {
		deletedByte, _ := json.Marshal(deletedList)
		deletedStr := string(deletedByte)
		changedContent = fmt.Sprintf(T("label.deleted"), changedContent, deletedStr)
	}
	if len(changedList) > 0 {
		changedByte, _ := json.Marshal(changedList)
		changedStr := string(changedByte)

		oldByte, _ := json.Marshal(oldList)
		oldStr := string(oldByte)

		changedContent = fmt.Sprintf(T("label.before_change"), changedContent, oldStr)
		changedContent = fmt.Sprintf(T("label.after_change"), changedContent, changedStr)
	}

	return
}

func UpdateApiChangeLog(appApiChange AppApiChange) {
	models.Orm.Table("api_definition").Where("app = ? and api_status in ('1','3','4','30','31','32','33','34')", appApiChange.App).Count(&appApiChange.CurApiSum)
	models.Orm.Table("api_definition").Where("app = ? and api_status = '1'", appApiChange.App).Count(&appApiChange.NewApiSum)
	models.Orm.Table("api_definition").Where("app = ? and api_status = '2'", appApiChange.App).Count(&appApiChange.DeletedApiSum)
	models.Orm.Table("api_definition").Where("app = ? and api_status in ('3','30','31','32','33','34')", appApiChange.App).Count(&appApiChange.ChangedApiSum)
	models.Orm.Table("api_definition").Where("app = ? and api_status in ('3','4','30','31','32','33','34')", appApiChange.App).Count(&appApiChange.ExistApiSum)
	models.Orm.Table("api_definition").Where("app = ? and `check` = 'fail'and api_status != '2'", appApiChange.App).Count(&appApiChange.CheckFailApiSum)

	var newList, deletedList, changedList, apiCheckFailList []string
	models.Orm.Table("api_definition").Where("app = ? and api_status = '1'", appApiChange.App).Pluck("api_id", &newList)
	models.Orm.Table("api_definition").Where("app = ? and api_status = '2'", appApiChange.App).Pluck("api_id", &deletedList)
	models.Orm.Table("api_definition").Where("app = ? and api_status in ('3','30','31','32','33','34')", appApiChange.App).Pluck("api_id", &changedList)
	models.Orm.Table("api_definition").Where("app = ? and `check` = 'fail' and api_status != '2'", appApiChange.App).Pluck("api_check_fail_reason", &apiCheckFailList)

	if len(newList) > 0 {
		for index, item := range newList {
			if index == 0 {
				appApiChange.NewApiContent = fmt.Sprintf("%s", item)
			} else {
				appApiChange.NewApiContent = fmt.Sprintf("%s\n%s", appApiChange.NewApiContent, item)
			}
		}
		appApiChange.NewApiContent = fmt.Sprintf("<pre><code>%s</code></pre>", appApiChange.NewApiContent)
	}

	if len(deletedList) > 0 {
		for index, item := range deletedList {
			if index == 0 {
				appApiChange.DeletedApiContent = fmt.Sprintf("%s", item)
			} else {
				appApiChange.DeletedApiContent = fmt.Sprintf("%s\n%s", appApiChange.DeletedApiContent, item)
			}
		}
		appApiChange.DeletedApiContent = fmt.Sprintf("<pre><code>%s</code></pre>", appApiChange.DeletedApiContent)
	}

	if len(changedList) > 0 {
		for index, item := range changedList {
			if index == 0 {
				appApiChange.ChangedApiContent = fmt.Sprintf("%s", item)
			} else {
				appApiChange.ChangedApiContent = fmt.Sprintf("%s\n%s", appApiChange.ChangedApiContent, item)
			}
		}
		appApiChange.ChangedApiContent = fmt.Sprintf("<pre><code>%s</code></pre>", appApiChange.ChangedApiContent)
	}

	if len(apiCheckFailList) > 0 {
		for index, item := range apiCheckFailList {
			if index == 0 {
				appApiChange.ApiCheckFailContent = fmt.Sprintf("%s", item)
			} else {
				appApiChange.ApiCheckFailContent = fmt.Sprintf("%s\n%s", appApiChange.ApiCheckFailContent, item)
			}
		}
		appApiChange.ApiCheckFailContent = fmt.Sprintf("<pre><code>%s</code></pre>", appApiChange.ApiCheckFailContent)
	}

	_, counts, _, _ := GetAPISpecCount("app", appApiChange.App)

	allCount := counts[0] + counts[1] + counts[2]
	if appApiChange.CurApiSum > 0 || allCount > 0 {
		appApiChange.ApiCheckResult = fmt.Sprintf(T("label.spec_check"), allCount, counts[0], counts[1], counts[2])
	}

	if appApiChange.NewApiSum > 0 || appApiChange.DeletedApiSum > 0 || appApiChange.ChangedApiSum > 0 || len(appApiChange.ApiCheckResult) > 0 {
		err := models.Orm.Table("app_api_changelog").Create(&appApiChange).Error
		if err != nil {
			Logger.Error("%v", err)
		}
	}
	return
}
