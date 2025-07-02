package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func GetAnalysisFromReplyList(replyList []string) (issueList, testResultList, asserList []map[string]interface{}, err error) {
	for _, reply := range replyList {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(reply))
		if err != nil {
			Logger.Error("%s", err)
		}

		afterTxt := doc.Text()
		if len(afterTxt) == 0 {
			Logger.Warning("未找到有效信息，请核对~")
		}
		Logger.Debug("afterTxt: %v", afterTxt)

		issueSubList, testResultSubList, asserSubList, _ := GetAnalysisFromReply(reply)
		issueList = append(issueList, issueSubList...)
		testResultList = append(testResultList, testResultSubList...)
		asserList = append(asserList, asserSubList...)
	}

	return
}

func GetAnalysisFromReply(reply string) (issueList, testResultList, assertList []map[string]interface{}, err error) {
	issueList, errTmp1 := GetIssueFromReply(reply)
	if errTmp1 != nil {
		err = errTmp1
	}
	testResultList, errTmp2 := GetTestResultFromReply(reply)
	if errTmp2 != nil {
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp2)
		} else {
			err = errTmp2
		}
	}
	assertList, errTmp3 := GetAssertFromReply(reply)
	if errTmp3 != nil {
		if err != nil {
			err = fmt.Errorf("%s; %s", err, errTmp3)
		} else {
			err = errTmp3
		}
	}
	return
}

func GetIssueFromReply(reply string) (issueList []map[string]interface{}, err error) {
	issueReg := regexp.MustCompile(`"分析问题":[\s\S](\[[\s\n]*\{[\s\n]*[\s\S]*?[^{\w\S}]\}.?[\s\n]*\])`)

	// 查找所有匹配项
	issueMatch := issueReg.FindAllStringSubmatch(reply, -1)
	var targetStr string
	if len(issueMatch) > 0 {
		if len(issueMatch[0]) > 1 {
			targetStr = issueMatch[0][1]
		}

		if len(issueMatch[0]) > 2 {
			Logger.Warning("匹配到了多笔数据，请核对 ~")
			Logger.Debug("issueMatch: %s", issueMatch[0])
		}
	}

	err = json.Unmarshal([]byte(targetStr), &issueList)
	if err != nil {
		Logger.Debug("reply: %s", reply)
		Logger.Debug("matchIssue: %v", targetStr)
		Logger.Error("%s", err)
	}

	return
}

func GetTestResultFromReply(reply string) (testResultList []map[string]interface{}, err error) {
	testResultReg := regexp.MustCompile(`"判断结果":[\s\S](\[[\s\n]*\{[\s\n]*[\s\S]*?[^{\w\S}]\}.?[\s\n]*\])`)

	// 查找所有匹配项
	testResultMatch := testResultReg.FindAllStringSubmatch(reply, -1)

	var targetStr string
	if len(testResultMatch) > 0 {
		if len(testResultMatch[0]) > 1 {
			targetStr = testResultMatch[0][1]
		}
		if len(testResultMatch[0]) > 2 {
			Logger.Warning("匹配到了多笔数据，请核对 ~")
			Logger.Debug("testResultMatch: %s", testResultMatch[0])
		}
	}

	err = json.Unmarshal([]byte(targetStr), &testResultList)
	if err != nil {
		Logger.Debug("reply: %s", reply)
		Logger.Debug("matchTestResult: %v", targetStr)
		Logger.Error("%s", err)
	}

	return
}

func GetAssertFromReply(reply string) (assertList []map[string]interface{}, err error) {
	assertReg := regexp.MustCompile(`"断言设置":[\s\S](\[[\s\n]*\{[\s\n]*[\s\S]*?[^{\w\S}]\}.?[\s\n]*\])`)

	// 查找所有匹配项
	assertMatch := assertReg.FindAllStringSubmatch(reply, -1)

	var targetStr string
	if len(assertMatch) > 0 {
		if len(assertMatch[0]) > 1 {
			targetStr = assertMatch[0][1]
		}
		if len(assertMatch[0]) > 2 {
			Logger.Warning("匹配到了多笔数据，请核对 ~")
			Logger.Debug("assertMatch: %s", assertMatch[0])
		}
	}

	err = json.Unmarshal([]byte(targetStr), &assertList)
	if err != nil {
		Logger.Debug("reply: %s", reply)
		Logger.Debug("matchAssert: %v", targetStr)
		Logger.Error("%s", err)
	}

	return
}

