package main

import (
    "flag"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/config"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/fs/history"
    "github.com/cuhsat/cu/pkg/fs/limit"
    "github.com/cuhsat/cu/pkg/ui"
    "github.com/cuhsat/cu/pkg/ui/mode"
)

func usage() {
    fs.Usage("usage: cu [-h | -t] [-n # | -c #] [-x | -e PATTERN] [PATH ...]")
}

func main() {
    var l limit.Limit
    var c limit.Count
    var e heap.Filters

    // config
    cfg := config.Load()

    // flags
    m := mode.Less
    x := flag.Bool("x", false, "Hex mode")

    // limits
    h := flag.Bool("h", false, "Head limit")
    t := flag.Bool("t", false, "Tail limit")

    // counts
    flag.IntVar(&c.Lines, "n", 0, "Lines count")
    flag.IntVar(&c.Bytes, "c", 0, "Bytes count")

    // filters
    flag.Var(&e, "e", "Pattern")

    flag.Usage = usage
    flag.Parse()

    args := flag.Args()

    if len(args) == 0 {
        args = append(args, ".")
    }

    if *h && *t {
        fs.Usage("head or tail")
    }

    if c.Lines > 0 && c.Bytes > 0 {
        fs.Usage("lines or bytes")
    }

    if *x && len(e) > 0 {
        fs.Usage("hex or pattern")
    }

    if *h {
        l.Head = c
    }

    if *t {
        l.Tail = c
    }

    if *x {
        m = mode.Hex
    }

    if len(e) > 0 {
        m = mode.Grep
    } 

    hs := heapset.NewHeapSet(args, l, e...)
    defer hs.ThrowAway()

    hi := history.NewHistory()
    defer hi.Close()

    ui := ui.NewUI(cfg, m)
    defer ui.Close()

    ui.Run(hs, hi)
}
