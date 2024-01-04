package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

func RepeatRunPlaybook(id, mode, product, source string) (err error) {
	var productInfo DbProduct
	playbookInfo, productSceneInfo, err := GetPlRunInfo(source, id)

	if len(product) > 0 {
		productInfo, _ = GetProductInfo(product)
	} else {
		productInfo = productSceneInfo
	}

	if productInfo.Threading == "yes" && playbookInfo.RunTime > 1 && productInfo.ThreadNumber > 1 {
		if playbookInfo.RunTime > productInfo.ThreadNumber {
			loopNum := playbookInfo.RunTime/productInfo.ThreadNumber + 1
			count := 1
			for i := 0; i < loopNum; i++ {
				Logger.Info("并发模式-最大执行数:%d,总循环次数:%d,当前循环第%d次", productInfo.ThreadNumber, loopNum, i+1)
				wg := sync.WaitGroup{}
				for j := 0; j < productInfo.ThreadNumber; j++ {
					if count > playbookInfo.RunTime {
						break
					}
					Logger.Info("并发模式-执行次数:%d", count)
					wg.Add(1)
					go func(playbookInfo DbScene, productInfo DbProduct) {
						err1 := playbookInfo.RunPlaybook(mode, source, productInfo)
						if err1 != nil {
							err = err1
							Logger.Error("%s", err)
							return
						}
						wg.Done()
					}(playbookInfo, productInfo)
					count++
				}
				wg.Wait()
			}
		} else {
			wg := sync.WaitGroup{}
			for i := 0; i < playbookInfo.RunTime; i++ {
				Logger.Info("并发模式-执行次数:%d", i+1)
				wg.Add(1)
				go func(playbookInfo DbScene) {
					err1 := playbookInfo.RunPlaybook(mode, source, productInfo)
					if err1 != nil {
						err = err1
						Logger.Error("%s", err)
						return
					}
					wg.Done()
				}(playbookInfo)
			}
			wg.Wait()
		}
	} else {
		for i := 0; i < playbookInfo.RunTime; i++ {
			Logger.Info("串行模式-执行次数:%d", i+1)
			err1 := playbookInfo.RunPlaybook(mode, source, productInfo)
			if err1 != nil {
				err = err1
				return
			}
		}
	}
	return
}

func CompareResult(apis []string, mode string) (result string, err error) {
	baseChk := make(map[string][]interface{})
	isPass := 0

	for _, item := range apis {
		var sceneFileBase DataFile
		content, err1 := ioutil.ReadFile(item)
		if err1 != nil {
			err = err1
			Logger.Error("Error: %s, filePath: %s", err, item)
			return
		}
		err = yaml.Unmarshal([]byte(content), &sceneFileBase)
		for k, v := range sceneFileBase.Output {
			baseChk[k] = v
		}
	}
	for _, filePath := range apis {
		var sceneFile DataFile
		content, err1 := ioutil.ReadFile(filePath)
		if err1 != nil {
			err = err1
			Logger.Error("Error: %s, filePath: %s", err, filePath)
			return
		}
		var err2 error
		if mode == "json" {
			err2 = json.Unmarshal([]byte(content), &sceneFile)
		} else if mode == "yaml" {
			err2 = yaml.Unmarshal([]byte(content), &sceneFile)
		} else {
			err2 = yaml.Unmarshal([]byte(content), &sceneFile)
		}
		if err2 != nil {
			err = err2
			Logger.Error("Err: %s, filePath: %s", err, filePath)
			return
		}

		for baseK, baseV := range sceneFile.Output {
			if len(baseV) > 1 {
				for i := 0; i < len(baseV); i++ {
					if len(baseChk[baseK]) > i {
						str1 := fmt.Sprintf("%v", baseV[i])
						str2 := fmt.Sprintf("%v", baseChk[baseK][i])
						if str1 == str2 {
							Logger.Info(" 数据对比%s: [%v]=[%v] 结果:pass", baseK, baseV[i], baseChk[baseK][i])
						} else {
							isPass++
							errTmp := fmt.Errorf("数据对比%s: [%v]=[%v] 结果:fail", baseK, baseV[i], baseChk[baseK][i])
							Logger.Error("%s", errTmp)
							if err != nil {
								err = fmt.Errorf("%s; %s", err, errTmp)
							} else {
								err = errTmp
							}
						}
					}
				}
			} else {
				str1 := fmt.Sprintf("%v", baseV[0])
				str2 := fmt.Sprintf("%v", baseChk[baseK][0])
				if str1 == str2 {
					Logger.Info(" 数据对比%s: [%v]=[%v] 结果:pass", baseK, baseV[0], baseChk[baseK][0])
				} else {
					isPass++
					errTmp := fmt.Errorf("数据对比%s: [%v]=[%v] 结果:fail", baseK, baseV[0], baseChk[baseK][0])
					Logger.Error("%s", errTmp)
					if err != nil {
						err = fmt.Errorf("%s; %s", err, errTmp)
					} else {
						err = errTmp
					}
				}
			}
		}
	}

	if isPass == 0 {
		result = "pass"
	} else {
		result = "fail"
	}

	return
}

