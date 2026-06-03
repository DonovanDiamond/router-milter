package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Protocol   string `yaml:"protocol"`
	Address    string `yaml:"address"`
	ScriptPath string `yaml:"script_path"`
}

var protocol = flag.String("proto", "", "Protocol family (unix or tcp)")
var address = flag.String("addr", "", "Bind to address or unix domain socket")
var scriptPath = flag.String("script", "", "Path to the script.")
var configPath = flag.String("config", "config.yaml", "Path to configuration file (yaml)")

func (config *Config) LoadFromFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(raw, &config); err != nil {
		return err
	}
	return nil
}

func (config *Config) LoadFromFlags() {
	if *protocol != "" {
		config.Protocol = *protocol
	}
	if *address != "" {
		config.Address = *address
	}
	if *scriptPath != "" {
		config.ScriptPath = *scriptPath
	}
}

func (config *Config) Validate() error {
	if config.Protocol != "unix" && config.Protocol != "tcp" {
		return fmt.Errorf("invalid protocol: '%s', should be 'unix' or 'tcp'", config.Protocol)
	}
	if config.Address == "" {
		return fmt.Errorf("missing bind address")
	}
	if config.ScriptPath == "" {
		return fmt.Errorf("missing script path")
	}
	return nil
}

func LoadConfig() (Config, error) {
	flag.Parse()
	config := Config{}
	fileErr := config.LoadFromFile(*configPath)
	config.LoadFromFlags()
	validateErr := config.Validate()

	// the below is done so that if you do not specify a valid config file,
	// but you specify valid flags, ignore the invalid config file and just
	// use the flags.

	// if the config is valid, even if fileErr != nil, return the valid config
	if validateErr == nil {
		return config, nil
	}
	// otherwise return file error if it exists
	if fileErr != nil {
		return config, fmt.Errorf("failed to load config from '%s': %v", *configPath, fileErr)
	}
	// or return validate error if there is no file error
	return config, fmt.Errorf("invalid config: %v", validateErr)
}
