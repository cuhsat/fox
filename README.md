![](docs/fox.png "Forensic Examiner")

The Swiss Army Knife for examining text files. Combining the power of many traditional tools like **grep**, **hexdump** and **strings** with the abilities of modern **LLMs**, to leverage your forensic examination process. Standalone native binaries are available for Windows, Linux and macOS.

![](docs/images/terminal.png)

## Key Features
* Read-only in-memory filesystem abstraction
* Multibyte support with [bidirectional character](https://nvd.nist.gov/vuln/detail/CVE-2021-42574) detection
* Built-in `grep`, `hexdump`, `diff` and `strings` like abilities
* Built-in parsing of Linux Journals and Windows Event Logs
* Built-in popular cryptography and similarity hashes
* Deflation and extraction of many archive formats
* Evidence streaming using [Splunk HEC](https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTinput) or [ECS](https://www.elastic.co/docs/reference/ecs)
* Evidence bag with Chain of Custody signing
* Integrated plugin support for [Dissect](https://docs.dissect.tools) or [Eric Zimmerman's tools](https://ericzimmerman.github.io/)
* Integrated agent using [Ollama LLMs](https://ollama.com/search) like *DeepSeek R1*

## Install
Install directly using Go:
```console
go install github.com/cuhsat/fox@latest
```

## Build
Build a full-featured version:
```console
go build -o fox main.go
```

Build a minimal version with stripped AI and UI:
```console
go build -o fox -tags minimal main.go
```

## License
🦊 [Forensic Examiner](https://forensic-examiner.eu) is released under the [GPL-3.0](LICENSE.md).
