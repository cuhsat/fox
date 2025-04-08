package heapset

import (
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/config"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/cuhsat/cu/pkg/fs/limit"
    "github.com/fsnotify/fsnotify"
)

type action func(h *heap.Heap) string

type HeapSet struct {
    watcher *fsnotify.Watcher // file watcher
    watcher_fn Callback       // file watcher callback

    config  config.Config     // user config
    limit   limit.Limit       // heap limit
    
    heaps   []*heap.Heap      // set heaps
    index   int               // set index
}

func NewHeapSet(c config.Config, l limit.Limit, p []string, f ...string) *HeapSet {
    w, err := fsnotify.NewWatcher()

    if err != nil {
        fs.Panic(err)
    }

    hs := HeapSet{
        watcher: w,
        config: c,
        limit: l,
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

func (hs *HeapSet) Counts() {
    hs.buffer("wc", func(h *heap.Heap) string {
        return fmt.Sprintf("%8d %8d %s\n", h.Length(), len(h.MMap), h.Path)
    })
}

func (hs *HeapSet) Hashes() {
    hs.buffer(fmt.Sprintf("%ssum", hs.config.CU.Hash), func(h *heap.Heap) string {
        return fmt.Sprintf("%x  %s\n", h.Hash(hs.config.CU.Hash), h.Path)
    })
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

func (hs *HeapSet) buffer(t string, fn action) {
    f := fs.Stdout()

    for _, h := range hs.heaps {
        if h.Flag != heap.Normal {
            continue
        }

        if !h.Loaded() {
            h.Reload()
        }

        _, err := io.WriteString(f, fn(h))

        if err != nil {
            fs.Panic(err)
        }
    }

    f.Close()

    hs.heaps = append(hs.heaps, &heap.Heap{
        Title: t,
        Flag: heap.StdOut,
        Path: f.Name(),
    })

    hs.index = len(hs.heaps)-1

    hs.loadLazy()
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
            Limit: hs.limit,
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
                Limit: hs.limit,
            })
        }
    }
}

func (hs *HeapSet) loadLazy() *heap.Heap {
    h := hs.heaps[hs.index]

    if !h.Loaded() {
        h.Reload()

        hs.notifyHeap(h)

        hs.heaps[hs.index] = h
    }

    return h
}
