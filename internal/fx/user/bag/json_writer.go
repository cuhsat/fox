package bag

import (
    "encoding/json"
    "fmt"
    "os"
    "os/user"
    "time"

    "github.com/cuhsat/fx/internal/fx/sys"
)

const (
    indent = "  "
)

type JsonWriter struct {
    file *os.File    // file handle
    pretty bool      // export pretty
    title string     // export title
    entry *jsonEntry // current entry
}

type jsonEntry struct {
    Title string    `json:"_comment"`
    Meta jsonMeta   `json:"meta"`
    Data []jsonLine `json:"data"`
}

type jsonMeta struct {
    File  jsonFile `json:"file"`
    User  jsonUser `json:"user"`
    Time  jsonTime `json:"time"`
    Hash  string   `json:"hash"`    
}

type jsonFile struct {
    Path    string   `json:"path"`
    Filters []string `json:"filters"`
}

type jsonUser struct {
    Login string `json:"login"`
    Name  string `json:"name"`
}

type jsonTime struct {
    UTC   time.Time `json:"utc"`
    Local time.Time `json:"local"`
}

type jsonLine struct {
    Line int    `json:"line"`
    Data string `json:"data"`
}

func NewJsonWriter(pretty bool) *JsonWriter {
    return &JsonWriter{
        pretty: pretty,
    }
}

func (w *JsonWriter) Init(f *os.File, _ bool, t string) {
    w.file = f
    w.title = t
}

func (w *JsonWriter) Start() {
    w.entry = &jsonEntry{
        Title: w.title,
    }
}

func (w *JsonWriter) Finalize() {
    var buf []byte
    var err error

    if w.pretty {
        buf, err = json.MarshalIndent(w.entry, "", indent)
    } else {
        buf, err = json.Marshal(w.entry)
    }

    if err != nil {
        sys.Error(err)
        return
    }

    writeln(w.file, string(buf))
}

func (w *JsonWriter) WriteFile(p string, f []string) {
    w.entry.Meta.File = jsonFile{
        Path: p, Filters: f,
    }
}

func (w *JsonWriter) WriteUser(u *user.User) {
    w.entry.Meta.User = jsonUser{
        Login: u.Username, Name: u.Name,
    }
}

func (w *JsonWriter) WriteTime(t time.Time) {
    w.entry.Meta.Time = jsonTime{
        UTC: t.UTC(), Local: t,
    }
}

func (w *JsonWriter) WriteHash(b []byte) {
    w.entry.Meta.Hash = fmt.Sprintf("%x", b)
}

func (w *JsonWriter) WriteLine(n int, s string) {
    w.entry.Data = append(w.entry.Data, jsonLine{
        Line: n, Data: s,
    })
}
