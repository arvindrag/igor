package main

import (
	"context"
	"log"

	"github.com/ollama/ollama/api"
)

type OllamaClient struct {
	client *api.Client
	model  string
	ctx    context.Context
}

func BuildOllamaClient(model string) OllamaClient {
	ctx := context.Background()
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Failed to create Ollama client: %v", err)
	}
	return OllamaClient{
		client: client,
		ctx:    ctx,
		model:  model,
	}
}

func (o OllamaClient) Generate(prompt string, handler func(resp api.GenerateResponse) error, stream bool) {
	true_ := true
	_ = o.client.Generate(o.ctx, &api.GenerateRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: &true_,
	}, handler)
}
