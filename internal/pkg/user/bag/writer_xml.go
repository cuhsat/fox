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

func (w *XmlWriter) SetFile(path string, fs []string) {
	w.entry.Metadata.File = xmlFile{
		Path: path, Filters: fs,
	}
}

func (w *XmlWriter) SetUser(usr *user.User) {
	w.entry.Metadata.User = xmlUser{
		Login: usr.Username, Name: usr.Name,
	}
}

func (w *XmlWriter) SetTime(bag, mod time.Time) {
	w.entry.Metadata.Time = xmlTime{
		Bagged:   bag.UTC().Format(time.RFC3339),
		Modified: mod.UTC().Format(time.RFC3339),
	}
}

func (w *XmlWriter) SetHash(sum []byte) {
	w.entry.Metadata.Hash = fmt.Sprintf("%x", sum)
}

func (w *XmlWriter) SetLine(nr int, s string) {
	w.entry.Lines.Line = append(w.entry.Lines.Line, xmlLine{
		Line: nr, Data: s,
	})
}
