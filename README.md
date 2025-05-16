![](assets/logo.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the power of **cat**, **less**, **grep**, **hexdump**, **head**, **tail**, **wc**, and various decompression and hash algorithms into one performant standalone binary. As a forensic tool, no write access is made. For more information simply type `fox --help`.

![](assets/demo.png "Demo")

# Install

```console
go install -tags=ui,ai github.com/cuhsat/fox@latest
```

# Features
* Standalone native binary for AMD64 and ARM64
* Binary customizable via feature flags
* Memory mapped lazy loading in support of big files
* Multi core data handling for fast response times
* Canonical hex view of binary files
* Regular expression filtering
* Unicode multi-byte support
* Unicode bidirectional character filtering (CVE-2021-42574)
* Build-in decompression of: bzip2, gzip, tar, xz, zip, zlib, zstd
* Build-in cryptographic hashes: MD5, SHA1, SHA256, SHA3, SHA3-XXX
* Build-in AI with RAG support via LangChain and Ollama
* Build-in file statistics support
* Build-in input and output history
* Configurable plugin support for tools like Dissect
* Evidence bag formats: Text, Markdown, JSON, JSONL, XML, SQLite
* Evidence bag cryptographic signing via HMAC-SHA256
* Graphical user interface
* Suspend to shell support
* Mouse scrolling support
* Full clipboard support
* Theme support including:
  *  Examiner-Light
  *  Examiner-Dark
  *  Catppuccin-Latte
  *  Catppuccin-Frappe
  *  Catppuccin-Macchiato
  *  Catppuccin-Mocha
  *  Solarized-Light
  *  Solarized-Dark
  *  VSCode-Light
  *  VSCode-Dark
  *  Monokai
  *  Darcula
  *  Nord
  *  Corporate
  *  Matrix
  *  Ansi16
  *  Monochrome
