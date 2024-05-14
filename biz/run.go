package biz

import (
	"crypto/tls"
	"data4perf/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	// "bytes"
	"strings"
	"sync"
	// "net/url"
	//"gopkg.in/yaml.v2"
)

func ChkUniVar(name, app string) (b bool) {
	var ParamDef ParameterDefinition
	uniVarName := "uniVars"
	models.Orm.Table("parameter_definition").Where("app = ? and name = ? ", app, uniVarName).Find(&ParamDef)
	uniVars := GetListFromHtml(ParamDef.Value)
	for _, v := range uniVars {
		if strings.TrimSpace(v) == strings.TrimSpace(name) {
			b = true
			return
		}
	}

	return
}

func (api ApiDefinition) GetHeader() (header map[string]string, err error) {
	header = make(map[string]string)
	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", api.App).Find(&envConfig)
	if len(envConfig.Auth) == 0 {
		err = errors.New("Get Header: Not Found token value in Host")
		return
	}

	strTmp := GetStrFromHtml(envConfig.Auth)
	err = json.Unmarshal([]byte(strTmp), &header)
	if err != nil {
		return
	}
	return
}

func (handle ApiFuzzingData) GetHeader() (header map[string]string, err error) {
	header = make(map[string]string)
	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", handle.App).Find(&envConfig)
	if len(envConfig.Auth) == 0 {
		err = errors.New("Get Header: Not Found token value in Host")
		return
	}

	strTmp := GetStrFromHtml(envConfig.Auth)
	err = json.Unmarshal([]byte(strTmp), &header)
	if err != nil {
		return
	}
	return
}

func (apiRelation ApiRelation) GetHeader() (header map[string]string, err error) {
	header = make(map[string]string)
	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", apiRelation.App).Find(&envConfig)
	if len(envConfig.Auth) == 0 {
		err = errors.New("Get Header: Not Found")
		return
	}
	strTmp := GetStrFromHtml(envConfig.Auth)
	err = json.Unmarshal([]byte(strTmp), &header)
	if err != nil {
		return
	}
	return
}

func GetListFromHtml(rawStr string) (strList []string) {
	newTxt := strings.Replace(rawStr, "<br>", ",", -1)
	newTxt = strings.Replace(newTxt, "\n", ",", -1)
	newTxt = strings.Replace(newTxt, " ", ",", -1)
	newTxt = strings.Replace(newTxt, "<p>", " ", -1)
	newTxt = strings.Replace(newTxt, "</p>", ",", -1)
	newTxt = strings.Replace(newTxt, "<div>", " ", -1)
	newTxt = strings.Replace(newTxt, "</div>", ",", -1)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(newTxt))
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	afterTxt := dom.Text()

	if len(afterTxt) == 0 {
		return
	}
	tmpList := strings.Split(afterTxt, ",")
	for _, rawItem := range tmpList {
		item := strings.Trim(rawItem, " ")
		if len(item) == 0 {
			continue
		}
		strList = append(strList, item)
	}
	return
}

func GetStrFromHtml(rawStr string) (afterTxt string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawStr))
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	//handle := doc.Find("code")
	//afterTxt = handle.Text()
	//if len(afterTxt) == 0 {
	//	Logger.Warning("未找到有效信息，请核对~")
	//}

	afterTxt = doc.Text()
	if len(afterTxt) == 0 {
		Logger.Warning("未找到有效信息，请核对~")
	}

	return
}

func (apiTestData ApiTestData) GetHeader() (header map[string]string, err error) {
	header = make(map[string]string)
	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", apiTestData.App).Find(&envConfig)
	if len(envConfig.Auth) == 0 {
		err = errors.New("Get Header: Not Found")
		return
	}
	strTmp := GetStrFromHtml(envConfig.Auth)
	err = json.Unmarshal([]byte(strTmp), &header)
	if err != nil {
		return
	}
	return
}

