package biz

import (
	// "bytes"
	// "crypto/tls"
	"encoding/json"
	// "errors"
	"fmt"
	// "io/ioutil"
	"math/rand"
	// "net/http"
	"regexp"
	"strconv"
	"strings"

	// "github.com/ernestosuarez/itertools"

	"data4perf/models"
)

func (apiDefinition ApiDefinition) GetTestDataOutVars(customDepVars map[string]string) (retOutVars map[string]string, err error) {
	var allVar map[string]interface{}
	allVar = make(map[string]interface{})
	retOutVars = make(map[string]string)
	for _, item := range apiDefinition.PathVariable {
		allVar[item.Name] = item.ValueType
	}

	for _, item := range apiDefinition.QueryParameter {
		allVar[item.Name] = item.ValueType
	}

	for _, item := range apiDefinition.Body {
		allVar[item.Name] = item.ValueType
	}

	for k, v := range allVar {
		if value, ok := customDepVars[k]; ok {
			retOutVars[k] = value
		} else {
			switch v {
			case "integer", "int":
				retOutVars[k] = strconv.Itoa(rand.Intn(100))
			case "string":
				retOutVars[k] = GetRandomStr(8, "")
			case "bool":
				retOutVars[k] = "true,false"
			default:
				retOutVars[k] = GetRandomStr(8, "")
			}
		}
	}

	return
}

func CreateTestData(id string) (err error) {
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

	var apiTestData ApiTestData
	apiTestData.ApiDesc = apiDefinition.ApiDesc
	apiTestData.ApiId = apiDefinition.ApiId
	apiTestData.ApiModule = apiDefinition.ApiModule
	apiTestData.RunNum = 1
	apiTestData.App = apiDefinition.App
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

	for _, url := range qUrls {
		apiTestData.UrlQuery = url
		if len(bodys) > 0 {
			for _, body := range bodys {
				bodyStr, err := json.Marshal(body)
				apiTestData.Body = string(bodyStr)
				err = models.Orm.Table("api_test_data").Create(&apiTestData).Error
				if err != nil {
					Logger.Error("%s", err)
				}
			}
		}
	}
	return
}

func (apiCase ApiRelation) GetLastValue(mapped map[string]string, data map[string]interface{}) (map[string]string, map[string]interface{}) {
	if data == nil {
		return mapped, nil
	}
	var mKeys, dKeys []string
	for k := range mapped {
		mKeys = append(mKeys, k)
	}

	var tmpMap map[string]interface{}
	tmpMap = make(map[string]interface{})
	var depDefs map[string]string
	depDefs = make(map[string]string)
	i := 0
	for _, v := range mKeys {
		if strings.Contains(mapped[v], "-") || strings.Contains(mapped[v], "*") {
			depDefs[v] = mapped[v]
			dKeys = append(dKeys, v)
			i++
		} else {
			varType := fmt.Sprintf("%T", data[mapped[v]])
			var strValue string
			switch varType {
			case "float64":
				tmpVar := data[mapped[v]].(float64)
				strValue = strconv.FormatFloat(tmpVar, 'f', 0, 64)
			case "string":
				strValue = data[mapped[v]].(string)
			case "bool":
				strValue = strconv.FormatBool(data[mapped[v]].(bool))
			default:
				Logger.Error("A Problem had Occured at Get Out Var: %s", v)
			}
			mapped[v] = strValue
		}
	}
	if i == 0 {
		Logger.Debug("mapped: %+v", mapped)
		return mapped, nil
	}

	j := 0
	for _, v := range dKeys {
		items := strings.SplitN(mapped[v], "-", 2)
		mapped[v] = items[1]
		tmpMap = data[items[0]].(map[string]interface{})
		j++
		if j == len(dKeys) {
			return apiCase.GetLastValue(mapped, tmpMap)
		}

	}
	return mapped, nil
}

func (apiCase ApiRelation) GetAPI() (api ApiDefinition, err error) {
	models.Orm.Table("api_detail").Where("app = ? and api_id = ? ", apiCase.App, apiCase.ApiId).Find(&api)
	if len(api.ApiId) == 0 {
		err = fmt.Errorf("Not Found %s Api Relation", apiCase.ApiId)
		Logger.Error("%s", err)
		return
	}
	return
}

