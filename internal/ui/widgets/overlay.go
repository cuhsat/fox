package widgets

import (
    "fmt"
    "sync"
    "time"

    "github.com/cuhsat/fx/internal/ui/themes"
    "github.com/gdamore/tcell/v2"
)

const (
    short = 1 // seconds
    long  = 3 // seconds
)

type Overlay struct {
    base
    m sync.RWMutex
    ch chan message
    buffer *message
}

type message struct {
    v string
    s tcell.Style
    t time.Duration
} 

func NewOverlay(ctx *Context, term tcell.Screen) *Overlay {
    return &Overlay{
        base: base{ctx, term},

        ch: make(chan message, 64),
    }
}

func (o *Overlay) Render(x, y, w, h int) {
    o.m.RLock()
    msg := o.buffer
    o.m.RUnlock()

    if msg != nil {
        o.print(x, y, fmt.Sprintf("%-*s", w, msg.v), msg.s)
    }
}

func (o *Overlay) SendError(err string) {
    o.ch <- message{err, themes.Overlay0, short}
}

func (o *Overlay) SendInfo(msg string) {
    o.ch <- message{msg, themes.Overlay1, short}
}

func (o *Overlay) Watch() {
    for msg := range o.ch {
        o.m.Lock()
        o.buffer = &msg
        o.m.Unlock()

        time.Sleep(msg.t * time.Second)

        o.m.Lock()
        o.buffer = nil
        o.m.Unlock()

        o.term.PostEvent(tcell.NewEventInterrupt(nil))
    }
}

func (o *Overlay) Close() {
    close(o.ch)
}
