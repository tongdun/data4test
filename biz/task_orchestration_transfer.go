package biz

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"data4test/models"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v3"
)

// ExportTaskOrchestration 导出编排及其任务完整依赖，返回可下载的归档文件名
func ExportTaskOrchestration(ids, userName string) (fileName string, err error) {
	var idList []string
	for _, id := range strings.Split(ids, ",") {
		id = strings.TrimSpace(id)
		if len(id) > 0 {
			idList = append(idList, id)
		}
	}
	if len(idList) == 0 {
		return "", fmt.Errorf(T("error.orchestration_select_first"))
	}

	var orchestrations []DbTaskOrchestration
	if err = models.Orm.Table("task_orchestration").
		Where("id in (?) and deleted_at is null", idList).
		Order("id asc").
		Find(&orchestrations).Error; err != nil {
		return "", err
	}
	if len(orchestrations) == 0 {
		return "", fmt.Errorf(T("error.orchestration_not_found"), ids)
	}

	// 收集任务名称并解析为任务 id
	nameSet := make(map[string]bool)
	var taskNames []string
	manifest := make([]KOrchestration, 0, len(orchestrations))
	for _, orch := range orchestrations {
		names := ParseTaskList(orch.TaskList)
		for _, name := range names {
			if !nameSet[name] {
				nameSet[name] = true
				taskNames = append(taskNames, name)
			}
		}
		manifest = append(manifest, KOrchestration{
			Name:     orch.OrchestrationName,
			Tasks:    names,
			Products: orch.ProductList,
			Remark:   orch.Remark,
		})
	}

	var tasks []DbSchedule
	if err = models.Orm.Table("schedule").
		Where("task_name in (?) and deleted_at is null", taskNames).
		Find(&tasks).Error; err != nil {
		return "", err
	}
	if len(tasks) != len(taskNames) {
		existing := make(map[string]bool, len(tasks))
		for _, task := range tasks {
			existing[task.TaskName] = true
		}
		for _, name := range taskNames {
			if !existing[name] {
				return "", fmt.Errorf(T("error.orchestration_task_missing"), orchestrations[0].OrchestrationName, name)
			}
		}
	}

	// 依赖闭包校验：任务/场景/数据文件缺失终止导出；产品与环境配置缺失只告警
	if err = validateOrchestrationDependencies(orchestrations, tasks); err != nil {
		return "", err
	}

	taskIds := make([]string, 0, len(tasks))
	for _, task := range tasks {
		taskIds = append(taskIds, task.Id)
	}

	productConfigs, err := collectOrchestrationProductConfigs(orchestrations, tasks)
	if err != nil {
		return "", err
	}
	var extraProducts []string
	for _, cfg := range productConfigs {
		extraProducts = append(extraProducts, cfg.Product.Name)
	}
	files, err := exportSchedule(strings.Join(taskIds, ","), userName, extraProducts, true)
	if err != nil {
		return "", err
	}

	content, err := yaml.Marshal(manifest)
	if err != nil {
		return "", err
	}

	fileName = fmt.Sprintf("orchestration_%s.zip", time.Now().Format("20060102150405"))
	zipPath := filepath.Join(DownloadBasePath, fileName)
	zf, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	entry, err := zw.Create("orchestration_info.yaml")
	if err != nil {
		return "", err
	}
	if _, err = entry.Write(content); err != nil {
		return "", err
	}

	configBytes, err := yaml.Marshal(productConfigs)
	if err != nil {
		return "", err
	}
	configEntry, err := zw.Create("products_info.yaml")
	if err != nil {
		return "", err
	}
	if _, err = configEntry.Write(configBytes); err != nil {
		return "", err
	}

	for _, name := range strings.Split(files, ",") {
		name = strings.TrimSpace(name)
		if len(name) == 0 {
			continue
		}
		data, readErr := ioutil.ReadFile(filepath.Join(DownloadBasePath, name))
		if readErr != nil {
			return "", readErr
		}
		entry, err = zw.Create("bundle/" + name)
		if err != nil {
			return "", err
		}
		if _, err = entry.Write(data); err != nil {
			return "", err
		}
	}

	if err = zw.Close(); err != nil {
		return "", err
	}
	return fileName, nil
}

