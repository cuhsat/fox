package types

type Heap int

const (
	Regular Heap = iota
	Stdin
	Stdout
	Stderr
	Deflate
)

type Output int

const (
	File Output = iota
	Grep
	Hex
	Hash
	Count
)
