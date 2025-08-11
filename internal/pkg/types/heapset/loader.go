package heapset

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/hiforensics/fox/internal/pkg/files"
	"github.com/hiforensics/fox/internal/pkg/files/archive/rar"
	"github.com/hiforensics/fox/internal/pkg/files/archive/tar"
	"github.com/hiforensics/fox/internal/pkg/files/archive/zip"
	"github.com/hiforensics/fox/internal/pkg/files/compress/br"
	"github.com/hiforensics/fox/internal/pkg/files/compress/bzip2"
	"github.com/hiforensics/fox/internal/pkg/files/compress/gzip"
	"github.com/hiforensics/fox/internal/pkg/files/compress/lz4"
	"github.com/hiforensics/fox/internal/pkg/files/compress/xz"
	"github.com/hiforensics/fox/internal/pkg/files/compress/zlib"
	"github.com/hiforensics/fox/internal/pkg/files/compress/zstd"
	"github.com/hiforensics/fox/internal/pkg/files/format/csv"
	"github.com/hiforensics/fox/internal/pkg/files/parser/evtx"
	"github.com/hiforensics/fox/internal/pkg/flags"
	"github.com/hiforensics/fox/internal/pkg/sys"
	"github.com/hiforensics/fox/internal/pkg/types"
	"github.com/hiforensics/fox/internal/pkg/types/heap"
)

func (hs *HeapSet) loadPath(path string) {
	base := path

	fi, err := os.Stat(path)

	if err != nil {
		sys.Error(err)
		return
	}

	if fi.IsDir() {
		hs.loadDir(path)
		return
	}

	if !flags.Get().Opt.NoDeflate {
		path = hs.deflate(path, base)

		if len(path) == 0 {
			return
		}
	}

	path = hs.process(path, base)

	if len(path) == 0 {
		return
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

func (hs *HeapSet) loadArchive(fn files.Deflate, path, base, pass string) {
	for _, i := range fn(path, pass) {
		i.Path = hs.deflate(i.Path, base)

		if len(i.Path) == 0 {
			return
		}

		i.Path = hs.process(i.Path, base)

		if len(i.Path) == 0 {
			return
		}

		hs.atomicAdd(heap.New(
			filepath.Join(base, i.Name),
			i.Path,
			base,
			types.Deflate,
		))
	}
}

func (hs *HeapSet) loadPlugin(path, base, name string) {
	hs.atomicAdd(heap.New(
		fmt.Sprintf("%s : %s", base, name),
		path,
		base,
		types.Plugin,
	))
}

func (hs *HeapSet) loadPipe() {
	pipe := sys.Stdin().Name()

	hs.atomicAdd(heap.New(
		"",
		pipe,
		pipe,
		types.Stdin,
	))
}

func (hs *HeapSet) load() *heap.Heap {
	h := hs.atomicGet(atomic.LoadInt32(hs.index))

	hs.watch(h.Ensure().Path)

	return h
}

func (hs *HeapSet) deflate(path, base string) string {
	pass := flags.Get().Deflate.Pass

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
	switch {
	case rar.Detect(path):
		hs.loadArchive(rar.Deflate, path, base, pass)
		return ""
	case tar.Detect(path):
		hs.loadArchive(tar.Deflate, path, base, pass)
		return ""
	case zip.Detect(path):
		hs.loadArchive(zip.Deflate, path, base, pass)
		return ""
	}

	return path
}

func (hs *HeapSet) process(path, base string) string {
	if !flags.Get().Opt.NoConvert {
		// check for parser
		if evtx.Detect(path) {
			path = evtx.Parse(path)
		}

		// check for format
		if csv.Detect(path) {
			path = csv.Format(path)
		}
	}

	if !flags.Get().Opt.NoPlugins {
		// check for plugin
		for _, p := range hs.plugins {
			if p.Match(path) {
				p.Execute(path, base, func(path, base, dir string) {
					if len(dir) > 0 {
						hs.Open(dir)
					}

					hs.loadPlugin(path, base, p.Name)
				})

				return ""
			}
		}
	}

	return path
}
