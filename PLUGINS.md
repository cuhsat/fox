# Disassembler

## objdump
```toml
[Autostart.Disassembler]
Name = "objdump"
Pattern = '.*\.(bin|dll|exe|scr|sys)'
Commands = [
  'objdump -D "{{file}}"',
]
```

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

# Eric Zimmerman's tools

## PECmd
```toml
[Autostart.Prefetch]
Name = "PECmd"
Pattern = '.*\.pf'
Commands = [
  'dotnet PECmd.dll -f "{{file}}" --mp',
]
```

## JLECmd
```toml
[Autostart.Jumplists]
Name = "JLECmd"
Pattern = '.*\.(automatic|custom)Destination-ms'
Commands = [
  'dotnet JLECmd.dll -f "{{file}}" --mp --fd',
]
```
