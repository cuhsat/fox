package fs

import (
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
)

var Watcher, _ = fsnotify.NewBufferedWatcher(2048)

var ffs = NewForensicFs(
	afero.NewReadOnlyFs(
		afero.NewOsFs(),
	),
	NewNotifyFs(
		afero.NewMemMapFs(),
		Watcher,
	),
)

func Open(path string) afero.File {
	f, _ := ffs.Open(path)
	return f
}

func Create(path string) afero.File {
	f, _ := ffs.Create(path)
	return f
}

func Exists(path string) bool {
	return ffs.Exists(path)
}

type ForensicFs struct {
	base  afero.Fs
	layer afero.Fs
}

func NewForensicFs(base, layer afero.Fs) *ForensicFs {
	return &ForensicFs{base: base, layer: layer}
}

func (fs *ForensicFs) Chmod(name string, mode os.FileMode) error {
	return fs.layer.Chmod(name, mode)
}

func (fs *ForensicFs) Chown(name string, uid, gid int) error {
	return fs.layer.Chown(name, uid, gid)
}

func (fs *ForensicFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.layer.Chtimes(name, atime, mtime)
}

func (fs *ForensicFs) Create(name string) (afero.File, error) {
	return fs.layer.Create(name)
}

func (fs *ForensicFs) Mkdir(name string, perm os.FileMode) error {
	return fs.layer.Mkdir(name, perm)
}

func (fs *ForensicFs) MkdirAll(path string, perm os.FileMode) error {
	return fs.layer.MkdirAll(path, perm)
}

func (fs *ForensicFs) Name() string {
	return "ForensicFs"
}

func (fs *ForensicFs) Open(name string) (afero.File, error) {
	lf, err := fs.layer.Open(name)

	if err == nil {
		return lf, nil // layer file
	}

	bf, err := fs.base.Open(name)

	if err == nil {
		return bf, nil // base file
	}

	return nil, err
}

func (fs *ForensicFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	lf, err := fs.layer.OpenFile(name, flag, perm)

	if err == nil {
		return lf, nil // layer file
	}

	bf, err := fs.base.OpenFile(name, flag, perm)

	if err == nil {
		return bf, nil // base file
	}

	return nil, err
}

func (fs *ForensicFs) Remove(name string) error {
	return fs.layer.Remove(name)
}

func (fs *ForensicFs) RemoveAll(path string) error {
	return fs.layer.Remove(path)
}

func (fs *ForensicFs) Rename(oldname, newname string) error {
	return fs.layer.Rename(oldname, newname)
}

func (fs *ForensicFs) Stat(name string) (os.FileInfo, error) {
	return fs.layer.Stat(name)
}

func (fs *ForensicFs) Exists(name string) bool {
	_, err := fs.layer.Stat(name)

	if err == nil {
		return true // layer stat
	}

	_, err = fs.base.Stat(name)

	if err == nil {
		return true // base stat
	}

	return false
}
