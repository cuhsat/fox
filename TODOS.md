# TODOS

## Bug Fixes
- Fix hex bug when tab to smaller file
- Fix wrap to wrap with rune width
- Fix partial highlighting

## Features
- Add auto pulling for models (embed and model)
- Add alternative key binding for goto
- Refactor UI to state machine
  - Add key chains
  - Remap key bindings
- Add `find` command abilities?
- Add `CMD` command mode like vim?
  - Also per plugin
  - Examples
    - `open ...`
    - `set model mistral`
- Add super timeline per group regex
  - Define *(Minimal) Common Log Format*
  - Close all other regular heaps
  - Sort lines by timestamp
- Add scan with *Yara* rules
- Add documentation
  - Manpage
  - *Bash* and *Zsh* autocompletion files
    - https://applejag.eu/blog/go-spf13-cobra-custom-flag-types/
    - https://cobra.dev/docs/how-to-guides/shell-completion/
    - https://github.com/spf13/cobra/blob/v1.8.0/site/content/completions/_index.md#completions-for-flags

## Ideas
- Add color to output
  - https://github.com/logrusorgru/aurora
  - https://github.com/cyucelen/marker
- Add readline standards
  - https://github.com/chzyer/readline
- Use reflow algos?
  - https://github.com/muesli/reflow
- Horizontal scrollable input field
- Parallel, multiple filters?
- Autocomplete inputs from history?
- Add search to hex view
- Render while still loading?
- Add persistence to RAG
  - Add flag `--persist=DB`?
- Add possibility to hash many formats at once?
  - Flag `--types=MD5,SHA1`
- Watch configs for changes?
  ```
  viper.WatchConfig()
  viper.OnConfigChange(func(e fsnotify.Event) {
    // do something
  })
  ```
- Generic syntax highlighting?
  - `Start Color [ … ] End Color`
  - `{}[]<>()““‘‘:;`

## Misc
- Optimize speed
  - https://dev.to/moseeh_52/efficient-file-reading-in-go-mastering-bufionewscanner-vs-osreadfile-4h05
  - https://dave.cheney.net/high-performance-json.html
- Add something about *MITRE* to the readme
- Add more debug prints?

## Resources

### Icon
https://thenounproject.com/icon/fox-1486590/

### Font
RobotoCondensed-bold

### Colors
`#ffffff white`
`#0f88cd blue`
`#333333 black`

### Domains
`forensic.examiner.rocks`
`hinterland.tools`
`forensik.jetzt`
`forensik.wtf`
`fox.cu`

### Quotes
> Quaere et invenies.
> Übersetzung: Suche und du wirst finden.

> Veritas vincit.
> Die Wahrheit trägt den Sieg davon.
> Walther, Proverbia sententiaeque 33157s
