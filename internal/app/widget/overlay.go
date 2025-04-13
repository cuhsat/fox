package widget

import (
    "fmt"
    "time"

    "github.com/cuhsat/cu/internal/app/themes"
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

func NewOverlay(screen tcell.Screen, status *Status) *Overlay {
    return &Overlay{
        widget: widget{
            screen: screen,
            status: status,
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
        style: themes.Error,
        delay: DelayLong,
    }
}

func (o *Overlay) SendStatus(msg string) {
    o.ch <- message{
        value: msg,
        style: themes.Info,
        delay: DelayShort,
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