func (p Playbook) RunPlaybookContent(product string, privateParameter map[string]interface{}, source string) (result, historyApi string, err error) {
	filePath := p.Apis[p.Tag]

	var sceneFile DataFile

	errTmp := CreateFileByFileName(filePath)
	if errTmp != nil {
		Logger.Error("%s", errTmp)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Logger.Debug("filePath: %v", filePath)
		Logger.Error("Error: %s, filePath: %s", err, filePath)
		result = "fail"
		historyApi = path.Base(filePath)
		return
	}

	suffix := GetStrSuffix(filePath)
	if suffix == ".json" {
		err = json.Unmarshal([]byte(content), &sceneFile)
	} else if suffix == ".yaml" || suffix == ".yml" {
		err = yaml.Unmarshal([]byte(content), &sceneFile)
	}

	if err != nil {
		Logger.Error("Error: %s, filePath: %s", err, filePath)
		historyApi = path.Base(filePath)
		return
	}

	if len(product) == 0 {
		product = p.Product
	}

	envConfig, err := GetEnvConfig(product, "scene")
	envType, _ := GetEnvTypeByName(product)

	dataEnvConfig, err := GetEnvConfig(sceneFile.Api.App, "data")
	if len(dataEnvConfig.Prepath) > 0 {
		envConfig.Prepath = dataEnvConfig.Prepath
	} else if len(sceneFile.Env.Prepath) > 0 {
		envConfig.Prepath = sceneFile.Env.Prepath
	}

	header, err := sceneFile.GetHeader(envConfig)
	if err != nil {
		historyApi = path.Base(filePath)
		return
	}

	lang := GetRequestLangage(header)

	urls, err := sceneFile.GetUrl(envConfig)
	if len(urls) == 0 {
		Logger.Error("url为空，请检查")
		historyApi = path.Base(filePath)
		return
	}
	sceneFile.Urls = urls

	if err != nil {
		Logger.Debug("filePath: %v", filePath)
		Logger.Debug("urlQuerys: %+v", urls)
		Logger.Error("%s", err)
		historyApi = path.Base(filePath)
		return
	}

	outputDict, err := p.GetPalybookDepParams(privateParameter)
	//Logger.Debug("outputDict: %v", outputDict)
	querys, err := sceneFile.GetQuery(lang, outputDict)
	if err != nil {
		Logger.Debug("filePath: %v", filePath)
		Logger.Debug("querys: %+v", querys)
		historyApi = path.Base(filePath)
		return
	}

	bodys, bodyList, err := sceneFile.GetBody(lang, outputDict)
	_ = sceneFile.CreateDataOrderByKey(lang, outputDict) // 执行动作：create_xxx
	_ = sceneFile.RecordDataOrderByKey(bodys)            // 执行动作：record_xxx
	_ = sceneFile.ModifyFileWithData(bodys)              // 执行动作：modify_file

	if err != nil {
		Logger.Debug("outputDict: %v", outputDict)
		Logger.Debug("filePath: %v", filePath)
		historyApi = path.Base(filePath)
		return
	}

	var resList [][]byte
	var errs []error
	tag := 0

	if sceneFile.GetIsParallel() {
		wg := sync.WaitGroup{}
		for _, url := range urls {
			if len(querys) > 0 {
				wg.Add(len(querys))
				for _, data := range querys {
					dJson, _ := json.Marshal(data)
					if tag == 0 {
						sceneFile.Request = []string{string(dJson)}
					} else {
						sceneFile.Request = append(sceneFile.Request, string(dJson))
					}
					tag++
					go func(method, url string, data map[string]interface{}, header map[string]interface{}) {
						defer wg.Add(-1)
						res, err := RunHttp(method, url, data, header)
						resList = append(resList, res)
						errs = append(errs, err)
					}(sceneFile.Api.Method, url, data, header)
				}
			} else if len(bodys) > 0 {
				wg.Add(len(bodys)) // 一次把全部需要等待的任务加上
				for _, data := range bodys {
					dJson, errTmp := json.Marshal(data)
					if errTmp != nil {
						var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
						dJsonNew, subErr := jsonNew.Marshal(&data)
						if subErr != nil {
							Logger.Error("%v", subErr)
						}
						dJson = dJsonNew
					}

					if tag == 0 {
						sceneFile.Request = []string{string(dJson)}
					} else {
						sceneFile.Request = append(sceneFile.Request, string(dJson))
					}
					tag++
					//wg.Add(1)  // 定时任务执行过程中，会概率性发生panic
					go func(method, url string, data, header map[string]interface{}) (errs []error) {
						defer wg.Add(-1)
						res, err := RunHttp(method, url, data, header)
						resList = append(resList, res)
						errs = append(errs, err)
						return
					}(sceneFile.Api.Method, url, data, header)
				}
			} else {
				wg.Add(1)
				sceneFile.Request = []string{}
				go func(method, url string, header map[string]interface{}) (errs []error) {
					defer wg.Add(-1)
					res, err := RunHttp(method, url, nil, header)
					resList = append(resList, res)
					errs = append(errs, err)
					return
				}(sceneFile.Api.Method, url, header)
			}
			wg.Wait()
		}
	} else {
		for _, url := range urls {
			if len(querys) > 0 {
				for _, data := range querys {
					dJson, _ := json.Marshal(data)
					if tag == 0 {
						sceneFile.Request = []string{string(dJson)}
					} else {
						sceneFile.Request = append(sceneFile.Request, string(dJson))
					}
					tag++

					if sceneFile.Api.Method == "delete" {
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
						res, err := RunHttp(sceneFile.Api.Method, url, nil, header)
						resList = append(resList, res)
						errs = append(errs, err)
					} else {
						res, err := RunHttp(sceneFile.Api.Method, url, data, header)
						resList = append(resList, res)
						errs = append(errs, err)
					}

					_ = sceneFile.SetSleepAction()
				}
			} else if len(bodys) > 0 || len(bodyList) > 0 {
				if len(bodyList) > 0 {
					var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
					readerNew, subErr := jsonNew.Marshal(&bodyList)
					if subErr != nil {
						Logger.Error("%v", subErr)
					}
					sceneFile.Request = []string{string(readerNew)}
					res, err := RunHttpJsonList(sceneFile.Api.Method, url, bodyList, header)
					if err != nil {
						Logger.Debug("%s", err)
					}
					resList = append(resList, res)
					errs = append(errs, err)

					_ = sceneFile.SetSleepAction()
				} else {
					for _, data := range bodys {
						dJson, errTmp := json.Marshal(data)
						if errTmp != nil {
							var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
							dJsonNew, _ := jsonNew.Marshal(&data)
							dJson = dJsonNew
						}
						if tag == 0 {
							sceneFile.Request = []string{string(dJson)}
						} else {
							sceneFile.Request = append(sceneFile.Request, string(dJson))
						}
						tag++
						res, err := RunHttp(sceneFile.Api.Method, url, data, header)
						resList = append(resList, res)
						errs = append(errs, err)
						_ = sceneFile.SetSleepAction()
					}
				}
			} else {
				sceneFile.Request = []string{}
				res, err := RunHttp(sceneFile.Api.Method, url, nil, header)
				resList = append(resList, res)
				errs = append(errs, err)
				_ = sceneFile.SetSleepAction()
			}
		}
	}

	result, dst, err1 := sceneFile.GetResult(lang, source, filePath, header, p.IsThread, resList, outputDict, errs)
	b, _ := IsStrEndWithTimeFormat(filePath)
	if b {
		dst = filePath
	}

	historyApi = strings.Split(dst, "/history/")[1]
	err = WriteDataResultByFile(filePath, result, dst, product, envType, err1)

	if err1 != nil {
		if err != nil {
			err = fmt.Errorf("%s, %s", err, err1)
		} else {
			err = err1
		}
	}

	return
}

