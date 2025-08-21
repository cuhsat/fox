package heapset

import (
	"path/filepath"
	"sync/atomic"

	mem "github.com/cuhsat/memfile"
	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

func (hs *HeapSet) SetCallbacks(fn1, fn2 Callback) {
	hs.error = fn1
	hs.watch = fn2
}

func (hs *HeapSet) notifyHeap(name string) {
	if name == sys.Log.Name() {
		if hs.error != nil {
			hs.error() // raise error
		}

		return
	}

	idx, ok := hs.findByPath(name)

	if ok && idx == atomic.LoadInt32(hs.index) {
		h := hs.atomicGet(idx)
		h.Reload()

		if hs.watch != nil {
			hs.watch() // raise watch
		}
	}
}

func (hs *HeapSet) watchFile(name string) {
	switch f := sys.Open(name); f.(type) {

	// regular file
	case nil:
		if err := hs.watcher.Add(filepath.Dir(name)); err != nil {
			sys.Error(err)
		}

	// memory file
	case *mem.File:
		f.(*mem.File).SetNotify(func(name string) {
			hs.notifyHeap(name)
		})
	}
}

func (hs *HeapSet) watchFiles() {
	for {
		select {
		case ev, ok := <-hs.watcher.Events:
			if ok && ev.Has(fsnotify.Write) {
				hs.notifyHeap(ev.Name)
			}

		case err, ok := <-hs.watcher.Errors:
			if ok {
				sys.Error(err)
			}
		}
	}
}
