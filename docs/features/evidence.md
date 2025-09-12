# Evidence
By pressing <kbd>Ctrl</kbd> + <kbd>S</kbd> all filtered lines of the current file will be saved into an **Evidence Bag** and cryptographically signed. The bag can be (re)opened by pressing <kbd>Ctrl</kbd> + <kbd>B</kbd>.

> All timestamps will be normalized to UTC.

## Evidence Bag
Available evidence bag formats:

| Type      | Details                                          |
|-----------|--------------------------------------------------|
| `plain`   | Plain data (*default*)                           |
| `text`    | Plain data with metadata                         |
| `json`    | [JSON Schema](../../../api/evidence.schema.json) |
| `jsonl`   | [JSON Schema](../../../api/evidence.schema.json) |
| `xml`     | [XML Schema](../../../api/evidence.schema.xsd)   |
| `sqlite3` | [SQL Schema](../../../api/evidence.schema.sql)   |

## Evidence Signing
The evidence bag will always be signed using `SHA256` to guarantee the **Chain of Custody**.

> To use `HMAC-SHA256` signing, specific a **Key Phrase** using the `--key` flag.

## Evidence Streaming
All evidence can additionally be streamed to a server using the `--url` flag.

### Raw
If no schema was specified, the raw data will be streamed alongside with custom `x-evidence-*` HTTP-headers:

```
Hello World
```

### ECS Schema
Evidence can be streamed to an URL using the [Elastic Common Schema](https://www.elastic.co/docs/reference/ecs).

> To use a local running **Logstash** instance, set the `-L` flag.

```json
{
  "@timestamp": "2025-09-07T15:45:00.190311Z",
  "message": "Hello World\n",
  "labels": {
    "case": "2025-09-07",
    "filters": ""
  },
  "agent": {
    "type": "Forensic Examiner",
    "version": "1.2.3"
  },
  "ecs": {
    "version": "9.1.0"
  },
  "file": {
    "mtime": "2025-09-07T15:45:00.190311Z",
    "path": "/hello.txt",
    "size": 6,
    "hash": {
      "sha256": "d2a84f4b8b650937ec8f73cd8be2c74add5a911ba64df27458ed8229da804a26"
    }
  },
  "user": {
    "name": "jd",
    "full_name": "John Doe"
  }
}
```

### HEC Schema
Evidence can be streamed to a [Splunk HTTP Event Collector](https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTinput), an **Authorization Token** is required for this.

> To use a local running **Splunk** instance, set the `-S` flag.

```json
{
  "time": 1757260053699,
  "source": "Forensic Examiner",
  "sourcetype": "_json",
  "index": "2025-09-07",
  "event": {
    "user": "jd (John Doe)",
    "path": "/hello.txt",
    "hash": "d2a84f4b8b650937ec8f73cd8be2c74add5a911ba64df27458ed8229da804a26",
    "time": 1757260053699,
    "size": 6,
    "lines": [
      "Hello World\n"
    ]
  }
}
```
