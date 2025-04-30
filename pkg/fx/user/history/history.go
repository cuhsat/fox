package history

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cuhsat/fx/pkg/fx/sys"
	"github.com/cuhsat/fx/pkg/fx/user"
)

const (
	filename = ".fx_history"
)

type History struct {
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
		t := strings.SplitN(s.Text(), ";", 1)

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

func (h *History) AddCommand(cmd string) {
	defer h.Reset()

	h.lines = append(h.lines, cmd)

	if h.file == nil {
		return
	}

	l := fmt.Sprintf("%10d;%s", time.Now().Unix(), cmd)

	_, err := fmt.Fprintln(h.file, l)

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
	if h.file != nil {
		h.file.Close()
	}
}
