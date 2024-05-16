package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

//func RepeatRunPlaybook(productInfo DbProduct, playbook Playbook, exeNum int, mode, source string) (err error) {
//	if productInfo.Threading == "yes" && playbookInfo.RunTime > 1 && productInfo.ThreadNumber > 1 {
//		if playbookInfo.RunTime > productInfo.ThreadNumber {
//			loopNum := playbookInfo.RunTime/productInfo.ThreadNumber + 1
//			count := 1
//			for i := 0; i < loopNum; i++ {
//				Logger.Info("并发模式-最大执行数:%d,总循环次数:%d,当前循环第%d次", productInfo.ThreadNumber, loopNum, i+1)
//				wg := sync.WaitGroup{}
//				for j := 0; j < productInfo.ThreadNumber; j++ {
//					if count > playbookInfo.RunTime {
//						break
//					}
//					Logger.Info("并发模式-执行次数:%d", count)
//					wg.Add(1)
//					go func(playbookInfo DbScene, productInfo DbProduct) {
//						playbook := playbookInfo.GetPlaybook()
//						err1 := playbook.RunPlaybook(playbookInfo.Id, mode, source, productInfo)
//						if err1 != nil {
//							err = err1
//							Logger.Error("%s", err)
//							return
//						}
//						wg.Done()
//					}(playbookInfo, productInfo)
//					count++
//				}
//				wg.Wait()
//			}
//		} else {
//			wg := sync.WaitGroup{}
//			for i := 0; i < playbookInfo.RunTime; i++ {
//				Logger.Info("并发模式-执行次数:%d", i+1)
//				wg.Add(1)
//				go func(playbookInfo DbScene) {
//					playbook := playbookInfo.GetPlaybook()
//					err1 := playbook.RunPlaybook(playbookInfo.Id, mode, source, productInfo)
//					if err1 != nil {
//						err = err1
//						Logger.Error("%s", err)
//						return
//					}
//					wg.Done()
//				}(playbookInfo)
//			}
//			wg.Wait()
//		}
//	} else {
//		for i := 0; i < playbookInfo.RunTime; i++ {
//			if playbookInfo.RunTime > 1 {
//				Logger.Info("串行模式-执行次数:%d", i+1)
//			}
//			playbook := playbookInfo.GetPlaybook()
//			//err1 := playbookInfo.RunPlaybook(mode, source, productInfo)
//			err1 := playbook.RunPlaybook(playbookInfo.Id, mode, source, productInfo)
//			if err1 != nil {
//				if err != nil {
//					err = fmt.Errorf("%v; %v", err, err1)
//				} else {
//					err = err1
//				}
//				break // 循环进行多次测试，如果遇错即退出
//			}
//		}
//	}
//
//	return
//}
//

