package generate

type chatResponse struct {
	Text         string `json:"text"`
	GenerationID string `json:"generation_id"`
	ResponseID   string `json:"response_id"`
	FinishReason string `json:"finish_reason"`
	ChatHistory  []struct {
		Message string `json:"message"`
		Role    string `json:"role"`
	} `json:"chat_history"`
	Meta struct {
		APIVersion struct {
			Version string `json:"version"`
		} `json:"api_version"`
		BilledUnits struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"billed_units"`
		Tokens struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"tokens"`
		CachedTokens int `json:"cached_tokens"`
	} `json:"meta"`
}
