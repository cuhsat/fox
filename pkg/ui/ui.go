package ui

import (
    "strings"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/mattn/go-runewidth"
    "github.com/nsf/termbox-go"
    "golang.design/x/clipboard"
)

const Delta = 1

// Buffer colors
const BufferFg = termbox.Attribute(248)
const BufferBg = termbox.Attribute(235)

// Prompt colors
const PromptFg = termbox.Attribute(248)
const PromptBg = termbox.Attribute(236)

var width, height, data, page int

type UI struct {
    buffer *Buffer
    prompt *Prompt
}

func NewUI() *UI {
    err := termbox.Init()
    
    if err != nil {
        fs.Panic(err)
    }
    
    err = clipboard.Init()

    if err != nil {
        fs.Panic(err)
    }

    termbox.SetInputMode(termbox.InputEsc)
    termbox.SetOutputMode(termbox.Output256)

    width, height = termbox.Size()

    return &UI{
        buffer: NewBuffer(),
        prompt: NewPrompt(),
    }
}

func (ui *UI) Loop(hs *fs.HeapSet) {
    ui.render(hs.Heap())

    for {
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            switch ev.Key {
            case termbox.KeyEsc, termbox.KeyCtrlQ:
                return

            case termbox.KeyCtrlC:
                clipboard.Write(clipboard.FmtText, hs.Heap().Copy())

            // case termbox.KeyCtrlS:
            //     // TODO: Save filtered results

            // case termbox.KeyCtrlSpace:
            //     // TODO: Toggle line numbers

            case termbox.KeyHome:
                ui.buffer.GoToBegin()

            case termbox.KeyEnd:
                ui.buffer.GoToEnd()

            case termbox.KeyPgup:
                ui.buffer.PageUp()

            case termbox.KeyPgdn:
                ui.buffer.PageDown()

            case termbox.KeyArrowUp:
                ui.buffer.ScrollUp(Delta)
                
            case termbox.KeyArrowDown:
                ui.buffer.ScrollDown(Delta)

            case termbox.KeyArrowLeft:
                ui.buffer.ScrollLeft(Delta)

            case termbox.KeyArrowRight:
                ui.buffer.ScrollRight(Delta)

            case termbox.KeyEnter:
                s := ui.prompt.Accept()

                if len(s) > 0 {
                    ui.buffer.Reset()
                    hs.Heap().AddFilter(s)
                }

            case termbox.KeyTab:
                c := hs.Heap().Chain

                ui.buffer.Reset()
                hs.Cycle()

                h := hs.Heap()
                h.NoFilter()

                for _, l := range c {
                    h.AddFilter(l.Name)
                }

            case termbox.KeyBackspace, termbox.KeyBackspace2:
                h := hs.Heap()

                if len(ui.prompt.Value) > 0 {
                    ui.prompt.DelChar()
                } else if len(h.Chain) > 0 {
                    ui.buffer.Reset()
                    h.DelFilter()
                }

            case termbox.KeySpace:
                ui.prompt.AddChar(' ')

            default:
                if ev.Ch != 0 {
                    ui.prompt.AddChar(ev.Ch)
                }
            }

        case termbox.EventError:
            fs.Error(ev.Err)
        }

        ui.render(hs.Heap())
    }
}

func (ui *UI) Close() {
    termbox.Close()
}

func (ui *UI) render(h *fs.Heap) {
    termbox.Clear(BufferFg, BufferBg)
    termbox.HideCursor()

    width, height = termbox.Size()

    data = len(h.SMap)
    page = height - 1

    line := strings.Repeat(" ", width)

    b := height - 1

    print(0, b, line, PromptFg, PromptBg)

    ui.buffer.Render(0, 0, h)
    ui.prompt.Render(0, b, h)

    termbox.Flush()
}

func length(s string) (l int) {
    for _, r := range []rune(s) {
        l += space(r)
    }

    return
}

func space(r rune) int {
    w := runewidth.RuneWidth(r)

    if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
        return 1
    } else {
        return w
    }
}

func print(x, y int, s string, fg, bg termbox.Attribute) {
    for _, r := range s {
        termbox.SetCell(x, y, r, fg, bg)

        x += space(r)
    }
}
