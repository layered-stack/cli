package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	parser "layered_stack_cli/internal"
	aws_provider "layered_stack_cli/providers/aws_provider"
)

func init() {
	// Load .env file if it exists (silent fail if not found)
	_ = godotenv.Load()
}

func main() {
	config, err := parser.LoadConfig("layered_stack.yml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("- Layered Stack -")
	fmt.Printf("Found \"%s\" in layered_stack.yml\n", config.Name)

	// Process SSH keys by provider
	if len(config.SSHKeys) > 0 {
		awsKeys := make([]map[string]interface{}, 0)
		for _, key := range config.SSHKeys {
			if key.Provider == "aws" {
				awsKeys = append(awsKeys, map[string]interface{}{
					"name":       key.Name,
					"provider":   key.Provider,
					"public_key": key.PublicKey,
				})
			}
		}

		if len(awsKeys) > 0 {
			if err := aws_provider.SyncSshKeys(awsKeys, nil); err != nil {
				log.Fatalf("Failed to sync AWS SSH keys: %v", err)
			}
			fmt.Printf("Processed %d AWS SSH keys\n", len(awsKeys))
		}
	}
}
