package rd

import (
	"encoding/json"
)

// Endpoints
const (
	UnrestrictUrl = ApiBaseUrl + "/unrestrict/link"
)

type (
	DownloadInfo struct {
		ID         string `json:"id"`
		Filename   string `json:"filename"`
		MimeType   string `json:"mimeType"`
		Filesize   int    `json:"filesize"`
		Link       string `json:"link"`
		Host       string `json:"host"`
		HostIcon   string `json:"host_icon"`
		Chunks     int    `json:"chunks"`
		CRC        int    `json:"crc"`
		Download   string `json:"download"`
		Streamable int    `json:"streamable"`
	}

	UnrestrictService interface {
		SimpleUnrestrict(link string) (info DownloadInfo, err error)
	}

	UnrestrictClient struct {
		HTTPDoer
	}
)

func (c *UnrestrictClient) SimpleUnrestrict(link string) (info DownloadInfo, err error) {
	resp, err := PostForm(c, UnrestrictUrl, map[string]string{"link": link})
	if err != nil {
		return info, err
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}
