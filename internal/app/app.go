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
    delta = 1 // lines
)

const (
    bracketL = "ESC[200~" // bracketed paste start
    bracketR = "ESC[201~" // bracketed paste end
)

type App struct {
    ctx *widget.Context

    term tcell.Screen

    themes *themes.Themes

    header  *widget.Header
    output  *widget.Output
    prompt  *widget.Prompt
    overlay *widget.Overlay
}

func New(m mode.Mode) *App {
    encoding.Register()

    term, err := tcell.NewScreen()

    if err != nil {
        sys.Fatal(err)
    }

    err = term.Init()

    if err != nil {
        sys.Fatal(err)
    }

    term.EnableMouse()
    term.EnablePaste()
    term.HideCursor()

    ctx := widget.NewContext()

    app := App{
        ctx: ctx,

        term: term,

        themes: themes.New(ctx.Theme),

        header:  widget.NewHeader(ctx, term),
        output:  widget.NewOutput(ctx, term),
        prompt:  widget.NewPrompt(ctx, term),
        overlay: widget.NewOverlay(ctx, term),
    }

    app.State(m)

    return &app
}

func (app *App) Run(hs *heapset.HeapSet, hi *history.History, bag *bag.Bag) {
    hs.SetCallback(func() {
        app.term.PostEvent(tcell.NewEventInterrupt(app.ctx.Follow))
    })

    go app.overlay.Watch()

    for {
        _, heap := hs.Current()

        w, h := app.render(hs)

        ev := app.term.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventInterrupt:
            v, ok := ev.Data().(bool)

            if ok && v {
                app.output.ScrollEnd()
            }

            continue

        case *tcell.EventClipboard:
            if app.ctx.Mode == mode.Hex {
                continue
            }

            v := string(ev.Data())

            v = strings.TrimPrefix(v, bracketL)
            v = strings.TrimSuffix(v, bracketR)

            app.prompt.Value = v

        case *tcell.EventResize:
            app.term.Sync()
            app.output.Reset()

        case *tcell.EventError:
            app.overlay.SendError(ev.Error())

        case *tcell.EventMouse:
            switch ev.Buttons() {
            case tcell.WheelUp:
                app.output.ScrollUp(delta)

            case tcell.WheelDown:
                app.output.ScrollDown(delta)

            case tcell.WheelLeft:
                app.output.ScrollLeft(delta)

            case tcell.WheelRight:
                app.output.ScrollRight(delta)

            case tcell.ButtonMiddle:
                app.term.GetClipboard()
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
                if app.ctx.Mode == mode.Hex {
                    continue
                }

                app.term.GetClipboard()

            case tcell.KeyCtrlC:
                if app.ctx.Mode == mode.Hex {
                    continue
                }

                app.term.SetClipboard(heap.Bytes())

                app.overlay.SendInfo("Copied to clipboard")

            case tcell.KeyCtrlS:
                if app.ctx.Mode == mode.Hex {
                    continue
                }

                bag.Put(heap)

                app.overlay.SendInfo(fmt.Sprintf("Saved to %s", bag.Path))

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
                app.ctx.Theme = app.themes.Cycle()

                app.term.Fill(' ', themes.Base)
                app.term.Show()

                app.overlay.SendInfo(fmt.Sprintf("Theme %s", app.ctx.Theme))

            case tcell.KeyCtrlF:
                app.ctx.ToggleFollow()

            case tcell.KeyCtrlN:
                app.ctx.ToggleNumbers()

            case tcell.KeyCtrlW:
                app.ctx.ToggleWrap()

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
                    app.output.ScrollUp(delta)
                }

            case tcell.KeyDown:
                if mods & tcell.ModAlt != 0 {
                    app.prompt.Value = hi.NextCommand()
                } else if mods & tcell.ModCtrl != 0 && mods & tcell.ModShift != 0 {
                    app.output.ScrollEnd()
                } else if mods & tcell.ModShift != 0 {
                    app.output.ScrollDown(page_h)
                } else {
                    app.output.ScrollDown(delta)
                }

            case tcell.KeyLeft:
                if mods & tcell.ModShift != 0 {
                    app.output.ScrollLeft(page_w)
                } else {
                    app.output.ScrollLeft(delta)
                }

            case tcell.KeyRight:
                if mods & tcell.ModShift != 0 {
                    app.output.ScrollRight(page_w)
                } else {
                    app.output.ScrollRight(delta)
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

                switch app.ctx.Mode {
                case mode.Goto:
                    app.output.Goto(v)

                    app.State(app.ctx.Last)

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
                } else if app.ctx.Mode == mode.Goto {
                    app.State(app.ctx.Last)
                } else if len(*types.GetFilters()) > 0 {
                    app.output.Reset()
                    heap.DelFilter()
                } else if app.ctx.Mode == mode.Grep {
                    app.State(mode.Less)
                }

            default:
                r := ev.Rune()

                switch r {
                case 0: // error
                    continue

                case 32: // space
                    if app.ctx.Mode == mode.Less {
                        app.output.ScrollDown(page_h)
                    } else {
                        app.prompt.AddRune(r)
                    }

                default: // all other keys
                    if app.ctx.Mode == mode.Less {
                        app.State(mode.Grep)
                    }

                    app.prompt.AddRune(r)
                }
            }
        }
    }
}

func (app *App) State(m mode.Mode) {
    if !app.ctx.SwitchMode(m) {
        return
    }

    switch m {
    case mode.Less, mode.Hex: // static modes
        app.prompt.Lock = true

    case mode.Grep, mode.Goto: // interactive modes
        app.prompt.Lock = false
    }

    if app.ctx.Last == mode.Hex || m == mode.Hex {
        app.output.Reset()
    }
}

func (app *App) Close() {
    r := recover()

    defer app.ctx.Save()

    defer app.term.Fini()

    defer app.overlay.Close()
    
    if r != nil {
        sys.Fatal(r)
    }
}

func (app *App) render(hs *heapset.HeapSet) (w int, h int) {
    defer app.term.Show()

    _, heap := hs.Current()

    app.term.SetTitle(fmt.Sprintf("Forensic Examiner - %s", heap))
    app.term.SetStyle(themes.Base)
    app.term.Clear()

    x, y := 0, 0
    w, h = app.term.Size()

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
