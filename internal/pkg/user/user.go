package user

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/sys/fs"
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
	path := Config(name)

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

func TempDir(prefix string) string {
	tmp, err := os.MkdirTemp(Cache(), fmt.Sprintf("%s-*", prefix))

	if err != nil {
		sys.Panic(err)
	}

	return tmp
}

func TempFile(prefix string) *os.File {
	tmp, err := os.CreateTemp(Cache(), fmt.Sprintf("%s-*", prefix))

	if err != nil {
		sys.Panic(err)
	}

	return tmp
}

func Persist(name string) string {
	f, ok := fs.Open(name).(fs.File)

	if !ok {
		return name // regular file
	}

	t := TempFile("fox")

	_, err := t.WriteTo(f)

	if err != nil {
		sys.Error(err)
	}

	return f.Name()
}

func Config(name string) string {
	dir, err := os.UserHomeDir()

	if err != nil {
		sys.Panic(err)
	}

	return filepath.Join(dir, ".config", "fox", name)
}

func Cache() string {
	dir, err := os.UserHomeDir()

	if err != nil {
		sys.Panic(err)
	}

	tmp := filepath.Join(dir, ".cache", "fox")

	err = os.MkdirAll(tmp, 0700)

	if err != nil {
		sys.Panic(err)
	}

	return tmp
}
