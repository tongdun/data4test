package main

import (
	"encoding/json"
	"github.com/mritd/chinaid"
	"runtime/debug"
	"strconv"
	"time"

	//"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"

	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"

	"data4test/login"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	//_ "github.com/mkevac/debugcharts"
	//_ "net/http/pprof"

	ada "github.com/GoAdminGroup/go-admin/adapter/gin"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // sql driver
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
	"github.com/GoAdminGroup/go-admin/plugins"
	_ "github.com/GoAdminGroup/themes/sword" // ui theme

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/auth"

	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"

	"data4test/biz"
	"data4test/models"
	"data4test/pages"
	"data4test/tables"

	"github.com/GoAdminGroup/filemanager"
	plugModels "github.com/GoAdminGroup/go-admin/plugins/admin/models"
	"github.com/GoAdminGroup/librarian"
)

type Args struct {
	Config string
	Upload string
}

var (
	args Args
	ctx  *gin.Context
)

func InstallFlags() {
	flag.StringVar(&args.Config, "c", "./config.json", "Config file")
	flag.StringVar(&args.Upload, "u", "./uploads", "Upload file")
	return
}

func init() {
	InstallFlags()
	if _, err := os.Stat(args.Config); os.IsNotExist(err) {
		panic(err)
	}

	if _, err := os.Stat(args.Upload); os.IsNotExist(err) {
		os.Mkdir(args.Upload, os.ModePerm)
	}
	return
}

