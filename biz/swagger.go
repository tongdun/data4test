package biz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"data4perf/models"
)

const (
	apiNew = iota
	apiDeleted
	apiChanged
	apiKeep
)

func GetLastOne(item string) (value string) {
	tmps := strings.Split(item, "/")
	value = tmps[len(tmps)-1]
	return
}

func (swagger Swagger) GetAllDefinition() (defAllDict VarMapList) {
	defAllDict = make(VarMapList)
	var definiKeys []string
	j := 0
	if len(swagger.Definitions) > 0 {
		for _, v := range definiKeys {
			var paramList []VarDefModel
			subProperty := swagger.Definitions[v].Properties

			subProKeys := make([]string, len(subProperty))
			j = 0
			for k := range subProperty {
				subProKeys[j] = k
				j++
			}

			for _, sv := range subProKeys {
				var varDef VarDefModel
				varDef.ValueType = subProperty[sv].Type
				varDef.Desc = subProperty[sv].Description
				varDef.Name = sv
				if len(varDef.ValueType) == 0 && len(subProperty[sv].Ref) > 0 {
					resDef := GetLastOne(subProperty[sv].Ref)
					respList, _ := defAllDict[resDef]
					varDef.ValueType = "object"
					mapTemp := make(map[string]string)
					for _, item := range respList {
						mapTemp[item.Name] = item.ValueType
					}
					if len(mapTemp) != 0 {
						mapByte, errTmp := json.Marshal(mapTemp)
						if errTmp != nil {
							Logger.Error("%s", errTmp)
						}
						varDef.EgValue = string(mapByte)
					}
				}

				paramList = append(paramList, varDef)
			}

			defAllDict[v] = paramList
		}
	} else if len(swagger.Components.ComSchema) > 0 {
		definiKeys = make([]string, len(swagger.Components.ComSchema))
		for k := range swagger.Components.ComSchema {
			definiKeys[j] = k
			j++
		}
	}

	if len(swagger.Definitions) > 0 {
		definiKeys = make([]string, len(swagger.Definitions))
		for k := range swagger.Definitions {
			definiKeys[j] = k
			j++
		}
	} else if len(swagger.Components.ComSchema) > 0 {
		for _, v := range definiKeys {
			var paramList []VarDefModel
			subProperty := swagger.Components.ComSchema[v].Properties

			subProKeys := make([]string, len(subProperty))
			j = 0
			for k := range subProperty {
				subProKeys[j] = k
				j++
			}

			for _, sv := range subProKeys {
				var varDef VarDefModel
				varDef.ValueType = subProperty[sv].Type
				varDef.Desc = subProperty[sv].Description
				varDef.Name = sv
				if len(varDef.ValueType) == 0 && len(subProperty[sv].Ref) > 0 {
					resDef := GetLastOne(subProperty[sv].Ref)
					respList, _ := defAllDict[resDef]
					//Logger.Debug("resDef: %s", resDef)
					//Logger.Debug("respList: %#v", respList)
					varDef.ValueType = "object"
					mapTemp := make(map[string]string)
					for _, item := range respList {
						mapTemp[item.Name] = item.ValueType
					}
					if len(mapTemp) != 0 {
						mapByte, errTmp := json.Marshal(mapTemp)
						if errTmp != nil {
							Logger.Error("%s", errTmp)
						}
						varDef.EgValue = string(mapByte)
					}
				}

				paramList = append(paramList, varDef)
			}

			defAllDict[v] = paramList
		}
	}

	return

}

