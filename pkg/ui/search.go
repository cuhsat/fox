package ui

import (
    "github.com/nsf/termbox-go"
)

type Search struct {
    value string // value
    cx    int    // cursor
}

func NewSearch() *Search {
    return &Search{
        value: "",
        cx: 0,
    }
}

func (s *Search) Render(x, y int) {
    termbox.SetCursor(s.cx + x, y)

    printEx(x, y, s.value, termbox.ColorLightGreen | termbox.AttrBold, termbox.ColorDefault)
}

func (s *Search) AddChar(r rune) {
    s.cx++
    s.value += string(r)
}

func (s *Search) DelChar() {
    if s.cx == 0 {
        return
    }

    s.cx--
    s.value = s.value[:len(s.value)-1]
}

func (s *Search) GetValue() (v string) {
    v, s.cx, s.value = s.value, 0, ""

    return
}
