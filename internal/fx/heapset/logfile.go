package heapset

import (
	"errors"

	"github.com/cuhsat/fx/internal/fx"
    "github.com/cuhsat/fx/internal/fx/heap"
    "github.com/cuhsat/fx/internal/fx/types"
)

func (hs *HeapSet) Raise(msg string) {
    for i, h := range hs.heaps {
        if h.Type == types.StdErr {
            hs.index = i
	        return
        }
    }

    hs.heaps = append(hs.heaps, &heap.Heap{
    	Title: "error",
        Path: fx.Logfile.Name(),
        Type: types.StdErr,
    })

    hs.index = len(hs.heaps)-1

    hs.load()

    if hs.error_fn != nil {
    	hs.error_fn(errors.New(msg))
    }
}
