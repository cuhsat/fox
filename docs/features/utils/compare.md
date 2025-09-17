# Compare
Two files can be compared with the [`fox compare`](../../basics/usage/compare.md) command.

## Example
```console
$ fox compare -p testdata/test.ioc testdata.diff
- 1 Hello World
+ 1 Hello Earth
  2 127.0.0.1
- 3 2001:0db8:85a3:08d3:1319:8a2e:0370:7344
  4 00:80:41:ae:fd:7e
  5 test@example.org
  6 https://example.org
- 7 550e8400-e29b-11d4-a716-446655440000
+ 6 EOF
```
