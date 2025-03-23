# cu
[![Go Report Card](https://goreportcard.com/badge/github.com/cuhsat/cu?style=flat-square)](https://goreportcard.com/report/github.com/cuhsat/cu)
[see:you] is a format agnostic, performance optimized text file viewer.

TODO: Gif

## Usage
```sh
$ cu [-hv] PATH ...
```

Available options:
* `-h` Show usage
* `-v` Show version

## Keyboard
| Shortcut                        | Action              |
| ------------------------------- | ------------------- |
| <kbd>Enter</kbd>                | Apply filter        |
| <kbd>Backspace</kbd>            | Delete filter       |
| <kbd>Tab</kbd>                  | Load next file      |
| <kbd>Shift</kbd>+<kbd>Tab</kbd> | Load previous file  |
| <kbd>Ctrl</kbd>+<kbd>c</kbd>    | Copy filtered lines |
| <kbd>Ctrl</kbd>+<kbd>s</kbd>    | Save filtered lines |
| <kbd>Ctrl</kbd>+<kbd>n</kbd>    | Toggle line numbers |
| <kbd>Ctrl</kbd>+<kbd>w</kbd>    | Toggle line wrap    |
| <kbd>Ctrl</kbd>+<kbd>q</kbd>    | Exit                |

## License
Released under the [MIT License](LICENSE).
