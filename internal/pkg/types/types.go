package types

type Run int
type Heap int

const (
	MD5     = "md5"
	SHA1    = "sha1"
	SHA256  = "sha256"
	SHA3    = "sha3"
	SHA3224 = "sha3-224"
	SHA3256 = "sha3-256"
	SHA3384 = "sha3-384"
	SHA3512 = "sha3-512"
	SDHASH  = "sdhash"
	SSDEEP  = "ssdeep"
	TLSH    = "tlsh"
)

const (
	File Run = iota
	Grep
	Hex
	Hash
	Stats
	Strings
)

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
