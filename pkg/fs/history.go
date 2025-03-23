package fs

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
)

const (
    FileHistory = ".cu_history"    
)

type History struct {
    file  *os.File // file handle
    lines []string // buffer lines
    index int      // buffer index
}

func NewHistory() *History {
    dir, err := os.UserHomeDir()

    if err != nil {
        Panic(err)
    }

    f, err := os.OpenFile(filepath.Join(dir, FileHistory), FLAG_LOG, MODE_FILE)

    if err != nil {
        Panic(err)
    }

    var l []string

    s := bufio.NewScanner(f)
    
    for s.Scan() {
        l = append(l, s.Text())
    }
    
    err = s.Err()

    if err != nil {
        Panic(err)
    }

    return &History{
        file: f,
        lines: l,
        index: len(l),
    }
}

func (h *History) AddCommand(s string) {
    defer h.Reset()

    h.lines = append(h.lines, s)

    _, err := fmt.Fprintln(h.file, s)

    if err != nil {
        Error(err)
    }
}

func (h *History) PrevCommand() string {
    if h.index > 0 {
        h.index--
    }

    return h.lines[h.index]
}

func (h *History) NextCommand() string {
    if h.index < len(h.lines)-1 {
        h.index++
    } else {
        return ""
    }

    return h.lines[h.index]
}

func (h *History) Reset() {
    h.index = len(h.lines)
}

func (h *History) Close() {
    h.file.Close()
}
