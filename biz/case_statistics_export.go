package biz

import (
	"archive/zip"
	"data4test/models"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

// ExportCaseStatistics2Excel 将选中统计定义覆盖的全量用例按模块导出：一个模块一个 xlsx，最终打包为单个 zip。
// 列结构复用现有导出模板（buildTestCaseXlsx）；数据过滤应用各定义自身的 intro_versions / product。
func ExportCaseStatistics2Excel(idStr, templateName, lang, screenshotMode, excelType, host string) (fileName string, err error) {
	ids := parseIds(idStr)
	if len(ids) == 0 {
		return "", E("common.btn_select_first")
	}

	template, err := GetTestCaseExportTemplateByName(templateName)
	if err != nil {
		return
	}

	var css []CaseStatistics
	models.Orm.Table("case_statistics").Where("id IN (?)", ids).Find(&css)
	if len(css) == 0 {
		return "", E("case_statistics.not_found")
	}

	// module → rows：跨定义去重合并（按 case_number 去重），order 记录模块首次出现顺序
	moduleRows := make(map[string][]TestCaseExportRow)
	moduleSeen := make(map[string]map[string]bool)
	var order []string

	for _, cs := range css {
		def, perr := ParseCaseStatisticsDefinition(cs.Definition)
		if perr != nil {
			Logger.Error("导出统计定义 %d 解析失败: %v", cs.Id, perr)
			continue
		}
		versions := splitTrim(cs.IntroVersions)
		products := splitTrim(cs.Product)
		for _, m := range collectLeafModules(def) {
			q := models.Orm.Table("test_case").Where("deleted_at IS NULL AND module = ?", m)
			if len(versions) > 0 {
				q = q.Where("intro_version IN (?)", versions)
			}
			if len(products) > 0 {
				q = q.Where("product IN (?)", products)
			}
			var rows []TestCaseExportRow
			if qerr := q.Order("id asc").Find(&rows).Error; qerr != nil {
				Logger.Error("导出模块 %s 用例失败: %v", m, qerr)
				continue
			}
			if len(rows) == 0 {
				continue
			}
			if _, ok := moduleSeen[m]; !ok {
				moduleSeen[m] = make(map[string]bool)
				order = append(order, m)
			}
			for _, r := range rows {
				key := strings.TrimSpace(r.CaseNumber)
				if key == "" {
					// 空编号按整行追加（无法去重），避免误丢弃
					moduleRows[m] = append(moduleRows[m], r)
					continue
				}
				if moduleSeen[m][key] {
					continue
				}
				moduleSeen[m][key] = true
				moduleRows[m] = append(moduleRows[m], r)
			}
		}
	}

	if len(order) == 0 {
		return "", E("case_statistics.no_module")
	}

	curTime := time.Now().Format("20060102150405")
	entries := make(map[string]string) // zip 条目名 → 本地绝对路径
	dirRename := map[string]string{}   // 语种导出时模块目录名翻译映射（合并各模块）
	tempXlsx := make([]string, 0, len(order))

	for i, m := range order {
		rows := moduleRows[m]
		if len(rows) == 0 {
			continue
		}
		displayName := m
		if lang != "" && lang != "zh-CN" {
			if v := GetCaseCountLocalized(m, lang); v != "" {
				displayName = v
			}
		}
		entryName := sanitizeDirName(displayName) + ".xlsx"
		// 不同模块清洗后重名时追加序号，避免 zip 内条目覆盖
		if _, exists := entries[entryName]; exists {
			entryName = fmt.Sprintf("%s_%d.xlsx", sanitizeDirName(displayName), i+1)
		}

		tempBase := fmt.Sprintf("case_stat_export_%s_%d", curTime, i+1)
		xlsxPath, imagePaths, embedImgs, dr, berr := buildTestCaseXlsx(template, rows, tempBase, lang, screenshotMode, excelType, host)
		if berr != nil {
			return "", berr
		}
		if screenshotMode == "embed" && excelType == "wps" {
			if ierr := injectWpsCellImages(xlsxPath, embedImgs); ierr != nil {
				Logger.Error("inject wps cell images failed: %s", ierr)
				return "", ierr
			}
		}
		for k, v := range dr {
			dirRename[k] = v
		}

		entries[entryName] = xlsxPath
		tempXlsx = append(tempXlsx, xlsxPath)

		// path 模式：把该模块截图一并打入 zip
		if screenshotMode == "path" {
			for _, p := range imagePaths {
				abs := resolveUploadFilePath(p)
				if _, statErr := os.Stat(abs); statErr == nil {
					entries[imageEntryName(abs, dirRename)] = abs
				}
			}
		}
	}

	fileName = fmt.Sprintf("case_statistics_%s.zip", curTime)
	if err = zipMultipleFiles(fileName, entries); err != nil {
		return "", err
	}
	for _, p := range tempXlsx {
		_ = os.Remove(p)
	}
	return fileName, nil
}

// zipMultipleFiles 将多文件按指定条目名打包为 zip，存到 CaseFilePath/<fileName>。
// entries: zip 条目名(相对路径) → 本地绝对路径；同名条目取先出现者，条目名排序保证输出确定性。
func zipMultipleFiles(fileName string, entries map[string]string) (err error) {
	filePath := fmt.Sprintf("%s/%s", CaseFilePath, fileName)
	fz, err := os.Create(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	defer fz.Close()
	zw := zip.NewWriter(fz)

	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true
		src := entries[name]
		fr, e := os.Open(src)
		if e != nil {
			Logger.Error("open %s failed: %v", src, e)
			continue
		}
		w, e := zw.Create(name)
		if e != nil {
			fr.Close()
			return e
		}
		_, e = io.Copy(w, fr)
		fr.Close()
		if e != nil {
			return e
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