// ImportTaskOrchestrationCheck 上传检查：解包编排清单，复用任务导入检查，返回冲突信息
func ImportTaskOrchestrationCheck(zipPath, userName string) (data map[string]interface{}, err error) {
	importId := uuid.NewV4().String()
	tempDir := filepath.Join(DownloadBasePath, "orch_import_"+importId)
	defer func() {
		if err != nil {
			os.RemoveAll(tempDir)
		}
	}()

	if err = unzipToDir(zipPath, tempDir); err != nil {
		return nil, err
	}

	manifestPath := filepath.Join(tempDir, "orchestration_info.yaml")
	raw, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf(T("error.orchestration_manifest_missing"))
	}
	var orchestrations []KOrchestration
	if err = yaml.Unmarshal(raw, &orchestrations); err != nil {
		return nil, err
	}
	if len(orchestrations) == 0 {
		return nil, fmt.Errorf(T("error.orchestration_manifest_empty"))
	}

	bundlePath, err := pickTaskBundle(filepath.Join(tempDir, "bundle"))
	if err != nil {
		return nil, err
	}

	scheduleResult, err := ImportScheduleCheck(bundlePath, userName)
	if err != nil {
		return nil, err
	}

	var productConfigs []CliConfig
	if raw, readErr := ioutil.ReadFile(filepath.Join(tempDir, "products_info.yaml")); readErr == nil {
		if err = yaml.Unmarshal(raw, &productConfigs); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(readErr) {
		return nil, readErr
	}

	// 编排同名冲突
	conflicts := make(map[string]string)
	creators := make(map[string]string)
	var names []string
	for _, orch := range orchestrations {
		names = append(names, orch.Name)
	}
	var existing []DbTaskOrchestration
	models.Orm.Table("task_orchestration").
		Select("orchestration_name, user_name").
		Where("orchestration_name in (?) and deleted_at is null", names).
		Find(&existing)
	for _, row := range existing {
		conflicts[row.OrchestrationName] = "orchestration"
		creators[row.OrchestrationName] = row.UserName
	}

	state := OrchestrationImportState{
		ImportId:         importId,
		UserName:         userName,
		TempDir:          tempDir,
		ScheduleImportId: scheduleResult.ImportId,
		Orchestrations:   orchestrations,
		Status:           "pending",
		ProductConfigs:   productConfigs,
	}
	if err = saveOrchestrationImportState(state); err != nil {
		return nil, err
	}

	var conflictHtml []string
	for _, orch := range orchestrations {
		if _, ok := conflicts[orch.Name]; ok {
			conflictHtml = append(conflictHtml, fmt.Sprintf(`<div class="conflict-item"><span class="type-tag">%s</span> %s</div>`, T("task_orchestration.conflict_orchestration"), html.EscapeString(orch.Name)))
		}
	}

	checkConfigs := append([]CliConfig(nil), productConfigs...)
	if envPath := findFileInDir(scheduleResult.StagingDir, "env_config.yaml"); envPath != "" {
		var cfg CliConfig
		if readErr := readCliConfigFile(envPath, &cfg); readErr != nil {
			return nil, readErr
		}
		checkConfigs = append(checkConfigs, cfg)
	}
	metaConflicts := orchestrationMetaConflicts(checkConfigs)
	conflictHtml = append(conflictHtml, metaConflicts...)

	// 任务映射与缺失项
	importedNames := make(map[string]bool, len(scheduleResult.Tasks))
	taskByName := make(map[string]KTask, len(scheduleResult.Tasks))
	for _, task := range scheduleResult.Tasks {
		importedNames[task.Name] = true
		taskByName[task.Name] = task
	}
	importedPlaybooks := make(map[string]bool, len(scheduleResult.Playbooks))
	for _, item := range scheduleResult.Playbooks {
		importedPlaybooks[item.Name] = true
	}
	importedData := make(map[string]bool, len(scheduleResult.DataFiles))
	for _, item := range scheduleResult.DataFiles {
		if len(item.FileName) > 0 {
			importedData[item.FileName] = true
		}
		if len(item.Name) > 0 {
			importedData[item.Name] = true
		}
	}

	details := make([]OrchestrationCheckDetail, 0, len(orchestrations))
	var detailHtml strings.Builder
	blocked := false
	for _, orch := range orchestrations {
		detail := OrchestrationCheckDetail{Name: orch.Name}
		if _, ok := conflicts[orch.Name]; ok {
			detail.Exists = true
		}
		for _, name := range orch.Tasks {
			mapping := OrchestrationTaskMapping{TaskName: name}
			var rows []struct {
				Id        string
				DeletedAt string
			}
			models.Orm.Table("schedule").
				Select("id, deleted_at").
				Where("task_name = ?", name).
				Find(&rows)
			switch {
			case len(rows) > 0 && len(rows[0].DeletedAt) == 0:
				mapping.Source = "existing"
				mapping.TaskId = rows[0].Id
			case importedNames[name] && len(rows) == 0:
				mapping.Source = "import"
			default:
				// 既不在包内、又无可用同号任务（或同名任务已被软删除占用唯一键）
				mapping.Source = "missing"
				detail.Missing = append(detail.Missing, name)
			}
			detail.Tasks = append(detail.Tasks, mapping)

			task, ok := taskByName[name]
			if !ok {
				continue
			}
			for _, ref := range task.PlaybookList {
				ref = strings.TrimSpace(ref)
				if len(ref) == 0 {
					continue
				}
				if task.TaskType == "playbook" || task.TaskType == "scene" {
					if importedPlaybooks[ref] {
						continue
					}
					var count int64
					models.Orm.Table("playbook").Where("name = ? and deleted_at is null", ref).Count(&count)
					if count == 0 {
						detail.MissingScenes = append(detail.MissingScenes, ref)
					}
				} else {
					if importedData[ref] {
						continue
					}
					var count int64
					models.Orm.Table("scene_data").
						Where("(file_name = ? or name = ?) and deleted_at is null", ref, ref).
						Count(&count)
					if count == 0 {
						detail.MissingData = append(detail.MissingData, ref)
					}
				}
			}
		}

		// 产品缺失只提示，不阻断（导入会按包内配置创建）
		productRefs := splitDependencyNames(orch.Products)
		productRefs = append(productRefs, taskProductRefs(taskByName, orch.Tasks)...)
		for _, product := range productRefs {
			var count int64
			models.Orm.Table("product").Where("product = ?", product).Count(&count)
			if count == 0 {
				detail.MissingProducts = append(detail.MissingProducts, product)
			}
		}

		if len(detail.Missing) > 0 || len(detail.MissingScenes) > 0 || len(detail.MissingData) > 0 {
			blocked = true
		}
		details = append(details, detail)

		detailHtml.WriteString(`<div class="conflict-item" style="padding:8px 0">`)
		detailHtml.WriteString(`<div style="font-weight:600;margin-bottom:4px">` + html.EscapeString(detail.Name))
		if detail.Exists {
			detailHtml.WriteString(` <span class="type-tag">` + T("task_orchestration.check_existing") + `</span>`)
		}
		detailHtml.WriteString(`</div>`)
		for _, mapping := range detail.Tasks {
			label := T("task_orchestration.check_import")
			if mapping.Source == "existing" {
				label = T("task_orchestration.check_existing") + " #" + mapping.TaskId
			} else if mapping.Source == "missing" {
				label = `<span style="color:#c0392b">` + T("task_orchestration.check_missing") + `</span>`
			}
			detailHtml.WriteString(fmt.Sprintf(`<div style="font-size:13px;color:#555;margin-left:12px">- %s：%s</div>`, html.EscapeString(mapping.TaskName), label))
		}
		for _, name := range detail.MissingScenes {
			detailHtml.WriteString(fmt.Sprintf(`<div style="font-size:13px;color:#c0392b;margin-left:12px">- %s：%s</div>`, T("task_orchestration.check_missing_scene"), html.EscapeString(name)))
		}
		for _, name := range detail.MissingData {
			detailHtml.WriteString(fmt.Sprintf(`<div style="font-size:13px;color:#c0392b;margin-left:12px">- %s：%s</div>`, T("task_orchestration.check_missing_data"), html.EscapeString(name)))
		}
		for _, name := range detail.MissingProducts {
			detailHtml.WriteString(fmt.Sprintf(`<div style="font-size:13px;color:#b8860b;margin-left:12px">- %s：%s</div>`, T("task_orchestration.check_missing_product"), html.EscapeString(name)))
		}
		detailHtml.WriteString(`</div>`)
	}

	for i := range metaConflicts {
		conflicts[fmt.Sprintf("meta:%d", i)] = "metadata"
	}
	return map[string]interface{}{
		"import_id":                    importId,
		"orchestrations":               orchestrations,
		"orchestration_conflicts":      conflicts,
		"orchestration_conflicts_html": strings.Join(conflictHtml, ""),
		"orchestration_details":        details,
		"orchestration_details_html":   detailHtml.String(),
		"blocked":                      blocked,
		"schedule":                     BuildImportCheckResult(scheduleResult),
	}, nil
}

