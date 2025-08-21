package agent

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/app/ai/agent/llm"
	"github.com/cuhsat/fox/internal/app/ai/agent/rag"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

type Agent struct {
	File sys.File // agent chat file

	buf *heap.Heap // buffered heap

	llm *llm.LLM // agent llm
	rag *rag.RAG // agent rag

	ch chan string // answer channel
}

func New() *Agent {
	a := &Agent{
		File: sys.Create("agent"),

		llm: llm.New(),
		rag: rag.New(),

		ch: make(chan string, 64),
	}

	go a.listen()

	return a
}

func (a *Agent) Close() {
	close(a.ch)
}

func (a *Agent) PS1(query string) {
	_, _ = a.File.WriteString(fmt.Sprintf("%s %s\n", text.PS1, query))
}

func (a *Agent) Ask(query string, h *heap.Heap) {
	if h.Type != types.Agent {
		a.buf = h // buffer last valid heap
	}

	col := a.rag.Embed(a.buf)

	if col == nil {
		return
	}

	ctx := a.rag.Query(query, col)

	if len(ctx) == 0 {
		return
	}

	a.llm.Ask(query, ctx, func(res api.ChatResponse) error {
		if len(res.Message.Content) > 0 {
			a.ch <- res.Message.Content
		} else {
			a.ch <- "\n\n"
		}

		return nil
	})
}

func (a *Agent) listen() {
	flg, end := flags.Get(), true

	var sb strings.Builder

	for s := range a.ch {
		// response start
		if end {
			s = strings.TrimSpace(s)
		}

		s = strings.Replace(s, "  ", "", 1)

		// response chunk
		if !flg.Print {
			_, _ = a.File.WriteString(s)
		} else {
			_, _ = fmt.Print(s)
		}

		// response end
		end = s == "\n\n"

		sb.WriteString(s)

		if end {
			a.llm.AddSystem(sb.String())
			sb.Reset()
		}
	}
}
