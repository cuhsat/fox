package heapset

import (
    "os"
    "path/filepath"
    "slices"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/fsnotify/fsnotify"
)

type HeapSet struct {
    watcher *fsnotify.Watcher // file watcher
    watcher_fn Callback       // file watcher callback

    heaps   []*heap.Heap      // set heaps
    index   int               // set index
}

func NewHeapSet(p []string, f ...string) *HeapSet {
    w, err := fsnotify.NewWatcher()

    if err != nil {
        fs.Panic(err)
    }

    hs := HeapSet{
        watcher: w,
        index: 0,
    }

    go hs.notify()

    for _, pp := range p {
        hs.loadPath(pp)
    }

    if len(hs.heaps) == 0 {
        fs.Panic("no files in directory")
    }

    hs.loadLazy()

    for _, ff := range f {
        hs.heaps[0].AddFilter(ff)
    }

    return &hs
}

func (hs *HeapSet) Length() int {
    return len(hs.heaps)
}

func (hs *HeapSet) Current() (int, *heap.Heap) {
    return hs.index+1, hs.heaps[hs.index]
}

func (hs *HeapSet) PrevHeap() *heap.Heap {
    hs.index -= 1

    if hs.index < 0 {
        hs.index = len(hs.heaps)-1
    }

    return hs.loadLazy()
}

func (hs *HeapSet) NextHeap() *heap.Heap {
    hs.index += 1

    if hs.index >= len(hs.heaps) {
        hs.index = 0
    }

    return hs.loadLazy()
}

func (hs *HeapSet) CloseHeap() *heap.Heap {
    if len(hs.heaps) == 1 {
        return nil
    }    

    hs.heaps = slices.Delete(hs.heaps, hs.index, hs.index+1)
    hs.index -= 1

    return hs.NextHeap()
}

func (hs *HeapSet) ThrowAway() {
    hs.watcher.Close()

    for _, h := range hs.heaps {
        h.ThrowAway()

        if h.Flag != heap.Normal {
            os.Remove(h.Path)
        }
    }

    hs.heaps = hs.heaps[:0]
    hs.index = -1
}

func (hs *HeapSet) loadPath(p string) {
    var f heap.Flag

    // read stdin
    if p == "-" {
        p, f = fs.Stdin(), heap.StdIn
    }

    fi, err := os.Stat(p)

    if err != nil {
        fs.Panic(err)
    }

    // load file
    if !fi.IsDir() {
        hs.heaps = append(hs.heaps, &heap.Heap{
            Path: p,
            Flag: f,
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
            hs.heaps = append(hs.heaps, &heap.Heap{
                Path: filepath.Join(p, e.Name()),
            })
        }
    }
}

func (hs *HeapSet) loadLazy() *heap.Heap {
    h := hs.heaps[hs.index]

    if !h.Loaded() {
        h.Reload()

        hs.notifyHeap(h)
    }

    return h
}
