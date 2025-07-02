package biz

type StatusMgmt struct {
	UseStatus    int `gorm:"column:use_status" json:"use_status"`
	ModifyStatus int `gorm:"column:modify_status" json:"modify_status"`
}

type ImportCommon struct {
	ConversationId string `gorm:"column:conversation_id" json:"conversation_id"`
	RawReply       string `gorm:"column:raw_reply" json:"raw_reply"`
	CommonExtend
}

type CommonExtend struct {
	InputBase
	CreateBase
}

type InputBase struct {
	IntroVersion string `gorm:"column:intro_version" json:"intro_version"`
	Product      string `gorm:"column:product" json:"product"`
}

type CreateBase struct {
	CreatePlatform string `gorm:"column:source" json:"create_platform"`
	CreateUser     string `gorm:"column:create_user" json:"create_user"`
}
