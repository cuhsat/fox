# Linux Systemd Journals
Linux systemd journals (`.journal`) will automatically be detected by its file extension, parsed and converted to [JSON Lines](https://jsonlines.org/) output for further processing.

> To prevent the automatic parsing and converting of journals, use either the `-r` or `--no-convert` flag. 

## Example
```console
$ fox -p testdata/test.journal
{"System":{"Seq":1,"Timestamp":"2022-10-24T02:00:01.009823+02:00","_SOURCE_MONOTONIC_TIMESTAMP":"0","_TRANSPORT":"kernel","_BOOT_ID":"35e8501129134edd9df5267c49f744a4","_MACHINE_ID":"46eff2c0526d485fa8103457ae6f7146","_HOSTNAME":"fedora"},"EventData":{"PRIORITY":6,"SYSLOG_FACILITY":0,"SYSLOG_IDENTIFIER":"kernel","MESSAGE":"Booting Linux on physical CPU 0x0000000000 [0x410fd083]"}}
{"System":{"Seq":2,"Timestamp":"2022-10-24T02:00:01.009952+02:00","_SOURCE_MONOTONIC_TIMESTAMP":"0","_TRANSPORT":"kernel","_BOOT_ID":"35e8501129134edd9df5267c49f744a4","_MACHINE_ID":"46eff2c0526d485fa8103457ae6f7146","_HOSTNAME":"fedora"},"EventData":{"SYSLOG_FACILITY":0,"SYSLOG_IDENTIFIER":"kernel","PRIORITY":5,"MESSAGE":"Linux version 6.0.7-301.fc37.aarch64 (mockbuild@buildvm-a64-05.iad2.fedoraproject.org) (gcc (GCC) 12.2.1 20220819 (Red Hat 12.2.1-2), GNU ld version 2.38-24.fc37) #1 SMP PREEMPT_DYNAMIC Fri Nov 4 18:13:35 UTC 2022"}}
{"System":{"Seq":3,"Timestamp":"2022-10-24T02:00:01.010016+02:00","_SOURCE_MONOTONIC_TIMESTAMP":"0","_TRANSPORT":"kernel","_BOOT_ID":"35e8501129134edd9df5267c49f744a4","_MACHINE_ID":"46eff2c0526d485fa8103457ae6f7146","_HOSTNAME":"fedora"},"EventData":{"PRIORITY":6,"SYSLOG_FACILITY":0,"SYSLOG_IDENTIFIER":"kernel","MESSAGE":"Machine model: Raspberry Pi 4 Model B Rev 1.1"}}
...
```
