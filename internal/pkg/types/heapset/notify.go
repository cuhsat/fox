package heapset

import (
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/cuhsat/memfile"
	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

var (
	files = make(chan string)
)

func (hs *HeapSet) watchFile(name string) {
	switch f := sys.Open(name); f.(type) {

	// virtual file
	case *memfile.FileData:
		f.(*memfile.FileData).Notify(files)

	// regular file
	case nil:
		err := hs.watcher.Add(filepath.Dir(name))

		if err != nil {
			sys.Error(err)
		}
	}
}

func (hs *HeapSet) observe() {
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

			hs.notify(ev.Name)

		case name := <-files:
			hs.notify(name)
		}
	}
}

func (hs *HeapSet) notify(name string) {
	if name == sys.Log.Name() {
		if hs.error != nil {
			hs.error() // bound callback
		}

		return
	}

	hs.RLock()

	for i, h := range hs.heaps {
		if !strings.HasSuffix(h.Path, name) {
			return
		}

		h.Reload()

		idx := int(atomic.LoadInt32(hs.index))

		if hs.watch != nil && idx == i {
			hs.watch() // bound callback
		}

		return
	}

	hs.RUnlock()
}
