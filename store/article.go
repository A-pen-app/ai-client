package store

import (
	"context"
	"fmt"

	"github.com/A-pen-app/ai-clients/models"
)

type ArticleConfig struct {
	MaxToken int64
}

type articleStore struct {
	aiClient AIClient
	cfg      *ArticleConfig
}

func NewArticleStore(aiClient AIClient, config *ArticleConfig) Article {
	if config == nil {
		config = &ArticleConfig{
			MaxToken: 2048,
		}
	}

	return &articleStore{
		aiClient: aiClient,
		cfg:      config,
	}
}

func (s *articleStore) ExtractTags(ctx context.Context, content string, professionType models.PlatformType) (string, error) {
	if s.aiClient == nil {
		return "", fmt.Errorf("AI client is not initialized")
	}

	systemPrompt := models.GetExtractTagsSystemPrompt(professionType)

	message := models.AIChatMessage{
		SystemPrompt: systemPrompt,
		Text:         content,
		ImageUrls:    []string{},
	}

	opts := models.AIClientOptions{
		MaxTokens:      s.cfg.MaxToken,
		ResponseFormat: models.ResponseFormatJSON,
	}

	resp, err := s.aiClient.Generate(ctx, message, opts)
	if err != nil {
		return "", err
	}

	if resp == "" {
		return "", fmt.Errorf("empty response content from AI client")
	}

	return resp, nil
}

func (s *articleStore) Polish(ctx context.Context, content string, professionType models.PlatformType) (string, error) {
	if s.aiClient == nil {
		return "", fmt.Errorf("AI client is not initialized")
	}

	systemPrompt := models.GetPolishArticleSystemPrompt(professionType)

	message := models.AIChatMessage{
		SystemPrompt: systemPrompt,
		Text:         content,
		ImageUrls:    []string{},
	}

	opts := models.AIClientOptions{
		MaxTokens:      s.cfg.MaxToken,
		ResponseFormat: models.ResponseFormatText,
	}

	resp, err := s.aiClient.Generate(ctx, message, opts)
	if err != nil {
		return "", err
	}

	if resp == "" {
		return "", fmt.Errorf("empty response content from AI client")
	}

	return resp, nil
}
