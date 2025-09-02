package text

import (
	"fmt"
	"os"
	"strings"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/files/evidence"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

const Ext = ".bag"

const width = 78

type Writer struct {
	file *os.File // file handle
}

func New() *Writer {
	return new(Writer)
}

func (w *Writer) Open(f *os.File, old bool, title string) {
	w.file = f

	if !old {
		w.write("%s\n%s\n", app.Art, title)
	}
}

func (w *Writer) Begin() {
	w.write("\n")
}

func (w *Writer) Flush() {
	w.write("%s\n", strings.Repeat("=", width))
}

func (w *Writer) WriteMeta(meta evidence.Meta) {
	var sb strings.Builder

	for _, f := range meta.Filters {
		sb.WriteString("> ")
		sb.WriteString(f)
	}

	w.write("%s\n", strings.Repeat("=", width))
	w.write("File: %s %s (%d bytes)\n", meta.Path, sb.String(), meta.Size)
	w.write("User: %s (%s)\n", meta.User.Username, meta.User.Name)
	w.write("Time: %s / %s\n", evidence.Utc(meta.Bagged), evidence.Utc(meta.Modified))
	w.write("Hash: %x\n", meta.Hash)
	w.write("%s\n", strings.Repeat("-", width))
}

func (w *Writer) WriteLine(nr, grp int, str string) {
	w.write("%d:%d:%s\n", nr, grp, str)
}

func (w *Writer) write(format string, a ...any) {
	_, err := fmt.Fprintf(w.file, format, a...)

	if err != nil {
		sys.Error(err)
	}
}
