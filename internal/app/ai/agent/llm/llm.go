package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/app/ai"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user/config"
)

type LLM struct {
	client  *api.Client   // chat client
	history []api.Message // chat history
}

func New() *LLM {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		sys.Panic(err)
	}

	return &LLM{
		client:  client,
		history: make([]api.Message, 0),
	}
}

func (llm *LLM) Ask(query, lines string, fn api.ChatResponseFunc) {
	llm.AddUser(fmt.Sprintf(app.Prompt, query, lines))

	cfg := config.Get()
	ctx := context.Background()
	req := &api.ChatRequest{
		Model:     ai.Model,
		KeepAlive: &api.Duration{Duration: time.Minute * 10},
		Messages:  llm.history,
		Options: map[string]any{
			"num_ctx":     cfg.GetInt("ai.num_ctx"),
			"temperature": cfg.GetFloat64("ai.temp"),
			"seed":        cfg.GetInt("ai.seed"),
			"top_k":       cfg.GetInt("ai.top_k"),
			"top_p":       cfg.GetFloat64("ai.top_p"),
		},
	}

	err := llm.client.Chat(ctx, req, fn)

	if err != nil {
		sys.Error(err)
	}
}

func (llm *LLM) AddUser(content string) {
	llm.history = append(llm.history, api.Message{
		Role:    "user",
		Content: content,
	})
}

func (llm *LLM) AddSystem(content string) {
	llm.history = append(llm.history, api.Message{
		Role:    "system",
		Content: content,
	})
}
