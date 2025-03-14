package main

import (
    "os"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/ui"
)

func main() {
    if len(os.Args) < 2 {
        fs.Usage("cu FILE")
    }

    heap := fs.NewHeap(os.Args[1])

    defer heap.ThrowAway()

    ui := ui.NewUI()

    defer ui.Close()

    ui.Loop(heap)
}
