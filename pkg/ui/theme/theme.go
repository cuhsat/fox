package theme

import (
    "github.com/gdamore/tcell/v2"
)

const (
    Default = "monokai"
)

var (
    Title  tcell.Style
    Output tcell.Style
    Info   tcell.Style
    Mode   tcell.Style
    Hint   tcell.Style
    Input  tcell.Style
    Colors []tcell.Style
)

func Load(name string) {
    t := map[string][]int32{
        "monokai": Monokai,
    }[name]

    Title = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[8])).
        Background(tcell.NewHexColor(t[7]))

    Output = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[0])).
        Background(tcell.NewHexColor(t[1]))

    Info = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[2])).
        Background(tcell.NewHexColor(t[3]))

    Mode = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[4])).
        Background(tcell.NewHexColor(t[5]))

    Hint = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[6])).
        Background(tcell.NewHexColor(t[7]))

    Input = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[8])).
        Background(tcell.NewHexColor(t[9]))

    for i := 10; i < 16; i++ {
        s := tcell.StyleDefault.
            Foreground(tcell.NewHexColor(t[i])).
            Background(tcell.NewHexColor(t[1]))

        Colors = append(Colors, s)
    }
}
