# fox compare
Compare two files.

## Usage
```console
fox compare [FLAG ...] FILE1 FILE2
```

### Aliases
`diff`, `cmp`, `ce`

### Arguments
Files to open

### Global
- `-p`, `--print` — print directly to console
- `--no-file` — don't print filenames
- `--no-line` — don't print line numbers

## Compare
- `-g`, `--git` — use the unified git diff format

## Example
```console
$ fox compare server.log mirror.log
```