func (api ApiDefinition) GetFormatDepVars(depOutVars map[string]string) (retOutVars map[string]string, err error) {
	var allVar map[string]interface{}
	allVar = make(map[string]interface{})
	retOutVars = make(map[string]string)
	if len(api.PathVariable) > 0 {
		for _, item := range api.PathVariable {
			allVar[item.Name] = item.ValueType
		}
	}

	if len(api.QueryParameter) > 0 {
		for _, item := range api.QueryParameter {
			allVar[item.Name] = item.ValueType
		}
	}

	if len(api.Body) > 0 {
		for _, item := range api.Body {
			allVar[item.Name] = item.ValueType
		}
	}
	retOutVars = depOutVars
	var envConfig EnvConfig
	abDef := GetAbDef()
	models.Orm.Table("env_config").Where("app = ?", api.App).Find(&envConfig)
	for k, v := range allVar {
		if envConfig.Testmode == "abnormal" {
			if value, ok := depOutVars[k]; !ok {
				if v == "integer" || v == "int" {
					retOutVars[k] = abDef["intAb"]
				} else if v == "string" {
					retOutVars[k] = abDef["strAb"]
				} else if v == "array" {
					retOutVars[k] = abDef["arrAb"]
				} else if v == "bool" {
					retOutVars[k] = abDef["bool"]
				} else {
					retOutVars[k] = abDef["strAb"]
				}

			} else {
				retOutVars[k] = value
			}
		} else if envConfig.Testmode == "normal" {
			if value, ok := depOutVars[k]; !ok {
				if v == "integer" || v == "int" {
					retOutVars[k] = abDef["intNor"]
				} else if v == "string" {
					retOutVars[k] = abDef["strNor"]
				} else if v == "array" {
					retOutVars[k] = abDef["arrNor"]
				} else if v == "bool" {
					retOutVars[k] = abDef["bool"]
				} else {
					retOutVars[k] = abDef["strNor"]
				}
			} else {
				retOutVars[k] = value
			}
		}
	}
	Logger.Debug("retOutVars: %q", retOutVars)
	return
}

func (apiRelation ApiRelation) GetTestDepVars() (depOutVars map[string]string, err error) {
	var depCases []string
	depOutVars = make(map[string]string)
	if len(apiRelation.ParamApis) > 0 {
		pds := strings.Split(apiRelation.ParamApis, ",")
		depCases = append(depCases, pds...)
	}
	if len(apiRelation.PreApis) > 0 {
		bcs := strings.Split(apiRelation.PreApis, ",")
		depCases = append(depCases, bcs...)
	}
	depCases = append(depCases, apiRelation.ApiId)

	var testResult ApiTestResult
	var parameterDefinition ParameterDefinition
	for _, dep := range depCases {
		if dep == apiRelation.ApiId {
			continue
		}
		models.Orm.Table("api_test_result").Where("app = ? and api_id = ?", apiRelation.App, dep).Find(&testResult)
		if len(testResult.OutVars) > 0 {
			err = json.Unmarshal([]byte(testResult.OutVars), &depOutVars)
			if err != nil {
				Logger.Error("%s", err)
				return
			}
		}
	}

	for _, dep := range depCases {
		models.Orm.Table("parameter_definition").Where("app = ? and name = ?", apiRelation.App, dep).Find(&parameterDefinition)
		afterStr := GetStrFromHtml(parameterDefinition.Value)
		if len(afterStr) > 0 {
			err = json.Unmarshal([]byte(afterStr), &depOutVars)
			if err != nil {
				Logger.Error("%s", err)
				return
			}
		}

	}

	return
}

