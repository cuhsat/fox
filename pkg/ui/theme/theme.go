package theme

import (
    "github.com/gdamore/tcell/v2"
)

const (
    Default = "monokai"
)

var (
    Output tcell.Style
    Info   tcell.Style
    File   tcell.Style
    Number tcell.Style
    Filter tcell.Style
    Colors []tcell.Style
)

func Load(name string) {
    t := map[string][]int32{
        "monokai": Monokai,
    }[name]

    Output = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[0])).
        Background(tcell.NewHexColor(t[1]))

    Info = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[2])).
        Background(tcell.NewHexColor(t[3]))

    File = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[4])).
        Background(tcell.NewHexColor(t[5]))

    Number = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[6])).
        Background(tcell.NewHexColor(t[7]))

    Filter = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[8])).
        Background(tcell.NewHexColor(t[9]))

    for i := 10; i < 16; i++ {
        s := tcell.StyleDefault.
            Foreground(tcell.NewHexColor(t[i])).
            Background(tcell.NewHexColor(t[1]))

        Colors = append(Colors, s)
    }
}
