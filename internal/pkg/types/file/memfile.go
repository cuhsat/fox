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
	io.ReadCloser
	io.ReadSeeker
	io.ReadSeekCloser
	io.ReadWriteCloser
	io.ReadWriteSeeker
	io.ReadWriter
	io.Seeker
	io.StringWriter
	io.Writer
	io.WriterAt
	io.WriterTo
	io.WriteCloser
	io.WriteSeeker

	Name() string
	Stat() (fs.FileInfo, error)
	Truncate(size int64) error
}

type Data struct {
	sync.RWMutex

	name string
	buf  []byte
	mod  time.Time
	pos  atomic.Int64
	evt  chan fsnotify.Event
}

type Info struct {
	data *Data
}

var (
	vfs = make(map[string]File)
)

var (
	ErrorInvalidOffset = errors.New("invalid offset")
)

func New(name string) File {
	f := &Data{
		name: fmt.Sprintf("fox://%d/%s", rand.Uint64(), name),
	}

	vfs[f.Name()] = f

	return f
}

func Open(name string) File {
	return vfs[name]
}

func Create(name, data string) File {
	f := New(name)

	_, _ = f.WriteString(data)

	return f
}

func (d *Data) Bytes() []byte {
	d.RLock()
	defer d.RUnlock()
	return d.buf
}

func (d *Data) Watch(ch chan fsnotify.Event) {
	d.Lock()
	d.evt = ch
	d.Unlock()
}

func (d *Data) Close() error {
	d.pos.Store(0)
	return nil
}

func (d *Data) Name() string {
	return d.name
}

func (d *Data) Read(b []byte) (n int, err error) {
	n, err = d.ReadAt(b, d.pos.Load())

	d.pos.Add(int64(n))

	return
}

func (d *Data) ReadAt(b []byte, off int64) (n int, err error) {
	d.RLock()
	defer d.RUnlock()

	if off < 0 || int64(int(off)) < off {
		return 0, ErrorInvalidOffset
	}

	if off > int64(len(d.buf)) {
		return 0, io.EOF
	}

	n = copy(b, d.buf[off:])

	if n < len(b) {
		return n, io.EOF
	}

	return n, nil
}

func (d *Data) ReadFrom(r io.Reader) (n int64, err error) {
	b, err := io.ReadAll(r)

	if err != nil {
		return 0, err
	}

	i, err := d.Write(b)

	return int64(i), err
}

func (d *Data) Seek(offset int64, whence int) (int64, error) {
	d.RLock()
	defer d.RUnlock()

	var abs int64

	switch whence {
	case io.SeekStart:
		abs = offset

	case io.SeekCurrent:
		abs = d.pos.Load() + offset

	case io.SeekEnd:
		abs = int64(len(d.buf)) + offset

	default:
		return 0, ErrorInvalidOffset
	}

	if abs < 0 {
		return 0, ErrorInvalidOffset
	}

	d.pos.Store(abs)

	return abs, nil
}

func (d *Data) Stat() (fs.FileInfo, error) {
	return &Info{data: d}, nil
}

func (d *Data) Truncate(size int64) error {
	d.Lock()
	defer d.Unlock()

	switch {
	case size < 0 || int64(int(size)) < size:
		return ErrorInvalidOffset

	case size <= int64(len(d.buf)):
		d.buf = d.buf[:size]

	default:
		d.buf = append(d.buf, make([]byte, int(size)-len(d.buf))...)
	}

	d.mod = time.Now()

	d.notify()

	return nil
}

func (d *Data) Write(b []byte) (n int, err error) {
	n, err = d.WriteAt(b, d.pos.Load())

	d.pos.Add(int64(n))

	return
}

func (d *Data) WriteAt(b []byte, off int64) (n int, err error) {
	d.Lock()
	defer d.Unlock()

	if off < 0 || int64(int(off)) < off {
		return 0, ErrorInvalidOffset
	}

	if off > int64(len(d.buf)) {
		_ = d.Truncate(off)
	}

	n = copy(d.buf[off:], b)

	d.buf = append(d.buf, b[n:]...)
	d.mod = time.Now()

	d.notify()

	return len(b), nil
}

func (d *Data) WriteTo(w io.Writer) (n int64, err error) {
	d.RLock()
	defer d.RUnlock()

	i, err := w.Write(d.buf)

	return int64(i), err
}

func (d *Data) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

func (d *Data) notify() {
	if d.evt != nil {
		d.evt <- fsnotify.Event{
			Name: d.name,
			Op:   fsnotify.Write,
		}
	}
}

func (i *Info) Name() string {
	return i.data.Name()
}

func (i *Info) Size() int64 {
	i.data.RLock()
	defer i.data.RUnlock()
	return int64(len(i.data.buf))
}

func (i *Info) Mode() fs.FileMode {
	return 0 // regular
}

func (i *Info) ModTime() time.Time {
	i.data.RLock()
	defer i.data.RUnlock()
	return i.data.mod
}

func (i *Info) IsDir() bool {
	return false
}

func (i *Info) Sys() any {
	return nil
}
