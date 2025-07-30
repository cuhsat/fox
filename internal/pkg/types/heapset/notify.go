package heapset

import (
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"

	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types/file"
)

func (hs *HeapSet) watch(name string) {
	switch f := file.Open(name); f.(type) {

	// virtual file
	case *file.Data:
		f.(*file.Data).Watch(hs.watcher.Events)

	// regular file
	case nil:
		err := hs.watcher.Add(filepath.Dir(name))

		if err != nil {
			sys.Error(err)
		}
	}
}

func (hs *HeapSet) notify() {
	for {
		select {
		case err, ok := <-hs.watcher.Errors:
			if !ok {
				continue
			}

			sys.Error(err)

		case ev, ok := <-hs.watcher.Events:
			if !ok || !ev.Has(fsnotify.Write) {
				continue
			}

			if ev.Name == sys.Log.Name() {
				if hs.fnError != nil {
					hs.fnError() // bound callback
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

				if hs.fnWatch != nil && idx == i {
					hs.fnWatch() // bound callback
				}

				break
			}

			hs.RUnlock()
		}
	}
}
