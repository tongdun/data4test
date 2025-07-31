package biz

import (
	"bytes"
	"data4test/models"
	"encoding/json"
	"fmt"
	"github.com/GoAdminGroup/go-admin/template/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func GetKnowledgeType() (knowTypes []types.FieldOption) {
	var kType types.FieldOption
	chTypes := []string{"场景", "数据", "任务", "用例", "接口", "文档", "全部"}
	enTypes := []string{"playbook", "data", "task", "testcase", "api", "doc", "all"}
	for index, item := range chTypes {
		kType.Value = enTypes[index]
		kType.Text = item
		knowTypes = append(knowTypes, kType)
	}

	return
}

func EmptyFileContent(filePath string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func WritePlaybookKnowledge() {
	fileName := "playbook_from_history.txt"
	filePath := fmt.Sprintf("%s/%s", KnowledgeBasePath, fileName)
	var playbookList []Scene
	var playbookNameList []string
	models.Orm.Table("scene_test_history").Where("result = ?", "pass").Group("name").Pluck("name", &playbookNameList)
	if len(playbookNameList) == 0 {
		return
	}
	models.Orm.Table("playbook").Where("name in (?)", playbookNameList).Find(&playbookList)
	if len(playbookList) == 0 {
		return
	}
	EmptyFileContent(filePath)
	var kPlaybook KPlaybook
	for _, playbook := range playbookList {
		kPlaybook.Name = playbook.Name
		kPlaybook.DataList = strings.Split(playbook.DataFileList, ",")
		kByte, _ := json.MarshalIndent(kPlaybook, "", "    ")
		WriteDataInCommonFile(filePath, string(kByte))
	}

	return
}

func WriteDataKnowledge() {
	fileName := "data_from_history.txt"
	filePath := fmt.Sprintf("%s/%s", KnowledgeBasePath, fileName)
	var dataList []DataBase
	var dataNameList []string
	models.Orm.Table("scene_data_test_history").Where("result = ?", "pass").Group("name").Pluck("name", &dataNameList)
	if len(dataNameList) == 0 {
		return
	}
	models.Orm.Table("scene_data").Where("name in (?)", dataNameList).Find(&dataList)
	if len(dataNameList) == 0 {
		return
	}

	EmptyFileContent(filePath)

	for _, data := range dataList {
		remark := fmt.Sprintf("// %s", data.FileName)
		WriteDataInCommonFile(filePath, remark)
		WriteDataInCommonFile(filePath, data.Content)
	}

	return
}

func WriteTaskKnowledge() {
	fileName := "task.txt"
	filePath := fmt.Sprintf("%s/%s", KnowledgeBasePath, fileName)
	var taskList []Schedule
	models.Orm.Table("schedule").Where("task_status = ? and task_type = ?", "finished", "scene").Group("task_name").Find(&taskList)
	if len(taskList) == 0 {
		return
	}

	EmptyFileContent(filePath)

	var kTask KTask
	for _, task := range taskList {
		kTask.Name = task.TaskName
		kTask.PlaybookList = strings.Split(task.SceneList, ",")
		kByte, _ := json.MarshalIndent(kTask, "", "    ")
		WriteDataInCommonFile(filePath, string(kByte))
	}

	return
}

func WriteTestCaseKnowledge() {
	fileName := "testcase.txt"
	filePath := fmt.Sprintf("%s/%s", KnowledgeBasePath, fileName)
	var caseList []TestCase
	// 待有数据再变更策略
	//models.Orm.Table("test_case").Where("test_result in (?)", "finished", "scene").Group("case_number").Find(&caseList)
	models.Orm.Table("test_case").Find(&caseList)
	if len(caseList) == 0 {
		return
	}

	EmptyFileContent(filePath)

	var kCase KCase
	for _, testCase := range caseList {
		kCase.CaseNumber = testCase.CaseNumber
		kCase.CaseName = testCase.CaseName
		kCase.TestSteps = testCase.TestSteps
		kCase.TestRange = testCase.TestRange
		kCase.Module = testCase.Module
		kCase.PreCondition = testCase.PreCondition
		kCase.ExpectResult = testCase.ExpectResult
		kCase.Priority = testCase.Priority
		kCase.Auto = testCase.Auto
		kByte, _ := json.MarshalIndent(kCase, "", "    ")
		WriteDataInCommonFile(filePath, string(kByte))
	}

	return
}

func GetExistAppDefineFile() (appExistList []string) {
	var appList []string
	models.Orm.Table("env_config").Pluck("app", &appList)
	if len(appList) == 0 {
		return
	}

	for _, appName := range appList {
		fileName := fmt.Sprintf("%s.txt", appName)
		filePath := fmt.Sprintf("%s/%s", KnowledgeBasePath, fileName)
		_, err := os.Stat(filePath)
		if err == nil {
			appExistList = append(appExistList, filePath)
		}

	}

	return
}

func WriteApiDefineKnowledge() {
	var appList []EnvConfig
	models.Orm.Table("env_config").Find(&appList)
	if len(appList) == 0 {
		return
	}

	for _, app := range appList {
		fileName := fmt.Sprintf("%s.txt", app.App)
		filePath := fmt.Sprintf("%s/%s", KnowledgeBasePath, fileName)
		var content []byte
		if len(strings.TrimSpace(app.SwaggerPath)) > 0 {
			resp, errTemp := http.Get(app.SwaggerPath)
			// 一次性读取
			if errTemp == nil {
				content, errTemp = ioutil.ReadAll(resp.Body)
				if errTemp != nil {
					Logger.Error("%s", errTemp)
					content = []byte{}
				}
				resp.Body.Close()
			}
		}

		if len(content) == 0 {
			cmdStr := fmt.Sprintf("ls -lt %s/%s*.json | head -n 1 | awk '{print $9}'", ApiFilePath, app.App)
			fileName, errTmp := ExecCommand(cmdStr)
			if errTmp != nil {
				Logger.Error("%s", errTmp)
			} else {
				content, _ = ioutil.ReadFile(fileName)
				if errTmp != nil {
					Logger.Error("%s", errTmp)
					content = []byte{}
				}
			}
		}

		if len(content) == 0 {
			Logger.Warning("未获取到[%s]的接口信息，请核对 ~", app.App)
			continue // 没有获取到数据，继续下一个应用的接口获取
		}

		EmptyFileContent(filePath)

		var strContent bytes.Buffer

		json.Indent(&strContent, []byte(content), "", "    ")

		WriteDataInCommonFile(filePath, strContent.String())

	}

	return
}

func GetAllKnowledge(kType, syncUser string) (err error) {
	if syncUser != "admin" {
		err = fmt.Errorf("当前用户无[知识库同步]权限, 请联系管理员处理 ~")
		return
	}

	switch kType {
	case "场景", "playbook":
		WritePlaybookKnowledge()
	case "任务", "task":
		WriteTaskKnowledge()
	case "数据", "data":
		WriteDataKnowledge()
	case "接口", "api":
		WriteApiDefineKnowledge()
	case "用例", "testcase":
		WriteTestCaseKnowledge()
	case "文档", "doc":
		return
	case "全部", "all":
		WritePlaybookKnowledge()
		WriteDataKnowledge()
		WriteTaskKnowledge()
		WriteTestCaseKnowledge()
		WriteApiDefineKnowledge()
	default:
		err = fmt.Errorf("知识库类型[%s]未知，请核对~ ", kType)
	}

	return
}

func GetAIRAGConnectInfo() (aiConnect DataSetConnect, err error) {
	var sysParameter SysParameter
	parameterName := "aiRAGEngine"
	models.Orm.Table("sys_parameter").Where("name = ?", parameterName).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		err = fmt.Errorf("系统参数中未定义[%s]参数的值，请核对~", parameterName)
		Logger.Error("%s", err)
		return
	}

	json.Unmarshal([]byte(sysParameter.ValueList), &aiConnect)

	if len(aiConnect.BaseUrl) == 0 || len(aiConnect.ApiKey) == 0 {
		err = fmt.Errorf("[%s]参数定义连接信息不全: %v，请核对~", parameterName, aiConnect)
		Logger.Error("%s", err)
	}

	return
}

