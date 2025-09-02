package raw

import (
	"encoding/json"

	"github.com/cuhsat/fox/internal/pkg/files/evidence"
)

type Raw evidence.Evidence

func New() *Raw {
	return new(Raw)
}

func (raw *Raw) String() string {
	buf, err := json.Marshal(raw)

	if err == nil {
		return string(buf)
	} else {
		return err.Error()
	}
}

func (raw *Raw) Headers() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

func (raw *Raw) SetMeta(meta evidence.Meta) {
	raw.Meta = meta
}

func (raw *Raw) AddLine(nr, grp int, str string) {
	raw.Lines = append(raw.Lines, evidence.Line{Nr: nr, Grp: grp, Str: str})
}
