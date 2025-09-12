# fox entropy
Display file entropy.

## Usage
```console
fox entropy [FLAG ...] PATH ...
```

### Aliases
`en`

### Arguments
Path(s) to open

### Global
- `-p`, `--print` — print directly to console

### Entropy
- `-n`, `--min[=DECIMAL]` — minimum score (*default:* `0.8`)
- `-m`, `--max[=DECIMAL]` — maximum score (*default:* `0.8`)

## Example
```console
$ fox entropy -n ./**/*
```