func (p Playbook) GetHistoryApiList(dbScene DbScene) (apiStr string) {
	var lastFileName string
	baseName := path.Base(dbScene.LastFile)
	b, _ := IsStrEndWithTimeFormat(baseName)
	if b {
		dirName := GetHistoryDataDirName(baseName)
		suffix := GetStrSuffix(baseName)
		lastFileName = fmt.Sprintf("%s%s", dirName, suffix)
	} else {
		lastFileName = baseName
	}

	rawApiList := GetListFromHtml(dbScene.ApiList)
	for index, item := range p.HistoryApis {
		// 如果执行过的数据中有空值，则跳过
		if len(item) == 0 {
			continue
		}
		b, _ := IsStrEndWithTimeFormat(item)
		if index == 0 {
			if b {
				apiStr = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s\">%s</a>", item, path.Base(item))
			} else {
				apiStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", item, path.Base(item))
			}
		} else {
			if b {
				apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/history/preview?path=/%s\">%s</a>", apiStr, item, path.Base(item))
			} else {
				apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", apiStr, item, path.Base(item))
			}
		}
	}

	lastFileTag := len(rawApiList)
	for index, item := range rawApiList {
		if item == lastFileName {
			lastFileTag = index
		}
		if index > lastFileTag {
			if len(apiStr) > 0 {
				apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", apiStr, item, path.Base(item))
			} else {
				apiStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", item, path.Base(item))
			}
		}
	}

	return
}

