package keyvault

import "context"

type EncryptionResult struct {
	Result     string
	KeyVersion string
}

type KeyVault interface {
	// Encrypt will encrypt plaint text to base64 string.
	// If keyVersion is empty, it will use latest key version.
	Encrypt(ctx context.Context, plaintext string, keyName string, keyVersion string) (result EncryptionResult, err error)

	// Decrypt will decrtypt encrypted text based on key and key version.
	// If keyVersion is empty, it will use latest key version.
	Decrypt(ctx context.Context, encrypted string, keyName string, keyVersion string) (decrypted string, err error)
}
