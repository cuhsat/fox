package llm

import (
	"context"
	"fmt"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/app/ai"
	"github.com/cuhsat/fox/internal/pkg/sys"
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
	llm.AddUser(fmt.Sprintf(app.Base, query, lines))

	flg := flags.Get().AI
	ctx := context.Background()
	req := &api.ChatRequest{
		Model:     ai.Model,
		KeepAlive: ai.Alive,
		Messages:  llm.history,
		Options: map[string]any{
			"num_ctx":     flg.NumCtx,
			"temperature": flg.Temp,
			"seed":        flg.Seed,
			"top_k":       flg.TopK,
			"top_p":       flg.TopP,
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
