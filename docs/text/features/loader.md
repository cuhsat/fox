# Loader
The file loader applies path [matching and globbing](https://github.com/bmatcuk/doublestar?tab=readme-ov-file#patterns) before opening any files.

## Deflate
The file loader is also capable of reading nested archives/compressed files in many different formats. This is done by parsing the file for **magic numbers**, like the `file` utility does.

> To prevent automatic deflation, use either the `-r` or `--no-deflate` flag.

## Convert
The file loader also detects and converts [Windows Event Logs](files/evtx.md) and [Linux Systemd Journals](files/journal.md) to [JSON Lines](https://jsonlines.org/) upon loading. 

> To prevent automatic conversion, use either the `-r` or `--no-convert` flag.

## Plugins
If the files path matches a `path` defined by an `auto` [plugin](../plugins.md), the plugin is automatically executed and its **STDOUT** output is shown instead of the files contents. 

> To prevent automatic plugins, use either the `-r` or `--no-plugins` flag.

## Archive Formats
Built-in extraction of:

| Format | Specification                                                                           |
|--------|-----------------------------------------------------------------------------------------|
| `cab`  | [Microsoft Cabinet Format](https://msdn.microsoft.com/en-us/library/bb267310.aspx)      |
| `rar`  | [RAR Archive Format](https://www.rarlab.com/technote.htm)                               |
| `tar`  | [Basic Tar Format](https://www.gnu.org/software/tar/manual/html_node/Standard.html)     |
| `zip`  | [ZIP File Format](https://www.loc.gov/preservation/digital/formats/fdd/fdd000354.shtml) |

## Compression Formats
Built-in deflation of:

| Format   | Specification                                                       |
|----------|---------------------------------------------------------------------|
| `brotli` | [Brotli Compression Format](https://github.com/google/brotli)       |
| `bzip2`  | [bzip2 Data Compressor](https://sourceware.org/bzip2/)              |
| `gzip`   | [Gzip Data Compression](https://www.gnu.org/software/gzip/)         |
| `lz4`    | [LZ4 Compression](https://github.com/lz4/lz4)                       |
| `xz`     | [xz File Format](https://tukaani.org/xz/format.html)                |
| `zlib`   | [zlib Compression](https://zlib.net/)                               |
| `zstd`   | [Zstandard Compression Algorithm](https://github.com/facebook/zstd) |

## Examples
Load all files with `log` extension:
```console
$ fox *.log
```

Load all files with `evtx` extension, in all subfolders:
```console
$ fox ./**/*.evtx
```

Load all files in all subfolders:
```console
$ fox ./**/*
```
