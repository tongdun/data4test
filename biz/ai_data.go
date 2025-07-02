package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/titanous/json5"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime/debug"
	"strings"
)

func UseAiData(ids, userName string) (err error) {
	idList := strings.Split(ids, ",")
	var aiDataList []AiData
	models.Orm.Table("ai_data").Where("id in (?)", idList).Find(&aiDataList)
	if len(aiDataList) == 0 {
		Logger.Error("数据存在异常:%s，请核对~ ", ids)
		err = fmt.Errorf("据存在异常，请核对~")
		return
	}

	for _, item := range aiDataList {
		var tmpData SceneData
		models.Orm.Table("scene_data").Where("name = ?", item.Name).Find(&tmpData)
		tmpData.ApiId = item.ApiId
		tmpData.App = item.App
		tmpData.FileName = item.FileName
		tmpData.FileType = item.FileType
		tmpData.Content = item.Content
		tmpData.Result = item.Result
		tmpData.FailReason = item.FailReason
		tmpData.UserName = userName // 设置为当前操作用户
		tmpData.RunTime = 1
		if len(tmpData.Name) == 0 {
			tmpData.Name = item.Name
			err = models.Orm.Table("scene_data").Create(tmpData).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		} else {
			tmpData.Name = item.Name
			err = models.Orm.Table("scene_data").Where("name = ?", item.Name).Update(tmpData).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		}

		dst := fmt.Sprintf("%s/%s", DataBasePath, tmpData.FileName)

		errTmp := ioutil.WriteFile(dst, []byte(tmpData.Content), 0644)

		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%v,%v", err, errTmp)
			} else {
				err = errTmp
			}
		}
	}

	// 数据落库成功后再设置取用状态
	err = models.Orm.Table("ai_data").Where("id in (?)", idList).UpdateColumn(&StatusMgmt{UseStatus: 2}).Error
	if err != nil {
		Logger.Error("%v", err)
		return
	}

	return
}

func (input InputData) AICreateDataAndPlaybookByApiDefine(ids string) (err error) {
	content, appendContent, err := GetAiTemplateByName(input.AiTemplate, "2")
	if err != nil {
		return
	}
	apiDefs := GetApiDefineById(ids)
	var apiDefStr string
	for index, apiDef := range apiDefs {
		if index == 0 {
			tmpStr := fmt.Sprintf("接口定义如下：\n接口ID: %s, 所属应用: %s, 所属模块: %s, 接口描述: %s, 请求方法: %s, 请求路径: %s, Header参数: %v, Path参数：%v, Query参数: %v, Body参数: %v, Resp参数: %v;\n", apiDef.ApiId, apiDef.App, apiDef.ApiModule, apiDef.ApiDesc, apiDef.HttpMethod, apiDef.Path, apiDef.Header, apiDef.Path, apiDef.QueryParameter, apiDef.Body, apiDef.Response)
			apiDefStr = tmpStr
		} else {
			tmpStr := fmt.Sprintf("接口ID: %s, 所属应用: %s, 所属模块: %s, 接口描述: %v, 请求方法: %s, 请求路径: %s, Header参数: %v, Path参数：%v, Query参数: %v, Body参数: %v, Resp参数: %v;\n", apiDef.ApiId, apiDef.App, apiDef.ApiModule, apiDef.ApiDesc, apiDef.HttpMethod, apiDef.Path, apiDef.Header, apiDef.Path, apiDef.QueryParameter, apiDef.Body, apiDef.Response)
			apiDefStr = apiDefStr + tmpStr
		}

	}
	query := strings.Replace(content, "{需求定义}", apiDefStr, -1)
	go func(query, appendContent string, input InputData) {
		defer func() { // 如果遇到panic， 不影响主程序的运行
			if e := recover(); e != nil {
				Logger.Error("panic: %v", e)
				Logger.Error("stack: %v", string(debug.Stack()))
			}
		}()
		replyList, err := ConnetAIModel(query, appendContent, "", input.CreatePlatform, input.CreateUser)
		if err != nil {
			Logger.Debug("query: %v", query)
			Logger.Debug("replyList: %v", replyList)
		}

		dataList, playbookList, err := GetDataAndPlaybookFromReplyList(replyList)
		errTmp := input.SaveAIDataByDataListMap(dataList)

		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}

		errTmp = input.SaveAIPlaybookByPlaybookListMap(playbookList)

		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}

	}(query, appendContent, input)

	return
}

