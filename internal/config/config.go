package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	ScriptsDir string
	DocsDir    string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadConfig() Config {
	home, _ := os.UserHomeDir()
	return Config{
		ScriptsDir: getEnv("SHORTCUTS_SCRIPTS_DIR", filepath.Join(home, "dotfiles", "scripts")),
		DocsDir:    getEnv("SHORTCUTS_DOCS_DIR", "./docs"),
	}
}
