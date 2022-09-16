package azurekeyvault

import (
	"context"
	"encoding/base64"
	"fmt"
	"keyvault-demo/internal/pkg/keyvault"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

var algo = azkeys.JSONWebKeyEncryptionAlgorithmRSAOAEP256

type azureKeyVault struct {
	vaultURL string
	client   *azkeys.Client
}

func (az *azureKeyVault) getClient() *azkeys.Client {
	if az.client == nil {
		fmt.Println("Creating Azure client...")
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatalf("failed to obtain a credential: %v", err)
		}

		az.client = azkeys.NewClient(az.vaultURL, cred, nil)
	}

	return az.client
}

func NewAzureKeyVault(vaultURL string) keyvault.KeyVault {
	return &azureKeyVault{vaultURL, nil}
}

func (az *azureKeyVault) Encrypt(ctx context.Context, plaintext string, keyName string, keyVersion string) (result keyvault.EncryptionResult, err error) {
	if keyVersion == "" {
		keyVersion, err = az.getLatestKeyVersion(ctx, keyName)
		if err != nil {
			return result, err
		}
	}
	encryptRes, err := az.getClient().Encrypt(ctx, keyName, keyVersion, azkeys.KeyOperationsParameters{
		Algorithm: &algo,
		Value:     []byte(plaintext),
	}, nil)

	if err != nil {
		return result, err
	}

	result.Result = base64.StdEncoding.EncodeToString(encryptRes.Result)
	result.KeyVersion = keyVersion
	return result, nil
}

func (az *azureKeyVault) Decrypt(ctx context.Context, encryptedBase64 string, keyName string, keyVersion string) (decryptedBase64 string, err error) {
	if keyVersion == "" {
		keyVersion, err = az.getLatestKeyVersion(ctx, keyName)
		if err != nil {
			return encryptedBase64, err
		}
	}
	encryptedByte, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return decryptedBase64, err
	}
	decryptedRes, err := az.getClient().Decrypt(ctx, keyName, keyVersion, azkeys.KeyOperationsParameters{
		Algorithm: &algo,
		Value:     encryptedByte,
	}, nil)

	if err != nil {
		return "", err
	}

	return string(decryptedRes.Result), nil
}

func (az *azureKeyVault) getLatestKeyVersion(ctx context.Context, keyName string) (keyVersion string, err error) {
	key, err := az.getClient().GetKey(ctx, keyName, "", nil)
	if err != nil {
		return keyVersion, err
	}
	return key.Key.KID.Version(), nil
}
