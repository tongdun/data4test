package biz

// Dataset 数据集结构体
type Dataset struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Document 文档结构体
type Document struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DatasetsResponse 数据集响应结构体
type DatasetsResponse struct {
	Data  []Dataset `json:"data"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}

// DocumentsResponse 文档响应结构体
type DocumentsResponse struct {
	Data []Document `json:"data"`
}

// fileDocument 文件文档关联信息
type fileDocument struct {
	documentID string
	datasetID  string
}

type DataSetConnect struct {
	DataSetName string `json:"dataSetName"`
	AIConnect
}
