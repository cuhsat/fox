package sys

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"

	foxfs "github.com/cuhsat/fox/internal/pkg/sys/fs"
)

var Watcher, _ = fsnotify.NewBufferedWatcher(1024)

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

func Open(path string) File {
	f, err := ffs.Open(path)

	if err != nil {
		Error(err)
	}

	return f
}

func Create(path string) File {
	f, err := ffs.Create(path)

	if err != nil {
		Error(err)
	}

	return f
}

func Exists(path string) bool {
	return ffs.Exists(path)
}
