package ui

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui/theme"
    "github.com/cuhsat/cu/pkg/ui/widget"
    "github.com/gdamore/tcell/v2"
    "github.com/gdamore/tcell/v2/encoding"
)

const (
    Delta = 1 // lines
)

type UI struct {
    screen tcell.Screen
    output *widget.Output
    input  *widget.Input
    info   *widget.Info
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

    theme.Load(theme.Default)

    scr.SetStyle(theme.Output)
    scr.EnableMouse()
    scr.EnablePaste()
    scr.HideCursor()

    return &UI{
        screen: scr,
        output: widget.NewOutput(scr),
        input:  widget.NewInput(scr),
        info:   widget.NewInfo(scr),
    }
}

func (ui *UI) Run(hs *data.HeapSet, hi *fs.History) {
    go ui.info.Watch()

    for {
        heap := hs.Heap()
        w, h := ui.render(heap)

        ev := ui.screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventInterrupt:
            continue

        case *tcell.EventClipboard:
            ui.input.Value = string(ev.Data())

        case *tcell.EventResize:
            ui.screen.Sync()

        case *tcell.EventError:
            fs.Error(ev.Error())

        case *tcell.EventMouse:
            switch ev.Buttons() {
            case tcell.WheelUp:
                ui.output.ScrollUp(Delta)

            case tcell.WheelDown:
                ui.output.ScrollDown(Delta)

            case tcell.WheelLeft:
                ui.output.ScrollLeft(Delta)

            case tcell.WheelRight:
                ui.output.ScrollRight(Delta)
            }

        case *tcell.EventKey:
            page_w := w
            page_h := h - 1

            switch ev.Key() {
            case tcell.KeyCtrlQ, tcell.KeyEscape:
                return

            case tcell.KeyCtrlP:
                ui.screen.GetClipboard()

            case tcell.KeyCtrlC:
                ui.screen.SetClipboard(heap.Copy())

                ui.info.SendMessage(fmt.Sprintf("%s copied", heap.Path))

            case tcell.KeyCtrlS:
                path := heap.Save()
                
                ui.info.SendMessage(fmt.Sprintf("%s saved", path))

            case tcell.KeyCtrlN:
                ui.output.ToggleNumbers()

            case tcell.KeyCtrlW:
                ui.output.ToggleWrap()

            case tcell.KeyHome:
                ui.output.ScrollBegin()

            case tcell.KeyEnd:
                ui.output.ScrollEnd()

            case tcell.KeyUp:
                if ev.Modifiers() & tcell.ModCtrl == 1 {
                    ui.input.Value = hi.PrevCommand()
                } else if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollUp(page_h)
                } else {
                    ui.output.ScrollUp(Delta)
                }

            case tcell.KeyDown:
                if ev.Modifiers() & tcell.ModCtrl == 1 {
                    ui.input.Value = hi.NextCommand()
                } else if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollDown(page_h)
                } else {
                    ui.output.ScrollDown(Delta)
                }

            case tcell.KeyLeft:
                if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollLeft(page_w)
                } else {
                    ui.output.ScrollLeft(Delta)
                }

            case tcell.KeyRight:
                if ev.Modifiers() & tcell.ModShift == 1 {
                    ui.output.ScrollRight(page_w)
                } else {
                    ui.output.ScrollRight(Delta)
                }

            case tcell.KeyPgUp:
                ui.output.ScrollPageUp(page_h)

            case tcell.KeyPgDn:
                ui.output.ScrollPageDown(page_h)

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

func (ui *UI) render(heap *data.Heap) (w int, h int) {
    defer ui.screen.Show()

    ui.screen.SetTitle(fmt.Sprintf("cu - %s", heap.Path))
    ui.screen.Clear()

    w, h = ui.screen.Size()

    ui.output.Render(heap, 0, 0, w, h-1)
    ui.input.Render(heap, 0, h-1, w)
    ui.info.Render(0, h-1, w)

    return
}
