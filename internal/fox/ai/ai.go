//go:build ai

package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

const (
	Build = true
)

const (
	Default = "mistral"
)

var (
	llm *ollama.LLM = nil
)

type Chat struct {
	sync.RWMutex

	file *os.File              // chat file
	msgs []llms.MessageContent // chat messages
	ch   chan string           // chat channel
}

func Init(model string) bool {
	var err error

	if len(model) == 0 {
		model = Default
	}

	if strings.ToLower(model) == "default" {
		model = Default
	}

	llm, err = ollama.New(ollama.WithModel(model))

	if err != nil {
		sys.Error(err)
		return false
	}

	return true
}

func NewChat() *Chat {
	return &Chat{
		file: sys.TempFile("chat", ".txt"),
		msgs: make([]llms.MessageContent, 0),
		ch:   make(chan string, 16),
	}
}

func (o *Chat) Path() string {
	return o.file.Name()
}

func (o *Chat) Close() {
	_ = o.file.Close()
}

func (o *Chat) Prompt(s string, b []byte) {
	o.write(fmt.Sprintf("%s %s\n", text.Chevron, s))
	o.human(s)

	if _, err := llm.GenerateContent(
		context.Background(),
		o.msgs,
		llms.WithSeed(0),
		llms.WithTemperature(0),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) > 0 {
				o.ch <- string(chunk)
			} else {
				o.ch <- "\n\n"
			}
			return nil
		}),
	); err != nil {
		sys.Error(err)
	}
}

func (o *Chat) Listen(hi *history.History) {
	var buf strings.Builder

	for s := range o.ch {
		// response start
		if buf.Len() == 0 {
			s = strings.TrimLeft(s, " ")
		}

		// response chunk
		o.write(s)
		buf.WriteString(s)

		// response end
		if s == "\n\n" {
			s = buf.String()

			o.system(s)
			hi.AddSystem(s)
			buf.Reset()
		}
	}
}

func (o *Chat) write(s string) {
	o.Lock()

	_, err := o.file.WriteString(s)

	if err != nil {
		sys.Error(err)
	}

	err = o.file.Sync()

	if err != nil {
		sys.Error(err)
	}

	o.Unlock()
}

func (o *Chat) human(s string) {
	o.history(llms.ChatMessageTypeHuman, s)
}

func (o *Chat) system(s string) {
	o.history(llms.ChatMessageTypeSystem, s)
}

func (o *Chat) history(r llms.ChatMessageType, s string) {
	o.Lock()
	o.msgs = append(o.msgs, llms.TextParts(r, s))
	o.Unlock()
}
