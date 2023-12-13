package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func CreateSceneDataFromRaw(id, mode string) (err error) {
	var apiDefinition ApiDefinition
	models.Orm.Table("api_definition").Where("id = ?", id).Find(&apiDefinition)
	Logger.Debug("apiDefinition: %#v", apiDefinition)
	if len(apiDefinition.ApiId) == 0 {
		err = fmt.Errorf("未找到API的定义信息，请核对~")
		Logger.Error("%s", err)
		return
	}

	var apiRelation ApiRelation
	models.Orm.Table("api_relation").Where("api_id = ? and app = ?", apiDefinition.ApiId, apiDefinition.App).Find(&apiRelation)
	if len(apiRelation.ApiId) == 0 {
		Logger.Debug("apiRelation: %v", apiRelation)
	}

	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", apiDefinition.App).Find(&envConfig)
	if len(envConfig.App) == 0 {
		Logger.Debug("apiRelation: %v", apiRelation)
	}

	var sceneData SceneData
	sceneData.App = apiDefinition.App
	sceneData.Name = apiDefinition.ApiDesc
	sceneData.ApiId = apiDefinition.ApiId
	sceneData.Name = fmt.Sprintf("通用-%s-%s-%s", apiDefinition.ApiModule, apiDefinition.ApiDesc, GetRandomStr(4, ""))
	if mode == "json" {
		sceneData.FileName = fmt.Sprintf("%s.json", sceneData.Name)
	} else {
		sceneData.FileName = fmt.Sprintf("%s.yml", sceneData.Name)
	}

	sceneData.RunTime = 1

	var dataFile DataFile
	dataFile.Name = sceneData.Name
	dataFile.Version = 1.0
	dataFile.IsRunPostApis = "no"
	dataFile.IsRunPreApis = "no"
	dataFile.IsParallel = "no"
	dataFile.IsUseEnvConfig = "yes"
	dataFile.ApiId = apiDefinition.ApiId
	dataFile.Env.Protocol = envConfig.Protocol
	dataFile.Env.Host = envConfig.Ip
	dataFile.Env.Prepath = envConfig.Prepath
	dataFile.Api.App = apiDefinition.App
	dataFile.Api.Description = apiDefinition.ApiDesc
	dataFile.Api.Module = apiDefinition.ApiModule
	dataFile.Api.Method = apiDefinition.HttpMethod
	dataFile.Api.Path = apiDefinition.Path

	if len(apiRelation.OutVars) > 0 {
		var outVars map[string]string
		outVars = make(map[string]string)
		err = json.Unmarshal([]byte(apiRelation.OutVars), &outVars)
		if err != nil {
			Logger.Error("%s", err)
			return
		}
		for k, v := range outVars {
			var sceneAssert SceneAssert
			sceneAssert.Type = "output"
			sceneAssert.Source = v
			sceneAssert.Value = k
			dataFile.Assert = append(dataFile.Assert, sceneAssert)
		}
	}

	if len(apiRelation.ParamApis) > 0 {
		dataFile.Api.ParamApis = strings.Split(apiRelation.ParamApis, ",")
	}

	if len(envConfig.Auth) > 0 {
		header := make(map[string]interface{})
		err = json.Unmarshal([]byte(envConfig.Auth), &header)
		if err != nil {
			Logger.Error("%s", err)
		}
		dataFile.Single.Header = header
	}

	if len(apiDefinition.QueryParameter) > 0 {
		queryParam := make(map[string]interface{})
		for _, item := range apiDefinition.QueryParameter {
			queryParam[item.Name] = item.ValueType
		}

		dataFile.Single.Query = queryParam
	}

	if len(apiDefinition.PathVariable) > 0 {
		pathVar := make(map[string]interface{})
		for _, item := range apiDefinition.PathVariable {
			pathVar[item.Name] = item.ValueType
		}
		dataFile.Single.Path = pathVar
	}

	if len(apiDefinition.Body) > 0 {
		body := make(map[string]interface{})
		for _, item := range apiDefinition.PathVariable {
			body[item.Name] = item.ValueType
		}
		dataFile.Single.Body = body
	}
	var dataInfo []byte
	if mode == "json" {
		dataInfo, err = json.MarshalIndent(dataFile, "", "    ")
	} else {
		dataInfo, err = yaml.Marshal(dataFile)
	}

	if err != nil {
		Logger.Error("%s", err)
	}

	var dbSceneData DbSceneData
	models.Orm.Table("scene_data").Where("api_id = ? and app = ? and name = ?", apiDefinition.ApiId, apiDefinition.App, sceneData.Name).Find(&dbSceneData)
	sceneData.Content = fmt.Sprintf("<pre><code>%s</code></pre>", dataInfo)

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

	return
}

