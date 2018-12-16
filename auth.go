package realdebrid

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const (
	AuthBaseUrl = "https://api.real-debrid.com/oauth/v2"

	DeviceUrl      = AuthBaseUrl + "/device/code"
	CredentialsUrl = AuthBaseUrl + "/device/credentials"
	TokenUrl       = AuthBaseUrl + "/token"

	OpenSourceClientId = "X245A4XAIBGVM"
)

// StartAuthentication starts the authentication flow for the service
// RealDebrid API information: https://api.real-debrid.com/#device_auth_no_secret
func (c *Client) StartAuthentication(clientID string) (v Verification, err error) {
	authUrl, err := url.Parse(DeviceUrl)
	if err != nil {
		return v, err
	}

	query := authUrl.Query()
	query.Add("client_id", clientID)
	query.Add("new_credentials", "yes")
	authUrl.RawQuery = query.Encode()

	resp, err := c.get(authUrl.String())
	if err != nil {
		return v, err
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}

// ObtainSecret returns the Client ID and Client secret that are used for
// obtaining a valid token in the next step
func (c *Client) ObtainSecret(deviceCode, clientID string) (secrets Secrets, err error) {
	authUrl, err := url.Parse(CredentialsUrl)
	if err != nil {
		return secrets, err
	}

	query := authUrl.Query()
	query.Add("client_id", clientID)
	query.Add("code", deviceCode)
	authUrl.RawQuery = query.Encode()

	resp, err := c.get(authUrl.String())
	if err != nil {
		return secrets, err
	}

	err = json.NewDecoder(resp.Body).Decode(&secrets)
	if secrets.ClientID == "" || secrets.ClientSecret == "" {
		return secrets, fmt.Errorf("secrets not authorized")
	}
	return secrets, err
}

// ObtainAccessToken tries to get the client ID
func (c *Client) ObtainAccessToken(clientID, secret, code string) (t Token, err error) {
	resp, err := c.postForm(TokenUrl, url.Values{
		"client_id":     {clientID},
		"client_secret": {secret},
		"code":          {code},
		"grant_type":    {"http://oauth.net/grant_type/device/1.0"},
	})

	if err != nil {
		return t, err
	}

	t.obtainedAt = time.Now()
	err = json.NewDecoder(resp.Body).Decode(&t)
	return t, err
}

// Authenticate sets a given token for authentication for further requests
func (c *Client) Authenticate(t Token) {
	c.token = t
}

// Reauthenticate tries to get a new token from the service, and if successful
// it sets the new token
func (c *Client) Reauthenticate() error {
	if c.token.RefreshToken == "" {
		return fmt.Errorf("cannot reauthorize without refresh token")
	}

	secrets, err := c.ObtainSecret(c.token.RefreshToken, OpenSourceClientId)
	if err != nil {
		return err
	}

	t, err := c.ObtainAccessToken(secrets.ClientID, secrets.ClientSecret, c.token.RefreshToken)
	if err != nil {
		return err
	}

	c.Authenticate(t)
	return nil
}

// IsAuthorized checks if the current token is still valid
func (c *Client) IsAuthorized() bool {
	if c.token.AccessToken == "" || c.token.RefreshToken == "" {
		return false
	}

	// We expire the token 10 seconds before, so we don't send a request and risk to have it failing during transport
	tokenExpiry := c.token.obtainedAt.Add(time.Second * time.Duration(c.token.ExpiresIn-10))
	return time.Now().Before(tokenExpiry)
}
