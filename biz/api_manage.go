package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	//"regexp"
	//"strconv"
	"strings"
	//"sync"
)

func GetAppList() (appList []string) {
	models.Orm.Table("env_config").Order("created_at desc").Pluck("app", &appList)
	return
}

func GetProductList() (productList []string) {
	models.Orm.Table("product").Order("created_at desc").Group("product").Pluck("product", &productList)
	return
}

func GetHistoryDateList() (dateList []string) {
	models.Orm.Table("scene_data_test_history").Select("DATE(created_at)").Group("DATE(created_at)").Order("DATE(created_at) desc").Limit(30).Pluck("DATE(created_at)", &dateList)
	return
}

func GetSceneHistoryDateList() (dateList []string) {
	models.Orm.Table("scene_test_history").Select("DATE(created_at)").Group("DATE(created_at)").Order("DATE(created_at) desc").Limit(30).Pluck("DATE(created_at)", &dateList)
	return
}

func GetDataList() (dataList []string) {
	// 控制台控制列表和数据列表只关联标准类型的数据用例
	models.Orm.Table("scene_data").Group("name").Where("file_type = 1").Order("created_at desc").Pluck("name", &dataList)
	return
}

func GetDataFileList() (dataFileList []DepDataModel) {
	var fileTmp []string
	var depDataFile DepDataModel
	models.Orm.Table("scene_data").Order("created_at desc").Pluck("file_name", &fileTmp)
	for _, item := range fileTmp {
		depDataFile.DataFile = item
		dataFileList = append(dataFileList, depDataFile)
	}
	return
}

func GetEnvList() (envList []EnvModel) {
	var fileTmp []string
	var envModel EnvModel
	models.Orm.Table("product").Order("created_at desc").Pluck("product", &fileTmp)
	for _, item := range fileTmp {
		envModel.Product = item
		envList = append(envList, envModel)
	}
	return
}

func GetAppInfo(appName string) (appModel AppModel) {
	appModel.AppName = appName
	var prefixList []string
	models.Orm.Table("env_config").Select("prepath").Where("app = ?", appName).Limit(1).Pluck("prepath", &prefixList)
	if len(prefixList) > 0 {
		appModel.Prefix = prefixList[0]
	}
	models.Orm.Table("api_definition").Group("http_method").Where("app = ?", appName).Pluck("http_method", &appModel.Methods)

	models.Orm.Table("api_definition").Group("api_module").Where("app = ?", appName).Pluck("api_module", &appModel.Modules)

	models.Orm.Table("api_definition").Group("path").Where("app = ?", appName).Pluck("path", &appModel.Apis)

	models.Orm.Table("api_definition").Group("api_desc").Where("app = ?", appName).Pluck("api_desc", &appModel.ApisDesc)

	models.Orm.Table("scene_data").Order("created_at desc").Where("app = ?", appName).Pluck("name", &appModel.DatasDesc)

	return
}

func GetApiInfo(appName, method, path string) (apiModel ApiInfoModel) {
	apiModel.App = appName
	apiModel.Method = method
	apiModel.Path = path

	var prefixList []string
	models.Orm.Table("env_config").Select("prepath").Where("app = ?", appName).Pluck("prepath", &prefixList)
	if len(prefixList) > 0 {
		apiModel.Prefix = prefixList[0]
	}

	apiId := fmt.Sprintf("%s_%s", method, path)
	models.Orm.Table("scene_data").Group("name").Select("name").Where("app = ? and api_id = ?", appName, apiId).Pluck("name", &apiModel.DatasDesc)

	var apiDef ApiDefinition
	models.Orm.Table("api_definition").Where("app = ? and http_method = ? and path = ?", appName, method, path).Find(&apiDef)
	apiModel.ApiDesc = apiDef.ApiDesc
	apiModel.Module = apiDef.ApiModule
	apiModel.PathVars = apiDef.PathVariable
	apiModel.QueryVars = apiDef.QueryParameter
	apiModel.BodyVars = apiDef.Body
	apiModel.HeaderVars = apiDef.Header

	return
}

func GetApiDataByApiDesc(appName, module, apiDesc string) (apiModel ApiInfoModel) {
	apiModel.App = appName
	apiModel.Module = module
	apiModel.ApiDesc = apiDesc

	var apiDef ApiDefDB
	models.Orm.Table("api_definition").Where("app = ? and api_module = ? and api_desc = ?", appName, module, apiDesc).Find(&apiDef)
	apiModel.Method = apiDef.HttpMethod
	apiModel.Path = apiDef.Path

	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", appName).Find(&envConfig)

	apiModel.Prefix = envConfig.Prepath
	apiId := fmt.Sprintf("%s_%s", apiDef.HttpMethod, apiDef.Path)
	models.Orm.Table("scene_data").Group("name").Select("name").Where("app = ? and api_id = ?", appName, apiId).Pluck("name", &apiModel.DatasDesc)

	json.Unmarshal([]byte(apiDef.PathVariable), &apiModel.PathVars)
	json.Unmarshal([]byte(apiDef.QueryParameter), &apiModel.QueryVars)
	json.Unmarshal([]byte(apiDef.Body), &apiModel.BodyVars)
	json.Unmarshal([]byte(apiDef.Header), &apiModel.HeaderVars)

	for _, value := range apiModel.HeaderVars {
		if value.Name == "Content-Type" {
			apiModel.BodyMode = value.EgValue
			break
		}
	}

	json.Unmarshal([]byte(apiDef.Response), &apiModel.RespVars)

	return
}

func GetModuleInfo(appName, module string) (moduleModel ModuleModel) {
	moduleModel.AppName = appName
	moduleModel.Module = module

	models.Orm.Table("api_definition").Group("http_method").Select("http_method").Where("app = ? and api_module = ?", appName, module).Pluck("http_method", &moduleModel.Methods)
	models.Orm.Table("api_definition").Group("path").Select("path").Where("app = ? and api_module = ?", appName, module).Pluck("path", &moduleModel.Apis)
	models.Orm.Table("api_definition").Group("api_desc").Select("api_desc").Where("app = ? and api_module = ?", appName, module).Pluck("api_desc", &moduleModel.ApisDesc)

	return
}

func GetMethodInfo(appName, method string) (methodModel MethodModel) {
	methodModel.AppName = appName
	methodModel.Method = method

	models.Orm.Table("api_definition").Group("api_module").Select("api_module").Where("app = ? and http_method = ?", appName, method).Pluck("api_module", &methodModel.Modules)
	models.Orm.Table("api_definition").Group("path").Select("path").Where("app = ? and  http_method = ?", appName, method).Pluck("path", &methodModel.Apis)
	models.Orm.Table("api_definition").Group("api_desc").Select("api_desc").Where("app = ? and http_method = ?", appName, method).Pluck("api_desc", &methodModel.ApisDesc)

	return
}

