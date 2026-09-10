package biz

import (
	"archive/zip"
	"data4test/models"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

// caseImportFields 可导入/回写的系统字段（case_number 作为标识键单独处理，test_process 走图片逻辑）
var caseImportFields = []string{
	"case_name", "case_type", "priority", "pre_condition", "test_range",
	"test_steps", "expect_result", "auto", "scene", "fun_developer",
	"case_designer", "case_executor", "test_time", "test_result", "module",
	"intro_version", "product", "remark", "ext_info",
}

// imageNameRegexp 图片包内文件名后缀约束：N.format
var imageNameRegexp = regexp.MustCompile(`^\d+\.[A-Za-z0-9]+$`)

// sampleCaseNumberPrefix 导入模板示例行的用例编号前缀，导入时自动跳过
const sampleCaseNumberPrefix = "示例"

// testResultDropList 测试结果下拉选项（中文标签，导入时由 normalizeTestResult 归一化）
var testResultDropList = []string{"通过", "失败", "未测", "已废弃", "暂未合入"}

// sampleFieldValues 导入模板示例数据行的取值，key 为模板列 field（含 ext_info.xxx 子键）
var sampleFieldValues = map[string]string{
	"case_number":     "示例-001",
	"module":          "登录模块",
	"case_name":       "正常登录成功",
	"case_type":       "功能测试",
	"priority":        "高",
	"pre_condition":   "已注册有效账号",
	"test_range":      "Web 登录页",
	"test_steps":      "1. 打开登录页\n2. 输入账号密码\n3. 点击登录",
	"expect_result":   "登录成功，跳转首页",
	"auto":            "1",
	"scene":           "登录场景",
	"fun_developer":   "张三",
	"case_designer":   "李四",
	"case_executor":   "王五",
	"test_time":       "2026-09-08 10:00:00",
	"test_result":     "通过",
	"intro_version":   "v1.2.0",
	"product":         "示例产品",
	"remark":          "示例数据，导入时自动跳过",
	"ext_info":        "amount: 100\nnote: 首轮验证",
	"ext_info.amount": "100",
	"ext_info.note":   "首轮验证",
}

// normalizeTestResult 将导入的测试结果原始值归一化为系统标准值（中英文均支持，空值默认未测试）
func normalizeTestResult(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "pass", "通过", "passed", "成功":
		return "pass"
	case "fail", "失败", "failed", "不通过":
		return "fail"
	case "untest", "未测", "未测试", "未执行", "待测试":
		return "untest"
	case "n/a", "na", "已废弃", "废弃", "不适用":
		return "deprecated"
	case "unmerged", "暂未合入", "未合入":
		return "unmerged"
	case "":
		return "untest"
	default:
		return value
	}
}

// normalizeTestTime 将导入的测试时间归一化为标准格式 "2006-01-02 15:04:05"
// 兼容 Excel 日期序列号（如 46262.41667）与常见日期字符串格式；无法识别时原样保留
func normalizeTestTime(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return ""
	}

	// 纯数字（无日期分隔符）且落在 Excel 序列号合理区间(约 1955~2065 年) → 按序列号转时间
	if !strings.ContainsAny(v, "/-:.") {
		if serial, err := strconv.ParseFloat(v, 64); err == nil && serial >= 20000 && serial <= 60000 {
			return excelSerialToTime(serial).Format("2006-01-02 15:04:05")
		}
	}

	layouts := []string{
		"2006-01-02 15:04:05", "2006-01-02 15:04", "2006-01-02",
		"2006/01/02 15:04:05", "2006/01/02 15:04", "2006/01/02",
		"2006-1-2 15:04:05", "2006-1-2 15:04", "2006-1-2",
		"2006/1/2 15:04:05", "2006/1/2 15:04", "2006/1/2",
		"2006.01.02 15:04:05", "2006.01.02",
		"01/02/2006 15:04:05", "01/02/2006",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, v); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}
	return value
}

// excelSerialToTime 将 Excel 日期序列号（1900 纪元，含闰年 bug）转为时间
func excelSerialToTime(serial float64) time.Time {
	days := int(serial)
	frac := serial - float64(days)
	t := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC).AddDate(0, 0, days)
	return t.Add(time.Duration(frac * 24 * float64(time.Hour)))
}

