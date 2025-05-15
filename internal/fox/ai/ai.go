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
	"github.com/cuhsat/fox/internal/pkg/types/heap"
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

	file  *os.File              // chat file
	parts []llms.MessageContent // chat parts
	ch    chan string           // chat channel
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

	// TODO: Add embedding model with chain

	if err != nil {
		sys.Error(err)
		return false
	}

	return true
}

func NewChat() *Chat {
	return &Chat{
		file:  sys.TempFile("chat", ".txt"),
		parts: make([]llms.MessageContent, 0),
		ch:    make(chan string, 16),
	}
}

func (c *Chat) Path() string {
	return c.file.Name()
}

func (c *Chat) Close() {
	_ = c.file.Close()
}

func (c *Chat) Prompt(s string, h *heap.Heap) {
	c.write(fmt.Sprintf("%s %s\n", text.Chevron, s))
	c.human(s)

	em := make([]string, h.Lines())

	for _, str := range *h.SMap() {
		em = append(em, h.Unmap(&str))
	}

	if _, err := llm.CreateEmbedding(
		context.Background(),
		em,
	); err != nil {
		sys.Error(err)
	}

	if _, err := llm.GenerateContent(
		context.Background(),
		c.parts,
		llms.WithSeed(0),
		llms.WithTemperature(0),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) > 0 {
				c.ch <- string(chunk)
			} else {
				c.ch <- "\n\n"
			}
			return nil
		}),
	); err != nil {
		sys.Error(err)
	}
}

func (c *Chat) Listen(hi *history.History) {
	var buf strings.Builder

	for s := range c.ch {
		// response start
		if buf.Len() == 0 {
			s = strings.TrimLeft(s, " ")
		}

		// response chunk
		c.write(s)
		buf.WriteString(s)

		// response end
		if s == "\n\n" {
			s = buf.String()

			c.system(s)
			hi.AddSystem(s)
			buf.Reset()
		}
	}
}

func (c *Chat) write(s string) {
	c.Lock()

	_, err := c.file.WriteString(s)

	if err != nil {
		sys.Error(err)
	}

	err = c.file.Sync()

	if err != nil {
		sys.Error(err)
	}

	c.Unlock()
}

func (c *Chat) human(s string) {
	c.history(llms.ChatMessageTypeHuman, s)
}

func (c *Chat) system(s string) {
	c.history(llms.ChatMessageTypeSystem, s)
}

func (c *Chat) history(r llms.ChatMessageType, s string) {
	c.Lock()
	c.parts = append(c.parts, llms.TextParts(r, s))
	c.Unlock()
}
