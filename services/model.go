package services

// 定义模型结构
type Endpoint struct {
	CheckHealthURL    string `json:"check_health_url"`
	EnableCheckHealth bool   `json:"enable_check_health"`
	EnableTestflight  bool   `json:"enable_testflight"`
	EndpointID        string `json:"endpoint_id"`
	ModelName         string `json:"model_name"`
	ModelNameOverride string `json:"model_name_override"`
	Provider          string `json:"provider"`
	Status            string `json:"status"`
	URL               string `json:"url"`
	Weight            int    `json:"weight"`
	SupportedApiFlag  string `json:"supported_api_flag"`
}

type Model struct {
	ContextSize                int      `json:"context_size"`
	DefaultStops               []string `json:"default_stops"`
	Description                string   `json:"description"`
	InputTokenPricePerMillion  int      `json:"input_token_price_per_million"`
	IsPrivate                  bool     `json:"is_private"`
	ModelName                  string   `json:"model_name"`
	OutputTokenPricePerMillion int      `json:"output_token_price_per_million"`
	Rank                       int      `json:"rank"`
	Status                     string   `json:"status"`
	Tags                       []string `json:"tags"`
	UserEmail                  string   `json:"user_email"`
}

type LLMModelInfo struct {
	Description      string     `json:"description"`
	Endpoints        []Endpoint `json:"endpoints"`
	InputTokenPrice  int        `json:"input_token_price"`
	MaxTokens        int        `json:"max_tokens"`
	Model            Model      `json:"model"`
	ModelName        string     `json:"model_name"`
	Note             string     `json:"note"`
	OutputTokenPrice int        `json:"output_token_price"`
	Private          bool       `json:"private"`
	Rank             int        `json:"rank"`
	Starred          bool       `json:"starred"`
	Status           string     `json:"status"`
	Tags             []string   `json:"tags"`
}

type ModelResponse struct {
	Models []LLMModelInfo `json:"models"`
}