func SaveApiDef(apiModel ApiDefSaveModel) (err error) {
	appIsExist := 0
	models.Orm.Table("env_config").Where("app = ?", apiModel.App).Count(&appIsExist)
	if appIsExist == 0 {
		var envConfig EnvConfig
		envConfig.App = apiModel.App
		envConfig.Prepath = apiModel.Prefix
		envConfig.Protocol = "http"
		envConfig.Threading = "no"
		envConfig.Testmode = "custom"
		err = models.Orm.Table("env_config").Create(&envConfig).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	var dbApiDef, apiDef ApiDefDB

	apiDef.App = apiModel.App

	apiDef.HttpMethod = apiModel.Method
	apiDef.Path = apiModel.Path
	apiDef.ApiId = fmt.Sprintf("%s_%s", apiModel.Method, apiModel.Path)
	apiDef.ApiModule = apiModel.Module
	apiDef.ApiDesc = apiModel.ApiDesc

	pathVarStr, _ := json.Marshal(apiModel.PathVars)
	apiDef.PathVariable = string(pathVarStr)

	queryVarStr, _ := json.Marshal(apiModel.QueryVars)
	apiDef.QueryParameter = string(queryVarStr)

	if len(apiModel.BodyStr) > 0 {
		apiModel.BodyVars, err = Str2DefModel(apiModel.BodyStr)
	}

	bodyVarStr, _ := json.Marshal(apiModel.BodyVars)
	apiDef.Body = string(bodyVarStr)

	if len(apiModel.BodyMode) > 0 && len(apiModel.HeaderVars) == 0 {
		var tmpVar VarDefModel
		tmpVar.Name = "Content-Type"
		if apiModel.BodyMode == "json" || apiModel.BodyMode == "form-data" || apiModel.BodyMode == "x-www-form-urlencoded" {
			tmpVar.EgValue = fmt.Sprintf("application/%s", apiModel.BodyMode)
		} else {
			tmpVar.EgValue = apiModel.BodyMode
		}

		apiModel.HeaderVars = append(apiModel.HeaderVars, tmpVar)
	}

	headerVarStr, _ := json.Marshal(apiModel.HeaderVars)
	apiDef.Header = string(headerVarStr)

	respVarStr, _ := json.Marshal(apiModel.RespVars)
	apiDef.Response = string(respVarStr)

	models.Orm.Table("api_definition").Where("app = ? and http_method = ? and path = ?", apiModel.App, apiModel.Method, apiModel.Path).Find(&dbApiDef)

	checkTag, errTmp := ApiDetailCheck(apiDef.HttpMethod, apiDef.Path, apiModel.BodyVars, apiModel.PathVars, apiModel.QueryVars, apiModel.RespVars)
	if !checkTag {
		apiDef.Check = "pass"
	} else {
		apiDef.Check = "fail"
	}

	if len(dbApiDef.ApiId) > 0 {
		apiDef.Version = dbApiDef.Version + 1
		if errTmp != nil {
			apiDef.ChangeContent = fmt.Sprintf("%s\n%v", dbApiDef.ChangeContent, errTmp)
		}
		err = models.Orm.Table("api_definition").Where("api_id = ? and app = ?", apiDef.ApiId, apiDef.App).Update(&apiDef).Error
	} else {
		apiDef.Version = 1
		if errTmp != nil {
			apiDef.ChangeContent = fmt.Sprintf("%v", errTmp)
		}
		err = models.Orm.Table("api_definition").Create(&apiDef).Error
	}

	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func SaveApiData(apiModel ApiDataSaveModel) (err error) {
	var sceneData SceneData
	sceneData.App = apiModel.App
	sceneData.Name = apiModel.ApiDesc
	sceneData.ApiId = fmt.Sprintf("%s_%s", apiModel.Method, apiModel.Path)

	var fullName string
	tagNum := strings.Count(apiModel.DataDesc, "-")
	if tagNum >= 2 {
		fullName = apiModel.DataDesc
	} else {
		tagNum = strings.Count(apiModel.DataDesc, "_")
		if tagNum >= 2 {
			fullName = apiModel.DataDesc
		} else {
			fullName = fmt.Sprintf("%s-%s-%s", apiModel.Module, apiModel.ApiDesc, apiModel.DataDesc)
		}
	}

	sceneData.Name = fullName
	if apiModel.BodyMode == "json" {
		sceneData.FileName = fmt.Sprintf("%s.json", sceneData.Name)
	} else {
		sceneData.FileName = fmt.Sprintf("%s.yml", sceneData.Name)
	}

	sceneData.RunTime = 1

	var dataFile DataFile
	dataFile.Name = sceneData.Name

	if len(apiModel.Other) > 0 {
		dataFile.Version = apiModel.Other[0].Version
		dataFile.IsParallel = apiModel.Other[0].IsParallel
		dataFile.IsUseEnvConfig = apiModel.Other[0].IsUseEnvConfig
	} else {
		dataFile.Version = 1.0
		dataFile.IsParallel = "no"
		dataFile.IsUseEnvConfig = "yes"
	}
	dataFile.IsRunPostApis = "no"
	dataFile.IsRunPreApis = "no"
	dataFile.ApiId = sceneData.ApiId

	dataFile.Env.Prepath = apiModel.Prefix
	dataFile.Api.App = apiModel.App
	dataFile.Api.Description = apiModel.ApiDesc
	dataFile.Api.Module = apiModel.Module
	dataFile.Api.Method = apiModel.Method
	dataFile.Api.Path = apiModel.Path
	dataFile.Assert = apiModel.Asserts
	dataFile.Action = apiModel.Actions

	if len(apiModel.PreApis) > 0 {
		for _, item := range apiModel.PreApis {
			dataFile.Api.PreApi = append(dataFile.Api.PreApi, item.DataFile)
		}

	}

	if len(apiModel.PostApis) > 0 {
		for _, item := range apiModel.PostApis {
			dataFile.Api.PostApis = append(dataFile.Api.PostApis, item.DataFile)
		}

	}

	if len(apiModel.HeaderVars) > 0 {
		mapData := make(map[string]interface{})
		for _, item := range apiModel.HeaderVars {
			if len(item.TestValue) >= 1 {
				mapData[item.Name] = item.TestValue[0]
			} else {
				mapData[item.Name] = ""
			}
		}
		dataFile.Single.Header = mapData
	}

	if len(apiModel.BodyMode) > 0 {
		if dataFile.Single.Header == nil {
			mapData := make(map[string]interface{})
			if apiModel.BodyMode == "json" || apiModel.BodyMode == "form-data" || apiModel.BodyMode == "x-www-form-urlencoded" {
				mapData["Content-Type"] = fmt.Sprintf("application/%s", apiModel.BodyMode)
			} else {
				mapData["Content-Type"] = apiModel.BodyMode
			}
			dataFile.Single.Header = mapData
		} else {
			dataFile.Single.Header["Content-Type"] = apiModel.BodyMode
			if apiModel.BodyMode == "json" || apiModel.BodyMode == "form-data" || apiModel.BodyMode == "x-www-form-urlencoded" {
				dataFile.Single.Header["Content-Type"] = fmt.Sprintf("application/%s", apiModel.BodyMode)
			} else {
				dataFile.Single.Header["Content-Type"] = apiModel.BodyMode
			}
		}

	}

	if len(apiModel.QueryVars) > 0 {
		mapData := make(map[string]interface{})
		multiMapData := make(map[string][]interface{})
		for _, item := range apiModel.QueryVars {
			if len(item.TestValue) == 1 {
				mapData[item.Name] = item.TestValue[0]
			} else if len(item.TestValue) > 1 {
				for _, subItem := range item.TestValue {
					multiMapData[item.Name] = append(multiMapData[item.Name], subItem)
				}
			} else {
				mapData[item.Name] = ""
			}
		}
		dataFile.Single.Query = mapData
		dataFile.Multi.Query = multiMapData
	}

	if len(apiModel.PathVars) > 0 {
		mapData := make(map[string]interface{})
		multiMapData := make(map[string][]interface{})
		for _, item := range apiModel.PathVars {
			if len(item.TestValue) == 1 {
				mapData[item.Name] = item.TestValue[0]
			} else if len(item.TestValue) > 1 {
				for _, subItem := range item.TestValue {
					multiMapData[item.Name] = append(multiMapData[item.Name], subItem)
				}
			} else {
				mapData[item.Name] = ""
			}
		}
		dataFile.Single.Path = mapData
		dataFile.Multi.Path = multiMapData
	}

	if len(apiModel.BodyVars) > 0 {
		mapData := make(map[string]interface{})
		multiMapData := make(map[string][]interface{})
		for _, item := range apiModel.BodyVars {
			if len(item.TestValue) == 1 {
				mapData[item.Name] = item.TestValue[0]
			} else if len(item.TestValue) > 1 {
				for _, subItem := range item.TestValue {
					multiMapData[item.Name] = append(multiMapData[item.Name], subItem)
				}
			} else {
				mapData[item.Name] = ""
			}
		}
		dataFile.Single.Body = mapData
		dataFile.Multi.Body = multiMapData
	}

	var dataInfo []byte
	if apiModel.BodyMode == "json" {
		dataInfo, err = json.MarshalIndent(dataFile, "", "    ")
	} else {
		dataInfo, err = yaml.Marshal(dataFile)
	}

	if err != nil {
		Logger.Error("%s", err)
	}

	var dbSceneData DbSceneData
	models.Orm.Table("scene_data").Where("api_id = ? and app = ? and name = ?", sceneData.ApiId, sceneData.App, sceneData.Name).Find(&dbSceneData)
	sceneData.Content = fmt.Sprintf("<pre><code>%s</code></pre>", dataInfo)
	sceneData.FileType = 1

	if len(dbSceneData.ApiId) == 0 {
		err = models.Orm.Table("scene_data").Create(&sceneData).Error
	} else {
		err = models.Orm.Table("scene_data").Where("id = ?", dbSceneData.Id).Update(&sceneData).Error
	}

	ymlFilePath := fmt.Sprintf("%s/%s", DataBasePath, sceneData.FileName)
	err = ioutil.WriteFile(ymlFilePath, dataInfo, 0644)
	if err != nil {
		Logger.Error("%s", err)
	}

	if len(apiModel.Product) > 0 {
		var apiList, products []string
		if len(dataFile.Api.PreApi) > 0 {
			apiList = append(apiList, dataFile.Api.PreApi...)
		}
		apiList = append(apiList, sceneData.FileName)
		if len(dataFile.Api.PostApis) > 0 {
			apiList = append(apiList, dataFile.Api.PostApis...)
		}

		products = append(products, apiModel.Product)

		// 如果数据域调试增加了前置或后置接口数据，运行后自动进行场景保存
		if len(dataFile.Api.PreApi) > 0 || len(dataFile.Api.PostApis) > 0 {
			err = SavePlaybook(fullName, apiList, products)
		}

	}

	return
}

func RunApiDebugData(apiModel ApiDataSaveModel) (runResp RunRespModel, err error) {
	urlStr, headerStr, requestStr, responseStr, outputStr, result, dst, err := RunSceneDebugContent(apiModel)
	runResp.Request = requestStr
	runResp.TestResult = result
	runResp.Url = urlStr
	runResp.Header = headerStr
	runResp.Response = responseStr

	if err != nil || result != "pass" {
		runResp.TestResult = "fail"
		runResp.FailReason = fmt.Sprintf("%v", err)
	}

	envType, errTmp := GetEnvTypeByName(apiModel.Product)
	if errTmp != nil {
		Logger.Warning("未取环境信息异常: %v", errTmp)
		//return
		if err != nil {
			err = fmt.Errorf("%s;%s", err, errTmp)
		} else {
			err = errTmp
		}
	}

	runResp.Output = outputStr

	err1 := apiModel.WriteDataFileHistoryResult(result, dst, envType, err)

	if err1 != nil {
		err = fmt.Errorf("%s, %s", err, err1)
	}

	return
}

func RunHistoryData(apiModel HistorySaveModel) (runResp RunRespModel, err error) {
	url, headerStr, requestStr, responseStr, outputStr, result, dst, err := RunHistoryContent(apiModel)
	runResp.Request = requestStr
	runResp.TestResult = result
	runResp.Url = url
	runResp.Header = headerStr
	runResp.Response = responseStr
	if result != "pass" {
		runResp.FailReason = fmt.Sprintf("%s", err)
	}

	envType, _ := GetEnvTypeByName(apiModel.Product)

	runResp.Output = outputStr

	err1 := apiModel.WriteDataFileHistoryResult(result, dst, envType, err)

	if err1 != nil {
		Logger.Error("%v", err1)
		err = fmt.Errorf("%s, %s", err, err1)
	}

	return
}

func GetDataFileFromModel(apiModel ApiDataSaveModel) (dataFile DataFile) {
	dataFile.Name = apiModel.DataDesc
	dataFile.Env.Host = apiModel.Host
	dataFile.Env.Protocol = apiModel.Prototype

	dataFile.Api.Method = apiModel.Method

	dataFile.Env.Prepath = apiModel.Prefix

	dataFile.Api.App = apiModel.App
	dataFile.Api.Description = apiModel.ApiDesc
	dataFile.Api.Module = apiModel.Module
	dataFile.Api.Path = apiModel.Path
	dataFile.Version = 1
	dataFile.ApiId = fmt.Sprintf("%s_%s", apiModel.Method, apiModel.Path)
	dataFile.IsUseEnvConfig = "no"
	dataFile.IsRunPostApis = "no"
	dataFile.IsRunPreApis = "no"
	dataFile.IsParallel = "yes"

	for _, item := range apiModel.PreApis {
		dataFile.Api.PreApi = append(dataFile.Api.PreApi, item.DataFile)
	}

	for _, item := range apiModel.PostApis {
		dataFile.Api.PostApis = append(dataFile.Api.PostApis, item.DataFile)
	}

	dataFile.Single.Header = make(map[string]interface{})
	for _, item := range apiModel.HeaderVars {
		if len(item.TestValue) > 0 {
			dataFile.Single.Header[item.Name] = item.TestValue[0]
		}
	}

	if len(apiModel.BodyMode) > 0 {
		if apiModel.BodyMode == "json" || apiModel.BodyMode == "form-data" || apiModel.BodyMode == "x-www-form-urlencoded" {
			dataFile.Single.Header["Content-Type"] = fmt.Sprintf("application/%s", apiModel.BodyMode)
		} else {
			dataFile.Single.Header["Content-Type"] = apiModel.BodyMode
		}
	}

	dataFile.Single.Query = make(map[string]interface{})
	dataFile.Multi.Query = make(map[string][]interface{})
	for _, item := range apiModel.QueryVars {
		if len(item.TestValue) > 0 && len(item.TestValue) == 1 {
			dataFile.Single.Query[item.Name] = item.TestValue[0]
		} else if len(item.TestValue) > 1 {
			dataFile.Multi.Query[item.Name] = item.TestValue
		}
	}

	dataFile.Single.Body = make(map[string]interface{})
	dataFile.Multi.Body = make(map[string][]interface{})
	for _, item := range apiModel.BodyVars {
		if len(item.TestValue) > 0 && len(item.TestValue) == 1 {
			dataFile.Single.Body[item.Name] = item.TestValue[0]
		} else if len(item.TestValue) > 1 {
			dataFile.Multi.Body[item.Name] = item.TestValue
		}
	}

	dataFile.Single.Path = make(map[string]interface{})
	dataFile.Multi.Path = make(map[string][]interface{})
	for _, item := range apiModel.PathVars {
		if len(item.TestValue) > 0 && len(item.TestValue) == 1 {
			dataFile.Single.Path[item.Name] = item.TestValue[0]
		} else if len(item.TestValue) > 1 {
			dataFile.Multi.Path[item.Name] = item.TestValue
		}
	}

	for _, item := range apiModel.Asserts {
		if len(item.Type) > 0 {
			dataFile.Assert = append(dataFile.Assert, item)
		}
	}

	return
}

func RunSceneDebugContent(apiModel ApiDataSaveModel) (urlStr, headerStr, requestStr, responseStr, outputStr, result, dst string, err error) {
	df := GetDataFileFromModel(apiModel)
	fileName := fmt.Sprintf("%s-%s-%s.yml", apiModel.Module, apiModel.ApiDesc, apiModel.DataDesc)
	filePath := fmt.Sprintf("%s/%s", DataBasePath, fileName)

	urlStr, headerStr, requestStr, responseStr, outputStr, result, dst, err = df.RunDataFileStruct(apiModel.App, apiModel.Product, filePath, "common", "consoleData", nil)

	return
}

func (df DataFile) GetResponseStr() (urlStr, headerStr, requestStr, responseStr, outputStr string, err error) {
	if len(df.Response) > 0 {
		for _, item := range df.Response {
			responseMap := make(map[string]interface{})
			errTmp := json.Unmarshal([]byte(item), &responseMap)
			resJson, errTmp := json.MarshalIndent(responseMap, "", "    ")
			if errTmp != nil {
				Logger.Error("%s", errTmp)
				if err != nil {
					err = fmt.Errorf("%s; %s", err, errTmp)
				} else {
					err = errTmp
				}
			}
			if len(responseStr) == 0 {
				responseStr = string(resJson)
			} else {
				responseStr = fmt.Sprintf("%s\n%s", responseStr, string(resJson))
			}
		}
	}
	if len(df.Urls) > 0 {
		for _, item := range df.Urls {
			if len(urlStr) == 0 {
				urlStr = item
			} else {
				urlStr = fmt.Sprintf("%s\n%s", urlStr, item)
			}
		}
	}

	hJson, errTmp := json.MarshalIndent(df.Single.Header, "", "    ")
	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}
	}
	headerStr = string(hJson)

	for _, item := range df.Request {
		requestMap := make(map[string]interface{})
		json.Unmarshal([]byte(item), &requestMap)                   // 将JSON字符串反序列化为结构体
		rJson, errTmp := json.MarshalIndent(requestMap, "", "    ") // 将结构体再格式化为有层级的JSON字符串
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}
		if len(requestStr) == 0 {
			requestStr = string(rJson)
		} else {
			requestStr = fmt.Sprintf("%s\n%s", requestStr, string(rJson))
		}
	}

	if df.Output != nil {
		for k, v := range df.Output {
			tmp, _ := json.Marshal(v)
			strTmp := fmt.Sprintf("%s: %s", k, tmp)
			if len(outputStr) == 0 {
				outputStr = strTmp
			} else {
				outputStr = fmt.Sprintf("%s\n%s", outputStr, strTmp)
			}

		}

	}
	return
}

