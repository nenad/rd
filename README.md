# rd

![status](https://travis-ci.com/nenadstojanovikj/rd.svg?branch=master)

RealDebrid API client for the Go language

### Under development, minimal API at the moment

Supported API:

| Torrents  | Description
| ------------- | -----|
| POST /torrents/addMagnet  | Prepares a magnet link for download
| GET /torrents/info/<ID>  | Gets info for torrent ID
| POST /torrents/selectFiles/<ID>| Selects files from a torrent
| GET /torrents | Gets list of torrents in your account
| DELETE /torrent/delete/<ID> | Deletes a torrent from your account

| Unrestrict  | Description
| ------------- | -----|
| POST /unrestrict/link | Unrestricts a link

| Downloads  | Description
| ------------- | -----|
| GET /downloads | Lists downloads on your account
| DELETE /downloads/delete/<ID> | Deletes a download from your account

| Authentication |
| --- |
| GET /device/code |
| GET /device/credentials |
| POST /token |
