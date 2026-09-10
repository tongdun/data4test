package biz

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"data4test/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gopkg.in/yaml.v2"
)

// GetLocalizedValue 从可能为多语种 JSON 的字段值中按语种取值，回退默认
func GetLocalizedValue(raw, lang string) (value string) {
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 || raw[0] != '{' {
		return raw
	}

	m := make(map[string]string)
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return raw
	}

	for _, k := range []string{lang, "default", "zh-CN"} {
		if v, ok := m[k]; ok && len(strings.TrimSpace(v)) > 0 {
			return v
		}
	}
	// 兜底：按键名排序后取首个非空，保证确定性
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if v := m[k]; len(strings.TrimSpace(v)) > 0 {
			return v
		}
	}
	return ""
}

// GetExtInfoValue 从扩展信息 YAML 取子键，值若为多语种 JSON 则按语种取词
func GetExtInfoValue(extInfo, key, lang string) (value string) {
	extInfo = strings.TrimSpace(extInfo)
	if len(extInfo) == 0 {
		return ""
	}
	m := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(extInfo), &m); err != nil {
		return ""
	}
	v, ok := m[key]
	if !ok {
		return ""
	}
	return GetLocalizedValue(stringifyYamlValue(v), lang)
}

// stringifyYamlValue 将 YAML 值转字符串
func stringifyYamlValue(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", t)
	}
}

// queryTestCaseExportRows 按勾选 ids 或条件筛选查询用例
func queryTestCaseExportRows(ids, product, module, introVersion, caseDesigner, createdAtStart, createdAtEnd string) (rows []TestCaseExportRow, err error) {
	dbHandle := models.Orm.Table("test_case")

	if len(ids) > 0 && ids != "," {
		idList := strings.Split(strings.Trim(ids, ","), ",")
		dbHandle = dbHandle.Where("id in (?)", idList)
	} else {
		if len(product) > 0 {
			dbHandle = dbHandle.Where("product = ?", product)
		}
		if len(module) > 0 {
			moduleList := strings.Split(module, ",")
			dbHandle = dbHandle.Where("module in (?)", moduleList)
		}
		if len(introVersion) > 0 {
			versionList := strings.Split(introVersion, ",")
			dbHandle = dbHandle.Where("intro_version in (?)", versionList)
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
	}

	err = dbHandle.Order("id asc").Find(&rows).Error
	return
}

// fieldValue 取系统列的原始值
func fieldValue(row TestCaseExportRow, field string) string {
	switch field {
	case "case_number":
		return row.CaseNumber
	case "case_name":
		return row.CaseName
	case "case_type":
		return row.CaseType
	case "priority":
		return row.Priority
	case "pre_condition":
		return row.PreCondition
	case "test_range":
		return row.TestRange
	case "test_steps":
		return row.TestSteps
	case "expect_result":
		return row.ExpectResult
	case "auto":
		return row.Auto
	case "scene":
		return row.Scene
	case "fun_developer":
		return row.FunDeveloper
	case "case_designer":
		return row.CaseDesigner
	case "case_executor":
		return row.CaseExecutor
	case "test_time":
		return row.TestTime
	case "test_result":
		return row.TestResult
	case "module":
		return row.Module
	case "intro_version":
		return row.IntroVersion
	case "product":
		return row.Product
	case "remark":
		return row.Remark
	}
	return ""
}

// columnName 列索引(0 起)转 Excel 列名
func columnName(n int) string {
	name := ""
	for n >= 0 {
		name = string(rune('A'+n%26)) + name
		n = n/26 - 1
	}
	return name
}

// imgSrcRegexp 匹配 HTML 中 <img ... src="..."> 的 src 值
var imgSrcRegexp = regexp.MustCompile(`<img[^>]+src\s*=\s*["']([^"']+)["']`)

// extractImgSrcs 从 HTML 提取所有图片 src（测试过程粘贴图片）
func extractImgSrcs(raw string) []string {
	matches := imgSrcRegexp.FindAllStringSubmatch(raw, -1)
	srcs := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
			srcs = append(srcs, strings.TrimSpace(m[1]))
		}
	}
	return srcs
}

// resolveUploadFilePath 将图片引用解析为本地绝对路径：/uploads/xxx 映射到 UploadBasePath 下，其余按相对路径兜底
func resolveUploadFilePath(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	if strings.HasPrefix(ref, "/uploads/") {
		return filepath.Join(UploadBasePath, strings.TrimPrefix(ref, "/uploads/"))
	}
	return filepath.Join(UploadBasePath, ref)
}

// resolveScreenshotValue 处理 test_process 截图列；返回单元格文本与需打包的图片路径
func resolveScreenshotValue(raw, screenshotMode string, xlsxFile *excelize.File, sheet, cell string) (text string, imagePaths []string) {
	if screenshotMode == "" {
		return "", nil
	}
	raw = strings.TrimSpace(raw)
	var paths []string
	if len(raw) > 0 {
		// 优先从 HTML 提取 <img src>（粘贴图片）；兼容旧 JSON 路径数组
		paths = extractImgSrcs(raw)
		if len(paths) == 0 {
			if err := json.Unmarshal([]byte(raw), &paths); err != nil {
				paths = []string{raw}
			}
		}
	}

	if screenshotMode == "embed" {
		if len(paths) > 0 {
			absPath := resolveUploadFilePath(paths[0])
			if _, statErr := os.Stat(absPath); statErr == nil {
				_ = xlsxFile.AddPicture(sheet, cell, absPath, `{"x_scale":0.4,"y_scale":0.4}`)
			}
			if len(paths) > 1 {
				text = strings.Join(paths[1:], "\n")
			}
		}
		return text, nil
	}

	// path 模式：写路径文本并收集图片用于打包
	return strings.Join(paths, "\n"), paths
}

