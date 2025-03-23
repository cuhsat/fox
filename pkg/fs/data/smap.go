package data

import (
    "github.com/edsrzf/mmap-go"
)

type SMap []*SEntry

type SEntry struct {
    Start, End, Len, Nr int
}

func smap(m mmap.MMap) (s SMap) {
    var j int = 0

    for i := 0; i < len(m); i++ {
        if m[i] == '\n' {
            s = append(s, &SEntry{
                Start: j,
                End: i,
                Len: i - j,
                Nr: len(s) + 1,
            })

            j = i + 1
        }
    }

    return
}
