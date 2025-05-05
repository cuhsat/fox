package gui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"

	"github.com/cuhsat/fx/internal/app/fx"
	"github.com/cuhsat/fx/internal/pkg/gui/context"
	"github.com/cuhsat/fx/internal/pkg/gui/themes"
	"github.com/cuhsat/fx/internal/pkg/gui/widgets"
	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/text"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/heapset"
	"github.com/cuhsat/fx/internal/pkg/types/mode"
	"github.com/cuhsat/fx/internal/pkg/user/bag"
	"github.com/cuhsat/fx/internal/pkg/user/history"
	"github.com/cuhsat/fx/internal/pkg/user/plugins"
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

type GUI struct {
	ctx *context.Context

	root tcell.Screen

	themes *themes.Themes

	title   *widgets.Title
	view    *widgets.View
	status  *widgets.Status
	overlay *widgets.Overlay

	plugins *plugins.Plugins
}

func New(m mode.Mode) *GUI {
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

	gui := GUI{
		ctx: ctx,

		root: root,

		themes: themes.New(ctx.Theme()),

		title:   widgets.NewTitle(ctx),
		view:    widgets.NewView(ctx),
		status:  widgets.NewStatus(ctx),
		overlay: widgets.NewOverlay(ctx),

		plugins: plugins.New(),
	}

	root.SetCursorStyle(cursor, themes.Cursor)
	root.SetStyle(themes.Base)
	root.Sync()

	gui.change(m)

	return &gui
}

