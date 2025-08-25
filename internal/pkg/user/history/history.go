package history

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cuhsat/fox/internal/pkg/flags"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

type History struct {
	sync.RWMutex

	file  *os.File     // file handle
	lines []string     // buffer lines
	index atomic.Int64 // buffer index
}

func New() *History {
	var err error

	h := History{
		lines: make([]string, 0),
	}

	var m int

	if !flags.Get().Opt.Readonly {
		m = os.O_CREATE | os.O_APPEND | os.O_RDWR
	} else {
		m = os.O_RDONLY
	}

	h.file, err = os.OpenFile(sys.Config("history"), m, 0600)

	if err != nil {
		sys.Error(err)
		return &h
	}

	s := bufio.NewScanner(h.file)

	for s.Scan() {
		t := strings.SplitN(s.Text(), ";", 2)

		if len(t) > 1 {
			h.lines = append(h.lines, t[1])
		}
	}

	err = s.Err()

	if err != nil {
		sys.Error(err)
	}

	h.index.Store(int64(len(h.lines)))

	return &h
}

func (h *History) AddLine(s string) {
	defer h.Reset()

	// prepare string
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)

	h.Lock()
	h.lines = append(h.lines, s)
	h.Unlock()

	if h.file == nil {
		return
	}

	l := fmt.Sprintf("%10d;%s", time.Now().Unix(), s)

	h.Lock()
	_, _ = fmt.Fprintln(h.file, l)
	h.Unlock()
}

func (h *History) PrevLine() string {
	var d int64 = 0

	if h.index.Load() > 0 {
		d = -1
	}

	return h.get(h.index.Add(d))
}

func (h *History) NextLine() string {
	if h.index.Load() >= h.len()-1 {
		return ""
	}

	return h.get(h.index.Add(1))
}

func (h *History) Reset() {
	h.index.Store(h.len())
}

func (h *History) Close() {
	h.Lock()

	if h.file != nil {
		_ = h.file.Close()
	}

	h.Unlock()
}

func (h *History) len() int64 {
	h.RLock()
	defer h.RUnlock()
	return int64(len(h.lines))
}

func (h *History) get(idx int64) string {
	h.RLock()
	defer h.RUnlock()
	return h.lines[idx]
}
