package heapset

// import (
// 	"path/filepath"
// 	"strings"
// 	"sync/atomic"

// 	"github.com/cuhsat/memfile"
// 	"github.com/fsnotify/fsnotify"

// 	"github.com/cuhsat/fox/internal/pkg/sys"
// )

// type Callback func()

// type Watcher struct {
// 	watcher *fsnotify.Watcher // file watcher
// 	files   chan string

// 	watch Callback // watcher callback
// 	error Callback // error callback
// }

// func NewWatcher() *Watcher {
// 	w, err := fsnotify.NewWatcher()

// 	if err != nil {
// 		sys.Error(err)
// 	}

// 	return &Watcher{
// 		watcher: w,
// 		files:   make(chan string),
// 	}
// }

// func (w *Watcher) watchFile(name string) {
// 	switch f := sys.Open(name); f.(type) {

// 	// virtual file
// 	case *memfile.FileData:
// 		f.(*memfile.FileData).Notify(files)

// 	// regular file
// 	case nil:
// 		err := w.watcher.Add(filepath.Dir(name))

// 		if err != nil {
// 			sys.Error(err)
// 		}
// 	}
// }

// func (w *Watcher) observe() {
// 	for {
// 		select {
// 		case err, ok := <-w.watcher.Errors:
// 			if !ok {
// 				continue
// 			}

// 			sys.Error(err)

// 		case ev, ok := <-w.watcher.Events:
// 			if !ok || !ev.Has(fsnotify.Write) {
// 				continue
// 			}

// 			w.notify(ev.Name)

// 		case name := <-files:
// 			w.notify(name)
// 		}
// 	}
// }

// func (w *Watcher) notify(name string) {
// 	if name == sys.Log.Name() {
// 		if w.error != nil {
// 			w.error() // bound callback
// 		}

// 		return
// 	}

// 	hs.RLock()

// 	for i, h := range hs.heaps {
// 		if !strings.HasSuffix(h.Path, name) {
// 			return
// 		}

// 		h.Reload()

// 		idx := int(atomic.LoadInt32(hs.index))

// 		if hs.watch != nil && idx == i {
// 			hs.watch() // bound callback
// 		}

// 		return
// 	}

// 	hs.RUnlock()
// }
