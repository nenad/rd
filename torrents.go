package realdebrid

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	AddMagnetUrl      = ApiBaseUrl + "/torrents/addMagnet"
	GetTorrentInfoUrl = ApiBaseUrl + "/torrents/info/%s"
	SelectFilesUrl    = ApiBaseUrl + "/torrents/selectFiles/%s"
)

func (c *Client) AddMagnetLinkSimple(magnet string) (info TorrentUrlInfo, err error) {
	resp, err := c.postForm(AddMagnetUrl, url.Values{"magnet": {magnet}})
	if err != nil {
		return info, err
	}

	err = json.NewDecoder(resp.Body).Decode(&info)
	return info, err
}

func (c *Client) SelectFilesFromTorrent(id string, fileIds []int) error {
	_, err := c.postForm(fmt.Sprintf(SelectFilesUrl, id), url.Values{"files": {joinInts(fileIds)}})
	return err
}

func (c *Client) GetTorrent(id string) (info TorrentInfo, err error) {
	resp, err := c.get(fmt.Sprintf(GetTorrentInfoUrl, id))
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
