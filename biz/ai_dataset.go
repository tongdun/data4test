package biz

import (
	"encoding/json"
	"fmt"
)

// CreateDataset 创建一个新的数据集
func (aiConnect DataSetConnect) CreateDataset(name string) (string, error) {
	data := make(map[string]interface{})
	data["name"] = name
	data["permission"] = "all_team_members"
	data["description"] = "盾测线上环境使用"
	data["indexing_technique"] = "high_quality"

	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/datasets", aiConnect.BaseUrl)

	respBody, err := RunHttpJson("POST", url, aiConnect.Timeout, data, header)

	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(respBody), &result)
	if err != nil {
		Logger.Error("err: %s", err)
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("获取数据集ID失败")
}

// DeleteDataset 删除一个已有知识库
func (aiConnect DataSetConnect) DeleteDataset(dataSetId string) {
	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/datasets/%s", aiConnect.BaseUrl, dataSetId)

	_, _ = RunHttpJson("DELETE", url, aiConnect.Timeout, nil, header)

}

// CreateDocument 在指定数据集中创建新文档
func (aiConnect DataSetConnect) CreateDocument(datasetID, name, text string) (string, error) {
	data := map[string]interface{}{
		"name":               name,
		"text":               text,
		"indexing_technique": "high_quality",
		"process_rule": map[string]string{
			"mode": "automatic",
		},
	}

	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/datasets/%s/document/create-by-text", aiConnect.BaseUrl, datasetID)

	respBody, err := RunHttpJson("POST", url, aiConnect.Timeout, data, header)

	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(respBody), &result)
	if err != nil {
		Logger.Error("err: %s", err)
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}
	return "", nil
}

// UpdateDocument 更新指定数据集中的文档
func (aiConnect DataSetConnect) UpdateDocument(datasetID, documentID, name, text string) error {
	//data := make(map[string]interface{})
	data := map[string]interface{}{
		"name": name,
		"text": text,
		"process_rule": map[string]string{
			"mode": "automatic",
		},
	}

	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/datasets/%s/documents/%s/update-by-text", aiConnect.BaseUrl, datasetID, documentID)

	respBody, err := RunHttpJson("POST", url, aiConnect.Timeout, data, header)

	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(respBody), &result)
	if err != nil {
		Logger.Error("err: %s", err)
	}

	return nil
}

func (aiConnect DataSetConnect) GetDatasetByName(name string) (dataSet Dataset, err error) {
	data := make(map[string]interface{})
	data["page"] = 1
	data["limit"] = 30

	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/datasets", aiConnect.BaseUrl)

	respBody, err := RunHttpJson("GET", url, aiConnect.Timeout, data, header)

	var result DatasetsResponse
	err = json.Unmarshal([]byte(respBody), &result)
	if err != nil {
		Logger.Error("err: %s", err)
	}

	if len(result.Data) == 0 {
		err = fmt.Errorf("未获取到知识库[%s], 请核对 ~", name)
		return
	}

	for index, item := range result.Data {
		if item.Name == name {
			dataSet = item
			break
		}

		if index == len(result.Data)-1 {
			err = fmt.Errorf("未获取到知识库[%s], 请核对 ~", name)
		}
	}

	return
}

func (aiConnect DataSetConnect) GetDocumentByName(datasetID, fileName string) (doc Document, err error) {
	data := make(map[string]interface{})
	data["keyword"] = fileName

	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/datasets/%s/documents", aiConnect.BaseUrl, datasetID)
	var result DocumentsResponse
	respBody, err := RunHttpJson("GET", url, aiConnect.Timeout, data, header)
	if err != nil {
		Logger.Error("err: %s", err)
	}
	err = json.Unmarshal([]byte(respBody), &result)
	if err != nil {
		Logger.Error("err: %s", err)
	}

	if len(result.Data) == 0 {
		err = fmt.Errorf("未获取到文档[%s], 请核对 ~", fileName)
		return
	}

	doc = result.Data[0]

	return
}

// DeleteDocument 删除指定数据集中的文档
func (aiConnect DataSetConnect) DeleteDocument(datasetID, documentID string) error {
	header := make(map[string]interface{})
	header["Authorization"] = fmt.Sprintf("Bearer %s", aiConnect.ApiKey)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/datasets/%s/documents/%s", aiConnect.BaseUrl, datasetID, documentID)

	_, err := RunHttpJson("DELETE", url, aiConnect.Timeout, nil, header)
	if err != nil {
		Logger.Error("err: %s", err)
	}

	return err

}
