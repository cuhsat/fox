package widget

import (
    "unicode"

    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/gdamore/tcell/v2"
    "github.com/mattn/go-runewidth"
)

const (
    MinAscii = 0x20
    MaxAscii = 0x7f
)

const (
    Abbreviation = "…"
    NotPrintable = '.' // tcell.RuneBullet
)

type Stackable interface {
    Render(hs *heapset.HeapSet, x, y, w, h int) int
}

type widget struct {
    screen tcell.Screen
}

func (wi *widget) printBlank(x, y, w int, sty tcell.Style) {
    for i := 0; i < w; i++ {
        wi.screen.SetContent(x + i, y, ' ', nil, sty)
    }
}

func (wi *widget) printAscii(x, y int, s string, sty tcell.Style) {
    for _, r := range s {
        if r < MinAscii || r > MaxAscii {
            r = NotPrintable
        }

        wi.screen.SetContent(x, y, r, nil, sty)
        
        x += 1
    }
}

func (wi *widget) print(x, y int, s string, sty tcell.Style) {
    for _, r := range s {
        switch r {
        case '\t':
            r = tcell.RuneRArrow
        case '\r':
            r = tcell.RuneLArrow
        default:
            if !unicode.IsPrint(r) {
                r = NotPrintable
            }
        }

        wi.screen.SetContent(x, y, r, nil, sty)
        
        x += runewidth.RuneWidth(r)
    }
}

func abbrev(s string, x, w int) string {
    if x + length(s) > w + 1 {
        s = string([]rune(s)[:(w-x)-1]) + Abbreviation
    }

    return s
}

func length(s string) (l int) {
    for _, r := range s {
        l += runewidth.RuneWidth(r)
    }

    return
}