func (apiDetail ApiDetail) GetRequestData(allDefini VarMapList) (headerList, bodyList, pathList, queryList, respList []VarDefModel) {
	if len(apiDetail.Responses.R200.Schema.Ref) > 0 {
		resDef := GetLastOne(apiDetail.Responses.R200.Schema.Ref)
		respList, _ = allDefini[resDef]
		parameters := apiDetail.Parameters
		for _, pv := range parameters {
			var varDefModel VarDefModel
			var targetValue string
			var tmpMap []VarDefModel
			if len(pv.Schema.Ref) > 0 {
				depDef := GetLastOne(pv.Schema.Ref)
				tmpMap, _ = allDefini[depDef]
			} else {
				if pv.In == "body" {
					targetValue = pv.Schema.Type
				} else {
					targetValue = pv.Type
				}
			}

			varDefModel.Name = pv.Name
			varDefModel.ValueType = pv.Type
			varDefModel.Desc = pv.Description
			if pv.Required {
				varDefModel.IsMust = "yes"
			} else {
				varDefModel.IsMust = "no"
			}

			if len(targetValue) > 0 || len(varDefModel.Name) > 0 {
				switch pv.In {
				case "header":
					headerList = append(headerList, varDefModel)
				case "query":
					queryList = append(queryList, varDefModel)
				case "path":
					pathList = append(pathList, varDefModel)
				case "body", "formData":
					bodyList = append(bodyList, varDefModel)
				}
			} else {
				switch pv.In {
				case "header":
					headerList = tmpMap
				case "query":
					queryList = tmpMap
				case "path":
					pathList = tmpMap
				case "body", "formData":
					bodyList = tmpMap
				}
			}

		}
	} else if len(apiDetail.Responses.R200.Content.Star.Schema.Ref) > 0 {
		resDef := GetLastOne(apiDetail.Responses.R200.Content.Star.Schema.Ref)
		respList, _ = allDefini[resDef]
		if len(apiDetail.RequestBody.RequestContent.JSONContent.Schema.Ref) > 0 {
			bodyDef := GetLastOne(apiDetail.RequestBody.RequestContent.JSONContent.Schema.Ref)
			bodyList, _ = allDefini[bodyDef]
		} else if apiDetail.RequestBody.RequestContent.MultipartContent.Schema.Properties != nil {
			for pk, pv := range apiDetail.RequestBody.RequestContent.MultipartContent.Schema.Properties {
				var varDefModel VarDefModel
				varDefModel.Name = pk
				varDefModel.ValueType = pv.Type
				varDefModel.Desc = pv.Description
				bodyList = append(bodyList, varDefModel)
			}
		}

		parameters := apiDetail.Parameters
		for _, pv := range parameters {
			var varDefModel VarDefModel
			var targetValue string
			var tmpMap []VarDefModel
			if len(pv.Schema.Ref) > 0 {
				depDef := GetLastOne(pv.Schema.Ref)
				tmpMap, _ = allDefini[depDef]
			} else {
				if pv.In == "body" {
					targetValue = pv.Schema.Type
				} else {
					targetValue = pv.Type
				}
			}

			varDefModel.Name = pv.Name
			varDefModel.ValueType = pv.Type
			varDefModel.Desc = pv.Description
			if pv.Required {
				varDefModel.IsMust = "yes"
			} else {
				varDefModel.IsMust = "no"
			}

			if len(targetValue) > 0 || len(varDefModel.Name) > 0 {
				switch pv.In {
				case "header":
					headerList = append(headerList, varDefModel)
				case "query":
					queryList = append(queryList, varDefModel)
				case "path":
					pathList = append(pathList, varDefModel)
				case "body", "formData":
					bodyList = append(bodyList, varDefModel)
				}
			} else {
				switch pv.In {
				case "header":
					headerList = tmpMap
				case "query":
					queryList = tmpMap
				case "path":
					pathList = tmpMap
				case "body", "formData":
					bodyList = tmpMap
				}
			}

		}
	}

	return
}

