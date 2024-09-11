package biz

import (
	"archive/tar"
	"compress/gzip"
	"data4perf/models"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var CronHandle = NewCrontab()

func Init() {
	ServiceRestart()
}

func ServiceRestart() {
	ids := GetRunningTask()
	_ = RunTask(ids, "init")
}

// NewCrontab new crontab
func NewCrontab() *Crontab {
	return &Crontab{
		inner: cron.New(),
		ids:   make(map[string]cron.EntryID),
	}
}

func (r *RunSchecdule) Run() {
	OneTask(r.Id)
	CronHandle.UpdateTime(r.Id)
}

func (c *Crontab) UpdateTime(id string) (err error) {
	var dbSchedule DbSchedule
	if value, ok := c.ids[id]; ok {
		lastTime, nextTime := c.GetRunTime(value)
		models.Orm.Table("schedule").Where("id = ?", id).Find(&dbSchedule)
		if len(dbSchedule.TaskName) == 0 {
			err = fmt.Errorf("未找到场景信息，请核对: %v", id)
			Logger.Error("%s", err)
			return
		}

		if !lastTime.IsZero() {
			dbSchedule.LastAt = lastTime.Format("2006-01-02 15:04:05")
		}

		if !nextTime.IsZero() {
			dbSchedule.NextAt = nextTime.Format("2006-01-02 15:04:05")
		}

		err = models.Orm.Table("schedule").Where("id = ?", id).Update(&dbSchedule).Error
		if err != nil {
			Logger.Error("%s", err)
		}
	}
	return
}

func (c *Crontab) GetRunTime(id cron.EntryID) (lastTime, nextTime time.Time) {
	lastTime = c.inner.Entry(id).Prev
	nextTime = c.inner.Entry(id).Next
	return
}

// IDs ...
func (c *Crontab) IDs() map[string]cron.EntryID {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	validIDs := make([]string, 0, len(c.ids))
	invalidIDs := make([]string, 0)
	for sid, eid := range c.ids {
		if e := c.inner.Entry(eid); e.ID != eid {
			invalidIDs = append(invalidIDs, sid)
			continue
		}
		validIDs = append(validIDs, sid)
	}
	for _, id := range invalidIDs {
		delete(c.ids, id)
	}

	return c.ids
}

// Start start the crontab engine
func (c *Crontab) Start() {
	c.inner.Start()
}

// Stop stop the crontab engine
func (c *Crontab) Stop() {
	c.inner.Stop()
}

// DelByID remove one crontab task
func (c *Crontab) DelByID(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	eid, ok := c.ids[id]
	if !ok {
		return
	}
	c.inner.Remove(eid)
	delete(c.ids, id)
}

// AddByID add one crontab task
// id is unique
// spec is the crontab expression
func (c *Crontab) AddByID(id string, spec string, cmd cron.Job) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.ids[id]; ok {
		err = errors.Errorf("crontab id exists")
		return
	}
	eid, err := c.inner.AddJob(spec, cmd)
	if err != nil {
		return
	}
	c.ids[id] = eid
	return
}

// AddByFunc add function as crontab task
func (c *Crontab) AddByFunc(id string, spec string, f func()) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.ids[id]; ok {
		return errors.Errorf("crontab id exists")
	}
	eid, err := c.inner.AddFunc(spec, f)
	if err != nil {
		return err
	}
	c.ids[id] = eid
	return nil
}

// IsExists check the crontab task whether existed with job id
func (c *Crontab) IsExists(jid string) bool {
	_, exist := c.ids[jid]
	return exist
}

func OneTask(id string) (err error) {
	var dbSchedule DbSchedule
	models.Orm.Table("schedule").Where("id = ?", id).Find(&dbSchedule)
	if len(dbSchedule.TaskName) == 0 {
		return
	}

	productTaskList, _ := GetProductInfo(dbSchedule.ProductList)
	productTaskInfo := productTaskList[0]

	var err1 error
	dataList, _ := dbSchedule.GetDataIds()
	sceneList, _, _ := dbSchedule.GetSceneIds()
	if dbSchedule.TaskType == "data" {
		for _, dataId := range dataList {
			err1 = RunSceneDataOnce(dataId, dbSchedule.ProductList, "task")
			if err1 != nil {
				err = err1
				Logger.Error("%s", err)
				return
			}
		}
	} else if dbSchedule.TaskType == "scene" {
		for _, sceneId := range sceneList {
			playbookInfo, productList, err := GetPlRunInfo("task", sceneId)
			playbook := playbookInfo.GetPlaybook()
			productSceneInfo := productList[0]
			if len(dbSchedule.ProductList) == 0 {
				_, _, err1 = playbook.RunPlaybook(playbookInfo.Id, "start", "task", productSceneInfo)
			} else {
				_, _, err1 = playbook.RunPlaybook(playbookInfo.Id, "start", "task", productTaskInfo)
			}

			if err1 != nil {
				err = err1
				Logger.Error("%s", err)
				break
			}
		}
	}

	return
}

