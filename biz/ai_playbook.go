package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func GetPlaybookFromReply(reply string) (playbookList []map[string]interface{}, err error) {
	//playbookReg := regexp.MustCompile(`"测试场景":[\s\S](\[[\s\S]\{[\s\S]*?[^{\w\S}]\}[\s\S]\])[\s\S]`)
	//playbookReg := regexp.MustCompile(`"测试场景":[\s\S](\[[\s\S]\{[\s\S]*?[^{\w\S}]\}[\s\S]\])[\s\S]`)
	playbookReg := regexp.MustCompile(`"测试场景":[\s\S](\[[\s\n]*\{[\s\n]*[\s\S]*?[^{\w\S}]\}.?[\s\n]*\])`)

	// 查找所有匹配项
	playbookMatch := playbookReg.FindAllStringSubmatch(reply, -1)

	var targetStr string
	if len(playbookMatch) > 0 {
		if len(playbookMatch[0]) > 1 {
			targetStr = playbookMatch[0][1]
		}
		if len(playbookMatch[0]) > 2 {
			Logger.Warning("匹配到了多笔数据，请核对 ~")
			Logger.Debug("dataMatch: %s", playbookMatch[0])
		}
	}

	err = json.Unmarshal([]byte(targetStr), &playbookList)
	if err != nil {
		Logger.Debug("matchPlaybook: %v", targetStr)
		Logger.Debug("playbookList: %v", playbookList)
		Logger.Error("%s", err)
	}

	return
}

func GetPlaybookFromReplyList(replyList []string) (playbookList []map[string]interface{}, err error) {
	for index, reply := range replyList {
		playbookSubList, errTmp := GetPlaybookFromReply(reply)
		if errTmp != nil {
			if err != nil {
				err = fmt.Errorf("%s;%s", err, errTmp)
			} else {
				err = errTmp
			}
		}
		if len(playbookSubList) == 0 {
			var errTmp error
			if len(replyList) == 1 {
				errTmp = fmt.Errorf("返回信息未匹配到场景，请核对~")
			} else {
				errTmp = fmt.Errorf("第%d步返回信息未匹配到场景，请核对~", index+1)
			}
			Logger.Error("%s", errTmp)
			if err == nil {
				err = errTmp
			} else {
				err = fmt.Errorf("%s;%s", err, errTmp)
			}
		}
		playbookList = append(playbookList, playbookSubList...)

	}
	return
}

