package types

const (
	CRC32IEEE = "crc32-ieee"
	CRC64ECMA = "crc64-ecma"
	CRC64ISO  = "crc64-iso"
	SDHASH    = "sdhash"
	SSDEEP    = "ssdeep"
	TLSH      = "tlsh"
	MD5       = "md5"
	SHA1      = "sha1"
	SHA256    = "sha256"
	SHA3      = "sha3"
	SHA3224   = "sha3-224"
	SHA3256   = "sha3-256"
	SHA3384   = "sha3-384"
	SHA3512   = "sha3-512"
)

type Heap int

const (
	Regular Heap = iota
	Deflate
	Ignore
	Stdin
	Stdout
	Stderr
	Plugin
	Agent
)

type Invoke int

const (
	None Invoke = iota
	Counts
	Entropy
	Strings
	Hash
)
