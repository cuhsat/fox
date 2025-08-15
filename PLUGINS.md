# Dissect

## target-info
```toml
[Autostart.Info]
Name = "target-info"
Pattern = '.*\.(dd|img|raw|ad1|asdf|E0?|00?)'
Commands = [
  'target-info "{{FILE}}"',
]
```

## target-query
```toml
[Hotkey.F7]
Name = "target-query"
Prompt = "Query"
Commands = [
  'target-query -j -f "{{VALUE}}" "{{BASE}}"',
]
```

## target-shell
```toml
[Hotkey.F8]
Name = "target-shell"
Prompt = "Shell"
Commands = [
  'target-shell -c="{{VALUE}}" "{{BASE}}"',
]
```

# Eric Zimmerman's tools

## JLECmd
```toml
[Autostart.Jumplist]
Name = "JLECmd"
Pattern = '.*\.(automatic|custom)Destination-ms'
Commands = [
  'dotnet JLECmd.dll -f "{{FILE}}" --json "{{TEMP}}"',
]
```

## LECmd
```toml
[Autostart.Link]
Name = "LECmd"
Pattern = '.*\.lnk'
Commands = [
  'dotnet LECmd.dll -f "{{FILE}}" --json "{{TEMP}}"',
]
```

## MFTECmd
```toml
[Autostart.Filesystem]
Name = "MFTECmd"
Pattern = '\$(Boot|LogFile|J|MFT|SDS)'
Commands = [
  'dotnet MFTECmd.dll -f "{{FILE}}" --json "{{TEMP}}"',
]
```

## PECmd
```toml
[Autostart.Prefetch]
Name = "PECmd"
Pattern = '.*\.pf'
Commands = [
  'dotnet PECmd.dll -f "{{FILE}}" --json "{{TEMP}}"',
]
```

## RBCmd
```toml
[Autostart.Trash]
Name = "RBCmd"
Pattern = '(INFO2|\$[0-9A-Z]{7}(\..+)?)$'
Commands = [
  'dotnet RBCmd.dll -f "{{FILE}}" --csv "{{TEMP}}"',
]
```

## RECmd
```toml
[Autostart.Registry]
Name = "RECmd"
Pattern = '.*\.dat'
Commands = [
  'dotnet RECmd.dll -f "{{FILE}}" --json "{{TEMP}}"',
]
```

## SQLECmd
```toml
[Autostart.Database]
Name = "SQLECmd"
Pattern = '.*\.db'
Commands = [
  'dotnet SQLECmd.dll -f "{{FILE}}" --json "{{TEMP}}"',
]
```

## SrumECmd
```toml
[Autostart.Energy]
Name = "SrumECmd"
Pattern = 'SRUDB.dat'
Commands = [
  'dotnet SrumECmd.dll -f "{{FILE}}" --csv "{{TEMP}}"',
]
```

## WxTCmd
```toml
[Autostart.Timeline]
Name = "WxTCmd"
Pattern = '.*\ActivitiesCache.db'
Commands = [
  'dotnet WxTCmd.dll -f "{{FILE}}" --csv "{{TEMP}}"',
]
```

# Reverse Engineering

## capa
```toml
[Autostart.Capa]
Name = "capa"
Pattern = '.*\.(bin|dll|exe|scr|sys)'
Commands = [
  'capa "{{FILE}}"',
]
```

## objdump
```toml
[Autostart.Dump]
Name = "objdump"
Pattern = '.*\.(bin|dll|exe|scr|sys)'
Commands = [
  'objdump --disassemble "{{FILE}}"',
]
```
