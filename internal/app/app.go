package app

import (
    "fmt"
    "strings"

    "github.com/cuhsat/fx/internal/app/themes"
    "github.com/cuhsat/fx/internal/app/widget"
    "github.com/cuhsat/fx/internal/sys"
    "github.com/cuhsat/fx/internal/sys/files/bag"
    "github.com/cuhsat/fx/internal/sys/files/history"
    "github.com/cuhsat/fx/internal/sys/heapset"
    "github.com/cuhsat/fx/internal/sys/types"
    "github.com/cuhsat/fx/internal/sys/types/mode"
    "github.com/gdamore/tcell/v2"
    "github.com/gdamore/tcell/v2/encoding"
)

const (
    Delta = 1 // lines
)

const (
    LBracket = "ESC[200~" // bracketed paste start
    RBracket = "ESC[201~" // bracketed paste end
)

type App struct {
    screen  tcell.Screen

    themes  *themes.Themes

    status  *widget.Status
    header  *widget.Header
    output  *widget.Output
    prompt  *widget.Prompt
    overlay *widget.Overlay
}

func NewApp(m mode.Mode) *App {
    encoding.Register()

    scr, err := tcell.NewScreen()

    if err != nil {
        sys.Fatal(err)
    }

    err = scr.Init()

    if err != nil {
        sys.Fatal(err)
    }

    scr.EnableMouse()
    scr.EnablePaste()
    scr.HideCursor()

    stt := widget.NewStatus()

    app := App{
        screen:  scr,
        status:  stt,
        themes:  themes.NewThemes(stt.Theme),
        header:  widget.NewHeader(scr, stt),
        output:  widget.NewOutput(scr, stt),
        prompt:  widget.NewPrompt(scr, stt),
        overlay: widget.NewOverlay(scr, stt),
    }

    app.State(m)

    return &app
}

