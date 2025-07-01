package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

var (
	//go:embed prompt.txt
	Prompt string
)

var (
	client *api.Client
)

type Rag struct {
	sync.RWMutex

	model string        // agent model
	file  sys.File      // agent file
	keep  *api.Duration // agent keep alive
	msgs  []api.Message // agent history
	ch    chan string   // agent channel
}

func Init() bool {
	var err error

	client, err = api.ClientFromEnvironment()

	if err != nil {
		sys.Error(err)
	}

	return err == nil
}

func NewRag(model string) *Rag {
	return &Rag{
		model: model,
		file:  sys.TempFile("RAG"),
		keep:  &api.Duration{time.Minute * 10},
		msgs:  make([]api.Message, 0),
		ch:    make(chan string, 16),
	}
}

func (rag *Rag) Path() string {
	return rag.file.Name()
}

func (rag *Rag) Load() {
	go client.Chat(
		context.Background(),
		&api.ChatRequest{
			Model:     rag.model,
			KeepAlive: rag.keep,
		},
		func(_ api.ChatResponse) error {
			return nil
		})
}

func (rag *Rag) Prompt(s string, h *heap.Heap) {
	var sb strings.Builder

	rag.write(fmt.Sprintf("%s %s\n", text.User, s))

	for _, str := range *h.SMap() {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", str.Nr, str.Str))
	}

	rag.Lock()

	rag.msgs = append(rag.msgs, api.Message{
		Content: fmt.Sprintf(Prompt, sb.String(), s),
		Role:    "user",
	})

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:     rag.model,
		Messages:  rag.msgs,
		KeepAlive: rag.keep,
		Options: map[string]any{
			"num_ctx":     8192,
			"temperature": 0.2,
			"top_p":       0.5,
			"top_k":       10,
			"seed":        8211,
		},
	}

	rag.Unlock()

	if err := client.Chat(ctx, req, func(cr api.ChatResponse) error {
		if s := cr.Message.Content; len(s) > 0 {
			rag.ch <- s
		} else {
			rag.ch <- "\n\n"
		}

		return nil
	}); err != nil {
		sys.Error(err)
	}
}

func (rag *Rag) Listen(_ *history.History) {
	var buf strings.Builder

	for s := range rag.ch {
		// response start
		if buf.Len() == 0 {
			s = strings.TrimSpace(s)
		}

		// response chunk
		rag.write(s)
		buf.WriteString(s)

		// response end
		if s == "\n\n" {
			s = buf.String()

			buf.Reset()
		}
	}
}

func (rag *Rag) write(s string) {
	rag.Lock()

	_, err := rag.file.WriteString(s)

	if err != nil {
		sys.Error(err)
	}

	rag.Unlock()
}