func SavePlaybook(name string, apiList, product []string) (err error) {
	var dbScene DbScene
	var scene SceneWithNoUpdateTime
	scene.Name = name
	var apisStr string
	for _, item := range apiList {
		if len(apisStr) == 0 {
			apisStr = item
		} else {
			apisStr = fmt.Sprintf("%s %s", apisStr, item)
		}
	}
	scene.ApiList = apisStr
	scene.RunTime = 1
	scene.SceneType = 1
	scene.Priority = 999
	scene.UserName = "控制台"

	for _, item := range product {
		scene.Product = item
		models.Orm.Table("playbook").Where("name = ? and product = ?", name, item).Find(&dbScene)
		if len(dbScene.ApiList) == 0 {
			err = models.Orm.Table("playbook").Create(scene).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		} else {
			err = models.Orm.Table("playbook").Where("name = ? and product = ?", name, item).Update(scene).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		}
	}

	return
}

func GetDataInfo(appName, method, path, dataDesc string) (apiModel ApiDataSaveModel) {
	apiModel.App = appName
	apiModel.Method = method
	apiModel.Path = path
	apiId := fmt.Sprintf("%s_%s", method, path)

	var apiDef ApiDefinition
	models.Orm.Table("api_definition").Where("app = ? and http_method = ? and path = ?", appName, method, path).Find(&apiDef)
	apiModel.ApiDesc = apiDef.ApiDesc
	apiModel.Module = apiDef.ApiModule
	apiModel.DataDesc = apiDef.ApiDesc
	tagNum := strings.Count(dataDesc, "-")
	var sceneData SceneData
	var fullName string
	if tagNum >= 3 {
		fullName = dataDesc
	} else {
		fullName = fmt.Sprintf("控制台-%s-%s-%s", apiDef.ApiModule, apiDef.ApiDesc, dataDesc)
	}
	models.Orm.Table("scene_data").Where("app = ? and api_id = ? and name = ?", appName, apiId, fullName).Find(&sceneData)
	apiModel.DataDesc = dataDesc

	var pathVarList, queryVarList, bodyVarList, headerVarList, respVarList []VarDefModel

	pathVarList = apiDef.PathVariable
	queryVarList = apiDef.QueryParameter
	bodyVarList = apiDef.Body
	headerVarList = apiDef.Header
	respVarList = apiDef.Response

	var dataFile DataFile

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sceneData.Content))
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	handle := doc.Find("code")
	afterTxt := handle.Text()
	if len(afterTxt) == 0 {
		return
	}

	err = yaml.Unmarshal([]byte(afterTxt), &dataFile)
	if err != nil {
		Logger.Error("%s", err)
	}

	apiModel.Asserts = dataFile.Assert

	for _, item := range dataFile.Api.PreApi {
		var depDataModel DepDataModel
		depDataModel.DataFile = item
		apiModel.PreApis = append(apiModel.PreApis, depDataModel)
	}

	for _, item := range dataFile.Api.PostApis {
		var depDataModel DepDataModel
		depDataModel.DataFile = item
		apiModel.PostApis = append(apiModel.PostApis, depDataModel)
	}

	for _, item := range pathVarList {
		var varData VarDataModel
		varData.ValueType = item.ValueType
		varData.Name = item.Name
		varData.EgValue = item.EgValue
		varData.Desc = item.Desc
		varData.IsMust = item.IsMust
		if value, ok := dataFile.Single.Path[item.Name]; ok {
			varData.TestValue = append(varData.TestValue, value)
		}
		if value, ok := dataFile.Multi.Path[item.Name]; ok {
			varData.TestValue = append(varData.TestValue, value...)
		}
		apiModel.PathVars = append(apiModel.PathVars, varData)
	}

	for _, item := range queryVarList {
		var varData VarDataModel
		varData.ValueType = item.ValueType
		varData.Name = item.Name
		varData.EgValue = item.EgValue
		varData.Desc = item.Desc
		varData.IsMust = item.IsMust
		if value, ok := dataFile.Single.Query[item.Name]; ok {
			valueStr := Interface2Str(value)
			varData.TestValue = append(varData.TestValue, valueStr)
		}

		if value, ok := dataFile.Multi.Query[item.Name]; ok {
			varData.TestValue = append(varData.TestValue, value...)
		}

		apiModel.QueryVars = append(apiModel.QueryVars, varData)
	}

	for _, item := range respVarList {
		var varData VarDataModel
		varData.ValueType = item.ValueType
		varData.Name = item.Name
		varData.EgValue = item.EgValue
		varData.Desc = item.Desc
		varData.IsMust = item.IsMust
		apiModel.RespVars = append(apiModel.RespVars, varData)
	}

	for _, item := range bodyVarList {
		var varData VarDataModel
		varData.ValueType = item.ValueType
		varData.Name = item.Name
		varData.EgValue = item.EgValue
		varData.Desc = item.Desc
		varData.IsMust = item.IsMust
		if value, ok := dataFile.Single.Body[item.Name]; ok {
			varData.TestValue = append(varData.TestValue, value)
		}
		if value, ok := dataFile.Multi.Body[item.Name]; ok {
			varData.TestValue = append(varData.TestValue, value...)
		}
		apiModel.BodyVars = append(apiModel.BodyVars, varData)
	}

	for _, item := range headerVarList {
		var varData VarDataModel
		varData.ValueType = item.ValueType
		varData.Name = item.Name
		varData.EgValue = item.EgValue
		varData.Desc = item.Desc
		varData.IsMust = item.IsMust
		if value, ok := dataFile.Single.Header[item.Name]; ok {
			varData.TestValue = append(varData.TestValue, value)
		}
		apiModel.HeaderVars = append(apiModel.HeaderVars, varData)
		if item.Name == "Content-Type" {
			valueStr := Interface2Str(item.EgValue)
			apiModel.BodyMode = valueStr
		}
	}

	var other OtherModel
	other.IsParallel = dataFile.IsParallel
	other.Version = dataFile.Version
	other.IsUseEnvConfig = dataFile.IsUseEnvConfig
	other.ApiId = dataFile.ApiId
	apiModel.Other = append(apiModel.Other, other)

	return
}

