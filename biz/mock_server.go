package biz

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mritd/chinaid"

	"data4test/models"
)

// RegisterMockRoutes 注册三方接口 Mock 路由（Web 与 CLI 共用）。
// 数据来源均为本地文件（DownloadBasePath/UploadBasePath）、sys_parameter 表与 chinaid 纯函数，离线可用。
func RegisterMockRoutes(r *gin.Engine) {
	// Mock 数据
	r.GET("/mock/file/:name", func(c *gin.Context) {
		fileName := c.Param("name")
		lang := c.Query("lang")
		data := GetFileContent(lang, fileName)
		c.String(http.StatusOK, data)
	})

	r.GET("/mock/file", func(c *gin.Context) {
		fileName := c.Query("name")
		lang := c.Query("lang")
		data := GetFileContent(lang, fileName)
		c.String(http.StatusOK, data)
	})

	r.GET("/mock/systemParameter/:name", func(c *gin.Context) {
		name := c.Param("name")
		lang := c.Query("lang")
		data := GetValueFromSysParameter(lang, name)
		c.String(http.StatusOK, data)
	})

	r.GET("/mock/data/slow", func(c *gin.Context) {
		sleep := c.Query("sleep")
		sleepSecond, err := strconv.Atoi(sleep)
		if err != nil {
			Logger.Warning("mock slow sleep parse: %s", err)
		} else {
			time.Sleep(time.Duration(sleepSecond) * time.Second)
		}
		c.JSON(http.StatusOK, mockPersonData())
	})

	r.GET("/mock/data/certid/:idno", func(c *gin.Context) {
		c.JSON(http.StatusOK, mockIDNoData(c.Param("idno")))
	})

	r.GET("/mock/data/certid", func(c *gin.Context) {
		c.JSON(http.StatusOK, mockIDNoData(c.Query("idno")))
	})

	r.GET("/mock/data/quick", func(c *gin.Context) {
		c.JSON(http.StatusOK, mockPersonData())
	})
}

// mockPersonData 生成随机个人信息（slow 与 quick 共用）。
func mockPersonData() gin.H {
	return gin.H{
		"name":     chinaid.Name(),
		"age":      strconv.Itoa(GetRandomInt(2, 130)),
		"sex":      GetValueFromSysParameter("", "Sex"),
		"id_no":    chinaid.IDNo(),
		"addr":     chinaid.Address(),
		"weight":   strconv.Itoa(GetRandomInt(80, 250)),
		"mobile":   chinaid.Mobile(),
		"grade":    strconv.Itoa(GetRandomInt(1, 12)),
		"bool":     GetValueFromSysParameter("", "Bool"),
		"yesOrNo":  GetValueFromSysParameter("", "YorN"),
		"bankCard": chinaid.BankNo(),
	}
}

// mockIDNoData 根据身份证号解析返回身份信息。
func mockIDNoData(idno string) gin.H {
	data := GetInfoFromIDNo(idno)
	return gin.H{
		"city":     data.City,
		"sex":      data.Sex,
		"code":     data.Code,
		"birthday": data.Birthday,
		"district": data.District,
		"province": data.Province,
		"address":  data.Addr,
		"country":  "中国",
	}
}

// startCliMock 启动仅注册 mock 路由的本地 Gin 服务，监听 0.0.0.0:port（所有网卡），
// 供被测系统经本机局域网 IP 调用（不能绑 127.0.0.1，否则被测系统无法访问）。返回 server 供优雅关闭。
func startCliMock(port int) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	RegisterMockRoutes(r)
	srv := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Warning("cli mock server error: %s", err)
		}
	}()
	return srv
}

// shutdownCliMock 优雅关闭本地 mock 服务。
func shutdownCliMock(srv *http.Server) {
	if srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

// detectCliMockUsage 检测本地库数据文件 / 系统参数是否引用了 /mock/ 地址。
func detectCliMockUsage() bool {
	var cnt int
	models.Orm.Table("scene_data").Where("content like ?", "%/mock/%").Count(&cnt)
	if cnt > 0 {
		return true
	}
	models.Orm.Table("sys_parameter").Where("value_list like ?", "%/mock/%").Count(&cnt)
	return cnt > 0
}

// replaceCliMockURLs 把数据文件 / 系统参数中指向 Web 端 mock 主机的地址替换为本机 mock 地址，
// 返回替换后的本机 mock 地址（http://<local-ip>:<port>）供对外宣告。baseURL 为空则返回空串。
// 替换地址用本机局域网 IP（供被测系统调用），而非 127.0.0.1。REPLACE 为 MySQL 与 SQLite 通用函数。
func replaceCliMockURLs(baseURL string, port int) string {
	if baseURL == "" {
		return ""
	}
	local := fmt.Sprintf("http://%s:%d", getLocalIP(), port)
	models.Orm.Exec("UPDATE scene_data SET content = REPLACE(content, ?, ?)", baseURL, local)
	models.Orm.Exec("UPDATE sys_parameter SET value_list = REPLACE(value_list, ?, ?)", baseURL, local)
	return local
}

// getLocalIP 返回本机对外可达的 IPv4 地址，供 mock 服务对外宣告使用。
// 自动探测首个非回环 IPv4；探测失败则回退 127.0.0.1。
// 注：mock 地址如需指定，由自动化任务的 private_parameter 私有参数自行维护，不从 config.json 读取。
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ip4 := ipnet.IP.To4(); ip4 != nil {
			return ip4.String()
		}
	}
	return "127.0.0.1"
}
