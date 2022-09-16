package hashicorptransit

import (
	"context"
	"fmt"
	"keyvault-demo/config"
	"keyvault-demo/internal/pkg/keyvault"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var hte keyvault.KeyVault

func TestMain(m *testing.M) {
	config.ReadConfig("../../../../.env")
	hte = NewHashicorpTransitEngine(config.Configuration().VaultURL, config.Configuration().VaultUsername, config.Configuration().VaultPassword)
	m.Run()
}

func TestEncryptDecrypt(t *testing.T) {
	ctx := context.TODO()
	text := "hello-key-vault"
	keyName := "demo"

	start := time.Now()
	encrypted, err := hte.Encrypt(ctx, text, keyName, "")
	fmt.Printf("Time to encrypt: %d ms", time.Now().UnixMilli()-start.UnixMilli())
	assert.Nil(t, err)

	fmt.Printf("\nencrypted: %s", encrypted.Result)

	start = time.Now()
	decrypted, err := hte.Decrypt(ctx, encrypted.Result, keyName, encrypted.KeyVersion)
	fmt.Printf("\nTime to decrypt: %d ms", time.Now().UnixMilli()-start.UnixMilli())
	assert.Nil(t, err)
	assert.Equal(t, text, decrypted)

	fmt.Printf("\ndecrypted: %s\n", decrypted)
}

func BenchmarkEncrypt(b *testing.B) {
	ctx := context.TODO()
	text := "hello-key-vault"
	keyName := "demo"
	for i := 0; i < b.N; i++ {
		hte.Encrypt(ctx, text, keyName, "")
	}
}

func BenchmarkDecrypt(b *testing.B) {
	ctx := context.TODO()
	keyName := "demo"
	encrypted := "vault:v1:j0cIunUd3fRcWlJAx6i2zqV5CPue1iSVeSIQiuDHzukJf7kCoyaYW7X7rCJKE5vaR5bH8nLPw8gi5mWKIhueJNX6GeWBHRIPaDXNz2shaLxFKuN1QlLe6GuRR58wjlIikhRbQfzmG9Je/59pQD63iiFU1fvzk/QfAhkNcl5CvImntbvRaodQKzbM3AOwJ9VOhPitK1dG9HgKDfmdlsxxLQhDjoiO9iEWZHVrz5bxPLbt+7eeRMLrU84sf4NPkICb4Y69IkBuurGqTOi4XMcRwsoqT2asqjtdbCnWSnN68i1AuYSkB4mp1jnnJ2TTs48w70lLJ+bnKqDEu6sIjxhO9A=="
	for i := 0; i < b.N; i++ {
		hte.Decrypt(ctx, encrypted, keyName, "")
	}
}
