package ui

import (
    "github.com/nsf/termbox-go"
)

var Themes = map[string][]int{
    "markdown": {248, 235, 231, 235, 186, 235, 231, 235},
}

var (
    BufferFg termbox.Attribute
    BufferBg termbox.Attribute
    SearchFg termbox.Attribute
    SearchBg termbox.Attribute
    PromptFg termbox.Attribute
    PromptBg termbox.Attribute
    CursorFg termbox.Attribute
    CursorBg termbox.Attribute
)

func setTheme(name string) {
    var t = Themes[name]

    BufferFg = termbox.Attribute(t[0])
    BufferBg = termbox.Attribute(t[1])

    SearchFg = termbox.Attribute(t[2]) | termbox.AttrBold
    SearchBg = termbox.Attribute(t[3])

    PromptFg = termbox.Attribute(t[4])
    PromptBg = termbox.Attribute(t[5])

    CursorFg = termbox.Attribute(t[6]) // | termbox.AttrBold
    CursorBg = termbox.Attribute(t[7])
}
