//go:build !minimal

package ai

import (
	"context"
	"time"

	"github.com/ollama/ollama/api"
)

var (
	Model = "" // disabled by default
	Alive = &api.Duration{Duration: time.Minute * 10}
)

func IsAvailable() bool {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = client.Heartbeat(ctx)

	return err == nil
}

func Load(model string) {
	if len(model) == 0 || !IsAvailable() {
		return
	}

	Model = model

	go func() {
		client, err := api.ClientFromEnvironment()

		if err != nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_ = client.Chat(ctx, &api.ChatRequest{
			Model:     Model,
			KeepAlive: Alive,
		}, func(cr api.ChatResponse) error {
			return nil // preloaded model
		})
	}()
}
