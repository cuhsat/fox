package types

type Print int

const (
	File Print = iota
	Grep
	Hex
	Hash
	Stats
	Strings
)

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
