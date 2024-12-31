package main

import (
	"os"
	"testing"

	parser "layered_stack_cli/internal"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary test config file
	testConfig := `
name: test-stack
ssh_keys:
  - name: test-key
    provider: aws
    public_key: "ssh-rsa AAAA... test@example.com"
`
	tmpfile, err := os.CreateTemp("", "layered_stack*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(testConfig)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading the config
	config, err := parser.LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the loaded configuration
	if config.Name != "test-stack" {
		t.Errorf("Expected stack name 'test-stack', got '%s'", config.Name)
	}

	if len(config.SSHKeys) != 1 {
		t.Errorf("Expected 1 SSH key, got %d", len(config.SSHKeys))
	}

	if config.SSHKeys[0].Provider != "aws" {
		t.Errorf("Expected provider 'aws', got '%s'", config.SSHKeys[0].Provider)
	}
}
