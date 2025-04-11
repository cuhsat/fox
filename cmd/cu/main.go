package main

import (
    "flag"
    "os"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/heapset"
    "github.com/cuhsat/cu/pkg/fs/history"
    "github.com/cuhsat/cu/pkg/fs/limit"
    "github.com/cuhsat/cu/pkg/ui"
    "github.com/cuhsat/cu/pkg/ui/mode"
    "github.com/cuhsat/cu/pkg/ui/status"
)

func usage() {
    fs.Usage("usage: cu [-r] [-h | -t] [-n # | -c #] [-x | -e PATTERN] [-o FILE] [- | PATH ...]")
}

func main() {
    var c limit.Count
    var e heap.Filters

    // flags
    m := status.DefaultMode
    r := flag.Bool("r", false, "Raw mode")
    x := flag.Bool("x", false, "Hex mode")

    // limits
    h := flag.Bool("h", false, "Head limit")
    t := flag.Bool("t", false, "Tail limit")

    // output
    o := flag.String("o", "", "Output file")

    // counts
    flag.IntVar(&c.Lines, "n", 0, "Lines count")
    flag.IntVar(&c.Bytes, "c", 0, "Bytes count")

    // filters
    flag.Var(&e, "e", "Pattern")

    flag.Usage = usage
    flag.Parse()

    a := flag.Args()

    if len(a) == 0 {
        a = append(a, ".")
    }

    if *h && *t {
        fs.Usage("head or tail")
    }

    if c.Lines > 0 && c.Bytes > 0 {
        fs.Usage("lines or bytes")
    }

    if !*x && c.Bytes > 0 {
        fs.Usage("bytes needs hex")
    }

    if *x && len(e) > 0 {
        fs.Usage("hex or pattern")
    }

    if *h {
        limit.SetHead(c)
    }

    if *t {
        limit.SetTail(c)
    }

    if *x {
        m = mode.Hex
    }

    if len(e) > 0 {
        m = mode.Grep
    }

    if len(*o) > 0 {
        *r = true
    }

    hs := heapset.NewHeapSet(a, e...)
    defer hs.ThrowAway()

    if fs.IsCharDev(os.Stdout) || *r {
        hs.Print(*o, *x)
        os.Exit(0)
    }

    hi := history.NewHistory()
    defer hi.Close()

    ui := ui.NewUI(m)
    defer ui.Close()

    ui.Run(hs, hi)
}
