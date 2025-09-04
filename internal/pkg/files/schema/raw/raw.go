package raw

import (
	"fmt"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/files/evidence"
)

type Raw struct {
	headers map[string]string
	content strings.Builder
}

func New() *Raw {
	return &Raw{
		headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
}

func (raw *Raw) String() string {
	return raw.content.String()
}

func (raw *Raw) Headers() map[string]string {
	return raw.headers
}

func (raw *Raw) SetMeta(meta evidence.Meta) {
	raw.headers["x-evidence-path"] = meta.Path
	raw.headers["x-evidence-hash"] = fmt.Sprintf("%x", meta.Hash)
}

func (raw *Raw) AddLine(_, _ int, str string) {
	raw.content.WriteString(str)
	raw.content.WriteRune('\n')
}
