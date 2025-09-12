# fox hash
Display file hash or checksums.

## Usage
```console
fox hash [FLAG ...] PATH ...
```

### Aliases
`sum`, `ha`

### Arguments
Path(s) to open

### Global
- `-p`, `--print` — print directly to console

### Hash
- `-t`, `--type=ALGORITHM` — hash algorithm (*default:* `SHA256`)

Cryptographic hash algorithms:
> MD5, SHA1, SHA256, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

Similarity hash algorithms:
> SDHASH, SSDEEP, TLSH

Checksum algorithms:
> CRC32-IEEE, CRC64-ECMA, CRC64-ISO

## Example
```console
$ fox hash -t=tlsh artifacts.zip
```
