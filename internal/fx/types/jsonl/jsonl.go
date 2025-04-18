package jsonl

import (
    "bytes"
    "encoding/json"
    "path/filepath"
    "strings"
)

const (
    indent = "  "
)

func Detect(path string) bool {
    ext := strings.ToLower(filepath.Ext(path))

    return ext == ".jsonl"
}

func Pretty(s string) []string {
    var buf bytes.Buffer

    if len(s) == 0 {
        return []string{""}
    }

    err := json.Indent(&buf, []byte(s), "", indent)

    if err != nil {
        return []string{s}
    }

    return strings.Split(buf.String(), "\n")
}
