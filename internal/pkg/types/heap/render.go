package heap

import (
	"encoding/json"
)

func (h *Heap) Render() {
	h.RLock()
	cached := h.last().rmap != nil
	h.RUnlock()

	if cached {
		return // use cache
	}

	l := h.last()

	h.Lock()
	l.rmap = l.smap.Render()
	h.Unlock()
}

func (h *Heap) Reset() {
	l := h.last()

	h.Lock()
	l.rmap = nil
	h.Unlock()
}

func (h *Heap) Wrap(w int) {
	l, s := h.last(), ""

	if len(*l.smap) > 0 {
		s = (*l.smap)[0].Str
	}

	h.Lock()

	if len(s) > 0 && json.Valid([]byte(s)) {
		l.rmap = l.smap.Indent()
	} else {
		l.rmap = l.smap.Wrap(w)
	}

	h.Unlock()
}
