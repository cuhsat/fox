//go:build !minimal

package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"

	_ "github.com/gdamore/tcell/v2/encoding"

	"github.com/cuhsat/fox/internal"
	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/app/ai"
	"github.com/cuhsat/fox/internal/app/ai/agent"
	"github.com/cuhsat/fox/internal/app/ui/themes"
	"github.com/cuhsat/fox/internal/app/ui/widgets"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
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
	ctx *app.Context

	root tcell.Screen

	agent   *agent.Agent
	themes  *themes.Themes
	plugins *plugins.Plugins

	title   *widgets.Title
	view    *widgets.View
	prompt  *widgets.Prompt
	overlay *widgets.Overlay
}

func Start(args []string, invoke types.Invoke) {
	hs := heapset.New(args)
	defer hs.ThrowAway()

	hi := history.New()
	defer hi.Close()

	bg := bag.New()
	defer bg.Close()

	ui := create()
	defer ui.delete()

	ui.run(hs, hi, bg, invoke)
}

func create() *UI {
	runewidth.CreateLUT()

	root, err := tcell.NewScreen()

	if err != nil {
		sys.Panic(err)
	}

	err = root.Init()

	if err != nil {
		sys.Panic(err)
	}

	if !flags.Get().Opt.NoMouse {
		root.EnableMouse(tcell.MouseDragEvents)
	}

	root.EnablePaste()

	ctx := app.NewContext(root)

	ui := UI{
		ctx: ctx,

		root: root,

		themes:  themes.New(ctx.Theme()),
		plugins: plugins.New(),

		title:   widgets.NewTitle(ctx),
		view:    widgets.NewView(ctx),
		prompt:  widgets.NewPrompt(ctx),
		overlay: widgets.NewOverlay(ctx),
	}

	root.SetCursorStyle(widgets.Cursor, themes.Cursor)
	root.SetStyle(themes.Terminal)
	root.Sync()

	ui.render(nil)
	ui.change(flags.Get().UI.Mode)

	ui.agent = ai.NewAgent(ctx)

	return &ui
}

func (ui *UI) delete() {
	if ui.agent != nil {
		ui.agent.Close()
	}

	if ui.plugins != nil {
		plugins.Close()
	}

	ui.overlay.Close()
	ui.root.Fini()
	ui.ctx.Save()
}

