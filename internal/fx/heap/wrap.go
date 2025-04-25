package heap

import (
    "github.com/cuhsat/fx/internal/fx/file"
)

func (h *Heap) Wrap(w int) {
    if h.RMap != nil {
        return // use cache
    }

    if file.CanIndent(h.Path) {
        h.RMap = h.SMap.Indent(h.MMap)
    } else {
        h.RMap = h.SMap.Wrap(w)
    }
}

func (h *Heap) Reset() {
    h.RMap = nil
}
