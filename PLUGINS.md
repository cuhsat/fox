# Plugins

## JLECmd
```toml
[Autostart.Jumplists]
Name = "JLECmd"
Pattern = '.*\.(automatic|custom)Destination-ms'
Commands = [
  'dotnet JLECmd.dll -f "$FILE" --mp --fd',
]
```

## PECmd
```toml
[Autostart.Prefetch]
Name = "PECmd"
Pattern = '.*\.pf'
Commands = [
  'dotnet PECmd.dll -f "$FILE" --mp',
]
```

## SBECmd
```toml
[Plugin.F1]
Name = "SBECmd"
Commands = [
  'dotnet SBECmd.dll -d "$PARENT"',
]
```
