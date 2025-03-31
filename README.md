# cu
[![Go Report Card](https://goreportcard.com/badge/github.com/cuhsat/cu?style=flat-square)](https://goreportcard.com/report/github.com/cuhsat/cu)
[see:you] is a format agnostic, performance optimized text file viewer.

![](assets/cu.gif)

## Usage
```console
$ cu [-xhv] PATH ...
```

Available options:
* `-x` Hex mode
* `-h` Show usage
* `-v` Show version

## Keyboard
| Shortcut                        | Action              |
| ------------------------------- | ------------------- |
| <kbd>Enter</kbd>                | Apply filter        |
| <kbd>Backspace</kbd>            | Delete filter       |
| <kbd>Tab</kbd>                  | Load next file      |
| <kbd>Shift</kbd>+<kbd>Tab</kbd> | Load previous file  |
| <kbd>Ctrl</kbd>+<kbd>r</kbd>    | Reload actual file  |
| <kbd>Ctrl</kbd>+<kbd>s</kbd>    | Save filtered lines |
| <kbd>Ctrl</kbd>+<kbd>c</kbd>    | Copy filtered lines |
| <kbd>Ctrl</kbd>+<kbd>v</kbd>    | Paste filter input  |
| <kbd>Ctrl</kbd>+<kbd>l</kbd>    | Toggle line numbers |
| <kbd>Ctrl</kbd>+<kbd>w</kbd>    | Toggle line wrap    |
| <kbd>Ctrl</kbd>+<kbd>x</kbd>    | Toggle hex mode     |
| <kbd>Ctrl</kbd>+<kbd>y</kbd>    | Show file hash      |
| <kbd>Ctrl</kbd>+<kbd>q</kbd>    | Exit                |

## License
Released under the [MIT License](LICENSE).
