package rd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Endpoints
const (
	AddMagnetUrl      = ApiBaseUrl + "/torrents/addMagnet"
	GetTorrentInfoUrl = ApiBaseUrl + "/torrents/info/%s"
	SelectFilesUrl    = ApiBaseUrl + "/torrents/selectFiles/%s"
)

// Possible torrent states
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
	TorrentService interface {
		AddMagnetLinkSimple(magnet string) (info TorrentUrlInfo, err error)
		SelectFilesFromTorrent(id string, fileIds []int) error
		GetTorrent(id string) (info TorrentInfo, err error)
	}

	TorrentClient struct {
		HTTPDoer
	}

	Status string

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
)

func (c *TorrentClient) AddMagnetLinkSimple(magnet string) (info TorrentUrlInfo, err error) {
	resp, err := PostForm(c, AddMagnetUrl, map[string]string{"magnet": magnet})
	if err != nil {
		return info, err
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}

func (c *TorrentClient) SelectFilesFromTorrent(id string, fileIds []int) error {
	_, err := PostForm(c, fmt.Sprintf(SelectFilesUrl, id), map[string]string{"files": joinInts(fileIds)})
	return err
}

func (c *TorrentClient) GetTorrent(id string) (info TorrentInfo, err error) {
	resp, err := Get(c, fmt.Sprintf(GetTorrentInfoUrl, id))
	if err != nil {
		return info, err
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}

func joinInts(slice []int) string {
	b := make([]string, len(slice))
	for i, v := range slice {
		b[i] = strconv.Itoa(v)
	}

	return strings.Join(b, ",")
}
