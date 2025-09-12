# CSV
Comma-separated value files (`.csv`) will automatically be detected by file extension. The built-in sniffer will detect different delimiters and will pad the output left aligned for further processing.  

> Use <kbd>Ctrl</kbd> + <kbd>P</kbd> to (un)pin the first line while in the **UI**.

## Delimiters
Detected delimiters are:

| Character | Name      |
|-----------|-----------|
| `,`       | Comma     |
| `:`       | Colon     |
| `;`       | Semicolon |
| `|`       | Pipe      |
| `⇥`       | Tabulator |

## Example
```console
$ fox -p testdata/test.csv
MinTemp MaxTemp Rainfall Evaporation Sunshine WindGustDir WindGustSpeed WindDir9am WindDir3pm WindSpeed9am WindSpeed3pm Humidity9am Humidity3pm Pressure9am Pressure3pm Cloud9am Cloud3pm Temp9am Temp3pm RainToday RISK_MM RainTomorrow
8       24.3    0        3.4         6.3      NW          30            SW         NW         6            20           68          29          1019.7      1015        7        7        14.4    23.6    No        3.6     Yes
14      26.9    3.6      4.4         9.7      ENE         39            E          W          4            17           80          36          1012.4      1008.4      5        3        17.5    25.7    Yes       3.6     Yes
...
```
