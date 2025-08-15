package examiner

import (
	"fmt"
	"strings"

	"github.com/hiforensics/fox/internal/fox/ai"
	"github.com/hiforensics/fox/internal/fox/ai/examiner/llm"
	"github.com/hiforensics/fox/internal/fox/ai/examiner/rag"
	"github.com/hiforensics/fox/internal/pkg/text"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/file"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/ollama/ollama/api"
)

type Examiner struct {
	File file.File // examiner chat file

	buf *heap.Heap // buffered heap
	llm *llm.LLM   // examiner llm
	rag *rag.RAG   // examiner rag

	ch chan string // answer channel
}

func New() *Examiner {
	e := &Examiner{
		File: file.New("Examiner"),

		llm: llm.New(),
		rag: rag.New(),

		ch: make(chan string, 16),
	}

	go e.listen()

	return e
}

func (e *Examiner) PS1(query string) {
	_, _ = e.File.WriteString(fmt.Sprintf("%s %s\n", text.PS1, query))
}

func (e *Examiner) Ask(query string, h *heap.Heap) {
	if !ai.IsInit() {
		return
	}

	if h.Type != types.Prompt {
		e.buf = h // buffer last regular heap
	}

	col := e.rag.Embed(e.buf)

	if col == nil {
		return
	}

	ctx := e.rag.Query(query, col)

	if len(ctx) == 0 {
		return
	}

	e.llm.Ask(query, ctx, func(res api.ChatResponse) error {
		if len(res.Message.Content) > 0 {
			e.ch <- res.Message.Content
		} else {
			e.ch <- "\n\n"
		}

		return nil
	})
}

func (e *Examiner) listen() {
	t := true

	for s := range e.ch {
		// response start
		if t {
			s = strings.TrimSpace(s)
		}

		s = strings.Replace(s, "  ", "", 1)

		// response chunk
		_, _ = e.File.WriteString(s)

		// response end
		t = s == "\n\n"
	}
}
