package jsonl

import (
    "bytes"
    "encoding/json"
    "path/filepath"
    "strings"
)

const (
    JSONL = ".jsonl"
)

const (
    Indent = "  "
)

func Detect(p string) bool {
    return strings.ToLower(filepath.Ext(p)) == JSONL
}

func Pretty(s string) []string {
    var b bytes.Buffer

    if len(s) == 0 {
        return []string{""}
    }

    err := json.Indent(&b, []byte(s), "", Indent)

    if err != nil {
        return []string{err.Error()}
    }

    return strings.Split(b.String(), "\n")
}
