package biz

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"data4test/models"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/url"
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

// resolveUploadFilePath 将图片引用解析为本地绝对路径：/uploads/xxx 映射到 UploadBasePath 下，其余按相对路径兜底。
// 保存时 src 经 url.PathEscape 编码，磁盘文件名为未编码形式，需反向解码后才能定位。
func resolveUploadFilePath(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	if strings.HasPrefix(ref, "/uploads/") {
		rest := strings.TrimPrefix(ref, "/uploads/")
		if q := strings.IndexByte(rest, '?'); q >= 0 {
			rest = rest[:q]
		}
		if decoded, err := url.PathUnescape(rest); err == nil {
			return filepath.Join(UploadBasePath, decoded)
		}
		return filepath.Join(UploadBasePath, rest)
	}
	return filepath.Join(UploadBasePath, ref)
}

// resolveScreenshotValue 处理 test_process 截图列；返回单元格文本与需打包的图片路径。
// col/rowNum 为当前单元格的列下标与行号，用于嵌入模式把多张图依次放到同行右侧单元格。
// embedImgs 收集 WPS 嵌入模式的图片（用于导出后注入 cellimages 部件）。
// excelType 区分嵌入实现：wps 用 DISPIMG 单元格内嵌，office 用 IMAGE() 函数内嵌（需 host 构造图片 URL）。
func resolveScreenshotValue(raw, screenshotMode, excelType, host string, xlsxFile *excelize.File, sheet string, col, rowNum int, embedImgs *[]wpsEmbedImage) (text string, imagePaths []string) {
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
		if excelType == "office" {
			// OFFICE：用 IMAGE() 函数把图片作为单元格值嵌入（Excel 365+），第一张放 test_process 单元格，后续依次放同行右侧单元格
			for i, p := range paths {
				imgURL := toHTTPImageURL(p, host)
				if imgURL == "" {
					continue
				}
				cell := columnName(col+i) + fmt.Sprintf("%d", rowNum)
				// <f> 存 _xlfn.IMAGE（公式），<v> 存 =IMAGE（缓存值）
				xlsxFile.SetCellFormula(sheet, cell, fmt.Sprintf(`_xlfn.IMAGE("%s","",0)`, imgURL))
				xlsxFile.SetCellStr(sheet, cell, fmt.Sprintf(`=IMAGE("%s","",0)`, imgURL))
			}
			return "", nil
		}
		// WPS：嵌入单元格，每张图写一个 DISPIMG 公式，第一张放 test_process 单元格，后续依次放同行右侧单元格
		for i, p := range paths {
			absPath := resolveUploadFilePath(p)
			if _, statErr := os.Stat(absPath); statErr != nil {
				continue
			}
			id := fmt.Sprintf("ID_%08X", len(*embedImgs)+1)
			cell := columnName(col+i) + fmt.Sprintf("%d", rowNum)
			// <f> 存 _xlfn.DISPIMG（公式），<v> 存 =DISPIMG（缓存值，与 WPS 一致）
			xlsxFile.SetCellFormula(sheet, cell, fmt.Sprintf(`_xlfn.DISPIMG("%s",1)`, id))
			xlsxFile.SetCellStr(sheet, cell, fmt.Sprintf(`=DISPIMG("%s",1)`, id))
			*embedImgs = append(*embedImgs, wpsEmbedImage{id: id, absPath: absPath})
		}
		return "", nil
	}

	// path 模式：收集图片用于打包，单元格文本留空（打包图片不写路径文本）
	return "", paths
}

// toHTTPImageURL 将图片引用（/uploads/xxx 或相对路径）转为 http 可访问地址，供 IMAGE() 函数使用。
func toHTTPImageURL(ref, host string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		return ref
	}
	if strings.HasPrefix(ref, "/") {
		return "http://" + host + ref
	}
	return "http://" + host + "/" + ref
}

