package ai

import (
	"context"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

const (
	Default = "" // disabled
)

var (
	Client *api.Client
)

var (
	Model = Default
	Alive = &api.Duration{Duration: time.Minute * 10}
)

func Init(model string) {
	Model = model

	if len(model) == 0 || strings.ToLower(model) == "default" {
		return
	}

	llm, err := api.ClientFromEnvironment()

	if err != nil {
		return
	}

	// preload model in the background
	go func(*api.Client) {
		if llm.Chat(context.Background(), &api.ChatRequest{
			Model:     Model,
			KeepAlive: Alive,
		}, null) == nil {
			Client = llm
		}
	}(llm)
}

func IsInit() bool {
	return Client != nil
}

func null(_ api.ChatResponse) error {
	return nil
}
