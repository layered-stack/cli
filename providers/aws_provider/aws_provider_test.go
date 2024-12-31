package aws_provider

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// MockEC2Client implements EC2ClientAPI for testing
type MockEC2Client struct {
	existingKeys     []types.KeyPairInfo
	importedKeyPairs []ec2.ImportKeyPairInput
}

func (m *MockEC2Client) DescribeKeyPairs(ctx context.Context, params *ec2.DescribeKeyPairsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeKeyPairsOutput, error) {
	return &ec2.DescribeKeyPairsOutput{
		KeyPairs: m.existingKeys,
	}, nil
}

func (m *MockEC2Client) ImportKeyPair(ctx context.Context, params *ec2.ImportKeyPairInput, optFns ...func(*ec2.Options)) (*ec2.ImportKeyPairOutput, error) {
	m.importedKeyPairs = append(m.importedKeyPairs, *params)
	return &ec2.ImportKeyPairOutput{}, nil
}

func TestSyncSshKeys(t *testing.T) {
	tests := []struct {
		name         string
		existingKeys []types.KeyPairInfo
		inputKeys    []map[string]interface{}
		wantImports  int
	}{
		{
			name: "new key should be imported",
			existingKeys: []types.KeyPairInfo{
				{KeyName: aws.String("existing-key")},
			},
			inputKeys: []map[string]interface{}{
				{
					"name":       "new-key",
					"public_key": "ssh-rsa AAAAB...",
				},
			},
			wantImports: 1,
		},
		{
			name: "existing key should not be imported",
			existingKeys: []types.KeyPairInfo{
				{KeyName: aws.String("existing-key")},
			},
			inputKeys: []map[string]interface{}{
				{
					"name":       "existing-key",
					"public_key": "ssh-rsa AAAAB...",
				},
			},
			wantImports: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockEC2Client{
				existingKeys: tt.existingKeys,
			}

			err := SyncSshKeys(tt.inputKeys, mockClient)
			if err != nil {
				t.Errorf("SyncSshKeys() error = %v", err)
				return
			}

			if got := len(mockClient.importedKeyPairs); got != tt.wantImports {
				t.Errorf("SyncSshKeys() imported %d keys, want %d", got, tt.wantImports)
			}
		})
	}
}
