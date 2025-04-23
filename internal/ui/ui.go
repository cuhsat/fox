package ui

import (
    "fmt"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heapset"
    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/mode"
    "github.com/cuhsat/fx/internal/fx/user/bag"
    "github.com/cuhsat/fx/internal/fx/user/history"
    "github.com/cuhsat/fx/internal/ui/lib"
    "github.com/cuhsat/fx/internal/ui/themes"
    "github.com/gdamore/tcell/v2"
    "github.com/gdamore/tcell/v2/encoding"
    "github.com/mattn/go-runewidth"
)

const (
    delta = 1 // lines
)

const (
    bracketL = "ESC[200~" // bracketed paste start
    bracketR = "ESC[201~" // bracketed paste end
)

type UI struct {
    ctx *lib.Context

    term tcell.Screen

    themes *themes.Themes

    title  *lib.Title
    buffer *lib.Buffer
    prompt *lib.Prompt
    
    overlay *lib.Overlay
}

func New(m mode.Mode) *UI {
    encoding.Register()

    runewidth.CreateLUT()

    term, err := tcell.NewScreen()

    if err != nil {
        fx.Panic(err)
    }

    err = term.Init()

    if err != nil {
        fx.Panic(err)
    }

    term.EnableMouse()
    term.EnablePaste()

    term.HideCursor()
    term.SetCursorStyle(tcell.CursorStyleBlinkingBar)

    ctx := lib.NewContext()

    ui := UI{
        ctx: ctx,

        term: term,

        themes: themes.New(ctx.Theme),

        title:   lib.NewTitle(ctx, term),
        buffer:  lib.NewBuffer(ctx, term),
        prompt:  lib.NewPrompt(ctx, term),
        overlay: lib.NewOverlay(ctx, term),
    }

    ui.State(m)

    return &ui
}

