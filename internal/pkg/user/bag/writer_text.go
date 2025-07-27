package bag

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hiforensics/fox/internal/pkg/sys"
)

type TextWrite struct {
	file *os.File // file handle
}

func NewTextWriter() *TextWrite {
	return &TextWrite{
		file: nil,
	}
}

func (w *TextWrite) Init(f *os.File, old bool, title string) {
	w.file = f

	if old {
		return
	}

	_, err := fmt.Fprintf(w.file, "%s\n", title)

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) Start() {
	_, err := fmt.Fprintf(w.file, "\n")

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) Flush() {
	_, err := fmt.Fprintf(w.file, "\n")

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) WriteMeta(meta meta) {
	f := strings.Join(meta.filters, " > ")

	_, err := fmt.Fprintf(w.file, "File: %s %s (%d bytes)\n", meta.path, f, meta.size)

	if err != nil {
		sys.Error(err)
	}

	_, err = fmt.Fprintf(w.file, "User: %s (%s)\n", meta.user.Username, meta.user.Name)

	if err != nil {
		sys.Error(err)
	}

	_, err = fmt.Fprintf(w.file, "Time: %s / %s\n",
		meta.bagged.UTC().Format(time.RFC3339),
		meta.modified.UTC().Format(time.RFC3339),
	)

	if err != nil {
		sys.Error(err)
	}

	_, err = fmt.Fprintf(w.file, "Hash: %x\n\n", meta.hash)

	if err != nil {
		sys.Error(err)
	}
}

func (w *TextWrite) WriteLine(nr, grp int, s string) {
	_, err := fmt.Fprintf(w.file, "%d:%d:%s\n", nr, grp, s)

	if err != nil {
		sys.Error(err)
	}
}