func (app *App) Run(hs *heapset.HeapSet, hi *history.History, bag *bag.Bag) {
    hs.SetCallback(func() {
        app.screen.PostEvent(tcell.NewEventInterrupt(app.status.Follow))
    })

    go app.overlay.Watch()

    for {
        _, heap := hs.Current()

        w, h := app.render(hs)

        ev := app.screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventInterrupt:
            v, ok := ev.Data().(bool)

            if ok && v {
                app.output.ScrollEnd()
            }

            continue

        case *tcell.EventClipboard:
            if app.status.Mode == mode.Hex {
                continue
            }

            v := string(ev.Data())

            v = strings.TrimPrefix(v, LBracket)
            v = strings.TrimSuffix(v, RBracket)

            app.prompt.Value = v

        case *tcell.EventResize:
            app.screen.Sync()
            app.output.Reset()

        case *tcell.EventError:
            app.overlay.SendError(ev.Error())

        case *tcell.EventMouse:
            switch ev.Buttons() {
            case tcell.WheelUp:
                app.output.ScrollUp(Delta)

            case tcell.WheelDown:
                app.output.ScrollDown(Delta)

            case tcell.WheelLeft:
                app.output.ScrollLeft(Delta)

            case tcell.WheelRight:
                app.output.ScrollRight(Delta)

            case tcell.ButtonMiddle:
                app.screen.GetClipboard()
            }

        case *tcell.EventKey:
            mods := ev.Modifiers()

            page_w := w
            page_h := h-2

            switch ev.Key() {
            case tcell.KeyEscape:
                return

            case tcell.KeyCtrlL, tcell.KeyF1:
                heap.ClearFilters()

                app.State(mode.Less)

            case tcell.KeyCtrlG, tcell.KeyF2:
                app.State(mode.Grep)

            case tcell.KeyCtrlX, tcell.KeyF3:
                app.State(mode.Hex)

            case tcell.KeyCtrlSpace, tcell.KeyF4:
                app.State(mode.Goto)

            case tcell.KeyF9:
                hs.Word()

            case tcell.KeyF10:
                hs.Md5()

            case tcell.KeyF11:
                hs.Sha1()

            case tcell.KeyF12:
                hs.Sha256()

            case tcell.KeyCtrlV:
                if app.status.Mode == mode.Hex {
                    continue
                }

                app.screen.GetClipboard()

            case tcell.KeyCtrlC:
                if app.status.Mode == mode.Hex {
                    continue
                }

                app.screen.SetClipboard(heap.Bytes())

                app.overlay.SendStatus("Copied to clipboard")

            case tcell.KeyCtrlS:
                if app.status.Mode == mode.Hex {
                    continue
                }

                bag.Put(heap)

                app.overlay.SendStatus(fmt.Sprintf("Saved to %s", bag.Path))

            case tcell.KeyCtrlQ:
                app.output.Reset()

                heap = hs.CloseHeap()

                if heap == nil {
                    return // exit
                }

            case tcell.KeyCtrlR:
                app.output.Reset()
                
                heap.Reload()

            case tcell.KeyCtrlT:
                app.status.Theme = app.themes.Cycle()

                app.screen.Fill(' ', themes.Output)
                app.screen.Show()

                app.overlay.SendStatus(fmt.Sprintf("Theme %s", app.status.Theme))

            case tcell.KeyCtrlF:
                app.status.ToggleFollow()

            case tcell.KeyCtrlN:
                app.status.ToggleNumbers()

            case tcell.KeyCtrlW:
                app.status.ToggleWrap()

                app.output.Reset()

            case tcell.KeyHome:
                app.output.ScrollStart()

            case tcell.KeyEnd:
                app.output.ScrollEnd()

            case tcell.KeyUp:
                if mods & tcell.ModAlt != 0 {
                    app.prompt.Value = hi.PrevCommand()
                } else if mods & tcell.ModCtrl != 0 && mods & tcell.ModShift != 0 {
                    app.output.ScrollStart()
                } else if mods & tcell.ModShift != 0 {
                    app.output.ScrollUp(page_h)
                } else {
                    app.output.ScrollUp(Delta)
                }

            case tcell.KeyDown:
                if mods & tcell.ModAlt != 0 {
                    app.prompt.Value = hi.NextCommand()
                } else if mods & tcell.ModCtrl != 0 && mods & tcell.ModShift != 0 {
                    app.output.ScrollEnd()
                } else if mods & tcell.ModShift != 0 {
                    app.output.ScrollDown(page_h)
                } else {
                    app.output.ScrollDown(Delta)
                }

            case tcell.KeyLeft:
                if mods & tcell.ModShift != 0 {
                    app.output.ScrollLeft(page_w)
                } else {
                    app.output.ScrollLeft(Delta)
                }

            case tcell.KeyRight:
                if mods & tcell.ModShift != 0 {
                    app.output.ScrollRight(page_w)
                } else {
                    app.output.ScrollRight(Delta)
                }

            case tcell.KeyPgUp:
                app.output.ScrollUp(page_h)

            case tcell.KeyPgDn:
                app.output.ScrollDown(page_h)

            case tcell.KeyEnter:
                v := app.prompt.Accept()

                if len(v) == 0 {
                    continue
                }

                hi.AddCommand(v)

                switch app.status.Mode {
                case mode.Goto:
                    app.output.Goto(v)

                    app.State(app.status.Last)

                default:
                    app.output.Reset()
                
                    heap.AddFilter(v)
                }

            case tcell.KeyTab:
                app.output.Reset()

                if mods & tcell.ModShift != 0 {
                    heap = hs.PrevHeap()
                } else {
                    heap = hs.NextHeap()
                }

            case tcell.KeyBackspace2:
                if len(app.prompt.Value) > 0 {
                    app.prompt.DelRune()
                } else if app.status.Mode == mode.Goto {
                    app.State(app.status.Last)
                } else if len(*types.GetFilters()) > 0 {
                    app.output.Reset()
                    heap.DelFilter()
                } else if app.status.Mode == mode.Grep {
                    app.State(mode.Less)
                }

            default:
                r := ev.Rune()

                switch r {
                case 0: // error
                    continue

                case 32: // space
                    if app.status.Mode == mode.Less {
                        app.output.ScrollDown(page_h)
                    } else {
                        app.prompt.AddRune(r)
                    }

                default: // all other keys
                    if app.status.Mode == mode.Less {
                        app.State(mode.Grep)
                    }

                    app.prompt.AddRune(r)
                }
            }
        }
    }
}

func (app *App) State(m mode.Mode) {
    if !app.status.SwitchMode(m) {
        return
    }

    switch m {
    case mode.Less, mode.Hex: // static modes
        app.prompt.Lock = true

    case mode.Grep, mode.Goto: // interactive modes
        app.prompt.Lock = false
    }

    if app.status.Last == mode.Hex || m == mode.Hex {
        app.output.Reset()
    }
}

func (app *App) Close() {
    r := recover()

    defer app.status.Save()

    defer app.screen.Fini()

    defer app.overlay.Close()
    
    if r != nil {
        sys.Fatal(r)
    }
}

func (app *App) render(hs *heapset.HeapSet) (w int, h int) {
    defer app.screen.Show()

    _, heap := hs.Current()

    app.screen.SetTitle(fmt.Sprintf("fx - %s", heap))
    app.screen.SetStyle(themes.Output)
    app.screen.Clear()

    x, y := 0, 0
    w, h = app.screen.Size()

    for _, widget := range [...]widget.Queueable{
        app.header,
        app.output,
        app.prompt,
    }{
        y += widget.Render(hs, x, y, w, h-y)
    }

    app.overlay.Render(0, 0, w, h)

    return
}
