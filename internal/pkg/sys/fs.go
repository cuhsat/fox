package sys

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"

	foxfs "github.com/cuhsat/fox/internal/pkg/sys/fs"
)

var Watcher, _ = fsnotify.NewBufferedWatcher(2048)

var ffs = foxfs.NewForensicFs(
	afero.NewReadOnlyFs(
		afero.NewOsFs(),
	),
	foxfs.NewNotifyFs(
		afero.NewMemMapFs(),
		Watcher,
	),
)

type File = afero.File

func OpenThrough(path string) File {
	f, err := ffs.Open(path)

	if err != nil {
		Error(err)
		return nil
	}

	return f
}

func Exists(path string) bool {
	return ffs.Exists(path)
}

func CreateMem(path string) File {
	dir := filepath.Dir(path)

	err := ffs.MkdirAll(dir, 0x600)

	if err != nil {
		Error(err)
		return nil
	}

	f, err := ffs.Create(path)

	if err != nil {
		Error(err)
		return nil
	}

	return f
}
