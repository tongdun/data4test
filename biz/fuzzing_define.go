package biz

type FuzzingDefinition struct {
	Name    string `gorm:"column:name" json:"name"`
	Value   string `gorm:"column:value" json:"value"`
	Type string `gorm:"column:type" json:"type"`
}
