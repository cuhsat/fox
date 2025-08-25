package user

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/cuhsat/fox/internal/pkg/sys"
)

func LoadConfig(cfg *viper.Viper, name string) bool {
	var ok bool

	cfg.SetConfigPermissions(0600)
	cfg.SetConfigType("toml")
	cfg.SetConfigName(name)

	for _, path := range []string{
		// system config
		"/etc/fox",
		// local config
		"/usr/local/etc/fox",
		// user config
		"$HOME/.config/fox",
	} {
		cfg.AddConfigPath(path)

		if cfg.MergeInConfig() == nil {
			ok = true
		}
	}

	return ok
}

func SaveConfig(cfg *viper.Viper, name string) bool {
	path := sys.Config(name)

	err := os.MkdirAll(filepath.Dir(path), 0700)

	if err != nil {
		return false
	}

	err = cfg.WriteConfigAs(path)

	if err != nil {
		return false
	}

	return true
}
