package bag

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/cuhsat/fox/internal/pkg/sys"
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

func (w *RawWrite) Start() {}

func (w *RawWrite) Flush() {}

func (w *RawWrite) SetFile(_ string, _ int64, _ []string) {}

func (w *RawWrite) SetUser(_ *user.User) {}

func (w *RawWrite) SetTime(_, _ time.Time) {}

func (w *RawWrite) SetHash(_ []byte) {}

func (w *RawWrite) SetLine(_, _ int, s string) {
	_, err := fmt.Fprintln(w.file, s)

	if err != nil {
		sys.Error(err)
	}
}
