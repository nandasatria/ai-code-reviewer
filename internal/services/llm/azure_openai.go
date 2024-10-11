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
	endpoint2  string
	key2       string
)

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

	endpoint2 = os.Getenv("AZURE_OPENAI_ENDPOINT2")
	key2 = os.Getenv("AZURE_OPENAI_KEY2")

	if endpoint == "" || key == "" || model == "" || apiVersion == "" {
		log.Fatalf("Environment variables are not set properly.")
	}
}

func ChatMPN1(messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	log.Printf("Start Running with endpoint: %s\n", endpoint)
	res, err := Chat(messages, endpoint, key)
	if err != nil {
		log.Panicf("Error from endpoint: %s", endpoint)
	} else {
		log.Printf("Success Get Request from endpoint : %s", endpoint)
	}

	return res, err
}

func ChatMPN2(messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	log.Printf("Start Running with endpoint: %s\n", endpoint2)

	res, err := Chat(messages, endpoint2, key2)
	if err != nil {
		log.Panicf("Error from endpoint: %s", endpoint)
	} else {
		log.Printf("Success Get Request from endpoint : %s", endpoint)
	}
	return res, err
}

func Chat(messages []openai.ChatCompletionMessageParamUnion, endpointlocal string, keylocal string) (string, error) {
	client := openai.NewClient(
		azure.WithEndpoint(endpointlocal, apiVersion),
		azure.WithAPIKey(keylocal),
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
