# Windows Event Logs
Windows event logs (`.evtx`) will automatically be detected by its magic bytes, parsed and converted to [JSON Lines](https://jsonlines.org/) output for further processing.

> To prevent the automatic parsing and converting of event logs, use either the `-r` or `--no-convert` flag. 

## Example
```console
$ fox -p testdata/test.evtx
{"Event":{"EventData":{"Binary":null,"Data":["WIN-TIP3N90KK74","IEUSER-S9CSARBA"]},"System":{"Channel":"System","Computer":"WIN-TIP3N90KK74","EventID":{"Qualifiers":"32768","Value":"6011"},"EventLog":"","EventRecordID":"1","Keywords":"0x80000000000000","Level":"4","Security":{},"Task":"0","TimeCreated":{"SystemTime":"2018-01-03T06:28:19Z"}}}}
{"Event":{"EventData":{"Binary":null,"Data":["6.03.","9600","Multiprocessor Free","16384"]},"System":{"Channel":"System","Computer":"WIN-TIP3N90KK74","EventID":{"Qualifiers":"32768","Value":"6009"},"EventLog":"","EventRecordID":"2","Keywords":"0x80000000000000","Level":"4","Security":{},"Task":"0","TimeCreated":{"SystemTime":"2018-01-03T06:28:19Z"}}}}
{"Event":{"EventData":{"Binary":"E20701000300030006001C001300A3030000000000000000","Data":null},"System":{"Channel":"System","Computer":"WIN-TIP3N90KK74","EventID":{"Qualifiers":"32768","Value":"6005"},"EventLog":"","EventRecordID":"3","Keywords":"0x80000000000000","Level":"4","Security":{},"Task":"0","TimeCreated":{"SystemTime":"2018-01-03T06:28:19Z"}}}}
...
```
