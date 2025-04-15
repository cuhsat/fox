package json

import (
    "bytes"
    "encoding/json"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/sys/text"
)

const (
    Json  = ".json"
    JsonL = ".jsonl"
)

const (
    Indent = "  "
)

func Detect(p string) bool {
    e := strings.ToLower(filepath.Ext(p))

    return e == JsonL || (e == Json && text.Lines(p) <= 1)
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
