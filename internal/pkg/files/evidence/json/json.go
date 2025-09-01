package json

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cuhsat/fox/internal/pkg/files/evidence"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

const (
	Ext       = ".json"
	ExtPretty = ".jsonl"
)

type Writer struct {
	file   *os.File // file handle
	pretty bool     // export pretty
	title  string   // export title
	entry  *entry   // current entry
}

type entry struct {
	Title string `json:"_comment"`

	Metadata struct {
		File struct {
			Path    string   `json:"Path"`
			Size    int64    `json:"Size"`
			Filters []string `json:"Filters"`
		} `json:"file"`

		User struct {
			Login string `json:"login"`
			Name  string `json:"Name"`
		} `json:"User"`

		Time struct {
			Bagged   time.Time `json:"Bagged"`
			Modified time.Time `json:"Modified"`
		} `json:"time"`

		Hash string `json:"Hash"`
	} `json:"metadata"`

	Lines []line `json:"lines"`
}

type line struct {
	Nr   int    `json:"nr"`
	Grp  int    `json:"grp"`
	Data string `json:"data"`
}

func New(pretty bool) *Writer {
	return &Writer{
		pretty: pretty,
	}
}

func (w *Writer) Open(file *os.File, _ bool, title string) {
	w.file = file
	w.title = title
}

func (w *Writer) Begin() {
	w.entry = &entry{
		Title: w.title,
	}
}

func (w *Writer) Flush() {
	var buf []byte
	var err error

	if w.pretty {
		buf, err = json.MarshalIndent(w.entry, "", "  ")
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

func (w *Writer) WriteMeta(meta evidence.Meta) {
	w.entry.Metadata.File.Path = meta.Path
	w.entry.Metadata.File.Size = meta.Size
	w.entry.Metadata.File.Filters = meta.Filters

	w.entry.Metadata.Hash = fmt.Sprintf("%x", meta.Hash)

	w.entry.Metadata.Time.Bagged = meta.Bagged.UTC()
	w.entry.Metadata.Time.Modified = meta.Modified.UTC()

	w.entry.Metadata.User.Login = meta.User.Username
	w.entry.Metadata.User.Name = meta.User.Name
}

func (w *Writer) WriteLine(nr, grp int, s string) {
	w.entry.Lines = append(w.entry.Lines, line{nr, grp, s})
}
