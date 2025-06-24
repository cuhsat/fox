package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"

	_ "github.com/gdamore/tcell/v2/encoding"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/fox/ui/context"
	"github.com/cuhsat/fox/internal/fox/ui/themes"
	"github.com/cuhsat/fox/internal/fox/ui/widgets"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
	"github.com/cuhsat/fox/internal/pkg/types/mode"
	"github.com/cuhsat/fox/internal/pkg/user/bag"
	"github.com/cuhsat/fox/internal/pkg/user/history"
	"github.com/cuhsat/fox/internal/pkg/user/plugins"
)

const (
	delta = 1 // lines
	wheel = 5 // lines
)

const (
	brPrefix = "ESC[200~" // bracketed paste start
	brSuffix = "ESC[201~" // bracketed paste end
)

type UI struct {
	ctx *context.Context

	root tcell.Screen

	plugins *plugins.Plugins
	themes  *themes.Themes

	title   *widgets.Title
	view    *widgets.View
	prompt  *widgets.Prompt
	overlay *widgets.Overlay
}

func New(m mode.Mode) *UI {
	runewidth.CreateLUT()

	root, err := tcell.NewScreen()

	if err != nil {
		sys.Panic(err)
	}

	err = root.Init()

	if err != nil {
		sys.Panic(err)
	}

	root.EnableMouse(tcell.MouseDragEvents)
	root.EnablePaste()

	ctx := context.New(root)

	ui := UI{
		ctx: ctx,

		root: root,

		plugins: plugins.New(),
		themes:  themes.New(ctx.Theme()),

		title:   widgets.NewTitle(ctx),
		view:    widgets.NewView(ctx),
		prompt:  widgets.NewPrompt(ctx),
		overlay: widgets.NewOverlay(ctx),
	}

	root.SetCursorStyle(widgets.Cursor, themes.Cursor)
	root.SetStyle(themes.Base)
	root.Sync()

	ui.render(nil)
	ui.change(m)

	return &ui
}

