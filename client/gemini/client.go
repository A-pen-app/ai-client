package gemini

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/A-pen-app/ai-clients/models"
	"github.com/A-pen-app/ai-clients/store"
	"github.com/A-pen-app/ai-clients/util"
	"github.com/google/generative-ai-go/genai"
)

// Client wraps Gemini client and implements the AIClient interface
type Client struct {
	client       *genai.Client
	defaultModel string
}

// NewClient creates a new Gemini client
func NewClient(client *genai.Client, model string) store.Client {
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &Client{
		client:       client,
		defaultModel: model,
	}
}

func (c *Client) Generate(ctx context.Context, message models.Message, opts models.Options) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("gemini client is not initialized")
	}

	modelName := c.defaultModel
	if opts.Model != "" {
		modelName = opts.Model
	}

	model := c.client.GenerativeModel(modelName)

	if message.SystemPrompt != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(message.SystemPrompt)},
		}
	}

	if opts.ResponseFormat == models.ResponseFormatJSON {
		model.ResponseMIMEType = "application/json"
	}

	if opts.MaxTokens > 0 {
		model.SetMaxOutputTokens(int32(opts.MaxTokens))
	}

	var promptParts []genai.Part

	if message.Text != "" {
		promptParts = append(promptParts, genai.Text(message.Text))
	}

	for _, url := range message.ImageUrls {
		imageData, err := util.DownloadImage(ctx, url)
		if err != nil {
			return "", fmt.Errorf("failed to download image: %w", err)
		}
		mimeType := http.DetectContentType(imageData)
		promptParts = append(promptParts, genai.ImageData(mimeType, imageData))
	}

	resp, err := model.GenerateContent(ctx, promptParts...)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("empty content in Gemini response")
	}

	var resultText strings.Builder
	for _, part := range candidate.Content.Parts {
		if text, ok := part.(genai.Text); ok {
			resultText.WriteString(string(text))
		}
	}
	return resultText.String(), nil
}
