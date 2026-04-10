# fincut

A CLI tool for trimming and diffing structured log files with regex-based filter pipelines.

---

## Installation

```bash
go install github.com/yourname/fincut@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/fincut.git && cd fincut && go build -o fincut .
```

---

## Usage

```bash
# Trim a log file using a regex filter pipeline
fincut trim --input app.log --filter "ERROR|WARN" --output trimmed.log

# Diff two structured log files
fincut diff --before before.log --after after.log

# Chain multiple filters
fincut trim --input app.log --filter "user_id=42" --filter "status=5\d\d" | fincut diff --before - --after latest.log
```

### Flags

| Flag | Description |
|------|-------------|
| `--input` | Path to the input log file |
| `--filter` | Regex pattern to filter log lines (repeatable) |
| `--output` | Path to write trimmed output |
| `--before` | Baseline log file for diff |
| `--after` | Comparison log file for diff |
| `--format` | Log format hint: `json`, `logfmt`, or `plain` (default: `json`) |

---

## Example

```bash
$ fincut trim --input service.log --filter "level=error" --format json
[2024-01-15T10:23:01Z] ERROR  connection timeout  host=db-primary latency=3002ms
[2024-01-15T10:23:45Z] ERROR  auth failure        user_id=99 ip=10.0.0.5
```

---

## License

MIT © 2024 yourname