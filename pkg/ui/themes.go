package ui

import (
    "github.com/gdamore/tcell/v2"
)

const (
    ThemeDefault = "monokai"
)

var Themes = map[string][]int32{
    "monokai": {
        0x7f8490, // foreground
        0x222327, // background
        0x181819, // info foreground
        0x85d3f2, // info background
        0x181819, // file foreground
        0xa7df78, // file background
        0xe2e2e3, // filter foreground
        0x3e3b48, // filter background
        0xfc5d7c, // highlight 1
        0xf39660, // highlight 2
        0xe7c664, // highlight 3
        0x9ed072, // highlight 4
        0x76cce0, // highlight 5
        0xb39df3, // highlight 6
    },
}

var (
    StyleOutput     tcell.Style
    StyleInfo       tcell.Style
    StyleFile       tcell.Style
    StyleFilter     tcell.Style
    StyleHighlights []tcell.Style
)

func setTheme(name string) {
    var t = Themes[name]

    StyleOutput = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[0])).
        Background(tcell.NewHexColor(t[1]))

    StyleInfo = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[2])).
        Background(tcell.NewHexColor(t[3]))

    StyleFile = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[4])).
        Background(tcell.NewHexColor(t[5]))

    StyleFilter = tcell.StyleDefault.
        Foreground(tcell.NewHexColor(t[6])).
        Background(tcell.NewHexColor(t[7]))

    for i := 8; i < 14; i++ {
        style := tcell.StyleDefault.
            Foreground(tcell.NewHexColor(t[i])).
            Background(tcell.NewHexColor(t[1]))

        StyleHighlights = append(StyleHighlights, style)
    }
}
