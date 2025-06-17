package main

import (
	"context"
	"log"

	"github.com/ollama/ollama/api"
)

type OllamaEngine struct {
	client *api.Client
	model  string
	ctx    context.Context
	memory []int
	stream bool
}

func InitOllamaEngine(model string) *OllamaEngine {
	ctx := context.Background()
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Failed to create Ollama client: %v", err)
	}
	return &OllamaEngine{
		client: client,
		ctx:    ctx,
		model:  model,
		stream: true,
		memory: []int{},
	}
}

func (o *OllamaEngine) Generate(prompt string, handler func(resp api.GenerateResponse) error, stream bool) {
	_ = o.client.Generate(o.ctx, &api.GenerateRequest{
		Model:   o.model,
		Prompt:  prompt,
		Stream:  &stream,
		Context: o.memory,
	}, func(resp api.GenerateResponse) error {
		o.memory = resp.Context
		return handler(resp)
	})
}
