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
    MMap mmap.MMap // memory map
    SMap SMap      // string map current
    smap SMap      // string map original
    chain []SMap   // filter chain
    file *os.File  // file handle
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
        smap: s,
        file: f,
    }
}

func (h *Heap) AddFilter(value string) {
    h.SMap = h.filter([]byte(value))
}

func (h *Heap) DelFilter() {
    h.SMap = h.smap
}

func (h *Heap) ThrowAway() {
    h.MMap.Unmap()
    h.file.Close()
}

func (h *Heap) filter(b []byte) (f SMap) {
    ls := len(h.SMap)
    lc := min(runtime.GOMAXPROCS(0), ls)
    
    ch := make(chan SEntry, lc)

    go func() {
        for s := range ch {
            f = append(f, s)
        }        
    }()

    var wg sync.WaitGroup

    for c := 0; c < lc; c++ {
        min := c * ls / lc
        max := ((c+1) * ls) / lc

        wg.Add(1)

        go func() {
            defer wg.Done()

            h.search(min, max, b, ch)
        }()
    }

    wg.Wait()

    close(ch)

    h.sort(f)

    return
}

func (h *Heap) search(min, max int, b []byte, ch chan<- SEntry) {
    for _, s := range h.SMap[min:max] {
        if bytes.Contains(h.MMap[s.Start:s.End], b) {
            ch <- s
        }
    }
}

func (h *Heap) sort(s SMap) {
    slices.SortFunc(s, func(a, b SEntry) int {
        return cmp.Compare(a.Nr, b.Nr)
    })
}