func GetRunningTask() (ids string) {
	var idList []string
	status := "running"
	models.Orm.Table("schedule").Where("task_status = ?", status).Pluck("id", &idList)
	if len(idList) == 0 {
		return
	} else {
		for _, item := range idList {
			if len(ids) == 0 {
				ids = item
			} else {
				ids = ids + "," + item
			}
		}
	}
	return
}

func GetTaskCron(id, mode string) (cronStr string, dbSchedule DbSchedule, err error) {
	curTime := fmt.Sprintf(time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
	idInt, _ := strconv.Atoi(id)
	models.Orm.Table("schedule").Where("id = ?", idInt).Find(&dbSchedule)
	if len(dbSchedule.TaskName) == 0 {
		err = errors.Errorf("未对找对应任务信息，请核查: %s", dbSchedule.TaskName)
		return
	}

	if mode != "init" {
		if dbSchedule.TaskStatus == "running" {
			err = fmt.Errorf("任务已在运行中")
			return
		}
	}

	switch dbSchedule.TaskMode {
	case "cron":
		cronStr = dbSchedule.Crontab
	case "once":
		cronStr = "once"
		dbSchedule.LastAt = curTime
	case "day":
		cronStr = fmt.Sprintf("0 %s * * *", dbSchedule.Time4day)
	case "week":
		cronStr = fmt.Sprintf("0 %s * * %s", dbSchedule.Time4week, dbSchedule.Week)
	}
	dbSchedule.TaskStatus = "running"
	return
}

func (s DbSchedule) GetDataIds() (ids []string, err error) {
	dataList := GetListFromHtml(s.DataList)

	for _, dataItem := range dataList {
		var dbSceneData DbSceneData
		if len(dataItem) == 0 {
			continue
		}
		models.Orm.Table("scene_data").Where("name = ?", dataItem).Find(&dbSceneData)

		if len(dbSceneData.Name) == 0 {
			err1 := errors.Errorf("未找到名称[%s]的数据，请核对", dataItem)
			err = err1
			Logger.Error("%s", err)
			return
		}
		ids = append(ids, dbSceneData.Id)
	}
	return
}

func (s DbSchedule) GetSceneIds() (ids []string, sceneTypes []int, err error) {
	sceneList := GetListFromHtml(s.SceneList)

	var dbPlaybooks []DbScene
	models.Orm.Table("playbook").Where("name in (?)", sceneList).Order("priority").Find(&dbPlaybooks)

	for _, sceneItem := range dbPlaybooks {
		ids = append(ids, sceneItem.Id)
		sceneTypes = append(sceneTypes, sceneItem.SceneType)
	}
	//Logger.Debug("Run Order: %v", ids)
	//for _, sceneItem := range sceneList {
	//	if len(sceneItem) == 0 {
	//		continue
	//	}
	//	var dbPlaybook DbScene
	//	models.Orm.Table("playbook").Where("name = ?", sceneItem).Find(&dbPlaybook)
	//	ids = append(ids,dbPlaybook.Id)
	//	sceneTypes = append(sceneTypes, dbPlaybook.SceneType)
	//}

	return
}

func (s DbSchedule) UpdateTaskStatus(status string) (err error) {
	s.TaskStatus = status
	if status == "running" {
		curTime := time.Now().Format("2006-01-02 15:04:05")
		s.LastAt = curTime
	}

	err = models.Orm.Table("schedule").Where("id = ?", s.Id).Update(&s).Error
	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func StopTask(ids string) (err error) {
	idList := strings.Split(ids, ",")
	for _, id := range idList {
		if len(id) == 0 {
			continue
		}
		var dbSchedule DbSchedule
		models.Orm.Table("schedule").Where("id = ?", id).Find(&dbSchedule)
		if dbSchedule.TaskMode != "once" {
			CronHandle.DelByID(id)
		}
		dbSchedule.UpdateTaskStatus("stopped")
		Logger.Info("任务已暂停：%v", dbSchedule.TaskName)
	}
	return
}

func RunOnceTask(id string) (err error) {
	var task DbSchedule
	models.Orm.Table("schedule").Where("id = ?", id).Find(&task)
	if len(task.TaskName) == 0 {
		err = errors.Errorf("未对找对应任务信息，请核查: %s", task.TaskName)
		return
	}
	switch task.TaskType {
	case "data":
		dataIds, _ := task.GetDataIds()
		for _, dataId := range dataIds {
			err1 := RepeatRunDataFile(dataId, task.ProductList, "task")
			if err1 != nil {
				err = err1
				Logger.Error("%v", err)
				break // 数据执行完后，后续数据不再执行，后续可以做错误精细化管理
			}
		}
	case "scene":
		sceneIds, _, _ := task.GetSceneIds()
		for _, sceneId := range sceneIds {
			var err1 error
			err1 = RunPlaybookFromMgmt(sceneId, "start", task.ProductList, "task")
			if err1 != nil {
				if err != nil {
					err = fmt.Errorf("%v; %v", err, err1)
				} else {
					err = err1
				}
				//break // 场景执行完后，后续场景不再执行，后续可以做错误精细化管理
			}
		}
	}
	task.TaskStatus = "finished"
	err = models.Orm.Table("schedule").Where("id = ?", task.Id).Update(&task).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func RunTask(ids, mode string) (err error) {
	idList := strings.Split(ids, ",")
	for _, id := range idList {
		if len(id) == 0 {
			continue
		}

		cronStr, retSchedule, err1 := GetTaskCron(id, mode)
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			continue
		}
		// 更新运行中任务的状态和时间
		err1 = models.Orm.Table("schedule").Where("id = ?", retSchedule.Id).Update(&retSchedule).Error
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}

		if cronStr == "once" {
			Logger.Info("添加单次任务:[%s]", retSchedule.TaskName)
			go func(id string) {
				RunOnceTask(id)
			}(id)
		} else {
			task := &RunSchecdule{id}
			isExist := CronHandle.IsExists(id)
			if isExist {
				err1 := errors.Errorf("任务[%s]已添加", retSchedule.TaskName)
				err = err1
				Logger.Warning("%s", err)
			} else {
				Logger.Info("添加定时任务:[%s]", retSchedule.TaskName)
				err1 := CronHandle.AddByID(id, cronStr, task)
				if err1 != nil {
					err = errors.Errorf("添加任务[%s]失败: %s", retSchedule.TaskName, err1)
					Logger.Error("%s", err)
				}
			}
		}
	}

	CronHandle.Start()
	return
}

func CopySchedule(id, userName string) (err error) {
	var dbSchedule DbSchedule
	models.Orm.Table("schedule").Where("id = ?", id).Find(&dbSchedule)
	if len(dbSchedule.TaskName) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", id)
		Logger.Error("%s", err)
		return
	}

	var schedule Schedule4Copy
	schedule.TaskName = fmt.Sprintf("%s_复制", dbSchedule.TaskName)
	schedule.TaskType = dbSchedule.TaskType
	schedule.SceneList = dbSchedule.SceneList
	schedule.DataList = dbSchedule.DataList
	schedule.Crontab = dbSchedule.Crontab
	schedule.Week = dbSchedule.Week
	schedule.Time4day = dbSchedule.Time4day
	schedule.Time4week = dbSchedule.Time4week
	schedule.Threading = dbSchedule.Threading
	schedule.Remark = dbSchedule.Remark
	schedule.TaskMode = dbSchedule.TaskMode
	schedule.TaskStatus = "not_started"
	schedule.UserName = userName
	schedule.DataNumber = dbSchedule.DataNumber
	schedule.SceneNumber = dbSchedule.SceneNumber
	schedule.ProductList = dbSchedule.ProductList

	err = models.Orm.Table("schedule").Create(schedule).Error
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func ExportSchedule(id, userName string) (fileName string, err error) {
	sqlRule1 := "# SQL生成规则: "
	sqlRule2 := "# (1)产品配置/应用配置，存在则跳过; 其他不存在则插入，存在则更新"
	sqlRule3 := "# (2)创建人统一为导出操作人的信息，所属产品或关联产品统一为导出任务第一个产品的信息"

	curTime := time.Now().Format("20060102150405")
	fileName = fmt.Sprintf("%s_%s.sql", "数据SQL", curTime)
	importFileName := fmt.Sprintf("%s_%s.tgz", "导入文件", curTime)
	filePath := fmt.Sprintf("%s/%s", DownloadBasePath, fileName)
	importFilePath := fmt.Sprintf("%s/%s", DownloadBasePath, importFileName)
	_ = WriteDataInCommonFile(filePath, sqlRule1)
	_ = WriteDataInCommonFile(filePath, sqlRule2)
	_ = WriteDataInCommonFile(filePath, sqlRule3)
	_ = WriteDataInCommonFile(filePath, "")

	_ = WriteDataInCommonFile(filePath, "# 选择数据库, 默认data4test")
	_ = WriteDataInCommonFile(filePath, "use data4test;")
	_ = WriteDataInCommonFile(filePath, "set names utf8;")
	_ = WriteDataInCommonFile(filePath, "")

	dataMapFromTask, playbookMap, productMap, err := GetTaskSQL(userName, id, filePath)
	if err != nil {
		return
	}

	productName, appMapFromProduct, err := GetProductSQL(filePath, productMap)
	if err != nil {
		return
	}

	dataMap, err := GetPlaybookSQL(userName, productName, filePath, playbookMap)
	if err != nil {
		return
	}
	for k, _ := range dataMapFromTask {
		if _, ok := dataMap[k]; !ok {
			dataMap[k] = true
		}
	}

	assertMap, sysParameterMap, appMap, fileNameImportMap, err := GetDataSQL(userName, filePath, dataMap)
	if err != nil {
		return
	}

	for k, _ := range appMapFromProduct {
		if _, ok := appMap[k]; !ok {
			appMap[k] = true
		}
	}

	err = GetAppSQL(productName, filePath, appMap)
	if err != nil {
		return
	}

	err = GetAssertTemplateSQL(filePath, assertMap)
	if err != nil {
		return
	}

	err = GetSysParameterSQL(filePath, sysParameterMap)
	if err != nil {
		return
	}

	isExist, err := GetImportFilePackage(importFileName, importFilePath, fileNameImportMap)
	if err != nil {
		return
	}

	if isExist {
		fileName = fmt.Sprintf("%s, %s", fileName, importFileName)
	}

	return
}

func GetAppSQL(productName, filePath string, appMap map[string]bool) (err error) {
	appSQLDesc := "# 应用配置"
	var appNames []string
	for k, _ := range appMap {
		if len(k) == 0 {
			continue
		}
		appNames = append(appNames, k)
	}

	var appList []EnvConfig
	models.Orm.Table("env_config").Where("app in (?)", appNames).Find(&appList)
	if len(appList) == 0 {
		err = fmt.Errorf("未找到[%v]应用配置信息", appList)
		Logger.Error("%s", err)
		return
	}

	var appValueStr, appNameDesc string
	for index, item := range appList {
		if strings.Contains(item.Remark, "'") {
			item.Remark = strings.Replace(item.Remark, "'", "\\'", -1)
		}
		if index == 0 {
			appValueStr = fmt.Sprintf("('%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')", productName, item.App, item.Ip, item.Protocol, item.Prepath, item.Threading, item.Auth, item.Testmode, item.SwaggerPath, item.Remark)
			appNameDesc = fmt.Sprintf("# 导出的应用为: %s", item.App)
		} else {
			appValueStr = fmt.Sprintf("%s, ('%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')", appValueStr, productName, item.App, item.Ip, item.Protocol, item.Prepath, item.Threading, item.Auth, item.Testmode, item.SwaggerPath, item.Remark)
			appNameDesc = fmt.Sprintf("%s / %s", appNameDesc, item.App)
		}
	}

	appSQL := fmt.Sprintf("INSERT IGNORE INTO `env_config`(product, app, ip, protocol, prepath, threading, auth, testmode, swagger_path, remark) VALUES %s;", appValueStr)
	_ = WriteDataInCommonFile(filePath, appSQLDesc)
	_ = WriteDataInCommonFile(filePath, appNameDesc)
	_ = WriteDataInCommonFile(filePath, appSQL)
	_ = WriteDataInCommonFile(filePath, "")

	return
}

func GetProductSQL(filePath string, productMap map[string]bool) (productName string, appMap map[string]bool, err error) {
	productSQLDesc := "# 产品配置"
	var productNameList []string
	for k, _ := range productMap {
		if len(k) == 0 {
			continue
		}
		productNameList = append(productNameList, k)
	}

	var productList []DbProduct
	models.Orm.Table("product").Where("product in (?)", productNameList).Find(&productList)
	if len(productList) == 0 {
		err = fmt.Errorf("未找到: %v 环境信息，请核对", productNameList)
		Logger.Error("%s", err)
		return
	}

	var productValueStr, productNameDesc string
	var appNameList []string
	for index, item := range productList {
		if strings.Contains(item.Remark, "'") {
			item.Remark = strings.Replace(item.Remark, "'", "\\'", -1)
		}
		if index == 0 {
			productName = item.Product.Name
			productValueStr = fmt.Sprintf("('%s','%s','%s','%s',%d,'%s','%s','%s',%d,'%s','%s')", item.Product.Name, item.Ip, item.Protocol, item.Threading, item.ThreadNumber, item.Auth, item.Testmode, item.Apps, item.EnvType, item.PrivateParameter, item.Remark)
			appNameList = strings.Split(item.Apps, ",")
			productNameDesc = fmt.Sprintf("# 导出的产品为: %s", item.Name)
		} else {
			productValueStr = fmt.Sprintf("%s, ('%s','%s','%s','%s',%d,'%s','%s','%s',%d,'%s','%s')", productValueStr, item.Product.Name, item.Ip, item.Protocol, item.Threading, item.ThreadNumber, item.Auth, item.Testmode, item.Apps, item.EnvType, item.PrivateParameter, item.Remark)
			appTmp := strings.Split(item.Apps, ",")
			appNameList = append(appNameList, appTmp...)
			productNameDesc = fmt.Sprintf("%s / %s", productNameDesc, item.Name)
		}
	}

	appMap = make(map[string]bool)
	for _, item := range appNameList {
		if _, ok := appMap[item]; !ok {
			appMap[item] = true
		}
	}

	productSQL := fmt.Sprintf("INSERT IGNORE INTO `product`(product, ip, protocol, threading, thread_number, auth,testmode, apps, env_type, private_parameter, remark) VALUES %s;", productValueStr)

	_ = WriteDataInCommonFile(filePath, productSQLDesc)
	_ = WriteDataInCommonFile(filePath, productNameDesc)
	_ = WriteDataInCommonFile(filePath, productSQL)
	_ = WriteDataInCommonFile(filePath, "")

	return
}

func GetTaskSQL(userName, taskId, filePath string) (dataMap, playbookMap, productMap map[string]bool, err error) {
	taskSQLDesc := "# 任务信息"
	var taskList []DbSchedule
	ids := strings.Split(taskId, ",")
	models.Orm.Table("schedule").Where("id in (?)", ids).Find(&taskList)
	if len(taskList) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", taskId)
		Logger.Error("%s", err)
		return
	}

	var productNames []string
	var taskValueStr, taskNameDesc string
	for index, item := range taskList {
		if strings.Contains(item.Remark, "'") {
			item.Remark = strings.Replace(item.Remark, "'", "\\'", -1)
		}
		if index == 0 {
			taskValueStr = fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", item.TaskName, item.TaskMode, item.Threading, item.TaskType, item.DataNumber, item.DataList, item.SceneNumber, item.SceneList, item.ProductList, item.Time4week, item.Time4day, item.Week, item.Remark, userName)
			productNames = strings.Split(item.ProductList, ",")
			taskNameDesc = fmt.Sprintf("# 导出的任务为: %s", item.TaskName)
		} else {
			taskValueStr = fmt.Sprintf("%s, ('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", taskValueStr, item.TaskName, item.TaskMode, item.Threading, item.TaskType, item.DataNumber, item.DataList, item.SceneNumber, item.SceneList, item.ProductList, item.Time4week, item.Time4day, item.Week, item.Remark, userName)
			productTmp := strings.Split(item.ProductList, ",")
			productNames = append(productNames, productTmp...)
			taskNameDesc = fmt.Sprintf("%s / %s", taskNameDesc, item.TaskName)
		}
	}

	productMap = make(map[string]bool)
	for _, item := range productNames {
		if _, ok := productMap[item]; !ok {
			productMap[item] = true
		}
	}

	playbookMap = make(map[string]bool)
	dataMap = make(map[string]bool)
	for _, item := range taskList {
		if item.TaskType == "data" {
			dataTmp := strings.Split(item.DataList, ",")
			for _, item := range dataTmp {
				if _, ok := dataMap[item]; !ok {
					dataMap[item] = true
				}
			}
		} else if item.TaskType == "scene" {
			playbookTmp := strings.Split(item.SceneList, ",")
			for _, item := range playbookTmp {
				if _, ok := playbookMap[item]; !ok {
					playbookMap[item] = true
				}
			}
		}
	}

	taskSQL := fmt.Sprintf("REPLACE INTO `schedule`(task_name, task_mode, threading, task_type, data_number, data_list, scene_number, scene_list, product_list, time4week, time4day, week, remark, user_name) VALUES %s;", taskValueStr)
	_ = WriteDataInCommonFile(filePath, taskSQLDesc)
	_ = WriteDataInCommonFile(filePath, taskNameDesc)
	_ = WriteDataInCommonFile(filePath, taskSQL)
	_ = WriteDataInCommonFile(filePath, "")

	return
}

