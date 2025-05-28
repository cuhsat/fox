![](assets/logo.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the power of **cat**, **less**, **grep**, **hexdump**, **head**, **tail**, **strings**, **wc** and various decompression and cryptographic hash algorithms into one performant standalone binary. For more information simply type `fox --help`

```console
go install github.com/cuhsat/fox@latest
```

![](assets/demo.png "Demo")

# Features
* Standalone native binary for AMD64 and ARM64
* Binary size customizable via feature flags
* Multi core data handling for fast response times
* Memory mapped lazy loaded files
* File write access is prohibited
* Unicode multi-byte support
* Unicode bidirectional character filtering (CVE-2021-42574)
* Build-in canonical hex view of files
* Build-in regular expression filtering
* Build-in ASCII and Unicode string carving
* Build-in parser for Windows event log EVTX format
* Build-in decompression of: brotli, bzip2, gzip, tar, xz, zip, zlib, zstd
* Build-in cryptographic hashes: MD5, SHA1, SHA256, SHA3, SHA3-XXX
* Build-in RAG agent provided by an Ollama local LLM
* Build-in timestamped input and agent history
* Build-in wc like file content statistics
* Build-in plugin support for tools like Fox-IT's Dissect
* Evidence bag formats: Text, JSON, JSONL, XML, SQLite
* Evidence bag cryptographic signing via HMAC-SHA256
* Terminal interface supporting a vast amount of terminals
  * With support for copy & bracketed paste
  * With support for mouse scrolling
  * With suspend to shell capability
  * With configureable color themes
  * And many popular themes already build-in:
    * Examiner-Light
    * Examiner-Dark
    * Catppuccin-Latte
    * Catppuccin-Frappe
    * Catppuccin-Macchiato
    * Catppuccin-Mocha
    * Solarized-Light
    * Solarized-Dark
    * VSCode-Light
    * VSCode-Dark
    * Monokai
    * Darcula
    * Nord
    * Corpo
    * Matrix
    * Ansi16
    * Monochrome

![](assets/tail.png "Blue Tail")
