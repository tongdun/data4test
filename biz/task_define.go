package biz

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type Schedule struct {
	TaskName    string `gorm:"column:task_name" json:"task_name" yaml:"task_name"`
	TaskMode    string `gorm:"column:task_mode" json:"task_mode" yaml:"task_mode"`
	Crontab     string `gorm:"column:crontab" json:"crontab" yaml:"crontab"`
	Threading   string `gorm:"column:threading" json:"threading" yaml:"threading"`
	TaskType    string `gorm:"column:task_type" json:"task_type" yaml:"task_type"`
	DataNumber  string `gorm:"column:data_number" json:"data_number" yaml:"data_number"`
	DataList    string `gorm:"column:data_list" json:"data_list" yaml:"data_list"`
	SceneNumber string `gorm:"column:scene_number" json:"scene_number" yaml:"scene_number"`
	SceneList   string `gorm:"column:scene_list" json:"scene_list" yaml:"scene_list"`
	ProductList string `gorm:"column:product_list" json:"product_list" yaml:"product_list"`
	TaskStatus  string `gorm:"column:task_status" json:"task_status" yaml:"task_status"`
	Time4week   string `gorm:"column:time4week" json:"time4week" yaml:"time4week"`
	Time4day    string `gorm:"column:time4day" json:"time4day" yaml:"time4day"`
	Week        string `gorm:"column:week" json:"week" yaml:"week"`
	LastAt      string `gorm:"column:last_at" json:"last_at" yaml:"last_at"`
	NextAt      string `gorm:"column:next_at" json:"next_at" yaml:"next_at"`
	Remark      string `gorm:"column:remark" json:"remark" yaml:"remark"`
	UserName    string `gorm:"column:user_name" json:"remark" yaml:"user_name"`
}

type Schedule4Copy struct {
	TaskName    string `gorm:"column:task_name" json:"task_name" yaml:"task_name"`
	TaskMode    string `gorm:"column:task_mode" json:"task_mode" yaml:"task_mode"`
	Crontab     string `gorm:"column:crontab" json:"crontab" yaml:"crontab"`
	Threading   string `gorm:"column:threading" json:"threading" yaml:"threading"`
	TaskType    string `gorm:"column:task_type" json:"task_type" yaml:"task_type"`
	DataNumber  string `gorm:"column:data_number" json:"data_number" yaml:"data_number"`
	DataList    string `gorm:"column:data_list" json:"data_list" yaml:"data_list"`
	SceneNumber string `gorm:"column:scene_number" json:"scene_number" yaml:"scene_number"`
	SceneList   string `gorm:"column:scene_list" json:"scene_list" yaml:"scene_list"`
	ProductList string `gorm:"column:product_list" json:"product_list" yaml:"product_list"`
	TaskStatus  string `gorm:"column:task_status" json:"task_status" yaml:"task_status"`
	Time4week   string `gorm:"column:time4week" json:"time4week" yaml:"time4week"`
	Time4day    string `gorm:"column:time4day" json:"time4day" yaml:"time4day"`
	Week        string `gorm:"column:week" json:"week" yaml:"week"`
	Remark      string `gorm:"column:remark" json:"remark" yaml:"remark"`
	UserName    string `gorm:"column:user_name" json:"remark" yaml:"user_name"`
}

type DbSchedule struct {
	Id      string `gorm:"column:id" json:"id"`
	StopTag chan struct{}
	Schedule
}

type Crontab struct {
	inner *cron.Cron
	ids   map[string]cron.EntryID
	mutex sync.Mutex
}

type RunSchecdule struct {
	Id string
}
