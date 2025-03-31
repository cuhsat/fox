package widget

import (
    "fmt"
    "time"

    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    OverlayDelay = 2 // seconds
)

type Overlay struct {
    widget

    ch chan message
    buffer *message
}

type message struct {
    value string
    style tcell.Style
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

func (o *Overlay) SendMessage(msg string) {
    o.ch <- message{
        value: msg,
        style: theme.Info,
    }
}

func (o *Overlay) Watch() {
    for m := range o.ch {
        o.buffer = &m

        time.Sleep(OverlayDelay * time.Second)

        o.buffer = nil

        o.screen.PostEvent(tcell.NewEventInterrupt(nil))
    }
}

func (o *Overlay) Close() {
    close(o.ch)
}
