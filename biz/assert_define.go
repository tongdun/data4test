package biz

type AssertValueDefine struct {
	Id    string `gorm:"column:id" json:"id"`
	Name  string `gorm:"column:name" json:"name"`
	Value string `gorm:"column:value" json:"value"`
}
