package bag

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"
)

type MarkdownWriter struct {
	file *os.File // file handle
}

func NewMarkdownWriter() *MarkdownWriter {
	return &MarkdownWriter{
		file: nil,
	}
}

func (w *MarkdownWriter) Init(f *os.File, n bool, t string) {
	w.file = f

	if n {
		writeln(w.file, fmt.Sprintf("# %s", t))
	}
}

func (w *MarkdownWriter) Start() {
	writeln(w.file, "")
}

func (w *MarkdownWriter) Finalize() {
	writeln(w.file, "")
}

func (w *MarkdownWriter) WriteFile(p string, fs []string) {
	var sb strings.Builder

	for _, f := range fs {
		sb.WriteString(fmt.Sprintf(" > `%s`", f))
	}

	writeln(w.file, fmt.Sprintf("## `%s`%s", p, sb.String()))
}

func (w *MarkdownWriter) WriteUser(u *user.User) {
	writeln(w.file, fmt.Sprintf("* %s (%s)", u.Username, u.Name))
}

func (w *MarkdownWriter) WriteTime(t, f time.Time) {
	writeln(w.file, fmt.Sprintf("* _%s_", t.UTC().String()))
	writeln(w.file, fmt.Sprintf("* _%s_", f.UTC().String()))
}

func (w *MarkdownWriter) WriteHash(b []byte) {
	writeln(w.file, fmt.Sprintf("* _%x_\n", b))
}

func (w *MarkdownWriter) WriteLines(ns []int, ss []string) {
	writeln(w.file, "```")

	for i := 0; i < len(ss); i++ {
		writeln(w.file, fmt.Sprintf("%08d  %v", ns[i], ss[i]))
	}

	writeln(w.file, "```")
}
