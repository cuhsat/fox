package main

import (
    "flag"
    "io"
    "os"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/data"
    "github.com/cuhsat/cu/pkg/ui"
)

var Version string = "dev"

func main() {
    h := flag.Bool("h", false, "Show help")
    v := flag.Bool("v", false, "Show version")

    flag.CommandLine.SetOutput(io.Discard)
    flag.Parse()

    if *h || len(os.Args) < 2 {
        fs.Usage("cu [-hv] PATH ...")
    }

    if *v {
        fs.Print("cu", Version)
    }

    hs := data.NewHeapSet(os.Args[1:])
    defer hs.ThrowAway()

    hi := fs.NewHistory()
    defer hi.Close()

    ui := ui.NewUI()
    defer ui.Close()

    ui.Run(hs, hi)
}
