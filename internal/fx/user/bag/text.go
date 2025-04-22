package bag

import (
    "fmt"
    "os"
    "os/user"
    "strings"
    "time"
)

type TextExporter struct {
    file *os.File // file handle
}

func NewTextExporter() *TextExporter {
    return &TextExporter{
        file: nil,
    }
}

func (te *TextExporter) Init(f *os.File, n bool, t string) {
    te.file = f

    if n {
        writeln(te.file, t)
    }
}

func (te *TextExporter) Start() {
    writeln(te.file, "")
}

func (te *TextExporter) Finalize() {
    writeln(te.file, "")
}

func (te *TextExporter) ExportFile(p string, f []string) {
    writeln(te.file, fmt.Sprintf("%s > %s", p, strings.Join(f, " > ")))
}

func (te *TextExporter) ExportUser(u *user.User) {
    writeln(te.file, fmt.Sprintf("%s (%s)", u.Username, u.Name))
}

func (te *TextExporter) ExportTime(t time.Time) {
    writeln(te.file, t.UTC().String())
    writeln(te.file, t.String())
}

func (te *TextExporter) ExportHash(b []byte) {
    writeln(te.file, fmt.Sprintf("%x\n", b))
}

func (te *TextExporter) ExportLine(n int, s string) {
    writeln(te.file, fmt.Sprintf("%08d  %v", n, s))
}
