// usage: cu [-x] [-h # | -t #] [-f FILTER] [PATH ...]
package main

import (
    "flag"
    "io"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/fs/history"
    "github.com/cuhsat/cu/pkg/ui"
    "github.com/cuhsat/cu/pkg/ui/mode"
)

func main() {
    var l heap.Limit
    var f heap.Filters

    // flags
    m := mode.Normal
    x := flag.Bool("x", false, "Hex mode")

    // limits
    flag.IntVar(&l.Head, "h", 0, "Head count")
    flag.IntVar(&l.Tail, "t", 0, "Tail count")

    // filters
    flag.Var(&f, "f", "Filter")
    
    flag.CommandLine.SetOutput(io.Discard)
    flag.Parse()

    a := flag.Args()

    if len(a) == 0 {
        a = append(a, ".")
    }

    if *x {
        m = mode.Hex
    }

    if l.Head > 0 && l.Tail > 0 {
        fs.Panic("either head or tail")
    }

    hs := heapset.NewHeapSet(a, l, f...)
    defer hs.ThrowAway()

    hi := history.NewHistory()
    defer hi.Close()

    ui := ui.NewUI(m)
    defer ui.Close()

    ui.Run(hs, hi)
}
