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
	"sync"
)

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
	return
}

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
	sceneData.Content = fmt.Sprintf("%s", dataInfo)

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

func WriteContent2File(fileName, content string) (err error) {
	filePath := fmt.Sprintf("%s/%s", DataBasePath, fileName)
	err = ioutil.WriteFile(filePath, []byte(content), 0644)
	//_, err = os.Stat(filePath)
	//if _, err = os.Stat(filePath); err != nil {
	//	if os.IsNotExist(err) {
	//		err = ioutil.WriteFile(filePath, []byte(content), 0644)
	//	}
	//}
	return
}

func GetFileHistoryVersion(fileName, fileType string) (vOld int) {
	cmdStr := fmt.Sprintf("ls -lt %s/%s/%s_V*.%s | head -n 1 | awk '{print $9}'", OldFilePath, fileName, fileName, fileType)
	oldFileName := ExecCommandWithOutput(cmdStr)
	if len(oldFileName) == 0 {
		return
	}

	tmps := strings.Split(oldFileName, "_V")
	tStr := tmps[len(tmps)-1]
	suffixStr := fmt.Sprintf(".%s", fileType)
	vStr := strings.TrimSuffix(tStr, suffixStr)

	vOld, err := strconv.Atoi(vStr)
	if err != nil {
		Logger.Error("%v", err)
	}

	return
}

func BakOldVer(id, content, fileName string) (err error) {
	filePath := fmt.Sprintf("%s/%s", DataBasePath, fileName)

	fileTypeTmp := strings.Split(fileName, ".")
	fileType := fileTypeTmp[len(fileTypeTmp)-1]
	fileNameOnly := fileName[:len(fileName)-len(fileType)-1]

	curVer := GetFileHistoryVersion(fileNameOnly, fileType) + 1

	// 判断目录是否存在，若不存在，则自动新建
	targetDirName := fmt.Sprintf("%s/%s", OldFilePath, fileNameOnly)
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

	var historyFilePath string
	switch fileType {
	case "yml", "yaml", "json":
		strReg := regexp.MustCompile(`version\:\s{1}([0-9.]+)`)
		strMatch := strReg.FindAllSubmatch([]byte(content), -1)

		for _, item := range strMatch {
			rawStr := string(item[0])
			newVer := fmt.Sprintf("version: %d", curVer)
			content = strings.Replace(content, rawStr, newVer, 1)

			var dbContents []string
			models.Orm.Table("scene_data").Where("id = ?", id).Pluck("content", &dbContents)
			if len(dbContents) == 0 {
				err = fmt.Errorf("未找到[%v]数据，请核对", id)
				Logger.Error("%s", err)
				return
			}
			dbContent := dbContents[0]
			dbContent = strings.Replace(dbContent, rawStr, newVer, 1)
			models.Orm.Table("scene_data").Where("id = ?", id).Update("content", dbContent)
			break
		}
	}

	historyFilePath = fmt.Sprintf("%s/%s/%s_V%d.%s", OldFilePath, fileNameOnly, fileNameOnly, curVer, fileType)

	// 写到最新的文件中
	errTmp := ioutil.WriteFile(filePath, []byte(content), 0644)
	if errTmp != nil {
		Logger.Error("%s", errTmp)
		return errTmp
	}

	// 写到old历史版本中
	errTmp = ioutil.WriteFile(historyFilePath, []byte(content), 0644)
	if errTmp != nil {
		Logger.Error("%s", errTmp)
		return errTmp
	}

	return
}

func ModifyEditedData(id, fileName string) (err error) {
	var dbData SceneDataRecord
	models.Orm.Table(" scene_data_test_history").Where("id = ?", id).Find(&dbData)
	dirName := GetHistoryDataDirName(fileName)
	dbData.Content = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, fileName, fileName)
	err = models.Orm.Table("scene_data_test_history").Where("id = ?", id).Update(&dbData).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

