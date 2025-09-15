// Package hec specification: https://docs.splunk.com/Documentation/Splunk/latest/Data/FormateventsforHTTPEventCollector
package hec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cuhsat/fox/internal"
	"github.com/cuhsat/fox/internal/pkg/files/evidence"
	"github.com/cuhsat/fox/internal/pkg/flags"
)

type Hec struct {
	Time       int64  `json:"time"`
	Source     string `json:"source"`
	Sourcetype string `json:"sourcetype"`
	Index      string `json:"index,omitempty"`

	Event struct {
		User string `json:"user"`
		Path string `json:"path"`
		Hash string `json:"hash"`
		Time int64  `json:"time"`
		Size int64  `json:"size"`

		Lines []string `json:"lines"`
	} `json:"event"`
}

func New() *Hec {
	return &Hec{
		Source:     info.Product,
		Sourcetype: "_json",
	}
}

func (hec *Hec) String() string {
	buf, err := json.Marshal(hec)

	if err == nil {
		return string(buf)
	} else {
		return err.Error()
	}
}

func (hec *Hec) Headers() map[string]string {
	token := strings.ToLower(flags.Get().Bag.Auth)

	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Splunk %s", token),
	}
}

func (hec *Hec) SetMeta(meta evidence.Meta) {
	hec.Time = meta.Bagged.UnixMilli()
	hec.Index = meta.Name

	hec.Event.User = fmt.Sprintf("%s (%s)", meta.User.Username, meta.User.Name)
	hec.Event.Hash = fmt.Sprintf("%x", meta.Hash)
	hec.Event.Time = meta.Modified.UnixMilli()
	hec.Event.Size = meta.Size
	hec.Event.Path = meta.Path
}

func (hec *Hec) AddLine(_, _ int, str string) {
	hec.Event.Lines = append(hec.Event.Lines, fmt.Sprintf("%s\n", str))
}
