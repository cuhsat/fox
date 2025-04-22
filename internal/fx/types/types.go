package types

const (
    Regular Heap = iota
    Stdin
    Stdout
    Stderr
    Deflate
)

type Heap int

type Format func(string) []string

type FileEntry struct {
    Path string
    Name string
}