func (apiRelation ApiRelation) GetFuzzingDepVars() (depOutVars map[string]string, err error) {
	var depCases []string
	depOutVars = make(map[string]string)
	if len(apiRelation.ParamApis) > 0 {
		pds := strings.Split(apiRelation.ParamApis, ",")
		depCases = append(depCases, pds...)
	}
	if len(apiRelation.PreApis) > 0 {
		bcs := strings.Split(apiRelation.PreApis, ",")
		depCases = append(depCases, bcs...)
	}
	depCases = append(depCases, apiRelation.ApiId)

	var testResult ApiTestResult
	var parameterDefinition ParameterDefinition
	for _, dep := range depCases {
		if dep == apiRelation.ApiId {
			continue
		}
		models.Orm.Table("api_test_result").Where("app = ? and api_id = ?", apiRelation.App, dep).Find(&testResult)
		if len(testResult.OutVars) > 0 {
			Logger.Debug("testResult: %q", testResult)
			err = json.Unmarshal([]byte(testResult.OutVars), &depOutVars)
			if err != nil {
				Logger.Error("%s", err)
				return
			}
		}
	}

	for _, dep := range depCases {
		models.Orm.Table("parameter_definition").Where("app_id = ? and name = ?", apiRelation.ApiId, dep).Find(&parameterDefinition)
		if len(parameterDefinition.Value) > 0 {
			err = json.Unmarshal([]byte(parameterDefinition.Value), &depOutVars)
			if err != nil {
				Logger.Error("%s", err)
				return
			}
		}

	}
	return
}

func (apiCase ApiRelation) GetRawUrl() (url string) {
	var envConfig EnvConfig
	models.Orm.Table("env_config").Where("app = ?", apiCase.App).Find(&envConfig)
	paths := strings.Split(apiCase.ApiId, "_")
	url = envConfig.Protocol + "://" + envConfig.Ip + envConfig.Prepath + paths[1]
	return
}

func GetUrl(rawUrl string, depOutVars map[string]string) (urls []string, err error) {
	pathVarsReg := regexp.MustCompile(`{[:alpha:]+}`)
	var pathVars []string
	if !strings.Contains(rawUrl, "{") {
		urls = append(urls, rawUrl)
	} else {
		pathVars = pathVarsReg.FindAllString(rawUrl, -1)
		for _, v := range pathVars {
			str1 := v[1 : len(v)-1]
			if value, ok := depOutVars[str1]; ok {
				if !strings.Contains(value, ",") {
					url := strings.Replace(rawUrl, v, value, -1)
					urls = append(urls, url)
				} else {
					strList := strings.Split(value, ",")
					for _, sv := range strList {
						url := strings.Replace(rawUrl, v, sv, -1)
						urls = append(urls, url)
					}
				}
			} else {
				err = fmt.Errorf("Get URL: Not Found [%q] value in DepOutVars", v)
				Logger.Error("%s", err)
				return
			}

		}
	}

	return

}

func (api ApiDefinition) GetQuery(depOutVars map[string]string) (querys []string, err error) {
	if len(api.QueryParameter) == 0 {
		return
	}
	var mapQuery map[string]string
	mapQuery = make(map[string]string)
	for _, item := range api.QueryParameter {
		mapQuery[item.Name] = item.ValueType
	}

	var qKeys []string
	for k := range mapQuery {
		qKeys = append(qKeys, k)
	}
	var parameterDefinition ParameterDefinition
	for _, v := range qKeys {
		if value, ok := depOutVars[v]; ok {
			if !strings.Contains(value, ",") {
				queryStr := fmt.Sprintf("%s=%s", v, value)
				querys = append(querys, queryStr)
			} else {
				strList := strings.Split(value, ",")
				for _, sv := range strList {
					queryStr := fmt.Sprintf("%s=%s", v, sv)
					querys = append(querys, queryStr)
				}
			}
		} else {
			models.Orm.Table("parameter_definition").Where("app = ? and name = ?", api.App, v).Find(&parameterDefinition)
			if len(parameterDefinition.Value) == 0 {
				err = fmt.Errorf("Get Query: Not Found [%q] value in DepOutVars and ComVars", v)
				return
			} else {
				if !strings.Contains(parameterDefinition.Value, ",") {
					queryStr := fmt.Sprintf("%s=%s", v, parameterDefinition.Value)
					querys = append(querys, queryStr)
				} else {
					strList := strings.Split(value, ",")
					for _, sv := range strList {
						queryStr := fmt.Sprintf("%s=%s", v, sv)
						querys = append(querys, queryStr)
					}
				}
			}

		}

	}

	return
}

