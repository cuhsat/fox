package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

type Examiner struct {
	sync.RWMutex

	file sys.File      // examiner file
	keep *api.Duration // examiner keep alive
	msgs []api.Message // examiner history
	ch   chan string   // examiner channel
}

func NewExaminer() *Examiner {
	return &Examiner{
		file: sys.TempFile("Examiner"),
		keep: &api.Duration{time.Minute * 10},
		msgs: make([]api.Message, 0),
		ch:   make(chan string, 16),
	}
}

func (e *Examiner) Path() string {
	return e.file.Name()
}

func (e *Examiner) Prompt(s string, h *heap.Heap) {
	var sb strings.Builder

	e.write(fmt.Sprintf("%s %s\n", text.User, s))

	for _, str := range *h.SMap() {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", str.Nr, str.Str))
	}

	e.Lock()

	e.msgs = append(e.msgs, api.Message{
		Content: fmt.Sprintf(fox.Prompt, sb.String(), s),
		Role:    "user",
	})

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:     model,
		Messages:  e.msgs,
		KeepAlive: e.keep,
		Options: map[string]any{
			"num_ctx":     8192,
			"temperature": 0.2,
			"top_p":       0.5,
			"top_k":       10,
			"seed":        8211,
		},
	}

	e.Unlock()

	if err := client.Chat(ctx, req, func(cr api.ChatResponse) error {
		if s := cr.Message.Content; len(s) > 0 {
			e.ch <- s
		} else {
			e.ch <- "\n\n"
		}

		return nil
	}); err != nil {
		sys.Error(err)
	}
}

func (e *Examiner) Listen(_ *history.History) {
	var buf strings.Builder

	for s := range e.ch {
		// response start
		if buf.Len() == 0 {
			s = strings.TrimSpace(s)
		}

		// response chunk
		e.write(s)

		buf.WriteString(s)

		// response end
		if s == "\n\n" {
			s = buf.String()

			buf.Reset()
		}
	}
}

func (e *Examiner) write(s string) {
	e.Lock()

	_, err := e.file.WriteString(s)

	if err != nil {
		sys.Error(err)
	}

	e.Unlock()
}
