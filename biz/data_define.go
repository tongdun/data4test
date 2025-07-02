package biz

type SceneData struct {
	FileName string `gorm:"column:file_name" json:"file_name" yaml:"file_name"`
	RunTime  int    `gorm:"column:run_time" json:"run_time" yaml:"run_time"`
	Remark   string `gorm:"column:remark" json:"remark" yaml:"remark"`
	UserName string `gorm:"column:user_name" json:"remark" yaml:"user_name"`
	CommonDataBase
}

type DbSceneData struct {
	Id string `gorm:"column:id" json:"id"`
	SceneData
}

type DataFile struct {
	Name             string                   `json:"name" yaml:"name"`
	Version          float64                  `json:"version" yaml:"version"`
	ApiId            string                   `json:"api_id" yaml:"api_id"`
	IsRunPreApis     string                   `json:"is_run_pre_apis" yaml:"is_run_pre_apis"`
	IsRunPostApis    string                   `json:"is_run_post_apis" yaml:"is_run_post_apis"`
	IsParallel       string                   `json:"is_parallel" yaml:"is_parallel"`
	IsUseEnvConfig   string                   `json:"is_use_env_config" yaml:"is_use_env_config"`
	IsVarStrongCheck string                   `json:"is_var_strong_check,omitempty" yaml:"is_var_strong_check,omitempty"`
	Env              SceneEnv                 `json:"env" yaml:"env"`
	Api              SceneApi                 `json:"api" yaml:"api"`
	Single           SceneSingle              `json:"single" yaml:"single"`
	Multi            SceneMulti               `json:"multi" yaml:"multi"`
	Action           []SceneAction            `json:"action" yaml:"action"`
	Assert           []SceneAssert            `json:"assert" yaml:"assert"`
	Output           map[string][]interface{} `json:"output" yaml:"output"`
	TestResult       []string                 `json:"test_result" yaml:"test_result"`
	FailReason       []string                 `json:"failReason,omitempty" yaml:"failReason,omitempty"`
	Urls             []string                 `json:"urls" yaml:"urls"`
	Request          []string                 `json:"request" yaml:"request"`
	Response         []string                 `json:"response" yaml:"response"`
}

type SceneAssert struct {
	Source string      `json:"source" yaml:"source"`
	Type   string      `json:"type" yaml:"type"`
	Value  interface{} `json:"value" yaml:"value"`
}

type SceneAction struct {
	Type  string      `json:"type" yaml:"type"`
	Value interface{} `json:"value" yaml:"value"`
}

type SceneEnv struct {
	Protocol string `json:"protocol" yaml:"protocol"`
	Host     string `json:"host" yaml:"host"`
	Prepath  string `json:"prepath" yaml:"prepath"`
}

type SceneApi struct {
	Description string   `json:"description" yaml:"description"`
	Module      string   `json:"module" yaml:"module"`
	App         string   `json:"app" yaml:"app"`
	Method      string   `json:"method" yaml:"method"`
	Path        string   `json:"path" yaml:"path"`
	PreApi      []string `json:"pre_apis" yaml:"pre_apis"`
	ParamApis   []string `json:"param_apis" yaml:"param_apis"`
	PostApis    []string `json:"post_apis" yaml:"post_apis"`
}

type SceneSingle struct {
	Header     map[string]interface{} `json:"header" yaml:"header"`
	RespHeader map[string]interface{} `json:"respHeader,omitempty" yaml:"respHeader,omitempty"`
	Path       map[string]interface{} `json:"path" yaml:"path"`
	Query      map[string]interface{} `json:"query" yaml:"query"`
	Body       map[string]interface{} `json:"body" yaml:"body"`
	BodyList   []interface{}          `json:"bodyList,omitempty" yaml:"bodyList,omitempty"`
}

type SceneMulti struct {
	//Header     map[string][]string `json:"header" yaml:"header"`  // 如有需要，再适配
	Path  map[string][]interface{} `json:"path" yaml:"path"`
	Query map[string][]interface{} `json:"query" yaml:"query"`
	Body  map[string][]interface{} `json:"body" yaml:"body"`
}

type SceneDataRecord struct {
	Name       string `gorm:"column:name" json:"name"`
	ApiId      string `gorm:"column:api_id" json:"api_id"`
	App        string `gorm:"column:app" json:"app"`
	Content    string `gorm:"column:content" json:"content" yaml:"content"`
	Result     string `gorm:"column:result" json:"result" yaml:"result"`
	FailReason string `gorm:"column:fail_reason" json:"fail_reason" yaml:"fail_reason"`
	EnvType    int    `gorm:"column:env_type" json:"env_type" yaml:"env_type"`
	Product    string `gorm:"column:product" json:"product" yaml:"product"`
}

type SceneHistoryRecord struct {
	Id   string `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type HistoryModel struct {
	Tag  string   `json:"tag"`
	List []string `json:"list"`
}
