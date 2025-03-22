package ui

import (
    "github.com/gdamore/tcell/v2"
)

var Themes = map[string][]int32{
    "monokai": {
        0x7f8490, // foreground
        0x222327, // background
        0xe7c664, // highlight
        0x181819, // info foreground
        0x85d3f2, // info background
        0x181819, // file foreground
        0xa7df78, // file background
        0xe2e2e3, // filter foreground
        0x3e3b48, // filter background
    },
}

var (
    StyleOutput    tcell.Style
    StyleHighlight tcell.Style
    StyleInfo      tcell.Style
    StyleFile      tcell.Style
    StyleFilter    tcell.Style
)

func setTheme(name string) {
    var t = Themes[name]

    StyleOutput = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[0])).
        Background(tcell.NewHexColor(t[1]))

    StyleHighlight = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[2])).
        Background(tcell.NewHexColor(t[1]))

    StyleInfo = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[3])).
        Background(tcell.NewHexColor(t[4]))

    StyleFile = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[5])).
        Background(tcell.NewHexColor(t[6]))

    StyleFilter = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[7])).
        Background(tcell.NewHexColor(t[8]))
}
