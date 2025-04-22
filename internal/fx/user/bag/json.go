package bag

import (
    "encoding/json"
    "fmt"
    "os"
    "os/user"
    "time"

    "github.com/cuhsat/fx/internal/fx"
)

const (
    indent = "  "
)

type JsonExporter struct {
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

func NewJsonExporter(pretty bool) *JsonExporter {
    return &JsonExporter{
        pretty: pretty,
    }
}

func (je *JsonExporter) Init(f *os.File, _ bool, t string) {
    je.file = f
    je.title = t
}

func (je *JsonExporter) Start() {
    je.entry = &jsonEntry{
        Title: je.title,
    }
}

func (je *JsonExporter) Finalize() {
    var buf []byte
    var err error

    if je.pretty {
        buf, err = json.MarshalIndent(je.entry, "", indent)
    } else {
        buf, err = json.Marshal(je.entry)
    }

    if err != nil {
        fx.Error(err)
        return
    }

    writeln(je.file, string(buf))
}

func (je *JsonExporter) ExportFile(p string, f []string) {
    je.entry.Meta.File = jsonFile{
        Path: p, Filters: f,
    }
}

func (je *JsonExporter) ExportUser(u *user.User) {
    je.entry.Meta.User = jsonUser{
        Login: u.Username, Name: u.Name,
    }
}

func (je *JsonExporter) ExportTime(t time.Time) {
    je.entry.Meta.Time = jsonTime{
        UTC: t.UTC(), Local: t,
    }
}

func (je *JsonExporter) ExportHash(b []byte) {
    je.entry.Meta.Hash = fmt.Sprintf("%x", b)
}

func (je *JsonExporter) ExportLine(n int, s string) {
    je.entry.Data = append(je.entry.Data, jsonLine{
        Line: n, Data: s,
    })
}
