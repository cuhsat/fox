![](assets/logo.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the most useful functionalities from `(z)cat`, `(z)less`, `head`, `tail`, `grep`, `hexdump`, `sha256sum`, `sha1sum`, `md5sum`, `wc` and `jq` into one performant standalone binary.

![](assets/grep.png "Grep Mode")
![](assets/hex.png "Hex Mode")

> As this is a forensic tool, no write actions will be performed.

# Install
```console
$ make install
```

# Usage
```
$ fx [-r] [-h | -t] [-n # | -c #] [-x | -e PATTERN] [-o FILE] [PATH ... | -]
```

Special options:
* `-r` Raw output
* `-h` Head limit
* `-t` Tail limit
* `-n` Lines count
* `-c` Bytes count
* `-x` Hexdump mode
* `-e` Pattern value
* `-o` Evidence file

Default options:
* `--help` Usage information
* `--version` Version number

## Examples
Reading all files in the current directory:
```console
$ fx
```

Reading directly from stdin:
```console
$ fx -
```

Reading `gzip` compressed files:
```console
$ fx foo.gz bar.gz
```

Reading all `.jsonl` files in all subdirectories:
```console
$ fx ./*/*.jsonl
```

Writing all lines with `John Doe` from all files to stdout:
```console
$ fx -r -e "John Doe"
```

Writing the first `3` lines of `foo` to `bar`:
```console
$ fx -h -n 3 foo > bar
```

Writing the last `8` bytes of `foo` to `bar` in hex:
```console
$ fx -t -c 8 -x foo > bar
```

# Keyboard

## General
| Shortcut                                             | Action                |
| ---------------------------------------------------- | --------------------- |
| <kbd>Esc</kbd>                                       | Exit                  |
| <kbd>F1</kbd> / <kbd>Ctrl</kbd> + <kbd>l</kbd>       | Less mode             |
| <kbd>F2</kbd> / <kbd>Ctrl</kbd> + <kbd>g</kbd>       | Grep mode             |
| <kbd>F3</kbd> / <kbd>Ctrl</kbd> + <kbd>x</kbd>       | Hex mode              |
| <kbd>F4</kbd> / <kbd>Ctrl</kbd> + <kbd>Space</kbd>   | Goto mode             |
| <kbd>F9</kbd>                                        | Show file(s) counts   |
| <kbd>F10</kbd>                                       | Show file(s) MD5      |
| <kbd>F11</kbd>                                       | Show file(s) SHA1     |
| <kbd>F12</kbd>                                       | Show file(s) SHA256   |
| <kbd>Tab</kbd>                                       | Load next file        |
| <kbd>Shift</kbd> + <kbd>Tab</kbd>                    | Load prev file        |
| <kbd>Shift</kbd> + <kbd>Up</kbd>                     | Scroll page up        |
| <kbd>Shift</kbd> + <kbd>Down</kbd>                   | Scroll page down      |
| <kbd>Shift</kbd> + <kbd>Left</kbd>                   | Scroll page left      |
| <kbd>Shift</kbd> + <kbd>Right</kbd>                  | Scroll page right     |
| <kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Up</kbd>   | Scroll to start       |
| <kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>Down</kbd> | Scroll to end         |
| <kbd>Ctrl</kbd> + <kbd>r</kbd>                       | Reload file           |
| <kbd>Ctrl</kbd> + <kbd>q</kbd>                       | Close file            |
| <kbd>Ctrl</kbd> + <kbd>t</kbd>                       | Cycle themes          |
| <kbd>Ctrl</kbd> + <kbd>f</kbd>                       | Toggle file follow    |
| <kbd>Ctrl</kbd> + <kbd>n</kbd>                       | Toggle line numbers   |
| <kbd>Ctrl</kbd> + <kbd>w</kbd>                       | Toggle wrap text      |
| <kbd>Ctrl</kbd> + <kbd>s</kbd>                       | Save as evidence      |
| <kbd>Ctrl</kbd> + <kbd>c</kbd>                       | Copy to clipboard     |

## Less Mode
| Shortcut                                             | Action                |
| ---------------------------------------------------- | --------------------- |
| <kbd>Space</kbd>                                     | Scroll page down      |

## Grep Mode
| Shortcut                                             | Action                |
| ---------------------------------------------------- | --------------------- |
| <kbd>Enter</kbd>                                     | Append filter         |
| <kbd>Backspace</kbd>                                 | Delete filter         |
| <kbd>Alt</kbd> + <kbd>Up</kbd>                       | Prev input in history |
| <kbd>Alt</kbd> + <kbd>Down</kbd>                     | Next input in history |
| <kbd>Ctrl</kbd> + <kbd>v</kbd>                       | Paste as input        |
| <kbd>Any Key</kbd>                                   | Filter content        |

# Config
> Located under `~/.fxrc`.

```toml
Theme = "Monokai"
Follow = true  # Follow file
Line = true    # Line numbers
Wrap = true    # Wrap text
```

## Environment
```bash
FX_THEME=Monokai
```

## Themes
* `Monokai`
* `Catppuccin-Latte`
* `Catppuccin-Frappe`
* `Catppuccin-Macchiato`
* `Catppuccin-Mocha`
* `Ansi`
* `Matrix`
* `Monochrome`

---
Made with ❤️ in Go