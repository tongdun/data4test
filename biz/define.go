package biz

type ParameterDefinition struct {
	Name  string `gorm:"column:name" json:"name"`
	Value string `gorm:"column:value" json:"value"`
	App   string `gorm:"column:app" json:"app"`
}

type ChinaData struct {
	AddName  string       `json:"addName"`
	Children []*ChinaData `json:"children"`
}
