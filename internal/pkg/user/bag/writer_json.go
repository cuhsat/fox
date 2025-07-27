package bag

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hiforensics/fox/internal/pkg/sys"
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
	Title string `json:"_comment"`

	Metadata struct {
		File struct {
			Path    string   `json:"path"`
			Size    int64    `json:"size"`
			Filters []string `json:"filters"`
		} `json:"file"`

		User struct {
			Login string `json:"login"`
			Name  string `json:"name"`
		} `json:"user"`

		Time struct {
			Bagged   time.Time `json:"bagged"`
			Modified time.Time `json:"modified"`
		} `json:"time"`

		Hash string `json:"hash"`
	} `json:"metadata"`

	Lines []jsonLine `json:"lines"`
}

type jsonLine struct {
	Nr   int    `json:"nr"`
	Grp  int    `json:"grp"`
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

func (w *JsonWriter) WriteMeta(meta meta) {
	w.entry.Metadata.File.Path = meta.path
	w.entry.Metadata.File.Size = meta.size
	w.entry.Metadata.File.Filters = meta.filters

	w.entry.Metadata.Hash = fmt.Sprintf("%x", meta.hash)

	w.entry.Metadata.Time.Bagged = meta.bagged.UTC()
	w.entry.Metadata.Time.Modified = meta.modified.UTC()

	w.entry.Metadata.User.Login = meta.user.Username
	w.entry.Metadata.User.Name = meta.user.Name
}

func (w *JsonWriter) WriteLine(nr, grp int, s string) {
	w.entry.Lines = append(w.entry.Lines, jsonLine{nr, grp, s})
}
