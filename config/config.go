package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

const EnvPrefix = "2PC"

type Config struct {
	RestServer struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"rest_server"`

	SCS  GRPCClient `mapstructure:"scs"`
	SMS  GRPCClient `mapstructure:"sms"`
	OMS  GRPCClient `mapstructure:"oms"`
	IIMS GRPCClient `mapstructure:"iims"`
}

type GRPCClient struct {
	Address  string        `mapstructure:"address"`
	Insecure bool          `mapstructure:"insecure"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Tries    int           `mapstructure:"tries"`
}

func Get() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AddConfigPath("./config/")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	for _, k := range v.AllKeys() {
		val := v.GetString(k)
		v.Set(k, os.ExpandEnv(val))
	}

	var cfg *Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.RestServer.Port == "" {
		cfg.RestServer.Port = "8080"
	}

	return cfg, nil
}
