package theme

import (
    "strings"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/ui/theme/palette"
    "github.com/gdamore/tcell/v2"
)

const (
    Default = "monokai"
)

var palettes = map[string][]int32{
    "latte": palette.Latte,
    "frappe": palette.Frappe,
    "macchiato": palette.Macchiato,
    "mocha": palette.Mocha,
    "matrix": palette.Matrix,
    "monokai": palette.Monokai,
    "monochrome": palette.Monochrome,
}

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

func Load(name string) {
    p, ok := palettes[strings.ToLower(name)]

    if !ok {
        fs.Panic("theme not found")
    }

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

    for i := 18; i < 24; i++ {
        Colors = append(Colors, tcell.StyleDefault.
            Foreground(tcell.NewHexColor(p[i])).
            Background(tcell.NewHexColor(p[1])))
    }
}
