package ui

import (
    "strings"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/nsf/termbox-go"
)

const Delta = 1

const CommonFg = termbox.Attribute(248)
const CommonBg = termbox.Attribute(235)

const StatusFg = termbox.Attribute(248)
const StatusBg = termbox.Attribute(237)

const SearchFg = termbox.Attribute(248)
const SearchBg = termbox.Attribute(237)

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
    termbox.SetOutputMode(termbox.Output256)

    width, height = termbox.Size()

    return &UI{
        status: NewStatus(),
        buffer: NewBuffer(),
        search: NewSearch(),
    }
}

func (ui *UI) Loop(hs *fs.HeapSet) {
    ui.render(hs.Heap())

    for {
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            switch ev.Key {
            case termbox.KeyEsc:
                h := hs.Heap()

                if len(h.Chain) > 0 {
                    ui.buffer.Reset()
                    h.DelFilter()
                } else {
                    return
                }

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
                v := ui.search.GetValue()

                if len(v) > 0 {
                    ui.buffer.Reset()
                    hs.Heap().AddFilter(v)
                }

            case termbox.KeyTab:
                c := hs.Heap().Chain

                ui.buffer.Reset()
                hs.Cycle()

                h := hs.Heap()
                h.NoFilter()

                for _, l := range c {
                    h.AddFilter(l.Name)
                }

            case termbox.KeyBackspace, termbox.KeyBackspace2:
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

        ui.render(hs.Heap())
    }
}

func (ui *UI) Close() {
    termbox.Close()
}

func (ui *UI) render(h *fs.Heap) {
    termbox.Clear(CommonFg, CommonBg)
    termbox.HideCursor()

    width, height = termbox.Size()

    data = len(h.SMap)
    page = height - 2

    line := strings.Repeat(" ", width)

    z := height - 1

    printEx(0, 0, line, StatusFg, StatusBg)
    printEx(0, z, line, SearchFg, SearchBg)

    ui.status.Render(0, 0, h)
    ui.buffer.Render(0, 1, h)
    ui.search.Render(0, z)

    termbox.Flush()
}

func print(x, y int, s string) {
    for x, c := range s {
        termbox.SetChar(x, y, c)
    }
}

func printEx(x, y int, s string, fg, bg termbox.Attribute) {
    for x, c := range s {
        termbox.SetCell(x, y, c, fg, bg)
    }
}
