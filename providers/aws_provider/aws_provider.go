package aws_provider

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2ClientAPI defines the interface for EC2 client operations we need
type EC2ClientAPI interface {
	DescribeKeyPairs(ctx context.Context, params *ec2.DescribeKeyPairsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeKeyPairsOutput, error)
	ImportKeyPair(ctx context.Context, params *ec2.ImportKeyPairInput, optFns ...func(*ec2.Options)) (*ec2.ImportKeyPairOutput, error)
}

// newEC2Client creates a new EC2 client using environment variables
func newEC2Client() (EC2ClientAPI, error) {
	// Check required environment variables
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" || region == "" {
		return nil, fmt.Errorf("missing required AWS environment variables: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION")
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	return ec2.NewFromConfig(cfg), nil
}

// SyncSshKeys ensures SSH keys from the config exist in AWS
func SyncSshKeys(sshKeys []map[string]interface{}, client EC2ClientAPI) error {
	if client == nil {
		var err error
		client, err = newEC2Client()
		if err != nil {
			return fmt.Errorf("failed to create AWS client: %w", err)
		}
	}

	// Get existing key pairs
	existing, err := client.DescribeKeyPairs(context.Background(), &ec2.DescribeKeyPairsInput{})
	if err != nil {
		return fmt.Errorf("failed to get existing key pairs: %w", err)
	}

	for _, keyConfig := range sshKeys {
		keyName, ok := keyConfig["name"].(string)
		if !ok {
			return fmt.Errorf("ssh key name not found or not a string")
		}

		publicKey, ok := keyConfig["public_key"].(string)
		if !ok {
			return fmt.Errorf("public key not found or not a string for key '%s'", keyName)
		}

		// Check if key exists
		keyExists := false
		for _, existing := range existing.KeyPairs {
			if *existing.KeyName == keyName {
				keyExists = true
				fmt.Printf("SSH key '%s' already exists\n", keyName)
				break
			}
		}

		if !keyExists {
			fmt.Printf("Creating SSH key '%s'...\n", keyName)
			_, err := client.ImportKeyPair(context.Background(), &ec2.ImportKeyPairInput{
				KeyName:           aws.String(keyName),
				PublicKeyMaterial: []byte(publicKey),
			})
			if err != nil {
				return fmt.Errorf("failed to import key pair '%s': %w", keyName, err)
			}
		}
	}

	return nil
}