// 功能先进行屏蔽
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
	var copyName string
	sceneData.App = dbSceneData.App
	sceneData.Name = fmt.Sprintf("%s_复制", dbSceneData.Name)
	sceneData.RunTime = dbSceneData.RunTime
	tmpList := strings.SplitN(dbSceneData.FileName, ".", 2)
	if len(tmpList) == 2 {
		copyName = fmt.Sprintf("%s_复制.%s", tmpList[0], tmpList[1])
	} else {
		copyName = fmt.Sprintf("%s_复制.yml", dbSceneData.Name)
	}

	sceneData.FileName = copyName
	sceneData.Remark = dbSceneData.Remark
	sceneData.ApiId = dbSceneData.ApiId
	sceneData.Content = dbSceneData.Content
	sceneData.UserName = userName
	if dbSceneData.FileType == 0 {
		sceneData.FileType = 1
	} else {
		sceneData.FileType = dbSceneData.FileType
	}

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
func InitDataFileByName(filePath string) (rawFilePath string, fileType int, err error) {
	// 如果数据文件为历史数据，则不进行数据文件的创建
	fileName := path.Base(filePath)

	var prefixName, newFileName string
	if strings.HasSuffix(fileName, ".log") {
		contentLog, err2 := ioutil.ReadFile(filePath)
		if err2 != nil {
			Logger.Error("%s", err2)
			return "", fileType, err2
		}
		tmpList := strings.Split(string(contentLog), "\n")
		if len(tmpList) > 0 {
			if strings.HasPrefix(tmpList[0], "cmd:") {
				splitFilePath := strings.Split(tmpList[0], " ")
				rawFilePath = splitFilePath[len(splitFilePath)-1]
				fileName = path.Base(rawFilePath)
				filePath = rawFilePath
			} else {
				prefixName = GetHistoryDataDirName(fileName)
			}
		} else {
			prefixName = GetHistoryDataDirName(fileName)
		}
	}

	timeReg, err := regexp.Compile(`_[0-9]{8}_[0-9]{6}\.[0-9]+\.`)
	if err != nil {
		Logger.Error("%v", err)
	}

	timeMatch := timeReg.FindString(fileName)
	if len(timeMatch) > 0 {
		newFileName = strings.Replace(fileName, timeMatch, ".", -1)
	}

	var sceneData SceneData
	var content string
	var byteContent []byte
	if len(prefixName) > 0 {
		models.Orm.Table("scene_data").Where("file_name like ?", prefixName+"%").Find(&sceneData)
	} else if len(newFileName) > 0 {
		models.Orm.Table("scene_data").Where("file_name = ?", newFileName).Find(&sceneData)
	} else {
		models.Orm.Table("scene_data").Where("file_name = ?", fileName).Find(&sceneData)
	}

	if len(sceneData.FileName) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", fileName)
		Logger.Error("%s", err)
		return
	} else {
		rawFilePath = fmt.Sprintf("%s/%s", DataBasePath, sceneData.FileName)
	}

	fileType = sceneData.FileType

	switch fileType {
	case 1:
		if strings.Contains(sceneData.Content, "<pre><code>") {
			content = strings.Replace(sceneData.Content, "<pre><code>", "", -1)
			content = strings.Replace(content, "</code></pre>", "", -1)
		} else {
			content = sceneData.Content
		}

		var dataFile DataFile

		if strings.HasSuffix(fileName, ".json") {
			err = json.Unmarshal([]byte(content), &dataFile)
			byteContent, err = json.MarshalIndent(dataFile, "", "    ")
		} else {
			byteContent = []byte(content) // YML未进行格式化，无需重新校正
			//err = yaml.Unmarshal([]byte(content), &dataFile)
			//byteContent, err = yaml.Marshal(dataFile)
		}

		if err != nil {
			Logger.Error("%s", err)
			return
		}

	default:
		byteContent = []byte(sceneData.Content)
	}

	// 当原文件是历史文件时，不对历史原文件进行覆盖
	if len(newFileName) == 0 {
		err = ioutil.WriteFile(filePath, byteContent, 0644)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	return
}

