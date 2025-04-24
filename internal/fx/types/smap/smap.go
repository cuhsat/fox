package smap

import (
    "github.com/edsrzf/mmap-go"
)

const (
    Space = 2
)

type SMap []*String

type String struct {
    Nr    int
    Start int
    End   int
    Len   int
    Off   int
}

func Map(m mmap.MMap) (s SMap) {
    i, j := 0, 0

    for ; i < len(m); i++ {
        if m[i] != '\n' {
            continue
        }

        s = append(s, &String{
            Nr: len(s)+1,
            Start: j,
            End: i,
            Len: i - j,
        })

        j = i+1
    }

    s = append(s, &String{
        Nr: len(s)+1,
        Start: j,
        End: len(m),
        Len: len(m) - i,
    })

    return
}

func (s SMap) Indent(m mmap.MMap) (r SMap) {
    for _, str := range s {
        i, off, brk := str.Start, 0, true

        for j := str.Start; j < str.End; j++ {
            switch m[j] {
            case '{', '[':
                off += Space
                brk = true
            case '}', ']':
                off -= Space
                brk = true
            }

            if brk {
                brk = false

                r = append(r, &String{
                    Nr: str.Nr,
                    Start: j,
                    End: j+1,
                    Len: 1,
                    Off: off - Space,
                })
            }
        }

        off = 0

        r = append(r, &String{
            Nr: str.Nr,
            Start: str.Start,
            End: i,
            Len: i - str.Start,
            Off: off,
        })
    }

    return
}

func (s SMap) Wrap(w int) (r SMap) {
    for _, str := range s {
        s := str.Start
        l := str.Len

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
