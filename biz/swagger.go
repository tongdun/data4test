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
		definiKeys = make([]string, len(swagger.Definitions))
		for k := range swagger.Definitions {
			definiKeys[j] = k
			j++
		}
	} else if len(swagger.Components.ComSchema) > 0 {
		definiKeys = make([]string, len(swagger.Components.ComSchema))
		for k := range swagger.Components.ComSchema {
			definiKeys[j] = k
			j++
		}
	}

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
		if len(apiDetail.Consumes) > 0 {
			var varHeaderDefModel VarDefModel
			varHeaderDefModel.Name = "Content-Type"
			varHeaderDefModel.ValueType = "string"
			varHeaderDefModel.EgValue = apiDetail.Consumes[0]
			varHeaderDefModel.Desc = "请求类型"
			headerList = append(headerList, varHeaderDefModel)
		}

		for _, pv := range parameters {
			var varDefModel VarDefModel
			var tmpMap []VarDefModel
			if len(pv.Schema.Ref) > 0 {
				depDef := GetLastOne(pv.Schema.Ref)
				tmpMap, _ = allDefini[depDef]
			} else {
				if pv.In == "body" {
					varDefModel.ValueType = pv.Schema.Type
				} else {
					varDefModel.ValueType = pv.Type
				}
				varDefModel.Name = pv.Name
				varDefModel.Desc = pv.Description
				if pv.Required {
					varDefModel.IsMust = "yes"
				} else {
					varDefModel.IsMust = "no"
				}
			}

			if len(tmpMap) > 0 {
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
			} else {
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

			if len(apiDetail.Consumes) == 0 && (len(targetValue) > 0 || len(varDefModel.Name) > 0) {
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
		Logger.Info("从Swagger路径获取接口文档信息: %s", envConfig.SwaggerPath)
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
	var apiIds []string
	for _, v := range pathKeys {
		pathDef = swagger.Paths[v]
		if !reflect.DeepEqual(pathDef.Put, ApiDetail{}) {
			chekTag, errTmp := pathDef.GetApiDetail("put", v, pathDef.Put.Summary, app, allDefine)
			if !chekTag {
				checkFailCount++
			}
			apiId := fmt.Sprintf("put_%s", v)
			apiIds = append(apiIds, apiId)
			if errTmp != nil {
				err = fmt.Errorf("%s,%s", err, errTmp)
			}
		}
		if !reflect.DeepEqual(pathDef.Post, ApiDetail{}) {
			chekTag, errTmp := pathDef.GetApiDetail("post", v, pathDef.Post.Summary, app, allDefine)
			if !chekTag {
				checkFailCount++
			}
			apiId := fmt.Sprintf("put_%s", v)
			apiIds = append(apiIds, apiId)
			if errTmp != nil {
				err = fmt.Errorf("%s,%s", err, errTmp)
			}
		}
		if !reflect.DeepEqual(pathDef.Delete, ApiDetail{}) {
			chekTag, errTmp := pathDef.GetApiDetail("delete", v, pathDef.Delete.Summary, app, allDefine)
			if !chekTag {
				checkFailCount++
			}
			apiId := fmt.Sprintf("put_%s", v)
			apiIds = append(apiIds, apiId)
			if errTmp != nil {
				err = fmt.Errorf("%s,%s", err, errTmp)
			}
		}
		if !reflect.DeepEqual(pathDef.Get, ApiDetail{}) {
			chekTag, errTmp := pathDef.GetApiDetail("get", v, pathDef.Get.Summary, app, allDefine)
			if !chekTag {
				checkFailCount++
			}
			apiId := fmt.Sprintf("put_%s", v)
			apiIds = append(apiIds, apiId)
			if errTmp != nil {
				err = fmt.Errorf("%s,%s", err, errTmp)
			}
		}

	}

	var dbApiIds []string
	models.Orm.Table("api_definition").Where("app = ?", app).Pluck("api_id", &dbApiIds)
	LenOfApiIds := len(apiIds)
	for _, item := range dbApiIds {
		for index, subItem := range apiIds {
			if item == subItem {
				break
			}
			if index == LenOfApiIds-1 {
				var dbApiDef DbApiStringDefinition

				_ = models.Orm.Table("api_definition").Where("app = ? and api_id = ?", app, item).Find(&dbApiDef)
				dbApiDef.ApiStatus = 2
				errTmp := models.Orm.Table("api_definition").Where("id = ?", dbApiDef.Id).Update(&dbApiDef).Error
				if errTmp != nil {
					Logger.Error("%s", errTmp)
				}
				return
			}
		}
	}

	return

}

func (pathDef PathDef) GetApiDetail(method, path, desc, app string, allDefini VarMapList) (checkTag bool, err error) {
	var apiDefinition ApiDefinition
	var apiStrDef ApiStringDefinition
	var dbApiStrDef DbApiStringDefinition
	var apiId string
	var apiDetail ApiDetail
	var dataCount int
	switch method {
	case "put":
		apiDetail = pathDef.Put
		apiId = "put_" + path
		apiDefinition.HttpMethod = "put"
		apiStrDef.HttpMethod = "put"
	case "post":
		apiDetail = pathDef.Post
		apiId = "post_" + path
		apiDefinition.HttpMethod = "post"
		apiStrDef.HttpMethod = "post"
	case "delete":
		apiDetail = pathDef.Delete
		apiId = "delete_" + path
		apiDefinition.HttpMethod = "delete"
		apiStrDef.HttpMethod = "delete"
	case "get":
		apiDetail = pathDef.Get
		apiId = "get_" + path
		apiDefinition.HttpMethod = "get"
		apiStrDef.HttpMethod = "get"
	}
	apiDefinition.ApiId = apiId
	apiStrDef.ApiId = apiId
	if len(apiDetail.Tags) > 0 {
		apiDefinition.ApiModule = apiDetail.Tags[0]
		apiStrDef.ApiModule = apiDetail.Tags[0]
	} else {
		apiDefinition.ApiModule = "other"
		apiStrDef.ApiModule = "other"
	}

	var headerList, bodyList, pathList, queryList, respList []VarDefModel
	headerList, bodyList, pathList, queryList, respList = apiDetail.GetRequestData(allDefini)

	apiDefinition.ApiDesc = desc
	apiStrDef.ApiDesc = desc
	apiDefinition.Path = path
	apiStrDef.Path = path
	apiDefinition.App = app
	apiStrDef.App = app

	mh, _ := json.Marshal(headerList)
	apiStrDef.Header = string(mh)
	apiDefinition.Header = headerList

	mp, _ := json.Marshal(pathList)
	apiStrDef.PathVariable = string(mp)
	apiDefinition.PathVariable = pathList

	mq, _ := json.Marshal(queryList)
	apiStrDef.QueryParameter = string(mq)
	apiDefinition.QueryParameter = queryList

	mb, _ := json.Marshal(bodyList)
	apiStrDef.Body = string(mb)
	apiDefinition.Body = bodyList

	ms, _ := json.Marshal(respList)
	apiStrDef.Response = string(ms)
	apiDefinition.Response = respList

	checkTag, checkResult := ApiDetailCheck(method, path, bodyList, pathList, queryList, respList)

	if !checkTag {
		apiDefinition.Check = "fail"
		apiStrDef.Check = "fail"
		apiStrDef.ApiCheckFailReason = fmt.Sprintf("%v", checkResult)
		apiDefinition.ApiCheckFailReason = fmt.Sprintf("%v", checkResult)
	} else {
		apiDefinition.Check = "pass"
		apiStrDef.Check = "pass"
	}

	models.Orm.Table("api_definition").Where("app = ? and api_id = ?", app, apiId).Find(&dbApiStrDef)
	if len(dbApiStrDef.ApiId) == 0 {
		apiStrDef.ApiStatus = 1
		apiStrDef.Version = 1
		apiStrDef.IsNeedAuto = "1"
		models.Orm.Table("scene_data").Where("app = ? and api_id = ?", app, apiId).Count(&dataCount)
		if dataCount > 0 {
			apiStrDef.IsAuto = "1"
		} else {
			apiStrDef.IsAuto = "-1"
		}
		err = models.Orm.Table("api_definition").Create(&apiStrDef).Error
	} else {
		var oldHeaderList, oldBodyList, oldPathList, oldQueryList, oldRespList []VarDefModel
		if dbApiStrDef.Header != "null" {
			json.Unmarshal([]byte(dbApiStrDef.Header), &oldHeaderList)
		}
		if dbApiStrDef.Body != "null" {
			json.Unmarshal([]byte(dbApiStrDef.Body), &oldBodyList)
		}
		if dbApiStrDef.PathVariable != "null" {
			json.Unmarshal([]byte(dbApiStrDef.PathVariable), &oldPathList)
		}
		if dbApiStrDef.QueryParameter != "null" {
			json.Unmarshal([]byte(dbApiStrDef.QueryParameter), &oldQueryList)
		}
		if dbApiStrDef.Response != "null" {
			json.Unmarshal([]byte(dbApiStrDef.Response), &oldRespList)
		}

		var allChanged, headerChanged, bodyChanged, pathChanged, queryChanged, respChanged string

		if len(headerList) > 0 || len(oldHeaderList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(headerList, oldHeaderList)
			if isChanged {
				//apiStrDef.Check = "fail"  //变更的状态是否置为检查失败，再考虑
				headerChanged = GetChangedContent("Header", newList, deletedList, changedList, oldList)
			}
		}

		if len(bodyList) > 0 || len(oldBodyList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(bodyList, oldBodyList)
			if isChanged {
				//apiStrDef.Check = "fail"   //变更的状态是否置为检查失败，再考虑
				bodyChanged = GetChangedContent("Body", newList, deletedList, changedList, oldList)
			}
		}

		if len(pathList) > 0 || len(oldPathList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(pathList, oldPathList)
			if isChanged {
				//apiStrDef.Check = "fail"   //变更的状态是否置为检查失败，再考虑
				pathChanged = GetChangedContent("Path", newList, deletedList, changedList, oldList)
			}
		}

		if len(queryList) > 0 || len(oldQueryList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(queryList, oldQueryList)
			if isChanged {
				//apiStrDef.Check = "fail"   //变更的状态是否置为检查失败，再考虑
				queryChanged = GetChangedContent("Query", newList, deletedList, changedList, oldList)
			}
		}

		if len(respList) > 0 || len(oldRespList) > 0 {
			isChanged, newList, deletedList, changedList, oldList := CompareParameterDef(respList, oldRespList)
			if isChanged {
				//apiStrDef.Check = "fail"  //变更的状态是否置为检查失败，再考虑
				respChanged = GetChangedContent("Resp", newList, deletedList, changedList, oldList)
			}
		}

		if len(headerChanged) > 0 {
			allChanged = fmt.Sprintf("%s", headerChanged)
		}
		if len(pathChanged) > 0 {
			if len(allChanged) > 0 {
				allChanged = fmt.Sprintf("%s\n%s", allChanged, pathChanged)
			} else {
				allChanged = fmt.Sprintf("%s", pathChanged)
			}
		}
		if len(queryChanged) > 0 {
			if len(allChanged) > 0 {
				allChanged = fmt.Sprintf("%s\n%s", allChanged, queryChanged)
			} else {
				allChanged = fmt.Sprintf("%s", queryChanged)
			}
		}
		if len(bodyChanged) > 0 {
			if len(allChanged) > 0 {
				allChanged = fmt.Sprintf("%s\n%s", allChanged, bodyChanged)
			} else {
				allChanged = fmt.Sprintf("%s", bodyChanged)
			}
		}
		if len(respChanged) > 0 {
			if len(allChanged) > 0 {
				allChanged = fmt.Sprintf("%s\n%s", allChanged, respChanged)
			} else {
				allChanged = fmt.Sprintf("%s", respChanged)
			}

		}
		// 接口状态, 1:新增,2:被删除,3:被修改,4:保持原样
		if len(allChanged) > 0 {
			apiStrDef.ApiStatus = 3
			curTime := time.Now().Format("2006/01/02 15:04:05")
			if len(dbApiStrDef.ChangeContent) > 0 {
				apiStrDef.ChangeContent = fmt.Sprintf("%s\n%s变更记录:\n%s", dbApiStrDef.ChangeContent, curTime, allChanged)
			} else {
				apiStrDef.ChangeContent = fmt.Sprintf("%s变更记录:\n%s", curTime, allChanged)
			}
		} else {
			apiStrDef.ApiStatus = 4
		}

		models.Orm.Table("scene_data").Where("api_id = ?", dbApiStrDef.ApiId).Count(&dataCount)
		if dataCount > 0 {
			apiStrDef.IsAuto = "1"
		}

		apiStrDef.IsNeedAuto = dbApiStrDef.IsNeedAuto
		apiStrDef.Version = dbApiStrDef.Version
		apiStrDef.Remark = dbApiStrDef.Remark

		err = models.Orm.Table("api_definition").Where("id = ?", dbApiStrDef.Id).Update(&apiStrDef).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	//先进行屏蔽
	//var apiRelation ApiRelation
	//apiRelation.ApiId = apiId
	//apiRelation.App = app
	//apiRelation.ApiModule = apiDefinition.ApiModule
	//apiRelation.ApiDesc = apiDefinition.ApiDesc
	//apiRelation.Auto = "yes"
	//curTime := time.Now()
	//apiRelation.CreatedAt = curTime.Format(baseFormat)
	//var dbApiRelation DbApiRelation
	//models.Orm.Table("api_relation").Where("app = ? and api_id = ?", app, apiId).Find(&dbApiRelation)
	//if len(dbApiRelation.ApiId) == 0 {
	//	err = models.Orm.Table("api_relation").Create(&apiRelation).Error
	//} else {
	//	err = models.Orm.Table("api_relation").Where("id = ?", dbApiRelation.Id).Update(&apiRelation).Error
	//}

	return
}
