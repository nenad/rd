package rd_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/nenadstojanovikj/rd"

	"github.com/stretchr/testify/assert"
)

func NewUnrestrictTestClient(fn TestRoundTripFunc) rd.UnrestrictService {
	c := &http.Client{
		Transport: fn,
	}
	return rd.NewRealDebrid(
		rd.Token{ExpiresIn: 3600, TokenType: "Bearer", AccessToken: "VALID_TOKEN", RefreshToken: "REFRESH_TOKEN"},
		c).Unrestrict
}

func TestClient_UnrestrictLinkSimple(t *testing.T) {
	client := NewUnrestrictTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/rest/1.0/unrestrict/link", req.URL.String())
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "test-link-here", req.FormValue("link"))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewBufferString(
				`{
    "id": "AIO2UCGAIAQMD",
    "filename": "helloworld.mkv",
    "mimeType": "video/x-matroska",
    "filesize": 1361435465,
    "link": "https://real-debrid.com/d/HIMWA4NP4ZLGY",
    "host": "real-debrid.com",
    "host_icon": "https://fcdn.real-debrid.com/0754/images/hosters/realdebrid.png",
    "chunks": 16,
    "crc": 1,
    "download": "https://30.rdeb.io/d/AIO2UCGAIAQMD/helloworld.mkv",
    "streamable": 1
}`,
			)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	urlInfo, err := client.SimpleUnrestrict("test-link-here")
	assert.NoError(t, err)
	assert.Equal(t, rd.DownloadInfo{
		ID:         "AIO2UCGAIAQMD",
		Filename:   "helloworld.mkv",
		Filesize:   1361435465,
		Link:       "https://real-debrid.com/d/HIMWA4NP4ZLGY",
		MimeType:   "video/x-matroska",
		Host:       "real-debrid.com",
		HostIcon:   "https://fcdn.real-debrid.com/0754/images/hosters/realdebrid.png",
		Chunks:     16,
		CRC:        1,
		Download:   "https://30.rdeb.io/d/AIO2UCGAIAQMD/helloworld.mkv",
		Streamable: 1,
	}, urlInfo)
}
