# Plugin
> Must be located under `~/fx_plugins`.

```toml
[Plugin.F1]
Name = "Echo Plugin"
Exec = "echo $?"
Input = true
```

## Variabes
* `$?` User input
* `$+` Current file path
* `$*` All open file paths
