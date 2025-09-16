package llm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user/config"
)

type LLM struct {
	sync.RWMutex

	client  *api.Client   // chat client
	alive   *api.Duration // chat alive
	history []api.Message // chat history
}

func New(model string, keep time.Duration) *LLM {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		sys.Panic(err)
	}

	alive := &api.Duration{Duration: keep}

	// preload model
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		_ = client.Chat(ctx, &api.ChatRequest{
			Model:     model,
			KeepAlive: alive,
		}, func(cr api.ChatResponse) error {
			return nil // preloaded model
		})
	}()

	return &LLM{
		client:  client,
		alive:   alive,
		history: make([]api.Message, 0),
	}
}

func (llm *LLM) Use(model string, fn api.PullProgressFunc) error {
	ctx := context.Background()
	req := &api.PullRequest{
		Model: model,
	}

	return llm.client.Pull(ctx, req, fn)
}

func (llm *LLM) Ask(model, query, lines string, fn api.ChatResponseFunc) error {
	llm.AddSystem(fmt.Sprintf(fox.Prompt, lines))
	llm.AddUser(query)

	llm.RLock()

	cfg := config.Get()
	ctx := context.Background()
	req := &api.ChatRequest{
		Model:     model,
		KeepAlive: llm.alive,
		Messages:  llm.history,
		Options: map[string]any{
			"num_ctx":     cfg.GetInt("ai.num_ctx"),
			"temperature": cfg.GetFloat64("ai.temp"),
			"seed":        cfg.GetInt("ai.seed"),
			"top_k":       cfg.GetInt("ai.top_k"),
			"top_p":       cfg.GetFloat64("ai.top_p"),
		},
	}

	llm.RUnlock()

	return llm.client.Chat(ctx, req, fn)
}

func (llm *LLM) AddUser(content string) {
	llm.addMessage("user", content)
}

func (llm *LLM) AddSystem(content string) {
	llm.addMessage("system", content)
}

func (llm *LLM) AddAssistant(content string) {
	llm.addMessage("assistant", content)
}

func (llm *LLM) addMessage(role, content string) {
	llm.Lock()
	llm.history = append(llm.history, api.Message{
		Role:    role,
		Content: content,
	})
	llm.Unlock()
}