func RepeatRunPlaybook(productInfo DbProduct, playbook Playbook, runNum int, mode, source, dbId string) (result, lastFile string, err error) {
	if productInfo.Threading == "yes" && runNum > 1 && productInfo.ThreadNumber > 1 {
		if runNum > productInfo.ThreadNumber {
			loopNum := runNum/productInfo.ThreadNumber + 1
			count := 1
			for i := 0; i < loopNum; i++ {
				Logger.Info("并发模式-最大执行数:%d,总循环次数:%d,当前循环第%d次", productInfo.ThreadNumber, loopNum, i+1)
				wg := sync.WaitGroup{}
				for j := 0; j < productInfo.ThreadNumber; j++ {
					if count > runNum {
						break
					}
					Logger.Info("并发模式-执行次数:%d", count)
					wg.Add(1)
					go func(playbook Playbook, dbId string) {
						subResult, subLastFile, err1 := playbook.RunPlaybook(dbId, mode, source, productInfo)
						result = subResult
						lastFile = subLastFile
						if err1 != nil {
							err = err1
							Logger.Error("%s", err)
							return
						}
						wg.Done()
					}(playbook, dbId)
					count++
				}
				wg.Wait()
			}
		} else {
			wg := sync.WaitGroup{}
			for i := 0; i < runNum; i++ {
				Logger.Info("并发模式-执行次数:%d", i+1)
				wg.Add(1)
				go func(playbook Playbook, dbId string) {
					subResult, subLastFile, err1 := playbook.RunPlaybook(dbId, mode, source, productInfo)
					result = subResult
					lastFile = subLastFile
					if err1 != nil {
						err = err1
						Logger.Error("%s", err)
						return
					}
					wg.Done()
				}(playbook, dbId)
			}
			wg.Wait()
		}
	} else {
		for i := 0; i < runNum; i++ {
			if runNum > 1 {
				Logger.Info("串行模式-执行次数:%d", i+1)
			}

			subResult, subLastFile, err1 := playbook.RunPlaybook(dbId, mode, source, productInfo)
			result = subResult
			lastFile = subLastFile
			if err1 != nil {
				if err != nil {
					err = fmt.Errorf("%v; %v", err, err1)
				} else {
					err = err1
				}
				break // 循环进行多次测试，如果遇错即退出
			}
		}
	}

	return
}

func RunPlaybookFromMgmt(id, mode, product, source string) (err error) {
	var productList []DbProduct
	playbookInfo, productSceneInfo, err := GetPlRunInfo(source, id)
	playbook := playbookInfo.GetPlaybook()
	if err != nil {
		return
	}

	if len(product) > 0 {
		productList, _ = GetProductInfo(product)
	} else {
		productList = productSceneInfo
	}

	for _, productInfo := range productList {
		_, _, err = RepeatRunPlaybook(productInfo, playbook, playbookInfo.RunTime, mode, source, playbookInfo.Id)
		if err != nil {
			Logger.Error("%s", err)
		}
	}
	return
}

func RunPlaybookFromConsole(sceneModel SceneSaveModel) (runResp RunSceneRespModel, err error) {
	var productInfo DbProduct
	var playbook Playbook

	productList, _ := GetProductInfo(sceneModel.Product)
	productInfo = productList[0]

	playbook.Product = sceneModel.Product
	if len(sceneModel.Name) == 0 {
		tmpStr := GetRandomStr(4, "")
		playbook.Name = fmt.Sprintf("临时测试场景-%s", tmpStr)
	} else {
		playbook.Name = sceneModel.Name
	}

	playbook.SceneType = sceneModel.SceneType

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
	}

	runResp.TestResult = "pass"

	runResp.TestResult, runResp.LastFile, err = RepeatRunPlaybook(productInfo, playbook, sceneModel.RunNum, "start", "consolePlaybook", "")
	if runResp.TestResult != "pass" {
		if err != nil {
			runResp.FailReason = fmt.Sprintf("%s", err)
		}
	}

	return
}