func GetPlaybookSQL(userName, productName, filePath string, playbookMap map[string]bool) (dataMap map[string]bool, err error) {
	playbookSQLDesc := "# 场景信息"
	var playbookNameList []string
	var playbookValueStr string
	for k, _ := range playbookMap {
		playbookNameList = append(playbookNameList, k)
	}

	var playbookList []DbScene
	models.Orm.Table("playbook").Where("name in (?)", playbookNameList).Find(&playbookList)
	if len(playbookList) == 0 {
		err = fmt.Errorf("未找到对应场景，请核对: %s", playbookNameList)
		Logger.Error("%s", err)
		return
	}

	dataMap = make(map[string]bool)
	for index, item := range playbookList {
		if strings.Contains(item.Remark, "'") {
			item.Remark = strings.Replace(item.Remark, "'", "\\'", -1)
		}
		if index == 0 {
			playbookValueStr = fmt.Sprintf("('%s','%s','%s',%d,%d,'%s','%s','%s')", item.Name, item.DataNumber, item.ApiList, item.SceneType, item.RunTime, item.Remark, userName, productName)
		} else {
			playbookValueStr = fmt.Sprintf("%s, ('%s','%s','%s',%d,%d,'%s','%s','%s')", playbookValueStr, item.Name, item.DataNumber, item.ApiList, item.SceneType, item.RunTime, item.Remark, userName, productName)
		}
		dataTmp := GetListFromHtml(item.ApiList)
		for _, dataItem := range dataTmp {
			if _, ok := dataMap[dataItem]; !ok {
				dataMap[dataItem] = true
			}
		}
	}

	playbookNoDesc := fmt.Sprintf("# 导出的场景数量为: %d", len(playbookList))

	playbookSQL := fmt.Sprintf("REPLACE INTO `playbook`(name, data_number, api_list, scene_type, run_time, remark, user_name, product) VALUES %s;", playbookValueStr)
	_ = WriteDataInCommonFile(filePath, playbookSQLDesc)
	_ = WriteDataInCommonFile(filePath, playbookNoDesc)
	_ = WriteDataInCommonFile(filePath, playbookSQL)
	_ = WriteDataInCommonFile(filePath, "")

	return
}