func (ui *UI) Run(hs *heapset.HeapSet, hi *history.History, bag *bag.Bag) {
	hs.Bind(func() {
		_ = ui.root.PostEvent(tcell.NewEventInterrupt(ui.ctx.IsFollow()))
	}, func() {
		_ = ui.root.PostEvent(tcell.NewEventError(nil))
	})

	events := make(chan tcell.Event, 128)
	closed := make(chan struct{})

	go ui.root.ChannelEvents(events, closed)
	go ui.overlay.Listen()

	esc := false

	for {
		select {
		case _ = <-closed:
			return // channels closed

		case ev := <-events:
			if ev == nil {
				return // term closed
			}

			w, h := ui.root.Size()

			_, heap := hs.Heap()

			switch ev := ev.(type) {
			case *tcell.EventInterrupt:
				v, ok := ev.Data().(bool)

				if ok && v {
					ui.view.ScrollEnd()
				}

			case *tcell.EventClipboard:
				if ui.ctx.Mode().Static() {
					continue
				}

				v := string(ev.Data())

				v = strings.TrimPrefix(v, brPrefix)
				v = strings.TrimSuffix(v, brSuffix)

				ui.prompt.Enter(v)

			case *tcell.EventResize:
				ui.root.Sync()
				ui.view.Reset()

			case *tcell.EventError:
				ui.change(mode.Less)

				hs.OpenLog()

				ui.overlay.SendError("An error occurred")

			case *tcell.EventMouse:
				btns := ev.Buttons()

				if btns&tcell.ButtonMiddle != 0 {
					ui.root.GetClipboard()
				} else if btns&tcell.WheelUp != 0 {
					ui.view.ScrollUp(wheel)
				} else if btns&tcell.WheelDown != 0 {
					ui.view.ScrollDown(wheel)
				} else if btns&tcell.WheelLeft != 0 {
					ui.view.ScrollLeft(wheel)
				} else if btns&tcell.WheelRight != 0 {
					ui.view.ScrollRight(wheel)
				}

			case *tcell.EventKey:
				mods := ev.Modifiers()

				pageW := w - 1 // minus text abbreviation
				pageH := h - 2 // minus title and status

				if ui.ctx.IsNumbers() {
					pageW -= text.Dec(heap.Count()) + 1
				}

				if ev.Key() != tcell.KeyEscape {
					esc = false
				}

				switch ev.Key() {
				case tcell.KeyEscape:
					if esc {
						return
					}

					if ui.ctx.Mode().Prompt() {
						if !ui.ctx.Last().Prompt() {
							ui.change(ui.ctx.Last())
						} else {
							ui.change(mode.Default)
						}
					}

					esc = true

				case tcell.KeyTab:
					if hs.Len() <= 1 {
						continue
					}

					ui.view.SaveState(heap.Path)

					if mods&tcell.ModShift != 0 {
						heap = hs.PrevHeap()
					} else {
						heap = hs.NextHeap()
					}

					ui.view.LoadState(heap.Path)

				case tcell.KeyF1:
					hs.Counts()
					ui.change(mode.Default)

				case tcell.KeyF2:
					hs.Strings()
					ui.change(mode.Default)

				case tcell.KeyF3:
					hs.Md5()
					ui.change(mode.Default)

				case tcell.KeyF4:
					hs.Sha1()
					ui.change(mode.Default)

				case tcell.KeyF5:
					hs.Sha256()
					ui.change(mode.Default)

				case tcell.KeyF6:
					hs.Sha3()
					ui.change(mode.Default)

				case tcell.KeyF7:
					fallthrough
				case tcell.KeyF8:
					fallthrough
				case tcell.KeyF9:
					fallthrough
				case tcell.KeyF10:
					fallthrough
				case tcell.KeyF11:
					fallthrough
				case tcell.KeyF12:
					if ui.plugins == nil {
						continue
					}

					p, ok := ui.plugins.Hotkey[ev.Name()]

					if !ok {
						continue
					}

					go p.Execute(heap.Path, heap.Base, func(file sys.File, base string) {
						name := fmt.Sprintf("%s (%s)", base, p.Name)
						hs.OpenFile(file.Name(), base, name, types.Plugin)
						ui.ctx.ForceRender()
					})

					ui.overlay.SendInfo(fmt.Sprintf("%s executed", p.Name))

					if len(p.Prompt) > 0 {
						ui.change(mode.Mode(p.Prompt))
					}

				case tcell.KeyUp:
					if ui.ctx.Mode().Prompt() {
						ui.prompt.Enter(hi.PrevLine())
					} else if mods&tcell.ModCtrl != 0 && mods&tcell.ModShift != 0 {
						ui.view.ScrollStart()
					} else if mods&tcell.ModShift != 0 {
						ui.view.ScrollUp(pageH)
					} else {
						ui.view.ScrollUp(delta)
					}

				case tcell.KeyDown:
					if ui.ctx.Mode().Prompt() {
						ui.prompt.Enter(hi.NextLine())
					} else if mods&tcell.ModCtrl != 0 && mods&tcell.ModShift != 0 {
						ui.view.ScrollEnd()
					} else if mods&tcell.ModShift != 0 {
						ui.view.ScrollDown(pageH)
					} else {
						ui.view.ScrollDown(delta)
					}

				case tcell.KeyLeft:
					if ui.ctx.Mode().Prompt() {
						if mods&tcell.ModCtrl != 0 {
							ui.prompt.MoveStart()
						} else {
							ui.prompt.Move(-1)
						}
					} else if mods&tcell.ModShift != 0 {
						ui.view.ScrollLeft(pageW)
					} else {
						ui.view.ScrollLeft(delta)
					}

				case tcell.KeyRight:
					if ui.ctx.Mode().Prompt() {
						if mods&tcell.ModCtrl != 0 {
							ui.prompt.MoveEnd()
						} else {
							ui.prompt.Move(+1)
						}
					} else if mods&tcell.ModShift != 0 {
						ui.view.ScrollRight(pageW)
					} else {
						ui.view.ScrollRight(delta)
					}

				case tcell.KeyHome:
					ui.view.ScrollStart()

				case tcell.KeyPgUp:
					ui.view.ScrollUp(pageH)

				case tcell.KeyPgDn:
					ui.view.ScrollDown(pageH)

				case tcell.KeyEnd:
					ui.view.ScrollEnd()

				case tcell.KeyCtrlCarat:
					ui.ctx.ChangeTheme(ui.themes.Cycle())

					ui.root.Fill(' ', themes.Base)
					ui.root.Show()
					ui.root.SetCursorStyle(widgets.Cursor, themes.Cursor)

					ui.overlay.SendInfo(fmt.Sprintf("Theme %s", ui.ctx.Theme()))

				case tcell.KeyCtrlSpace:
					ui.change(mode.Goto)

				case tcell.KeyCtrlO:
					ui.change(mode.Open)

				case tcell.KeyCtrlL:
					ui.change(mode.Less)

				case tcell.KeyCtrlF:
					ui.change(mode.Grep)

				case tcell.KeyCtrlX:
					ui.change(mode.Hex)

				case tcell.KeyCtrlT:
					ui.ctx.ToggleFollow()

				case tcell.KeyCtrlN:
					ui.ctx.ToggleNumbers()
					ui.view.Preserve()

				case tcell.KeyCtrlW:
					ui.ctx.ToggleWrap()
					ui.view.Preserve()

				case tcell.KeyCtrlV:
					ui.root.GetClipboard()

				case tcell.KeyCtrlC:
					if ui.ctx.Mode().Static() {
						continue
					}

					ui.root.SetClipboard(heap.Bytes())

					ui.overlay.SendInfo(fmt.Sprintf("%s copied to clipboard", heap.String()))

				case tcell.KeyCtrlS:
					if ui.ctx.Mode().Static() {
						continue
					}

					if bag.Put(heap) {
						ui.overlay.SendInfo(fmt.Sprintf("%s saved to %s", heap.String(), bag.Path))
					}

				case tcell.KeyCtrlB:
					if sys.Exists(bag.Path) {
						ui.view.Reset()
						hs.OpenFile(bag.Path, bag.Path, bag.Path, types.Ignore)
					} else {
						ui.overlay.SendError(fmt.Sprintf("%s not found", bag.Path))
					}

				case tcell.KeyCtrlH:
					ui.view.Reset()
					hs.OpenHelp()

				case tcell.KeyCtrlD:
					ui.view.Reset()
					hs.OpenLog()

				case tcell.KeyCtrlQ:
					ui.view.Reset()

					if hs.CloseHeap() == nil {
						return // exit
					}

				case tcell.KeyCtrlZ:
					err := ui.root.Suspend()

					if err != nil {
						sys.Error(err)
						continue
					}

					sys.Shell()

					err = ui.root.Resume()

					if err != nil {
						sys.Panic(err)
					}

				case tcell.KeyEnter:
					v := ui.prompt.ReadLine()
					m := ui.ctx.Mode()

					if m.Prompt() && len(v) == 0 {
						continue
					}

					hi.AddLine(v)

					switch m {
					case mode.Less, mode.Hex:
						ui.view.ScrollLine()

					case mode.Grep:
						ui.view.Reset()
						heap.AddFilter(v)
						ui.change(mode.Less)

					case mode.Goto:
						ui.view.Goto(v)
						ui.change(ui.ctx.Last())

					case mode.Open:
						ui.ctx.Background(func() {
							hs.Open(v)
						})
						ui.change(ui.ctx.Last())

					default:
						plugins.Input <- v
						ui.change(ui.ctx.Last())
					}

				case tcell.KeyDelete, tcell.KeyCtrlK:
					ui.prompt.DelRune(false)

				case tcell.KeyBackspace2:
					if len(ui.prompt.Value()) > 0 {
						ui.prompt.DelRune(true)
					} else if ui.ctx.Mode().Prompt() {
						if !ui.ctx.Last().Prompt() {
							ui.change(ui.ctx.Last())
						} else {
							ui.change(mode.Default)
						}
					} else if len(heap.Patterns()) > 0 {
						if !ui.ctx.Mode().Static() {
							ui.view.Reset()
							heap.DelFilter()
						}
					}

				default:
					r := ev.Rune()

					switch r {
					case 0: // error
						continue

					case 32: // space
						if ui.prompt.Locked() {
							ui.view.ScrollDown(pageH)
						} else {
							ui.prompt.AddRune(r)
						}

					default: // all other keys
						if ui.ctx.Mode() == mode.Less {
							ui.change(mode.Grep)
						}

						ui.prompt.AddRune(r)
					}
				}
			}

			ui.render(hs)
		}
	}
}

