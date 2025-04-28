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

Print all lines with `John Doe` from all files

```console
fx -pe "John Doe"
```

Print the first `512` bytes from `nist.dd` to `mbr` in hex

```console
fx -xhc 512 nist.dd > mbr
```