func GetSwagger(id string) (checkFailCount int, err error) {
	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("id = ?", id).Find(&envConfig)
	if len(envConfig.App) == 0 {
		err = fmt.Errorf("Not found related app id:%s", id)
		return
	}
	app := envConfig.App

	var content []byte
	if len(strings.TrimSpace(envConfig.SwaggerPath)) == 0 {
		cmdStr := fmt.Sprintf("ls -lt %s/%s*.json | head -n 1 | awk '{print $9}'", ApiFilePath, app)

		fileName, errTmp := ExecCommand(cmdStr)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			err = errTmp
			return
		}

		Logger.Info("从%s文件获取接口文档信息", fileName)
		content, errTmp = ioutil.ReadFile(fileName)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			err = errTmp
			return
		}
	} else {
		Logger.Info("从Swagger路径获取接口文档信息: %v", envConfig.SwaggerPath)
		for i := 1; i < 30; i++ {
			resp, errTemp := http.Get(envConfig.SwaggerPath)
			Logger.Info("第%d次获取Swagger接口信息", i)
			if errTemp != nil {
				if i == 29 {
					Logger.Error("%s", errTemp)
					return 0, errTemp
				} else {
					time.Sleep(100000)
				}
			} else {
				// 一次性读取
				content, errTemp = ioutil.ReadAll(resp.Body)
				if errTemp != nil {
					Logger.Error("%s", errTemp)
					return 0, errTemp
				}
				resp.Body.Close()
				break
			}
		}

		var fileNameOld string

		cmdStr := fmt.Sprintf("ls -lt %s/%s*.json | head -n 1 | awk '{print $9}'", ApiFilePath, app)
		fileNameOld, err = ExecCommand(cmdStr)
		if err != nil {
			goto nextHere
		} else if len(fileNameOld) == 0 {
			curTime := time.Now().Format(timeFormatWithNoSpecial)
			fileNameNew := fmt.Sprintf("%s/%s_%s.json", ApiFilePath, app, curTime)
			errTmp := ioutil.WriteFile(fileNameNew, content, 0644)
			if errTmp != nil {
				Logger.Error("%s", errTmp)
				goto nextHere
			}
		} else {
			curTime := time.Now().Format(timeFormatWithNoSpecial)
			fileNameNew := fmt.Sprintf("%s/%s_%s.json", ApiFilePath, app, curTime)
			errTmp := WriteJson(content, fileNameNew)

			if errTmp != nil {
				goto nextHere
			}
			cmdStr = fmt.Sprintf("diff %s %s", fileNameOld, fileNameNew)
			diffContent, errTmp := ExecCommand(cmdStr)
			if errTmp != nil {
				goto nextHere
			}
			if len(diffContent) > 0 {
				Logger.Warning("%s", diffContent)
			} else {
				// 新旧文件无差异，把新文件实时进行删除
				cmdStr := fmt.Sprintf("rm -f %s", fileNameNew)
				_, _ = ExecCommand(cmdStr)
			}
		}
	}

nextHere:
	var swagger Swagger
	err = json.Unmarshal([]byte(content), &swagger)
	allDefine := swagger.GetAllDefinition()

	pathKeys := make([]string, len(swagger.Paths))
	j := 0
	for k := range swagger.Paths {
		pathKeys[j] = k
		j++
	}
	var pathDef PathDef
	for _, v := range pathKeys {
		pathDef = swagger.Paths[v]
		if !reflect.DeepEqual(pathDef.Put, ApiDetail{}) {
			chekTag, errTmp := pathDef.GetApiDetail("put", v, pathDef.Put.Summary, app, allDefine)
			if !chekTag {
				checkFailCount++
			}
			if errTmp != nil {
				err = fmt.Errorf("%s,%s", err, errTmp)
			}
		}
		if !reflect.DeepEqual(pathDef.Post, ApiDetail{}) {
			chekTag, errTmp := pathDef.GetApiDetail("post", v, pathDef.Post.Summary, app, allDefine)
			if !chekTag {
				checkFailCount++
			}
			if errTmp != nil {
				err = fmt.Errorf("%s,%s", err, errTmp)
			}
		}
		if !reflect.DeepEqual(pathDef.Delete, ApiDetail{}) {
			chekTag, errTmp := pathDef.GetApiDetail("delete", v, pathDef.Delete.Summary, app, allDefine)
			if !chekTag {
				checkFailCount++
			}
			if errTmp != nil {
				err = fmt.Errorf("%s,%s", err, errTmp)
			}
		}
		if !reflect.DeepEqual(pathDef.Get, ApiDetail{}) {
			chekTag, errTmp := pathDef.GetApiDetail("get", v, pathDef.Get.Summary, app, allDefine)
			if !chekTag {
				checkFailCount++
			}
			if errTmp != nil {
				err = fmt.Errorf("%s,%s", err, errTmp)
			}
		}

	}

	return

}

