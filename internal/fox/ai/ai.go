package ai

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ollama/ollama/api"
)

var (
	ErrNotAvailable = errors.New("AI is not available")
)

var (
	Alive = &api.Duration{Duration: time.Minute * 10}
	Model = ""
)

var (
	mutex  sync.RWMutex
	client *api.Client
)

func Init(model string) {
	if len(model) == 0 {
		return
	}

	llm, err := api.ClientFromEnvironment()

	if err != nil {
		return
	}

	Model = model

	go func(*api.Client) {
		if llm.Chat(context.Background(), &api.ChatRequest{
			Model:     Model,
			KeepAlive: Alive,
		}, func(cr api.ChatResponse) error {
			return nil // preload model
		}) == nil {
			mutex.Lock()
			client = llm
			mutex.Unlock()
		}
	}(llm)
}

func IsInit() bool {
	ch := make(chan bool, 1)

	go func() {
		for GetClient() == nil {
			time.Sleep(time.Millisecond * 100)
		}

		ch <- true
	}()

	select {
	case <-ch:
		return true // ready

	case <-time.After(time.Second * 2):
		return false // timeout
	}
}

func GetClient() *api.Client {
	mutex.RLock()
	defer mutex.RUnlock()
	return client
}