func CreateSingleSceneDataFromData(apiId string, id int, source, mode string) (err error) {
	var dataFile DataFile
	var sceneData SceneData
	if source == "data" {
		var apiTestData ApiTestData
		models.Orm.Table("api_test_data").Where("id = ? and api_id = ?", id, apiId).Find(&apiTestData)

		if len(apiTestData.ApiId) == 0 {
			err1 := fmt.Errorf("未找到接口[%s]测试数据", apiId)
			err = err1
			Logger.Error("%s", err)
			return
		}
		sceneData.Name = fmt.Sprintf("通用-%s-%s-%s", apiTestData.ApiModule, apiTestData.ApiDesc, GetRandomStr(4, ""))
		if strings.Contains(apiTestData.UrlQuery, "?") {
			var queryParam map[string]interface{}
			queryParam = make(map[string]interface{})
			queryStrs := strings.Split(apiTestData.UrlQuery, "?")
			if strings.Contains(queryStrs[1], "&") {
				items := strings.Split(queryStrs[1], "&")
				for _, item := range items {
					kv := strings.Split(item, "=")
					queryParam[kv[0]] = kv[1]
				}
			} else {
				kv := strings.Split(queryStrs[1], "=")
				queryParam[kv[0]] = kv[1]
			}
			dataFile.Single.Query = queryParam

		}

		if len(apiTestData.Body) > 0 {
			var body map[string]interface{}
			body = make(map[string]interface{})
			err = json.Unmarshal([]byte(apiTestData.Body), &body)
			if err != nil {
				Logger.Error("%s", err)
			}
			dataFile.Single.Body = body
		}
	} else {
		var apiFuzzingData ApiFuzzingData
		models.Orm.Table("api_fuzzing_data").Where("id = ? and api_id = ?", id, apiId).Find(&apiFuzzingData)

		if len(apiFuzzingData.ApiId) == 0 {
			err1 := fmt.Errorf("未找到接口[%s]的模糊数据", apiId)
			Logger.Error("%s", err)
			err = err1
			return
		}

		sceneData.Name = fmt.Sprintf("模糊-%s-%s-%s", apiFuzzingData.ApiModule, apiFuzzingData.ApiDesc, GetRandomStr(4, ""))
		if strings.Contains(apiFuzzingData.UrlQuery, "?") {
			var queryParam map[string]interface{}
			queryParam = make(map[string]interface{})
			queryStrs := strings.Split(apiFuzzingData.UrlQuery, "?")
			if strings.Contains(queryStrs[1], "&") {
				items := strings.Split(queryStrs[1], "&")
				for _, item := range items {
					kv := strings.Split(item, "=")
					queryParam[kv[0]] = kv[1]
				}
			} else {
				kv := strings.Split(queryStrs[1], "=")
				queryParam[kv[0]] = kv[1]
			}
			dataFile.Single.Query = queryParam

		}

		if len(apiFuzzingData.Body) > 0 {
			var body map[string]interface{}
			body = make(map[string]interface{})
			err = json.Unmarshal([]byte(apiFuzzingData.Body), &body)
			if err != nil {
				Logger.Error("%s", err)
			}
			dataFile.Single.Body = body
		}
	}

	var apiDefinition ApiDefinition
	models.Orm.Table("api_definition").Where("api_id = ?", apiId).Find(&apiDefinition)

	var apiRelation ApiRelation
	models.Orm.Table("api_relation").Where("api_id = ? and app = ?", apiId, apiDefinition.App).Find(&apiRelation)
	if len(apiRelation.ApiId) == 0 {
		Logger.Debug("apiRelation: %v", apiRelation)
	}

	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", apiDefinition.App).Find(&envConfig)
	if len(envConfig.App) == 0 {
		Logger.Debug("apiRelation: %v", apiRelation)

	}

	sceneData.App = apiDefinition.App
	sceneData.ApiId = apiDefinition.ApiId
	if mode == "josn" {
		sceneData.FileName = fmt.Sprintf("%s.json", sceneData.Name)
	} else {
		sceneData.FileName = fmt.Sprintf("%s.yml", sceneData.Name)
	}

	sceneData.RunTime = 1

	dataFile.Name = sceneData.Name
	dataFile.Version = 1.0
	dataFile.IsRunPostApis = "no"
	dataFile.IsRunPreApis = "no"
	dataFile.IsParallel = "no"
	dataFile.IsUseEnvConfig = "yes"
	dataFile.ApiId = apiDefinition.ApiId
	dataFile.Env.Protocol = envConfig.Protocol
	dataFile.Env.Host = envConfig.Ip
	dataFile.Env.Prepath = envConfig.Prepath
	dataFile.Api.App = apiDefinition.App
	dataFile.Api.Description = apiDefinition.ApiDesc
	dataFile.Api.Module = apiDefinition.ApiModule
	dataFile.IsUseEnvConfig = "yes"
	tmps := strings.Split(apiDefinition.ApiId, "_")
	if len(tmps) == 2 {
		dataFile.Api.Method = tmps[0]
		dataFile.Api.Path = tmps[1]
	} else {
		err = fmt.Errorf("关联接口信息不符合标准格式httpMethod_/path, 实际：%s", apiDefinition.ApiId)
		Logger.Error("%s", err)
		return
	}

	if len(apiRelation.OutVars) > 0 {
		var outVars map[string]string
		outVars = make(map[string]string)
		err = json.Unmarshal([]byte(apiRelation.OutVars), &outVars)
		if err != nil {
			Logger.Error("%s", err)
		}
		for k, v := range outVars {
			var sceneAssert SceneAssert
			sceneAssert.Type = "output"
			sceneAssert.Source = v
			sceneAssert.Value = k
			dataFile.Assert = append(dataFile.Assert, sceneAssert)
		}
	}

	if len(apiRelation.ParamApis) > 0 {
		dataFile.Api.ParamApis = strings.Split(apiRelation.ParamApis, ",")
	}

	if len(envConfig.Auth) > 0 {
		header := make(map[string]interface{})
		err = json.Unmarshal([]byte(envConfig.Auth), &header)
		if err != nil {
			Logger.Error("%s", err)
		}
		dataFile.Single.Header = header
	}

	if len(apiDefinition.PathVariable) > 0 {
		pathVar := make(map[string]interface{})
		//err = json.Unmarshal([]byte(apiDefinition.PathVariable), &pathVar)
		//if err != nil {
		//	Logger.Error("%s", err)
		//}
		for _, item := range apiDefinition.QueryParameter {
			pathVar[item.Name] = item.ValueType
		}
		dataFile.Single.Path = pathVar
	}
	var dataInfo []byte
	if mode == "json" {
		dataInfo, err = json.MarshalIndent(dataFile, "", "    ")
	} else {
		dataInfo, err = yaml.Marshal(dataFile)
	}

	if err != nil {
		Logger.Error("%s", err)
	}

	sceneData.Content = fmt.Sprintf("<pre><code>%s</code></pre>", dataInfo)

	var dbSceneData DbSceneData
	models.Orm.Table("scene_data").Where("api_id = ? and app = ? and name = ?", apiDefinition.ApiId, apiDefinition.App, sceneData.Name).Find(&dbSceneData)
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
	return
}