func GetDataSQL(userName, filePath string, dataMap map[string]bool) (assertMap, sysParameterMap, appMap, fileNameImportMap map[string]bool, err error) {
	dataSQLDesc := "# 数据信息"
	var dataNameList []string
	var dataValueStr string
	for k, _ := range dataMap {
		dataNameList = append(dataNameList, k)
	}

	var dataList []DbSceneData
	models.Orm.Table("scene_data").Where("file_name in (?)", dataNameList).Find(&dataList)
	if len(dataList) == 0 {
		err = fmt.Errorf("未找到[%v]数据，请核对", dataNameList)
		Logger.Error("%s", err)
		return
	}

	var appNames []string
	var apiNum int
	models.Orm.Table("scene_data").Where("file_name in (?)", dataNameList).Group("api_id").Count(&apiNum)

	models.Orm.Table("scene_data").Where("file_name in (?)", dataNameList).Group("app").Select("app").Pluck("app", &appNames)

	if len(appNames) == 0 {
		err = fmt.Errorf("未找到[%v]应用信息，请核对", dataNameList)
		Logger.Error("%s", err)
		return
	}

	appMap = make(map[string]bool)
	for _, value := range appNames {
		appMap[value] = true
	}

	fileNameDef := GetValuesFromSysParameter("default", "fileName")

	assertMap = make(map[string]bool)
	sysParameterMap = make(map[string]bool)
	fileNameImportMap = make(map[string]bool)

	var isGetMockFile bool
	for index, item := range dataList {
		// 历史数据修正
		if strings.Contains(item.Content, "<pre><code>") {
			item.Content = strings.Replace(item.Content, "<pre><code>", "", -1)
			item.Content = strings.Replace(item.Content, "</code></pre>", "", -1)
		} else {
			item.Content = item.Content
		}

		if strings.Contains(item.Content, "multipart/form-data") {
			var df DataFile
			if strings.HasSuffix(item.FileName, ".json") {
				err = json.Unmarshal([]byte(item.Content), &df)
			} else {
				err = yaml.Unmarshal([]byte(item.Content), &df)
			}

			if err != nil {
				Logger.Debug("fileName: %s", item.FileName)
				Logger.Error("%s", err)
				return
			}

			for _, value := range fileNameDef {
				if _, ok := df.Single.Body[value]; ok {
					fileNameImport := Interface2Str(df.Single.Body[value])
					fileNameImportMap[fileNameImport] = true
				}
				if _, ok := df.Multi.Body[value]; ok {
					for _, subV := range df.Multi.Body[value] {
						fileNameImport := Interface2Str(subV)
						fileNameImportMap[fileNameImport] = true
					}
				}
			}

		}

		if strings.Contains(item.Content, "/mock/file/") && !isGetMockFile {
			isGetMockFile = true
			sysParameterMap["MockTemplateFile"] = true
			values := GetValuesFromSysParameter("", "MockTemplateFile")
			for _, v := range values {
				fileNameImportMap[v] = true
			}
		}

		if strings.Contains(item.Content, "'") {
			item.Content = strings.Replace(item.Content, "'", "\\'", -1)
		}
		if strings.Contains(item.Remark, "'") {
			item.Remark = strings.Replace(item.Remark, "'", "\\'", -1)
		}
		if index == 0 {
			dataValueStr = fmt.Sprintf("('%s','%s','%s','%s',%d,'%s',%d,'%s','%s')", item.Name, item.ApiId, item.App, item.FileName, item.FileType, item.Content, item.RunTime, item.Remark, userName)
		} else {
			dataValueStr = fmt.Sprintf("%s, ('%s','%s','%s','%s',%d,'%s',%d,'%s','%s')", dataValueStr, item.Name, item.ApiId, item.App, item.FileName, item.FileType, item.Content, item.RunTime, item.Remark, userName)
		}

		indexStr, _, _, _ := GetStrByIndex(item.Content, "assert:\n", "output: {}")
		if len(indexStr) > 0 {
			assertVars, _ := GetInStrDef(indexStr)
			for item, _ := range assertVars {
				if _, ok := assertMap[item]; !ok {
					assertMap[item] = true
				}
			}
		}

		indexStr, _, _, _ = GetStrByIndex(item.Content, "env:\n", "action:")
		if len(indexStr) > 0 {
			systemParameterVars, _ := GetInStrDef(indexStr)
			for item, _ := range systemParameterVars {
				if _, ok := sysParameterMap[item]; !ok {
					sysParameterMap[item] = true
				}
			}
		}
	}

	dataNoDesc := fmt.Sprintf("# 导出的数据数为: %d, 覆盖的接口数为: %d", len(dataList), apiNum)

	dataSQL := fmt.Sprintf("REPLACE INTO `scene_data`(name, api_id, app, file_name, file_type, content, run_time, remark, user_name) VALUES %s;", dataValueStr)
	_ = WriteDataInCommonFile(filePath, dataSQLDesc)
	_ = WriteDataInCommonFile(filePath, dataNoDesc)
	_ = WriteDataInCommonFile(filePath, dataSQL)
	_ = WriteDataInCommonFile(filePath, "")

	return
}

