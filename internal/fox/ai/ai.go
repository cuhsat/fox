//go:build !no_ai

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
	Model = "mistral"
)

const (
	input = `
As an logfile analyst, you are tasked with answering questions about log files,
or other line based text files.
Use only the following context to answer the question.
If you can't find the answer the provided context,
answer with: "This information is not available in the provided context."

Context (line numbers are in square brackets):
%s

Question:
%s

Answer precise, topic oriented and cite relevent lines.
If you refer to a specific line within the context,
give also the according line number i.e. "In line [123] ...".
`
)

var (
	llm *ollama.LLM
)

type Agent struct {
	sync.RWMutex

	file  *os.File              // agent file
	parts []llms.MessageContent // agent history
	ch    chan string           // agent channel
}

func Init(model string) bool {
	var err error

	if len(model) == 0 {
		model = Model
	}

	if strings.ToLower(model) == "default" {
		model = Model
	}

	llm, err = ollama.New(ollama.WithModel(model))

	if err != nil {
		sys.Error(err)
		return false
	}

	return true
}

func NewAgent() *Agent {
	return &Agent{
		file:  sys.TempFile("rag", ".txt"),
		parts: make([]llms.MessageContent, 0),
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
	a.human(fmt.Sprintf(input, string(h.Bytes()), s))

	if _, err := llm.GenerateContent(
		context.Background(),
		a.parts,
		llms.WithSeed(0),
		llms.WithTemperature(0),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) > 0 {
				a.ch <- string(chunk)
			} else {
				a.ch <- "\n\n"
			}
			return nil
		}),
	); err != nil {
		sys.Error(err)
	}
}

func (a *Agent) Listen(hi *history.History) {
	var buf strings.Builder

	for s := range a.ch {
		// response start
		if buf.Len() == 0 {
			s = strings.TrimLeft(s, " ")
		}

		// response chunk
		a.write(s)
		buf.WriteString(s)

		// response end
		if s == "\n\n" {
			s = buf.String()

			a.agent(s)
			hi.AddSystem(s)
			buf.Reset()
		}
	}
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

func (a *Agent) human(s string) {
	a.history(llms.ChatMessageTypeHuman, s)
}

func (a *Agent) agent(s string) {
	a.history(llms.ChatMessageTypeSystem, s)
}

func (a *Agent) history(r llms.ChatMessageType, s string) {
	a.Lock()
	a.parts = append(a.parts, llms.TextParts(r, s))
	a.Unlock()
}