func (p Playbook) WritePlaybookResult(id, result, lastFile, product, source string, envType int, errIn error) (err error) {
	var dbScene DbScene
	var sceneRecode SceneRecord
	s, _ := strconv.Atoi(id)
	models.Orm.Table("playbook").Where("id = ?", s).Find(&dbScene)
	if len(dbScene.Name) == 0 {
		return
	}

	apiStr := p.GetHistoryApiList(dbScene)

	sceneRecode.ApiList = apiStr

	if source == "playbook" {
		curTime := time.Now()
		dbScene.UpdatedAt = curTime.Format(baseFormat)
		if len(result) == 0 {
			dbScene.Result = " "
		} else {
			dbScene.Result = result
		}
		if len(lastFile) == 0 {
			dbScene.LastFile = " " // 用空格字符串刷新数据
		} else {
			b, _ := IsStrEndWithTimeFormat(path.Base(lastFile))
			if b {
				dirName := GetHistoryDataDirName(path.Base(lastFile))
				suffix := GetStrSuffix(path.Base(lastFile))
				dbScene.LastFile = fmt.Sprintf("%s%s", dirName, suffix)
			} else {
				dbScene.LastFile = path.Base(lastFile)
			}

		}
		if errIn != nil {
			dbScene.FailReason = fmt.Sprintf("%s", errIn)
		} else {
			dbScene.FailReason = " "
		}

		err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	sceneRecode.Name = dbScene.Name
	b, _ := IsStrEndWithTimeFormat(path.Base(lastFile))
	if b {
		dirName := GetHistoryDataDirName(path.Base(lastFile))
		sceneRecode.LastFile = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, path.Base(lastFile), path.Base(lastFile))
	} else {
		sceneRecode.LastFile = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", path.Base(lastFile), path.Base(lastFile))
	}
	sceneRecode.SceneType = dbScene.SceneType
	sceneRecode.Result = result
	if errIn != nil {
		sceneRecode.FailReason = fmt.Sprintf("%s", errIn)
	}

	if len(product) > 0 {
		sceneRecode.Product = product
	} else {
		sceneRecode.Product = product
	}

	sceneRecode.EnvType = envType

	err = WritePlaybookRecord(sceneRecode)

	return
}

// id为历史记录中的ID, mode为继续还是再来一次，batchTag为执行批次，envType为执行环境的类型
func WritePlaybookHistoryResult(id, result, lastFile, mode string, envType int, errIn error) (err error) {
	var dbScene DbSceneRecord
	s, _ := strconv.Atoi(id)
	models.Orm.Table("scene_test_history").Where("id = ?", s).Find(&dbScene)
	if len(dbScene.Name) == 0 {
		return
	}
	if len(result) == 0 {
		dbScene.Result = " "
	} else {
		dbScene.Result = result
	}

	baseLastFile := path.Base(lastFile)
	b, _ := IsStrEndWithTimeFormat(baseLastFile)
	var dirName string
	if b {
		dirName = GetHistoryDataDirName(baseLastFile)
	}
	if len(lastFile) == 0 {
		dbScene.LastFile = " " // 用空格字符串刷新数据
	} else {
		if b {
			dbScene.LastFile = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, baseLastFile, baseLastFile)
		} else {
			dbScene.LastFile = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", baseLastFile, baseLastFile)
		}
	}

	if errIn != nil {
		dbScene.FailReason = fmt.Sprintf("%s", errIn)
	} else {
		dbScene.FailReason = " "
	}

	if mode == "continue" {
		curTime := time.Now()
		dbScene.UpdatedAt = curTime.Format(baseFormat)
		err = UpdateSceneRecord(dbScene)
	} else {
		var sceneRecode SceneRecord
		sceneRecode.Name = dbScene.Name
		sceneRecode.ApiList = dbScene.ApiList
		if b {
			sceneRecode.LastFile = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, baseLastFile, baseLastFile)
		} else {
			sceneRecode.LastFile = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", baseLastFile, baseLastFile)
		}
		sceneRecode.SceneType = dbScene.SceneType
		sceneRecode.Result = dbScene.Result
		sceneRecode.FailReason = dbScene.FailReason
		sceneRecode.Product = dbScene.Product
		sceneRecode.EnvType = envType
		err = WritePlaybookRecord(sceneRecode)
	}

	return
}

