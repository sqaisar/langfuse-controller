package langfuse

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add other fields as needed
}

type APIKey struct {
	ID        string `json:"id"`
	PublicKey string `json:"publicKey"`
	SecretKey string `json:"secretKey"`
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
}

type CreateProjectRequest struct {
	Name string `json:"name"`
}

type CreateAPIKeyRequest struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
}

type Model struct {
	ID              string  `json:"id"`
	ModelName       string  `json:"modelName"`
	MatchPattern    string  `json:"matchPattern"`
	StartDate       string  `json:"startDate,omitempty"`
	Unit            string  `json:"unit"`
	InputPrice      float64 `json:"inputPrice,omitempty"`
	OutputPrice     float64 `json:"outputPrice,omitempty"`
	TotalPrice      float64 `json:"totalPrice,omitempty"`
	TokenizerId     string  `json:"tokenizerId,omitempty"`
	TokenizerConfig string  `json:"tokenizerConfig,omitempty"`
}

// Add other types for LlmConnection, Prompt, ScoreConfig
