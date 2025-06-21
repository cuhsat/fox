package heapset

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/cuhsat/fox/internal/pkg/file/archive"
	"github.com/cuhsat/fox/internal/pkg/file/compress/br"
	"github.com/cuhsat/fox/internal/pkg/file/compress/bzip2"
	"github.com/cuhsat/fox/internal/pkg/file/compress/gzip"
	"github.com/cuhsat/fox/internal/pkg/file/compress/lz4"
	"github.com/cuhsat/fox/internal/pkg/file/compress/xz"
	"github.com/cuhsat/fox/internal/pkg/file/compress/zlib"
	"github.com/cuhsat/fox/internal/pkg/file/compress/zstd"
	"github.com/cuhsat/fox/internal/pkg/file/format/csv"
	"github.com/cuhsat/fox/internal/pkg/file/parser/evtx"
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

	// check for compression
	switch {
	case br.Detect(path):
		path = br.Deflate(path)
	case bzip2.Detect(path):
		path = bzip2.Deflate(path)
	case gzip.Detect(path):
		path = gzip.Deflate(path)
	case lz4.Detect(path):
		path = lz4.Deflate(path)
	case xz.Detect(path):
		path = xz.Deflate(path)
	case zlib.Detect(path):
		path = zlib.Deflate(path)
	case zstd.Detect(path):
		path = zstd.Deflate(path)
	}

	// check for archive
	if archive.Detect(path) {
		hs.loadArchive(path, base)
		return
	}

	// check for parser
	if evtx.Detect(path) {
		path = evtx.Parse(path)
	}

	// check for format
	if csv.Detect(path) {
		path = csv.Format(path)
	}

	// check for plugin
	for _, p := range hs.plugins {
		if p.Match(path) {
			p.Execute(path, base, func(path, base string) {
				hs.loadPlugin(path, base, p.Name)
			})
			return
		}
	}

	hs.loadFile(path, base)
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

func (hs *HeapSet) loadFile(path, base string) {
	h := heap.New(
		base,
		path,
		base,
		types.Regular,
	)

	if path != base {
		h.Type = types.Deflate
	}

	hs.atomicAdd(h)
}

func (hs *HeapSet) loadArchive(path, base string) {
	for _, i := range archive.Deflate(path) {
		if evtx.Detect(i.Path) {
			i.Path = evtx.Parse(i.Path)
		}

		hs.atomicAdd(heap.New(
			filepath.Join(base, i.Name),
			i.Path,
			i.Path,
			types.Deflate,
		))
	}
}

func (hs *HeapSet) loadPlugin(path, base, name string) {
	hs.atomicAdd(heap.New(
		fmt.Sprintf("%s (%s)", base, name),
		path,
		base,
		types.Plugin,
	))
}

func (hs *HeapSet) loadPipe() {
	pipe := sys.Stdin()

	hs.atomicAdd(heap.New(
		"",
		pipe,
		pipe,
		types.Stdin,
	))
}

func (hs *HeapSet) load() *heap.Heap {
	h := hs.atomicGet(atomic.LoadInt32(hs.index))

	hs.watchHeap(h.Ensure())

	return h
}

func (hs *HeapSet) unload(h *heap.Heap) {
	h.ThrowAway()

	// clean up temporary files
	if h.Type != types.Regular && h.Type != types.Ignore {
		_ = os.Remove(h.Path)
	}
}
