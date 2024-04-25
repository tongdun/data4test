package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	//"sync"
)

func HistoryDataRunAgain(id string) (err error) {
	var hData HistoryDataDetail
	models.Orm.Table("scene_data_test_history").Where("id = ?", id).Find(&hData)

	if len(hData.ApiId) == 0 {
		err = fmt.Errorf("未找到历史数据的详情信息")
		Logger.Error("Error: %s, ID: %s", err, id)
		return
	}

	basePath := GetStrFromHtml(hData.Content)
	dirName := GetHistoryDataDirName(basePath)

	filePath := fmt.Sprintf("%s/%s/%s", HistoryBasePath, dirName, basePath)

	var dataFile DataFile
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if strings.HasSuffix(filePath, ".json") {
		err = json.Unmarshal([]byte(content), &dataFile)
	} else {
		err = yaml.Unmarshal([]byte(content), &dataFile)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	_, _, _, _, _, result, dst, err := dataFile.RunDataFileStruct(hData.App, hData.Product, filePath, "again", "mgmt", nil)

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
