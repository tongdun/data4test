package biz

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	netUrl "net/url"
	"os"
	"path/filepath"
	"strings"
)

func RunHttpFormData(method, url string, data map[string]interface{}, header map[string]interface{}) (res []byte, err error) {
	var req *http.Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if _, ok := header["Content-Type"]; !ok {
		err = fmt.Errorf("header 未正常定义，请核对")
		return
	}

	methodUpper := strings.ToUpper(method)
	client := &http.Client{Transport: tr}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for k, v := range data {
		strValue, err1 := Interface2Str(v)
		if err1 != nil {
			err = err1
			return
		}
		b, err2 := IsValueInSysParameter("fileName", k)
		if err2 != nil {
			err = err2
			return
		}
		if header["Content-Type"] == "multipart/form-data" && b {
			filePath := fmt.Sprintf("%s/%s", UploadBasePath, strValue)
			_, errTmp := os.Stat(filePath)
			if errTmp != nil {
				if os.IsNotExist(errTmp) {
					err = fmt.Errorf("[%s]文件不存在，请核对", filePath)
					Logger.Error("%s", err)
					return
				} else {
					Logger.Error("%s", errTmp)
					err = errTmp
					return
				}
			}
			file, errFile2 := os.Open(filePath)
			defer file.Close()

			part2, errFile2 := writer.CreateFormFile(k, filepath.Base(filePath))
			_, errFile2 = io.Copy(part2, file)
			if errFile2 != nil {
				err = errFile2
				Logger.Error("%s", err)
				return
			}
		} else {
			_ = writer.WriteField(k, strValue)
		}
	}

	err = writer.Close()
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if methodUpper == "GET" {
		if len(data) > 0 {
			uri, err1 := netUrl.Parse(url)
			if err1 != nil {
				err = err1
				Logger.Error("%s", err)
				return
			}

			tmpData := make(netUrl.Values)
			for k, v := range data {
				strK, _ := Interface2Str(v)
				tmpData[k] = []string{strK}
				uri.RawQuery = tmpData.Encode()
			}

			req, err = http.NewRequest(methodUpper, uri.String(), nil)
			Logger.Info("method: %s, url: %s, data: %v", methodUpper, uri.Scheme, data)
		} else {
			req, err = http.NewRequest(methodUpper, url, nil)
			Logger.Info("method: %s, url: %s, data: %v", methodUpper, url, data)

		}

	} else {
		req, err = http.NewRequest(methodUpper, url, payload)
		Logger.Info("method: %s, url: %s, data: %v", methodUpper, url, data)
		//LogHandle.Printf("payload: %v",payload)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	for k, v := range header {
		if k == "Content-Type" {
			continue
		}
		valueStr, _ := Interface2Str(v)
		req.Header.Add(k, valueStr)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Error("%s", err)
	}
	Logger.Debug("resBody: %s", string(resBody))

	// 返回500的是否需要拦截
	//if resp.StatusCode != 200 || resp.StatusCode != 500 {
	if resp.StatusCode != 200 {
		err = fmt.Errorf("请求失败，返回码: %d, 返回信息: %s", resp.StatusCode, string(resBody))
		Logger.Error("%s", err)
	}
	return resBody, err
}

func RunHttpUrlencoded(method, url string, data map[string]interface{}, header map[string]interface{}) (res []byte, err error) {
	var req *http.Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if _, ok := header["Content-Type"]; !ok {
		err = fmt.Errorf("header 未正常定义，请核对")
		Logger.Error("%s", err)
		return
	}

	client := &http.Client{Transport: tr}
	methodUpper := strings.ToUpper(method)

	if methodUpper == "GET" && len(data) > 0 {
		uri, err1 := netUrl.Parse(url)
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}

		tmpData := make(netUrl.Values)
		for k, v := range data {
			strK, _ := Interface2Str(v)
			tmpData[k] = []string{strK}
			uri.RawQuery = tmpData.Encode()
		}

		req, err = http.NewRequest(methodUpper, uri.String(), nil)
	} else {
		dataPayload := netUrl.Values{}
		for k, v := range data {
			strValue, err1 := Interface2Str(v)
			if err1 != nil {
				err = err1
				return
			}

			dataPayload.Add(k, strValue)
		}

		payload := strings.NewReader(dataPayload.Encode())
		req, err = http.NewRequest(methodUpper, url, payload)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	for k, v := range header {
		vStr, _ := Interface2Str(v)
		req.Header.Add(k, vStr)

	}

	resp, err := client.Do(req)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Error("%s", err)
	}

	// 返回500的是否需要拦截
	//if resp.StatusCode != 200 || resp.StatusCode != 500 {
	if resp.StatusCode != 200 {
		err = fmt.Errorf("请求失败，返回码: %d, 返回信息: %s", resp.StatusCode, string(resBody))
		Logger.Error("%s", err)
	}
	return resBody, err
}