func GetDataAndPlaybookFromReplyList(replyList []string) (dataList, playbookList []map[string]interface{}, err error) {
	for index, reply := range replyList {
		dataSubList, playbookSubList, errTmp := GetDataAndPlaybookFromReply(reply)
		if errTmp != nil {
			if err != nil {
				err = fmt.Errorf("%s;%s", err, errTmp)
			} else {
				err = errTmp
			}
		}

		if len(dataSubList) == 0 {
			var errTmp error
			if len(replyList) == 1 {
				errTmp = fmt.Errorf("返回信息未匹配到数据，请核对~")
			} else {
				errTmp = fmt.Errorf("第%d步返回信息未匹配到数据，请核对~", index+1)
			}
			Logger.Error("%s", errTmp)
			if err == nil {
				err = errTmp
			} else {
				err = fmt.Errorf("%s; %s", err, errTmp)
			}
		} else {
			dataList = append(dataList, dataSubList...)
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
				err = fmt.Errorf("%s; %s", err, errTmp)
			}
		} else {
			playbookList = append(playbookList, playbookSubList...)
		}

	}
	return
}

func (input CommonExtend) SaveAIDataByDataListMap(dataList []map[string]interface{}) (err error) {
	if len(dataList) == 0 {
		return fmt.Errorf("无数据信息，请核对~")
	}

	for _, item := range dataList {
		aiData := input.AssembleAIData(item)
		errTmp := aiData.AddAiData()
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

func (input CommonExtend) AssembleAIData(aiRawData map[string]interface{}) (aiData AiData) {
	if len(input.IntroVersion) > 0 {
		aiData.Name = fmt.Sprintf("%s_%s", aiRawData["数据描述"].(string), input.IntroVersion)
	} else {
		aiData.Name = aiRawData["数据描述"].(string)
	}

	aiData.FileName = fmt.Sprintf("%s.yml", aiData.Name)
	aiData.App = aiRawData["所属应用"].(string)
	aiData.ApiId = aiRawData["接口ID"].(string)
	aiData.FileType = 1

	aiData.CreateUser = input.CreateUser
	aiData.CreatePlatform = input.CreatePlatform

	var dataFile DataFile
	dataFile.Name = aiData.Name
	dataFile.ApiId = aiData.ApiId
	dataFile.Api.App = aiData.App
	dataFile.Api.Module = aiRawData["所属模块"].(string)
	dataFile.Api.Description = aiRawData["接口描述"].(string)
	dataFile.Api.Method = aiRawData["请求方法"].(string)
	dataFile.Api.Path = aiRawData["请求路径"].(string)
	dataFile.Version = 1
	dataFile.IsParallel = "no"
	dataFile.IsRunPostApis = "no"
	dataFile.IsUseEnvConfig = "yes"
	dataFile.IsRunPreApis = "no"
	dataFile.Env.Protocol = "http"

	appModel := GetAppInfo(aiData.App)
	dataFile.Env.Prepath = appModel.Prefix

	dataFile.Single.Header = make(map[string]interface{})
	header := aiRawData["Header参数"].(map[string]interface{})
	for k, v := range header {
		varType := fmt.Sprintf("%T", v)
		var vList []interface{}
		if varType == "string" {
			dataFile.Single.Header[k] = v
		} else if varType == "[]interface {}" {
			vList = v.([]interface{})
			if len(vList) > 0 {
				dataFile.Single.Header[k] = vList[0]
			}
		}
	}

	dataFile.Single.Body = make(map[string]interface{})
	dataFile.Multi.Body = make(map[string][]interface{})
	body := aiRawData["Body参数"].(map[string]interface{})
	for k, v := range body {
		varType := fmt.Sprintf("%T", v)
		var vList []interface{}
		if varType == "[]interface {}" {
			vList = v.([]interface{})
		}
		if len(vList) == 1 {
			dataFile.Single.Body[k] = vList[0]
		} else if len(vList) > 1 {
			dataFile.Multi.Body[k] = vList
		}
	}
	dataFile.Single.Query = make(map[string]interface{})
	dataFile.Multi.Query = make(map[string][]interface{})
	query := aiRawData["Query参数"].(map[string]interface{})
	for k, v := range query {
		varType := fmt.Sprintf("%T", v)
		var vList []interface{}
		if varType == "[]interface {}" {
			vList = v.([]interface{})
		}
		if len(vList) == 1 {
			dataFile.Single.Query[k] = vList[0]
		} else if len(vList) > 1 {
			dataFile.Multi.Query[k] = vList
		}
	}

	dataFile.Single.Path = make(map[string]interface{})
	dataFile.Multi.Path = make(map[string][]interface{})
	path := aiRawData["Path参数"].(map[string]interface{})
	for k, v := range path {
		varType := fmt.Sprintf("%T", v)
		var vList []interface{}
		if varType == "[]interface {}" {
			vList = v.([]interface{})
		}
		if len(vList) == 1 {
			dataFile.Single.Path[k] = vList[0]
		} else if len(vList) > 1 {
			dataFile.Multi.Path[k] = vList
		}
	}

	dataContent, errTmp := yaml.Marshal(dataFile)
	if errTmp != nil {
		Logger.Debug("dataContent: %s", dataContent)
		Logger.Error("%s", errTmp)

	}

	aiData.Content = string(dataContent)

	return
}

func (aiData AiData) AddAiData() (err error) {
	var tmpData AiData
	models.Orm.Table("ai_data").Where("name = ? and  app= ? and source = ?", aiData.Name, aiData.App, aiData.CreatePlatform).Find(&tmpData)
	if len(tmpData.Name) > 0 {
		err = models.Orm.Table("ai_data").Where("name = ? and  app= ? and source = ?", aiData.Name, aiData.App, aiData.CreatePlatform).Update(aiData).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	} else {
		err = models.Orm.Table("ai_data").Create(&aiData).Error
		if err != nil {
			Logger.Error("aiData: %s", aiData)
			Logger.Error("%s", err)
		}
	}

	dst := fmt.Sprintf("%s/%s", AiDataBasePath, aiData.FileName)

	errTmp := ioutil.WriteFile(dst, []byte(aiData.Content), 0644)

	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%v,%v", err, errTmp)
		} else {
			err = errTmp
		}
	}

	return
}