func (ds DbScene) GetHistoryApiList(lastFile, batchTag string) (apiStr, lastFileAfter string) {
	tmps := strings.Split(lastFile, "/")
	lastFileName := tmps[len(tmps)-1]
	if len(lastFile) > 0 {
		if strings.HasSuffix(lastFileName, ".yml") {
			lastFileAfter = strings.Replace(tmps[len(tmps)-1], ".yml", fmt.Sprintf("_%s.yml", batchTag), -1)
		} else if strings.HasSuffix(lastFileName, ".json") {
			lastFileAfter = strings.Replace(tmps[len(tmps)-1], ".json", fmt.Sprintf("_%s.json", batchTag), -1)
		}
	} else {
		lastFileAfter = " " // 用空格字符串刷新数据
	}

	//apiList := GetListFromHtml(ds.ApiList)
	apiList := strings.Split(ds.ApiList, ",")
	lastFileTag := len(apiList)
	for index, item := range apiList {
		var apiAfter, dirName string
		if item == lastFileName {
			lastFileTag = index
		}
		if index > lastFileTag {
			if index == 0 {
				apiStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", item, item)
			} else {
				apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", apiStr, item, item)
			}
		} else {
			if strings.HasSuffix(item, ".yml") {
				dirName = strings.Split(item, ".yml")[0]
				apiAfter = strings.Replace(item, ".yml", fmt.Sprintf("_%s.yml", batchTag), -1)
			} else if strings.HasSuffix(item, ".json") {
				dirName = strings.Split(item, ".json")[0]
				apiAfter = strings.Replace(item, ".json", fmt.Sprintf("_%s.json", batchTag), -1)
			}
			if index == 0 {
				apiStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s/%s\">%s</a>", dirName, apiAfter, apiAfter)
			} else {
				apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s/%s\">%s</a>", apiStr, dirName, apiAfter, apiAfter)
			}
		}
	}
	return
}

func (playbook Playbook) RunPlaybook(playbookId, mode, source string, dbProduct DbProduct) (result, lastFile string, err error) {
	if len(dbProduct.Name) > 0 {
		playbook.Product = dbProduct.Name
	}

	var runApis []string
	var tag int

	// 数据置为最初状态
	if mode != "continue" {
		runApis = playbook.Apis
		tag = 0
	} else {
		if len(playbook.LastFile) != 0 {
			index := GetSliceIndex(playbook.Apis, playbook.LastFile)
			if index != -1 {
				runApis = playbook.Apis[index:]
				tag = index
			} else {
				runApis = playbook.Apis
			}
		}
	}

	envType := GetEnvTypeByName(playbook.Product)
	isFail := 0
	result = "fail"

	switch playbook.SceneType {
	case 1, 2:
		for k := range runApis {
			playbook.Tag = tag + k
			subResult, historyApi, errTmp := playbook.RunPlaybookContent(envType, source)
			if errTmp != nil {
				if err != nil {
					err = fmt.Errorf("%s; %s", err, errTmp)
				} else {
					err = errTmp
				}
			}

			playbook.HistoryApis = append(playbook.HistoryApis, historyApi)
			playbook.LastFile = historyApi
			if subResult == "fail" {
				errTmp = playbook.WritePlaybookResult(playbookId, subResult, source, envType, errTmp)
				if errTmp != nil {
					Logger.Error("%s", errTmp)
					if err != nil {
						err = fmt.Errorf("%s; %s", err, errTmp)
					} else {
						err = errTmp
					}
				}
				lastFile = runApis[k]
				return
			}
		}
	case 3:
		for k := range runApis {
			playbook.Tag = tag + k
			subResult, historyApi, errTmp := playbook.RunPlaybookContent(envType, source)
			if errTmp != nil {
				Logger.Error("%v", errTmp)
				if err != nil {
					err = errTmp
				} else {
					err = fmt.Errorf("%v; %v", err, errTmp)
				}
			}
			playbook.HistoryApis = append(playbook.HistoryApis, historyApi)
			if subResult == "fail" {
				isFail++
			}
		}

		if isFail > 0 {
			tmpResult := "fail"
			errTmp := playbook.WritePlaybookResult(playbookId, tmpResult, source, envType, err) // 串行继续时，无最近执行的文件
			if errTmp != nil {
				Logger.Error("%v", err)
				if err != nil {
					err = errTmp
				} else {
					err = fmt.Errorf("%v; %v", err, errTmp)
				}
				return
			}
			return // 如果验证全部执行完后，存在失败项，则不再进行后续操作
		}
	case 4, 5:
		wg := sync.WaitGroup{}
		for k := range runApis {
			wg.Add(1)
			go func(inPlaybook Playbook, id string, startIndex, index, envType int, errIn error) {
				inPlaybook.Tag = startIndex + index
				subResult, historyApi, errTmp := inPlaybook.RunPlaybookContent(envType, source)
				if errTmp != nil {
					Logger.Error("%v", errTmp)
					if errIn != nil {
						errIn = errTmp
					} else {
						errIn = fmt.Errorf("%v; %v", errIn, errTmp)
					}
				}
				inPlaybook.HistoryApis = append(inPlaybook.HistoryApis, historyApi)
				if subResult == "fail" {
					isFail++
				}
				if errIn != nil {
					Logger.Error("%v", errIn)
				}
			}(playbook, playbookId, tag, k, envType, err)
		}

		wg.Wait()
		if isFail > 0 {
			tmpResult := "fail"
			errTmp := playbook.WritePlaybookResult(playbookId, tmpResult, source, envType, err) // 并发模式时，无最近执行的文件
			if errTmp != nil {
				Logger.Error("%v", err)
				if err != nil {
					err = errTmp
				} else {
					err = fmt.Errorf("%v; %v", err, errTmp)
				}
				return
			}
			return // 如果验证全部执行完后，存在失败项，则不再进行后续操作
		}
	}

	if err == nil && isFail == 0 {
		result = "pass"
		lastFile = ""
	}

	if playbook.SceneType == 2 || playbook.SceneType == 5 {
		//Logger.Debug("开始比较")
		result, err = CompareResult(playbook.HistoryApis, mode)
	}
	playbook.LastFile = lastFile
	err = playbook.WritePlaybookResult(playbookId, result, source, envType, err)
	if err != nil {
		Logger.Error("%v", err)
		return
	}

	if result != "pass" {
		err = fmt.Errorf("测试 %v", result)
	}
	return
}

