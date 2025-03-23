package data

import (
    "os"
    "path/filepath"

    "github.com/cuhsat/cu/pkg/fs"
)

const (
    Stdin = "STDIN"
)

type HeapSet struct {
    heaps []*Heap // set heaps
    index int     // set index
}

func NewHeapSet(path string) *HeapSet {
    hs := HeapSet{
        index: 0,
    }

    if path == "-" {
        path = Stdin
        
        fs.Stdin(path)
    }

    fi, err := os.Stat(path)

    if err != nil {
        fs.Panic(err)
    }

    if !fi.IsDir() {
        hs.heaps = append(hs.heaps, NewHeap(path))

        return &hs
    }

    dir, err := os.ReadDir(path)

    if err != nil {
        fs.Panic(err)
    }
 
    for _, e := range dir {
        if e.IsDir() {
            continue
        }

        f := filepath.Join(path, e.Name())

        // lazy load all but first
        if len(hs.heaps) == 0 {
            hs.heaps = append(hs.heaps, NewHeap(f))
        } else {
            hs.heaps = append(hs.heaps, &Heap{
                Path: f,
                file: nil,
            })
        }
    }

    if len(hs.heaps) == 0 {
        fs.Panic("no files in directory")
    }

    return &hs
}

func (hs *HeapSet) Heap() *Heap {
    return hs.heaps[hs.index]
}

func (hs *HeapSet) Prev() *Heap {
    hs.index -= 1

    if hs.index < 0 {
        hs.index = len(hs.heaps)-1
    }

    return hs.lazyLoad()
}

func (hs *HeapSet) Next() *Heap {
    hs.index += 1

    if hs.index >= len(hs.heaps) {
        hs.index = 0
    }

    return hs.lazyLoad()
}

func (hs *HeapSet) ThrowAway() {
    for _, h := range hs.heaps {
        h.ThrowAway()
    }

    hs.heaps = hs.heaps[:0]
    hs.index = -1
}

func (hs *HeapSet) lazyLoad() *Heap {
    h := hs.Heap()

    if h.file == nil {
        hs.heaps[hs.index] = NewHeap(h.Path)
    }

    return hs.heaps[hs.index]
}
