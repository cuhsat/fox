package fs

import (
    "bytes"
    "cmp"
    "os"
    "runtime"
    "slices"
    "sync"
 
    "github.com/edsrzf/mmap-go"
)

type Heap struct {
    Path string    // file path
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
    f, err := os.OpenFile(path, os.O_RDONLY, MODE_FILE)

    if err != nil {
        Panic(err)
    }

    m, err := mmap.Map(f, mmap.RDONLY, 0)

    if err != nil {
        Panic(err)
    }

    s := smap(m)

    return &Heap{
        Path: path,
        MMap: m,
        SMap: s,
        rmap: s,
        file: f,
    }
}

func (h *Heap) Copy() []byte {
    return h.data()
}

func (h* Heap) Save() string {
    fn := h.Path

    for _, l := range h.Chain {
        fn += "-" + l.Name
    }

    err := os.WriteFile(fn, h.data(), MODE_FILE)

    if err != nil {
        Error(err)
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

func (h *Heap) data() []byte {
    var b bytes.Buffer

    for _, s := range h.SMap {
        _, err := b.Write([]byte(h.MMap[s.Start:s.End]))

        if err != nil {
            Error(err)
        }

        err = b.WriteByte('\n')

        if err != nil {
            Error(err)
        }
    }

    return b.Bytes()
}