// taskProductRefs 收集编排内各任务引用的产品
func taskProductRefs(taskByName map[string]KTask, taskNames []string) []string {
	seen := make(map[string]bool)
	var refs []string
	for _, name := range taskNames {
		task, ok := taskByName[name]
		if !ok || task.ProductList == nil {
			continue
		}
		for _, product := range splitDependencyNames(*task.ProductList) {
			if seen[product] {
				continue
			}
			seen[product] = true
			refs = append(refs, product)
		}
	}
	return refs
}

// ImportTaskOrchestrationConfirm 确认导入：先落任务依赖，再按名称映射创建/更新编排
func ImportTaskOrchestrationConfirm(importId, mode, userName string) (err error) {
	state, err := loadOrchestrationImportState(importId)
	if err != nil {
		return err
	}
	// 防止伪造的 state 文件造成任意目录删除或伪造依赖导入
	if err = ValidateImportId(state.ScheduleImportId); err != nil {
		return err
	}
	if !withinDownloadBase(state.TempDir) {
		return fmt.Errorf(T("error.orchestration_archive_unsafe"))
	}
	defer func() {
		if err == nil {
			state.Status = "done"
		} else {
			state.Status = "error"
			state.Message = err.Error()
		}
		_ = saveOrchestrationImportState(state)
	}()
	state.Status = "processing"
	state.Message = ""
	_ = saveOrchestrationImportState(state)
	defer os.RemoveAll(state.TempDir)

	if err = ImportScheduleConfirm(state.ScheduleImportId, mode, userName); err != nil {
		return err
	}

	if mode == "cancel" {
		return nil
	}
	for _, cfg := range state.ProductConfigs {
		if err = upsertCliConfigMode(cfg, cfg.SysParameter, mode); err != nil {
			return err
		}
	}

	// 多个编排互相独立：单个失败只回滚本条，不影响其他编排
	var failures []string
	for _, k := range state.Orchestrations {
		if e := importSingleOrchestration(k, mode, userName); e != nil {
			failures = append(failures, fmt.Sprintf("%s: %s", k.Name, e))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("%s", strings.Join(failures, "; "))
	}
	return nil
}

// importSingleOrchestration 导入单个编排；失败时只回滚本条新建的记录
func importSingleOrchestration(k KOrchestration, mode, userName string) (err error) {
	names := make([]string, 0, len(k.Tasks))
	for _, name := range k.Tasks {
		name = strings.TrimSpace(name)
		if len(name) > 0 {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return fmt.Errorf(T("error.orchestration_no_tasks"), k.Name)
	}

	var createdId string
	defer func() {
		if err != nil && len(createdId) > 0 {
			models.Orm.Table("task_orchestration").Where("id = ?", createdId).Delete(&DbTaskOrchestration{})
		}
	}()

	var tasks []DbSchedule
	if err = models.Orm.Table("schedule").
		Select("id, task_name").
		Where("task_name in (?) and deleted_at is null", names).
		Find(&tasks).Error; err != nil {
		return err
	}
	existing := make(map[string]bool, len(tasks))
	for _, task := range tasks {
		existing[task.TaskName] = true
	}
	for _, name := range names {
		if !existing[name] {
			return fmt.Errorf(T("error.orchestration_task_missing"), k.Name, name)
		}
	}

	var current DbTaskOrchestration
	models.Orm.Table("task_orchestration").
		Where("orchestration_name = ? and deleted_at is null", k.Name).
		Find(&current)

	switch {
	case len(current.Id) > 0 && mode == "skip":
		// 同名编排内容不一致时静默跳过会丢弃包内步骤，需让用户显式选择覆盖
		if NormalizeTaskList(current.TaskList) != strings.Join(names, "\n") {
			return fmt.Errorf(T("error.orchestration_skip_conflict"), k.Name)
		}
		return nil
	case len(current.Id) > 0 && mode == "overwrite":
		return models.Orm.Table("task_orchestration").Where("id = ?", current.Id).
			Updates(map[string]interface{}{
				"task_list":    strings.Join(names, "\n"),
				"product_list": k.Products,
				"remark":       k.Remark,
				"updated_at":   time.Now().Format(baseFormat),
			}).Error
	default:
		orch := TaskOrchestration{
			OrchestrationName:   k.Name,
			TaskList:            strings.Join(names, "\n"),
			ProductList:         k.Products,
			OrchestrationStatus: OrchestrationStatusNotStarted,
			Remark:              k.Remark,
			UserName:            userName,
			CreatedAt:           time.Now().Format(baseFormat),
			UpdatedAt:           time.Now().Format(baseFormat),
		}
		if err = models.Orm.Table("task_orchestration").Omit("deleted_at", "last_at").Create(&orch).Error; err != nil {
			return err
		}
		var saved DbTaskOrchestration
		models.Orm.Table("task_orchestration").
			Where("orchestration_name = ? and deleted_at is null", k.Name).Find(&saved)
		createdId = saved.Id
		return nil
	}
}

func orchestrationImportStatePath(importId string) string {
	return filepath.Join(DownloadBasePath, "orchestration_import_"+importId+".json")
}

func saveOrchestrationImportState(state OrchestrationImportState) error {
	if err := ValidateImportId(state.ImportId); err != nil {
		return err
	}
	content, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(orchestrationImportStatePath(state.ImportId), content, 0644)
}

func loadOrchestrationImportState(importId string) (state OrchestrationImportState, err error) {
	if err = ValidateImportId(importId); err != nil {
		return state, err
	}
	content, err := ioutil.ReadFile(orchestrationImportStatePath(importId))
	if err != nil {
		return state, fmt.Errorf(T("error.orchestration_import_not_found"), importId)
	}
	err = json.Unmarshal(content, &state)
	return
}

// importIdPattern 导入任务 id 只允许 UUID 格式，杜绝路径穿越
var importIdPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ValidateImportId 校验导入任务 id 格式
func ValidateImportId(importId string) error {
	if !importIdPattern.MatchString(strings.TrimSpace(importId)) {
		return fmt.Errorf(T("error.orchestration_import_not_found"), importId)
	}
	return nil
}

// withinDownloadBase 判断路径是否位于 DownloadBasePath 目录内
func withinDownloadBase(target string) bool {
	target = strings.TrimSpace(target)
	if len(target) == 0 {
		return false
	}
	base := filepath.Clean(DownloadBasePath)
	clean := filepath.Clean(target)
	return clean == base || strings.HasPrefix(clean, base+string(os.PathSeparator))
}

// GetTaskOrchestrationImportStatus 查询导入任务状态
func GetTaskOrchestrationImportStatus(importId string) (status, message string, err error) {
	state, err := loadOrchestrationImportState(importId)
	if err != nil {
		return "", "", err
	}
	return state.Status, state.Message, nil
}

// unzipToDir 解压 zip 到指定目录（防止路径穿越，并限制文件数与总大小）
func unzipToDir(zipPath, destDir string) error {
	data, err := ioutil.ReadFile(zipPath)
	if err != nil {
		return err
	}
	return extractBoundedZip(data, destDir, archiveMaxFiles, archiveMaxBytes)
}

// pickTaskBundle 选择真正包含任务清单的依赖包（导出包中可能同时包含附件包与数据包）
func pickTaskBundle(bundleDir string) (string, error) {
	files, err := ioutil.ReadDir(bundleDir)
	if err != nil {
		return "", fmt.Errorf(T("error.orchestration_manifest_missing"))
	}

	var candidates []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.HasSuffix(name, ".tgz") || strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".zip") {
			candidates = append(candidates, filepath.Join(bundleDir, name))
		}
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf(T("error.orchestration_manifest_missing"))
	}

	for _, path := range candidates {
		if archiveHasTaskInfo(path) {
			return path, nil
		}
	}
	return "", fmt.Errorf(T("error.orchestration_manifest_missing"))
}

