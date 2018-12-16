package realdebrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
)

type (
	HTTPDoer interface {
		Do(r *http.Request) (*http.Response, error)
	}
)

func PostForm(doer HTTPDoer, url string, values map[string]string) (resp *http.Response, err error) {
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

func Get(doer HTTPDoer, path string, params ...map[string]string) (resp *http.Response, err error) {
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
