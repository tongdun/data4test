package biz

type AppModel struct {
	AppName   string   `json:"appName"`
	Modules   []string `json:"modules"`
	Apis      []string `json:"apis"`
	Prefix    string   `json:"prefix"`
	Methods   []string `json:"methods"`
	ApisDesc  []string `json:"apisDesc"`
	DatasDesc []string `json:"datasDesc"`
}

type RunRespModel struct {
	Response   string `json:"response"`
	Url        string `json:"url"`
	Request    string `json:"request"`
	Header     string `json:"header"`
	TestResult string `json:"testResult"`
	FailReason string `json:"failReason"`
	Output     string `json:"output"`
}

type RunSceneRespModel struct {
	LastFile   string `json:"lastDataFile"`
	TestResult string `gorm:"column:result" json:"testResult"`
	FailReason string `json:"failReason"`
}

type ModuleModel struct {
	AppName  string   `json:"appName"`
	Module   string   `json:"module"`
	Apis     []string `json:"apis"`
	Methods  []string `json:"methods"`
	ApisDesc []string `json:"apisDesc"`
}

type MethodModel struct {
	AppName  string   `json:"appName"`
	Method   string   `json:"method"`
	Apis     []string `json:"apis"`
	Modules  []string `json:"modules"`
	ApisDesc []string `json:"apisDesc"`
	//DatasDesc []string `json:"datasDesc"`
}

type VarDefModel struct {
	Name      string `json:"name"`
	ValueType string `json:"valueType"`
	IsMust    string `json:"isMust"`
	EgValue   string `json:"egValue"`
	Desc      string `json:"desc"`
}

type VarDataModel struct {
	VarDefModel
	TestValue []interface{} `json:"testValue"`
}

type DepDataModel struct {
	DataFile string `json:"dataFile"`
}

type EnvModel struct {
	Product string `json:"product"`
}

type ApiInfoModel struct {
	ApiDefSaveModel
	DatasDesc []string `json:"datasDesc"`
}

type SceneInfoModel struct {
	Name      string         `json:"name"`
	Product   string         `json:"product"`
	SceneType string         `json:"type"`
	RunNum    int            `json:"runNum"`
	DataList  []DepDataModel `json:"dataList"`
}

type ApiDefSaveModel struct {
	App        string        `json:"app"`
	Module     string        `json:"module"`
	ApiDesc    string        `json:"apiDesc"`
	Method     string        `json:"method"`
	BodyMode   string        `json:"bodyMode"`
	Path       string        `json:"path"`
	Prefix     string        `json:"prefix"`
	BodyStr    string        `json:"bodyStr"`
	PathVars   []VarDefModel `json:"pathVars"`
	QueryVars  []VarDefModel `json:"queryVars"`
	BodyVars   []VarDefModel `json:"bodyVars"`
	HeaderVars []VarDefModel `json:"headerVars"`
	RespVars   []VarDefModel `json:"respVars"`
}

type OtherModel struct {
	Version        float64 `json:"version"`
	ApiId          string  `json:"apiId"`
	IsParallel     string  `json:"isParallel"`
	IsUseEnvConfig string  `json:"isUseEnvConfig"`
}

type ApiDataSaveModel struct {
	App        string         `json:"app"`
	Module     string         `json:"module"`
	ApiDesc    string         `json:"apiDesc"`
	DataDesc   string         `json:"dataDesc"`
	Method     string         `json:"method"`
	Path       string         `json:"path"`
	Host       string         `json:"host"`
	Prototype  string         `json:"prototype"`
	BodyMode   string         `json:"bodyMode"`
	Prefix     string         `json:"prefix"`
	Product    string         `json:"product"`
	PathVars   []VarDataModel `json:"pathVars"`
	QueryVars  []VarDataModel `json:"queryVars"`
	BodyVars   []VarDataModel `json:"bodyVars"`
	HeaderVars []VarDataModel `json:"headerVars"`
	RespVars   []VarDataModel `json:"respVars"`
	Actions    []SceneAction  `json:"actions"`
	Asserts    []SceneAssert  `json:"asserts"`
	PreApis    []DepDataModel `json:"preApis"`
	PostApis   []DepDataModel `json:"postApis"`
	Other      []OtherModel   `json:"otherConfig"`
}

type HistorySaveModel struct {
	App        string         `json:"app"`
	Module     string         `json:"module"`
	ApiDesc    string         `json:"apiDesc"`
	DataDesc   string         `json:"dataDesc"`
	Method     string         `json:"method"`
	Path       string         `json:"path"`
	BodyMode   string         `json:"bodyMode"`
	Prefix     string         `json:"prefix"`
	Product    string         `json:"product"`
	Prototype  string         `json:"prototype"`
	Host       string         `json:"host"`
	FileName   string         `json:"fileName"`
	PathVars   []VarDataModel `json:"pathVars"`
	QueryVars  []VarDataModel `json:"queryVars"`
	BodyVars   []VarDataModel `json:"bodyVars"`
	HeaderVars []VarDataModel `json:"headerVars"`
	RespVars   []VarDataModel `json:"respVars"`
	Actions    []SceneAction  `json:"actions"`
	Asserts    []SceneAssert  `json:"asserts"`
	PreApis    []DepDataModel `json:"preApis"`
	PostApis   []DepDataModel `json:"postApis"`
	Other      []OtherModel   `json:"otherConfig"`
	RunRespModel
	//Response string `json:"response"`
	//Url string `json:"url"`
	//Header string `json:"header"`
	//Request string `json:"request"`
	//FailReason string `json:"failReason"`
}

type SceneHistorySaveModel struct {
	//Id        string         `gorm:"column:id" json:"id"`
	Name      string         `json:"name"`
	Product   string         `json:"product"`
	SceneType string         `json:"type"`
	RunNum    int            `json:"runNum"`
	DataList  []DepDataModel `json:"dataList"`
	RunSceneRespModel
}

type HistoryDataDetail struct {
	Id      string `gorm:"column:id" json:"id"`
	Name    string `gorm:"column:name" json:"name"`
	ApiId   string `gorm:"column:api_id" json:"api_id"`
	App     string `gorm:"column:app" json:"app"`
	Content string `gorm:"column:content" json:"content"`
	EnvType int    `gorm:"column:env_type" json:"env_type"`
	Product string `gorm:"column:product" json:"product"`
}
