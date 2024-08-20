package biz

type AssertValueDefine struct {
	Id      string `gorm:"column:id" json:"id"`
	Name    string `gorm:"column:name" json:"name"`
	Value   string `gorm:"column:value" json:"value"`
	Remark  string `gorm:"column:remark" json:"remark"`
	UerName string `gorm:"column:user_name" json:"user_name"`
}

type DepAssertModel struct {
	Name string `gorm:"column:name" json:"name"`
}