func WritePlaybookRecord(sceneRecode SceneRecord) (err error) {
	err = models.Orm.Table("scene_test_history").Create(sceneRecode).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func UpdateSceneRecord(sceneRecode DbSceneRecord) (err error) {
	err = models.Orm.Table("scene_test_history").Where("id = ?", sceneRecode.Id).Update(sceneRecode).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func GetHistoryPlaybook(id string) (playbook Playbook, err error) {
	var dbScene DbScene
	var filePaths []string
	s, _ := strconv.Atoi(id)
	models.Orm.Table("scene_test_history").Where("id = ?", s).Find(&dbScene)
	if len(dbScene.ApiList) == 0 {
		err = fmt.Errorf("未找到[%v]场景，请核对", s)
		Logger.Error("%s", err)
		return
	}

	var filePath string
	fileNames := GetListFromHtml(dbScene.ApiList)

	for _, item := range fileNames {

		var dirName string

		b, num := IsStrEndWithTimeFormat(item)
		suffix := GetStrSuffix(item)
		if b {
			dirName = item[:len(item)-num-len(suffix)]
			filePath = fmt.Sprintf("%s/%s/%s", HistoryBasePath, dirName, item)
		} else {
			dirName = item[:len(item)-len(suffix)]
			filePath = fmt.Sprintf("%s/%s", DataBasePath, item)

		}
		filePaths = append(filePaths, filePath)
	}
	playbook.Apis = filePaths
	playbook.Name = dbScene.Name
	playbook.LastFile = dbScene.LastFile
	playbook.Product = dbScene.Product
	return
}

func GetPlRunInfo(source, id string) (dbScene DbScene, dbProduct DbProduct, err error) {
	s, _ := strconv.Atoi(id)
	if source == "history" {
		models.Orm.Table("scene_test_history").Where("id = ?", s).Find(&dbScene)
	} else {
		models.Orm.Table("playbook").Where("id = ?", s).Find(&dbScene)
	}

	if len(dbScene.ApiList) == 0 {
		err = fmt.Errorf("未找到对应场景，请核对: %s", id)
		Logger.Error("%s", err)
		return
	}

	dbProduct, err = GetProductInfo(dbScene.Product)
	if err != nil {
		Logger.Error("%v", err)
	}

	return
}

func GetProductInfo(product string) (dbProduct DbProduct, err error) {
	models.Orm.Table("product").Where("product = ?", product).Find(&dbProduct)
	if len(dbProduct.Name) == 0 {
		err = fmt.Errorf("未找到: %v 环境信息，请核对", product)
	}
	return
}

func (p Playbook) GetPalybookDepParams(privateParameter map[string]interface{}) (outputDict map[string][]interface{}, err error) {
	var preApis []string
	for _, item := range p.HistoryApis {
		fullPathApi := fmt.Sprintf("%s/%s", HistoryBasePath, item)
		preApis = append(preApis, fullPathApi)
	}

	outputDict = make(map[string][]interface{})

	for _, filePath := range preApis {
		var sceneFile DataFile
		content, err1 := ioutil.ReadFile(filePath)
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}
		if strings.HasSuffix(filePath, ".json") {
			err1 = json.Unmarshal([]byte(content), &sceneFile)
		} else if strings.HasSuffix(filePath, ".yml") {
			err1 = yaml.Unmarshal([]byte(content), &sceneFile)
		}

		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}

		if len(sceneFile.Output) > 0 {
			for k, v := range sceneFile.Output {
				outputDict[k] = v
			}
		}
		if len(sceneFile.Request) > 0 && len(sceneFile.Request[0]) > 0 {
			requestMap := make(map[string]interface{})
			errTmp := json.Unmarshal([]byte(sceneFile.Request[0]), &requestMap)
			if errTmp != nil {
				Logger.Debug("%v", sceneFile.Request[0])
				Logger.Error("%v", errTmp)
			} else {
				for k, v := range requestMap {
					if _, ok := outputDict[k]; !ok {
						outputDict[k] = append(outputDict[k], v)
					}

				}
			}

		}
	}

	if len(p.Apis) == 1 {
		selfApi := p.Apis[p.Tag]
		var selfScene DataFile
		var otherApis []string
		content, err1 := ioutil.ReadFile(selfApi)
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}
		if strings.HasSuffix(selfApi, ".json") {
			err1 = json.Unmarshal([]byte(content), &selfScene)
		} else if strings.HasSuffix(selfApi, ".yml") {
			err1 = yaml.Unmarshal([]byte(content), &selfScene)
		}
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}
		if len(selfScene.Api.ParamApis) > 0 {
			otherApis = append(otherApis, selfScene.Api.ParamApis...)
		}
		if len(selfScene.Api.PreApi) > 0 {
			otherApis = append(otherApis, selfScene.Api.PreApi...)
		}

		for _, fileName := range otherApis {
			var sceneFile DataFile
			filePath := fmt.Sprintf("%s/%s", DataBasePath, fileName)
			content, err1 := ioutil.ReadFile(filePath)
			if err1 != nil {
				err = err1
				Logger.Error("%s", err)
			}
			if strings.HasSuffix(filePath, ".json") {
				err1 = json.Unmarshal([]byte(content), &sceneFile)
			} else if strings.HasSuffix(filePath, ".yml") {
				err1 = yaml.Unmarshal([]byte(content), &sceneFile)
			}
			if err1 != nil {
				err = err1
				Logger.Error("%s", err)
			}

			if len(sceneFile.Output) > 0 {
				for k, v := range sceneFile.Output {
					outputDict[k] = v
				}
			}

			if len(sceneFile.Request) > 0 {
				requestMap := make(map[string]string)
				errTmp := json.Unmarshal([]byte(sceneFile.Request[0]), &requestMap)
				if errTmp != nil {
					Logger.Error("%v", errTmp)
				} else {
					for k, v := range requestMap {
						if _, ok := outputDict[k]; !ok {
							outputDict[k] = append(outputDict[k], v)
						}
					}
				}

			}
		}
	}

	for k, v := range privateParameter {
		outputDict[k] = append(outputDict[k], v)
	}

	return
}

