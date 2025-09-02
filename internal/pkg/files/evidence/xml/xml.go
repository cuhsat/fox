package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/files/evidence"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

const Ext = ".xml"

type Writer struct {
	file  *os.File // file handle
	bag   *Bag     // root element
	entry *entry   // current entry
}

type Bag struct {
	Title    string   `xml:",comment"`
	XMLName  xml.Name `xml:"bag"`
	Evidence []entry  `xml:"evidence"`
}

type entry struct {
	Metadata struct {
		File struct {
			Path    string   `xml:"Path"`
			Size    int64    `xml:"Size"`
			Filters []string `xml:"filter"`
		} `xml:"file"`

		User struct {
			Login string `xml:"login"`
			Name  string `xml:"Name"`
		} `xml:"User"`

		Time struct {
			Bagged   string `xml:"Bagged"`
			Modified string `xml:"Modified"`
		} `xml:"time"`

		Hash string `xml:"Hash"`
	} `xml:"metadata"`

	Lines struct {
		Line []line `xml:"line"`
	} `xml:"lines"`
}

type line struct {
	Nr   int    `xml:"nr,attr"`
	Grp  int    `xml:"grp,attr"`
	Data string `xml:",cdata"`
}

func New() *Writer {
	return &Writer{
		file: nil,
	}
}

func (w *Writer) Open(file *os.File, _ bool, title string) {
	w.file = file

	w.bag = &Bag{Title: title}

	buf, err := io.ReadAll(w.file)

	if err != nil {
		sys.Panic(err)
	}

	err = xml.Unmarshal(buf, &w.bag)

	if err != nil && err != io.EOF {
		sys.Panic(err)
	}
}

func (w *Writer) Begin() {
	w.entry = new(entry)
}

func (w *Writer) Flush() {
	var buf []byte
	var err error

	w.bag.Evidence = append(w.bag.Evidence, *w.entry)

	buf, err = xml.MarshalIndent(w.bag, "", "  ")

	if err != nil {
		sys.Error(err)
		return
	}

	_, err = w.file.Seek(0, 0)

	if err != nil {
		sys.Error(err)
		return
	}

	err = w.file.Truncate(0)

	if err != nil {
		sys.Error(err)
		return
	}

	var sb strings.Builder

	sb.WriteString(xml.Header)
	sb.Write(buf)

	_, err = fmt.Fprintln(w.file, sb.String())

	if err != nil {
		sys.Error(err)
	}
}

func (w *Writer) WriteMeta(meta evidence.Meta) {
	w.entry.Metadata.File.Path = meta.Path
	w.entry.Metadata.File.Size = meta.Size
	w.entry.Metadata.File.Filters = meta.Filters

	w.entry.Metadata.Hash = fmt.Sprintf("%x", meta.Hash)

	w.entry.Metadata.Time.Bagged = evidence.Utc(meta.Bagged)
	w.entry.Metadata.Time.Modified = evidence.Utc(meta.Modified)

	w.entry.Metadata.User.Login = meta.User.Username
	w.entry.Metadata.User.Name = meta.User.Name
}

func (w *Writer) WriteLine(nr, grp int, str string) {
	w.entry.Lines.Line = append(w.entry.Lines.Line, line{nr, grp, str})
}
