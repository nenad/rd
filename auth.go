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

func (c *Client) Authorize(clientID, secret, code string) (err error) {
	resp, err := c.postForm(TokenUrl, url.Values{
		"client_id":     {clientID},
		"client_secret": {secret},
		"code":          {code},
		"grant_type":    {"http://oauth.net/grant_type/device/1.0"},
	})

	if err != nil {
		return err
	}

	t := Token{}
	t.obtainedAt = time.Now()
	err = json.NewDecoder(resp.Body).Decode(&t)
	c.token = t
	return err
}

func (c *Client) Reauthorize() error {
	if c.token.RefreshToken == "" {
		return fmt.Errorf("cannot reauthorize without refresh token")
	}

	secrets, err := c.ObtainSecret(c.token.RefreshToken, OpenSourceClientId)
	if err != nil {
		return err
	}

	return c.Authorize(secrets.ClientID, secrets.ClientSecret, c.token.RefreshToken)
}

func (c *Client) IsAuthorized() bool {
	if c.token.AccessToken == "" || c.token.RefreshToken == "" {
		return false
	}

	// We expire the token 10 seconds before, so we don't send a request and risk to have it failed mid-transport
	tokenExpiry := c.token.obtainedAt.Add(time.Second * time.Duration(c.token.ExpiresIn-10))
	return time.Now().Before(tokenExpiry)
}
