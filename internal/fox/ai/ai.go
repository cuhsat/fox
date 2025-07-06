package ai

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ollama/ollama/api"
)

const (
	Default = "" // disabled
)

var (
	Model = Default
	Alive = &api.Duration{Duration: time.Minute * 10}
)

var (
	mutex  sync.RWMutex
	client *api.Client
)

func Init(model string) {
	Model = model

	if len(model) == 0 || strings.ToLower(model) == "default" {
		return
	}

	llm, err := api.ClientFromEnvironment()

	if err != nil {
		return
	}

	// preload model in the background
	go func(*api.Client) {
		if llm.Chat(context.Background(), &api.ChatRequest{
			Model:     Model,
			KeepAlive: Alive,
		}, null) == nil {
			mutex.Lock()
			client = llm
			mutex.Unlock()
		}
	}(llm)
}

func IsInit() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return client != nil
}

func GetClient() *api.Client {
	mutex.RLock()
	defer mutex.RUnlock()
	return client
}

func null(_ api.ChatResponse) error {
	return nil
}
