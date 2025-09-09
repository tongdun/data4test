package biz

type AiTemplateDefine struct {
	TemplateName       string `gorm:"column:template_name" json:"template_name"`
	TemplateType       string `gorm:"column:template_type" json:"template_type"`
	TemplateContent    string `gorm:"column:template_content" json:"template_content"`
	AppendConversion   string `gorm:"force;column:append_conversion" json:"append_conversion"`
	UseStatus          string `gorm:"column:use_status" json:"use_status"`
	ApplicablePlatform string `gorm:"column:applicable_platform" json:"applicable_platform"`
	CreateUser         string `gorm:"column:create_user" json:"create_user" `
	//ModifyUser         string `gorm:"column:modify_user" json:"modify_user" `
}

type DbAiTemplateDefine struct {
	Id string `gorm:"column:id" json:"id"`
	AiTemplateDefine
}
