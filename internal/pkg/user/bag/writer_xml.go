package bag

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

const (
	xmlIndent = "  "
)

type XmlWriter struct {
	file  *os.File     // file handle
	bag   *xmlBag      // root element
	entry *xmlEvidence // current entry
}

type xmlBag struct {
	Title    string        `xml:",comment"`
	XMLName  xml.Name      `xml:"bag"`
	Evidence []xmlEvidence `xml:"evidence"`
}

type xmlEvidence struct {
	Metadata struct {
		File struct {
			Path    string   `xml:"path"`
			Size    int64    `xml:"size"`
			Filters []string `xml:"filter"`
		} `xml:"file"`

		User struct {
			Login string `xml:"login"`
			Name  string `xml:"name"`
		} `xml:"user"`

		Time struct {
			Bagged   string `xml:"bagged"`
			Modified string `xml:"modified"`
		} `xml:"time"`

		Hash string `xml:"hash"`
	} `xml:"metadata"`

	Lines struct {
		Line []xmlLine `xml:"line"`
	} `xml:"lines"`
}

type xmlLine struct {
	Nr   int    `xml:"nr,attr"`
	Grp  int    `xml:"grp,attr"`
	Data string `xml:",cdata"`
}

func NewXmlWriter() *XmlWriter {
	return &XmlWriter{
		file: nil,
	}
}

func (w *XmlWriter) Init(file *os.File, _ bool, title string) {
	w.file = file

	w.bag = &xmlBag{Title: title}

	buf, err := io.ReadAll(w.file)

	if err != nil {
		sys.Panic(err)
	}

	err = xml.Unmarshal(buf, &w.bag)

	if err != nil && err != io.EOF {
		sys.Panic(err)
	}
}

func (w *XmlWriter) Start() {
	w.entry = new(xmlEvidence)
}

func (w *XmlWriter) Flush() {
	var buf []byte
	var err error

	w.bag.Evidence = append(w.bag.Evidence, *w.entry)

	buf, err = xml.MarshalIndent(w.bag, "", xmlIndent)

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

func (w *XmlWriter) WriteMeta(meta meta) {
	w.entry.Metadata.File.Path = meta.path
	w.entry.Metadata.File.Size = meta.size
	w.entry.Metadata.File.Filters = meta.filters

	w.entry.Metadata.Hash = fmt.Sprintf("%x", meta.hash)

	w.entry.Metadata.Time.Bagged = utc(meta.bagged)
	w.entry.Metadata.Time.Modified = utc(meta.modified)

	w.entry.Metadata.User.Login = meta.user.Username
	w.entry.Metadata.User.Name = meta.user.Name
}

func (w *XmlWriter) WriteLine(nr, grp int, s string) {
	w.entry.Lines.Line = append(w.entry.Lines.Line, xmlLine{nr, grp, s})
}