func GetDirFilePathList(dirName string) (filePathList []string) {
	cmdStr := fmt.Sprintf("ls -lt %s | awk '{print $9}'", dirName)
	output, errTmp := ExecCommandRaw(cmdStr)
	if errTmp != nil {
		Logger.Debug("cmdStr: %s", cmdStr)
		Logger.Error("%s", errTmp)
		return
	}
	fileNameList := strings.Split(output, ",")
	for _, item := range fileNameList {
		if len(strings.TrimSpace(item)) == 0 {
			continue
		}
		filePath := fmt.Sprintf("%s/%s", dirName, item)
		filePathList = append(filePathList, filePath)
	}

	return
}

func GetDocFilePathList() (fileDocPathList []string) {
	fileDocPathList = []string{
		fmt.Sprintf("%s/ai/ai_analysis_design.md", DocFilePath),
		fmt.Sprintf("%s/ai/ai_data_design.md", DocFilePath),
		fmt.Sprintf("%s/ai/ai_playbook_design.md", DocFilePath),
		fmt.Sprintf("%s/ai/ai_template_design.md", DocFilePath),
		fmt.Sprintf("%s/ai/ai_testcase_design.md", DocFilePath),
		fmt.Sprintf("%s/arch/arch.md", DocFilePath),
		fmt.Sprintf("%s/design/action_design.md", DocFilePath),
		fmt.Sprintf("%s/design/api_mgmt_design.md", DocFilePath),
		fmt.Sprintf("%s/design/assert_design.md", DocFilePath),
		fmt.Sprintf("%s/design/console_design.md", DocFilePath),
		fmt.Sprintf("%s/design/data_file_design.md", DocFilePath),
		fmt.Sprintf("%s/design/mock_design.md", DocFilePath),
		fmt.Sprintf("%s/design/parameter_design.md", DocFilePath),
		fmt.Sprintf("%s/design/playbook_design.md", DocFilePath),
		fmt.Sprintf("%s/design/script_design.md", DocFilePath),
		fmt.Sprintf("%s/design/task_design.md", DocFilePath),
		fmt.Sprintf("%s/design/testcase_design.md", DocFilePath),
		fmt.Sprintf("%s/function/feature_introduction.md", DocFilePath),
		fmt.Sprintf("%s/function/module_function.md", DocFilePath),
		fmt.Sprintf("%s/question/FAQ.md", DocFilePath),
	}
	return
}