func (api ApiDefinition) GetBody(depOutVars map[string]string) (bodys []map[string]interface{}, err error) {
	if len(api.Body) == 0 {
		return
	}
	var mapBody map[string]interface{}
	mapBody = make(map[string]interface{})
	for _, item := range api.Body {
		mapBody[item.Name] = item.ValueType
	}

	var bKeys []string
	var multiVars []string
	var multiKeys []string
	var multiDict map[string]string
	multiDict = make(map[string]string)
	for k := range mapBody {
		bKeys = append(bKeys, k)
	}
	var intValue int
	var boolValue bool

	for _, v := range bKeys {
		if value, ok := depOutVars[v]; ok {
			if mapBody[v] == "array" {
				strList := strings.Split(value, ",")
				var intList []int
				for _, sv := range strList {
					intValue, err := strconv.Atoi(sv)
					if err == nil {
						intList = append(intList, intValue)
					}
				}
				if len(intList) > 0 {
					mapBody[v] = intList
				} else {
					mapBody[v] = strList
				}
			} else {
				if strings.Contains(value, ",") {
					multiDict[v] = value
					multiKeys = append(multiKeys, v)
					multiVars = append(multiVars, value)
				} else {
					if mapBody[v] == "int" || mapBody[v] == "integer" {
						intValue, err = strconv.Atoi(value)
						if err != nil {
							return
						}
						mapBody[v] = intValue
					} else if mapBody[v] == "boolean" || mapBody[v] == "bool" {
						boolValue, err = strconv.ParseBool(value)
						if err != nil {
							return
						}
						mapBody[v] = boolValue
					} else {
						mapBody[v] = value
					}
				}
			}
		} else {
			err = fmt.Errorf("Get Body: Not Found [%q] value in DepOutVars", v)
			Logger.Error("%s", err)
			return
		}
	}

	if len(multiVars) == 0 {
		bodys = append(bodys, mapBody)
		return
	}
	if len(multiVars) == 1 {
		strList1 := strings.Split(multiVars[0], ",")
		for _, sv := range strList1 {
			var body map[string]interface{}
			body = make(map[string]interface{})
			for mk, mv := range mapBody {
				body[mk] = mv
			}
			if mapBody[multiKeys[0]] == "int" || mapBody[multiKeys[0]] == "integer" {
				intValue, err := strconv.Atoi(sv)
				if err != nil {
					Logger.Error("%s", err)
				}
				body[multiKeys[0]] = intValue
			} else if mapBody[multiKeys[0]] == "boolean" || mapBody[multiKeys[0]] == "bool" {
				boolValue, err = strconv.ParseBool(sv)
				body[multiKeys[0]] = boolValue
			} else {
				body[multiKeys[0]] = sv
			}
			bodys = append(bodys, body)
		}

	} else if len(multiVars) == 2 {
		strList1 := strings.Split(multiVars[0], ",")

		strList2 := strings.Split(multiVars[1], ",")
		for _, sv := range strList1 {
			for _, ssv := range strList2 {
				var body map[string]interface{}
				body = make(map[string]interface{})

				for mk, mv := range mapBody {
					body[mk] = mv
				}

				for i := 0; i < 2; i++ {
					if mapBody[multiKeys[i]] == "int" || mapBody[multiKeys[i]] == "integer" {
						var intValue int
						if i == 0 {
							intValue, err = strconv.Atoi(sv)
							if err != nil {
								Logger.Error("%s", err)
							}
						} else if i == 1 {
							intValue, err = strconv.Atoi(ssv)
							if err != nil {
								Logger.Error("%s", err)
							}
						}
						body[multiKeys[i]] = intValue
					} else if mapBody[multiKeys[i]] == "boolean" || mapBody[multiKeys[i]] == "bool" {
						var boolValue bool
						if i == 0 {
							boolValue, err = strconv.ParseBool(sv)
							if err != nil {
								Logger.Error("%s", err)
							}
						} else if i == 1 {
							boolValue, err = strconv.ParseBool(ssv)
							if err != nil {
								Logger.Error("%s", err)
							}

						}
						body[multiKeys[i]] = boolValue

					}

				}
				bodys = append(bodys, body)
			}

		}

	}
	return
}
