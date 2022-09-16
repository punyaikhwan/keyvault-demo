package hashicorp

import (
	"context"
	"encoding/base64"
	"fmt"
	"keyvault-demo/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ctxTest     = context.TODO()
	keyNameTest = "demo"
	algoTest    = JSONWebKeyEncryptionAlgorithmRSAOAEP256
	plainText   = "password"
)

func TestMain(m *testing.M) {
	config.ReadConfig("../../../.env")
	m.Run()
}

func TestLogin(t *testing.T) {
	cred := NewUsernamePasswordCredential(config.Configuration().VaultUsername, config.Configuration().VaultPassword)
	client := Client{
		VaultURL: config.Configuration().VaultURL,
	}
	cred.SetClient(&client)

	accessToken, err := cred.GetToken()
	assert.Nil(t, err)
	assert.NotEmpty(t, accessToken.Token)
}

func TestEncrypt(t *testing.T) {
	cred := NewUsernamePasswordCredential(config.Configuration().VaultUsername, config.Configuration().VaultPassword)
	client := NewClient(config.Configuration().VaultURL, cred)

	encResult, err := client.Encrypt(ctxTest, keyNameTest, "", KeyOperationsParameters{
		Algorithm: &algoTest,
		Value:     []byte(plainText),
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, encResult.Base64Result)
	assert.NotEmpty(t, encResult.Version)
}

func TestEncryptDecrypt(t *testing.T) {
	cred := NewUsernamePasswordCredential(config.Configuration().VaultUsername, config.Configuration().VaultPassword)
	client := NewClient(config.Configuration().VaultURL, cred)

	encResult, err := client.Encrypt(ctxTest, keyNameTest, "", KeyOperationsParameters{
		Algorithm: &algoTest,
		Value:     []byte(plainText),
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, encResult.Base64Result)
	assert.NotEmpty(t, encResult.Version)

	fmt.Printf("Encrypted: %s\n", encResult.Base64Result)

	decryptResult, err := client.Decrypt(ctxTest, keyNameTest, KeyOperationsParameters{
		Algorithm:   &algoTest,
		Base64Value: encResult.Base64Result,
	})

	assert.Nil(t, err)

	textResult, err := base64.StdEncoding.DecodeString(decryptResult.Base64Result)
	assert.Nil(t, err)
	assert.Equal(t, string(textResult), plainText)

	fmt.Printf("Decrypted: %s\n", string(textResult))
}
