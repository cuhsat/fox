# JSON & JSON Lines
JavaScript Object Notation files (`.json` or `.jsonl`) will automatically be detected by evaluating the first line for valid JSON.

> Use <kbd>Ctrl</kbd> + <kbd>W</kbd> to wrap and indent [JSON Lines](https://jsonlines.org/) while in the **UI**.

## Example
```console
$ fox -p testdata/test.json
[
  {
    "name": "Adeel Solangi",
    "language": "Sindhi",
    "id": "V59OF92YF627HFY0",
    "bio": "Donec lobortis eleifend condimentum. Cras dictum dolor lacinia lectus vehicula rutrum. Maecenas quis nisi nunc. Nam tristique feugiat est vitae mollis. Maecenas quis nisi nunc.",
    "version": 6.1
  },
...
```
