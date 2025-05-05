package widgets

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/cuhsat/fx/internal/app/ui/context"
	"github.com/cuhsat/fx/internal/app/ui/themes"
)

const (
	short = 1 // seconds
	long  = 2 // seconds
)

type Overlay struct {
	base
	m      sync.RWMutex
	ch     chan message
	buffer *message
}

type message struct {
	v string
	s tcell.Style
	t time.Duration
}

func NewOverlay(ctx *context.Context) *Overlay {
	return &Overlay{
		base: base{ctx},

		ch: make(chan message, 64),
	}
}

func (o *Overlay) Render(x, y, w, h int) {
	o.m.RLock()
	msg := o.buffer
	o.m.RUnlock()

	if msg != nil {
		o.print(x, y, fmt.Sprintf(" %-*s", w-1, msg.v), msg.s)
	}
}

func (o *Overlay) Listen() {
	for msg := range o.ch {
		o.m.Lock()
		o.buffer = &msg
		o.m.Unlock()

		time.Sleep(msg.t * time.Second)

		o.m.Lock()
		o.buffer = nil
		o.m.Unlock()

		o.ctx.Interrupt()
	}
}

func (o *Overlay) SendError(err string) {
	o.ch <- message{err, themes.Overlay0, long}
}

func (o *Overlay) SendInfo(msg string) {
	o.ch <- message{msg, themes.Overlay1, short}
}

func (o *Overlay) Close() {
	close(o.ch)
}
