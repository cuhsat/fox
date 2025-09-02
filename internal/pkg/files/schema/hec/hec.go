// Package hec specification:
// https://docs.splunk.com/Documentation/Splunk/latest/Data/FormateventsforHTTPEventCollector
package hec

import (
	"encoding/json"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/files/evidence"
)

type Hec struct {
	Time       int    `json:"time"`
	Source     string `json:"source"`
	Sourcetype string `json:"sourcetype"`
	Index      string `json:"index"`

	Event struct {
		Lines []line `json:"lines"`
	} `json:"event"`
}

type line []struct {
	Nr   int    `json:"nr"`
	Grp  int    `json:"grp"`
	Data string `json:"data"`
}

func New() *Hec {
	hec := new(Hec)

	hec.Source = app.Product

	return hec
}

func (hec *Hec) String() string {
	buf, err := json.Marshal(hec)

	if err == nil {
		return string(buf)
	} else {
		return err.Error()
	}
}

func (hec *Hec) SetMeta(meta evidence.Meta) {
	hec.Time = 0
	hec.Sourcetype = meta.Path
	hec.Index = meta.Name
}

func (hec *Hec) AddLine(nr, grp int, str string) {
	hec.Event.Lines = append(hec.Event.Lines, line{nr, grp, str})
}
