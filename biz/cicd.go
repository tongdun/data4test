package biz

import (
	"data4perf/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func CallBack4Cicd(url, nodeId, status string) {
	time.Sleep(time.Duration(5) * time.Second)
	fullUrl := fmt.Sprintf("%s?nodeId=%s&status=%s", url, nodeId, status)
	Logger.Info("开始回调: %s", fullUrl)
	resp, err := http.Get(fullUrl)
	if err != nil {
		Logger.Error("%s", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		Logger.Info("回调结束：%s", string(body))
		resp.Body.Close()
	}
}

func GetApiCount4Cicd(appName, branch string) (resp string) {
	var appApiChange AppApiChange
	var apiChangeStr string

	models.Orm.Table("app_api_changelog").Order("created_at desc").Where("app = ? and branch = ?", appName, branch).Limit(1).Find(&appApiChange)

	if appApiChange.NewApiSum > 0 || appApiChange.DeletedApiSum > 0 || appApiChange.ChangedApiSum > 0 {
		apiChangeStr = fmt.Sprintf("接口变动: 总数: %v, 保持原样: %v, 新增: %v, 被删除: %v, 被修改: %v", appApiChange.CurApiSum, appApiChange.ExistApiSum, appApiChange.NewApiSum, appApiChange.DeletedApiSum, appApiChange.ChangedApiSum)
	}

	if appApiChange.CurApiSum > 0 || len(appApiChange.ApiCheckResult) > 0 {
		if len(apiChangeStr) > 0 {
			resp = fmt.Sprintf("%s\n%s", apiChangeStr, appApiChange.ApiCheckResult)
		} else {
			resp = appApiChange.ApiCheckResult
		}
		id, _ := GetAppId(appName)
		apiCheckResultUrl := fmt.Sprintf("查看详情: http://%s:%s/admin/app_dashboard?id=%d", HOST_IP, SERVER_PORT, id)
		resp = fmt.Sprintf("%s%s\n", resp, apiCheckResultUrl)
	} else {
		resp = fmt.Sprintf("[%s] 应用暂无接口数据，请配置Swagger接口文档或在 http://%s:%s 控制台进行手动维护", appName, HOST_IP, SERVER_PORT)
	}

	return
}
