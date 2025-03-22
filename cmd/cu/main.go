package main

import (
    "os"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/ui"
)

func main() {
    if len(os.Args) < 2 {
        fs.Usage("cu <DIRECTORY|FILE|->")
    }

    hs := fs.NewHeapSet(os.Args[1])
    defer hs.ThrowAway()

    hi := fs.NewHistory()
    defer hi.Close()

    ui := ui.NewUI()
    defer ui.Close()

    ui.Run(hs, hi)
}
