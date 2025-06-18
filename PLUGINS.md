# Disassembler

## objdump
```toml
[Autostart.objdump]
Name = "objdump"
Pattern = '.*\.(bin|dll|exe|scr|sys)'
Commands = [
  'objdump -D "$FILE"',
]
```

# Dissect

## target-info
```toml
[Autostart.Info]
Name = "target-info"
Pattern = '.*\.(dd|img|raw|ad1|asdf|E0?|00?)'
Commands = [
  'target-info "$FILE"',
]
```

## target-query
```toml
[Plugin.F7]
Name = "target-query"
Prompt = "Query"
Commands = [
  'target-query -j -f "$INPUT" "$BASE"',
]
```

# Eric Zimmerman's tools

## PECmd
```toml
[Autostart.PECmd]
Name = "PECmd"
Pattern = '.*\.pf'
Commands = [
  'dotnet PECmd.dll -f "$FILE" --mp',
]
```

## JLECmd
```toml
[Autostart.JLECmd]
Name = "JLECmd"
Pattern = '.*\.(automatic|custom)Destination-ms'
Commands = [
  'dotnet JLECmd.dll -f "$FILE" --mp --fd',
]
```

## SBECmd
```toml
[Plugin.F8]
Name = "SBECmd"
Commands = [
  'dotnet SBECmd.dll -d "$PARENT"',
]
```
