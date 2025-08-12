# Dissect

## target-info
```toml
[Autostart.Info]
Name = "target-info"
Pattern = '.*\.(dd|img|raw|ad1|asdf|E0?|00?)'
Commands = [
  'target-info "{{file}}"',
]
```

## target-query
```toml
[Hotkey.F7]
Name = "target-query"
Prompt = "Query"
Commands = [
  'target-query -j -f "{{value}}" "{{base}}"',
]
```

## target-shell
```toml
[Hotkey.F8]
Name = "target-shell"
Prompt = "Shell"
Commands = [
  'target-shell -c="{{value}}" "{{base}}"',
]
```

# Eric Zimmerman's tools

## JLECmd
```toml
[Autostart.Jumplist]
Name = "JLECmd"
Pattern = '.*\.(automatic|custom)Destination-ms'
Commands = [
  'dotnet JLECmd.dll -f "{{file}}" --json "{{dir}}"',
]
```

## LECmd
```toml
[Autostart.Link]
Name = "LECmd"
Pattern = '.*\.lnk'
Commands = [
  'dotnet LECmd.dll -f "{{file}}" --json "{{dir}}"',
]
```

## MFTECmd
```toml
[Autostart.Filesystem]
Name = "MFTECmd"
Pattern = '\$(Boot|LogFile|J|MFT|SDS)'
Commands = [
  'dotnet MFTECmd.dll -f "{{file}}" --json "{{dir}}"',
]
```

## PECmd
```toml
[Autostart.Prefetch]
Name = "PECmd"
Pattern = '.*\.pf'
Commands = [
  'dotnet PECmd.dll -f "{{file}}" --json "{{dir}}"',
]
```

## RBCmd
```toml
[Autostart.Trash]
Name = "RBCmd"
Pattern = '(INFO2|\$[0-9A-Z]{7}(\..+)?)$'
Commands = [
  'dotnet RBCmd.dll -f "{{file}}" --csv "{{dir}}"',
]
```

## RECmd
```toml
[Autostart.Registry]
Name = "RECmd"
Pattern = '.*\.dat'
Commands = [
  'dotnet RECmd.dll -f "{{file}}" --json "{{dir}}"',
]
```

## SQLECmd
```toml
[Autostart.Database]
Name = "SQLECmd"
Pattern = '.*\.db'
Commands = [
  'dotnet SQLECmd.dll -f "{{file}}" --json "{{dir}}"',
]
```

## SrumECmd
```toml
[Autostart.Energy]
Name = "SrumECmd"
Pattern = 'SRUDB.dat'
Commands = [
  'dotnet SrumECmd.dll -f "{{file}}" --csv "{{dir}}"',
]
```

## WxTCmd
```toml
[Autostart.Timeline]
Name = "WxTCmd"
Pattern = '.*\ActivitiesCache.db'
Commands = [
  'dotnet WxTCmd.dll -f "{{file}}" --csv "{{dir}}"',
]
```

# Reverse Engineering

## capa
```toml
[Autostart.Capa]
Name = "capa"
Pattern = '.*\.(bin|dll|exe|scr|sys)'
Commands = [
  'capa "{{file}}"',
]
```

## objdump
```toml
[Autostart.Dump]
Name = "objdump"
Pattern = '.*\.(bin|dll|exe|scr|sys)'
Commands = [
  'objdump --disassemble "{{file}}"',
]
```
