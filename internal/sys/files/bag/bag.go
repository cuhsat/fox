package bag

import (
    "bufio"
    "fmt"
    "io"
    "math"
    "os"
    "strings"
    "time"

    "github.com/cuhsat/fx/internal/sys"
    "github.com/cuhsat/fx/internal/sys/heap"
    "github.com/cuhsat/fx/internal/sys/text"
    "github.com/cuhsat/fx/internal/sys/types"
)

const (
    File = "EVIDENCE"
)

const (
    Header = "Forensic Examiner - Evidence Bag"
)

type Bag struct {
    path string   // file path
    file *os.File // file handle
}

func NewBag(p string) *Bag {
    return &Bag{
        path: p,
    }
}

func (b *Bag) Init() {
    var err error

    if len(b.path) == 0 {
        b.path = File
    }

    b.file, err = os.OpenFile(b.path, sys.O_EVIDENCE, 0644)

    if err != nil {
        sys.Fatal(err)
    }

    _, err = b.file.WriteString(Header + "\n")

    if err != nil {
        sys.Fatal(err)
    }
}

func (b *Bag) Put(h *heap.Heap) {
    if b.file == nil {
        b.Init()
    }

    t := fmt.Sprintf("%s :: %s", time.Now().UTC(), h.String())
    f := types.GetFilters()

    if len(*f) > 0 {
        t = fmt.Sprintf("%s > %s", t, f)
    }

    _, err := b.file.WriteString(text.Block(t, len(t)+4) + "\n")

    if err != nil {
        sys.Fatal(err)
    }

    len_nr := int(math.Log10(float64(h.Length()))) + 1

    for _, s := range h.SMap {
        l := fmt.Sprintf("%0*d: %v\n", len_nr, s.Nr, string(h.MMap[s.Start:s.End]))

        _, err = b.file.WriteString(l)

        if err != nil {
            sys.Fatal(err)
        }
    }
}

func (b *Bag) Close() {
    if b.file != nil {
        b.file.Close()        
    }
}

func IsEvidence(p string) bool {
     f, err := os.Open(p)

     if err != nil {
         sys.Fatal(err)
     }

     defer f.Close()

     r := bufio.NewReader(f)

     l, _, err := r.ReadLine()

     switch err {
     case io.EOF:
         return false
     case nil:
         return strings.Compare(string(l), Header) == 0
     default:
         sys.Fatal(err)
     }

     return false
}
