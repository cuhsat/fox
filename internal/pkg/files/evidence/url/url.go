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
	c   *http.Client  // http client
	url string        // export url
	scm schema.Schema // export schema
}

func New(url string, scm schema.Schema) *Url {
	return &Url{new(http.Client), url, scm}
}

func (w *Url) Open(_ *os.File, _ bool, _ string) {}

func (w *Url) Begin() {}

func (w *Url) Flush() {
	body := strings.NewReader(w.scm.String())

	req, err := http.NewRequest("POST", w.url, body)

	if err != nil {
		sys.Error(err)
		return
	}

	for k, v := range w.scm.Headers() {
		req.Header.Set(k, v)
	}

	res, err := w.c.Do(req)

	if err != nil {
		sys.Error(err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		sys.Error(http.StatusText(res.StatusCode))
	}
}

func (w *Url) WriteMeta(meta evidence.Meta) {
	w.scm.SetMeta(meta)
}

func (w *Url) WriteLine(nr, grp int, str string) {
	w.scm.AddLine(nr, grp, str)
}
