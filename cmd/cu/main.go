package main

import (
    "os"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/ui"
)

func main() {
    if len(os.Args) < 2 {
        fs.Usage("cu - or PATH")
    }

    hs := fs.NewHeapSet(os.Args[1])

    defer hs.ThrowAway()

    ui := ui.NewUI("monokai")

    defer ui.Close()

    ui.Loop(hs)
}
