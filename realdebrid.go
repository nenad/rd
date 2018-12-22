package rd

import (
	"net/http"
)

const apiBaseUrl = "https://api.real-debrid.com/rest/1.0"

type RealDebrid struct {
	Torrents   TorrentService
	Unrestrict UnrestrictService
	Downloads  DownloadService

	httpClient *HTTPClient
}

func NewRealDebrid(token Token, client *http.Client, options ...func(*HTTPClient)) *RealDebrid {
	if client == nil {
		client = http.DefaultClient
	}

	c := &HTTPClient{client: client, token: token}

	for _, option := range options {
		option(c)
	}

	return &RealDebrid{
		httpClient: c,
		Torrents:   &TorrentClient{c},
		Unrestrict: &UnrestrictClient{c},
		Downloads:  &DownloadClient{c},
	}
}

func (c *RealDebrid) IsTokenValid() bool {
	return c.httpClient.token.IsValid()
}
