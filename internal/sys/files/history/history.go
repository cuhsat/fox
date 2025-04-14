package history

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "time"
    "strings"

    "github.com/cuhsat/fx/internal/sys"
)

const (
    File = ".fx_history"    
)

type History struct {
    file *os.File  // file handle
    lines []string // buffer lines
    index int      // buffer index
}

func NewHistory() *History {
    dir, err := os.UserHomeDir()

    if err != nil {
        sys.Fatal(err)
    }

    f, err := os.OpenFile(filepath.Join(dir, File), sys.O_HISTORY, 0644)

    if err != nil {
        sys.Fatal(err)
    }

    var l []string

    s := bufio.NewScanner(f)
    
    for s.Scan() {
        t := strings.SplitN(s.Text(), ";", 1)
        
        if len(t) > 1 {
            l = append(l, t[1])            
        }
    }
    
    err = s.Err()

    if err != nil {
        sys.Fatal(err)
    }

    return &History{
        file: f,
        lines: l,
        index: len(l),
    }
}

func (h *History) AddCommand(cmd string) {
    defer h.Reset()

    s := fmt.Sprintf("%10d;%s", time.Now().Unix(), cmd)

    _, err := fmt.Fprintln(h.file, s)

    h.lines = append(h.lines, cmd)

    if err != nil {
        sys.Error(err)
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
