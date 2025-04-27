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
fx -p -e "John Doe"
```

Print the first `512` bytes to `mbr` in hex

```console
fx -x -h -c 512 nist.dd > mbr
```
