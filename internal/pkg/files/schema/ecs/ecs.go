// Package ecs specification: https://www.elastic.co/docs/reference/ecs/ecs-field-reference
package ecs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cuhsat/fox/internal"
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
		Path  string    `json:"path"`
		Size  int64     `json:"size"`

		Hash struct {
			Sha256 string `json:"sha256"`
		} `json:"hash"`
	} `json:"file"`

	User struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"user"`
}

func New() *Ecs {
	ecs := new(Ecs)
	ecs.Labels = make(map[string]string)

	ecs.Ecs.Version = Version

	ecs.Agent.Type = info.Product
	ecs.Agent.Version = info.Version[1:]

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

func (ecs *Ecs) Headers() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

func (ecs *Ecs) SetMeta(meta evidence.Meta) {
	ecs.Labels["case"] = meta.Name
	ecs.Labels["filters"] = strings.Join(meta.Filters, " > ")

	ecs.Timestamp = meta.Bagged.UTC()

	ecs.File.Path = meta.Path
	ecs.File.Size = meta.Size
	ecs.File.Mtime = meta.Modified.UTC()
	ecs.File.Hash.Sha256 = fmt.Sprintf("%x", meta.Hash)

	ecs.User.Name = meta.User.Username
	ecs.User.FullName = meta.User.Name
}

func (ecs *Ecs) AddLine(_, _ int, str string) {
	ecs.Message += fmt.Sprintf("%s\n", str)
}
