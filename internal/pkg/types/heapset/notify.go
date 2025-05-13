package heapset

import (
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

func (hs *HeapSet) watchHeap(h *heap.Heap) {
	hs.watchPath(h.Path)
}

func (hs *HeapSet) watchPath(path string) {
	err := hs.watch.Add(filepath.Dir(path))

	if err != nil {
		sys.Error(err)
	}
}

func (hs *HeapSet) notify() {
	for {
		select {
		case err, ok := <-hs.watch.Errors:
			if !ok {
				continue
			}

			sys.Error(err)

		case ev, ok := <-hs.watch.Events:
			if !ok || !ev.Has(fsnotify.Write) {
				continue
			}

			if ev.Name == sys.Log.Name {
				if hs.errorFn != nil {
					hs.errorFn() // bound callback
				}

				continue
			}

			hs.RLock()

			for i, h := range hs.heaps {
				if !strings.HasSuffix(h.Path, ev.Name) {
					continue
				}

				h.Reload()

				idx := int(atomic.LoadInt32(hs.index))

				if hs.watchFn != nil && idx == i {
					hs.watchFn() // bound callback
				}

				break
			}

			hs.RUnlock()
		}
	}
}
