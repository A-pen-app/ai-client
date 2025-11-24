package models

type AIChatMessage struct {
	SystemPrompt string
	Text         string
	ImageUrls    []string
}

type ResponseFormat string

const (
	ResponseFormatJSON ResponseFormat = "json"
	ResponseFormatText ResponseFormat = "text"
)

type AIClientOptions struct {
	MaxTokens      int64
	Model          string
	ResponseFormat ResponseFormat
}
