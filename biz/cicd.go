package biz

import (
	"data4test/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func CallBack4Cicd(url, nodeId, status string) {
	time.Sleep(time.Duration(5) * time.Second)
	fullUrl := fmt.Sprintf("%s?nodeId=%s&status=%s", url, nodeId, status)
	Logger.Info(T("info.callback_start"), fullUrl)
	resp, err := http.Get(fullUrl)
	if err != nil {
		Logger.Error("%s", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		Logger.Info(T("info.callback_end"), string(body))
		resp.Body.Close()
	}
}

func GetApiCount4Cicd(appName, branch string) (resp string) {
	var appApiChange AppApiChange
	var apiChangeStr string

	models.Orm.Table("app_api_changelog").Order("created_at desc").Where("app = ? and branch = ?", appName, branch).Limit(1).Find(&appApiChange)

	if appApiChange.NewApiSum > 0 || appApiChange.DeletedApiSum > 0 || appApiChange.ChangedApiSum > 0 {
		apiChangeStr = fmt.Sprintf(T("info.api_change_summary"), appApiChange.CurApiSum, appApiChange.ExistApiSum, appApiChange.NewApiSum, appApiChange.DeletedApiSum, appApiChange.ChangedApiSum)
	}

	if appApiChange.CurApiSum > 0 || len(appApiChange.ApiCheckResult) > 0 {
		if len(apiChangeStr) > 0 {
			resp = fmt.Sprintf("%s\n%s", apiChangeStr, appApiChange.ApiCheckResult)
		} else {
			resp = appApiChange.ApiCheckResult
		}
		id, _ := GetAppId(appName)
		apiCheckResultUrl := fmt.Sprintf(T("info.view_details"), HOST_IP, SERVER_PORT, id)
		resp = fmt.Sprintf("%s%s\n", resp, apiCheckResultUrl)
	} else {
		resp = fmt.Sprintf(T("error.no_api_data"), appName, HOST_IP, SERVER_PORT)
	}

	return
}
