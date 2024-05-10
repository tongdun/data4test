package biz

type Playbook struct {
	Name        string
	Apis        []string
	LastFile    string
	Tag         int
	Product     string
	IsThread    string
	SceneType   int
	HistoryApis []string
}

type Scene struct {
	Name       string `gorm:"column:name" json:"name" yaml:"name"`
	DataNumber string `gorm:"column:data_number" json:"data_number" yaml:"data_number"`
	ApiList    string `gorm:"column:api_list" json:"api_list" yaml:"api_list"`
	LastFile   string `gorm:"column:last_file" json:"last_file" yaml:"last_file"`
	Result     string `gorm:"column:result" json:"result" yaml:"result"`
	RunTime    int    `gorm:"column:run_time" json:"run_time" yaml:"run_time"`
	SceneType  int    `gorm:"column:scene_type" json:"scene_type" yaml:"scene_type"`
	Product    string `gorm:"column:product" json:"product" yaml:"product"`
	Remark     string `gorm:"column:remark" json:"remark" yaml:"remark"`
	UserName   string `gorm:"column:user_name" json:"remark" yaml:"user_name"`
	Priority   int    `gorm:"column:priority" json:"priority" yaml:"priority"`
	FailReason string `gorm:"column:fail_reason" json:"fail_reason" yaml:"fail_reason"`
	UpdatedAt  string `gorm:"column:updated_at" json:"updated_at"`
}

type DbScene struct {
	Id string `gorm:"column:id" json:"id"`
	Scene
}

type DbSceneWithNoUpdateTime struct {
	Id string `gorm:"column:id" json:"id"`
	SceneWithNoUpdateTime
}

type SceneWithNoUpdateTime struct {
	Name       string `gorm:"column:name" json:"name" yaml:"name"`
	DataNumber string `gorm:"column:data_number" json:"data_number" yaml:"data_number"`
	ApiList    string `gorm:"column:api_list" json:"api_list" yaml:"api_list"`
	LastFile   string `gorm:"column:last_file" json:"last_file" yaml:"last_file"`
	Result     string `gorm:"column:result" json:"result" yaml:"result"`
	RunTime    int    `gorm:"column:run_time" json:"run_time" yaml:"run_time"`
	SceneType  int    `gorm:"column:scene_type" json:"scene_type" yaml:"scene_type"`
	Product    string `gorm:"column:product" json:"product" yaml:"product"`
	Remark     string `gorm:"column:remark" json:"remark" yaml:"remark"`
	UserName   string `gorm:"column:user_name" json:"remark" yaml:"user_name"`
	Priority   int    `gorm:"column:priority" json:"priority" yaml:"priority"`
	FailReason string `gorm:"column:fail_reason" json:"fail_reason" yaml:"fail_reason"`
}

type SceneRecord struct {
	Name       string `gorm:"column:name" json:"name" yaml:"name"`
	ApiList    string `gorm:"column:api_list" json:"api_list" yaml:"api_list"`
	LastFile   string `gorm:"column:last_file" json:"last_file" yaml:"last_file"`
	Result     string `gorm:"column:result" json:"result" yaml:"result"`
	SceneType  int    `gorm:"column:scene_type" json:"scene_type" yaml:"scene_type"`
	Product    string `gorm:"column:product" json:"product" yaml:"product"`
	FailReason string `gorm:"column:fail_reason" json:"fail_reason" yaml:"fail_reason"`
	EnvType    int    `gorm:"column:env_type" json:"env_type" yaml:"env_type"`
}

type DbSceneRecord struct {
	Id        string `gorm:"column:id" json:"id"`
	UpdatedAt string `gorm:"column:updated_at" json:"updated_at"`
	Product   string `gorm:"column:product" json:"product" yaml:"product"`
	SceneRecord
}

type SceneSaveModel struct {
	Product   string         `json:"product"`
	Name      string         `json:"name"`
	DataList  []DepDataModel `json:"dataList"`
	SceneType int            `json:"type"`
	RunNum    int            `json:"runNum"`
}
