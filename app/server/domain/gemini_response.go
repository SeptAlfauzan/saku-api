package domain

// /RAW GEMINI RESPONSE
type GeminiResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Usage       Usage  `json:"usage"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	ServiceTier string `json:"service_tier"`
	Steps       []Step `json:"steps"`
	Object      string `json:"object"`
	Model       string `json:"model"`
}

type Usage struct {
	TotalTokens                int               `json:"total_tokens"`
	TotalInputTokens           int               `json:"total_input_tokens"`
	InputTokensByModality      []TokenByModality `json:"input_tokens_by_modality"`
	TotalCachedTokens          int               `json:"total_cached_tokens"`
	TotalOutputTokens          int               `json:"total_output_tokens"`
	TotalToolUseTokens         int               `json:"total_tool_use_tokens"`
	TotalThoughtTokens         int               `json:"total_thought_tokens"`
	RawPromptToken             int               `json:"raw_prompt_token"`
	ModelInvocationTokenCounts []ModelInvocation `json:"model_invocation_token_counts"`
}

type TokenByModality struct {
	Modality string `json:"modality"`
	Tokens   int    `json:"tokens"`
}

type ModelInvocation struct {
	PromptTokensDetails     []TokenByModality `json:"prompt_tokens_details"`
	CandidatesTokensDetails []TokenByModality `json:"candidates_tokens_details"`
}

type Step struct {
	Signature string    `json:"signature,omitempty"`
	Type      string    `json:"type"`
	Content   []Content `json:"content,omitempty"`
}

type Content struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type"`
}