func GetPriority(ids string) (idList []string, err error) {
	idsTmp := strings.Split(ids, ",")

	models.Orm.Table("playbook").Where("id in (?)", idsTmp).Order("priority").Pluck("id", &idList)
	if len(idList) == 0 {
		err = fmt.Errorf("未找到对应数据，请核对: %s", ids)
		Logger.Error("%s", err)
		return
	}
	return
}

func CopyPlaybook(id, userName string) (err error) {
	var dbScene DbScene
	models.Orm.Table("playbook").Where("id = ?", id).Find(&dbScene)
	if len(dbScene.ApiList) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", id)
		Logger.Error("%s", err)
		return
	}

	var scene SceneWithNoUpdateTime
	scene.Product = dbScene.Product
	scene.Name = fmt.Sprintf("%s_复制", dbScene.Name)
	scene.ApiList = dbScene.ApiList
	scene.RunTime = dbScene.RunTime
	scene.SceneType = dbScene.SceneType
	scene.Remark = dbScene.Remark
	scene.Priority = dbScene.Priority
	scene.UserName = userName
	scene.DataNumber = dbScene.DataNumber
	//curTime := time.Now()
	//dbScene.UpdatedAt = curTime.Format(baseFormat)

	err = models.Orm.Table("playbook").Create(scene).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func ImportPlaybook(id string) (newCount, oldCount int, err error) {
	product, _ := GetProductName(id)
	fileName := fmt.Sprintf("%s/%s.xlsx", CommonFilePath, product)
	sceneList, err := ReadSceneFromExcel(fileName)
	if err != nil {
		return
	}
	for _, item := range sceneList {
		var dbScene DbSceneWithNoUpdateTime
		models.Orm.Table("playbook").Where("name = ? and product = ?", item.Name, item.Product).Find(&dbScene)
		if len(dbScene.Name) == 0 {
			err = models.Orm.Table("playbook").Create(item).Error
			if err != nil {
				Logger.Error("%s", err)
				return
			}
			newCount++
		} else {
			oldCount++
			Logger.Info("产品:[%v]下场景:[%s]已存在，如有需要，请手动更新", item.Product, item.Name)
		}
	}

	return
}

func GetPlaybookByName(name, product string) (sceneInfo SceneInfoModel, err error) {
	var dbScene DbScene
	if len(product) == 0 {
		models.Orm.Table("playbook").Where("name = ?", name).Find(&dbScene)
	} else {
		models.Orm.Table("playbook").Where("name = ? and product = ?", name, product).Find(&dbScene)
	}
	if len(dbScene.Name) == 0 {
		models.Orm.Table("playbook").Where("name = ?", name).Find(&dbScene)
		if len(dbScene.Name) == 0 {
			err = fmt.Errorf("未找到[%v]场景，请核对", name)
			Logger.Warning("%s", err)
			return
		}
	}
	sceneInfo.Name = name
	sceneInfo.Product = dbScene.Product
	sceneInfo.RunNum = dbScene.RunTime
	if dbScene.SceneType == 2 {
		sceneInfo.SceneType = "比较"
	} else {
		sceneInfo.SceneType = "默认"
	}

	dataList := GetListFromHtml(dbScene.ApiList)
	for _, item := range dataList {
		var dataModel DepDataModel
		item = strings.Replace(item, "\n", "", -1)
		dataModel.DataFile = item
		sceneInfo.DataList = append(sceneInfo.DataList, dataModel)
	}

	return
}

