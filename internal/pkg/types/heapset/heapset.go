package heapset

import (
	"os"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fox/internal/fox"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/user/plugins"
)

const (
	Stdin = "-"
)

type callback func()

type HeapSet struct {
	sync.RWMutex
	plugins []plugins.Plugin // automatic plugins

	watch   *fsnotify.Watcher // file watcher
	watchFn callback          // file watcher callback
	errorFn callback          // error callback

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

	if ps := plugins.New(); ps != nil {
		hs.plugins = ps.Autostarts()
	}

	go hs.notify()

	hs.watchPath(sys.Log.Name)

	if sys.IsPiped(os.Stdin) {
		paths = append(paths, Stdin)
	}

	for _, path := range paths {
		if path == Stdin {
			hs.loadPipe()
			break
		}

		hs.Open(path)
	}

	if len(paths) == 0 {
		hs.OpenHelp()
	} else if hs.Len() == 0 {
		sys.Exit("no files found")
	}

	hs.load()

	return &hs
}

func (hs *HeapSet) Len() int32 {
	hs.RLock()
	defer hs.RUnlock()
	return int32(len(hs.heaps))
}

func (hs *HeapSet) Bind(fn1, fn2 callback) {
	hs.watchFn = fn1
	hs.errorFn = fn2
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
	idx, ok := hs.findByPath(sys.Log.Name)

	if !ok {
		idx = hs.Len()

		hs.atomicAdd(&heap.Heap{
			Title: "Log",
			Path:  sys.Log.Name,
			Base:  sys.Log.Name,
			Type:  types.Stderr,
		})
	}

	atomic.StoreInt32(hs.index, idx)

	hs.atomicGet(idx).Reload()
}

func (hs *HeapSet) OpenHelp() {
	idx, ok := hs.findByName("Help")

	if !ok {
		idx = hs.Len()

		p := sys.DumpStr(fox.Help)

		hs.atomicAdd(&heap.Heap{
			Title: "Help",
			Path:  p,
			Base:  p,
			Type:  types.Stdout,
		})
	}

	atomic.StoreInt32(hs.index, idx)

	hs.atomicGet(idx).Reload()
}

func (hs *HeapSet) OpenFile(path, base, title string, tp types.Heap) {
	if !sys.Exists(path) {
		return
	}

	idx, ok := hs.findByPath(path)

	if !ok {
		idx = hs.Len()

		hs.atomicAdd(&heap.Heap{
			Title: title,
			Path:  path,
			Base:  base,
			Type:  tp,
		})
	}

	atomic.StoreInt32(hs.index, idx)

	hs.load()
}

func (hs *HeapSet) PrevHeap() *heap.Heap {
	idx := atomic.AddInt32(hs.index, -1)

	if idx < 0 {
		atomic.StoreInt32(hs.index, hs.Len()-1)
	}

	return hs.load()
}

func (hs *HeapSet) NextHeap() *heap.Heap {
	idx := atomic.AddInt32(hs.index, 1)

	if idx >= hs.Len() {
		atomic.StoreInt32(hs.index, 0)
	}

	return hs.load()
}

func (hs *HeapSet) CloseHeap() *heap.Heap {
	if hs.Len() == 1 {
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
	_ = hs.watch.Close()

	hs.Lock()

	for _, h := range hs.heaps {
		hs.unload(h)
	}

	hs.heaps = hs.heaps[:0]

	hs.Unlock()

	atomic.AddInt32(hs.index, -1)
}

func (hs *HeapSet) findByPath(path string) (int32, bool) {
	hs.RLock()
	defer hs.RUnlock()

	for i, h := range hs.heaps {
		if h.Base == path {
			return int32(i), true
		}
	}

	return 0, false
}

func (hs *HeapSet) findByName(name string) (int32, bool) {
	hs.RLock()
	defer hs.RUnlock()

	for i, h := range hs.heaps {
		if h.Title == name {
			return int32(i), true
		}
	}

	return 0, false
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
