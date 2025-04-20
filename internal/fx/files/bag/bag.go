package bag

import (
    "fmt"
    "os"
    "os/user"
    "strings"
    "time"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/text"
    "github.com/cuhsat/fx/internal/fx/types"
)

const (
    filename = "EVIDENCE"
)

type Bag struct {
    Path string   // file path
    file *os.File // file handle
}

func New(path string) *Bag {
    return &Bag{
        Path: path,
    }
}

func (bag *Bag) Put(h *heap.Heap) {
    if bag.file == nil {
        bag.mustInit()
    }

    f := *types.GetFilters()

    d := text.Dec(h.Length())

    usr, err := user.Current()

    if err != nil {
        fx.Error(err)
    }

    for _, l := range [...]string{
        // filters
        fmt.Sprintf("%s > %s", h, strings.Join(f, " > ")),
        
        // username
        fmt.Sprintf("%s (%s)", usr.Username, usr.Name),
        
        // global time
        time.Now().UTC().String(),
        
        // local time
        time.Now().String(),
        
        // file hashsum
        fmt.Sprintf("%x", h.Sha256()),
    } {
        bag.write(fmt.Sprintf("// %s", l))
    }

    bag.write("")

    for _, s := range h.SMap {
        str := string(h.MMap[s.Start:s.End])

        bag.write(fmt.Sprintf("[%0*d] %v", d, s.Nr, str))
    }

    bag.write("")

    return
}

func (bag *Bag) Close() {
    if bag.file == nil {
        bag.file.Close()
    }
}

func (bag *Bag) mustInit() {
    var err error

    if len(bag.Path) == 0 {
        bag.Path = filename
    }

    is := fx.Exists(bag.Path)

    bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

    if err != nil {
        fx.Fatal(err)
    }

    if !is {
        bag.write("// Forensic Examiner - Evidence Bag\n")
    }

    return
}

func (bag *Bag) write(s string) {
    _, err := fmt.Fprintln(bag.file, s)

    if err != nil {
        fx.Error(err)
    }
}