func GetAllPlaybook() (names []string, err error) {
	models.Orm.Table("playbook").Order("created_at desc").Group("name").Pluck("name", &names)
	if len(names) == 0 {
		Logger.Warning("暂无场景数据")
		return
	}
	return
}

func SaveScene(sceneSave SceneSaveModel) (err error) {
	var dbScene DbScene
	var scene SceneWithNoUpdateTime
	var apiStr, numStr string
	scene.Product = sceneSave.Product
	scene.Name = sceneSave.Name
	for index, value := range sceneSave.DataList {
		if index == 0 {
			apiStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", value.DataFile, value.DataFile)
			numStr = fmt.Sprintf("%v", index+1)
		} else {
			apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", apiStr, value.DataFile, value.DataFile)
			numStr = fmt.Sprintf("%s,%v", numStr, index+1)
		}
	}

	scene.ApiList = apiStr
	scene.DataNumber = numStr
	scene.RunTime = sceneSave.RunNum
	scene.SceneType = sceneSave.SceneType
	scene.Remark = ""
	scene.Priority = 999
	scene.UserName = ""

	models.Orm.Table("playbook").Where("name = ? and product = ?", sceneSave.Name, sceneSave.Product).Find(&dbScene)
	if len(dbScene.Name) == 0 {
		err = models.Orm.Table("playbook").Create(scene).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	} else {
		dbScene.ApiList = apiStr
		dbScene.DataNumber = numStr
		dbScene.SceneType = sceneSave.SceneType
		dbScene.RunTime = sceneSave.RunNum
		err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	return
}

func RunPlaybookDebugData(sceneModel SceneSaveModel) (runResp RunSceneRespModel, err error) {
	var playbook Playbook
	playbook.Product = sceneModel.Product
	playbook.Name = sceneModel.Name
	var apisStr string
	for _, item := range sceneModel.DataList {
		var dirName, filePath string
		b, _ := IsStrEndWithTimeFormat(item.DataFile)
		if b {
			dirName = GetHistoryDataDirName(item.DataFile)
			filePath = fmt.Sprintf("%s/%s/%s", HistoryBasePath, dirName, item.DataFile)
		} else {
			filePath = fmt.Sprintf("%s/%s", DataBasePath, item.DataFile)
		}

		playbook.Apis = append(playbook.Apis, filePath)
		if len(apisStr) == 0 {
			apisStr = item.DataFile
		} else {
			apisStr = fmt.Sprintf("%s %s", apisStr, item.DataFile)
		}
	}

	var dbProduct DbProduct
	models.Orm.Table("product").Where("product = ?", playbook.Product).Find(&dbProduct)
	privateParameter := make(map[string]interface{})
	if len(dbProduct.Name) > 0 {
		if len(dbProduct.PrivateParameter) > 2 {
			err = json.Unmarshal([]byte(dbProduct.PrivateParameter), &privateParameter)
			if err != nil {
				Logger.Error("%s", err)
				return
			}
		}
	}

	runResp.TestResult = "pass"
	for k := range playbook.Apis {
		playbook.Tag = k
		result, historyApi, errTmp := playbook.RunPlaybookContent(playbook.Product, privateParameter, "console")
		playbook.HistoryApis = append(playbook.HistoryApis, historyApi)
		if errTmp != nil || result != "pass" {
			runResp.TestResult = result

			if errTmp != nil {
				err = errTmp
				runResp.FailReason = fmt.Sprintf("%s", errTmp)
			}

			runResp.LastFile = historyApi
			playbook.LastFile = historyApi
			goto Record
			return
		}
	}

	if sceneModel.SceneType == 2 {
		result, errTmp := CompareResult(playbook.Apis, "")
		runResp.TestResult = result
		if errTmp != nil {
			err = errTmp
		}
	}

	if runResp.TestResult != "pass" {
		if err != nil {
			runResp.FailReason = fmt.Sprintf("%s", err)
		}
		runResp.LastFile = path.Base(playbook.LastFile)
	}

Record:
	var sceneRecode SceneRecord
	sceneRecode.Name = sceneModel.Name

	var dbScene DbScene
	baseLastFile := path.Base(playbook.LastFile)
	dbScene.LastFile = baseLastFile
	b, _ := IsStrEndWithTimeFormat(baseLastFile)
	dirName := GetHistoryDataDirName(baseLastFile)
	var lastFileStr string
	if b {
		lastFileStr = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, baseLastFile, baseLastFile)
	} else {
		lastFileStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", baseLastFile, baseLastFile)
	}
	dbScene.ApiList = apisStr
	apiAfter := playbook.GetHistoryApiList(dbScene)
	sceneRecode.ApiList = apiAfter
	sceneRecode.LastFile = lastFileStr
	sceneRecode.SceneType = sceneModel.SceneType
	sceneRecode.Result = runResp.TestResult
	sceneRecode.FailReason = runResp.FailReason
	sceneRecode.Product = sceneModel.Product
	sceneRecode.EnvType = dbProduct.EnvType

	err = WritePlaybookRecord(sceneRecode)

	return
}

