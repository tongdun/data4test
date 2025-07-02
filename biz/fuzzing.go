package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"strconv"
)

func GetFuzzingData() (fuzzDefs []FuzzingDefinition, err error) {
	models.Orm.Table("fuzzing_definition").Find(&fuzzDefs)
	return
}

func CreateFuzzingBody(body map[string]interface{}, fuzzValue string) (bodys []map[string]interface{}, err error) {
	var bKeys []string
	for k, _ := range body {
		bKeys = append(bKeys, k)
	}

	for _, k := range bKeys {
		bodyTemp := CopyMap(body)
		bodyTemp[k] = fuzzValue
		bodys = append(bodys, bodyTemp)
	}

	return

}

func CreateFuzzingData(id string) (err error) {
	var apiDefinition ApiDefinition
	var apiRelation ApiRelation
	var qUrls []string

	s, _ := strconv.Atoi(id)
	models.Orm.Table("api_definition").Where("id = ?", s).Find(&apiDefinition)

	if len(apiDefinition.ApiId) == 0 {
		err1 := fmt.Errorf("Not Found API Definition info")
		err = err1
		Logger.Error("%s", err)
		return
	}

	models.Orm.Table("api_relation").Where("api_id = ?", apiDefinition.ApiId).Find(&apiRelation)

	var apiFuzzingData ApiFuzzingData
	apiFuzzingData.ApiDesc = apiDefinition.ApiDesc
	apiFuzzingData.ApiId = apiDefinition.ApiId
	apiFuzzingData.ApiModule = apiDefinition.ApiModule
	apiFuzzingData.RunNum = 1
	apiFuzzingData.App = apiDefinition.App
	rawUrl := apiRelation.GetRawUrl()
	customDepVars, _ := apiRelation.GetTestDepVars()

	depOutVars, err := apiDefinition.GetTestDataOutVars(customDepVars)

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	urls, err := GetUrl(rawUrl, depOutVars)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	querys, err := apiDefinition.GetQuery(depOutVars)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	bodys, err := apiDefinition.GetBody(depOutVars)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if len(querys) > 0 {
		for _, url := range urls {
			for _, query := range querys {
				qUrl := url + "?" + query
				qUrls = append(qUrls, qUrl)
			}
		}
	} else {
		qUrls = urls
	}

	fuzzDefs, err := GetFuzzingData()
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	for _, url := range qUrls {
		apiFuzzingData.UrlQuery = url
		if len(bodys) > 0 {
			for _, body := range bodys {
				for _, fuzzDef := range fuzzDefs {
					apiFuzzingData.DataDesc = fuzzDef.Name
					fuzzBodys, err := CreateFuzzingBody(body, fuzzDef.Value)
					if err != nil {
						Logger.Error("%s", err)
					}
					for _, fuzzBody := range fuzzBodys {
						bodyStr, err := json.Marshal(fuzzBody)
						if err != nil {
							Logger.Error("%s", err)
						}
						apiFuzzingData.Body = string(bodyStr)
						err = models.Orm.Table("api_fuzzing_data").Create(&apiFuzzingData).Error
						if err != nil {
							Logger.Error("%s", err)
						}
					}
				}
			}
		}
	}
	return
}
