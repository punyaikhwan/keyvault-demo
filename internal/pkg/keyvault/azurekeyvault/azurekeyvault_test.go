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
	encrypted := "dA8wWf6BiK4BYs/xgOP0iVTWIgeomUUubBbn/ihk4LYsAMa/z6LqFmdhiBz04i3ky3jyIpUbQK1tBGskfApCFJ+4nzwAQN1EzYKWe0QRT+ixm5KqHB1G79UROJtSiTjurxY4Epkyp665QcNSOatj+QCXEV4eSO85do6YM2JJR3L4pU5PnJ/3MzbBeKFwEnau90icQ3CSsYHuH3Ff65bgzNdoWbW5aV+F771JwR74H/TE3YsI5XzFvvDXlg8c1TKCpOj5wamTPP0PU2DPgM/MyhVTaMqyNnbGEqxEk3zTnt/ZRfWGSmwQxAPOP7ssK7BII+0JUMAMQv35Lgw0SpikhQ=="
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
	encrypted := "dA8wWf6BiK4BYs/xgOP0iVTWIgeomUUubBbn/ihk4LYsAMa/z6LqFmdhiBz04i3ky3jyIpUbQK1tBGskfApCFJ+4nzwAQN1EzYKWe0QRT+ixm5KqHB1G79UROJtSiTjurxY4Epkyp665QcNSOatj+QCXEV4eSO85do6YM2JJR3L4pU5PnJ/3MzbBeKFwEnau90icQ3CSsYHuH3Ff65bgzNdoWbW5aV+F771JwR74H/TE3YsI5XzFvvDXlg8c1TKCpOj5wamTPP0PU2DPgM/MyhVTaMqyNnbGEqxEk3zTnt/ZRfWGSmwQxAPOP7ssK7BII+0JUMAMQv35Lgw0SpikhQ=="
	for i := 0; i < b.N; i++ {
		az.Decrypt(ctx, encrypted, keyName, "")
	}
}
