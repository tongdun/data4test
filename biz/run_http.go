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
	"time"
)

func RunHttpUrlencodedRawResp(method, url string, data map[string]interface{}, acceptHeader, responseHeader map[string]interface{}) (res io.Reader, err error) {
	var req *http.Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
			strK := Interface2Str(v)
			if len(strK) == 0 { // 为GET请求时，入参值为空时，直接过滤
				continue
			}
			tmpData[k] = []string{strK}
			uri.RawQuery = tmpData.Encode()
		}

		req, err = http.NewRequest(methodUpper, uri.String(), nil)
	} else {
		dataPayload := netUrl.Values{}
		for k, v := range data {
			strValue := Interface2Str(v)

			dataPayload.Add(k, strValue)
		}

		payload := strings.NewReader(dataPayload.Encode())
		req, err = http.NewRequest(methodUpper, url, payload)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	for k, v := range acceptHeader {
		vStr := Interface2Str(v)
		req.Header.Add(k, vStr)
	}

	resp, err := client.Do(req)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	defer resp.Body.Close()

	return resp.Body, err
}

//func GetHttpHandle(timeout int64, data map[string]interface{}, header map[string]interface{}) (req *http.Request, client *http.Client) {
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client = &http.Client{Transport: tr, Timeout: time.Duration(timeout) * time.Second}
//	payload := &bytes.Buffer{}
//	writer := multipart.NewWriter(payload)
//	for k, v := range header {
//		if k == "Content-Type" {
//			continue
//		}
//		valueStr := Interface2Str(v)
//		req.Header.Add(k, valueStr)
//	}
//
//	return
//}

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
		strValue := Interface2Str(v)
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
					filePath = fmt.Sprintf("%s/%s", DownloadBasePath, strValue)
					_, errTmp := os.Stat(filePath)
					if errTmp != nil {
						err = fmt.Errorf("[%s]在[文件-上传文件]和[文件-下载文件]目录下均不存在，请核对", filePath)
						Logger.Error("%s", err)
						return
					}
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
				strK := Interface2Str(v)
				if len(strK) == 0 { // 为GET请求时，入参值为空时，直接过滤
					continue
				}
				tmpData[k] = []string{strK}
				uri.RawQuery = tmpData.Encode()
			}

			req, err = http.NewRequest(methodUpper, uri.String(), nil)

		} else {
			req, err = http.NewRequest(methodUpper, url, nil)

		}
	} else {
		req, err = http.NewRequest(methodUpper, url, payload)
	}

	if err != nil {
		Logger.Debug("method: %s, url: %s, data: %v", methodUpper, url, data)
		Logger.Error("%s", err)
		return
	}

	for k, v := range header {
		if k == "Content-Type" {
			continue
		}
		valueStr := Interface2Str(v)
		req.Header.Add(k, valueStr)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		Logger.Debug("req: %s", req)
		Logger.Debug("resp: %s", resp)
		Logger.Error("%s", err)
		return
	}

	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Debug("resBody: %s", string(resBody))
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

func RunHttpFormDataOld(method, url string, data map[string]interface{}, header map[string]interface{}) (res []byte, err error) {
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
		strValue := Interface2Str(v)
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
					filePath = fmt.Sprintf("%s/%s", DownloadBasePath, strValue)
					_, errTmp := os.Stat(filePath)
					if errTmp != nil {
						err = fmt.Errorf("[%s]在[文件-上传文件]和[文件-下载文件]目录下均不存在，请核对", filePath)
						Logger.Error("%s", err)
						return
					}
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
				strK := Interface2Str(v)
				if len(strK) == 0 { // 为GET请求时，入参值为空时，直接过滤
					continue
				}
				tmpData[k] = []string{strK}
				uri.RawQuery = tmpData.Encode()
			}

			req, err = http.NewRequest(methodUpper, uri.String(), nil)

		} else {
			req, err = http.NewRequest(methodUpper, url, nil)

		}
	} else {
		req, err = http.NewRequest(methodUpper, url, payload)
	}

	if err != nil {
		Logger.Debug("method: %s, url: %s, data: %v", methodUpper, url, data)
		Logger.Error("%s", err)
		return
	}

	for k, v := range header {
		if k == "Content-Type" {
			continue
		}
		valueStr := Interface2Str(v)
		req.Header.Add(k, valueStr)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		Logger.Debug("req: %s", req)
		Logger.Debug("resp: %s", resp)
		Logger.Error("%s", err)
		return
	}

	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Debug("resBody: %s", string(resBody))
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