func GetAssertTemplateSQL(filePath string, assertMap map[string]bool) (err error) {
	assertTemplateSQLDesc := "# 断言模板信息"
	var assertNameList []string
	var assertValueStr string

	for k, _ := range assertMap {
		assertNameList = append(assertNameList, k)
	}
	//无断言模板信息，直接跳过
	if len(assertNameList) == 0 {
		return
	}
	var assertList []AssertValueDefine
	models.Orm.Table("assert_template").Where("name in (?)", assertNameList).Find(&assertList)
	if len(assertList) == 0 {
		err = fmt.Errorf("未关联到断言值定义:%v", assertNameList)
		Logger.Error("%s", err)
		return
	}

	for index, item := range assertList {
		if strings.Contains(item.Remark, "'") {
			item.Remark = strings.Replace(item.Remark, "'", "\\'", -1)
		}
		if index == 0 {
			assertValueStr = fmt.Sprintf("('%s','%s','%s','%s')", item.Name, item.Value, item.Remark, item.UerName)
		} else {
			assertValueStr = fmt.Sprintf("%s, ('%s','%s','%s','%s')", assertValueStr, item.Name, item.Value, item.Remark, item.UerName)
		}
	}

	assertNoDesc := fmt.Sprintf("# 导出的断言模板数量为: %d", len(assertList))

	assertTemplateSQL := fmt.Sprintf("REPLACE INTO `assert_template`(name, value, remark, user_name) VALUES %s;", assertValueStr)
	_ = WriteDataInCommonFile(filePath, assertTemplateSQLDesc)
	_ = WriteDataInCommonFile(filePath, assertNoDesc)
	_ = WriteDataInCommonFile(filePath, assertTemplateSQL)
	_ = WriteDataInCommonFile(filePath, "")

	return
}

