package fs

import (
	"io/fs"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
)

type NotifyFile struct {
	base afero.File

	watcher *fsnotify.Watcher
}

func (f *NotifyFile) Close() error {
	return f.base.Close()
}

func (f *NotifyFile) Name() string {
	return f.base.Name()
}

func (f *NotifyFile) Read(p []byte) (n int, err error) {
	return f.base.Read(p)
}

func (f *NotifyFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.base.ReadAt(p, off)
}

func (f *NotifyFile) Readdir(count int) ([]fs.FileInfo, error) {
	return f.base.Readdir(count)
}

func (f *NotifyFile) Readdirnames(n int) ([]string, error) {
	return f.base.Readdirnames(n)
}

func (f *NotifyFile) Seek(offset int64, whence int) (int64, error) {
	return f.base.Seek(offset, whence)
}

func (f *NotifyFile) Stat() (fs.FileInfo, error) {
	return f.base.Stat()
}

func (f *NotifyFile) Sync() error {
	return f.base.Sync()
}

func (f *NotifyFile) Truncate(size int64) error {
	err := f.base.Truncate(size)

	if err == nil {
		f.watcher.Events <- fsnotify.Event{
			Name: f.base.Name(),
			Op:   fsnotify.Write,
		}
	}

	return err
}

func (f *NotifyFile) Write(p []byte) (n int, err error) {
	n, err = f.base.Write(p)

	if err == nil {
		f.watcher.Events <- fsnotify.Event{
			Name: f.base.Name(),
			Op:   fsnotify.Write,
		}
	}

	return
}

func (f *NotifyFile) WriteAt(p []byte, off int64) (n int, err error) {
	n, err = f.base.WriteAt(p, off)

	if err == nil {
		f.watcher.Events <- fsnotify.Event{
			Name: f.base.Name(),
			Op:   fsnotify.Write,
		}
	}

	return
}

func (f *NotifyFile) WriteString(s string) (ret int, err error) {
	ret, err = f.base.WriteString(s)

	if err == nil {
		f.watcher.Events <- fsnotify.Event{
			Name: f.base.Name(),
			Op:   fsnotify.Write,
		}
	}

	return
}
