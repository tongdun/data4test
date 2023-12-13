package main

import (
	"encoding/json"
	"github.com/mritd/chinaid"
	"strconv"
	"time"

	//"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"

	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"

	"data4perf/login"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/gin-contrib/pprof"
	//_ "github.com/mkevac/debugcharts"  
	//_ "net/http/pprof"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin" // web framework adapter
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

	"data4perf/biz"
	"data4perf/models"
	"data4perf/pages"
	"data4perf/tables"

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
	pprof.Register(r)

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
		ctx.Redirect(http.StatusMovedPermanently, "/admin/likePostman")
		return
	})

	r.GET("/admin", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/admin/likePostman")
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
			data["data"] = appList
		} else {
			dataInfo := biz.GetDataInfoByDataDesc("", "", "", dataDesc)
			data["data"] = dataInfo
		}
		data["code"] = 200
		data["msg"] = "操作成功"
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
				data["code"] = 400
				data["msg"] = err
				c.JSON(http.StatusOK, data)
			}
		}

		err := biz.SaveApiData(apiDataSave)

		if err != nil {
			data["code"] = 400
			data["msg"] = err
		} else {
			data["code"] = 200
			data["msg"] = "操作成功"
		}

		c.JSON(http.StatusOK, data)
	})

	r.POST("/sceneSave", func(c *gin.Context) {
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
		if typeTag == "比较" {
			sceneSave.SceneType = 2
		} else {
			sceneSave.SceneType = 1
		}

		json.Unmarshal([]byte(c.PostForm("dataList")), &sceneSave.DataList)
		data := make(map[string]interface{})

		if len(sceneSave.Name) == 0 {
			data["code"] = 400
			data["msg"] = "场景名称不能为空"
		} else {
			err := biz.SaveScene(sceneSave)
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

		err := biz.SaveApiData(apiDataSave)
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
		sceneSave.Name = c.PostForm("name")[1 : len(c.PostForm("name"))-1]
		typeTag := c.PostForm("type")[1 : len(c.PostForm("type"))-1]
		runNumTag := c.PostForm("runNum")
		runNum, err := strconv.Atoi(runNumTag)
		if err != nil {
			sceneSave.RunNum = 1
		} else {
			sceneSave.RunNum = runNum
		}
		if typeTag == "比较" {
			sceneSave.SceneType = 2
		} else {
			sceneSave.SceneType = 1
		}

		json.Unmarshal([]byte(c.PostForm("dataList")), &sceneSave.DataList)

		reqDataResps, err := biz.RunPlaybookDebugData(sceneSave)
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

		//json.Unmarshal([]byte(c.PostForm("product")), &apiDataSave.Envs)

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
		data := biz.GetFileContent(fileName)
		c.String(http.StatusOK, data)
	})

	r.GET("/mock/systemParameter/:name", func(c *gin.Context) {
		name := c.Param("name")
		data := biz.GetValueFromSysParameter("", name)
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
		idno := c.Param("name")
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
