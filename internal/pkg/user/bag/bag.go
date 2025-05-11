package bag

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

const (
	Text     = "text"
	Xml      = "xml"
	Json     = "json"
	Jsonl    = "jsonl"
	Markdown = "markdown"
)

const (
	filename = "EVIDENCE"
)

const (
	header = "FORENSIC EXAMINER EVIDENCE BAG"
)

type Bag struct {
	Path string   // file path
	file *os.File // file handle
	key  string   // key phrase
	w    writer   // writer
}

type writer interface {
	Init(f *os.File, n bool, t string)
	Start()
	Finalize()
	WriteFile(p string, fs []string)
	WriteUser(u *user.User)
	WriteTime(t, f time.Time)
	WriteHash(b []byte)
	WriteLines(ns []int, ss []string)
}

func New(path, key, wt string) *Bag {
	var w writer
	var e string

	switch strings.ToLower(wt) {
	case Markdown:
		w = NewMarkdownWriter()
		e = ".md"
	case Jsonl:
		w = NewJsonWriter(false)
		e = ".jsonl"
	case Json:
		w = NewJsonWriter(true)
		e = ".json"
	case Xml:
		w = NewXmlWriter()
		e = ".xml"
	case Text:
		fallthrough
	default:
		w = NewTextWriter()
	}

	if len(path) == 0 {
		path = filename
	}

	if len(e) > 0 {
		path += e
	}

	return &Bag{
		Path: path,
		file: nil,
		key:  key,
		w:    w,
	}
}

func (bag *Bag) Put(h *heap.Heap) bool {
	if bag.file == nil && !bag.init() {
		return false
	}

	usr, err := user.Current()

	if err != nil {
		sys.Error(err)
	}

	sum, err := h.Sha256()

	if err != nil {
		sys.Error(err)
	}

	fi, err := os.Stat(h.Path)

	if err != nil {
		sys.Error(err)
	}

	bag.w.Start()

	bag.w.WriteFile(h.String(), *types.Filters())
	bag.w.WriteUser(usr)
	bag.w.WriteTime(time.Now(), fi.ModTime())
	bag.w.WriteHash(sum)

	smap := *h.SMap()
	mmap := *h.MMap()

	var ns []int
	var ss []string

	for _, s := range smap {
		ns = append(ns, s.Nr)
		ss = append(ss, string(mmap[s.Start:s.End]))
	}

	bag.w.WriteLines(ns, ss)

	bag.w.Finalize()

	bag.sign()

	return true
}

func (bag *Bag) Close() {
	if bag.file == nil {
		bag.file.Close()
	}
}

func (bag *Bag) init() bool {
	var err error

	is := sys.Exists(bag.Path)

	bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

	if err != nil {
		sys.Error(err)
		return false
	}

	bag.w.Init(bag.file, !is, header)

	return true
}

func (bag *Bag) sign() {
	var imp hash.Hash

	if len(bag.key) > 0 {
		imp = hmac.New(sha256.New, []byte(bag.key))
	} else {
		imp = sha256.New()
	}

	buf, err := os.ReadFile(bag.Path)

	if err != nil {
		sys.Error(err)
		return
	}

	imp.Write(buf)

	sum := fmt.Appendf(nil, "%x", imp.Sum(nil))

	err = os.WriteFile(bag.Path+".sha256", sum, 0600)

	if err != nil {
		sys.Error(err)
	}

	return
}

func writeln(f *os.File, s string) {
	_, err := fmt.Fprintln(f, s)

	if err != nil {
		sys.Error(err)
	}
}
