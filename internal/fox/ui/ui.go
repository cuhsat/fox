//go:build !no_ui

package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/fox/ai"
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
	Build = true
)

const (
	delta = 1 // lines
)

const (
	brPrefix = "ESC[200~" // bracketed paste start
	brSuffix = "ESC[201~" // bracketed paste end
)

const (
	cursor = tcell.CursorStyleBlinkingBar
)

type UI struct {
	ctx *context.Context

	root tcell.Screen

	agent   *ai.Agent
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

	root.EnableMouse()
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

	if ai.Build && ai.Init(ctx.Model()) {
		ui.agent = ai.NewAgent()
	}

	root.SetCursorStyle(cursor, themes.Cursor)
	root.SetStyle(themes.Base)
	root.Sync()

	ui.change(m)

	return &ui
}

func (ui *UI) Run(hs *heapset.HeapSet, hi *history.History, bag *bag.Bag) {
	hs.Bind(func() {
		_ = ui.root.PostEvent(tcell.NewEventInterrupt(ui.ctx.IsTail()))
	}, func() {
		_ = ui.root.PostEvent(tcell.NewEventError(nil))
	})

	events := make(chan tcell.Event, 128)
	closed := make(chan struct{})

	go ui.root.ChannelEvents(events, closed)
	go ui.overlay.Listen()

	if ui.agent != nil {
		go ui.agent.Listen(hi)
	}

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
				if ui.ctx.Mode() == mode.Hex {
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
				hs.OpenLog()

				ui.overlay.SendError("An error occurred")

			case *tcell.EventMouse:
				btns := ev.Buttons()

				if btns&tcell.ButtonMiddle != 0 {
					ui.root.GetClipboard()
				} else if btns&tcell.WheelUp != 0 {
					ui.view.ScrollUp(delta)
				} else if btns&tcell.WheelDown != 0 {
					ui.view.ScrollDown(delta)
				} else if btns&tcell.WheelLeft != 0 {
					ui.view.ScrollLeft(delta)
				} else if btns&tcell.WheelRight != 0 {
					ui.root.GetClipboard()
				}

			case *tcell.EventKey:
				mods := ev.Modifiers()

				pageW := w - 1 // minus text abbreviation
				pageH := h - 2 // minus title and status

				if ui.ctx.IsLine() {
					pageW -= text.Dec(heap.Total()) + 1
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
							ui.change(mode.Less)
						}
					}

					esc = true

				case tcell.KeyTab:
					ui.view.Reset()

					if mods&tcell.ModShift != 0 {
						heap = hs.PrevHeap()
					} else {
						heap = hs.NextHeap()
					}

				case tcell.KeyF1:
					fallthrough
				case tcell.KeyF2:
					fallthrough
				case tcell.KeyF3:
					fallthrough
				case tcell.KeyF4:
					fallthrough
				case tcell.KeyF5:
					fallthrough
				case tcell.KeyF6:
					fallthrough
				case tcell.KeyF7:
					if ui.plugins == nil {
						continue
					}

					pl, ok := ui.plugins.Plugins[ev.Name()]

					if !ok {
						continue
					}

					go pl.Execute(heap.Path, heap.Base, hs.Files(), func(p, b, t string) {
						hs.OpenFile(p, b, t, types.Stdout)
						ui.ctx.Interrupt()
					})

					ui.overlay.SendInfo(fmt.Sprintf("%s executed", pl.Name))

					if len(pl.Prompt) > 0 {
						ui.change(mode.Mode(pl.Prompt))
					}

				case tcell.KeyF8:
					hs.Counts()

				case tcell.KeyF9:
					hs.Md5()

				case tcell.KeyF10:
					hs.Sha1()

				case tcell.KeyF11:
					hs.Sha256()

				case tcell.KeyF12:
					hs.Sha3()

				case tcell.KeyUp:
					if ui.ctx.Mode().Prompt() {
						ui.prompt.Enter(hi.PrevCommand())
					} else if mods&tcell.ModCtrl != 0 && mods&tcell.ModShift != 0 {
						ui.view.ScrollStart()
					} else if mods&tcell.ModShift != 0 {
						ui.view.ScrollUp(pageH)
					} else {
						ui.view.ScrollUp(delta)
					}

				case tcell.KeyDown:
					if ui.ctx.Mode().Prompt() {
						ui.prompt.Enter(hi.NextCommand())
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

				case tcell.KeyCtrlSpace:
					ui.change(mode.Goto)

				case tcell.KeyCtrlO:
					ui.change(mode.Open)

				case tcell.KeyCtrlL:
					ui.change(mode.Less)

				case tcell.KeyCtrlG:
					ui.change(mode.Grep)

				case tcell.KeyCtrlX:
					ui.change(mode.Hex)

				case tcell.KeyCtrlA:
					ui.change(mode.Rag)

				case tcell.KeyCtrlT:
					ui.ctx.ChangeTheme(ui.themes.Cycle())

					ui.root.Fill(' ', themes.Base)
					ui.root.Show()

					ui.root.SetCursorStyle(cursor, themes.Cursor)

					ui.overlay.SendInfo(fmt.Sprintf("Theme %s", ui.ctx.Theme()))

				case tcell.KeyCtrlV:
					if ui.ctx.Mode() == mode.Hex {
						continue
					}

					ui.root.GetClipboard()

				case tcell.KeyCtrlC:
					if ui.ctx.Mode() == mode.Hex {
						continue
					}

					ui.root.SetClipboard(heap.Bytes())

					ui.overlay.SendInfo(fmt.Sprintf("%s copied to clipboard", heap.String()))

				case tcell.KeyCtrlS:
					if ui.ctx.Mode() == mode.Hex {
						continue
					}

					if !bag.Put(heap) {
						continue
					}

					ui.overlay.SendInfo(fmt.Sprintf("%s saved to %s", heap.String(), bag.Path))

				case tcell.KeyCtrlE:
					if sys.Exists(bag.Path) {
						hs.OpenFile(bag.Path, bag.Path, bag.Path, types.Regular)
					} else {
						ui.overlay.SendError(fmt.Sprintf("%s not found", bag.Path))
					}

				case tcell.KeyCtrlH:
					hs.OpenHelp()

				case tcell.KeyCtrlD:
					hs.OpenLog()

				case tcell.KeyCtrlQ:
					ui.view.Reset()

					if hs.CloseHeap() == nil {
						return // exit
					}

				case tcell.KeyCtrlF:
					if ui.ctx.Mode() != mode.Hex {
						ui.ctx.ToggleFollow()
					}

				case tcell.KeyCtrlN:
					if ui.ctx.Mode() != mode.Hex {
						ui.ctx.ToggleNumbers()
					}

				case tcell.KeyCtrlW:
					if ui.ctx.Mode() != mode.Hex {
						ui.ctx.ToggleWrap()
						ui.view.Reset()
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

					if len(v) == 0 {
						continue
					}

					hi.AddCommand(v)

					m := ui.ctx.Mode()

					switch m {
					case mode.Grep:
						_ = types.GetFilters().Set(v)
						ui.view.Reset()
						ui.ctx.Background(func() {
							heap.AddFilter(v)
						})
						ui.change(mode.Less)

					case mode.Goto:
						ui.view.Goto(v)
						ui.change(ui.ctx.Last())

					case mode.Open:
						ui.ctx.Background(func() {
							hs.Open(v)
						})
						ui.change(ui.ctx.Last())

					case mode.Rag:
						if ui.agent != nil {
							ui.view.Reset()
							ui.ctx.Background(func() {
								ui.prompt.Lock(true)
								ui.agent.Prompt(v, heap)
								ui.prompt.Lock(false)
							})
							hs.OpenChat(ui.agent.Path())
						}

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
							ui.change(mode.Less)
						}
					} else if len(*types.GetFilters()) > 0 {
						if ui.ctx.Mode() != mode.Hex {
							types.GetFilters().Pop()
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

	if ui.agent != nil {
		ui.agent.Close()
	}

	ui.overlay.Close()
	ui.root.Fini()
	ui.ctx.Save()
}

func (ui *UI) change(m mode.Mode) {
	// check for RAG support
	if m == mode.Rag && ui.agent == nil {
		ui.overlay.SendError("RAG agent not available")
		return
	}

	if !ui.ctx.SwitchMode(m) {
		return
	}

	// former mode
	if ui.ctx.Last().Prompt() {
		ui.prompt.Enter("")
	}

	// actual mode
	ui.prompt.Lock(!m.Prompt())

	// force cursor off
	if ui.prompt.Locked() {
		ui.ctx.Root.HideCursor()
	}

	if ui.ctx.Last() == mode.Hex || m == mode.Hex {
		ui.view.Reset()
	}
}

func (ui *UI) render(hs *heapset.HeapSet) {
	defer ui.root.Show()

	_, heap := hs.Heap()

	if heap.Type == types.Stdin {
		ui.root.Sync() // prevent hiccups
	}

	ui.root.SetTitle(fmt.Sprintf("%s - %s", fox.Product, heap))
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
}
