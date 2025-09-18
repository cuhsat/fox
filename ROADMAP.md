# Roadmap

## 1. Bug Fixes
* Fix tab to smaller file in HEX mode
* Fix SMap wrap to wrap with rune width
* Fix partial highlighting

## 2. Refactorings
* Implement State Machine for UI
  * Add keychain support
  * Remap some key bindings

## 3. Features
* Add alternative key binding for GOTO mode
* Add horizontal scroll support to input field
* Add support for Fortinet log files
  * https://cyber.wtf/2024/08/30/parsing-fortinet-binary-firewall-logs/
  * https://github.com/GDATAAdvancedAnalytics/FortilogDecoder/blob/main/fortilog_decoder.py
* Add persistence to RAG (`--persist=DB` flag)
* Add possibility to hash with many formats at once
  * Usage `--types=MD5,SHA1`
* Add additional documentation
  * Add manpage(s)
  * Add `bash` and `zsh` autocompletion files
    * https://applejag.eu/blog/go-spf13-cobra-custom-flag-types/
    * https://cobra.dev/docs/how-to-guides/shell-completion/
    * https://github.com/spf13/cobra/blob/v1.8.0/site/content/completions/_index.md#completions-for-flags
* Add *Yara* rules scan
* Add Super Timeline ability
  1. Specify *Common Log Format* (CLF)
  2. Extract timestamps per RegEx 
  3. Sort lines by timestamp

## 4. Optimizations
* Add debug output (`-v` flag)
* SMap speed
  * https://dev.to/moseeh_52/efficient-file-reading-in-go-mastering-bufionewscanner-vs-osreadfile-4h05
  * https://dave.cheney.net/high-performance-json.html

## 5. Ideas
* Colorize CLI output?
  * https://github.com/logrusorgru/aurora
  * https://github.com/cyucelen/marker
* Watch configs for changes?
  * `viper.WatchConfig()`
  * `viper.OnConfigChange(func(e fsnotify.Event)`
* Add auto-completion from history?
* Add search to HEX mode?
* Add `readline` config support?
  * https://github.com/chzyer/readline
* Add `find` command ability?
* Use multiple parallel filters?
* Use reflow algos?
  * https://github.com/muesli/reflow
* Generic syntax highlighting?
  * `Start Color [ … ] End Color`
  * `{}[]<>()““‘‘:;`
