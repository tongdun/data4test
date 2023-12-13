package biz

type Swagger struct {
	Paths       map[string]PathDef   `json:"paths"`
	Definitions map[string]DefiniDef `json:"definitions,omitempty"`
	Components  ComponentDef         `json:"components,omitempty"`
}

type PathDef struct {
	Put    ApiDetail `json:"put"`
	Post   ApiDetail `json:"post"`
	Get    ApiDetail `json:"get"`
	Delete ApiDetail `json:"delete"`
}

type ApiDetail struct {
	Description string         `json:"description"`
	Consumes    []string       `json:"consumes"`
	Produces    []string       `json:"produces"`
	Tags        []string       `json:"tags"`
	Summary     string         `json:"summary"`
	Parameters  []ParamDetail  `json:"parameters,omitempty"`
	RequestBody RequestContent `json:"requestBody,omitempty"`
	Responses   ResponseDetail `json:"responses"`
}

type ParamDetail struct {
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Name        string       `json:"name"`
	In          string       `json:"in"`
	Required    bool         `json:"required"`
	Schema      SchemaDetail `json:"schema"`
}

type ComponentDef struct {
	ComSchema map[string]ComSchemaDetail `json:"schemas"`
}

type ComSchemaDetail struct {
	Type       string                  `json:"type"`
	Required   []string                `json:"required,omitempty"`
	Properties map[string]ProperDetail `json:"properties"`
}

type DefiniDef struct {
	Type       string                  `json:"type"`
	Required   []string                `json:"required"`
	Properties map[string]ProperDetail `json:"properties"`
}

type ProperDetail struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Ref         string `json:"$ref,omitempty"`
}

type SchemaDetail struct {
	Type       string                  `json:"type"`
	Ref        string                  `json:"$ref,omitempty"`
	Properties map[string]ProperDetail `json:"properties,omitempty"`
}

type ContentDetail struct {
	Star StarDetail `json:"*/*"`
}

type RequestContent struct {
	RequestContent RequestBodyDetail `json:"content"`
}

type RequestBodyDetail struct {
	JSONContent      StarDetail `json:"application/json"`
	MultipartContent StarDetail `json:"multipart/form-data"`
}

type StarDetail struct {
	Schema SchemaDetail `json:"schema"`
}

type ResponseDetail struct {
	R200 R200Detail `json:"200"`
}

type R200Detail struct {
	Description string        `json:"description"`
	Schema      SchemaDetail  `json:"schema,omitempty"`
	Content     ContentDetail `json:"content,omitempty"`
}

type Response struct {
	Status        string                 `json:"status,omitempty"`
	Message       string                 `json:"message,omitempty"`
	Content       map[string]interface{} `json:"content,omitempty"`
	IsSuccess     bool                   `json:"isSuccess,omitempty"`
	ResultMessage string                 `json:"resultMessage,omitempty"`
	ResultObject  map[string]interface{} `json:"resultObject,omitempty"`
}

type DbApiDefinition struct {
	Id string `gorm:"column:id" json:"id"`
	ApiDefinition
}

type DbApiStringDefinition struct {
	Id string `gorm:"column:id" json:"id"`
	ApiStringDefinition
}

type DbApiRelation struct {
	Id string `gorm:"column:id" json:"id"`
	ApiRelation
}

type VarMapList map[string][]VarDefModel
type DoubleMap map[string]map[string]string
type SingleMap map[string]string

type SimpleCase struct {
	Module     string `gorm:"column:module" json:"module"`
	TestResult string `gorm:"column:test_result" json:"test_result"`
	App        string `gorm:"column:app" json:"app"`
}
