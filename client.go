package realdebrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
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

type (
	HTTPGetter interface {
		Get(url string) (*http.Response, error)
	}
	HTTPFormPoster interface {
		PostForm(url string, data url.Values) (resp *http.Response, err error)
	}
	HTTPDoer interface {
		Do(r *http.Request) (*http.Response, error)
	}
)

type AuthClient struct {
	internal *http.Client
}

type Client struct {
	token    Token
	internal *http.Client
}

func NewClient(token Token, client *http.Client) *Client {
	return &Client{token, client}
}

func NewNonAuthorizedClient(client *http.Client) *Client {
	return &Client{Token{}, client}
}

func (c *Client) do(r *http.Request) (resp *http.Response, err error) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

	resp, err = c.internal.Do(r)
	if err != nil {
		return resp, err
	}

	return resp, getError(resp)
}

func (c *Client) post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)

	return c.do(req)
}

func (c *Client) get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) postForm(url string, values url.Values) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.SetBoundary("realdebrid-boundary")
	for key, value := range values {
		if err := writer.WriteField(key, strings.Join(value, "\n")); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return c.do(req)
}

func getError(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		return nil
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("Content-Type should be application/json, got %s and status code %d for %s",
			r.Header.Get("Content-Type"), r.StatusCode, r.Request.URL.String())
	}

	return extractError(r)
}

func extractError(r *http.Response) httpError {
	defer r.Body.Close()
	e := &httpError{}
	e.StatusCode = r.StatusCode
	e.ErrorCodeDescription = errorCodeDescription[e.ErrorCode]
	err := json.NewDecoder(r.Body).Decode(e)
	if err != nil {
		e.ErrorDetails = "decoding error"
		e.ErrorMessage = err.Error()
	}

	return *e
}
