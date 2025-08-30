package sys

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"

	"github.com/cuhsat/fox/internal/pkg/sys/fs"
)

var Watcher, _ = fsnotify.NewBufferedWatcher(1024)

var mem = fs.NewNotifyFs(afero.NewMemMapFs(), Watcher)

type File = afero.File

func Open(path string) File {
	mf, err := mem.Open(path)

	if err == nil {
		return mf // memory file
	}

	rf, err := os.OpenFile(path, os.O_RDONLY, 0400)

	if err == nil {
		return rf // regular file
	}

	Panic(err)
	return nil
}

func Create(path string) File {
	f, err := mem.Create(path)

	if err != nil {
		Panic(err)
	}

	return f
}

func Exists(path string) bool {
	_, err := mem.Stat(path)

	if err == nil {
		return true
	}

	_, err = os.Stat(path)

	if err == nil {
		return true
	}

	return false
}
