![](docs/img/fox.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the power of many traditional tools like **grep**, **hexdump** and **strings** with the possibilities of modern **LLMs**, to leverage your forensic examination process. Standalone native **AMD64** and **ARM64** binaries are available for **Windows**, **Linux** and **macOS**.

❯ Visit [forensic-examiner.eu](https://forensic-examiner.eu) for more information.

![](docs/img/demo.png)

## Features
* In-memory read-only forensic filesystem
* Unicode multibyte support with [bidirectional character](https://nvd.nist.gov/vuln/detail/CVE-2021-42574) filtering
* Built-in regular expression filtering with dynamic context window
* Built-in canonical `hexdump` view of files
* Built-in `wc` like counts with Shannon entropy
* Built-in ASCII and Unicode string carving
* Built-in detection of simple indicators of compromise
* Built-in parsing of Linux Journals and Windows Event Logs
* Built-in formating of CSV, JSON and JSON Lines
* Built-in extraction and deflation of many formats
* Built-in cryptographic and similarity hashes
* Built-in AI agent using local [Ollama LLMs](https://ollama.com/search) like *DeepSeek R1*
* Plugin support for e.g. [Dissect](https://docs.dissect.tools) or [Eric Zimmerman's tools](https://ericzimmerman.github.io/)
* Evidence bag with chain of custody signing
* Evidence streaming using [ECS](https://www.elastic.co/docs/reference/ecs)
or [Splunk HEC](https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTinput)
* Terminal UI compatible with many terminals

## Install
Install directly using Go:
```console
go install github.com/cuhsat/fox@latest
```

## Build
Full-featured version:
```console
go build -o fox main.go
```

Minimal version with AI and UI stripped:
```console
go build -o fox -tags minimal main.go
```

## License
🦊 is released under the [GPL-3.0](LICENSE.md).
