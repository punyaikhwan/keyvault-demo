package hashicorptransit

import (
	"context"
	"encoding/base64"
	"fmt"
	"keyvault-demo/internal/pkg/hashicorp"
	"keyvault-demo/internal/pkg/keyvault"
)

var algo = hashicorp.JSONWebKeyEncryptionAlgorithmRSAOAEP256

type hashicorpTransitEngine struct {
	vaultURL string
	username string
	password string
	client   *hashicorp.Client
}

func (hte *hashicorpTransitEngine) getClient() *hashicorp.Client {
	if hte.client == nil {
		fmt.Println("Creating Hashicorp client...")
		cred := hashicorp.NewUsernamePasswordCredential(hte.username, hte.password)

		hte.client = hashicorp.NewClient(hte.vaultURL, cred)
	}

	return hte.client
}

func NewHashicorpTransitEngine(vaultURL string, username string, password string) keyvault.KeyVault {
	return &hashicorpTransitEngine{vaultURL, username, password, nil}
}

func (hte *hashicorpTransitEngine) Encrypt(ctx context.Context, plaintext string, keyName string, keyVersion string) (result keyvault.EncryptionResult, err error) {
	encryptRes, err := hte.getClient().Encrypt(ctx, keyName, keyVersion, hashicorp.KeyOperationsParameters{
		Algorithm: &algo,
		Value:     []byte(plaintext),
	})

	if err != nil {
		return result, err
	}

	result.Result = encryptRes.Base64Result
	result.KeyVersion = encryptRes.Version
	return result, nil
}

func (hte *hashicorpTransitEngine) Decrypt(ctx context.Context, encryptedBase64 string, keyName string, keyVersion string) (decryptedBase64 string, err error) {
	decryptedRes, err := hte.getClient().Decrypt(ctx, keyName, hashicorp.KeyOperationsParameters{
		Algorithm:   &algo,
		Base64Value: encryptedBase64,
	})

	if err != nil {
		return "", err
	}

	result, err := base64.StdEncoding.DecodeString(decryptedRes.Base64Result)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
