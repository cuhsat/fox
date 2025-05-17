//go:build no_ai

package ai

import (
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

const (
	Build = false
)

type Agent struct {
	// stub
}

func Init() bool {
	return false
}

func NewAgent(model string) *Agent {
	return nil
}

func (_ *Agent) Path() string {
	return ""
}

func (_ *Agent) Close() {
	// stub
}

func (_ *Agent) Prompt(s string, h *heap.Heap) {
	// stub
}

func (_ *Agent) Listen(hi *history.History) {
	// stub
}