func GetSysParameterSQL(filePath string, systemParameterMap map[string]bool) (err error) {
	sysParameterSQLDesc := "# 系统参数信息"
	var sysParameterNameList []string
	var sysParameterValueStr string

	for k, _ := range systemParameterMap {
		sysParameterNameList = append(sysParameterNameList, k)
	}

	//无系统定义，无需拦截
	if len(sysParameterNameList) == 0 {
		return
	}

	var sysParameterList []SysParameter
	models.Orm.Table("sys_parameter").Where("name in (?)", sysParameterNameList).Find(&sysParameterList)

	//无系统参数，无需拦截
	if len(sysParameterList) == 0 {
		return
	}

	for index, item := range sysParameterList {
		if strings.Contains(item.Remark, "'") {
			item.Remark = strings.Replace(item.Remark, "'", "\\'", -1)
		}
		if index == 0 {
			sysParameterValueStr = fmt.Sprintf("('%s','%s','%s')", item.Name, item.ValueList, item.Remark)
		} else {
			sysParameterValueStr = fmt.Sprintf("%s, ('%s','%s','%s')", sysParameterValueStr, item.Name, item.ValueList, item.Remark)
		}
	}

	sysParamterNoDesc := fmt.Sprintf("# 导出的系统参数数量为: %d", len(sysParameterList))

	sysParameterSQL := fmt.Sprintf("REPLACE INTO `sys_parameter`(name, value_list, remark) VALUES %s;", sysParameterValueStr)
	_ = WriteDataInCommonFile(filePath, sysParameterSQLDesc)
	_ = WriteDataInCommonFile(filePath, sysParamterNoDesc)
	_ = WriteDataInCommonFile(filePath, sysParameterSQL)
	_ = WriteDataInCommonFile(filePath, "")

	return
}

