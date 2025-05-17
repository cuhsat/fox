//go:build no_ai

package ai

import (
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

const (
	Build = false
)

type Agent struct {
	// stub
}

func Init(model string) bool {
	return false
}

func NewAgent() *Agent {
	return nil
}

func (_ *Agent) Path() string {
	return ""
}

func (_ *Agent) Close() {
	// stub
}

func (_ *Agent) Prompt(s string, b []byte) {
	// stub
}

func (_ *Agent) Listen(hi *history.History) {
	// stub
}
