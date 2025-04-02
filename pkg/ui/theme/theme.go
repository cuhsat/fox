package theme

import (
    "github.com/gdamore/tcell/v2"
)

const (
    Default = "monokai"
)

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

func Load(theme string) {
    t := map[string][]int32{
        "monokai": Monokai,
    }[theme]

    Output = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[0])).
        Background(tcell.NewHexColor(t[1]))

    Header = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[2])).
        Background(tcell.NewHexColor(t[3]))

    Input = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[4])).
        Background(tcell.NewHexColor(t[5]))

    Error = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[6])).
        Background(tcell.NewHexColor(t[7]))

    Mode = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[8])).
        Background(tcell.NewHexColor(t[9]))

    Hint = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[10])).
        Background(tcell.NewHexColor(t[11]))

    Rule = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[12])).
        Background(tcell.NewHexColor(t[13]))

    Info = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[14])).
        Background(tcell.NewHexColor(t[15]))

    Line = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[16])).
        Background(tcell.NewHexColor(t[17]))

    for i := 18; i < 24; i++ {
        s := tcell.StyleDefault.
            Foreground(tcell.NewHexColor(t[i])).
            Background(tcell.NewHexColor(t[1]))

        Colors = append(Colors, s)
    }
}
