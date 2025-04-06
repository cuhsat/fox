package heapset

import (
    "path/filepath"

    "github.com/cuhsat/cu/pkg/fs"
    "github.com/cuhsat/cu/pkg/fs/heap"
    "github.com/fsnotify/fsnotify"
)

type Callback func()

func (hs *HeapSet) SetCallback(fn Callback) {
    hs.watcher_fn = fn
}

func (hs *HeapSet) notifyHeap(h *heap.Heap) {
    err := hs.watcher.Add(filepath.Dir(h.Path))

    if err != nil {
        fs.Panic(err)
    }
}

func (hs *HeapSet) notify() {
    for {
        select {
        case err, ok := <-hs.watcher.Errors:
            if !ok {
                continue
            }
            
            fs.Error(err)

        case ev, ok := <-hs.watcher.Events:
            if !ok || !ev.Has(fsnotify.Write) {
                continue
            }

            for i, h := range hs.heaps {
                if h.Path != ev.Name {
                    continue
                }

                h.Reload()

                if hs.watcher_fn != nil && hs.index == i {
                    hs.watcher_fn()
                }

                break
            }
        }
    }
}
