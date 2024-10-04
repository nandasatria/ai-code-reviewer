package llm

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

var (
	endpoint   string
	key        string
	model      string
	apiVersion string
)

func sampleUse() {
	var messages []openai.ChatCompletionMessageParamUnion
	messages = append(messages, openai.UserMessage("Hello"))

	response, err := Chat(messages)
	if err != nil {
		log.Fatalf("Error while chat: %v\n", err)
		return
	}
	fmt.Println(response)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}

	endpoint = os.Getenv("AZURE_OPENAI_ENDPOINT")
	key = os.Getenv("AZURE_OPENAI_KEY")
	model = os.Getenv("AZURE_OPENAI_MODEL")
	apiVersion = os.Getenv("AZURE_OPENAI_API_VERSION")

	if endpoint == "" || key == "" || model == "" || apiVersion == "" {
		log.Fatalf("Environment variables are not set properly.")
	}
}

func Chat(messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	client := openai.NewClient(
		azure.WithEndpoint(endpoint, apiVersion),
		azure.WithAPIKey(key),
	)

	ctx := context.Background()

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:    openai.F(messages),
		Seed:        openai.Int(1),
		Model:       openai.F(model),
		MaxTokens:   openai.Int(4000),
		Temperature: openai.Float(0.5),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get chat completion: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from completion")
	}

	return completion.Choices[0].Message.Content, nil
}
