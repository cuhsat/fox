package widget

import (
    "fmt"
    "time"

    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    DelayShort = 2 // seconds
    DelayLong  = 5 // seconds
)

type Overlay struct {
    widget

    ch chan message
    buffer *message
}

type message struct {
    value string
    style tcell.Style
    delay time.Duration
} 

func NewOverlay(screen tcell.Screen) *Overlay {
    return &Overlay{
        widget: widget{
            screen: screen,
        },
        
        ch: make(chan message),
        buffer: nil,
    }
}

func (o *Overlay) Render(x, y, w, h int) {
    if o.buffer != nil {
        s := fmt.Sprintf("%-*s", w, o.buffer.value)

        o.print(x, y, s, o.buffer.style)
    }
}

func (o *Overlay) SendError(err error) {
    o.ch <- message{
        value: err.Error(),
        style: theme.Error,
        delay: DelayLong,
    }
}

func (o *Overlay) SendStatus(msg string) {
    o.ch <- message{
        value: msg,
        style: theme.Info,
        delay: DelayShort,
    }
}

func (o *Overlay) SendMessage(msg string) {
    o.ch <- message{
        value: msg,
        style: theme.Info,
        delay: DelayLong,
    }
}

func (o *Overlay) Watch() {
    for m := range o.ch {
        o.buffer = &m

        time.Sleep(m.delay * time.Second)

        o.buffer = nil

        o.screen.PostEvent(tcell.NewEventInterrupt(nil))
    }
}

func (o *Overlay) Close() {
    close(o.ch)
}
