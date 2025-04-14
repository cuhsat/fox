package heapset

import (
    "os"
    "path/filepath"
    "slices"

    "github.com/cuhsat/fx/internal/sys"
    "github.com/cuhsat/fx/internal/sys/heap"
    "github.com/cuhsat/fx/internal/sys/types"
    "github.com/cuhsat/fx/internal/sys/types/gzip"
    "github.com/cuhsat/fx/internal/sys/types/json"
    "github.com/fsnotify/fsnotify"
)

type HeapSet struct {
    watcher *fsnotify.Watcher // file watcher
    watcher_fn Callback       // file watcher callback

    heaps   []*heap.Heap      // set heaps
    index   int               // set index
}

func NewHeapSet(p []string) *HeapSet {
    w, err := fsnotify.NewWatcher()

    if err != nil {
        sys.Fatal(err)
    }

    hs := HeapSet{
        watcher: w,
        index: 0,
    }

    go hs.notify()

    for _, pe := range p {
        if pe == "-" {
            hs.loadPipe()
            break
        }

        m, err := filepath.Glob(pe)

        if err != nil {
            sys.Fatal(err)
        }

        for _, me := range m {
            hs.loadPath(me)
        }
    }

    if len(hs.heaps) == 0 {
        sys.Fatal("no files in directory")
    }

    hs.load()

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

    return hs.load()
}

func (hs *HeapSet) NextHeap() *heap.Heap {
    hs.index += 1

    if hs.index >= len(hs.heaps) {
        hs.index = 0
    }

    return hs.load()
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

        if h.Flag != heap.Regular {
            os.Remove(h.Path)
        }
    }

    hs.heaps = hs.heaps[:0]
    hs.index = -1
}

func (hs *HeapSet) loadPath(p string) {
    fi, err := os.Stat(p)

    if err != nil {
        sys.Fatal(err)
    }

    if !fi.IsDir() {
        hs.loadFile(p)
    } else {
        hs.loadDir(p)
    }
}

func (hs *HeapSet) loadPipe() {
    hs.heaps = append(hs.heaps, &heap.Heap{
        Path: sys.Stdin(),
        Flag: heap.StdIn,
    })  
}

func (hs *HeapSet) loadFile(p string) {
    var fn types.Format

    b, f, fn := p, heap.Regular, nil

    if gzip.Detect(p) {
        p = gzip.Deflate(p, sys.TempFile("gzip"))
        f = heap.Deflate
    }

    if json.Detect(p) {
        fn = json.Pretty
    }

    hs.heaps = append(hs.heaps, &heap.Heap{
        Path: p,
        Base: b,
        Flag: f,
        Fmt: fn,
    })
}

func (hs *HeapSet) loadDir(p string) {
    dir, err := os.ReadDir(p)

    if err != nil {
        sys.Fatal(err)
    }
 
    for _, e := range dir {
        if !e.IsDir() {
            hs.loadFile(filepath.Join(p, e.Name()))
        }
    }
}

func (hs *HeapSet) load() *heap.Heap {
    h := hs.heaps[hs.index]

    if !h.Loaded() {
        h.Reload()

        hs.notifyHeap(h)
    }

    h.ApplyFilters()

    return h
}
