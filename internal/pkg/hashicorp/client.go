package hashicorp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	VaultURL    string
	Token       AccessToken
	TransitPath string
}

func NewClient(vaultURL string, transitPath string, cred TokenCredential) *Client {
	client := &Client{
		VaultURL:    vaultURL,
		TransitPath: transitPath,
	}

	cred.SetClient(client)

	client.Token, _ = cred.GetToken()

	return client
}

type Key struct {
	Type     string      `json:"type"`
	Versions KeyVersions `json:"keys"`
	Name     string      `json:"name"`
}

// KeyVersions is map key version to creation time
type KeyVersions map[string]int64

func (k KeyVersions) Last() (version string, createdAt time.Time) {
	unixCreatedAt := int64(0)
	for keyVersion, unixCreationTime := range k {
		if unixCreationTime > unixCreatedAt {
			version = keyVersion
			unixCreatedAt = unixCreationTime
		}
	}

	createdAt = time.Unix(unixCreatedAt, 0)

	return version, createdAt
}

func (c *Client) GetKey(ctx context.Context, name string) (key Key, err error) {
	url := fmt.Sprintf("%s/%s/keys/%s", c.VaultURL, c.TransitPath, name)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return key, err
	}

	client := &http.Client{}
	req.Header.Set("X-Vault-Token", c.Token.Token)

	resp, err := client.Do(req)
	if err != nil {
		return key, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &key)

	return key, err
}

// HashicorpKeyType - algorithm identifier
type HashicorpKeyType string

const (
	HashicorpAES128GCM96     HashicorpKeyType = "aes128-gcm96"
	HashicorpAES256GCM96     HashicorpKeyType = "aes256-gcm96"
	HashicorpCACHA20POLY1305 HashicorpKeyType = "chacha20-poly1305"
	HashicorpED25519         HashicorpKeyType = "ed25519"
	HashicorpECDSAP256       HashicorpKeyType = "ecdsa-p256"
	HashicorpECDSAP384       HashicorpKeyType = "ecdsa-p384"
	HashicorpECDSAP521       HashicorpKeyType = "ecdsa-p521"
	HashicorpRSA2048         HashicorpKeyType = "rsa-2048"
	HashicorpRSA3072         HashicorpKeyType = "rsa-3072"
	HashicorpRSA4096         HashicorpKeyType = "rsa-4096"
	HashicorpHMAC            HashicorpKeyType = "hmac"
)

type KeyOperationsParameters struct {
	Algorithm   *HashicorpKeyType `json:"type,omitempty"`
	Nonce       string            `json:"nonce"`
	Value       []byte            `json:"-"`
	Base64Value string            `json:"plaintext"`
}

// EncryptResponse contains the response from method Client.Encrypt.
type EncryptResponse struct {
	KeyOperationResult
}

// DecryptResponse contains the response from method Client.Decrypt.
type DecryptResponse struct {
	KeyOperationResult
}

// KeyOperationResult - The key operation result.
type KeyOperationResult struct {
	// READ-ONLY; Key identifier
	Version string `json:"version"`

	// READ-ONLY
	Base64Result string `json:"value,omitempty"`
}

type encryptResult struct {
	Data struct {
		CipherText string  `json:"ciphertext"`
		KeyVersion float64 `json:"key_version"`
	} `json:"data"`
}

type decryptResult struct {
	Data struct {
		PlainText string `json:"plaintext"`
	} `json:"data"`
}

func (c *Client) Encrypt(ctx context.Context, name string, version string, params KeyOperationsParameters) (EncryptResponse, error) {
	url := fmt.Sprintf("%s/%s/encrypt/%s", c.VaultURL, c.TransitPath, name)
	value := params.Base64Value
	if value == "" {
		value = base64.StdEncoding.EncodeToString(params.Value)
	}
	bodyParam := map[string]interface{}{
		"plaintext": value,
		"nonce":     params.Nonce,
		"type":      params.Algorithm,
	}

	jsonByte, _ := json.Marshal(bodyParam)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return EncryptResponse{}, err
	}

	client := &http.Client{}
	req.Header.Set("X-Vault-Token", c.Token.Token)

	resp, err := client.Do(req)
	if err != nil {
		return EncryptResponse{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string][]string
		json.Unmarshal(body, &errorBody)
		return EncryptResponse{}, errors.New(errorBody["errors"][0])
	}

	var encRes encryptResult
	err = json.Unmarshal(body, &encRes)

	var result EncryptResponse
	result.Base64Result = encRes.Data.CipherText
	result.Version = strconv.Itoa(int(encRes.Data.KeyVersion))
	return result, err
}

func (c *Client) Decrypt(ctx context.Context, name string, params KeyOperationsParameters) (DecryptResponse, error) {
	url := fmt.Sprintf("%s/%s/decrypt/%s", c.VaultURL, c.TransitPath, name)
	value := params.Base64Value
	if value == "" {
		value = base64.StdEncoding.EncodeToString(params.Value)
	}
	bodyParam := map[string]interface{}{
		"ciphertext": value,
		"nonce":      params.Nonce,
		"type":       params.Algorithm,
	}

	jsonByte, _ := json.Marshal(bodyParam)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return DecryptResponse{}, err
	}

	client := &http.Client{}
	req.Header.Set("X-Vault-Token", c.Token.Token)

	resp, err := client.Do(req)
	if err != nil {
		return DecryptResponse{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string][]string
		json.Unmarshal(body, &errorBody)
		return DecryptResponse{}, errors.New(errorBody["errors"][0])
	}

	var decRes decryptResult

	json.Unmarshal(body, &decRes)

	var result DecryptResponse
	result.Base64Result = decRes.Data.PlainText
	return result, err
}
