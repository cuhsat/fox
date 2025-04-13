package block

import (
    "github.com/cuhsat/fx/internal/sys/heap"
)

type Block struct {
    W int
    H int
}

type Line struct {
    Nr  string
    Str string
}

type Context struct {
    Heap *heap.Heap

    Line bool
    Wrap bool

    X int
    Y int
    W int
    H int
}
