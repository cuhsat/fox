package ui

import (
    "github.com/cuhsat/cu/pkg/fs"
)

type Status struct {
}

func NewStatus() *Status {
    return &Status{
    }
}

func (s *Status) Render(x, y int, heap *fs.Heap) {
    info := heap.Path

    for _, l := range heap.Chain {
        info += " > " + l.Name
    }

    if len(info) > width {
        info = info[:width-2] + "…"
    }

    printEx(x, y, info, StatusFg, StatusBg)
}
