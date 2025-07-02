package biz

type AIConnect struct {
	BaseUrl string `json:"baseUrl"`
	ApiKey  string `json:"apiKey"`
	Timeout int64  `json:"timeout"`
}

// 请求体结构
type ChatRequest struct {
	Inputs         map[string]string `json:"inputs"`
	Query          string            `json:"query"`
	User           string            `json:"user"`
	ResponseMode   string            `json:"response_mode"`
	ConversationId string            `json:"conversation_id"`
	Files          []ChatFile        `json:"files"`
}

type ChatFile struct {
	Type           string `json:"type"`
	TransferMethod string `json:"transfer_method"`
	Url            string `json:"url"`
}

// 响应结构（简化）
//type ChatResponse struct {
//	Message struct {
//		Content string `json:"content"`
//	} `json:"message"`
//}

// 响应结构（全量）
type ChatResponse struct {
	Event          string         `json:"event"`
	MessageId      string         `json:"message_id"`
	ConversationId string         `json:"conversation_id"`
	Mode           string         `json:"mode"`
	Answer         string         `json:"answer"`
	Metadata       MetaDataDetail `json:"metadata"`
	CreatedAt      int            `json:"created_at"`
}

type MetaDataDetail struct {
	Usage              UsageDetail         `json:"usage"`
	RetrieverResources []RetrieverResource `json:"retriever_resources"`
}

type UsageDetail struct {
	PromptTokens        int     `json:"prompt_tokens"`
	PromptUnitPrice     string  `json:"prompt_unit_price"`
	PromptPriceUnit     string  `json:"prompt_price_unit"`
	PromptPrice         string  `json:"prompt_price"`
	CompletionTokens    int     `json:"completion_tokens"`
	CompletionUnitPrice string  `json:"completion_unit_price"`
	CompletionPriceUnit string  `json:"completion_price_unit"`
	CompletionPrice     string  `json:"completion_price"`
	TotalTokens         int     `json:"total_tokens"`
	TotalPrice          string  `json:"total_price"`
	Currency            string  `json:"currency"`
	Latency             float64 `json:"latency"`
}

type ConversationsResponse struct {
	Limit   int            `json:"limit"`
	HasMore bool           `json:"has_more"`
	Data    []ResponseData `json:"data"`
}

type ResponseData struct {
	Id             bool     `json:"id"`
	ConversationId string   `json:"conversation_id"`
	Inputs         string   `json:"inputs"`
	Query          string   `json:"query"`
	Answer         string   `json:"answer"`
	MessageFiles   []string `json:"message_files"`
	//Feedback           string              `json:"feedback"`
	RetrieverResources []RetrieverResource `json:"retriever_resources"`
	CreatedAt          int                 `json:"created_at"`
}

type RetrieverResource struct {
	Position     int     `json:"position"`
	DatasetId    string  `json:"dataset_id"`
	DatasetName  string  `json:"dataset_name"`
	DocumentId   string  `json:"document_id"`
	DocumentName string  `json:"document_name"`
	SegmentId    string  `json:"segment_id"`
	Score        float64 `json:"score"`
	Content      string  `json:"content"`
}
