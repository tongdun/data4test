package biz

type EnvConfig struct {
	Id   string `gorm:"column:id" json:"id"`
	App   string `gorm:"column:app" json:"app"`
	Product   string `gorm:"column:product" json:"product"`
	Ip        string `gorm:"column:ip" json:"ip"`
	Protocol  string `gorm:"column:protocol" json:"protocol"`
	Auth      string `gorm:"column:auth" json:"auth"`
	Prepath   string `gorm:"column:prepath" json:"prepath"`
	Threading string `gorm:"column:threading" json:"threading"`
	MaxThreadNum  int `gorm:"column:thread_number" json:"thread_number"`
	//Usermode  string `gorm:"column:usermode" json:"usermode"`
	Testmode  string `gorm:"column:testmode" json:"testmode"`
	SwaggerPath  string `gorm:"column:swagger_path" json:"swagger_path"`
}