// archiveHasTaskInfo 判断归档中是否包含 task_info.yaml
func archiveHasTaskInfo(path string) bool {
	if strings.HasSuffix(path, ".zip") {
		reader, err := zip.OpenReader(path)
		if err != nil {
			return false
		}
		defer reader.Close()
		if len(reader.File) > archiveMaxFiles {
			return false
		}
		var size uint64
		for _, f := range reader.File {
			if f.UncompressedSize64 > uint64(archiveMaxBytes)-size {
				return false
			}
			size += f.UncompressedSize64
		}
		for _, f := range reader.File {
			if filepath.Base(f.Name) == "task_info.yaml" {
				return true
			}
		}
		return false
	}

	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	var reader io.Reader = f
	if strings.HasSuffix(path, ".tgz") || strings.HasSuffix(path, ".tar.gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return false
		}
		defer gz.Close()
		reader = gz
	}

	tr := tar.NewReader(io.LimitReader(reader, archiveMaxBytes+archiveMaxFiles*1024+1024))
	var size int64
	for count := 0; count < archiveMaxFiles; count++ {
		header, err := tr.Next()
		if err != nil {
			return false
		}
		if header.Size < 0 || header.Size > archiveMaxBytes-size {
			return false
		}
		size += header.Size
		if filepath.Base(header.Name) == "task_info.yaml" {
			return true
		}
	}
	return false
}

