package rd_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/nenadstojanovikj/rd"

	"github.com/stretchr/testify/assert"
)

type TestRoundTripFunc func(req *http.Request) *http.Response

func (rt TestRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req), nil
}

func NewAuthTestClient(fn TestRoundTripFunc) *rd.AuthClient {
	c := &http.Client{
		Transport: fn,
	}
	return &rd.AuthClient{c}
}

func TestAuthClient_CanStartAuthenticationFlowSuccessfully(t *testing.T) {
	client := NewAuthTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/oauth/v2/device/code?client_id=X245A4XAIBGVM&new_credentials=yes", req.URL.String())
		assert.Equal(t, "GET", req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
    "device_code": "MRHRIKNQJZMCIXQ7ZAWMMURIIFBTJTMJYACGWBQUDMZVED3ODAGQ", "user_code": "EBKEG2RR", "interval": 5, "expires_in": 600, "verification_url": "https://real-debrid.com/device",
    "direct_verification_url": "https://real-debrid.com/authorize?client_id=X245A4XAIBGVM&device_id=MRHRIKNQJZMCIXQ7ZAWMMURIIFBTJTMJYACGWBQUDMZVED3ODAGQ"
}`,
			)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	verification, err := client.StartAuthentication(rd.OpenSourceClientId)
	assert.NoError(t, err)
	assert.Equal(t, rd.Verification{
		ExpiresIn:             600,
		Interval:              5,
		VerificationURL:       "https://real-debrid.com/device",
		DeviceCode:            "MRHRIKNQJZMCIXQ7ZAWMMURIIFBTJTMJYACGWBQUDMZVED3ODAGQ",
		DirectVerificationURL: "https://real-debrid.com/authorize?client_id=X245A4XAIBGVM&device_id=MRHRIKNQJZMCIXQ7ZAWMMURIIFBTJTMJYACGWBQUDMZVED3ODAGQ",
		UserCode:              "EBKEG2RR",
	}, verification)
}

func TestAuthClient_CanObtainSecretsSuccessfully(t *testing.T) {
	client := NewAuthTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/oauth/v2/device/credentials?client_id=X245A4XAIBGVM&code=YD7HNOMEJOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOQ", req.URL.String())
		assert.Equal(t, "GET", req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"client_id":"3N4RHGK5OKNIH","client_secret":"135d1b6dc60dddbdaa2e5d41772c85d56c54790b"}`)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	secrets, err := client.ObtainSecret("YD7HNOMEJOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOQ", rd.OpenSourceClientId)
	assert.NoError(t, err)
	assert.Equal(t, rd.Secrets{
		ClientSecret: "135d1b6dc60dddbdaa2e5d41772c85d56c54790b",
		ClientID:     "3N4RHGK5OKNIH",
	}, secrets)
}

func TestAuthClient_ErrorsOnWrongOrEmptySecrets(t *testing.T) {
	client := NewAuthTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/oauth/v2/device/credentials?client_id=X245A4XAIBGVM&code=YD7HNOMEJOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOQ", req.URL.String())
		assert.Equal(t, "GET", req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"client_id":null,"client_secret":null}`)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	_, err := client.ObtainSecret("YD7HNOMEJOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOQ", rd.OpenSourceClientId)
	assert.EqualError(t, err, "secrets not authorized")
}

func TestAuthClient_CanObtainValidToken(t *testing.T) {
	client := NewAuthTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/oauth/v2/token", req.URL.String())
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "0N2RHHK5OKNIX", req.FormValue("client_id"))
		assert.Equal(t, "135d1b6dc60dddbcaa2e5dc1772c85d56c5479ba", req.FormValue("client_secret"))
		assert.Equal(t, "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP", req.FormValue("code"))
		assert.Equal(t, "http://oauth.net/grant_type/device/1.0", req.FormValue("grant_type"))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token": "QMUJ32Q4S3X57D3NC354V4OI62JZ74KV5H3DZX7JJEOZSLXCWVYA", "expires_in": 3600, "refresh_token": "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP", "token_type": "Bearer" }`)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	token, err := client.ObtainAccessToken("0N2RHHK5OKNIX", "135d1b6dc60dddbcaa2e5dc1772c85d56c5479ba", "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP")
	assert.NoError(t, err)
	expectedToken := rd.Token{
		ExpiresIn:    3600,
		AccessToken:  "QMUJ32Q4S3X57D3NC354V4OI62JZ74KV5H3DZX7JJEOZSLXCWVYA",
		TokenType:    "Bearer",
		RefreshToken: "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP",
	}

	assert.Equal(t, expectedToken.ExpiresIn, token.ExpiresIn)
	assert.Equal(t, expectedToken.AccessToken, token.AccessToken)
	assert.Equal(t, expectedToken.TokenType, token.TokenType)
	assert.Equal(t, expectedToken.RefreshToken, token.RefreshToken)
}

