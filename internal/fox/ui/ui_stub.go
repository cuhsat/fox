//go:build no_ui

package ui

import (
	"github.com/cuhsat/fox/internal/pkg/types/heapset"
	"github.com/cuhsat/fox/internal/pkg/types/mode"
	"github.com/cuhsat/fox/internal/pkg/user/bag"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

const (
	Build = false
)

type UI struct {
	// stub
}

func New(_ mode.Mode) *UI {
	return nil
}

func (_ *UI) Run(hs *heapset.HeapSet, hi *history.History, bag *bag.Bag) {
	// stub
}

func (_ *UI) Close() {
	// stub
}
