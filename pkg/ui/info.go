package ui

import (
    "fmt"
    "time"

    "github.com/gdamore/tcell/v2"
)

const (
    Delay = 2 // seconds
)

type Info struct {
    buffer chan string
    value  string
}

func NewInfo() *Info {
    return &Info{
        buffer: make(chan string),
        value: "",
    }
}

func (i *Info) Render(x, y, w int) {
    if len(i.value) > 0 {
        print(x, y, fmt.Sprintf(" %-*s", w, i.value), StyleInfo)
    }
}

func (i *Info) SendInfo(s string) {
    i.buffer <- s
}

func (i *Info) Watch() {
    for i.value = range i.buffer {
        time.Sleep(Delay * time.Second)

        i.value = ""

        screen.PostEvent(tcell.NewEventInterrupt(nil))
    }
}

func (i *Info) Close() {
    close(i.buffer)
}
