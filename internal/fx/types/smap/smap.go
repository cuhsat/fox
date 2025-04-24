package smap

import (
    "github.com/edsrzf/mmap-go"
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
        i, spc, brk := str.Start, 0, true

        for j := str.Start; j < str.End; j++ {
            switch m[j] {
            case '{', '[':
                spc += 2
                brk = true
            case '}', ']':
                spc -= 2
                brk = true
            }

            if brk {
                brk = false

                r = append(r, &String{
                    Nr: str.Nr,
                    Start: i,
                    End: j,
                    Len: j - i,
                    Off: spc,
                })
            }
        }

        r = append(r, &String{
            Nr: str.Nr,
            Start: str.Start,
            End: i,
            Len: i - str.Start,
            Off: spc,
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
