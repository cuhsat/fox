package buffer

import (
    "github.com/cuhsat/fx/pkg/fx/types/heap"
)

type Buffer struct {
    W int
    H int
}

type Line struct {
    Nr  string
    Str string
}

type Part struct {
    X int
    Y int
    Str string
}

type Context struct {
    Heap *heap.Heap

    Wrap bool
    Line bool

    X int
    Y int
    W int
    H int
}
