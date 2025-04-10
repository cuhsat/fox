package heapset

import (
    "os"
    "path/filepath"
    "slices"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/gzip"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/fsnotify/fsnotify"
)

const (
    Stdin = "-"
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

    for _, pe := range p {
        if pe == Stdin {
            hs.loadPipe()
            break
        }

        m, err := filepath.Glob(pe)

        if err != nil {
            fs.Panic(err)
        }

        for _, me := range m {
            hs.loadPath(me)
        }
    }

    if len(hs.heaps) == 0 {
        fs.Panic("no files in directory")
    }

    hs.loadLazy()

    for _, fe := range f {
        hs.heaps[0].AddFilter(fe)
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
    fi, err := os.Stat(p)

    if err != nil {
        fs.Panic(err)
    }

    if !fi.IsDir() {
        hs.loadFile(p)
    } else {
        hs.loadDir(p)
    }
}

func (hs *HeapSet) loadPipe() {
    hs.heaps = append(hs.heaps, &heap.Heap{
        Title: Stdin,
        Path: fs.Stdin(),
        Flag: heap.StdIn,
    })  
}

func (hs *HeapSet) loadFile(p string) {
    f := heap.Normal
    t := p

    if gzip.Detect(p) {
        p = gzip.Deflate(p, fs.TempFile("gzip"))
        f = heap.Deflate
    }

    hs.heaps = append(hs.heaps, &heap.Heap{
        Title: t,
        Path: p,
        Flag: f,
    })
}

func (hs *HeapSet) loadDir(p string) {
    dir, err := os.ReadDir(p)

    if err != nil {
        fs.Panic(err)
    }
 
    for _, e := range dir {
        if !e.IsDir() {
            hs.loadFile(filepath.Join(p, e.Name()))
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
