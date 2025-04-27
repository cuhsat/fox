package ui

import (
    "fmt"
    "strings"

    "github.com/cuhsat/fx/pkg/fx"
    "github.com/cuhsat/fx/pkg/fx/sys"
    "github.com/cuhsat/fx/pkg/fx/text"
    "github.com/cuhsat/fx/pkg/fx/types"
    "github.com/cuhsat/fx/pkg/fx/types/buffer"
    "github.com/cuhsat/fx/pkg/fx/types/heapset"
    "github.com/cuhsat/fx/pkg/fx/types/mode"
    "github.com/cuhsat/fx/pkg/fx/user/bag"
    "github.com/cuhsat/fx/pkg/fx/user/history"
    "github.com/cuhsat/fx/pkg/ui/context"
    "github.com/cuhsat/fx/pkg/ui/themes"
    "github.com/cuhsat/fx/pkg/ui/widgets"
    "github.com/gdamore/tcell/v2"
    "github.com/gdamore/tcell/v2/encoding"
    "github.com/mattn/go-runewidth"
)

const (
    delta = 1 // lines
)

const (
    bracketPrefix = "ESC[200~" // bracketed paste start
    bracketSuffix = "ESC[201~" // bracketed paste end
)

const (
    cursor = tcell.CursorStyleBlinkingBar // cursor style
)

type UI struct {
    ctx *context.Context

    term tcell.Screen

    themes *themes.Themes

    title  *widgets.Title
    view   *widgets.View
    status *widgets.Status
    
    overlay *widgets.Overlay
}

func New(m mode.Mode) *UI {
    encoding.Register()

    runewidth.CreateLUT()

    term, err := tcell.NewScreen()

    if err != nil {
        sys.Panic(err)
    }

    err = term.Init()

    if err != nil {
        sys.Panic(err)
    }

    term.EnableMouse()
    term.EnablePaste()

    term.HideCursor()

    ctx := context.New()

    ui := UI{
        ctx: ctx,

        term: term,

        themes: themes.New(ctx.Theme),

        title:   widgets.NewTitle(ctx, term),
        view:    widgets.NewView(ctx, term),
        status:  widgets.NewStatus(ctx, term),
        overlay: widgets.NewOverlay(ctx, term),
    }

    term.SetCursorStyle(cursor, themes.Cursor)
    term.SetStyle(themes.Base)
    term.Sync()

    ui.state(m)

    return &ui
}

