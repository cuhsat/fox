package heapset

import (
	"fmt"
	"os"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"

	"github.com/hiforensics/fox/internal/fox"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/file"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
	"github.com/hiforensics/fox/internal/pkg/user/plugins"
)

const (
	Stdin = "-"
)

type Call func()

type Each func(*heap.Heap)

type HeapSet struct {
	sync.RWMutex
	plugins []plugins.Plugin // automatic plugins

	watcher *fsnotify.Watcher // file watcher

	fnWatch Call // watcher callback
	fnError Call // error callback

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
		index:   new(int32),
	}

	if ps := plugins.New(); ps != nil {
		hs.plugins = ps.Autostarts()
	}

	go hs.notify()

	hs.watch(sys.Log.Name())

	if sys.Piped(os.Stdin) {
		paths = append(paths, Stdin)
	}

	for _, path := range paths {
		if path == Stdin {
			hs.loadPipe()
			break
		}

		hs.Open(path)
	}

	if hs.Len() == 0 {
		hs.OpenHelp()
	}

	hs.load()

	return &hs
}

func (hs *HeapSet) Len() int32 {
	hs.RLock()
	defer hs.RUnlock()
	return int32(len(hs.heaps))
}

func (hs *HeapSet) Each(fn Each) {
	hs.RLock()

	for _, h := range hs.heaps {
		fn(h.Ensure())
	}

	hs.RUnlock()
}

func (hs *HeapSet) Bind(fn1, fn2 Call) {
	hs.fnWatch = fn1
	hs.fnError = fn2
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

func (hs *HeapSet) OpenFox(path string) {
	idx, ok := hs.findByName("Examine")

	if !ok {
		idx = hs.Len()

		hs.atomicAdd(&heap.Heap{
			Title: "Examine",
			Path:  path,
			Base:  path,
			Type:  types.Prompt,
		})
	}

	atomic.StoreInt32(hs.index, idx)

	hs.load()
}

func (hs *HeapSet) OpenLog() {
	idx, ok := hs.findByPath(sys.Log.Name())

	if !ok {
		idx = hs.Len()

		hs.atomicAdd(&heap.Heap{
			Title: "Log",
			Path:  sys.Log.Name(),
			Base:  sys.Log.Name(),
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

		f := file.Create("Help", fmt.Sprintf(fox.Help, fox.Version))

		hs.atomicAdd(&heap.Heap{
			Title: "Help",
			Path:  f.Name(),
			Base:  f.Name(),
			Type:  types.Stdout,
		})
	}

	atomic.StoreInt32(hs.index, idx)

	hs.atomicGet(idx).Reload()
}

func (hs *HeapSet) OpenPlugin(path, base, title string) {
	idx, ok := hs.findByName(title)

	if !ok {
		idx = hs.Len()

		hs.atomicAdd(&heap.Heap{
			Title: title,
			Path:  path,
			Base:  base,
			Type:  types.Plugin,
		})
	} else {
		old := hs.atomicGet(idx)

		hs.atomicMod(&heap.Heap{
			Title: title,
			Path:  path,
			Base:  base,
			Type:  types.Plugin,
		}, idx)

		old.ThrowAway()
	}

	atomic.StoreInt32(hs.index, idx)

	hs.load()
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
		return nil // close program
	}

	idx := atomic.LoadInt32(hs.index)

	h := hs.atomicGet(idx)

	hs.atomicDel(idx)

	h.ThrowAway()

	atomic.AddInt32(hs.index, -1)

	return hs.NextHeap()
}

func (hs *HeapSet) Aggregate() bool {
	f := file.New("Aggregated")

	hs.RLock()

	var heaps []*heap.Heap

	for _, h := range hs.heaps {
		switch h.Type {
		case types.Deflate:
			fallthrough

		case types.Regular:
			_, _ = f.Write(h.Ensure().Bytes())
			_, _ = f.WriteString("\n")

			h.ThrowAway()

		default:
			heaps = append(heaps, h)
		}
	}

	hs.RUnlock()

	fi, _ := f.Stat()

	if fi.Size() == 0 {
		return false
	}

	hs.Lock()

	hs.heaps = append(heaps, &heap.Heap{
		Title: "Aggregated",
		Path:  f.Name(),
		Base:  f.Name(),
		Type:  types.Ignore,
	})

	hs.Unlock()

	idx := hs.Len() - 1

	atomic.StoreInt32(hs.index, idx)

	hs.atomicGet(idx).Reload()

	return true
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

func (hs *HeapSet) atomicMod(h *heap.Heap, idx int32) {
	hs.Lock()
	hs.heaps[idx] = h
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
