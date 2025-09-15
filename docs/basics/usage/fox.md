# fox
The Swiss Army Knife for examining text files.

## Usage
```console
fox [ACTION] [FLAG ...] [PATH ...]
```

### Actions
`counts`, `deflate`, `entropy`, `hash`, `strings`

### Arguments
Path(s) to open or `-` for **STDIN**

### Print
- `-p`, `--print` — print directly to console
- `--no-file` — don't print filenames
- `--no-line` — don't print line numbers

### Deflate
- `-P`, `--pwd=PASSWORD` — password for decryption (only **RAR**, **ZIP**)

### Hex display
- `-x`, `--hex` — show file in canonical hex

### File limits
- `-`, `--head` — limit head of file by *...*
- `-t`, `--tail` — limit tail of file by *...*
- `-n`, `--lines[=NUMBER]` — number of lines (*default:* `10`)
- `-c`, `--bytes[=NUMBER]` — number of bytes (*default:* `16`)

### Line filter
- `-e`, `--regexp=PATTERN` — filter for lines that match pattern
- `-C`, `--context=NUMBER` — number of lines surrounding context of match
- `-B`, `--before=NUMBER` — number of lines leading context before match
- `-A`, `--after=NUMBER` — number of lines trailing context after match

### AI agent
- `-m`, `--model=MODEL` — model for the agent to use
- `-q`, `--query=QUERY` — query for the agent to process

### AI model
- `--num-ctx=SIZE` — context window length (*default:* `4096`)
- `--temp=DECIMAL` — option for temperature (*default:* `0.2`)
- `--topp=DECIMAL` — option for model top_p (*default:* `0.5`)
- `--topk=NUMBER` — option for model top_k (*default:* `10`)
- `--seed=NUMBER` — option for random seed (*default:* `8211`)

### UI flags
- `--state={N|W|T|-}` — sets the used UI state flags
- `--theme=THEME` — sets the used UI theme
- `--legacy` — don't use any unicode decorations (**ISO 8859-1**)

### Evidence bag
- `-N`, `--case=NAME` — evidence bag case name (*default:* `YYYY-MM-DD`)
- `-f`, `--file=FILE` — evidence bag file name (*default:* `evidence`)
- `--mode=MODE` — evidence bag file mode (*default:* `plain`)
- `-k, --key=PHRASE` — key phrase to sign evidence bag via **HMAC-SHA256**

Available evidence bag modes:
> none, plain, text, json, jsonl, xml, sqlite

### Evidence URL
- `-u`, `--url=SERVER` — forward evidence to server address
- `-a`, `--auth=TOKEN` — forward evidence using auth token
- `--ecs` — use **[ECS](https://forensic-examiner.eu/features/evidence.html#ecs-schema)** schema for evidence
- `--hec` — use **[HEC](https://forensic-examiner.eu/features/evidence.html#hec-schema)** schema for evidence

### Disable
- `-R`, `--readonly` — don't write any new files
- `-r`, `--raw` — don't process files at all
- `--no-convert` — don't convert automatically
- `--no-deflate` — don't deflate automatically
- `--no-plugins` — don't run any plugins
- `--no-mouse` — don't use the mouse

### Aliases
- `-L`, `--logstash` — short for: `--ecs --url=http://localhost:8080`
- `-S`, `--splunk` — short for: `--hec --url=http://localhost:8088/...`
- `-T`, `--text` — short for: `--mode=text`
- `-j`, `--json` — short for: `--mode=json`
- `-J`, `--jsonl` — short for: `--mode=jsonl`
- `-s`, `--sqlite` — short for: `--mode=sqlite`
- `-X`, `--xml` — short for: `--mode=xml`

### Standard
- `--help` — shows this message
- `--credits` — shows the credits
- `--version` — shows the version

## Examples
```console
$ fox -pe "John Doe" ./**/*
```

```console
$ fox -pxhc=512 image.dd > mbr
```

```console
$ fox -pq="analyse this" log.evtx
```
