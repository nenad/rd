package rd

import (
	"encoding/json"
	"fmt"
	"time"
)

// Endpoints
const (
	downloadsUrl      = apiBaseUrl + "/downloads"
	downloadDeleteUrl = apiBaseUrl + "/downloads/delete/%s"
)

type (
	DownloadInfo struct {
		ID         string    `json:"id"`
		Filename   string    `json:"filename"`
		MimeType   string    `json:"mimeType"`
		Filesize   int64     `json:"filesize"`
		Link       string    `json:"link"`
		Host       string    `json:"host"`
		HostIcon   string    `json:"host_icon"`
		Chunks     int       `json:"chunks"`
		Download   string    `json:"download"`
		Streamable int       `json:"streamable"`
		Generated  time.Time `json:"generated"`
	}

	DownloadService interface {
		List() ([]DownloadInfo, error)
		Delete(id string) error
	}

	DownloadClient struct {
		HTTPDoer
	}
)

func (s *DownloadClient) List() (items []DownloadInfo, err error) {
	resp, err := httpGet(s.HTTPDoer, downloadsUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&items)

	return items, err
}

func (s *DownloadClient) Delete(id string) error {
	_, err := httpDelete(s.HTTPDoer, fmt.Sprintf(downloadDeleteUrl, id))
	return err
}