func CompareResult(apis []string, mode string) (result string, err error) {
	baseChk := make(map[string][]interface{})
	isPass := 0
	for _, item := range apis {
		var sceneFileBase DataFile
		item := fmt.Sprintf("%s/%s", HistoryBasePath, item)
		content, err1 := ioutil.ReadFile(item)
		if err1 != nil {
			err = err1
			Logger.Error("Error: %s, filePath: %s", err, item)
			return
		}
		err = yaml.Unmarshal([]byte(content), &sceneFileBase)
		if err != nil {
			Logger.Error("Error: %s, filePath: %s", err, item)
		}
		for k, v := range sceneFileBase.Output {
			baseChk[k] = v
		}
	}

	for _, filePath := range apis {
		var sceneFile DataFile
		filePath = fmt.Sprintf("%s/%s", HistoryBasePath, filePath)
		content, err1 := ioutil.ReadFile(filePath)
		if err1 != nil {
			err = err1
			Logger.Error("Error: %s, filePath: %s", err, filePath)
			return
		}
		var err2 error
		switch mode {
		case "json":
			err2 = json.Unmarshal([]byte(content), &sceneFile)
		case "yaml", "yml":
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

func (p Playbook) RunPlaybookContent(envType int, source string) (result, historyApi string, err error) {
	filePath := p.Apis[p.Tag]
	depOutVars, err := p.GetPalybookDepParams()
	if err != nil {
		Logger.Error("%s", err)
	}

	result, dst, errTmp := RunDataFile("", filePath, p.Product, source, depOutVars)
	if errTmp != nil {
		result = "fail"
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}
	}

	b, _ := IsStrEndWithTimeFormat(filePath)
	if b {
		dst = filePath
	}

	hApiTmps := strings.Split(dst, "/history/")

	if len(hApiTmps) > 1 {
		historyApi = hApiTmps[1]
	} else {
		historyApi = path.Base(filePath) // 如果未写到history中，
		dst = filePath
	}

	err = WriteDataResultByFile(filePath, result, dst, p.Product, envType, errTmp)

	if errTmp != nil {
		if err != nil {
			err = fmt.Errorf("%s, %s", err, errTmp)
		} else {
			err = errTmp
		}
	}

	return
}

func (p Playbook) GetHistoryApiList() (apiStr string) {
	var lastFileName string
	baseName := path.Base(p.LastFile)
	b, _ := IsStrEndWithTimeFormat(baseName)
	if b {
		dirName := GetHistoryDataDirName(baseName)
		suffix := GetStrSuffix(baseName)
		lastFileName = fmt.Sprintf("%s%s", dirName, suffix)
	} else {
		lastFileName = baseName
	}

	//rawApiList := GetListFromHtml(p.Apis)
	rawApiList := p.Apis
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

func (p Playbook) WritePlaybookResult(id, result, source, lastFile string, envType int, errIn error) (err error) {
	var sceneRecode SceneRecord

	apiStr := p.GetHistoryApiList()

	sceneRecode.ApiList = apiStr

	if source == "playbook" {
		var dbScene DbScene
		s, _ := strconv.Atoi(id)
		models.Orm.Table("playbook").Where("id = ?", s).Find(&dbScene)
		if len(dbScene.Name) == 0 {
			return
		}

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
			Logger.Error("%v", err)
		}
	} else if source == "task" {
		//Logger.Debug("lastFile: %v", lastFile)
	}

	sceneRecode.Name = p.Name
	if len(lastFile) > 0 {
		b, _ := IsStrEndWithTimeFormat(path.Base(lastFile))
		if b {
			dirName := GetHistoryDataDirName(path.Base(lastFile))
			sceneRecode.LastFile = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, path.Base(lastFile), path.Base(lastFile))
		} else {
			sceneRecode.LastFile = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", path.Base(lastFile), path.Base(lastFile))
		}
	}

	sceneRecode.SceneType = p.SceneType
	sceneRecode.Result = result
	if errIn != nil {
		sceneRecode.FailReason = fmt.Sprintf("%v", errIn)
	}

	sceneRecode.Product = p.Product
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
	Logger.Debug("b: %v", b)
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
	// 场景类型若未设置，置为默认值，串行中断: 1
	if sceneRecode.SceneType == 0 {
		sceneRecode.SceneType = 1
	}

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

func GetPlRunInfo(source, id string) (dbScene DbScene, dbProduct []DbProduct, err error) {
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

	return
}

func GetProductInfo(product string) (productList []DbProduct, err error) {
	productTmp := strings.Split(product, ",")
	models.Orm.Table("product").Where("product in (?)", productTmp).Find(&productList)
	if len(productList) == 0 {
		err = fmt.Errorf("未找到: %v 环境信息，请核对", product)
	}
	return
}

func (p Playbook) GetPalybookDepParams() (outputDict map[string][]interface{}, err error) {
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
			Logger.Debug("filePath: %s", filePath)
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

	var tmpStr []string
	models.Orm.Table("product").Where("product = ?", p.Product).Pluck("auth", &tmpStr)
	privateParameter := make(map[string]interface{})
	if len(tmpStr) > 0 {
		if len(tmpStr[0]) > 2 {
			err = json.Unmarshal([]byte(tmpStr[0]), &privateParameter)
			if err != nil {
				Logger.Error("%s", err)
				return
			}
		}
	}

	for k, v := range privateParameter {
		vStr := Interface2Str(v)
		var values []string
		if strings.Contains(vStr, ",") {
			values = strings.Split(vStr, ",")
			for _, subV := range values {
				strings.TrimSpace(subV)
				outputDict[k] = append(outputDict[k], subV)
			}
		} else {
			outputDict[k] = append(outputDict[k], v)
		}
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

	switch dbScene.SceneType {
	case 1:
		sceneInfo.SceneType = "串行中断"
	case 2:
		sceneInfo.SceneType = "串行比较"
	case 3:
		sceneInfo.SceneType = "串行继续"
	case 4:
		sceneInfo.SceneType = "普通并发"
	case 5:
		sceneInfo.SceneType = "并发比较"
	default:
		sceneInfo.SceneType = "串行中断"
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
			if values[3] == "串行中断" {
				scene.SceneType = 1
			} else if values[3] == "串行比较" {
				scene.SceneType = 2
			} else if values[3] == "串行继续" {
				scene.SceneType = 3
			} else if values[3] == "普通并发" {
				scene.SceneType = 4
			} else if values[3] == "并发比较" {
				scene.SceneType = 5
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

// 同接口一键生成场景用例
func CreatePlaybookByAPIId(id, userName string) (err error) {
	var apiDef ApiDefinition
	models.Orm.Table("api_definition").Where("id = ?", id).Find(&apiDef)
	if len(apiDef.ApiId) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", id)
		Logger.Error("%s", err)
		return
	}

	var dfNames []string
	models.Orm.Table("scene_data").Where("api_id = ?", apiDef.ApiId).Group("file_name").Pluck("file_name", &dfNames)
	if len(dfNames) == 0 {
		err = fmt.Errorf("未找到[%v]接口关联的数据，请核对", apiDef.ApiId)
		Logger.Error("%s", err)
		return
	}

	var dbScene DbScene
	var scene SceneWithNoUpdateTime
	var apiStr, numStr string
	scene.Product = ""
	scene.Name = fmt.Sprintf("全量-%s-%s-集合-%s", apiDef.ApiModule, apiDef.ApiDesc, GetRandomStr(4, ""))
	for index, value := range dfNames {
		if index == 0 {
			apiStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", value, value)
			numStr = fmt.Sprintf("%v", index+1)
		} else {
			apiStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", apiStr, value, value)
			numStr = fmt.Sprintf("%s,%v", numStr, index+1)
		}
	}

	scene.ApiList = apiStr
	scene.DataNumber = numStr
	scene.RunTime = 1
	scene.SceneType = 1
	scene.Remark = "自动生成"
	scene.Priority = 999
	scene.UserName = userName

	models.Orm.Table("playbook").Where("name = ? and product = ?", scene.Name, scene.Product).Find(&dbScene)
	if len(dbScene.Name) == 0 {
		err = models.Orm.Table("playbook").Create(scene).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	} else {
		dbScene.ApiList = apiStr
		dbScene.DataNumber = numStr
		dbScene.SceneType = scene.SceneType
		dbScene.RunTime = scene.RunTime
		err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	return
}

func GetLastFileLink(filePath string) (linkStr string) {
	lastFile := path.Base(filePath)
	b, _ := IsStrEndWithTimeFormat(lastFile)
	dirName := GetHistoryDataDirName(lastFile)
	if b {
		linkStr = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, lastFile, lastFile)
	} else {
		linkStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", lastFile, lastFile)
	}

	return
}
