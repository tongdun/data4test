package biz

type ApiDefinition struct {
	ApiId      string `gorm:"column:api_id" json:"api_id"`
	ApiModule  string `gorm:"column:api_module" json:"api_module"`
	ApiDesc    string `gorm:"column:api_desc" json:"api_desc"`
	HttpMethod string `gorm:"column:http_method" json:"http_method"`
	Path       string `gorm:"column:path" json:"path"`
	QueryParameterModel
	Version            int    `gorm:"column:version" json:"version"`
	Check              string `gorm:"column:check" json:"check"`
	ApiStatus          int    `gorm:"column:api_status" json:"api_status"`
	IsAuto             int    `gorm:"column:is_auto" json:"is_auto"`
	IsNeedAuto         int    `gorm:"column:is_need_auto" json:"is_need_auto"`
	ChangeContent      string `gorm:"column:change_content" json:"change_content"`
	ApiCheckFailReason string `gorm:"column:api_check_fail_reason" json:"api_check_fail_reason"`
	App                string `gorm:"column:app" json:"app"`
	Remark             string `gorm:"column:remark" json:"remark"`
}

type ApiStringDefinition struct {
	ApiId      string `gorm:"column:api_id" json:"api_id"`
	ApiModule  string `gorm:"column:api_module" json:"api_module"`
	ApiDesc    string `gorm:"column:api_desc" json:"api_desc"`
	HttpMethod string `gorm:"column:http_method" json:"http_method"`
	Path       string `gorm:"column:path" json:"path"`
	QueryParameterString
	Version            int    `gorm:"column:version" json:"version"`
	Check              string `gorm:"column:check" json:"check"`
	ApiStatus          int    `gorm:"column:api_status" json:"api_status"`
	IsAuto             string `gorm:"column:is_auto" json:"is_auto"`
	IsNeedAuto         string `gorm:"column:is_need_auto" json:"is_need_auto"`
	ChangeContent      string `gorm:"column:change_content" json:"change_content"`
	ApiCheckFailReason string `gorm:"column:api_check_fail_reason" json:"api_check_fail_reason"`
	App                string `gorm:"column:app" json:"app"`
	Remark             string `gorm:"column:remark" json:"remark"`
}

type QueryParameterModel struct {
	Header         []VarDefModel `gorm:"column:header" json:"header"`
	QueryParameter []VarDefModel `gorm:"column:query_parameter" json:"query_parameter"`
	PathVariable   []VarDefModel `gorm:"column:path_variable" json:"path_variable"`
	Body           []VarDefModel `gorm:"column:body" json:"body"`
	Response       []VarDefModel `gorm:"column:response" json:"response"`
}

type QueryParameterString struct {
	Header         string `gorm:"column:header" json:"header"`
	QueryParameter string `gorm:"column:query_parameter" json:"query_parameter"`
	PathVariable   string `gorm:"column:path_variable" json:"path_variable"`
	Body           string `gorm:"column:body" json:"body"`
	Response       string `gorm:"column:response" json:"response"`
}

type ApiDefDB struct {
	ApiId          string `gorm:"column:api_id" json:"api_id"`
	ApiModule      string `gorm:"column:api_module" json:"api_module"`
	ApiDesc        string `gorm:"column:api_desc" json:"api_desc"`
	HttpMethod     string `gorm:"column:http_method" json:"http_method"`
	Path           string `gorm:"column:path" json:"path"`
	Header         string `gorm:"column:header" json:"header"`
	QueryParameter string `gorm:"column:query_parameter" json:"query_parameter"`
	PathVariable   string `gorm:"column:path_variable" json:"path_variable"`
	Body           string `gorm:"column:body" json:"body"`
	Response       string `gorm:"column:response" json:"response"`
	Version        int    `gorm:"column:version" json:"version"`
	Check          string `gorm:"column:check" json:"check"`
	ApiStatus      int    `gorm:"column:api_status" json:"api_status"`
	ChangeContent  string `gorm:"column:change_content" json:"change_content"`
	App            string `gorm:"column:app" json:"app"`
	Remark         string `gorm:"column:remark" json:"remark"`
}

