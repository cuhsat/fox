package ui

import (
    "github.com/nsf/termbox-go"
)

const Prompt = "> "

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
    termbox.SetCursor(len(Prompt) + x + s.cx, y)

    printEx(x, y, Prompt + s.value, SearchFg, SearchBg)
}

func (s *Search) AddChar(r rune) {
    if s.cx >= width - (len(Prompt)+1) {
        return
    }

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