func GetDataInfoByDataDesc(appName, module, apiDesc, dataDesc string) (apiModel ApiDataSaveModel) {
	var sceneData SceneData
	var fullName string
	if len(appName) == 0 {
		models.Orm.Table("scene_data").Where("name = ? and file_type = 1", dataDesc).Find(&sceneData)
		appName = sceneData.App
	} else {
		models.Orm.Table("scene_data").Where("app = ? and name = ?  and file_type = 1", appName, dataDesc).Find(&sceneData)
	}
	if len(sceneData.ApiId) == 0 {
		fullName = fmt.Sprintf("控制台-%s-%s-%s", module, apiDesc, dataDesc)
		models.Orm.Table("scene_data").Where("app = ? and name = ?  and file_type = 1", appName, fullName).Find(&sceneData)
	}

	if len(sceneData.ApiId) > 0 {
		pathTmps := strings.Split(sceneData.ApiId, "_")
		apiModel.Method = pathTmps[0]
		apiModel.Path = pathTmps[1]
	}

	var apiDef ApiStringDefinition
	models.Orm.Table("api_definition").Where("app = ? and http_method = ? and path = ?", appName, apiModel.Method, apiModel.Path).Find(&apiDef)
	if len(appName) == 0 {
		apiModel.App = apiDef.App
	} else {
		apiModel.App = appName
	}

	if len(dataDesc) == 0 {
		apiModel.DataDesc = apiDef.ApiDesc
	} else {
		apiModel.DataDesc = dataDesc
	}

	if len(module) == 0 {
		apiModel.Module = apiDef.ApiModule
	} else {
		apiModel.Module = module
	}

	var pathVarList, queryVarList, bodyVarList, headerVarList, respVarList []VarDefModel
	json.Unmarshal([]byte(apiDef.PathVariable), &pathVarList)
	json.Unmarshal([]byte(apiDef.QueryParameter), &queryVarList)
	json.Unmarshal([]byte(apiDef.Body), &bodyVarList)
	json.Unmarshal([]byte(apiDef.Header), &headerVarList)

	//json.Unmarshal([]byte(apiDef.Response), &respVarList)  // 数据域无需展示Resp信息

	var dataFile DataFile

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sceneData.Content))
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	handle := doc.Find("code")
	afterTxt := handle.Text()
	if len(afterTxt) == 0 {
		return
	}

	err = yaml.Unmarshal([]byte(afterTxt), &dataFile)
	if err != nil {
		Logger.Error("%s", err)
	}

	if len(apiModel.App) == 0 {
		apiModel.App = dataFile.Api.App
	}

	if len(apiModel.Module) == 0 {
		apiModel.Module = dataFile.Api.Module
	}

	if len(apiModel.ApiDesc) == 0 {
		apiModel.ApiDesc = dataFile.Api.Description
	}

	apiModel.Asserts = dataFile.Assert
	apiModel.Actions = dataFile.Action
	apiModel.Prefix = dataFile.Env.Prepath

	if len(appName) == 0 || len(module) == 0 || len(apiDesc) == 0 {
		apiModel.Path = dataFile.Api.Path
		apiModel.Method = dataFile.Api.Method
		apiModel.Prefix = dataFile.Env.Prepath
	}

	for _, item := range dataFile.Api.PreApi {
		var depDataModel DepDataModel
		depDataModel.DataFile = item
		apiModel.PreApis = append(apiModel.PreApis, depDataModel)
	}

	for _, item := range dataFile.Api.PostApis {
		var depDataModel DepDataModel
		depDataModel.DataFile = item
		apiModel.PostApis = append(apiModel.PostApis, depDataModel)
	}

	for key, value := range dataFile.Single.Path {
		var varData VarDataModel
		for _, item := range pathVarList {
			if item.Name == key {
				varData.Desc = item.Desc
				varData.EgValue = item.EgValue
				varData.IsMust = item.IsMust
				varData.ValueType = item.ValueType
			}
		}
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value)
		apiModel.PathVars = append(apiModel.PathVars, varData)
	}
	for key, value := range dataFile.Multi.Path {
		var varData VarDataModel
		for _, item := range pathVarList {
			if item.Name == key {
				varData.Desc = item.Desc
				varData.EgValue = item.EgValue
				varData.IsMust = item.IsMust
				varData.ValueType = item.ValueType
			}
		}
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value...)
		apiModel.PathVars = append(apiModel.PathVars, varData)
	}

	for key, value := range dataFile.Single.Query {
		var varData VarDataModel
		for _, item := range queryVarList {
			if item.Name == key {
				varData.Desc = item.Desc
				varData.EgValue = item.EgValue
				varData.IsMust = item.IsMust
				varData.ValueType = item.ValueType
			}
		}
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value)
		apiModel.QueryVars = append(apiModel.QueryVars, varData)
	}
	for key, value := range dataFile.Multi.Query {
		var varData VarDataModel
		for _, item := range queryVarList {
			if item.Name == key {
				varData.Desc = item.Desc
				varData.EgValue = item.EgValue
				varData.IsMust = item.IsMust
				varData.ValueType = item.ValueType
			}
		}
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value...)
		apiModel.QueryVars = append(apiModel.QueryVars, varData)
	}

	for _, item := range respVarList {
		var varData VarDataModel
		varData.Desc = item.Desc
		varData.EgValue = item.EgValue
		varData.IsMust = item.IsMust
		varData.ValueType = item.ValueType
		apiModel.RespVars = append(apiModel.RespVars, varData)
	}

	for key, value := range dataFile.Single.Body {
		var varData VarDataModel
		for _, item := range bodyVarList {
			if item.Name == key {
				varData.Desc = item.Desc
				varData.EgValue = item.EgValue
				varData.IsMust = item.IsMust
				varData.ValueType = item.ValueType
			}
		}
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value)
		apiModel.BodyVars = append(apiModel.BodyVars, varData)
	}
	for key, value := range dataFile.Multi.Body {
		var varData VarDataModel
		for _, item := range bodyVarList {
			if item.Name == key {
				varData.Desc = item.Desc
				varData.EgValue = item.EgValue
				varData.IsMust = item.IsMust
				varData.ValueType = item.ValueType
			}
		}
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value...)
		apiModel.BodyVars = append(apiModel.BodyVars, varData)
	}

	for key, value := range dataFile.Single.Header {
		var varData VarDataModel
		for _, item := range headerVarList {
			if item.Name == key {
				varData.Desc = item.Desc
				varData.EgValue = item.EgValue
				varData.IsMust = item.IsMust
				varData.ValueType = item.ValueType
			}
		}
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value)
		apiModel.HeaderVars = append(apiModel.HeaderVars, varData)

		if key == "Content-Type" {
			valueStr := Interface2Str(value)
			apiModel.BodyMode = valueStr
		}

	}

	var other OtherModel
	other.IsParallel = dataFile.IsParallel
	other.Version = dataFile.Version
	other.IsUseEnvConfig = dataFile.IsUseEnvConfig
	other.ApiId = dataFile.ApiId
	apiModel.Other = append(apiModel.Other, other)

	return
}

