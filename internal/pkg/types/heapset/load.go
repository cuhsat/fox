package heapset

import (
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/cuhsat/fox/internal/pkg/file"
	"github.com/cuhsat/fox/internal/pkg/file/archive/tar"
	"github.com/cuhsat/fox/internal/pkg/file/archive/zip"
	"github.com/cuhsat/fox/internal/pkg/file/compress/bzip2"
	"github.com/cuhsat/fox/internal/pkg/file/compress/gzip"
	"github.com/cuhsat/fox/internal/pkg/file/compress/xz"
	"github.com/cuhsat/fox/internal/pkg/file/compress/zstd"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
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

	switch {
	case bzip2.Detect(path):
		path = bzip2.Deflate(path)
	case gzip.Detect(path):
		path = gzip.Deflate(path)
	case xz.Detect(path):
		path = xz.Deflate(path)
	case zstd.Detect(path):
		path = zstd.Deflate(path)
	case tar.Detect(path):
		hs.loadTar(path, base)
		return
	case zip.Detect(path):
		hs.loadZip(path, base)
		return
	}

	for _, p := range hs.plugins {
		if p.Match(path) {
			path, title := p.Execute(path, base, hs.Files())
			hs.loadAuto(path, base, title)
			return
		}
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

func (hs *HeapSet) loadAuto(path, base, title string) {
	hs.atomicAdd(&heap.Heap{
		Title: title,
		Path:  path,
		Base:  base,
		Type:  types.Plugin,
	})
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
		_ = os.Remove(h.Path)
	}
}
