package realdebrid

import (
	"fmt"
	"net/http"
)

var errorCodeDescription = map[int]string{
	-1: "Internal error",
	0:  "Unknown error",
	1:  "Missing parameter",
	2:  "Bad parameter value",
	3:  "Unknown method",
	4:  "Method not allowed",
	5:  "Slow down",
	6:  "Resource unreachable",
	7:  "Resource not found",
	8:  "Bad token",
	9:  "Permission denied",
	10: "Two-Factor authentication needed",
	11: "Two-Factor authentication pending",
	12: "Invalid login",
	13: "Invalid password",
	14: "Account locked",
	15: "Account not activated",
	16: "Unsupported hoster",
	17: "Hoster in maintenance",
	18: "Hoster limit reached",
	19: "Hoster temporarily unavailable",
	20: "Hoster not available for free users",
	21: "Too many active downloads",
	22: "IP Address not allowed",
	23: "Traffic exhausted",
	24: "File unavailable",
	25: "Service unavailable",
	26: "Upload too big",
	27: "Upload error",
	28: "File not allowed",
	29: "Torrent too big",
	30: "Torrent file invalid",
	31: "Action already done",
	32: "Image resolution error",
}

const ApiBaseUrl = "https://api.real-debrid.com/rest/1.0"

type Client struct {
	AuthService    AuthService
	TorrentService TorrentService
	httpClient     *RDClient
}

type RDClient struct {
	httpClient HTTPDoer
	token      Token
	refresher  TokenRefresher
}

func NewClient(token Token, client *http.Client, options ...func(*RDClient)) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	auth := &AuthClient{client}
	c := &RDClient{client, token, auth}

	for _, option := range options {
		option(c)
	}

	return &Client{
		httpClient:     c,
		AuthService:    auth,
		TorrentService: &TorrentClient{c},
	}
}

type AutoRefreshClient struct {
	httpClient HTTPDoer
	token      Token
	refresher  TokenRefresher
}

func (c *AutoRefreshClient) Do(r *http.Request) (resp *http.Response, err error) {
	if !c.token.IsValid() {
		t, err := c.refresher.RefreshAccessToken(c.token)
		if err != nil {
			return nil, err
		}
		c.token = t
	}

	return c.httpClient.Do(r)
}

func AutoRefresh(c *RDClient) {
	arClient := AutoRefreshClient{httpClient: c.httpClient, token: c.token, refresher: c.refresher}
	c.httpClient = &arClient
}

func (c *Client) IsTokenValid() bool {
	return c.httpClient.token.IsValid()
}

func (c *RDClient) Do(r *http.Request) (resp *http.Response, err error) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

	resp, err = c.httpClient.Do(r)
	if err != nil {
		return resp, err
	}

	return resp, parseErrorResponse(resp)
}
