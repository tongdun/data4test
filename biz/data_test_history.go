package biz

import (
	"data4perf/models"
	"fmt"
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

	result, dst, err := RunDataFile(hData.App, filePath, hData.Product, "historyAgain", nil)
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