func RunHttpUrlencoded(method, url string, data map[string]interface{}, acceptHeader, responseHeader map[string]interface{}) (res []byte, err error) {
	var req *http.Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
			strK := Interface2Str(v)
			if len(strK) == 0 { // 为GET请求时，入参值为空时，直接过滤
				continue
			}
			tmpData[k] = []string{strK}
			uri.RawQuery = tmpData.Encode()
		}

		req, err = http.NewRequest(methodUpper, uri.String(), nil)
	} else {
		dataPayload := netUrl.Values{}
		for k, v := range data {
			strValue := Interface2Str(v)

			dataPayload.Add(k, strValue)
		}

		payload := strings.NewReader(dataPayload.Encode())
		req, err = http.NewRequest(methodUpper, url, payload)
	}

	if err != nil {
		Logger.Error("%s", err)
		return
	}

	for k, v := range acceptHeader {
		vStr := Interface2Str(v)
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

	downloadRawInfo := resp.Header.Get("Content-Disposition")
	downloadInfo, _ := netUrl.QueryUnescape(downloadRawInfo)

	// 返回500的是否需要拦截
	//if resp.StatusCode != 200 || resp.StatusCode != 500 {
	if resp.StatusCode != 200 {
		err = fmt.Errorf("请求失败，返回码: %d, 返回信息: %s", resp.StatusCode, string(resBody))
	}

	var downloadFileName, downloadFilePath string
	if len(downloadInfo) > 0 {
		tmps := strings.Split(downloadInfo, "=")
		if len(tmps) > 1 {
			downloadFileName = tmps[1]
			if strings.Contains(downloadFileName, "\"") {
				downloadFileName = strings.Replace(downloadFileName, "\"", "", -1)
			}
			if strings.Contains(downloadFileName, "'") {
				downloadFileName = strings.Replace(downloadFileName, "'", "", -1)
			}

			downloadFilePath = fmt.Sprintf("%s/%s", DownloadBasePath, downloadFileName)
		}
	} else {
		for k, v := range responseHeader {
			vStr := Interface2Str(v)
			if k == "Content-Disposition" {
				tmps := strings.Split(vStr, "=")
				if len(tmps) > 1 {
					downloadFileName = tmps[1]
					downloadFilePath = fmt.Sprintf("%s/%s", DownloadBasePath, downloadFileName)
				}
				break
			}
		}
	}

	if len(downloadFilePath) > 0 {
		fh, errTmp := os.Create(downloadFilePath)
		if errTmp != nil {
			Logger.Error("%v", errTmp)
			if err != nil {
				err = fmt.Errorf("%s;%s", err, errTmp)
			} else {
				err = errTmp
			}
			return []byte(downloadFileName), err
		}
		defer fh.Close()

		_, errTmp = fh.Write(resBody)
		if errTmp != nil {
			Logger.Error("%v", errTmp)
			if err != nil {
				err = fmt.Errorf("%s;%s", err, errTmp)
			} else {
				err = errTmp
			}
			return resBody, err
		}
		return []byte(downloadFileName), err
	}

	return resBody, err
}

func RunHttpJson(method, url string, timeout int64, data map[string]interface{}, header map[string]interface{}) (res []byte, err error) {
	var req *http.Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if _, ok := header["Content-Type"]; !ok {
		err = fmt.Errorf("header 未正常定义，请核对")
		return
	}

	methodUpper := strings.ToUpper(method)
	client := &http.Client{Transport: tr, Timeout: time.Duration(timeout) * time.Second}

	if methodUpper == "GET" && len(data) > 0 {
		uri, err1 := netUrl.Parse(url)
		if err1 != nil {
			err = err1
			Logger.Error("%s", err)
			return
		}

		tmpData := make(netUrl.Values)
		for k, v := range data {
			strK := Interface2Str(v)
			if len(strK) == 0 { // 为GET请求时，入参值为空时，直接过滤
				continue
			}
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
		vStr := Interface2Str(v)
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
		vStr := Interface2Str(v)
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

func RunHttp(method, url string, data map[string]interface{}, acceptHeader, responseHeader map[string]interface{}) (res []byte, err error) {
	contentTypeRaw := Interface2Str(acceptHeader["Content-Type"])
	var contentType string

	if strings.Contains(contentTypeRaw, "application/x-www-form-urlencoded") {
		contentType = "application/x-www-form-urlencoded"
	} else if strings.Contains(contentTypeRaw, "multipart/form-data") {
		contentType = "multipart/form-data"
	} else if strings.Contains(contentTypeRaw, "form-data") {
		contentType = "form-data"
	} else if strings.Contains(contentTypeRaw, "application/json") {
		contentType = "application/json"
	} else {
		contentType = "application/x-www-form-urlencoded"
	}

	switch contentType {
	case "application/x-www-form-urlencoded":
		res, err = RunHttpUrlencoded(method, url, data, acceptHeader, responseHeader)
	case "form-data", "multipart/form-data":
		res, err = RunHttpFormData(method, url, data, acceptHeader)
	case "application/json":
		res, err = RunHttpJson(method, url, 180, data, acceptHeader)
	default:
		res, err = RunHttpUrlencoded(method, url, data, acceptHeader, responseHeader)
	}

	return
}
