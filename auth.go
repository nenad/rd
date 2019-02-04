package rd

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	authBaseUrl = "https://api.real-debrid.com/oauth/v2"

	deviceUrl      = authBaseUrl + "/device/code"
	credentialsUrl = authBaseUrl + "/device/credentials"
	tokenUrl       = authBaseUrl + "/token"

	defaultClientID = "X245A4XAIBGVM"
)

type (
	AuthClient struct {
		HTTPDoer
	}

	Verification struct {
		DeviceCode            string `json:"device_code"`
		UserCode              string `json:"user_code"`
		Interval              int    `json:"interval"`
		ExpiresIn             int64  `json:"expires_in"`
		VerificationURL       string `json:"verification_url"`
		DirectVerificationURL string `json:"direct_verification_url"`
	}
	Secrets struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}
)

func NewAuthClient(doer HTTPDoer) *AuthClient {
	return &AuthClient{doer}
}

// StartAuthentication starts the authentication flow for the service
// RealDebrid API information: https://api.real-debrid.com/#device_auth_no_secret
func (c *AuthClient) StartAuthentication(clientID string) (v Verification, err error) {
	resp, err := httpGet(c, deviceUrl, map[string]string{"client_id": clientID, "new_credentials": "yes"})
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}

// ObtainSecret returns the HTTPClient ID and HTTPClient secret that are used for
// obtaining a valid token in the next step
func (c *AuthClient) ObtainSecret(deviceCode, clientID string) (secrets Secrets, err error) {
	resp, err := httpGet(c, credentialsUrl, map[string]string{"client_id": clientID, "code": deviceCode})
	if err != nil {
		return secrets, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&secrets)
	if secrets.ClientID == "" || secrets.ClientSecret == "" {
		return secrets, fmt.Errorf("secrets not authorized")
	}
	return secrets, err
}

// ObtainAccessToken tries to get a new token from the service
func (c *AuthClient) ObtainAccessToken(clientID, secret, code string) (t Token, err error) {
	resp, err := httpPostForm(c, tokenUrl, map[string]string{
		"client_id":     clientID,
		"client_secret": secret,
		"code":          code,
		"grant_type":    "http://oauth.net/grant_type/device/1.0",
	})
	if err != nil {
		return t, err
	}

	t.ObtainedAt = time.Now()
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&t)
	return t, err
}

// RefreshAccessToken tries to refresh the given token and get a new one
func (c *AuthClient) RefreshAccessToken(token Token) (t Token, err error) {
	if token.RefreshToken == "" {
		return t, fmt.Errorf("cannot reauthorize without refresh token")
	}

	secrets, err := c.ObtainSecret(token.RefreshToken, defaultClientID)
	if err != nil {
		return t, err
	}

	return c.ObtainAccessToken(secrets.ClientID, secrets.ClientSecret, token.RefreshToken)
}
