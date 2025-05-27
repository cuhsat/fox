//go:build no_ui || no_ai

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

func NewAgent(_ string) *Agent {
	return nil
}

func (_ *Agent) Path() string {
	return ""
}

func (_ *Agent) Load() {
	// stub
}

func (_ *Agent) Close() {
	// stub
}

func (_ *Agent) Prompt(_ string, _ *heap.Heap) {
	// stub
}

func (_ *Agent) Listen(_ *history.History) {
	// stub
}
