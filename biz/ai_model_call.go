package biz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// 调用 Dify 接口的函数
func CallDifyChat(aiConnect AIConnect, query string, userId, cId, fileId, fileType string) (reply, conversionId string, err error) {
	// 构造请求体
	var reqBody ChatRequest
	if len(fileId) > 0 {
		reqBody = ChatRequest{
			Inputs:         map[string]string{},
			Query:          query,
			User:           userId,
			ConversationId: cId,
			ResponseMode:   "blocking",
			Files: []ChatFile{{
				Type:           fileType,
				TransferMethod: "local_file",
				Url:            fileId,
			}},
		}
	} else {
		reqBody = ChatRequest{
			Inputs:         map[string]string{},
			Query:          query,
			User:           userId,
			ConversationId: cId,
			ResponseMode:   "blocking",
		}
	}

	data := make(map[string]interface{})
	jsonBytes, _ := json.Marshal(reqBody)
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return "", "", err
	}

	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/chat-messages", aiConnect.BaseUrl)

	respBody, err := RunHttpJson("POST", url, aiConnect.Timeout, data, header)

	// 解析 JSON
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", "", fmt.Errorf("响应解析失败: %v, 原始响应: %s", err, string(respBody))
	}

	if len(cId) == 0 {
		conversionId = chatResp.ConversationId
	}

	reply = chatResp.Answer
	return chatResp.Answer, conversionId, nil
}

func CallDify2UploadFile(aiConnect AIConnect, filePath, userCode string) (fileID, fileType string, err error) {
	// 创建表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return fileID, fileType, err
	}
	io.Copy(part, file)

	// 添加其他表单字段
	tFileType, err := DetectFileType(filePath)
	if err != nil {
		Logger.Error("%v", err)
		return fileID, fileType, err
	}

	documentType := []string{"TXT", "MD", "MARKDOWN", "PDF", "HTML", "XLSX", "XLS", "DOCX", "CSV", "EML", "MSG", "PPTX", "PPT", "XML", "EPUB"}
	imageType := []string{"JPG", "JPEG", "PNG", "GIF", "WEBP", "SVG"}

	for _, k := range imageType {
		if k == tFileType {
			fileType = "image"
			break
		}
	}

	if len(fileType) == 0 {
		for _, k := range documentType {
			if k == tFileType {
				fileType = "document"
				break
			}
		}
	}

	if len(fileType) == 0 {
		err = fmt.Errorf("当前文件类型[%s]不支持解析，请更换文档载体~", tFileType)
		return
	}

	_ = writer.WriteField("user", userCode)
	_ = writer.WriteField("type", fileType) // 根据文件类型调整（如PDF、XLSX）

	writer.Close()

	// 创建请求
	url := fmt.Sprintf("%s/files/upload", aiConnect.BaseUrl)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fileID, fileType, err
	}
	req.Header.Set("Authorization", "Bearer "+aiConnect.ApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fileID, fileType, err
	}
	defer resp.Body.Close()

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fileID, fileType, err
	}

	if fileID, ok := result["id"].(string); ok {
		return fileID, fileType, nil
	}
	return fileID, fileType, fmt.Errorf("获取fileID失败")
}

func ConnetAIModel(query, appendQuery, uploadFilePath, aiPlatform, userCode string) (replyList []string, err error) {
	aiConnect, err := GetAIModelConnectInfo(aiPlatform)
	if err != nil {
		return
	}

	var fileId, fileType string
	if len(uploadFilePath) > 0 {
		fileId, fileType, err = CallDify2UploadFile(aiConnect, uploadFilePath, userCode)
		if err != nil {
			Logger.Error("上传文件到大模型失败: %s", err)
			return
		}
	}

	reply1, cId, err := CallDifyChat(aiConnect, query, userCode, "", fileId, fileType)
	if err != nil {
		Logger.Error("调用%s失败: %v", aiPlatform, err)
		return
	}
	replyList = append(replyList, reply1)

	appendQueryList := strings.Split(appendQuery, "|**|")
	for _, item := range appendQueryList {
		if item == " " {
			continue
		}
		answer, _, errTmp := CallDifyChat(aiConnect, item, userCode, cId, fileId, fileType)
		if errTmp != nil {
			Logger.Error("调用%s失败: %v", aiPlatform, err)
			err = errTmp
			return
		}
		replyList = append(replyList, answer)
	}

	return
}

func (input ImportCommon) ConnectModel2GetMessage() (replyList []string, err error) {
	aiConnect, err := GetAIModelConnectInfo(input.CreatePlatform)
	if err != nil {
		return
	}

	replyList, err = aiConnect.CallModel2GetMessage(input.CreateUser, input.ConversationId)
	if err != nil {
		Logger.Error("调用%s失败: %v", input.CreatePlatform, err)
		return
	}

	return
}

func (aiConnect AIConnect) CallModel2GetMessage(userCode, conversationID string) (replyList []string, err error) {
	url := fmt.Sprintf("%s/messages", aiConnect.BaseUrl)
	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	data := make(map[string]interface{})
	data["user"] = userCode
	data["conversation_id"] = conversationID

	replyList, err = aiConnect.CallModel2Get("GET", url, data, header)
	if err != nil {
		Logger.Debug("url: %s", url)
		Logger.Debug("header: %s", header)
		Logger.Debug("data: %s", data)
		Logger.Error("%s", err)

	}

	return
}

func (aiConnect AIConnect) CallModel2GetMessageOld(userCode, conversationID string) (replyList []string, err error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/messages", aiConnect.BaseUrl)
	req, _ := http.NewRequest("GET", url, nil)

	// 设置查询参数
	q := req.URL.Query()
	q.Add("user", userCode)
	q.Add("conversation_id", conversationID)
	//q.Add("limit", "20") // 单次最大消息数
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", "Bearer "+aiConnect.ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var cResp ConversationsResponse

	if errTmp := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		Logger.Error("%s", errTmp)
		err = errTmp
		return
	}

	if len(cResp.Data) > 0 {
		for _, item := range cResp.Data {
			replyList = append(replyList, item.Answer)
		}
	}

	return
}

func (aiConnect AIConnect) CallModel2Get(method, url string, data, header map[string]interface{}) (replyList []string, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	// 设置查询参数
	q := req.URL.Query()
	for k, v := range data {
		q.Add(k, v.(string))
	}

	req.URL.RawQuery = q.Encode()
	for k, v := range header {
		req.Header.Set(k, v.(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var cResp ConversationsResponse

	if errTmp := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		Logger.Error("%s", errTmp)
		err = errTmp
		return
	}

	if len(cResp.Data) > 0 {
		for _, item := range cResp.Data {
			replyList = append(replyList, item.Answer)
		}
	}

	return
}
