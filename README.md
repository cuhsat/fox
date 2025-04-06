# cu [see:you]
The Swiss Army Knife for viewing lined based text files. Combining `cat`, `zcat`, `grep`, `hexdump`, `head`, `tail`, `wc` and `sha256` into one performance oriented tool.

## Usage
```
$ cu [-x] [-h # | -t #] [-f FILTER] [PATH ...]
```

Available options:
* `-x` Hex mode
* `-h` Head count
* `-t` Tail count
* `-f` Filter value

## Files
* `~/.cu_history` Input history
* `~/.cin` Input buffer
* `~/.cout` Output buffer

## Keyboard
| Shortcut                                           | Action                |
| -------------------------------------------------- | --------------------- |
| <kbd>Esc</kbd> / <kbd>Ctrl</kbd> + <kbd>q</kbd>    | Exit                  |
| <kbd>F1</kbd> / <kbd>Ctrl</kbd> + <kbd>g</kbd>     | Grep mode             |
| <kbd>F2</kbd> / <kbd>Ctrl</kbd> + <kbd>x</kbd>     | Hex mode              |
| <kbd>F3</kbd> / <kbd>Ctrl</kbd> + <kbd>Space</kbd> | Goto mode             |
| <kbd>Enter</kbd>                                   | Append filter         |
| <kbd>Backspace</kbd>                               | Delete filter         |
| <kbd>Tab</kbd>                                     | Load next file        |
| <kbd>Shift</kbd> + <kbd>Tab</kbd>                  | Load prev file        |
| <kbd>Shift</kbd> + <kbd>Up</kbd>                   | Scroll page up        |
| <kbd>Shift</kbd> + <kbd>Down</kbd>                 | Scroll page down      |
| <kbd>Shift</kbd> + <kbd>Left</kbd>                 | Scroll page left      |
| <kbd>Shift</kbd> + <kbd>Right</kbd>                | Scroll page right     |
| <kbd>Alt</kbd> + <kbd>Up</kbd>                     | Prev input in history |
| <kbd>Alt</kbd> + <kbd>Down</kbd>                   | Next input in history |
| <kbd>Ctrl</kbd> + <kbd>r</kbd>                     | Reload file           |
| <kbd>Ctrl</kbd> + <kbd>h</kbd>                     | Show file hash        |
| <kbd>Ctrl</kbd> + <kbd>n</kbd>                     | Toggle line numbers   |
| <kbd>Ctrl</kbd> + <kbd>w</kbd>                     | Toggle line wrap      |
| <kbd>Ctrl</kbd> + <kbd>s</kbd>                     | Save buffer content   |
| <kbd>Ctrl</kbd> + <kbd>c</kbd>                     | Copy buffer content   |
| <kbd>Ctrl</kbd> + <kbd>v</kbd>                     | Paste as input        |

## License
Released under the [MIT License](LICENSE).
