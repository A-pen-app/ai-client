package models

// InputFile is one attachment with its type stated by the caller. Sniffing is
// not enough — Go's table has no HEIC, so an iPhone photo reads as nothing.
type InputFile struct {
	URL      string
	MimeType string
}

type AIChatMessage struct {
	SystemPrompt string
	Text         string
	Files        []InputFile
	// Deprecated: use Files, which states the type instead of sniffing it.
	ImageUrls []string
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
