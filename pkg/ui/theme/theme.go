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
    Info   tcell.Style
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

    Info = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[12])).
        Background(tcell.NewHexColor(t[13]))

    for i := 14; i < 20; i++ {
        s := tcell.StyleDefault.
            Foreground(tcell.NewHexColor(t[i])).
            Background(tcell.NewHexColor(t[1]))

        Colors = append(Colors, s)
    }
}
