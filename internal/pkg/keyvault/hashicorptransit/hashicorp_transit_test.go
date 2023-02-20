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
	hte = NewHashicorpTransitEngine(config.Configuration().VaultURL, config.Configuration().VaultUsername, config.Configuration().VaultPassword, config.Configuration().HashicorpTransitPath)
	m.Run()
}

func TestEncryptDecrypt(t *testing.T) {
	ctx := context.TODO()
	text := "hello-key-vault"
	keyName := "9999-demo"

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
	keyName := "9999-demo"
	for i := 0; i < b.N; i++ {
		hte.Encrypt(ctx, text, keyName, "")
	}
}

func BenchmarkDecrypt(b *testing.B) {
	ctx := context.TODO()
	keyName := "demo"
	encrypted := "vault:v4:EejDgr9PGOGk/aLBzxdnOFvOwGnkolfS4CTsqau11V22xYeNpbN3AqdWXUbLBu8Hp+UAEL2imETwBVILgnBGFhjYUXjgPbq4OoR0T/i/12FgZYo5BeDk7eZbW40w68uCsAcnEgB3CA4fBRdEdooLCAsRS34SgWnhmLuP/iObchTulAh9HYZfJ1pW34dULiCKd/Kb0gqwG9tBly6RYHp45Rxi0aFF+Dz5TKcUnK0m9CAqRIYPAc0Zl0PgrAPQ0IyEtkU5BoUBmANnIAY1rdI4zraDVb6+u7K2IWJAY4mPEmn6oy4wS15x+7WN254OFg9wbnjpuf9yjOHITodDDsSnug=="
	for i := 0; i < b.N; i++ {
		hte.Decrypt(ctx, encrypted, keyName, "")
	}
}
