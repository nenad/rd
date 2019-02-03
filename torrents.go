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
	magnetAddUrl          = apiBaseUrl + "/torrents/addMagnet"
	torrentInfoUrl        = apiBaseUrl + "/torrents/info/%s"
	torrentsUrl           = apiBaseUrl + "/torrents"
	torrentDeleteUrl      = apiBaseUrl + "/torrents/delete/%s"
	torrentSelectFilesUrl = apiBaseUrl + "/torrents/selectFiles/%s"
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
		GetTorrents() (infos []TorrentInfo, err error)
		Delete(id string) error
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
		Bytes            int64     `json:"bytes"`
		OriginalBytes    int64     `json:"original_bytes"`
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
	resp, err := httpPostForm(c, magnetAddUrl, map[string]string{"magnet": magnet})
	if err != nil {
		return info, err
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}

func (c *TorrentClient) SelectFilesFromTorrent(id string, fileIds []int) error {
	_, err := httpPostForm(c, fmt.Sprintf(torrentSelectFilesUrl, id), map[string]string{"files": joinInts(fileIds)})
	return err
}

func (c *TorrentClient) GetTorrent(id string) (info TorrentInfo, err error) {
	resp, err := httpGet(c, fmt.Sprintf(torrentInfoUrl, id))
	if err != nil {
		return info, err
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}

func (c *TorrentClient) Delete(id string) error {
	_, err := httpDelete(c, fmt.Sprintf(torrentDeleteUrl, id))
	return err
}

func (c *TorrentClient) GetTorrents() (infos []TorrentInfo, err error) {
	resp, err := httpGet(c, torrentsUrl)
	if err != nil {
		return infos, err
	}

	err = json.NewDecoder(resp.Body).Decode(&infos)
	return infos, err
}

func joinInts(slice []int) string {
	b := make([]string, len(slice))
	for i, v := range slice {
		b[i] = strconv.Itoa(v)
	}

	return strings.Join(b, ",")
}
