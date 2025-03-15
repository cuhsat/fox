package ui

import (
    "github.com/cuhsat/cu/pkg/fs"
    "github.com/nsf/termbox-go"
)

const Delta = 1

const CClear  = termbox.ColorDefault
const CStatus = termbox.ColorWhite | termbox.AttrBold
const CSearch = termbox.ColorLightGreen | termbox.AttrBold

var width, height, data, page int

type UI struct {
    status *Status
    buffer *Buffer
    search *Search
}

func NewUI() *UI {
    err := termbox.Init()
    
    if err != nil {
        fs.Panic(err)
    }
    
    termbox.SetInputMode(termbox.InputEsc)

    width, height = termbox.Size()

    return &UI{
        status: NewStatus(),
        buffer: NewBuffer(),
        search: NewSearch(),
    }
}

func (ui *UI) Render(heap *fs.Heap) {
    termbox.Clear(CClear, CClear)

    width, height = termbox.Size()

    data = len(heap.SMap)
    page = height - 2

    ui.status.Render(0, 0, heap)
    ui.buffer.Render(0, 1, heap)
    ui.search.Render(0, height - 1)

    termbox.Flush()
}

func (ui *UI) Loop(heap *fs.Heap) {
    ui.Render(heap)

    for {
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            switch ev.Key {
            case termbox.KeyEsc:
                return

            case termbox.KeyHome:
                ui.buffer.GoToBegin()

            case termbox.KeyEnd:
                ui.buffer.GoToEnd()

            case termbox.KeyPgup:
                ui.buffer.PageUp()

            case termbox.KeyPgdn:
                ui.buffer.PageDown()

            case termbox.KeyArrowUp:
                ui.buffer.ScrollUp(Delta)
                
            case termbox.KeyArrowDown:
                ui.buffer.ScrollDown(Delta)

            case termbox.KeyArrowLeft:
                ui.buffer.ScrollLeft(Delta)

            case termbox.KeyArrowRight:
                ui.buffer.ScrollRight(Delta)

            case termbox.KeyEnter:
                value := ui.search.GetValue()

                if len(value) > 0 {
                    ui.buffer.Reset()
                    heap.AddFilter(value)
                }

            case termbox.KeyTab:
                ui.buffer.Reset()
                heap.DelFilter()

            case termbox.KeyBackspace2:
            case termbox.KeyBackspace:
                ui.search.DelChar()

            case termbox.KeySpace:
                ui.search.AddChar(' ')

            default:
                if ev.Ch != 0 {
                    ui.search.AddChar(ev.Ch)
                }
            }

        case termbox.EventError:
            fs.Error(ev.Err)
        }

        ui.Render(heap)
    }
}

func (ui *UI) Close() {
    termbox.Close()
}

func printEx(x, y int, s string, fg, bg termbox.Attribute) {
    for x, c := range s {
        termbox.SetCell(x, y, c, fg, bg)
    }
}

func print(x, y int, s string) {
    for x, c := range s {
        termbox.SetChar(x, y, c)
    }
}
