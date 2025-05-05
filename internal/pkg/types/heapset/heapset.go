package heapset

import (
	"slices"
	"sync"
	"sync/atomic"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/heap"
)

type callback func()

type HeapSet struct {
	sync.RWMutex

	watch    *fsnotify.Watcher // file watcher
	watch_fn callback          // file watcher callback
	error_fn callback          // error callback

	heaps []*heap.Heap // set heaps
	index *int32       // set index
}

func New(paths []string) *HeapSet {
	w, err := fsnotify.NewWatcher()

	if err != nil {
		sys.Error(err)
	}

	hs := HeapSet{
		watch: w,
		index: new(int32),
	}

	go hs.notify()

	hs.watchPath(sys.Log.Name)

	for _, path := range paths {
		if path == "-" {
			hs.loadPipe()
			break
		}

		hs.Open(path)
	}

	if hs.Size() == 0 {
		sys.Exit("no files found")
	}

	hs.load()

	return &hs
}

func (hs *HeapSet) Bind(fn1, fn2 callback) {
	hs.watch_fn = fn1
	hs.error_fn = fn2
}

func (hs *HeapSet) Size() int32 {
	hs.RLock()
	defer hs.RUnlock()
	return int32(len(hs.heaps))
}

func (hs *HeapSet) Files() []string {
	hs.RLock()

	fs := make([]string, 0, len(hs.heaps))

	for _, h := range hs.heaps {
		if h.Type == types.Regular {
			fs = append(fs, h.Path)
		}
	}

	hs.RUnlock()

	return fs
}

func (hs *HeapSet) Heap() (int32, *heap.Heap) {
	idx := atomic.LoadInt32(hs.index)
	return idx + 1, hs.atomicGet(idx)
}

func (hs *HeapSet) Open(path string) {
	match, err := doublestar.FilepathGlob(path)

	if err != nil {
		sys.Error(err)
	}

	for _, m := range match {
		hs.loadPath(m)
	}
}

func (hs *HeapSet) OpenLog() {
	idx := hs.findByPath(sys.Log.Name)

	if idx < 0 {
		idx = hs.Size()

		hs.atomicAdd(&heap.Heap{
			Title: "log",
			Path:  sys.Log.Name,
			Base:  sys.Log.Name,
			Type:  types.Stderr,
		})
	}

	atomic.StoreInt32(hs.index, idx)

	hs.atomicGet(idx).Reload()
}

func (hs *HeapSet) OpenFile(path, title string, tp types.Heap) {
	if !sys.Exists(path) {
		return
	}

	idx := hs.findByPath(path)

	if idx < 0 {
		idx = hs.Size()

		hs.atomicAdd(&heap.Heap{
			Title: title,
			Path:  path,
			Base:  path,
			Type:  tp,
		})
	}

	atomic.StoreInt32(hs.index, idx)

	hs.load()
}

func (hs *HeapSet) PrevHeap() *heap.Heap {
	idx := atomic.AddInt32(hs.index, -1)

	if idx < 0 {
		atomic.StoreInt32(hs.index, hs.Size()-1)
	}

	return hs.load()
}

func (hs *HeapSet) NextHeap() *heap.Heap {
	idx := atomic.AddInt32(hs.index, 1)

	if idx >= hs.Size() {
		atomic.StoreInt32(hs.index, 0)
	}

	return hs.load()
}

func (hs *HeapSet) CloseHeap() *heap.Heap {
	if hs.Size() == 1 {
		return nil
	}

	idx := atomic.LoadInt32(hs.index)

	h := hs.atomicGet(idx)

	hs.atomicDel(idx)
	hs.unload(h)

	atomic.AddInt32(hs.index, -1)

	return hs.NextHeap()
}

func (hs *HeapSet) ThrowAway() {
	hs.watch.Close()

	hs.Lock()

	for _, h := range hs.heaps {
		hs.unload(h)
	}

	hs.heaps = hs.heaps[:0]

	hs.Unlock()

	atomic.AddInt32(hs.index, -1)
}

func (hs *HeapSet) findByPath(path string) int32 {
	hs.RLock()
	defer hs.RUnlock()

	for i, h := range hs.heaps {
		if h.Base == path {
			return int32(i)
		}
	}

	return -1
}

func (hs *HeapSet) findByName(name string) int32 {
	hs.RLock()
	defer hs.RUnlock()

	for i, h := range hs.heaps {
		if h.Title == name {
			return int32(i)
		}
	}

	return -1
}

func (hs *HeapSet) atomicAdd(h *heap.Heap) {
	hs.Lock()
	hs.heaps = append(hs.heaps, h)
	hs.Unlock()
}

func (hs *HeapSet) atomicGet(idx int32) *heap.Heap {
	hs.RLock()
	defer hs.RUnlock()
	return hs.heaps[idx]
}

func (hs *HeapSet) atomicDel(idx int32) {
	hs.Lock()
	hs.heaps = slices.Delete(hs.heaps, int(idx), int(idx)+1)
	hs.Unlock()
}
