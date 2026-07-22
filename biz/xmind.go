package biz

import (
	"archive/zip"
	"data4test/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
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
	StatusDef := map[int]string{1: T("status.draft"), 2: T("status.pending_review"), 3: T("status.in_review"), 4: T("status.redo"), 5: T("status.deprecated"), 6: T("status.feature"), 7: T("status.final")}
	PriorityDef := map[int]string{1: "P1", 2: "P2", 3: "P3"}
	//AutoDef := map[int]string{1: "是", 2: "否"}
	TestTypeDef := map[int]string{1: T("test_type.smoke"), 2: T("test_type.scenario"), 3: T("test_type.exception")}
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
		if item.ExecutionType == 2 {
			testcase.Auto = "0"
		} else {
			testcase.Auto = "1"
		}
		//testcase.Auto = AutoDef[item.ExecutionType]
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

func UpdateTestCaseDB(productName, introVersion string, caseMap map[string]string) (err error) {
	var testcase, testCaseDb TestCase
	models.Orm.Table("test_case").Where("product = ? AND intro_version = ? AND case_name = ?", productName, introVersion, caseMap[T("field.case_name")]).Find(&testCaseDb)
	if len(testCaseDb.CaseNumber) == 0 {
		testcase.TestRange = caseMap[T("field.test_range")]
		testcase.PreCondition = caseMap[T("field.precondition")]
		testcase.CaseType = caseMap[T("field.case_type")]
		testcase.Auto = caseMap[T("field.auto_support")]
		testcase.Module = caseMap[T("field.module")]
		testcase.Product = productName
		testcase.Priority = caseMap[T("field.priority")]
		testcase.ExpectResult = caseMap[T("field.expected_result")]
		testcase.TestSteps = caseMap[T("field.test_steps")]
		testcase.IntroVersion = introVersion
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

	return
}

func GetContentFromXmindFile(filePath string) (xmindContent XMindContent, err error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	defer r.Close()

	var content []byte

	for _, f := range r.File {
		if f.Name == "content.json" {
			rc, errTmp := f.Open()
			if errTmp != nil {
				Logger.Error("%s", err)
				err = errTmp
				return
			}

			defer rc.Close()

			content, errTmp = ioutil.ReadAll(rc)
			if errTmp != nil {
				Logger.Error("%s", err)
				err = errTmp
				return
			}
			break
		}
	}

	if content == nil {
		err = fmt.Errorf(T("error.content_json_not_found"))
		Logger.Error("%s", err)
		return
	}

	err = json.Unmarshal(content, &xmindContent)

	return
}

func Xmind2Excel(filePath string) (fileName string, err error) {
	xmindContent, err := GetContentFromXmindFile(filePath)
	if err != nil {
		return "", err
	}
	timeFormat := "20060102150405"
	curTime := time.Now().Format(timeFormat)
	fileTmp := strings.Split(path.Base(filePath), ".")[0]
	fileName = fmt.Sprintf("%s_%s.xls", fileTmp, curTime)
	fileCasePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	if len(xmindContent) == 0 {
		err = fmt.Errorf(T("error.case_info_not_found"))
		return
	}

	target := xmindContent[0]
	firstTitle := target.RootTopic.Title
	var caseList []map[string]string

	priorityMap := map[string]string{
		"priority-1": "P1",
		"priority-2": "P2",
		"priority-3": "P3",
		"priority-4": "P4",
		"priority-5": "P5",
		"priority-6": "P6",
		"priority-7": "P7",
	}

	testResultMap := map[string]string{
		"c_symbol_like":        "Pass",
		"c_symbol_dislike":     "Fail",
		"c_symbol_music":       "Music",
		"c_symbol_pen":         "Pen",
		"c_symbol_telephone":   "Telephone",
		"c_symbol_hourglass":   "Hourglass",
		"c_symbol_flight":      "Flight",
		"c_symbol_heart":       "Heart",
		"symbol-lightning":     "Lightning",
		"symbol-idea":          "Idea",
		"symbol-question":      "Question",
		"symbol-pin":           "Pin",
		"symbol-run":           "Run",
		"symbol-entertainment": "Entertainment",
		"symbol-100":           "100%",
	}

	for _, secondItem := range target.RootTopic.Children.Attached {
		secondTitle := secondItem.Title
		for _, thirdItem := range secondItem.Children.Attached {
			thirdTitle := thirdItem.Title
			caseMap := make(map[string]string)
			var priorityTag string
			var testResultTag string
			var otherTag string
			for _, thirdMark := range thirdItem.Markers {
				markString := thirdMark.MarkerID
				if strings.HasPrefix(markString, "priority") {
					priorityTag = priorityMap[markString]
				} else if strings.Contains(markString, "symbol_") {
					if markString == "c_symbol_like" || markString == "c_symbol_dislike" {
						testResultTag = testResultMap[markString]
					} else {
						if len(otherTag) == 0 {
							otherTag = testResultMap[markString]
						} else if !strings.Contains(otherTag, testResultMap[markString]) {
							otherTag = fmt.Sprintf("%s,%s", otherTag, testResultMap[markString])
						}
					}
				} else {
					if len(otherTag) == 0 {
						otherTag = markString
					} else if !strings.Contains(otherTag, markString) {
						otherTag = fmt.Sprintf("%s,%s", otherTag, markString)
					}
				}

			}

			for index3, fourthItem := range thirdItem.Children.Attached {
				fourthTitle := fourthItem.Title
				if fourthTitle == T("field.test_steps") {
					var stepStr, expectStr string
					var subResultTag string
					for _, fifthItem := range fourthItem.Children.Attached {
						if len(stepStr) == 0 {
							stepStr = fifthItem.Title
						} else {
							stepStr = fmt.Sprintf("%s;\n%s", stepStr, fifthItem.Title)
						}

						if len(fifthItem.Children.Attached) > 0 {
							for _, sixthItem := range fifthItem.Children.Attached {
								if len(expectStr) == 0 {
									expectStr = sixthItem.Title
								} else {
									expectStr = fmt.Sprintf("%s;\n%s", expectStr, sixthItem.Title)

								}
								if subResultTag == "Fail" || subResultTag == "PartTest" || len(testResultTag) > 0 {
									continue
								}

								//if len(sixthItem.Markers) == 0 {
								//	subResultTag = "UnTest"
								//} // 代码看着像写错了

								for _, sixthMark := range sixthItem.Markers {
									sixMarkString := sixthMark.MarkerID
									if strings.Contains(sixMarkString, "symbol_") && len(testResultTag) == 0 {
										if sixMarkString == "c_symbol_like" || sixMarkString == "c_symbol_dislike" {
											subResultTag = testResultMap[sixMarkString]
										} else {
											if len(otherTag) == 0 {
												otherTag = testResultMap[sixMarkString]
											} else if !strings.Contains(otherTag, testResultMap[sixMarkString]) {
												otherTag = fmt.Sprintf("%s,%s", otherTag, testResultMap[sixMarkString])
											}
										}
									} else {
										if len(otherTag) == 0 {
											otherTag = sixMarkString
										} else if !strings.Contains(otherTag, sixMarkString) {
											otherTag = fmt.Sprintf("%s,%s", otherTag, sixMarkString)
										}
									}
								}
							}
						}
					}
					caseMap[T("field.test_steps")] = stepStr
					caseMap[T("field.expected_result")] = expectStr
					if len(testResultTag) == 0 {
						if len(subResultTag) == 0 {
							testResultTag = "UnTest"
						} else {
							testResultTag = subResultTag
						}
					}
					caseMap[T("field.test_result")] = testResultTag
					caseMap[T("field.tags")] = otherTag

				} else {
					infos := strings.Split(fourthTitle, ":")
					if len(infos) > 1 {
						caseMap[infos[0]] = infos[1]
					}
				}

				if index3 == len(thirdItem.Children.Attached)-1 {
					caseMap[T("field.product")] = firstTitle
					caseMap[T("field.module")] = secondTitle
					caseMap[T("field.case_name")] = thirdTitle

					if strings.HasPrefix(thirdTitle, "#") {
						caseMap[T("field.test_result")] = "Deprecated"
					}

					if len(priorityTag) > 0 {
						caseMap[T("field.priority")] = priorityTag
					}

					caseList = append(caseList, caseMap)
				}
			}

		}
	}

	titleList := []string{T("field.product"), T("field.module"), T("field.case_number"), T("field.case_name"), T("field.priority"), T("field.test_range"), T("field.precondition"), T("field.test_steps"), T("field.expected_result"), T("field.test_result"), T("field.tags")}
	for index1, item := range caseList {
		var valueList []string

		if index1 == 0 {
			WriteDataInXls(fileCasePath, titleList)
		}

		for _, k := range titleList {
			if _, ok := item[k]; ok {
				valueList = append(valueList, item[k])
			} else {
				valueList = append(valueList, " ")
			}

		}
		WriteDataInXls(fileCasePath, valueList)
	}

	return
}

func Xmind2Import(product, introVersion, filePath string) (err error) {
	xmindContent, err := GetContentFromXmindFile(filePath)
	if err != nil {
		return err
	}
	if len(xmindContent) == 0 {
		err = fmt.Errorf(T("error.case_info_not_found"))
		return
	}

	target := xmindContent[0]
	firstTitle := target.RootTopic.Title
	var caseList []map[string]string

	for _, secondItem := range target.RootTopic.Children.Attached {
		secondTitle := secondItem.Title
		for _, thirdItem := range secondItem.Children.Attached {
			thirdTitle := thirdItem.Title
			caseMap := make(map[string]string)
			var priorityTag string

			if len(thirdItem.Markers) > 0 {
				priorityMark := thirdItem.Markers[0].MarkerID
				switch priorityMark {
				case "priority-1":
					priorityTag = "P1"
				case "priority-2":
					priorityTag = "P2"
				case "priority-3":
					priorityTag = "P3"
				case "priority-4":
					priorityTag = "P4"
				case "priority-5":
					priorityTag = "P5"
				}
			}
			for index3, fourthItem := range thirdItem.Children.Attached {
				fourthTitle := fourthItem.Title
				if fourthTitle == T("field.test_steps") {
					var stepStr, expectStr string
					for _, fifthItem := range fourthItem.Children.Attached {
						if len(stepStr) == 0 {
							stepStr = fifthItem.Title
						} else {
							stepStr = fmt.Sprintf("%s;\n%s", stepStr, fifthItem.Title)
						}

						if len(fifthItem.Children.Attached) > 0 {
							for _, sixthItem := range fifthItem.Children.Attached {
								if len(expectStr) == 0 {
									expectStr = sixthItem.Title
								} else {
									expectStr = fmt.Sprintf("%s;\n%s", expectStr, sixthItem.Title)

								}
							}
						}

					}
					caseMap[T("field.test_steps")] = stepStr
					caseMap[T("field.expected_result")] = expectStr
				} else {
					infos := strings.Split(fourthTitle, ":")
					if len(infos) > 1 {
						caseMap[infos[0]] = infos[1]
					}
				}

				if index3 == len(thirdItem.Children.Attached)-1 {
					caseMap[T("field.product")] = firstTitle
					caseMap[T("field.module")] = secondTitle
					caseMap[T("field.case_name")] = thirdTitle
					if len(priorityTag) > 0 {
						caseMap[T("field.priority")] = priorityTag
					}
					caseList = append(caseList, caseMap)
				}
			}

		}
	}

	//titleList := []string{T("field.product"), T("field.module"), T("field.case_name"), T("field.priority"), T("field.test_range"), T("field.precondition"), T("field.test_steps"), T("field.expected_result")}
	//for index1, item := range caseList {
	//	var valueList []string

	//if index1 == 0 {
	//	WriteDataInXls(fileCasePath, titleList)
	//}
	//
	//for _, k := range titleList {
	//	if _, ok := item[k]; ok {
	//		valueList = append(valueList, item[k])
	//	} else {
	//		valueList = append(valueList, " ")
	//	}
	//
	//}
	//WriteDataInXls(fileCasePath, valueList)
	//}

	return
}

func GetCaseFromXmind(filePath string) (caseList []map[string]string, err error) {
	xmindContent, err := GetContentFromXmindFile(filePath)
	if err != nil {
		Logger.Error("err: %v", err)
		return
	}

	if len(xmindContent) == 0 {
		err = fmt.Errorf(T("error.case_info_not_found"))
		return
	}

	target := xmindContent[0]
	firstTitle := target.RootTopic.Title
	//var caseList []map[string]string

	for _, secondItem := range target.RootTopic.Children.Attached {
		secondTitle := secondItem.Title
		for _, thirdItem := range secondItem.Children.Attached {
			thirdTitle := thirdItem.Title
			caseMap := make(map[string]string)
			var priorityTag string

			if len(thirdItem.Markers) > 0 {
				priorityMark := thirdItem.Markers[0].MarkerID
				switch priorityMark {
				case "priority-1":
					priorityTag = "P1"
				case "priority-2":
					priorityTag = "P2"
				case "priority-3":
					priorityTag = "P3"
				case "priority-4":
					priorityTag = "P4"
				case "priority-5":
					priorityTag = "P5"
				}
			}
			for index3, fourthItem := range thirdItem.Children.Attached {
				fourthTitle := fourthItem.Title
				if fourthTitle == T("field.test_steps") {
					var stepStr, expectStr string
					for _, fifthItem := range fourthItem.Children.Attached {
						if len(stepStr) == 0 {
							stepStr = fifthItem.Title
						} else {
							stepStr = fmt.Sprintf("%s;\n%s", stepStr, fifthItem.Title)
						}

						if len(fifthItem.Children.Attached) > 0 {
							for _, sixthItem := range fifthItem.Children.Attached {
								if len(expectStr) == 0 {
									expectStr = sixthItem.Title
								} else {
									expectStr = fmt.Sprintf("%s;\n%s", expectStr, sixthItem.Title)

								}
							}
						}

					}
					caseMap[T("field.test_steps")] = stepStr
					caseMap[T("field.expected_result")] = expectStr
				} else {
					infos := strings.Split(fourthTitle, ":")
					if len(infos) > 1 {
						caseMap[infos[0]] = infos[1]
					}
				}

				if index3 == len(thirdItem.Children.Attached)-1 {
					caseMap[T("field.product")] = firstTitle
					caseMap[T("field.module")] = secondTitle
					caseMap[T("field.case_name")] = thirdTitle
					if len(priorityTag) > 0 {
						caseMap[T("field.priority")] = priorityTag
					}
					caseList = append(caseList, caseMap)
				}
			}
		}
	}

	return
}

func Xmind2ImportAndUse(userName, product, introVersion, filePath string) (err error) {
	//xmindContent, err := GetContentFromXmindFile(filePath)
	//if err != nil {
	//	return err
	//}
	//if len(xmindContent) == 0 {
	//	err = fmt.Errorf(T("error.case_info_not_found"))
	//	return
	//}
	//
	//target := xmindContent[0]
	//firstTitle := target.RootTopic.Title
	//var caseList []map[string]string
	//
	//for _, secondItem := range target.RootTopic.Children.Attached {
	//	secondTitle := secondItem.Title
	//	for _, thirdItem := range secondItem.Children.Attached {
	//		thirdTitle := thirdItem.Title
	//		caseMap := make(map[string]string)
	//		var priorityTag string
	//
	//		if len(thirdItem.Markers) > 0 {
	//			priorityMark := thirdItem.Markers[0].MarkerID
	//			switch priorityMark {
	//			case "priority-1":
	//				priorityTag = "P1"
	//			case "priority-2":
	//				priorityTag = "P2"
	//			case "priority-3":
	//				priorityTag = "P3"
	//			case "priority-4":
	//				priorityTag = "P4"
	//			case "priority-5":
	//				priorityTag = "P5"
	//			}
	//		}
	//		for index3, fourthItem := range thirdItem.Children.Attached {
	//			fourthTitle := fourthItem.Title
	//			if fourthTitle == T("field.test_steps") {
	//				var stepStr, expectStr string
	//				for _, fifthItem := range fourthItem.Children.Attached {
	//					if len(stepStr) == 0 {
	//						stepStr = fifthItem.Title
	//					} else {
	//						stepStr = fmt.Sprintf("%s;\n%s", stepStr, fifthItem.Title)
	//					}
	//
	//					if len(fifthItem.Children.Attached) > 0 {
	//						for _, sixthItem := range fifthItem.Children.Attached {
	//							if len(expectStr) == 0 {
	//								expectStr = sixthItem.Title
	//							} else {
	//								expectStr = fmt.Sprintf("%s;\n%s", expectStr, sixthItem.Title)
	//
	//							}
	//						}
	//					}
	//
	//				}
	//				caseMap[T("field.test_steps")] = stepStr
	//				caseMap[T("field.expected_result")] = expectStr
	//			} else {
	//				infos := strings.Split(fourthTitle, ":")
	//				if len(infos) > 1 {
	//					caseMap[infos[0]] = infos[1]
	//				}
	//			}
	//
	//			if index3 == len(thirdItem.Children.Attached)-1 {
	//				caseMap[T("field.product")] = firstTitle
	//				caseMap[T("field.module")] = secondTitle
	//				caseMap[T("field.case_name")] = thirdTitle
	//				if len(priorityTag) > 0 {
	//					caseMap[T("field.priority")] = priorityTag
	//				}
	//				caseList = append(caseList, caseMap)
	//			}
	//		}
	//
	//	}
	//}
	caseList, err := GetCaseFromXmind(filePath)
	for _, testCase := range caseList {
		var DbAiCase DbAiCase
		err = models.Orm.Table("ai_case").Where("product = ? and case_number = ? and intro_version = ?", product, testCase["case_number"], introVersion).Find(&DbAiCase).Error
		if err != nil {
			Logger.Error("%s", err)
		}
		UseAiData(DbAiCase.Id, userName)
	}

	return
}
