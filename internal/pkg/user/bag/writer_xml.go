package bag

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
	"time"

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
	Bagged   string `xml:"bagged"`
	Modified string `xml:"modified"`
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

func (w *XmlWriter) Init(f *os.File, _ bool, t string) {
	w.file = f

	w.bag = &xmlBag{
		Title: t,
	}

	buf, err := io.ReadAll(w.file)

	if err != nil {
		sys.Error(err)
		return
	}

	err = xml.Unmarshal(buf, &w.bag)

	if err != nil && err != io.EOF {
		sys.Error(err)
		return
	}
}

func (w *XmlWriter) Start() {
	w.entry = new(xmlEvidence)
}

func (w *XmlWriter) Finalize() {
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

	writeln(w.file, sb.String())
}

func (w *XmlWriter) WriteFile(p string, fs []string) {
	w.entry.Metadata.File = xmlFile{
		Path: p, Filters: fs,
	}
}

func (w *XmlWriter) WriteUser(u *user.User) {
	w.entry.Metadata.User = xmlUser{
		Login: u.Username, Name: u.Name,
	}
}

func (w *XmlWriter) WriteTime(t, f time.Time) {
	w.entry.Metadata.Time = xmlTime{
		Bagged:   t.UTC().Format(time.RFC3339),
		Modified: f.UTC().Format(time.RFC3339),
	}
}

func (w *XmlWriter) WriteHash(b []byte) {
	w.entry.Metadata.Hash = fmt.Sprintf("%x", b)
}

func (w *XmlWriter) WriteLines(ns []int, ss []string) {
	for i := 0; i < len(ss); i++ {
		w.entry.Lines.Line = append(w.entry.Lines.Line, xmlLine{
			Line: ns[i], Data: ss[i],
		})
	}
}
