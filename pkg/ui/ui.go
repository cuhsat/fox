package ui

import (
    "fmt"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/config"
    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/fs/history"
    "github.com/cuhsat/cu/pkg/ui/mode"
    "github.com/cuhsat/cu/pkg/ui/status"
    "github.com/cuhsat/cu/pkg/ui/themes"
    "github.com/cuhsat/cu/pkg/ui/widget"
    "github.com/gdamore/tcell/v2"
    "github.com/gdamore/tcell/v2/encoding"
)

const (
    Delta = 1 // lines
)

type UI struct {
    screen  tcell.Screen

    status  *status.Status
    config  *config.Config
    themes  *themes.Themes

    header  *widget.Header
    output  *widget.Output
    input   *widget.Input
    overlay *widget.Overlay
}

func NewUI(m mode.Mode) *UI {
    encoding.Register()

    cfg := config.GetConfig()
    sts := status.NewStatus()

    scr, err := tcell.NewScreen()

    if err != nil {
        fs.Panic(err)
    }

    err = scr.Init()

    if err != nil {
        fs.Panic(err)
    }

    scr.EnableMouse()
    scr.EnablePaste()
    scr.HideCursor()

    ui := UI{
        screen:  scr,
        status:  sts,
        config:  cfg,
        themes:  themes.NewThemes(cfg.UI.Theme),
        header:  widget.NewHeader(scr, sts),
        output:  widget.NewOutput(scr, sts),
        input:   widget.NewInput(scr, sts),
        overlay: widget.NewOverlay(scr, sts),
    }

    ui.State(m)

    return &ui
}

