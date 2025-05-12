package ai

import (
	"strings"

	"github.com/tmc/langchaingo/llms/ollama"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

const (
	Default = "mistral"
)

var (
	llm *ollama.LLM = nil
)

func Init(model string) bool {
	var err error

	if len(model) == 0 {
		model = Default
	}

	if strings.ToLower(model) == "default" {
		model = Default
	}

	llm, err = ollama.New(ollama.WithModel(model))

	if err != nil {
		sys.Error(err)
		return false
	}

	return true
}
