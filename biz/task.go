package biz

import (
	"data4perf/models"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

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

func UpdateTaskInfoList(id string, dataList, dataNumList, sceneList, sceneNumList []string) (err error) {
	var dbSchedule DbSchedule
	models.Orm.Table("schedule").Where("id = ?", id).Find(&dbSchedule)
	var dataStr, dataNumStr, sceneStr, sceneNumStr string
	if len(dbSchedule.TaskName) == 0 {
		return
	} else {
		for index, value := range dataList {
			if len(dataList) > index {
				numValue := dataList[index]
				if len(numValue) == 0 {
					dataNumList[index] = fmt.Sprintf("%d", index+1)
				}
			} else {
				dataNumList = append(dataNumList, fmt.Sprintf("%d", index+1))
			}
			if index == 0 {
				dataStr = value
				dataNumStr = fmt.Sprintf("%v", dataNumList[index])
			} else {
				dataStr = fmt.Sprintf("%s,%v", dataStr, value)
				dataNumStr = fmt.Sprintf("%s,%v", dataNumStr, dataNumList[index])
			}
		}

		for index, value := range sceneList {
			if len(sceneList) > index {
				numValue := sceneNumList[index]
				if len(numValue) == 0 {
					sceneNumList[index] = fmt.Sprintf("%d", index+1)
				}
			} else {
				sceneNumList = append(sceneNumList, fmt.Sprintf("%d", index+1))
			}
			if index == 0 {
				sceneStr = value
				sceneNumStr = fmt.Sprintf("%v", sceneNumList[index])
			} else {
				sceneStr = fmt.Sprintf("%s,%v", sceneStr, value)
				sceneNumStr = fmt.Sprintf("%s,%v", sceneNumStr, sceneNumList[index])
			}
		}

		dbSchedule.DataList = dataStr
		dbSchedule.DataNumber = dataNumStr
		dbSchedule.SceneList = sceneStr
		dbSchedule.SceneNumber = sceneNumStr
		err = models.Orm.Table("schedule").Where("id = ?", dbSchedule.Id).Update(dbSchedule).Error
		if err != nil {
			Logger.Error("%s", err)
		}
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
