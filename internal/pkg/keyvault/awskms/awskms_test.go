package awskms

import (
	"context"
	"fmt"
	"keyvault-demo/config"
	"keyvault-demo/internal/pkg/keyvault"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var az keyvault.KeyVault

func TestMain(m *testing.M) {
	config.ReadConfig("../../../../.env")
	az = NewAWSKeyVault(config.Configuration().VaultURL)
	m.Run()
}

func TestEncryptDecrypt(t *testing.T) {
	ctx := context.TODO()
	text := "hello-key-vault"
	keyName := "ed6d750e-a181-4a9e-a9f0-8bff2b32d7c7"

	start := time.Now()
	encrypted, err := az.Encrypt(ctx, text, keyName, "")
	fmt.Printf("Time to encrypt: %d ms", time.Now().UnixMilli()-start.UnixMilli())
	assert.Nil(t, err)

	fmt.Printf("\nencrypted: %s", encrypted.Result)

	start = time.Now()
	decrypted, err := az.Decrypt(ctx, encrypted.Result, keyName, encrypted.KeyVersion)
	fmt.Printf("\nTime to decrypt: %d ms", time.Now().UnixMilli()-start.UnixMilli())
	assert.Nil(t, err)
	assert.Equal(t, decrypted, text)

	fmt.Printf("\ndecrypted: %s", decrypted)
}

func TestDecrypt(t *testing.T) {
	ctx := context.TODO()
	encrypted := "AQICAHjlQV9R4FkJ30+g2Sr5g6lmpoTNgSDt5efuqQJXWgjzSAFJKJURj+e0p3JpHdGKqejVAAAAbTBrBgkqhkiG9w0BBwagXjBcAgEAMFcGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQM1fl01nftn7JOgt9WAgEQgCqXLcJOlgdhu3cdOCvacSRI9xa8jgZW/fYnCw6qBdZ5E29xvW7hKa2poTA="
	keyName := "ed6d750e-a181-4a9e-a9f0-8bff2b32d7c7"

	decrypted, err := az.Decrypt(ctx, encrypted, keyName, "")
	assert.Nil(t, err)
	fmt.Printf("decrypted: %s\n", decrypted)
}

func BenchmarkEncrypt(b *testing.B) {
	ctx := context.TODO()
	text := "hello-key-vault"
	keyName := "ed6d750e-a181-4a9e-a9f0-8bff2b32d7c7"
	for i := 0; i < b.N; i++ {
		az.Encrypt(ctx, text, keyName, "")
	}
}

func BenchmarkDecrypt(b *testing.B) {
	ctx := context.TODO()
	keyName := "ed6d750e-a181-4a9e-a9f0-8bff2b32d7c7"
	encrypted := "AQICAHjlQV9R4FkJ30+g2Sr5g6lmpoTNgSDt5efuqQJXWgjzSAFJKJURj+e0p3JpHdGKqejVAAAAbTBrBgkqhkiG9w0BBwagXjBcAgEAMFcGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQM1fl01nftn7JOgt9WAgEQgCqXLcJOlgdhu3cdOCvacSRI9xa8jgZW/fYnCw6qBdZ5E29xvW7hKa2poTA="
	for i := 0; i < b.N; i++ {
		az.Decrypt(ctx, encrypted, keyName, "")
	}
}
