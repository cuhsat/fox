package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/cuhsat/fox/internal/pkg/files"
	"github.com/cuhsat/fox/internal/pkg/files/archive/cab"
	"github.com/cuhsat/fox/internal/pkg/files/archive/rar"
	"github.com/cuhsat/fox/internal/pkg/files/archive/tar"
	"github.com/cuhsat/fox/internal/pkg/files/archive/zip"
	"github.com/cuhsat/fox/internal/pkg/files/compress/br"
	"github.com/cuhsat/fox/internal/pkg/files/compress/bzip2"
	"github.com/cuhsat/fox/internal/pkg/files/compress/gzip"
	"github.com/cuhsat/fox/internal/pkg/files/compress/lz4"
	"github.com/cuhsat/fox/internal/pkg/files/compress/xz"
	"github.com/cuhsat/fox/internal/pkg/files/compress/zlib"
	"github.com/cuhsat/fox/internal/pkg/files/compress/zstd"
	"github.com/cuhsat/fox/internal/pkg/files/format/csv"
	"github.com/cuhsat/fox/internal/pkg/files/parser/evtx"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/user/plugins"
)

const (
	Stdin = "-"
)

type Loader struct {
	plugins []plugins.Plugin // automatic plugins
	entries []Entry          // temporary entries
}

type Entry struct {
	Name string     // entry name
	Path string     // entry file path
	Base string     // entry file base
	Type types.Heap // entry heap type
}

func New() *Loader {
	l := new(Loader)

	if ps := plugins.New(); ps != nil {
		l.plugins = ps.Autostarts()
	}

	return l
}

func (l *Loader) Init(paths []string) []Entry {
	if sys.Piped(os.Stdin) {
		paths = append(paths, Stdin)
	}

	for _, path := range paths {
		if path == Stdin {
			l.addPipe()
			break
		}

		l.loadPath(path)
	}

	return l.entries
}

func (l *Loader) Load(path string) []Entry {
	l.entries = l.entries[:0] // reset

	l.loadPath(path)

	return l.entries
}

func (l *Loader) loadPath(path string) {
	match, err := doublestar.FilepathGlob(path)

	if err != nil {
		sys.Error(err)
	}

	for _, m := range match {
		fi, err := os.Stat(m)

		if err != nil {
			sys.Error(err)
			return
		}

		if fi.IsDir() {
			l.loadDir(m)
		} else {
			l.loadFile(m)
		}
	}
}

func (l *Loader) loadDir(path string) {
	dir, err := os.ReadDir(path)

	if err != nil {
		sys.Error(err)
		return
	}

	for _, f := range dir {
		if !f.IsDir() {
			l.loadFile(filepath.Join(path, f.Name()))
		}
	}
}

func (l *Loader) loadFile(path string) {
	base := path

	if !flags.Get().Opt.NoDeflate {
		path = l.deflate(path, base)

		if len(path) == 0 {
			return
		}
	}

	path = l.process(path, base)

	if len(path) == 0 {
		return
	}

	l.addFile(path, base)
}

func (l *Loader) loadArchive(path, base string, fn files.Deflate) {
	defer func() {
		if err := recover(); err != nil {
			sys.Error("corrupted archive or wrong password")
			return
		}
	}()

	items := fn(path, flags.Get().Deflate.Pass)

	if len(items) == 0 {
		panic("no item(s)")
	}

	for _, i := range items {
		i.Path = l.deflate(i.Path, base)

		if len(i.Path) == 0 {
			continue
		}

		i.Path = l.process(i.Path, base)

		if len(i.Path) == 0 {
			continue
		}

		l.addItem(i.Path, base, i.Name)
	}
}

func (l *Loader) addPipe() {
	pipe := sys.Stdin().Name()

	l.entries = append(l.entries, Entry{
		"",
		pipe,
		pipe,
		types.Stdin,
	})
}

func (l *Loader) addFile(path, base string) {
	var t = types.Regular

	if path != base {
		t = types.Deflate
	}

	l.entries = append(l.entries, Entry{
		base,
		path,
		base,
		t,
	})
}

func (l *Loader) addItem(path, base, name string) {
	l.entries = append(l.entries, Entry{
		filepath.Join(base, name),
		path,
		base,
		types.Deflate,
	})
}

func (l *Loader) addPlugin(path, base, name string) {
	l.entries = append(l.entries, Entry{
		fmt.Sprintf("%s : %s", base, name),
		path,
		base,
		types.Plugin,
	})
}

func (l *Loader) deflate(path, base string) string {
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
	case cab.Detect(path):
		l.loadArchive(path, base, cab.Deflate)
		return ""
	case rar.Detect(path):
		l.loadArchive(path, base, rar.Deflate)
		return ""
	case tar.Detect(path):
		l.loadArchive(path, base, tar.Deflate)
		return ""
	case zip.Detect(path):
		l.loadArchive(path, base, zip.Deflate)
		return ""
	}

	return path
}

func (l *Loader) process(path, base string) string {
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
		for _, p := range l.plugins {
			if p.Match(path) {
				p.Execute(path, base, func(path, base, dir string) {
					if len(dir) > 0 {
						l.loadDir(dir) // load dir results
					}

					l.addPlugin(path, base, p.Name)
				})

				return ""
			}
		}
	}

	return path
}
