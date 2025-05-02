# Usage Examples

Examine the current dir

```console
fx
```

Examine directly from stdin

```console
fx -
```

Examine all `.jsonl` files in all sub dirs

```console
fx ./**/*.jsonl
```

Print the `MD5` hashsum from all files in `nist.zip`

```console
fx -ps=md5 nist.zip
```

Print all lines containing `John Doe` from all files

```console
fx -pe "John Doe"
```

Print the first `512` bytes from `img.dd` to `mbr` in hex

```console
fx -pxhc=512 img.dd > mbr
```
