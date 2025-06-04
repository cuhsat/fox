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

func (w *JsonWriter) Init(file *os.File, _ bool, title string) {
	w.file = file
	w.title = title
}

func (w *JsonWriter) Start() {
	w.entry = &jsonEvidence{
		Title: w.title,
	}
}

func (w *JsonWriter) Flush() {
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

	_, err = fmt.Fprintln(w.file, string(buf))

	if err != nil {
		sys.Error(err)
	}
}

func (w *JsonWriter) SetFile(path string, fs []string) {
	w.entry.Metadata.File = jsonFile{
		Path: path, Filters: fs,
	}
}

func (w *JsonWriter) SetUser(usr *user.User) {
	w.entry.Metadata.User = jsonUser{
		Login: usr.Username, Name: usr.Name,
	}
}

func (w *JsonWriter) SetTime(bag, mod time.Time) {
	w.entry.Metadata.Time = jsonTime{
		Bagged: bag.UTC(), Modified: mod.UTC(),
	}
}

func (w *JsonWriter) SetHash(sum []byte) {
	w.entry.Metadata.Hash = fmt.Sprintf("%x", sum)
}

func (w *JsonWriter) SetLine(nr int, s string) {
	w.entry.Lines = append(w.entry.Lines, jsonLine{
		Line: nr, Data: s,
	})
}
