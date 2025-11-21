package openai

import (
	"context"
	"fmt"

	"github.com/A-pen-app/ai-clients/models"
	"github.com/A-pen-app/ai-clients/store"
	"github.com/openai/openai-go/v2"
)

type Client struct {
	client       *openai.Client
	defaultModel openai.ChatModel
}

func NewClient(client *openai.Client, model openai.ChatModel) store.Client {
	if model == "" {
		model = openai.ChatModelGPT4o
	}
	return &Client{
		client:       client,
		defaultModel: model,
	}
}

func (c *Client) Generate(ctx context.Context, message models.Message, opts models.Options) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("openai client is not initialized")
	}

	model := c.defaultModel
	if opts.Model != "" {
		model = openai.ChatModel(opts.Model)
	}

	var userContentParts []openai.ChatCompletionContentPartUnionParam

	if message.Text != "" {
		userContentParts = append(userContentParts, openai.TextContentPart(message.Text))
	}

	for _, url := range message.ImageUrls {
		userContentParts = append(userContentParts, openai.ImageContentPart(
			openai.ChatCompletionContentPartImageImageURLParam{
				URL: url,
			},
		))
	}

	messages := []openai.ChatCompletionMessageParamUnion{}
	if message.SystemPrompt != "" {
		messages = append(messages, openai.SystemMessage(message.SystemPrompt))
	}
	messages = append(messages, openai.UserMessage(userContentParts))

	params := openai.ChatCompletionNewParams{
		Model:    model,
		Messages: messages,
	}

	if opts.ResponseFormat == models.ResponseFormatJSON {
		params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		}
	}

	if opts.MaxTokens > 0 {
		params.MaxTokens = openai.Int(opts.MaxTokens)
	}

	resp, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty response choices from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
