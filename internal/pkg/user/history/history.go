package history

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user"
)

const (
	filename = ".fox_history"
)

type History struct {
	sync.RWMutex

	file  *os.File // file handle
	lines []string // buffer lines
	index int      // buffer index
}

func New() *History {
	var err error

	h := History{
		lines: make([]string, 0),
	}

	_, p := user.Config(filename)

	h.file, err = os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

	if err != nil {
		sys.Error(err)
		return &h
	}

	s := bufio.NewScanner(h.file)

	for s.Scan() {
		t := strings.SplitN(s.Text(), ":", 2)

		if len(t) > 1 {
			h.lines = append(h.lines, t[1])
		}
	}

	err = s.Err()

	if err != nil {
		sys.Error(err)
	}

	h.index = len(h.lines)

	return &h
}

func (h *History) AddEntry(r, s string) {
	defer h.Reset()

	h.Lock()
	h.lines = append(h.lines, s)
	h.Unlock()

	if h.file == nil {
		return
	}

	l := fmt.Sprintf("%10d;%s:%s", time.Now().Unix(), r, s)

	_, err := fmt.Fprintln(h.file, l)

	if err != nil {
		sys.Error(err)
	}
}

func (h *History) AddCommand(cmd string) {
	h.AddEntry("user", cmd)
}

func (h *History) PrevCommand() string {
	h.RLock()
	defer h.RUnlock()

	if h.index > 0 {
		h.index--
	} else {
		return ""
	}

	return h.lines[h.index]
}

func (h *History) NextCommand() string {
	h.RLock()
	defer h.RUnlock()

	if h.index < len(h.lines)-1 {
		h.index++
	} else {
		return ""
	}

	return h.lines[h.index]
}

func (h *History) Reset() {
	h.RLock()
	h.index = len(h.lines)
	h.RUnlock()
}

func (h *History) Close() {
	if h.file != nil {
		h.file.Close()
	}
}