type ApiTestDetail struct {
	ApiId      string `gorm:"column:api_id" json:"api_id"`
	ApiDesc    string `gorm:"column:api_desc" json:"api_desc"`
	DataDesc   string `gorm:"column:data_desc" json:"data_desc"`
	Header     string `gorm:"column:header" json:"header"`
	Url        string `gorm:"column:url" json:"url"`
	Body       string `gorm:"column:body" json:"body"`
	Response   string `gorm:"column:response" json:"response"`
	TestResult string `gorm:"column:test_result" json:"test_result"`
	FailReason string `gorm:"column:fail_reason" json:"fail_reason"`
	CreatedAt  string `gorm:"column:created_at" json:"created_at"`
	App        string `gorm:"column:app" json:"app"`
}

type ApiTestData struct {
	DataDesc       string `gorm:"column:data_desc" json:"data_desc"`
	ApiDesc        string `gorm:"column:api_desc" json:"api_desc"`
	ApiId          string `gorm:"column:api_id" json:"api_id"`
	ApiModule      string `gorm:"column:api_module" json:"api_module"`
	Header         string `gorm:"column:header" json:"header"`
	UrlQuery       string `gorm:"column:url_query" json:"url_query"`
	Body           string `gorm:"column:body" json:"body"`
	RunNum         int    `gorm:"column:run_num" json:"run_num"`
	ExpectedResult string `gorm:"column:expected_result" json:"expected_result"`
	ActualResult   string `gorm:"column:actual_result" json:"actual_result"`
	Result         string `gorm:"column:result" json:"result"`
	FailReason     string `gorm:"column:fail_reason" json:"fail_reason"`
	Response       string `gorm:"column:response" json:"response"`
	App            string `gorm:"column:app" json:"app"`
}

type ApiFuzzingData struct {
	DataDesc       string `gorm:"column:data_desc" json:"data_desc"`
	ApiDesc        string `gorm:"column:api_desc" json:"api_desc"`
	ApiId          string `gorm:"column:api_id" json:"api_id"`
	ApiModule      string `gorm:"column:api_module" json:"api_module"`
	Header         string `gorm:"column:header" json:"header"`
	UrlQuery       string `gorm:"column:url_query" json:"url_query"`
	Body           string `gorm:"column:body" json:"body"`
	RunNum         int    `gorm:"column:run_num" json:"run_num"`
	ExpectedResult string `gorm:"column:expected_result" json:"expected_result"`
	ActualResult   string `gorm:"column:actual_result" json:"actual_result"`
	Result         string `gorm:"column:result" json:"result"`
	FailReason     string `gorm:"column:fail_reason" json:"fail_reason"`
	Response       string `gorm:"column:response" json:"response"`
	App            string `gorm:"column:app" json:"app"`
}

type DbApiTestResult struct {
	Id string `gorm:"column:id" json:"id"`
	ApiTestResult
}

type ApiTestResult struct {
	ApiId   string `gorm:"column:api_id" json:"api_id"`
	OutVars string `gorm:"column:out_vars" json:"out_vars"`
	// Result      string `gorm:"column:result" json:"result"`
	UpdatedAt string `gorm:"column:updated_at" json:"updated_at"`
	App       string `gorm:"column:app" json:"app"`
}

type ApiRelation struct {
	// gorm.Model
	ApiId     string `gorm:"column:api_id" json:"api_id"`
	ApiDesc   string `gorm:"column:api_desc" json:"api_desc"`
	ApiModule string `gorm:"column:api_module" json:"api_module"`
	App       string `gorm:"column:app" json:"app"`
	Auto      string `gorm:"column:auto" json:"auto"`
	PreApis   string `gorm:"column:pre_apis" json:"pre_apis"`
	OutVars   string `gorm:"column:out_vars" json:"out_vars"`
	CheckVars string `gorm:"column:check_vars" json:"check_vars"`
	ParamApis string `gorm:"column:param_apis" json:"param_apis"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}