func (input CommonExtend) SaveAIPlaybookByRepyList(replyList []string) (err error) {
	playbookList, err := GetPlaybookFromReplyList(replyList)
	if err != nil {
		Logger.Error("%s", err)
	}

	for _, item := range playbookList {
		aiPlaybook := input.AssembleAIPlaybook(item)
		errTmp := aiPlaybook.AddAiPlaybook()
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

func (input CommonExtend) SaveAIPlaybookByPlaybookListMap(playbookList []map[string]interface{}) (err error) {
	if len(playbookList) == 0 {
		return fmt.Errorf("未匹配到场景信息，请核对~")
	}

	for _, item := range playbookList {
		aiPlaybook := input.AssembleAIPlaybook(item)
		errTmp := aiPlaybook.AddAiPlaybook()
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

func (input CommonExtend) AssembleAIPlaybook(aiRawPlaybook map[string]interface{}) (aiPlaybook AiPlaybook) {
	if len(input.IntroVersion) > 0 {
		aiPlaybook.PlaybookDesc = fmt.Sprintf("%s_%s", aiRawPlaybook["场景描述"].(string), input.IntroVersion)
	} else {
		aiPlaybook.PlaybookDesc = aiRawPlaybook["场景描述"].(string)
	}

	rawDataList := aiRawPlaybook["关联数据"].([]interface{})

	var dataStr string
	for index, item := range rawDataList {
		var dataFileName string
		if len(input.IntroVersion) > 0 {
			dataFileName = fmt.Sprintf("%s_%s.yml", item, input.IntroVersion)
		} else {
			dataFileName = fmt.Sprintf("%s.yml", item)
		}

		if index == 0 {
			dataStr = dataFileName
		} else {
			dataStr = fmt.Sprintf("%s,%s", dataStr, dataFileName)
		}

	}
	aiPlaybook.PlaybookType = "1"
	aiPlaybook.Priority = "999"
	aiPlaybook.DataFileList = dataStr
	aiPlaybook.Product = input.Product
	aiPlaybook.CreateUser = input.CreateUser
	aiPlaybook.CreatePlatform = input.CreatePlatform

	return
}

func (aiPlaybook AiPlaybook) AddAiPlaybook() (err error) {
	var tmpPlaybook AiPlaybook
	models.Orm.Table("ai_playbook").Where("playbook_desc = ? and  product = ? and source = ?", aiPlaybook.PlaybookDesc, aiPlaybook.Product, aiPlaybook.CreatePlatform).Find(&tmpPlaybook)
	if len(tmpPlaybook.PlaybookDesc) > 0 {
		err = models.Orm.Table("ai_playbook").Where("playbook_desc = ? and  product = ? and source = ?", aiPlaybook.PlaybookDesc, aiPlaybook.Product, aiPlaybook.CreatePlatform).Update(aiPlaybook).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	} else {
		err = models.Orm.Table("ai_playbook").Create(&aiPlaybook).Error
		if err != nil {
			Logger.Error("aiPlaybook: %s", aiPlaybook)
			Logger.Error("%s", err)
		}
	}

	return
}

func (input ImportCommon) SaveAIPlaybookByRawRepy() (err error) {
	playbookList, err := GetPlaybookFromReply(input.RawReply)
	if err != nil {
		Logger.Error("%s", err)
	}

	for _, item := range playbookList {
		aiPlaybook := input.CommonExtend.AssembleAIPlaybook(item)
		errTmp := aiPlaybook.AddAiPlaybook()
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

func GetAiDataFileLinkByDataStr(pStr string) (linkStr string) {
	pList := strings.Split(pStr, ",")
	for _, item := range pList {
		if len(item) == 0 {
			continue
		}
		if len(linkStr) == 0 {
			linkStr = fmt.Sprintf("<a href=\"/admin/fm/ai_data/preview?path=/%s\">%s</a>", item, item) //跳详情，可自动点击编辑进行改写
		} else {
			linkStr = fmt.Sprintf("%s<br><a href=\"/admin/fm/ai_data/preview?path=/%s\">%s</a>", linkStr, item, item)
		}
	}
	return
}

func (input InputPlaybook) AICreateDataAndPlaybookByCreateDesc(createDesc, uploadFilePath string) (err error) {
	content, appendContent, err := GetAiTemplateByName(input.AiTemplate, "3")
	if err != nil {
		return
	}

	query := strings.Replace(content, "{需求定义}", createDesc, -1)
	go input.CommonExtend.GetDataAndPlaybookFromModel(query, appendContent, uploadFilePath)
	//go func(query, appendContent, uploadFilePath string, input InputPlaybook) {
	//	defer func() { // 如果遇到panic， 不影响主程序的运行
	//		if e := recover(); e != nil {
	//			Logger.Error("panic: %v", e)
	//			Logger.Error("stack: %v", string(debug.Stack()))
	//		}
	//	}()
	//	replyList, err := ConnetAIModel(query, appendContent, uploadFilePath, input.CreatePlatform, input.CreateUser)
	//	if err != nil {
	//		Logger.Debug("query: %v", query)
	//		Logger.Debug("replyList: %v", replyList)
	//	}
	//
	//	dataList, playbookList, err := GetDataAndPlaybookFromReplyList(replyList)
	//
	//	errTmp := input.SaveAIDataByDataListMap(dataList)
	//	if errTmp != nil {
	//		Logger.Error("%s", errTmp)
	//		if err != nil {
	//			err = fmt.Errorf("%s; %s", err, errTmp)
	//		} else {
	//			err = errTmp
	//		}
	//	}
	//
	//	errTmp = input.SaveAIPlaybookByPlaybookListMap(playbookList)
	//	if errTmp != nil {
	//		Logger.Error("%s", errTmp)
	//		if err != nil {
	//			err = fmt.Errorf("%s; %s", err, errTmp)
	//		} else {
	//			err = errTmp
	//		}
	//	}
	//}(query, appendContent, uploadFilePath, input)

	return
}

func AiPlaybookTest(idStr, source string, analysisInput AnalysisDataInput) (err error) {
	idList := strings.Split(idStr, ",")
	for _, id := range idList {
		if len(id) == 0 {
			continue
		}
		var aiPlaybook AiPlaybook
		models.Orm.Table("ai_playbook").Where("id = (?)", id).Find(&aiPlaybook)
		var playbook Playbook
		playbook.Name = aiPlaybook.PlaybookDesc
		playbook.LastFile = aiPlaybook.LastFile
		playbook.Product = aiPlaybook.Product
		playbook.Apis = strings.Split(aiPlaybook.DataFileList, ",")
		playbook.SceneType, _ = strconv.Atoi(aiPlaybook.PlaybookType)
		_, _, errTmp := playbook.RunAiPlaybook(id, source, analysisInput)

		if errTmp != nil {
			if err == nil {
				err = errTmp
			} else {
				err = fmt.Errorf("%s, %s", err, errTmp)
			}
		}
		//aiPlaybook.Result = result
		//aiPlaybook.LastFile = lastFile
		//aiPlaybook.FailReason = fmt.Sprintf("%s", errTmp)

	}

	return
}

func (playbook Playbook) RunAiPlaybook(dbId, source string, analysisInput AnalysisDataInput) (result, lastFile string, err error) {
	var runApis []string
	var tag int
	envType := GetEnvTypeByName(playbook.Product)
	isFail := 0
	result = "fail"
	runApis = playbook.Apis
	switch playbook.SceneType {
	case 1, 2:
		for k := range runApis {
			playbook.Tag = tag + k
			subResult, historyApi, errTmp := playbook.RunAiPlaybookContent(source, analysisInput)
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
				errTmp = playbook.WritePlaybookResult(dbId, subResult, source, envType, errTmp)
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
			subResult, historyApi, errTmp := playbook.RunAiPlaybookContent(source, analysisInput)
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
			errTmp := playbook.WritePlaybookResult(dbId, tmpResult, source, envType, err) // 串行继续时，无最近执行的文件
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
				subResult, historyApi, errTmp := inPlaybook.RunAiPlaybookContent(source, analysisInput)
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
			}(playbook, dbId, tag, k, envType, err)
		}

		wg.Wait()
		if isFail > 0 {
			tmpResult := "fail"
			errTmp := playbook.WritePlaybookResult(dbId, tmpResult, source, envType, err) // 并发模式时，无最近执行的文件
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
		result, err = CompareResult(playbook.HistoryApis, "yaml")
	}
	playbook.LastFile = lastFile
	err = playbook.WritePlaybookResult(dbId, result, source, envType, err)
	if err != nil {
		Logger.Error("%v", err)
		return
	}

	if result != "pass" {
		err = fmt.Errorf("测试 %v", result)
	}
	return
}

func (playbook Playbook) RunAiPlaybookContent(source string, analysisInput AnalysisDataInput) (result, historyApi string, err error) {
	filePath := playbook.Apis[playbook.Tag]
	depOutVars, err := playbook.GetPlaybookDepParams()
	if err != nil {
		Logger.Error("%s", err)
	}

	dbData, err := GetDataByFileName(filePath, source)
	if err != nil {
		return
	}

	df, err1 := RunAiData(dbData.App, analysisInput.Product, filePath, depOutVars, []byte(dbData.Content))
	if err1 != nil {
		Logger.Error("%s", err1)
		if err != nil {
			err = fmt.Errorf("%s; %s", err, err1)
		} else {
			err = err1
		}
	}

	dst, err := GetResultFilePath(filePath)
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
	result, err = df.AnalysisWithLLM(dst, analysisInput)

	return
}
