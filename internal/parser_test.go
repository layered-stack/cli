package parser

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary YAML file for testing
	content := []byte(`
name: "Test Config"
ssh_keys:
  - name: "test-key"
    provider: "aws"
    public_key: "ssh-rsa AAAA..."
`)
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test successful config loading
	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Errorf("LoadConfig() error = %v, want nil", err)
	}

	// Verify the parsed content
	if config.Name != "Test Config" {
		t.Errorf("config.Name = %v, want %v", config.Name, "Test Config")
	}

	// Verify SSH keys
	if len(config.SSHKeys) != 1 {
		t.Errorf("len(config.SSHKeys) = %v, want 1", len(config.SSHKeys))
	}
	if config.SSHKeys[0].Name != "test-key" {
		t.Errorf("config.SSHKeys[0].Name = %v, want %v", config.SSHKeys[0].Name, "test-key")
	}
	if config.SSHKeys[0].Provider != "aws" {
		t.Errorf("config.SSHKeys[0].Provider = %v, want %v", config.SSHKeys[0].Provider, "aws")
	}

	// Test loading non-existent file
	_, err = LoadConfig("non-existent.yaml")
	if err == nil {
		t.Error("LoadConfig() error = nil, want error for non-existent file")
	}
}
