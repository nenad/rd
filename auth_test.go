package realdebrid

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_NewCredentialsCanBeObtained(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
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

	verification, err := client.StartAuthentication(OpenSourceClientId)
	assert.NoError(t, err)
	assert.Equal(t, Verification{
		ExpiresIn:             600,
		Interval:              5,
		VerificationURL:       "https://real-debrid.com/device",
		DeviceCode:            "MRHRIKNQJZMCIXQ7ZAWMMURIIFBTJTMJYACGWBQUDMZVED3ODAGQ",
		DirectVerificationURL: "https://real-debrid.com/authorize?client_id=X245A4XAIBGVM&device_id=MRHRIKNQJZMCIXQ7ZAWMMURIIFBTJTMJYACGWBQUDMZVED3ODAGQ",
		UserCode:              "EBKEG2RR",
	}, verification)
}

func TestClient_SecretsCanBeObtainedSuccessfully(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
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

	secrets, err := client.ObtainSecret("YD7HNOMEJOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOQ", OpenSourceClientId)
	assert.NoError(t, err)
	assert.Equal(t, Secrets{
		ClientSecret: "135d1b6dc60dddbdaa2e5d41772c85d56c54790b",
		ClientID:     "3N4RHGK5OKNIH",
	}, secrets)
}

func TestClient_FailToGetAuthorizedSecrets(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
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

	_, err := client.ObtainSecret("YD7HNOMEJOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOQ", OpenSourceClientId)
	assert.EqualError(t, err, "secrets not authorized")
}

func TestClient_AuthorizeSuccess(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
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

	err := client.Authorize("0N2RHHK5OKNIX", "135d1b6dc60dddbcaa2e5dc1772c85d56c5479ba", "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP")
	assert.NoError(t, err)
	expectedToken := Token{
		ExpiresIn:    3600,
		AccessToken:  "QMUJ32Q4S3X57D3NC354V4OI62JZ74KV5H3DZX7JJEOZSLXCWVYA",
		TokenType:    "Bearer",
		RefreshToken: "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP",
	}

	assert.Equal(t, expectedToken.ExpiresIn, client.token.ExpiresIn)
	assert.Equal(t, expectedToken.AccessToken, client.token.AccessToken)
	assert.Equal(t, expectedToken.TokenType, client.token.TokenType)
	assert.Equal(t, expectedToken.RefreshToken, client.token.RefreshToken)
	assert.True(t, client.IsAuthorized())
}

func TestClient_AuthorizeFailure(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
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

	err := client.Authorize("0N2RHHK5OKNIX", "135d1b6dc60dddbcaa2e5dc1772c85d56c5479ba", "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP")
	assert.EqualError(t, err, "error_message: wrong_parameter\nerror_code: 2 - Unknown error\nerror_details: \nstatus_code: 400\n")
	assert.False(t, client.IsAuthorized())
}

func TestClient_Reauthorize(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
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
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"access_token": "QMUJ32Q4S3X57D3NC354V4OI62JZ74KV5H3DZX7JJEOZSLXCWVYA", "expires_in": 3600, "refresh_token": "ZD7HNOMEXOJY7P2FP4XIJA5E634RWZKWWQ6RZNJJT235G4RNCAOP", "token_type": "Bearer" }`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			}
		}
		t.Error("the expected routes were not called")
		return nil
	})

	err := client.Reauthorize()
	assert.NoError(t, err)
	assert.True(t, client.IsAuthorized())
}

func TestClient_ExpiredAuthorization(t *testing.T) {
	client := NewClient(
		Token{
			AccessToken:  "ACCESS",
			RefreshToken: "REFRESH",
			obtainedAt:   time.Now().Add(-3600 * time.Second),
		},
		http.DefaultClient)

	assert.False(t, client.IsAuthorized())
}
