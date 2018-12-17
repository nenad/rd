# rd

![status](https://travis-ci.com/nenadstojanovikj/realdebrid-go.svg?branch=master)

RealDebrid API client for the Go language

### Under development, minimal API at the moment

Supported API:

| Torrents  | 
| ------------- | 
| POST /torrents/addMagnet  | 
| GET /torrents/info/<ID>  | 
| POST /torrents/selectFiles/<ID>|

| Authentication |
| --- |
| GET /device/code |
| GET /device/credentials |
| POST /token |
