package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Proxy  ProxyConfig `mapstructure:"proxy"`
	Routes []Route     `mapstructure:"routes"`
}

type ProxyConfig struct {
	Bind string `mapstructure:"bind"`
}

type Route struct {
	Prefix string `mapstructure:"prefix"`
	Target string `mapstructure:"target"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("goproxi")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