func GetDataFromReply(reply string) (dataList []map[string]interface{}, err error) {
	// 编译正则表达式（含跨行匹配）
	//answerReg := regexp.MustCompile(`(\{[\s\S]*?\})`)  //匹配用例
	//dataReg := regexp.MustCompile(`"测试数据":[\s\S]\[[\s\S](\{[\s\S]*?[^{\w\S}]\})[\s\S]\][\s\S]`) // 1st
	//dataReg := regexp.MustCompile(`"测试数据":[\s\S](\[[\s\S]\{[\s\S]*?[^{\w\S}]\}[\s\S]\])[\s\S]`) // all
	//dataReg := regexp.MustCompile(`"测试数据":[\s\S](\[[\s\S]\{[\s\S]*?[^{\w\S}]\}[\s\S]\])[\s\S]`) // all,开始没有空格
	//dataReg := regexp.MustCompile(`"测试数据":[\s\S](\[[\s\S]\{[\s\S]*?[^{\w\S}]\}[\s\S]\])[\s\S]`) // all， 开始有空格
	dataReg := regexp.MustCompile(`"测试数据":[\s\S](\[[\s\n]*\{[\s\n]*[\s\S]*?[^{\w\S}]\}.?[\s\n]*\])`) // all， 开始有空格
	//dataReg := regexp.MustCompile(`[\n]\{[\s\S]*[^{\w\S}]\}`) // 同时匹配数据和场景

	// 查找所有匹配项
	dataMatch := dataReg.FindAllStringSubmatch(reply, -1)
	var targetStr string
	if len(dataMatch) > 0 {
		if len(dataMatch[0]) > 1 {
			targetStr = dataMatch[0][1]
		}
		if len(dataMatch[0]) > 2 {
			Logger.Warning("匹配到了多笔数据，请核对 ~")
			Logger.Debug("dataMatch: %s", dataMatch[0])
		}
	}

	err = json.Unmarshal([]byte(targetStr), &dataList)
	if err != nil {
		Logger.Debug("matchData: %v", targetStr)
		Logger.Error("%s", err)
	}

	return
}

