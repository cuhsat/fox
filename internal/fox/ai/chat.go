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

	response  chan string           // response channel
	responses []llms.MessageContent // responses
}

func NewChat(model string) *Chat {
	if len(model) == 0 {
		model = Default
	}

	if strings.ToLower(model) == "default" {
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

		response:  make(chan string),
		responses: make([]llms.MessageContent, 0),
	}
}

func (o *Chat) Path() string {
	return o.file.Name()
}

func (o *Chat) Close() {
	o.file.Close()
}

func (o *Chat) Embed(h *heap.Heap) {
	// TODO: TEST
	ctx := context.Background()
	str := string(h.Bytes())

	o.llm.CreateEmbedding(ctx, strings.Split(str, "\n"))
}

func (o *Chat) Prompt(s string, h *heap.Heap) {
	o.write(fmt.Sprintf("%s %s\n", text.Chevron, s))
	o.human(s)

	if _, err := o.llm.GenerateContent(
		context.Background(),
		o.responses,
		llms.WithSeed(0),
		llms.WithTemperature(0),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) > 0 {
				o.response <- string(chunk)
			} else {
				o.response <- "\n\n"
			}
			return nil
		}),
	); err != nil {
		sys.Error(err)
	}
}

func (o *Chat) Listen(hi *history.History) {
	var buf strings.Builder

	for s := range o.response {
		// response start
		if buf.Len() == 0 {
			s = strings.TrimLeft(s, " ")
		}

		// reponse chunk
		o.write(s)
		buf.WriteString(s)

		// response end
		if s == "\n\n" {
			s = buf.String()

			o.system(s)

			hi.AddEntry("system", s)
			buf.Reset()
		}
	}
}

func (o *Chat) human(s string) {
	o.history(llms.ChatMessageTypeHuman, s)
}

func (o *Chat) system(s string) {
	o.history(llms.ChatMessageTypeSystem, s)
}

func (o *Chat) history(r llms.ChatMessageType, s string) {
	o.Lock()
	o.responses = append(o.responses, llms.TextParts(r, s))
	o.Unlock()
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
