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

func (w *TextWriter) WriteFile(p string, f []string) {
    if len(f) > 0 {
        writeln(w.file, fmt.Sprintf("%s > %s", p, strings.Join(f, " > ")))
    } else {
        writeln(w.file, p)
    }
}

func (w *TextWriter) WriteUser(u *user.User) {
    writeln(w.file, fmt.Sprintf("%s (%s)", u.Username, u.Name))
}

func (w *TextWriter) WriteTime(t time.Time) {
    writeln(w.file, t.UTC().String())
    writeln(w.file, t.String())
}

func (w *TextWriter) WriteHash(b []byte) {
    writeln(w.file, fmt.Sprintf("%x\n", b))
}

func (w *TextWriter) WriteLine(n int, s string) {
    writeln(w.file, fmt.Sprintf("%08d  %v", n, s))
}
