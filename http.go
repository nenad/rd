package rd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
)

type (
	PaginatedResponse struct {
		Page         int
		Offset       int
		CountPerPage int
		TotalCount   int
		Items        []interface{}
	}

	HTTPDoer interface {
		Do(r *http.Request) (*http.Response, error)
	}

	HTTPClient struct {
		client    HTTPDoer
		token     Token
		refresher TokenRefresher
	}
)

func (c *HTTPClient) Do(r *http.Request) (resp *http.Response, err error) {
	if c.refresher != nil && !c.token.IsValid() {
		if err := c.refreshToken(); err != nil {
			return nil, err
		}
	}

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

	resp, err = c.client.Do(r)
	if err != nil {
		return resp, err
	}

	return resp, parseErrorResponse(resp)
}

func AutoRefresh(c *HTTPClient) {
	c.refresher = NewAuthClient(c.client)
}

func (c *HTTPClient) refreshToken() error {
	token, err := c.refresher.RefreshAccessToken(c.token)
	if err != nil {
		return err
	}
	c.token = token
	return nil
}

func httpPostForm(doer HTTPDoer, url string, values map[string]string) (resp *http.Response, err error) {
	formBytes := &bytes.Buffer{}
	writer := multipart.NewWriter(formBytes)
	_ = writer.SetBoundary("realdebrid-boundary")
	for key, value := range values {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, formBytes)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err = doer.Do(req)
	return resp, parseErrorResponse(resp)
}

func httpGet(doer HTTPDoer, path string, params ...map[string]string) (resp *http.Response, err error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	for _, p := range params {
		for k, v := range p {
			query.Add(k, v)
		}
	}

	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err = doer.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, parseErrorResponse(resp)
}

func httpDelete(doer HTTPDoer, path string) (resp *http.Response, err error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err = doer.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, parseErrorResponse(resp)
}

func parseErrorResponse(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		return nil
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("Content-Type should be application/json, got %s and status code %d for %s",
			r.Header.Get("Content-Type"), r.StatusCode, r.Request.URL.String())
	}

	return extractError(r)
}

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

type httpError struct {
	ErrorMessage string `json:"error"`
	ErrorCode    int    `json:"error_code"`
	ErrorDetails string `json:"error_details"`

	StatusCode           int
	ErrorCodeDescription string
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

func (e httpError) Error() (msg string) {
	msg += fmt.Sprintf("error_message: %s\n", e.ErrorMessage)
	msg += fmt.Sprintf("error_code: %d - %s\n", e.ErrorCode, e.ErrorCodeDescription)
	msg += fmt.Sprintf("error_details: %s\n", e.ErrorDetails)
	msg += fmt.Sprintf("status_code: %d\n", e.StatusCode)
	return msg
}
