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

type ViewConfig struct {
	Name string   `mapstructure:"name"`
	Dirs []string `mapstructure:"dirs"`
	Type string   `mapstructure:"type"` // "alias" or "doc"
}

type Config struct {
	Views []ViewConfig `mapstructure:"views"`
	Theme Theme        `mapstructure:"theme"`
}

func LoadConfig() Config {
	home, _ := os.UserHomeDir()
	
	// Set default views
	defaultViews := []ViewConfig{
		{
			Name: "Aliases",
			Type: "alias",
			Dirs: []string{filepath.Join(home, "dotfiles", "scripts")},
		},
		{
			Name: "Docs",
			Type: "doc",
			Dirs: []string{"./docs", filepath.Join(home, ".local", "share", "shortcuts-tui", "docs")},
		},
	}

	viper.SetDefault("views", defaultViews)
	// Catppuccin Mocha Defaults
	viper.SetDefault("theme.primary", "#a6e3a1")   // Green
	viper.SetDefault("theme.secondary", "#6c7086") // Overlay0
	viper.SetDefault("theme.text", "#cdd6f4")      // Text

	// Setup config file search paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(home, ".config", "shortcuts"))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Error reading config file:", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("SHORTCUTS")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Unable to decode config into struct:", err)
	}

	return cfg
}
