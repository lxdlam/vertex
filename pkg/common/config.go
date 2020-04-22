package common

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

// Config is a simple struct that contains all necessary options.
type Config struct {
	LogPath  string `toml:"log_path"`
	LogLevel string `toml:"log_level"`
}

// NewConfig will return a config instance with default value
func NewConfig() *Config {
	return &Config{
		LogPath:  "./log/vertex.log",
		LogLevel: "INFO",
	}
}

// Parse will parse a config toml
func (c *Config) Parse(path string) error {
	tree, err := toml.LoadFile(path)

	if err != nil {
		return fmt.Errorf("parse config from file failed%w", err)
	}

	err = tree.Unmarshal(c)

	if err != nil {
		return fmt.Errorf("parse config from file failed%w", err)
	}

	return nil
}
