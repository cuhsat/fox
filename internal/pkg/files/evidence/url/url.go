package url

import (
	"net/http"
	"os"
	"strings"

	"github.com/cuhsat/fox/internal/pkg/files/evidence"
	"github.com/cuhsat/fox/internal/pkg/files/schema"
	"github.com/cuhsat/fox/internal/pkg/sys"
)

type Url struct {
	url string        // export url
	scm schema.Schema // export schema
}

func New(url string, scm schema.Schema) *Url {
	return &Url{
		url: url,
		scm: scm,
	}
}

func (w *Url) Open(_ *os.File, _ bool, _ string) {}

func (w *Url) Begin() {}

func (w *Url) Flush() {
	res, err := http.Post(w.url, "application/json", strings.NewReader(w.scm.String()))

	if err != nil {
		sys.Error(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		sys.Error(http.StatusText(res.StatusCode))
	}
}

func (w *Url) WriteMeta(meta evidence.Meta) {
	w.scm.SetMeta(meta)
}

func (w *Url) WriteLine(nr, grp int, str string) {
	w.scm.AddLine(nr, grp, str)
}
