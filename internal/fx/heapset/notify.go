package heapset

import (
    "path/filepath"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/fsnotify/fsnotify"
)

func (hs *HeapSet) notifyHeap(h *heap.Heap) {
    err := hs.watch.Add(filepath.Dir(h.Path))

    if err != nil {
        fx.Error(err)
    }
}

func (hs *HeapSet) notify() {
    for {
        select {
        case err, ok := <-hs.watch.Errors:
            if !ok {
                continue
            }
            
            fx.Error(err)

        case ev, ok := <-hs.watch.Events:
            if !ok || !ev.Has(fsnotify.Write) {
                continue
            }

            for i, h := range hs.heaps {
                if h.Path != ev.Name {
                    continue
                }

                h.Reload()

                if hs.watch_fn != nil && hs.index == i {
                    hs.watch_fn() // callback
                }

                break
            }
        }
    }
}
