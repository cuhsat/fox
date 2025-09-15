package bag

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/cuhsat/fox/internal"
	"github.com/cuhsat/fox/internal/pkg/files/evidence"
	"github.com/cuhsat/fox/internal/pkg/files/evidence/json"
	"github.com/cuhsat/fox/internal/pkg/files/evidence/plain"
	"github.com/cuhsat/fox/internal/pkg/files/evidence/sqlite"
	"github.com/cuhsat/fox/internal/pkg/files/evidence/text"
	"github.com/cuhsat/fox/internal/pkg/files/evidence/url"
	"github.com/cuhsat/fox/internal/pkg/files/evidence/xml"
	"github.com/cuhsat/fox/internal/pkg/files/schema/ecs"
	"github.com/cuhsat/fox/internal/pkg/files/schema/hec"
	"github.com/cuhsat/fox/internal/pkg/files/schema/raw"
	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
	"github.com/cuhsat/fox/internal/pkg/types"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
)

type Bag struct {
	Path string        // file path
	Mode flags.BagMode // file mode

	file *os.File          // file handle
	name string            // case name
	key  string            // key phrase
	url  string            // url address
	ws   []evidence.Writer // writers
}

func New() *Bag {
	var ws []evidence.Writer

	flg := flags.Get().Bag

	path := flg.File

	if len(path) == 0 {
		path = flags.BagFile
	}

	switch flg.Mode {
	case flags.BagModeNone:
	case flags.BagModeSqlite:
		ws = append(ws, sqlite.New())
		path += sqlite.Ext

	case flags.BagModeJsonl:
		ws = append(ws, json.New(false))
		path += json.ExtPretty

	case flags.BagModeJson:
		ws = append(ws, json.New(true))
		path += json.Ext

	case flags.BagModeXml:
		ws = append(ws, xml.New())
		path += xml.Ext

	case flags.BagModeText:
		ws = append(ws, text.New())
		path += text.Ext

	default:
		ws = append(ws, plain.New())
		path += plain.Ext
	}

	if len(flg.Url) > 0 {
		if flg.Ecs {
			ws = append(ws, url.New(flg.Url, ecs.New()))
		} else if flg.Hec {
			ws = append(ws, url.New(flg.Url, hec.New()))
		} else {
			ws = append(ws, url.New(flg.Url, raw.New()))
		}
	}

	return &Bag{
		Path: path,
		Mode: flg.Mode,
		name: flg.Case,
		key:  flg.Key,
		url:  flg.Url,
		file: nil,
		ws:   ws,
	}
}

func (bag *Bag) String() string {
	if bag.file != nil {
		return bag.Path
	} else {
		return bag.url
	}
}

func (bag *Bag) Close() {
	if bag.file != nil {
		_ = bag.file.Close()
	}
}

func (bag *Bag) Put(h *heap.Heap) bool {
	bag.init()

	usr, err := user.Current()

	if err != nil {
		sys.Error(err)
	}

	sum, err := h.HashSum(types.SHA256)

	if err != nil {
		sys.Error(err)
	}

	abs, err := filepath.Abs(h.String())

	if err != nil {
		sys.Error(err)
	}

	for _, w := range bag.ws {
		w.Begin()

		w.WriteMeta(evidence.Meta{
			User:     usr,
			Name:     bag.name,
			Path:     abs,
			Size:     h.Len(),
			Hash:     sum,
			Filters:  h.Patterns(),
			Bagged:   now(),
			Modified: mod(h),
		})

		for _, str := range *h.FMap() {
			w.WriteLine(str.Nr, str.Grp, str.Str)
		}

		w.Flush()
	}

	if bag.file != nil {
		bag.sign()
	}

	return len(bag.ws) > 0
}

func (bag *Bag) init() {
	old := fs.Exists(bag.Path)

	if bag.Mode != flags.BagModeNone {
		var err error

		bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

		if err != nil {
			sys.Panic(err)
		}
	}

	title := fmt.Sprintf("Forensic Examiner Evidence Bag (%s)", info.Version)

	for _, w := range bag.ws {
		w.Open(bag.file, old, title)
	}
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

	sum := base64.StdEncoding.EncodeToString(imp.Sum(nil))

	err = os.WriteFile(bag.Path+".sig", []byte(sum), 0600)

	if err != nil {
		sys.Error(err)
	}

	return
}

func now() time.Time {
	return time.Now().UTC()
}

func mod(h *heap.Heap) time.Time {
	mt := now()

	if h.Type == types.Regular {
		fi, err := os.Stat(h.Base)

		if err == nil {
			mt = fi.ModTime()
		} else {
			sys.Error(err)
		}
	}

	return mt
}
