package biz

type AiData struct {
	DataBase
	CreateBase
}

type InputData struct {
	AiTemplate string `gorm:"column:ai_template" json:"ai_template"`
	CommonExtend
}

type DataBase struct {
	Name       string `gorm:"column:name" json:"name"`
	ApiId      string `gorm:"column:api_id" json:"api_id"`
	App        string `gorm:"column:app" json:"app"`
	Content    string `gorm:"column:content" json:"content"`
	FileName   string `gorm:"column:file_name" json:"file_name"`
	FileType   int    `gorm:"column:file_type" json:"file_type"`
	Result     string `gorm:"column:result" json:"result"`
	FailReason string `gorm:"column:fail_reason" json:"fail_reason"`
}

type AnalysisDataInput struct {
	AiTemplate string `gorm:"column:ai_template" json:"ai_template"`
	Product    string `gorm:"column:product" json:"product"`
	CreateBase
}
