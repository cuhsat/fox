package ui

import (
    "strings"

    "github.com/cuhsat/cu/pkg/fs"
    //"github.com/nsf/termbox-go"
)

type Buffer struct {
    dx int // delta x
    dy int // delta y
}

func NewBuffer() *Buffer {
    return &Buffer{
        dx: 0,
        dy: 0,
    }
}

func (b *Buffer) Render(x, y int, heap *fs.Heap) {
    for l, s := range heap.SMap[b.dy:] {
        if l > page-1 {
            break
        }

        str := string(heap.MMap[s.Start:s.End+1])

        n := min(b.dx, len(str))

        printLine(x, y + l, str[n:], BufferFg, BufferBg)

        for z, c := range heap.Chain {
            mark(x, y + l, z, str[n:], c.Name)
        }
    }

    for l := len(heap.SMap[b.dy:]); l < height; l++ {
        printLine(x, y + l, "", BufferFg, BufferBg)
    }
}

func (b *Buffer) Reset() {
    b.dx = 0
    b.dy = 0
}

func (b *Buffer) GoToBegin() {
    b.dy = 0
}

func (b *Buffer) GoToEnd() { //
    b.dy = min(data - page, data)
}

func (b *Buffer) PageUp() {
    b.dy = max(b.dy - page, 0)
}

func (b *Buffer) PageDown() { //
    b.dy = min(b.dy + page, data - page, data)
}

func (b *Buffer) ScrollUp(delta int) {
    b.dy = max(b.dy - delta, 0)
}

func (b *Buffer) ScrollDown(delta int) {
    b.dy = min(b.dy + delta, max(data - page, 0))
}

func (b *Buffer) ScrollLeft(delta int) {
    b.dx = max(b.dx - delta, 0)
}

func (b *Buffer) ScrollRight(delta int) {
    b.dx = min(b.dx + delta, width)
}

func mark(x, y, z int, s, n string) {
    i := strings.Index(s, n)

    // c := termbox.Attribute(10 - z)

    if i > -1 {
        print(x + i, y, n, SearchFg, SearchBg)

        i++

        mark(x + i, y, z, s[i:], n)
    }
}
