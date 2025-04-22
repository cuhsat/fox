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
    "github.com/cuhsat/fx/internal/fx/types/smap"
)

const (
    filename = "EVIDENCE"
)

type writer interface {
    WriteTitle(s string)
    WriteMetas(p, f, u string, t, l time.Time, b []byte)
    WriteLines(smap smap.SMap)
}

type Bag struct {
    Path string   // file path
    file *os.File // file handle

    w writer      // writer
}

func New(path string, md bool) *Bag {
    if len(path) == 0 {
        path = filename
    }

    return &Bag{
        Path: path,
    }
}

func (bag *Bag) Put(h *heap.Heap) bool {
    if bag.file == nil && !bag.init() {
        return false
    }

    f := *types.GetFilters()

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
        bag.write(fmt.Sprintf("%s", l))
    }

    bag.write("")

    d := text.Dec(h.Length())

    for _, s := range h.SMap {
        str := string(h.MMap[s.Start:s.End])

        // line
        bag.write(fmt.Sprintf("%0*d  %v", d, s.Nr, str))
    }

    bag.write("")

    return true
}

func (bag *Bag) Close() {
    if bag.file == nil {
        bag.file.Close()
    }
}

func (bag *Bag) init() bool {
    var err error

    is := fx.Exists(bag.Path)

    bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

    if err != nil {
        fx.Error(err)
        return false
    }

    if !is {
        bag.write("FORENSIC EXAMINER EVIDENCE BAG\n")
    }

    return true
}

func (bag *Bag) write(s string) {
    _, err := fmt.Fprintln(bag.file, s)

    if err != nil {
        fx.Error(err)
    }
}
