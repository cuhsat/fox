package types

type Heap int

const (
	Regular Heap = iota
	Deflate
	Ignore
	Stdin
	Stdout
	Stderr
	Prompt
	Plugin
)

type Output int

const (
	File Output = iota
	Grep
	Hex
	Hash
	Stats
	Carve
)
