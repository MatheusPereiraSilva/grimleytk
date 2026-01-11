package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads a grimley.yaml file and parses it into a Config struct
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml file '%s': %w", path, err)
	}

	return &cfg, nil
}