func GetApiDataList(appName, method, path string) (dataDescs, dataFileList []string) {
	var dataDescTmp []string
	apiId := fmt.Sprintf("%s_%s", method, path)
	models.Orm.Table("scene_data").Group("name").Select("name").Where("app = ? and api_id = ?", appName, apiId).Pluck("name", &dataDescTmp)
	for _, item := range dataDescTmp {
		names := strings.SplitN(item, "-", 4)
		dataDescs = append(dataDescs, names[3])
	}

	models.Orm.Table("scene_data").Group("file_name").Select("file_name").Pluck("file_name", &dataFileList)

	return
}

func GetProductPlaybook(product string) (playbooks []string) {
	models.Orm.Table("playbook").Order("created_at desc").Where("product = ?", product).Pluck("name", &playbooks)
	return
}

func GetDateHistory(dateName string) (historyDatas []string) {
	dayBegin := dateName + " 00:00:00"
	dayEnd := dateName + " 23:59:59"
	var rawDatas []string
	models.Orm.Table("scene_data_test_history").Order("created_at desc").Limit(50).Select("content").Where("created_at > ? and created_at < ? and content <> ''", dayBegin, dayEnd).Pluck("content", &rawDatas)
	for _, item := range rawDatas {
		hData := GetStrFromHtml(item)
		historyDatas = append(historyDatas, hData)
	}
	return
}

