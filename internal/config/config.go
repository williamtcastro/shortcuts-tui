package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Theme struct {
	PrimaryColor   string `mapstructure:"primary"`
	SecondaryColor string `mapstructure:"secondary"`
	TextColor      string `mapstructure:"text"`
}

type Config struct {
	ScriptsDirs []string `mapstructure:"scripts_dirs"`
	DocsDirs    []string `mapstructure:"docs_dirs"`
	Theme       Theme    `mapstructure:"theme"`
}

func LoadConfig() Config {
	home, _ := os.UserHomeDir()
	
	// Set defaults
	viper.SetDefault("scripts_dirs", []string{filepath.Join(home, "dotfiles", "scripts")})
	viper.SetDefault("docs_dirs", []string{"./docs", filepath.Join(home, ".local", "share", "shortcuts-tui", "docs")})
	viper.SetDefault("theme.primary", "#25A065")
	viper.SetDefault("theme.secondary", "#545454")
	viper.SetDefault("theme.text", "#FFFDF5")

	// Setup config file search paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(home, ".config", "shortcuts"))
	viper.AddConfigPath(".")

	// Attempt to read config (ignore error if not found, rely on defaults)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Error reading config file:", err)
		}
	}

	// Environment variable overrides
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SHORTCUTS")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Unable to decode config into struct:", err)
	}

	return cfg
}
