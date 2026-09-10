package biz

import (
	"data4test/models"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/vivian0517/goxmind"
)

func ExportTestCase2Markdown(ids, source string) (fileName string, err error) {
	idList := strings.Split(ids, ",")
	var productNames, modules []string
	models.Orm.Table(source).Where("id in (?)", idList).Group("product").Pluck("product", &productNames)
	if len(productNames) == 0 {
		errTmp := fmt.Errorf(T("error.product_count_invalid"), len(productNames))
		return "", errTmp
	} else if len(productNames) > 1 {
		Logger.Warning(T("warn.product_count_default"), len(productNames))
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
				WriteDataInCommonFile(filePath, fmt.Sprintf(T("md.priority"), itemCase.Priority))
				WriteDataInCommonFile(filePath, fmt.Sprintf(T("md.precondition"), itemCase.PreCondition))
				WriteDataInCommonFile(filePath, fmt.Sprintf(T("md.test_range"), itemCase.TestRange))
				WriteDataInCommonFile(filePath, fmt.Sprintf(T("md.test_steps")))

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

func ExportTestCase2XmindById(ids, source string) (fileName string, err error) {
	idList := strings.Split(ids, ",")
	var productNames, modules []string
	models.Orm.Table(source).Where("id in (?)", idList).Group("product").Pluck("product", &productNames)
	if len(productNames) == 0 {
		errTmp := fmt.Errorf(T("error.product_count_invalid"), len(productNames))
		return "", errTmp
	} else if len(productNames) > 1 {
		Logger.Warning(T("warn.product_count_default"), len(productNames))
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
				case "P1", T("priority.high"):
					priorityMark = goxmind.Priority1
				case "P2", T("priority.medium"):
					priorityMark = goxmind.Priority2
				case "P3", T("priority.low"):
					priorityMark = goxmind.Priority3
				case "P4":
					priorityMark = goxmind.Priority4
				default:
					priorityMark = goxmind.Priority2
				}
				thirdNode.AddMaker(priorityMark)
				thirdNode.AddNode(fmt.Sprintf(T("xmind.case_number"), itemCase.CaseNumber))
				thirdNode.AddNode(fmt.Sprintf(T("xmind.precondition"), itemCase.PreCondition))
				thirdNode.AddNode(fmt.Sprintf(T("xmind.test_range"), itemCase.TestRange))
				if itemCase.Auto == "1" {
					thirdNode.AddNode(fmt.Sprintf(T("xmind.auto_yes")))
				} else if itemCase.Auto == "0" {
					thirdNode.AddNode(fmt.Sprintf(T("xmind.auto_no")))
				}

				fourthNode := thirdNode.AddNode(fmt.Sprintf(T("xmind.test_steps")))
				var testSteps, expectResults []string
				if strings.Contains(itemCase.TestSteps, ";") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", "", -1), ";")
				} else if strings.Contains(itemCase.TestSteps, "；") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", "", -1), "；")
				} else if strings.Contains(itemCase.TestSteps, "\n") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", ";", -1), ";")
				} else if strings.Contains(itemCase.TestSteps, " ") { // 部分数据有幻觉，会拆的格式不正常
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, " ", ";", -1), ";")
				} else {
					testSteps = []string{itemCase.TestSteps}
				}

				if strings.Contains(itemCase.ExpectResult, ";") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", "", -1), ";")
				} else if strings.Contains(itemCase.ExpectResult, "；") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", "", -1), "；")
				} else if strings.Contains(itemCase.ExpectResult, "\n") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", ";", -1), ";")
				} else if strings.Contains(itemCase.ExpectResult, " ") { // 部分数据有幻觉，会拆的格式不正常
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, " ", ";", -1), ";")
				} else {
					expectResults = []string{itemCase.ExpectResult}
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

func ExportAiCase2XmindByCondition(product, introVersion, createPlatform, module, createUser, createdAtStart, createdAtEnd, source string) (fileName string, err error) {
	var productNames, modules []string
	dbHandle := models.Orm.Table(source)
	if len(product) > 0 {
		dbHandle = dbHandle.Where("product = ?", product)
	}

	if len(createPlatform) > 0 {
		dbHandle = dbHandle.Where("source = ?", createPlatform)
	}

	if len(createUser) > 0 {
		dbHandle = dbHandle.Where("create_user = ?", createUser)
	}

	if len(createdAtStart) > 0 {
		dbHandle = dbHandle.Where("created_at >= ?", createdAtStart)
	}

	if len(createdAtEnd) > 0 {
		dbHandle = dbHandle.Where("created_at < ?", createdAtEnd)
	}

	if len(module) > 0 {
		moduleList := strings.Split(module, ",")
		dbHandle = dbHandle.Where("module in (?)", moduleList)
	}

	if len(introVersion) > 0 {
		introVersionList := strings.Split(introVersion, ",")
		dbHandle = dbHandle.Where("intro_version in (?)", introVersionList)
	}

	if len(product) == 0 {
		dbHandle.Group("product").Pluck("product", &productNames)
	} else {
		productNames = append(productNames, product)
	}

	if len(productNames) == 0 {
		errTmp := fmt.Errorf(T("error.product_count_invalid"), len(productNames))
		return "", errTmp
	} else if len(productNames) > 1 {
		Logger.Warning(T("warn.product_count_default"), len(productNames))
	}

	if len(module) == 0 {
		dbHandle.Group("module").Pluck("module", &modules)
	} else {
		modules = strings.Split(module, ",")
	}

	timeFormat := "20060102150405"
	curTime := time.Now().Format(timeFormat)
	fileName = fmt.Sprintf("%s_%s.xmind", productNames[0], curTime)
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	xmind := goxmind.New()
	firstNode := xmind.AddSheet("Sheet title", productNames[0])

	for _, module := range modules {
		var testCases []TestCase
		dbHandle.Where("module = ?", module).
			Find(&testCases)
		if len(testCases) > 0 {
			secondNode := firstNode.AddNode(module)
			for _, itemCase := range testCases {
				thirdNode := secondNode.AddNode(itemCase.CaseName)
				var priorityMark goxmind.MarkerId
				switch itemCase.Priority {
				case "P0":
					priorityMark = goxmind.Priority1
				case "P1", T("priority.high"):
					priorityMark = goxmind.Priority1
				case "P2", T("priority.medium"):
					priorityMark = goxmind.Priority2
				case "P3", T("priority.low"):
					priorityMark = goxmind.Priority3
				case "P4":
					priorityMark = goxmind.Priority4
				default:
					priorityMark = goxmind.Priority2
				}
				thirdNode.AddMaker(priorityMark)
				thirdNode.AddNode(fmt.Sprintf(T("xmind.case_number"), itemCase.CaseNumber))
				thirdNode.AddNode(fmt.Sprintf(T("xmind.precondition"), itemCase.PreCondition))
				thirdNode.AddNode(fmt.Sprintf(T("xmind.test_range"), itemCase.TestRange))
				if itemCase.Auto == "1" {
					thirdNode.AddNode(fmt.Sprintf(T("xmind.auto_yes")))
				} else if itemCase.Auto == "0" {
					thirdNode.AddNode(fmt.Sprintf(T("xmind.auto_no")))
				}

				fourthNode := thirdNode.AddNode(fmt.Sprintf(T("xmind.test_steps")))
				var testSteps, expectResults []string
				if strings.Contains(itemCase.TestSteps, ";") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", "", -1), ";")
				} else if strings.Contains(itemCase.TestSteps, "；") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", "", -1), "；")
				} else if strings.Contains(itemCase.TestSteps, "\n") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", ";", -1), ";")
				} else if strings.Contains(itemCase.TestSteps, " ") { // 部分数据有幻觉，会拆的格式不正常
					itemCase.TestSteps = strings.Replace(itemCase.TestSteps, ". ", ".", -1)
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, " ", ";", -1), ";")
				} else {
					testSteps = []string{itemCase.TestSteps}
				}

				if strings.Contains(itemCase.ExpectResult, ";") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", "", -1), ";")
				} else if strings.Contains(itemCase.ExpectResult, "；") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", "", -1), "；")
				} else if strings.Contains(itemCase.ExpectResult, "\n") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", ";", -1), ";")
				} else if strings.Contains(itemCase.ExpectResult, " ") { // 部分数据有幻觉，会拆的格式不正常
					itemCase.ExpectResult = strings.Replace(itemCase.ExpectResult, ". ", ".", -1)
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, " ", ";", -1), ";")
				} else {
					expectResults = []string{itemCase.ExpectResult}
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

func ExportTestCase2XmindByCondition(product, introVersion, module, caseDesigner, createdAtStart, createdAtEnd, source string) (fileName string, err error) {
	var productNames, modules []string
	dbHandle := models.Orm.Table(source)
	if len(product) > 0 {
		dbHandle = dbHandle.Where("product = ?", product)
	}

	if len(caseDesigner) > 0 {
		dbHandle = dbHandle.Where("case_designer = ?", caseDesigner)
	}

	if len(createdAtStart) > 0 {
		dbHandle = dbHandle.Where("created_at >= ?", createdAtStart)
	}

	if len(createdAtEnd) > 0 {
		dbHandle = dbHandle.Where("created_at < ?", createdAtEnd)
	}

	if len(module) > 0 {
		moduleList := strings.Split(module, ",")
		dbHandle = dbHandle.Where("module in (?)", moduleList)
	}

	if len(introVersion) > 0 {
		introVersionList := strings.Split(introVersion, ",")
		dbHandle = dbHandle.Where("intro_version in (?)", introVersionList)
	}

	if len(product) == 0 {
		dbHandle.Group("product").Pluck("product", &productNames)
	} else {
		productNames = append(productNames, product)
	}

	if len(productNames) == 0 {
		errTmp := fmt.Errorf(T("error.product_count_invalid"), len(productNames))
		return "", errTmp
	} else if len(productNames) > 1 {
		Logger.Warning(T("warn.product_count_default"), len(productNames))
	}

	if len(module) == 0 {
		dbHandle.Group("module").Pluck("module", &modules)
	} else {
		modules = strings.Split(module, ",")
	}

	timeFormat := "20060102150405"
	curTime := time.Now().Format(timeFormat)
	fileName = fmt.Sprintf("%s_%s.xmind", productNames[0], curTime)
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	xmind := goxmind.New()
	firstNode := xmind.AddSheet("Sheet title", productNames[0])

	for _, module := range modules {
		var testCases []TestCase
		dbHandle.Where("module = ?", module).
			Find(&testCases)
		if len(testCases) > 0 {
			secondNode := firstNode.AddNode(module)
			for _, itemCase := range testCases {
				thirdNode := secondNode.AddNode(itemCase.CaseName)
				var priorityMark goxmind.MarkerId
				switch itemCase.Priority {
				case "P0":
					priorityMark = goxmind.Priority1
				case "P1", T("priority.high"):
					priorityMark = goxmind.Priority1
				case "P2", T("priority.medium"):
					priorityMark = goxmind.Priority2
				case "P3", T("priority.low"):
					priorityMark = goxmind.Priority3
				case "P4":
					priorityMark = goxmind.Priority4
				default:
					priorityMark = goxmind.Priority2
				}
				thirdNode.AddMaker(priorityMark)
				thirdNode.AddNode(fmt.Sprintf(T("xmind.case_number"), itemCase.CaseNumber))
				thirdNode.AddNode(fmt.Sprintf(T("xmind.precondition"), itemCase.PreCondition))
				thirdNode.AddNode(fmt.Sprintf(T("xmind.test_range"), itemCase.TestRange))
				if itemCase.Auto == "1" {
					thirdNode.AddNode(fmt.Sprintf(T("xmind.auto_yes")))
				} else if itemCase.Auto == "0" {
					thirdNode.AddNode(fmt.Sprintf(T("xmind.auto_no")))
				}

				fourthNode := thirdNode.AddNode(fmt.Sprintf(T("xmind.test_steps")))
				var testSteps, expectResults []string
				if strings.Contains(itemCase.TestSteps, ";") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", "", -1), ";")
				} else if strings.Contains(itemCase.TestSteps, "；") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", "", -1), "；")
				} else if strings.Contains(itemCase.TestSteps, "\n") {
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, "\n", ";", -1), ";")
				} else if strings.Contains(itemCase.TestSteps, " ") { // 部分数据有幻觉，会拆的格式不正常
					testSteps = strings.Split(strings.Replace(itemCase.TestSteps, " ", ";", -1), ";")
				} else {
					testSteps = []string{itemCase.TestSteps}
				}

				if strings.Contains(itemCase.ExpectResult, ";") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", "", -1), ";")
				} else if strings.Contains(itemCase.ExpectResult, "；") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", "", -1), "；")
				} else if strings.Contains(itemCase.ExpectResult, "\n") {
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, "\n", ";", -1), ";")
				} else if strings.Contains(itemCase.ExpectResult, " ") { // 部分数据有幻觉，会拆的格式不正常
					expectResults = strings.Split(strings.Replace(itemCase.ExpectResult, " ", ";", -1), ";")
				} else {
					expectResults = []string{itemCase.ExpectResult}
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

// AutoFillCaseExecutorTime 根据测试结果/测试时间变更返回需要自动填充的字段。
// isSingleUpdate 为 true 表示列表内联编辑（仅更新执行者）；false 表示编辑页保存。
func AutoFillCaseExecutorTime(id, testResult, testTime, userName string, isSingleUpdate bool) map[string]string {
	var orig struct {
		TestResult    string `gorm:"column:test_result"`
		TestTime      string `gorm:"column:test_time"`
		CaseExecutor  string `gorm:"column:case_executor"`
		ResultHistory string `gorm:"column:result_history"`
	}
	models.Orm.Table("test_case").Where("id = ?", id).
		Select("test_result, test_time, case_executor, result_history").Scan(&orig)

	updates := map[string]string{}
	if isSingleUpdate {
		// 列表内联编辑：操作测试结果 → 自动更新用例执行者
		if testResult != "" && testResult != orig.TestResult {
			updates["case_executor"] = userName
		}
	} else {
		// 编辑页保存：操作测试结果 → 自动更新测试时间 + 执行者
		if testResult != orig.TestResult {
			updates["test_time"] = time.Now().Format(baseFormat)
			updates["case_executor"] = userName
		} else if testTime != orig.TestTime {
			// 操作测试时间 → 自动更新执行者
			updates["case_executor"] = userName
		}
	}

	// 结果变化 → 把旧结果追加进历史（首次执行不写）
	if testResult != "" && testResult != orig.TestResult {
		if newHistory, ok := appendTestCaseResultHistory(orig.ResultHistory, orig.TestResult, orig.CaseExecutor, orig.TestTime); ok {
			updates["result_history"] = newHistory
		}
	}
	return updates
}

// appendTestCaseResultHistory 追加一条历史记录；首次执行（旧结果为空或默认未测试且无执行人）返回 false。
func appendTestCaseResultHistory(history, oldResult, oldExecutor, oldTime string) (string, bool) {
	if oldResult == "" || (oldResult == "untest" && oldExecutor == "") {
		return "", false
	}
	seq := 1
	if history != "" {
		seq = len(strings.Split(history, "\n")) + 1
	}
	line := fmt.Sprintf("%d|%s|%s|%s", seq, oldResult, oldExecutor, oldTime)
	if history == "" {
		return line, true
	}
	return history + "\n" + line, true
}

// BatchUpdateTestCase 批量修改选中用例；空字符串表示不修改该字段。
// test_result 变化时沿用单条编辑的历史记录逻辑（追加旧结果 + 自动记录执行者/时间）。
func BatchUpdateTestCase(ids []string, module, testResult, introVersion, funDeveloper, caseDesigner, caseExecutor, testTime, userName string) (err error) {
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		updates := map[string]interface{}{}
		if module != "" {
			// 模块变化时，迁移测试过程图片到新模块目录
			var orig struct {
				CaseNumber  string `gorm:"column:case_number"`
				Module      string `gorm:"column:module"`
				TestProcess string `gorm:"column:test_process"`
			}
			models.Orm.Table("test_case").Where("id = ?", id).
				Select("case_number, module, test_process").Scan(&orig)
			if orig.Module != "" && orig.Module != module {
				if newHTML, e := MigrateTestCaseProcessImages(orig.CaseNumber, orig.Module, module, orig.TestProcess); e != nil {
					Logger.Error("迁移测试过程图片失败 id=%s: %v", id, e)
				} else if newHTML != orig.TestProcess {
					updates["test_process"] = newHTML
				}
			}
			updates["module"] = module
		}
		if introVersion != "" {
			updates["intro_version"] = introVersion
		}
		if funDeveloper != "" {
			updates["fun_developer"] = funDeveloper
		}
		if caseDesigner != "" {
			updates["case_designer"] = caseDesigner
		}
		if testResult != "" {
			for k, v := range AutoFillCaseExecutorTime(id, testResult, testTime, userName, false) {
				updates[k] = v
			}
			updates["test_result"] = testResult
		}
		// 显式填写的执行者/时间覆盖自动值
		if caseExecutor != "" {
			updates["case_executor"] = caseExecutor
		}
		if testTime != "" {
			updates["test_time"] = testTime
		}
		if len(updates) == 0 {
			continue
		}
		if e := models.Orm.Table("test_case").Where("id = ?", id).Updates(updates).Error; e != nil {
			Logger.Error("BatchUpdateTestCase id=%s err=%v", id, e)
			return e
		}
	}
	return nil
}

