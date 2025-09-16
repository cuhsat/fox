# Plugins
Located under `~/.config/fox/plugins`.

Two different kind of plugins can be defined: 

- Automatic plugins (`[auto.*]`)
- Hotkey plugins (`[hotkey.*]`)

Hotkey plugins provide the ability to ask the user for a specific `INPUT`, which then can be used as a variable while executing a command.

> To prevent the automatic execution of plugins, use either the `-r` or `--no-plugins` flag.

## Configuration
Available variables:

- `FILE` path
- `BASE` path
- `TEMP` folder (will be created)

Additional variables for hotkeys:

- `INPUT` by user

Available hotkeys:

- <kbd>F8</kbd> to <kbd>F12</kbd>

```toml
[auto.info]
name = 'target-info'
path = '.*\.(dd|img|raw|ad1|asdf|E0?|00?)'
exec = [
  'target-info "FILE"'
]

[hotkey.f8]
name = 'target-query'
mode = 'query'
exec = [
  'target-query -j -f "INPUT" "BASE"'
]
```

## Dissect

### [target-info](https://docs.dissect.tools/en/latest/tools/target-info.html)
```toml
[auto.info]
name = 'target-info'
path = '.*\.(dd|img|raw|ad1|asdf|E0?|00?)'
exec = [
  'target-info "FILE"'
]
```

### [target-query](https://docs.dissect.tools/en/latest/tools/target-query.html)
```toml
[hotkey.f8]
name = 'target-query'
mode = 'query'
exec = [
  'target-query -j -f "INPUT" "BASE"'
]
```

### [target-shell](https://docs.dissect.tools/en/latest/tools/target-shell.html)
```toml
[hotkey.f9]
name = 'target-shell'
mode = 'shell'
exec = [
  'target-shell -c="INPUT" "BASE"'
]
```

## Eric Zimmerman's Tools

### [JLECmd](https://github.com/EricZimmerman/JLECmd)
```toml
[auto.jle]
name = 'JLECmd'
path = '.*\.(automatic|custom)Destination-ms'
exec = [
  'JLECmd.exe -f "FILE" --json "TEMP"'
]
```

### [LECmd](https://github.com/EricZimmerman/LECmd)
```toml
[auto.le]
name = 'LECmd'
path = '.*\.lnk'
exec = [
  'LECmd.exe -f "FILE" --json "TEMP"'
]
```

### [MFTECmd](https://github.com/EricZimmerman/MFTECmd)
```toml
[auto.mfte]
name = 'MFTECmd'
path = '\(Boot|LogFile|J|MFT|SDS)'
exec = [
  'MFTECmd.exe -f "FILE" --json "TEMP"'
]
```

### [PECmd](https://github.com/EricZimmerman/PECmd)
```toml
[auto.pe]
name = 'PECmd'
path = '.*\.pf'
exec = [
  'PECmd.exe -f "FILE" --json "TEMP"'
]
```

### [RBCmd](https://github.com/EricZimmerman/RBCmd)
```toml
[auto.rb]
name = 'RBCmd'
path = '(INFO2|\[0-9A-Z]{7}(\..+)?)'
exec = [
  'RBCmd.exe -f "FILE" --csv "TEMP"'
]
```

### [RECmd](https://github.com/EricZimmerman/RECmd)
```toml
[auto.re]
name = 'RECmd'
path = '.*\.dat'
exec = [
  'RECmd.exe -f "FILE" --json "TEMP"'
]
```

### [SQLECmd](https://github.com/EricZimmerman/SQLECmd)
```toml
[auto.sqle]
name = 'SQLECmd'
path = '.*\.db'
exec = [
  'SQLECmd.exe -f "FILE" --json "TEMP"'
]
```

### [SrumECmd](https://github.com/EricZimmerman/Srum)
```toml
[auto.srume]
name = 'SrumECmd'
path = 'SRUDB.dat'
exec = [
  'SrumECmd.exe -f "FILE" --csv "TEMP"'
]
```

### [WxTCmd](https://github.com/EricZimmerman/WxTCmd)
```toml
[auto.wxt]
name = 'WxTCmd'
path = '.*\ActivitiesCache.db'
exec = [
  'WxTCmd.exe -f "FILE" --csv "TEMP"'
]
```

## Forensics

### [FACT](https://github.com/cuhsat/fact)
```toml
[auto.fact]
name = 'fact'
path = '.*\.(dd|img|raw)'
exec = [
  'sudo fmount "FILE" | ffind | flog -D TEMP'
]
```

## Reverse Engineering

### [Mandiant capa](https://github.com/mandiant/capa)
```toml
[auto.capa]
name = 'capa'
path = '.*\.(bin|dll|exe|scr|sys)'
exec = [
  'capa "FILE"'
]
```

### [objdump](https://linux.die.net/man/1/objdump)
```toml
[auto.obj]
name = 'objdump'
path = '.*\.(bin|dll|exe|scr|sys)'
exec = [
  'objdump --disassemble "FILE"'
]
```
