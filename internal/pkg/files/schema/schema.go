package schema

import (
	"github.com/cuhsat/fox/internal/pkg/files/evidence"
)

type Schema interface {
	String() string
	SetMeta(meta evidence.Meta)
	AddLine(nr, grp int, str string)
}
