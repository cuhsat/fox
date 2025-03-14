package ui

import (
    "github.com/cuhsat/cu/pkg/fs"
    "github.com/nsf/termbox-go"
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
    for ty, s := range heap.SMap[b.dy:] {
        if ty > page-1 {
            break
        }

        str := string(heap.MMap[s.Start:s.End+1])

        n := min(b.dx, len(str))

        for tx, c := range str[n:] {
            if tx > width {
                break
            }

            termbox.SetChar(tx + x, ty + y, c)
        }
    }
}

func (b *Buffer) Reset() {
    b.dx = 0
    b.dy = 0
}

func (b *Buffer) GoToBegin() {
    b.dy = 0
}

func (b *Buffer) GoToEnd() {
    b.dy = min(data - page, data)
}

func (b *Buffer) PageUp() {
    b.dy = max(b.dy - page, 0)
}

func (b *Buffer) PageDown() {
    b.dy = min(b.dy + page, data - page, data)
}

func (b *Buffer) ScrollUp(delta int) {
    b.dy = max(b.dy - delta, 0)
}

func (b *Buffer) ScrollDown(delta int) {
    b.dy = min(b.dy + delta, data - page, data)
}

func (b *Buffer) ScrollLeft(delta int) {
    b.dx = max(b.dx - delta, 0)
}

func (b *Buffer) ScrollRight(delta int) { //
    b.dx = min(b.dx + delta, width)
}