func GetDataAndPlaybookFromReply(reply string) (dataList, playbookList []map[string]interface{}, err error) {
	dataReg := regexp.MustCompile(`[\n](\{[\s\S]*[^{\w\S}]\})`) // 同时匹配数据和场景

	// 查找所有匹配项
	dataMatch := dataReg.FindAllStringSubmatch(reply, -1)
	var targetStr string
	if len(dataMatch) > 0 {
		if len(dataMatch[0]) > 1 {
			targetStr = dataMatch[0][1]
		}
		if len(dataMatch[0]) > 2 {
			Logger.Warning("匹配到了多笔数据，请核对 ~")
			Logger.Debug("match: %s", dataMatch[0])
		}
	}
	//Logger.Debug("reply: %s", reply)
	//Logger.Debug("targetStr: %s", targetStr)
	respMatch := make(map[string]interface{})

	err = json5.Unmarshal([]byte(targetStr), &respMatch)
	if err != nil {
		Logger.Debug("matchData: %v", targetStr)
		Logger.Error("%s", err)
	}

	if v, ok := respMatch["测试数据"]; ok {
		dataTmpList := v.([]interface{})
		for _, subV := range dataTmpList {
			dataMap := subV.(map[string]interface{})
			dataList = append(dataList, dataMap)
		}
	}

	if v, ok := respMatch["测试场景"]; ok {
		playbookTmpList := v.([]interface{})
		for _, subV := range playbookTmpList {
			dataMap := subV.(map[string]interface{})
			playbookList = append(playbookList, dataMap)
		}
	}

	return
}

func (input ImportCommon) AICreateDataAndPlaybookByImport() (err error) {
	if len(input.RawReply) > 0 {
		err = input.SaveAIDataByRawRepy()
		if err != nil {
			Logger.Error("%s", err)
		}
	} else {
		replyList, errTmp := input.ConnectModel2GetMessage()
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			Logger.Debug("conversationId: %v", input.ConversationId)
			Logger.Debug("replyList: %v", replyList)
			return errTmp
		}

		if len(replyList) == 0 {
			return fmt.Errorf("未获取到回复信息，请核对")
		}

		dataList, playbookList, err := GetDataAndPlaybookFromReplyList(replyList)
		errTmp = input.SaveAIDataByDataListMap(dataList)

		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}

		errTmp = input.SaveAIPlaybookByPlaybookListMap(playbookList)

		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}
	}

	return
}

func (input ImportCommon) SaveAIDataByRawRepy() (err error) {
	dataList, err := GetDataFromReply(input.RawReply)
	if err != nil {
		Logger.Error("%s", err)
	}

	for _, item := range dataList {
		aiData := input.AssembleAIData(item)
		errTmp := aiData.AddAiData()
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

func (input CommonExtend) ConnectModel2GetMessage(conversationId string) (replyList []string, err error) {
	aiConnect, err := GetAIModelConnectInfo(input.CreatePlatform)
	if err != nil {
		return
	}

	replyList, err = aiConnect.CallModel2GetMessage(input.CreateUser, conversationId)
	if err != nil {
		Logger.Error("调用%s失败: %v", input.CreatePlatform, err)
		return
	}

	return
}

func UpdateAiDataStatus(id string) (err error) {
	err = models.Orm.Table("ai_data").Where("id = ?", id).UpdateColumn(&StatusMgmt{ModifyStatus: 2}).Error
	if err != nil {
		Logger.Error("%v", err)
		return
	}

	return
}

func BakOldAiDataVer(id, content, fileName string) (err error) {
	filePath := fmt.Sprintf("%s/%s", AiDataBasePath, fileName)

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
			models.Orm.Table("ai_data").Where("id = ?", id).Pluck("content", &dbContents)
			if len(dbContents) == 0 {
				err = fmt.Errorf("未找到[%v]数据，请核对", id)
				Logger.Error("%s", err)
				return
			}
			dbContent := dbContents[0]
			dbContent = strings.Replace(dbContent, rawStr, newVer, 1)
			models.Orm.Table("ai_data").Where("id = ?", id).Update("content", dbContent)
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

func (input InputData) AICreateDataByCreateDesc(createDesc, uploadFilePath string) (err error) {
	content, appendContent, err := GetAiTemplateByName(input.AiTemplate, "2")
	if err != nil {
		return
	}

	query := strings.Replace(content, "{需求定义}", createDesc, -1)
	go func(query, appendContent, uploadFilePath string, input InputData) {
		defer func() { // 如果遇到panic， 不影响主程序的运行
			if e := recover(); e != nil {
				Logger.Error("panic: %v", e)
				Logger.Error("stack: %v", string(debug.Stack()))
			}
		}()
		replyList, err := ConnetAIModel(query, appendContent, uploadFilePath, input.CreatePlatform, input.CreateUser)
		if err != nil {
			Logger.Debug("query: %v", query)
			Logger.Debug("replyList: %v", replyList)
		}

		dataList, playbookList, err := GetDataAndPlaybookFromReplyList(replyList)

		errTmp := input.SaveAIDataByDataListMap(dataList)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}

		errTmp = input.SaveAIPlaybookByPlaybookListMap(playbookList)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, errTmp)
			} else {
				err = errTmp
			}
		}

	}(query, appendContent, uploadFilePath, input)

	return
}

