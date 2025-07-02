package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

func (aiCase AITestCase) AddAiCase() (err error) {
	var tmpCase AITestCase
	models.Orm.Table("ai_case").Where("case_number = ? and  case_name= ? and product = ? and source = ?", aiCase.CaseNumber, aiCase.CaseName, aiCase.Product, aiCase.CreatePlatform).Find(&tmpCase)
	if len(tmpCase.CaseName) > 0 {
		err = models.Orm.Table("ai_case").Where("case_number = ? and  case_name= ? and product = ? and source = ?", aiCase.CaseNumber, aiCase.CaseName, aiCase.Product, aiCase.CreatePlatform).Update(aiCase).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	} else {
		err = models.Orm.Table("ai_case").Create(&aiCase).Error

		if err != nil {
			Logger.Error("aiCase: %s", aiCase)
			Logger.Error("%s", err)
		}
	}

	return
}

func (input ImportCommon) SaveAICaseByRawRepy() (err error) {
	caseList, err := GetCaseFromReply(input.RawReply)
	if err != nil {
		Logger.Error("%s", err)
	}

	for _, item := range caseList {
		aiCase := input.CommonExtend.AssembleAICase(item)
		errTmp := aiCase.AddAiCase()
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

func (input CommonExtend) SaveAICaseByRepyList(replyList []string) (err error) {
	caseList, err := GetCaseFromReplyList(replyList)
	if err != nil {
		Logger.Error("%s", err)
	}

	for _, item := range caseList {
		aiCase := input.AssembleAICase(item)
		errTmp := aiCase.AddAiCase()
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

func GetCaseFromReplyList(replyList []string) (caseList []map[string]string, err error) {
	for index, reply := range replyList {
		caseSubList, errTmp := GetCaseFromReply(reply)
		if errTmp != nil {
			if err != nil {
				err = fmt.Errorf("%s;%s", err, errTmp)
			} else {
				err = errTmp
			}
		}
		if len(caseSubList) == 0 {
			errTmp := fmt.Errorf("第%d步返回信息未匹配到用例，请核对~", index+1)
			Logger.Error("%s", errTmp)
			if err == nil {
				err = errTmp
			} else {
				err = fmt.Errorf("%s;%s", err, errTmp)
			}
		}
		caseList = append(caseList, caseSubList...)

	}
	return
}

func GetCaseFromReply(reply string) (caseList []map[string]string, err error) {
	// 编译正则表达式（含跨行匹配）
	//answerReg := regexp.MustCompile(`(\{[\s\S]*?\})`)
	answerReg := regexp.MustCompile(`(\{[\s\S]*?[^{\w}]\})`)

	// 查找所有匹配项
	answerMatch := answerReg.FindAllStringSubmatch(reply, -1)
	if len(answerMatch) == 0 {
	}
	for _, match := range answerMatch {
		for _, item := range match {
			testCase := make(map[string]string)
			err = json.Unmarshal([]byte(item), &testCase)
			if err != nil {
				Logger.Debug("matchCase: %v", item)
				Logger.Error("%s", err)
				continue
			}
			caseList = append(caseList, testCase)
		}

	}

	return
}

func (input CommonExtend) AssembleAICase(testCase map[string]string) (aiCase AITestCase) {
	rKeyList := []string{"用例编号", "用例名称", "用例类型", "优先级", "所属模块", "预置条件", "测试范围", "测试步骤", "预期结果", "是否支持自动化"}
	tKeyList := []string{"case_number", "case_name", "case_type", "priority", "module", "pre_condition", "test_range", "test_steps", "expect_result", "auto"}
	aiTmpCase := make(map[string]string)
	for index, k := range rKeyList {
		if v, ok := testCase[k]; ok {
			aiTmpCase[tKeyList[index]] = v
		}
	}

	aiCase.CaseNumber = aiTmpCase["case_number"]
	aiCase.CaseName = aiTmpCase["case_name"]
	aiCase.CaseType = aiTmpCase["case_type"]
	aiCase.Priority = aiTmpCase["priority"]
	aiCase.PreCondition = aiTmpCase["pre_condition"]
	aiCase.TestRange = aiTmpCase["test_range"]
	aiCase.TestSteps = aiTmpCase["test_steps"]
	aiCase.Module = aiTmpCase["module"]
	aiCase.ExpectResult = aiTmpCase["expect_result"]
	if aiTmpCase["auto"] == "否" {
		aiCase.Auto = "0"
	} else {
		aiCase.Auto = "1"
	}
	aiCase.CreateUser = input.CreateUser
	aiCase.IntroVersion = input.IntroVersion
	aiCase.Product = input.Product
	aiCase.CreatePlatform = input.CreatePlatform

	return
}

func UpdateAICaseByOptimize(ids, reply string) (err error) {
	idList := strings.Split(ids, ",")
	caseAfterList, err := GetCaseFromReply(reply)
	for index, afterCase := range caseAfterList {
		var caseBase CaseBase

		caseBase.CaseNumber = afterCase["用例编号"]
		caseBase.CaseName = afterCase["用例名称"]
		caseBase.CaseType = afterCase["用例类型"]
		caseBase.Priority = afterCase["优先级"]
		caseBase.Module = afterCase["所属模块"]
		caseBase.PreCondition = afterCase["预置条件"]
		caseBase.TestRange = afterCase["测试范围"]
		caseBase.TestSteps = afterCase["测试步骤"]
		caseBase.ExpectResult = afterCase["预期结果"]
		caseBase.Auto = afterCase["是否支持自动化"]

		errTmp := models.Orm.Table("ai_case").Where("id", idList[index]).Update(caseBase).Error
		if errTmp != nil {
			Logger.Debug("afterCase: %+v", afterCase)
			Logger.Debug("caseBase: %+v", caseBase)
			Logger.Error("%v", errTmp)
			if err != nil {
				err = fmt.Errorf("%s;%s", err, errTmp)
			} else {
				err = errTmp
			}

		}
	}
	return
}

func (input InputCase) AICreateCaseByApiDefine(ids string) (err error) {
	content, appendContent, err := GetAiTemplateByName(input.AiTemplate, "1")
	if err != nil {
		return
	}
	apiDefs := GetApiDefineById(ids)
	var apiDefStr string
	for index, apiDef := range apiDefs {
		if index == 0 {
			tmpStr := fmt.Sprintf("接口定义如下：\n接口ID: %s, 所属模块: %s, 接口描述: %v, 请求方法: %s, 请求路径: %s, Header参数: %v, Path参数：%v, Query参数: %v, Body参数: %v, Resp参数: %v;\n", apiDef.ApiId, apiDef.ApiModule, apiDef.ApiDesc, apiDef.HttpMethod, apiDef.Path, apiDef.Header, apiDef.Path, apiDef.QueryParameter, apiDef.Body, apiDef.Response)
			apiDefStr = tmpStr
		} else {
			tmpStr := fmt.Sprintf("接口ID: %s, 所属模块: %s, 接口描述: %v, 请求方法: %s, 请求路径: %s, Header参数: %v, Path参数：%v, Query参数: %v, Body参数: %v, Resp参数: %v;\n", apiDef.ApiId, apiDef.ApiModule, apiDef.ApiDesc, apiDef.HttpMethod, apiDef.Path, apiDef.Header, apiDef.Path, apiDef.QueryParameter, apiDef.Body, apiDef.Response)
			apiDefStr = apiDefStr + tmpStr
		}

	}
	query := strings.Replace(content, "{需求定义}", apiDefStr, -1)
	go func(query, appendContent string, input InputCase) {
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

		err = input.CommonExtend.SaveAICaseByRepyList(replyList)
		if err != nil {
			Logger.Error("%s", err)
		}

	}(query, appendContent, input)

	return
}

func (input ImportCommon) AICreateCaseByImport() (err error) {
	if len(input.RawReply) > 0 {
		err = input.SaveAICaseByRawRepy()
		if err != nil {
			Logger.Error("%s", err)
		}
	} else {
		replyList, errTmp := input.CommonExtend.ConnectModel2GetMessage(input.ConversationId)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			Logger.Debug("conversationId: %v", input.ConversationId)
			Logger.Debug("replyList: %v", replyList)
			return errTmp
		}
		errTmp = input.CommonExtend.SaveAICaseByRepyList(replyList)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			return errTmp
		}
	}

	return
}

func (input InputCase) AICreateCaseByCreateDesc(createDesc, uploadFilePath string) (err error) {
	content, appendContent, err := GetAiTemplateByName(input.AiTemplate, "1")
	if err != nil {
		return
	}

	query := strings.Replace(content, "{需求定义}", createDesc, -1)
	go func(query, appendContent, uploadFilePath string, input InputCase) {
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
		err = input.CommonExtend.SaveAICaseByRepyList(replyList)
		if err != nil {
			Logger.Error("%s", err)
		}
	}(query, appendContent, uploadFilePath, input)

	return
}

func AIOptimizeCase(ids, optimizeDesc, aiPlatform, createUser string) (err error) {
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

func UseAiCase(ids, userName string) (err error) {
	idList := strings.Split(ids, ",")

	var aiCaseList []AITestCase
	models.Orm.Table("ai_case").Where("id in (?)", idList).Find(&aiCaseList)
	if len(aiCaseList) == 0 {
		err = fmt.Errorf("未找到相关的用例，请核对~")
		Logger.Error("%s", err)
		return
	}

	for _, item := range aiCaseList {
		var tmpCase TestCase
		models.Orm.Table("test_case").Where("product = ? and  case_name= ? and case_number = ?", item.Product, item.CaseName, item.CaseNumber).Find(&tmpCase)
		tmpCase.CaseNumber = item.CaseNumber
		tmpCase.CaseType = item.CaseType
		tmpCase.Priority = item.Priority
		tmpCase.Auto = item.Auto
		tmpCase.Module = item.Module
		tmpCase.IntroVersion = item.IntroVersion
		tmpCase.PreCondition = item.PreCondition
		tmpCase.TestRange = item.TestRange
		tmpCase.TestSteps = item.TestSteps
		tmpCase.ExpectResult = item.ExpectResult
		tmpCase.CaseDesigner = userName
		tmpCase.Product = item.Product
		//curTime := time.Now()
		//testcase.UpdatedAt = curTime.Format(baseFormat)
		if len(tmpCase.CaseName) == 0 {
			tmpCase.CaseName = item.CaseName
			err = models.Orm.Table("test_case").Create(tmpCase).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		} else {
			tmpCase.CaseName = item.CaseName
			err = models.Orm.Table("test_case").Where("product = ? and  case_name= ? and case_number = ?", item.Product, item.CaseName, item.CaseNumber).Update(tmpCase).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		}
	}

	// 数据落库成功后再设置取用状态
	err = models.Orm.Table("ai_case").Where("id in (?)", idList).UpdateColumn(&StatusMgmt{UseStatus: 2}).Error
	if err != nil {
		Logger.Error("%v", err)
		return
	}

	return
}

func UpdateAiCaseStatus(id string) (err error) {
	err = models.Orm.Table("ai_case").Where("id = ?", id).UpdateColumn(&StatusMgmt{ModifyStatus: 2}).Error
	if err != nil {
		Logger.Error("%v", err)
		return
	}

	return
}

func GetDbAiCaseById(id string) (aiCases []DbAiCase) {
	ids := strings.Split(id, ",")
	models.Orm.Table("ai_case").Where("id in (?)", ids).Find(&aiCases)
	return
}