func main() {
	flag.Parse()
	defer func() {
		err := recover()
		if err != nil {
			biz.Logger.Error("panic: %v", err)
			biz.Logger.Error("stack: %v", string(debug.Stack()))
			startServer()
		}
	}()
	startServer()
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	const visitorRoleID = int64(3)
	var err error
	r := gin.Default()
	r.Use(cors.Default())
	//pprof.Register(r)  // 性能查看

	eng := engine.Default()

	template.AddLoginComp(login.Get())
	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(args.Config).
		AddGenerators(tables.Generators).
		AddPlugins(librarian.NewLibrarianWithConfig(librarian.Config{
			Path:           biz.DocFilePath,
			MenuUserRoleID: visitorRoleID,
			Prefix:         "librarian",
			BuildMenu:      true, // auto build menus
		})).
		Use(r); err != nil {
		panic(err)
	}

	r.Static("./upload", args.Upload)
	r.Static("./static", "./web/static")

	eng.HTMLFile("GET", "/admin/likePostman", "./html/index.html", nil)

	eng.Data("GET", "/admin/librarian", func(ctx *context.Context) {
		conn := eng.SqliteConnection()
		user := plugModels.User().SetConn(conn).Find(visitorRoleID)
		_ = auth.SetCookie(ctx, user, conn)
		ctx.Redirect("/admin/librarian/README")
		return
	}, true)

	r.GET("/", func(ctx *gin.Context) {
		user, _ := engine.User(ctx)
		var roleName string
		if len(user.Roles) > 0 {
			roleName = user.Roles[0].Name
		}
		if _, ok := biz.REDIRECT_PATH[roleName]; ok {
			ctx.Redirect(http.StatusMovedPermanently, biz.REDIRECT_PATH[roleName])
		} else {
			ctx.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
		}

		return
	})

	r.GET("/admin", func(ctx *gin.Context) {
		user, _ := engine.User(ctx)
		var roleName string
		if len(user.Roles) > 0 {
			roleName = user.Roles[0].Name
		}

		if _, ok := biz.REDIRECT_PATH[roleName]; ok {
			ctx.Redirect(http.StatusMovedPermanently, biz.REDIRECT_PATH[roleName])
		} else {
			ctx.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
		}

		return

	})

	r.GET("/admin/dashboard", ada.Content(pages.GetDashBoardContent))

	r.GET("/admin/product_dashboard", ada.Content(func(ctx *gin.Context) (panel types.Panel, e error) {
		return pages.GetDashBoard2Content(ctx)
	}))

	r.GET("/admin/app_dashboard", ada.Content(func(ctx *gin.Context) (panel types.Panel, e error) {
		return pages.GetDashBoard3Content(ctx)
	}))

	plug, _ := plugins.FindByName("filemanager")

	plug.(*filemanager.FileManager).SetPathValidator(func(path string) error {
		if !strings.Contains(path, "mgmt") {
			return errors.New("没有权限")
		}
		return nil
	})

	r.GET("/appList", func(c *gin.Context) {
		appList := biz.GetAppList()
		data := map[string]interface{}{
			"data": appList,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/dataList", func(c *gin.Context) {
		dataDesc := c.Query("dataDesc")
		data := make(map[string]interface{})
		if len(dataDesc) == 0 {
			appList := biz.GetDataList()
			if len(appList) == 0 {
				data["code"] = 400
				data["msg"] = "未关联到数据"
			} else {
				data["data"] = appList
				data["code"] = 200
				data["msg"] = "操作成功"
			}
		} else {
			dataInfo := biz.GetDataInfoByDataDesc("", "", "", dataDesc)
			if len(dataInfo.Path) == 0 {
				data["code"] = 400
				data["msg"] = "未关联到数据"
			} else {
				if dataInfo.BodyMode == "json" {
					//dataStr
				}
				data["data"] = dataInfo
				data["code"] = 200
				data["msg"] = "操作成功"
			}
		}

		c.JSON(http.StatusOK, data)
	})

	r.GET("/dataFileList", func(c *gin.Context) {
		dataFileList := biz.GetDataFileList()
		data := map[string]interface{}{
			"data": dataFileList,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/assertTemplateList", func(c *gin.Context) {
		assertTemplateList := biz.GetAssertTemplateList()
		data := map[string]interface{}{
			"data": assertTemplateList,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/envList", func(c *gin.Context) {
		dataFileList := biz.GetEnvList()
		data := map[string]interface{}{
			"data": dataFileList,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/envList/:product", func(c *gin.Context) {
		productName := c.Param("product")
		product, err := biz.GetProductEnv(productName)
		data := make(map[string]interface{})
		if err != nil {
			data = map[string]interface{}{
				"data": product,
				"code": 400,
				"msg":  err,
			}
		} else {
			data = map[string]interface{}{
				"data": product,
				"code": 200,
				"msg":  "操作成功",
			}
		}

		c.JSON(http.StatusOK, data)
	})

	r.GET("/allMenu", func(c *gin.Context) {
		appName := c.Query("app")
		allMenus := biz.GetApiMenu(appName)
		data := map[string]interface{}{
			"data": allMenus,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/dataMenu", func(c *gin.Context) {
		appName := c.Query("app")
		allMenus := biz.GetDataMenu(appName)
		data := map[string]interface{}{
			"data": allMenus,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/sceneMenu", func(c *gin.Context) {
		product := c.Query("product")
		allMenus := biz.GetSceneMenu(product)
		data := map[string]interface{}{
			"data": allMenus,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/historyMenu", func(c *gin.Context) {
		dateName := c.Query("dateName")
		allMenus := biz.GetHistoryMenu(dateName)
		data := map[string]interface{}{
			"data": allMenus,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/sceneHistoryMenu", func(c *gin.Context) {
		dateName := c.Query("dateName")
		allMenus := biz.GetSceneHistoryMenu(dateName)
		data := map[string]interface{}{
			"data": allMenus,
			"code": 200,
			"msg":  "操作成功",
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/appList/:appName", func(c *gin.Context) {
		appName := c.Param("appName")
		method := c.Query("method")
		path := c.Query("path")
		apiDesc := c.Query("apiDesc")
		dataDesc := c.Query("dataDesc")
		module := c.Query("module")
		mode := c.Query("mode")
		data := make(map[string]interface{})

		if mode == "def" {
			if len(apiDesc) > 0 {
				apiModel := biz.GetApiDataByApiDesc(appName, module, apiDesc)
				data["data"] = apiModel
			} else if len(appName) > 0 && len(method) > 0 && len(path) > 0 && len(appName) > 0 {
				apiModel := biz.GetApiInfo(appName, method, path)
				data["data"] = apiModel
			} else if len(appName) > 0 && len(module) > 0 && len(apiDesc) > 0 {
				apiModel := biz.GetApiDataByApiDesc(appName, module, apiDesc)
				data["data"] = apiModel
			} else if len(appName) > 0 && len(module) > 0 {
				moduleModel := biz.GetModuleInfo(appName, module)
				data["data"] = moduleModel
			} else if len(appName) > 0 && len(method) > 0 {
				methodModel := biz.GetMethodInfo(appName, method)
				data["data"] = methodModel
			} else if len(appName) > 0 {
				appModel := biz.GetAppInfo(appName)
				data["data"] = appModel
			}
		} else if mode == "run" {
			if len(dataDesc) > 0 {
				dataInfo := biz.GetDataInfoByDataDesc(appName, module, apiDesc, dataDesc)
				data["data"] = dataInfo
			} else if len(module) > 0 && len(apiDesc) > 0 && len(appName) > 0 && len(dataDesc) > 0 {
				dataInfo := biz.GetDataInfoByDataDesc(appName, module, apiDesc, dataDesc)
				data["data"] = dataInfo
			} else if len(method) > 0 && len(path) > 0 && len(appName) > 0 && len(dataDesc) > 0 {
				dataInfo := biz.GetDataInfo(appName, method, path, dataDesc)
				data["data"] = dataInfo
			} else if len(method) > 0 && len(path) > 0 && len(appName) > 0 && len(apiDesc) > 0 {
				dataDescList, dataFileList := biz.GetApiDataList(appName, method, path)
				dataMap := make(map[string]interface{})
				dataMap["dataDescList"] = dataDescList
				dataMap["dataFileList"] = dataFileList
				data["data"] = dataMap
			} else if len(method) > 0 && len(path) > 0 && len(appName) > 0 {
				apiModel := biz.GetApiInfo(appName, method, path)
				data["data"] = apiModel
			} else if len(dataDesc) > 0 {
				dataInfo := biz.GetDataInfoByDataDesc(appName, module, apiDesc, dataDesc)
				data["data"] = dataInfo
			} else if len(apiDesc) > 0 {
				apiModel := biz.GetApiDataByApiDesc(appName, module, apiDesc)
				data["data"] = apiModel
			} else if len(appName) > 0 {
				appModel := biz.GetAppInfo(appName)
				data["data"] = appModel
			}
		} else {
			appModel := biz.GetAppInfo(appName)
			data["data"] = appModel
		}

		data["code"] = 200
		data["msg"] = "操作成功"
		c.JSON(http.StatusOK, data)
	})

	r.GET("/historyList", func(c *gin.Context) {
		fileName := c.Query("fileName")
		data := make(map[string]interface{})

		apiModel, err := biz.GetHistoryByFileName(fileName)
		if err != nil {
			data["code"] = 400
			data["msg"] = err
		} else {
			data["data"] = apiModel
			data["code"] = 200
			data["msg"] = "操作成功"
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/sceneHistoryList", func(c *gin.Context) {
		name := c.Query("name")
		data := make(map[string]interface{})
		sceneModel, err := biz.GetSceneHistory(name)
		if err != nil {
			data["code"] = 400
			data["msg"] = err
		} else {
			data["data"] = sceneModel
			data["code"] = 200
			data["msg"] = "操作成功"
		}
		c.JSON(http.StatusOK, data)
	})

	r.GET("/sceneList", func(c *gin.Context) {
		data := make(map[string]interface{})
		dataList, err := biz.GetAllPlaybook()
		if err != nil {
			data["code"] = 400
			data["msg"] = "获取关联场景失败"
		} else {
			data["code"] = 200
			data["msg"] = "操作成功"
			data["data"] = dataList
		}

		c.JSON(http.StatusOK, data)
	})

	r.GET("/sceneList/:name", func(c *gin.Context) {
		name := c.Param("name")
		product := c.Query("product")

		data := make(map[string]interface{})

		dataList, err := biz.GetPlaybookByName(name, product)
		if err != nil {
			data["code"] = 400
			data["msg"] = "获取场景关联数据失败"
		} else {
			data["code"] = 200
			data["msg"] = "操作成功"
			data["data"] = dataList
		}

		c.JSON(http.StatusOK, data)
	})

	r.POST("/apiDefSave", func(c *gin.Context) {
		var apiDefSave biz.ApiDefSaveModel
		apiDefSave.App = c.PostForm("app")[1 : len(c.PostForm("app"))-1]
		apiDefSave.Module = c.PostForm("module")[1 : len(c.PostForm("module"))-1]
		apiDefSave.ApiDesc = c.PostForm("apiDesc")[1 : len(c.PostForm("apiDesc"))-1]
		apiDefSave.Path = c.PostForm("path")[1 : len(c.PostForm("path"))-1]
		apiDefSave.Prefix = c.PostForm("prefix")[1 : len(c.PostForm("prefix"))-1]
		apiDefSave.Method = c.PostForm("method")[1 : len(c.PostForm("method"))-1]
		apiDefSave.BodyMode = c.PostForm("bodyMode")[1 : len(c.PostForm("bodyMode"))-1]
		apiDefSave.BodyStr = c.PostForm("bodyStr")[1 : len(c.PostForm("bodyStr"))-1]

		json.Unmarshal([]byte(c.PostForm("pathVars")), &apiDefSave.PathVars)
		json.Unmarshal([]byte(c.PostForm("queryVars")), &apiDefSave.QueryVars)
		json.Unmarshal([]byte(c.PostForm("bodyVars")), &apiDefSave.BodyVars)
		json.Unmarshal([]byte(c.PostForm("headerVars")), &apiDefSave.HeaderVars)
		json.Unmarshal([]byte(c.PostForm("respVars")), &apiDefSave.RespVars)

		data := make(map[string]interface{})

		err = biz.SaveApiDef(apiDefSave)
		//biz.ModifyPlaybookApiList()   //更新数据使用，无需要时进行屏蔽
		//biz.UpdateDataAssertContent() //更新断言数据使用，无需要时进行屏蔽
		if err != nil {
			data["code"] = 400
			data["msg"] = err
		} else {
			data["code"] = 200
			data["msg"] = "操作成功"
		}

		c.JSON(http.StatusOK, data)
	})

	r.POST("/apiDataSave", func(c *gin.Context) {
		// 获取登录用户
		user, _ := engine.User(c)
		userName := user.Name

		var apiDataSave biz.ApiDataSaveModel
		apiDataSave.App = c.PostForm("app")[1 : len(c.PostForm("app"))-1]
		apiDataSave.Module = c.PostForm("module")[1 : len(c.PostForm("module"))-1]
		apiDataSave.ApiDesc = c.PostForm("apiDesc")[1 : len(c.PostForm("apiDesc"))-1]
		apiDataSave.DataDesc = c.PostForm("dataDesc")[1 : len(c.PostForm("dataDesc"))-1]
		apiDataSave.Path = c.PostForm("path")[1 : len(c.PostForm("path"))-1]
		apiDataSave.Prefix = c.PostForm("prefix")[1 : len(c.PostForm("prefix"))-1]
		apiDataSave.Method = c.PostForm("method")[1 : len(c.PostForm("method"))-1]
		apiDataSave.Product = c.PostForm("product")[1 : len(c.PostForm("product"))-1]
		apiDataSave.BodyMode = c.PostForm("bodyMode")[1 : len(c.PostForm("bodyMode"))-1]
		bodyStr := c.PostForm("bodyStr")[1 : len(c.PostForm("bodyStr"))-1]

		json.Unmarshal([]byte(c.PostForm("pathVars")), &apiDataSave.PathVars)
		json.Unmarshal([]byte(c.PostForm("queryVars")), &apiDataSave.QueryVars)
		json.Unmarshal([]byte(c.PostForm("bodyVars")), &apiDataSave.BodyVars)
		json.Unmarshal([]byte(c.PostForm("headerVars")), &apiDataSave.HeaderVars)
		json.Unmarshal([]byte(c.PostForm("respVars")), &apiDataSave.RespVars)
		json.Unmarshal([]byte(c.PostForm("actions")), &apiDataSave.Actions)
		json.Unmarshal([]byte(c.PostForm("asserts")), &apiDataSave.Asserts)
		json.Unmarshal([]byte(c.PostForm("preApis")), &apiDataSave.PreApis)
		json.Unmarshal([]byte(c.PostForm("postApis")), &apiDataSave.PostApis)
		json.Unmarshal([]byte(c.PostForm("otherConfig")), &apiDataSave.Other)

		data := make(map[string]interface{})
		if len(bodyStr) > 0 {
			apiDataSave.BodyVars, err = biz.Str2VarModel(bodyStr)
			if err != nil {
				biz.Logger.Debug("err: %v", err)
				data["code"] = 400
				data["msg"] = err
				biz.Logger.Debug("data: %v", data)
				c.JSON(http.StatusOK, data)
			} else {
				err := biz.SaveApiData(apiDataSave, userName)

				if err != nil {
					data["code"] = 400
					data["msg"] = err
				} else {
					data["code"] = 200
					data["msg"] = "操作成功"
				}

				c.JSON(http.StatusOK, data)
			}
		} else {
			err := biz.SaveApiData(apiDataSave, userName)

			if err != nil {
				data["code"] = 400
				data["msg"] = err
			} else {
				data["code"] = 200
				data["msg"] = "操作成功"
			}

			c.JSON(http.StatusOK, data)
		}

	})

	r.POST("/sceneSave", func(c *gin.Context) {
		// 获取登录用户
		user, _ := engine.User(c)
		userName := user.Name

		var sceneSave biz.SceneSaveModel
		sceneSave.Product = c.PostForm("product")[1 : len(c.PostForm("product"))-1]
		sceneSave.Name = c.PostForm("name")[1 : len(c.PostForm("name"))-1]
		typeTag := c.PostForm("type")[1 : len(c.PostForm("type"))-1]
		runNumTag := c.PostForm("runNum")
		runNum, err := strconv.Atoi(runNumTag)
		if err != nil {
			sceneSave.RunNum = 1
		} else {
			sceneSave.RunNum = runNum
		}

		switch typeTag {
		case "串行中断":
			sceneSave.SceneType = 1
		case "串行比较":
			sceneSave.SceneType = 2
		case "串行继续":
			sceneSave.SceneType = 3
		case "普通并发":
			sceneSave.SceneType = 4
		case "并发比较":
			sceneSave.SceneType = 5
		default:
			sceneSave.SceneType = 1
		}

		json.Unmarshal([]byte(c.PostForm("dataList")), &sceneSave.DataList)
		data := make(map[string]interface{})

		if len(sceneSave.Name) == 0 {
			data["code"] = 400
			data["msg"] = "场景名称不能为空"
		} else {
			err := biz.SaveScene(sceneSave, userName)
			if err != nil {
				data["code"] = 400
				data["msg"] = err
			} else {
				data["code"] = 200
				data["msg"] = "操作成功"
			}
		}

		c.JSON(http.StatusOK, data)
	})

	r.POST("/historySave", func(c *gin.Context) {
		// 获取登录用户
		user, _ := engine.User(c)
		userName := user.Name

		var apiDataSave biz.ApiDataSaveModel
		apiDataSave.App = c.PostForm("app")[1 : len(c.PostForm("app"))-1]
		apiDataSave.Module = c.PostForm("module")[1 : len(c.PostForm("module"))-1]
		apiDataSave.ApiDesc = c.PostForm("apiDesc")[1 : len(c.PostForm("apiDesc"))-1]
		apiDataSave.DataDesc = c.PostForm("dataDesc")[1 : len(c.PostForm("dataDesc"))-1]
		apiDataSave.Path = c.PostForm("path")[1 : len(c.PostForm("path"))-1]
		apiDataSave.Prefix = c.PostForm("prefix")[1 : len(c.PostForm("prefix"))-1]
		apiDataSave.Method = c.PostForm("method")[1 : len(c.PostForm("method"))-1]
		apiDataSave.Product = c.PostForm("product")[1 : len(c.PostForm("product"))-1]
		apiDataSave.Host = c.PostForm("host")[1 : len(c.PostForm("host"))-1]
		apiDataSave.Prototype = c.PostForm("prototype")[1 : len(c.PostForm("prototype"))-1]
		apiDataSave.BodyMode = c.PostForm("bodyMode")[1 : len(c.PostForm("bodyMode"))-1]

		json.Unmarshal([]byte(c.PostForm("pathVars")), &apiDataSave.PathVars)
		json.Unmarshal([]byte(c.PostForm("queryVars")), &apiDataSave.QueryVars)
		json.Unmarshal([]byte(c.PostForm("bodyVars")), &apiDataSave.BodyVars)
		json.Unmarshal([]byte(c.PostForm("headerVars")), &apiDataSave.HeaderVars)
		json.Unmarshal([]byte(c.PostForm("respVars")), &apiDataSave.RespVars)
		json.Unmarshal([]byte(c.PostForm("actions")), &apiDataSave.Actions)
		json.Unmarshal([]byte(c.PostForm("asserts")), &apiDataSave.Asserts)
		json.Unmarshal([]byte(c.PostForm("preApis")), &apiDataSave.PreApis)
		json.Unmarshal([]byte(c.PostForm("postApis")), &apiDataSave.PostApis)

		err := biz.SaveApiData(apiDataSave, userName)
		data := make(map[string]interface{})
		if err != nil {
			data["code"] = 400
			data["msg"] = err
		} else {
			data["code"] = 200
			data["msg"] = "操作成功"
		}

		c.JSON(http.StatusOK, data)
	})

	r.POST("/sceneRun", func(c *gin.Context) {
		var sceneSave biz.SceneSaveModel
		sceneSave.Product = c.PostForm("product")[1 : len(c.PostForm("product"))-1]
		data := make(map[string]interface{})
		if len(sceneSave.Product) == 0 {
			data["code"] = 400
			data["msg"] = "未配置环境信息，请先设置"
			c.JSON(http.StatusOK, data)
			return
		}
		sceneSave.Name = c.PostForm("name")[1 : len(c.PostForm("name"))-1]
		typeTag := c.PostForm("type")[1 : len(c.PostForm("type"))-1]
		var runNumTag string
		if strings.Contains(c.PostForm("runNum"), "\"") {
			runNumTag = c.PostForm("runNum")[1 : len(c.PostForm("runNum"))-1]
		} else {
			runNumTag = c.PostForm("runNum")
		}
		runNum, err := strconv.Atoi(runNumTag)
		if err != nil {
			biz.Logger.Error("%s", err)
			sceneSave.RunNum = 1
		} else {
			sceneSave.RunNum = runNum
		}

		switch typeTag {
		case "串行中断":
			sceneSave.SceneType = 1
		case "串行比较":
			sceneSave.SceneType = 2
		case "串行继续":
			sceneSave.SceneType = 3
		case "普通并发":
			sceneSave.SceneType = 4
		case "并发比较":
			sceneSave.SceneType = 5
		default:
			sceneSave.SceneType = 1
		}

		json.Unmarshal([]byte(c.PostForm("dataList")), &sceneSave.DataList)

		reqDataResps, err := biz.RunPlaybookFromConsole(sceneSave)

		if err != nil {
			data["code"] = 400
			data["msg"] = "操作失败"
		} else {
			data["code"] = 200
			data["msg"] = "操作成功"
		}

		data["data"] = reqDataResps
		c.JSON(http.StatusOK, data)
	})

	r.POST("/dataRun", func(c *gin.Context) {
		var apiDataSave biz.ApiDataSaveModel
		apiDataSave.App = c.PostForm("app")[1 : len(c.PostForm("app"))-1]
		apiDataSave.Module = c.PostForm("module")[1 : len(c.PostForm("module"))-1]
		apiDataSave.ApiDesc = c.PostForm("apiDesc")[1 : len(c.PostForm("apiDesc"))-1]
		apiDataSave.DataDesc = c.PostForm("dataDesc")[1 : len(c.PostForm("dataDesc"))-1]
		apiDataSave.Path = c.PostForm("path")[1 : len(c.PostForm("path"))-1]
		apiDataSave.Prefix = c.PostForm("prefix")[1 : len(c.PostForm("prefix"))-1]
		apiDataSave.Method = c.PostForm("method")[1 : len(c.PostForm("method"))-1]
		apiDataSave.Product = c.PostForm("product")[1 : len(c.PostForm("product"))-1]
		apiDataSave.Host = c.PostForm("host")[1 : len(c.PostForm("host"))-1]
		apiDataSave.Prototype = c.PostForm("prototype")[1 : len(c.PostForm("prototype"))-1]
		apiDataSave.BodyMode = c.PostForm("bodyMode")[1 : len(c.PostForm("bodyMode"))-1]

		json.Unmarshal([]byte(c.PostForm("pathVars")), &apiDataSave.PathVars)
		json.Unmarshal([]byte(c.PostForm("queryVars")), &apiDataSave.QueryVars)
		json.Unmarshal([]byte(c.PostForm("bodyVars")), &apiDataSave.BodyVars)
		json.Unmarshal([]byte(c.PostForm("headerVars")), &apiDataSave.HeaderVars)
		json.Unmarshal([]byte(c.PostForm("respVars")), &apiDataSave.RespVars)
		json.Unmarshal([]byte(c.PostForm("actions")), &apiDataSave.Actions)
		json.Unmarshal([]byte(c.PostForm("asserts")), &apiDataSave.Asserts)
		json.Unmarshal([]byte(c.PostForm("preApis")), &apiDataSave.PreApis)
		json.Unmarshal([]byte(c.PostForm("postApis")), &apiDataSave.PostApis)

		reqDataResps, err := biz.RunApiDebugData(apiDataSave)
		data := make(map[string]interface{})
		if err != nil {
			data["code"] = 400
			data["msg"] = "操作失败"
		} else {
			data["code"] = 200
			data["msg"] = "操作成功"
		}
		data["data"] = reqDataResps
		c.JSON(http.StatusOK, data)
	})

	r.POST("/historyRun", func(c *gin.Context) {
		var apiDataSave biz.HistorySaveModel
		apiDataSave.App = c.PostForm("app")[1 : len(c.PostForm("app"))-1]
		apiDataSave.Module = c.PostForm("module")[1 : len(c.PostForm("module"))-1]
		apiDataSave.ApiDesc = c.PostForm("apiDesc")[1 : len(c.PostForm("apiDesc"))-1]
		apiDataSave.DataDesc = c.PostForm("dataDesc")[1 : len(c.PostForm("dataDesc"))-1]
		apiDataSave.Prototype = c.PostForm("prototype")[1 : len(c.PostForm("prototype"))-1]
		apiDataSave.Host = c.PostForm("host")[1 : len(c.PostForm("host"))-1]
		apiDataSave.Path = c.PostForm("path")[1 : len(c.PostForm("path"))-1]
		apiDataSave.Prefix = c.PostForm("prefix")[1 : len(c.PostForm("prefix"))-1]
		apiDataSave.Method = c.PostForm("method")[1 : len(c.PostForm("method"))-1]
		apiDataSave.Product = c.PostForm("product")[1 : len(c.PostForm("product"))-1]
		apiDataSave.FileName = c.PostForm("fileName")[1 : len(c.PostForm("fileName"))-1]
		apiDataSave.BodyMode = c.PostForm("bodyMode")[1 : len(c.PostForm("bodyMode"))-1]

		json.Unmarshal([]byte(c.PostForm("pathVars")), &apiDataSave.PathVars)
		json.Unmarshal([]byte(c.PostForm("queryVars")), &apiDataSave.QueryVars)
		json.Unmarshal([]byte(c.PostForm("bodyVars")), &apiDataSave.BodyVars)
		json.Unmarshal([]byte(c.PostForm("headerVars")), &apiDataSave.HeaderVars)
		json.Unmarshal([]byte(c.PostForm("respVars")), &apiDataSave.RespVars)
		json.Unmarshal([]byte(c.PostForm("actions")), &apiDataSave.Actions)
		json.Unmarshal([]byte(c.PostForm("asserts")), &apiDataSave.Asserts)
		json.Unmarshal([]byte(c.PostForm("preApis")), &apiDataSave.PreApis)
		json.Unmarshal([]byte(c.PostForm("postApis")), &apiDataSave.PostApis)

		reqDataResps, err := biz.RunHistoryData(apiDataSave)
		data := make(map[string]interface{})
		if err != nil {
			data["code"] = 400
			data["msg"] = "操作失败"
		} else {
			data["code"] = 200
			data["msg"] = "操作成功"
		}
		data["data"] = reqDataResps
		c.JSON(http.StatusOK, data)
	})

	r.POST("/apiCheck", func(c *gin.Context) {
		var apiCheck biz.ApiCheck
		apiCheck.App = c.PostForm("app")
		apiCheck.NodeId = c.PostForm("nodeId")
		apiCheck.DevEnvHost = c.PostForm("devEnvHost")
		apiCheck.Branch = c.PostForm("branch")

		status, err := biz.GetApiChkResult(apiCheck)
		if err != nil {
			biz.Logger.Error("%s", err)
		}

		c.JSON(http.StatusOK, status)

	})

	r.GET("/apiCheckResult", func(c *gin.Context) {
		appName := c.Query("app")
		branch := c.Query("branch")
		data := make(map[string]interface{})
		if len(appName) == 0 {
			data["code"] = 400
			data["message"] = "操作失败"
			data["data"] = fmt.Sprintf("应用名称为空")
		} else {
			resp := biz.GetApiCount4Cicd(appName, branch)
			data["code"] = 200
			data["message"] = "操作成功"
			data["data"] = resp
		}

		c.JSON(http.StatusOK, data)
	})

	// Mock数据
	r.GET("/mock/file/:name", func(c *gin.Context) {
		fileName := c.Param("name")
		lang := c.Query("lang")
		data := biz.GetFileContent(lang, fileName)
		c.String(http.StatusOK, data)
	})

	r.GET("/mock/file", func(c *gin.Context) {
		fileName := c.Query("name")
		lang := c.Query("lang")
		data := biz.GetFileContent(lang, fileName)
		c.String(http.StatusOK, data)
	})

	r.GET("/mock/systemParameter/:name", func(c *gin.Context) {
		name := c.Param("name")
		lang := c.Query("lang")
		data := biz.GetValueFromSysParameter(lang, name)
		c.String(http.StatusOK, data)
	})

	r.GET("/mock/data/slow", func(c *gin.Context) {
		sleep := c.Query("sleep")
		sleepSecond, err := strconv.Atoi(sleep)
		if err != nil {
			fmt.Printf("Error: %s", err)
		} else {
			time.Sleep(time.Duration(sleepSecond) * time.Second)
		}
		sex := biz.GetValueFromSysParameter("", "Sex")
		age := strconv.Itoa(biz.GetRandomInt(2, 130))
		weight := strconv.Itoa(biz.GetRandomInt(80, 250))
		grade := strconv.Itoa(biz.GetRandomInt(1, 12))
		yesOrNo := biz.GetValueFromSysParameter("", "YorN")
		boolValue := biz.GetValueFromSysParameter("", "Bool")
		c.JSON(http.StatusOK, gin.H{
			"name":     chinaid.Name(),
			"age":      age,
			"sex":      sex,
			"id_no":    chinaid.IDNo(),
			"addr":     chinaid.Address(),
			"weight":   weight,
			"mobile":   chinaid.Mobile(),
			"grade":    grade,
			"bool":     boolValue,
			"yesOrNo":  yesOrNo,
			"bankCard": chinaid.BankNo(),
		})
	})

	r.GET("/mock/data/certid/:idno", func(c *gin.Context) {
		idno := c.Param("idno")
		data := biz.GetInfoFromIDNo(idno)
		c.JSON(http.StatusOK, gin.H{
			"city":     data.City,
			"sex":      data.Sex,
			"code":     data.Code,
			"birthday": data.Birthday,
			"district": data.District,
			"province": data.Province,
			"address":  data.Addr,
			"country":  "中国",
		})
	})

	r.GET("/mock/data/certid", func(c *gin.Context) {
		idno := c.Query("idno")
		data := biz.GetInfoFromIDNo(idno)
		c.JSON(http.StatusOK, gin.H{
			"city":     data.City,
			"sex":      data.Sex,
			"code":     data.Code,
			"birthday": data.Birthday,
			"district": data.District,
			"province": data.Province,
			"address":  data.Addr,
			"country":  "中国",
		})
	})

	r.GET("/mock/data/quick", func(c *gin.Context) {
		sex := biz.GetValueFromSysParameter("", "Sex")
		age := strconv.Itoa(biz.GetRandomInt(2, 130))
		weight := strconv.Itoa(biz.GetRandomInt(80, 250))
		grade := strconv.Itoa(biz.GetRandomInt(1, 12))
		yesOrNo := biz.GetValueFromSysParameter("", "YorN")
		boolValue := biz.GetValueFromSysParameter("", "Bool")
		c.JSON(http.StatusOK, gin.H{
			"name":     chinaid.Name(),
			"age":      age,
			"sex":      sex,
			"id_no":    chinaid.IDNo(),
			"addr":     chinaid.Address(),
			"weight":   weight,
			"mobile":   chinaid.Mobile(),
			"grade":    grade,
			"bool":     boolValue,
			"yesOrNo":  yesOrNo,
			"bankCard": chinaid.BankNo(),
		})
	})

	// 导入接口
	r.POST("/api_define_import", func(c *gin.Context) {
		ids := c.PostForm("ids")
		swagger_path := c.PostForm("swagger_path")
		uploadFile, errTmp := c.FormFile("upload_file")
		var uploadFilePath string
		if errTmp == nil {
			uploadFilePath = fmt.Sprintf("%s/%s", biz.UploadBasePath, uploadFile.Filename)
			c.SaveUploadedFile(uploadFile, uploadFilePath)
		}
		var status string
		var idList []string
		if len(ids) == 0 || ids == "," {
			status = "请先选择一个归属应用"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		} else {
			idList = strings.Split(ids, ",")
		}

		failCount, err := biz.GetSwaggerNew(idList[0], swagger_path, uploadFilePath)
		if failCount > 0 {
			status = fmt.Sprintf("导入完成，%d个接口不符合规范，请前往[接口定义]列表查看", failCount)
		} else {
			status = "导入完成，请前往[接口-接口定义]列表查看"
		}

		_ = biz.UpdateApiChangeByAppId(idList[0])
		if err != nil {
			status = fmt.Sprintf("%s; %s", status, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			status = "导入完成，请前往[接口-接口定义]列表查看"
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	// Xmind2Excel
	r.POST("/testcase_xmind2excel", func(c *gin.Context) {
		uploadFile, errTmp := c.FormFile("upload_file")
		var uploadFilePath string
		if errTmp == nil {
			uploadFilePath = fmt.Sprintf("%s/%s", biz.UploadBasePath, uploadFile.Filename)
			c.SaveUploadedFile(uploadFile, uploadFilePath)
		}
		var status string

		if len(uploadFilePath) == 0 {
			status = "请先上传Xmind用例文档"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		fileName, err := biz.Xmind2Excel(uploadFilePath)
		if err != nil {
			status = fmt.Sprintf("转换失败: %s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			//status = fmt.Sprintf("转换完成，请前往[文件-用例文件]下载, 文件名: %s", fileName)
			hostIp := c.Request.Host
			downloadUrl := fmt.Sprintf("http://%s/admin/fm/case/download?path=/%s", hostIp, fileName)
			status = fmt.Sprintf("转换完成\n请复制下述链接下载:\n%s", downloadUrl)
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	// Xmind2Import
	r.POST("/testcase_xmind2import", func(c *gin.Context) {
		introVersion := c.PostForm("intro_version")
		product := c.PostForm("product")
		uploadFile, errTmp := c.FormFile("upload_file")
		var uploadFilePath string
		if errTmp == nil {
			uploadFilePath = fmt.Sprintf("%s/%s", biz.UploadBasePath, uploadFile.Filename)
			c.SaveUploadedFile(uploadFile, uploadFilePath)
		}
		var status string

		if len(uploadFilePath) == 0 {
			status = "请先上传Xmind用例文档"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		err := biz.Xmind2Import(product, introVersion, uploadFilePath)
		if err != nil {
			status = fmt.Sprintf("转换失败: %s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			status = fmt.Sprintf("导入完成, 请刷新列表查看~ ")
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	// Xmind2ImportAndUse
	r.POST("/testcase_xmind_import_and_use", func(c *gin.Context) {
		user, _ := engine.User(c)
		userName := user.Name
		introVersion := c.PostForm("intro_version")
		product := c.PostForm("product")
		uploadFile, errTmp := c.FormFile("upload_file")
		var uploadFilePath string
		if errTmp == nil {
			uploadFilePath = fmt.Sprintf("%s/%s", biz.UploadBasePath, uploadFile.Filename)
			c.SaveUploadedFile(uploadFile, uploadFilePath)
		}
		var status string

		if len(uploadFilePath) == 0 {
			status = "请先上传Xmind用例文档"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		err := biz.Xmind2ImportAndUse(userName, product, introVersion, uploadFilePath)
		if err != nil {
			status = fmt.Sprintf("转换失败: %s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			status = fmt.Sprintf("导入完成, 请刷新列表查看~ ")
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	// AI用例
	r.POST("/ai_case_create_by_create_desc", func(c *gin.Context) {
		user, _ := engine.User(c)
		var inputCase biz.InputCase
		inputCase.CreateUser = user.Name
		inputCase.AiTemplate = c.PostForm("ai_template")
		inputCase.IntroVersion = c.PostForm("intro_version")
		inputCase.Product = c.PostForm("product")
		inputCase.CreatePlatform = c.PostForm("create_platform")
		createDesc := c.PostForm("create_desc")
		uploadFile, errTmp := c.FormFile("upload_file")
		var uploadFilePath string
		if errTmp == nil {
			uploadFilePath = fmt.Sprintf("%s/%s", biz.UploadBasePath, uploadFile.Filename)
			c.SaveUploadedFile(uploadFile, uploadFilePath)
			biz.Logger.Debug("uploadFilePath: %s", uploadFilePath)
		}
		status := "生成任务已在后台运行,请稍后刷新列表查看生成用例"
		if len(createDesc) == 0 {
			status = "请先输入生成需求"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		err := inputCase.AICreateCaseByCreateDesc(createDesc, uploadFilePath)
		if err != nil {
			status = fmt.Sprintf("生成失败：%s: %s", createDesc, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	r.POST("/ai_case_optimize", func(c *gin.Context) {
		user, _ := engine.User(c)
		createUser := user.Name
		ids := c.PostForm("ids")
		createPlatform := c.PostForm("optimize_platform")
		optimizeDesc := c.PostForm("optimize_desc")
		status := "优化任务已在后台运行,请稍后刷新列表查看优化用例"

		if len(optimizeDesc) == 0 || ids == "," {
			status = "请先选择优化用例和输入优化指令"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		ids = strings.Trim(ids, ",")
		err := biz.AIOptimizeCase(ids, optimizeDesc, createPlatform, createUser)
		if err != nil {
			status = fmt.Sprintf("优化失败：%s: %s", optimizeDesc, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	r.POST("/ai_create_case_by_api_define", func(c *gin.Context) {
		user, _ := engine.User(c)
		var inputCase biz.InputCase
		inputCase.CreateUser = user.Name
		ids := c.PostForm("ids")
		inputCase.AiTemplate = c.PostForm("ai_template")
		inputCase.IntroVersion = c.PostForm("intro_version")
		inputCase.Product = c.PostForm("product")
		inputCase.CreatePlatform = c.PostForm("create_platform")
		status := "生成任务已在后台运行,请稍后前往[助手-智能用例]列表查看生成用例"
		if ids == "," {
			status = "请先选择接口定义"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}
		ids = strings.Trim(ids, ",")
		err = inputCase.AICreateCaseByApiDefine(ids)
		if err != nil {
			status = fmt.Sprintf("生成失败：%s: %s", ids, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return

	})

	r.POST("/ai_case_import", func(c *gin.Context) {
		user, _ := engine.User(c)
		var input biz.ImportCommon
		input.CreateUser = user.Name
		input.IntroVersion = c.PostForm("intro_version")
		input.Product = c.PostForm("product")
		input.CreatePlatform = c.PostForm("create_platform")
		input.ConversationId = c.PostForm("conversation_id")
		input.RawReply = c.PostForm("raw_reply")
		status := "导入完成,请刷新列表查看生成用例"
		if len(input.ConversationId) == 0 && len(input.RawReply) == 0 {
			status = "请先输入会话ID或原生回复"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		err := input.AICreateCaseByImport()
		if err != nil {
			status = fmt.Sprintf("导入失败：%s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	// AI数据
	r.POST("/ai_data_create_by_create_desc", func(c *gin.Context) {
		user, _ := engine.User(c)
		var input biz.InputData
		input.CreateUser = user.Name
		input.AiTemplate = c.PostForm("ai_template")
		input.IntroVersion = c.PostForm("intro_version")
		input.Product = c.PostForm("product")
		input.CreatePlatform = c.PostForm("create_platform")
		createDesc := c.PostForm("create_desc")
		uploadFile, errTmp := c.FormFile("upload_file")
		var uploadFilePath string
		if errTmp == nil {
			uploadFilePath = fmt.Sprintf("%s/%s", biz.UploadBasePath, uploadFile.Filename)
			c.SaveUploadedFile(uploadFile, uploadFilePath)
		}
		status := "生成任务已在后台运行,请稍后刷新列表查看生成数据"
		if len(createDesc) == 0 {
			status = "请先输入生成需求"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}
		err := input.AICreateDataAndPlaybookByCreateDesc(createDesc, uploadFilePath)
		if err != nil {
			status = fmt.Sprintf("生成失败：%s: %s", createDesc, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	r.POST("/ai_data_optimize", func(c *gin.Context) {
		user, _ := engine.User(c)
		createUser := user.Name
		ids := c.PostForm("ids")
		createPlatform := c.PostForm("optimize_platform")
		optimizeDesc := c.PostForm("optimize_desc")
		status := "优化任务已在后台运行,请稍后刷新列表查看优化数据"

		if len(optimizeDesc) == 0 || ids == "," {
			status = "请先选择优化数据和输入优化指令"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		ids = strings.Trim(ids, ",")
		err := biz.AIOptimizeData(ids, optimizeDesc, createPlatform, createUser)
		if err != nil {
			status = fmt.Sprintf("优化失败：%s: %s", optimizeDesc, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	r.POST("/ai_create_data_by_api_define", func(c *gin.Context) {
		user, _ := engine.User(c)
		var input biz.InputData
		input.CreateUser = user.Name
		ids := c.PostForm("ids")
		input.AiTemplate = c.PostForm("ai_template")
		input.IntroVersion = c.PostForm("intro_version")
		input.CreatePlatform = c.PostForm("create_platform")
		input.Product = c.PostForm("product")
		status := "生成任务已在后台运行,请稍后前往[助手-智能数据]列表查看生成数据"
		if ids == "," {
			status = "请先选择接口定义"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}
		ids = strings.Trim(ids, ",")
		err = input.AICreateDataAndPlaybookByApiDefine(ids)
		if err != nil {
			status = fmt.Sprintf("生成失败：%s: %s", ids, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return

	})

	r.POST("/ai_data_import", func(c *gin.Context) {
		user, _ := engine.User(c)
		var importCommon biz.ImportCommon
		importCommon.CreateUser = user.Name
		importCommon.CreatePlatform = c.PostForm("create_platform")
		importCommon.IntroVersion = c.PostForm("intro_version")
		importCommon.ConversationId = c.PostForm("conversation_id")
		importCommon.RawReply = c.PostForm("raw_reply")
		importCommon.Product = c.PostForm("product")
		status := "导入完成,请刷新列表查看生成数据"
		if len(importCommon.ConversationId) == 0 && len(importCommon.RawReply) == 0 {
			status = "请先输入会话ID或原生回复"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		err := importCommon.AICreateDataAndPlaybookByImport()
		if err != nil {
			status = fmt.Sprintf("导入失败：%s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	r.POST("/ai_data_test_and_analysis", func(c *gin.Context) {
		var analysisInput biz.AnalysisDataInput
		user, _ := engine.User(c)
		analysisInput.CreateUser = user.Name
		ids := c.PostForm("ids")
		analysisInput.CreatePlatform = c.PostForm("analysis_platform")
		analysisInput.AiTemplate = c.PostForm("ai_template")
		analysisInput.Product = c.PostForm("product")
		status := "分析任务已在后台运行,请稍后查看执行结果，前往[智能分析]列表查看分析结果"
		err = biz.AiDataTest(ids, analysisInput)
		if err != nil {
			status = fmt.Sprintf("分析遇错：%s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	r.POST("/ai_data_test", func(c *gin.Context) {
		//user, _ := engine.User(c)
		//userName := user.Name
		idStr := c.PostForm("ids")
		product := c.PostForm("product")

		var status string
		if idStr == "," {
			status = "请先选择数据再测试"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			ids := strings.Split(idStr, ",")
			var errTag int
			var status string
			source := "ai_data"
			status = "测试完成，请刷新列表查看测试结果"
			for _, id := range ids {
				if len(id) == 0 {
					//status = "测试完成，请刷新列表查看测试结果"
					continue
				}

				err := biz.RepeatRunDataFile(id, product, source)
				if err != nil {
					if errTag == 0 {
						status = fmt.Sprintf("测试失败：%s: %s", id, err)
					} else {
						status = fmt.Sprintf("%s; %s: %s", status, id, err)
					}
					errTag = 1
				}
			}

			if errTag == 0 {
				c.JSON(http.StatusOK, map[string]interface{}{
					"code": 200,
					"msg":  status,
					"data": map[string]string{},
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]interface{}{
					"code": 400,
					"msg":  status,
					"data": map[string]string{},
				})
			}
		}

		return
	})

	r.POST("/data_batch_run", func(c *gin.Context) {
		//user, _ := engine.User(c)
		//userName := user.Name
		idStr := c.PostForm("ids")
		product := c.PostForm("product")

		var status string
		if idStr == "," {
			status = "请先选择数据再测试"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]interface{}{},
			})
		} else {
			ids := strings.Split(idStr, ",")
			var errTag int
			var status string
			source := "data"
			status = "测试完成，请刷新列表查看测试结果"
			for _, id := range ids {
				if len(id) == 0 {
					//status = "测试完成，请刷新列表查看测试结果"
					continue
				}
				err := biz.RepeatRunDataFile(id, product, source)
				if err != nil {
					if errTag == 0 {
						status = fmt.Sprintf("测试失败：%s: %s", id, err)
					} else {
						status = fmt.Sprintf("%s; %s: %s", status, id, err)
					}
					errTag = 1
				}
			}

			if errTag == 0 {
				c.JSON(http.StatusOK, map[string]interface{}{
					"code": 200,
					"msg":  status,
					"data": map[string]string{},
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]interface{}{
					"code": 400,
					"msg":  status,
					"data": map[string]string{},
				})
			}
		}

		return
	})

	// AI分析
	r.POST("/ai_issue_import", func(c *gin.Context) {
		user, _ := engine.User(c)
		var importCommon biz.ImportCommon
		importCommon.CreateUser = user.Name
		importCommon.CreatePlatform = c.PostForm("create_platform")
		importCommon.ConversationId = c.PostForm("conversation_id")
		importCommon.RawReply = c.PostForm("raw_reply")
		importCommon.Product = c.PostForm("product")
		status := "导入完成,请刷新列表查看分析数据"
		if len(importCommon.ConversationId) == 0 && len(importCommon.RawReply) == 0 {
			status = "请先输入会话ID或原生回复"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		err := importCommon.AIAnalysisDataByImport()
		if err != nil {
			status = fmt.Sprintf("导入失败：%s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	// AI场景
	r.POST("/ai_playbook_create_by_create_desc", func(c *gin.Context) {
		user, _ := engine.User(c)
		var input biz.InputPlaybook
		input.CreateUser = user.Name
		input.AiTemplate = c.PostForm("ai_template")
		input.IntroVersion = c.PostForm("intro_version")
		input.Product = c.PostForm("product")
		input.CreatePlatform = c.PostForm("create_platform")
		createDesc := c.PostForm("create_desc")
		uploadFile, errTmp := c.FormFile("upload_file")
		var uploadFilePath string
		if errTmp == nil {
			uploadFilePath = fmt.Sprintf("%s/%s", biz.UploadBasePath, uploadFile.Filename)
			c.SaveUploadedFile(uploadFile, uploadFilePath)
		}
		status := "生成任务已在后台运行,请稍后刷新列表查看生成场景"
		if len(createDesc) == 0 {
			status = "请先输入生成需求"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}
		err := input.AICreateDataAndPlaybookByCreateDesc(createDesc, uploadFilePath)
		if err != nil {
			status = fmt.Sprintf("生成失败：%s: %s", createDesc, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	r.POST("/ai_playbook_optimize", func(c *gin.Context) {
		user, _ := engine.User(c)
		createUser := user.Name
		ids := c.PostForm("ids")
		createPlatform := c.PostForm("optimize_platform")
		optimizeDesc := c.PostForm("optimize_desc")
		status := "优化任务已在后台运行,请稍后刷新列表查看优化场景"

		if len(optimizeDesc) == 0 || ids == "," {
			status = "请先选择优化场景和输入优化指令"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
			return
		}

		ids = strings.Trim(ids, ",")
		err := biz.AIOptimizeCase(ids, optimizeDesc, createPlatform, createUser)
		if err != nil {
			status = fmt.Sprintf("优化失败：%s: %s", optimizeDesc, err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
		return
	})

	r.POST("/ai_playbook_import", func(c *gin.Context) {
		user, _ := engine.User(c)
		var input biz.ImportCommon
		input.CreateUser = user.Name
		input.CreatePlatform = c.PostForm("create_platform")
		input.ConversationId = c.PostForm("conversation_id")
		input.RawReply = c.PostForm("raw_reply")
		input.Product = c.PostForm("product")
		input.IntroVersion = c.PostForm("intro_version")
		status := "导入完成,请刷新列表查看生成场景"
		if len(input.ConversationId) == 0 && len(input.RawReply) == 0 {
			status = "请先输入会话ID或原生回复"
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		}

		err := input.AICreateDataAndPlaybookByImport()
		if err != nil {
			status = fmt.Sprintf("导入失败：%s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
	})

	r.POST("/ai_playbook_test_and_analysis", func(c *gin.Context) {
		var analysisInput biz.AnalysisDataInput
		user, _ := engine.User(c)
		analysisInput.CreateUser = user.Name
		ids := c.PostForm("ids")
		analysisInput.CreatePlatform = c.PostForm("analysis_platform")
		analysisInput.AiTemplate = c.PostForm("ai_template")
		analysisInput.Product = c.PostForm("product")
		status := "分析任务已在后台运行,请稍后查看执行结果，前往[智能分析]列表查看分析结果"
		err = biz.AiPlaybookTest(ids, "ai_playbook", analysisInput)
		if err != nil {
			status = fmt.Sprintf("分析遇错：%s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{
					"reason": "操作失败",
				},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{
					"reason": "操作成功",
				},
			})
		}
	})

	// 同步知识库
	r.POST("/sync_knowledge", func(c *gin.Context) {
		user, _ := engine.User(c)
		syncUser := user.Name
		kType := c.PostForm("k_type")

		status := "同步任务已在后台运行,请稍后查看同步结果，前往[知识库]查看结果"

		err = biz.UpdateAssetKnowledge(kType, syncUser)
		if err != nil {
			status = fmt.Sprintf("同步遇错：%s", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": 400,
				"msg":  status,
				"data": map[string]string{},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": 200,
				"msg":  status,
				"data": map[string]string{},
			})
		}
	})

	models.Init(eng.MysqlConnection())
	biz.Init()
	listen := fmt.Sprintf(":%d", biz.SERVER_PORT)

	err1 := r.Run(listen)
	if err1 != nil {
		panic(err1)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.MysqlConnection().Close()
	eng.SqliteConnection().Close()
}
