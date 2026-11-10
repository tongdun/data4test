package main

import (
	"net/http"
	"net/url"
	"strings"

	"data4test/biz"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	adminModels "github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/gin-gonic/gin"
)

// Both join endpoints mutate task_list and require the existing orchestration edit permission.
func orchestrationJoinAuthorization(getUser func(interface{}) (adminModels.UserModel, bool)) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, loggedIn := getUser(c)
		if !loggedIn || user.Id == 0 || strings.TrimSpace(user.Name) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": biz.T("error.unauthorized")})
			return
		}
		params := url.Values{"id": {c.PostForm("orchestration_id")}}
		if !auth.CheckPermissions(user, config.Url("/edit/task_orchestration"), http.MethodPost, params) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "msg": biz.T("error.permission_denied")})
			return
		}
		c.Next()
	}
}

func registerOrchestrationJoinRoutes(router *gin.Engine, handler gin.HandlerFunc) {
	guard := orchestrationJoinAuthorization(engine.User)
	router.POST("/orchestration_join", guard, handler)
	router.POST("/orchestration_join_row", guard, handler)
}