func (input AnalysisDataInput) SaveAIIssueByRepyList(df DataFile, replyList []string) (testResult, failReason []string, assert []SceneAssert, err error) {
	issueList, testResultRaw, assertRaw, err := GetAnalysisFromReplyList(replyList)
	if err != nil {
		Logger.Error("%s", err)
	}

	for _, item := range issueList {
		aiIssue := input.AssembleAIIssue(df, item)
		errTmp := aiIssue.AddAiIssue()
		if errTmp != nil {
			if err != nil {
				err = fmt.Errorf("%s;%s", err, errTmp)
			} else {
				err = errTmp
			}
		}
	}

	for _, item := range testResultRaw {
		testResult = append(testResult, item["测试结果"].(string))
		vType := fmt.Sprintf("%T", item["失败原因"])
		if vType == "string" {
			failReason = append(failReason, item["失败原因"].(string))
		} else if vType == "[]interface {}" {
			tmpList := item["失败原因"].([]interface{})
			var tmpStr string
			for index, item := range tmpList {
				if index == 0 {
					tmpStr = item.(string)
				} else {
					tmpStr = fmt.Sprintf("%s,%s", tmpStr, item.(string))
				}
			}
			failReason = append(failReason, tmpStr)
		} else {
			Logger.Warning("值类型为: %s, 不支持", vType)
		}
	}

	for _, item := range assertRaw {
		var assertTmp SceneAssert
		assertTmp.Type = item["断言类型"].(string)
		assertTmp.Source = item["返回数据字段定位"].(string)
		assertTmp.Value = item["断言值"]
		assert = append(assert, assertTmp)
	}

	return
}

func (input AnalysisDataInput) AssembleAIIssue(df DataFile, issue map[string]interface{}) (aiIssue AIIssue) {
	rKeyList := []string{"问题名称", "问题级别", "问题类型", "问题详情", "问题原因推测", "影响范围分析", "受影响的场景推测", "受影响的数据推测"}
	tKeyList := []string{"issue_name", "issue_level", "issue_type", "issue_detail", "root_cause", "impact_scope_analysis", "impact_playbook", "impact_data"}

	aiTmpIssue := make(map[string]string)
	for index, k := range rKeyList {
		if v, ok := issue[k]; ok {
			vType := fmt.Sprintf("%T", v)
			if vType == "string" {
				aiTmpIssue[tKeyList[index]] = v.(string)
			} else if vType == "[]interface {}" {
				tmpList := v.([]interface{})
				var tmpStr string
				for index, item := range tmpList {
					if index == 0 {
						tmpStr = item.(string)
					} else {
						tmpStr = fmt.Sprintf("%s,%s", tmpStr, item.(string))
					}
				}
				aiTmpIssue[tKeyList[index]] = tmpStr
			} else {
				Logger.Warning("值类型为: %s, 不支持", vType)
			}

		}
	}

	aiIssue.IssueName = aiTmpIssue["issue_name"]
	aiIssue.IssueLevel = aiTmpIssue["issue_level"]
	aiIssue.IssueSource = "1" //  1: 数据
	aiIssue.SourceName = fmt.Sprintf("%s.yml", df.Name)
	if len(df.Request) > 0 {
		aiIssue.RequestData = df.Request[0]
	} else {
		aiIssue.RequestData = ""
	}
	if len(df.Response) > 0 {
		aiIssue.ResponseData = df.Response[0]
	} else {
		aiIssue.ResponseData = ""
	}

	aiIssue.IssueDetail = aiTmpIssue["issue_detail"]
	if aiTmpIssue["issue_type"] == "BUG" {
		aiIssue.ConfirmStatus = "1"
	} else {
		aiIssue.ConfirmStatus = "2"
	}

	aiIssue.RootCause = aiTmpIssue["root_cause"]
	aiIssue.ImpactScopeAnalysis = aiTmpIssue["impact_scope_analysis"]
	aiIssue.ImpactPlaybook = aiTmpIssue["impact_playbook"]
	aiIssue.ImpactData = aiTmpIssue["impact_data"]
	aiIssue.ResolutionStatus = "1" // 默认创建
	aiIssue.AgainTestResult = "2"  // 默认未知
	aiIssue.ImpactTestResult = "2" // 默认未知
	aiIssue.CreateUser = input.CreateUser
	aiIssue.ProductList = input.Product
	aiIssue.CreatePlatform = input.CreatePlatform

	return
}

