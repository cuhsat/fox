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
    if len(s) > 0 {
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

func (s SMap) Wrap(w int) (r SMap) {
    for _, str := range s {
        s := str.Start
        l := str.Len

        // break string
        for l > w {
            r = append(r, &String{
                Nr: str.Nr,
                Start: s,
                End: s + w,
                Len: w,
            })

            s += w
            l -= w
        }

        r = append(r, &String{
            Nr: str.Nr,
            Start: s,
            End: str.End,
            Len: l,
        })
    }

    return
}

func (s SMap) Find(nr int) int {
    for i, str := range s {
        if str.Nr == nr {
            return i
        }
    }

    return -1
}

func (s SMap) Size() (w, h int) {
    for _, str := range s {
        w = max(w, str.Len)
    }

    h = len(s)

    return
}
