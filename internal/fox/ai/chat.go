package ai

import (
	"context"
	"fmt"
	"os"
	// "strings"
	"sync"

	"github.com/tmc/langchaingo/llms"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	//"github.com/cuhsat/fox/internal/pkg/user/history"
)

type Chat struct {
	sync.RWMutex

	file *os.File    		   // chat file
	msgs []llms.MessageContent // chat messages
	ch   chan string 		   // chat channel
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
	o.file.Close()
}

func (o *Chat) Prompt(s string) {
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

func (o *Chat) Listen(/*hi *history.History*/) {
	// var buf strings.Builder

	for s := range o.ch {
		// response start
		// if buf.Len() == 0 {
		// 	s = strings.TrimLeft(s, " ")
		// }

		// reponse chunk
		o.write(s)
		// buf.WriteString(s)

		// response end
		// if s == "\n\n" {
		// 	s = buf.String()

		// 	//o.system(s)

		// 	//hi.AddEntry("system", s)
		// 	buf.Reset()
		// }
	}
}

func (o *Chat) human(s string) {
	o.history(llms.ChatMessageTypeHuman, s)
}

// func (o *Chat) system(s string) {
// 	o.history(llms.ChatMessageTypeSystem, s)
// }

func (o *Chat) history(r llms.ChatMessageType, s string) {
	o.Lock()
	o.msgs = append(o.msgs, llms.TextParts(r, s))
	o.Unlock()
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
