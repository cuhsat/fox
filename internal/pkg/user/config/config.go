package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/cuhsat/fox/configs"
	"github.com/cuhsat/fox/internal/pkg/user"
)

func Get() *viper.Viper {
	return viper.GetViper()
}

func Load(flg *pflag.FlagSet) {
	cfg := Get()

	// setup default config
	_ = cfg.ReadConfig(strings.NewReader(configs.Default))

	// setup command line flags
	_ = cfg.BindPFlag("ai.model", flg.Lookup("model"))
	_ = cfg.BindPFlag("ai.num_ctx", flg.Lookup("num-ctx"))
	_ = cfg.BindPFlag("ai.temp", flg.Lookup("temp"))
	_ = cfg.BindPFlag("ai.topp", flg.Lookup("topp"))
	_ = cfg.BindPFlag("ai.topk", flg.Lookup("topk"))
	_ = cfg.BindPFlag("ai.seed", flg.Lookup("seed"))
	_ = cfg.BindPFlag("ui.theme", flg.Lookup("theme"))

	// setup default values
	cfg.SetDefault("ui.state.n", true)
	cfg.SetDefault("ui.state.w", false)
	cfg.SetDefault("ui.state.t", false)

	// setup environment
	cfg.AutomaticEnv()
	cfg.SetEnvPrefix("FOX")
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	user.LoadConfig(cfg, "foxrc")
}

func Save() {
	cfg := Get()

	user.SaveConfig(cfg, "foxrc")
}