// setImportField 将字段值写入用例（field 为 DB 列名）
func setImportField(item *TestCaseImportItem, field, value string) {
	switch field {
	case "case_name":
		item.CaseName = value
	case "case_type":
		item.CaseType = value
	case "priority":
		item.Priority = value
	case "pre_condition":
		item.PreCondition = value
	case "test_range":
		item.TestRange = value
	case "test_steps":
		item.TestSteps = value
	case "expect_result":
		item.ExpectResult = value
	case "auto":
		item.Auto = value
	case "scene":
		item.Scene = value
	case "fun_developer":
		item.FunDeveloper = value
	case "case_designer":
		item.CaseDesigner = value
	case "case_executor":
		item.CaseExecutor = value
	case "test_time":
		item.TestTime = normalizeTestTime(value)
	case "test_result":
		item.TestResult = normalizeTestResult(value)
	case "module":
		item.Module = value
	case "intro_version":
		item.IntroVersion = value
	case "product":
		item.Product = value
	case "remark":
		item.Remark = value
	case "ext_info":
		item.ExtInfo = value
	}
}

// getImportField 读取用例字段值（field 为 DB 列名）
func getImportField(item *TestCaseImportItem, field string) string {
	switch field {
	case "case_name":
		return item.CaseName
	case "case_type":
		return item.CaseType
	case "priority":
		return item.Priority
	case "pre_condition":
		return item.PreCondition
	case "test_range":
		return item.TestRange
	case "test_steps":
		return item.TestSteps
	case "expect_result":
		return item.ExpectResult
	case "auto":
		return item.Auto
	case "scene":
		return item.Scene
	case "fun_developer":
		return item.FunDeveloper
	case "case_designer":
		return item.CaseDesigner
	case "case_executor":
		return item.CaseExecutor
	case "test_time":
		return item.TestTime
	case "test_result":
		return item.TestResult
	case "module":
		return item.Module
	case "intro_version":
		return item.IntroVersion
	case "product":
		return item.Product
	case "remark":
		return item.Remark
	case "ext_info":
		return item.ExtInfo
	case "test_process":
		return item.TestProcess
	}
	return ""
}

