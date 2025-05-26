//go:build !no_ai

package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	_ "embed"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

const (
	Build = true
)

const (
	Model = "mistral"
)

var (
	//go:embed prompt.txt
	Prompt string
)

var (
	rag *api.Client
)

type Agent struct {
	sync.RWMutex

	model string        // agent model
	file  *os.File      // agent file
	msgs  []api.Message // agent history
	ch    chan string   // agent channel
}

func Init() bool {
	var err error

	rag, err = api.ClientFromEnvironment()

	if err != nil {
		sys.Error(err)
	}

	return err == nil
}

func NewAgent(model string) *Agent {
	if len(model) == 0 {
		model = Model
	}

	if strings.ToLower(model) == "default" {
		model = Model
	}

	return &Agent{
		model: model,
		file:  sys.TempFile(),
		msgs:  make([]api.Message, 0),
		ch:    make(chan string, 16),
	}
}

func (a *Agent) Path() string {
	return a.file.Name()
}

func (a *Agent) Close() {
	_ = a.file.Close()
}

func (a *Agent) Prompt(s string, h *heap.Heap) {
	a.write(fmt.Sprintf("%s %s\n", text.Chevron, s))

	a.addHeap(h)
	a.addPrompt(fmt.Sprintf(Prompt, s))

	a.RLock()

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    a.model,
		Messages: a.msgs,
		Options: map[string]any{
			"num_ctx":     4096,
			"temperature": 0.2,
			"top_p":       0.5,
			"top_k":       10,
			"seed":        82,
		},
	}

	a.RUnlock()

	if err := rag.Chat(ctx, req, func(cr api.ChatResponse) error {
		if s := cr.Message.Content; len(s) > 0 {
			a.ch <- s
		} else {
			a.ch <- "\n\n"
		}

		return nil
	}); err != nil {
		sys.Error(err)
	}
}

func (a *Agent) Listen(hi *history.History) {
	var buf strings.Builder

	for s := range a.ch {
		// response start
		if buf.Len() == 0 {
			s = strings.TrimSpace(s)
		}

		// response chunk
		a.write(s)
		buf.WriteString(s)

		// response end
		if s == "\n\n" {
			s = buf.String()

			a.addAnswer(s)
			hi.AddSystem(s)
			buf.Reset()
		}
	}
}

func (a *Agent) addHeap(h *heap.Heap) {
	for _, str := range *h.SMap() {
		a.addMessage("tool", fmt.Sprintf("[%d] %s", str.Nr, str.Str))
	}
}

func (a *Agent) addPrompt(s string) {
	a.addMessage("user", s)
}

func (a *Agent) addAnswer(s string) {
	a.addMessage("assistant", s)
}

func (a *Agent) addMessage(r, s string) {
	a.Lock()
	a.msgs = append(a.msgs, api.Message{
		Role:    r,
		Content: s,
	})
	a.Unlock()
}

func (a *Agent) write(s string) {
	a.Lock()

	_, err := a.file.WriteString(s)

	if err != nil {
		sys.Error(err)
	}

	err = a.file.Sync()

	if err != nil {
		sys.Error(err)
	}

	a.Unlock()
}
