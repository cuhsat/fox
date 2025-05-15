//go:build !ai

package ai

import (
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

const (
	Build = false
)

type Chat struct {
	// stub
}

func Init(model string) bool {
	return false
}

func NewChat() *Chat {
	return nil
}

func (_ *Chat) Path() string {
	return ""
}

func (_ *Chat) Close() {
	// stub
}

func (_ *Chat) Prompt(s string, b []byte) {
	// stub
}

func (_ *Chat) Listen(hi *history.History) {
	// stub
}
