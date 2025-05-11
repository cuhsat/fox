package bag

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

const (
	jsonIndent = "  "
)

type JsonWriter struct {
	file   *os.File      // file handle
	pretty bool          // export pretty
	title  string        // export title
	entry  *jsonEvidence // current entry
}

type jsonEvidence struct {
	Title    string       `json:"_comment"`
	Metadata jsonMetadata `json:"metadata"`
	Lines    []jsonLine   `json:"lines"`
}

type jsonMetadata struct {
	File jsonFile `json:"file"`
	User jsonUser `json:"user"`
	Time jsonTime `json:"time"`
	Hash string   `json:"hash"`
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
	Bagged   time.Time `json:"bagged"`
	Modified time.Time `json:"modified"`
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
	w.entry = &jsonEvidence{
		Title: w.title,
	}
}

func (w *JsonWriter) Finalize() {
	var buf []byte
	var err error

	if w.pretty {
		buf, err = json.MarshalIndent(w.entry, "", jsonIndent)
	} else {
		buf, err = json.Marshal(w.entry)
	}

	if err != nil {
		sys.Error(err)
		return
	}

	writeln(w.file, string(buf))
}

func (w *JsonWriter) WriteFile(p string, fs []string) {
	w.entry.Metadata.File = jsonFile{
		Path: p, Filters: fs,
	}
}

func (w *JsonWriter) WriteUser(u *user.User) {
	w.entry.Metadata.User = jsonUser{
		Login: u.Username, Name: u.Name,
	}
}

func (w *JsonWriter) WriteTime(t, f time.Time) {
	w.entry.Metadata.Time = jsonTime{
		Bagged: t.UTC(), Modified: f.UTC(),
	}
}

func (w *JsonWriter) WriteHash(b []byte) {
	w.entry.Metadata.Hash = fmt.Sprintf("%x", b)
}

func (w *JsonWriter) WriteLines(ns []int, ss []string) {
	for i := 0; i < len(ss); i++ {
		w.entry.Lines = append(w.entry.Lines, jsonLine{
			Line: ns[i], Data: ss[i],
		})
	}
}
