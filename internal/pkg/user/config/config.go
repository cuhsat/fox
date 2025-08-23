package config

import (
	"errors"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/cuhsat/fox/configs"
	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/user"
)

const Filename = ".foxrc"

func Get() *viper.Viper {
	return viper.GetViper()
}

func Load(flg *pflag.FlagSet) {
	cfg := Get()

	// setup config file
	cfg.SetConfigPermissions(0600)
	cfg.SetConfigName(Filename)
	cfg.SetConfigType("toml")
	cfg.AddConfigPath("$HOME")

	// setup command line flags
	cfg.BindPFlag("ai.model", flg.Lookup("model"))
	cfg.BindPFlag("ai.num_ctx", flg.Lookup("num-ctx"))
	cfg.BindPFlag("ai.temp", flg.Lookup("temp"))
	cfg.BindPFlag("ai.topp", flg.Lookup("topp"))
	cfg.BindPFlag("ai.topk", flg.Lookup("topk"))
	cfg.BindPFlag("ai.seed", flg.Lookup("seed"))

	cfg.BindPFlag("ui.theme", flg.Lookup("theme"))
	cfg.SetDefault("ui.state.n", true)
	cfg.SetDefault("ui.state.w", false)
	cfg.SetDefault("ui.state.t", false)

	// setup environment
	cfg.AutomaticEnv()
	cfg.SetEnvPrefix("FOX")
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := cfg.ReadInConfig()

	if err != nil {
		var e viper.ConfigFileNotFoundError

		if errors.Is(err, &e) {
			cfg.ReadConfig(strings.NewReader(configs.Config))
		} else {
			sys.Panic(err)
		}
	}
}

func Save() {
	path := Get().ConfigFileUsed()

	if len(path) == 0 {
		_, path = user.File(Filename)
	}

	err := Get().WriteConfigAs(path)

	if err != nil {
		sys.Error(err)
	}
}