func CreateMultiSceneDataFromData(apiId string, ids []int, source, mode string) (err error) {
	var dataFile DataFile
	var sceneData SceneData
	var pathMulti map[string][]interface{}
	var queryMulti, bodyMulti map[string][]interface{}
	queryMulti = make(map[string][]interface{})
	pathMulti = make(map[string][]interface{})
	bodyMulti = make(map[string][]interface{})
	if source == "data" {
		var apiTestDatas []ApiTestData
		models.Orm.Table("api_test_data").Where("id in (?)", ids).Find(&apiTestDatas)
		Logger.Debug("apiTestDatas: %v", apiTestDatas)

		for index, apiTestData := range apiTestDatas {
			if index == 0 {
				sceneData.Name = fmt.Sprintf("通用-%s-%s-%s", apiTestData.ApiModule, apiTestData.ApiDesc, GetRandomStr(4, ""))
			}
			if strings.Contains(apiTestData.UrlQuery, "?") {
				queryStrs := strings.Split(apiTestData.UrlQuery, "?")
				if strings.Contains(queryStrs[1], "&") {
					items := strings.Split(queryStrs[1], "&")
					for _, item := range items {
						kv := strings.Split(item, "=")
						queryMulti[kv[0]] = append(queryMulti[kv[0]], kv[1])
					}
				} else {
					kv := strings.Split(queryStrs[1], "=")
					queryMulti[kv[0]] = append(queryMulti[kv[0]], kv[1])
				}
			}
			if len(apiTestData.Body) > 0 {
				var body map[string]interface{}
				body = make(map[string]interface{})
				err = json.Unmarshal([]byte(apiTestData.Body), &body)
				if err != nil {
					Logger.Error("%s", err)
					return
				}
				for k, v := range body {
					bodyMulti[k] = append(bodyMulti[k], v)
				}
			}
		}
	} else {
		var apiFuzzingDatas []ApiFuzzingData
		models.Orm.Table("api_fuzzing_data").Where("id in (?)", ids).Find(&apiFuzzingDatas)
		for index, apiFuzzingData := range apiFuzzingDatas {
			if index == 0 {
				sceneData.Name = fmt.Sprintf("模糊-%s-%s-%s", apiFuzzingData.ApiModule, apiFuzzingData.ApiDesc, GetRandomStr(4, ""))
			}
			if strings.Contains(apiFuzzingData.UrlQuery, "?") {
				queryStrs := strings.Split(apiFuzzingData.UrlQuery, "?")
				if strings.Contains(queryStrs[1], "&") {
					items := strings.Split(queryStrs[1], "&")
					for _, item := range items {
						kv := strings.Split(item, "=")
						queryMulti[kv[0]] = append(queryMulti[kv[0]], kv[1])
					}
				} else {
					kv := strings.Split(queryStrs[1], "=")
					queryMulti[kv[0]] = append(queryMulti[kv[0]], kv[1])
				}
			}

			if len(apiFuzzingData.Body) > 0 {
				var body map[string]interface{}
				body = make(map[string]interface{})
				err = json.Unmarshal([]byte(apiFuzzingData.Body), &body)
				if err != nil {
					Logger.Error("%s", err)
				}
				for k, v := range body {
					bodyMulti[k] = append(bodyMulti[k], v)
				}

			}
		}

	}

	var apiDefinition ApiDefinition
	models.Orm.Table("api_definition").Where("api_id = ?", apiId).Find(&apiDefinition)

	var apiRelation ApiRelation
	models.Orm.Table("api_relation").Where("api_id = ? and app = ?", apiId, apiDefinition.App).Find(&apiRelation)
	if len(apiRelation.ApiId) == 0 {
		Logger.Debug("apiRelation: %v", apiRelation)

	}

	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", apiDefinition.App).Find(&envConfig)
	if len(envConfig.App) == 0 {
		Logger.Debug("apiRelation: %v", apiRelation)
	}

	sceneData.App = apiDefinition.App
	sceneData.ApiId = apiDefinition.ApiId
	sceneData.FileName = fmt.Sprintf("%s.yml", sceneData.Name)
	sceneData.RunTime = 1

	dataFile.Name = sceneData.Name
	dataFile.Version = 1.0
	dataFile.IsRunPostApis = "no"
	dataFile.IsRunPreApis = "no"
	dataFile.IsParallel = "no"
	dataFile.IsUseEnvConfig = "yes"
	dataFile.ApiId = apiDefinition.ApiId
	dataFile.Env.Protocol = envConfig.Protocol
	dataFile.Env.Host = envConfig.Ip
	dataFile.Env.Prepath = envConfig.Prepath
	dataFile.Api.App = apiDefinition.App
	dataFile.Api.Description = apiDefinition.ApiDesc
	dataFile.Api.Module = apiDefinition.ApiModule
	tmps := strings.Split(apiDefinition.ApiId, "_")
	if len(tmps) == 2 {
		dataFile.Api.Method = tmps[0]
		dataFile.Api.Path = tmps[1]
	} else {
		err = fmt.Errorf("关联接口信息不符合标准格式httpMethod_/path, 实际：%s", apiId)
		Logger.Error("%s", err)
		return
	}

	if len(apiRelation.OutVars) > 0 {
		var outVars map[string]string
		outVars = make(map[string]string)
		err = json.Unmarshal([]byte(apiRelation.OutVars), &outVars)
		if err != nil {
			Logger.Error("%s", err)
		}
		for k, v := range outVars {
			var sceneAssert SceneAssert
			sceneAssert.Type = "output"
			sceneAssert.Source = v
			sceneAssert.Value = k
			dataFile.Assert = append(dataFile.Assert, sceneAssert)
		}
	}

	if len(apiRelation.ParamApis) > 0 {
		dataFile.Api.ParamApis = strings.Split(apiRelation.ParamApis, ",")
	}

	if len(envConfig.Auth) > 0 {
		header := make(map[string]interface{})
		err = json.Unmarshal([]byte(envConfig.Auth), &header)
		if err != nil {
			Logger.Error("%s", err)
		}
		dataFile.Single.Header = header
	}

	if len(apiDefinition.PathVariable) > 0 {
		pathVar := make(map[string]interface{})
		for _, item := range apiDefinition.QueryParameter {
			pathVar[item.Name] = item.ValueType
		}

	}

	dataFile.Multi.Query = queryMulti
	dataFile.Multi.Path = pathMulti
	dataFile.Multi.Body = bodyMulti

	dataYAML, err := yaml.Marshal(dataFile)
	if err != nil {
		Logger.Error("%s", err)
	}

	sceneData.Content = fmt.Sprintf("<pre><code>%s</code></pre>", dataYAML)

	var dbSceneData DbSceneData
	models.Orm.Table("scene_data").Where("api_id = ? and app = ? and name = ?", apiId, apiDefinition.App, sceneData.Name).Find(&dbSceneData)
	if len(dbSceneData.ApiId) == 0 {
		err = models.Orm.Table("scene_data").Create(&sceneData).Error
	} else {
		err = models.Orm.Table("scene_data").Where("id = ?", dbSceneData.Id).Update(&sceneData).Error
	}
	var dataFilePath string
	if mode == "json" {
		dataFilePath = fmt.Sprintf("%s/%s.json", DataBasePath, sceneData.FileName)
	} else {
		dataFilePath = fmt.Sprintf("%s/%s.yml", DataBasePath, sceneData.FileName)
	}

	err = ioutil.WriteFile(dataFilePath, dataYAML, 0644)
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func GetApiClassified(ids, source string) (apiClass map[string][]int, err error) {
	apiClass = make(map[string][]int)
	var tableName string
	var idList, tmpIds []string
	if strings.Contains(ids, ",") {
		tmpIds = strings.Split(ids, ",")
	}

	for _, id := range tmpIds {
		if len(id) != 0 {
			idList = append(idList, id)
		}
	}
	if source == "data" {
		tableName = "api_test_data"
		var apiTestData ApiTestData
		for _, id := range idList {
			s, _ := strconv.Atoi(id)
			models.Orm.Table(tableName).Where("id = ?", s).Find(&apiTestData)

			if len(apiTestData.ApiId) == 0 {
				err = fmt.Errorf("Not Found API Test Data")
				Logger.Error("%s", err)
				return
			}
			apiClass[apiTestData.ApiId] = append(apiClass[apiTestData.ApiId], s)
		}
	} else {
		tableName = "api_fuzzing_data"
		var apiFuzzingData ApiFuzzingData
		for _, id := range idList {
			s, _ := strconv.Atoi(id)
			models.Orm.Table(tableName).Where("id = ?", s).Find(&apiFuzzingData)

			if len(apiFuzzingData.ApiId) == 0 {
				err = fmt.Errorf("Not Found API Fuzzing Data")
				Logger.Error("%s", err)
				return
			}
			apiClass[apiFuzzingData.ApiId] = append(apiClass[apiFuzzingData.ApiId], s)
		}
	}

	return apiClass, err
}

func CreateSceneData(ids, source, mode string) (err error) {
	var apiClass map[string][]int
	apiClass = make(map[string][]int)
	apiClass, err = GetApiClassified(ids, source)
	for k, v := range apiClass {
		if len(v) == 1 {
			err = CreateSingleSceneDataFromData(k, v[0], source, mode)
			if err != nil {
				Logger.Error("%s", err)
			}
		} else {
			err = CreateMultiSceneDataFromData(k, v, source, mode)
			if err != nil {
				Logger.Error("%s", err)
			}
		}
	}
	return
}

func BakOldVer(id, afterTxt, fileName string) (err error) {
	var dataFile DataFile

	if strings.HasSuffix(fileName, ".json") {
		err = json.Unmarshal([]byte(afterTxt), &dataFile)
	} else {
		err = yaml.Unmarshal([]byte(afterTxt), &dataFile)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	filePath := fmt.Sprintf("%s/%s", DataBasePath, fileName)
	_, err = os.Stat(filePath)
	if _, err = os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			dataFile.Version = 1
			var dataInfo []byte
			var err1 error
			if strings.HasSuffix(fileName, ".json") {
				dataInfo, err1 = json.MarshalIndent(dataFile, "", "    ")
			} else {
				dataInfo, err1 = yaml.Marshal(dataFile)
			}

			if err1 != nil {
				Logger.Error("%s", err1)
				return err1
			}
			err2 := ioutil.WriteFile(filePath, dataInfo, 0644)
			if err2 != nil {
				Logger.Error("%s", err2)
				return err2
			}
		}
		return nil
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	var oldSceneFile DataFile
	if strings.HasSuffix(fileName, ".json") {
		err = json.Unmarshal(content, &oldSceneFile)
	} else {
		err = yaml.Unmarshal(content, &oldSceneFile)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	oldVer := oldSceneFile.Version
	newVer := oldSceneFile.Version + 1

	dataFile.Version = newVer
	var dataInfo []byte
	if strings.HasSuffix(fileName, ".json") {
		dataInfo, err = json.MarshalIndent(dataFile, "", "    ")
	} else {
		dataInfo, err = yaml.Marshal(dataFile)
	}

	err = ioutil.WriteFile(filePath, dataInfo, 0644)
	if err != nil {
		Logger.Error("%s", err)
	}

	var dbSceneData DbSceneData
	models.Orm.Table("scene_data").Where("id = ?", id).Find(&dbSceneData)
	dbSceneData.Content = fmt.Sprintf("<pre><code>%s</code></pre>", dataInfo)
	err = models.Orm.Table("scene_data").Where("id = ?", id).Update(&dbSceneData).Error
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	basePath := path.Base(filePath)
	var tmpStr, oldFileName string
	if strings.HasSuffix(fileName, ".json") {
		tmpStr = fmt.Sprintf("_V%.0f.json", oldVer)
		oldFileName = strings.ReplaceAll(basePath, ".json", tmpStr)
	} else {
		tmpStr = fmt.Sprintf("_V%.0f.yml", oldVer)
		oldFileName = strings.ReplaceAll(basePath, ".yml", tmpStr)
	}

	Logger.Debug("oldFileName: %s", oldFileName)

	dataFile.Version = oldVer
	oldVerPath := fmt.Sprintf("%s/%s", OldFilePath, oldFileName)
	var oldData []byte
	if strings.HasSuffix(fileName, ".json") {
		oldData, err = yaml.Marshal(dataFile)
	} else {
		oldData, err = yaml.Marshal(dataFile)
	}

	err = ioutil.WriteFile(oldVerPath, oldData, 0644)
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func SyncSceneData() (newTag, modTag int, err error) {
	dirHandle, err := ioutil.ReadDir(DataBasePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	var fileList []string
	for _, fileH := range dirHandle {
		if !fileH.IsDir() {
			fileName := fileH.Name()
			fileList = append(fileList, fileName)
		}
	}
	var errInfo string
	for index, fileName := range fileList {
		if !strings.HasSuffix(fileName, ".yml") && !strings.HasSuffix(fileName, ".yaml") {
			continue
		}

		var editScene DataFile
		dataFile, err1 := GetSceneContent(fileName)
		if err1 != nil {
			if len(errInfo) == 0 {
				errInfo = fmt.Sprintf("%s, 请核对文件: %s", err, fileName)
			} else {
				err = fmt.Errorf("%s, 请核对文件: %s", err, fileName)
				errInfo = fmt.Sprintf("%s\n%s", errInfo, err)
			}
			continue
		}

		editScene.Name = dataFile.Name
		editScene.Version = dataFile.Version
		editScene.ApiId = dataFile.ApiId
		editScene.IsRunPreApis = dataFile.IsRunPreApis
		editScene.IsRunPostApis = dataFile.IsRunPostApis
		editScene.IsParallel = dataFile.IsParallel
		editScene.IsUseEnvConfig = dataFile.IsUseEnvConfig
		editScene.Env = dataFile.Env
		editScene.Api = dataFile.Api
		editScene.Single = dataFile.Single
		editScene.Multi = dataFile.Multi
		editScene.Assert = dataFile.Assert
		dataYAML, err := yaml.Marshal(editScene)
		if err != nil {
			if len(errInfo) == 0 {
				errInfo = fmt.Sprintf("%s, 请核对文件: %s", err, fileName)
			} else {
				err = fmt.Errorf("%s, 请核对文件: %s", err, fileName)
				errInfo = fmt.Sprintf("%s\n%s", errInfo, err)
			}
			continue
		}
		if index == 0 {
			Logger.Debug("dataYAML: %v", string(dataYAML))
		}

		var dbSceneData DbSceneData
		var sceneData SceneData
		sceneData.Name = dataFile.Name
		sceneData.ApiId = dataFile.ApiId
		sceneData.App = dataFile.Api.App
		sceneData.FileName = fileName
		sceneData.Content = fmt.Sprintf("<pre><code>%s</code></pre>", dataYAML)
		models.Orm.Table("scene_data").Where("api_id = ? and app = ? and name = ?", dataFile.ApiId, dataFile.Api.App, dataFile.Name).Find(&dbSceneData)

		if len(sceneData.Name) == 0 {
			continue
		}

		if (dbSceneData == DbSceneData{}) {
			sceneData.RunTime = 1
			sceneData.App = dataFile.Api.App
			err = models.Orm.Table("scene_data").Create(&sceneData).Error
			newTag++
		} else {
			sceneData.App = dbSceneData.App
			err = models.Orm.Table("scene_data").Where("id = ?", dbSceneData.Id).Update(&sceneData).Error
			modTag++
		}
	}

	if len(errInfo) > 0 {
		err = fmt.Errorf("%s", errInfo)
	}

	return
}

func CopySceneData(id, userName string) (err error) {
	var dbSceneData DbSceneData
	models.Orm.Table("scene_data").Where("id = ?", id).Find(&dbSceneData)
	if len(dbSceneData.Name) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", id)
		Logger.Error("%s", err)
		return
	}

	var sceneData SceneData
	sceneData.App = dbSceneData.App
	sceneData.Name = fmt.Sprintf("%s_复制", dbSceneData.Name)
	sceneData.RunTime = dbSceneData.RunTime
	copyName := strings.Replace(dbSceneData.FileName, ".yml", "_复制.yml", -1)
	sceneData.FileName = copyName
	sceneData.Remark = dbSceneData.Remark
	sceneData.ApiId = dbSceneData.ApiId
	sceneData.Content = dbSceneData.Content
	sceneData.UserName = userName
	err = models.Orm.Table("scene_data").Create(sceneData).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func UpdateApiDefVer(id string) (err error) {
	var apiDef ApiDefinition
	models.Orm.Table("api_definition").Where("id = ?", id).Select("version").Find(&apiDef)
	apiDef.Version = apiDef.Version + 1
	err = models.Orm.Table("api_definition").Where("id = ?", id).UpdateColumn(ApiDefinition{Version: apiDef.Version}).Error
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	return
}

// 执行场景，如果数据文件不存在，则进行创建
func CreateFileByFileName(filePath string) (err error) {
	// 如果数据文件为历史数据，则不进行数据文件的创建
	fileName := path.Base(filePath)
	timeReg, err := regexp.Compile(`_[0-9]{8}_[0-9]{6}\.[0-9]+\.`)
	if err != nil {
		Logger.Error("%v", err)
	}
	timeMatch := timeReg.FindString(fileName)

	if len(timeMatch) > 0 {
		return
	}

	var sceneData SceneData
	var content string
	models.Orm.Table("scene_data").Where("file_name = ?", fileName).Find(&sceneData)
	if len(sceneData.FileName) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", fileName)
		Logger.Error("%s", err)
		return
	}
	content = GetStrFromHtml(sceneData.Content)

	var dataFile DataFile
	var dataInfo []byte
	var err1 error
	if strings.HasSuffix(fileName, ".json") {
		err = json.Unmarshal([]byte(content), &dataFile)
		dataInfo, err1 = json.MarshalIndent(dataFile, "", "    ")
	} else {
		err = yaml.Unmarshal([]byte(content), &dataFile)
		dataInfo, err1 = yaml.Marshal(dataFile)
	}

	if err1 != nil {
		Logger.Error("%s", err1)
		return err1
	}
	err2 := ioutil.WriteFile(filePath, dataInfo, 0644)
	if err2 != nil {
		Logger.Error("%s", err2)
		return err2
	}

	return
}
