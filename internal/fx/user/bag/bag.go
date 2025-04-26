package bag

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "hash"
    "os"
    "os/user"
    "time"

    "github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/types"
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
    key string    // key phrase
    exp exporter  // exporter
}

type exporter interface {
    Init(f *os.File, n bool, t string)
    Start()
    Finalize()
    ExportFile(p string, f []string)
    ExportUser(u *user.User)
    ExportTime(t time.Time)
    ExportHash(b []byte)
    ExportLine(n int, s string)
}

func New(path, key string, e types.Export) *Bag {
    var exp exporter
    var ext string

    switch e {
    case types.Jsonl:
        exp = NewJsonExporter(false)
        ext = ".jsonl"
    case types.Json:
        exp = NewJsonExporter(true)
        ext = ".json"
    default:
        exp = NewTextExporter()
    }

    if len(path) == 0 {
        path = filename + ext
    }

    return &Bag{
        Path: path,
        file: nil,
        key: key,
        exp: exp,
    }
}

func (bag *Bag) Put(h *heap.Heap) bool {
    if bag.file == nil && !bag.init() {
        return false
    }

    usr, err := user.Current()

    if err != nil {
        fx.Error(err)
    }

    bag.exp.Start()

    bag.exp.ExportFile(h.String(), *types.GetFilters())
    bag.exp.ExportUser(usr)
    bag.exp.ExportTime(time.Now())
    bag.exp.ExportHash(h.Sha256())

    for _, s := range h.SMap {
        bag.exp.ExportLine(s.Nr, string(h.MMap[s.Start:s.End]))
    }

    bag.exp.Finalize()

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

    is := fx.Exists(bag.Path)

    bag.file, err = os.OpenFile(bag.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

    if err != nil {
        fx.Error(err)
        return false
    }

    bag.exp.Init(bag.file, !is, header)

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
        fx.Error(err)
        return
    }

    imp.Write(buf)

    sum := []byte(fmt.Sprintf("%x", imp.Sum(nil)))

    err = os.WriteFile(bag.Path + ".sha256", sum, 0600)

    if err != nil {
        fx.Error(err)
    }

    return
}

func writeln(f *os.File, s string) {
    _, err := fmt.Fprintln(f, s)

    if err != nil {
        fx.Error(err)
    }
}