func (ds DbScene) GetPlaybook() (playbook Playbook) {
	var filePaths []string
	var filePath string
	//fileNames := GetListFromHtml(ds.ApiList)
	fileNames := strings.Split(ds.ApiList, ",")
	for _, fileName := range fileNames {
		if len(fileName) == 0 {
			continue
		}
		fileName = strings.Replace(fileName, "\n", "", -1)
		filePath = fmt.Sprintf("%s/%s", DataBasePath, fileName)
		filePaths = append(filePaths, filePath)
	}

	playbook.Apis = filePaths
	playbook.Name = ds.Name
	playbook.LastFile = ds.LastFile
	playbook.Product = ds.Product
	playbook.SceneType = ds.SceneType
	return
}

func GetAfterContent(lang, in string, depOutVars map[string][]interface{}) (out string, err error) {
	var allFalseCount int
	var notDefVars, tmpDefVars map[string]string
	in, notDefVars, falseCount, err := GetIndexStr(lang, in, "env:", "api:", depOutVars)
	allFalseCount = allFalseCount + falseCount

	in, tmpDefVars, falseCount, errTmp := GetIndexStr(lang, in, "single:", "multi:", depOutVars)
	allFalseCount = allFalseCount + falseCount
	for k, v := range tmpDefVars {
		notDefVars[k] = v
	}

	if errTmp != nil {
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}

	}

	in, tmpDefVars, falseCount, errTmp = GetIndexStr(lang, in, "multi:", "action:", depOutVars)
	allFalseCount = allFalseCount + falseCount
	for k, v := range tmpDefVars {
		notDefVars[k] = v
	}

	if errTmp != nil {
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}

	}

	in, tmpDefVars, _, _ = GetIndexStr(lang, in, "action:", "assert:", depOutVars)
	//allFalseCount = allFalseCount + falseCount  // 如果是自身的变量，无法替换，不进入未定义拦截
	//for k, v := range tmpDefVars {
	//	notDefVars[k] = v
	//}
	if len(tmpDefVars) > 0 {
		Logger.Warning("action part not def vars: %v", tmpDefVars)
	}

	in, tmpDefVars, _, _ = GetIndexStr(lang, in, "assert:", "output:", depOutVars) // 因为output为自动生成的数据，初始化时，为空
	//allFalseCount = allFalseCount + falseCount                                 // 如果是断言值模板的定义，先不进入未定义拦截，后续优化
	//for k, v := range tmpDefVars {
	//	notDefVars[k] = v
	//}

	if len(tmpDefVars) > 0 {
		Logger.Warning("assert part not def vars: %v", tmpDefVars)
	}

	out = in

	if allFalseCount > 0 {
		if err != nil {
			if len(notDefVars) > 0 {
				err = fmt.Errorf("%s; 存在未定义参数: %v，请先定义或关联", err, notDefVars)
			}
		} else {
			if len(notDefVars) > 0 {
				err = fmt.Errorf("存在未定义参数: %v，请先定义或关联", notDefVars)
			}
			Logger.Error("%s", err)
		}

		return
	}

	return
}