func (apiCase ApiRelation) SaveTestResult(requestData, response map[string]interface{}) (err error) {
	var testResult ApiTestResult
	var dbResult DbApiTestResult
	testResult.ApiId = apiCase.ApiId

	if len(apiCase.OutVars) > 0 {
		var mapOutVar map[string]string
		mapOutVar = make(map[string]string)
		err = json.Unmarshal([]byte(apiCase.OutVars), &mapOutVar)
		retOut, _ := apiCase.GetLastValue(mapOutVar, response)
		var outByte []byte
		outByte, err = json.Marshal(retOut)
		testResult.OutVars = string(outByte)
		testResult.OutVars = string(outByte)
	}
	testResult.App = apiCase.App
	curTime := time.Now()
	testResult.UpdatedAt = curTime.Format(baseFormat)
	models.Orm.Table("api_test_result").Where("app = ? and api_id = ?", apiCase.App, apiCase.ApiId).Find(&dbResult)

	if len(dbResult.ApiId) == 0 {
		err = models.Orm.Table("api_test_result").Create(testResult).Error
	} else {
		err = models.Orm.Table("api_test_result").Where("id = ?", dbResult.Id).Update(testResult).Error
	}
	return
}

func (apiCase ApiRelation) IsThread() (b bool) {
	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", apiCase.App).Find(&envConfig)
	if envConfig.Threading == "yes" {
		b = true
	}
	return
}

