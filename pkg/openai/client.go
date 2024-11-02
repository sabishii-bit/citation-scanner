package openai

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIClient is a struct that handles OpenAI API interactions.
type OpenAIClient struct {
	client      *openai.Client
	model       openai.ChatModel
	systemRole  string
	temperature float64
	maxTokens   int64
}

// NewClient creates and returns a new OpenAIClient with default settings.
func NewClient(apiKey string, opts ...func(*OpenAIClient)) *OpenAIClient {
	if apiKey == "" {
		log.Fatal("OpenAI API key is not set")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	aiClient := &OpenAIClient{
		client:      client,
		model:       openai.ChatModelGPT4o,          // Default model
		systemRole:  "You are a helpful assistant.", // Default role
		temperature: 0.05,                           // Default temperature
		maxTokens:   16300,                          // Default max tokens in the return
	}

	// Apply options to override defaults if provided
	for _, opt := range opts {
		opt(aiClient)
	}

	return aiClient
}

// WithModel is an option to set a custom model.
func WithModel(model openai.ChatModel) func(*OpenAIClient) {
	return func(c *OpenAIClient) {
		c.model = model
	}
}

// WithSystemRole is an option to set a custom system role.
func WithSystemRole(role string) func(*OpenAIClient) {
	return func(c *OpenAIClient) {
		c.systemRole = role
	}
}

// WithTemperature is an option to set a custom temperature.
func WithTemperature(temp float64) func(*OpenAIClient) {
	return func(c *OpenAIClient) {
		c.temperature = temp
	}
}

// WithMaxTokens is an option to set a custom max tokens value.
func WithMaxTokens(tokens int64) func(*OpenAIClient) {
	return func(c *OpenAIClient) {
		c.maxTokens = tokens
	}
}

// SendChatRequest sends a chat request to the OpenAI API and returns the response.
func (c *OpenAIClient) SendChatRequest(prompt string) (string, error) {
	chatCompletion, err := c.client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(c.systemRole),
				openai.UserMessage(prompt),
			}),
			Model:       openai.F(c.model),
			MaxTokens:   openai.F(c.maxTokens),
			Temperature: openai.F(c.temperature),
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	if len(chatCompletion.Choices) > 0 {
		return chatCompletion.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response received")
}
