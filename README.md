![](assets/logo.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the most useful functionalities from _zcat_, _zless_, _grep_, _hexdump_, _head_, _tail_, _jq_, _wc_ and various hashes into one performant standalone binary. As a forensic tool, <ins>no write actions</ins> are guaranteed to be made.

![](assets/grep.png "Forensic Examiner")

## Usage
```
fx [-p] [-h | -t] [-n # | -c #] [-x | -e PATTERN] [-j | -J] [-o FILE] [PATH ... | -]
```

Special options:
* `-p` Print raw
* `-h` Limit head
* `-t` Limit tail
* `-n` Lines count
* `-c` Bytes count
* `-x` Hexdump mode
* `-e` Pattern value
* `-j` JSON output
* `-J` JSONL output
* `-o` Evidence file

Default options:
* `--help` Usage information
* `--version` Version number

### Examples
Read all files in the current directory:
```console
fx
```

Read directly from stdin:
```console
fx -
```

Read `Zip` and `Tar` compressed files:
```console
fx foo.zip bar.tar.gz
```

Read all `.jsonl` files in all subdirectories:
```console
fx ./**/*.jsonl
```

Extract all users with [Dissect](https://docs.dissect.tools) from `foo.dd`:
```console
target-query -f users foo.dd | rdump -J | fx -
```

Write all lines with `John Doe` from all files to stdout:
```console
fx -p -e "John Doe"
```

Write the first `3` lines of `foo` to `bar`:
```console
fx -hn 3 foo > bar
```

Write the last `8` bytes of `foo` to `bar` in hex:
```console
fx -tc 8 -x foo > bar
```

## Install
```console
make install
```

## Keymap
| Shortcut                                             | Action                |
| ---------------------------------------------------- | --------------------- |
| <kbd>Esc</kbd>                                       | Exit                  |
| <kbd>F1</kbd> / <kbd>Ctrl</kbd> + <kbd>l</kbd>       | Less mode             |
| <kbd>F2</kbd> / <kbd>Ctrl</kbd> + <kbd>g</kbd>       | Grep mode             |
| <kbd>F3</kbd> / <kbd>Ctrl</kbd> + <kbd>x</kbd>       | Hex mode              |
| <kbd>F4</kbd> / <kbd>Ctrl</kbd> + <kbd>Space</kbd>   | Goto mode             |
| <kbd>F5</kbd> / <kbd>Ctrl</kbd> + <kbd>s</kbd>       | Save evidence         |
| <kbd>F6</kbd> / <kbd>Ctrl</kbd> + <kbd>c</kbd>       | Copy to clipboard     |
| <kbd>F7</kbd>                                        | Show counts           |
| <kbd>F8</kbd> / <kbd>Ctrl</kbd> + <kbd>d</kbd>       | Open debug log        |
| <kbd>F9</kbd> / <kbd>Ctrl</kbd> + <kbd>e</kbd>       | Open evidence         |
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
| <kbd>Ctrl</kbd> + <kbd>w</kbd>                       | Toggle text wrap      |

### F1 - Less Mode
| Shortcut                                             | Action                |
| ---------------------------------------------------- | --------------------- |
| <kbd>Space</kbd>                                     | Scroll page down      |

### F2 - Grep Mode
| Shortcut                                             | Action                |
| ---------------------------------------------------- | --------------------- |
| <kbd>Enter</kbd>                                     | Append filter         |
| <kbd>Backspace</kbd>                                 | Delete filter         |
| <kbd>Alt</kbd> + <kbd>Up</kbd>                       | Prev input in history |
| <kbd>Alt</kbd> + <kbd>Down</kbd>                     | Next input in history |
| <kbd>Ctrl</kbd> + <kbd>v</kbd>                       | Paste input           |
| <kbd>Any Key</kbd>                                   | Filter content        |

### F3 - Hex Mode
| Shortcut                                             | Action                |
| ---------------------------------------------------- | --------------------- |
| <kbd>Space</kbd>                                     | Scroll page down      |

### F4 - Goto Mode
| Shortcut                                             | Action                |
| ---------------------------------------------------- | --------------------- |
| <kbd>Enter</kbd>                                     | Goto line             |
| <kbd>Alt</kbd> + <kbd>Up</kbd>                       | Prev input in history |
| <kbd>Alt</kbd> + <kbd>Down</kbd>                     | Next input in history |
| <kbd>Ctrl</kbd> + <kbd>v</kbd>                       | Paste input           |
| <kbd>Any Key</kbd>                                   | Line number           |

## Evidence
> JSON specification under [spec.json](api/spec.json).

## Config
> Located under `~/.fxrc`.

```toml
Theme = "Default"
Follow = false
Line = false
Wrap = false
```

### Environment
```bash
FX_THEME=Default
```

### Themes
* `Default`
* `Monokai`
* `Catppuccin-Latte`
* `Catppuccin-Frappe`
* `Catppuccin-Macchiato`
* `Catppuccin-Mocha`
* `VSCode-Light`
* `VSCode-Dark`
* `Darcula`
* `Nord`
* `Ansi16`
* `Matrix`
* `Monochrome`

---
Made with ❤ in Go