func GetDateSceneHistory(dateName string) (historyDatas []string) {
	dayBegin := dateName + " 00:00:00"
	dayEnd := dateName + " 23:59:59"
	var sceneHistoryRecords []SceneHistoryRecord
	models.Orm.Table("scene_test_history").Order("created_at desc").Limit(30).Where("created_at > ? and created_at < ?", dayBegin, dayEnd).Find(&sceneHistoryRecords)
	for _, item := range sceneHistoryRecords {
		var data string
		data = fmt.Sprintf("%s:%s", item.Id, item.Name)
		historyDatas = append(historyDatas, data)
	}
	return
}

func GetHistoryByFileName(fileName string) (apiModel HistorySaveModel, err error) {
	var dataFile DataFile

	dirName := GetHistoryDataDirName(fileName)

	filePath := fmt.Sprintf("%s/%s/%s", HistoryBasePath, dirName, fileName)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if strings.HasSuffix(fileName, ".json") {
		err = json.Unmarshal([]byte(content), &dataFile)
	} else {
		err = yaml.Unmarshal([]byte(content), &dataFile)
	}
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	apiModel.App = dataFile.Api.App
	apiModel.Module = dataFile.Api.Module
	apiModel.ApiDesc = dataFile.Api.Description
	apiModel.Asserts = dataFile.Assert
	apiModel.Prefix = dataFile.Env.Prepath
	apiModel.DataDesc = dataFile.Name
	apiModel.Method = dataFile.Api.Method
	apiModel.Path = dataFile.Api.Path
	apiModel.Prefix = dataFile.Env.Prepath
	apiModel.Prototype = dataFile.Env.Protocol
	apiModel.Host = dataFile.Env.Host
	apiModel.ApiDesc = dataFile.Api.Description
	if len(dataFile.Output) > 0 {
		for k, v := range dataFile.Output {
			var valueTmp string
			for _, item := range v {
				valueTmp = fmt.Sprintf("%s %v", valueTmp, item)
			}
			strTmp := fmt.Sprintf("%s: %v", k, valueTmp)

			if len(apiModel.Output) == 0 {
				apiModel.Output = strTmp
			} else {
				apiModel.Output = fmt.Sprintf("%s\n%s", apiModel.Output, strTmp)
			}

		}
	}

	hJson, errTmp := json.MarshalIndent(dataFile.Single.Header, "", "    ")
	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}
	}
	apiModel.Header = string(hJson)

	for _, item := range dataFile.Urls {
		if len(apiModel.Url) == 0 {
			apiModel.Url = item
		} else {
			apiModel.Url = fmt.Sprintf("%s\n%s", apiModel.Url, item)
		}

	}

	for _, item := range dataFile.Request {
		requestMap := make(map[string]interface{})
		errTmp := json.Unmarshal([]byte(item), &requestMap)
		requstJson, errTmp := json.MarshalIndent(requestMap, "", "    ")
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}
		if len(apiModel.Request) == 0 {
			apiModel.Request = string(requstJson)
		} else {
			apiModel.Request = fmt.Sprintf("%s\n%s", apiModel.Request, string(requstJson))
		}

	}

	for _, item := range dataFile.TestResult {
		if len(apiModel.TestResult) == 0 {
			apiModel.TestResult = item
		} else {
			apiModel.TestResult = fmt.Sprintf("%s\n%s", apiModel.TestResult, item)
		}

	}

	for _, item := range dataFile.Response {
		responseMap := make(map[string]interface{})
		errTmp := json.Unmarshal([]byte(item), &responseMap)
		responseJson, errTmp := json.MarshalIndent(responseMap, "", "    ")
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}
		if len(apiModel.Response) == 0 {
			apiModel.Response = string(responseJson)
		} else {
			apiModel.Response = fmt.Sprintf("%s\n%s", apiModel.Response, string(responseJson))
		}

	}

	for _, item := range dataFile.Api.PreApi {
		var depDataModel DepDataModel
		depDataModel.DataFile = item
		apiModel.PreApis = append(apiModel.PreApis, depDataModel)
	}

	for _, item := range dataFile.Api.PostApis {
		var depDataModel DepDataModel
		depDataModel.DataFile = item
		apiModel.PostApis = append(apiModel.PostApis, depDataModel)
	}

	for key, value := range dataFile.Single.Path {
		var varData VarDataModel
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value)
		apiModel.PathVars = append(apiModel.PathVars, varData)
	}

	for key, value := range dataFile.Multi.Path {
		var varData VarDataModel
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value...)
		apiModel.PathVars = append(apiModel.PathVars, varData)
	}

	for key, value := range dataFile.Single.Query {
		var varData VarDataModel
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value)
		apiModel.QueryVars = append(apiModel.QueryVars, varData)
	}
	for key, value := range dataFile.Multi.Query {
		var varData VarDataModel
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value...)
		apiModel.QueryVars = append(apiModel.QueryVars, varData)
	}

	for key, value := range dataFile.Single.Body {
		var varData VarDataModel
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value)
		apiModel.BodyVars = append(apiModel.BodyVars, varData)
	}
	for key, value := range dataFile.Multi.Body {
		var varData VarDataModel
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value...)
		apiModel.BodyVars = append(apiModel.BodyVars, varData)
	}

	for key, value := range dataFile.Single.Header {
		var varData VarDataModel
		varData.Name = key
		varData.TestValue = append(varData.TestValue, value)
		apiModel.HeaderVars = append(apiModel.HeaderVars, varData)
		if key == "Content-Type" {
			valueStr := Interface2Str(value)
			apiModel.BodyMode = valueStr
		}
	}

	var other OtherModel
	other.IsParallel = dataFile.IsParallel
	other.Version = dataFile.Version
	other.IsUseEnvConfig = dataFile.IsUseEnvConfig
	other.ApiId = dataFile.ApiId
	apiModel.Other = append(apiModel.Other, other)

	return
}

