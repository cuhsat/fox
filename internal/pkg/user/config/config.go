package config

import (
	"github.com/spf13/viper"

	"github.com/cuhsat/fox/configs"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user"
)

const (
	Filename = ".foxrc"
)

var (
	Default = configs.Config
)

func Save() {
	_, path := user.File(Filename)

	err := viper.WriteConfigAs(path)

	if err != nil {
		sys.Error(err)
	}
}
