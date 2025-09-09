package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

	switch suffix {
	case ".yaml", ".yml", ".json":
		dst = fmt.Sprintf("%s/%s/%s_%s%s", HistoryBasePath, dirName, dirName, curTime, suffix)
	default:
		dst = fmt.Sprintf("%s/%s/%s_%s.log", HistoryBasePath, dirName, dirName, curTime)
	}

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

func WriteDataResultByFile(src, result, dst, product, source string, envType int, errIn error) (err error) {
	var df DataFile
	var content []byte

	switch source {
	case "task":
		content, err = ioutil.ReadFile(dst)
	case "data", "playbook":
		content, err = ioutil.ReadFile(src)
	default:
		content, err = ioutil.ReadFile(dst)
	}

	if err != nil {
		Logger.Error("src: %s", src)
		Logger.Error("dst: %s", dst)
		Logger.Error("Error: %s", err)

	}

	tmps := strings.Split(src, ".")
	fileType := tmps[len(tmps)-1]

	var sceneDataRecord SceneDataRecord
	var apiStr string

	switch fileType {
	case "yml", "yaml":
		err = yaml.Unmarshal(content, &df)
	case "json":
		err = json.Unmarshal(content, &df)
	}

	if err != nil {
		Logger.Debug("content: %s", string(content))
		Logger.Error("%s", err)
	}

	switch fileType {
	case "yml", "yaml", "json":
		sceneDataRecord.ApiId = df.ApiId
		sceneDataRecord.App = df.Api.App
	default:
		var dbSD DbSceneData
		fileName := path.Base(src)
		models.Orm.Table("scene_data").Where("file_name = ?", fileName).Find(&dbSD)
		sceneDataRecord.ApiId = dbSD.ApiId
		sceneDataRecord.App = dbSD.App
	}

	apiStr = GetLastFileLink(dst)

	sceneDataRecord.Content = apiStr
	if len(df.Name) == 0 {
		dirName := GetHistoryDataDirName(path.Base(dst))
		sceneDataRecord.Name = dirName
	} else {
		sceneDataRecord.Name = df.Name
	}

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

func (apiModel ApiDataSaveModel) WriteDataFileHistoryResult(result, dst string, envType int, errIn error) (err error) {
	var sceneDataRecord SceneDataRecord

	dirName := GetHistoryDataDirName(path.Base(dst))

	apiStr := fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, path.Base(dst), path.Base(dst))

	sceneDataRecord.Content = apiStr

	sceneDataRecord.Name = dirName
	sceneDataRecord.ApiId = fmt.Sprintf("%s_%s", apiModel.Method, apiModel.Path)
	sceneDataRecord.App = apiModel.App
	sceneDataRecord.Result = result
	sceneDataRecord.EnvType = envType
	sceneDataRecord.Product = apiModel.Product

	if errIn != nil {
		sceneDataRecord.FailReason = fmt.Sprintf("%s", errIn)
	}

	err = models.Orm.Table("scene_data_test_history").Create(sceneDataRecord).Error

	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func (apiModel HistorySaveModel) WriteDataFileHistoryResult(result, dst string, envType int, errIn error) (err error) {
	var sceneDataRecord SceneDataRecord

	dirName := GetHistoryDataDirName(path.Base(dst))

	apiStr := fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, path.Base(dst), path.Base(dst))

	sceneDataRecord.Content = apiStr
	sceneDataRecord.Name = fmt.Sprintf("%s-%s-%s", apiModel.Module, apiModel.ApiDesc, apiModel.DataDesc)

	sceneDataRecord.ApiId = fmt.Sprintf("%s_%s", apiModel.Method, apiModel.Path)
	sceneDataRecord.App = apiModel.App
	sceneDataRecord.Result = result
	sceneDataRecord.EnvType = envType
	sceneDataRecord.Product = apiModel.Product

	if errIn != nil {
		sceneDataRecord.FailReason = fmt.Sprintf("%s", errIn)
	}

	err = models.Orm.Table("scene_data_test_history").Create(sceneDataRecord).Error

	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func WriteSceneDataResult(id string, result, dst, product, source string, envType int, errIn error) (err error) {
	var dbSceneData DbSceneData
	s, _ := strconv.Atoi(id)
	if source == "ai_data" {
		models.Orm.Table("ai_data").Where("id = ?", s).Find(&dbSceneData)
	} else {
		models.Orm.Table("scene_data").Where("id = ?", s).Find(&dbSceneData)
		if len(dbSceneData.ApiId) == 0 {
			return
		}
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

	if source == "ai_data" {
		err = models.Orm.Table("ai_data").Where("id = ?", dbSceneData.Id).Update(&dbSceneData).Error
	} else {
		err = models.Orm.Table("scene_data").Where("id = ?", dbSceneData.Id).Update(&dbSceneData).Error
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	err = RecordDataHistory(dst, product, source, envType, dbSceneData)

	return
}

func RunNonStandard(app, rawFilePath, content, logFilePath, product, source string, depOutVars map[string][]interface{}) (result, dst string, err error) {
	header := make(map[string]interface{})

	if len(product) > 0 {
		sceneEnvConfig, errTmp := GetEnvConfig(product, "product")
		if errTmp != nil {
			Logger.Warning("%s", errTmp)
		}

		if depOutVars == nil {
			depOutVars = make(map[string][]interface{})
		}

		depOutVars["host"] = append(depOutVars["host"], sceneEnvConfig.Ip)
		depOutVars["protocol"] = append(depOutVars["protocol"], sceneEnvConfig.Protocol)

		if len(sceneEnvConfig.Auth) > 2 {
			err = json.Unmarshal([]byte(sceneEnvConfig.Auth), &header)
			if err != nil {
				Logger.Error("%s", err)
				return "fail", rawFilePath, err
			}

			for k, v := range header {
				depOutVars[k] = append(depOutVars[k], v)
			}
		}

		dbProductList, err := GetProductInfo(product)
		dbProduct := dbProductList[0]
		if err != nil {
			Logger.Error("%v", err)
		}

		privateParameter := dbProduct.GetPrivateParameter()
		for k, v := range privateParameter {
			if _, ok := depOutVars[k]; !ok {
				depOutVars[k] = append(depOutVars[k], v)
			}
		}
	} else if len(app) > 0 {
		envConfig, err := GetEnvConfig(app, "app")
		if err != nil {
			Logger.Warning("%s", err)
		}

		depOutVars["host"] = append(depOutVars["host"], envConfig.Ip)
		depOutVars["protocol"] = append(depOutVars["protocol"], envConfig.Protocol)

		if len(envConfig.Auth) > 2 {
			err = json.Unmarshal([]byte(envConfig.Auth), &header)
			if err != nil {
				Logger.Error("%s", err)
				return "fail", rawFilePath, err
			}

			for k, v := range header {
				depOutVars[k] = append(depOutVars[k], v)
			}
		}
	}

	lang := GetRequestLangage(header)

	contentStr, notDefVars, falseCount, errTmp := GetIndexStr(lang, string(content), "", "", depOutVars)
	if falseCount > 0 {
		if errTmp != nil {
			err = fmt.Errorf("%s; 存在未定义参数: %s，请先定义或关联", errTmp, notDefVars)
		} else {
			err = fmt.Errorf("存在未定义参数: %s，请先定义或关联", notDefVars)
		}

		Logger.Error("%s", err)
		return "fail", rawFilePath, err
	}

	// 全部更新脚本的内容
	err2 := ioutil.WriteFile(rawFilePath, []byte(contentStr), 0644)
	if err2 != nil {
		Logger.Error("%s", err2)
		return "fail", rawFilePath, err2
	}

	var runEngine string
	suffix := GetStrSuffix(rawFilePath)
	runEngine, err = GetValueFromMapDef("scriptRunEngine", suffix)
	if err != nil {
		runEngine, err = GetScriptRunEngin(rawFilePath)
		if err != nil {
			return "fail", rawFilePath, err
		}
	}

	var outputStr, resultStr, newFilePath string
	strCommand := fmt.Sprintf("%s %s", runEngine, rawFilePath)
	outputStrTmp, errTmp0 := exec.Command(runEngine, rawFilePath).CombinedOutput()
	outputStr = string(outputStrTmp)
	if errTmp0 != nil {
		Logger.Info("cmd: %s", strCommand)
		Logger.Debug("output: \n%s", outputStr)
		Logger.Error("err: %s", errTmp0)
		if err != nil {
			err = fmt.Errorf("%s;%s%s", err, outputStr, errTmp0)
		} else {
			err = fmt.Errorf("%s%s", outputStr, errTmp0)
		}
	}

	switch source {
	case "historyContinue":
		if len(logFilePath) > 0 {
			dst = logFilePath
		}
	}

	if len(dst) == 0 {
		dstTmp, errTmp1 := GetResultFilePath(rawFilePath)
		if errTmp1 != nil {
			Logger.Error("%s", errTmp1)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp1)
			} else {
				err = errTmp1
			}
		}
		dst = dstTmp
	}

	var errTmp2 error
	if err != nil {
		resultStrTmp := fmt.Sprintf("cmd: %s\nerr: %s", strCommand, err)
		resultStr = resultStrTmp
		errTmp2 = ioutil.WriteFile(dst, []byte(resultStr), 0644)
	} else {
		errTmp2 = ioutil.WriteFile(dst, []byte(outputStr), 0644)
	}

	if errTmp2 != nil {
		Logger.Debug("dst: %s", dst)
		Logger.Error("%s", errTmp2)
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp2)
		} else {
			err = errTmp2
		}
	}

	if err != nil {
		result = "fail"
	} else {
		result = "pass"
	}

	if source == "data" || source == "playbook" {
		prefix := GetHistoryDataDirName(path.Base(dst))
		suffix := GetStrSuffix(path.Base(dst))
		newFilePath = fmt.Sprintf("%s/%s%s", DataBasePath, prefix, suffix)

		var errTmp error
		if err != nil {
			errTmp = ioutil.WriteFile(newFilePath, []byte(resultStr), 0644) // 执行失败才记录log
		}

		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%v,%v", err, errTmp)
			} else {
				err = errTmp
			}
		}
	}

	return
}

