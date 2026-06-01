package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Protocol        string   `yaml:"protocol"`
	Address         string   `yaml:"address"`
	RejectFrom      []string `yaml:"reject_from"`
	RejectTo        []string `yaml:"reject_to"`
	RejectToRegex   []string `yaml:"reject_to_regex"`
	RejectToSha256  []string `yaml:"reject_to_sha256"`
}

func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(raw, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
