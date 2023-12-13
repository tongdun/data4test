package biz

import (
	"data4perf/models"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"sync"
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

	var dataFile, dataFileNew DataFile
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

	urls := dataFile.Urls
	header := dataFile.Single.Header
	datas := dataFile.Request
	tag := 0
	var resList [][]byte
	var errs []error

	if dataFile.GetIsParallel() {
		wg := sync.WaitGroup{}
		for _, url := range urls {
			if len(datas) > 0 {
				if dataFile.Api.Method == "get" {
					for _, data := range datas {
						if tag == 0 {
							dataFileNew.Request = []string{data}
						} else {
							dataFileNew.Request = append(dataFile.Request, data)
						}
						dataMap := make(map[string]interface{})
						errTmp := json.Unmarshal([]byte(data), &dataMap)
						if errTmp != nil {
							Logger.Error("%s", errTmp)
							return
						}
						tag++
						wg.Add(1)
						go func(method, url string, data map[string]interface{}, header map[string]interface{}) {
							defer wg.Add(-1)
							res, err := RunHttp(method, url, data, header)
							resList = append(resList, res)
							errs = append(errs, err)
						}(dataFile.Api.Method, url, dataMap, header)
					}
				} else {
					for _, data := range datas {
						if tag == 0 {
							dataFile.Request = []string{data}
						} else {
							dataFile.Request = append(dataFile.Request, data)
						}
						dataMap := make(map[string]interface{})
						errTmp := json.Unmarshal([]byte(data), &dataMap)
						if errTmp != nil {
							Logger.Error("%s", errTmp)
							return
						}
						tag++
						wg.Add(1)
						go func(method, url string, data map[string]interface{}, header map[string]interface{}) {
							defer wg.Add(-1)
							res, err := RunHttp(method, url, data, header)
							resList = append(resList, res)
							errs = append(errs, err)
						}(dataFile.Api.Method, url, dataMap, header)
					}
				}
			} else {
				dataFileNew.Request = []string{}
				wg.Add(1)
				go func(method, url string, header map[string]interface{}) {
					res, err := RunHttp(method, url, nil, header)
					resList = append(resList, res)
					errs = append(errs, err)
				}(dataFile.Api.Method, url, header)
			}
			wg.Wait()
		}
	} else {
		for _, url := range urls {
			if len(datas) > 0 {
				if dataFile.Api.Method == "get" {
					for _, data := range datas {
						if tag == 0 {
							dataFileNew.Request = []string{data}
						} else {
							dataFileNew.Request = append(dataFileNew.Request, data)
						}
						dataMap := make(map[string]interface{})
						errTmp := json.Unmarshal([]byte(data), &dataMap)
						if errTmp != nil {
							Logger.Error("%s", errTmp)
							return
						}
						tag++
						res, err := RunHttp(dataFile.Api.Method, url, dataMap, header)
						resList = append(resList, res)
						errs = append(errs, err)
					}
				} else {
					for _, data := range datas {
						if tag == 0 {
							dataFileNew.Request = []string{data}
						} else {
							dataFileNew.Request = append(dataFileNew.Request, data)
						}
						dataMap := make(map[string]interface{})
						errTmp := json.Unmarshal([]byte(data), &dataMap)
						if errTmp != nil {
							Logger.Error("%s", errTmp)
							return
						}
						tag++
						res, err := RunHttp(dataFile.Api.Method, url, dataMap, header)
						resList = append(resList, res)
						errs = append(errs, err)
					}
				}

			} else {
				dataFileNew.Request = []string{}
				res, err := RunHttp(dataFile.Api.Method, url, nil, header)
				resList = append(resList, res)
				errs = append(errs, err)
			}
		}
	}

	var sceneDataRecord SceneDataRecord
	lang := GetRequestLangage(dataFile.Single.Header)
	result, dst, err := dataFile.GetResult(lang, "again", filePath, header, dataFile.IsParallel, resList, nil, errs)

	if err != nil {
		sceneDataRecord.FailReason = fmt.Sprintf("%s", err)
	}

	err = WriteDataResultByFile(filePath, result, dst, hData.Product, hData.EnvType, err)

	if err != nil {
		Logger.Error("%s", err)
	}
	return
}
