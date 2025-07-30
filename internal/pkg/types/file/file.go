package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.ReaderFrom
	io.ReadSeeker
	io.Seeker
	io.Writer
	io.WriterAt
	io.WriterTo
	io.StringWriter

	Name() string
	Stat() (fs.FileInfo, error)
	Truncate(size int64) error
}

type FileData struct {
	sync.RWMutex

	name string
	buf  []byte
	mod  time.Time
	pos  atomic.Int64
	evt  chan fsnotify.Event
}

type FileInfo struct {
	fd *FileData
}

var (
	vfs = make(map[string]File)
)

var (
	ErrorInvalidOffset = errors.New("invalid offset")
)

func Open(name string) File {
	return vfs[name]
}

func New(name string) File {
	rand.Seed(time.Now().UnixNano())

	f := NewFileData(name)

	vfs[f.Name()] = f

	return f
}

func NewFileData(name string) *FileData {
	return &FileData{name: fmt.Sprintf("fox://%d/%s", rand.Uint64(), name)}
}

func NewFileInfo(fd *FileData) *FileInfo {
	return &FileInfo{fd: fd}
}

func (fd *FileData) Bytes() []byte {
	fd.RLock()
	defer fd.RUnlock()
	return fd.buf
}

func (fd *FileData) Watch(ch chan fsnotify.Event) {
	fd.Lock()
	fd.evt = ch
	fd.Unlock()
}

func (fd *FileData) Close() error {
	fd.pos.Store(0)
	return nil
}

func (fd *FileData) Name() string {
	return fd.name
}

func (fd *FileData) Read(b []byte) (n int, err error) {
	n, err = fd.ReadAt(b, fd.pos.Load())

	fd.pos.Add(int64(n))

	return
}

func (fd *FileData) ReadAt(b []byte, off int64) (n int, err error) {
	fd.RLock()
	defer fd.RUnlock()

	if off < 0 || int64(int(off)) < off {
		return 0, ErrorInvalidOffset
	}

	if off > int64(len(fd.buf)) {
		return 0, io.EOF
	}

	n = copy(b, fd.buf[off:])

	if n < len(b) {
		return n, io.EOF
	}

	return n, nil
}

func (fd *FileData) ReadFrom(r io.Reader) (n int64, err error) {
	b, err := io.ReadAll(r)

	if err != nil {
		return 0, err
	}

	i, err := fd.Write(b)

	return int64(i), err
}

func (fd *FileData) Seek(offset int64, whence int) (int64, error) {
	fd.RLock()
	defer fd.RUnlock()

	var abs int64

	switch whence {
	case io.SeekStart:
		abs = offset

	case io.SeekCurrent:
		abs = fd.pos.Load() + offset

	case io.SeekEnd:
		abs = int64(len(fd.buf)) + offset

	default:
		return 0, ErrorInvalidOffset
	}

	if abs < 0 {
		return 0, ErrorInvalidOffset
	}

	fd.pos.Store(abs)

	return abs, nil
}

func (fd *FileData) Stat() (fs.FileInfo, error) {
	return NewFileInfo(fd), nil
}

func (fd *FileData) Truncate(size int64) error {
	fd.Lock()
	defer fd.Unlock()

	switch {
	case size < 0 || int64(int(size)) < size:
		return ErrorInvalidOffset

	case size <= int64(len(fd.buf)):
		fd.buf = fd.buf[:size]

	default:
		fd.buf = append(fd.buf, make([]byte, int(size)-len(fd.buf))...)
	}

	fd.mod = time.Now()

	fd.notify()

	return nil
}

func (fd *FileData) Write(b []byte) (n int, err error) {
	n, err = fd.WriteAt(b, fd.pos.Load())

	fd.pos.Add(int64(n))

	return
}

func (fd *FileData) WriteAt(b []byte, off int64) (n int, err error) {
	fd.Lock()
	defer fd.Unlock()

	if off < 0 || int64(int(off)) < off {
		return 0, ErrorInvalidOffset
	}

	if off > int64(len(fd.buf)) {
		_ = fd.Truncate(off)
	}

	n = copy(fd.buf[off:], b)

	fd.buf = append(fd.buf, b[n:]...)
	fd.mod = time.Now()

	fd.notify()

	return len(b), nil
}

func (fd *FileData) WriteTo(w io.Writer) (n int64, err error) {
	fd.RLock()
	defer fd.RUnlock()

	i, err := w.Write(fd.buf)

	return int64(i), err
}

func (fd *FileData) WriteString(s string) (n int, err error) {
	return fd.Write([]byte(s))
}

func (fd *FileData) notify() {
	if fd.evt != nil {
		fd.evt <- fsnotify.Event{
			Name: fd.name,
			Op:   fsnotify.Write,
		}
	}
}

func (fi *FileInfo) Name() string {
	return fi.fd.Name()
}

func (fi *FileInfo) Size() int64 {
	return int64(len(fi.fd.buf))
}

func (fi *FileInfo) Mode() fs.FileMode {
	return 0
}

func (fi *FileInfo) ModTime() time.Time {
	return fi.fd.mod
}

func (fi *FileInfo) IsDir() bool {
	return false
}

func (fi *FileInfo) Sys() any {
	return nil
}
