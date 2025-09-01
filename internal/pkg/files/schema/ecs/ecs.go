// Package ecs specification:
// https://www.elastic.co/docs/reference/ecs/ecs-field-reference
package ecs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cuhsat/fox/internal/app"
	"github.com/cuhsat/fox/internal/pkg/files/evidence"
)

const Version = "9.1.0"

type Ecs struct {
	Timestamp time.Time         `json:"@timestamp"`
	Message   string            `json:"message"`
	Labels    map[string]string `json:"labels"`

	Agent struct {
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"agent"`

	Ecs struct {
		Version string `json:"version"`
	} `json:"ecs"`

	File struct {
		Mtime time.Time `json:"mtime"`
		Path  string    `json:"Path"`
		Size  int64     `json:"Size"`

		Hash struct {
			Sha256 string `json:"sha256"`
		} `json:"Hash"`
	} `json:"file"`

	User struct {
		Name     string `json:"Name"`
		FullName string `json:"full_name"`
	} `json:"User"`
}

func New() *Ecs {
	ecs := new(Ecs)
	ecs.Labels = make(map[string]string)

	ecs.Ecs.Version = Version

	ecs.Agent.Type = app.Product
	ecs.Agent.Version = app.Version[1:]

	return ecs
}

func (ecs *Ecs) String() string {
	buf, err := json.Marshal(ecs)

	if err == nil {
		return string(buf)
	} else {
		return err.Error()
	}
}

func (ecs *Ecs) SetMeta(meta evidence.Meta) {
	ecs.Labels["case"] = meta.Name
	ecs.Labels["Filters"] = strings.Join(meta.Filters, " > ")

	ecs.Timestamp = meta.Bagged.UTC()

	ecs.File.Path = meta.Path
	ecs.File.Size = meta.Size
	ecs.File.Mtime = meta.Modified.UTC()
	ecs.File.Hash.Sha256 = fmt.Sprintf("%x", meta.Hash)

	ecs.User.Name = meta.User.Username
	ecs.User.FullName = meta.User.Name
}

func (ecs *Ecs) AddLine(nr, grp int, str string) {
	ecs.Message += fmt.Sprintf("%d:%d: %s\n", nr, grp, str)
}
