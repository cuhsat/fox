# Plugin
> Must be located under `~/fx_plugins`.

```toml
[Plugin.F1]
Name = "target-info"
Exec = "target-info $+"

[Plugin.F2]
Name = "target-query"
Exec = "target-query -f $? $+"
Input = "Func"
```

## Variabes
* `$?` User input
* `$+` Current file path
* `$*` All open file paths
