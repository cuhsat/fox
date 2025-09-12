# fox deflate
Deflate compressed files.

## Usage
```console
fox deflate [FLAG ...] PATH...
```

### Aliases
`extract`, `unzip`, `de`

### Arguments
Path(s) to open

### Global
- `--no-file` — don't print filenames

### Deflate
- `-l`, `--list` — don't deflate, only list files
- `-d`, `--dir[=PATH]` — deflate into directory (*default:* `.`)
- `-P`, `--pwd=PASSWORD` — password for decryption (only **RAR**, **ZIP**)

## Example
```console
$ fox deflate --pwd=infected ioc.zip
```
