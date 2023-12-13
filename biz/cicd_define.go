package biz

type AppApiChange struct {
	App     string `gorm:"column:app" json:"app" yaml:"app"`
	CurApiSum int `gorm:"column:curApiSum" json:"curApiSum" yaml:"curApiSum"`
	ExistApiSum int `gorm:"column:existApiSum" json:"exsitApiSum" yaml:"existApiSum"`
	NewApiSum int `gorm:"column:newApiSum" json:"newApiSum" yaml:"newApiSum"`
	DeletedApiSum int `gorm:"column:deletedApiSum" json:"deletedApiSum" yaml:"deletedApiSum"`
	ChangedApiSum    int `gorm:"column:changedApiSum" json:"changedApiSum" yaml:"changedApiSum"`
	CheckFailApiSum    int `gorm:"column:checkFailApiSum" json:"checkFailApiSum" yaml:"checkFailApiSum"`
	NewApiContent    string `gorm:"column:newApiContent" json:"fail_reason" yaml:"newApiContent"`
	DeletedApiContent    string `gorm:"column:deletedApiContent" json:"deletedApiContent" yaml:"deletedApiContent"`
	ChangedApiContent    string `gorm:"column:changedApiContent" json:"changedApiContent" yaml:"changedApiContent"`
	Branch    string `gorm:"column:branch" json:"branch" yaml:"branch"`
	ApiCheckResult    string `gorm:"column:apiCheckResult" json:"apiCheckResult" yaml:"apiCheckResult"`
	ApiCheckFailContent    string `gorm:"column:apiCheckFailContent" json:"apiCheckFailContent" yaml:"apiCheckFailContent"`
	Remark    string `gorm:"column:remark" json:"remark" yaml:"remark"`
}

