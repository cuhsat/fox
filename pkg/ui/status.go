package ui

import (
    "github.com/cuhsat/cu/pkg/fs"
)

type Status struct {
    Search string
}

func NewStatus() *Status {
    return &Status{
        Search: "",
    }
}

func (s *Status) Render(x, y int, heap *fs.Heap) {
    info := heap.Path

    if len(s.Search) > 0 {
        info += s.Search
    }

    printEx(x, y, info, CStatus, CClear)
}
