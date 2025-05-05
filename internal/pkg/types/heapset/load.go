package heapset

import (
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/cuhsat/fx/internal/pkg/file"
	"github.com/cuhsat/fx/internal/pkg/file/bzip2"
	"github.com/cuhsat/fx/internal/pkg/file/gzip"
	"github.com/cuhsat/fx/internal/pkg/file/tar"
	"github.com/cuhsat/fx/internal/pkg/file/zip"
	"github.com/cuhsat/fx/internal/pkg/sys"
	"github.com/cuhsat/fx/internal/pkg/types"
	"github.com/cuhsat/fx/internal/pkg/types/heap"
)

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

	hs.watchHeap(h.Ensure().Filter())

	return h
}

func (hs *HeapSet) unload(h *heap.Heap) {
	h.ThrowAway()

	// clean up temporary files
	if h.Type == types.Stdin || h.Type == types.Deflate {
		os.Remove(h.Path)
	}
}