func ReadSceneFromExcel(fileName string) (sceneList []SceneWithNoUpdateTime, err error) {
	xlsFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	for index, sheet := range xlsFile.Sheets {
		if index > 0 {
			break
		}

		for subIndex, row := range sheet.Rows {
			if subIndex == 0 {
				continue
			}
			var values []string
			var scene SceneWithNoUpdateTime
			for _, cell := range row.Cells {
				values = append(values, cell.String())
			}
			scene.Name = values[0]
			scene.ApiList = values[1]
			if values[3] == "比较" {
				scene.SceneType = 2
			} else {
				scene.SceneType = 1
			}
			prioInt, err := strconv.Atoi(values[4])
			if err != nil {
				Logger.Error("%s", err)
			} else {
				scene.Priority = prioInt
			}

			scene.RunTime, _ = strconv.Atoi(values[5])
			scene.Remark = values[8]
			scene.Product = values[9]
			sceneList = append(sceneList, scene)
		}
	}
	return
}

func UpdatePlaybookApiList(id string, apiList, numList []string) (err error) {
	var dbScene DbScene
	models.Orm.Table("playbook").Where("id = ?", id).Find(&dbScene)
	var apiStr, numStr string
	if len(dbScene.Name) == 0 {
		return
	} else {
		for index, value := range apiList {
			if len(numList) > index {
				numValue := numList[index]
				if len(numValue) == 0 {
					numList[index] = fmt.Sprintf("%d", index+1)
				}
			} else {
				numList = append(numList, fmt.Sprintf("%d", index+1))
			}
			if index == 0 {
				apiStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", value, value)
				numStr = fmt.Sprintf("%v", numList[index])
			} else {
				apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", apiStr, value, value)
				numStr = fmt.Sprintf("%s,%v", numStr, numList[index])
			}
		}

		dbScene.ApiList = apiStr
		dbScene.DataNumber = numStr
		err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}
	return
}

// 场景数据升级函数
func ModifyPlaybookContent() (err error) {
	var ids []string
	models.Orm.Table("playbook").Order("id ASC").Where("id > 246").Pluck("id", &ids)
	for _, id := range ids {
		var dbScene DbScene
		models.Orm.Table("playbook").Where("id = ?", id).Find(&dbScene)
		doc, errTmp := goquery.NewDocumentFromReader(strings.NewReader(dbScene.ApiList))
		if errTmp != nil {
			Logger.Error("%v", errTmp)
			err = errTmp
		}
		var numList []string
		handle := doc.Text()
		afterTxt1 := strings.Replace(handle, ".yml", ".yml,", -1)
		afterTxt2 := strings.Replace(afterTxt1, ".json", ".json,", -1)
		afterTxt := strings.Replace(afterTxt2, ".yaml", ".yaml,", -1)
		apiTmpList := strings.Split(afterTxt, ",")
		var apiList []string
		for _, value := range apiTmpList {
			if len(value) > 0 {
				apiList = append(apiList, value)
			}
		}

		var apiStr, numStr string
		if len(dbScene.Name) == 0 {
			Logger.Debug("未找到id: %v场景数据，请核对", id)
			return
		} else {
			for index, value := range apiList {
				if len(value) == 0 {
					continue
				}
				if len(numList) > index {
					numValue := numList[index]
					if len(numValue) == 0 {
						numList[index] = fmt.Sprintf("%d", index+1)
					}
				} else {
					numList = append(numList, fmt.Sprintf("%d", index+1))
				}
				if index == 0 {
					apiStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", value, value)
					numStr = fmt.Sprintf("%v", numList[index])
				} else {
					apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", apiStr, value, value)
					numStr = fmt.Sprintf("%s,%v", numStr, numList[index])
				}
			}

			dbScene.ApiList = apiStr
			dbScene.DataNumber = numStr
			err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		}
	}

	return
}
