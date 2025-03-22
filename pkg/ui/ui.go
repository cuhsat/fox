package ui

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/gdamore/tcell/v2"
    "github.com/gdamore/tcell/v2/encoding"
    "github.com/mattn/go-runewidth"
)

const (
    Delta = 1
)

var (
    screen tcell.Screen
)

type UI struct {
    info   string
    input  *Input
    output *Output
}

func NewUI(theme string) *UI {
    encoding.Register()

    scr, err := tcell.NewScreen()

    if err != nil {
        fs.Panic(err)
    }

    err = scr.Init()

    if err != nil {
        fs.Panic(err)
    }

    screen = scr

    setTheme(theme)

    screen.HideCursor()
    screen.SetStyle(StyleOutput)

    return &UI{
        input:  NewInput(),
        output: NewOutput(),
    }
}

func (ui *UI) Loop(hs *fs.HeapSet) {
    for {
        heap := hs.Heap()
        w, h := ui.render(heap)

        ev := screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventResize:
            screen.Sync()

        case *tcell.EventError:
            fs.Error(ev.Error())

        case *tcell.EventKey:
            switch ev.Key() {
            case tcell.KeyCtrlQ, tcell.KeyEscape:
                return

            case tcell.KeyCtrlC:
                screen.SetClipboard(heap.Copy())

                ui.info = fmt.Sprintf("%s copied [%d]", heap.Path, len(heap.SMap))

            case tcell.KeyCtrlS:
                path := heap.Save()

                ui.info = fmt.Sprintf("%s saved [%d]", path, len(heap.SMap))

            case tcell.KeyHome:
                ui.output.ScrollBegin()

            case tcell.KeyEnd:
                ui.output.ScrollEnd(h-1)

            case tcell.KeyUp:
                if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollUp(h-1)
                } else {
                    ui.output.ScrollUp(Delta)
                }

            case tcell.KeyDown:
                if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollDown(h-1)
                } else {
                    ui.output.ScrollDown(Delta)
                }

            case tcell.KeyLeft:
                if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollLeft(w)
                } else {
                    ui.output.ScrollLeft(Delta)
                }

            case tcell.KeyRight:
                if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollRight(w)
                } else {
                    ui.output.ScrollRight(Delta)
                }

            case tcell.KeyPgUp:
                ui.output.ScrollPageUp(h-1)

            case tcell.KeyPgDn:
                ui.output.ScrollPageDown(h-1)

            case tcell.KeyEnter:
                v := ui.input.Accept()

                if len(v) > 0 {
                    ui.output.Reset()
                    heap.AddFilter(v)
                }

            case tcell.KeyTab:
                chain := heap.Chain

                ui.output.Reset()

                if ev.Modifiers() & tcell.ModShift == 1 {
                    heap = hs.Prev()
                } else {
                    heap = hs.Next()
                }

                heap.NoFilter()

                for _, f := range chain {
                    heap.AddFilter(f.Name)
                }

            case tcell.KeyBackspace, tcell.KeyBackspace2:
                if len(ui.input.Value) > 0 {
                    ui.input.DelRune()
                } else if len(heap.Chain) > 0 {
                    ui.output.Reset()
                    heap.DelFilter()
                }

            default:
                if ev.Rune() != 0 {
                    ui.input.AddRune(ev.Rune())
                }
            }

        case *tcell.EventInterrupt:
            ui.info = ""
        }
    }
}

func (ui *UI) Close() {
    r := recover()

    screen.Fini()

    if r != nil {
        fs.Panic(r)
    }
}

func (ui *UI) render(heap *fs.Heap) (w int, h int) {
    defer screen.Show()

    screen.Clear()

    w, h = screen.Size()

    ui.output.Render(heap, 0, 0, h-1)
    ui.input.Render(heap, 0, h-1, w)

    if len(ui.info) > 0 {
        ui.setInfo(0, h-1)
    }

    return
}

func length(s string) (l int) {
    for _, r := range s {
        l += runewidth.RuneWidth(r)
    }

    return
}

func print(x, y int, s string, sty tcell.Style) {
    for _, r := range s {
        if r == '\t' {
            r = tcell.RuneRArrow
        }

        screen.SetContent(x, y, r, nil, sty)
        x += runewidth.RuneWidth(r)
    }
}
