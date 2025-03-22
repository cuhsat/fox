package ui

import (
    "fmt"
    "time"

    "github.com/gdamore/tcell/v2"
)

const (
    Delay = 3
)

func (ui *UI) setInfo(x, y int) {
    w, _ := screen.Size()

    print(x, y, fmt.Sprintf(" %-*s", w, ui.info), StyleInfo)

    go func() {
        time.Sleep(Delay * time.Second)
        
        screen.PostEvent(tcell.NewEventInterrupt(""))
    }()
}
