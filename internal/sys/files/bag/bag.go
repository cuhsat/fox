package bag

import (
    "fmt"
    "os"

    "github.com/cuhsat/fx/internal/sys"
    "github.com/cuhsat/fx/internal/sys/heap"
    "github.com/cuhsat/fx/internal/sys/text"
)

const (
    File = "EVIDENCE"
)

type Bag struct {
    file *os.File // file handle
}

func NewBag(p string) *Bag {
    if len(p) == 0 {
        p = File
    }

    f, err := os.OpenFile(p, sys.O_EVIDENCE, 0644)

    if err != nil {
        sys.Fatal(err)
    }

    return &Bag{
        file: f,
    }
}

func (b *Bag) Put(h *heap.Heap) {
    t := text.Block(h.Base, 78)

    _, err := b.file.WriteString(fmt.Sprintf("%s\n", t))

    if err != nil {
        sys.Fatal(err)
    }

    for _, s := range h.SMap {
        _, err := b.file.Write(h.MMap[s.Start:s.End])

        if err != nil {
            sys.Fatal(err)
        }

        _, err = b.file.Write([]byte{'\n'})

        if err != nil {
            sys.Fatal(err)
        }
    }
}

func (b *Bag) Close() {
    b.file.Close()
}