func (ui *UI) Close() {
	if ui.plugins != nil {
		plugins.Close()
	}

	ui.overlay.Close()
	ui.root.Fini()
	ui.ctx.Save()
}

func (ui *UI) change(m mode.Mode) {
	if !ui.ctx.SwitchMode(m) {
		return
	}

	// former mode
	if ui.ctx.Last().Prompt() {
		ui.prompt.Enter("")
	}

	// actual mode
	ui.prompt.Lock(!m.Prompt())

	// force the cursor off
	if ui.prompt.Locked() {
		ui.ctx.Root.HideCursor()
	}

	if ui.ctx.Last().Static() || m.Static() {
		ui.view.Reset()
	}
}

func (ui *UI) render(hs *heapset.HeapSet) {
	title := fox.Product

	if hs != nil {
		_, heap := hs.Heap()

		if heap.Type == types.Stdin {
			ui.root.Sync() // prevent hiccups
		}

		title = fmt.Sprintf("%s － %s", title, heap)
	}

	ui.root.SetTitle(title)
	ui.root.SetStyle(themes.Base)
	ui.root.Clear()

	x, y := 0, 0
	w, h := ui.root.Size()

	for _, base := range [...]widgets.Queueable{
		ui.title,
		ui.view,
		ui.prompt,
	} {
		y += base.Render(hs, x, y, w, h-y)
	}

	ui.overlay.Render(0, 0, w, h)

	ui.root.Show()
}
