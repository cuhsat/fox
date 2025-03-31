package main

import (
    "flag"
    "io"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui"
)

var Version string = "dev"

func main() {
    x := flag.Bool("x", false, "Hex mode")
    h := flag.Bool("h", false, "Show help")
    v := flag.Bool("v", false, "Show version")

    flag.CommandLine.SetOutput(io.Discard)
    flag.Parse()

    var mode int

    if *h || len(flag.Args()) < 1 {
        fs.Usage("cu [-xhv] PATH ...")
    }

    if *v {
        fs.Print("cu", Version)
    }

    if *x {
        mode = ui.ModeHex
    } else {
        mode = ui.ModeText
    }

    hs := data.NewHeapSet(flag.Args())
    defer hs.ThrowAway()

    hi := fs.NewHistory()
    defer hi.Close()

    ui := ui.NewUI(mode)
    defer ui.Close()

    ui.Run(hs, hi)
}
