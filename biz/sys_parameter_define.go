package biz

type SysParameter struct {
	Name      string `gorm:"column:name" json:"name"`
	ValueList string `gorm:"column:value_list" json:"value_list"`
	Remark    string `gorm:"column:remark" json:"remark"`
}
