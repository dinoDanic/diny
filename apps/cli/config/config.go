package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Tone string
type Length string

const (
	Professional Tone = "professional"
	Casual       Tone = "casual"
	Friendly     Tone = "friendly"
)

const (
	Short  Length = "short"
	Normal Length = "normal"
	Long   Length = "long"
)

type Config struct {
	Theme  string       `mapstructure:"theme"`
	Commit CommitConfig `mapstructure:"commit"`
}

type CommitConfig struct {
	Conventional bool   `mapstructure:"conventional"`
	Emoji        bool   `mapstructure:"emoji"`
	Tone         Tone   `mapstructure:"tone"`
	Length       Length `mapstructure:"length"`
}

var cfg *Config

func Load(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}

		viper.AddConfigPath(home + "/.config/diny")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetDefault("theme", "catppuccin")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg = &c
	return &c, nil
}

func Get() *Config {
	if cfg == nil {
		return &Config{}
	}
	return cfg
}