// validateOrchestrationDependencies 导出前依赖校验：
// 任务、场景、数据文件缺失直接终止；产品、环境配置缺失仅告警。
func validateOrchestrationDependencies(orchestrations []DbTaskOrchestration, tasks []DbSchedule) error {
	ownerOf := func(taskName string) string {
		for _, orch := range orchestrations {
			for _, name := range ParseTaskList(orch.TaskList) {
				if name == taskName {
					return orch.OrchestrationName
				}
			}
		}
		return orchestrations[0].OrchestrationName
	}

	sceneOwners := make(map[string]string)
	dataOwners := make(map[string]string)
	productSet := make(map[string]bool)

	for _, orch := range orchestrations {
		for _, p := range splitDependencyNames(orch.ProductList) {
			productSet[p] = true
		}
	}
	for _, task := range tasks {
		owner := ownerOf(task.TaskName)
		for _, p := range splitDependencyNames(task.ProductList) {
			productSet[p] = true
		}
		switch task.TaskType {
		case "scene":
			for _, name := range splitDependencyNames(task.SceneList) {
				sceneOwners[name] = owner
			}
		case "data":
			for _, name := range splitDependencyNames(task.DataList) {
				dataOwners[name] = owner
			}
		}
	}

	if len(sceneOwners) > 0 {
		names := mapKeys(sceneOwners)
		var scenes []DbScene
		if err := models.Orm.Table("playbook").
			Where("name in (?) and deleted_at is null", names).
			Find(&scenes).Error; err != nil {
			return err
		}
		have := make(map[string]bool)
		for _, scene := range scenes {
			have[scene.Name] = true
			for _, data := range splitDependencyNames(scene.DataFileList) {
				dataOwners[data] = sceneOwners[scene.Name]
			}
		}
		for _, name := range names {
			if !have[name] {
				return fmt.Errorf(T("error.orchestration_scene_missing"), sceneOwners[name], name)
			}
		}
	}

	if len(dataOwners) > 0 {
		names := mapKeys(dataOwners)
		var existing []string
		if err := models.Orm.Table("scene_data").
			Where("file_name in (?) and deleted_at is null", names).
			Pluck("file_name", &existing).Error; err != nil {
			return err
		}
		have := sliceToSet(existing)
		for _, name := range names {
			if !have[name] {
				return fmt.Errorf(T("error.orchestration_data_missing"), dataOwners[name], name)
			}
		}
	}

	// 产品：缺失只告警
	if len(productSet) > 0 {
		products := setKeys(productSet)
		var rows []Product
		models.Orm.Table("product").Where("product in (?)", products).Find(&rows)
		found := make(map[string]Product, len(rows))
		for _, row := range rows {
			found[row.Name] = row
		}
		var apps []string
		for _, name := range products {
			product, ok := found[name]
			if !ok {
				Logger.Warning(T("warn.orchestration_product_missing"), name)
				continue
			}
			for _, app := range splitDependencyNames(product.Apps) {
				apps = append(apps, app)
			}
		}

		// 环境配置：按 app 判断，缺失只告警
		if len(apps) > 0 {
			var existingApps []string
			models.Orm.Table("env_config").Where("app in (?)", apps).Pluck("app", &existingApps)
			have := sliceToSet(existingApps)
			for _, app := range apps {
				if !have[app] {
					Logger.Warning(T("warn.orchestration_env_missing"), app)
				}
			}
		}
	}

	return nil
}

