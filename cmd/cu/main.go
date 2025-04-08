package main

import (
    "flag"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/config"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/fs/history"
    "github.com/cuhsat/cu/pkg/ui"
    "github.com/cuhsat/cu/pkg/ui/mode"
)

func usage() {
    fs.Usage("usage: cu [-h # | -t #] [-x | -e PATTERN] [PATH ...]")
}

func main() {
    var l heap.Limit
    var e heap.Filters

    // config
    c := config.Load()

    // flags
    m := mode.Less
    x := flag.Bool("x", false, "Hex mode")

    // limits
    flag.IntVar(&l.Head, "h", 0, "Head lines")
    flag.IntVar(&l.Tail, "t", 0, "Tail lines")

    // filters
    flag.Var(&e, "e", "Pattern")

    flag.Usage = usage
    flag.Parse()

    a := flag.Args()

    if len(a) == 0 {
        a = append(a, ".")
    }

    if *x {
        m = mode.Hex
    }

    if len(e) > 0 {
        m = mode.Grep
    } 

    if l.Head > 0 && l.Tail > 0 {
        fs.Usage("either head or tail")
    }

    if *x && len(e) > 0 {
        fs.Usage("either hex or pattern")
    }

    hs := heapset.NewHeapSet(a, l, e...)
    defer hs.ThrowAway()

    hi := history.NewHistory()
    defer hi.Close()

    ui := ui.NewUI(c, m)
    defer ui.Close()

    ui.Run(hs, hi)
}
