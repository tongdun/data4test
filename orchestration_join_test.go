package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	adminModels "github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/gin-gonic/gin"
)

func TestOrchestrationJoinRequiresEditPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	permission := func(path, method string) adminModels.PermissionModel {
		return adminModels.PermissionModel{HttpPath: []string{path}, HttpMethod: []string{method}}
	}
	for _, test := range []struct {
		name        string
		loggedIn    bool
		permissions []adminModels.PermissionModel
		status      int
	}{
		{name: "anonymous", status: http.StatusUnauthorized},
		{name: "read-only", loggedIn: true, permissions: []adminModels.PermissionModel{permission("/info/task_orchestration", "GET")}, status: http.StatusForbidden},
		{name: "schedule-editor", loggedIn: true, permissions: []adminModels.PermissionModel{permission("/edit/schedule", "POST")}, status: http.StatusForbidden},
		{name: "orchestration-editor", loggedIn: true, permissions: []adminModels.PermissionModel{permission("/edit/task_orchestration", "POST")}, status: http.StatusOK},
		{name: "super-admin", loggedIn: true, permissions: []adminModels.PermissionModel{permission("*", "")}, status: http.StatusOK},
	} {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			writes := 0
			guard := orchestrationJoinAuthorization(func(interface{}) (adminModels.UserModel, bool) {
				return adminModels.UserModel{Id: 1, Name: "tester", Permissions: test.permissions}, test.loggedIn
			})
			router.POST("/join", guard, func(c *gin.Context) { writes++; c.Status(http.StatusOK) })
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/join", nil))
			if response.Code != test.status {
				t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
			}
			if test.status != http.StatusOK && writes != 0 {
				t.Fatal("unauthorized request reached mutation handler")
			}
			if test.status == http.StatusOK && writes != 1 {
				t.Fatal("authorized request was blocked")
			}
		})
	}
}

func TestBothOrchestrationJoinRoutesRejectAnonymousRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	writes := 0
	registerOrchestrationJoinRoutes(router, func(c *gin.Context) { writes++; c.Status(http.StatusOK) })
	for _, path := range []string{"/orchestration_join", "/orchestration_join_row"} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, path, nil))
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("%s status=%d", path, response.Code)
		}
	}
	if writes != 0 {
		t.Fatal("anonymous request reached mutation handler")
	}
}