// TestCaseImageUploadHandler 测试过程图片上传（保存到模块目录，按 用例编号_N 命名）
var TestCaseImageUploadHandler context.Handler = func(ctx *context.Context) {
	if ctx.Request.MultipartForm == nil || len(ctx.Request.MultipartForm.File["file"]) == 0 {
		ctx.JSON(200, map[string]interface{}{"errno": 400})
		return
	}

	fileHeader := ctx.Request.MultipartForm.File["file"][0]
	f, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(200, map[string]interface{}{"errno": 500})
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		ctx.JSON(200, map[string]interface{}{"errno": 500})
		return
	}

	caseNumber, module := "", ""
	if v := ctx.Request.MultipartForm.Value["case_number"]; len(v) > 0 {
		caseNumber = v[0]
	}
	if v := ctx.Request.MultipartForm.Value["module"]; len(v) > 0 {
		module = v[0]
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".png"
	}

	imgURL, err := saveTestCaseProcessImage(caseNumber, module, ext, data)
	if err != nil {
		Logger.Error("保存测试过程图片失败: %s", err)
		ctx.JSON(200, map[string]interface{}{"errno": 500})
		return
	}

	ctx.JSON(200, map[string]interface{}{"errno": 0, "data": []string{imgURL}})
}

// saveTestCaseProcessImage 保存图片到模块目录，命名为 用例编号_N.ext，返回可访问 URL
func saveTestCaseProcessImage(caseNumber, module, ext string, data []byte) (string, error) {
	moduleDir := sanitizeDirName(module)
	caseNo := sanitizeFileName(caseNumber)
	if caseNo == "" {
		caseNo = "case"
	}

	dir := filepath.Join(UploadBasePath, moduleDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	idx := nextImageIndex(dir, caseNo)
	filename := fmt.Sprintf("%s_%d%s", caseNo, idx, ext)
	if err := os.WriteFile(filepath.Join(dir, filename), data, 0644); err != nil {
		return "", err
	}

	return "/uploads/" + url.PathEscape(moduleDir) + "/" + url.PathEscape(filename), nil
}

// sanitizeDirName 清洗模块名为安全目录名
func sanitizeDirName(s string) string {
	s = strings.TrimSpace(s)
	r := strings.NewReplacer("/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_")
	s = r.Replace(s)
	if s == "" {
		s = "ungrouped"
	}
	return s
}

// sanitizeFileName 清洗用例编号为安全文件名
func sanitizeFileName(s string) string {
	s = strings.TrimSpace(s)
	r := strings.NewReplacer("/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_")
	return r.Replace(s)
}

// nextImageIndex 计算下一个图片序号（模块目录内 用例编号_N 的最大 N + 1）
func nextImageIndex(dir, caseNumber string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 1
	}
	prefix := caseNumber + "_"
	max := 0
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		base := strings.TrimSuffix(name, filepath.Ext(name))
		numStr := strings.TrimPrefix(base, prefix)
		if n, err := strconv.Atoi(numStr); err == nil && n > max {
			max = n
		}
	}
	return max + 1
}

