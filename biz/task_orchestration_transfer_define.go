package biz

// KOrchestration 编排导出清单条目
type KOrchestration struct {
	Name     string   `yaml:"name" json:"name"`
	Tasks    []string `yaml:"tasks" json:"tasks"`
	Products string   `yaml:"products" json:"products"`
	Remark   string   `yaml:"remark" json:"remark"`
}

// OrchestrationTaskMapping 检查阶段的任务映射结果
type OrchestrationTaskMapping struct {
	TaskName string `json:"task_name"`
	Source   string `json:"source"` // existing / import / missing
	TaskId   string `json:"task_id"`
}

// OrchestrationCheckDetail 单个编排的检查明细
type OrchestrationCheckDetail struct {
	Name            string                     `json:"name"`
	Exists          bool                       `json:"exists"`
	Tasks           []OrchestrationTaskMapping `json:"tasks"`
	Missing         []string                   `json:"missing"`          // 缺失任务（阻断）
	MissingScenes   []string                   `json:"missing_scenes"`   // 缺失场景（阻断）
	MissingData     []string                   `json:"missing_data"`     // 缺失数据文件（阻断）
	MissingProducts []string                   `json:"missing_products"` // 缺失产品（仅提示）
}

// OrchestrationImportState 编排导入检查结果（confirm 阶段使用）
type OrchestrationImportState struct {
	Blocked          bool             `json:"blocked"`
	DependencyError  string           `json:"dependency_error,omitempty"`
	ProductConfigs   []CliConfig      `json:"product_configs,omitempty"`
	ImportId         string           `json:"import_id"`
	UserName         string           `json:"user_name"`
	TempDir          string           `json:"temp_dir"`
	ScheduleImportId string           `json:"schedule_import_id"`
	Orchestrations   []KOrchestration `json:"orchestrations"`
	Status           string           `json:"status"` // pending / processing / done / error
	Message          string           `json:"message"`
}