func (df DataFile) RunStandard(product, filePath, mode, source, dataContent string, depOutVars map[string][]interface{}) (urlStr, headerStr, requestStr, responseStr, outputStr, result, dst string, err error) {
	if depOutVars == nil {
		depOutVars = make(map[string][]interface{})
	}

	//if len(df.Api.PreApi) > 0 && df.IsRunPreApis == "yes" {
	//	for _, preApiFile := range df.Api.PreApi {
	//		preFilePath := fmt.Sprintf("%s/%s", DataBasePath, preApiFile)
	//		Logger.Debug("开始执行前置用例: %v", preFilePath)
	//		result, dst, err = RunStandard(df.Api.App, preFilePath, product, source, depOutVars)
	//		if err != nil {
	//			Logger.Error("%s", err)
	//			return
	//		}
	//		if result == "fail" {
	//			return
	//		}
	//	}
	//}   // 功能待完善

	var envConfig EnvConfig

	envConfig, _ = GetEnvConfig(df.Api.App, "app")

	depOutVarsTmp, err1 := df.GetDepParams()
	if err1 != nil {
		Logger.Error("%s", err1)
		if err != nil {
			err = fmt.Errorf("%s;%s", err, err1)
		} else {
			err = err1
		}
	}

	for k, v := range depOutVarsTmp {
		if _, ok := depOutVars[k]; !ok {
			depOutVars[k] = v
		}
	}

	if len(product) > 0 {
		sceneEnvConfig, errTmp := GetEnvConfig(product, "product")
		if errTmp != nil {
			Logger.Warning("%s", errTmp)
		}
		envConfig.Ip = sceneEnvConfig.Ip
		envConfig.Auth = sceneEnvConfig.Auth
		envConfig.Product = product
		envConfig.Protocol = sceneEnvConfig.Protocol

		dbProductList, err := GetProductInfo(product)
		dbProduct := dbProductList[0]
		if err != nil {
			Logger.Error("%v", err)
		}
		privateParameter := dbProduct.GetPrivateParameter()
		for k, v := range privateParameter {
			if _, ok := depOutVars[k]; !ok {
				depOutVars[k] = append(depOutVars[k], v)
			}
		}
	}

	header, err := df.GetHeader(envConfig)
	df.Single.Header = header
	if err != nil {
		urlStr, headerStr, requestStr, responseStr, outputStr, _ = df.GetResponseStr()
		return
	}

	lang := GetRequestLangage(header)

	var querys, bodys []map[string]interface{}
	var bodyList []interface{}
	var urls []string
	var rHeader map[string]interface{}

	if mode == "again" {
		urls = df.Urls
		if df.Api.Method == "get" && len(df.Request) > 0 {
			for _, item := range df.Request {
				queryDict := make(map[string]interface{})
				errTmp := json.Unmarshal([]byte(item), &queryDict)
				if errTmp != nil {
					Logger.Error("%v", errTmp)
				} else {
					querys = append(querys, queryDict)
				}
			}
		} else {
			if df.Single.BodyList != nil {
				errTmp := json.Unmarshal([]byte(df.Request[0]), &bodyList)
				if errTmp != nil {
					Logger.Error("%v", errTmp)
				}
			} else {
				for _, item := range df.Request {
					bodyDict := make(map[string]interface{})
					errTmp := json.Unmarshal([]byte(item), &bodyDict)
					if errTmp != nil {
						Logger.Error("%v", errTmp)
					} else {
						bodys = append(bodys, bodyDict)
					}
				}
			}
		}
	} else {
		_, errTmp := yaml.Marshal(df)
		if errTmp != nil {
			Logger.Error("%v", errTmp)
			err = errTmp
			urlStr, headerStr, requestStr, responseStr, outputStr, _ = df.GetResponseStr()
			return
		}

		contentStr, errTmp := GetAfterContent(lang, dataContent, depOutVars)
		if strings.Contains(contentStr, "is_var_strong_check: \"no\"") {
			Logger.Warning("%s数据开启参数弱校验，请自行保证所需依赖参数的定义", filePath)
			errTmp = nil
		}
		if errTmp != nil {
			Logger.Debug("afterContent:\n%s", contentStr)
			err = errTmp
			urlStr, headerStr, requestStr, responseStr, outputStr, _ = df.GetResponseStr()
			return
		}
		errTmp = yaml.Unmarshal([]byte(contentStr), &df)
		if errTmp != nil {
			Logger.Debug("\nrawContent: %s", dataContent)
			Logger.Debug("\nafterContent: %s", contentStr)
			Logger.Error("%v", errTmp)
			err = errTmp
			urlStr, headerStr, requestStr, responseStr, outputStr, _ = df.GetResponseStr()
			return
		}

		df.Single.Header = header

		urls, errTmp = df.GetUrl(envConfig)
		if errTmp != nil {
			Logger.Debug("fileName: %s", path.Base(filePath))
			Logger.Error("%v", errTmp)
			err = errTmp
			urlStr, headerStr, requestStr, responseStr, outputStr, _ = df.GetResponseStr()
			return
		}
		df.Urls = urls

		querys = df.GetQuery()

		bodys, bodyList = df.GetBody()
		rHeader = df.Single.RespHeader
	}

	// 后续可优化，有依赖和无依赖进行控制
	go df.CreateDataOrderByKey(lang, filePath, depOutVars) // 无依赖，异步执行生成动作：create_xxx
	_ = df.RecordDataOrderByKey(bodys)                     // 有依赖，同步执行记录动作：record_xxx
	_ = df.ModifyFileWithData(bodys)                       // 有依赖，同步执行模板动作：modify_file

	if err != nil {
		Logger.Error("%s", err)
		urlStr, headerStr, requestStr, responseStr, outputStr, _ = df.GetResponseStr()
		return
	}

	var resList [][]byte
	var errs []error
	tag := 0
	if df.GetIsParallel() { //控制台过来的并发有bug
		wg := sync.WaitGroup{}
		for _, url := range urls {
			if len(querys) > 0 {
				for _, data := range querys {
					dJson, _ := json.Marshal(data)
					if tag == 0 {
						df.Request = []string{string(dJson)}
					} else {
						df.Request = append(df.Request, string(dJson))
					}
					tag++
					wg.Add(1)
					go func(method, url string, data map[string]interface{}, header map[string]interface{}) {
						defer wg.Add(-1)
						res, err := RunHttp(method, url, data, header, rHeader)
						resList = append(resList, res)
						df.Response = append(df.Response, string(res))
						errs = append(errs, err)
					}(df.Api.Method, url, data, header)
				}
			} else if len(bodys) > 0 || len(bodyList) > 0 {
				if len(bodyList) > 0 {
					if len(bodyList) > 0 {
						var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
						readerNew, _ := jsonNew.Marshal(&bodyList)
						df.Request = []string{string(readerNew)}
						res, err := RunHttpJsonList(df.Api.Method, url, bodyList, header)
						if err != nil {
							Logger.Debug("%s", err)
						}
						resList = append(resList, res)
						df.Response = append(df.Response, string(res))
						errs = append(errs, err)
					}
				} else {
					wg.Add(len(bodys)) // 一次把全部需要等待的任务加上
					for _, data := range bodys {
						dJson, _ := json.Marshal(data)
						if tag == 0 {
							df.Request = []string{string(dJson)}
						} else {
							df.Request = append(df.Request, string(dJson))
						}
						tag++
						//wg.Add(1)  // 定时任务执行过程中，会概率性发生panic
						go func(method, url string, data map[string]interface{}, header map[string]interface{}) {
							defer wg.Add(-1)
							res, err := RunHttp(method, url, data, header, rHeader)
							resList = append(resList, res)
							df.Response = append(df.Response, string(res))
							errs = append(errs, err)
						}(df.Api.Method, url, data, header)
					}
				}
			} else {
				df.Request = []string{} // 没有请求参数，默认置空
				wg.Add(1)
				go func(method, url string, header map[string]interface{}) {
					res, err := RunHttp(method, url, nil, header, rHeader)
					resList = append(resList, res)
					df.Response = append(df.Response, string(res))
					errs = append(errs, err)
				}(df.Api.Method, url, header)
			}
			wg.Wait()
		}
	} else {
		for _, url := range urls {
			if len(querys) > 0 {
				for _, data := range querys {
					var tUrl string
					dJson, _ := json.Marshal(data)
					if tag == 0 {
						df.Request = []string{string(dJson)}
					} else {
						df.Request = append(df.Request, string(dJson))
					}
					tag++
					if df.Api.Method == "delete" {
						subTag := 0
						for k, v := range data {
							strV := Interface2Str(v)
							if subTag == 0 {
								tUrl = fmt.Sprintf("%s?%s=%s", url, k, strV)
							} else {
								tUrl = fmt.Sprintf("%s&%s=%s", tUrl, k, strV)
							}
							subTag++
						}
						res, err := RunHttp(df.Api.Method, tUrl, nil, header, rHeader)
						resList = append(resList, res)
						df.Response = append(df.Response, string(res))
						errs = append(errs, err)
					} else {
						res, err := RunHttp(df.Api.Method, url, data, header, rHeader)
						resList = append(resList, res)
						df.Response = append(df.Response, string(res))
						errs = append(errs, err)
					}
					_ = df.SetSleepAction()
				}
			} else if len(bodys) > 0 || len(bodyList) > 0 {
				if len(bodyList) > 0 {
					if len(bodyList) > 0 {
						var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
						readerNew, _ := jsonNew.Marshal(&bodyList)
						df.Request = []string{string(readerNew)}
						res, err := RunHttpJsonList(df.Api.Method, url, bodyList, header)
						if err != nil {
							Logger.Debug("%s", err)
						}
						resList = append(resList, res)
						df.Response = append(df.Response, string(res))
						errs = append(errs, err)
					}
				} else {
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
							df.Request = []string{string(dJson)}
						} else {
							df.Request = append(df.Request, string(dJson))
						}
						tag++
						res, err := RunHttp(df.Api.Method, url, data, header, rHeader)
						resList = append(resList, res)
						df.Response = append(df.Response, string(res))
						errs = append(errs, err)
						_ = df.SetSleepAction()
					}
				}
			} else {
				df.Request = []string{} // 当只有路由时，请求数据默认设置为空
				res, err := RunHttp(df.Api.Method, url, nil, header, rHeader)
				if err != nil {
					Logger.Error("%s", err)
				}
				resList = append(resList, res)
				df.Response = append(df.Response, string(res))
				errs = append(errs, err)
				_ = df.SetSleepAction()
			}
		}
	}

	result, dst, df.Output, err = df.GetResult(source, filePath, resList, depOutVars, errs)

	if result != "pass" {
		for _, item := range errs {
			if item != nil {
				if err != nil {
					err = fmt.Errorf("%s; %s", err, item)
				} else {
					err = fmt.Errorf("%s", item)
				}
			}
		}
	}

	urlStr, headerStr, requestStr, responseStr, outputStr, _ = df.GetResponseStr()

	return
}

