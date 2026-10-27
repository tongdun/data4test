package biz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CliAuthMiddleware CLI 认证中间件
// 从 Authorization Header 或 form 字段中提取 JWT token 并验证
// 无 token 或 token 无效时返回 401 + 登录方法说明
func CliAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""
		// 1. 优先从 Authorization: Bearer <token> 取
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}
		// 2. 其次从 form 字段 token 取
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}

		if tokenStr == "" {
			// 无 token → 返回 401，提供登录方法说明
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  T("error.cli_token_required"),
				"data": gin.H{
					"login_url":  "/cli/login",
					"login_hint": T("cli.login_hint"),
				},
			})
			c.Abort()
			return
		}

		claims, err := ParseCliToken(tokenStr)
		if err != nil {
			Logger.Warning("CLI auth: token parse failed: %v", err)
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  T("error.cli_token_invalid"),
				"data": gin.H{
					"login_url":  "/cli/login",
					"login_hint": T("cli.login_hint"),
				},
			})
			c.Abort()
			return
		}

		c.Set("cliUserName", claims.UserName)
		c.Set("cliUserID", claims.UserID)
		c.Set("cliName", claims.Name)
		c.Next()
	}
}
