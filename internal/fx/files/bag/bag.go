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

const (
    header = "Forensic Examiner - Evidence Bag"
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

    usr, err := user.Current()

    if err != nil {
        fx.Error(err)
    }

    f := *types.GetFilters()
    t := append([]string{h.String()}, f...)

    var blk []string

    blk = append(blk, strings.Join(t, " > "))    
    blk = append(blk, fmt.Sprintf("%s (%s)", usr.Username, usr.Name))
    blk = append(blk, fmt.Sprintf("%s", time.Now().UTC()))
    blk = append(blk, fmt.Sprintf("%s", time.Now()))
    blk = append(blk, fmt.Sprintf("%x", h.Sha256()))

    _, err = fmt.Fprintln(bag.file, text.Block(blk, -1, header))

    if err != nil {
        fx.Error(err)
    } 

    d := text.Dec(h.Length())

    for _, s := range h.SMap {
        str := string(h.MMap[s.Start:s.End])

        l := fmt.Sprintf(" %0*d %v", d, s.Nr, str)

        _, err = fmt.Fprintln(bag.file, l)

        if err != nil {
            fx.Error(err)
        } 
    }

    _, err = fmt.Fprintln(bag.file, "")

    if err != nil {
        fx.Error(err)
    }

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

    bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

    if err != nil {
        fx.Fatal(err)
    }

    return
}
