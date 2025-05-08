package types

type Heap int

const (
	Regular Heap = iota
	Stdin
	Stdout
	Stderr
	Plugin
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