func (ui *UI) Run(hs *heapset.HeapSet, hi *history.History, bag *bag.Bag) {
    hs.Bind(func() {
        ui.term.PostEvent(tcell.NewEventInterrupt(ui.ctx.Follow))
    }, func() {
        ui.term.PostEvent(tcell.NewEventError(nil))
    })

    events, quit := make(chan tcell.Event, 16), make(chan struct{})

    go ui.term.ChannelEvents(events, quit)

    go ui.overlay.Listen()

    for {
        select {
        case _ = <-quit:
            return // channels closed

        case ev := <-events:
            if ev == nil {
                return // term closed
            }

            w, h := ui.term.Size()

            _, heap := hs.Current()

            switch ev := ev.(type) {
            case *tcell.EventInterrupt:
                v, ok := ev.Data().(bool)

                if ok && v {
                    ui.view.ScrollEnd()
                }

            case *tcell.EventClipboard:
                if ui.ctx.Mode == mode.Hex {
                    continue
                }

                v := string(ev.Data())

                v = strings.TrimPrefix(v, bracketPrefix)
                v = strings.TrimSuffix(v, bracketSuffix)

                ui.status.Value = v

            case *tcell.EventResize:
                ui.term.Sync()
                ui.view.Reset()

            case *tcell.EventError:
                hs.OpenLog()

                ui.overlay.SendError("An error occurred")

            case *tcell.EventMouse:
                btns := ev.Buttons()

                if btns & tcell.ButtonMiddle != 0 {
                    ui.term.GetClipboard()
                } else if btns & tcell.WheelUp != 0 {
                    ui.view.ScrollUp(delta)
                } else if btns & tcell.WheelDown != 0 {
                    ui.view.ScrollDown(delta)
                } else if btns & tcell.WheelLeft != 0 {
                    ui.view.ScrollLeft(delta)
                } else if btns & tcell.WheelRight != 0 {
                    ui.term.GetClipboard()
                }

            case *tcell.EventKey:
                mods := ev.Modifiers()

                page_w := w-1 // minus text abbreviation
                page_h := h-2 // minus title and status

                if ui.ctx.Line {
                    page_w -= text.Dec(heap.Length()) + buffer.SpaceText
                }

                switch ev.Key() {
                case tcell.KeyEscape:
                    return

                case tcell.KeyCtrlL, tcell.KeyF1:
                    ui.state(mode.Less)

                case tcell.KeyCtrlG, tcell.KeyF2:
                    ui.state(mode.Grep)

                case tcell.KeyCtrlX, tcell.KeyF3:
                    ui.state(mode.Hex)

                case tcell.KeyCtrlSpace, tcell.KeyF4:
                    ui.state(mode.Goto)

                case tcell.KeyCtrlO, tcell.KeyF5:
                    ui.state(mode.Open)

                case tcell.KeyF9:
                    hs.Stats()

                case tcell.KeyF10:
                    hs.Md5()

                case tcell.KeyF11:
                    hs.Sha1()

                case tcell.KeyF12:
                    hs.Sha256()

                case tcell.KeyCtrlV:
                    if ui.ctx.Mode == mode.Hex {
                        continue
                    }

                    ui.term.GetClipboard()

                case tcell.KeyCtrlC:
                    if ui.ctx.Mode == mode.Hex {
                        continue
                    }

                    ui.term.SetClipboard(heap.Bytes())

                    ui.overlay.SendInfo("Copied to clipboard")

                case tcell.KeyCtrlS, tcell.KeyPrint:
                    if ui.ctx.Mode == mode.Hex {
                        continue
                    }

                    if !bag.Put(heap) {
                        continue
                    }

                    ui.overlay.SendInfo(fmt.Sprintf("Saved to %s", bag.Path))

                case tcell.KeyCtrlE:
                    if sys.Exists(bag.Path) {
                        hs.OpenHeap(bag.Path)
                    } else {
                        ui.overlay.SendError(fmt.Sprintf("%s not found", bag.Path))
                    }

                case tcell.KeyCtrlD:
                    hs.OpenLog()

                case tcell.KeyCtrlQ:
                    ui.view.Reset()

                    heap = hs.CloseHeap()

                    if heap == nil {
                        return // exit
                    }

                case tcell.KeyCtrlZ:
                    err := ui.term.Suspend()

                    if err != nil {
                        sys.Error(err)
                        continue
                    }

                    sys.Shell()

                    err = ui.term.Resume()

                    if err != nil {
                        sys.Panic(err)
                    }

                case tcell.KeyCtrlT:
                    ui.ctx.Theme = ui.themes.Cycle()

                    ui.term.Fill(' ', themes.Base)
                    ui.term.Show()

                    ui.term.SetCursorStyle(cursor, themes.Cursor)

                    ui.overlay.SendInfo(fmt.Sprintf("Theme %s", ui.ctx.Theme))

                case tcell.KeyCtrlF:
                    if ui.ctx.Mode != mode.Hex {
                        ui.ctx.ToggleFollow()
                    }

                case tcell.KeyCtrlN:
                    if ui.ctx.Mode != mode.Hex {
                        ui.ctx.ToggleNumbers()
                    }

                case tcell.KeyCtrlW:
                    if ui.ctx.Mode != mode.Hex {
                        ui.ctx.ToggleWrap()
                        ui.view.Reset()
                    }

                case tcell.KeyHome:
                    ui.view.ScrollStart()

                case tcell.KeyEnd:
                    ui.view.ScrollEnd()

                case tcell.KeyUp:
                    if mods & tcell.ModAlt != 0 {
                        ui.status.Value = hi.PrevCommand()
                    } else if mods & tcell.ModCtrl != 0 && mods & tcell.ModShift != 0 {
                        ui.view.ScrollStart()
                    } else if mods & tcell.ModShift != 0 {
                        ui.view.ScrollUp(page_h)
                    } else {
                        ui.view.ScrollUp(delta)
                    }

                case tcell.KeyDown:
                    if mods & tcell.ModAlt != 0 {
                        ui.status.Value = hi.NextCommand()
                    } else if mods & tcell.ModCtrl != 0 && mods & tcell.ModShift != 0 {
                        ui.view.ScrollEnd()
                    } else if mods & tcell.ModShift != 0 {
                        ui.view.ScrollDown(page_h)
                    } else {
                        ui.view.ScrollDown(delta)
                    }

                case tcell.KeyLeft:
                    if mods & tcell.ModShift != 0 {
                        ui.view.ScrollLeft(page_w)
                    } else {
                        ui.view.ScrollLeft(delta)
                    }

                case tcell.KeyRight:
                    if mods & tcell.ModShift != 0 {
                        ui.view.ScrollRight(page_w)
                    } else {
                        ui.view.ScrollRight(delta)
                    }

                case tcell.KeyPgUp:
                    ui.view.ScrollUp(page_h)

                case tcell.KeyPgDn:
                    ui.view.ScrollDown(page_h)

                case tcell.KeyEnter:
                    v := ui.status.Accept()

                    if len(v) == 0 {
                        continue
                    }

                    hi.AddCommand(v)

                    switch ui.ctx.Mode {
                    case mode.Goto:
                        ui.view.Goto(v)
                        ui.state(ui.ctx.Last)

                    case mode.Open:
                        hs.Open(v)
                        ui.state(ui.ctx.Last)

                    default:
                        ui.view.Reset()
                        heap.AddFilter(v)
                        ui.state(mode.Less)
                    }

                case tcell.KeyTab:
                    ui.view.Reset()

                    if mods & tcell.ModShift != 0 {
                        heap = hs.PrevHeap()
                    } else {
                        heap = hs.NextHeap()
                    }

                case tcell.KeyBackspace2:
                    if len(ui.status.Value) > 0 {
                        ui.status.DelRune()
                    } else if ui.ctx.Mode == mode.Goto {
                        ui.state(ui.ctx.Last)
                    } else if ui.ctx.Mode == mode.Open {
                        ui.state(ui.ctx.Last)
                    } else if len(*types.Filters()) > 0 {
                        ui.view.Reset()
                        heap.DelFilter()
                    } else if ui.ctx.Mode == mode.Grep {
                        ui.state(mode.Less)
                    }

                default:
                    r := ev.Rune()

                    switch r {
                    case 0: // error
                        continue

                    case 32: // space
                        if ui.status.Lock {
                            ui.view.ScrollDown(page_h)
                        } else {
                            ui.status.AddRune(r)
                        }

                    default: // all other keys
                        if ui.ctx.Mode == mode.Less {
                            ui.state(mode.Grep)
                        }

                        ui.status.AddRune(r)
                    }
                }
            }

            ui.render(hs)
        }
    }
}

func (ui *UI) Close() {
    ui.overlay.Close()
    ui.term.Fini()
    ui.ctx.Save()
}

func (ui *UI) state(m mode.Mode) {
    if !ui.ctx.SwitchMode(m) {
        return
    }

    switch m {
    // static modes
    case mode.Less, mode.Hex:
        ui.status.Lock = true

    // input modes
    case mode.Grep, mode.Goto, mode.Open:
        ui.status.Lock = false
    }

    if ui.ctx.Last == mode.Hex || m == mode.Hex {
        ui.view.Reset()
    }
}

func (ui *UI) render(hs *heapset.HeapSet) {
    defer ui.term.Show()

    _, heap := hs.Current()

    if heap.Type == types.Stdin {
        ui.term.Sync() // prevent hickups
    }

    ui.term.SetTitle(fmt.Sprintf("%s - %s", fx.Product, heap))
    ui.term.SetStyle(themes.Base)
    ui.term.Clear()

    x, y := 0, 0
    w, h := ui.term.Size()

    for _, base := range [...]widgets.Queueable{
        ui.title,
        ui.view,
        ui.status,
    } {
        y += base.Render(hs, x, y, w, h-y)
    }

    ui.overlay.Render(0, 0, w, h)
}
