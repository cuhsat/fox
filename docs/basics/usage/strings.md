# fox strings
Display ASCII and Unicode strings.

## Usage
```console
fox strings [FLAG ...] PATH ...
```

### Aliases
`carve`, `st`

### Arguments
Path(s) to open

### Global
- `-p`, `--print` — print directly to console
- `--no-file` — don't print filenames
- `--no-line` — don't print line numbers

### Strings
- `-i`, `--ioc` — detect built-in IoCs
- `-e`, `--regexp=PATTERN` — search for pattern
- `-n`, `--min=NUMBER` — minimum length (*default:* `3`)
- `-m`, `--max=NUMBER` — maximum length (*default:* Unlimited)
- `-a`, `--ascii` — only carve ASCII strings

Built-in IoC patterns:
> UUID, IPv4, IPv6, MAC, URL, Mail

## Example
```console
$ fox strings -in=8 malware.exe
```
