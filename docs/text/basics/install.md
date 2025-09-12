# Install
[Releases](https://github.com/cuhsat/fox/releases) are available for both `AMD64` and `ARM64` architectures, signed SBOMs included.

## From source
```console
go install github.com/cuhsat/fox@latest
```

## From binaries
> No installation is required, as the binaries are standalone native executables.

Binaries are available for:
 
- Linux
- macOS
- Windows

## From packages
Packages are available in:

- `.apk` format, use `sudo apk add --allow-untrusted fox*.apk`
- `.deb` format, use `sudo dpkg -i fox*.deb`
- `.rpm` format, use  `sudo rpm -i fox*.rpm`

## Dependencies
If you wish to use the LLM capabilities, connection to a local or remote running [Ollama](https://ollama.com) server is required.
