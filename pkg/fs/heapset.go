package fs

import (
    "os"
    "path/filepath"
)

type HeapSet struct {
    heaps []*Heap // set heaps
    index int     // set index
}

func NewHeapSet(path string) *HeapSet {
    hs := HeapSet{
        index: 0,
    }

    fi, err := os.Stat(path)

    if err != nil {
        Panic(err)
    }

    if !fi.IsDir() {
        hs.heaps = append(hs.heaps, NewHeap(path))

        return &hs
    }

    dir, err := os.ReadDir(path)

    if err != nil {
        Panic(err)
    }
 
    for _, e := range dir {
        if e.IsDir() {
            continue
        }

        f := filepath.Join(path, e.Name())

        // lazy load but first
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
        Panic("no files in directory")
    }

    return &hs
}

func (hs *HeapSet) Heap() *Heap {
    return hs.heaps[hs.index]
}

func (hs *HeapSet) Cycle() {
    hs.index = (hs.index+1) % len(hs.heaps)

    h := hs.Heap()

    // lazy loading
    if h.file == nil {
        hs.heaps[hs.index] = NewHeap(h.Path)
    }
}

func (hs *HeapSet) ThrowAway() {
    for _, h := range hs.heaps {
        h.ThrowAway()
    }

    hs.heaps = hs.heaps[:0]
    hs.index = -1
}