func (pathDef PathDef) GetApiDetail(method, path, desc, app string, allDefini VarMapList) (checkTag bool, err error) {
	var apiDefinition ApiDefinition
	var apiStringDefinition ApiStringDefinition
	var dbApiStringDefinition DbApiStringDefinition
	var apiId string
	var apiDetail ApiDetail
	//var sceneData SceneData
	var dataCount int
	switch method {
	case "put":
		apiDetail = pathDef.Put
		apiId = "put_" + path
		apiDefinition.HttpMethod = "put"
		apiStringDefinition.HttpMethod = "put"
	case "post":
		apiDetail = pathDef.Post
		apiId = "post_" + path
		apiDefinition.HttpMethod = "post"
		apiStringDefinition.HttpMethod = "post"
	case "delete":
		apiDetail = pathDef.Delete
		apiId = "delete_" + path
		apiDefinition.HttpMethod = "delete"
		apiStringDefinition.HttpMethod = "delete"
	case "get":
		apiDetail = pathDef.Get
		apiId = "get_" + path
		apiDefinition.HttpMethod = "get"
		apiStringDefinition.HttpMethod = "get"
	}
	apiDefinition.ApiId = apiId
	apiStringDefinition.ApiId = apiId
	if len(apiDetail.Tags) > 0 {
		apiDefinition.ApiModule = apiDetail.Tags[0]
		apiStringDefinition.ApiModule = apiDetail.Tags[0]
	} else {
		apiDefinition.ApiModule = "other"
		apiStringDefinition.ApiModule = "other"
	}

	var headerList, bodyList, pathList, queryList, respList []VarDefModel
	headerList, bodyList, pathList, queryList, respList = apiDetail.GetRequestData(allDefini)

	apiDefinition.ApiDesc = desc
	apiStringDefinition.ApiDesc = desc
	apiDefinition.Path = path
	apiStringDefinition.Path = path
	apiDefinition.App = app
	apiStringDefinition.App = app

	mh, _ := json.Marshal(headerList)
	apiStringDefinition.Header = string(mh)
	apiDefinition.Header = headerList

	mp, _ := json.Marshal(pathList)
	apiStringDefinition.PathVariable = string(mp)
	apiDefinition.PathVariable = pathList

	mq, _ := json.Marshal(queryList)
	apiStringDefinition.QueryParameter = string(mq)
	apiDefinition.QueryParameter = queryList

	mb, _ := json.Marshal(bodyList)
	apiStringDefinition.Body = string(mb)
	apiDefinition.Body = bodyList

	ms, _ := json.Marshal(respList)
	apiStringDefinition.Response = string(ms)
	apiDefinition.Response = respList

	checkTag, checkResult := ApiDetailCheck(method, path, bodyList, pathList, queryList, respList)

	if !checkTag {
		apiDefinition.Check = "fail"
		apiStringDefinition.Check = "fail"
		apiStringDefinition.ApiCheckFailReason = fmt.Sprintf("%v", checkResult)
		apiDefinition.ApiCheckFailReason = fmt.Sprintf("%v", checkResult)
	} else {
		apiDefinition.Check = "pass"
		apiStringDefinition.Check = "pass"
	}

	models.Orm.Table("api_definition").Where("app = ? and api_id = ?", app, apiId).Find(&dbApiStringDefinition)
	if len(dbApiStringDefinition.ApiId) == 0 {
		apiStringDefinition.ApiStatus = 1
		apiStringDefinition.Version = 1
		apiStringDefinition.IsAuto = 0
		apiStringDefinition.IsNeedAuto = 1
		err = models.Orm.Table("api_definition").Create(&apiStringDefinition).Error
	} else {
		var oldHeaderList, oldBodyList, oldPathList, oldQueryList, oldRespList []VarDefModel
		if dbApiStringDefinition.Header != "null" {
			json.Unmarshal([]byte(dbApiStringDefinition.Header), &oldHeaderList)
		}
		if dbApiStringDefinition.Body != "null" {
			json.Unmarshal([]byte(dbApiStringDefinition.Body), &oldBodyList)
		}
		if dbApiStringDefinition.PathVariable != "null" {
			json.Unmarshal([]byte(dbApiStringDefinition.PathVariable), &oldPathList)
		}
		if dbApiStringDefinition.QueryParameter != "null" {
			json.Unmarshal([]byte(dbApiStringDefinition.QueryParameter), &oldQueryList)
		}
		if dbApiStringDefinition.Response != "null" {
			json.Unmarshal([]byte(dbApiStringDefinition.Response), &oldRespList)
		}

		var allChanged, headerChanged, bodyChanged, pathChanged, queryChanged, respChanged string

		if len(headerList) > 0 || len(oldHeaderList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(headerList, oldHeaderList)
			if isChanged {
				apiStringDefinition.Check = "fail"
				headerChanged = GetChangedContent("Header", newList, deletedList, changedList, oldList)
			}
		}

		if len(bodyList) > 0 || len(oldBodyList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(headerList, oldHeaderList)
			if isChanged {
				apiStringDefinition.Check = "fail"
				bodyChanged = GetChangedContent("Body", newList, deletedList, changedList, oldList)
			}
		}

		if len(pathList) > 0 || len(oldPathList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(headerList, oldHeaderList)
			if isChanged {
				apiStringDefinition.Check = "fail"
				pathChanged = GetChangedContent("Path", newList, deletedList, changedList, oldList)
			}
		}

		if len(queryList) > 0 || len(oldQueryList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(headerList, oldHeaderList)
			if isChanged {
				apiStringDefinition.Check = "fail"
				queryChanged = GetChangedContent("Query", newList, deletedList, changedList, oldList)
			}
		}

		if len(respList) > 0 || len(oldRespList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(headerList, oldHeaderList)
			if isChanged {
				apiStringDefinition.Check = "fail"
				respChanged = GetChangedContent("Resp", newList, deletedList, changedList, oldList)
			}
		}

		if len(headerChanged) > 0 {
			allChanged = fmt.Sprintf("%s\n%s\n", allChanged, headerChanged)
		}
		if len(bodyChanged) > 0 {
			allChanged = fmt.Sprintf("%s\n%s\n", allChanged, bodyChanged)
		}
		if len(pathChanged) > 0 {
			allChanged = fmt.Sprintf("%s\n%s\n", allChanged, pathChanged)
		}
		if len(queryChanged) > 0 {
			allChanged = fmt.Sprintf("%s\n%s\n", allChanged, queryChanged)
		}

		if len(respChanged) > 0 {
			allChanged = fmt.Sprintf("%s\n%s\n", allChanged, respChanged)
		}

		if len(allChanged) > 0 {
			err = models.Orm.Table("api_definition").Where("id = ?", dbApiStringDefinition.Id).UpdateColumn(&ApiStringDefinition{ApiStatus: 2, ChangeContent: allChanged}).Error
			if err != nil {
				Logger.Error("%v", err)
			}
		} else {
			apiStringDefinition.ApiStatus = 4
			err = models.Orm.Table("api_definition").Where("id = ?", dbApiStringDefinition.Id).Update(&apiStringDefinition.ApiStatus).Error
		}

		models.Orm.Table("scene_data").Where("id = ?", dbApiStringDefinition.ApiId).Count(&dataCount)
		if dataCount > 0 {
			apiStringDefinition.IsAuto = 1
			err = models.Orm.Table("api_definition").Where("id = ?", dbApiStringDefinition.Id).Update(&apiStringDefinition.IsAuto).Error
		}
	}

	var apiRelation ApiRelation
	apiRelation.ApiId = apiId
	apiRelation.App = app
	apiRelation.ApiModule = apiDefinition.ApiModule
	apiRelation.ApiDesc = apiDefinition.ApiDesc
	apiRelation.Auto = "yes"
	curTime := time.Now()
	apiRelation.CreatedAt = curTime.Format(baseFormat)
	var dbApiRelation DbApiRelation
	models.Orm.Table("api_relation").Where("app = ? and api_id = ?", app, apiId).Find(&dbApiRelation)
	if len(dbApiRelation.ApiId) == 0 {
		err = models.Orm.Table("api_relation").Create(&apiRelation).Error
	} else {
		err = models.Orm.Table("api_relation").Where("id = ?", dbApiRelation.Id).Update(&apiRelation).Error
	}

	return
}