func GetSceneHistory(name string) (sceneModel SceneHistorySaveModel, err error) {
	tmpList := strings.Split(name, ":")
	if len(tmpList) < 2 {
		err = fmt.Errorf("数据有误，请核对， name: %v", name)
		return
	}
	var sceneRecord SceneRecord
	models.Orm.Table("scene_test_history").Where("id = ?", tmpList[0]).Find(&sceneRecord)
	sceneModel.RunNum = 1
	sceneModel.Name = sceneRecord.Name
	switch sceneRecord.SceneType {
	case 1:
		sceneModel.SceneType = "串行中断"
	case 2:
		sceneModel.SceneType = "串行比较"
	case 3:
		sceneModel.SceneType = "串行继续"
	case 4:
		sceneModel.SceneType = "普通并发"
	case 5:
		sceneModel.SceneType = "并发比较"
	default:
		sceneModel.SceneType = "串行中断"
	}

	sceneModel.FailReason = sceneRecord.FailReason
	sceneModel.TestResult = sceneRecord.Result
	sceneModel.LastFile = sceneRecord.LastFile
	sceneModel.Product = sceneRecord.Product
	var dataList []string
	dataList = GetListFromHtml(sceneRecord.ApiList)
	for _, item := range dataList {
		var depData DepDataModel
		depData.DataFile = item
		sceneModel.DataList = append(sceneModel.DataList, depData)
	}
	return
}

