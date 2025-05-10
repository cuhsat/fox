package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/types/heap"
	"github.com/cuhsat/fx/internal/pkg/user/history"
)

const (
	Default = "llama3.2"
)

const (
	role = "assistant"
	seed = 821119
	temp = 0
)

type Ollama struct {
	sync.RWMutex

	model   string        // ollama model
	file    *os.File      // ollama chat file
	client  *api.Client   // ollama client
	history []api.Message // ollama history
	content chan string   // ollama channel
}

func NewOllama(model string) *Ollama {
	if len(model) == 0 || strings.ToLower(model) == "default" {
		model = Default
	}

	client, err := api.ClientFromEnvironment()

	if err != nil {
		sys.Error(err)
		return nil
	}

	return &Ollama{
		model:   model,
		file:    sys.TempFile("ollama", ".txt"),
		client:  client,
		history: make([]api.Message, 0),
		content: make(chan string),
	}
}

func (o *Ollama) Path() string {
	return o.file.Name()
}

func (o *Ollama) Close() {
	o.file.Close()
}

func (o *Ollama) Prompt(s string, h *heap.Heap) {
	s = fmt.Sprintf("Context: %s", h.Bytes())

	o.Lock()
	o.history = append(o.history, api.Message{
		Role:    "user",
		Content: s,
	})
	o.Unlock()

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    o.model,
		Messages: o.history,
		Options: map[string]any{
			"temperature": temp,
			"seed":        seed,
		},
	}

	fn := func(r api.ChatResponse) error {
		if !r.Done {
			o.content <- r.Message.Content
		} else {
			o.content <- "\n"
		}

		return nil
	}

	err := o.client.Chat(ctx, req, fn)

	if err != nil {
		sys.Error(err)
	}
}

func (o *Ollama) Listen(hi *history.History) {
	var buf strings.Builder

	for s := range o.content {
		buf.WriteString(s)

		_, err := o.file.WriteString(s)

		if err != nil {
			sys.Error(err)
		}

		err = o.file.Sync()

		if err != nil {
			sys.Error(err)
		}

		if s == "\n" {
			b := strings.TrimSpace(buf.String())

			o.Lock()
			o.history = append(o.history, api.Message{
				Content: b,
				Role:    role,
			})
			o.Unlock()

			hi.AddEntry(role, strings.ReplaceAll(b, "\n", ""))

			buf.Reset()
		}
	}
}
