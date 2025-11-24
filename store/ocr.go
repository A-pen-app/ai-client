package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/A-pen-app/ai-clients/models"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/mq/v2"
	"github.com/tidwall/sjson"
)

type Config struct {
	MaxToken int64
	IsProd   bool
}

type ocrStore struct {
	mq       mq.MQ
	aiClient AIClient
	cfg      *Config
}

func NewOcrStore(mq mq.MQ, aiClient AIClient, config *Config) OCR {
	if config == nil {
		config = &Config{
			MaxToken: 1024,
			IsProd:   false,
		}
	}

	return &ocrStore{
		mq:       mq,
		aiClient: aiClient,
		cfg:      config,
	}
}

func (s *ocrStore) ScanName(ctx context.Context, link string) (string, error) {
	if s.aiClient == nil {
		return "", fmt.Errorf("AI client is not initialized")
	}

	message := models.AIChatMessage{
		SystemPrompt: models.SystemContent,
		Text:         models.NamePrompt,
		ImageUrls:    []string{link},
	}

	opts := models.AIClientOptions{
		MaxTokens:      s.cfg.MaxToken,
		ResponseFormat: models.ResponseFormatJSON,
	}

	resp, err := s.aiClient.Generate(ctx, message, opts)
	if err != nil {
		return "", err
	}

	result := struct {
		Name string `json:"name"`
	}{}

	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return "", err
	}

	return result.Name, nil
}

func (s *ocrStore) ScanRawInfo(ctx context.Context, userID string, link string, platformType models.PlatformType) (*models.OCRRawInfo, error) {
	if s.aiClient == nil {
		return nil, fmt.Errorf("AI client is not initialized")
	}

	prompt := models.GetInfoPrompt(platformType)

	message := models.AIChatMessage{
		SystemPrompt: models.SystemContent,
		Text:         prompt,
		ImageUrls:    []string{link},
	}

	opts := models.AIClientOptions{
		MaxTokens:      s.cfg.MaxToken,
		ResponseFormat: models.ResponseFormatJSON,
	}

	resp, err := s.aiClient.Generate(ctx, message, opts)
	if err != nil {
		return nil, err
	}

	if resp == "" {
		return nil, fmt.Errorf("empty response content from AI client")
	}

	modifiedJSON, err := sjson.Set(resp, "identify_url", link)
	if err != nil {
		return nil, err
	}

	ocrTopic := models.OCRTopicDev
	if s.cfg.IsProd {
		ocrTopic = models.OCRTopicProd
	}

	if err := s.mq.Send(string(ocrTopic), models.OCREventMessage{
		UserID:    userID,
		Payload:   modifiedJSON,
		CreatedAt: time.Now(),
		Type:      string(models.OCRMessageTypeIdentifyOCR),
		Source:    string(platformType),
	}); err != nil {
		logging.Errorw(ctx, "Failed to send ocr result", "error", err)
	}

	ocr := models.OCRRawInfo{}
	if err := json.Unmarshal([]byte(modifiedJSON), &ocr); err != nil {
		return nil, err
	}

	return &ocr, nil
}