func Str2VarModel(bodyStr string) (bodyVar []VarDataModel, err error) {
	bodyStr = strings.Replace(bodyStr, "\\", "", -1)
	var tempMap map[string]interface{}

	err = json.Unmarshal([]byte(bodyStr), &tempMap)

	if err != nil {
		Logger.Error("%s", err)
	}
	for k, v := range tempMap {
		var varModel VarDataModel
		varModel.Name = k
		varModel.TestValue = append(varModel.TestValue, v)
		bodyVar = append(bodyVar, varModel)
	}
	return
}

func Str2DefModel(bodyStr string) (bodyVar []VarDefModel, err error) {
	bodyStr = strings.Replace(bodyStr, "\\", "", -1)
	var tempMap map[string]interface{}

	err = json.Unmarshal([]byte(bodyStr), &tempMap)

	if err != nil {
		Logger.Error("%s", err)
	}
	for k, v := range tempMap {
		var varModel VarDefModel
		varModel.Name = k
		varType := fmt.Sprintf("%T", v)
		varModel.ValueType = varType
		bodyVar = append(bodyVar, varModel)
	}
	return
}

func (dataFile DataFile) UpdateDataFileFromHistoryModel(apiModel HistorySaveModel) {
	dataFile.Name = apiModel.DataDesc
	if len(apiModel.Host) != 0 {
		dataFile.Env.Host = apiModel.Host

	}
	if len(apiModel.Prototype) != 0 {
		dataFile.Env.Protocol = apiModel.Prototype
	}

	if len(apiModel.Method) != 0 {
		dataFile.Api.Method = apiModel.Method
	}

	if len(apiModel.Prefix) != 0 {
		dataFile.Env.Prepath = apiModel.Prefix
	}

	dataFile.Api.App = apiModel.App
	dataFile.Api.Description = apiModel.ApiDesc
	dataFile.Api.Module = apiModel.Module
	dataFile.Api.Path = apiModel.Path

	for _, item := range apiModel.PreApis {
		dataFile.Api.PreApi = append(dataFile.Api.PreApi, item.DataFile)
	}

	for _, item := range apiModel.PostApis {
		dataFile.Api.PostApis = append(dataFile.Api.PostApis, item.DataFile)
	}

	dataFile.Single.Header = make(map[string]interface{})
	for _, item := range apiModel.HeaderVars {
		if len(item.TestValue) > 0 {
			dataFile.Single.Header[item.Name] = item.TestValue[0]
		}
	}

	if len(apiModel.BodyMode) > 0 {
		if apiModel.BodyMode == "json" || apiModel.BodyMode == "form-data" || apiModel.BodyMode == "x-www-form-urlencoded" {
			dataFile.Single.Header["Content-Type"] = fmt.Sprintf("application/%s", apiModel.BodyMode)
		} else {
			dataFile.Single.Header["Content-Type"] = apiModel.BodyMode
		}
	}

	dataFile.Single.Query = make(map[string]interface{})
	dataFile.Multi.Query = make(map[string][]interface{})
	for _, item := range apiModel.QueryVars {
		if len(item.TestValue) > 0 && len(item.TestValue) == 1 {
			dataFile.Single.Query[item.Name] = item.TestValue[0]
		} else if len(item.TestValue) > 1 {
			dataFile.Multi.Query[item.Name] = item.TestValue
		}
	}

	dataFile.Multi.Body = make(map[string][]interface{})
	dataFile.Single.Body = make(map[string]interface{})
	for _, item := range apiModel.BodyVars {
		if len(item.TestValue) > 0 && len(item.TestValue) == 1 {
			dataFile.Single.Body[item.Name] = item.TestValue[0]
		} else if len(item.TestValue) > 1 {
			if item.TestValue != nil {
				dataFile.Multi.Body[item.Name] = item.TestValue
			}
		}
	}

	dataFile.Single.Path = make(map[string]interface{})
	dataFile.Multi.Path = make(map[string][]interface{})
	for _, item := range apiModel.PathVars {
		if len(item.TestValue) > 0 && len(item.TestValue) == 1 {
			dataFile.Single.Path[item.Name] = item.TestValue[0]
		} else if len(item.TestValue) > 1 {
			dataFile.Multi.Path[item.Name] = item.TestValue
		}
	}

	for _, item := range apiModel.Asserts {
		if len(item.Type) > 0 {
			dataFile.Assert = append(dataFile.Assert, item)
		}
	}

	return
}

func RunHistoryContent(apiModel HistorySaveModel) (urlStr, headerStr, requestStr, responseStr, outputStr, result, dst string, err error) {
	var dataFile DataFile
	dirName := GetHistoryDataDirName(apiModel.FileName)
	historyFilePath := fmt.Sprintf("%s/%s/%s", HistoryBasePath, dirName, apiModel.FileName) //

	content, err := ioutil.ReadFile(historyFilePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if strings.HasSuffix(apiModel.FileName, ".json") {
		err = json.Unmarshal([]byte(content), &dataFile)
	} else {
		err = yaml.Unmarshal([]byte(content), &dataFile)
	}

	dataFile.UpdateDataFileFromHistoryModel(apiModel)

	urlStr, headerStr, requestStr, responseStr, outputStr, result, dst, err = dataFile.RunDataFileStruct(apiModel.App, apiModel.Product, historyFilePath, "again", "console", nil) // again表示历史数据再来一次

	return
}
