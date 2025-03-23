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

func NewHeapSet(paths []string) *HeapSet {
    hs := HeapSet{
        index: 0,
    }

    for _, path := range paths {
        hs.loadPath(path)
    }

    if len(hs.heaps) == 0 {
        fs.Panic("no files in directory")
    }

    hs.loadLazy()

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

    return hs.loadLazy()
}

func (hs *HeapSet) Next() *Heap {
    hs.index += 1

    if hs.index >= len(hs.heaps) {
        hs.index = 0
    }

    return hs.loadLazy()
}

func (hs *HeapSet) ThrowAway() {
    for _, h := range hs.heaps {
        h.ThrowAway()
    }

    hs.heaps = hs.heaps[:0]
    hs.index = -1
}

func (hs *HeapSet) loadPath(p string) {
    // read stdin
    if p == "-" {
        p = Stdin
        
        fs.Stdin(p)
    }

    fi, err := os.Stat(p)

    if err != nil {
        fs.Panic(err)
    }

    // load file
    if !fi.IsDir() {
        hs.heaps = append(hs.heaps, &Heap{
            Path: p,
            file: nil,
        })

        return
    }

    // load directory
    dir, err := os.ReadDir(p)

    if err != nil {
        fs.Panic(err)
    }
 
    for _, e := range dir {
        if !e.IsDir() {
            hs.heaps = append(hs.heaps, &Heap{
                Path: filepath.Join(p, e.Name()),
                file: nil,
            })
        }
    }
}

func (hs *HeapSet) loadLazy() *Heap {
    h := hs.Heap()

    if h.file == nil {
        hs.heaps[hs.index] = NewHeap(h.Path)
    }

    return hs.heaps[hs.index]
}
