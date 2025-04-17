package library

import (
    "fmt"
    "time"

    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/gdamore/tcell/v2"
)

const (
    delayShort = 2 // seconds
    delayLong  = 5 // seconds
)

type Overlay struct {
    base

    ch chan message
    buffer *message
}

type message struct {
    value string
    style tcell.Style
    delay time.Duration
} 

func NewOverlay(ctx *Context, term tcell.Screen) *Overlay {
    return &Overlay{
        base: base{
            ctx: ctx,
            term: term,
        },
        
        ch: make(chan message, 64),
        buffer: nil,
    }
}

func (o *Overlay) Render(x, y, w, h int) {
    if o.buffer != nil {
        s := fmt.Sprintf("%-*s", w, o.buffer.value)

        o.print(x, y, s, o.buffer.style)
    }
}

func (o *Overlay) SendError(err string) {
    o.ch <- message{
        value: err,
        style: themes.Overlay0,
        delay: delayLong,
    }
}

func (o *Overlay) SendInfo(msg string) {
    o.ch <- message{
        value: msg,
        style: themes.Overlay1,
        delay: delayShort,
    }
}

func (o *Overlay) Watch() {
    for m := range o.ch {
        o.buffer = &m

        time.Sleep(m.delay * time.Second)

        o.buffer = nil

        o.term.PostEvent(tcell.NewEventInterrupt(nil))
    }
}

func (o *Overlay) Close() {
    close(o.ch)
}