// resolveExportCellValue 按模板列 field 求值
func resolveExportCellValue(row TestCaseExportRow, field, lang, screenshotMode, excelType, host string, xlsxFile *excelize.File, sheet string, col, rowNum int, embedImgs *[]wpsEmbedImage) (string, []string) {
	switch {
	case strings.HasPrefix(field, "ext_info."):
		return GetExtInfoValue(row.ExtInfo, strings.TrimPrefix(field, "ext_info."), lang), nil
	case field == "ext_info":
		return row.ExtInfo, nil
	case field == "test_process":
		return resolveScreenshotValue(row.TestProcess, screenshotMode, excelType, host, xlsxFile, sheet, col, rowNum, embedImgs)
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
func ExportTestCase2ExcelByTemplate(ids, product, module, introVersion, caseDesigner, createdAtStart, createdAtEnd, templateName, lang, screenshotMode, excelType, packFormat, host string) (fileName string, err error) {
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
	var embedImgs []wpsEmbedImage
	// 语种导出时收集「原始模块目录 → 译文目录」映射，用于打包时翻译图片所在模块目录名
	dirRename := map[string]string{}
	for r, row := range rows {
		rowNum := r + 2
		if lang != "" && lang != "zh-CN" {
			if oldDir := sanitizeDirName(row.Module); oldDir != "" {
				if newName := GetCaseCountLocalized(row.Module, lang); newName != "" && newName != row.Module {
					if newDir := sanitizeDirName(newName); newDir != "" && newDir != oldDir {
						dirRename[oldDir] = newDir
					}
				}
			}
		}
		for c, col := range template.Columns {
			cell := columnName(c) + fmt.Sprintf("%d", rowNum)
			text, imgs := resolveExportCellValue(row, col.Field, lang, screenshotMode, excelType, host, xlsxFile, sheet, c, rowNum, &embedImgs)
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

	if screenshotMode == "embed" && excelType == "wps" {
		if err = injectWpsCellImages(xlsxPath, embedImgs); err != nil {
			Logger.Error("inject wps cell images failed: %s", err)
			return
		}
	}

	if screenshotMode == "embed" || screenshotMode == "" || packFormat == "" {
		return baseName + ".xlsx", nil
	}

	switch packFormat {
	case "zip":
		fileName = baseName + ".zip"
		err = zipExportFiles(xlsxPath, imagePaths, fileName, dirRename)
	case "tgz":
		fileName = baseName + ".tgz"
		err = tgzExportFiles(xlsxPath, imagePaths, fileName, dirRename)
	default:
		fileName = baseName + ".xlsx"
	}
	if err == nil && fileName != baseName+".xlsx" {
		_ = os.Remove(xlsxPath)
	}
	return
}

// imageEntryName 计算图片在归档内的条目名：优先保留 /uploads 下的相对目录（模块目录/文件名），
// 避免不同模块下同名图片在归档内冲突、丢失层级。
func imageEntryName(absPath string, dirRename map[string]string) string {
	rel := filepath.Base(absPath)
	if r, err := filepath.Rel(UploadBasePath, absPath); err == nil &&
		!strings.HasPrefix(r, "..") && !filepath.IsAbs(r) {
		rel = filepath.ToSlash(r)
	}
	// 语种导出时翻译模块目录段（rel 首段为模块目录）
	if len(dirRename) > 0 {
		if i := strings.IndexByte(rel, '/'); i > 0 {
			if newDir, ok := dirRename[rel[:i]]; ok && newDir != "" {
				rel = newDir + rel[i:]
			}
		}
	}
	return rel
}

// writeTarEntry 将文件以指定条目名写入 tar，逻辑同 WriteTarFile，仅条目名可定制。
func writeTarEntry(tw *tar.Writer, filePath, entryName string) (err error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	fr, err := os.Open(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	defer fr.Close()

	h := new(tar.Header)
	h.Name = entryName
	h.Size = fi.Size()
	h.Mode = int64(fi.Mode())
	h.ModTime = fi.ModTime()

	if err = tw.WriteHeader(h); err != nil {
		Logger.Error("%s", err)
		return
	}
	_, err = io.Copy(tw, fr)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	return
}

// tgzExportFiles 打包 xlsx 与图片为 tar.gz
func tgzExportFiles(xlsxPath string, imagePaths []string, fileName string, dirRename map[string]string) (err error) {
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
			if err = writeTarEntry(tw, absPath, imageEntryName(absPath, dirRename)); err != nil {
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
func zipExportFiles(xlsxPath string, imagePaths []string, fileName string, dirRename map[string]string) (err error) {
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	fz, err := os.Create(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	zw := zip.NewWriter(fz)

	addToZip := func(src, entryName string) error {
		fr, e := os.Open(src)
		if e != nil {
			return e
		}
		defer fr.Close()
		w, e := zw.Create(entryName)
		if e != nil {
			return e
		}
		_, e = io.Copy(w, fr)
		return e
	}

	if err = addToZip(xlsxPath, filepath.Base(xlsxPath)); err != nil {
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
			if err = addToZip(absPath, imageEntryName(absPath, dirRename)); err != nil {
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

// wpsEmbedImage 一张待嵌入单元格的图片（WPS DISPIMG）
type wpsEmbedImage struct {
	id      string // DISPIMG ID，如 "ID_00000001"
	absPath string // 本地绝对路径
}

// buildCellImagesXML 生成 WPS xl/cellimages.xml（嵌入单元格图片清单）。
// 结构对齐 WPS 官方写入格式：nvPicPr 含 cNvPicPr，blipFill 含 stretch，spPr 含 xfrm/prstGeom/noFill/ln。
func buildCellImagesXML(entries []wpsCellImageEntry) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	b.WriteString(`<etc:cellImages xmlns:etc="http://www.wps.cn/officeDocument/2017/etCustomData" xmlns:xdr="http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)
	for _, e := range entries {
		widthEMU := e.widthPx * 9525
		heightEMU := e.heightPx * 9525
		fmt.Fprintf(&b, `<etc:cellImage><xdr:pic><xdr:nvPicPr><xdr:cNvPr id="%s" name="%s"/><xdr:cNvPicPr><a:picLocks noChangeAspect="1"/></xdr:cNvPicPr></xdr:nvPicPr><xdr:blipFill><a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></xdr:blipFill><xdr:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom><a:noFill/><a:ln w="9525"><a:noFill/></a:ln></xdr:spPr></xdr:pic></etc:cellImage>`,
			e.id, e.name, e.relID, widthEMU, heightEMU)
	}
	b.WriteString(`</etc:cellImages>`)
	return []byte(b.String())
}

// buildCellImagesRels 生成 WPS xl/_rels/cellimages.xml.rels（rId → media 映射）
func buildCellImagesRels(entries []wpsCellImageEntry) []byte {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	b.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for _, e := range entries {
		fmt.Fprintf(&b, `<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/%s"/>`, e.relID, e.mediaName)
	}
	b.WriteString(`</Relationships>`)
	return []byte(b.String())
}

// wpsCellImageEntry 单张 cellImage 的构造信息
type wpsCellImageEntry struct {
	name      string // cNvPr@name，DISPIMG ID
	id        string // cNvPr@id
	relID     string // blip@r:embed
	mediaName string // media 文件名，如 image1.png
	widthPx   int    // 图片宽度（像素）
	heightPx  int    // 图片高度（像素）
}

// injectWpsCellImages 向已生成的 xlsx 注入 WPS「嵌入单元格图片」部件：
// xl/cellimages.xml、xl/_rels/cellimages.xml.rels、xl/media/imageN.ext，
// 并注册 [Content_Types].xml 与 xl/_rels/workbook.xml.rels。
func injectWpsCellImages(xlsxPath string, images []wpsEmbedImage) (err error) {
	if len(images) == 0 {
		return nil
	}

	r, err := zip.OpenReader(xlsxPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// 读入图片字节并构造 cellImage 条目与 media 文件
	media := make(map[string][]byte) // mediaName -> bytes
	entries := make([]wpsCellImageEntry, 0, len(images))
	for i, img := range images {
		data, rerr := os.ReadFile(img.absPath)
		if rerr != nil {
			Logger.Error("read embed image failed: %s, %s", img.absPath, rerr)
			continue
		}
		ext := strings.ToLower(filepath.Ext(img.absPath))
		if ext == "" {
			ext = ".png"
		}
		widthPx, heightPx := 0, 0
		if cfg, _, derr := image.DecodeConfig(bytes.NewReader(data)); derr == nil {
			widthPx, heightPx = cfg.Width, cfg.Height
		}
		mediaName := fmt.Sprintf("image%d%s", i+1, ext)
		media[mediaName] = data
		entries = append(entries, wpsCellImageEntry{
			name:      img.id,
			id:        fmt.Sprintf("%d", i+1),
			relID:     fmt.Sprintf("rId%d", i+1),
			mediaName: mediaName,
			widthPx:   widthPx,
			heightPx:  heightPx,
		})
	}
	if len(entries) == 0 {
		return nil
	}

	cellimagesXML := buildCellImagesXML(entries)
	cellimagesRels := buildCellImagesRels(entries)

	// 写新 zip：原条目（含补丁后的 content types / workbook rels）+ 新部件
	tmpPath := xlsxPath + ".tmp"
	fz, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(fz)

	copyEntry := func(name string, data []byte) error {
		w, e := zw.Create(name)
		if e != nil {
			return e
		}
		_, e = w.Write(data)
		return e
	}

	for _, f := range r.File {
		data, rerr := readZipFile(f)
		if rerr != nil {
			return rerr
		}
		switch f.Name {
		case "[Content_Types].xml":
			data = patchContentTypes(data)
		case "xl/_rels/workbook.xml.rels":
			data = patchWorkbookRels(data)
		}
		if e := copyEntry(f.Name, data); e != nil {
			return e
		}
	}

	if e := copyEntry("xl/cellimages.xml", cellimagesXML); e != nil {
		return e
	}
	if e := copyEntry("xl/_rels/cellimages.xml.rels", cellimagesRels); e != nil {
		return e
	}
	for name, data := range media {
		if e := copyEntry("xl/media/"+name, data); e != nil {
			return e
		}
	}

	if e := zw.Close(); e != nil {
		return e
	}
	if e := fz.Close(); e != nil {
		return e
	}
	return os.Rename(tmpPath, xlsxPath)
}

// patchContentTypes 在 [Content_Types].xml 中注册 /xl/cellimages.xml
func patchContentTypes(data []byte) []byte {
	override := `<Override PartName="/xl/cellimages.xml" ContentType="application/vnd.wps-officedocument.cellimage+xml"/>`
	return []byte(strings.Replace(string(data), "</Types>", override+"</Types>", 1))
}

// patchWorkbookRels 在 xl/_rels/workbook.xml.rels 中关联 cellimages.xml
func patchWorkbookRels(data []byte) []byte {
	rel := `<Relationship Id="rIdWpsCellImages" Type="http://www.wps.cn/officeDocument/2020/cellImage" Target="cellimages.xml"/>`
	return []byte(strings.Replace(string(data), "</Relationships>", rel+"</Relationships>", 1))
}
