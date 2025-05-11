![](assets/fox.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the power of **cat**, **less**, **grep**, **hexdump**, **head**, **tail**, **wc**, and various decompression and hash algorithms into one performant standalone binary. As a forensic tool, no write access is made. For more information type `fox --help`.

![](assets/demo.png "Demo")

# Install

```console
go install github.com/cuhsat/fox@latest
```

# Features
* Memory mapped lazy loading in support of big files
* Multi core data handling for fast response time
* Canonical hex view of binary files
* Regular expression filtering
* Unicode multi-byte support
* Unicode bi-directional character filtering
* Build-in decompression of: bzip2, gzip, tar, zip
* Build-in cryptographic hashes: md5, sha1, sha256, sha3
* Build-in AI support via Ollama
* Build-in file statisics support
* Build-in input history
* Evidence bag formats: text, xml, json, jsonl, markdown
* Evidence bag cryptographic signing
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