func RunHttpJson(method, url string, data map[string]interface{}, header map[string]interface{}) (res []byte, err error) {
	var req *http.Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if _, ok := header["Content-Type"]; !ok {
		err = fmt.Errorf("header 未正常定义，请核对")
		return
	}

	methodUpper := strings.ToUpper(method)
	client := &http.Client{Transport: tr}

	if methodUpper == "GET" && len(data) > 0 {
		uri, err1 := netUrl.Parse(url)
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}

		tmpData := make(netUrl.Values)
		for k, v := range data {
			strK, _ := Interface2Str(v)
			tmpData[k] = []string{strK}
			uri.RawQuery = tmpData.Encode()
		}

		req, err = http.NewRequest(methodUpper, uri.String(), nil)
	} else {
		var reader []byte
		reader, err = json.Marshal(data)
		if err != nil {
			var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
			readerNew, err2 := jsonNew.Marshal(&data)
			if err2 != nil {
				Logger.Error("%s", err2)
				err = err2
				return
			}
			reader = readerNew
		}

		payload := strings.NewReader(string(reader))

		req, err = http.NewRequest(methodUpper, url, payload)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	for k, v := range header {
		vStr, _ := Interface2Str(v)
		req.Header.Add(k, vStr)
	}

	resp, err := client.Do(req)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Error("%s", err)
	}

	// 返回500的是否需要拦截
	//if resp.StatusCode != 200 || resp.StatusCode != 500 {
	if resp.StatusCode != 200 {
		err = fmt.Errorf("请求失败，返回码: %d, 返回信息: %s", resp.StatusCode, string(resBody))
		Logger.Error("%s", err)
	}
	return resBody, err
}

func RunHttpJsonList(method, url string, data []interface{}, header map[string]interface{}) (res []byte, err error) {
	var req *http.Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if _, ok := header["Content-Type"]; !ok {
		err = fmt.Errorf("header 未正常定义，请核对")
		return
	}

	methodUpper := strings.ToUpper(method)
	client := &http.Client{Transport: tr}

	if methodUpper == "GET" {
		return
	} else {
		var jsonNew = jsoniter.ConfigCompatibleWithStandardLibrary
		readerNew, err2 := jsonNew.Marshal(&data)
		if err2 != nil {
			Logger.Error("%s", err2)
			err = err2
			return
		}
		//reader = readerNew
		payload := strings.NewReader(string(readerNew))

		req, err = http.NewRequest(methodUpper, url, payload)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	for k, v := range header {
		vStr, _ := Interface2Str(v)
		req.Header.Add(k, vStr)
	}

	resp, err := client.Do(req)
	if err != nil {
		Logger.Error("%s", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Error("%s", err)
	}
	Logger.Debug("resBody: %s", string(resBody))

	// 返回500的是否需要拦截
	//if resp.StatusCode != 200 || resp.StatusCode != 500 {
	if resp.StatusCode != 200 {
		err = fmt.Errorf("请求失败，返回码: %d, 返回信息: %s", resp.StatusCode, string(resBody))
		Logger.Error("%s", err)
	}
	return resBody, err
}

func RunHttp(method, url string, data map[string]interface{}, header map[string]interface{}) (res []byte, err error) {
	if _, ok := header["Content-Type"]; !ok {
		err = fmt.Errorf("header 未正常定义，请核对")
		return
	}

	contentTypeRaw, _ := Interface2Str(header["Content-Type"])

	var contentType string

	if strings.Contains(contentTypeRaw, "application/x-www-form-urlencoded") {
		contentType = "application/x-www-form-urlencoded"
	} else if strings.Contains(contentTypeRaw, "multipart/form-data") {
		contentType = "multipart/form-data"
	} else if strings.Contains(contentTypeRaw, "form-data") {
		contentType = "form-data"
	} else if strings.Contains(contentTypeRaw, "application/json") {
		contentType = "application/json"
	}

	switch contentType {
	case "application/x-www-form-urlencoded":
		res, err = RunHttpUrlencoded(method, url, data, header)
	case "form-data", "multipart/form-data":
		res, err = RunHttpFormData(method, url, data, header)
	case "application/json":
		res, err = RunHttpJson(method, url, data, header)
	default:
		res, err = RunHttpUrlencoded(method, url, data, header)
	}
	return
}
