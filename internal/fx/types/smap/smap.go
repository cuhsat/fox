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
        Len: len(m) - j,
    })

    return
}

func (s SMap) Indent(m mmap.MMap) (r SMap) {
    for _, str := range s {
        l := len(r)

        // blank line
        if str.Len == 0 {
            r = append(r, str)
        }

        pos := make(stack, 0)
        dqt := 0
        off := 0

        for i := str.Start; i < str.End; i++ {
            switch m[i] {
            case '{', '[':
                j := pos.Pop()

                if j >= 0 && j < i {
                    add(&r, str.Nr, j, i, off)
                }

                pos.Push(i+1)

                // bracket line
                add(&r, str.Nr, i, i+1, off)

                off += Space

            case '}', ']':
                j := pos.Pop()

                if j >= 0 && j < i {
                    add(&r, str.Nr, j, i, off)
                }

                off -= Space

                d := 1

                if i < str.Len && m[i+1] == ',' {
                    // off = 12
                    d += 1
                }

                // bracket line
                add(&r, str.Nr, i, i+d, off)

                pos.Push(i+d)

            case ',':
                if dqt % 2 != 0 {
                    continue
                }

                j := pos.Pop()

                if j >= 0 {
                    add(&r, str.Nr, j, i+1, off)
                }

                pos.Push(i+1)

            case '"':
                // parser look back
                if  m[max(i-1, 0)] != '\\' {
                    dqt += 1
                }
            }
        }

        // normal line
        if len(r) == l {
            r = append(r, str)
        }
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

func add(s *SMap, n, i, j, o int) {
    *s = append(*s, &String{
        Nr: n, Start: i, End: j, Len: j-i, Off: o,
    })
}
