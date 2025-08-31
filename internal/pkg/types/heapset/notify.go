package heapset

import (
	"path/filepath"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
)

func (hs *HeapSet) SetCallbacks(fn1, fn2 Callback) {
	hs.error = fn1
	hs.watch = fn2
}

func (hs *HeapSet) addFile(path string) {
	err := fs.Watcher.Add(filepath.Dir(path))

	if err != nil {
		sys.Error(err)
	}
}

func (hs *HeapSet) watchFiles() {
	for {
		select {
		case ev, ok := <-fs.Watcher.Events:
			if !ok || !ev.Has(fsnotify.Write) {
				continue
			}

			if ev.Name == sys.Log.Name() {
				if hs.error != nil {
					hs.error() // raise error
				}

				continue
			}

			idx, ok := hs.findByPath(ev.Name)

			if ok && idx == atomic.LoadInt32(hs.index) {
				h := hs.atomicGet(idx)
				h.Reload()

				if hs.watch != nil {
					hs.watch() // raise watch
				}

				continue
			}

		case err, ok := <-fs.Watcher.Errors:
			if ok {
				sys.Error(err)
			}
		}
	}
}