func splitDependencyNames(value string) []string {
	seen := make(map[string]bool)
	var names []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if len(item) == 0 || seen[item] {
			continue
		}
		seen[item] = true
		names = append(names, item)
	}
	return names
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func setKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func sliceToSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, item := range items {
		set[item] = true
	}
	return set
}

// collectOrchestrationProductConfigs includes task products and orchestration-only overrides.
func collectOrchestrationProductConfigs(orchestrations []DbTaskOrchestration, tasks []DbSchedule) ([]CliConfig, error) {
	names := make(map[string]bool)
	for _, orch := range orchestrations {
		for _, name := range splitDependencyNames(orch.ProductList) {
			names[name] = true
		}
	}
	for _, task := range tasks {
		for _, name := range splitDependencyNames(task.ProductList) {
			names[name] = true
		}
	}
	keys := setKeys(names)
	sort.Strings(keys)
	var configs []CliConfig
	for _, name := range keys {
		var products []Product
		if err := models.Orm.Table("product").Where("product = ?", name).Find(&products).Error; err != nil {
			return nil, err
		}
		if len(products) == 0 {
			continue
		}
		product := products[0]
		apps := make(map[string]bool)
		for _, app := range splitDependencyNames(product.Apps) {
			apps[app] = true
		}
		raw, err := buildCliConfigYaml(name, apps, nil, nil)
		if err != nil {
			return nil, err
		}
		var cfg CliConfig
		if err = yaml.Unmarshal(raw, &cfg); err != nil {
			return nil, err
		}
		configs = append(configs, cfg)
	}
	return configs, nil
}

func orchestrationMetaConflicts(configs []CliConfig) []string {
	var conflicts []string
	seen := make(map[string]bool)
	for _, cfg := range configs {
		objects := map[string][]string{"product": {cfg.Product.Name}}
		for app := range cfg.App {
			objects["env_config"] = append(objects["env_config"], app)
		}
		for name := range cfg.SysParameter {
			objects["sys_parameter"] = append(objects["sys_parameter"], name)
		}
		for name := range cfg.AssertTemplates {
			objects["assert_template"] = append(objects["assert_template"], name)
		}
		for table, names := range objects {
			column := "name"
			if table == "product" {
				column = "product"
			}
			if table == "env_config" {
				column = "app"
			}
			for _, name := range names {
				if name == "" {
					continue
				}
				key := table + ":" + name
				if seen[key] {
					continue
				}
				seen[key] = true
				var count int
				models.Orm.Table(table).Where(column+" = ?", name).Count(&count)
				if count > 0 {
					conflicts = append(conflicts, `<div class="conflict-item"><span class="type-tag">`+html.EscapeString(table)+`</span> `+html.EscapeString(name)+`</div>`)
				}
			}
		}
	}
	sort.Strings(conflicts)
	return conflicts
}