// resolveExportCellValue 按模板列 field 求值
func resolveExportCellValue(row TestCaseExportRow, field, lang, screenshotMode string, xlsxFile *excelize.File, sheet, cell string) (string, []string) {
	switch {
	case strings.HasPrefix(field, "ext_info."):
		return GetExtInfoValue(row.ExtInfo, strings.TrimPrefix(field, "ext_info."), lang), nil
	case field == "ext_info":
		return row.ExtInfo, nil
	case field == "test_process":
		return resolveScreenshotValue(row.TestProcess, screenshotMode, xlsxFile, sheet, cell)
	case field == FieldCaseName || field == FieldCaseModule || field == FieldCaseType ||
		field == FieldPreCondition || field == FieldTestRange || field == FieldTestSteps || field == FieldExpectResult:
		if v := GetCaseLocalized(row.CaseNumber, row.Module, lang, field); v != "" {
			return v, nil
		}
		return fieldValue(row, field), nil
	default:
		return GetLocalizedValue(fieldValue(row, field), lang), nil
	}
}

// ExportTestCase2ExcelByTemplate 按模板导出用例为 Excel
func ExportTestCase2ExcelByTemplate(ids, product, module, introVersion, caseDesigner, createdAtStart, createdAtEnd, templateName, lang, screenshotMode, packFormat string) (fileName string, err error) {
	template, err := GetTestCaseExportTemplateByName(templateName)
	if err != nil {
		return
	}

	rows, err := queryTestCaseExportRows(ids, product, module, introVersion, caseDesigner, createdAtStart, createdAtEnd)
	if err != nil {
		return
	}
	if len(rows) == 0 {
		err = fmt.Errorf(T("error.case_info_not_found"))
		return
	}

	curTime := time.Now().Format("20060102150405")
	baseName := fmt.Sprintf("%s_%s", templateName, curTime)

	xlsxFile := excelize.NewFile()
	sheet := "Sheet1"
	for c, col := range template.Columns {
		xlsxFile.SetCellValue(sheet, columnName(c)+"1", col.Title)
	}

	var imagePaths []string
	for r, row := range rows {
		rowNum := r + 2
		for c, col := range template.Columns {
			cell := columnName(c) + fmt.Sprintf("%d", rowNum)
			text, imgs := resolveExportCellValue(row, col.Field, lang, screenshotMode, xlsxFile, sheet, cell)
			if len(text) > 0 {
				xlsxFile.SetCellValue(sheet, cell, text)
			}
			imagePaths = append(imagePaths, imgs...)
		}
		if screenshotMode == "embed" {
			xlsxFile.SetRowHeight(sheet, rowNum, 80)
		}
	}

	xlsxPath := fmt.Sprintf("%s/%s.xlsx", CaseFilePath, baseName)
	if err = xlsxFile.SaveAs(xlsxPath); err != nil {
		Logger.Error("%s", err)
		return
	}

	if screenshotMode == "embed" || screenshotMode == "" || packFormat == "" {
		return baseName + ".xlsx", nil
	}

	switch packFormat {
	case "zip":
		fileName = baseName + ".zip"
		err = zipExportFiles(xlsxPath, imagePaths, fileName)
	case "tgz":
		fileName = baseName + ".tgz"
		err = tgzExportFiles(xlsxPath, imagePaths, fileName)
	default:
		fileName = baseName + ".xlsx"
	}
	if err == nil && fileName != baseName+".xlsx" {
		_ = os.Remove(xlsxPath)
	}
	return
}

// tgzExportFiles 打包 xlsx 与图片为 tar.gz
func tgzExportFiles(xlsxPath string, imagePaths []string, fileName string) (err error) {
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	fw, err := os.Create(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	gw := gzip.NewWriter(fw)
	tw := tar.NewWriter(gw)

	if err = WriteTarFile(tw, xlsxPath); err != nil {
		Logger.Error("%s", err)
		return
	}
	seen := map[string]bool{}
	for _, p := range imagePaths {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}
		absPath := resolveUploadFilePath(p)
		if seen[absPath] {
			continue
		}
		seen[absPath] = true
		if _, statErr := os.Stat(absPath); statErr == nil {
			if err = WriteTarFile(tw, absPath); err != nil {
				Logger.Error("%s", err)
				return
			}
		}
	}

	if err = tw.Close(); err != nil {
		Logger.Error("%s", err)
		return
	}
	if err = gw.Close(); err != nil {
		Logger.Error("%s", err)
		return
	}
	if err = fw.Close(); err != nil {
		Logger.Error("%s", err)
		return
	}
	return
}

// zipExportFiles 打包 xlsx 与图片为 zip
func zipExportFiles(xlsxPath string, imagePaths []string, fileName string) (err error) {
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	fz, err := os.Create(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	zw := zip.NewWriter(fz)

	addToZip := func(src string) error {
		fr, e := os.Open(src)
		if e != nil {
			return e
		}
		defer fr.Close()
		w, e := zw.Create(filepath.Base(src))
		if e != nil {
			return e
		}
		_, e = io.Copy(w, fr)
		return e
	}

	if err = addToZip(xlsxPath); err != nil {
		Logger.Error("%s", err)
		return
	}
	seen := map[string]bool{}
	for _, p := range imagePaths {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}
		absPath := resolveUploadFilePath(p)
		if seen[absPath] {
			continue
		}
		seen[absPath] = true
		if _, statErr := os.Stat(absPath); statErr == nil {
			if err = addToZip(absPath); err != nil {
				Logger.Error("%s", err)
				return
			}
		}
	}

	if err = zw.Close(); err != nil {
		Logger.Error("%s", err)
		return
	}
	if err = fz.Close(); err != nil {
		Logger.Error("%s", err)
		return
	}
	return
}
