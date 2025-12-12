package config

import (
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	InstallDir  string
	DefaultJava string
	CacheDir    string
}

// GetConfig returns the application configuration
func GetConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	jvtDir := filepath.Join(homeDir, ".jvt")

	return &Config{
		InstallDir:  filepath.Join(jvtDir, "versions"),
		CacheDir:    filepath.Join(jvtDir, "cache"),
		DefaultJava: "",
	}, nil
}

// EnsureDirectories creates necessary directories if they don't exist
func (c *Config) EnsureDirectories() error {
	dirs := []string{c.InstallDir, c.CacheDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
