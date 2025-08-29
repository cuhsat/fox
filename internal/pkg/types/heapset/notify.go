package heapset

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

var watched = make(map[sys.File]time.Time)

func (hs *HeapSet) SetCallbacks(fn1, fn2 Callback) {
	hs.error = fn1
	hs.watch = fn2
}

func (hs *HeapSet) notifyHeap(path string) {
	if path == sys.Log.Name() {
		if hs.error != nil {
			hs.error() // raise error
		}

		return
	}

	idx, ok := hs.findByPath(path)

	if ok && idx == atomic.LoadInt32(hs.index) {
		h := hs.atomicGet(idx)
		h.Reload()

		if hs.watch != nil {
			hs.watch() // raise watch
		}
	}
}

func (hs *HeapSet) watchFile(path string) {
	switch f := sys.Open(path); f.(type) {

	// regular file
	case *os.File:
		if err := hs.watcher.Add(filepath.Dir(path)); err != nil {
			sys.Error(err)
		}

	// memory file
	case sys.File:
		fi, _ := f.Stat()
		watched[f] = fi.ModTime()
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

func (hs *HeapSet) pollFiles() {
	for {
		for f, t := range watched {
			fi, _ := f.Stat()

			if fi.ModTime().After(t) {
				watched[f] = t
				hs.notifyHeap(f.Name())
			}
		}

		time.Sleep(time.Millisecond * 200)
	}
}
