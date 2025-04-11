package themes

import (
    "strings"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/ui/themes/palette"
    "github.com/cuhsat/cu/pkg/ui/themes/palette/catppuccin"
    "github.com/cuhsat/cu/pkg/ui/themes/palette/monokai"
    "github.com/gdamore/tcell/v2"
)

const (
    Default = "monokai"
)

// global styles
var (
    Output tcell.Style
    Header tcell.Style
    Input  tcell.Style
    Error  tcell.Style
    Mode   tcell.Style
    Hint   tcell.Style
    Rule   tcell.Style
    Info   tcell.Style
    Line   tcell.Style
    Colors []tcell.Style
)

type Palettes [][]int32

type Themes struct {
    palettes Palettes
    names []string
    index int
}

func NewThemes(name string) *Themes {
    t := Themes{
        palettes: Palettes{
            // monokai
            monokai.Monokai,

            // catppuccin
            catppuccin.Latte,
            catppuccin.Frappe,
            catppuccin.Macchiato,
            catppuccin.Mocha,    

            // misc palettes
            palette.Matrix,
            palette.Monochrome,
        },
        names: []string{
            "Monokai",
            "Catppuccin-Latte",
            "Catppuccin-Frappe",
            "Catppuccin-Macchiato",
            "Catppuccin-Mocha",
            "Matrix",
            "Monochrome",
        },
        index: 0,
    }

    t.Load(name)

    return &t
}

func (t *Themes) Cycle() string {
    t.index += 1
    t.index %= len(t.names)

    n := t.names[t.index]

    t.Load(n)

    return n
}

func (t *Themes) Load(name string) {
    t.index = -1

    for i, n := range t.names {
        if strings.ToLower(n) == strings.ToLower(name) {
            t.index = i
        }
    }

    if t.index == -1 {
        fs.Panic("theme not found")
    }

    p := t.palettes[t.index]

    Output = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[0])).
        Background(tcell.NewHexColor(p[1]))

    Header = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[2])).
        Background(tcell.NewHexColor(p[3]))

    Input = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[4])).
        Background(tcell.NewHexColor(p[5]))

    Error = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[6])).
        Background(tcell.NewHexColor(p[7]))

    Mode = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[8])).
        Background(tcell.NewHexColor(p[9]))

    Hint = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[10])).
        Background(tcell.NewHexColor(p[11]))

    Rule = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[12])).
        Background(tcell.NewHexColor(p[13]))

    Info = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[14])).
        Background(tcell.NewHexColor(p[15]))

    Line = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(p[16])).
        Background(tcell.NewHexColor(p[17]))

    Colors = Colors[:0] // reset

    for i := 18; i < 24; i++ {
        Colors = append(Colors, tcell.StyleDefault.
            Foreground(tcell.NewHexColor(p[i])).
            Background(tcell.NewHexColor(p[1])))
    }
}