func (input InputData) AICreateDataAndPlaybookByCreateDesc(createDesc, uploadFilePath string) (err error) {
	content, appendContent, err := GetAiTemplateByName(input.AiTemplate, "2")
	if err != nil {
		return
	}

	query := strings.Replace(content, "{需求定义}", createDesc, -1)
	go input.CommonExtend.GetDataAndPlaybookFromModel(query, appendContent, uploadFilePath)

	return
}

func (input CommonExtend) GetDataAndPlaybookFromModel(query, appendContent, uploadFilePath string) {
	defer func() { // 如果遇到panic， 不影响主程序的运行
		if e := recover(); e != nil {
			Logger.Error("panic: %v", e)
			Logger.Error("stack: %v", string(debug.Stack()))
		}
	}()

	replyList, err := ConnetAIModel(query, appendContent, uploadFilePath, input.CreatePlatform, input.CreateUser)
	if err != nil {
		Logger.Debug("query: %v", query)
		Logger.Debug("replyList: %v", replyList)
	}
	dataList, playbookList, err := GetDataAndPlaybookFromReplyList(replyList)

	errTmp := input.SaveAIDataByDataListMap(dataList)
	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}
	}

	errTmp = input.SaveAIPlaybookByPlaybookListMap(playbookList)
	if errTmp != nil {
		Logger.Error("%s", errTmp)
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp)
		} else {
			err = errTmp
		}
	}
}

func AIOptimizeData(ids, optimizeDesc, aiPlatform, createUser string) (err error) {
	dbAiCases := GetDbAiCaseById(ids)
	var rawCaseList []map[string]string
	for _, item := range dbAiCases {
		tmpCase := make(map[string]string)
		tmpCase["用例编号"] = item.CaseNumber
		tmpCase["用例名称"] = item.CaseName
		tmpCase["用例类型"] = item.CaseType
		tmpCase["优先级"] = item.Priority
		tmpCase["所属模块"] = item.Module
		tmpCase["预置条件"] = item.PreCondition
		tmpCase["测试范围"] = item.TestRange
		tmpCase["测试步骤"] = item.TestSteps
		tmpCase["预期结果"] = item.ExpectResult
		tmpCase["是否支持自动化"] = item.Auto
		rawCaseList = append(rawCaseList, tmpCase)
	}
	data, _ := json.MarshalIndent(rawCaseList, "", "  ")
	aiCaseStr := string(data)
	query := fmt.Sprintf("%s\n初始用例如上，%s，并严格按原格式返回", aiCaseStr, optimizeDesc)
	go func(ids, optimizeDesc, aiPlatform string) {
		defer func() { // 如果遇到panic， 不影响主程序的运行
			if e := recover(); e != nil {
				Logger.Error("panic: %v", e)
				Logger.Error("stack: %v", string(debug.Stack()))
			}
		}()

		replyList, err := ConnetAIModel(query, "", "", aiPlatform, createUser)
		if err != nil {
			Logger.Debug("query: %v", query)
			Logger.Debug("replyList: %v", replyList)
		}
		if len(rawCaseList) > 0 {
			_ = UpdateAICaseByOptimize(ids, replyList[0])
		} else {
			Logger.Warning("未获取到回复信息，请核对~")
		}

	}(ids, optimizeDesc, aiPlatform)

	return
}

