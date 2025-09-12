# Hashes
Hashes of all open files can be generated with the [`fox hash`](../../basics/usage/hash.md) command or using a [hotkey](../ui/keymap.md).

## Cryptographic Hashes
Built-in cryptographic hashes *(for chain-of-custody)*:

- `MD5`
- `SHA1`
- `SHA256`
- `SHA3`
- `SHA3-224`
- `SHA3-256`
- `SHA3-384`
- `SHA3-512`

## Similarity Hashes
Built-in similarity hashes *(for malware detection)*:
 
- `SDHASH`
- `SSDEEP`
- `TLSH`

## Checksums
Built-in checksums *(for file verification)*:

- `CRC32-IEEE`
- `CRC64-ECMA`
- `CRC64-ISO`

## Example
```console
$ fox hashes -pt=md5 testdata/test.bin
534a6a08a693b374803b6eda2bf8baab  testdata/test.bin
```
