package biz

import (
	"data4test/models"
	"fmt"
	"github.com/vivian0517/goxmind"
	"strings"
	"time"
)

func ExportTestCase2Markdown(ids, source string) (fileName string, err error) {
	idList := strings.Split(ids, ",")
	var productNames, modules []string
	models.Orm.Table(source).Where("id in (?)", idList).Group("product").Pluck("product", &productNames)
	if len(productNames) == 0 {
		errTmp := fmt.Errorf("关联产品数为: %d, 不支持导出", len(productNames))
		return "", errTmp
	} else if len(productNames) > 1 {
		Logger.Warning("关联产品数为: %d, 默认导出第一个作为产品", len(productNames))
	}

	timeFormat := "20060102150405"
	curTime := time.Now().Format(timeFormat)
	fileName = fmt.Sprintf("%s_%s.md", productNames[0], curTime)
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	fistTitle := fmt.Sprintf("# %s", productNames[0])
	WriteDataInCommonFile(filePath, fistTitle)
	models.Orm.Table(source).Where("id in (?)", idList).Group("module").Pluck("module", &modules)

	for _, module := range modules {
		var testCases []TestCase
		models.Orm.Table(source).Where(" module = ? and id in (?)", module, idList).Find(&testCases)
		if len(testCases) > 0 {
			secondTitle := fmt.Sprintf("## %s", module)
			WriteDataInCommonFile(filePath, secondTitle)
			for _, itemCase := range testCases {
				thirdTitle := fmt.Sprintf("### %s", itemCase.CaseName)
				WriteDataInCommonFile(filePath, thirdTitle)
				WriteDataInCommonFile(filePath, fmt.Sprintf("#### 优先级: %s", itemCase.Priority))
				WriteDataInCommonFile(filePath, fmt.Sprintf("#### 前置条件: %s", itemCase.PreCondition))
				WriteDataInCommonFile(filePath, fmt.Sprintf("#### 测试范围: %s", itemCase.TestRange))
				WriteDataInCommonFile(filePath, fmt.Sprintf("#### 测试步骤"))

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
					fourthTitle := fmt.Sprintf("##### %s", itemStep)
					WriteDataInCommonFile(filePath, fourthTitle)
					if index < len(expectResults) {
						fifthContent := fmt.Sprintf("- %s", expectResults[index])
						WriteDataInCommonFile(filePath, fifthContent)
					}
				}

				if len(expectResults) > len(testSteps) {
					for _, itemExpect := range expectResults[len(testSteps):] {
						fifthContent := fmt.Sprintf("- %s", itemExpect)
						WriteDataInCommonFile(filePath, fifthContent)
					}
				}
			}
		}
	}

	return
}

func ExportTestCase2Xmind(ids, source string) (fileName string, err error) {
	idList := strings.Split(ids, ",")
	var productNames, modules []string
	models.Orm.Table(source).Where("id in (?)", idList).Group("product").Pluck("product", &productNames)
	if len(productNames) == 0 {
		errTmp := fmt.Errorf("关联产品数为: %d, 不支持导出", len(productNames))
		return "", errTmp
	} else if len(productNames) > 1 {
		Logger.Warning("关联产品数为: %d, 默认导出第一个作为产品", len(productNames))
	}

	timeFormat := "20060102150405"
	curTime := time.Now().Format(timeFormat)
	fileName = fmt.Sprintf("%s_%s.xmind", productNames[0], curTime)
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	xmind := goxmind.New()
	firstNode := xmind.AddSheet("Sheet title", productNames[0])
	models.Orm.Table(source).Where("id in (?)", idList).Group("module").Pluck("module", &modules)

	for _, module := range modules {
		var testCases []TestCase
		models.Orm.Table(source).Where(" module = ? and id in (?)", module, idList).Find(&testCases)
		if len(testCases) > 0 {
			secondNode := firstNode.AddNode(module)
			for _, itemCase := range testCases {
				thirdNode := secondNode.AddNode(itemCase.CaseName)
				var priorityMark goxmind.MarkerId
				switch itemCase.Priority {
				case "P0":
					priorityMark = goxmind.Priority1
				case "P1", "高":
					priorityMark = goxmind.Priority1
				case "P2", "中":
					priorityMark = goxmind.Priority2
				case "P3", "低":
					priorityMark = goxmind.Priority3
				case "P4":
					priorityMark = goxmind.Priority4
				default:
					priorityMark = goxmind.Priority2
				}
				thirdNode.AddMaker(priorityMark)
				thirdNode.AddNode(fmt.Sprintf("前置条件: %s", itemCase.PreCondition))
				thirdNode.AddNode(fmt.Sprintf("测试范围: %s", itemCase.TestRange))
				if itemCase.Auto == "1" {
					thirdNode.AddNode(fmt.Sprintf("是否支持自动化: 是"))
				} else if itemCase.Auto == "0" {
					thirdNode.AddNode(fmt.Sprintf("是否支持自动化: 否"))
				}

				fourthNode := thirdNode.AddNode(fmt.Sprintf("测试步骤"))
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

				var lastNode *goxmind.Node
				for index, itemStep := range testSteps {
					fifthNode := fourthNode.AddNode(itemStep)
					if len(expectResults) > index {
						fifthNode.AddNode(expectResults[index])
					}
					if index == len(testSteps)-1 {
						lastNode = fifthNode
					}

				}

				if len(expectResults) > len(testSteps) {
					for _, itemExpect := range expectResults[len(testSteps):] {
						lastNode.AddNode(itemExpect)
					}
				}

			}
		}
	}

	xmind.Save(filePath)

	return
}