func (ui *UI) run(hs *heapset.HeapSet, hi *history.History, bg *bag.Bag, util types.Invoke) {
	hs.SetCallbacks(func() {
		_ = ui.root.PostEvent(tcell.NewEventError(nil))
	}, func() {
		_ = ui.root.PostEvent(tcell.NewEventInterrupt(ui.ctx.IsFollow()))
	})

	events := make(chan tcell.Event, 128)
	closed := make(chan struct{})

	go ui.root.ChannelEvents(events, closed)
	go ui.overlay.Listen()

	ui.invoke(hs, util)

	flg := flags.Get()
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
					ui.view.ScrollLast()
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
				ui.view.Reset()

				hs.OpenLog()

				ui.root.Sync()
				ui.overlay.SendError("An error occurred")

			case *tcell.EventMouse:
				if flg.Opt.NoMouse {
					continue
				}

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

				if ui.ctx.IsNavi() {
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
					ui.view.Reset()
					hs.OpenHelp()

				case tcell.KeyF2:
					hs.Counts()
					ui.change(mode.Default)

				case tcell.KeyF3:
					hs.Entropy(0.0, 1.0)
					ui.change(mode.Default)

				case tcell.KeyF4:
					hs.Strings(3, math.MaxInt, true, nil)
					ui.change(mode.Default)

				case tcell.KeyF5:
					hs.HashSum(types.MD5)
					ui.change(mode.Default)

				case tcell.KeyF6:
					hs.HashSum(types.SHA1)
					ui.change(mode.Default)

				case tcell.KeyF7:
					hs.HashSum(types.SHA256)
					ui.change(mode.Default)

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

					if flg.Opt.NoPlugins {
						ui.overlay.SendError("Plugins deactivated")
						goto render
					}

					p, ok := ui.plugins.Hotkey[strings.ToLower(ev.Name())]

					if !ok {
						continue
					}

					go p.Execute(heap.Path, heap.Base, func(path, base, dir string) {
						name := fmt.Sprintf("%s %c %s", base, ui.ctx.Icon.HSep, p.Name)

						if len(dir) > 0 {
							hs.Open(dir)
						}

						hs.OpenPlugin(path, base, name)

						ui.ctx.ForceRender()
						ui.overlay.SendInfo(fmt.Sprintf("%s executed", p.Name))
					})

					if len(p.Mode) > 0 {
						ui.change(mode.Mode(p.Mode))
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
					} else if mods&tcell.ModCtrl != 0 {
						if hs.Len() > 1 {
							ui.view.SaveState(heap.Path)
							heap = hs.PrevHeap()
							ui.view.LoadState(heap.Path)
						}
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
					} else if mods&tcell.ModCtrl != 0 {
						if hs.Len() > 1 {
							ui.view.SaveState(heap.Path)
							heap = hs.NextHeap()
							ui.view.LoadState(heap.Path)
						}
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

					ui.root.Fill(' ', themes.Terminal)
					ui.root.Show()
					ui.root.SetCursorStyle(widgets.Cursor, themes.Cursor)

					ui.overlay.SendInfo(fmt.Sprintf("Theme %s", ui.ctx.Theme()))

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

				case tcell.KeyCtrlF:
					ui.change(mode.Fox)

				case tcell.KeyCtrlP:
					ui.ctx.TogglePinned()

				case tcell.KeyCtrlT:
					ui.ctx.ToggleFollow()

				case tcell.KeyCtrlN:
					ui.ctx.ToggleNavi()
					ui.view.Preserve()

				case tcell.KeyCtrlW:
					ui.ctx.ToggleWrap()
					ui.view.Preserve()

				case tcell.KeyCtrlJ:
					if heap.ModContext(-1) {
						ui.view.Reset()
					}

				case tcell.KeyCtrlK:
					if heap.ModContext(+1) {
						ui.view.Reset()
					}

				case tcell.KeyCtrlV:
					ui.root.GetClipboard()

				case tcell.KeyCtrlA:
					if ui.ctx.Mode().Static() {
						continue
					}

					if hs.Merge() {
						ui.overlay.SendInfo("Merged all open files")
					}

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

					if bg.Put(heap) {
						ui.overlay.SendInfo(fmt.Sprintf("%s saved to %s", heap.String(), bg.String()))
					}

				case tcell.KeyCtrlB:
					if fs.Exists(bg.Path) {
						ui.view.Reset()
						hs.OpenFile(bg.Path, bg.Path, bg.Path, types.Ignore)
					} else {
						ui.overlay.SendError(fmt.Sprintf("%s not found", bg.Path))
					}

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
						heap.AddFilter(v, 0, 0)
						ui.change(mode.Less)

					case mode.Goto:
						ui.view.Goto(v)
						ui.change(ui.ctx.Last())

					case mode.Open:
						ui.ctx.Background(func() {
							hs.Open(v)
						})
						ui.change(ui.ctx.Last())

					case mode.Fox:
						ui.view.Reset()
						ui.agent.Prompt(v)
						ui.ctx.Background(func() {
							ui.agent.Process(v, heap)
							heap.Title = ui.agent.String()
						})

						hs.OpenAgent(
							ui.agent.File.Name(),
							ui.agent.String(),
						)

					default:
						plugins.Input <- v
						ui.change(ui.ctx.Last())
					}

				case tcell.KeyDelete:
					ui.prompt.DelRune(false)

				case tcell.KeyBackspace, tcell.KeyBackspace2:
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

		render:
			ui.render(hs)
		}
	}
}

func (ui *UI) invoke(hs *heapset.HeapSet, util types.Invoke) {
	flg := flags.Get()

	switch util {
	case types.Compare:
		hs.Compare(
			flg.Compare.Git,
		).CloseOther()
	case types.Counts:
		hs.Counts().CloseOther()
	case types.Entropy:
		hs.Entropy(
			flg.Entropy.Min,
			flg.Entropy.Max,
		).CloseOther()
	case types.Strings:
		hs.Strings(
			flg.Strings.Min,
			flg.Strings.Max,
			flg.Strings.Ioc,
			flg.Strings.Re,
		).CloseOther()
	case types.Hash:
		hs.HashSum(
			flg.Hash.Algo.String(),
		).CloseOther()
	case types.None:
		// normal
	}
}

func (ui *UI) change(m mode.Mode) {
	// check for examiner support
	if m == mode.Fox && ui.agent == nil {
		ui.overlay.SendError("Agent is not available")
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

		title = fmt.Sprintf("%s %c %s", title, ui.ctx.Icon.HSep, heap)
	}

	ui.root.SetTitle(title)
	ui.root.SetStyle(themes.Terminal)
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
