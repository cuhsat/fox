package heapset

import (
	"os"
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"

	"github.com/cuhsat/fx/pkg/fx/file"
	"github.com/cuhsat/fx/pkg/fx/file/bzip2"
	"github.com/cuhsat/fx/pkg/fx/file/gzip"
	"github.com/cuhsat/fx/pkg/fx/file/tar"
	"github.com/cuhsat/fx/pkg/fx/file/zip"
	"github.com/cuhsat/fx/pkg/fx/sys"
	"github.com/cuhsat/fx/pkg/fx/types"
	"github.com/cuhsat/fx/pkg/fx/types/heap"
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

	if hs.Length() == 0 {
		sys.Exit("no files found")
	}

	hs.load()

	return &hs
}

func (hs *HeapSet) Bind(fn1, fn2 callback) {
	hs.watch_fn = fn1
	hs.error_fn = fn2
}

func (hs *HeapSet) Length() int32 {
	hs.RLock()
	defer hs.RUnlock()

	return int32(len(hs.heaps))
}

func (hs *HeapSet) Current() (int32, *heap.Heap) {
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
		idx = hs.Length()

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

func (hs *HeapSet) OpenHeap(path string) {
	if !sys.Exists(path) {
		return
	}

	idx := hs.findByPath(path)

	if idx < 0 {
		idx = hs.Length()

		hs.loadPath(path)
	}

	atomic.StoreInt32(hs.index, idx)

	hs.load()
}

func (hs *HeapSet) PrevHeap() *heap.Heap {
	idx := atomic.AddInt32(hs.index, -1)

	if idx < 0 {
		atomic.StoreInt32(hs.index, hs.Length()-1)
	}

	return hs.load()
}

func (hs *HeapSet) NextHeap() *heap.Heap {
	idx := atomic.AddInt32(hs.index, 1)

	if idx >= hs.Length() {
		atomic.StoreInt32(hs.index, 0)
	}

	return hs.load()
}

func (hs *HeapSet) CloseHeap() *heap.Heap {
	if hs.Length() == 1 {
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

func (hs *HeapSet) loadPath(path string) {
	fi, err := os.Stat(path)

	if err != nil {
		sys.Error(err)
		return
	}

	if fi.IsDir() {
		hs.loadDir(path)
		return
	}

	base := path

	if bzip2.Detect(path) {
		path = bzip2.Deflate(path)
	}

	if gzip.Detect(path) {
		path = gzip.Deflate(path)
	}

	if tar.Detect(path) {
		hs.loadTar(path, base)
		return
	}

	if zip.Detect(path) {
		hs.loadZip(path, base)
		return
	}

	hs.loadFile(path, base)
}

func (hs *HeapSet) loadPipe() {
	pipe := sys.Stdin()

	hs.atomicAdd(&heap.Heap{
		Path: pipe,
		Base: pipe,
		Type: types.Stdin,
	})
}

func (hs *HeapSet) loadDir(path string) {
	dir, err := os.ReadDir(path)

	if err != nil {
		sys.Error(err)
		return
	}

	for _, f := range dir {
		if !f.IsDir() {
			hs.loadPath(filepath.Join(path, f.Name()))
		}
	}
}

func (hs *HeapSet) loadTar(path, base string) {
	for _, i := range tar.Deflate(path) {
		hs.loadItem(i, base)
	}
}

func (hs *HeapSet) loadZip(path, base string) {
	for _, i := range zip.Deflate(path) {
		hs.loadItem(i, base)
	}
}

func (hs *HeapSet) loadFile(path, base string) {
	h := &heap.Heap{
		Title: base,
		Path:  path,
		Base:  base,
		Type:  types.Regular,
	}

	if path != base {
		h.Type = types.Deflate
	}

	hs.atomicAdd(h)
}

func (hs *HeapSet) loadItem(i *file.Item, base string) {
	hs.atomicAdd(&heap.Heap{
		Title: filepath.Join(base, i.Name),
		Path:  i.Path,
		Base:  i.Path,
		Type:  types.Deflate,
	})
}

func (hs *HeapSet) load() *heap.Heap {
	h := hs.atomicGet(atomic.LoadInt32(hs.index))

	if !h.Loaded() {
		h.Reload()

		hs.watchHeap(h)
	} else {
		h.Filter() // TODO: Resets rmap to nil
	}

	return h
}

func (hs *HeapSet) unload(h *heap.Heap) {
	h.ThrowAway()

	// clean up temporary files
	if h.Type == types.Stdin || h.Type == types.Deflate {
		os.Remove(h.Path)
	}
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
