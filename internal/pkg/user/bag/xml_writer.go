package bag

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/cuhsat/fx/internal/pkg/sys"
)

const (
	xmlIndent = "  "
)

type XmlWriter struct {
	file  *os.File // file handle
	title string   // export title
	bag   *xmlBag  // root element
}

type xmlBag struct {
	Title    string      `xml:",comment"`
	XMLName  xml.Name    `xml:"bag"`
	Evidence xmlEvidence `xml:"evidence"`
}

type xmlEvidence struct {
	Metadata xmlMetadata `xml:"metadata"`
	Lines    xmlLines    `xml:"lines"`
}

type xmlMetadata struct {
	File xmlFile `xml:"file"`
	User xmlUser `xml:"user"`
	Time xmlTime `xml:"time"`
	Hash string  `xml:"hash"`
}

type xmlFile struct {
	Path    string   `xml:"path"`
	Filters []string `xml:"filter"`
}

type xmlUser struct {
	Login string `xml:"login"`
	Name  string `xml:"name"`
}

type xmlTime struct {
	Bagged   time.Time `xml:"bagged"`
	Modified time.Time `xml:"modified"`
}

type xmlLines struct {
	Line []xmlLine `xml:"data"`
}

type xmlLine struct {
	Line int    `xml:"line,attr"`
	Data string `xml:",cdata"`
}

func NewXmlWriter() *XmlWriter {
	return &XmlWriter{
		file: nil,
	}
}

func (w *XmlWriter) Init(f *os.File, n bool, t string) {
	w.file = f
	w.title = t
	w.file = f
	w.title = t
}

func (w *XmlWriter) Start() {
	w.bag = &xmlBag{
		Title: w.title,
	}
}

func (w *XmlWriter) Finalize() {
	var buf []byte
	var err error

	buf, err = xml.MarshalIndent(w.bag, "", jsonIndent)

	if err != nil {
		sys.Error(err)
		return
	}

	writeln(w.file, xml.Header+string(buf))
}

func (w *XmlWriter) WriteFile(p string, fs []string) {
	w.bag.Evidence.Metadata.File.Filters = fs
}

func (w *XmlWriter) WriteUser(u *user.User) {
	w.bag.Evidence.Metadata.User = xmlUser{
		Login: u.Username, Name: u.Name,
	}
}

func (w *XmlWriter) WriteTime(t, f time.Time) {
	w.bag.Evidence.Metadata.Time = xmlTime{
		Bagged: t.UTC(), Modified: f.UTC(),
	}
}

func (w *XmlWriter) WriteHash(b []byte) {
	w.bag.Evidence.Metadata.Hash = fmt.Sprintf("%x", b)
}

func (w *XmlWriter) WriteLines(ns []int, ss []string) {
	for i := 0; i < len(ss); i++ {
		w.bag.Evidence.Lines.Line = append(w.bag.Evidence.Lines.Line, xmlLine{
			Line: ns[i], Data: ss[i],
		})
	}
}