func (dbData CommonDataBase) RunDataFile(filePath, product, source string, depOutVars map[string][]interface{}) (result, dst string, err error) {

	var mode string
	if source == "historyAgain" {
		mode = "again"
	}

	switch dbData.FileType {
	case 2, 3, 4, 5, 99:
		var logFilePath string
		if strings.HasSuffix(filePath, ".log") {
			logFilePath = filePath
		}
		rawFilePath := GetNoStandardFilePath(filePath)
		result, dst, err = RunNonStandard(dbData.App, rawFilePath, dbData.Content, logFilePath, product, source, depOutVars)
	default:
		var df DataFile
		if strings.HasSuffix(filePath, ".json") {
			err = json.Unmarshal([]byte(dbData.Content), &df)
		} else {
			err = yaml.Unmarshal([]byte(dbData.Content), &df)
		}

		if err != nil {
			Logger.Debug("filePath: %s", filePath)
			Logger.Debug("content:\n%s", dbData.Content)
			Logger.Error("%s", err)
			return
		}

		if len(dbData.App) > 0 {
			df.Api.App = dbData.App
		}

		_, _, _, _, _, result, dst, err = df.RunStandard(product, filePath, mode, source, dbData.Content, depOutVars)

	}

	return
}

func RepeatRunDataFile(id, product, source string) (err error) {
	dataInfo, appInfo, filePath, err := GetRunTimeData(id, source)
	var envType, maxThreadNum int
	var isThread string
	if len(product) > 0 {
		productList, _ := GetProductInfo(product)
		productTaskInfo := productList[0]
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

					go func(dbData DbSceneData) {
						result, dst, err1 := dbData.RunDataFile(filePath, product, source, nil)
						if err1 != nil {
							err = err1
							err = WriteSceneDataResult(id, result, dst, product, source, envType, err1)
							return
						}
						err = WriteSceneDataResult(id, result, dst, product, source, envType, err1)
						wg.Done()
					}(dataInfo)
					count++
				}
				wg.Wait()
			}
		} else {
			wg := sync.WaitGroup{}
			Logger.Info("共执行次数：%v", dataInfo.RunTime)
			for i := 0; i < dataInfo.RunTime; i++ {
				wg.Add(1)
				go func(dbData DbSceneData) {
					result, dst, err1 := dbData.RunDataFile(filePath, product, source, nil)
					if err1 != nil {
						err = WriteSceneDataResult(id, result, dst, product, source, envType, err1)
						return
					}
					err = WriteSceneDataResult(id, result, dst, product, source, envType, err1)
					wg.Done()
				}(dataInfo)
			}
			wg.Wait()
		}
	} else {
		for i := 0; i < dataInfo.RunTime; i++ {
			if i > 0 {
				Logger.Info("串行模式-执行次数:%d", i+1)
			}

			result, dst, err1 := dataInfo.RunDataFile(filePath, product, source, nil)
			if err1 != nil {
				Logger.Error("\n%s", err1)
				err = err1
			}

			err2 := WriteSceneDataResult(id, result, dst, product, source, envType, err1)
			if err2 != nil {
				Logger.Error("%s", err2)
				if err != nil {
					err = fmt.Errorf("%v; %v", err, err2)
				} else {
					err = err2
				}
				return
			}

			if result != "pass" {
				err = fmt.Errorf("test %s", result)
				return
			}
		}

	}

	return
}

