package hashicorp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type TokenCredential interface {
	GetToken() (AccessToken, error)
	SetClient(client *Client)
}

type UsernamePasswordCredential struct {
	username string
	password string
	client   *Client
}

type AccessToken struct {
	Token     string
	ExpiresOn time.Time
}

func NewUsernamePasswordCredential(username string, password string) *UsernamePasswordCredential {
	return &UsernamePasswordCredential{username, password, nil}
}

type loginPayload struct {
	Password string `json:"password"`
}

func (u *UsernamePasswordCredential) SetClient(client *Client) {
	u.client = client
}

func (u *UsernamePasswordCredential) GetToken() (AccessToken, error) {
	url := fmt.Sprintf("%s/auth/userpass/login/%s", u.client.VaultURL, u.username)
	params := loginPayload{
		Password: u.password,
	}

	jsonByte, _ := json.Marshal(params)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonByte))
	if err != nil {
		return AccessToken{}, err
	}

	defer resp.Body.Close()

	var bodyResponse map[string]interface{}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &bodyResponse)
	authResponse := bodyResponse["auth"].(map[string]interface{})
	leaseDuration := authResponse["lease_duration"].(float64)
	expiresOn := time.Now().Add(time.Duration(leaseDuration) * time.Second)
	return AccessToken{
		Token:     authResponse["client_token"].(string),
		ExpiresOn: expiresOn,
	}, err
}
