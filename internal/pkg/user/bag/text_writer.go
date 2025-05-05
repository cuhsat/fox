package bag

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"
)

type TextWriter struct {
	file *os.File // file handle
}

func NewTextWriter() *TextWriter {
	return &TextWriter{
		file: nil,
	}
}

func (w *TextWriter) Init(f *os.File, n bool, t string) {
	w.file = f

	if n {
		writeln(w.file, t)
	}
}

func (w *TextWriter) Start() {
	writeln(w.file, "")
}

func (w *TextWriter) Finalize() {
	writeln(w.file, "")
}

func (w *TextWriter) WriteFile(p string, fs []string) {
	var sb strings.Builder

	for _, f := range fs {
		sb.WriteString(fmt.Sprintf(" > %s", f))
	}

	writeln(w.file, fmt.Sprintf("%s%s", p, sb.String()))
}

func (w *TextWriter) WriteUser(u *user.User) {
	writeln(w.file, fmt.Sprintf("%s (%s)", u.Username, u.Name))
}

func (w *TextWriter) WriteTime(t, f time.Time) {
	writeln(w.file, t.UTC().String())
	writeln(w.file, f.UTC().String())
}

func (w *TextWriter) WriteHash(b []byte) {
	writeln(w.file, fmt.Sprintf("%x\n", b))
}

func (w *TextWriter) WriteLines(ns []int, ss []string) {
	for i := 0; i < len(ss); i++ {
		writeln(w.file, fmt.Sprintf("%08d  %v", ns[i], ss[i]))
	}
}
