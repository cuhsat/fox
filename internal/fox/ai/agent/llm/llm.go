package llm

import (
    "context"
    "fmt"

    "github.com/ollama/ollama/api"

    "github.com/hiforensics/fox/internal/fox"
    "github.com/hiforensics/fox/internal/fox/ai"
    "github.com/hiforensics/fox/internal/pkg/sys"
)

type LLM struct {
    history []api.Message // chat history
}

func New() *LLM {
    return &LLM{
        history: make([]api.Message, 0),
    }
}

func (llm *LLM) Ask(query, lines string, fn api.ChatResponseFunc) {
    llm.history = append(llm.history, api.Message{
        Role:    "user",
        Content: fmt.Sprintf(fox.Prompt, query, lines),
    })

    ctx := context.Background()
    req := &api.ChatRequest{
        Model:     ai.Model,
        KeepAlive: ai.Alive,
        Messages:  llm.history,
        Options: map[string]any{
            "temperature": 0.2,
            "top_p":       0.5,
            "top_k":       10,
            "seed":        8211,
        },
    }

    err := ai.GetClient().Chat(ctx, req, fn)

    if err != nil {
        sys.Error(err)
    }
}
