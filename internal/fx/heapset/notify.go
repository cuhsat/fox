package heapset

import (
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/fsnotify/fsnotify"
)

func (hs *HeapSet) watchHeap(h *heap.Heap) {
    hs.watchPath(h.Path)
}

func (hs *HeapSet) watchPath(path string) {
    err := hs.watch.Add(filepath.Dir(path))

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

            if ev.Name == fx.Log.Name {
                if hs.error_fn != nil {
                    hs.error_fn() // bound callback
                }

                continue
            }

            hs.Lock()

            for i, h := range hs.heaps {
                if !strings.HasSuffix(h.Path, ev.Name) {
                    continue
                }

                h.Reload()

                if hs.watch_fn != nil && hs.index == i {
                    hs.watch_fn() // bound callback
                }

                break
            }

            hs.Unlock()
        }
    }
}
