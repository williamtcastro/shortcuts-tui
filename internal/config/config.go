package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Theme struct {
	Primary   string `mapstructure:"primary"`
	Secondary string `mapstructure:"secondary"`
	Text      string `mapstructure:"text"`
	Accent    string `mapstructure:"accent"`
	Mauve     string `mapstructure:"mauve"`
	Flamingo  string `mapstructure:"flamingo"`
}

type ViewConfig struct {
	Name string   `mapstructure:"name"`
	Dirs []string `mapstructure:"dirs"`
	Type string   `mapstructure:"type"` // "alias" or "doc"
}

type Config struct {
	Views      []ViewConfig `mapstructure:"views"`
	Theme      Theme        `mapstructure:"theme"`
	Pagination string       `mapstructure:"pagination"` // "numeric" or "dots"
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
	viper.SetDefault("pagination", "numeric")
	
	// Catppuccin Mocha Defaults
	viper.SetDefault("theme.primary", "#a6e3a1")   // Green
	viper.SetDefault("theme.secondary", "#6c7086") // Overlay0
	viper.SetDefault("theme.text", "#cdd6f4")      // Text
	viper.SetDefault("theme.accent", "#f9e2af")    // Yellow
	viper.SetDefault("theme.mauve", "#cba6f7")     // Mauve
	viper.SetDefault("theme.flamingo", "#f2cdcd")  // Flamingo

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