func GetStrByIndex(rawStr, startStr, endStr string) (indexStr, targetStr string, startIndex, endIndex int) {
	if len(startStr) == 0 && len(endStr) == 0 {
		indexStr = rawStr
	} else {
		startIndex = strings.Index(rawStr, startStr)
		endIndex = strings.LastIndex(rawStr, endStr)
		if startIndex == -1 { // 开始未找到，相当于无相关定义，可直接跳过
			targetStr = rawStr
			return
		} else if endIndex == -1 {
			indexStr = rawStr[startIndex:]
		} else if startIndex > endIndex {
			Logger.Debug("rawStr: %s", rawStr)
			Logger.Debug("startStr: %s，endStr: %s", startStr, endStr)
			Logger.Error("rawStr[%d:%d], 索引有问题，请校对", startIndex, endIndex)
		} else {
			indexStr = rawStr[startIndex:endIndex]
		}
	}

	return
}

// 分区匹配和替换  env匹配区，single匹配区，multi匹配区，断言匹配区
func GetIndexStr(lang, rawStr, startStr, endStr string, depOutVars map[string][]interface{}) (targetStr string, notDefVars map[string]string, falseCount int, err error) {
	var indexStr string
	var startIndex, endIndex int
	indexStr, targetStr, startIndex, endIndex = GetStrByIndex(rawStr, startStr, endStr)
	if startIndex == -1 { //开始为-1说明没找到开始的信息，应该原样返回
		Logger.Warning("未找到开始值: %s, 结束值: %s", startStr, endStr)
	}

	strReg := regexp.MustCompile(`\{([-a-zA-Z0-9_]+)(\[(\W*\d+)\])*\}`)
	strMatch := strReg.FindAllSubmatch([]byte(indexStr), -1)

	//支持${xx}写法不替换
	strReg2 := regexp.MustCompile(`\$\{([-a-zA-Z0-9_]+)\}`)
	strMatch2 := strReg2.FindAllSubmatch([]byte(indexStr), -1)

	countMap := make(map[string]int) // 统计同一key出现的次数，按顺取用，如果引用次数超过定义次数，超过的取第0个值

	for _, item := range strMatch {
		key := string(item[1])
		if key == "self" {
			continue
		}

		rawStrDef := string(item[0])
		isReplace := 1
		for _, subItem := range strMatch2 {
			subRawStrDef := string(subItem[0])
			if rawStrDef == subRawStrDef {
				isReplace = 0
				break
			}
		}

		if isReplace == 0 {
			continue
		}

		if len(item[2]) > 0 {
			order, _ := strconv.Atoi(string(item[3]))
			if value, ok := depOutVars[key]; ok {
				if len(value) > order {
					var tmpKey string
					if order < 0 {
						tmpKey = Interface2Str(value[len(value)+order])
					} else {
						tmpKey = Interface2Str(value[order])
					}
					indexStr = strings.Replace(indexStr, rawStrDef, tmpKey, -1)
				} else {
					errTmp := fmt.Errorf("参数: %s定义参数不足%v，%s取值超出索引，请核对~", string(item[1]), value, rawStrDef)
					Logger.Error("%s", errTmp)
					if err != nil {
						err = fmt.Errorf("%s;%s", err, errTmp)
					} else {
						err = errTmp
					}
					falseCount++
				}
			}
		} else {
			if value, ok := depOutVars[key]; ok {
				if len(value) == 0 {
					Logger.Error("%s:%v", key, value)
					falseCount++
					continue
				}

				var vStr string
				if count, ok := countMap[key]; ok {
					if len(value) > count {
						vStr = Interface2Str(value[count])
					} else {
						vStr = Interface2Str(value[0])
					}
				} else {
					allValueDef1 := fmt.Sprintf("'*%s*'", rawStrDef)
					allValueDef2 := fmt.Sprintf("\"*%s*\"", rawStrDef)
					allValueDef3 := fmt.Sprintf("**%s**", rawStrDef)
					if strings.Contains(indexStr, allValueDef1) {
						rawStrDef = allValueDef1
						vStr = Interface2Str(value)
					} else if strings.Contains(indexStr, allValueDef2) {
						rawStrDef = allValueDef2
						vStr = Interface2Str(value)
					} else if strings.Contains(indexStr, allValueDef3) {
						rawStrDef = allValueDef3
						vStr = Interface2Str(value)
					} else {
						vStr = Interface2Str(value[0]) // 未定义，默认取第1个值
					}
				}

				indexStr = strings.Replace(indexStr, rawStrDef, vStr, 1)
				countMap[key] = countMap[key] + 1
			} else {
				value := GetValueFromSysParameter(lang, key)
				if len(value) > 0 {
					indexStr = strings.Replace(indexStr, rawStrDef, value, 1)
				} else {
					if strings.Contains(startStr, "assert:") || len(startStr) == 0 {
						value, errTmp := GetAssertTemplateValue(lang, key)
						if errTmp != nil {
							falseCount++
							if err != nil {
								err = fmt.Errorf("%s;%s", err, errTmp)
							} else {
								err = errTmp
							}
						}
						indexStr = strings.Replace(indexStr, rawStrDef, value, -1)
					}
				}
			}
		}
	}

	// Single下支持{self}占位符的使用，Multi不再支持，后续考虑停用该变量
	strReg3 := regexp.MustCompile(`(([-a-zA-Z0-9_]+)\:\s{1}['"]{1}\{self\}['"]{1})`)
	strMatch3 := strReg3.FindAllSubmatch([]byte(indexStr), -1)
	for _, item := range strMatch3 {
		key := string(item[2])
		rawStrDef := string(item[0])
		if value, ok := depOutVars[key]; ok {
			var vStr string
			vStr = Interface2Str(value[0])
			tmpStrDef := strings.Replace(rawStrDef, "{self}", vStr, -1)
			indexStr = strings.Replace(indexStr, rawStrDef, tmpStrDef, 1)
		} else {
			value := GetValueFromSysParameter(lang, key)
			if len(value) > 0 {
				tmpStrDef := strings.Replace(rawStrDef, "{self}", value, -1)
				indexStr = strings.Replace(indexStr, rawStrDef, tmpStrDef, 1)
			}
		}
	}

	indexStr = GetSpecialStr(lang, indexStr)

	indexStr, tmpCount := GetTreeData(indexStr)
	falseCount = falseCount + tmpCount

	notDefVars, _ = GetInStrDef(indexStr)
	if len(notDefVars) > 0 {
		falseCount = falseCount + len(notDefVars)
	}

	if len(startStr) == 0 && len(endStr) == 0 {
		targetStr = indexStr
	} else {
		if endIndex == -1 {
			targetStr = rawStr[:startIndex] + indexStr
		} else {
			targetStr = rawStr[:startIndex] + indexStr + rawStr[endIndex:]
		}
	}

	return
}

