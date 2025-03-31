package widget

import (
    "fmt"
    "time"

    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/gdamore/tcell/v2"
)

const (
    Delay = 2 // seconds
)

type Info struct {
    widget

    ch chan message
    buffer *message
}

type message struct {
    value string
    style tcell.Style
} 

func NewInfo(screen tcell.Screen) *Info {
    return &Info{
        widget: widget{
            screen: screen,
        },
        ch: make(chan message),
        buffer: nil,
    }
}

func (i *Info) Render(x, y, w int) {
    if i.buffer != nil {
        i.print(x, y, fmt.Sprintf(" %-*s", w, i.buffer.value), i.buffer.style)
    }
}

func (i *Info) SendMessage(s string) {
    i.ch <- message{
        value: s,
        style: theme.Info,
    }
}

func (i *Info) Watch() {
    for m := range i.ch {
        i.buffer = &m

        time.Sleep(Delay * time.Second)

        i.buffer = nil

        i.screen.PostEvent(tcell.NewEventInterrupt(nil))
    }
}

func (i *Info) Close() {
    close(i.ch)
}
