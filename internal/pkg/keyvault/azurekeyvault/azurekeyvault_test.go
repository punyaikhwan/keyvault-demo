package azurekeyvault

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
	az = NewAzureKeyVault(config.Configuration().VaultURL)
	m.Run()
}

func TestEncryptDecrypt(t *testing.T) {
	ctx := context.TODO()
	text := "hello-key-vault"
	keyName := "demo"

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
	encrypted := "BXxsSXh+VomQhBrSPyGI46cK2Y9lF+NQz5uD2uW8x8vSsTrGaroNZV+TJC9dhERVpX7r51oTsbRj/7hRrhrWf7mQ4gSOaSBM7IjkSibC+UHR9RWvt9NVIRuoueTrkjlpGn1E5wlzji7q+vji+KVa1JdUBt+dxqd1Fnw6xXgQFmrwyHbr1zkhBHnBGKA4DNhJ2ujnIo1h4EomXkA1EY5z7hqvF8EpQhuuAMTVRKtNhycSKk3K+0pqaH4h71bWqTVJ8vgSHnx+N9g2GJjtLwjrekVLnMFIJWjU6ZONuExMS26LsnR52TgtKDk9FDEnr4IZUka4WiJin2/8KrzcgQcgdQ=="
	keyName := "demo"

	decrypted, err := az.Decrypt(ctx, encrypted, keyName, "")
	assert.Nil(t, err)
	fmt.Printf("decrypted: %s\n", decrypted)
}

func BenchmarkEncrypt(b *testing.B) {
	ctx := context.TODO()
	text := "hello-key-vault"
	keyName := "demo"
	for i := 0; i < b.N; i++ {
		az.Encrypt(ctx, text, keyName, "")
	}
}

func BenchmarkDecrypt(b *testing.B) {
	ctx := context.TODO()
	keyName := "demo"
	encrypted := "fCmsij5S0GzpfXYQ2e3HSwQuTKS0RHQ/lRWcrOiyPSAVPKyPdQ3qlGSfVaX/0/bRAo63bxkckWYUH59uv3G1Y/kNUfMi7H1SL3q9hDe/PW+DYP0S8nI7phHJKVVMKNfO+pvm18k4Y3KEJRDw9cG3YJGFy5U+MSPYMUp127KMEYhI5XKmf/YI0I/qijiC8Xm8lWMq3EMLYxyqH8fBMneshHzT7JcN4nRD/22um9x9w5uEsXcbWQhikn5KSXpXEM8qKyrxvL9ZQ1YPxrIPJ7jLyw/6S74K6v5axefbmu+mtvqksu2zspoUSx3gAopPpjDuT/rTE0M7dilMNSW2Vq5JMg=="
	for i := 0; i < b.N; i++ {
		az.Decrypt(ctx, encrypted, keyName, "")
	}
}
