package types

type Heap int

const (
	Regular Heap = iota
	Stdin
	Stdout
	Stderr
	Deflate
)

type Export int

const (
	Text Export = iota
	Json
	Jsonl
)
