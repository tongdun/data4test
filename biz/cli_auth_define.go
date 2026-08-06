package biz

// GoadminUser GoAdmin 用户表结构
type GoadminUser struct {
	Id           uint   `gorm:"column:id"`
	Username     string `gorm:"column:username"`
	Password     string `gorm:"column:password"`
	Name         string `gorm:"column:name"`
	Avatar       string `gorm:"column:avatar"`
	RememberToken string `gorm:"column:remember_token"`
}

// CliLoginRequest CLI 登录请求
type CliLoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// CliLoginResponse CLI 登录响应
type CliLoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// CliClaims JWT 载荷
type CliClaims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	IssuedAt int64  `json:"iat"`
	ExpiresAt int64 `json:"exp"`
	Issuer   string `json:"iss"`
	Subject  string `json:"sub"`
}

// CliJWTConfig JWT 配置
type CliJWTConfig struct {
	Secret     string
	ExpireHour int
}
