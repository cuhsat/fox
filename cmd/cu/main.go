package main

import (
    "flag"
    "io"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/history"
    "github.com/cuhsat/cu/pkg/ui"
    "github.com/cuhsat/cu/pkg/ui/mode"
)

var Version string = "dev"

func main() {
    var f heap.Filter

    m := mode.Normal
    x := flag.Bool("x", false, "Hex mode")
    h := flag.Bool("h", false, "Show help")
    v := flag.Bool("v", false, "Show version")

    flag.Var(&f, "f", "Filter")
    
    flag.CommandLine.SetOutput(io.Discard)
    flag.Parse()

    if *h || len(flag.Args()) < 1 {
        fs.Usage("cu [-xhv] [-f FILTER] PATH ...")
    }

    if *v {
        fs.Print("cu", Version)
    }

    if *x {
        m = mode.Hex
    }

    hs := heap.NewHeapSet(flag.Args(), f...)
    defer hs.ThrowAway()

    hi := history.NewHistory()
    defer hi.Close()

    ui := ui.NewUI(m)
    defer ui.Close()

    ui.Run(hs, hi)
}
