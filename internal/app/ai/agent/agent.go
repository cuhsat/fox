package agent

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/ollama/ollama/api"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/app/ai/agent/llm"
	"github.com/cuhsat/fox/internal/app/ai/agent/rag"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

type Agent struct {
	File sys.File
	busy atomic.Bool

	buf *heap.Heap
	ctx *app.Context
	llm *llm.LLM
	rag *rag.RAG

	ch chan string
}

func New(ctx *app.Context) *Agent {
	a := &Agent{
		File: fs.Create("/fox/agent"),

		ctx: ctx,
		llm: llm.New(ctx.Model(), time.Minute*10),
		rag: rag.New(),

		ch: make(chan string, 64),
	}

	a.busy.Store(false)

	go a.listen()

	return a
}

func (a *Agent) IsBusy() bool {
	return a.busy.Load()
}

func (a *Agent) String() string {
	return fmt.Sprintf("Agent %c %s", text.Icons().HSep, a.ctx.Model())
}

func (a *Agent) Prompt(query string) {
	_, _ = a.File.WriteString(fmt.Sprintf("%c %s\n", text.Icons().Ps1, query))
}

func (a *Agent) Process(query string, h *heap.Heap) {
	if h.Type != types.Agent {
		a.buf = h // buffer last valid heap
	}

	if !a.load(query) {
		a.query(query)
	}
}

func (a *Agent) Close() {
	close(a.ch)
}

func (a *Agent) load(query string) bool {
	var model string

	if !strings.HasPrefix(query, "use model") {
		return false
	}

	model = strings.TrimPrefix(query, "use model")
	model = strings.TrimSpace(model)

	a.busy.Store(true)

	err := a.llm.Use(model, func(res api.ProgressResponse) error {
		if res.Completed >= res.Total {
			a.busy.Store(false)
		} else {
			a.ch <- "."
		}

		return nil
	})

	if err != nil {
		a.busy.Store(false)
		sys.Error(err)

		return true
	}

	_, _ = a.File.WriteString(fmt.Sprintf("Using model %s\n", model))

	a.ctx.ChangeModel(model)

	return true
}

func (a *Agent) query(query string) {
	col := a.rag.Embed(a.buf)

	if col == nil {
		return
	}

	ctx := a.rag.Query(query, col)

	if len(ctx) == 0 {
		return
	}

	a.busy.Store(true)

	err := a.llm.Ask(a.ctx.Model(), query, ctx, func(res api.ChatResponse) error {
		if len(res.Message.Content) == 0 {
			a.busy.Store(false)
			a.ch <- "\n\n"
		} else {
			a.ch <- res.Message.Content
		}

		return nil
	})

	if err != nil {
		a.busy.Store(false)
		sys.Error(err)
	}
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
			a.llm.AddAssistant(sb.String())
			sb.Reset()
		}
	}
}
