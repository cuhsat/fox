//go:build !minimal

package ai

import (
	"context"
	"time"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/app/ai/agent"
)

func NewAgent(ctx *app.Context) *agent.Agent {
	if len(ctx.Model()) == 0 {
		return nil // no model set
	}

	client, err := api.ClientFromEnvironment()

	if err != nil {
		return nil // no client found
	}

	ctxTo, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = client.Heartbeat(ctxTo)

	if err != nil {
		return nil // no server found
	}

	return agent.New(ctx)
}