// imgSrcRe 匹配 <img ... src="..." ...> 标签中的 src 属性（大小写不敏感）
var imgSrcRe = regexp.MustCompile(`(?i)<img[^>]*?\bsrc\s*=\s*["']([^"']+)["']`)

// SyncTestCaseProcessImages 保存用例时同步测试过程图片：
//  1. 删除富文本中已移除的图片文件（文件名形如 用例编号_N.ext）；
//  2. 将剩余图片按富文本中出现顺序重命名为 用例编号_1、_2 ...
//
// 返回重命名后的富文本；若无需变更则原样返回。
func SyncTestCaseProcessImages(html, caseNumber, module string) string {
	if html == "" || caseNumber == "" {
		return html
	}
	caseNo := sanitizeFileName(caseNumber)
	if caseNo == "" {
		return html
	}

	srcMatches := imgSrcRe.FindAllStringSubmatch(html, -1)

	// 受管图片：文件名形如 用例编号_N.ext，位于 /uploads/<目录>/ 下
	type imgRef struct{ dir, name, src string }
	var managed []imgRef
	for _, m := range srcMatches {
		dir, name, ok := parseUploadSrc(m[1])
		if !ok || !isManagedImage(name, caseNo) {
			continue
		}
		managed = append(managed, imgRef{dir, name, m[1]})
	}

	// 兜底：当前模块目录也纳入扫描，用于删除全部移除的图片
	dirs := map[string]bool{sanitizeDirName(module): true}
	surviveSet := map[string]map[string]bool{} // dir -> 存活文件名集合
	for _, r := range managed {
		dirs[r.dir] = true
		if surviveSet[r.dir] == nil {
			surviveSet[r.dir] = map[string]bool{}
		}
		surviveSet[r.dir][r.name] = true
	}

	// 1) 删除不再被引用的图片文件
	for dir := range dirs {
		for _, name := range caseImageFiles(dir, caseNo) {
			if !surviveSet[dir][name] {
				if err := os.Remove(filepath.Join(UploadBasePath, dir, name)); err != nil {
					Logger.Error("删除测试过程图片失败 %s/%s: %s", dir, name, err)
				}
			}
		}
	}

	// 2) 重命名剩余图片为 1..N（按富文本中出现顺序，去重）
	order := map[string][]string{} // dir -> 有序存活文件名
	for _, r := range managed {
		order[r.dir] = append(order[r.dir], r.name)
	}

	replace := map[string]string{} // 原 src -> 新 src
	for dir, ordered := range order {
		seen := map[string]bool{}
		var seq []string
		for _, name := range ordered {
			if !seen[name] {
				seen[name] = true
				seq = append(seq, name)
			}
		}
		targets := make([]string, len(seq))
		needRename := false
		for i, name := range seq {
			targets[i] = fmt.Sprintf("%s_%d%s", caseNo, i+1, filepath.Ext(name))
			if name != targets[i] {
				needRename = true
			}
		}
		if !needRename {
			continue
		}
		// 两阶段重命名，避免目标名被占用时覆盖
		for i, name := range seq {
			if name == targets[i] {
				continue
			}
			tmp := fmt.Sprintf("%s__tmp_%d%s", caseNo, i, filepath.Ext(name))
			if err := os.Rename(filepath.Join(UploadBasePath, dir, name), filepath.Join(UploadBasePath, dir, tmp)); err != nil {
				Logger.Error("重命名测试过程图片失败 %s/%s: %s", dir, name, err)
			}
		}
		for i, name := range seq {
			if name == targets[i] {
				continue
			}
			tmp := fmt.Sprintf("%s__tmp_%d%s", caseNo, i, filepath.Ext(name))
			if err := os.Rename(filepath.Join(UploadBasePath, dir, tmp), filepath.Join(UploadBasePath, dir, targets[i])); err != nil {
				Logger.Error("重命名测试过程图片失败 %s/%s: %s", dir, name, err)
			}
		}
		// 记录 src 替换（同一文件可能被多次引用）
		for i, name := range seq {
			if name == targets[i] {
				continue
			}
			for _, r := range managed {
				if r.dir == dir && r.name == name {
					replace[r.src] = "/uploads/" + url.PathEscape(dir) + "/" + url.PathEscape(targets[i])
				}
			}
		}
	}

	if len(replace) == 0 {
		return html
	}
	for old, new := range replace {
		html = strings.ReplaceAll(html, old, new)
	}
	return html
}

