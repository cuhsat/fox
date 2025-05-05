# Plugin
> Must be located under `~/fx_plugins`.

```toml
[Plugin.F1]
Name = "target-info"
Exec = "target-info $+ | unesc"

[Plugin.F2]
Name = "target-query"
Exec = "target-query -f $? $+ | unesc"
Input = "Func"
```

## Variabes
* `$?` User input
* `$+` Current file path
* `$*` All open file paths
