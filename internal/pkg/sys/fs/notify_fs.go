package fs

import (
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
)

type NotifyFs struct {
	base afero.Fs

	watcher *fsnotify.Watcher
}

func NewNotifyFs(base afero.Fs, watcher *fsnotify.Watcher) afero.Fs {
	return &NotifyFs{base: base, watcher: watcher}
}

func (fs *NotifyFs) Chmod(name string, mode os.FileMode) error {
	err := fs.base.Chmod(name, mode)

	if err == nil {
		fs.watcher.Events <- fsnotify.Event{
			Name: name,
			Op:   fsnotify.Chmod,
		}
	}

	return err
}

func (fs *NotifyFs) Chown(name string, uid, gid int) error {
	err := fs.base.Chown(name, uid, gid)

	if err == nil {
		fs.watcher.Events <- fsnotify.Event{
			Name: name,
			Op:   fsnotify.Chmod,
		}
	}

	return err
}

func (fs *NotifyFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	err := fs.base.Chtimes(name, atime, mtime)

	if err == nil {
		fs.watcher.Events <- fsnotify.Event{
			Name: name,
			Op:   fsnotify.Chmod,
		}
	}

	return err
}

func (fs *NotifyFs) Create(name string) (afero.File, error) {
	f, err := fs.base.Create(name)

	if err != nil {
		return nil, err
	} else {
		fs.watcher.Events <- fsnotify.Event{
			Name: name,
			Op:   fsnotify.Create,
		}
	}

	return &NotifyFile{f, fs.watcher}, nil
}

func (fs *NotifyFs) Mkdir(name string, perm os.FileMode) error {
	err := fs.base.Mkdir(name, perm)

	if err == nil {
		fs.watcher.Events <- fsnotify.Event{
			Name: name,
			Op:   fsnotify.Create,
		}
	}

	return err
}

func (fs *NotifyFs) MkdirAll(path string, perm os.FileMode) error {
	err := fs.base.MkdirAll(path, perm)

	if err == nil {
		fs.watcher.Events <- fsnotify.Event{
			Name: path,
			Op:   fsnotify.Remove,
		}
	}

	return err
}

func (fs *NotifyFs) Name() string {
	return "NotifyFs"
}

func (fs *NotifyFs) Open(name string) (afero.File, error) {
	f, err := fs.base.Open(name)

	if err != nil {
		return nil, err
	}

	return &NotifyFile{f, fs.watcher}, nil
}

func (fs *NotifyFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.base.OpenFile(name, flag, perm)

	if err != nil {
		return nil, err
	}

	return &NotifyFile{f, fs.watcher}, nil
}

func (fs *NotifyFs) Remove(name string) error {
	err := fs.base.Remove(name)

	if err == nil {
		fs.watcher.Events <- fsnotify.Event{
			Name: name,
			Op:   fsnotify.Remove,
		}
	}

	return err
}

func (fs *NotifyFs) RemoveAll(path string) error {
	err := fs.base.Remove(path)

	if err == nil {
		fs.watcher.Events <- fsnotify.Event{
			Name: path,
			Op:   fsnotify.Remove,
		}
	}

	return err
}

func (fs *NotifyFs) Rename(oldname, newname string) error {
	err := fs.base.Rename(oldname, newname)

	if err == nil {
		fs.watcher.Events <- fsnotify.Event{
			Name: oldname,
			Op:   fsnotify.Rename,
		}
	}

	return err
}

func (fs *NotifyFs) Stat(name string) (os.FileInfo, error) {
	return fs.base.Stat(name)
}