func GetTreeData(in string) (out string, falseCount int) {
	// 匹配字符串
	strReg := regexp.MustCompile(`\{TreeData_([a-zA-Z0-9]+)\[(\d+)\]\}`)
	strMatch := strReg.FindAllSubmatch([]byte(in), -1)

	var first, second, third string
	for _, item := range strMatch {
		rawStrDef := string(item[0])
		treeDataKey := string(item[1])
		deepStr := string(item[2])
		var err error
		deep, err := strconv.Atoi(deepStr)
		if err != nil {
			Logger.Error("%s", err)
		}
		if deep == 1 {
			first, second, third = GetTreeDataValue(treeDataKey, deep, first, second)
			in = strings.ReplaceAll(in, rawStrDef, first)
			break
		}
	}

	for _, item := range strMatch {
		rawStrDef := string(item[0])
		treeDataKey := string(item[1])
		deepStr := string(item[2])
		var err error
		deep, err := strconv.Atoi(deepStr)
		if err != nil {
			Logger.Error("%s", err)
		}
		if deep == 2 {
			first, second, third = GetTreeDataValue(treeDataKey, deep, first, second)
			in = strings.ReplaceAll(in, rawStrDef, second)
			break
		}
	}

	for _, item := range strMatch {
		rawStrDef := string(item[0])
		treeDataKey := string(item[1])
		deepStr := string(item[2])
		var err error
		deep, err := strconv.Atoi(deepStr)
		if err != nil {
			Logger.Error("%s", err)
		}
		if deep == 3 {
			first, second, third = GetTreeDataValue(treeDataKey, deep, first, second)
			in = strings.ReplaceAll(in, rawStrDef, third)
			break
		}
	}
	out = in
	return
}

