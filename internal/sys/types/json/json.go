package json

import (
    "bytes"
    "encoding/json"
    "path/filepath"
    "strings"

    "github.com/cuhsat/fx/internal/sys"
)

const (
    Json  = ".json"
    JsonL = ".jsonl"
)

const (
    Indent = "    "
)

func Detect(p string) bool {
    ext := strings.ToLower(filepath.Ext(p))

    return ext == Json || ext == JsonL
}

func Pretty(j string) (s []string) {
    var b bytes.Buffer

    err := json.Indent(&b, []byte(j), "", Indent)

    if err != nil {
        sys.Fatal(err)
    }

    return strings.Split(b.String(), "\n")
}
