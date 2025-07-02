package biz

import (
	"data4test/models"
	"encoding/json"
	"fmt"
	"github.com/GoAdminGroup/go-admin/template/types"
	"net/http"
	"os"
	"strings"
)

func GetAIModelConnectInfo(aiType string) (aiConnect AIConnect, err error) {
	var sysParameter SysParameter
	parameterName := "aiRunEngine"
	models.Orm.Table("sys_parameter").Where("name = ?", parameterName).Find(&sysParameter)
	if len(sysParameter.ValueList) == 0 {
		err = fmt.Errorf("系统参数中未定义[%s]参数的值，请核对~", parameterName)
		Logger.Error("%s", err)
		return
	}

	valueDefine := make(map[string]AIConnect)
	json.Unmarshal([]byte(sysParameter.ValueList), &valueDefine)

	if v, ok := valueDefine[aiType]; ok {
		aiConnect = v
	} else {
		err = fmt.Errorf("[%s]参数中未定义[%s]的值，请核对~", parameterName, aiType)
		Logger.Error("%s", err)
		return
	}

	if len(aiConnect.BaseUrl) == 0 || len(aiConnect.ApiKey) == 0 {
		err = fmt.Errorf("[%s]连接信息不全: %v，请核对~", aiType, aiConnect)
		Logger.Error("%s", err)
	}

	return
}

func GetNoSelectOption(desc string) (optionList []types.FieldOption) {
	var option types.FieldOption
	option.Value = desc
	option.Text = desc
	optionList = append(optionList, option)
	return
}

func DetectFileType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return "", err
	}
	// rawFileType返回信息为application/pdf格式
	rawFileType := http.DetectContentType(buffer)
	var tFileType string
	if strings.Contains(rawFileType, "/") {
		tFileType = strings.ToUpper(strings.Split(rawFileType, "/")[1])
	} else {
		tFileType = strings.ToUpper(rawFileType)
	}
	return tFileType, nil
}

func CallModel(aiConnect AIConnect, method, path string, data map[string]interface{}) (resBody string, err error) {
	url := fmt.Sprintf("%s/%s", aiConnect.BaseUrl, path)

	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)

	resByte, err := RunHttpJson(method, url, aiConnect.Timeout, data, header)

	resBody = string(resByte)
	return
}