func AiDataTest(ids string, analysisInput AnalysisDataInput) (err error) {
	idList := strings.Split(ids, ",")
	var aiDataList []AiData
	models.Orm.Table("ai_data").Where("id in (?)", idList).Find(&aiDataList)
	for _, item := range aiDataList {
		filePath := fmt.Sprintf("%s/%s", AiDataBasePath, item.FileName)
		df, err1 := RunAiData(item.App, analysisInput.Product, filePath, nil, []byte(item.Content))
		if err1 != nil {
			Logger.Error("%s", err1)
			if err != nil {
				err = fmt.Errorf("%s; %s", err, err1)
			} else {
				err = err1
			}
		}

		dst, _ := GetResultFilePath(filePath)

		b, _ := IsStrEndWithTimeFormat(filePath)
		if b {
			dst = filePath
		}

		if strings.Contains(dst, "/history/") {
		} else {
			dst = filePath
		}

		df.AnalysisWithLLM(dst, analysisInput)

	}

	return
}

func RunAiData(app, product, filePath string, depOutVars map[string][]interface{}, contentByte []byte) (df DataFile, err error) {
	if strings.HasSuffix(filePath, ".json") {
		err = json.Unmarshal(contentByte, &df)
	} else {
		err = yaml.Unmarshal(contentByte, &df)
	}

	if err != nil {
		return
	}

	var envConfig EnvConfig
	var targetApp string

	if len(app) > 0 {
		targetApp = app
	} else if len(df.Api.App) > 0 {
		targetApp = df.Api.App
	}

	envConfig, _ = GetEnvConfig(targetApp, "data")

	depOutVarsTmp, err1 := df.GetDepParams()
	if err1 != nil {
		Logger.Error("%s", err1)
		if err != nil {
			err = fmt.Errorf("%s;%s", err, err1)
		} else {
			err = err1
		}
	}
	if depOutVars == nil {
		depOutVars = make(map[string][]interface{})
	}

	for k, v := range depOutVarsTmp {
		if _, ok := depOutVars[k]; !ok {
			depOutVars[k] = v
		}
	}

	if len(product) > 0 {
		sceneEnvConfig, errTmp := GetEnvConfig(product, "scene")
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
		return
	}

	lang := GetRequestLangage(header)

	var querys, bodys []map[string]interface{}
	var bodyList []interface{}
	var urls []string
	var rHeader map[string]interface{}

	_, errTmp := yaml.Marshal(df)
	if errTmp != nil {
		Logger.Error("%v", errTmp)
		err = errTmp
		return
	}

	contentStr, errTmp := GetAfterContent(lang, string(contentByte), depOutVars)
	if strings.Contains(contentStr, "is_var_strong_check: \"no\"") {
		Logger.Warning("%s数据开启参数弱校验，请自行保证所需依赖参数的定义", filePath)
		errTmp = nil
	}
	if errTmp != nil {
		Logger.Debug("rawContent:\n%s", string(contentStr))
		Logger.Debug("afterContent:\n%s", contentStr)
		err = errTmp
		return
	}
	errTmp = yaml.Unmarshal([]byte(contentStr), &df)
	if errTmp != nil {
		Logger.Debug("\nrawContent: %s", string(contentByte))
		Logger.Debug("\nafterContent: %s", contentStr)
		Logger.Error("%v", errTmp)
		err = errTmp
		return
	}

	df.Single.Header = header

	urls, errTmp = df.GetUrl(envConfig)
	if errTmp != nil {
		Logger.Debug("fileName: %s", path.Base(filePath))
		Logger.Error("%v", errTmp)
		err = errTmp
		return
	}

	df.Urls = urls

	querys = df.GetQuery()

	bodys, bodyList = df.GetBody()
	rHeader = df.Single.RespHeader

	// 后续可优化，有依赖和无依赖进行控制
	go df.CreateDataOrderByKey(lang, filePath, depOutVars) // 无依赖，异步执行生成动作：create_xxx
	_ = df.RecordDataOrderByKey(bodys)                     // 有依赖，同步执行记录动作：record_xxx
	_ = df.ModifyFileWithData(bodys)                       // 有依赖，同步执行模板动作：modify_file

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	var resList [][]byte
	var errs []error
	tag := 0
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
			df.Response = append(df.Response, string(res))
			errs = append(errs, err)
			_ = df.SetSleepAction()
		}
	}

	return
}

