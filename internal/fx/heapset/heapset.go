package heapset

import (
    "os"
    "path/filepath"
    "slices"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/cuhsat/fx/internal/fx/types/gzip"
    "github.com/cuhsat/fx/internal/fx/types/jsonl"
    "github.com/fsnotify/fsnotify"
)

type callback func()

type HeapSet struct {
    watch *fsnotify.Watcher // file watcher
    watch_fn callback       // file watcher callback
    error_fn callback       // error callback

    heaps []*heap.Heap      // set heaps
    index int               // set index
}

func New(paths []string) *HeapSet {
    w, err := fsnotify.NewWatcher()

    if err != nil {
        fx.Error(err)
    }

    hs := HeapSet{
        watch: w,
        index: 0,
    }

    go hs.notify()

    hs.watchPath(fx.Logfile)

    for _, path := range paths {
        if path == "-" {
            hs.loadPipe()
            break
        }

        match, err := filepath.Glob(path)

        if err != nil {
            fx.Error(err)
        }

        for _, m := range match {
            hs.loadPath(m)
        }
    }

    if len(hs.heaps) == 0 {
        fx.Fatal("no files found")
    }

    hs.load()

    return &hs
}

func (hs *HeapSet) Bind(fn1, fn2 callback) {
    hs.watch_fn = fn1
    hs.error_fn = fn2
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
    hs.watch.Close()

    for _, h := range hs.heaps {
        h.ThrowAway()

        if h.Type > types.Regular {
            os.Remove(h.Path)
        }
    }

    hs.heaps = hs.heaps[:0]
    hs.index = -1
}

func (hs *HeapSet) loadPath(path string) {
    fi, err := os.Stat(path)

    if err != nil {
        fx.Error(err)
        return
    }

    if !fi.IsDir() {
        hs.loadFile(path)
    } else {
        hs.loadDir(path)
    }
}

func (hs *HeapSet) loadPipe() {
    hs.heaps = append(hs.heaps, &heap.Heap{
        Path: fx.Stdin(),
        Type: types.StdIn,
    })  
}

func (hs *HeapSet) loadFile(path string) {
    var fn types.Format

    base, tp := path, types.Regular

    if gzip.Detect(path) {
        path = gzip.Deflate(path)
        tp = types.Deflate
    }

    if jsonl.Detect(path) {
        fn = jsonl.Pretty
    }

    hs.heaps = append(hs.heaps, &heap.Heap{
        Path: path,
        Base: base,
        Type: tp,
        Fmt: fn,
    })
}

func (hs *HeapSet) loadDir(path string) {
    dir, err := os.ReadDir(path)

    if err != nil {
        fx.Error(err)
        return
    }
 
    for _, f := range dir {
        if !f.IsDir() {
            hs.loadFile(filepath.Join(path, f.Name()))
        }
    }
}

func (hs *HeapSet) load() *heap.Heap {
    h := hs.heaps[hs.index]

    if !h.Loaded() {
        h.Reload()

        hs.watchHeap(h)
    }

    h.ApplyFilters()

    return h
}