func GetImportFilePackage(fileName, filePath string, fileNameImportMap map[string]bool) (isExist bool, err error) {
	count := 0

	fw, err := os.Create(filePath)
	if err != nil {
		Logger.Error("%s", err)
	}
	defer fw.Close()
	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()
	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for item, _ := range fileNameImportMap {
		count++
		srcFilePath := fmt.Sprintf("%s/%s", UploadBasePath, item)
		fi, errTmp := os.Stat(srcFilePath)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			err = errTmp
			return
		}
		fr, errTmp := os.Open(srcFilePath)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
			err = errTmp
			return
		}

		h := new(tar.Header)
		h.Name = fi.Name()
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()
		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			panic(err)
		}
		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			panic(err)
		}
	}

	if count > 0 {
		isExist = true
	}

	return
}

func UpdateTaskInfoList(id, taskType string, dataList, dataNumList []string) (err error) {
	var dbSchedule DbSchedule
	models.Orm.Table("schedule").Where("id = ?", id).Find(&dbSchedule)
	var dataStr, dataNumStr string
	if len(dbSchedule.TaskName) == 0 {
		return
	}

	var tDataList []string

	for _, v := range dataList {
		if len(v) == 0 {
			continue
		}
		tDataList = append(tDataList, v)
	}

	for index, value := range tDataList {
		if index == 0 {
			dataStr = value
			if len(dataNumList) == 0 {
				dataNumList = append(dataNumList, fmt.Sprintf("%d", index+1))
			}
			dataNumStr = fmt.Sprintf("%v", dataNumList[index])
		} else {
			dataStr = fmt.Sprintf("%s,%s", dataStr, value)
			if len(dataNumList)-1 < index {
				dataNumList = append(dataNumList, fmt.Sprintf("%d", index+1))
			} else if len(dataNumList[index]) == 0 {
				dataNumList[index] = fmt.Sprintf("%d", index+1)
			}
			dataNumStr = fmt.Sprintf("%s,%v", dataNumStr, dataNumList[index])

		}
	}

	dbSchedule.TaskType = taskType

	if dbSchedule.TaskType == "data" {
		dbSchedule.DataList = dataStr
		dbSchedule.DataNumber = dataNumStr
	} else {
		dbSchedule.SceneList = dataStr
		dbSchedule.SceneNumber = dataNumStr
	}

	err = models.Orm.Table("schedule").Where("id = ?", dbSchedule.Id).Update(dbSchedule).Error
	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func AutoCreateSchedule(id string, taskType string) (err error) {
	var schedule Schedule4Copy
	curTimeStr := fmt.Sprintf(time.Unix(time.Now().Unix(), 0).Format("20060102"))
	idList := strings.Split(id, ",")
	var numStr string
	for index, item := range idList {
		if len(item) == 0 {
			continue
		}
		if index == 0 {
			numStr = fmt.Sprintf("%v", index+1)
		} else {
			numStr = fmt.Sprintf("%s,%v", numStr, index+1)
		}
	}
	var dbSceneDatas []DbSceneData
	var playbooks []Playbook
	if taskType == "data" {
		models.Orm.Table("scene_data").Where("id in (?)", idList).Find(&dbSceneDatas)
		if len(dbSceneDatas) == 0 {
			err = fmt.Errorf("未找到对应[%v]的场景数据，请核对", id)
			Logger.Error("%s", err)
			return
		}
		for _, item := range dbSceneDatas {
			if len(schedule.DataList) == 0 {
				schedule.DataList = item.Name
			} else {
				schedule.DataList = fmt.Sprintf("%s,%s", schedule.DataList, item.Name)
			}
		}
		schedule.DataNumber = numStr
		schedule.TaskName = fmt.Sprintf("数据任务_%s_%s", curTimeStr, GetRandomStr(4, ""))
	} else if taskType == "scene" {
		models.Orm.Table("playbook").Where("id in (?)", idList).Find(&playbooks)
		if len(playbooks) == 0 {
			err = fmt.Errorf("未找到对应[%v]的场景数据，请核对", id)
			Logger.Error("%s", err)
			return
		}
		for _, item := range playbooks {
			if len(schedule.SceneList) == 0 {
				schedule.SceneList = item.Name
			} else {
				schedule.SceneList = fmt.Sprintf("%s,%s", schedule.SceneList, item.Name)
			}
		}
		schedule.SceneNumber = numStr
		schedule.TaskName = fmt.Sprintf("场景任务_%s_%s", curTimeStr, GetRandomStr(4, ""))

	}

	schedule.TaskType = taskType
	schedule.Threading = "no"
	schedule.TaskMode = "once"
	schedule.TaskStatus = "not_started"
	schedule.UserName = "系统"

	err = models.Orm.Table("schedule").Create(schedule).Error
	if err != nil {
		Logger.Error("%s", err)
	}

	// 创建任务直接调任务先屏蔽，正常功能可能需要先创建任务，再单独发起执行
	//var dbSchedule DbSchedule
	//
	//models.Orm.Table("schedule").Where("task_name =?", schedule.TaskName).Find(&dbSchedule)
	//if len(dbSchedule.TaskName) == 0 {
	//	Logger.Error("未找到名称[%s]的任务，请核对", schedule.TaskName)
	//	return
	//}
	//
	//if taskType == "scene" {
	//	err = RunTask(dbSchedule.Id)
	//}

	return
}

func GetTaskDataStr(id string) (dataStr string) {
	var task DbSchedule
	models.Orm.Table("schedule").Where("id = ?", id).Find(&task)
	if len(task.TaskName) == 0 {
		return
	}

	if task.TaskType == "data" {
		dataStr = task.DataList
	} else {
		dataStr = task.SceneList
	}

	dataStr = strings.Replace(dataStr, ",", "\r\n", -1)

	return
}

func GetTaskEditTypeById(id string) (editType string) {
	var task DbSchedule
	models.Orm.Table("schedule").Where("id = ?", id).Find(&task)
	if len(task.TaskName) == 0 {
		return "select_scene"
	}

	if task.TaskType == "data" {
		if len(task.DataList) > 0 {
			editType = "input_data"
		} else {
			editType = "select_data"
		}
	} else if task.TaskType == "scene" {
		if len(task.SceneList) > 0 {
			editType = "input_scene"
		} else {
			editType = "select_scene"
		}
	} else {
		editType = "input_scene"
	}

	return
}

func GetPlaybookLinkByPlaybookStr(pStr string) (linkStr string) {
	pList := strings.Split(pStr, ",")
	for _, item := range pList {
		if len(item) == 0 {
			continue
		}
		var ids []int
		models.Orm.Table("playbook").Where("name = ?", item).Pluck("id", &ids)
		if len(ids) == 0 {
			Logger.Warning("未找到场景[%s], 请核对", item)
			if len(linkStr) == 0 {
				linkStr = item //跳详情，可点击编辑进行改写
			} else {
				linkStr = fmt.Sprintf("%s<br>%s", linkStr, item)
			}
		} else {
			if len(linkStr) == 0 {
				linkStr = fmt.Sprintf("<a href=\"/admin/info/playbook/detail?__goadmin_detail_pk=%d\">%s</a>", ids[0], item) //跳编辑区可直接改写
			} else {
				linkStr = fmt.Sprintf("%s<br><a href=\"/admin/info/playbook/detail?__goadmin_detail_pk=%d\">%s</a>", linkStr, ids[0], item)
			}
		}
	}
	return
}
