![](docs/img/logo.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the most useful functionalities from _zcat_, _zless_, _grep_, _hexdump_, _head_, _tail_, _wc_ and various cryptographic hash functions into one performant standalone binary.

![](docs/img/grep.png)

```console
go install github.com/cuhsat/fx@latest
```

## Examples
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

## Usage
```
fx [-x] [-p] [-h|t] [-n|c #] [-e PATTERN] [-j] [-J] [-o FILE] [PATH ...]

The Swiss Army Knife for examining text files

positional arguments:
  PATH to open (default: current dir)

mode:
  -x start in Hex mode

print:
  -p print to console (no UI)

limits:
  -h limit head of file by ...
  -t limit tail of file by ...
  -n # number of lines
  -c # number of bytes

filters:
  -e PATTERN to filter

evidence:
  -o FILE for evidence bag (default: EVIDENCE)
  -j output in JSON format
  -J output in JSON lines format

options:
  --help    show help message
  --version show version info
```
Made with ❤ in Go