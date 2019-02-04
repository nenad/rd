package rd

import (
	"encoding/json"
)

// Endpoints
const (
	unrestrictUrl = apiBaseUrl + "/unrestrict/link"
)

type (
	UnrestrictInfo struct {
		ID         string `json:"id"`
		Filename   string `json:"filename"`
		MimeType   string `json:"mimeType"`
		Filesize   int64  `json:"filesize"`
		Link       string `json:"link"`
		Host       string `json:"host"`
		HostIcon   string `json:"host_icon"`
		Chunks     int    `json:"chunks"`
		CRC        int    `json:"crc"`
		Download   string `json:"download"`
		Streamable int    `json:"streamable"`
	}

	UnrestrictService interface {
		SimpleUnrestrict(link string) (info UnrestrictInfo, err error)
	}

	UnrestrictClient struct {
		HTTPDoer
	}
)

func (c *UnrestrictClient) SimpleUnrestrict(link string) (info UnrestrictInfo, err error) {
	resp, err := httpPostForm(c, unrestrictUrl, map[string]string{"link": link})
	if err != nil {
		return info, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}
