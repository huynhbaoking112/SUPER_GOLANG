package utils

import (
	"go-backend-v2/pkg/utils"
	"testing"
)

func TestEncryptDecryptToken(t *testing.T) {
	// Test data - exactly 32 bytes key
	testKey := "MySecretEncryptionKey32BytesKey!"
	testToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// Test successful encryption and decryption
	encrypted, err := utils.EncryptToken(testToken, testKey)
	if err != nil {
		t.Fatalf("Failed to encrypt token: %v", err)
	}

	if encrypted == "" {
		t.Fatal("Encrypted token should not be empty")
	}

	if encrypted == testToken {
		t.Fatal("Encrypted token should be different from original token")
	}

	// Test decryption
	decrypted, err := utils.DecryptToken(encrypted, testKey)
	if err != nil {
		t.Fatalf("Failed to decrypt token: %v", err)
	}

	if decrypted != testToken {
		t.Fatalf("Decrypted token doesn't match original. Expected: %s, Got: %s", testToken, decrypted)
	}
}

func TestEncryptToken_InvalidKeyLength(t *testing.T) {
	testToken := "test-token"

	tests := []struct {
		name string
		key  string
	}{
		{"Too short key", "shortkey"},
		{"Too long key", "ThisKeyIsTooLongForAES256EncryptionAndShouldFail"},
		{"Empty key", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utils.EncryptToken(testToken, tt.key)
			if err == nil {
				t.Fatalf("Expected error for invalid key length, but got none")
			}
		})
	}
}

func TestDecryptToken_InvalidKeyLength(t *testing.T) {
	testEncrypted := "dGVzdA=="

	tests := []struct {
		name string
		key  string
	}{
		{"Too short key", "shortkey"},
		{"Too long key", "ThisKeyIsTooLongForAES256EncryptionAndShouldFail"},
		{"Empty key", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utils.DecryptToken(testEncrypted, tt.key)
			if err == nil {
				t.Fatalf("Expected error for invalid key length, but got none")
			}
		})
	}
}

func TestDecryptToken_InvalidEncryptedData(t *testing.T) {
	testKey := "MySecretEncryptionKey32BytesKey!"

	tests := []struct {
		name      string
		encrypted string
	}{
		{"Invalid base64", "invalid-base64-data!!!"},
		{"Too short data", "dGVzdA=="},
		{"Empty data", ""},
		{"Wrong encrypted data", "d2VpcmRkYXRhVGhhdENhbnROb3REZWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utils.DecryptToken(tt.encrypted, testKey)
			if err == nil {
				t.Fatalf("Expected error for invalid encrypted data, but got none")
			}
		})
	}
}

func TestValidateEncryptionKey(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		expectErr bool
	}{
		{"Valid key", "MySecretEncryptionKey32BytesKey!", false},
		{"Too short", "shortkey", true},
		{"Too long", "ThisKeyIsTooLongForAES256EncryptionAndShouldFailBecauseItExceedsThirtyTwoBytes", true},
		{"Empty", "", true},
		{"Exactly 32 bytes", "12345678901234567890123456789012", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateEncryptionKey(tt.key)
			if tt.expectErr && err == nil {
				t.Fatal("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestEncryptionRandomness(t *testing.T) {
	testKey := "MySecretEncryptionKey32BytesKey!"
	testToken := "same-token-content"

	// Encrypt the same token multiple times
	encrypted1, err1 := utils.EncryptToken(testToken, testKey)
	encrypted2, err2 := utils.EncryptToken(testToken, testKey)

	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to encrypt tokens: %v, %v", err1, err2)
	}

	// Results should be different (due to random nonce)
	if encrypted1 == encrypted2 {
		t.Fatal("Encrypted tokens should be different due to random nonce")
	}

	// But both should decrypt to the same original token
	decrypted1, err1 := utils.DecryptToken(encrypted1, testKey)
	decrypted2, err2 := utils.DecryptToken(encrypted2, testKey)

	if err1 != nil || err2 != nil {
		t.Fatalf("Failed to decrypt tokens: %v, %v", err1, err2)
	}

	if decrypted1 != testToken || decrypted2 != testToken {
		t.Fatal("Both decrypted tokens should match the original")
	}
}