// getCaseTypeList 返回用例类型下拉选项（来自 sys_parameter: TestCaseType，逗号分隔）
func getCaseTypeList() []string {
	var dbParams []SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", "TestCaseType").Find(&dbParams)
	if len(dbParams) == 0 {
		return nil
	}
	parts := strings.Split(dbParams[0].ValueList, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// addDropList 给指定列的数据区(第2行起)设置下拉数据验证
func addDropList(xlsxFile *excelize.File, sheet, colLetter string, values []string) {
	if len(values) == 0 {
		return
	}
	// 内联下拉列表受 Excel 255 字符限制，过长则跳过（避免下拉失效）
	if len(strings.Join(values, ",")) > 250 {
		return
	}
	dv := excelize.NewDataValidation(true)
	_ = dv.SetDropList(values)
	dv.SetSqref(colLetter + "2:" + colLetter + "1000")
	xlsxFile.AddDataValidation(sheet, dv)
}

// GenerateCaseImportTemplate 按模板名生成含表头 + 示例数据行的模板，供用户填写后导入
func GenerateCaseImportTemplate(templateName string) (fileName string, err error) {
	template, err := GetTestCaseExportTemplateByName(templateName)
	if err != nil {
		return
	}

	xlsxFile := excelize.NewFile()
	sheet := "Sheet1"
	textStyle, _ := xlsxFile.NewStyle(`{"custom_number_format":"@"}`)
	caseTypeValues := getCaseTypeList()
	productValues := GetCaseProductList()

	for c, col := range template.Columns {
		colLetter := columnName(c)
		xlsxFile.SetCellValue(sheet, colLetter+"1", col.Title)

		// 示例数据行
		if v, ok := sampleFieldValues[col.Field]; ok {
			xlsxFile.SetCellValue(sheet, colLetter+"2", v)
		}

		// 枚举列下拉
		switch col.Field {
		case "case_type":
			addDropList(xlsxFile, sheet, colLetter, caseTypeValues)
		case "test_result":
			addDropList(xlsxFile, sheet, colLetter, testResultDropList)
		case "product":
			addDropList(xlsxFile, sheet, colLetter, productValues)
		}

		// 测试时间列设文本格式，避免 Excel 自动转日期序列号
		if col.Field == "test_time" {
			xlsxFile.SetCellStyle(sheet, colLetter+"2", colLetter+"1000", textStyle)
		}
	}

	curTime := time.Now().Format("20060102150405")
	baseName := fmt.Sprintf("%s_%s", templateName, curTime)
	xlsxPath := fmt.Sprintf("%s/%s.xlsx", CaseFilePath, baseName)
	if err = xlsxFile.SaveAs(xlsxPath); err != nil {
		Logger.Error("%s", err)
		return
	}
	return baseName + ".xlsx", nil
}

// TestCaseImportCheck 解析上传的 Excel 与图片包，检查冲突，将完整结果保存到临时文件
func TestCaseImportCheck(excelPath, imagePkgPath, templateName, defaultsJSON, userName string) (result TestCaseImportResult, err error) {
	importId := uuid.NewV4().String()
	LastCaseImportId = importId
	result.ImportId = importId
	result.UserName = userName
	result.Template = templateName

	template, err := GetTestCaseExportTemplateByName(templateName)
	if err != nil {
		return
	}

	// 统一填写默认值（仅对未映射/无值字段生效）
	defaults := map[string]string{}
	if len(strings.TrimSpace(defaultsJSON)) > 0 {
		if jerr := json.Unmarshal([]byte(defaultsJSON), &defaults); jerr != nil {
			err = fmt.Errorf(T("test_case.import_defaults_invalid"), jerr)
			return
		}
	}

	xlsx, err := excelize.OpenFile(excelPath)
	if err != nil {
		err = fmt.Errorf(T("test_case.import_parse_error"), err)
		return
	}
	sheet := xlsx.GetSheetName(1)
	rows := xlsx.GetRows(sheet)
	if len(rows) < 2 {
		err = fmt.Errorf(T("test_case.import_no_data"))
		return
	}

	// 首行为表头：title → field（同名 title 对应多列时按模板顺序依次匹配，如 Designer）
	header := rows[0]
	fieldCol := map[string]int{} // field → 列索引
	usedTemplate := make([]bool, len(template.Columns))
	for c, title := range header {
		title = strings.TrimSpace(title)
		if title == "" {
			continue
		}
		for i, col := range template.Columns {
			if usedTemplate[i] {
				continue
			}
			if col.Title == title {
				fieldCol[col.Field] = c
				usedTemplate[i] = true
				break
			}
		}
	}

	caseNumCol, hasCaseNum := fieldCol["case_number"]

	// 逐行解析
	for r := 1; r < len(rows); r++ {
		row := rows[r]
		if isBlankRow(row) {
			continue
		}

		var item TestCaseImportItem

		// 编号仅来自 Excel（映射且非空）
		if hasCaseNum && caseNumCol < len(row) {
			item.CaseNumber = strings.TrimSpace(row[caseNumCol])
		}
		if item.CaseNumber == "" {
			err = fmt.Errorf(T("test_case.import_case_number_missing"), r+1)
			return
		}

		// 跳过模板附带的示例行（用例编号以「示例」开头）
		if strings.HasPrefix(item.CaseNumber, sampleCaseNumberPrefix) {
			continue
		}

		// 扩展信息子字段（ext_info.xxx）按列收集，供合并为 YAML
		extInfoSub := map[string]string{}
		for field, col := range fieldCol {
			if strings.HasPrefix(field, "ext_info.") && col < len(row) {
				if v := strings.TrimSpace(row[col]); v != "" {
					extInfoSub[strings.TrimPrefix(field, "ext_info.")] = v
				}
			}
		}

		// 其余字段：Excel 值 > 统一默认值 > 空
		for _, field := range caseImportFields {
			val := ""
			if col, ok := fieldCol[field]; ok && col < len(row) {
				val = strings.TrimSpace(row[col])
			}
			if val == "" {
				val = strings.TrimSpace(defaults[field])
			}
			if field == "ext_info" {
				item.ExtInfo = buildExtInfo(extInfoSub, val)
				continue
			}
			setImportField(&item, field, val)
		}

		// 测试过程：优先提取单元格内嵌图片（标准 / WPS DISPIMG），其次保留单元格文本（path 模式的路径/HTML）
		if tpCol, ok := fieldCol["test_process"]; ok && tpCol < len(row) {
			cell := columnName(tpCol) + strconv.Itoa(r+1)
			txt := strings.TrimSpace(row[tpCol])
			var parts []string
			if imgURL := extractCellImage(xlsx, sheet, cell, item.CaseNumber, item.Module); imgURL != "" {
				item.ImageFiles = append(item.ImageFiles, filepath.Base(imgURL))
				parts = append(parts, fmt.Sprintf(`<img src="%s">`, imgURL))
			}
			// WPS 嵌入单元格图片（=DISPIMG(...)）：测试过程列及其右侧溢出列逐列提取
			for c := tpCol; c < len(row); c++ {
				v := strings.TrimSpace(row[c])
				if imgURL := extractWpsDispimgImage(xlsx, v, item.CaseNumber, item.Module); imgURL != "" {
					item.ImageFiles = append(item.ImageFiles, filepath.Base(imgURL))
					parts = append(parts, fmt.Sprintf(`<img src="%s">`, imgURL))
				}
			}
			// 单元格文本：DISPIMG 公式是图片占位符，已按图片处理，不再当正文；其余（路径/HTML）保留
			if txt != "" && !dispimgRegexp.MatchString(txt) {
				parts = append(parts, txt)
			}
			item.TestProcess = strings.Join(parts, "\n")
		}

		result.Cases = append(result.Cases, item)
	}

	if len(result.Cases) == 0 {
		err = fmt.Errorf(T("test_case.import_no_data"))
		return
	}

	// module 为冲突键之一，必须非空，否则终止报错（补全信息后重新导入）
	for _, item := range result.Cases {
		if strings.TrimSpace(item.Module) == "" {
			err = fmt.Errorf(T("test_case.import_module_required"), item.CaseNumber)
			return
		}
	}

	// 截图：从图片包按「编号_N.format」归并
	if len(strings.TrimSpace(imagePkgPath)) > 0 {
		if ierr := importImagesFromPackage(imagePkgPath, result.Cases); ierr != nil {
			Logger.Error("import case images failed: %s", ierr)
		}
	}

	// 冲突检测：case_number + module
	if err = detectCaseConflicts(&result); err != nil {
		return
	}

	// 保存完整结果到临时文件，供确认阶段读取
	resultFilePath := fmt.Sprintf("%s/import_case_%s.json", DownloadBasePath, importId)
	data, _ := json.Marshal(result)
	if werr := ioutil.WriteFile(resultFilePath, data, 0644); werr != nil {
		err = fmt.Errorf(T("test_case.import_save_error"), werr)
		return
	}

	return
}

// TestCaseImportConfirm 根据用户选择（skip/overwrite/cancel）执行导入或取消
func TestCaseImportConfirm(importId, mode, userName string) (err error) {
	resultFilePath := fmt.Sprintf("%s/import_case_%s.json", DownloadBasePath, importId)

	if mode == "cancel" {
		os.Remove(resultFilePath)
		Logger.Info("Case import cancelled: id=%s, user=%s", importId, userName)
		return
	}

	data, readErr := ioutil.ReadFile(resultFilePath)
	if readErr != nil {
		err = readErr
		Logger.Error("%s", err)
		return
	}

	var result TestCaseImportResult
	if unmarshalErr := json.Unmarshal(data, &result); unmarshalErr != nil {
		err = unmarshalErr
		Logger.Error("unmarshal case import result failed: %s", err)
		return
	}

	isOverwrite := mode == "overwrite"

	// 冲突集合：key = case_number|module
	conflictSet := map[string]bool{}
	for _, c := range result.Conflicts {
		conflictSet[c.CaseNumber+"|"+c.Module] = true
	}

	importedCount := 0
	overwriteCount := 0
	skippedCount := 0

	for _, item := range result.Cases {
		key := item.CaseNumber + "|" + item.Module
		if !conflictSet[key] {
			// 无冲突：直接插入
			if cerr := models.Orm.Table("test_case").Create(&item).Error; cerr != nil {
				Logger.Error("create test_case failed: %s, %s", item.CaseNumber, cerr)
				continue
			}
			importedCount++
		} else if isOverwrite {
			// 覆盖：部分更新，仅写有值的字段（Excel 有值或统一默认值），未映射/空字段不动
			updates := map[string]interface{}{}
			for _, field := range caseImportFields {
				if v := getImportField(&item, field); v != "" {
					updates[field] = v
				}
			}
			if item.TestProcess != "" {
				updates["test_process"] = item.TestProcess
			}
			if len(updates) > 0 {
				if uerr := models.Orm.Table("test_case").Where("case_number = ? AND module = ?", item.CaseNumber, item.Module).Updates(updates).Error; uerr != nil {
					Logger.Error("update test_case failed: %s, %s", item.CaseNumber, uerr)
					continue
				}
			}
			overwriteCount++
		} else {
			skippedCount++
		}
	}

	Logger.Info("Case import done: id=%s, imported=%d, overwritten=%d, skipped=%d", importId, importedCount, overwriteCount, skippedCount)
	os.Remove(resultFilePath)
	return
}

// BuildCaseImportCheckResult 构造结构化导入检查结果，供前端页面渲染
func BuildCaseImportCheckResult(result TestCaseImportResult) map[string]interface{} {
	resp := map[string]interface{}{
		"import_id": result.ImportId,
		"template":  result.Template,
		"conflicts": len(result.Conflicts),
		"new_count": result.NewCount,
	}

	if len(result.Conflicts) > 0 {
		conflicts := make([]map[string]string, 0, len(result.Conflicts))
		for _, c := range result.Conflicts {
			conflicts = append(conflicts, map[string]string{
				"case_number": c.CaseNumber,
				"module":      c.Module,
				"exist_name":  c.ExistName,
				"import_name": c.ImportName,
			})
		}
		resp["conflict_details"] = conflicts
	}
	return resp
}

// detectCaseConflicts 按 case_number + module 检测冲突
func detectCaseConflicts(result *TestCaseImportResult) (err error) {
	if len(result.Cases) == 0 {
		return
	}

	caseNumbers := make([]string, 0, len(result.Cases))
	seen := map[string]bool{}
	for _, item := range result.Cases {
		if !seen[item.CaseNumber] {
			seen[item.CaseNumber] = true
			caseNumbers = append(caseNumbers, item.CaseNumber)
		}
	}

	type existingRow struct {
		CaseNumber string `gorm:"column:case_number"`
		Module     string `gorm:"column:module"`
		CaseName   string `gorm:"column:case_name"`
	}
	var existing []existingRow
	if qerr := models.Orm.Table("test_case").Where("case_number in (?)", caseNumbers).
		Select("case_number, module, case_name").Find(&existing).Error; qerr != nil {
		return qerr
	}

	existMap := map[string]string{} // "case_number|module" → case_name
	for _, e := range existing {
		existMap[e.CaseNumber+"|"+e.Module] = e.CaseName
	}

	for _, item := range result.Cases {
		key := item.CaseNumber + "|" + item.Module
		if existName, ok := existMap[key]; ok {
			result.Conflicts = append(result.Conflicts, ConflictInfo{
				CaseNumber: item.CaseNumber,
				Module:     item.Module,
				ExistName:  existName,
				ImportName: item.CaseName,
			})
		} else {
			result.NewCount++
		}
	}
	return
}

// extractCellImage 提取单元格内嵌图片并保存到模块目录，返回可访问 URL（无图片返回空串）
func extractCellImage(xlsx *excelize.File, sheet, cell, caseNumber, module string) (imgURL string) {
	name, data := xlsx.GetPicture(sheet, cell)
	if len(data) == 0 {
		return ""
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" {
		ext = ".png"
	}
	u, err := saveTestCaseProcessImage(caseNumber, module, ext, data)
	if err != nil {
		Logger.Error("save embedded test_process image failed: %s", err)
		return ""
	}
	return u
}

// dispimgRegexp 匹配 WPS 嵌入单元格图片公式 =DISPIMG("ID_xxx",n)（含 _xlfn.DISPIMG）
var dispimgRegexp = regexp.MustCompile(`(?i)DISPIMG\(\s*"([^"]+)"`)

// 以下结构用于解析 WPS 专有的 xl/cellimages.xml（嵌入单元格图片清单）
type wpsCellImages struct {
	XMLName   xml.Name       `xml:"cellImages"`
	CellImage []wpsCellImage `xml:"cellImage"`
}

type wpsCellImage struct {
	Pic wpsCellImagePic `xml:"pic"`
}

type wpsCellImagePic struct {
	NvPicPr  wpsNvPicPr  `xml:"nvPicPr"`
	BlipFill wpsBlipFill `xml:"blipFill"`
}

type wpsNvPicPr struct {
	CNvPr wpsCNvPr `xml:"cNvPr"`
}

type wpsCNvPr struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type wpsBlipFill struct {
	Blip wpsBlip `xml:"blip"`
}

type wpsBlip struct {
	Embed string `xml:"embed,attr"`
}

type wpsRelationships struct {
	XMLName      xml.Name          `xml:"Relationships"`
	Relationship []wpsRelationship `xml:"Relationship"`
}

type wpsRelationship struct {
	ID     string `xml:"Id,attr"`
	Type   string `xml:"Type,attr"`
	Target string `xml:"Target,attr"`
}

// extractWpsDispimgImage 提取 WPS「嵌入单元格图片」：单元格值为 =DISPIMG("ID_xxx",n) 时，
// 从 xl/cellimages.xml 反查 ID → rId → media 图片，保存到模块目录并返回可访问 URL；非 DISPIMG 返回空串
func extractWpsDispimgImage(xlsx *excelize.File, cellValue, caseNumber, module string) (imgURL string) {
	m := dispimgRegexp.FindStringSubmatch(cellValue)
	if len(m) < 2 {
		return ""
	}
	imgID := m[1]

	rawCellImages := xlsx.XLSX["xl/cellimages.xml"]
	if len(rawCellImages) == 0 {
		return ""
	}
	var cellImages wpsCellImages
	if err := xml.Unmarshal(rawCellImages, &cellImages); err != nil {
		return ""
	}
	embed := ""
	for _, ci := range cellImages.CellImage {
		if ci.Pic.NvPicPr.CNvPr.Name == imgID {
			embed = ci.Pic.BlipFill.Blip.Embed
			break
		}
	}
	if embed == "" {
		return ""
	}

	rawRels := xlsx.XLSX["xl/_rels/cellimages.xml.rels"]
	if len(rawRels) == 0 {
		return ""
	}
	var rels wpsRelationships
	if err := xml.Unmarshal(rawRels, &rels); err != nil {
		return ""
	}
	target := ""
	for _, rel := range rels.Relationship {
		if rel.ID == embed {
			target = rel.Target
			break
		}
	}
	if target == "" {
		return ""
	}

	mediaKey := "xl/" + strings.TrimPrefix(target, "../")
	data := xlsx.XLSX[mediaKey]
	if len(data) == 0 {
		return ""
	}
	ext := strings.ToLower(filepath.Ext(target))
	if ext == "" {
		ext = ".png"
	}
	u, err := saveTestCaseProcessImage(caseNumber, module, ext, data)
	if err != nil {
		Logger.Error("save wps dispimg image failed: %s", err)
		return ""
	}
	return u
}

// importImagesFromPackage 解压图片包，按「编号_N.format」归并到对应用例，写入模块目录并生成 HTML
func importImagesFromPackage(imagePkgPath string, items []TestCaseImportItem) (err error) {
	zr, err := zip.OpenReader(imagePkgPath)
	if err != nil {
		return
	}
	defer zr.Close()

	imageMap := map[string]map[string][]byte{} // case_number → {文件名: 字节}
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		base := filepath.Base(f.Name)
		idx := strings.LastIndex(base, "_")
		if idx <= 0 || idx == len(base)-1 {
			continue
		}
		caseNumber := base[:idx]
		suffix := base[idx+1:]
		if !imageNameRegexp.MatchString(suffix) {
			continue
		}
		data, rerr := readZipFile(f)
		if rerr != nil {
			Logger.Error("read import image failed: %s, %s", base, rerr)
			continue
		}
		if imageMap[caseNumber] == nil {
			imageMap[caseNumber] = map[string][]byte{}
		}
		imageMap[caseNumber][base] = data
	}

	for i := range items {
		files := imageMap[items[i].CaseNumber]
		if len(files) == 0 {
			continue
		}
		names := make([]string, 0, len(files))
		for name := range files {
			names = append(names, name)
		}
		sort.Strings(names)
		moduleDir := sanitizeDirName(items[i].Module)
		var parts []string
		for _, name := range names {
			dst := filepath.Join(UploadBasePath, moduleDir, name)
			if werr := os.WriteFile(dst, files[name], 0644); werr != nil {
				Logger.Error("save import image failed: %s, %s", name, werr)
				continue
			}
			items[i].ImageFiles = append(items[i].ImageFiles, name)
			imgURL := "/uploads/" + url.PathEscape(moduleDir) + "/" + url.PathEscape(name)
			parts = append(parts, fmt.Sprintf(`<img src="%s">`, imgURL))
		}
		items[i].TestProcess = strings.Join(parts, "\n")
	}
	return
}

// readZipFile 读取 zip 内文件内容
func readZipFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// buildExtInfo 合并扩展信息：整字段 YAML（Excel ext_info 或统一默认值）为底，子字段（ext_info.xxx）逐键覆盖
func buildExtInfo(sub map[string]string, whole string) string {
	whole = strings.TrimSpace(whole)
	if len(sub) == 0 {
		return whole
	}
	m := map[string]interface{}{}
	if len(whole) > 0 {
		_ = yaml.Unmarshal([]byte(whole), &m)
	}
	for k, v := range sub {
		m[k] = v
	}
	b, err := yaml.Marshal(m)
	if err != nil {
		return whole
	}
	return string(b)
}

// isBlankRow 判断整行为空
func isBlankRow(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}
