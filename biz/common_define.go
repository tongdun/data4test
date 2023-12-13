package biz

type IDNoInfo struct {
	Sex      int    `gorm:"column:sex" json:"sex" yaml:"sex"`
	Code     int    `gorm:"column:code" json:"code" yaml:"code"`
	District string `gorm:"column:district" json:"district" yaml:"district"`
	City     string `gorm:"column:city" json:"city" yaml:"city"`
	Province string `gorm:"column:province" json:"province" yaml:"province"`
	Birthday string `gorm:"column:birthday" json:"birthday" yaml:"birthday"`
}
