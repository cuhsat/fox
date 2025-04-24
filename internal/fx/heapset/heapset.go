package heapset

import (
    "os"
    "path/filepath"
    "slices"
    "sync"
    "sync/atomic"

    "github.com/bmatcuk/doublestar/v4"
    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/file"
    "github.com/cuhsat/fx/internal/fx/file/gzip"
    "github.com/cuhsat/fx/internal/fx/file/tar"
    "github.com/cuhsat/fx/internal/fx/file/zip"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/types"
    "github.com/fsnotify/fsnotify"
)

type callback func()

type HeapSet struct {
    sync.RWMutex

    watch *fsnotify.Watcher // file watcher
    watch_fn callback       // file watcher callback
    error_fn callback       // error callback

    heaps []*heap.Heap      // set heaps
    index *int32            // set index
}

func New(paths []string) *HeapSet {
    w, err := fsnotify.NewWatcher()

    if err != nil {
        fx.Error(err)
    }

    hs := HeapSet{
        watch: w,
        index: new(int32),
    }

    go hs.notify()

    hs.watchPath(fx.Log.Name)

    for _, path := range paths {
        if path == "-" {
            hs.loadPipe()
            break
        }

        match, err := doublestar.FilepathGlob(path)

        if err != nil {
            fx.Error(err)
        }

        for _, m := range match {
            hs.loadPath(m)
        }
    }

    if hs.Length() == 0 {
        fx.Exit("no files found")
    }

    hs.load()

    return &hs
}

func (hs *HeapSet) Bind(fn1, fn2 callback) {
    hs.watch_fn = fn1
    hs.error_fn = fn2
}

func (hs *HeapSet) Length() int32 {
    hs.RLock()
    defer hs.RUnlock()

    return int32(len(hs.heaps))
}

func (hs *HeapSet) Current() (int32, *heap.Heap) {
    idx := atomic.LoadInt32(hs.index)

    return idx+1, hs.get(idx)
}

func (hs *HeapSet) OpenHeap(path string) {
    if !fx.Exists(path) {
        return
    }

    idx := hs.findByPath(path)

    if idx < 0 {
        idx = hs.Length()

        hs.loadPath(path)
    }

    atomic.StoreInt32(hs.index, idx)

    hs.load()
}

func (hs *HeapSet) PrevHeap() *heap.Heap {
    idx := atomic.AddInt32(hs.index, -1)

    if idx < 0 {
        atomic.StoreInt32(hs.index, hs.Length()-1)
    }

    return hs.load()
}

func (hs *HeapSet) NextHeap() *heap.Heap {
    idx := atomic.AddInt32(hs.index, 1)

    if idx >= hs.Length() {
        atomic.StoreInt32(hs.index, 0)
    }

    return hs.load()
}

func (hs *HeapSet) CloseHeap() *heap.Heap {
    if hs.Length() == 1 {
        return nil
    }

    h := hs.get(atomic.LoadInt32(hs.index))

    hs.unload(h)

    hs.del()

    atomic.AddInt32(hs.index, -1)

    return hs.NextHeap()
}

func (hs *HeapSet) CloseHeaps() *heap.Heap {
    // hs.Lock()

    // for _, h := range hs.heaps {
    //     if len(h.SMap) == 0 {
    //         hs.CloseHeap()
    //     }
    // }

    // hs.Unlock()

    // if hs.Length() == 1 {
    //     return nil
    // }

    return hs.NextHeap()
}

func (hs *HeapSet) ThrowAway() {
    hs.watch.Close()

    hs.Lock()

    for _, h := range hs.heaps {
        hs.unload(h)
    }

    hs.heaps = hs.heaps[:0]

    hs.Unlock()

    atomic.AddInt32(hs.index, -1)
}

func (hs *HeapSet) findByPath(path string) int32 {
    hs.RLock()
    defer hs.RUnlock()

    for i, h := range hs.heaps {
        if h.Base == path {
            return int32(i)
        }
    }

    return -1
}

func (hs *HeapSet) findByName(name string) int32 {
    hs.RLock()
    defer hs.RUnlock()

    for i, h := range hs.heaps {
        if h.Title == name {
            return int32(i)
        }
    }

    return -1
}

func (hs *HeapSet) loadPath(path string) {
    fi, err := os.Stat(path)

    if err != nil {
        fx.Error(err)
        return
    }

    if fi.IsDir() {
        hs.loadDir(path)
        return
    }

    base := path

    if gzip.Detect(path) {
        path = gzip.Deflate(path)
    }

    if tar.Detect(path) {
        hs.loadTar(path, base)
        return
    }

    if zip.Detect(path) {
        hs.loadZip(path, base)
        return
    }

    hs.loadFile(path, base)
}

func (hs *HeapSet) loadPipe() {
    pipe := fx.Stdin()

    hs.add(&heap.Heap{
        Path: pipe,
        Base: pipe,
        Type: types.Stdin,
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
            hs.loadPath(filepath.Join(path, f.Name()))
        }
    }
}

func (hs *HeapSet) loadTar(path, base string) {
    for _, fe := range tar.Deflate(path) {
        hs.loadEntry(fe, base)
    }
}

func (hs *HeapSet) loadZip(path, base string) {
    for _, fe := range zip.Deflate(path) {
        hs.loadEntry(fe, base)
    }
}

func (hs *HeapSet) loadFile(path, base string) {
    h := &heap.Heap{
        Title: base,
        Path: path,
        Base: base,
        Type: types.Regular,
    }

    if path != base {
        h.Type = types.Deflate
    }

    hs.add(h)
}

func (hs *HeapSet) loadEntry(e *file.Entry, base string) {
    hs.add(&heap.Heap{
        Title: filepath.Join(base, e.Name),
        Path: e.Path,
        Base: e.Path,
        Type: types.Deflate,
    })
}

func (hs *HeapSet) load() *heap.Heap {
    h := hs.get(atomic.LoadInt32(hs.index))

    if !h.Loaded() {
        h.Reload()

        hs.watchHeap(h)
    }

    h.ApplyFilters()

    return h
}

func (hs *HeapSet) unload(h *heap.Heap) {
    h.ThrowAway()

    if h.Type > types.Regular {
        os.Remove(h.Path)
    }
}

func (hs *HeapSet) get(idx int32) *heap.Heap {
    hs.RLock()
    defer hs.RUnlock()

    return hs.heaps[idx]
}

func (hs *HeapSet) add(h *heap.Heap) {
    hs.Lock()

    hs.heaps = append(hs.heaps, h)

    hs.Unlock()
}

func (hs *HeapSet) del() {
    hs.Lock()

    idx := int(atomic.LoadInt32(hs.index))

    hs.heaps = slices.Delete(hs.heaps, idx, idx+1)

    hs.Unlock()
}
