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
	Logger.Info("连接大模型 URL: %s", url)

	_, respBody, err := RunHttpJson("POST", url, aiConnect.Timeout, data, header, nil)
	if err != nil {
		// 如果请求遇错，直接退出，不再进行后续会话
		return "", "", fmt.Errorf(T("error.response_parse_failed"), err, string(respBody))

	}

	// 解析 JSON
	var chatResp ChatResponse
	if err = json.Unmarshal(respBody, &chatResp); err != nil {
		return "", "", fmt.Errorf(T("error.response_parse_failed"), err, string(respBody))
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
		err = fmt.Errorf(T("error.file_type_not_supported"), tFileType)
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
	return fileID, fileType, fmt.Errorf(T("error.get_fileid_failed"))
}

// buildOpenAIChatURL 构造 OpenAI 兼容 chat/completions 地址，兼容 baseUrl 是否已带 /v1 前缀
func buildOpenAIChatURL(baseUrl string) string {
	base := strings.TrimRight(baseUrl, "/")
	if strings.HasSuffix(base, "/v1") {
		return base + "/chat/completions"
	}
	return base + "/v1/chat/completions"
}

// CallOpenAIChat 调用 OpenAI 兼容协议的 chat/completions 接口，返回首个回复内容
func CallOpenAIChat(aiConnect AIConnect, messages []OpenAIMessage) (reply string, err error) {
	reqBody := OpenAIChatRequest{Model: aiConnect.Model, Messages: messages}

	data := make(map[string]interface{})
	jsonBytes, _ := json.Marshal(reqBody)
	if err = json.Unmarshal(jsonBytes, &data); err != nil {
		return "", err
	}

	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"

	url := buildOpenAIChatURL(aiConnect.BaseUrl)
	Logger.Info("连接大模型 URL: %s", url)

	_, respBody, err := RunHttpJson("POST", url, aiConnect.Timeout, data, header, nil)
	if err != nil {
		return "", fmt.Errorf(T("error.response_parse_failed"), err, string(respBody))
	}

	var chatResp OpenAIChatResponse
	if err = json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf(T("error.response_parse_failed"), err, string(respBody))
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("code=%s type=%s: %s", chatResp.Error.Code, chatResp.Error.Type, chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf(T("error.no_reply_info"))
	}

	return chatResp.Choices[0].Message.Content, nil
}

// CallOpenAIModel 以 OpenAI 兼容协议对话：文本优先，不支持文件上传
func CallOpenAIModel(query, appendQuery, uploadFilePath string, aiConnect AIConnect) (replyList []string, err error) {
	if len(uploadFilePath) > 0 {
		return nil, fmt.Errorf(T("error.openai_no_file_upload"))
	}

	messages := []OpenAIMessage{{Role: "user", Content: query}}
	reply, err := CallOpenAIChat(aiConnect, messages)
	if err != nil {
		return nil, err
	}
	replyList = append(replyList, reply)

	for _, item := range strings.Split(appendQuery, "|**|") {
		if item == " " {
			continue
		}
		messages = append(messages, OpenAIMessage{Role: "assistant", Content: reply})
		messages = append(messages, OpenAIMessage{Role: "user", Content: item})
		reply, err = CallOpenAIChat(aiConnect, messages)
		if err != nil {
			return nil, err
		}
		replyList = append(replyList, reply)
	}

	return
}

func ConnetAIModel(query, appendQuery, uploadFilePath, aiPlatform, userCode string) (replyList []string, err error) {
	aiConnect, err := GetAIModelConnectInfo(aiPlatform)
	if err != nil {
		return
	}
	return ConnetAIModelWithConnect(query, appendQuery, uploadFilePath, userCode, aiPlatform, aiConnect)
}

// ConnetAIModelWithConnect 用已解析好的连接信息调用大模型（openai/dify 分流）。
// 独立于 ConnetAIModel，供「实时覆盖 baseUrl/apiKey」的场景复用：覆盖值已在 aiConnect 中，全程不再查库。
func ConnetAIModelWithConnect(query, appendQuery, uploadFilePath, userCode, platform string, aiConnect AIConnect) (replyList []string, err error) {
	if aiConnect.Protocol == AIProtocolOpenAI {
		return CallOpenAIModel(query, appendQuery, uploadFilePath, aiConnect)
	}

	var fileId, fileType string
	if len(uploadFilePath) > 0 {
		fileId, fileType, err = CallDify2UploadFile(aiConnect, uploadFilePath, userCode)
		if err != nil {
			Logger.Error(T("error.upload_file_failed"), err)
			return
		}
	}

	reply1, cId, err := CallDifyChat(aiConnect, query, userCode, "", fileId, fileType)
	if err != nil {
		Logger.Error(T("error.call_model_failed"), platform, err)
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
			Logger.Error(T("error.call_model_failed"), platform, err)
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
		Logger.Error(T("error.call_model_failed"), input.CreatePlatform, err)
		return
	}

	return
}

func (aiConnect AIConnect) CallModel2GetMessage(userCode, conversationID string) (replyList []string, err error) {
	if aiConnect.Protocol == AIProtocolOpenAI {
		return nil, fmt.Errorf(T("error.openai_no_history"))
	}
	url := fmt.Sprintf("%s/messages", aiConnect.BaseUrl)
	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	data := make(map[string]interface{})
	data["user"] = userCode
	data["conversation_id"] = conversationID

	replyList, err = aiConnect.CallModel2Get("GET", url, data, header)
	if err != nil {
		Logger.Debug("url: %s", url)
		Logger.Debug("data: %s", data)
		Logger.Error("%s", err)

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
