package rd

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRoundTripFunc func(req *http.Request) *http.Response

func (rt TestRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req), nil
}

func NewTestClient(fn TestRoundTripFunc) *HTTPClient {
	c := &http.Client{
		Transport: fn,
	}
	return &HTTPClient{token: Token{ExpiresIn: 3600, TokenType: "Bearer", AccessToken: "VALID_TOKEN", RefreshToken: "REFRESH_TOKEN"}, client: c}
}

func Test_AuthorizationHeaderIsPresent(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		auth := req.Header.Get("Authorization")
		assert.Equal(t, auth, "Bearer VALID_TOKEN", "Authorization header is not set correctly")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     map[string][]string{"Content-Type": {"application/json"}},
		}
	})

	req, _ := http.NewRequest("GET", "https://example.com", nil)
	_, err := client.Do(req)
	assert.NoError(t, err)
}

func Test_FormPost(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		contentType := req.Header.Get("Content-Type")
		assert.Equal(t, contentType, "multipart/form-data; boundary=realdebrid-boundary", "Content-Type is not set properly")

		body, _ := ioutil.ReadAll(req.Body)
		bodyString := string(body)
		expectedBody := "--realdebrid-boundary\r\nContent-Disposition: form-data; name=\"hello\"\r\n\r\nworld\r\n--realdebrid-boundary--\r\n"
		assert.Equal(t, bodyString, expectedBody)

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     map[string][]string{"Content-Type": {"application/json"}},
		}
	})

	_, err := httpPostForm(client, "https://example.com", map[string]string{"hello": "world"})
	assert.NoError(t, err)
}
