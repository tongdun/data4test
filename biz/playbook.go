package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
				if runNum > 1 {
					Logger.Warning("退出循环多次测试执行")
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

	if err != nil {
		Logger.Error("%v", err)
		return
	}
	playbook := playbookInfo.GetPlaybook(source)

	if len(product) > 0 {
		productList, _ = GetProductInfo(product)
	} else {
		productList = productSceneInfo
	}

	for _, productInfo := range productList {
		_, _, errTmp := RepeatRunPlaybook(productInfo, playbook, playbookInfo.RunTime, mode, source, playbookInfo.Id)
		if errTmp != nil {
			if err != nil {
				err = fmt.Errorf("%s;%s", err, errTmp)
			} else {
				err = errTmp
			}
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
	var lastFilePath string

	runResp.TestResult, lastFilePath, err = RepeatRunPlaybook(productInfo, playbook, sceneModel.RunNum, "start", "consolePlaybook", "")
	if runResp.TestResult != "pass" {
		if err != nil {
			runResp.FailReason = fmt.Sprintf("%s", err)
		}
	}

	runResp.LastFile = path.Base(lastFilePath)

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
		err = yaml.Unmarshal(content, &sceneFileBase)
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
			err2 = json.Unmarshal(content, &sceneFile)
		case "yaml", "yml":
			err2 = yaml.Unmarshal(content, &sceneFile)
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

func (playbook Playbook) RunPlaybookContent(envType int, source string) (result, historyApi string, err error) {
	filePath := playbook.Apis[playbook.Tag]
	depOutVars, err := playbook.GetPlaybookDepParams()
	if err != nil {
		Logger.Error("%s", err)
	}

	dbData, err := GetDataByFileName(filePath, source)
	if err != nil {
		return
	}

	result, dst, errTmp := dbData.RunDataFile(filePath, playbook.Product, source, depOutVars)

	if errTmp != nil {
		result = "fail"
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}
	}

	switch source {
	case "historyAgain", "historyContinue", "history":
		historyApi = dst
	default:
		b, _ := IsStrEndWithTimeFormat(filePath)
		if b {
			dst = filePath
		}

		if strings.Contains(dst, "/history/") {
			historyApi = dst
		} else {
			historyApi = filePath // 如果未写到history中，
			dst = filePath
		}
	}

	err = WriteDataResultByFile(filePath, result, dst, playbook.Product, envType, errTmp)

	if errTmp != nil {
		if err != nil {
			err = fmt.Errorf("%s, %s", err, errTmp)
		} else {
			err = errTmp
		}
	}

	return
}

func (playbook Playbook) GetHistoryApiList() (apiStr string) {
	var lastFileName, dirName, suffix string
	if len(playbook.LastFile) > 0 {
		baseName := path.Base(playbook.LastFile)
		b, _ := IsStrEndWithTimeFormat(baseName)
		if b {
			dirName = GetHistoryDataDirName(baseName)
			suffix = GetStrSuffix(baseName)
			lastFileName = fmt.Sprintf("%s%s", dirName, suffix)
		} else {
			lastFileName = baseName
		}
	}

	rawApiList := playbook.Apis
	for index, item := range playbook.HistoryApis {
		// 如果执行过的数据中有空值，则跳过
		if len(item) == 0 {
			continue
		}

		if index == 0 {
			apiStr = path.Base(item)
		} else {
			apiStr = fmt.Sprintf("%s,%s", apiStr, path.Base(item))
		}
	}

	lastFileTag := len(rawApiList)
	for index, item := range rawApiList {
		rawDir := GetHistoryDataDirName(path.Base(item))
		lastFileDir := GetHistoryDataDirName(lastFileName)
		if rawDir == lastFileDir {
			lastFileTag = index
			break
		}
	}

	if lastFileTag < len(rawApiList) {
		for _, item := range rawApiList[lastFileTag+1:] {
			if len(apiStr) > 0 {
				apiStr = fmt.Sprintf("%s,%s", apiStr, path.Base(item))
			} else {
				apiStr = path.Base(item)
			}

		}
	}

	return
}

func (playbook Playbook) WritePlaybookResult(id, result, source string, envType int, errIn error) (err error) {
	var sceneRecode SceneRecord
	apiStr := playbook.GetHistoryApiList()
	lastFile := playbook.LastFile

	sceneRecode.ApiList = apiStr
	err = UpdateDBPlaybook(id, source, result, lastFile, errIn)

	sceneRecode.Name = playbook.Name
	if len(lastFile) > 0 {
		sceneRecode.LastFile = path.Base(lastFile)
	}

	sceneRecode.SceneType = playbook.SceneType
	sceneRecode.Result = result
	if errIn != nil {
		sceneRecode.FailReason = fmt.Sprintf("%v", errIn)
	}

	sceneRecode.Product = playbook.Product
	sceneRecode.EnvType = envType

	errTmp := WritePlaybookRecord(sceneRecode)
	if errTmp != nil {
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}
	}
	return
}

func UpdateDBPlaybook(dbId, source, result, lastFile string, errIn error) (err error) {
	var dbScene DbScene
	s, _ := strconv.Atoi(dbId)

	if source == "playbook" {
		models.Orm.Table("playbook").Where("id = ?", s).Find(&dbScene)

	} else if source == "ai_playbook" {
		models.Orm.Table("ai_playbook").Where("id = ?", s).Find(&dbScene)
	}

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
		prefix := GetHistoryDataDirName(path.Base(lastFile))
		suffix := GetStrSuffix(path.Base(lastFile))
		dbScene.LastFile = fmt.Sprintf("%s%s", prefix, suffix)
	}
	if errIn != nil {
		dbScene.FailReason = fmt.Sprintf("%s", errIn)
	} else {
		dbScene.FailReason = " "
	}

	if source == "playbook" {
		err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
	} else if source == "ai_playbook" {
		err = models.Orm.Table("ai_playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
	}

	if err != nil {
		Logger.Error("%v", err)
	}

	return
}

// WritePlaybookHistoryResult id为历史记录中的ID, mode为继续还是再来一次，envType为执行环境的类型
func (playbook Playbook) WritePlaybookHistoryResult(id, result, mode string, envType int, errIn error) (err error) {
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

	baseLastFile := path.Base(playbook.LastFile)

	if len(baseLastFile) == 0 {
		dbScene.LastFile = " " // 用空格字符串刷新数据
	} else {
		dbScene.LastFile = baseLastFile
	}

	if errIn != nil {
		dbScene.FailReason = fmt.Sprintf("%s", errIn)
	} else {
		dbScene.FailReason = " "
	}

	if mode == "continue" {
		curTime := time.Now()
		dbScene.UpdatedAt = curTime.Format(baseFormat)
		dbScene.ApiList = playbook.GetHistoryApiList()
		err = UpdateSceneRecord(dbScene)
	} else {
		var sceneRecode SceneRecord
		sceneRecode.Name = dbScene.Name
		sceneRecode.LastFile = dbScene.LastFile
		sceneRecode.SceneType = dbScene.SceneType
		sceneRecode.Result = dbScene.Result
		sceneRecode.FailReason = dbScene.FailReason
		sceneRecode.Product = dbScene.Product
		sceneRecode.EnvType = envType
		sceneRecode.ApiList = playbook.GetHistoryApiList()

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
	if len(dbScene.DataFileList) == 0 {
		err = fmt.Errorf("未找到[%v]场景，请核对", s)
		Logger.Error("%s", err)
		return
	}

	var filePath string
	fileNames := strings.Split(dbScene.DataFileList, ",")

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
	if len(dbScene.LastFile) > 0 {
		//playbook.LastFile = GetStrFromHtml(dbScene.LastFile)
		playbook.LastFile = dbScene.LastFile
	}

	playbook.Product = dbScene.Product
	playbook.SceneType = dbScene.SceneType
	if len(playbook.LastFile) != 0 {
		index := GetSliceIndex(playbook.Apis, playbook.LastFile)
		playbook.Tag = index

	}

	if playbook.Tag != -1 {
		playbook.HistoryApis = playbook.Apis[:playbook.Tag]
	}

	return
}

func GetPlRunInfo(source, id string) (dbScene DbScene, dbProduct []DbProduct, err error) {
	s, _ := strconv.Atoi(id)
	if source == "history" {
		models.Orm.Table("scene_test_history").Where("id = ?", s).Find(&dbScene)
		dbScene.RunTime = 1 // 执行次数默认设置为1
	} else if source == "ai_playbook" {
		models.Orm.Table("ai_playbook").Where("id = ?", s).Find(&dbScene)
		dbScene.RunTime = 1 // 执行次数默认设置为1
	} else {
		models.Orm.Table("playbook").Where("id = ?", s).Find(&dbScene)
	}

	if len(dbScene.DataFileList) == 0 {
		err = fmt.Errorf("未找到对应场景，请核对: %s", id)
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

func (playbook Playbook) GetPlaybookDepParams() (outputDict map[string][]interface{}, err error) {
	outputDict = make(map[string][]interface{})

	// 专用参数的优先级低于动态执行的同名参数
	var tmpStr []string
	models.Orm.Table("product").Where("product = ?", playbook.Product).Pluck("private_parameter", &tmpStr)
	privateParameter := make(map[string]interface{})
	if len(tmpStr) > 0 {
		if len(tmpStr[0]) > 2 {
			err = json.Unmarshal([]byte(tmpStr[0]), &privateParameter)
			if err != nil {
				Logger.Error("解析专用参数异常: %s", err)
				return
			}
		}
	}

	for k, v := range privateParameter {
		vStr := Interface2Str(v)
		var values []string
		//if _, ok := outputDict[k]; ok {
		//	outputDict[k] = outputDict[k][:0] // 同名参数进行重置, 私有参数的优先级同名优先级最高
		//}

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

	// 动态过程执行参数的取用
	var preApis []string
	for _, item := range playbook.HistoryApis {
		//fullPathApi := fmt.Sprintf("%s/%s", HistoryBasePath, item)   // 历史执行过的接口使用全路径
		if len(item) == 0 {
			continue
		}
		preApis = append(preApis, item)
	}

	for _, filePath := range preApis {
		if len(filePath) == 0 {
			continue
		}
		var sceneFile DataFile
		content, err1 := ioutil.ReadFile(filePath)
		if err1 != nil {
			err = err1
			Logger.Debug("filePath: %s", filePath)
			//Logger.Error("%s", err)  //在上层链路展示错误信息
			return
		}

		suffix := GetStrSuffix(filePath)
		switch suffix {
		case ".json":
			err1 = json.Unmarshal(content, &sceneFile)
		case ".yml", ".yaml":
			err1 = yaml.Unmarshal(content, &sceneFile)
		default:
			continue // 如果非标准结构数据，直接进入下一个数据文件的参数取用
		}

		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}

		for index, request := range sceneFile.Request {
			if len(request) == 0 {
				continue
			}

			requestMap := make(map[string]interface{})
			errTmp := json.Unmarshal([]byte(request), &requestMap)
			if errTmp != nil {
				Logger.Debug("request: %v", request)
				Logger.Warning("%v", errTmp)
				continue
			}

			for k, v := range requestMap {
				vType := fmt.Sprintf("%T", v)
				if index == 0 {
					outputDict[k] = outputDict[k][:0] // 如果是不同数据文件的请求参数，同名参数进行重置

				}

				if vType == "[]interface {}" { // 如果请求数据本来就是list的，需要解耦一层
					vList := v.([]interface{})
					outputDict[k] = vList
				} else {
					outputDict[k] = append(outputDict[k], v)
				}
			}
		}
		// Output的参数在请求参数之后，优先级更高
		if len(sceneFile.Output) > 0 {
			for k, v := range sceneFile.Output {
				outputDict[k] = v
			}
		}
	}

	if len(playbook.Apis) == 1 {
		selfApi := playbook.Apis[playbook.Tag]
		var selfScene DataFile
		var otherApis []string
		content, err1 := ioutil.ReadFile(selfApi)
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}

		suffix := GetStrSuffix(selfApi)
		switch suffix {
		case ".json":
			err1 = json.Unmarshal(content, &selfScene)
		case ".yml", ".yaml":
			err1 = yaml.Unmarshal(content, &selfScene)
		default:
			return // 如果非标准结构数据，直接返回
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

			suffix := GetStrSuffix(filePath)
			switch suffix {
			case ".json":
				err1 = json.Unmarshal(content, &sceneFile)
			case ".yml", ".yaml":
				err1 = yaml.Unmarshal(content, &sceneFile)
			default:
				continue // 如果非标准结构数据，直接进入下一个数据文件的参数取用
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

			for index, request := range sceneFile.Request {
				requestMap := make(map[string]interface{})
				errTmp := json.Unmarshal([]byte(request), &requestMap)
				if errTmp != nil {
					Logger.Debug("request: %v", request)
					Logger.Warning("%v", errTmp)
					continue
				}

				for k, v := range requestMap {
					if index == 0 {
						outputDict[k] = outputDict[k][:0] // 如果是不同数据文件的请求参数，同名参数进行重置
						outputDict[k] = append(outputDict[k], v)
					} else {
						outputDict[k] = append(outputDict[k], v)
					}
				}
			}
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
	if len(dbScene.DataFileList) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", id)
		Logger.Error("%s", err)
		return
	}

	var scene SceneWithNoUpdateTime
	scene.Product = dbScene.Product
	scene.Name = fmt.Sprintf("%s_复制", dbScene.Name)
	scene.DataFileList = dbScene.DataFileList
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

func ImportPlaybookFromExcel(id string) (newCount, oldCount int, err error) {
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

	dataList := strings.Split(dbScene.DataFileList, ",")
	for _, item := range dataList {
		var dataModel DepDataModel
		if len(item) == 0 {
			continue
		}
		dataModel.DataFile = item
		sceneInfo.DataList = append(sceneInfo.DataList, dataModel)
	}

	return
}

func GetAllPlaybook() (names []string, err error) {
	//models.Orm.Table("playbook").Order("created_at desc").Group("name").Pluck("name", &names)  //去重再优化
	models.Orm.Table("playbook").Order("created_at desc").Pluck("name", &names)
	if len(names) == 0 {
		Logger.Warning("暂无场景数据")
		return
	}
	return
}

func SaveScene(sceneSave SceneSaveModel, userName string) (err error) {
	var dbScene DbScene
	var scene SceneWithNoUpdateTime
	var apiStr, numStr string
	scene.Product = sceneSave.Product
	scene.Name = sceneSave.Name
	scene.UserName = userName
	for index, value := range sceneSave.DataList {
		if index == 0 {
			apiStr = value.DataFile
			numStr = fmt.Sprintf("%v", index+1)
		} else {
			apiStr = fmt.Sprintf("%s,%s", apiStr, value.DataFile)
			numStr = fmt.Sprintf("%s,%v", numStr, index+1)
		}
	}

	scene.DataFileList = apiStr
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
		dbScene.DataFileList = apiStr
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
			scene.DataFileList = values[1]
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
			priorityInt, err := strconv.Atoi(values[4])
			if err != nil {
				Logger.Error("%s", err)
			} else {
				scene.Priority = priorityInt
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
				apiStr = value
				numStr = fmt.Sprintf("%v", numList[index])
			} else {
				apiStr = fmt.Sprintf("%s,%s", apiStr, value)
				numStr = fmt.Sprintf("%s,%v", numStr, numList[index])
			}
		}
		dbScene.DataFileList = apiStr
		dbScene.DataNumber = numStr
		err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}
	return
}

func GetPlaybookApiStr(id string) (apiStr string) {
	var dbScene DbScene
	models.Orm.Table("playbook").Where("id = ?", id).Find(&dbScene)
	if len(dbScene.Name) == 0 {
		return
	}

	//doc, _ := goquery.NewDocumentFromReader(strings.NewReader(dbScene.ApiList))
	//handle := doc.Text()
	//afterTxt1 := strings.Replace(handle, ".yml", ".yml\n", -1)
	//afterTxt2 := strings.Replace(afterTxt1, ".json", ".json\n", -1)
	//apiStr = strings.Replace(afterTxt2, ".yaml", ".yaml\n", -1)

	return strings.Replace(dbScene.DataFileList, ",", "\r\n", -1)
}

func GetPlaybookApiList(id string) (apiList []string) {
	var dbScene DbScene
	models.Orm.Table("playbook").Where("id = ?", id).Find(&dbScene)
	if len(dbScene.Name) == 0 {
		return
	}

	//doc, _ := goquery.NewDocumentFromReader(strings.NewReader(dbScene.ApiList))
	//handle := doc.Text()
	//afterTxt1 := strings.Replace(handle, ".yml", ".yml,", -1)
	//afterTxt2 := strings.Replace(afterTxt1, ".json", ".json,", -1)
	//afterTxt := strings.Replace(afterTxt2, ".yaml", ".yaml,", -1)
	//apiList = strings.Split(afterTxt, ",")

	return strings.Split(dbScene.DataFileList, ",")
}

// ModifyPlaybookContent 场景数据升级函数,增加序号/标签
func ModifyPlaybookContent() (err error) {
	var ids []string
	models.Orm.Table("playbook").Order("id ASC").Where("id > 246").Pluck("id", &ids)
	for _, id := range ids {
		var dbScene DbScene
		models.Orm.Table("playbook").Where("id = ?", id).Find(&dbScene)
		doc, errTmp := goquery.NewDocumentFromReader(strings.NewReader(dbScene.DataFileList))
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

			dbScene.DataFileList = apiStr
			dbScene.DataNumber = numStr
			err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		}
	}

	return
}

// ModifyPlaybookApiList 场景数据升级函数, 关联数据去链超链存储
func ModifyPlaybookApiList() (err error) {
	var ids []int
	Logger.Info("开始更新")
	models.Orm.Table("playbook").Order("id ASC").Pluck("id", &ids)
	for _, id := range ids {
		var dbScene DbScene
		models.Orm.Table("playbook").Where("id = ?", id).Find(&dbScene)
		if len(dbScene.Name) == 0 {
			Logger.Error("未找到id: %v场景数据，请核对", id)
			return
		}

		if !strings.Contains(dbScene.DataFileList, "</a>") {
			Logger.Info("%d:%s场景无需更新", id, dbScene.Name)
			continue
		}

		doc, errTmp := goquery.NewDocumentFromReader(strings.NewReader(dbScene.DataFileList))
		if errTmp != nil {
			Logger.Debug("dbScene.ApiList:%s", dbScene.DataFileList)
			Logger.Error("%d:%s场景信息获取异常，%s", id, dbScene.Name, errTmp)
			err = errTmp
			continue
		}

		handle := doc.Text()
		afterTxt1 := strings.Replace(handle, ".yml", ".yml,", -1)
		afterTxt2 := strings.Replace(afterTxt1, ".json", ".json,", -1)
		afterTxt3 := strings.Replace(afterTxt2, ".sh", ".sh,", -1)
		apiStr := strings.Replace(afterTxt3, ".yaml", ".yaml,", -1)
		dbScene.DataFileList = apiStr
		err = models.Orm.Table("playbook").Where("id = ?", dbScene.Id).Update(dbScene).Error
		if err != nil {
			Logger.Error("%s", err)
		}

	}
	Logger.Info("结束更新")
	return
}

// CreatePlaybookByAPIId 同接口一键生成场景用例
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

	scene.DataFileList = apiStr
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
		dbScene.DataFileList = apiStr
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

func GetPlaybookEditTypeById(id string) (editType string) {
	var p DbScene
	models.Orm.Table("playbook").Where("id = ?", id).Find(&p)

	if len(p.Name) > 0 {
		editType = "input"
	} else {
		editType = "select"
	}

	return
}

func GetDataFileLinkByDataStr(pStr string) (linkStr string) {
	pList := strings.Split(pStr, ",")
	for _, item := range pList {
		if len(item) == 0 {
			continue
		}
		if len(linkStr) == 0 {
			linkStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", item, item) //跳详情，可自动点击编辑进行改写
		} else {
			linkStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", linkStr, item, item)
		}
	}
	return
}

func GetApiDetailLinkByApiStr(appNameRaw interface{}, apiStr string) (linkStr string) {
	appName := Interface2Str(appNameRaw)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(apiStr))
	handle := doc.Text()
	apiIdList := strings.Split(handle, "\n")
	for _, item := range apiIdList {
		if len(item) == 0 {
			continue
		}
		var ids []int
		models.Orm.Table("api_definition").Where("app = ? and api_id = ?", appName, item).Pluck("id", &ids)
		if len(ids) == 0 {
			Logger.Warning("未找到接口[%s], 请核对", item)
			if len(linkStr) == 0 {
				linkStr = item //跳详情，可自动点击编辑进行改写
			} else {
				linkStr = fmt.Sprintf("%s<br>%s", linkStr, item) // 如果被删了，显示普通信息，无链接
			}
		} else {
			if len(linkStr) == 0 {
				linkStr = fmt.Sprintf("<a href=\"/admin/info/api_definition/detail?__goadmin_detail_pk=%d\">%s</a>", ids[0], item) //跳详情，可自动点击编辑进行改写
			} else {
				linkStr = fmt.Sprintf("%s<br><a href=\"/admin/info/api_definition/detail?__goadmin_detail_pk=%d\">%s</a>", linkStr, ids[0], item)
			}
		}
	}
	return
}

func GetApiDetailLinkByApiRaw(appNameRaw interface{}, apiStr string) (linkStr string) {
	appName := Interface2Str(appNameRaw)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(apiStr))
	handle := doc.Text()
	apiIdList := strings.Split(handle, "\n")
	for _, item := range apiIdList {
		if len(item) == 0 {
			continue
		}
		var ids []int
		var apiId, otherInfo string
		if strings.Contains(item, " ") {
			apiId = strings.Split(item, " ")[0]
			otherInfo = strings.Split(item, " ")[1]
		}

		models.Orm.Table("api_definition").Where("app = ? and api_id = ?", appName, apiId).Pluck("id", &ids)
		if len(ids) == 0 {
			Logger.Warning("未找到接口[%s], 请核对", apiId)
			if len(linkStr) == 0 {
				linkStr = item //跳详情，可自动点击编辑进行改写
			} else {
				linkStr = fmt.Sprintf("%s<br>%s", linkStr, apiId, otherInfo) // 如果被删了，显示普通信息，无链接
			}
		} else {
			if len(linkStr) == 0 {
				linkStr = fmt.Sprintf("<a href=\"/admin/info/api_definition/detail?__goadmin_detail_pk=%d\">%s</a> %s", ids[0], apiId, otherInfo) //跳详情，可自动点击编辑进行改写
			} else {
				linkStr = fmt.Sprintf("%s<br><a href=\"/admin/info/api_definition/detail?__goadmin_detail_pk=%d\">%s</a> %s", linkStr, ids[0], apiId, otherInfo)
			}
		}
	}
	return
}

func GetApiAutoDataList(apiId, pkId string) (linkStr string) {
	var apiDef ApiDefDB
	models.Orm.Table("api_definition").Where("id = ?", pkId).Find(&apiDef)
	if len(apiDef.ApiId) == 0 {
		Logger.Warning("未找到%s:%s接口定义, 请核对", pkId, apiId)
		return apiId
	}

	var dataCount int
	models.Orm.Table("scene_data").Where("api_id = ? and app = ?", apiId, apiDef.App).Limit(1).Count(&dataCount)
	if dataCount == 0 {
		return apiId
	}
	encodeApiId := url.QueryEscape(apiId)
	encodeApp := url.QueryEscape("app[]")

	linkStr = fmt.Sprintf("<a href=\"/admin/info/scene_data?api_id=%s&%s=%s\">%s</a>", encodeApiId, encodeApp, apiDef.App, apiId) // 直接跑数据列表进行过滤

	return
}

func GetDataUsedInPlaybookList(dataName, pkId string) (linkStr string) {
	var dataDef DbSceneData
	models.Orm.Table("scene_data").Where("id = ?", pkId).Find(&dataDef)
	if len(dataDef.Name) == 0 {
		Logger.Warning("未找到%s:%s数据定义, 请核对", pkId, dataName)
		return dataName
	}

	var playbookCount int
	matchStr := "%" + dataName + "%"
	models.Orm.Table("playbook").Where("data_file_list like ?", matchStr).Limit(1).Count(&playbookCount)
	if playbookCount == 0 {
		return dataName
	}

	var playbookIdList []int
	models.Orm.Table("playbook").Where("data_file_list like ?", matchStr).Group("id").Pluck("id", &playbookIdList)
	encodeId := url.QueryEscape("id[]")
	var queryStr string
	for index, id := range playbookIdList {
		if index == 0 {
			queryStr = fmt.Sprintf("%s=%d", encodeId, id)
		} else {
			queryStr = fmt.Sprintf("%s&%s=%d", queryStr, encodeId, id)
		}
	}
	linkStr = fmt.Sprintf("<a href=\"/admin/info/playbook?%s\">%s</a>", queryStr, dataName) // 直接跑数据列表进行过滤
	return
}

func GetPlaybookUsedInTaskList(playbookName, pkId string) (linkStr string) {
	var playbookDef DbScene
	models.Orm.Table("playbook").Where("id = ?", pkId).Find(&playbookDef)
	if len(playbookDef.Name) == 0 {
		Logger.Warning("未找到%s:%s场景定义, 请核对", pkId, playbookName)
		return playbookName
	}

	var playbookCount int
	matchStr := "%" + playbookName + "%"
	models.Orm.Table("schedule").Where("scene_list like ?", matchStr).Limit(1).Count(&playbookCount)
	if playbookCount == 0 {
		return playbookName
	}

	var taskIdList []int
	models.Orm.Table("schedule").Where("scene_list like ?", matchStr).Group("task_name").Pluck("id", &taskIdList)

	encodeId := url.QueryEscape("id[]")

	var queryStr string
	for index, id := range taskIdList {
		if index == 0 {
			queryStr = fmt.Sprintf("%s=%d", encodeId, id)
		} else {
			queryStr = fmt.Sprintf("%s&%s=%d", queryStr, encodeId, id)
		}
	}
	linkStr = fmt.Sprintf("<a href=\"/admin/info/schedule?%s\">%s</a>", queryStr, playbookName) // 直接跑数据列表进行过滤

	return
}

func GetHistoryDataLinkByDataStr(pStr string) (linkStr string) {
	pList := strings.Split(pStr, ",")
	for _, item := range pList {
		if len(item) == 0 {
			continue
		}

		var dirName, itemLinkStr string
		b, num := IsStrEndWithTimeFormat(item)
		suffix := GetStrSuffix(item)
		if b {
			dirName = item[:len(item)-num-len(suffix)]
			itemLinkStr = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, item, item)
		} else {
			itemLinkStr = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", item, item)
		}

		if len(linkStr) == 0 {
			linkStr = itemLinkStr //跳详情，可自动点击编辑进行改写
		} else {
			linkStr = fmt.Sprintf("%s<br>%s", linkStr, itemLinkStr)
		}
	}
	return
}

func GetDataDetailLinkByDataStr(dStr string) (linkStr string) {
	dList := strings.Split(dStr, ",")
	for _, item := range dList {
		if len(item) == 0 {
			continue
		}
		var ids []int
		models.Orm.Table("scene_data").Where("file_name = ?", item).Pluck("id", &ids)
		if len(dList) == 0 {
			Logger.Warning("未找到数据文件[%s], 请核对", item)
			if len(linkStr) == 0 {
				linkStr = item //跳详情，可自动点击编辑进行改写
			} else {
				linkStr = fmt.Sprintf("%s<br>%s", linkStr, item) // 如果被删了，显示普通信息，无链接
			}
		} else {
			if len(linkStr) == 0 {
				linkStr = fmt.Sprintf("<a href=\"/admin/info/scene_data/detail?__goadmin_detail_pk=%d\">%s</a>", ids[0], item) //跳详情，可自动点击编辑进行改写
			} else {
				linkStr = fmt.Sprintf("%s<br><a href=\"/admin/info/scene_data/detail?__goadmin_detail_pk=%d\">%s</a>", linkStr, ids[0], item)
			}
		}
	}
	return
}
