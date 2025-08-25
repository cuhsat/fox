package bag

import (
	"fmt"
	"os"
	"strings"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

const width = 78

type TextWrite struct {
	file *os.File // file handle
}

func NewTextWriter() *TextWrite {
	return new(TextWrite)
}

func (w *TextWrite) Init(f *os.File, old bool, title string) {
	w.file = f

	if !old {
		w.write("%s\n%s\n", app.Ascii, title)
	}
}

func (w *TextWrite) Start() {
	w.write("\n")
}

func (w *TextWrite) Flush() {
	//
}

func (w *TextWrite) WriteMeta(meta meta) {
	var sb strings.Builder

	for _, f := range meta.filters {
		sb.WriteString("> ")
		sb.WriteString(f)
	}

	w.write("%s\n", strings.Repeat("=", width))
	w.write("File: %s %s (%d bytes)\n", meta.path, sb.String(), meta.size)
	w.write("User: %s (%s)\n", meta.user.Username, meta.user.Name)
	w.write("Time: %s / %s\n", utc(meta.bagged), utc(meta.modified))
	w.write("Hash: %x\n", meta.hash)
	w.write("%s\n", strings.Repeat("-", width))
}

func (w *TextWrite) WriteLine(nr, grp int, s string) {
	w.write("%d:%d:%s\n", nr, grp, s)
}

func (w *TextWrite) write(format string, a ...any) {
	_, err := fmt.Fprintf(w.file, format, a...)

	if err != nil {
		sys.Error(err)
	}
}
