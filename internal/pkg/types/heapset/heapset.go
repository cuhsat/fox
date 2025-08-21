package heapset

import (
	"fmt"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/types/loader"
)

type Callback func()

type Each func(int, *heap.Heap)

type HeapSet struct {
	sync.RWMutex

	loader  *loader.Loader    // file loader
	watcher *fsnotify.Watcher // file watcher

	error Callback // error callback
	watch Callback // watch callback

	heaps []*heap.Heap // set heaps
	index *int32       // set index
}

func New(paths []string) *HeapSet {
	w, err := fsnotify.NewWatcher()

	if err != nil {
		sys.Error(err)
	}

	hs := HeapSet{
		watcher: w,
		loader:  loader.New(),
		index:   new(int32),
	}

	go hs.watchFiles()

	hs.watchFile(sys.Log.Name())

	for _, e := range hs.loader.Init(paths) {
		hs.atomicAdd(heap.New(e.Name, e.Path, e.Base, e.Type))
	}

	if hs.Len() == 0 {
		hs.OpenHelp()
	}

	// load first heap
	hs.LoadHeap()

	return &hs
}

func (hs *HeapSet) Len() int32 {
	hs.RLock()
	defer hs.RUnlock()
	return int32(len(hs.heaps))
}

func (hs *HeapSet) Each(fn Each) {
	hs.RLock()

	for i, h := range hs.heaps {
		fn(i, h.Ensure())
	}

	hs.RUnlock()
}

func (hs *HeapSet) Heap() (int32, *heap.Heap) {
	idx := atomic.LoadInt32(hs.index)
	return idx + 1, hs.atomicGet(idx)
}

func (hs *HeapSet) Open(path string) {
	for _, e := range hs.loader.Load(path) {
		hs.atomicAdd(heap.New(e.Name, e.Path, e.Base, e.Type))
	}
}

func (hs *HeapSet) OpenLog() {
	idx, ok := hs.findByPath(sys.Log.Name())

	if !ok {
		idx = hs.Len()

		hs.atomicAdd(heap.New(
			"Log",
			sys.Log.Name(),
			sys.Log.Name(),
			types.Stderr,
		))
	}

	atomic.StoreInt32(hs.index, idx)

	hs.atomicGet(idx).Reload()
}

func (hs *HeapSet) OpenHelp() {
	idx, ok := hs.findByName("Help")

	if !ok {
		idx = hs.Len()

		f := sys.Create("Help")
		_, _ = f.WriteString(fmt.Sprintf(app.Ascii+app.Help, app.Version))

		hs.atomicAdd(heap.New(
			"Keymap",
			f.Name(),
			f.Name(),
			types.Stdout,
		))
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

		hs.atomicAdd(heap.New(title, path, base, tp))
	}

	atomic.StoreInt32(hs.index, idx)

	hs.LoadHeap()
}

func (hs *HeapSet) OpenPlugin(path, base, title string) {
	idx, ok := hs.findByName(title)

	if !ok {
		idx = hs.Len()

		hs.atomicAdd(heap.New(title, path, base, types.Plugin))
	} else {
		old := hs.atomicGet(idx)

		hs.atomicAdd(heap.New(title, path, base, types.Plugin))

		old.ThrowAway()
	}

	atomic.StoreInt32(hs.index, idx)

	hs.LoadHeap()
}

func (hs *HeapSet) OpenAgent(path string) {
	idx, ok := hs.findByName("Agent")

	if !ok {
		idx = hs.Len()

		hs.atomicAdd(heap.New("Agent", path, path, types.Agent))
	}

	atomic.StoreInt32(hs.index, idx)

	hs.LoadHeap()
}

func (hs *HeapSet) PrevHeap() *heap.Heap {
	idx := atomic.AddInt32(hs.index, -1)

	if idx < 0 {
		atomic.StoreInt32(hs.index, hs.Len()-1)
	}

	return hs.LoadHeap()
}

func (hs *HeapSet) NextHeap() *heap.Heap {
	idx := atomic.AddInt32(hs.index, 1)

	if idx >= hs.Len() {
		atomic.StoreInt32(hs.index, 0)
	}

	return hs.LoadHeap()
}

func (hs *HeapSet) LoadHeap() *heap.Heap {
	h := hs.atomicGet(atomic.LoadInt32(hs.index))

	hs.watchFile(h.Ensure().Path)

	return h
}

func (hs *HeapSet) CloseHeap() *heap.Heap {
	if hs.Len() == 1 {
		return nil // close program
	}

	idx := atomic.LoadInt32(hs.index)

	h := hs.atomicGet(idx)

	hs.atomicDel(idx)

	h.ThrowAway()

	atomic.AddInt32(hs.index, -1)

	return hs.NextHeap()
}

func (hs *HeapSet) ThrowAway() {
	_ = hs.watcher.Close()

	hs.Lock()

	for _, h := range hs.heaps {
		h.ThrowAway()
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
