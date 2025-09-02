# Plugins

## Dissect Framework

### target-info
```toml
[auto.info]
name = 'target-info'
path = '.*\.(dd|img|raw|ad1|asdf|E0?|00?)'
exec = [
  'target-info "FILE"'
]
```

### target-query
```toml
[hotkey.f8]
name = 'target-query'
mode = 'query'
exec = [
  'target-query -j -f "INPUT" "BASE"'
]
```

### target-shell
```toml
[hotkey.f9]
name = 'target-shell'
mode = 'shell'
exec = [
  'target-shell -c="INPUT" "BASE"'
]
```

## Eric Zimmerman's tools

### JLECmd
```toml
[auto.jle]
name = 'JLECmd'
path = '.*\.(automatic|custom)Destination-ms'
exec = [
  'dotnet JLECmd.dll -f "FILE" --json "TEMP"'
]
```

### LECmd
```toml
[auto.le]
name = 'LECmd'
path = '.*\.lnk'
exec = [
  'dotnet LECmd.dll -f "FILE" --json "TEMP"'
]
```

### MFTECmd
```toml
[auto.mfte]
name = 'MFTECmd'
path = '\(Boot|LogFile|J|MFT|SDS)'
exec = [
  'dotnet MFTECmd.dll -f "FILE" --json "TEMP"'
]
```

### PECmd
```toml
[auto.pe]
name = 'PECmd'
path = '.*\.pf'
exec = [
  'dotnet PECmd.dll -f "FILE" --json "TEMP"'
]
```

### RBCmd
```toml
[auto.rb]
name = 'RBCmd'
path = '(INFO2|\[0-9A-Z]{7}(\..+)?)'
exec = [
  'dotnet RBCmd.dll -f "FILE" --csv "TEMP"'
]
```

### RECmd
```toml
[auto.re]
name = 'RECmd'
path = '.*\.dat'
exec = [
  'dotnet RECmd.dll -f "FILE" --json "TEMP"'
]
```

### SQLECmd
```toml
[auto.sqle]
name = 'SQLECmd'
path = '.*\.db'
exec = [
  'dotnet SQLECmd.dll -f "FILE" --json "TEMP"'
]
```

### SrumECmd
```toml
[auto.srume]
name = 'SrumECmd'
path = 'SRUDB.dat'
exec = [
  'dotnet SrumECmd.dll -f "FILE" --csv "TEMP"'
]
```

### WxTCmd
```toml
[auto.wxt]
name = 'WxTCmd'
path = '.*\ActivitiesCache.db'
exec = [
  'dotnet WxTCmd.dll -f "FILE" --csv "TEMP"'
]
```

## Forensic Artifacts Collecting Toolkit

### pipeline
```toml
[auto.fact]
name = 'fact'
path = '.*\.(dd|img|raw)'
exec = [
  'sudo fmount "FILE" | ffind | flog -D TEMP'
]
```

## Reverse Engineering

### capa
```toml
[auto.capa]
name = 'capa'
path = '.*\.(bin|dll|exe|scr|sys)'
exec = [
  'capa "FILE"'
]
```

### objdump
```toml
[auto.obj]
name = 'objdump'
path = '.*\.(bin|dll|exe|scr|sys)'
exec = [
  'objdump --disassemble "FILE"'
]
```
