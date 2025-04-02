# cu [see:you]
The Swiss Army Knife for viewing text or binary files.

## Usage
```console
$ cu [-xhv] [-f FILTER] PATH ...
```

Available options:
* `-f` Filter
* `-x` Hex mode
* `-h` Show usage
* `-v` Show version

## Files
* `~/.cu_history` Input history

## Keyboard
| Shortcut                                     | Action                |
| -------------------------------------------- | --------------------- |
| <kbd>Esc<kbd> / <kbd>Ctrl</kbd>+<kbd>q</kbd> | Exit                  |
| <kbd>F1</kbd> / <kbd>Ctrl</kbd>+<kbd>n</kbd> | Normal mode           |
| <kbd>F2</kbd> / <kbd>Ctrl</kbd>+<kbd>x</kbd> | Hex mode              |
| <kbd>F3</kbd> / <kbd>Ctrl</kbd>+<kbd>t</kbd> | Shell mode            |
| <kbd>Enter</kbd>                             | Apply filter          |
| <kbd>Backspace</kbd>                         | Erase filter          |
| <kbd>Tab</kbd>                               | Load next file        |
| <kbd>Shift</kbd>+<kbd>Tab</kbd>              | Load prev file        |
| <kbd>Ctrl</kbd>+<kbd>r</kbd>                 | Reload file           |
| <kbd>Ctrl</kbd>+<kbd>h</kbd>                 | Show file hash        |
| <kbd>Ctrl</kbd>+<kbd>l</kbd>                 | Toggle line numbers   |
| <kbd>Ctrl</kbd>+<kbd>w</kbd>                 | Toggle line wrap      |
| <kbd>Ctrl</kbd>+<kbd>s</kbd>                 | Save buffer content   |
| <kbd>Ctrl</kbd>+<kbd>c</kbd>                 | Copy buffer content   |
| <kbd>Ctrl</kbd>+<kbd>v</kbd>                 | Paste input           |
| <kbd>Alt</kbd>+<kbd>Up</kbd>                 | Prev input in history |
| <kbd>Alt</kbd>+<kbd>Down</kbd>               | Next input in history |
| <kbd>Shift</kbd>+<kbd>Up</kbd>               | Scroll page up        |
| <kbd>Shift</kbd>+<kbd>Down</kbd>             | Scroll page down      |
| <kbd>Shift</kbd>+<kbd>Left</kbd>             | Scroll page left      |
| <kbd>Shift</kbd>+<kbd>Right</kbd>            | Scroll page right     |

## License
Released under the [MIT License](LICENSE).
