package types

type Heap int

const (
	Regular Heap = iota
	Deflate
	Stdin
	Stdout
	Stderr
	Plugin
	Prompt
)

type Output int

const (
	File Output = iota
	Grep
	Hex
	Hash
	Reverse
	Stats
	Strings
)
