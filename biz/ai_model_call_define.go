package biz

// AI 连接协议类型
const (
	AIProtocolDify  = "dify"
	AIProtocolOpenAI = "openai"
)

type AIConnect struct {
	BaseUrl  string `json:"baseUrl"`
	ApiKey   string `json:"apiKey"`
	Timeout  int64  `json:"timeout"`
	Protocol string `json:"protocol"` // "" 或 "dify" = Dify（缺省向后兼容）；"openai" = OpenAI 兼容协议
	Model    string `json:"model"`    // 仅 openai 协议必需，如 "deepseek-chat" / "gpt-4o"
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

// OpenAI 兼容协议（chat/completions）请求/响应结构体
type OpenAIMessage struct {
	Role    string `json:"role"`    // "system" | "user" | "assistant"
	Content string `json:"content"`
}

type OpenAIChatRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIChatResponse struct {
	Id      string        `json:"id"`
	Choices []OpenAIChoice `json:"choices"`
	Error   *OpenAIError  `json:"error,omitempty"`
}

type OpenAIChoice struct {
	Message OpenAIMessage `json:"message"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}
