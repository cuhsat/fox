package bag

import (
	"fmt"
	"os"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

type RawWrite struct {
	file *os.File // file handle
}

func NewRawWriter() *RawWrite {
	return new(RawWrite)
}

func (w *RawWrite) Init(f *os.File, _ bool, _ string) {
	w.file = f
}

func (w *RawWrite) Start() {}

func (w *RawWrite) Flush() {}

func (w *RawWrite) WriteMeta(_ meta) {}

func (w *RawWrite) WriteLine(_, _ int, s string) {
	_, err := fmt.Fprintln(w.file, s)

	if err != nil {
		sys.Error(err)
	}
}
