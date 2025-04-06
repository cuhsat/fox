package smap

import (
    "github.com/edsrzf/mmap-go"
)

const (
    Break = '\n'
)

type SMap []*String

type String struct {
    Nr    int
    Start int
    End   int
    Len   int
}

func Map(m mmap.MMap) (s SMap) {
    var j int = 0

    // append strings
    for i := 0; i < len(m); i++ {
        if m[i] == Break {
            s = append(s, &String{
                Nr: len(s) + 1,
                Start: j,
                End: i,
                Len: i - j,
            })

            j = i + 1
        }
    }

    // append remaining string
    if len(s) > 0 && m[len(m)-1] != Break {
        l := s[len(s)-1]
        s = append(s, &String{
            Nr: l.Nr + 1,
            Start: l.End + 1,
            End: len(m),
            Len: len(m) - l.End,
        })
    }

    return
}

func (s SMap) Length() (l int) {
    for _, str := range s {
        l = max(l, str.Len)
    }

    return
}
