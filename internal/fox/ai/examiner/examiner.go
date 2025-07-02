package examiner

import (
	"context"
	"fmt"
	"strings"

	_ "embed"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/fox/ai"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

type Examiner struct {
	File sys.File // examiner chat file

	heap *heap.Heap // buffered heap

	hi []api.Message // examiner history
	ch chan string   // examiner channel
}

func New() *Examiner {
	return &Examiner{
		File: sys.TempFile("Examiner"),

		hi: make([]api.Message, 0),
		ch: make(chan string, 16),
	}
}

func (e *Examiner) User(s string) {
	e.File.WriteString(fmt.Sprintf("%s %s\n", text.User, s))
}

func (e *Examiner) Query(s string, h *heap.Heap) {
	if !ai.IsInit() {
		return
	}

	if h.Type != types.Prompt {
		e.heap = h // use last normal heap
	}

	var sb strings.Builder

	for _, str := range *e.heap.SMap() {
		sb.WriteString(fmt.Sprintf("Line %d: %s\n", str.Nr, str.Str))
	}

	e.hi = append(e.hi, api.Message{
		Role:    "user",
		Content: fmt.Sprintf(fox.Prompt, sb.String(), s),
	})

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:     ai.Model,
		KeepAlive: ai.Alive,
		Messages:  e.hi,
		Options: map[string]any{
			"temperature": 0.2,
			"top_p":       0.5,
			"top_k":       10,
			"seed":        8211,
		},
	}

	if err := ai.Client.Chat(ctx, req, func(cr api.ChatResponse) error {
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

func (e *Examiner) Listen() {
	t := true

	for s := range e.ch {
		// response start
		if t {
			s = strings.TrimSpace(s)
		}

		// response chunk
		e.File.WriteString(s)

		// response end
		t = (s == "\n\n")
	}
}
