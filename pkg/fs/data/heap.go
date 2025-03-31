package data

import (
    "bytes"
    "crypto/sha256"
    "cmp"
    "io"
    "os"
    "runtime"
    "slices"
    "sync"
 
    "github.com/cuhsat/cu/pkg/fs"
    "github.com/edsrzf/mmap-go"
)

type Heap struct {
    Path string    // file path
    Hash []byte    // file hash sum
    Chain []*SLink // filter chain
    MMap mmap.MMap // memory map
    SMap SMap      // string map current
    rmap SMap      // string map reserve
    file *os.File  // file handle
}

type SLink struct {
    Name string // filter name
    smap SMap   // filter string map
}

type chunk struct {
    min, max int
}

func NewHeap(path string) *Heap {
    f, err := os.OpenFile(path, os.O_RDONLY, fs.MODE_FILE)

    if err != nil {
        fs.Panic(err)
    }

    m, err := mmap.Map(f, mmap.RDONLY, 0)

    if err != nil {
        fs.Panic(err)
    }

    sha := sha256.New()

    _, err = io.Copy(sha, f)
    
    if err != nil {
        fs.Panic(err)
    }

    s := smap(m)

    return &Heap{
        Path: path,
        Hash: sha.Sum(nil),
        MMap: m,
        SMap: s,
        rmap: s,
        file: f,
    }
}

func (h *Heap) Reload() {
    h.ThrowAway()
    var err error

    h.file, err = os.OpenFile(h.Path, os.O_RDONLY, fs.MODE_FILE)

    if err != nil {
        fs.Panic(err)
    }

    h.MMap, err = mmap.Map(h.file, mmap.RDONLY, 0)

    if err != nil {
        fs.Panic(err)
    }

    sha := sha256.New()

    _, err = io.Copy(sha, h.file)
    
    if err != nil {
        fs.Panic(err)
    }

    h.Hash = sha.Sum(nil)
    h.SMap = smap(h.MMap)
    h.rmap = h.SMap
}

func (h *Heap) Lines() int {
    return len(h.rmap)
}

func (h *Heap) Copy() []byte {
    var b bytes.Buffer

    err := h.strings(&b)

    if err != nil {
        fs.Panic(err)
    }

    return b.Bytes()
}

func (h* Heap) Save() string {
    fn := h.Path

    for _, l := range h.Chain {
        fn += "-" + l.Name
    }

    f, err := os.OpenFile(fn, fs.FLAG_FILE, fs.MODE_FILE)

    if err != nil {
        fs.Panic(err)
    }

    defer f.Close()

    err = h.strings(f)

    if err != nil {
        fs.Panic(err)
    }

    return fn
}

func (h *Heap) AddFilter(value string) {
    h.SMap = h.filter([]byte(value))
    h.Chain = append(h.Chain, &SLink{
        Name: value,
        smap: h.SMap,
    })
}

func (h *Heap) DelFilter() {
    if len(h.Chain) > 0 {
        h.Chain = h.Chain[:len(h.Chain)-1]
    }

    if len(h.Chain) > 0 {
        h.SMap = h.Chain[len(h.Chain)-1].smap
    } else {
        h.SMap = h.rmap
    }
}

func (h *Heap) NoFilter() {
    h.Chain = h.Chain[:0]
    h.SMap = h.rmap
}

func (h *Heap) ThrowAway() {
    h.MMap.Unmap()
    h.file.Close()

    runtime.GC()
}

func (h *Heap) strings(w io.Writer) (err error) {
    for _, s := range h.SMap {
        _, err := w.Write([]byte(h.MMap[s.Start:s.End + 1]))

        if err != nil {
            return err
        }
    }

    return nil
}

func (h *Heap) filter(b []byte) (s SMap) {
    ch := make(chan *SEntry, len(h.SMap))

    defer close(ch)

    var wg sync.WaitGroup

    for _, c := range h.chunks() {
        wg.Add(1)

        go func() {
            h.search(ch, c, b)
            wg.Done()
        }()
    }

    wg.Wait()

    return h.gather(ch)
}

func (h *Heap) chunks() (c []*chunk) {
    n := len(h.SMap)
    m := min(runtime.GOMAXPROCS(0), n)
    
    for i := 0; i < m; i++ {
        c = append(c, &chunk{
            min: i * n / m,
            max: ((i+1) * n) / m,
        })
    }

    return
}

func (h *Heap) search(ch chan<- *SEntry, c *chunk, b []byte) {
    for _, s := range h.SMap[c.min:c.max] {
        if bytes.Contains(h.MMap[s.Start:s.End], b) {
            ch <- s
        }
    }
}

func (h *Heap) gather(ch <-chan *SEntry) (s SMap) {
    for len(ch) > 0 {
        s = append(s, <-ch)
    }

    slices.SortFunc(s, func(a, b *SEntry) int {
        return cmp.Compare(a.Nr, b.Nr)
    })

    return
}