func (ui *UI) Run(hs *heapset.HeapSet, hi *history.History) {
    hs.SetCallback(func() {
        ui.screen.PostEvent(tcell.NewEventInterrupt(ui.status.Follow))
    })

    go ui.overlay.Watch()

    for {
        _, heap := hs.Current()
        w, h := ui.render(hs)

        ev := ui.screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventInterrupt:
            v, ok := ev.Data().(bool)

            if ok && v {
                ui.output.ScrollEnd()
            }

            continue

        case *tcell.EventClipboard:
            ui.input.Value = string(ev.Data())

        case *tcell.EventResize:
            ui.screen.Sync()
            ui.output.Reset()

        case *tcell.EventError:
            ui.overlay.SendError(ev.Error())

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
            page_h := h-2

            switch ev.Key() {
            case tcell.KeyEscape:
                return

            case tcell.KeyCtrlL, tcell.KeyF1:
                ui.State(mode.Less)

            case tcell.KeyCtrlG, tcell.KeyF2:
                ui.State(mode.Grep)

            case tcell.KeyCtrlX, tcell.KeyF3:
                ui.State(mode.Hex)

            case tcell.KeyCtrlSpace, tcell.KeyF4:
                ui.State(mode.Goto)

            case tcell.KeyF9:
                hs.Word()

                ui.output.Reset()
                ui.State(mode.Less)

            case tcell.KeyF10:
                hs.Md5()

                ui.output.Reset()
                ui.State(mode.Less)

            case tcell.KeyF11:
                hs.Sha1()

                ui.output.Reset()
                ui.State(mode.Less)

            case tcell.KeyF12:
                hs.Sha256()

                ui.output.Reset()
                ui.State(mode.Less)

            case tcell.KeyCtrlV:
                if ui.status.Mode == mode.Hex {
                    continue
                }

                ui.screen.GetClipboard()

            case tcell.KeyCtrlC:
                if ui.status.Mode == mode.Hex {
                    continue
                }

                ui.screen.SetClipboard(heap.Copy())

                ui.overlay.SendStatus(fmt.Sprintf("%s copied", heap))

            case tcell.KeyCtrlS:
                if ui.status.Mode == mode.Hex {
                    continue
                }

                if heap.Flag == 0 && len(heap.Chain) == 0 {
                    continue
                }

                path := heap.Save()
                
                ui.overlay.SendStatus(fmt.Sprintf("%s saved", path))

            case tcell.KeyCtrlQ:
                chain := heap.Chain

                ui.output.Reset()

                heap = hs.CloseHeap()

                if heap == nil {
                    return // exit
                }

                heap.ResetFilter()

                for _, f := range chain {
                    heap.AddFilter(f.Name)
                }

            case tcell.KeyCtrlR:
                ui.output.Reset()
                
                heap.Reload()

            case tcell.KeyCtrlT:
                t := ui.themes.Cycle()

                ui.clear()

                ui.config.UI.Theme = t

                ui.overlay.SendStatus(fmt.Sprintf("Theme %s", t))

            case tcell.KeyCtrlF:
                ui.status.ToggleFollow()

            case tcell.KeyCtrlN:
                ui.status.ToggleNumbers()

            case tcell.KeyCtrlW:
                ui.status.ToggleWrap()

            case tcell.KeyHome:
                ui.output.ScrollBegin()

            case tcell.KeyEnd:
                ui.output.ScrollEnd()

            case tcell.KeyUp:
                if ev.Modifiers() & tcell.ModAlt != 0 {
                    ui.input.Value = hi.PrevCommand()
                } else if ev.Modifiers() & tcell.ModShift != 0 {
                    ui.output.ScrollUp(page_h)
                } else {
                    ui.output.ScrollUp(Delta)
                }

            case tcell.KeyDown:
                if ev.Modifiers() & tcell.ModAlt != 0 {
                    ui.input.Value = hi.NextCommand()
                } else if ev.Modifiers() & tcell.ModShift != 0 {
                    ui.output.ScrollDown(page_h)
                } else {
                    ui.output.ScrollDown(Delta)
                }

            case tcell.KeyLeft:
                if ev.Modifiers() & tcell.ModShift != 0 {
                    ui.output.ScrollLeft(page_w)
                } else {
                    ui.output.ScrollLeft(Delta)
                }

            case tcell.KeyRight:
                if ev.Modifiers() & tcell.ModShift != 0 {
                    ui.output.ScrollRight(page_w)
                } else {
                    ui.output.ScrollRight(Delta)
                }

            case tcell.KeyPgUp:
                ui.output.ScrollUp(page_h)

            case tcell.KeyPgDn:
                ui.output.ScrollDown(page_h)

            case tcell.KeyEnter:
                v := ui.input.Accept()

                if len(v) == 0 {
                    continue
                }

                hi.AddCommand(v)

                switch ui.status.Mode {
                case mode.Goto:
                    ui.output.Goto(v)

                    ui.State(ui.status.Last)

                default:
                    ui.output.Reset()
                
                    heap.AddFilter(v)
                }

            case tcell.KeyTab:
                chain := heap.Chain

                ui.output.Reset()

                if ev.Modifiers() & tcell.ModShift != 0 {
                    heap = hs.PrevHeap()
                } else {
                    heap = hs.NextHeap()
                }

                heap.ResetFilter()

                for _, f := range chain {
                    heap.AddFilter(f.Name)
                }

            case tcell.KeyBackspace2:
                if len(ui.input.Value) > 0 {
                    ui.input.DelRune()
                } else if ui.status.Mode == mode.Goto {
                    ui.State(ui.status.Last)
                } else if len(heap.Chain) > 0 {
                    ui.output.Reset()
                    heap.DelFilter()
                } else if ui.status.Mode == mode.Grep {
                    ui.State(mode.Less)
                }

            default:
                r := ev.Rune()

                switch r {
                case 0: // error
                    continue

                case 32: // space
                    if ui.status.Mode == mode.Less {
                        ui.output.ScrollDown(page_h)                        
                    }

                default: // all other keys
                    if ui.status.Mode == mode.Less {
                        ui.State(mode.Grep)                        
                    }

                    ui.input.AddRune(r)
                }
            }
        }
    }
}

func (ui *UI) State(m mode.Mode) {
    if !ui.status.SwitchMode(m) {
        return
    }

    switch m {
    case mode.Less, mode.Hex: // static modes
        ui.input.Lock = true

    case mode.Grep, mode.Goto: // interactive modes
        ui.input.Lock = false
    }

    if ui.status.Last == mode.Hex || m == mode.Hex {
        ui.output.Reset()
    }
}

func (ui *UI) Close() {
    r := recover()

    defer ui.config.Save()

    defer ui.screen.Fini()

    defer ui.overlay.Close()
    
    if r != nil {
        fs.Panic(r)
    }
}

func (ui *UI) clear() {
    ui.screen.Fill(' ', themes.Output)
    ui.screen.Show()
}

func (ui *UI) render(hs *heapset.HeapSet) (w int, h int) {
    defer ui.screen.Show()

    _, heap := hs.Current()

    ui.screen.SetTitle(fmt.Sprintf("cu - %s", heap))
    ui.screen.SetStyle(themes.Output)
    ui.screen.Clear()

    x, y := 0, 0
    w, h = ui.screen.Size()

    for _, widget := range [...]widget.Queueable{
        ui.header,
        ui.output,
        ui.input,
    }{
        y += widget.Render(hs, x, y, w, h-y)
    }

    ui.overlay.Render(0, 0, w, h)

    return
}
