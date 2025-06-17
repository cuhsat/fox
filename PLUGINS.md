# Fox-IT's Dissect tools

## Target-Info
```toml
[Autostart.Info]
Name = "target-info"
Pattern = '.*\.(dd|img|raw|ad1|asdf|E0?|00?)'
Commands = [
  'target-info "$FILE"',
]
```

## Target-Query
```toml
[Plugin.F7]
Name = "target-query"
Prompt = "Query"
Commands = [
  'target-query -j -f "$INPUT" "$BASE"',
]
```

# Eric Zimmerman's tools

## PECmd (Windows Prefetch files)
```toml
[Autostart.PECmd]
Name = "PECmd"
Pattern = '.*\.pf'
Commands = [
  'dotnet PECmd.dll -f "$FILE" --mp',
]
```

## JLECmd (Windows Jumplists)
```toml
[Autostart.JLECmd]
Name = "JLECmd"
Pattern = '.*\.(automatic|custom)Destination-ms'
Commands = [
  'dotnet JLECmd.dll -f "$FILE" --mp --fd',
]
```

## SBECmd (Windows Shellbags)
```toml
[Plugin.F8]
Name = "SBECmd"
Commands = [
  'dotnet SBECmd.dll -d "$PARENT"',
]
```

# Misc

## objdump (Disassembler)
```toml
[Autostart.objdump]
Name = "objdump"
Pattern = '.*\.(bin|dll|exe|scr|sys)'
Commands = [
  'objdump -D "$FILE"',
]
```
