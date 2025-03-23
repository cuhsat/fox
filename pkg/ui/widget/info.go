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

    buffer chan string
    value  string
}

func NewInfo(screen tcell.Screen) *Info {
    return &Info{
        widget: widget{
            screen: screen,
        },
        buffer: make(chan string),
        value: "",
    }
}

func (i *Info) Render(x, y, w int) {
    if len(i.value) > 0 {
        i.print(x, y, fmt.Sprintf(" %-*s", w, i.value), theme.Info)
    }
}

func (i *Info) SendInfo(s string) {
    i.buffer <- s
}

func (i *Info) Watch() {
    for i.value = range i.buffer {
        time.Sleep(Delay * time.Second)

        i.value = ""

        i.screen.PostEvent(tcell.NewEventInterrupt(nil))
    }
}

func (i *Info) Close() {
    close(i.buffer)
}
