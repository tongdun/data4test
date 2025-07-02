package biz

import (
	"data4test/models"
	"fmt"
	"strings"
	"time"
)

func ExportTestCase2Markdown(ids string) (fileName string, err error) {
	idList := strings.Split(ids, ",")
	var productNames, modules []string
	models.Orm.Table("test_case").Where("id in (?)", idList).Group("product").Pluck("product", &productNames)
	if len(productNames) == 0 || len(productNames) > 1 {
		errTmp := fmt.Errorf("关联产品数存在: %d, 不支持导出", len(productNames))
		Logger.Error("%s", errTmp)
		return "", errTmp
	}

	timeFormat := "20060102150405"
	curTime := time.Now().Format(timeFormat)
	fileName = fmt.Sprintf("%s_%s.md", productNames[0], curTime)
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	fistTitle := fmt.Sprintf("# %s", productNames[0])
	WriteDataInCommonFile(filePath, fistTitle)
	models.Orm.Table("test_case").Where("id in (?)", idList).Group("module").Pluck("module", &modules)
	for _, item := range modules {
		var testCases []TestCase
		models.Orm.Table("test_case").Where(" module = ? and id in (?)", item, idList).Find(&testCases)
		if len(testCases) > 0 {
			secondTitle := fmt.Sprintf("## %s", item)
			WriteDataInCommonFile(filePath, secondTitle)
			for _, itemCase := range testCases {
				thirdTitle := fmt.Sprintf("### %s", itemCase.CaseName)
				WriteDataInCommonFile(filePath, thirdTitle)
				WriteDataInCommonFile(filePath, itemCase.PreCondition)

				var testSteps, expectResults []string
				if strings.Contains(itemCase.TestSteps, ";") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", "", -1), ";")
				} else if strings.Contains(itemCase.TestSteps, "；") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", "", -1), "；")
				}

				if strings.Contains(itemCase.ExpectResult, ";") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", "", -1), ";")
				} else if strings.Contains(itemCase.ExpectResult, "；") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", "", -1), "；")
				}

				for index, itemStep := range testSteps {
					fourthTitle := fmt.Sprintf("- %s", itemStep)
					WriteDataInCommonFile(filePath, fourthTitle)
					if index < len(expectResults) {
						fifthContent := fmt.Sprintf("    - %s", expectResults[index])
						WriteDataInCommonFile(filePath, fifthContent)
					}
				}

				if len(expectResults) > len(testSteps) {
					for _, itemExpect := range expectResults[len(testSteps):] {
						fifthContent := fmt.Sprintf("    - %s", itemExpect)
						WriteDataInCommonFile(filePath, fifthContent)
					}
				}
			}
		}
	}

	return
}