func UpdateAssetKnowledge(kType, syncUser string) (err error) {

	err = GetAllKnowledge(kType, syncUser)
	if err != nil {
		Logger.Error("err: %v", err)
		return
	}

	var filePathList []string

	switch kType {
	case "场景", "playbook":
		filePathList = []string{
			fmt.Sprintf("%s/playbook_from_history.txt", KnowledgeBasePath),
		}
	case "任务", "task":
		filePathList = []string{
			fmt.Sprintf("%s/task.txt", KnowledgeBasePath),
		}
	case "数据", "data":
		filePathList = []string{
			fmt.Sprintf("%s/task_from_history.txt", KnowledgeBasePath),
		}
	case "接口", "api":
		filePathList = GetExistAppDefineFile()
	case "用例", "testcase":
		filePathList = []string{
			fmt.Sprintf("%s/testcase.txt", KnowledgeBasePath),
		}
	case "文档", "doc":
		filePathList = GetDocFilePathList()
	case "全部", "all":
		filePathList = GetDirFilePathList(KnowledgeBasePath)
		fileDocPathList := GetDocFilePathList()
		filePathList = append(filePathList, fileDocPathList...)
	default:
		err = fmt.Errorf("知识库类型[%s]未知，请核对~ ", kType)
	}

	aiConnect, err := GetAIRAGConnectInfo()
	if err != nil {
		return
	}

	var kName string
	if len(aiConnect.DataSetName) == 0 {
		kName = "Data4Test知识库"
	} else {
		kName = aiConnect.DataSetName
	}

	dataSet, err := aiConnect.GetDatasetByName(kName)
	if err != nil {
		Logger.Error("%s", err)

	}

	var dataSetId string
	if len(dataSet.Name) == 0 {
		dataSetId, err = aiConnect.CreateDataset(kName)
		if err != nil {
			Logger.Error("%s", err)
			return
		}
	} else {
		dataSetId = dataSet.ID
	}

	if len(dataSetId) == 0 {
		return
	}

	for _, filePath := range filePathList {
		content, errTmp := ioutil.ReadFile(filePath)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			err = errTmp
			return
		}

		doc, errTmp := aiConnect.GetDocumentByName(dataSetId, path.Base(filePath))
		if errTmp == nil {
			aiConnect.UpdateDocument(dataSetId, doc.ID, path.Base(filePath), string(content))
		} else {
			aiConnect.CreateDocument(dataSetId, path.Base(filePath), string(content))
		}
	}

	return
}
