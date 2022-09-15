package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

func main() {
	ctx := context.TODO()
	vaultURL := "https://demo-kms.vault.azure.net/"
	keyName := "demo"
	text := "password"
	// Create a credential using the NewDefaultAzureCredential type.
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	client := azkeys.NewClient(vaultURL, cred, nil)

	algo := azkeys.JSONWebKeyEncryptionAlgorithmRSAOAEP256
	encryptRes, err := client.Encrypt(ctx, keyName, "f75e27ba59ee4ce5baf74a8ef05ae2e2", azkeys.KeyOperationsParameters{
		Algorithm: &algo,
		Value:     []byte(text),
	}, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(base64.StdEncoding.EncodeToString([]byte(encryptRes.Result)))
}
