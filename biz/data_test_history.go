package biz

import (
	"data4test/models"
	"fmt"
	"path"
	"strings"
	//"sync"
)

func HistoryDataRunAgain(id string) (err error) {
	var hData HistoryDataDetail
	models.Orm.Table("scene_data_test_history").Where("id = ?", id).Find(&hData)

	if len(hData.Content) == 0 {
		err = fmt.Errorf("未找到历史数据的详情信息")
		Logger.Error("Error: %s, ID: %s", err, id)
		return
	}

	basePath := GetStrFromHtml(hData.Content)
	dirName := GetHistoryDataDirName(basePath)

	filePath := fmt.Sprintf("%s/%s/%s", HistoryBasePath, dirName, basePath)

	//var dbData DbSceneData
	//dbData.App = hData.App
	//dbData.Name = hData.Name
	//dbData.Content = hData.Content

	result, dst, err := hData.RunDataFile(filePath, hData.Product, "historyAgain", nil)
	//result, dst, err := RunDataFile(hData.App, filePath, hData.Product, "historyAgain", nil)
	var sceneDataRecord SceneDataRecord

	if err != nil {
		sceneDataRecord.FailReason = fmt.Sprintf("%s", err)
	}

	err = WriteDataResultByFile(filePath, result, dst, hData.Product, hData.EnvType, err)

	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

func RecordDataHistory(dst, product, source string, envType int, dbData DbSceneData) (err error) {
	var sceneDataRecord SceneDataRecord

	if len(dst) > 0 {
		b, _ := IsStrEndWithTimeFormat(path.Base(dst))
		if b {
			dirName := GetHistoryDataDirName(path.Base(dst))
			sceneDataRecord.Content = fmt.Sprintf("<a href=\"/admin/fm/history/preview?path=/%s/%s\">%s</a>", dirName, path.Base(dst), path.Base(dst))
		} else {
			if strings.HasPrefix(source, "ai") {
				sceneDataRecord.Content = fmt.Sprintf("<a href=\"/admin/fm/ai_data/preview?path=/%s\">%s</a>", path.Base(dst), path.Base(dst))
			} else {
				sceneDataRecord.Content = fmt.Sprintf("<a href=\"/admin/fm/data/preview?path=/%s\">%s</a>", path.Base(dst), path.Base(dst))

			}
		}
	}

	sceneDataRecord.Name = dbData.Name
	sceneDataRecord.ApiId = dbData.ApiId
	sceneDataRecord.App = dbData.App
	sceneDataRecord.Result = dbData.Result
	sceneDataRecord.FailReason = dbData.FailReason
	sceneDataRecord.EnvType = envType
	sceneDataRecord.Product = product

	err = models.Orm.Table("scene_data_test_history").Create(sceneDataRecord).Error

	if err != nil {
		Logger.Error("%s", err)
	}

	return
}
