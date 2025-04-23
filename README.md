![](assets/logo.png "Forensic Examiner")
The Swiss Army Knife for examining text files. Combining the most useful functionalities from _zcat_, _zless_, _grep_, _hexdump_, _head_, _tail_, _jq_, _wc_ and various cryptographic hashes into one performant standalone binary. As is this a forensic tool and not an editor, it is guaranteed, that <ins>no write actions</ins> will be made.

![](assets/grep.png "Screenshot")

# Install
```console
$ go install github.com/cuhsat/fx@latest
```

# Usage
```
fx [-x] [-p] [-h|t] [-n|c #] [-e PATTERN] [-j] [-J] [-o FILE] [PATH ...]
```
Positional arguments:
* `PATH` to open (default: current dir)

Mode:
* `-x` start in Hex mode

Print:
* `-p` print to console (no UI)

Limits:
* `-h` limit head of file by ...
* `-t` limit tail of file by ...
* `-n #` number of lines
* `-c #` number of bytes

Filters:
* `-e PATTERN` to filter

Evidence:
* `-o FILE` for evidence bag (default: EVIDENCE)
* `-j` output in JSON format
* `-J` output in JSON lines format

Options:
* `--help`    show help message
* `--version` show version info

## Examples
Examine the current dir:
```console
fx
```
Examine directly from stdin:
```console
fx -
```
Examine all `.jsonl` files in all sub dirs:
```console
fx ./**/*.jsonl
```
Print all lines with `John Doe` of all files:
```console
fx -p -e "John Doe"
```
Print the first `512` bytes to `mbr` in hex:
```console
fx -x -hc 512 nist.dd > mbr
```

# Basic Keymap
| Shortcut                           | Action        |
| ---------------------------------- | ------------- |
| <kbd>Esc</kbd>                     | Exit          |
| <kbd>Tab</kbd>                     | Next file     |
| <kbd>Ctrl</kbd> + <kbd>l</kbd>     | Less mode     |
| <kbd>Ctrl</kbd> + <kbd>g</kbd>     | Grep mode     |
| <kbd>Ctrl</kbd> + <kbd>x</kbd>     | Hex mode      |
| <kbd>Ctrl</kbd> + <kbd>Space</kbd> | Goto mode     |
| <kbd>Ctrl</kbd> + <kbd>s</kbd>     | Save evidence |
| <kbd>Enter</kbd>                   | Append filter |
| <kbd>Backspace</kbd>               | Delete filter |

---
Made with ❤ in Go