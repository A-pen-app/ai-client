package gemini

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/A-pen-app/ai-client/models"
	"github.com/A-pen-app/ai-client/store"
	"github.com/A-pen-app/ai-client/util"
	"google.golang.org/genai"
)

// Client wraps Gemini API client and implements the AIClient interface
type Client struct {
	client       *genai.Client
	defaultModel string
}

// NewClient creates a new Gemini API client
func NewClient(projectID string, location string, model string) (store.AIClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  projectID,
		Location: location,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &Client{
		client:       client,
		defaultModel: model,
	}, nil
}

func (c *Client) Generate(ctx context.Context, message models.AIChatMessage, opts models.AIClientOptions) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("gemini client is not initialized")
	}

	modelName := c.defaultModel
	if opts.Model != "" {
		modelName = opts.Model
	}

	var contentParts []*genai.Part

	if message.Text != "" {
		contentParts = append(contentParts, genai.NewPartFromText(message.Text))
	}

	for _, url := range message.ImageUrls {
		imageData, err := util.DownloadImage(ctx, url)
		if err != nil {
			return "", fmt.Errorf("failed to download image: %w", err)
		}
		mimeType := http.DetectContentType(imageData)
		contentParts = append(contentParts, genai.NewPartFromBytes(imageData, mimeType))
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(contentParts, genai.RoleUser),
	}

	// Build generation config
	config := &genai.GenerateContentConfig{}

	if message.SystemPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(message.SystemPrompt, genai.RoleUser)
	}

	if opts.ResponseFormat == models.ResponseFormatJSON {
		config.ResponseMIMEType = "application/json"
	}

	if opts.MaxTokens > 0 {
		config.MaxOutputTokens = int32(opts.MaxTokens)
	}

	resp, err := c.client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
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
		if part.Text != "" {
			resultText.WriteString(part.Text)
		}
	}

	return resultText.String(), nil
}
