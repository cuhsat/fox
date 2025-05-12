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
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

const (
	Default = "mistral"
)

type Chat struct {
	sync.RWMutex

	model string      // chat model
	file  *os.File    // chat file
	llm   *ollama.LLM // ollama llm
	ch    chan string // responses
}

func NewChat(model string) *Chat {
	if len(model) == 0 || strings.ToLower(model) == "default" {
		model = Default
	}

	llm, err := ollama.New(ollama.WithModel(model))

	if err != nil {
		sys.Error(err)
		return nil
	}

	return &Chat{
		model: model,
		file:  sys.TempFile("chat", ".txt"),
		llm:   llm,
		ch:    make(chan string),
	}
}

func (o *Chat) Path() string {
	return o.file.Name()
}

func (o *Chat) Close() {
	o.file.Close()
}

func (o *Chat) Prompt(s string, h *heap.Heap) {
	o.write(fmt.Sprintln(s))
	o.write(fmt.Sprintln(text.HSep))

	ctx := context.Background()

	// content := []llms.MessageContent{
	// 	llms.TextParts(llms.ChatMessageTypeSystem, "You are a company branding design wizard."),
	// 	llms.TextParts(llms.ChatMessageTypeHuman, "What would be a good company name a company that makes colorful socks?"),
	// }

	_, err := o.llm.Call(ctx, s,
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
	)

	if err != nil {
		sys.Error(err)
	}
}

func (o *Chat) Listen(hi *history.History) {
	var buf strings.Builder

	for s := range o.ch {
		if buf.Len() == 0 {
			s = strings.TrimLeft(s, " ")
		}

		o.write(s)

		buf.WriteString(s)

		if s == "\n\n" {
			// o.write(fmt.Sprintln(text.HSep))
			hi.AddEntry("assistant", buf.String())
			buf.Reset()
		}
	}
}

func (o *Chat) write(s string) {
	o.Lock()
	defer o.Unlock()

	_, err := o.file.WriteString(s)

	if err != nil {
		sys.Error(err)
	}

	err = o.file.Sync()

	if err != nil {
		sys.Error(err)
	}
}