func UpdateDataAssertContent() {
	var dataList []DbSceneData
	var err error
	models.Orm.Table("scene_data").Find(&dataList)
	count := 0
	errFileCount := 0
	doubleCount := 0
	lastCount := 0
	for _, dataInfo := range dataList {
		isChange := 0
		if len(dataInfo.Content) == 0 {
			err := fmt.Errorf("无断言信息%s", dataInfo.Name)
			Logger.Error("%s", err)
			continue
		}

		var content string
		content = dataInfo.Content
		if strings.Contains(content, "<pre><code>") {
			content = strings.Replace(content, "<pre><code>", "", -1)
			content = strings.Replace(content, "</code></pre>", "", -1)
		}

		var df DataFile
		if strings.HasSuffix(dataInfo.FileName, ".json") {
			err = json.Unmarshal([]byte(content), &df)
		} else if strings.HasSuffix(dataInfo.FileName, ".yml") || strings.HasSuffix(dataInfo.FileName, ".yaml") {
			err = yaml.Unmarshal([]byte(content), &df)
		} else {
			continue
		}
		if err != nil {
			errFileCount++
			continue
		}

		if len(df.Assert) > 0 {
			for index, item := range df.Assert {
				var newSource string
				if !strings.Contains(item.Source, "*") {
					if strings.Contains(item.Source, "-") {
						newSource = strings.Replace(item.Source, "-", ".", -1)
					} else {
						continue
					}

				}

				starIndex := strings.Index(item.Source, "*")
				starLastIndex := strings.LastIndex(item.Source, "*")
				boundIndex := strings.Index(item.Source, "[")
				if starIndex != starLastIndex {
					if !strings.Contains(item.Source, "[") {
						newSource = strings.Replace(item.Source, "*", "[0].", -1)
						newSource = strings.Replace(newSource, "-", ".", -1)
						Logger.Debug("rawSource: %v", item.Source)
						Logger.Debug("newSource: %v", newSource)
						if isChange == 0 {
							isChange++
						}
					} else {
						Logger.Debug("%d, 双list断言变更数据描述: %v", doubleCount, dataInfo.Name)
						Logger.Debug("rawSource: %v", item.Source)
						if isChange == 0 {
							isChange++
						}
						doubleCount++
						continue
					}
				}

				if starIndex < boundIndex && starIndex != -1 {
					items := strings.Split(item.Source, "[")
					if strings.Contains(items[0], "-") {
						items[0] = strings.Replace(items[0], "-", ".", -1)
					}
					newSource = strings.Replace(items[0], "*", fmt.Sprintf("[%s.", items[1]), 1)
					Logger.Debug("rawSource: %v", item.Source)
					Logger.Debug("newSource: %v", newSource)
					if isChange == 0 {
						isChange++
					}
					count++
				} else if boundIndex < starIndex && boundIndex != -1 {
					newSource = strings.Replace(item.Source, "*", ".", 1)
					if strings.Contains(newSource, "-") {
						newSource = strings.Replace(newSource, "-", ".", -1)
					}
					Logger.Debug("rawSource: %v", item.Source)
					Logger.Debug("newSource: %v", newSource)
					if isChange == 0 {
						isChange++
					}
					lastCount++
				} else if starIndex > 0 && boundIndex == -1 {
					if strings.Contains(item.Source, "-") && strings.Contains(item.Source, "*") {
						newSource = strings.Replace(item.Source, "*", "[0].", -1)
						newSource = strings.Replace(newSource, "-", ".", -1)
					} else if strings.Contains(item.Source, "-") {
						newSource = strings.Replace(item.Source, "-", ".", -1)
					} else if strings.Contains(item.Source, "*") {
						newSource = strings.Replace(item.Source, "*", "[0].", -1)
					}
					Logger.Debug("rawSource: %v", item.Source)
					Logger.Debug("newSource: %v", newSource)
					if isChange == 0 {
						isChange++
					}
				}

				if len(newSource) > 0 {
					df.Assert[index].Source = newSource
				}
			}
		}

		if isChange > 0 {
			Logger.Debug("afterSource: %v", df.Assert)
			afterByte, _ := yaml.Marshal(df)
			dataInfo.Content = string(afterByte)
			models.Orm.Table("scene_data").Where("id = ?", dataInfo.Id).Update(&dataInfo)
		}

	}
	Logger.Debug("单数组数量：%d", count)
	Logger.Debug("多数组数量：%d", doubleCount)
	return
}
