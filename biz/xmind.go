package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetFileName(dirName, product string) (fileName string, err error) {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		Logger.Error("%s", err)
	}
	var allNames []string
	for _, file := range files {
		tmpName := file.Name()
		if strings.HasPrefix(tmpName, product) && strings.HasSuffix(tmpName, ".xmind") {
			rawName := dirName + "/" + tmpName
			allNames = append(allNames, rawName)
		}
	}

	if len(allNames) > 0 {
		fileName = allNames[len(allNames)-1]
	}

	if len(fileName) == 0 {
		err = fmt.Errorf("Not Found file [%s*.xmind] in directory[%s]", product, dirName)
		Logger.Error("%s", err)
	}
	return
}

func GetJSON(id string) (err error) {
	StatusDef := map[int]string{1: "草稿", 2: "待评审", 3: "评审中", 4: "重做", 5: "废弃", 6: "特性", 7: "终稿"}
	PriorityDef := map[int]string{1: "P1", 2: "P2", 3: "P3"}
	AutoDef := map[int]string{1: "是", 2: "否"}
	TestTypeDef := map[int]string{1: "冒烟", 2: "场景", 3: "异常"}
	var product Product
	models.Orm.Table("product").Where("id = ?", id).Find(&product)
	if len(product.Name) == 0 {
		err = fmt.Errorf("Not found related product id:%s", id)
		return
	}
	productName := product.Name
	fileName, err := GetFileName(CaseFilePath, productName)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	var ms string
	if strings.Contains(fileName, "_v") || strings.Contains(fileName, "_V") {
		items := strings.Split(fileName, "_")
		tmpName := items[1]
		ms = strings.Replace(tmpName, ".xmind", "", -1)
	}

	output, err := exec.Command("xmind2case", fileName, "-json").Output()
	if err != nil {
		Logger.Debug("output: %s", output)
		Logger.Error("%s", err)
	}

	jsonFileName := fileName[:len(fileName)-len(".xmind")] + ".json"
	content, err := ioutil.ReadFile(jsonFileName)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	var xmindTestCases []XmindTestCase
	err = json.Unmarshal([]byte(content), &xmindTestCases)
	if err != nil {
		Logger.Error("%s", err)
	}

	var caseNumberPrefix, sufixNum string
	for _, item := range xmindTestCases {
		var testcase, testCaseDb, testCaseDb2 TestCase
		modules := strings.Split(item.Suite, "-")
		testcase.CaseName = item.Name
		testcase.Auto = AutoDef[item.ExecutionType]
		testcase.TestResult = StatusDef[item.Status]

		if strings.Contains(testcase.CaseName, ":") {
			tmps := strings.Split(testcase.CaseName, ":")
			testcase.CaseType = tmps[0]
		} else if strings.Contains(testcase.CaseName, "：") {
			tmps := strings.Split(testcase.CaseName, "：")
			testcase.CaseType = tmps[0]
		} else {
			testcase.CaseType = TestTypeDef[item.Importance]
		}

		testcase.Priority = PriorityDef[item.Importance]
		testcase.PreCondition = item.Preconditions
		testcase.Product = productName
		var stepStr, resultStr string
		for _, step := range item.Steps {
			tmpAction := fmt.Sprintf("%d. %s\n", step.StepNumber, step.Actions)
			tmpResult := fmt.Sprintf("%d. %s\n", step.StepNumber, step.ExpectedResults)
			stepStr = stepStr + tmpAction
			resultStr = resultStr + tmpResult
		}

		testcase.TestSteps = stepStr
		testcase.ExpectResult = resultStr
		if len(modules) > 1 {
			testcase.Module = modules[0]
			if len(ms) > 0 {
				caseNumberPrefix = modules[1] + "_" + ms
			} else {
				caseNumberPrefix = modules[1]
			}

		} else {
			testcase.CaseName = item.Name
			if len(ms) > 0 {
				caseNumberPrefix = item.Product + "_" + ms
			} else {
				caseNumberPrefix = item.Product + "_" + "other"
			}

		}

		chkStr := "%" + caseNumberPrefix + "%"
		models.Orm.Table("test_case").Where("product = ? AND case_number LIKE ?", productName, chkStr).Find(&testCaseDb)
		if len(testCaseDb.CaseNumber) == 0 {
			sufixNum = strconv.Itoa(1)
		} else {
			tmps := strings.Split(testCaseDb.CaseNumber, "_")
			numStr := tmps[len(tmps)-1]
			s, err := strconv.Atoi(numStr)
			if err != nil {
				Logger.Error("%s", err)
			}
			sufixNum = strconv.Itoa(s + 1)
		}

		testcase.CaseNumber = caseNumberPrefix + "_" + sufixNum

		curTime := time.Now()
		testcase.UpdatedAt = curTime.Format(baseFormat)
		models.Orm.Table("test_case").Where("product = ? and case_number = ?", productName, testcase.CaseNumber).Find(&testCaseDb2)
		if len(testCaseDb2.CaseNumber) == 0 {
			err = models.Orm.Table("test_case").Create(testcase).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		} else {
			err = models.Orm.Table("test_case").Where("product = ? and case_number = ?", productName, testcase.CaseNumber).Update(testcase).Error
			if err != nil {
				Logger.Error("%s", err)
			}
		}
	}

	return

}
