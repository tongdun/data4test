package biz

type ProductEnvModel struct {
	Name     string `gorm:"column:product" json:"product"`
	Ip       string `gorm:"column:ip" json:"ip"`
	Protocol string `gorm:"column:protocol" json:"protocol"`
	//Auth string `gorm:"column:auth" json:"auth"`
	Auth []VarDataModel `gorm:"column:auth" json:"auth"`
}

type Product struct {
	Name             string `gorm:"column:product" json:"product"`
	Ip               string `gorm:"column:ip" json:"ip"`
	Protocol         string `gorm:"column:protocol" json:"protocol"`
	Threading        string `gorm:"column:threading" json:"threading"`
	ThreadNumber     int    `gorm:"column:thread_number" json:"thread_number"`
	Auth             string `gorm:"column:auth" json:"auth"`
	Apps             string `gorm:"column:apps" json:"apps"`
	Testmode         string `gorm:"column:testmode" json:"testmode"`
	EnvType          int    `gorm:"column:env_type" json:"env_type"`
	PrivateParameter string `gorm:"column:private_parameter" json:"private_parameter"`
}

type DbProduct struct {
	Id string `gorm:"column:id" json:"id"`
	Product
}
