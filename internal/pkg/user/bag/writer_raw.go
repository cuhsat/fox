package bag

import (
	"os"
	"os/user"
	"time"
)

type RawWrite struct {
	file *os.File // file handle
}

func NewRawWriter() *RawWrite {
	return &RawWrite{
		file: nil,
	}
}

func (w *RawWrite) Init(f *os.File, _ bool, _ string) {
	w.file = f
}

func (w *RawWrite) Start() {
	//
}

func (w *RawWrite) Finalize() {
	//
}

func (w *RawWrite) WriteFile(_ string, _ []string) {
	//
}

func (w *RawWrite) WriteUser(_ *user.User) {
	//
}

func (w *RawWrite) WriteTime(_, _ time.Time) {
	//
}

func (w *RawWrite) WriteHash(_ []byte) {
	//
}

func (w *RawWrite) WriteLines(_ []int, ss []string) {
	for i := 0; i < len(ss); i++ {
		writeln(w.file, ss[i])
	}
}
