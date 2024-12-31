package parser

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SSHKey struct {
	Name      string `yaml:"name"`
	Provider  string `yaml:"provider"`
	PublicKey string `yaml:"public_key"`
}

type Config struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	SSHKeys     []SSHKey `yaml:"ssh_keys"`
}

// LoadConfig loads and parses the YAML configuration file
func LoadConfig(filepath string) (*Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Parse YAML into Config struct
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