func (ui *UI) Run(hs *heapset.HeapSet, hi *history.History, bag *bag.Bag) {
    hs.Bind(func() {
        ui.term.PostEvent(tcell.NewEventInterrupt(ui.ctx.Follow))
    }, func() {
        ui.term.PostEvent(tcell.NewEventError(nil))
    })

    go ui.overlay.Watch()

    for {
        _, heap := hs.Current()

        w, h := ui.render(hs)

        ev := ui.term.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventInterrupt:
            v, ok := ev.Data().(bool)

            if ok && v {
                ui.buffer.ScrollEnd()
            }

            continue

        case *tcell.EventClipboard:
            if ui.ctx.Mode == mode.Hex {
                continue
            }

            v := string(ev.Data())

            v = strings.TrimPrefix(v, bracketL)
            v = strings.TrimSuffix(v, bracketR)

            ui.prompt.Value = v

        case *tcell.EventResize:
            ui.term.Sync()
            ui.buffer.Reset()

        case *tcell.EventError:
            hs.OpenLog()

            ui.overlay.SendError("An error occurred")

        case *tcell.EventMouse:
            switch ev.Buttons() {
            case tcell.WheelUp:
                ui.buffer.ScrollUp(delta)

            case tcell.WheelDown:
                ui.buffer.ScrollDown(delta)

            case tcell.WheelLeft:
                ui.buffer.ScrollLeft(delta)

            case tcell.WheelRight:
                ui.buffer.ScrollRight(delta)

            case tcell.ButtonMiddle:
                ui.term.GetClipboard()
            }

        case *tcell.EventKey:
            mods := ev.Modifiers()

            page_w := w-1
            page_h := h-2

            if ui.ctx.Line {
                page_w -= text.Dec(heap.Length())+1
            }

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
                hs.OpenHeap(bag.Path)

            case tcell.KeyCtrlD:
                hs.OpenLog()

            case tcell.KeyCtrlR:
                ui.buffer.Reset()

                heap.Reload()

            case tcell.KeyCtrlQ:
                ui.buffer.Reset()

                heap = hs.CloseHeap()

                if heap == nil {
                    return // exit
                }

            case tcell.KeyCtrlT:
                ui.ctx.Theme = ui.themes.Cycle()

                ui.term.Fill(' ', themes.Base)
                ui.term.Show()

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
                    ui.buffer.Reset()
                }

            case tcell.KeyHome:
                ui.buffer.ScrollStart()

            case tcell.KeyEnd:
                ui.buffer.ScrollEnd()

            case tcell.KeyUp:
                if mods & tcell.ModAlt != 0 {
                    ui.prompt.Value = hi.PrevCommand()
                } else if mods & tcell.ModCtrl != 0 && mods & tcell.ModShift != 0 {
                    ui.buffer.ScrollStart()
                } else if mods & tcell.ModShift != 0 {
                    ui.buffer.ScrollUp(page_h)
                } else {
                    ui.buffer.ScrollUp(delta)
                }

            case tcell.KeyDown:
                if mods & tcell.ModAlt != 0 {
                    ui.prompt.Value = hi.NextCommand()
                } else if mods & tcell.ModCtrl != 0 && mods & tcell.ModShift != 0 {
                    ui.buffer.ScrollEnd()
                } else if mods & tcell.ModShift != 0 {
                    ui.buffer.ScrollDown(page_h)
                } else {
                    ui.buffer.ScrollDown(delta)
                }

            case tcell.KeyLeft:
                if mods & tcell.ModShift != 0 {
                    ui.buffer.ScrollLeft(page_w)
                } else {
                    ui.buffer.ScrollLeft(delta)
                }

            case tcell.KeyRight:
                if mods & tcell.ModShift != 0 {
                    ui.buffer.ScrollRight(page_w)
                } else {
                    ui.buffer.ScrollRight(delta)
                }

            case tcell.KeyPgUp:
                ui.buffer.ScrollUp(page_h)

            case tcell.KeyPgDn:
                ui.buffer.ScrollDown(page_h)

            case tcell.KeyEnter:
                v := ui.prompt.Accept()

                if len(v) == 0 {
                    continue
                }

                hi.AddCommand(v)

                switch ui.ctx.Mode {
                case mode.Goto:
                    ui.buffer.Goto(v)

                    ui.State(ui.ctx.Last)

                default:
                    ui.buffer.Reset()
                
                    heap.AddFilter(v)
                }

            case tcell.KeyTab:
                ui.buffer.Reset()

                if mods & tcell.ModShift != 0 {
                    heap = hs.PrevHeap()
                } else {
                    heap = hs.NextHeap()
                }

            case tcell.KeyBackspace2:
                if len(ui.prompt.Value) > 0 {
                    ui.prompt.DelRune()
                } else if ui.ctx.Mode == mode.Goto {
                    ui.State(ui.ctx.Last)
                } else if len(*types.GetFilters()) > 0 {
                    ui.buffer.Reset()
                    heap.DelFilter()
                } else if ui.ctx.Mode == mode.Grep {
                    ui.State(mode.Less)
                }

            default:
                r := ev.Rune()

                switch r {
                case 0: // error
                    continue

                case 32: // space
                    if ui.ctx.Mode == mode.Less || ui.ctx.Mode == mode.Hex {
                        ui.buffer.ScrollDown(page_h)
                    } else {
                        ui.prompt.AddRune(r)
                    }

                default: // all other keys
                    if ui.ctx.Mode == mode.Less {
                        ui.State(mode.Grep)
                    }

                    ui.prompt.AddRune(r)
                }
            }
        }
    }
}

func (ui *UI) State(m mode.Mode) {
    if !ui.ctx.SwitchMode(m) {
        return
    }

    switch m {
    case mode.Less, mode.Hex: // static modes
        ui.prompt.Lock = true

    case mode.Grep, mode.Goto: // interactive modes
        ui.prompt.Lock = false
    }

    if ui.ctx.Last == mode.Hex || m == mode.Hex {
        ui.buffer.Reset()
    }
}

func (ui *UI) Close() {
    defer ui.ctx.Save()
    defer ui.term.Fini()
    defer ui.overlay.Close()
}

func (ui *UI) render(hs *heapset.HeapSet) (w int, h int) {
    defer ui.term.Show()

    _, heap := hs.Current()

    if heap.Type == types.Stdin {
        ui.term.Sync() // prevent hickups
    }

    ui.term.SetTitle(fmt.Sprintf("Forensic Examiner - %s", heap))
    ui.term.SetStyle(themes.Base)
    ui.term.Clear()

    x, y := 0, 0
    w, h = ui.term.Size()

    for _, base := range [...]lib.Queueable{
        ui.title,
        ui.buffer,
        ui.prompt,
    } {
        y += base.Render(hs, x, y, w, h-y)
    }

    ui.overlay.Render(0, 0, w, h)

    return
}
