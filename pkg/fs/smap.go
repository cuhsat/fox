package fs

import (
    "github.com/edsrzf/mmap-go"
)

type SMap []*SEntry

type SEntry struct {
    Start, End, Nr int
}

func smap(m mmap.MMap) (s SMap) {
    var j int = 0

    for i := 0; i < len(m); i++ {
        if m[i] == '\n' {
            s = append(s, &SEntry{
                Start: j,
                End: i,
                Nr: len(s),
            })

            j = i+1
        }
    }

    return
}
