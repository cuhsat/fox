package types

const (
    Regular Heap = iota
    StdIn
    StdOut
    StdErr
    Deflate
)

type Heap int

type Format func(string) []string
