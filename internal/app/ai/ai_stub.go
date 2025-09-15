//go:build minimal

package ai

import (
	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/app/ai/agent"
)

func NewAgent(_ *app.Context) *agent.Agent {
	return nil
}