func RunSceneDataOnce(id, product, source string) (err error) {
	dataInfo, _, filePath, err := GetRunTimeData(id, source)
	var envType int
	if len(product) > 0 {
		productList, _ := GetProductInfo(product)
		productTaskInfo := productList[0]
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
	//app, filePath, err := GetFilePath(id, source)
	//if err != nil {
	//	Logger.Error("%s", err)
	//	return
	//}

	//result, dst, err1 := RunDataFile(app, filePath, product, source, nil)
	result, dst, err1 := dataInfo.RunDataFile(filePath, product, source, nil)
	if err1 != nil {
		Logger.Error("\n%s", err1)
		err = err1
	}
	err = WriteSceneDataResult(id, result, dst, product, source, envType, err1)
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

func (df DataFile) GetIsParallel() (b bool) {
	if df.IsParallel == "yes" {
		b = true
	} else {
		b = false
	}
	return
}

func (df DataFile) GetIsRunPreApis() (b bool) {
	if value := df.IsRunPreApis; value == "yes" {
		b = true
	} else {
		b = false
	}
	return
}

func (df DataFile) GetIsRunPostApis() (b bool) {
	if value := df.IsRunPostApis; value == "yes" {
		b = true
	} else {
		b = false
	}
	return
}

func GetCommonHeader(envConfig EnvConfig) (header map[string]interface{}, err error) {
	header = make(map[string]interface{})
	if len(envConfig.Auth) != 0 {
		strTmp := GetStrFromHtml(envConfig.Auth)
		err = json.Unmarshal([]byte(strTmp), &header)
		if err != nil {
			Logger.Error("解析header异常: %s", err)
			return
		}
	}

	return
}

func (df DataFile) GetHeader(envConfig EnvConfig) (header map[string]interface{}, err error) {
	header = make(map[string]interface{})
	headerData := CopyMapInterface(df.Single.Header)

	headerEnv, err := GetCommonHeader(envConfig)

	if df.IsUseEnvConfig == "yes" {
		header = headerEnv
		for k, v := range headerData {
			if _, ok := header[k]; !ok {
				header[k] = v
			}

			if k == "Content-Type" {
				header[k] = v
			}
		}
	} else {
		header = headerData
		for k, v := range headerEnv {
			if _, ok := header[k]; !ok {
				header[k] = v
			}
		}
	}

	return
}

func (df DataFile) GetUrl(envConfig EnvConfig) (rawUrls []string, err error) {
	var rawUrl string
	envInfo := []string{envConfig.Protocol, "://", envConfig.Ip, envConfig.Prepath, df.Api.Path}
	sceneInfo := []string{df.Env.Protocol, "://", df.Env.Host, df.Env.Prepath, df.Api.Path}
	tag := 0
	for i := 0; i < 5; i++ {
		if df.IsUseEnvConfig == "no" { //不使用公共环境，自身数据文件配置优先级高
			if len(sceneInfo[i]) > 0 {
				rawUrl = rawUrl + sceneInfo[i]
			} else if len(envInfo[i]) > 0 {
				rawUrl = rawUrl + envInfo[i]
			} else {
				if i != 3 { //允许未设置公共路径
					tag++
				}
			}
		} else { //使用公共环境，公共配置优先级高
			if len(envInfo[i]) > 0 {
				rawUrl = rawUrl + envInfo[i]
			} else if len(sceneInfo[i]) > 0 {
				rawUrl = rawUrl + sceneInfo[i]
			} else {
				if i != 3 { //允许未设置公共路径
					tag++
				}
			}
		}
	}

	if tag != 0 {
		err1 := fmt.Errorf("环境信息不完善,请检查, URL: %s", rawUrl)
		Logger.Debug("appEnvWithDataEnv: %v", envInfo)
		Logger.Debug("playbookEnvWithDataEnv: %v", sceneInfo)
		err = err1
		return
	}

	if strings.Contains(rawUrl, "{") {
		if df.Single.Path == nil && df.Multi.Path == nil {
			err = fmt.Errorf("未进行Path变量定义，请先定义")
			return
		}
		pathVarsReg := regexp.MustCompile(`{[[:alpha:]]+}`)
		var pathVars []string
		pathVars = pathVarsReg.FindAllString(rawUrl, -1)
		for _, v := range pathVars {
			str1 := v[1 : len(v)-1]
			var tag int
			tag = 0
			if value, ok := df.Single.Path[str1]; ok {
				valueStr := Interface2Str(value)
				rawUrl = strings.Replace(rawUrl, v, valueStr, -1)
				tag = tag + 1
			}

			if values, ok := df.Multi.Path[str1]; ok {
				for _, value := range values {
					valueStr := Interface2Str(value)
					rawUrl = strings.Replace(rawUrl, v, valueStr, -1)
				}
				tag = tag + 1
			}

			if tag == 0 {
				err = fmt.Errorf("未找到Path:%s变量的定义值，请先进行定义", v)
				return
			}
		}
		rawUrls = append(rawUrls, rawUrl)
	} else {
		rawUrls = append(rawUrls, rawUrl)
	}

	return
}

func (df DataFile) GetQuery() (querys []map[string]interface{}) {
	var query map[string]interface{}
	query = make(map[string]interface{})

	if len(df.Single.Query) > 0 {
		for k, v := range df.Single.Query {
			query[k] = v
		}
	}

	if len(df.Multi.Query) > 0 {
		minLen := GetSliceMinLen(df.Multi.Query)
		for i := 0; i < minLen; i++ {
			for k, v := range df.Multi.Query {
				query[k] = v[i]
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

func (df DataFile) GetBody() (bodys []map[string]interface{}, bodyAfterList []interface{}) {
	if df.Api.Method == "get" {
		return
	}

	if df.Single.BodyList != nil {
		bodyAfterList = make([]interface{}, 0, len(df.Single.BodyList))
		bodyAfterList = df.Single.BodyList
	} else {
		var body map[string]interface{}
		if df.Single.Body != nil {
			body = CopyMap(df.Single.Body)
		}

		if len(df.Multi.Body) > 0 {
			minLen := GetSliceMinLen(df.Multi.Body)
			for i := 0; i < minLen; i++ {
				for k, v := range df.Multi.Body {
					body[k] = v[i]
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

func (df DataFile) CreateActionData() (err error) {
	if len(df.Action) > 0 {
		fileName := ""
		dataCount := 0
		//Logger.Debug("开始生成文件")
		for _, item := range df.Action {
			if item.Type == "create_csv" {
				tmpValue := Interface2Str(item.Value)
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
				bodys, _ := df.GetBody()
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
							_ = WriteDataInCommonFile(filePath, keyStr)
							_ = WriteDataInCommonFile(filePath, valueStr)
						} else {
							_ = WriteDataInCommonFile(filePath, valueStr)
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
						_ = WriteDataInCommonFile(filePath, valueStr)

					}
				}
			}
		}
	}

	if len(df.Action) > 0 {
		fileName := ""
		dataCount := 0
		for _, item := range df.Action {
			if item.Type == "create_excel" || item.Type == "create_xls" || item.Type == "create_xlsx" {
				tmpValue := Interface2Str(item.Value)
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
				bodys, _ := df.GetBody()
				if i == 0 {
					for _, item := range bodys {
						var valueList []string
						for k, v := range item {
							keyList = append(keyList, k)
							vStr := Interface2Str(v)
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, keyList)
						_ = WriteDataInXls(filePath, valueList)
					}
				} else {
					for _, item := range bodys {
						var valueList []string
						for _, key := range keyList {
							vStr := Interface2Str(item[key])
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, valueList)
					}
				}
			}
		}
		//Logger.Debug("结束生成文件")
	}

	if len(df.Action) > 0 {
		fileName := ""
		dataCount := 0
		for _, item := range df.Action {
			if item.Type == "create_excel" || item.Type == "create_xls" || item.Type == "create_xlsx" {
				tmpValue := Interface2Str(item.Value)
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
				bodys, _ := df.GetBody()
				if i == 0 {
					for _, item := range bodys {
						var valueList []string
						for k, v := range item {
							keyList = append(keyList, k)
							vStr := Interface2Str(v)
							valueList = append(valueList, vStr)
						}
						_ = WriteDataInXls(filePath, keyList)
						_ = WriteDataInXls(filePath, valueList)
					}
				} else {
					for _, item := range bodys {
						var valueList []string
						for _, key := range keyList {
							vStr := Interface2Str(item[key])
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

func (df DataFile) GetDepParams() (depOutDict map[string][]interface{}, err error) {
	depOutDict = make(map[string][]interface{})
	var tmpDict map[string][]interface{}
	tmpDict = make(map[string][]interface{})
	var fileNames []string

	fileNames = append(df.Api.ParamApis, df.Api.PreApi...)

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

func (df DataFile) GetResult(source, filePath string, res [][]byte, inOutPutDict map[string][]interface{}, errs []error) (result, dst string, outputDict map[string][]interface{}, err error) {
	outputDict = make(map[string][]interface{})
	isPass := 0
	dst, err = GetResultFilePath(filePath)
	if err != nil {
		Logger.Error("获取目标文件目录: %s", err)
		return
	}

	for i := 0; i < len(res); i++ {
		inIsPass := 0
		if i == 0 {
			df.Response = []string{string(res[i])}
			// 请求时有返回 Error 信息，结果设置为失败，不再走后续流程
			if len(errs) > i {
				if errs[i] != nil {
					df.TestResult = []string{"fail"}
					isPass++
					if err != nil {
						err = fmt.Errorf("%v,%v", err, errs[i])
					} else {
						err = errs[i]
					}
					failReasonStr := fmt.Sprintf("%v", err)
					df.FailReason = []string{failReasonStr}
					continue
				}
			}
		} else {
			df.Response = append(df.Response, string(res[i]))
			if len(errs) > i {
				if errs[i] != nil {
					if err != nil {
						err = fmt.Errorf("%v,%v", err, errs[i])
					} else {
						err = errs[i]
					}
					failReason := fmt.Sprintf("%v", err)

					if len(df.TestResult) < i+1 {
						df.TestResult = append(df.TestResult, "fail")
						df.FailReason = append(df.FailReason, failReason)
					} else {
						df.TestResult[i] = "fail"
						df.FailReason[i] = failReason
					}
					isPass++
					continue
				}
			}
		}

		// 未设置断言时，结果设置为成功，不再走后续流程
		if len(df.Assert) == 0 {
			if i == 0 {
				df.TestResult = []string{"pass"}
				df.FailReason = []string{""}
			} else {
				df.TestResult = append(df.TestResult, "pass")
				df.FailReason = append(df.FailReason, "")
			}
			continue
		}

		for _, assert := range df.Assert {
			aType := assert.Type
			// 若返回断言已经失败了，不再进行output动作
			if inIsPass != 0 && (aType == "output" || aType == "output_re") {
				continue
			}

			if assert.Source == "raw" || assert.Source == "ResponseBody" {
				err1 := assert.AssertValueCompare(string(res[i]))
				if err1 != nil {
					Logger.Error("\n%v", err1)
					if err != nil {
						err = fmt.Errorf("%s\n%s", err, err1)
					} else {
						err = err1
					}
					failReason := fmt.Sprintf("%v", err)
					if len(df.TestResult) < i+1 {
						df.TestResult = append(df.TestResult, "fail")
						df.FailReason = append(df.FailReason, failReason)
					} else {
						df.TestResult[i] = "fail"
						df.FailReason[i] = failReason
					}
					isPass++
					inIsPass++
					continue // 遇到失败，进入下一个断言值的校对
				}
			} else if strings.HasPrefix(assert.Source, "File:") || strings.HasPrefix(assert.Source, "FILE:") {
				targetList, errTmp := assert.GetValueFromFile(df.Response[i])
				if errTmp != nil {
					Logger.Error("%v", errTmp)
					if err != nil {
						err = fmt.Errorf("%v;%v", err, errTmp)
					} else {
						err = errTmp
					}
					failReason := fmt.Sprintf("%v", err)
					if len(df.TestResult) < i+1 {
						df.TestResult = append(df.TestResult, "fail")
						df.FailReason = append(df.FailReason, failReason)
					} else {
						df.TestResult[i] = "fail"
						df.FailReason[i] = failReason
					}
					isPass++
					continue
				}

				switch aType {
				case "output", "output_re":
					k := Interface2Str(assert.Value)
					for _, item := range targetList {
						outputDict[k] = append(outputDict[k], item)
					}
				default:
					for _, item := range targetList {
						vStr := Interface2Str(item)
						errTmp = assert.AssertValueCompare(vStr)
						if errTmp != nil {
							Logger.Error("\n%v", errTmp)
							if err != nil {
								err = fmt.Errorf("%v;%v", err, errTmp)
							} else {
								err = errTmp
							}
						}
					}

					if err != nil {
						failReason := fmt.Sprintf("%v", err)

						if len(df.TestResult) < i+1 {
							df.TestResult = append(df.TestResult, "fail")
							df.FailReason = append(df.FailReason, failReason)
						} else {
							df.TestResult[i] = "fail"
							if len(df.FailReason) < i+1 {
								df.FailReason = append(df.FailReason, failReason)
							} else {
								df.FailReason[i] = failReason
							}
						}
						isPass++
						continue
					}
				}
			} else {
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
					Logger.Error("err: %v", errTmp)
					if err != nil {
						err = fmt.Errorf("%v,%v", err, errTmp)
					} else {
						err = errTmp
					}
					failReason := fmt.Sprintf("%v", err)

					if len(df.TestResult) < i+1 {
						df.TestResult = append(df.TestResult, "fail")
						df.FailReason = append(df.FailReason, failReason)
					} else {
						df.TestResult[i] = "fail"
						df.FailReason[i] = failReason
					}
					isPass++
					inIsPass++
					continue
				}

				switch aType {
				case "output":
					keyName, values, err1 := assert.GetOutput(resDict)
					if err1 != nil {
						if err != nil {
							err = fmt.Errorf("%s, %s", err, err1)
						} else {
							err = err1
						}
						failReason := fmt.Sprintf("%v", err)

						if len(df.TestResult) < i+1 {
							df.TestResult = append(df.TestResult, "fail")
							df.FailReason = append(df.FailReason, failReason)
						} else {
							df.TestResult[i] = "fail"
							df.FailReason[i] = failReason
						}
						isPass++
						inIsPass++
						break
					}

					outputDict[keyName] = append(outputDict[keyName], values...)

				case "output_re":
					keyName, values, err1 := assert.GetOutputRe(res[i])
					if err1 != nil {
						Logger.Error("err1: %v", err1)
						if err != nil {
							err = fmt.Errorf("%s, %s", err, err1)
						} else {
							err = err1
						}
						failReason := fmt.Sprintf("%v", err)

						if len(df.TestResult) < i+1 {
							df.TestResult = append(df.TestResult, "fail")
							df.FailReason = append(df.FailReason, failReason)
						} else {
							df.TestResult[i] = "fail"
							df.FailReason[i] = failReason
						}
						isPass++
						inIsPass++
						break
					}
					outputDict[keyName] = append(outputDict[keyName], values...)
				default:
					_, err1 := assert.AssertResult(resDict, inOutPutDict)
					if err1 != nil {
						Logger.Error("\n%v", err1)
						if err != nil {
							err = fmt.Errorf("%s;\n%s", err, err1)
						} else {
							err = err1
						}
						failReason := fmt.Sprintf("%v", err)
						if len(df.TestResult) < i+1 {
							df.TestResult = append(df.TestResult, "fail")
							df.FailReason = append(df.FailReason, failReason)
						} else {
							df.TestResult[i] = "fail"
							df.FailReason[i] = failReason
						}
						isPass++
						inIsPass++
						break
					}
				}
			}
		}

		if inIsPass > 0 {
			df.TestResult[i] = "fail"
			if len(df.TestResult) < i+1 {
				df.FailReason = append(df.FailReason, "测试失败")
			}
		} else {
			if len(df.TestResult) < i+1 {
				df.TestResult = append(df.TestResult, "pass")
				df.FailReason = append(df.FailReason, "")
			} else {
				df.TestResult[i] = "pass"
				if len(df.FailReason) < i+1 {
					df.FailReason = append(df.FailReason, "")
				} else {
					df.FailReason[i] = ""
				}

			}
		}
	}

	if isPass != 0 {
		result = "fail"
	} else {
		result = "pass"
		df.FailReason = []string{}
	}
	var dataInfo, dataWithHeader []byte
	var errTmp error

	outputDict, errTmp = df.ChangeOutputValue(outputDict)
	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%v,%v", err, errTmp)
		} else {
			err = errTmp
		}
	}

	df.Output = outputDict

	if strings.HasSuffix(filePath, ".json") {
		dataInfo, errTmp = json.MarshalIndent(df, "", "    ")
	} else {
		dataInfo, errTmp = yaml.Marshal(df)
	}

	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%v,%v", err, errTmp)
		} else {
			err = errTmp
		}
	}

	if strings.HasSuffix(filePath, ".json") {
		dataWithHeader, errTmp = json.MarshalIndent(df, "", "    ")
	} else {
		dataWithHeader, errTmp = yaml.Marshal(df)
	}

	errTmp = ioutil.WriteFile(dst, dataWithHeader, 0644)

	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%v,%v", err, errTmp)
		} else {
			err = errTmp
		}
	}

	switch source {
	case "data", "playbook", "ai_data", "ai_playbook":
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

	return
}

func GetDataByFileName(fileName, source string) (dbData SceneData, err error) {
	baseName := path.Base(fileName)
	if strings.HasPrefix(source, "ai") {
		models.Orm.Table("ai_data").Where("file_name = ?", baseName).Find(&dbData)
	} else if strings.HasPrefix(source, "history") {
		var dbHData HistoryDataDetail
		Logger.Debug("baseName: %v", baseName)
		models.Orm.Table("scene_data_test_history").Where("content = ?", baseName).Find(&dbHData)
		dbData.Name = dbHData.Name
		dbData.FileName = dbHData.Content
		dbData.App = dbHData.App
		dbData.RunTime = 1
		dbData.FileType = dbHData.FileType
		dbData.ApiId = dbHData.ApiId
	} else {
		models.Orm.Table("scene_data").Where("file_name = ?", baseName).Find(&dbData)
	}

	if len(dbData.FileName) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", baseName)
		Logger.Error("%s", err)
		return
	}

	return
}

func GetDataFileRawContent(fileName string) (content string, err error) {
	var sceneData SceneData
	baseName := path.Base(fileName)
	models.Orm.Table("scene_data").Where("file_name = ?", baseName).Find(&sceneData)
	if len(sceneData.FileName) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", baseName)
		Logger.Error("%s", err)
		return
	}

	// 历史数据修正
	if strings.Contains(sceneData.Content, "<pre><code>") {
		sceneData.Content = strings.Replace(sceneData.Content, "<pre><code>", "", -1)
		sceneData.Content = strings.Replace(sceneData.Content, "</code></pre>", "", -1)
	}

	content = sceneData.Content
	return
}

func GetBodyFromRawContent(lang, fileName, content string, depOutVars map[string][]interface{}) (bodys []map[string]interface{}, err error) {
	contentStr, err := GetAfterContent(lang, string(content), depOutVars)
	if err != nil {
		Logger.Debug("rawContent: %s", content)
		Logger.Debug("afterContent: %s", contentStr)
		Logger.Error("%v", err)
		return
	}

	var df DataFile
	if strings.HasSuffix(fileName, ".json") {
		err = json.Unmarshal([]byte(contentStr), &df)
	} else {
		err = yaml.Unmarshal([]byte(contentStr), &df)
	}
	if err != nil {
		Logger.Debug("afterContent%s", contentStr)
		Logger.Error("%s", err)
		return
	}

	bodys, _ = df.GetBody()

	return
}
