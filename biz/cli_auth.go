package biz

import (
	"crypto/hmac"
	"crypto/sha256"
	"data4test/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GetCliJWTConfig 获取 JWT 配置
// Secret 优先从环境变量 CLI_JWT_SECRET 读取，fallback 为默认值
func GetCliJWTConfig() CliJWTConfig {
	secret := os.Getenv("CLI_JWT_SECRET")
	if secret == "" {
		secret = "data4test-cli-default-secret"
	}
	return CliJWTConfig{
		Secret:     secret,
		ExpireHour: 24,
	}
}

// base64URLEncode 安全的 URL 编码（无填充）
func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

// base64URLDecode 解码 URL 安全的 base64
func base64URLDecode(s string) ([]byte, error) {
	// 补齐填充
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

// CliLogin CLI 登录：验证用户名密码，返回 JWT token
func CliLogin(username, password string) (*CliLoginResponse, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, fmt.Errorf(T("error.cli_username_password_required"))
	}

	var user GoadminUser
	if err := models.Orm.Table("goadmin_users").Where("username = ?", username).First(&user).Error; err != nil {
		Logger.Warning("CLI login: user not found: %s", username)
		return nil, fmt.Errorf(T("error.cli_login_failed"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		Logger.Warning("CLI login: password mismatch for user: %s", username)
		return nil, fmt.Errorf(T("error.cli_login_failed"))
	}

	cfg := GetCliJWTConfig()
	now := time.Now()
	expiresAt := now.Add(time.Duration(cfg.ExpireHour) * time.Hour)

	claims := CliClaims{
		UserID:    user.Id,
		UserName:  user.Username,
		Name:      user.Name,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
		Issuer:    "data4test-cli",
		Subject:   user.Username,
	}

	tokenStr, err := signJWT(claims, cfg.Secret)
	if err != nil {
		Logger.Error("CLI login: failed to sign token: %v", err)
		return nil, fmt.Errorf(T("error.cli_login_failed"))
	}

	return &CliLoginResponse{
		Token:     tokenStr,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

// signJWT 手动签名 JWT：header.payload.signature
func signJWT(claims CliClaims, secret string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerB64 := base64URLEncode(headerJSON)
	payloadB64 := base64URLEncode(payloadJSON)
	signingInput := headerB64 + "." + payloadB64

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	signature := base64URLEncode(mac.Sum(nil))

	return signingInput + "." + signature, nil
}

// ParseCliToken 解析并验证 JWT token
func ParseCliToken(tokenStr string) (*CliClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	signature := parts[2]

	// 验证签名
	cfg := GetCliJWTConfig()
	mac := hmac.New(sha256.New, []byte(cfg.Secret))
	mac.Write([]byte(signingInput))
	expectedSig := base64URLEncode(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return nil, fmt.Errorf("invalid token signature")
	}

	// 解析 payload
	payloadJSON, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token payload encoding: %v", err)
	}

	var claims CliClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("invalid token payload: %v", err)
	}

	// 验证过期时间
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}