func (gui *GUI) Run(hs *heapset.HeapSet, hi *history.History, bag *bag.Bag) {
	hs.Bind(func() {
		gui.root.PostEvent(tcell.NewEventInterrupt(gui.ctx.IsFollow()))
	}, func() {
		gui.root.PostEvent(tcell.NewEventError(nil))
	})

	events := make(chan tcell.Event, 128)
	closed := make(chan struct{})

	go gui.root.ChannelEvents(events, closed)
	go gui.overlay.Listen()

	esc := false

	for {
		select {
		case _ = <-closed:
			return // channels closed

		case ev := <-events:
			if ev == nil {
				return // term closed
			}

			w, h := gui.root.Size()

			_, heap := hs.Heap()

			switch ev := ev.(type) {
			case *tcell.EventInterrupt:
				v, ok := ev.Data().(bool)

				if ok && v {
					gui.view.ScrollEnd()
				}

			case *tcell.EventClipboard:
				if gui.ctx.Mode() == mode.Hex {
					continue
				}

				v := string(ev.Data())

				v = strings.TrimPrefix(v, brPrefix)
				v = strings.TrimSuffix(v, brSuffix)

				gui.status.Value = v

			case *tcell.EventResize:
				gui.root.Sync()
				gui.view.Reset()

			case *tcell.EventError:
				hs.OpenLog()

				gui.overlay.SendError("An error occurred")

			case *tcell.EventMouse:
				btns := ev.Buttons()

				if btns&tcell.ButtonMiddle != 0 {
					gui.root.GetClipboard()
				} else if btns&tcell.WheelUp != 0 {
					gui.view.ScrollUp(delta)
				} else if btns&tcell.WheelDown != 0 {
					gui.view.ScrollDown(delta)
				} else if btns&tcell.WheelLeft != 0 {
					gui.view.ScrollLeft(delta)
				} else if btns&tcell.WheelRight != 0 {
					gui.root.GetClipboard()
				}

			case *tcell.EventKey:
				mods := ev.Modifiers()

				page_w := w - 1 // minus text abbreviation
				page_h := h - 2 // minus title and status

				if gui.ctx.IsLine() {
					page_w -= text.Dec(heap.Total()) + 1
				}

				if ev.Key() != tcell.KeyEscape {
					esc = false
				}

				switch ev.Key() {
				case tcell.KeyEscape:
					if esc {
						return
					}

					if gui.ctx.Mode().Prompt() {
						if !gui.ctx.Last().Prompt() {
							gui.change(gui.ctx.Last())
						} else {
							gui.change(mode.Less)
						}
					}

					esc = true

				case tcell.KeyTab:
					gui.view.Reset()

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
					fallthrough
				case tcell.KeyF8:
					if gui.plugins == nil {
						continue
					}

					pl, ok := gui.plugins.Plugins[ev.Name()]

					if !ok {
						continue
					}

					go pl.Execute(hs, gui.ctx.Interrupt)

					if len(pl.Input) > 0 {
						gui.change(mode.Mode(pl.Input))
					}

				case tcell.KeyF9:
					hs.Counts()

				case tcell.KeyF10:
					hs.Md5()

				case tcell.KeyF11:
					hs.Sha1()

				case tcell.KeyF12:
					hs.Sha256()

				case tcell.KeyUp:
					if mods&tcell.ModAlt != 0 {
						gui.status.Value = hi.PrevCommand()
					} else if mods&tcell.ModCtrl != 0 && mods&tcell.ModShift != 0 {
						gui.view.ScrollStart()
					} else if mods&tcell.ModShift != 0 {
						gui.view.ScrollUp(page_h)
					} else {
						gui.view.ScrollUp(delta)
					}

				case tcell.KeyDown:
					if mods&tcell.ModAlt != 0 {
						gui.status.Value = hi.NextCommand()
					} else if mods&tcell.ModCtrl != 0 && mods&tcell.ModShift != 0 {
						gui.view.ScrollEnd()
					} else if mods&tcell.ModShift != 0 {
						gui.view.ScrollDown(page_h)
					} else {
						gui.view.ScrollDown(delta)
					}

				case tcell.KeyLeft:
					if mods&tcell.ModShift != 0 {
						gui.view.ScrollLeft(page_w)
					} else {
						gui.view.ScrollLeft(delta)
					}

				case tcell.KeyRight:
					if mods&tcell.ModShift != 0 {
						gui.view.ScrollRight(page_w)
					} else {
						gui.view.ScrollRight(delta)
					}

				case tcell.KeyHome:
					gui.view.ScrollStart()

				case tcell.KeyPgUp:
					gui.view.ScrollUp(page_h)

				case tcell.KeyPgDn:
					gui.view.ScrollDown(page_h)

				case tcell.KeyEnd:
					gui.view.ScrollEnd()

				case tcell.KeyCtrlSpace:
					gui.change(mode.Goto)

				case tcell.KeyCtrlL:
					gui.change(mode.Less)

				case tcell.KeyCtrlG:
					gui.change(mode.Grep)

				case tcell.KeyCtrlX:
					gui.change(mode.Hex)

				case tcell.KeyCtrlO:
					gui.change(mode.Open)

				case tcell.KeyCtrlT:
					gui.ctx.ChangeTheme(gui.themes.Cycle())

					gui.root.Fill(' ', themes.Base)
					gui.root.Show()

					gui.root.SetCursorStyle(cursor, themes.Cursor)

					gui.overlay.SendInfo(fmt.Sprintf("Theme %s", gui.ctx.Theme()))

				case tcell.KeyCtrlV:
					if gui.ctx.Mode() == mode.Hex {
						continue
					}

					gui.root.GetClipboard()

				case tcell.KeyCtrlC:
					if gui.ctx.Mode() == mode.Hex {
						continue
					}

					gui.root.SetClipboard(heap.Bytes())

					gui.overlay.SendInfo("Copied to clipboard")

				case tcell.KeyCtrlS:
					if gui.ctx.Mode() == mode.Hex {
						continue
					}

					if !bag.Put(heap) {
						continue
					}

					gui.overlay.SendInfo(fmt.Sprintf("Saved to %s", bag.Path))

				case tcell.KeyCtrlE:
					if sys.Exists(bag.Path) {
						hs.OpenFile(bag.Path, bag.Path, types.Regular)
					} else {
						gui.overlay.SendError(fmt.Sprintf("%s not found", bag.Path))
					}

				case tcell.KeyCtrlD:
					hs.OpenLog()

				case tcell.KeyCtrlQ:
					gui.view.Reset()

					if hs.CloseHeap() == nil {
						return // exit
					}

				case tcell.KeyCtrlF:
					if gui.ctx.Mode() != mode.Hex {
						gui.ctx.ToggleFollow()
					}

				case tcell.KeyCtrlN:
					if gui.ctx.Mode() != mode.Hex {
						gui.ctx.ToggleNumbers()
					}

				case tcell.KeyCtrlW:
					if gui.ctx.Mode() != mode.Hex {
						gui.ctx.ToggleWrap()
						gui.view.Reset()
					}

				case tcell.KeyCtrlZ:
					err := gui.root.Suspend()

					if err != nil {
						sys.Error(err)
						continue
					}

					sys.Shell()

					err = gui.root.Resume()

					if err != nil {
						sys.Panic(err)
					}

				case tcell.KeyEnter:
					v := gui.status.Accept()

					if len(v) == 0 {
						continue
					}

					hi.AddCommand(v)

					m := gui.ctx.Mode()

					switch m {
					case mode.Grep:
						types.Filters().Set(v)
						gui.view.Reset()
						gui.ctx.Background(func() {
							heap.AddFilter(v)
						})
						gui.change(mode.Less)

					case mode.Goto:
						gui.view.Goto(v)
						gui.change(gui.ctx.Last())

					case mode.Open:
						gui.ctx.Background(func() {
							hs.Open(v)
						})
						gui.change(gui.ctx.Last())

					default:
						plugins.Input <- v
						gui.change(gui.ctx.Last())
					}

				case tcell.KeyBackspace2:
					if len(gui.status.Value) > 0 {
						gui.status.DelRune()
					} else if gui.ctx.Mode().Prompt() {
						if !gui.ctx.Last().Prompt() {
							gui.change(gui.ctx.Last())
						} else {
							gui.change(mode.Less)
						}
					} else if len(*types.Filters()) > 0 {
						types.Filters().Pop()
						gui.view.Reset()
						heap.DelFilter()
					}

				default:
					r := ev.Rune()

					switch r {
					case 0: // error
						continue

					case 32: // space
						if gui.status.Lock {
							gui.view.ScrollDown(page_h)
						} else {
							gui.status.AddRune(r)
						}

					default: // all other keys
						if gui.ctx.Mode() == mode.Less {
							gui.change(mode.Grep)
						}

						gui.status.AddRune(r)
					}
				}
			}

			gui.render(hs)
		}
	}
}

func (gui *GUI) Close() {
	if gui.plugins != nil {
		plugins.Close()
	}

	gui.overlay.Close()
	gui.root.Fini()
	gui.ctx.Save()
}

func (gui *GUI) change(m mode.Mode) {
	if !gui.ctx.SwitchMode(m) {
		return
	}

	// former mode
	if gui.ctx.Last().Prompt() {
		gui.status.Value = ""
	}

	// actual mode
	gui.status.Lock = !m.Prompt()

	if gui.ctx.Last() == mode.Hex || m == mode.Hex {
		gui.view.Reset()
	}
}

func (gui *GUI) render(hs *heapset.HeapSet) {
	defer gui.root.Show()

	_, heap := hs.Heap()

	if heap.Type == types.Stdin {
		gui.root.Sync() // prevent hickups
	}

	gui.root.SetTitle(fmt.Sprintf("%s - %s", fx.Product, heap))
	gui.root.SetStyle(themes.Base)
	gui.root.Clear()

	x, y := 0, 0
	w, h := gui.root.Size()

	for _, base := range [...]widgets.Queueable{
		gui.title,
		gui.view,
		gui.status,
	} {
		y += base.Render(hs, x, y, w, h-y)
	}

	gui.overlay.Render(0, 0, w, h)
}