func (aiIssue AIIssue) AddAiIssue() (err error) {
	var tmpIssue AIIssue
	models.Orm.Table("ai_issue").Where("issue_name = ? and  source_name= ?", aiIssue.IssueName, aiIssue.SourceName).Find(&tmpIssue)
	if len(tmpIssue.IssueName) > 0 {
		err = models.Orm.Table("ai_issue").Where("issue_name = ? and  source_name= ?", aiIssue.IssueName, aiIssue.SourceName).Update(aiIssue).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	} else {
		err = models.Orm.Table("ai_issue").Create(&aiIssue).Error

		if err != nil {
			Logger.Error("aiIssue: %s", aiIssue)
			Logger.Error("%s", err)
		}
	}

	return
}

func (input ImportCommon) AIAnalysisDataByImport() (err error) {
	if len(input.RawReply) > 0 {
		err = input.SaveAIIssueByRawRepy()
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
		issueList, _, _, _ := GetAnalysisFromReplyList(replyList)

		errTmp = input.SaveAIIssueByIssueListMap(issueList)
	}

	return
}

func (input ImportCommon) SaveAIIssueByIssueListMap(issueList []map[string]interface{}) (err error) {
	if len(issueList) == 0 {
		return fmt.Errorf("无数据信息，请核对~")
	}
	var inputA AnalysisDataInput
	inputA.Product = input.Product
	inputA.CreateUser = input.CreateUser
	inputA.CreatePlatform = input.CreatePlatform

	var df DataFile
	for _, item := range issueList {
		aiIssue := inputA.AssembleAIIssue(df, item)
		errTmp := aiIssue.AddAiIssue()
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

func (input ImportCommon) SaveAIIssueByRawRepy() (err error) {
	issueList, err := GetIssueFromReply(input.RawReply)
	if err != nil {
		Logger.Error("%s", err)
	}
	var inputA AnalysisDataInput
	inputA.Product = input.Product
	inputA.CreatePlatform = input.CreatePlatform
	inputA.CreateUser = input.CreateUser

	var df DataFile
	for _, item := range issueList {
		aiIssue := inputA.AssembleAIIssue(df, item)
		errTmp := aiIssue.AddAiIssue()
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

func AiDataTestAgain(aiIssue DbAIIssue) (err error) {
	filePath := fmt.Sprintf("%s/%s", AiDataBasePath, aiIssue.SourceName)
	var dbData DbSceneData
	Logger.Debug("item.SourceName: %v", aiIssue.SourceName)
	models.Orm.Table("ai_data").Where("file_name = (?)", aiIssue.SourceName).Find(&dbData)
	if len(dbData.Name) == 0 {
		newErr := fmt.Errorf("在[智能数据]列表未找到数据: %s，请核对", aiIssue.SourceName)
		if err == nil {
			err = newErr
		} else {
			err = fmt.Errorf("%v; %v", err, newErr)
		}
	}
	againTest := "2"
	result, dst, err1 := dbData.RunDataFile(filePath, aiIssue.ProductList, "ai_data", nil)
	if err1 != nil {
		againTest = "0"
		Logger.Error("\n%s", err1)
		err = err1
	}

	envType := GetEnvTypeByName(aiIssue.ProductList)
	err2 := WriteSceneDataResult(dbData.Id, result, dst, aiIssue.ProductList, "ai_data", envType, err1)
	if err2 != nil {
		Logger.Error("%s", err2)
		if err != nil {
			err = fmt.Errorf("%v; %v", err, err2)
		} else {
			err = err2
		}
	}

	if result != "pass" {
		err = fmt.Errorf("test %s", result)
	} else {
		againTest = "1"
		_ = models.Orm.Table("ai_issue").Where("id = (?)", aiIssue.Id).UpdateColumn(&IssueBase{ResolutionStatus: "4"}).Error
	}
	// 回归结果状态更新
	err3 := models.Orm.Table("ai_issue").Where("id = (?)", aiIssue.Id).UpdateColumn(&IssueBase{AgainTestResult: againTest}).Error
	if err3 != nil {
		Logger.Error("%v", err3)
	}

	return
}

func SourceAgainTest(ids string) (err error) {
	idList := strings.Split(ids, ",")
	var aiIssueList []DbAIIssue
	models.Orm.Table("ai_issue").Where("id in (?)", idList).Find(&aiIssueList)

	for _, item := range aiIssueList {
		if item.IssueSource == "1" {
			errTmp := AiDataTestAgain(item)
			if errTmp != nil {
				if err != nil {
					err = fmt.Errorf("%v; %v", err, errTmp)
				} else {
					err = errTmp
				}
			}
		} else if item.IssueSource == "2" {
			// 待开发
		}

	}

	return
}