// parseUploadSrc 解析图片 src，返回去转义后的目录名与文件名。
// 仅处理 /uploads/<目录>/<文件名> 形式的引用（旧版 uuid 图片直接位于 upload 根目录，返回 false）。
func parseUploadSrc(src string) (dir, name string, ok bool) {
	i := strings.Index(src, "/uploads/")
	if i < 0 {
		return "", "", false
	}
	rest := strings.TrimPrefix(src[i:], "/uploads/")
	if q := strings.IndexByte(rest, '?'); q >= 0 {
		rest = rest[:q]
	}
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	d, err1 := url.PathUnescape(parts[0])
	n, err2 := url.PathUnescape(parts[1])
	if err1 != nil || err2 != nil {
		return "", "", false
	}
	return d, n, true
}

// isManagedImage 判断文件名是否属于指定用例的受管图片（用例编号_数字.扩展名）
func isManagedImage(filename, caseNo string) bool {
	if !strings.HasPrefix(filename, caseNo+"_") {
		return false
	}
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	numStr := strings.TrimPrefix(base, caseNo+"_")
	if numStr == "" {
		return false
	}
	for _, c := range numStr {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// caseImageFiles 列出目录下属于指定用例编号的图片文件（文件名形如 用例编号_N.ext）
func caseImageFiles(dir, caseNo string) []string {
	entries, err := os.ReadDir(filepath.Join(UploadBasePath, dir))
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isManagedImage(e.Name(), caseNo) {
			names = append(names, e.Name())
		}
	}
	return names
}

// MigrateTestCaseProcessImages 迁移测试过程图片：用例模块变化时，把旧模块目录下属于该用例的图片文件移动到新模块目录，
// 并重写富文本中的 /uploads/<旧目录>/ 引用为 /uploads/<新目录>/。返回新的富文本；无变化则原样返回。
func MigrateTestCaseProcessImages(caseNumber, oldModule, newModule, html string) (newHTML string, err error) {
	if html == "" || oldModule == newModule {
		return html, nil
	}
	caseNo := sanitizeFileName(caseNumber)
	if caseNo == "" {
		return html, nil
	}
	oldDir := sanitizeDirName(oldModule)
	newDir := sanitizeDirName(newModule)
	if oldDir == newDir {
		return html, nil
	}

	newDirPath := filepath.Join(UploadBasePath, newDir)
	if err := os.MkdirAll(newDirPath, os.ModePerm); err != nil {
		return html, err
	}

	oldDirPath := filepath.Join(UploadBasePath, oldDir)
	for _, name := range caseImageFiles(oldDir, caseNo) {
		if err := os.Rename(filepath.Join(oldDirPath, name), filepath.Join(newDirPath, name)); err != nil {
			Logger.Error("迁移测试过程图片失败 %s/%s -> %s/%s: %s", oldDir, name, newDir, name, err)
		}
	}

	oldPrefix := "/uploads/" + url.PathEscape(oldDir) + "/"
	newPrefix := "/uploads/" + url.PathEscape(newDir) + "/"
	return strings.ReplaceAll(html, oldPrefix, newPrefix), nil
}