func (df DataFile) AnalysisWithLLM(dst string, input AnalysisDataInput) (result string, err error) {
	var analysisDesc string
	analysisDesc = fmt.Sprintf("接口ID: %s, 所属应用: %s, 所属模块: %s, 接口描述: %s, 请求方法: %s, 请求路径: %s", df.ApiId, df.Api.App, df.Api.Module, df.Api.Description, df.Api.Method, df.Api.Path)

	for index, item := range df.Response {
		if len(df.Request) == 0 {
			analysisDesc = fmt.Sprintf("%s,\n请求数据%d: 空,\n返回数据%d: %s", analysisDesc, index+1, index+1, item)
		} else if len(df.Request) > index {
			analysisDesc = fmt.Sprintf("%s,\n请求数据%d: %s,\n返回数据%d: %s", analysisDesc, index+1, df.Request[index], index+1, item)
		} else {
			Logger.Warning("请求数据与返回数据数量不匹配，请核对~")
		}
	}

	content, appendContent, err := GetAiTemplateByName(input.AiTemplate, "7")
	if err != nil {
		return
	}

	query := strings.Replace(content, "{分析定义}", analysisDesc, -1)

	go func(query, appendContent, dst string, input AnalysisDataInput, df DataFile) {
		defer func() { // 如果遇到panic， 不影响主程序的运行
			if e := recover(); e != nil {
				Logger.Error("panic: %v", e)
				Logger.Error("stack: %v", string(debug.Stack()))
			}
		}()
		replyList, err := ConnetAIModel(query, appendContent, "", input.CreatePlatform, input.CreateUser)
		if err != nil {
			Logger.Debug("query: %v", query)
			Logger.Debug("replyList: %v", replyList)
		}
		testResult, failReason, assert, err := input.SaveAIIssueByRepyList(df, replyList)
		if err != nil {
			Logger.Error("%s", err)
		}

		df.Assert = assert
		df.TestResult = testResult
		df.FailReason = failReason
		dbData, err := df.UpdateAiDataContent()
		if err != nil {
			return
		}
		df.UpdateAiDataFileResult(dst, dbData.FileName)
		envType := GetEnvTypeByName(input.Product)

		err = RecordDataHistory(dst, input.Product, "ai", envType, dbData)
	}(query, appendContent, dst, input, df)

	return
}

func (df DataFile) UpdateAiDataContent() (dbData DbSceneData, err error) {
	models.Orm.Table("ai_data").Where("name = ? and  app= ?", df.Name, df.Api.App).Find(&dbData)
	if len(dbData.Name) > 0 {
		var rawDf DataFile
		if strings.HasSuffix(dbData.FileName, ".json") {
			err = json.Unmarshal([]byte(dbData.Content), &rawDf)
		} else {
			err = yaml.Unmarshal([]byte(dbData.Content), &rawDf)
		}

		if err != nil {
			Logger.Error("%s", err)
		}

		rawDf.Assert = df.Assert

		afterContent, _ := yaml.Marshal(rawDf)
		dbData.Content = string(afterContent)

		var testResult string
		for _, item := range df.TestResult {
			if item == "fail" || item == "失败" {
				testResult = "fail"
				break
			}
		}

		if len(df.TestResult) > 0 && len(testResult) == 0 {
			testResult = "pass"
		}
		dbData.Result = testResult

		for index, item := range df.FailReason {
			if index == 0 {
				dbData.FailReason = item
			} else {
				dbData.FailReason = fmt.Sprintf("%s; %s", dbData.FailReason, item)
			}
		}
		err = models.Orm.Table("ai_data").Where("name = ? and  app= ?", df.Name, df.Api.App).Update(dbData).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	return
}

func (df DataFile) UpdateAiDataFileResult(dst, fileName string) {
	filePath := fmt.Sprintf("%s/%s", AiDataBasePath, fileName)
	afterContent, _ := yaml.Marshal(df)
	// 更新数据文件
	errTmp := ioutil.WriteFile(filePath, afterContent, 0644)
	if errTmp != nil {
		Logger.Error("%s", errTmp)

	}

	// 写历史数据
	errTmp = ioutil.WriteFile(dst, afterContent, 0644)

	if errTmp != nil {
		Logger.Error("%s", errTmp)
	}
}
