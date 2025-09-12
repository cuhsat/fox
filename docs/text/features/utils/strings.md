# Strings
**ASCII** and **Unicode** strings from all open files can be carved with the [`fox strings`](../../basics/usage/strings.md) command or using a [hotkey](../ui/keymap.md).

## Indicator of Compromises
Built-in detection of:

- `UUID`
- `IPv4`
- `IPv6`
- `MAC`
- `URL`
- `Mail`

## Example
```console
$ fox strings -pi testdata/test.ioc
00000000  data  Hello World
0000000c  ipv4  127.0.0.1
00000016  ipv6  2001:0db8:85a3:08d3:1319:8a2e:0370:7344
0000003e  mac   00:80:41:ae:fd:7e
00000050  mail  test@example.org
00000061  url   https://example.org
00000075  uuid  550e8400-e29b-11d4-a716-446655440000
```
