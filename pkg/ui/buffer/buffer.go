package buffer

import (
    "github.com/cuhsat/cu/pkg/fs/heap"
)

type Buffer struct {
    W, H int
}

type Line struct {
    Nr, Str string
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