func TestAuthClient_ErrorsOnFailedObtainingToken(t *testing.T) {
	client := NewAuthTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/oauth/v2/token", req.URL.String())
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "0N2RHHK5OKNIX", req.FormValue("client_id"))
		assert.Equal(t, "135d1b6dc60dddbcaa2e5dc1772c85d56c5479ba", req.FormValue("client_secret"))
		assert.Equal(t, "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP", req.FormValue("code"))
		assert.Equal(t, "http://oauth.net/grant_type/device/1.0", req.FormValue("grant_type"))

		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{ "error": "wrong_parameter", "error_code": 2 }`)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	token, err := client.ObtainAccessToken("0N2RHHK5OKNIX", "135d1b6dc60dddbcaa2e5dc1772c85d56c5479ba", "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP")
	assert.Empty(t, token)
	assert.EqualError(t, err, "error_message: wrong_parameter\nerror_code: 2 - Unknown error\nerror_details: \nstatus_code: 400\n")
}

func TestAuthClient_CanObtainAccessTokenFromRefreshToken(t *testing.T) {
	client := NewAuthTestClient(func(req *http.Request) *http.Response {
		if req.URL.String() == "https://api.real-debrid.com/oauth/v2/device/credentials?client_id=X245A4XAIBGVM&code=REFRESH_TOKEN" {
			assert.Equal(t, "GET", req.Method)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"client_id":"0N2RHHK5OKNIX","client_secret":"135d1b6dc60dddbcaa2e5dc1772c85d56c5479ba"}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			}
		}
		if req.URL.String() == "https://api.real-debrid.com/oauth/v2/token" {
			assert.Equal(t, "https://api.real-debrid.com/oauth/v2/token", req.URL.String())
			assert.Equal(t, "POST", req.Method)
			assert.Equal(t, "0N2RHHK5OKNIX", req.FormValue("client_id"))
			assert.Equal(t, "135d1b6dc60dddbcaa2e5dc1772c85d56c5479ba", req.FormValue("client_secret"))
			assert.Equal(t, "REFRESH_TOKEN", req.FormValue("code"))
			assert.Equal(t, "http://oauth.net/grant_type/device/1.0", req.FormValue("grant_type"))

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token": "ZMUJ32Q4S3X57D3NC354V4OI62JZ74KV5H3DZX7JJEOZSLXCWVYA", "expires_in": 3600, "refresh_token": "QD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP", "token_type": "Bearer" }`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			}
		}
		t.Error("the expected routes were not called")
		return nil
	})

	token, err := client.RefreshAccessToken(rd.Token{RefreshToken: "REFRESH_TOKEN"})
	assert.NoError(t, err)
	expectedToken := rd.Token{
		ExpiresIn:    3600,
		AccessToken:  "ZMUJ32Q4S3X57D3NC354V4OI62JZ74KV5H3DZX7JJEOZSLXCWVYA",
		TokenType:    "Bearer",
		RefreshToken: "QD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP",
	}

	assert.Equal(t, expectedToken.ExpiresIn, token.ExpiresIn)
	assert.Equal(t, expectedToken.AccessToken, token.AccessToken)
	assert.Equal(t, expectedToken.TokenType, token.TokenType)
	assert.Equal(t, expectedToken.RefreshToken, token.RefreshToken)
}
