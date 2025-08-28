![](assets/fox.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the power of many traditional tools like **grep**, **hexdump** or **strings** with the possibilities of modern **LLMs**, to leverage your forensic examining process.

![](assets/live.png "Live Demo")

```console
go install github.com/cuhsat/fox@latest
```

## Key Features
* Available for [Windows, Linux, macOS](https://github.com/cuhsat/fox/releases)
* Standalone native `AMD64` and `ARM64` binaries
* Multi core data handling for fast response times
* Memory mapped lazy loaded file handling
* In-memory virtual filesystem
* Prohibited file write access
* Path matching and globbing
* Unicode multi-byte support
* [Bidirectional character](https://nvd.nist.gov/vuln/detail/CVE-2021-42574) filtering
* Timestamps are normalized to UTC
* Built-in scroll past end ability
* Built-in timestamped input history
* Built-in regular expression filtering
* Built-in dynamic context window
* Built-in canonical `hexdump` of files
* Built-in `wc` like counts with Shannon entropy
* Built-in `ASCII` and `Unicode` string carving
* Built-in IoC detector for `UUID`, `IPv4`, `IPv6`, `MAC`, `URL`, `Mail`
* Built-in parser for Windows event log `EVTX` files
* Built-in sniffer for `CSV` delimiter formats
* Built-in formating of: `CSV`, `JSON`, `JSONL` data
* Built-in extraction of: `cab`, `rar`, `tar`, `zip`
* Built-in deflation of: `brotli`, `bzip2`, `gzip`, `lz4`, `xz`, `zlib`, `zstd`
* Built-in cryptographic hashes: `MD5`, `SHA1`, `SHA256`, `SHA3`, `SHA3-224`, `SHA3-256`, `SHA3-384`, `SHA3-512`
* Built-in similarity hashes: `SDHASH`, `SSDEEP`, `TLSH`
* Built-in checksums: `CRC32-IEEE`, `CRC64-ECMA`, `CRC64-ISO`
* Built-in in-memory RAG database for document embeddings
* Built-in AI agent using local [Ollama LLMs](https://ollama.com/search) like *Mistral* or *DeepSeek R1*
* [Plugin](PLUGINS.md) support for e.g. the [Dissect](https://docs.dissect.tools) framework or [Eric Zimmerman's tools](https://ericzimmerman.github.io/)
* Evidence bag formats: `raw`, `text`, `JSON`, `JSONL`, `XML`, `SQLite3`
* Evidence bag chain of custody signing via `HMAC-SHA256`
* Evidence streaming to server in [Elastic Common Schema 9.1](https://www.elastic.co/docs/reference/ecs)
* Terminal interface compatible with many terminals
  * Support for copy and bracketed paste
  * Support for mouse scrolling
  * Suspend to shell capability
  * Configurable color [themes](THEMES.md)

## License
🦊 is released under the [GPL-3.0](LICENSE.md).
