package bag

import (
    "fmt"
    "os"
    "os/user"
    "strings"
    "time"

    "github.com/cuhsat/fx/internal/sys"
    "github.com/cuhsat/fx/internal/sys/heap"
    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/cuhsat/fx/internal/sys/types"
)

const (
    File = "EVIDENCE"
)

type Bag struct {
    Path string   // file path
    
    file *os.File // file handle
}

func NewBag(p string) *Bag {
    return &Bag{
        Path: p,
    }
}

func (bag *Bag) Init() {
    var err error

    if len(bag.Path) == 0 {
        bag.Path = File
    }

    bag.file, err = os.OpenFile(bag.Path, sys.O_EVIDENCE, 0600)

    if err != nil {
        sys.Fatal(err)
    }
}

func (bag *Bag) Put(h *heap.Heap) {
    if bag.file == nil {
        bag.Init()
    }

    var b []string

    u, err := user.Current()

    if err != nil {
        sys.Fatal(err)
    }

    f := *types.GetFilters()
    t := append([]string{h.String()}, f...)

    b = append(b, strings.Join(t, " > "))
    b = append(b, fmt.Sprintf("%s @ %s", u.Username, time.Now().UTC()))
    b = append(b, fmt.Sprintf("%x", h.Sha256()))

    _, err = fmt.Fprintln(bag.file, text.Block(b, -1))

    if err != nil {
        sys.Error(err)
    } 

    d := text.Dec(h.Length())

    for _, s := range h.SMap {
        str := string(h.MMap[s.Start:s.End])

        l := fmt.Sprintf(" %0*d  %v", d, s.Nr, str)

        _, err = fmt.Fprintln(bag.file, l)

        if err != nil {
            sys.Error(err)
        } 
    }

    _, err = fmt.Fprintln(bag.file, "")
}

func (bag *Bag) Close() {
    if bag.file == nil {
        bag.file.Close()
    }
}
