package realdebrid

import (
	"fmt"
	"strings"
	"time"
)

const (
	StatusMagnetError      Status = "magnet_error"
	StatusWaitingFiles     Status = "waiting_files_selection"
	StatusMagnetConversion Status = "magnet_conversion"
	StatusQueued           Status = "queued"
	StatusDownloading      Status = "downloading"
	StatusDownloaded       Status = "downloaded"
	StatusError            Status = "error"
	StatusVirus            Status = "virus"
	StatusCompressing      Status = "compressing"
	StatusUploading        Status = "uploading"
	StatusDead             Status = "dead"
)

type (
	Status string

	httpError struct {
		ErrorMessage string `json:"error"`
		ErrorCode    int    `json:"error_code"`
		ErrorDetails string `json:"error_details"`

		StatusCode           int
		ErrorCodeDescription string
	}

	TorrentUrlInfo struct {
		ID  string `json:"id"`
		URI string `json:"uri"`
	}

	File struct {
		ID       int    `json:"id"`
		Path     string `json:"path"`
		Bytes    int    `json:"bytes"`
		Selected int    `json:"selected"`
	}

	TorrentInfo struct {
		ID               string    `json:"id"`
		Filename         string    `json:"filename"`
		OriginalFilename string    `json:"original_filename"`
		Hash             string    `json:"hash"`
		Bytes            int       `json:"bytes"`
		OriginalBytes    int       `json:"original_bytes"`
		Host             string    `json:"host"`
		Split            int       `json:"split"`
		Progress         int       `json:"progress"`
		Status           Status    `json:"status"`
		Added            time.Time `json:"added"`
		Files            []File    `json:"files"`
		Links            []string  `json:"links"`
		Ended            time.Time `json:"ended"`
		Speed            int       `json:"speed"`
		Seeders          int       `json:"seeders"`
	}

	Verification struct {
		DeviceCode            string `json:"device_code"`
		UserCode              string `json:"user_code"`
		Interval              int    `json:"interval"`
		ExpiresIn             int    `json:"expires_in"`
		VerificationURL       string `json:"verification_url"`
		DirectVerificationURL string `json:"direct_verification_url"`
	}

	Secrets struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	Token struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		obtainedAt   time.Time
	}
)

func (e httpError) Error() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("error_message: %s\n", e.ErrorMessage))
	builder.WriteString(fmt.Sprintf("error_code: %d - %s\n", e.ErrorCode, e.ErrorCodeDescription))
	builder.WriteString(fmt.Sprintf("error_details: %s\n", e.ErrorDetails))
	builder.WriteString(fmt.Sprintf("status_code: %d\n", e.StatusCode))
	return builder.String()
}
