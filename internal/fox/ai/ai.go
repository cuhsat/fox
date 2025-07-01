package ai

import (
	"context"
	"strings"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

const (
	Default = "mistral"
)

var (
	client *api.Client
)

var (
	model string
)

func Init(model string) bool {
	var err error

	if len(model) == 0 || strings.ToLower(model) == "default" {
		model = Default
	}

	client, err = api.ClientFromEnvironment()

	if err != nil {
		sys.Error(err)
	}

	// preload model
	go client.Chat(
		context.Background(),
		&api.ChatRequest{
			Model: model,
			// KeepAlive: ,
		},
		func(_ api.ChatResponse) error {
			return nil
		})

	return err == nil
}
