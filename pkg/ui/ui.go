package ui

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui/comp"
    "github.com/cuhsat/cu/pkg/ui/themes"
    "github.com/gdamore/tcell/v2"
    "github.com/gdamore/tcell/v2/encoding"
)

const (
    Delta = 1 // lines
)

type UI struct {
    screen tcell.Screen
    output *comp.Output
    input  *comp.Input
    info   *comp.Info
}

func NewUI() *UI {
    encoding.Register()

    scr, err := tcell.NewScreen()

    if err != nil {
        fs.Panic(err)
    }

    err = scr.Init()

    if err != nil {
        fs.Panic(err)
    }

    themes.Load(themes.Default)

    scr.HideCursor()
    scr.SetStyle(themes.Output)

    return &UI{
        screen: scr,
        output: comp.NewOutput(scr),
        input:  comp.NewInput(scr),
        info:   comp.NewInfo(scr),
    }
}

func (ui *UI) Run(hs *data.HeapSet, hi *fs.History) {
    go ui.info.Watch()

    for {
        heap := hs.Heap()
        w, h := ui.render(heap)

        ui.setTitle(heap.Path)

        ev := ui.screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventInterrupt:
            continue

        case *tcell.EventResize:
            ui.screen.Sync()

        case *tcell.EventError:
            fs.Error(ev.Error())

        case *tcell.EventKey:
            switch ev.Key() {
            case tcell.KeyCtrlQ, tcell.KeyEscape:
                return

            case tcell.KeyCtrlC:
                ui.screen.SetClipboard(heap.Copy())

                ui.info.SendInfo(fmt.Sprintf("%s copied", heap.Path))

            case tcell.KeyCtrlS:
                path := heap.Save()
                
                ui.info.SendInfo(fmt.Sprintf("%s saved", path))

            case tcell.KeyHome:
                ui.output.ScrollBegin()

            case tcell.KeyEnd:
                ui.output.ScrollEnd(h-1)

            case tcell.KeyUp:
                if ev.Modifiers() & tcell.ModCtrl == 1 {
                    ui.input.Value = hi.PrevCommand()
                } else if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollUp(h-1)
                } else {
                    ui.output.ScrollUp(Delta)
                }

            case tcell.KeyDown:
                if ev.Modifiers() & tcell.ModCtrl == 1 {
                    ui.input.Value = hi.NextCommand()
                } else if ev.Modifiers() & tcell.ModShift == 1 {
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
                    hi.AddCommand(v)
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
        }
    }
}

func (ui *UI) Close() {
    r := recover()

    ui.info.Close()

    ui.screen.Fini()

    if r != nil {
        fs.Panic(r)
    }
}

func (ui *UI) setTitle(s string) {
    ui.screen.SetTitle(fmt.Sprintf("cu - %s", s))
}

func (ui *UI) render(heap *data.Heap) (w int, h int) {
    defer ui.screen.Show()

    ui.screen.Clear()

    w, h = ui.screen.Size()

    ui.output.Render(heap, 0, 0, h-1)
    ui.input.Render(heap, 0, h-1, w)
    ui.info.Render(0, h-1, w)

    return
}