func (apiCase ApiRelation) SaveDetailResult(url string, requestData map[string]interface{}, response map[string]interface{}) (err error) {
	var testDetail ApiTestDetail
	if value, ok := response["status"]; ok {
		varType := fmt.Sprintf("%T", value)
		var strValue string
		switch varType {
		case "float64":
			tmpVar := value.(float64)
			strValue = strconv.FormatFloat(tmpVar, 'f', 0, 64)
		case "string":
			strValue = value.(string)
		case "bool":
			strValue = strconv.FormatBool(value.(bool))
		default:
			Logger.Error("A Problem had Occured at Get Out Var: %s", value)
		}

		testDetail.TestResult = strValue
	}

	if value, ok := response["message"]; ok {
		strValue := value.(string)
		testDetail.FailReason = strValue
	}

	header, err := apiCase.GetHeader()
	tmpHeader, _ := json.Marshal(header)
	testDetail.Header = string(tmpHeader)
	testDetail.ApiId = apiCase.ApiId
	testDetail.ApiDesc = apiCase.ApiDesc
	testDetail.Url = url

	reqByte, err := json.Marshal(requestData)
	testDetail.Body = string(reqByte)

	resByte, err := json.Marshal(response)
	testDetail.Response = string(resByte)

	testDetail.App = apiCase.App
	curTime := time.Now()
	testDetail.CreatedAt = curTime.Format(baseFormat)

	err = models.Orm.Table("api_test_detail").Create(testDetail).Error
	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func (apiTestData ApiTestData) SaveDetailResult(requestData map[string]interface{}, response map[string]interface{}) (err error) {
	var testDetail ApiTestDetail
	if value, ok := response["status"]; ok {
		varType := fmt.Sprintf("%T", value)
		var strValue string
		switch varType {
		case "float64":
			tmpVar := value.(float64)
			strValue = strconv.FormatFloat(tmpVar, 'f', 0, 64)
		case "string":
			strValue = value.(string)
		case "bool":
			strValue = strconv.FormatBool(value.(bool))
		default:
			Logger.Error("A Problem had Occured at Get Out Var: %s", value)
		}

		testDetail.TestResult = strValue
	}

	if value, ok := response["message"]; ok {
		strValue := value.(string)
		testDetail.FailReason = strValue
	}

	header, err := apiTestData.GetHeader()
	tmpHeader, _ := json.Marshal(header)
	testDetail.Header = string(tmpHeader)
	testDetail.ApiId = apiTestData.ApiId
	testDetail.ApiDesc = apiTestData.ApiDesc
	testDetail.Url = apiTestData.UrlQuery
	testDetail.DataDesc = apiTestData.DataDesc

	reqByte, err := json.Marshal(requestData)
	testDetail.Body = string(reqByte)

	resByte, err := json.Marshal(response)
	testDetail.Response = string(resByte)

	testDetail.App = apiTestData.App
	curTime := time.Now()
	testDetail.CreatedAt = curTime.Format(baseFormat)

	err = models.Orm.Table("api_test_detail").Create(testDetail).Error
	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func (testData ApiFuzzingData) SaveDetailResult(requestData map[string]interface{}, response map[string]interface{}) (err error) {
	var testDetail ApiTestDetail
	if value, ok := response["status"]; ok {
		varType := fmt.Sprintf("%T", value)
		var strValue string
		switch varType {
		case "float64":
			tmpVar := value.(float64)
			strValue = strconv.FormatFloat(tmpVar, 'f', 0, 64)
		case "string":
			strValue = value.(string)
		case "bool":
			strValue = strconv.FormatBool(value.(bool))
		default:
			Logger.Error("A Problem had Occured at Get Out Var: %s", value)
		}

		testDetail.TestResult = strValue
	}

	if value, ok := response["message"]; ok {
		strValue := value.(string)
		testDetail.FailReason = strValue
	}

	header, err := testData.GetHeader()
	tmpHeader, _ := json.Marshal(header)
	testDetail.Header = string(tmpHeader)
	testDetail.ApiId = testData.ApiId
	testDetail.ApiDesc = testData.ApiDesc
	testDetail.DataDesc = testData.DataDesc
	testDetail.Url = testData.UrlQuery

	reqByte, err := json.Marshal(requestData)
	testDetail.Body = string(reqByte)

	resByte, err := json.Marshal(response)
	testDetail.Response = string(resByte)

	testDetail.App = testData.App
	curTime := time.Now()
	testDetail.CreatedAt = curTime.Format(baseFormat)

	err = models.Orm.Table("api_test_detail").Create(testDetail).Error
	if err != nil {
		Logger.Error("%s", err)
	}

	return
}

func (api ApiDefinition) Run(UrlQuery string, data map[string]interface{}) (resRaw map[string]interface{}, err error) {
	var req *http.Request
	Logger.Debug("Raw Data: %+v", data)
	i := 0
	for k := range data {
		i++
		if ChkUniVar(k, api.App) {
			data[k] = GetRandomStr(12, "")
			if i == len(data) {
				Logger.Debug("After Data: %+v", data)
			}

		}
	}
	var kvlist []string

	if len(data) > 0 {
		for k, value := range data {
			varType := fmt.Sprintf("%T", value)
			var strValue string
			switch varType {
			case "float64":
				tmpVar := value.(float64)
				strValue = strconv.FormatFloat(tmpVar, 'f', 0, 64)
			case "string":
				strValue = value.(string)
			case "bool":
				strValue = strconv.FormatBool(value.(bool))
			default:
				Logger.Error("A Problem had Occured at Get Out Var: %s", value)
			}
			kvstr := fmt.Sprintf("%s=%v", k, strValue)
			kvlist = append(kvlist, kvstr)
		}
	}
	payload := strings.NewReader(strings.Join(kvlist, "&"))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	switch api.HttpMethod {
	case "get":
		req, _ = http.NewRequest("GET", UrlQuery, nil)
	case "post":
		req, _ = http.NewRequest("POST", UrlQuery, payload)
	case "delete":
		req, _ = http.NewRequest("DELETE", UrlQuery, payload)
	case "put":
		req, _ = http.NewRequest("PUT", UrlQuery, payload)
	}
	var hKeys []string
	header, err := api.GetHeader()
	for k := range header {
		hKeys = append(hKeys, k)
	}
	for _, v := range hKeys {
		req.Header.Add(v, header[v])
	}

	resp, err := client.Do(req)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Error("%s", err)
	}

	Logger.Debug("response str: %s", string(body))
	resRaw = make(map[string]interface{})
	err = json.Unmarshal(body, &resRaw)
	if err != nil {
		Logger.Error("%s", err)
	}
	if string(body) == "404 page not found" {
		return
	}
	return
}

func GetUniBody(body map[string]interface{}, app string) (bodyAfter map[string]interface{}, err error) {
	bodyAfter = CopyMap(body)
	for k := range body {
		if ChkUniVar(k, app) {
			bodyAfter[k] = GetRandomStr(12, "")
		}
	}
	return
}

func GetAfterBody(lang string, body map[string]interface{}, depOutVars map[string][]interface{}) (bodyAfter map[string]interface{}, err error) {
	treeData := make(map[string]interface{})
	for k, v := range body {
		varType := fmt.Sprintf("%T", v)
		if varType == "[]interface {}" {
			rawList := v.([]interface{})
			var afterList []interface{}
			count := 0
			for _, item := range rawList {
				subVarType := fmt.Sprintf("%T", item)
				if subVarType == "string" {
					strK := Interface2Str(item)
					//if err1 != nil {
					//	Logger.Error("%s", err1)
					//	return body, err1
					//}
					afterStr, err1 := RawStr2MadeStr(lang, "", strK, count, depOutVars)

					if err1 != nil {
						Logger.Error("%s", err1)
						return body, err1
					}
					afterList = append(afterList, afterStr)
				} else if subVarType == "map[interface {}]interface {}" {
					afterV := item.(map[interface{}]interface{})
					for subK, subV := range afterV {
						subVarType := fmt.Sprintf("%T", subV)
						strAfterK := Interface2Str(subK)
						//if err2 != nil {
						//	Logger.Error("%s", err2)
						//	break
						//}

						if subVarType == "[]interface {}" {
							sub4Item := item.(map[interface{}]interface{})
							Logger.Debug("subItem: %v", sub4Item)
							//for _, sub5Item := range sub4Item {
							//	sub5VarType := fmt.Sprintf("%T", sub5Item)
							//	str5K, _ := Interface2Str(sub5Item)
							//	afterStr, err3 := RawStr2MadeStr(lang, strAfterK, str5K, count, depOutVars)
							//	if err3 != nil {
							//		Logger.Error("%s", err3)
							//		return body, err3
							//	}
							//	afterV[strAfterK] = afterStr
							//}
							break
						} else {
							strK := Interface2Str(subV)
							//if err1 != nil {
							//	break
							//}
							afterStr, err3 := RawStr2MadeStr(lang, strAfterK, strK, count, depOutVars)
							if err3 != nil {
								Logger.Error("%s", err3)
								return body, err3
							}
							afterV[strAfterK] = afterStr
						}
					}
					afterList = append(afterList, afterV)
				} else if subVarType == "map[string]interface{}" {
					afterV := item.(map[string]interface{})
					for subK, subV := range afterV {
						strK := Interface2Str(subV)
						//if err1 != nil {
						//	Logger.Error("%s", err1)
						//	return body, err1
						//}
						afterStr, err2 := RawStr2MadeStr(lang, subK, strK, count, depOutVars)
						if err2 != nil {
							Logger.Error("%s", err2)
							return body, err2
						}
						afterV[subK] = afterStr
					}
					afterList = append(afterList, afterV)
				} else {
					Logger.Warning("不支持varType: %v的数据类型变量替换", varType)
				}
				body[k] = afterList
				count++
			}
		} else if varType == "interface {}" {
			afterV := v.(map[string]interface{})
			for subK, subV := range afterV {
				strK := Interface2Str(subV)
				//if err1 != nil {
				//	Logger.Error("%s", err1)
				//	return body, err1
				//}
				afterStr, err1 := RawStr2MadeStr(lang, subK, strK, 0, depOutVars)
				if err1 != nil {
					Logger.Error("%s", err1)
					return body, err1
				}
				afterV[subK] = afterStr
			}
			body[k] = afterV
		} else if varType == "string" {
			strK := Interface2Str(v)
			var order int
			var afterStr interface{}
			//if err1 != nil {
			//	Logger.Info("body key: %v", v)
			//	Logger.Error("%s", err1)
			//	return body, err1
			//}

			if strings.HasPrefix(strK, "{TreeData_") {
				treeData[k] = strK
				continue
			}

			afterStr, err1 := RawStr2MadeStr(lang, k, strK, order, depOutVars)
			if err1 != nil {
				return body, err1
			}

			body[k] = afterStr
		} else if varType == "map[interface {}]interface {}" {
			afterV := v.(map[interface{}]interface{})
			for subK, subV := range afterV {
				varType := fmt.Sprintf("%T", subV)
				if varType != "string" {
					afterV[subK] = subV
					continue
				}

				strK := Interface2Str(subV)
				//if err1 != nil {
				//	Logger.Error("%s", err1)
				//	return body, err1
				//}
				afterStr, err1 := RawStr2MadeStr(lang, subK.(string), strK, 0, depOutVars)
				if err1 != nil {
					Logger.Error("%s", err1)
					return body, err1
				}
				afterV[subK] = afterStr
			}
			body[k] = afterV
		} else if varType == "bool" || varType == "int" || varType == "float64" || varType == "float32" {
			body[k] = v
		} else {
			Logger.Warning("不支持varType: %v的数据类型变量替换", varType)
		}
	}
	var first, second, third string
	if len(treeData) > 0 {
		for _, v := range treeData {
			vStr := Interface2Str(v)
			treeDataKey, deep := GetTreeDataTag(vStr)
			if deep == 1 {
				first, second, third = GetTreeDataValue(treeDataKey, deep, first, second)
				break
			}
		}
		for _, v := range treeData {
			vStr := Interface2Str(v)
			treeDataKey, deep := GetTreeDataTag(vStr)
			if deep == 2 {
				first, second, third = GetTreeDataValue(treeDataKey, deep, first, second)
				break
			}

		}
		for _, v := range treeData {
			vStr := Interface2Str(v)
			treeDataKey, deep := GetTreeDataTag(vStr)
			if deep == 3 {
				first, second, third = GetTreeDataValue(treeDataKey, deep, first, second)
				break
			}
		}

		for k, v := range treeData {
			vStr := Interface2Str(v)
			_, deep := GetTreeDataTag(vStr)
			if deep == 1 {
				body[k] = first
			} else if deep == 2 {
				body[k] = second
			} else if deep == 3 {
				body[k] = third
			}
		}
	}

	return body, err
}

func GetAfterListBody(lang string, body []interface{}, depOutVars map[string][]interface{}) (bodyAfter []interface{}, err error) {
	for _, v := range body {
		if v == nil {
			continue
		}
		varType := fmt.Sprintf("%T", v)
		if varType == "interface {}" {
			afterV := v.(map[string]interface{})
			for subK, subV := range afterV {
				strK := Interface2Str(subV)
				//if err1 != nil {
				//	Logger.Error("%s", err1)
				//	return body, err1
				//}
				afterStr, err1 := RawStr2MadeStr(lang, subK, strK, 0, depOutVars)
				if err1 != nil {
					Logger.Error("%s", err1)
					return body, err1
				}
				afterV[subK] = afterStr
			}
			bodyAfter = append(bodyAfter, afterV)
		} else if varType == "string" {
			strK := Interface2Str(v)
			//if err1 != nil {
			//	Logger.Info("body key: %v", v)
			//	Logger.Error("%s", err1)
			//	return body, err1
			//}
			afterStr, err1 := RawStr2MadeStr(lang, "", strK, 0, depOutVars)
			if err1 != nil {
				Logger.Error("%s", err1)
				return body, err1
			}
			bodyAfter = append(bodyAfter, afterStr)

		} else if varType == "map[interface {}]interface {}" {
			afterV := v.(map[interface{}]interface{})
			for subK, subV := range afterV {
				strK := Interface2Str(subV)
				//if err1 != nil {
				//	Logger.Error("%s", err1)
				//	return body, err1
				//}
				afterStr, err1 := RawStr2MadeStr(lang, subK.(string), strK, 0, depOutVars)
				if err1 != nil {
					Logger.Error("%s", err1)
					return body, err1
				}
				afterV[subK] = afterStr
			}
			bodyAfter = append(bodyAfter, afterV)
		} else if varType == "map[string]interface {}" {
			vMap := v.(map[string]interface{})
			afterV := make(map[string]interface{})

			for subK, subV := range vMap {
				varType := fmt.Sprintf("%T", subV)
				if varType != "string" {
					Logger.Warning("不支持varType: %v的数据类型变量替换", varType)
					afterV[subK] = subV
				} else {
					strK := Interface2Str(subV)
					//if err1 != nil {
					//	Logger.Error("%s", err1)
					//	return body, err1
					//}

					afterStr, err1 := RawStr2MadeStr(lang, subK, strK, 0, depOutVars)
					if err1 != nil {
						Logger.Error("%s", err1)
						return body, err1
					}
					afterV[subK] = afterStr
				}
			}
			bodyAfter = append(bodyAfter, afterV)
		} else {
			Logger.Warning("不支持varType: %v的数据类型变量替换", varType)
			Logger.Debug("body: %v", body)
			return body, err
		}
	}
	return
}

func RunTestData(id string) (err error) {
	var testData ApiTestData
	s, _ := strconv.Atoi(id)
	models.Orm.Table("api_test_data").Where("id = ?", s).Find(&testData)
	if len(testData.ApiId) == 0 {
		err = fmt.Errorf("Not Found API Test Data")
		Logger.Error("%s", err)
		return
	}
	var api ApiDefinition
	models.Orm.Table("api_definition").Where("api_id = ?", testData.ApiId).Find(&api)
	if len(api.ApiId) == 0 {
		err = fmt.Errorf("Not Found API[%s] info", testData.ApiId)
		Logger.Error("%s", err)
		return
	}

	var apiRelation ApiRelation
	models.Orm.Table("api_relation").Where("api_id = ?", testData.ApiId).Find(&apiRelation)
	if len(apiRelation.ApiId) == 0 {
		err = fmt.Errorf("Not Found API[%s] info", testData.ApiId)
		Logger.Error("%s", err)
		return
	}
	var data, response map[string]interface{}
	response = make(map[string]interface{})
	data = make(map[string]interface{})
	if err = json.Unmarshal([]byte(testData.Body), &data); err != nil {
		Logger.Error("%s", err)
		return
	}

	if testData.RunNum > 1 {
		if api.HttpMethod == "get" {
			if apiRelation.IsThread() {
				wg := sync.WaitGroup{}
				for i := 0; i < testData.RunNum; i++ {
					wg.Add(1)
					go func(times int) {
						Logger.Info("并发执行次数: %d", times)
						response, err = api.Run(testData.UrlQuery, nil)
						if err != nil {
							Logger.Error("%s", err)
						}
						apiRelation.SaveTestResult(nil, response)
						testData.SaveDetailResult(nil, response)
						wg.Done()
					}(i)
				}
			} else {
				for i := 0; i < testData.RunNum; i++ {
					Logger.Info("串行执行次数: %d", i+1)
					response, err = api.Run(testData.UrlQuery, nil)
					err = apiRelation.SaveTestResult(nil, response)
					if err != nil {
						Logger.Error("%s", err)
					}
					err = testData.SaveDetailResult(nil, response)
					if err != nil {
						Logger.Error("%s", err)
					}
				}
			}
		} else {
			if apiRelation.IsThread() {
				wg := sync.WaitGroup{}
				for i := 0; i < testData.RunNum; i++ {
					afterData, _ := GetUniBody(data, testData.App)
					wg.Add(1)
					go func(times int, inBody map[string]interface{}) {
						Logger.Info("串行执行次数: %d", times+1)
						response, err = api.Run(testData.UrlQuery, inBody)
						if err != nil {
							Logger.Error("%s", err)
						}
						apiRelation.SaveTestResult(inBody, response)
						testData.SaveDetailResult(inBody, response)
						wg.Done()
					}(i, afterData)
				}
			} else {
				for i := 0; i < testData.RunNum; i++ {
					afterData, err := GetUniBody(data, testData.App)
					Logger.Info("串行执行次数: %d", i+1)
					response, err = api.Run(testData.UrlQuery, afterData)
					err = apiRelation.SaveTestResult(afterData, response)
					if err != nil {
						Logger.Error("%s", err)
					}
					err = testData.SaveDetailResult(afterData, response)
					if err != nil {
						Logger.Error("%s", err)
					}
				}
			}
		}
	} else {
		afterData, err := GetUniBody(data, testData.App)
		response, err = api.Run(testData.UrlQuery, afterData)
		err = apiRelation.SaveTestResult(afterData, response)
		if err != nil {
			Logger.Error("%s", err)
		}
		err = testData.SaveDetailResult(afterData, response)
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	return
}

func RunFuzzData(id string) (err error) {
	var testData ApiFuzzingData
	s, _ := strconv.Atoi(id)
	models.Orm.Table("api_fuzzing_data").Where("id = ?", s).Find(&testData)
	if len(testData.ApiId) == 0 {
		err = fmt.Errorf("Not Found API Test Data")
		Logger.Error("%s", err)
		return
	}
	var api ApiDefinition
	models.Orm.Table("api_definition").Where("api_id = ?", testData.ApiId).Find(&api)
	if len(api.ApiId) == 0 {
		err = fmt.Errorf("Not Found API[%s] info", testData.ApiId)
		Logger.Error("%s", err)
		return
	}

	var apiRelation ApiRelation
	models.Orm.Table("api_relation").Where("api_id = ?", testData.ApiId).Find(&apiRelation)
	if len(apiRelation.ApiId) == 0 {
		err = fmt.Errorf("Not Found API[%s] info", testData.ApiId)
		Logger.Error("%s", err)
		return
	}
	var data, response map[string]interface{}
	response = make(map[string]interface{})
	data = make(map[string]interface{})
	if err = json.Unmarshal([]byte(testData.Body), &data); err != nil {
		Logger.Error("%s", err)
		return
	}

	if testData.RunNum > 1 {
		if api.HttpMethod == "get" {
			if apiRelation.IsThread() {
				wg := sync.WaitGroup{}
				for i := 0; i < testData.RunNum; i++ {
					wg.Add(1)
					go func(times int) {
						Logger.Info("Thread Run Times: %d", times)
						response, err = api.Run(testData.UrlQuery, nil)
						if err != nil {
							Logger.Error("%s", err)
						}
						apiRelation.SaveTestResult(nil, response)
						testData.SaveDetailResult(nil, response)
						wg.Done()
					}(i)
				}
			} else {
				for i := 0; i < testData.RunNum; i++ {
					Logger.Info("Serial Run Times: %d", i+1)
					response, err = api.Run(testData.UrlQuery, nil)
					err = apiRelation.SaveTestResult(nil, response)
					if err != nil {
						Logger.Error("%s", err)
					}
					err = testData.SaveDetailResult(nil, response)
					if err != nil {
						Logger.Error("%s", err)
					}
				}
			}
		} else {
			if apiRelation.IsThread() {
				wg := sync.WaitGroup{}
				for i := 0; i < testData.RunNum; i++ {
					afterData, _ := GetUniBody(data, testData.App)
					wg.Add(1)
					go func(times int, inBody map[string]interface{}) {
						Logger.Info("Thread Run Times: %d", times+1)
						response, err = api.Run(testData.UrlQuery, inBody)
						if err != nil {
							Logger.Error("%s", err)
						}
						apiRelation.SaveTestResult(inBody, response)
						testData.SaveDetailResult(inBody, response)
						wg.Done()
					}(i, afterData)
				}
			} else {
				for i := 0; i < testData.RunNum; i++ {
					afterData, err := GetUniBody(data, testData.App)
					Logger.Info("Serial Run Times: %d", i+1)
					response, err = api.Run(testData.UrlQuery, afterData)
					err = apiRelation.SaveTestResult(afterData, response)
					if err != nil {
						Logger.Error("%s", err)
					}
					err = testData.SaveDetailResult(afterData, response)
					if err != nil {
						Logger.Error("%s", err)
					}
				}
			}
		}
	} else {
		afterData, err := GetUniBody(data, testData.App)
		response, err = api.Run(testData.UrlQuery, afterData)
		err = apiRelation.SaveTestResult(afterData, response)
		if err != nil {
			Logger.Error("%s", err)
		}
		err = testData.SaveDetailResult(afterData, response)
		if err != nil {
			Logger.Error("%s", err)
		}
	}

	return
}

func RunAgain(id string) (err error) {
	var apiTestDetail ApiTestDetail
	models.Orm.Table("api_test_detail").Where("id = ?", id).Find(&apiTestDetail)

	if len(apiTestDetail.ApiId) == 0 {
		err = fmt.Errorf("Not Found API test detail info")
		Logger.Error("Error: %s, ID: %s", err, id)
		return
	}
	var api ApiDefinition
	models.Orm.Table("api_definition").Where("api_id = ?", apiTestDetail.ApiId).Find(&api)
	if len(api.ApiId) == 0 {
		err = fmt.Errorf("Not Found API[%s] info", apiTestDetail.ApiId)
		Logger.Error("Error: %s, ID: %s", err, id)
		return
	}

	var apiCase ApiRelation
	models.Orm.Table("api_relation").Where("api_id = ?", apiTestDetail.ApiId).Find(&apiCase)
	if len(apiCase.ApiId) == 0 {
		err = fmt.Errorf("Not Found API[%s] info", apiTestDetail.ApiId)
		Logger.Error("%s", err)
		return
	}
	var data, response map[string]interface{}
	response = make(map[string]interface{})
	if err = json.Unmarshal([]byte(GetStrFromHtml(apiTestDetail.Body)), &data); err != nil {
		Logger.Error("%s", err)
		return
	}
	response, err = api.Run(apiTestDetail.Url, data)
	apiCase.SaveTestResult(data, response)
	apiCase.SaveDetailResult(apiTestDetail.Url, data, response)
	return
}
