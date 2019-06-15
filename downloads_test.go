package rd_test

import (
	"bytes"
	"github.com/nenad/rd"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func NewDownloadsTestClient(fn TestRoundTripFunc) rd.DownloadService {
	c := &http.Client{
		Transport: fn,
	}
	return rd.NewRealDebrid(
		rd.Token{ExpiresIn: 3600, TokenType: "Bearer", AccessToken: "VALID_TOKEN", RefreshToken: "REFRESH_TOKEN"},
		c).Downloads
}

func TestClient_GetDownloads(t *testing.T) {
	client := NewDownloadsTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/rest/1.0/downloads", req.URL.String())
		assert.Equal(t, "GET", req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewBufferString(`
[
    {
        "id": "EFX5LLAJYDR6B",
        "filename": "helloworld.mkv",
        "mimeType": "video/x-matroska",
        "filesize": 123123123,
        "link": "https://real-debrid.com/d/HIMWA4NP4ZLGP",
        "host": "real-debrid.com",
        "host_icon": "https://fcdn.real-debrid.com/0754/images/hosters/realdebrid.png",
        "chunks": 16,
        "download": "http://37.rdeb.io/d/EFX5LLAJYDR6B/helloworld.mkv",
        "streamable": 1,
        "generated": "2018-12-20T00:47:53.000Z"
    }
]
`,
			)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	info, err := client.List()
	assert.NoError(t, err)
	assert.Equal(t, []rd.DownloadInfo{{
		ID:         "EFX5LLAJYDR6B",
		Filename:   "helloworld.mkv",
		Filesize:   123123123,
		Link:       "https://real-debrid.com/d/HIMWA4NP4ZLGP",
		MimeType:   "video/x-matroska",
		Host:       "real-debrid.com",
		HostIcon:   "https://fcdn.real-debrid.com/0754/images/hosters/realdebrid.png",
		Chunks:     16,
		Download:   "http://37.rdeb.io/d/EFX5LLAJYDR6B/helloworld.mkv",
		Streamable: 1,
		Generated:  time.Date(2018, 12, 20, 0, 47, 53, 0, time.UTC),
	}}, info)
}

func TestClient_DeleteDownload(t *testing.T) {
	client := NewDownloadsTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/rest/1.0/downloads/delete/XCBYL4ZIYPU42", req.URL.String())
		assert.Equal(t, "DELETE", req.Method)

		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	err := client.Delete("XCBYL4ZIYPU42")
	assert.NoError(t, err)
}
